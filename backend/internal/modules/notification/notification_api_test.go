package notification_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/epmp/backend/internal/modules"
	"github.com/epmp/backend/internal/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type notifTestEnv struct {
	server    *httptest.Server
	db        *pgxpool.Pool
	token     string
	orgID     string
	tenantID  string
	invoiceID string
}

type notifEnvelope struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type notifEntry struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Link      string    `json:"link"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type notifList struct {
	Data        []notifEntry `json:"data"`
	Total       int64        `json:"total"`
	UnreadCount int64        `json:"unread_count"`
}

func setupNotifTestEnv(t *testing.T) *notifTestEnv {
	t.Helper()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/epmp?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Skipf("skipping integration test: cannot connect to DB: %v", err)
	}
	if err := db.Ping(ctx); err != nil {
		t.Skipf("skipping integration test: DB ping failed: %v", err)
	}

	log := logger.New()
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	jwtSecret := "test-jwt-secret-for-integration-tests"
	accessTTL := 2 * time.Hour
	refreshTTL := 30 * 24 * time.Hour
	if err := modules.Register(e, db, log, jwtSecret, accessTTL, refreshTTL); err != nil {
		t.Fatalf("failed to register modules: %v", err)
	}

	server := httptest.NewServer(e)
	token := notifTestToken(t, server)
	orgID := notifCreateOrg(t, server, token)

	propertyID := notifCreate(t, server, token, orgID, "/api/v1/properties",
		fmt.Sprintf(`{"name":"Notif Prop %d","description":"t","address":"t","property_type":"boarding_house","is_active":true}`, time.Now().UnixNano()))
	buildingID := notifCreate(t, server, token, orgID, "/api/v1/buildings",
		fmt.Sprintf(`{"property_id":"%s","name":"Notif Bldg %d","total_floors":3}`, propertyID, time.Now().UnixNano()))
	floorID := notifCreate(t, server, token, orgID, "/api/v1/floors",
		fmt.Sprintf(`{"building_id":"%s","name":"NF1 %d","floor_number":1,"is_active":true}`, buildingID, time.Now().UnixNano()))
	roomID := notifCreate(t, server, token, orgID, "/api/v1/rooms",
		fmt.Sprintf(`{"property_id":"%s","floor_id":"%s","name":"NR-%d","capacity":2,"price":500000,"is_available":true}`, propertyID, floorID, time.Now().UnixNano()))
	tenantID := notifCreate(t, server, token, orgID, "/api/v1/tenants",
		fmt.Sprintf(`{"full_name":"Notif Tenant %d","email":"ntf%d@example.com","phone":"0812","identity_number":"KTP-ntf-%d","is_active":true}`,
			time.Now().UnixNano(), time.Now().UnixNano(), time.Now().UnixNano()))
	contractID := notifCreate(t, server, token, orgID, "/api/v1/contracts",
		fmt.Sprintf(`{"tenant_id":"%s","property_id":"%s","room_id":"%s","status":"Active","start_date":"2026-01-01T00:00:00Z","end_date":"2026-12-31T00:00:00Z","monthly_rent":500000}`,
			tenantID, propertyID, roomID))
	invoiceID := notifCreate(t, server, token, orgID, "/api/v1/invoices",
		fmt.Sprintf(`{"contract_id":"%s","tenant_id":"%s","amount":500000,"currency":"IDR","status":"Unpaid","due_date":"%s"}`, contractID, tenantID, time.Now().AddDate(0, 0, 30).Format(time.RFC3339)))

	env := &notifTestEnv{server: server, db: db, token: token, orgID: orgID, tenantID: tenantID, invoiceID: invoiceID}
	t.Cleanup(func() {
		server.Close()
		db.Close()
	})
	return env
}

func notifTestToken(t *testing.T, server *httptest.Server) string {
	t.Helper()
	email := fmt.Sprintf("test_notif_%d@epmp-test.com", time.Now().UnixNano())
	body := fmt.Sprintf(`{"email":"%s","password":"testpass123","name":"Test User"}`, email)
	resp := notifDo(t, server, "POST", "/api/v1/auth/register", body, "", "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		resp2 := notifDo(t, server, "POST", "/api/v1/auth/login", body, "", "")
		defer resp2.Body.Close()
		if resp2.StatusCode != http.StatusOK {
			t.Fatalf("register/login failed: %d/%d", resp.StatusCode, resp2.StatusCode)
		}
		return notifToken(t, resp2.Body)
	}
	return notifToken(t, resp.Body)
}

func notifToken(t *testing.T, body io.Reader) string {
	t.Helper()
	var env notifEnvelope
	var raw struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(body).Decode(&env); err != nil {
		t.Fatalf("decode token: %v", err)
	}
	if err := json.Unmarshal(env.Data, &raw); err != nil {
		t.Fatalf("unmarshal token: %v", err)
	}
	return raw.AccessToken
}

func notifCreateOrg(t *testing.T, server *httptest.Server, token string) string {
	t.Helper()
	body := fmt.Sprintf(`{"name":"Notif Org %d","domain":"%d.notif.test","is_active":true}`, time.Now().UnixNano(), time.Now().UnixNano())
	resp := notifDo(t, server, "POST", "/api/v1/organizations", body, token, "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create org failed: %d", resp.StatusCode)
	}
	return notifID(t, resp.Body)
}

func notifCreate(t *testing.T, server *httptest.Server, token, orgID, path, body string) string {
	t.Helper()
	resp := notifDo(t, server, "POST", path, body, token, orgID)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("create %s failed: %d, body: %s", path, resp.StatusCode, string(b))
	}
	return notifID(t, resp.Body)
}

func notifID(t *testing.T, body io.Reader) string {
	t.Helper()
	var env notifEnvelope
	var raw struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(body).Decode(&env); err != nil {
		t.Fatalf("decode id: %v", err)
	}
	if err := json.Unmarshal(env.Data, &raw); err != nil {
		t.Fatalf("unmarshal id: %v", err)
	}
	return raw.ID
}

func notifDo(t *testing.T, server *httptest.Server, method, path, body, token, orgID string) *http.Response {
	t.Helper()
	var reqBody io.Reader
	if body != "" {
		reqBody = bytes.NewBufferString(body)
	}
	req, err := http.NewRequest(method, server.URL+path, reqBody)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if orgID != "" {
		req.Header.Set("X-Organization-ID", orgID)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	return resp
}

func (env *notifTestEnv) listNotifications(t *testing.T, orgID string) notifList {
	t.Helper()
	resp := notifDo(t, env.server, "GET", "/api/v1/notifications", "", env.token, orgID)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list notifications: expected 200, got %d", resp.StatusCode)
	}
	var env2 notifEnvelope
	var list notifList
	if err := json.NewDecoder(resp.Body).Decode(&env2); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if err := json.Unmarshal(env2.Data, &list); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return list
}

// ─── Tests ───────────────────────────────────────────────────────────────────

func TestNotificationAPI_PaymentCreatesNotification(t *testing.T) {
	env := setupNotifTestEnv(t)

	// Recording a payment should emit an org-wide notification.
	payBody := fmt.Sprintf(`{"invoice_id":"%s","tenant_id":"%s","amount":500000,"payment_date":"%s","payment_method":"Transfer","status":"Success","reference_number":"REF-%d"}`,
		env.invoiceID, env.tenantID, time.Now().Format(time.RFC3339), time.Now().UnixNano())
	payResp := notifDo(t, env.server, "POST", "/api/v1/payments", payBody, env.token, env.orgID)
	payResp.Body.Close()
	if payResp.StatusCode != http.StatusCreated {
		t.Fatalf("payment create failed: %d", payResp.StatusCode)
	}

	list := env.listNotifications(t, env.orgID)
	if list.UnreadCount < 1 {
		t.Fatalf("expected unread_count>=1, got %d", list.UnreadCount)
	}
	found := false
	for _, n := range list.Data {
		if n.Type == "payment" && n.Link != "" {
			found = true
		}
	}
	if !found {
		t.Error("expected a payment notification after recording a payment")
	}
}

func TestNotificationAPI_MarkRead(t *testing.T) {
	env := setupNotifTestEnv(t)

	payBody := fmt.Sprintf(`{"invoice_id":"%s","tenant_id":"%s","amount":100000,"payment_date":"%s","payment_method":"Cash","status":"Success"}`,
		env.invoiceID, env.tenantID, time.Now().Format(time.RFC3339))
	payResp := notifDo(t, env.server, "POST", "/api/v1/payments", payBody, env.token, env.orgID)
	payResp.Body.Close()

	list := env.listNotifications(t, env.orgID)
	if len(list.Data) == 0 {
		t.Fatal("expected at least one notification")
	}

	target := list.Data[0]
	readResp := notifDo(t, env.server, "POST", "/api/v1/notifications/"+target.ID+"/read", "", env.token, env.orgID)
	readResp.Body.Close()
	if readResp.StatusCode != http.StatusNoContent {
		t.Fatalf("mark read: expected 204, got %d", readResp.StatusCode)
	}

	list = env.listNotifications(t, env.orgID)
	for _, n := range list.Data {
		if n.ID == target.ID && !n.IsRead {
			t.Error("expected notification marked as read")
		}
	}
}

func TestNotificationAPI_MarkAllRead(t *testing.T) {
	env := setupNotifTestEnv(t)

	payBody := fmt.Sprintf(`{"invoice_id":"%s","tenant_id":"%s","amount":100000,"payment_date":"%s","payment_method":"Cash","status":"Success"}`,
		env.invoiceID, env.tenantID, time.Now().Format(time.RFC3339))
	payResp := notifDo(t, env.server, "POST", "/api/v1/payments", payBody, env.token, env.orgID)
	payResp.Body.Close()

	resp := notifDo(t, env.server, "POST", "/api/v1/notifications/read-all", "", env.token, env.orgID)
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("mark all read: expected 204, got %d", resp.StatusCode)
	}

	list := env.listNotifications(t, env.orgID)
	if list.UnreadCount != 0 {
		t.Errorf("expected unread_count=0, got %d", list.UnreadCount)
	}
}

func TestNotificationAPI_CrossOrg(t *testing.T) {
	env := setupNotifTestEnv(t)

	payBody := fmt.Sprintf(`{"invoice_id":"%s","tenant_id":"%s","amount":100000,"payment_date":"%s","payment_method":"Cash","status":"Success"}`,
		env.invoiceID, env.tenantID, time.Now().Format(time.RFC3339))
	payResp := notifDo(t, env.server, "POST", "/api/v1/payments", payBody, env.token, env.orgID)
	payResp.Body.Close()

	otherOrg := notifCreateOrg(t, env.server, env.token)
	list := env.listNotifications(t, otherOrg)
	for _, n := range list.Data {
		if n.Type == "payment" {
			t.Error("cross-org notification list leaked payment notification")
		}
	}
}

func TestNotificationAPI_Unauthorized(t *testing.T) {
	env := setupNotifTestEnv(t)

	resp := notifDo(t, env.server, "GET", "/api/v1/notifications", "", "", "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}
