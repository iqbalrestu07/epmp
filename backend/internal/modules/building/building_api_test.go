package building_test

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
	server     *httptest.Server
	db         *pgxpool.Pool
	token      string
	propertyID string
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

type buildingData struct {
	OrganizationID string    `json:"organization_id"`
	ID             string    `json:"id"`
	PropertyID     string    `json:"property_id"`
	Name           string    `json:"name"`
	TotalFloors    int       `json:"total_floors"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type buildingListData struct {
	Data       []buildingData `json:"data"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
	TotalPages int            `json:"total_pages"`
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
	propertyID := createTestProperty(t, server, token)

	env := &testEnv{
		server:     server,
		db:         db,
		token:      token,
		propertyID: propertyID,
	}

	t.Cleanup(func() {
		server.Close()
		db.Close()
	})

	return env
}

func getTestToken(t *testing.T, server *httptest.Server) string {
	t.Helper()

	email := fmt.Sprintf("test_bldg_%d@epmp-test.com", time.Now().UnixNano())
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

func createTestProperty(t *testing.T, server *httptest.Server, token string) string {
	t.Helper()

	name := fmt.Sprintf("Test Property for Building %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"name":"%s","description":"test","address":"test addr","property_type":"boarding_house","is_active":true}`, name)

	resp := doRequest(t, server, "POST", "/api/v1/properties", body, token)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("failed to create test property: %d", resp.StatusCode)
	}

	var env responseEnvelope
	var prop struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		t.Fatalf("failed to decode property response: %v", err)
	}
	if err := json.Unmarshal(env.Data, &prop); err != nil {
		t.Fatalf("failed to unmarshal property data: %v", err)
	}
	return prop.ID
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

func parseBuilding(t *testing.T, body io.Reader) buildingData {
	t.Helper()
	var env responseEnvelope
	if err := json.NewDecoder(body).Decode(&env); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	var b buildingData
	if err := json.Unmarshal(env.Data, &b); err != nil {
		t.Fatalf("failed to unmarshal building data: %v", err)
	}
	return b
}

func parseBuildingList(t *testing.T, body io.Reader) buildingListData {
	t.Helper()
	var env responseEnvelope
	if err := json.NewDecoder(body).Decode(&env); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	var list buildingListData
	if err := json.Unmarshal(env.Data, &list); err != nil {
		t.Fatalf("failed to unmarshal building list data: %v", err)
	}
	return list
}

// ─── Tests ───────────────────────────────────────────────────────────────────

func TestBuildingAPI_Create_Success(t *testing.T) {
	env := setupTestEnv(t)

	name := fmt.Sprintf("Test Building %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"property_id":"%s","name":"%s","total_floors":5}`, env.propertyID, name)

	resp := env.doAuthRequest(t, "POST", "/api/v1/buildings", body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	b := parseBuilding(t, resp.Body)
	if b.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if b.Name != name {
		t.Fatalf("expected name %s, got %s", name, b.Name)
	}
	if b.TotalFloors != 5 {
		t.Fatalf("expected total_floors 5, got %d", b.TotalFloors)
	}
	if b.PropertyID != env.propertyID {
		t.Fatalf("expected property_id %s, got %s", env.propertyID, b.PropertyID)
	}
	if b.CreatedAt.IsZero() {
		t.Fatal("expected non-zero created_at")
	}
}

func TestBuildingAPI_Create_MissingName(t *testing.T) {
	env := setupTestEnv(t)

	body := fmt.Sprintf(`{"property_id":"%s","total_floors":3}`, env.propertyID)
	resp := env.doAuthRequest(t, "POST", "/api/v1/buildings", body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestBuildingAPI_Create_MissingPropertyId(t *testing.T) {
	env := setupTestEnv(t)

	body := `{"name":"Test Building No Prop","total_floors":3}`
	resp := env.doAuthRequest(t, "POST", "/api/v1/buildings", body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestBuildingAPI_GetByID_Success(t *testing.T) {
	env := setupTestEnv(t)

	name := fmt.Sprintf("Get Building %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"property_id":"%s","name":"%s","total_floors":3}`, env.propertyID, name)

	createResp := env.doAuthRequest(t, "POST", "/api/v1/buildings", body)
	defer createResp.Body.Close()
	b := parseBuilding(t, createResp.Body)

	getResp := env.doAuthRequest(t, "GET", "/api/v1/buildings/"+b.ID, "")
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", getResp.StatusCode)
	}

	fetched := parseBuilding(t, getResp.Body)
	if fetched.ID != b.ID {
		t.Fatalf("expected ID %s, got %s", b.ID, fetched.ID)
	}
}

func TestBuildingAPI_GetByID_NotFound(t *testing.T) {
	env := setupTestEnv(t)

	resp := env.doAuthRequest(t, "GET", "/api/v1/buildings/nonexistent-id", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestBuildingAPI_List_Pagination(t *testing.T) {
	env := setupTestEnv(t)

	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("List Building %d-%d", i, time.Now().UnixNano())
		body := fmt.Sprintf(`{"property_id":"%s","name":"%s","total_floors":2}`, env.propertyID, name)
		resp := env.doAuthRequest(t, "POST", "/api/v1/buildings", body)
		resp.Body.Close()
	}

	resp := env.doAuthRequest(t, "GET", "/api/v1/buildings?page=1&per_page=2", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	list := parseBuildingList(t, resp.Body)
	if len(list.Data) > 2 {
		t.Fatalf("expected at most 2 items, got %d", len(list.Data))
	}
	if list.TotalPages < 1 {
		t.Fatalf("expected total_pages >= 1, got %d", list.TotalPages)
	}
}

func TestBuildingAPI_Update_Success(t *testing.T) {
	env := setupTestEnv(t)

	name := fmt.Sprintf("Update Building %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"property_id":"%s","name":"%s","total_floors":3}`, env.propertyID, name)

	createResp := env.doAuthRequest(t, "POST", "/api/v1/buildings", body)
	defer createResp.Body.Close()
	b := parseBuilding(t, createResp.Body)

	updatedName := name + " Updated"
	updateBody := fmt.Sprintf(`{"property_id":"%s","name":"%s","total_floors":10}`, env.propertyID, updatedName)
	updateResp := env.doAuthRequest(t, "PUT", "/api/v1/buildings/"+b.ID, updateBody)
	defer updateResp.Body.Close()

	if updateResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", updateResp.StatusCode)
	}

	updated := parseBuilding(t, updateResp.Body)
	if updated.Name != updatedName {
		t.Fatalf("expected name %s, got %s", updatedName, updated.Name)
	}
	if updated.TotalFloors != 10 {
		t.Fatalf("expected total_floors 10, got %d", updated.TotalFloors)
	}
}

func TestBuildingAPI_Delete_Success(t *testing.T) {
	env := setupTestEnv(t)

	name := fmt.Sprintf("Delete Building %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"property_id":"%s","name":"%s","total_floors":1}`, env.propertyID, name)

	createResp := env.doAuthRequest(t, "POST", "/api/v1/buildings", body)
	defer createResp.Body.Close()
	b := parseBuilding(t, createResp.Body)

	delResp := env.doAuthRequest(t, "DELETE", "/api/v1/buildings/"+b.ID, "")
	defer delResp.Body.Close()

	if delResp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", delResp.StatusCode)
	}

	getResp := env.doAuthRequest(t, "GET", "/api/v1/buildings/"+b.ID, "")
	defer getResp.Body.Close()
	if getResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", getResp.StatusCode)
	}
}

func TestBuildingAPI_Unauthorized(t *testing.T) {
	env := setupTestEnv(t)

	resp := doRequest(t, env.server, "GET", "/api/v1/buildings", "", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}
