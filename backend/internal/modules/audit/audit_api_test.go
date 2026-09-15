package audit_test

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

type auditTestEnv struct {
	server *httptest.Server
	db     *pgxpool.Pool
	token  string
	orgID  string
}

type auditEnvelope struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type auditLogEntry struct {
	ID          string          `json:"id"`
	UserEmail   string          `json:"user_email"`
	Action      string          `json:"action"`
	Module      string          `json:"module"`
	EntityID    string          `json:"entity_id"`
	Method      string          `json:"method"`
	Path        string          `json:"path"`
	StatusCode  int             `json:"status_code"`
	RequestBody json.RawMessage `json:"request_body"`
	CreatedAt   time.Time       `json:"created_at"`
}

type auditListData struct {
	Data       []auditLogEntry `json:"data"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PerPage    int             `json:"per_page"`
	TotalPages int             `json:"total_pages"`
}

func setupAuditTestEnv(t *testing.T) *auditTestEnv {
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
	token := auditTestToken(t, server)
	orgID := auditCreateOrg(t, server, token)

	env := &auditTestEnv{server: server, db: db, token: token, orgID: orgID}
	t.Cleanup(func() {
		server.Close()
		db.Close()
	})
	return env
}

func auditTestToken(t *testing.T, server *httptest.Server) string {
	t.Helper()
	email := fmt.Sprintf("test_audit_%d@epmp-test.com", time.Now().UnixNano())
	body := fmt.Sprintf(`{"email":"%s","password":"testpass123","name":"Test User"}`, email)
	resp := auditDo(t, server, "POST", "/api/v1/auth/register", body, "", "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		resp2 := auditDo(t, server, "POST", "/api/v1/auth/login", body, "", "")
		defer resp2.Body.Close()
		if resp2.StatusCode != http.StatusOK {
			t.Fatalf("register/login failed: %d/%d", resp.StatusCode, resp2.StatusCode)
		}
		return auditToken(t, resp2.Body)
	}
	return auditToken(t, resp.Body)
}

func auditToken(t *testing.T, body io.Reader) string {
	t.Helper()
	var env auditEnvelope
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

func auditCreateOrg(t *testing.T, server *httptest.Server, token string) string {
	t.Helper()
	body := fmt.Sprintf(`{"name":"Audit Org %d","domain":"%d.audit.test","is_active":true}`, time.Now().UnixNano(), time.Now().UnixNano())
	resp := auditDo(t, server, "POST", "/api/v1/organizations", body, token, "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create org failed: %d", resp.StatusCode)
	}
	var env auditEnvelope
	var org struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		t.Fatalf("decode org: %v", err)
	}
	if err := json.Unmarshal(env.Data, &org); err != nil {
		t.Fatalf("unmarshal org: %v", err)
	}
	return org.ID
}

func auditDo(t *testing.T, server *httptest.Server, method, path, body, token, orgID string) *http.Response {
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

func parseAuditList(t *testing.T, body io.Reader) auditListData {
	t.Helper()
	var env auditEnvelope
	if err := json.NewDecoder(body).Decode(&env); err != nil {
		t.Fatalf("decode: %v", err)
	}
	var list auditListData
	if err := json.Unmarshal(env.Data, &list); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return list
}

// ─── Tests ───────────────────────────────────────────────────────────────────

func TestAuditAPI_RecordsMutation(t *testing.T) {
	env := setupAuditTestEnv(t)

	// Perform a mutating request that should be audited.
	propName := fmt.Sprintf("Audited Prop %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"name":"%s","description":"t","address":"t","property_type":"boarding_house","is_active":true}`, propName)
	createResp := auditDo(t, env.server, "POST", "/api/v1/properties", body, env.token, env.orgID)
	createResp.Body.Close()
	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("setup create failed: %d", createResp.StatusCode)
	}

	resp := auditDo(t, env.server, "GET", "/api/v1/audit-logs?module=properties&action=CREATE", "", env.token, env.orgID)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	list := parseAuditList(t, resp.Body)
	found := false
	for _, e := range list.Data {
		if e.Module == "properties" && e.Action == "CREATE" && e.StatusCode == 201 {
			found = true
			if e.Path != "/api/v1/properties" {
				t.Errorf("unexpected path %s", e.Path)
			}
		}
	}
	if !found {
		t.Error("expected audit entry for POST /properties")
	}
}

func TestAuditAPI_DoesNotRecordReads(t *testing.T) {
	env := setupAuditTestEnv(t)

	// Reads are not audited; ensure no bogus rows for this org.
	resp := auditDo(t, env.server, "GET", "/api/v1/audit-logs?module=properties", "", env.token, env.orgID)
	defer resp.Body.Close()
	list := parseAuditList(t, resp.Body)
	for _, e := range list.Data {
		if e.Method == "GET" {
			t.Error("GET request was audited — only mutations should be recorded")
		}
	}
}

func TestAuditAPI_List_CrossOrg(t *testing.T) {
	env := setupAuditTestEnv(t)

	// Create a property under the primary org.
	body := fmt.Sprintf(`{"name":"OrgA Prop %d","description":"t","address":"t","property_type":"boarding_house","is_active":true}`, time.Now().UnixNano())
	createResp := auditDo(t, env.server, "POST", "/api/v1/properties", body, env.token, env.orgID)
	createResp.Body.Close()

	otherOrg := auditCreateOrg(t, env.server, env.token)
	resp := auditDo(t, env.server, "GET", "/api/v1/audit-logs?module=properties", "", env.token, otherOrg)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	list := parseAuditList(t, resp.Body)
	for _, e := range list.Data {
		if e.Module == "properties" && e.Path == "/api/v1/properties" && e.Action == "CREATE" {
			t.Error("cross-org audit list leaked entry from another organization")
		}
	}
}

func TestAuditAPI_Unauthorized(t *testing.T) {
	env := setupAuditTestEnv(t)

	resp := auditDo(t, env.server, "GET", "/api/v1/audit-logs", "", "", "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}
