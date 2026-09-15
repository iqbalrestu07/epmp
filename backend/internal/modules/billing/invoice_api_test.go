package billing_test

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

type invoiceTestEnv struct {
	server     *httptest.Server
	db         *pgxpool.Pool
	token      string
	orgID      string
	tenantID   string
	contractID string
}

type invoiceResponseEnvelope struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   *invoiceRespErr `json:"error"`
}

type invoiceRespErr struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type invoiceData struct {
	OrganizationID string    `json:"organization_id"`
	ID             string    `json:"id"`
	ContractID     string    `json:"contract_id"`
	TenantID       string    `json:"tenant_id"`
	Amount         float64   `json:"amount"`
	Currency       string    `json:"currency"`
	Status         string    `json:"status"`
	DueDate        time.Time `json:"due_date"`
	PaidDate       time.Time `json:"paid_date"`
	PaymentMethod  string    `json:"payment_method"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type invoiceListData struct {
	Data       []invoiceData `json:"data"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	PerPage    int           `json:"per_page"`
	TotalPages int           `json:"total_pages"`
}

func setupInvoiceTestEnv(t *testing.T) *invoiceTestEnv {
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
	token := getInvoiceTestToken(t, server)
	orgID := createInvoiceTestOrg(t, server, token)
	propertyID := createInvoiceTestProperty(t, server, token, orgID)
	buildingID := createInvoiceTestBuilding(t, server, token, orgID, propertyID)
	floorID := createInvoiceTestFloor(t, server, token, orgID, buildingID)
	roomID := createInvoiceTestRoom(t, server, token, orgID, propertyID, floorID)
	tenantID := createInvoiceTestTenant(t, server, token, orgID)
	contractID := createInvoiceTestContract(t, server, token, orgID, tenantID, propertyID, roomID)

	env := &invoiceTestEnv{
		server:     server,
		db:         db,
		token:      token,
		orgID:      orgID,
		tenantID:   tenantID,
		contractID: contractID,
	}

	t.Cleanup(func() {
		server.Close()
		db.Close()
	})

	return env
}

func getInvoiceTestToken(t *testing.T, server *httptest.Server) string {
	t.Helper()

	email := fmt.Sprintf("test_invoice_%d@epmp-test.com", time.Now().UnixNano())
	body := fmt.Sprintf(`{"email":"%s","password":"testpass123","name":"Test User"}`, email)

	resp := invoiceDoRequest(t, server, "POST", "/api/v1/auth/register", body, "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		resp2 := invoiceDoRequest(t, server, "POST", "/api/v1/auth/login", body, "")
		defer resp2.Body.Close()
		if resp2.StatusCode != http.StatusOK {
			t.Fatalf("failed to register or login test user: register=%d, login=%d", resp.StatusCode, resp2.StatusCode)
		}
		return invoiceParseToken(t, resp2.Body)
	}
	return invoiceParseToken(t, resp.Body)
}

func createInvoiceTestOrg(t *testing.T, server *httptest.Server, token string) string {
	t.Helper()

	orgName := fmt.Sprintf("Test Org Invoice %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"name":"%s","domain":"%d.invoice.test","is_active":true}`, orgName, time.Now().UnixNano())

	resp := invoiceDoRequest(t, server, "POST", "/api/v1/organizations", body, token)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("failed to create test org: expected 201, got %d", resp.StatusCode)
	}

	var env invoiceResponseEnvelope
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

func createInvoiceTestProperty(t *testing.T, server *httptest.Server, token, orgID string) string {
	t.Helper()

	name := fmt.Sprintf("Test Property for Invoice %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"name":"%s","description":"test","address":"test addr","property_type":"boarding_house","is_active":true}`, name)

	resp := invoiceDoRequestWithOrg(t, server, "POST", "/api/v1/properties", body, token, orgID)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("failed to create test property: %d", resp.StatusCode)
	}
	return invoiceParseID(t, resp.Body)
}

func createInvoiceTestBuilding(t *testing.T, server *httptest.Server, token, orgID, propertyID string) string {
	t.Helper()

	name := fmt.Sprintf("Test Building for Invoice %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"property_id":"%s","name":"%s","total_floors":5}`, propertyID, name)

	resp := invoiceDoRequestWithOrg(t, server, "POST", "/api/v1/buildings", body, token, orgID)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("failed to create test building: %d", resp.StatusCode)
	}
	return invoiceParseID(t, resp.Body)
}

func createInvoiceTestFloor(t *testing.T, server *httptest.Server, token, orgID, buildingID string) string {
	t.Helper()

	name := fmt.Sprintf("Test Floor for Invoice %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"building_id":"%s","name":"%s","floor_number":1,"is_active":true}`, buildingID, name)

	resp := invoiceDoRequestWithOrg(t, server, "POST", "/api/v1/floors", body, token, orgID)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("failed to create test floor: %d", resp.StatusCode)
	}
	return invoiceParseID(t, resp.Body)
}

func createInvoiceTestRoom(t *testing.T, server *httptest.Server, token, orgID, propertyID, floorID string) string {
	t.Helper()

	name := fmt.Sprintf("Test Room for Invoice %d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"property_id":"%s","floor_id":"%s","name":"%s","capacity":2,"price":500000,"is_available":true}`, propertyID, floorID, name)

	resp := invoiceDoRequestWithOrg(t, server, "POST", "/api/v1/rooms", body, token, orgID)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("failed to create test room: %d", resp.StatusCode)
	}
	return invoiceParseID(t, resp.Body)
}

func createInvoiceTestTenant(t *testing.T, server *httptest.Server, token, orgID string) string {
	t.Helper()

	name := fmt.Sprintf("Test Tenant for Invoice %d", time.Now().UnixNano())
	email := fmt.Sprintf("inv_tenant_%d@example.com", time.Now().UnixNano())
	identity := fmt.Sprintf("KTP-inv-%d", time.Now().UnixNano())
	body := fmt.Sprintf(`{"full_name":"%s","email":"%s","phone":"08123456789","identity_number":"%s","is_active":true}`, name, email, identity)

	resp := invoiceDoRequestWithOrg(t, server, "POST", "/api/v1/tenants", body, token, orgID)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("failed to create test tenant: %d", resp.StatusCode)
	}
	return invoiceParseID(t, resp.Body)
}

func createInvoiceTestContract(t *testing.T, server *httptest.Server, token, orgID, tenantID, propertyID, roomID string) string {
	t.Helper()

	body := fmt.Sprintf(`{"tenant_id":"%s","property_id":"%s","room_id":"%s","status":"Active","start_date":"2026-01-01T00:00:00Z","end_date":"2026-12-31T00:00:00Z","monthly_rent":500000,"deposit_amount":500000}`, tenantID, propertyID, roomID)

	resp := invoiceDoRequestWithOrg(t, server, "POST", "/api/v1/contracts", body, token, orgID)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("failed to create test contract: %d, body: %s", resp.StatusCode, string(b))
	}
	return invoiceParseID(t, resp.Body)
}

func invoiceParseToken(t *testing.T, body io.Reader) string {
	t.Helper()
	var env invoiceResponseEnvelope
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

func invoiceParseID(t *testing.T, body io.Reader) string {
	t.Helper()
	var env invoiceResponseEnvelope
	var raw struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(body).Decode(&env); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if err := json.Unmarshal(env.Data, &raw); err != nil {
		t.Fatalf("failed to unmarshal id: %v", err)
	}
	return raw.ID
}

func invoiceDoRequest(t *testing.T, server *httptest.Server, method, path, body, token string) *http.Response {
	t.Helper()
	return invoiceDoRequestWithOrg(t, server, method, path, body, token, "")
}

func invoiceDoRequestWithOrg(t *testing.T, server *httptest.Server, method, path, body, token, orgID string) *http.Response {
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

func (env *invoiceTestEnv) doAuthRequest(t *testing.T, method, path, body string) *http.Response {
	t.Helper()
	return invoiceDoRequestWithOrg(t, env.server, method, path, body, env.token, env.orgID)
}

func parseInvoice(t *testing.T, body io.Reader) invoiceData {
	t.Helper()
	var env invoiceResponseEnvelope
	if err := json.NewDecoder(body).Decode(&env); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	var r invoiceData
	if err := json.Unmarshal(env.Data, &r); err != nil {
		t.Fatalf("failed to unmarshal invoice data: %v", err)
	}
	return r
}

func parseInvoiceList(t *testing.T, body io.Reader) invoiceListData {
	t.Helper()
	var env invoiceResponseEnvelope
	if err := json.NewDecoder(body).Decode(&env); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	var list invoiceListData
	if err := json.Unmarshal(env.Data, &list); err != nil {
		t.Fatalf("failed to unmarshal invoice list data: %v", err)
	}
	return list
}

func (env *invoiceTestEnv) cleanupInvoice(t *testing.T, id string) {
	t.Helper()
	if id == "" {
		return
	}
	resp := env.doAuthRequest(t, "DELETE", "/api/v1/invoices/"+id, "")
	resp.Body.Close()
}

func (env *invoiceTestEnv) createInvoice(t *testing.T, status string) invoiceData {
	t.Helper()
	body := fmt.Sprintf(`{"contract_id":"%s","tenant_id":"%s","amount":500000,"currency":"IDR","status":"%s","due_date":"2026-12-31T00:00:00Z","payment_method":"Transfer","notes":"test invoice"}`, env.contractID, env.tenantID, status)

	resp := env.doAuthRequest(t, "POST", "/api/v1/invoices", body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("failed to create invoice: expected 201, got %d, body: %s", resp.StatusCode, string(b))
	}
	return parseInvoice(t, resp.Body)
}

// ─── Tests ───────────────────────────────────────────────────────────────────

func TestInvoiceAPI_Create_Success(t *testing.T) {
	env := setupInvoiceTestEnv(t)

	body := fmt.Sprintf(`{"contract_id":"%s","tenant_id":"%s","amount":750000,"currency":"USD","status":"Unpaid","due_date":"2026-12-31T00:00:00Z","payment_method":"Transfer","notes":"first invoice"}`, env.contractID, env.tenantID)

	resp := env.doAuthRequest(t, "POST", "/api/v1/invoices", body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 201, got %d, body: %s", resp.StatusCode, string(b))
	}

	r := parseInvoice(t, resp.Body)
	if r.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if r.OrganizationID != env.orgID {
		t.Errorf("expected organization_id '%s', got '%s'", env.orgID, r.OrganizationID)
	}
	if r.ContractID != env.contractID {
		t.Errorf("expected contract_id '%s', got '%s'", env.contractID, r.ContractID)
	}
	if r.Amount != 750000 {
		t.Errorf("expected amount 750000, got %f", r.Amount)
	}
	if r.Currency != "USD" {
		t.Errorf("expected currency 'USD', got '%s'", r.Currency)
	}
	if r.CreatedAt.IsZero() {
		t.Error("expected non-zero created_at")
	}
	if r.UpdatedAt.IsZero() {
		t.Error("expected non-zero updated_at")
	}

	env.cleanupInvoice(t, r.ID)
}

func TestInvoiceAPI_Create_MissingOrgHeader(t *testing.T) {
	env := setupInvoiceTestEnv(t)

	body := fmt.Sprintf(`{"contract_id":"%s","tenant_id":"%s","amount":500000,"status":"Unpaid","due_date":"2026-12-31T00:00:00Z"}`, env.contractID, env.tenantID)
	resp := invoiceDoRequestWithOrg(t, env.server, "POST", "/api/v1/invoices", body, env.token, "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestInvoiceAPI_Create_IgnoresBodyOrgID(t *testing.T) {
	env := setupInvoiceTestEnv(t)

	// Even if client sends a foreign organization_id, the header org must win.
	body := fmt.Sprintf(`{"organization_id":"00000000-0000-0000-0000-000000000000","contract_id":"%s","tenant_id":"%s","amount":500000,"status":"Unpaid","due_date":"2026-12-31T00:00:00Z"}`, env.contractID, env.tenantID)

	resp := env.doAuthRequest(t, "POST", "/api/v1/invoices", body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 201, got %d, body: %s", resp.StatusCode, string(b))
	}

	r := parseInvoice(t, resp.Body)
	defer env.cleanupInvoice(t, r.ID)

	if r.OrganizationID != env.orgID {
		t.Errorf("expected organization_id '%s' (from header), got '%s'", env.orgID, r.OrganizationID)
	}
}

func TestInvoiceAPI_GetByID_Success(t *testing.T) {
	env := setupInvoiceTestEnv(t)

	r := env.createInvoice(t, "Unpaid")
	defer env.cleanupInvoice(t, r.ID)

	getResp := env.doAuthRequest(t, "GET", "/api/v1/invoices/"+r.ID, "")
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", getResp.StatusCode)
	}

	fetched := parseInvoice(t, getResp.Body)
	if fetched.ID != r.ID {
		t.Fatalf("expected ID '%s', got '%s'", r.ID, fetched.ID)
	}
	if fetched.ContractID != env.contractID {
		t.Errorf("expected contract_id '%s', got '%s'", env.contractID, fetched.ContractID)
	}
}

func TestInvoiceAPI_GetByID_NotFound(t *testing.T) {
	env := setupInvoiceTestEnv(t)

	resp := env.doAuthRequest(t, "GET", "/api/v1/invoices/nonexistent-id", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestInvoiceAPI_GetByID_CrossOrg(t *testing.T) {
	env := setupInvoiceTestEnv(t)

	r := env.createInvoice(t, "Unpaid")
	defer env.cleanupInvoice(t, r.ID)

	// A second organization must not be able to read this invoice.
	otherOrgID := createInvoiceTestOrg(t, env.server, env.token)

	resp := invoiceDoRequestWithOrg(t, env.server, "GET", "/api/v1/invoices/"+r.ID, "", env.token, otherOrgID)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 for cross-org access, got %d", resp.StatusCode)
	}
}

func TestInvoiceAPI_List_Pagination(t *testing.T) {
	env := setupInvoiceTestEnv(t)

	var ids []string
	for i := 0; i < 3; i++ {
		r := env.createInvoice(t, "Unpaid")
		ids = append(ids, r.ID)
	}
	defer func() {
		for _, id := range ids {
			env.cleanupInvoice(t, id)
		}
	}()

	resp := env.doAuthRequest(t, "GET", "/api/v1/invoices?per_page=2&page=1", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	list := parseInvoiceList(t, resp.Body)
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

func TestInvoiceAPI_List_Search(t *testing.T) {
	env := setupInvoiceTestEnv(t)

	r := env.createInvoice(t, "Paid")
	defer env.cleanupInvoice(t, r.ID)

	searchResp := env.doAuthRequest(t, "GET", "/api/v1/invoices?search=paid", "")
	defer searchResp.Body.Close()

	if searchResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", searchResp.StatusCode)
	}

	list := parseInvoiceList(t, searchResp.Body)
	found := false
	for _, item := range list.Data {
		if item.ID == r.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find created invoice in search results")
	}
}

func TestInvoiceAPI_List_CrossOrgIsolation(t *testing.T) {
	env := setupInvoiceTestEnv(t)

	r := env.createInvoice(t, "Unpaid")
	defer env.cleanupInvoice(t, r.ID)

	otherOrgID := createInvoiceTestOrg(t, env.server, env.token)
	resp := invoiceDoRequestWithOrg(t, env.server, "GET", "/api/v1/invoices", "", env.token, otherOrgID)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	list := parseInvoiceList(t, resp.Body)
	for _, item := range list.Data {
		if item.ID == r.ID {
			t.Error("cross-org list leaked invoice from another organization")
		}
	}
}

func TestInvoiceAPI_Update_Success(t *testing.T) {
	env := setupInvoiceTestEnv(t)

	r := env.createInvoice(t, "Unpaid")
	defer env.cleanupInvoice(t, r.ID)

	updateBody := fmt.Sprintf(`{"contract_id":"%s","tenant_id":"%s","amount":900000,"currency":"SGD","status":"Overdue","due_date":"2027-01-31T00:00:00Z","payment_method":"Cash","notes":"updated notes"}`, env.contractID, env.tenantID)
	updateResp := env.doAuthRequest(t, "PUT", "/api/v1/invoices/"+r.ID, updateBody)
	defer updateResp.Body.Close()

	if updateResp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(updateResp.Body)
		t.Fatalf("expected 200, got %d, body: %s", updateResp.StatusCode, string(b))
	}

	updated := parseInvoice(t, updateResp.Body)
	if updated.Amount != 900000 {
		t.Errorf("expected amount 900000, got %f", updated.Amount)
	}
	if updated.Currency != "SGD" {
		t.Errorf("expected currency 'SGD', got '%s'", updated.Currency)
	}
	if updated.Status != "Overdue" {
		t.Errorf("expected status 'Overdue', got '%s'", updated.Status)
	}
	if updated.Notes != "updated notes" {
		t.Errorf("expected notes 'updated notes', got '%s'", updated.Notes)
	}
}

func TestInvoiceAPI_Delete_Success(t *testing.T) {
	env := setupInvoiceTestEnv(t)

	r := env.createInvoice(t, "Unpaid")

	delResp := env.doAuthRequest(t, "DELETE", "/api/v1/invoices/"+r.ID, "")
	defer delResp.Body.Close()

	if delResp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", delResp.StatusCode)
	}

	getResp := env.doAuthRequest(t, "GET", "/api/v1/invoices/"+r.ID, "")
	defer getResp.Body.Close()
	if getResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", getResp.StatusCode)
	}
}

func TestInvoiceAPI_Unauthorized(t *testing.T) {
	env := setupInvoiceTestEnv(t)

	resp := invoiceDoRequest(t, env.server, "GET", "/api/v1/invoices", "", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}
