package room_test

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
	floorID    string
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

type roomData struct {
	OrganizationID string    `json:"organization_id"`
	ID             string    `json:"id"`
	PropertyID     string    `json:"property_id"`
	FloorID        string    `json:"floor_id,omitempty"`
	Name           string    `json:"name"`
	Capacity       int       `json:"capacity"`
	Price          float64   `json:"price"`
	IsAvailable    bool      `json:"is_available"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type roomListData struct {
	Data       []roomData `json:"data"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	PerPage    int        `json:"per_page"`
	TotalPages int        `json:"total_pages"`
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
	floorID := createTestFloor(t, server, token, orgID, buildingID)

	env := &testEnv{
		server:     server,
		db:         db,
		token:      token,
		orgID:      orgID,
		propertyID: propertyID,
		buildingID: buildingID,
		floorID:    floorID,
	}

	t.Cleanup(func() {
		server.Close()
		db.Close()
	})

	return env
}

func getTestToken(t *testing.T, server *httptest.Server) string {
	t.Helper()

	email := fmt.Sprintf("test_room_%d@epmp-test.com", time.Now().UnixNano())
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

	orgName := fmt.Sprintf("Test Org Room %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"name":"%s","domain":"%d.room.test","is_active":true}`, orgName, time.Now().UnixNano())

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

	name := fmt.Sprintf("Test Property for Room %d", time.Now().UnixNano())
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

	name := fmt.Sprintf("Test Building for Room %d", time.Now().UnixNano())
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

func createTestFloor(t *testing.T, server *httptest.Server, token, orgID, buildingID string) string {
	t.Helper()

	name := fmt.Sprintf("Test Floor for Room %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"building_id":"%s","name":"%s","floor_number":1,"is_active":true}`, buildingID, name)

	resp := doRequestWithOrg(t, server, "POST", "/api/v1/floors", body, token, orgID)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("failed to create test floor: %d", resp.StatusCode)
	}

	var env responseEnvelope
	var floor struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		t.Fatalf("failed to decode floor response: %v", err)
	}
	if err := json.Unmarshal(env.Data, &floor); err != nil {
		t.Fatalf("failed to unmarshal floor data: %v", err)
	}
	return floor.ID
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

func parseRoom(t *testing.T, body io.Reader) roomData {
	t.Helper()
	var env responseEnvelope
	if err := json.NewDecoder(body).Decode(&env); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	var r roomData
	if err := json.Unmarshal(env.Data, &r); err != nil {
		t.Fatalf("failed to unmarshal room data: %v", err)
	}
	return r
}

func parseRoomList(t *testing.T, body io.Reader) roomListData {
	t.Helper()
	var env responseEnvelope
	if err := json.NewDecoder(body).Decode(&env); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	var list roomListData
	if err := json.Unmarshal(env.Data, &list); err != nil {
		t.Fatalf("failed to unmarshal room list data: %v", err)
	}
	return list
}

func (env *testEnv) cleanupRoom(t *testing.T, id string) {
	t.Helper()
	if id == "" {
		return
	}
	resp := env.doAuthRequest(t, "DELETE", "/api/v1/rooms/"+id, "")
	resp.Body.Close()
}

// ─── Tests ───────────────────────────────────────────────────────────────────

func TestRoomAPI_Create_Success(t *testing.T) {
	env := setupTestEnv(t)

	name := fmt.Sprintf("Room 101 %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"property_id":"%s","floor_id":"%s","name":"%s","capacity":2,"price":500000,"is_available":true}`, env.propertyID, env.floorID, name)

	resp := env.doAuthRequest(t, "POST", "/api/v1/rooms", body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	r := parseRoom(t, resp.Body)
	if r.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if r.Name != name {
		t.Errorf("expected name '%s', got '%s'", name, r.Name)
	}
	if r.Capacity != 2 {
		t.Errorf("expected capacity 2, got %d", r.Capacity)
	}
	if r.Price != 500000 {
		t.Errorf("expected price 500000, got %f", r.Price)
	}
	if !r.IsAvailable {
		t.Error("expected is_available=true")
	}
	if r.CreatedAt.IsZero() {
		t.Error("expected non-zero created_at")
	}
	if r.UpdatedAt.IsZero() {
		t.Error("expected non-zero updated_at")
	}

	env.cleanupRoom(t, r.ID)
}

func TestRoomAPI_Create_MissingName(t *testing.T) {
	env := setupTestEnv(t)

	body := fmt.Sprintf(`{"property_id":"%s","floor_id":"%s","capacity":2,"price":500000,"is_available":true}`, env.propertyID, env.floorID)
	resp := env.doAuthRequest(t, "POST", "/api/v1/rooms", body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestRoomAPI_Create_MissingPropertyId(t *testing.T) {
	env := setupTestEnv(t)

	body := `{"name":"Room No Prop","capacity":2,"price":500000,"is_available":true}`
	resp := env.doAuthRequest(t, "POST", "/api/v1/rooms", body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestRoomAPI_GetByID_Success(t *testing.T) {
	env := setupTestEnv(t)

	name := fmt.Sprintf("Get Room %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"property_id":"%s","floor_id":"%s","name":"%s","capacity":4,"price":750000,"is_available":true}`, env.propertyID, env.floorID, name)

	createResp := env.doAuthRequest(t, "POST", "/api/v1/rooms", body)
	defer createResp.Body.Close()
	r := parseRoom(t, createResp.Body)
	defer env.cleanupRoom(t, r.ID)

	getResp := env.doAuthRequest(t, "GET", "/api/v1/rooms/"+r.ID, "")
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", getResp.StatusCode)
	}

	fetched := parseRoom(t, getResp.Body)
	if fetched.ID != r.ID {
		t.Fatalf("expected ID '%s', got '%s'", r.ID, fetched.ID)
	}
	if fetched.Name != name {
		t.Errorf("expected name '%s', got '%s'", name, fetched.Name)
	}
}

func TestRoomAPI_GetByID_NotFound(t *testing.T) {
	env := setupTestEnv(t)

	resp := env.doAuthRequest(t, "GET", "/api/v1/rooms/nonexistent-id", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestRoomAPI_List_Pagination(t *testing.T) {
	env := setupTestEnv(t)

	var ids []string
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("List Room %d-%d", i, time.Now().UnixNano())
		body := fmt.Sprintf(`{"property_id":"%s","floor_id":"%s","name":"%s","capacity":2,"price":300000,"is_available":true}`, env.propertyID, env.floorID, name)
		resp := env.doAuthRequest(t, "POST", "/api/v1/rooms", body)
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("create %d failed: expected 201, got %d", i, resp.StatusCode)
		}
		r := parseRoom(t, resp.Body)
		resp.Body.Close()
		ids = append(ids, r.ID)
	}
	defer func() {
		for _, id := range ids {
			env.cleanupRoom(t, id)
		}
	}()

	resp := env.doAuthRequest(t, "GET", "/api/v1/rooms?per_page=2&page=1", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	list := parseRoomList(t, resp.Body)
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

func TestRoomAPI_List_FilterByFloorId(t *testing.T) {
	env := setupTestEnv(t)

	name := fmt.Sprintf("Filter Room %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"property_id":"%s","floor_id":"%s","name":"%s","capacity":1,"price":200000,"is_available":true}`, env.propertyID, env.floorID, name)
	resp := env.doAuthRequest(t, "POST", "/api/v1/rooms", body)
	r := parseRoom(t, resp.Body)
	resp.Body.Close()
	defer env.cleanupRoom(t, r.ID)

	listResp := env.doAuthRequest(t, "GET", fmt.Sprintf("/api/v1/rooms?floor_id=%s", env.floorID), "")
	defer listResp.Body.Close()

	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", listResp.StatusCode)
	}

	list := parseRoomList(t, listResp.Body)
	for _, item := range list.Data {
		if item.FloorID != env.floorID {
			t.Errorf("expected all items to have floor_id '%s', got '%s'", env.floorID, item.FloorID)
		}
	}
}

func TestRoomAPI_Update_Success(t *testing.T) {
	env := setupTestEnv(t)

	name := fmt.Sprintf("Update Room %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"property_id":"%s","floor_id":"%s","name":"%s","capacity":2,"price":500000,"is_available":true}`, env.propertyID, env.floorID, name)

	createResp := env.doAuthRequest(t, "POST", "/api/v1/rooms", body)
	defer createResp.Body.Close()
	r := parseRoom(t, createResp.Body)
	defer env.cleanupRoom(t, r.ID)

	updatedName := name + " Updated"
	updateBody := fmt.Sprintf(`{"property_id":"%s","floor_id":"%s","name":"%s","capacity":6,"price":900000,"is_available":false}`, env.propertyID, env.floorID, updatedName)
	updateResp := env.doAuthRequest(t, "PUT", "/api/v1/rooms/"+r.ID, updateBody)
	defer updateResp.Body.Close()

	if updateResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", updateResp.StatusCode)
	}

	updated := parseRoom(t, updateResp.Body)
	if updated.Name != updatedName {
		t.Errorf("expected name '%s', got '%s'", updatedName, updated.Name)
	}
	if updated.Capacity != 6 {
		t.Errorf("expected capacity 6, got %d", updated.Capacity)
	}
	if updated.Price != 900000 {
		t.Errorf("expected price 900000, got %f", updated.Price)
	}
	if updated.IsAvailable {
		t.Error("expected is_available=false")
	}
}

func TestRoomAPI_Delete_Success(t *testing.T) {
	env := setupTestEnv(t)

	name := fmt.Sprintf("Delete Room %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"property_id":"%s","floor_id":"%s","name":"%s","capacity":1,"price":100000,"is_available":true}`, env.propertyID, env.floorID, name)

	createResp := env.doAuthRequest(t, "POST", "/api/v1/rooms", body)
	defer createResp.Body.Close()
	r := parseRoom(t, createResp.Body)

	delResp := env.doAuthRequest(t, "DELETE", "/api/v1/rooms/"+r.ID, "")
	defer delResp.Body.Close()

	if delResp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", delResp.StatusCode)
	}

	getResp := env.doAuthRequest(t, "GET", "/api/v1/rooms/"+r.ID, "")
	defer getResp.Body.Close()
	if getResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", getResp.StatusCode)
	}
}

func TestRoomAPI_Unauthorized(t *testing.T) {
	env := setupTestEnv(t)

	resp := doRequest(t, env.server, "GET", "/api/v1/rooms", "", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}
