// Package testutil provides a shared harness for API integration tests:
// it boots the full module stack against a real Postgres database and
// exposes small helpers to create orgs and fixture entities.
package testutil

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

// Env is a running test server + a fresh organization.
type Env struct {
	Server *httptest.Server
	DB     *pgxpool.Pool
	Token  string
	OrgID  string
}

type envelope struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
}

// Setup boots the full application on an httptest server with a fresh user+org.
// Skips the test when no database is reachable.
func Setup(t *testing.T) *Env {
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
	if err := modules.Register(e, db, log, jwtSecret, 2*time.Hour, 30*24*time.Hour); err != nil {
		t.Fatalf("failed to register modules: %v", err)
	}

	server := httptest.NewServer(e)
	env := &Env{Server: server, DB: db, Token: registerUser(t, server)}
	env.OrgID = env.CreateOrg(t)

	t.Cleanup(func() {
		server.Close()
		db.Close()
	})
	return env
}

// Do performs an HTTP request; pass empty token for unauthenticated calls.
func (e *Env) Do(t *testing.T, method, path, body, token, orgID string) *http.Response {
	t.Helper()
	var reqBody io.Reader
	if body != "" {
		reqBody = bytes.NewBufferString(body)
	}
	req, err := http.NewRequest(method, e.Server.URL+path, reqBody)
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

// Create posts a JSON body to path under the env org and returns the new id.
func (e *Env) Create(t *testing.T, path, body string) string {
	t.Helper()
	resp := e.Do(t, "POST", path, body, e.Token, e.OrgID)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("create %s failed: %d, body: %s", path, resp.StatusCode, string(b))
	}
	return decodeID(t, resp.Body)
}

// CreateOrg creates a fresh organization and returns its id.
func (e *Env) CreateOrg(t *testing.T) string {
	t.Helper()
	n := time.Now().UnixNano()
	body := fmt.Sprintf(`{"name":"Test Org %d","domain":"%d.test","is_active":true}`, n, n)
	resp := e.Do(t, "POST", "/api/v1/organizations", body, e.Token, "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("create org failed: %d, body: %s", resp.StatusCode, string(b))
	}
	return decodeID(t, resp.Body)
}

// List decodes a paginated list response and returns the raw items.
func (e *Env) List(t *testing.T, path, orgID string) []map[string]interface{} {
	t.Helper()
	resp := e.Do(t, "GET", path, "", e.Token, orgID)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list %s failed: %d", path, resp.StatusCode)
	}
	var env envelope
	var list struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if err := json.Unmarshal(env.Data, &list); err != nil {
		t.Fatalf("unmarshal list: %v", err)
	}
	return list.Data
}

// ListArray decodes a list response whose data is a raw JSON array
// (used by the communication module, which does not paginate).
func (e *Env) ListArray(t *testing.T, path, orgID string) []map[string]interface{} {
	t.Helper()
	resp := e.Do(t, "GET", path, "", e.Token, orgID)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list %s failed: %d", path, resp.StatusCode)
	}
	var env envelope
	var items []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if err := json.Unmarshal(env.Data, &items); err != nil {
		t.Fatalf("unmarshal list array: %v", err)
	}
	return items
}

// ─── Fixture chains ─────────────────────────────────────────────────────────

// RoomChain creates property → building → floor → room and returns all ids.
func (e *Env) RoomChain(t *testing.T) (propertyID, buildingID, floorID, roomID string) {
	t.Helper()
	n := time.Now().UnixNano()
	propertyID = e.Create(t, "/api/v1/properties",
		fmt.Sprintf(`{"name":"Prop %d","description":"t","address":"t","property_type":"boarding_house","is_active":true}`, n))
	buildingID = e.Create(t, "/api/v1/buildings",
		fmt.Sprintf(`{"property_id":"%s","name":"Bldg %d","total_floors":3}`, propertyID, n))
	floorID = e.Create(t, "/api/v1/floors",
		fmt.Sprintf(`{"building_id":"%s","name":"F1 %d","floor_number":1,"is_active":true}`, buildingID, n))
	roomID = e.Create(t, "/api/v1/rooms",
		fmt.Sprintf(`{"property_id":"%s","floor_id":"%s","name":"R-%d","capacity":2,"price":500000,"is_available":true}`, propertyID, floorID, n))
	return
}

// Tenant creates a tenant and returns its id.
func (e *Env) Tenant(t *testing.T) string {
	t.Helper()
	n := time.Now().UnixNano()
	return e.Create(t, "/api/v1/tenants",
		fmt.Sprintf(`{"full_name":"Tenant %d","email":"t%d@example.com","phone":"0812","identity_number":"KTP-%d","is_active":true}`, n, n, n))
}

// ContractChain creates tenant + room chain + contract; returns ids.
func (e *Env) ContractChain(t *testing.T) (contractID, tenantID, propertyID, roomID string) {
	t.Helper()
	tenantID = e.Tenant(t)
	var buildingID, floorID string
	propertyID, buildingID, floorID, roomID = e.RoomChain(t)
	_ = buildingID
	_ = floorID
	contractID = e.Create(t, "/api/v1/contracts",
		fmt.Sprintf(`{"tenant_id":"%s","property_id":"%s","room_id":"%s","status":"Active","start_date":"2026-01-01T00:00:00Z","end_date":"2026-12-31T00:00:00Z","monthly_rent":500000}`,
			tenantID, propertyID, roomID))
	return
}

// InvoiceChain creates contract chain + invoice; returns ids.
func (e *Env) InvoiceChain(t *testing.T) (invoiceID, tenantID, contractID string) {
	t.Helper()
	contractID, tenantID, _, _ = e.ContractChain(t)
	invoiceID = e.Create(t, "/api/v1/invoices",
		fmt.Sprintf(`{"contract_id":"%s","tenant_id":"%s","amount":500000,"currency":"IDR","status":"Unpaid","due_date":"%s"}`,
			contractID, tenantID, time.Now().AddDate(0, 0, 30).Format(time.RFC3339)))
	return
}

// PaymentChain creates invoice chain + payment; returns ids.
func (e *Env) PaymentChain(t *testing.T) (paymentID, invoiceID, tenantID string) {
	t.Helper()
	invoiceID, tenantID, _ = e.InvoiceChain(t)
	paymentID = e.Create(t, "/api/v1/payments",
		fmt.Sprintf(`{"invoice_id":"%s","tenant_id":"%s","amount":100000,"payment_date":"%s","payment_method":"Transfer","status":"Success","reference_number":"REF-%d"}`,
			invoiceID, tenantID, time.Now().Format(time.RFC3339), time.Now().UnixNano()))
	return
}

// AssertCRUD runs the standard CRUD + org-isolation assertions for a module:
// create → get → list → update → delete, and verifies that the entity is
// invisible to another organization.
func AssertCRUD(t *testing.T, e *Env, path, createBody, updateBody string) {
	t.Helper()

	id := e.Create(t, path, createBody)

	resp := e.Do(t, "GET", path+"/"+id, "", e.Token, e.OrgID)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("get %s/%s: expected 200, got %d", path, id, resp.StatusCode)
	}

	found := false
	for _, item := range e.List(t, path, e.OrgID) {
		if item["id"] == id {
			found = true
		}
	}
	if !found {
		t.Fatalf("list %s: created entity %s not found", path, id)
	}

	// Cross-org isolation: another org must not see or mutate this entity.
	otherOrg := e.CreateOrg(t)
	for _, method := range []string{"GET", "PUT", "DELETE"} {
		resp = e.Do(t, method, path+"/"+id, updateBody, e.Token, otherOrg)
		resp.Body.Close()
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusNoContent {
			t.Fatalf("cross-org %s %s/%s: expected 404, got %d", method, path, id, resp.StatusCode)
		}
	}

	resp = e.Do(t, "PUT", path+"/"+id, updateBody, e.Token, e.OrgID)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("update %s/%s: expected 200, got %d", path, id, resp.StatusCode)
	}

	resp = e.Do(t, "DELETE", path+"/"+id, "", e.Token, e.OrgID)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		t.Fatalf("delete %s/%s: expected 200/204, got %d", path, id, resp.StatusCode)
	}

	resp = e.Do(t, "GET", path+"/"+id, "", e.Token, e.OrgID)
	resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("get-after-delete %s/%s: expected 404, got %d", path, id, resp.StatusCode)
	}
}

// AssertUnauthorized checks that the endpoint rejects unauthenticated requests.
func AssertUnauthorized(t *testing.T, e *Env, path string) {
	t.Helper()
	resp := e.Do(t, "GET", path, "", "", "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unauthenticated GET %s: expected 401, got %d", path, resp.StatusCode)
	}
}

// ─── internals ──────────────────────────────────────────────────────────────

func registerUser(t *testing.T, server *httptest.Server) string {
	t.Helper()
	email := fmt.Sprintf("test_%d@epmp-test.com", time.Now().UnixNano())
	body := fmt.Sprintf(`{"email":"%s","password":"testpass123","name":"Test User"}`, email)

	resp, err := http.Post(server.URL+"/api/v1/auth/register", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		return decodeToken(t, resp.Body)
	}

	resp2, err := http.Post(server.URL+"/api/v1/auth/login", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("register/login failed: %d/%d", resp.StatusCode, resp2.StatusCode)
	}
	return decodeToken(t, resp2.Body)
}

func decodeToken(t *testing.T, body io.Reader) string {
	t.Helper()
	var env envelope
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

func decodeID(t *testing.T, body io.Reader) string {
	t.Helper()
	var env envelope
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
