package floor_test

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
	orgID      string
	propertyID string
	buildingID string
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

type floorData struct {
	OrganizationID string    `json:"organization_id"`
	ID             string    `json:"id"`
	BuildingID     string    `json:"building_id"`
	Name           string    `json:"name"`
	FloorNumber    int       `json:"floor_number"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type floorListData struct {
	Data       []floorData `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalPages int         `json:"total_pages"`
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
	orgID := createTestOrg(t, server, token)
	propertyID := createTestProperty(t, server, token, orgID)
	buildingID := createTestBuilding(t, server, token, orgID, propertyID)

	env := &testEnv{
		server:     server,
		db:         db,
		token:      token,
		orgID:      orgID,
		propertyID: propertyID,
		buildingID: buildingID,
	}

	t.Cleanup(func() {
		server.Close()
		db.Close()
	})

	return env
}

func getTestToken(t *testing.T, server *httptest.Server) string {
	t.Helper()

	email := fmt.Sprintf("test_floor_%d@epmp-test.com", time.Now().UnixNano())
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

func createTestOrg(t *testing.T, server *httptest.Server, token string) string {
	t.Helper()

	orgName := fmt.Sprintf("Test Org Floor %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"name":"%s","domain":"%d.floor.test","is_active":true}`, orgName, time.Now().UnixNano())

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

func createTestProperty(t *testing.T, server *httptest.Server, token, orgID string) string {
	t.Helper()

	name := fmt.Sprintf("Test Property for Floor %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"name":"%s","description":"test","address":"test addr","property_type":"boarding_house","is_active":true}`, name)

	resp := doRequestWithOrg(t, server, "POST", "/api/v1/properties", body, token, orgID)
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

func createTestBuilding(t *testing.T, server *httptest.Server, token, orgID, propertyID string) string {
	t.Helper()

	name := fmt.Sprintf("Test Building for Floor %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"property_id":"%s","name":"%s","total_floors":5}`, propertyID, name)

	resp := doRequestWithOrg(t, server, "POST", "/api/v1/buildings", body, token, orgID)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("failed to create test building: %d", resp.StatusCode)
	}

	var env responseEnvelope
	var bldg struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		t.Fatalf("failed to decode building response: %v", err)
	}
	if err := json.Unmarshal(env.Data, &bldg); err != nil {
		t.Fatalf("failed to unmarshal building data: %v", err)
	}
	return bldg.ID
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

func (env *testEnv) doAuthRequest(t *testing.T, method, path, body string) *http.Response {
	t.Helper()
	return doRequestWithOrg(t, env.server, method, path, body, env.token, env.orgID)
}

func parseFloor(t *testing.T, body io.Reader) floorData {
	t.Helper()
	var env responseEnvelope
	if err := json.NewDecoder(body).Decode(&env); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	var f floorData
	if err := json.Unmarshal(env.Data, &f); err != nil {
		t.Fatalf("failed to unmarshal floor data: %v", err)
	}
	return f
}

func parseFloorList(t *testing.T, body io.Reader) floorListData {
	t.Helper()
	var env responseEnvelope
	if err := json.NewDecoder(body).Decode(&env); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	var list floorListData
	if err := json.Unmarshal(env.Data, &list); err != nil {
		t.Fatalf("failed to unmarshal floor list data: %v", err)
	}
	return list
}

func (env *testEnv) cleanupFloor(t *testing.T, id string) {
	t.Helper()
	if id == "" {
		return
	}
	resp := env.doAuthRequest(t, "DELETE", "/api/v1/floors/"+id, "")
	resp.Body.Close()
}

// ─── Tests ───────────────────────────────────────────────────────────────────

func TestFloorAPI_Create_Success(t *testing.T) {
	env := setupTestEnv(t)

	name := fmt.Sprintf("Floor 1 %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"building_id":"%s","name":"%s","floor_number":1,"is_active":true}`, env.buildingID, name)

	resp := env.doAuthRequest(t, "POST", "/api/v1/floors", body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	f := parseFloor(t, resp.Body)
	if f.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if f.Name != name {
		t.Errorf("expected name '%s', got '%s'", name, f.Name)
	}
	if f.FloorNumber != 1 {
		t.Errorf("expected floor_number 1, got %d", f.FloorNumber)
	}
	if f.BuildingID != env.buildingID {
		t.Errorf("expected building_id '%s', got '%s'", env.buildingID, f.BuildingID)
	}
	if !f.IsActive {
		t.Error("expected is_active=true")
	}
	if f.CreatedAt.IsZero() {
		t.Error("expected non-zero created_at")
	}
	if f.UpdatedAt.IsZero() {
		t.Error("expected non-zero updated_at")
	}

	env.cleanupFloor(t, f.ID)
}

func TestFloorAPI_Create_MissingName(t *testing.T) {
	env := setupTestEnv(t)

	body := fmt.Sprintf(`{"building_id":"%s","floor_number":1,"is_active":true}`, env.buildingID)
	resp := env.doAuthRequest(t, "POST", "/api/v1/floors", body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestFloorAPI_Create_MissingBuildingId(t *testing.T) {
	env := setupTestEnv(t)

	body := `{"name":"Floor No Building","floor_number":1,"is_active":true}`
	resp := env.doAuthRequest(t, "POST", "/api/v1/floors", body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestFloorAPI_GetByID_Success(t *testing.T) {
	env := setupTestEnv(t)

	name := fmt.Sprintf("Get Floor %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"building_id":"%s","name":"%s","floor_number":2,"is_active":true}`, env.buildingID, name)

	createResp := env.doAuthRequest(t, "POST", "/api/v1/floors", body)
	defer createResp.Body.Close()
	f := parseFloor(t, createResp.Body)
	defer env.cleanupFloor(t, f.ID)

	getResp := env.doAuthRequest(t, "GET", "/api/v1/floors/"+f.ID, "")
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", getResp.StatusCode)
	}

	fetched := parseFloor(t, getResp.Body)
	if fetched.ID != f.ID {
		t.Fatalf("expected ID '%s', got '%s'", f.ID, fetched.ID)
	}
	if fetched.Name != name {
		t.Errorf("expected name '%s', got '%s'", name, fetched.Name)
	}
}

func TestFloorAPI_GetByID_NotFound(t *testing.T) {
	env := setupTestEnv(t)

	resp := env.doAuthRequest(t, "GET", "/api/v1/floors/nonexistent-id", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestFloorAPI_List_Pagination(t *testing.T) {
	env := setupTestEnv(t)

	var ids []string
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("List Floor %d-%d", i, time.Now().UnixNano())
		body := fmt.Sprintf(`{"building_id":"%s","name":"%s","floor_number":%d,"is_active":true}`, env.buildingID, name, i+10)
		resp := env.doAuthRequest(t, "POST", "/api/v1/floors", body)
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("create %d failed: expected 201, got %d", i, resp.StatusCode)
		}
		f := parseFloor(t, resp.Body)
		resp.Body.Close()
		ids = append(ids, f.ID)
	}
	defer func() {
		for _, id := range ids {
			env.cleanupFloor(t, id)
		}
	}()

	resp := env.doAuthRequest(t, "GET", "/api/v1/floors?per_page=2&page=1", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	list := parseFloorList(t, resp.Body)
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

func TestFloorAPI_List_FilterByBuildingId(t *testing.T) {
	env := setupTestEnv(t)

	name := fmt.Sprintf("Filter Floor %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"building_id":"%s","name":"%s","floor_number":99,"is_active":true}`, env.buildingID, name)
	resp := env.doAuthRequest(t, "POST", "/api/v1/floors", body)
	f := parseFloor(t, resp.Body)
	resp.Body.Close()
	defer env.cleanupFloor(t, f.ID)

	listResp := env.doAuthRequest(t, "GET", fmt.Sprintf("/api/v1/floors?building_id=%s", env.buildingID), "")
	defer listResp.Body.Close()

	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", listResp.StatusCode)
	}

	list := parseFloorList(t, listResp.Body)
	for _, item := range list.Data {
		if item.BuildingID != env.buildingID {
			t.Errorf("expected all items to have building_id '%s', got '%s'", env.buildingID, item.BuildingID)
		}
	}
}

func TestFloorAPI_Update_Success(t *testing.T) {
	env := setupTestEnv(t)

	name := fmt.Sprintf("Update Floor %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"building_id":"%s","name":"%s","floor_number":3,"is_active":true}`, env.buildingID, name)

	createResp := env.doAuthRequest(t, "POST", "/api/v1/floors", body)
	defer createResp.Body.Close()
	f := parseFloor(t, createResp.Body)
	defer env.cleanupFloor(t, f.ID)

	updatedName := name + " Updated"
	updateBody := fmt.Sprintf(`{"building_id":"%s","name":"%s","floor_number":7,"is_active":false}`, env.buildingID, updatedName)
	updateResp := env.doAuthRequest(t, "PUT", "/api/v1/floors/"+f.ID, updateBody)
	defer updateResp.Body.Close()

	if updateResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", updateResp.StatusCode)
	}

	updated := parseFloor(t, updateResp.Body)
	if updated.Name != updatedName {
		t.Errorf("expected name '%s', got '%s'", updatedName, updated.Name)
	}
	if updated.FloorNumber != 7 {
		t.Errorf("expected floor_number 7, got %d", updated.FloorNumber)
	}
	if updated.IsActive {
		t.Error("expected is_active=false")
	}
}

func TestFloorAPI_Delete_Success(t *testing.T) {
	env := setupTestEnv(t)

	name := fmt.Sprintf("Delete Floor %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"building_id":"%s","name":"%s","floor_number":1,"is_active":true}`, env.buildingID, name)

	createResp := env.doAuthRequest(t, "POST", "/api/v1/floors", body)
	defer createResp.Body.Close()
	f := parseFloor(t, createResp.Body)

	delResp := env.doAuthRequest(t, "DELETE", "/api/v1/floors/"+f.ID, "")
	defer delResp.Body.Close()

	if delResp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", delResp.StatusCode)
	}

	getResp := env.doAuthRequest(t, "GET", "/api/v1/floors/"+f.ID, "")
	defer getResp.Body.Close()
	if getResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", getResp.StatusCode)
	}
}

func TestFloorAPI_Unauthorized(t *testing.T) {
	env := setupTestEnv(t)

	resp := doRequest(t, env.server, "GET", "/api/v1/floors", "", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}
