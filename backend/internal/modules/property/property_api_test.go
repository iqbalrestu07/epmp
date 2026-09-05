package property_test

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

// testEnv holds the shared test server and DB pool.
type testEnv struct {
	server *httptest.Server
	db     *pgxpool.Pool
	token  string
	orgID  string
}

// responseEnvelope is the standard JSON response wrapper.
type responseEnvelope struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   *responseError  `json:"error"`
}

type responseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// propertyData matches the PropertyResponse DTO.
type propertyData struct {
	OrganizationID string    `json:"organization_id"`
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Address        string    `json:"address"`
	PropertyType   string    `json:"property_type"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type propertyListData struct {
	Data       []propertyData `json:"data"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
	TotalPages int            `json:"total_pages"`
}

// setupTestEnv creates a test server with real DB connection.
// It registers a test user and obtains a JWT token for authenticated requests.
// Skips if DATABASE_URL is not set or DB is unreachable.
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

	// Check DB connectivity
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

	// Register or login a test user to get a token
	token := getTestToken(t, server)

	// Create a test organization and get its ID
	orgID := createTestOrg(t, server, token)

	env := &testEnv{
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

// getTestToken registers a unique test user and returns the access token.
func getTestToken(t *testing.T, server *httptest.Server) string {
	t.Helper()

	email := fmt.Sprintf("test_%d@epmp-test.com", time.Now().UnixNano())
	body := fmt.Sprintf(`{"email":"%s","password":"testpass123","name":"Test User"}`, email)

	resp := doRequest(t, server, "POST", "/api/v1/auth/register", body, "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		// User might already exist, try login
		resp2 := doRequest(t, server, "POST", "/api/v1/auth/login", body, "")
		defer resp2.Body.Close()
		if resp2.StatusCode != http.StatusOK {
			t.Fatalf("failed to register or login test user: register=%d, login=%d", resp.StatusCode, resp2.StatusCode)
		}
		return parseToken(t, resp2.Body)
	}
	return parseToken(t, resp.Body)
}

// createTestOrg creates a test organization and returns its ID.
func createTestOrg(t *testing.T, server *httptest.Server, token string) string {
	t.Helper()

	orgName := fmt.Sprintf("Test Org %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"name":"%s","domain":"%d.test","is_active":true}`, orgName, time.Now().UnixNano())

	resp := doRequest(t, server, "POST", "/api/v1/organizations", body, token)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("failed to create test org: expected 201, got %d", resp.StatusCode)
	}

	var env responseEnvelope
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

// doRequest performs an HTTP request against the test server.
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

// doRequestWithOrg performs an HTTP request with X-Organization-ID header.
func doRequestWithOrg(t *testing.T, server *httptest.Server, method, path, body, token, orgID string) *http.Response {
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

// doAuthRequest performs an authenticated HTTP request with org context.
func (env *testEnv) doAuthRequest(t *testing.T, method, path, body string) *http.Response {
	t.Helper()
	return doRequestWithOrg(t, env.server, method, path, body, env.token, env.orgID)
}

// parseProperty parses a property response envelope.
func parseProperty(t *testing.T, body io.Reader) propertyData {
	t.Helper()
	var env responseEnvelope
	if err := json.NewDecoder(body).Decode(&env); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	var p propertyData
	if err := json.Unmarshal(env.Data, &p); err != nil {
		t.Fatalf("failed to unmarshal property data: %v", err)
	}
	return p
}

// parsePropertyList parses a property list response envelope.
func parsePropertyList(t *testing.T, body io.Reader) propertyListData {
	t.Helper()
	var env responseEnvelope
	if err := json.NewDecoder(body).Decode(&env); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	var list propertyListData
	if err := json.Unmarshal(env.Data, &list); err != nil {
		t.Fatalf("failed to unmarshal property list data: %v", err)
	}
	return list
}

// cleanupProperty deletes a property by ID (soft delete) to keep tests isolated.
func (env *testEnv) cleanupProperty(t *testing.T, id string) {
	t.Helper()
	if id == "" {
		return
	}
	resp := env.doAuthRequest(t, "DELETE", "/api/v1/properties/"+id, "")
	resp.Body.Close()
}

// ─── Tests ──────────────────────────────────────────────────────────────────

func TestPropertyAPI_Create_Success(t *testing.T) {
	env := setupTestEnv(t)

	body := `{"name":"Test Property API","description":"Integration test property","address":"123 Test St","property_type":"boarding_house","is_active":true}`
	resp := env.doAuthRequest(t, "POST", "/api/v1/properties", body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	p := parseProperty(t, resp.Body)
	if p.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if p.Name != "Test Property API" {
		t.Errorf("expected name 'Test Property API', got '%s'", p.Name)
	}
	if p.PropertyType != "boarding_house" {
		t.Errorf("expected property_type 'boarding_house', got '%s'", p.PropertyType)
	}
	if !p.IsActive {
		t.Error("expected is_active=true")
	}
	if p.CreatedAt.IsZero() {
		t.Error("expected non-zero created_at")
	}
	if p.UpdatedAt.IsZero() {
		t.Error("expected non-zero updated_at")
	}

	env.cleanupProperty(t, p.ID)
}

func TestPropertyAPI_Create_InvalidBody(t *testing.T) {
	env := setupTestEnv(t)

	resp := env.doAuthRequest(t, "POST", "/api/v1/properties", "not json")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestPropertyAPI_GetByID_Success(t *testing.T) {
	env := setupTestEnv(t)

	// Create first
	body := `{"name":"Get Test Property","description":"Test","address":"456 Test Ave","property_type":"apartment","is_active":true}`
	resp := env.doAuthRequest(t, "POST", "/api/v1/properties", body)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create failed: expected 201, got %d", resp.StatusCode)
	}
	p := parseProperty(t, resp.Body)
	resp.Body.Close()
	defer env.cleanupProperty(t, p.ID)

	// Get by ID
	resp2 := env.doAuthRequest(t, "GET", "/api/v1/properties/"+p.ID, "")
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp2.StatusCode)
	}

	got := parseProperty(t, resp2.Body)
	if got.ID != p.ID {
		t.Errorf("expected ID '%s', got '%s'", p.ID, got.ID)
	}
	if got.Name != "Get Test Property" {
		t.Errorf("expected name 'Get Test Property', got '%s'", got.Name)
	}
}

func TestPropertyAPI_GetByID_NotFound(t *testing.T) {
	env := setupTestEnv(t)

	resp := env.doAuthRequest(t, "GET", "/api/v1/properties/00000000-0000-0000-0000-000000000000", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestPropertyAPI_List_Pagination(t *testing.T) {
	env := setupTestEnv(t)

	// Create 3 properties
	var ids []string
	for i := 0; i < 3; i++ {
		body := fmt.Sprintf(`{"name":"List Test Property %d","property_type":"villa","is_active":true}`, i)
		resp := env.doAuthRequest(t, "POST", "/api/v1/properties", body)
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("create %d failed: expected 201, got %d", i, resp.StatusCode)
		}
		p := parseProperty(t, resp.Body)
		resp.Body.Close()
		ids = append(ids, p.ID)
	}
	defer func() {
		for _, id := range ids {
			env.cleanupProperty(t, id)
		}
	}()

	// List with per_page=2
	resp := env.doAuthRequest(t, "GET", "/api/v1/properties?per_page=2&page=1", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	list := parsePropertyList(t, resp.Body)
	if len(list.Data) != 2 {
		t.Errorf("expected 2 items, got %d", len(list.Data))
	}
	if list.PerPage != 2 {
		t.Errorf("expected per_page=2, got %d", list.PerPage)
	}
	if list.Page != 1 {
		t.Errorf("expected page=1, got %d", list.Page)
	}
	if list.Total < 3 {
		t.Errorf("expected total>=3, got %d", list.Total)
	}
	if list.TotalPages < 2 {
		t.Errorf("expected total_pages>=2, got %d", list.TotalPages)
	}
}

func TestPropertyAPI_Update_Success(t *testing.T) {
	env := setupTestEnv(t)

	// Create
	body := `{"name":"Update Test Property","property_type":"boarding_house","is_active":true}`
	resp := env.doAuthRequest(t, "POST", "/api/v1/properties", body)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create failed: expected 201, got %d", resp.StatusCode)
	}
	p := parseProperty(t, resp.Body)
	resp.Body.Close()
	defer env.cleanupProperty(t, p.ID)

	// Update
	updateBody := `{"name":"Updated Property Name","description":"Updated description","address":"789 Updated St","property_type":"apartment","is_active":false}`
	resp2 := env.doAuthRequest(t, "PUT", "/api/v1/properties/"+p.ID, updateBody)
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp2.StatusCode)
	}

	updated := parseProperty(t, resp2.Body)
	if updated.Name != "Updated Property Name" {
		t.Errorf("expected name 'Updated Property Name', got '%s'", updated.Name)
	}
	if updated.PropertyType != "apartment" {
		t.Errorf("expected property_type 'apartment', got '%s'", updated.PropertyType)
	}
	if updated.IsActive {
		t.Error("expected is_active=false")
	}
	if updated.Description != "Updated description" {
		t.Errorf("expected description 'Updated description', got '%s'", updated.Description)
	}
}

func TestPropertyAPI_Delete_Success(t *testing.T) {
	env := setupTestEnv(t)

	// Create
	body := `{"name":"Delete Test Property","property_type":"warehouse","is_active":true}`
	resp := env.doAuthRequest(t, "POST", "/api/v1/properties", body)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create failed: expected 201, got %d", resp.StatusCode)
	}
	p := parseProperty(t, resp.Body)
	resp.Body.Close()

	// Delete
	resp2 := env.doAuthRequest(t, "DELETE", "/api/v1/properties/"+p.ID, "")
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", resp2.StatusCode)
	}

	// Verify deleted (should return 404)
	resp3 := env.doAuthRequest(t, "GET", "/api/v1/properties/"+p.ID, "")
	defer resp3.Body.Close()

	if resp3.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", resp3.StatusCode)
	}
}

func TestPropertyAPI_Unauthorized(t *testing.T) {
	env := setupTestEnv(t)

	// Request without token
	resp := doRequest(t, env.server, "GET", "/api/v1/properties", "", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}
