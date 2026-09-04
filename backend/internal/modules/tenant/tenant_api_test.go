package tenant_test

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

type tenantTestEnv struct {
	server *httptest.Server
	db     *pgxpool.Pool
	token  string
	orgID  string
}

type tenantResponseEnvelope struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   *tenantRespErr  `json:"error"`
}

type tenantRespErr struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type tenantData struct {
	OrganizationID string    `json:"organization_id"`
	ID             string    `json:"id"`
	FullName       string    `json:"full_name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	IdentityNumber string    `json:"identity_number"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type tenantListData struct {
	Data       []tenantData `json:"data"`
	Total      int64        `json:"total"`
	Page       int          `json:"page"`
	PerPage    int          `json:"per_page"`
	TotalPages int          `json:"total_pages"`
}

func setupTenantTestEnv(t *testing.T) *tenantTestEnv {
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
	token := getTenantTestToken(t, server)
	orgID := createTenantTestOrg(t, server, token)

	env := &tenantTestEnv{
		server: server,
		db:     db,
		token:  token,
		orgID:  orgID,
	}

	t.Cleanup(func() {
		server.Close()
		db.Close()
	})

	return env
}

func getTenantTestToken(t *testing.T, server *httptest.Server) string {
	t.Helper()

	email := fmt.Sprintf("test_tenant_%d@epmp-test.com", time.Now().UnixNano())
	body := fmt.Sprintf(`{"email":"%s","password":"testpass123","name":"Test User"}`, email)

	resp := tenantDoRequest(t, server, "POST", "/api/v1/auth/register", body, "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		resp2 := tenantDoRequest(t, server, "POST", "/api/v1/auth/login", body, "")
		defer resp2.Body.Close()
		if resp2.StatusCode != http.StatusOK {
			t.Fatalf("failed to register or login test user: register=%d, login=%d", resp.StatusCode, resp2.StatusCode)
		}
		return tenantParseToken(t, resp2.Body)
	}
	return tenantParseToken(t, resp.Body)
}

func createTenantTestOrg(t *testing.T, server *httptest.Server, token string) string {
	t.Helper()

	orgName := fmt.Sprintf("Test Org Tenant %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"name":"%s","domain":"%d.tenant.test","is_active":true}`, orgName, time.Now().UnixNano())

	resp := tenantDoRequest(t, server, "POST", "/api/v1/organizations", body, token)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("failed to create test org: expected 201, got %d", resp.StatusCode)
	}

	var env tenantResponseEnvelope
	var org struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		t.Fatalf("failed to decode org response: %v", err)
	}
	if err := json.Unmarshal(env.Data, &org); err != nil {
		t.Fatalf("failed to unmarshal org data: %v", err)
	}
	return org.ID
}

func tenantParseToken(t *testing.T, body io.Reader) string {
	t.Helper()
	var env tenantResponseEnvelope
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

func tenantDoRequest(t *testing.T, server *httptest.Server, method, path, body, token string) *http.Response {
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

func tenantDoRequestWithOrg(t *testing.T, server *httptest.Server, method, path, body, token, orgID string) *http.Response {
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
	if orgID != "" {
		req.Header.Set("X-Organization-ID", orgID)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to do request: %v", err)
	}
	return resp
}

func (env *tenantTestEnv) doAuthRequest(t *testing.T, method, path, body string) *http.Response {
	t.Helper()
	return tenantDoRequestWithOrg(t, env.server, method, path, body, env.token, env.orgID)
}

func parseTenant(t *testing.T, body io.Reader) tenantData {
	t.Helper()
	var env tenantResponseEnvelope
	if err := json.NewDecoder(body).Decode(&env); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	var r tenantData
	if err := json.Unmarshal(env.Data, &r); err != nil {
		t.Fatalf("failed to unmarshal tenant data: %v", err)
	}
	return r
}

func parseTenantList(t *testing.T, body io.Reader) tenantListData {
	t.Helper()
	var env tenantResponseEnvelope
	if err := json.NewDecoder(body).Decode(&env); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	var list tenantListData
	if err := json.Unmarshal(env.Data, &list); err != nil {
		t.Fatalf("failed to unmarshal tenant list data: %v", err)
	}
	return list
}

func (env *tenantTestEnv) cleanupTenant(t *testing.T, id string) {
	t.Helper()
	if id == "" {
		return
	}
	resp := env.doAuthRequest(t, "DELETE", "/api/v1/tenants/"+id, "")
	resp.Body.Close()
}

// ─── Tests ───────────────────────────────────────────────────────────────────

func TestTenantAPI_Create_Success(t *testing.T) {
	env := setupTenantTestEnv(t)

	name := fmt.Sprintf("John Doe %d", time.Now().UnixNano())
	email := fmt.Sprintf("john%d@example.com", time.Now().UnixNano())
	identity := fmt.Sprintf("KTP%d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"full_name":"%s","email":"%s","phone":"08123456789","identity_number":"%s","is_active":true}`, name, email, identity)

	resp := env.doAuthRequest(t, "POST", "/api/v1/tenants", body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	r := parseTenant(t, resp.Body)
	if r.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if r.FullName != name {
		t.Errorf("expected full_name '%s', got '%s'", name, r.FullName)
	}
	if r.Email != email {
		t.Errorf("expected email '%s', got '%s'", email, r.Email)
	}
	if r.Phone != "08123456789" {
		t.Errorf("expected phone '08123456789', got '%s'", r.Phone)
	}
	if !r.IsActive {
		t.Error("expected is_active=true")
	}
	if r.CreatedAt.IsZero() {
		t.Error("expected non-zero created_at")
	}
	if r.UpdatedAt.IsZero() {
		t.Error("expected non-zero updated_at")
	}

	env.cleanupTenant(t, r.ID)
}

func TestTenantAPI_Create_MissingFullName(t *testing.T) {
	env := setupTenantTestEnv(t)

	body := `{"email":"missing@example.com","phone":"08123456789","is_active":true}`
	resp := env.doAuthRequest(t, "POST", "/api/v1/tenants", body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestTenantAPI_Create_MissingEmail(t *testing.T) {
	env := setupTenantTestEnv(t)

	body := `{"full_name":"No Email","phone":"08123456789","is_active":true}`
	resp := env.doAuthRequest(t, "POST", "/api/v1/tenants", body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestTenantAPI_GetByID_Success(t *testing.T) {
	env := setupTenantTestEnv(t)

	name := fmt.Sprintf("Get Tenant %d", time.Now().UnixNano())
	email := fmt.Sprintf("get%d@example.com", time.Now().UnixNano())
	identity := fmt.Sprintf("KTP%d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"full_name":"%s","email":"%s","phone":"08123456789","identity_number":"%s","is_active":true}`, name, email, identity)

	createResp := env.doAuthRequest(t, "POST", "/api/v1/tenants", body)
	defer createResp.Body.Close()
	r := parseTenant(t, createResp.Body)
	defer env.cleanupTenant(t, r.ID)

	getResp := env.doAuthRequest(t, "GET", "/api/v1/tenants/"+r.ID, "")
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", getResp.StatusCode)
	}

	fetched := parseTenant(t, getResp.Body)
	if fetched.ID != r.ID {
		t.Fatalf("expected ID '%s', got '%s'", r.ID, fetched.ID)
	}
	if fetched.FullName != name {
		t.Errorf("expected full_name '%s', got '%s'", name, fetched.FullName)
	}
}

func TestTenantAPI_GetByID_NotFound(t *testing.T) {
	env := setupTenantTestEnv(t)

	resp := env.doAuthRequest(t, "GET", "/api/v1/tenants/nonexistent-id", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestTenantAPI_List_Pagination(t *testing.T) {
	env := setupTenantTestEnv(t)

	var ids []string
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("List Tenant %d-%d", i, time.Now().UnixNano())
		email := fmt.Sprintf("list%d-%d@example.com", i, time.Now().UnixNano())
		identity := fmt.Sprintf("KTP-list-%d-%d", i, time.Now().UnixNano())
		body := fmt.Sprintf(`{"full_name":"%s","email":"%s","phone":"08123456789","identity_number":"%s","is_active":true}`, name, email, identity)
		resp := env.doAuthRequest(t, "POST", "/api/v1/tenants", body)
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("create %d failed: expected 201, got %d", i, resp.StatusCode)
		}
		r := parseTenant(t, resp.Body)
		resp.Body.Close()
		ids = append(ids, r.ID)
	}
	defer func() {
		for _, id := range ids {
			env.cleanupTenant(t, id)
		}
	}()

	resp := env.doAuthRequest(t, "GET", "/api/v1/tenants?per_page=2&page=1", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	list := parseTenantList(t, resp.Body)
	if len(list.Data) != 2 {
		t.Errorf("expected 2 items, got %d", len(list.Data))
	}
	if list.PerPage != 2 {
		t.Errorf("expected per_page=2, got %d", list.PerPage)
	}
	if list.Total < 3 {
		t.Errorf("expected total>=3, got %d", list.Total)
	}
	if list.TotalPages < 2 {
		t.Errorf("expected total_pages>=2, got %d", list.TotalPages)
	}
}

func TestTenantAPI_List_Search(t *testing.T) {
	env := setupTenantTestEnv(t)

	uniqueName := fmt.Sprintf("ZxcUniqueTenant %d", time.Now().UnixNano())
	email := fmt.Sprintf("search%d@example.com", time.Now().UnixNano())
	identity := fmt.Sprintf("KTP-search-%d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"full_name":"%s","email":"%s","phone":"08123456789","identity_number":"%s","is_active":true}`, uniqueName, email, identity)

	resp := env.doAuthRequest(t, "POST", "/api/v1/tenants", body)
	r := parseTenant(t, resp.Body)
	resp.Body.Close()
	defer env.cleanupTenant(t, r.ID)

	searchResp := env.doAuthRequest(t, "GET", "/api/v1/tenants?search=ZxcUniqueTenant", "")
	defer searchResp.Body.Close()

	if searchResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", searchResp.StatusCode)
	}

	list := parseTenantList(t, searchResp.Body)
	found := false
	for _, item := range list.Data {
		if item.ID == r.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find created tenant in search results")
	}
}

func TestTenantAPI_Update_Success(t *testing.T) {
	env := setupTenantTestEnv(t)

	name := fmt.Sprintf("Update Tenant %d", time.Now().UnixNano())
	email := fmt.Sprintf("update%d@example.com", time.Now().UnixNano())
	identity := fmt.Sprintf("KTP%d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"full_name":"%s","email":"%s","phone":"08123456789","identity_number":"%s","is_active":true}`, name, email, identity)

	createResp := env.doAuthRequest(t, "POST", "/api/v1/tenants", body)
	defer createResp.Body.Close()
	r := parseTenant(t, createResp.Body)
	defer env.cleanupTenant(t, r.ID)

	updatedName := name + " Updated"
	updatedIdentity := fmt.Sprintf("KTP%d-upd", time.Now().UnixNano())
	updateBody := fmt.Sprintf(`{"full_name":"%s","email":"%s","phone":"08987654321","identity_number":"%s","is_active":false}`, updatedName, email, updatedIdentity)
	updateResp := env.doAuthRequest(t, "PUT", "/api/v1/tenants/"+r.ID, updateBody)
	defer updateResp.Body.Close()

	if updateResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", updateResp.StatusCode)
	}

	updated := parseTenant(t, updateResp.Body)
	if updated.FullName != updatedName {
		t.Errorf("expected full_name '%s', got '%s'", updatedName, updated.FullName)
	}
	if updated.Phone != "08987654321" {
		t.Errorf("expected phone '08987654321', got '%s'", updated.Phone)
	}
	if updated.IsActive {
		t.Error("expected is_active=false")
	}
}

func TestTenantAPI_Delete_Success(t *testing.T) {
	env := setupTenantTestEnv(t)

	name := fmt.Sprintf("Delete Tenant %d", time.Now().UnixNano())
	email := fmt.Sprintf("delete%d@example.com", time.Now().UnixNano())
	identity := fmt.Sprintf("KTP-delete-%d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"full_name":"%s","email":"%s","phone":"08123456789","identity_number":"%s","is_active":true}`, name, email, identity)

	createResp := env.doAuthRequest(t, "POST", "/api/v1/tenants", body)
	defer createResp.Body.Close()
	r := parseTenant(t, createResp.Body)

	delResp := env.doAuthRequest(t, "DELETE", "/api/v1/tenants/"+r.ID, "")
	defer delResp.Body.Close()

	if delResp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", delResp.StatusCode)
	}

	getResp := env.doAuthRequest(t, "GET", "/api/v1/tenants/"+r.ID, "")
	defer getResp.Body.Close()
	if getResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", getResp.StatusCode)
	}
}

func TestTenantAPI_Unauthorized(t *testing.T) {
	env := setupTenantTestEnv(t)

	resp := tenantDoRequest(t, env.server, "GET", "/api/v1/tenants", "", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}
