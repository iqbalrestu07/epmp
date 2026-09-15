package communication

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	mw "github.com/epmp/backend/internal/pkg/middleware"
)

// ─── Domain Types ─────────────────────────────────────────────────────────────

type WADevice struct {
	ID        string     `json:"id"`
	OrgID     string     `json:"org_id"`
	Label     string     `json:"label"`
	Phone     string     `json:"phone"`
	Status    string     `json:"status"` // connected | disconnected | qr_pending
	LastSeen  *time.Time `json:"last_seen"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type MessageTemplate struct {
	ID        string    `json:"id"`
	OrgID     string    `json:"org_id"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	Variables []string  `json:"variables"`
	Category  string    `json:"category"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BlastMessage struct {
	ID              string      `json:"id"`
	OrgID           string      `json:"org_id"`
	DeviceID        *string     `json:"device_id"`
	Title           string      `json:"title"`
	Template        string      `json:"template"`
	TargetType      string      `json:"target_type"`
	TargetFilter    interface{} `json:"target_filter"`
	Status          string      `json:"status"`
	TotalRecipients int         `json:"total_recipients"`
	SentCount       int         `json:"sent_count"`
	FailedCount     int         `json:"failed_count"`
	ScheduledAt     *time.Time  `json:"scheduled_at"`
	StartedAt       *time.Time  `json:"started_at"`
	CompletedAt     *time.Time  `json:"completed_at"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type BlastMessageLog struct {
	ID            string     `json:"id"`
	BlastID       string     `json:"blast_id"`
	TenantID      *string    `json:"tenant_id"`
	Phone         string     `json:"phone"`
	RecipientName string     `json:"recipient_name"`
	Message       string     `json:"message"`
	Status        string     `json:"status"`
	ErrorMessage  *string    `json:"error_message"`
	SentAt        *time.Time `json:"sent_at"`
	ReadAt        *time.Time `json:"read_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

// ─── Module ───────────────────────────────────────────────────────────────────

type Module struct {
	db    *pgxpool.Pool
	log   zerolog.Logger
	waMgr *WAManager
}

func NewModule(db *pgxpool.Pool, log zerolog.Logger) *Module {
	waMgr, err := NewWAManager(db, log)
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize whatsmeow WAManager")
	}
	return &Module{
		db:    db,
		log:   log,
		waMgr: waMgr,
	}
}

func success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"success": true, "data": data})
}

func created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, map[string]interface{}{"success": true, "data": data})
}

func fail(c echo.Context, code int, msg string) error {
	return c.JSON(code, map[string]interface{}{"success": false, "message": msg})
}

// ─── Routes ───────────────────────────────────────────────────────────────────

func (m *Module) RegisterRoutes(router *echo.Group) {
	g := router.Group("/communication")

	// Devices
	g.GET("/devices", m.listDevices)
	g.POST("/devices", m.createDevice)
	g.GET("/devices/:id/qr", m.getDeviceQR)
	g.PUT("/devices/:id/status", m.updateDeviceStatus)
	g.DELETE("/devices/:id", m.deleteDevice)

	// Templates
	g.GET("/templates", m.listTemplates)
	g.POST("/templates", m.createTemplate)
	g.PUT("/templates/:id", m.updateTemplate)
	g.DELETE("/templates/:id", m.deleteTemplate)

	// Blast messages
	g.GET("/blast", m.listBlasts)
	g.POST("/blast", m.createBlast)
	g.GET("/blast/:id", m.getBlast)
	g.GET("/blast/:id/logs", m.getBlastLogs)
	g.POST("/blast/:id/send", m.sendBlast)

	// Contacts preview (who will receive based on target_type)
	g.GET("/recipients", m.previewRecipients)
}

// ─── Devices ──────────────────────────────────────────────────────────────────

func (m *Module) listDevices(c echo.Context) error {
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return fail(c, 400, "X-Organization-ID header is required")
	}
	rows, err := m.db.Query(context.Background(),
		`SELECT id, org_id, label, COALESCE(phone,''), status, last_seen, created_at, updated_at
		 FROM wa_devices WHERE org_id=$1 ORDER BY created_at DESC`, orgID)
	if err != nil {
		m.log.Error().Err(err).Msg("listDevices")
		return fail(c, 500, "internal error")
	}
	defer rows.Close()

	devices := []WADevice{}
	for rows.Next() {
		var d WADevice
		if err := rows.Scan(&d.ID, &d.OrgID, &d.Label, &d.Phone, &d.Status, &d.LastSeen, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return fail(c, 500, "scan error")
		}
		devices = append(devices, d)
	}
	return success(c, devices)
}

func (m *Module) createDevice(c echo.Context) error {
	var req struct {
		Label string `json:"label"`
	}
	if err := c.Bind(&req); err != nil {
		return fail(c, 400, "invalid request")
	}
	if strings.TrimSpace(req.Label) == "" {
		return fail(c, 400, "label harus diisi")
	}
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return fail(c, 400, "X-Organization-ID header is required")
	}

	var d WADevice
	err := m.db.QueryRow(context.Background(),
		`INSERT INTO wa_devices (org_id, label, phone, status)
		 VALUES ($1, $2, '', 'disconnected')
		 RETURNING id, org_id, label, COALESCE(phone,''), status, last_seen, created_at, updated_at`,
		orgID, req.Label,
	).Scan(&d.ID, &d.OrgID, &d.Label, &d.Phone, &d.Status, &d.LastSeen, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		m.log.Error().Err(err).Msg("createDevice")
		return fail(c, 500, "gagal menyimpan device")
	}
	return created(c, d)
}

func (m *Module) getDeviceQR(c echo.Context) error {
	id := c.Param("id")
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return fail(c, 400, "X-Organization-ID header is required")
	}

	if m.waMgr == nil {
		return fail(c, 500, "WhatsApp manager tidak dapat diinisialisasi")
	}

	// Mark device as qr_pending in DB (org-scoped)
	tag, err := m.db.Exec(context.Background(),
		`UPDATE wa_devices SET status='qr_pending', updated_at=now() WHERE id=$1 AND org_id=$2`, id, orgID)
	if err != nil || tag.RowsAffected() == 0 {
		return fail(c, 404, "device not found")
	}

	qrCode, expiresIn, err := m.waMgr.GetQR(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, ErrAlreadyConnected) {
			return success(c, map[string]interface{}{
				"device_id": id,
				"status":    "connected",
				"message":   "Perangkat sudah terhubung",
			})
		}
		m.log.Error().Err(err).Str("device_id", id).Msg("failed to get real WhatsApp QR")
		return fail(c, 500, fmt.Sprintf("Gagal mendapatkan QR dari WhatsApp: %v", err))
	}

	return success(c, map[string]interface{}{
		"device_id":  id,
		"qr_string":  qrCode,
		"expires_in": expiresIn,
	})
}

func (m *Module) updateDeviceStatus(c echo.Context) error {
	id := c.Param("id")
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return fail(c, 400, "X-Organization-ID header is required")
	}
	var req struct {
		Status string `json:"status"`
		Phone  string `json:"phone"`
	}
	if err := c.Bind(&req); err != nil {
		return fail(c, 400, "invalid request")
	}

	if req.Status == "disconnected" && m.waMgr != nil {
		m.waMgr.Disconnect(id)
	}

	var lastSeen *time.Time
	if req.Status == "connected" {
		now := time.Now()
		lastSeen = &now
	}

	_, err := m.db.Exec(context.Background(),
		`UPDATE wa_devices SET status=$1, phone=COALESCE(NULLIF($2,''), phone), last_seen=COALESCE($3, last_seen), updated_at=now() WHERE id=$4 AND org_id=$5`,
		req.Status, req.Phone, lastSeen, id, orgID)
	if err != nil {
		return fail(c, 500, "failed to update status")
	}

	var d WADevice
	err = m.db.QueryRow(context.Background(),
		`SELECT id, org_id, label, COALESCE(phone,''), status, last_seen, created_at, updated_at FROM wa_devices WHERE id=$1 AND org_id=$2`, id, orgID,
	).Scan(&d.ID, &d.OrgID, &d.Label, &d.Phone, &d.Status, &d.LastSeen, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return fail(c, 404, "device not found")
	}
	return success(c, d)
}

func (m *Module) deleteDevice(c echo.Context) error {
	id := c.Param("id")
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return fail(c, 400, "X-Organization-ID header is required")
	}
	if m.waMgr != nil {
		_ = m.waMgr.Delete(c.Request().Context(), id)
	}
	tag, err := m.db.Exec(context.Background(), `DELETE FROM wa_devices WHERE id=$1 AND org_id=$2`, id, orgID)
	if err != nil {
		return fail(c, 500, "failed to delete")
	}
	if tag.RowsAffected() == 0 {
		return fail(c, 404, "device not found")
	}
	return success(c, map[string]string{"message": "device berhasil dihapus"})
}

// ─── Templates ────────────────────────────────────────────────────────────────

func (m *Module) listTemplates(c echo.Context) error {
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return fail(c, 400, "X-Organization-ID header is required")
	}
	rows, err := m.db.Query(context.Background(),
		`SELECT id, org_id, name, content, COALESCE(variables, '{}'), category, is_active, created_at, updated_at
		 FROM message_templates WHERE is_active=true AND org_id=$1 ORDER BY category, name`, orgID)
	if err != nil {
		return fail(c, 500, "internal error")
	}
	defer rows.Close()

	templates := []MessageTemplate{}
	for rows.Next() {
		var t MessageTemplate
		if err := rows.Scan(&t.ID, &t.OrgID, &t.Name, &t.Content, &t.Variables, &t.Category, &t.IsActive, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return fail(c, 500, "scan error")
		}
		templates = append(templates, t)
	}
	return success(c, templates)
}

func (m *Module) createTemplate(c echo.Context) error {
	var req struct {
		Name      string   `json:"name"`
		Content   string   `json:"content"`
		Category  string   `json:"category"`
		Variables []string `json:"variables"`
	}
	if err := c.Bind(&req); err != nil {
		return fail(c, 400, "invalid request")
	}
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Content) == "" {
		return fail(c, 400, "nama dan konten template harus diisi")
	}
	if req.Category == "" {
		req.Category = "general"
	}
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return fail(c, 400, "X-Organization-ID header is required")
	}

	var t MessageTemplate
	err := m.db.QueryRow(context.Background(),
		`INSERT INTO message_templates (org_id, name, content, variables, category)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, org_id, name, content, COALESCE(variables,'{}'), category, is_active, created_at, updated_at`,
		orgID, req.Name, req.Content, req.Variables, req.Category,
	).Scan(&t.ID, &t.OrgID, &t.Name, &t.Content, &t.Variables, &t.Category, &t.IsActive, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		m.log.Error().Err(err).Msg("createTemplate")
		return fail(c, 500, "gagal menyimpan template")
	}
	return created(c, t)
}

func (m *Module) updateTemplate(c echo.Context) error {
	id := c.Param("id")
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return fail(c, 400, "X-Organization-ID header is required")
	}
	var req struct {
		Name      string   `json:"name"`
		Content   string   `json:"content"`
		Category  string   `json:"category"`
		Variables []string `json:"variables"`
		IsActive  bool     `json:"is_active"`
	}
	if err := c.Bind(&req); err != nil {
		return fail(c, 400, "invalid request")
	}

	var t MessageTemplate
	err := m.db.QueryRow(context.Background(),
		`UPDATE message_templates SET name=$1, content=$2, category=$3, variables=$4, is_active=$5, updated_at=now()
		 WHERE id=$6 AND org_id=$7
		 RETURNING id, org_id, name, content, COALESCE(variables,'{}'), category, is_active, created_at, updated_at`,
		req.Name, req.Content, req.Category, req.Variables, req.IsActive, id, orgID,
	).Scan(&t.ID, &t.OrgID, &t.Name, &t.Content, &t.Variables, &t.Category, &t.IsActive, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fail(c, 404, "template tidak ditemukan")
		}
		return fail(c, 500, "gagal mengupdate template")
	}
	return success(c, t)
}

func (m *Module) deleteTemplate(c echo.Context) error {
	id := c.Param("id")
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return fail(c, 400, "X-Organization-ID header is required")
	}
	_, err := m.db.Exec(context.Background(),
		`UPDATE message_templates SET is_active=false, updated_at=now() WHERE id=$1 AND org_id=$2`, id, orgID)
	if err != nil {
		return fail(c, 500, "gagal menghapus template")
	}
	return success(c, map[string]string{"message": "template berhasil dihapus"})
}

// ─── Blast Messages ───────────────────────────────────────────────────────────

func (m *Module) listBlasts(c echo.Context) error {
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return fail(c, 400, "X-Organization-ID header is required")
	}
	rows, err := m.db.Query(context.Background(),
		`SELECT b.id, b.org_id, b.device_id, b.title, b.template, b.target_type,
		        b.status, b.total_recipients, b.sent_count, b.failed_count,
		        b.scheduled_at, b.started_at, b.completed_at, b.created_at, b.updated_at
		 FROM blast_messages b WHERE b.org_id=$1 ORDER BY b.created_at DESC LIMIT 100`, orgID)
	if err != nil {
		return fail(c, 500, "internal error")
	}
	defer rows.Close()

	blasts := []BlastMessage{}
	for rows.Next() {
		var b BlastMessage
		if err := rows.Scan(
			&b.ID, &b.OrgID, &b.DeviceID, &b.Title, &b.Template, &b.TargetType,
			&b.Status, &b.TotalRecipients, &b.SentCount, &b.FailedCount,
			&b.ScheduledAt, &b.StartedAt, &b.CompletedAt, &b.CreatedAt, &b.UpdatedAt,
		); err != nil {
			return fail(c, 500, "scan error")
		}
		blasts = append(blasts, b)
	}
	return success(c, blasts)
}

func (m *Module) createBlast(c echo.Context) error {
	var req struct {
		DeviceID     *string                `json:"device_id"`
		Title        string                 `json:"title"`
		Template     string                 `json:"template"`
		TargetType   string                 `json:"target_type"`
		TargetFilter map[string]interface{} `json:"target_filter"`
		ScheduledAt  *time.Time             `json:"scheduled_at"`
	}
	if err := c.Bind(&req); err != nil {
		return fail(c, 400, "invalid request")
	}
	if strings.TrimSpace(req.Template) == "" {
		return fail(c, 400, "template pesan harus diisi")
	}
	if req.TargetType == "" {
		req.TargetType = "all_tenants"
	}

	var filterJSON []byte
	if req.TargetFilter != nil {
		filterJSON, _ = json.Marshal(req.TargetFilter)
	}
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return fail(c, 400, "X-Organization-ID header is required")
	}

	var b BlastMessage
	err := m.db.QueryRow(context.Background(),
		`INSERT INTO blast_messages (org_id, device_id, title, template, target_type, target_filter, status, scheduled_at)
		 VALUES ($1, $2, $3, $4, $5, $6, 'draft', $7)
		 RETURNING id, org_id, device_id, title, template, target_type,
		           status, total_recipients, sent_count, failed_count,
		           scheduled_at, started_at, completed_at, created_at, updated_at`,
		orgID,
		req.DeviceID, req.Title, req.Template, req.TargetType,
		filterJSON, req.ScheduledAt,
	).Scan(
		&b.ID, &b.OrgID, &b.DeviceID, &b.Title, &b.Template, &b.TargetType,
		&b.Status, &b.TotalRecipients, &b.SentCount, &b.FailedCount,
		&b.ScheduledAt, &b.StartedAt, &b.CompletedAt, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		m.log.Error().Err(err).Msg("createBlast")
		return fail(c, 500, "gagal menyimpan blast message")
	}
	return created(c, b)
}

func (m *Module) getBlast(c echo.Context) error {
	id := c.Param("id")
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return fail(c, 400, "X-Organization-ID header is required")
	}
	var b BlastMessage
	err := m.db.QueryRow(context.Background(),
		`SELECT id, org_id, device_id, title, template, target_type,
		        status, total_recipients, sent_count, failed_count,
		        scheduled_at, started_at, completed_at, created_at, updated_at
		 FROM blast_messages WHERE id=$1 AND org_id=$2`, id, orgID,
	).Scan(
		&b.ID, &b.OrgID, &b.DeviceID, &b.Title, &b.Template, &b.TargetType,
		&b.Status, &b.TotalRecipients, &b.SentCount, &b.FailedCount,
		&b.ScheduledAt, &b.StartedAt, &b.CompletedAt, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fail(c, 404, "blast tidak ditemukan")
		}
		return fail(c, 500, "internal error")
	}
	return success(c, b)
}

func (m *Module) getBlastLogs(c echo.Context) error {
	id := c.Param("id")
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return fail(c, 400, "X-Organization-ID header is required")
	}
	rows, err := m.db.Query(context.Background(),
		`SELECT l.id, l.blast_id, l.tenant_id, l.phone, l.recipient_name, l.message, l.status, l.error_message, l.sent_at, l.read_at, l.created_at
		 FROM blast_message_logs l
		 JOIN blast_messages b ON b.id = l.blast_id
		 WHERE l.blast_id=$1 AND b.org_id=$2 ORDER BY l.created_at DESC`, id, orgID)
	if err != nil {
		return fail(c, 500, "internal error")
	}
	defer rows.Close()

	logs := []BlastMessageLog{}
	for rows.Next() {
		var l BlastMessageLog
		if err := rows.Scan(
			&l.ID, &l.BlastID, &l.TenantID, &l.Phone, &l.RecipientName,
			&l.Message, &l.Status, &l.ErrorMessage, &l.SentAt, &l.ReadAt, &l.CreatedAt,
		); err != nil {
			return fail(c, 500, "scan error")
		}
		logs = append(logs, l)
	}
	return success(c, logs)
}

// sendBlast: load recipients based on target_type, create logs, and send
func (m *Module) sendBlast(c echo.Context) error {
	id := c.Param("id")
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return fail(c, 400, "X-Organization-ID header is required")
	}

	// Load blast
	var b BlastMessage
	var templateContent string
	var targetType string
	var deviceID *string
	err := m.db.QueryRow(context.Background(),
		`SELECT id, template, target_type, device_id FROM blast_messages WHERE id=$1 AND org_id=$2`, id, orgID,
	).Scan(&b.ID, &templateContent, &targetType, &deviceID)
	if err != nil {
		return fail(c, 404, "blast tidak ditemukan")
	}

	// If no device specified, auto-detect any connected device in this org
	if deviceID == nil || *deviceID == "" {
		var activeDevID string
		err := m.db.QueryRow(context.Background(),
			`SELECT id FROM wa_devices WHERE status='connected' AND org_id=$1 ORDER BY updated_at DESC LIMIT 1`, orgID).Scan(&activeDevID)
		if err == nil && activeDevID != "" {
			deviceID = &activeDevID
			m.db.Exec(context.Background(), `UPDATE blast_messages SET device_id=$1 WHERE id=$2`, activeDevID, id)
		}
	}

	// Verify an explicitly-specified device belongs to this org.
	if deviceID != nil && *deviceID != "" {
		var devOrg string
		err := m.db.QueryRow(context.Background(),
			`SELECT org_id FROM wa_devices WHERE id=$1`, *deviceID).Scan(&devOrg)
		if err != nil || devOrg != orgID {
			return fail(c, 404, "device tidak ditemukan pada organisasi ini")
		}
	}

	if deviceID == nil || *deviceID == "" || m.waMgr == nil {
		m.db.Exec(context.Background(), `UPDATE blast_messages SET status='failed', updated_at=now() WHERE id=$1`, id)
		return fail(c, 400, "Tidak ada perangkat WhatsApp yang terhubung. Silakan hubungkan perangkat terlebih dahulu di menu Pengaturan Gateway.")
	}

	// Mark as sending
	now := time.Now()
	m.db.Exec(context.Background(),
		`UPDATE blast_messages SET status='sending', started_at=$1, updated_at=now() WHERE id=$2`, now, id)

	// Load tenants with room info based on target_type (org-scoped)
	query := `
		SELECT t.id, t.full_name, COALESCE(t.phone,'') as phone,
		       COALESCE(r.name,'') as room_name,
		       COALESCE(c.monthly_rent::text, '0') as rent_amount
		FROM tenants t
		LEFT JOIN contracts c ON c.tenant_id = t.id AND lower(c.status)='active' AND c.deleted_at IS NULL
		LEFT JOIN rooms r ON r.id = c.room_id AND r.deleted_at IS NULL
		WHERE t.is_active = true AND COALESCE(t.phone,'') != ''
		  AND t.organization_id = $1 AND t.deleted_at IS NULL`

	if targetType == "overdue" {
		query = `
		SELECT DISTINCT t.id, t.full_name, COALESCE(t.phone,'') as phone,
		       COALESCE(r.name,'') as room_name,
		       COALESCE(c.monthly_rent::text, '0') as rent_amount
		FROM tenants t
		JOIN invoices i ON i.tenant_id = t.id AND lower(i.status) != 'paid' AND i.deleted_at IS NULL
		LEFT JOIN contracts c ON c.tenant_id = t.id AND lower(c.status)='active' AND c.deleted_at IS NULL
		LEFT JOIN rooms r ON r.id = c.room_id AND r.deleted_at IS NULL
		WHERE t.is_active = true AND COALESCE(t.phone,'') != ''
		  AND t.organization_id = $1 AND t.deleted_at IS NULL`
	}

	rows, err := m.db.Query(context.Background(), query, orgID)
	if err != nil {
		m.log.Error().Err(err).Msg("sendBlast: query tenants")
		m.db.Exec(context.Background(),
			`UPDATE blast_messages SET status='failed', updated_at=now() WHERE id=$1`, id)
		return fail(c, 500, "gagal memuat penerima")
	}
	defer rows.Close()

	type recipient struct {
		TenantID   string
		Name       string
		Phone      string
		RoomName   string
		RentAmount string
	}
	var recipients []recipient
	for rows.Next() {
		var r recipient
		if err := rows.Scan(&r.TenantID, &r.Name, &r.Phone, &r.RoomName, &r.RentAmount); err == nil {
			recipients = append(recipients, r)
		}
	}
	rows.Close()

	if len(recipients) == 0 {
		m.db.Exec(context.Background(), `UPDATE blast_messages SET status='done', total_recipients=0, sent_count=0, failed_count=0, completed_at=now(), updated_at=now() WHERE id=$1`, id)
		return success(c, map[string]interface{}{
			"message":          "Tidak ada penerima aktif dengan nomor telepon yang ditemukan.",
			"total_recipients": 0,
			"sent_count":       0,
			"failed_count":     0,
		})
	}

	// Create log entries and send via whatsmeow
	sentCount := 0
	failedCount := 0
	for _, r := range recipients {
		// Render template variables
		msg := renderTemplate(templateContent, map[string]string{
			"tenant_name":    r.Name,
			"room_name":      r.RoomName,
			"invoice_amount": formatRupiah(r.RentAmount),
			"rent_amount":    formatRupiah(r.RentAmount),
			"due_date":       time.Now().AddDate(0, 0, 7).Format("02 Jan 2006"),
			"days_left":      "7",
			"start_date":     time.Now().Format("02 Jan 2006"),
			"message":        "Silakan hubungi pengelola untuk informasi lebih lanjut.",
		})

		logStatus := "sent"
		var errorMsg *string
		sentAt := time.Now()
		logID := uuid.New().String()
		tenantID := r.TenantID

		// Send real WhatsApp message with prior registration check
		sendErr := m.waMgr.SendMessage(context.Background(), *deviceID, r.Phone, msg)
		if sendErr != nil {
			m.log.Warn().Err(sendErr).Str("phone", r.Phone).Msg("whatsmeow send failed")
			logStatus = "failed"
			em := sendErr.Error()
			errorMsg = &em
			failedCount++
		} else {
			sentCount++
		}

		_, _ = m.db.Exec(context.Background(),
			`INSERT INTO blast_message_logs (id, blast_id, tenant_id, phone, recipient_name, message, status, error_message, sent_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			logID, id, tenantID, r.Phone, r.Name, msg, logStatus, errorMsg, sentAt,
		)
	}

	// Mark as done
	completedAt := time.Now()
	m.db.Exec(context.Background(),
		`UPDATE blast_messages 
		 SET status='done', total_recipients=$1, sent_count=$2, failed_count=$3, completed_at=$4, updated_at=now()
		 WHERE id=$5`,
		len(recipients), sentCount, failedCount, completedAt, id)

	return success(c, map[string]interface{}{
		"message":          fmt.Sprintf("Blast selesai. %d berhasil, %d gagal dari %d penerima.", sentCount, failedCount, len(recipients)),
		"total_recipients": len(recipients),
		"sent_count":       sentCount,
		"failed_count":     failedCount,
	})
}

func (m *Module) previewRecipients(c echo.Context) error {
	targetType := c.QueryParam("target_type")
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return fail(c, 400, "X-Organization-ID header is required")
	}

	query := `
		SELECT t.id, t.full_name, COALESCE(t.phone,'') as phone, COALESCE(r.name,'') as room_name
		FROM tenants t
		LEFT JOIN contracts c ON c.tenant_id = t.id AND lower(c.status)='active' AND c.deleted_at IS NULL
		LEFT JOIN rooms r ON r.id = c.room_id AND r.deleted_at IS NULL
		WHERE t.is_active = true AND t.organization_id = $1 AND t.deleted_at IS NULL`

	if targetType == "overdue" {
		query = `
		SELECT DISTINCT t.id, t.full_name, COALESCE(t.phone,'') as phone, COALESCE(r.name,'') as room_name
		FROM tenants t
		JOIN invoices i ON i.tenant_id = t.id AND lower(i.status) != 'paid' AND i.deleted_at IS NULL
		LEFT JOIN contracts c ON c.tenant_id = t.id AND lower(c.status)='active' AND c.deleted_at IS NULL
		LEFT JOIN rooms r ON r.id = c.room_id AND r.deleted_at IS NULL
		WHERE t.is_active = true AND t.organization_id = $1 AND t.deleted_at IS NULL`
	}

	rows, err := m.db.Query(context.Background(), query, orgID)
	if err != nil {
		return fail(c, 500, "internal error")
	}
	defer rows.Close()

	type contact struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Phone    string `json:"phone"`
		RoomName string `json:"room_name"`
	}
	contacts := []contact{}
	for rows.Next() {
		var contact contact
		if err := rows.Scan(&contact.ID, &contact.Name, &contact.Phone, &contact.RoomName); err == nil {
			contacts = append(contacts, contact)
		}
	}
	return success(c, map[string]interface{}{
		"recipients": contacts,
		"total":      len(contacts),
	})
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func renderTemplate(tmpl string, vars map[string]string) string {
	result := tmpl
	for k, v := range vars {
		result = strings.ReplaceAll(result, "{{"+k+"}}", v)
	}
	return result
}

func formatRupiah(s string) string {
	// Simple formatting: just return as-is for now
	// In production use a proper currency formatter
	return s
}
