package organization_test

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

type testEnv struct {
	server *httptest.Server
	db     *pgxpool.Pool
	token  string
}

type responseEnvelope struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   *responseError  `json:"error"`
}

type responseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type organizationData struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Domain    string    `json:"domain"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type organizationListData struct {
	Data       []organizationData `json:"data"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PerPage    int                `json:"per_page"`
	TotalPages int                `json:"total_pages"`
}

func setupTestEnv(t *testing.T) *testEnv {
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
	if err := modules.Register(e, db, log, jwtSecret); err != nil {
		t.Fatalf("failed to register modules: %v", err)
	}

	server := httptest.NewServer(e)
	token := getTestToken(t, server)

	env := &testEnv{
		server: server,
		db:     db,
		token:  token,
	}

	t.Cleanup(func() {
		server.Close()
		db.Close()
	})

	return env
}

func getTestToken(t *testing.T, server *httptest.Server) string {
	t.Helper()

	email := fmt.Sprintf("test_org_%d@epmp-test.com", time.Now().UnixNano())
	body := fmt.Sprintf(`{"email":"%s","password":"testpass123","name":"Test User"}`, email)

	resp := doRequest(t, server, "POST", "/api/v1/auth/register", body, "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		resp2 := doRequest(t, server, "POST", "/api/v1/auth/login", body, "")
		defer resp2.Body.Close()
		if resp2.StatusCode != http.StatusOK {
			t.Fatalf("failed to register or login test user: register=%d, login=%d", resp.StatusCode, resp2.StatusCode)
		}
		return parseToken(t, resp2.Body)
	}
	return parseToken(t, resp.Body)
}

func parseToken(t *testing.T, body io.Reader) string {
	t.Helper()
	var env responseEnvelope
	var raw struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(body).Decode(&env); err != nil {
		t.Fatalf("failed to decode token response: %v", err)
	}
	if err := json.Unmarshal(env.Data, &raw); err != nil {
		t.Fatalf("failed to unmarshal token data: %v", err)
	}
	return raw.AccessToken
}

func doRequest(t *testing.T, server *httptest.Server, method, path, body, token string) *http.Response {
	t.Helper()

	var reqBody io.Reader
	if body != "" {
		reqBody = bytes.NewBufferString(body)
	}

	req, err := http.NewRequest(method, server.URL+path, reqBody)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to do request: %v", err)
	}
	return resp
}

func (env *testEnv) doAuthRequest(t *testing.T, method, path, body string) *http.Response {
	t.Helper()
	return doRequest(t, env.server, method, path, body, env.token)
}

func parseOrganization(t *testing.T, body io.Reader) organizationData {
	t.Helper()
	var env responseEnvelope
	if err := json.NewDecoder(body).Decode(&env); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	var o organizationData
	if err := json.Unmarshal(env.Data, &o); err != nil {
		t.Fatalf("failed to unmarshal organization data: %v", err)
	}
	return o
}

func parseOrganizationList(t *testing.T, body io.Reader) organizationListData {
	t.Helper()
	var env responseEnvelope
	if err := json.NewDecoder(body).Decode(&env); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	var list organizationListData
	if err := json.Unmarshal(env.Data, &list); err != nil {
		t.Fatalf("failed to unmarshal organization list data: %v", err)
	}
	return list
}

// ─── Tests ───────────────────────────────────────────────────────────────────

func TestOrganizationAPI_Create_Success(t *testing.T) {
	env := setupTestEnv(t)

	name := fmt.Sprintf("Test Org %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"name":"%s","domain":"test%d.com","is_active":true}`, name, time.Now().UnixNano())

	resp := env.doAuthRequest(t, "POST", "/api/v1/organizations", body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	org := parseOrganization(t, resp.Body)
	if org.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if org.Name != name {
		t.Fatalf("expected name %s, got %s", name, org.Name)
	}
	if !org.IsActive {
		t.Fatal("expected is_active=true")
	}
	if org.CreatedAt.IsZero() {
		t.Fatal("expected non-zero created_at")
	}
	if org.UpdatedAt.IsZero() {
		t.Fatal("expected non-zero updated_at")
	}
}

func TestOrganizationAPI_Create_InvalidBody(t *testing.T) {
	env := setupTestEnv(t)

	resp := env.doAuthRequest(t, "POST", "/api/v1/organizations", `{}`)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestOrganizationAPI_GetByID_Success(t *testing.T) {
	env := setupTestEnv(t)

	name := fmt.Sprintf("Test Org Get %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"name":"%s","domain":"get%d.com","is_active":true}`, name, time.Now().UnixNano())

	createResp := env.doAuthRequest(t, "POST", "/api/v1/organizations", body)
	defer createResp.Body.Close()
	org := parseOrganization(t, createResp.Body)

	getResp := env.doAuthRequest(t, "GET", "/api/v1/organizations/"+org.ID, "")
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", getResp.StatusCode)
	}

	fetched := parseOrganization(t, getResp.Body)
	if fetched.ID != org.ID {
		t.Fatalf("expected ID %s, got %s", org.ID, fetched.ID)
	}
	if fetched.Name != org.Name {
		t.Fatalf("expected name %s, got %s", org.Name, fetched.Name)
	}
}

func TestOrganizationAPI_GetByID_NotFound(t *testing.T) {
	env := setupTestEnv(t)

	resp := env.doAuthRequest(t, "GET", "/api/v1/organizations/nonexistent-id", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestOrganizationAPI_List_Pagination(t *testing.T) {
	env := setupTestEnv(t)

	// Create multiple organizations
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("Pagination Org %d-%d", i, time.Now().UnixNano())
		body := fmt.Sprintf(`{"name":"%s","domain":"pag%d.com","is_active":true}`, name, time.Now().UnixNano())
		resp := env.doAuthRequest(t, "POST", "/api/v1/organizations", body)
		resp.Body.Close()
	}

	resp := env.doAuthRequest(t, "GET", "/api/v1/organizations?page=1&per_page=2", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	list := parseOrganizationList(t, resp.Body)
	if len(list.Data) > 2 {
		t.Fatalf("expected at most 2 items, got %d", len(list.Data))
	}
	if list.TotalPages < 1 {
		t.Fatalf("expected total_pages >= 1, got %d", list.TotalPages)
	}
	if list.Total < 1 {
		t.Fatalf("expected total >= 1, got %d", list.Total)
	}
}

func TestOrganizationAPI_Update_Success(t *testing.T) {
	env := setupTestEnv(t)

	name := fmt.Sprintf("Update Org %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"name":"%s","domain":"upd%d.com","is_active":true}`, name, time.Now().UnixNano())

	createResp := env.doAuthRequest(t, "POST", "/api/v1/organizations", body)
	defer createResp.Body.Close()
	org := parseOrganization(t, createResp.Body)

	updatedName := name + " Updated"
	updateBody := fmt.Sprintf(`{"name":"%s","domain":"upd%d.com","is_active":false}`, updatedName, time.Now().UnixNano())
	updateResp := env.doAuthRequest(t, "PUT", "/api/v1/organizations/"+org.ID, updateBody)
	defer updateResp.Body.Close()

	if updateResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", updateResp.StatusCode)
	}

	updated := parseOrganization(t, updateResp.Body)
	if updated.Name != updatedName {
		t.Fatalf("expected name %s, got %s", updatedName, updated.Name)
	}
	if updated.IsActive {
		t.Fatal("expected is_active=false after update")
	}
}

func TestOrganizationAPI_Delete_Success(t *testing.T) {
	env := setupTestEnv(t)

	name := fmt.Sprintf("Delete Org %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"name":"%s","domain":"del%d.com","is_active":true}`, name, time.Now().UnixNano())

	createResp := env.doAuthRequest(t, "POST", "/api/v1/organizations", body)
	defer createResp.Body.Close()
	org := parseOrganization(t, createResp.Body)

	delResp := env.doAuthRequest(t, "DELETE", "/api/v1/organizations/"+org.ID, "")
	defer delResp.Body.Close()

	if delResp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", delResp.StatusCode)
	}

	// Verify it's gone
	getResp := env.doAuthRequest(t, "GET", "/api/v1/organizations/"+org.ID, "")
	defer getResp.Body.Close()
	if getResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", getResp.StatusCode)
	}
}

func TestOrganizationAPI_Unauthorized(t *testing.T) {
	env := setupTestEnv(t)

	resp := doRequest(t, env.server, "GET", "/api/v1/organizations", "", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}
