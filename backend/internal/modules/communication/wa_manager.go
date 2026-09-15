package communication

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

var (
	ErrAlreadyConnected = errors.New("device already connected")
	ErrDeviceNotFound   = errors.New("device not found")
)

type cachedQRInfo struct {
	Code      string
	ExpiresAt time.Time
	Timeout   time.Duration
}

type clientEntry struct {
	client      *whatsmeow.Client
	deviceStore *store.Device
	cancelCtx   context.CancelFunc
}

type WAManager struct {
	mu        sync.RWMutex
	clients   map[string]*clientEntry  // deviceID -> clientEntry
	cachedQRs map[string]*cachedQRInfo // deviceID -> cachedQRInfo

	pgxPool   *pgxpool.Pool
	container *sqlstore.Container
	sqlDB     *sql.DB
	log       zerolog.Logger
}

func NewWAManager(pgxPool *pgxpool.Pool, log zerolog.Logger) (*WAManager, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/epmp?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgres for whatsmeow: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping postgres for whatsmeow: %w", err)
	}

	waLogger := waLog.Stdout("WhatsMeow", "INFO", true)
	container := sqlstore.NewWithDB(db, "postgres", waLogger)
	if err := container.Upgrade(context.Background()); err != nil {
		return nil, fmt.Errorf("upgrade whatsmeow schema: %w", err)
	}

	mgr := &WAManager{
		clients:   make(map[string]*clientEntry),
		cachedQRs: make(map[string]*cachedQRInfo),
		pgxPool:   pgxPool,
		container: container,
		sqlDB:     db,
		log:       log,
	}

	// Auto-reconnect previously paired devices in background
	go mgr.AutoReconnectAll(context.Background())

	return mgr, nil
}

// GetQR initiates pairing or returns the current active QR code.
func (m *WAManager) GetQR(ctx context.Context, deviceID string) (string, int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if already active and connected
	if entry, ok := m.clients[deviceID]; ok && entry.client.IsConnected() && entry.client.IsLoggedIn() {
		return "", 0, ErrAlreadyConnected
	}

	// Check if we have a valid cached QR that has not expired
	if q, ok := m.cachedQRs[deviceID]; ok && time.Now().Before(q.ExpiresAt) {
		remaining := int(time.Until(q.ExpiresAt).Seconds())
		if remaining > 3 {
			return q.Code, remaining, nil
		}
	}

	// Disconnect existing client if any
	if entry, ok := m.clients[deviceID]; ok {
		entry.client.Disconnect()
		if entry.cancelCtx != nil {
			entry.cancelCtx()
		}
		delete(m.clients, deviceID)
	}

	// Create fresh device store for new pairing
	deviceStore := m.container.NewDevice()
	waLogger := waLog.Stdout(fmt.Sprintf("WhatsMeow[%s]", deviceID[:8]), "INFO", true)
	client := whatsmeow.NewClient(deviceStore, waLogger)

	// Set event handler
	client.AddEventHandler(func(evt interface{}) {
		m.handleEvent(deviceID, evt)
	})

	qrCtx, cancel := context.WithCancel(context.Background())
	m.clients[deviceID] = &clientEntry{
		client:      client,
		deviceStore: deviceStore,
		cancelCtx:   cancel,
	}

	qrChan, err := client.GetQRChannel(qrCtx)
	if err != nil {
		cancel()
		delete(m.clients, deviceID)
		return "", 0, fmt.Errorf("get QR channel: %w", err)
	}

	if err := client.Connect(); err != nil {
		cancel()
		delete(m.clients, deviceID)
		return "", 0, fmt.Errorf("connect to WhatsApp: %w", err)
	}

	// Wait for first QR event
	select {
	case <-ctx.Done():
		return "", 0, ctx.Err()
	case <-time.After(15 * time.Second):
		return "", 0, errors.New("timeout menunggu QR code dari WhatsApp")
	case item, ok := <-qrChan:
		if !ok {
			return "", 0, errors.New("qr channel tertutup")
		}
		if item.Event == whatsmeow.QRChannelEventCode {
			expiresIn := int(item.Timeout.Seconds())
			if expiresIn <= 0 {
				expiresIn = 60
			}
			m.cachedQRs[deviceID] = &cachedQRInfo{
				Code:      item.Code,
				ExpiresAt: time.Now().Add(time.Duration(expiresIn) * time.Second),
				Timeout:   item.Timeout,
			}

			// Handle subsequent QR codes and success in background
			go m.listenQRChannel(deviceID, client, qrChan)

			return item.Code, expiresIn, nil
		}
		return "", 0, fmt.Errorf("unexpected event from WhatsApp: %s", item.Event)
	}
}

func (m *WAManager) listenQRChannel(deviceID string, client *whatsmeow.Client, qrChan <-chan whatsmeow.QRChannelItem) {
	for item := range qrChan {
		switch item.Event {
		case whatsmeow.QRChannelEventCode:
			m.mu.Lock()
			expiresIn := int(item.Timeout.Seconds())
			if expiresIn <= 0 {
				expiresIn = 60
			}
			m.cachedQRs[deviceID] = &cachedQRInfo{
				Code:      item.Code,
				ExpiresAt: time.Now().Add(time.Duration(expiresIn) * time.Second),
				Timeout:   item.Timeout,
			}
			m.mu.Unlock()
			m.log.Info().Str("device_id", deviceID).Msg("WhatsApp QR code updated")

		case "success":
			m.mu.Lock()
			delete(m.cachedQRs, deviceID)
			m.mu.Unlock()

			var phone string
			var jidStr string
			if client.Store.ID != nil {
				phone = client.Store.ID.User
				jidStr = client.Store.ID.String()
			}

			m.log.Info().Str("device_id", deviceID).Str("phone", phone).Msg("WhatsApp device successfully paired!")

			// Update wa_devices in database
			now := time.Now()
			_, err := m.pgxPool.Exec(context.Background(),
				`UPDATE wa_devices SET status='connected', phone=$1, session_data=$2, last_seen=$3, updated_at=now() WHERE id=$4`,
				phone, jidStr, now, deviceID,
			)
			if err != nil {
				m.log.Error().Err(err).Msg("failed to update wa_device in DB on pair success")
			}
			return

		case "timeout":
			m.mu.Lock()
			delete(m.cachedQRs, deviceID)
			m.mu.Unlock()
			m.log.Warn().Str("device_id", deviceID).Msg("WhatsApp QR pairing timed out")
			m.pgxPool.Exec(context.Background(),
				`UPDATE wa_devices SET status='disconnected', updated_at=now() WHERE id=$1 AND status='qr_pending'`,
				deviceID,
			)
			return

		default:
			m.log.Debug().Str("event", item.Event).Str("device_id", deviceID).Msg("WhatsApp QR event")
		}
	}
}

func (m *WAManager) handleEvent(deviceID string, evt interface{}) {
	switch v := evt.(type) {
	case *events.Connected:
		m.log.Info().Str("device_id", deviceID).Msg("whatsmeow: connected to server")
		m.pgxPool.Exec(context.Background(),
			`UPDATE wa_devices SET status='connected', last_seen=now(), updated_at=now() WHERE id=$1`,
			deviceID,
		)
	case *events.LoggedOut:
		m.log.Warn().Str("device_id", deviceID).Msg("whatsmeow: device logged out from phone")
		m.Disconnect(deviceID)
		m.pgxPool.Exec(context.Background(),
			`UPDATE wa_devices SET status='disconnected', session_data=NULL, updated_at=now() WHERE id=$1`,
			deviceID,
		)
	case *events.Disconnected:
		m.log.Warn().Str("device_id", deviceID).Msg("whatsmeow: disconnected from server")
	default:
		_ = v
	}
}

// Disconnect disconnects an active client.
func (m *WAManager) Disconnect(deviceID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.cachedQRs, deviceID)
	if entry, ok := m.clients[deviceID]; ok {
		entry.client.Disconnect()
		if entry.cancelCtx != nil {
			entry.cancelCtx()
		}
		delete(m.clients, deviceID)
	}
}

// Delete removes the device session entirely.
func (m *WAManager) Delete(ctx context.Context, deviceID string) error {
	m.Disconnect(deviceID)

	var sessionData string
	_ = m.pgxPool.QueryRow(ctx, `SELECT COALESCE(session_data,'') FROM wa_devices WHERE id=$1`, deviceID).Scan(&sessionData)
	if sessionData != "" {
		jid, err := types.ParseJID(sessionData)
		if err == nil && !jid.IsEmpty() {
			deviceStore, err := m.container.GetDevice(ctx, jid)
			if err == nil && deviceStore != nil {
				_ = m.container.DeleteDevice(ctx, deviceStore)
			}
		}
	}

	return nil
}

// getOrReconnectClient returns an active client or reconnects from saved session in DB.
func (m *WAManager) getOrReconnectClient(ctx context.Context, deviceID string) (*clientEntry, error) {
	m.mu.RLock()
	entry, ok := m.clients[deviceID]
	m.mu.RUnlock()

	if ok && entry.client.IsConnected() {
		return entry, nil
	}

	var sessionData string
	err := m.pgxPool.QueryRow(ctx, `SELECT COALESCE(session_data,'') FROM wa_devices WHERE id=$1`, deviceID).Scan(&sessionData)
	if err != nil || sessionData == "" {
		return nil, fmt.Errorf("perangkat WhatsApp %s tidak memiliki sesi atau belum terhubung", deviceID)
	}

	if err := m.reconnectDevice(ctx, deviceID, sessionData); err != nil {
		return nil, fmt.Errorf("gagal menghubungkan kembali perangkat %s: %w", deviceID, err)
	}

	m.mu.RLock()
	entry, ok = m.clients[deviceID]
	m.mu.RUnlock()

	if !ok || !entry.client.IsConnected() {
		return nil, fmt.Errorf("perangkat WhatsApp %s gagal terhubung", deviceID)
	}

	return entry, nil
}

// CheckNumberOnWhatsApp checks if a phone number is registered on WhatsApp.
func (m *WAManager) CheckNumberOnWhatsApp(ctx context.Context, deviceID, phone string) (bool, types.JID, error) {
	entry, err := m.getOrReconnectClient(ctx, deviceID)
	if err != nil {
		return false, types.EmptyJID, err
	}

	cleanPhone := normalizePhone(phone)
	if cleanPhone == "" {
		return false, types.EmptyJID, errors.New("nomor telepon kosong atau tidak valid")
	}

	// IsOnWhatsApp expects international phone format with `+` prefix
	resp, err := entry.client.IsOnWhatsApp(ctx, []string{"+" + cleanPhone})
	if err != nil {
		return false, types.EmptyJID, fmt.Errorf("gagal verifikasi nomor di WhatsApp: %w", err)
	}

	if len(resp) == 0 || !resp[0].IsIn {
		return false, types.EmptyJID, nil
	}

	targetJID := resp[0].JID
	if targetJID.IsEmpty() {
		targetJID = types.NewJID(cleanPhone, types.DefaultUserServer)
	}

	return true, targetJID, nil
}

// SendMessage checks if the recipient number is on WhatsApp first, and only sends if registered.
func (m *WAManager) SendMessage(ctx context.Context, deviceID, phone, text string) error {
	cleanPhone := normalizePhone(phone)
	if cleanPhone == "" {
		return errors.New("nomor telepon kosong atau tidak valid")
	}

	entry, err := m.getOrReconnectClient(ctx, deviceID)
	if err != nil {
		return err
	}

	// 1. Check if recipient number is on WhatsApp
	isOnWA, targetJID, err := m.CheckNumberOnWhatsApp(ctx, deviceID, cleanPhone)
	if err != nil {
		return err
	}
	if !isOnWA {
		return fmt.Errorf("nomor +%s tidak terdaftar di WhatsApp", cleanPhone)
	}

	// 2. Build and send message to the verified JID
	msg := &waE2E.Message{
		Conversation: proto.String(text),
	}

	_, err = entry.client.SendMessage(ctx, targetJID, msg)
	return err
}

func (m *WAManager) reconnectDevice(ctx context.Context, deviceID, sessionData string) error {
	jid, err := types.ParseJID(sessionData)
	if err != nil || jid.IsEmpty() {
		return fmt.Errorf("invalid JID in session data: %s", sessionData)
	}

	deviceStore, err := m.container.GetDevice(ctx, jid)
	if err != nil || deviceStore == nil {
		return fmt.Errorf("device store not found for JID: %s", jid.String())
	}

	waLogger := waLog.Stdout(fmt.Sprintf("WhatsMeow[%s]", deviceID[:8]), "INFO", true)
	client := whatsmeow.NewClient(deviceStore, waLogger)
	client.AddEventHandler(func(evt interface{}) {
		m.handleEvent(deviceID, evt)
	})

	if err := client.Connect(); err != nil {
		return fmt.Errorf("connect failed: %w", err)
	}

	m.mu.Lock()
	m.clients[deviceID] = &clientEntry{
		client:      client,
		deviceStore: deviceStore,
	}
	m.mu.Unlock()

	return nil
}

// AutoReconnectAll reconnects all devices that are marked as connected in the DB.
func (m *WAManager) AutoReconnectAll(ctx context.Context) {
	// Give backend a moment to boot
	time.Sleep(2 * time.Second)

	rows, err := m.pgxPool.Query(ctx,
		`SELECT id, session_data FROM wa_devices WHERE status='connected' AND session_data IS NOT NULL AND session_data != ''`)
	if err != nil {
		m.log.Error().Err(err).Msg("failed to query connected devices for auto-reconnect")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, sessionData string
		if err := rows.Scan(&id, &sessionData); err == nil {
			go func(devID, sData string) {
				if err := m.reconnectDevice(context.Background(), devID, sData); err != nil {
					m.log.Warn().Str("device_id", devID).Err(err).Msg("failed to auto-reconnect device")
				} else {
					m.log.Info().Str("device_id", devID).Msg("auto-reconnected device to WhatsApp")
				}
			}(id, sessionData)
		}
	}
}

func normalizePhone(phone string) string {
	result := ""
	for _, c := range phone {
		if c >= '0' && c <= '9' {
			result += string(c)
		}
	}
	if strings.HasPrefix(result, "0") {
		result = "62" + result[1:]
	}
	return result
}
