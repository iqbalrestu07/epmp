package report_test

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

type reportTestEnv struct {
	server     *httptest.Server
	db         *pgxpool.Pool
	token      string
	orgID      string
	tenantID   string
	contractID string
	invoiceID  string
}

type reportEnvelope struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type occupancyReport struct {
	Totals struct {
		TotalRooms    int64   `json:"total_rooms"`
		Occupied      int64   `json:"occupied"`
		Available     int64   `json:"available"`
		Reserved      int64   `json:"reserved"`
		OccupancyRate float64 `json:"occupancy_rate"`
	} `json:"totals"`
	Properties []struct {
		PropertyID   string `json:"property_id"`
		PropertyName string `json:"property_name"`
		TotalRooms   int64  `json:"total_rooms"`
		Available    int64  `json:"available"`
		Buildings    []struct {
			BuildingID string `json:"building_id"`
			TotalRooms int64  `json:"total_rooms"`
		} `json:"buildings"`
	} `json:"properties"`
}

type revenueReport struct {
	From string `json:"from"`
	To   string `json:"to"`
	Rows []struct {
		Period         string  `json:"period"`
		Currency       string  `json:"currency"`
		InvoiceCount   int64   `json:"invoice_count"`
		InvoicedAmount float64 `json:"invoiced_amount"`
		ReceivedAmount float64 `json:"received_amount"`
	} `json:"rows"`
}

type arAgingReport struct {
	Buckets []struct {
		Bucket       string  `json:"bucket"`
		Currency     string  `json:"currency"`
		InvoiceCount int64   `json:"invoice_count"`
		Outstanding  float64 `json:"outstanding"`
	} `json:"buckets"`
	Invoices []struct {
		InvoiceID   string  `json:"invoice_id"`
		TenantName  string  `json:"tenant_name"`
		Outstanding float64 `json:"outstanding"`
		Bucket      string  `json:"bucket"`
	} `json:"invoices"`
}

func setupReportTestEnv(t *testing.T) *reportTestEnv {
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
	token := reportTestToken(t, server)
	orgID := reportCreateOrg(t, server, token)
	propertyID := reportCreate(t, server, token, orgID, "/api/v1/properties",
		fmt.Sprintf(`{"name":"Report Prop %d","description":"t","address":"t","property_type":"boarding_house","is_active":true}`, time.Now().UnixNano()))
	buildingID := reportCreate(t, server, token, orgID, "/api/v1/buildings",
		fmt.Sprintf(`{"property_id":"%s","name":"Report Bldg %d","total_floors":3}`, propertyID, time.Now().UnixNano()))
	floorID := reportCreate(t, server, token, orgID, "/api/v1/floors",
		fmt.Sprintf(`{"building_id":"%s","name":"F1 %d","floor_number":1,"is_active":true}`, buildingID, time.Now().UnixNano()))
	reportCreate(t, server, token, orgID, "/api/v1/rooms",
		fmt.Sprintf(`{"property_id":"%s","floor_id":"%s","name":"R-%d","capacity":2,"price":500000,"is_available":true}`, propertyID, floorID, time.Now().UnixNano()))
	tenantID := reportCreate(t, server, token, orgID, "/api/v1/tenants",
		fmt.Sprintf(`{"full_name":"Report Tenant %d","email":"rpt%d@example.com","phone":"0812","identity_number":"KTP-rpt-%d","is_active":true}`,
			time.Now().UnixNano(), time.Now().UnixNano(), time.Now().UnixNano()))
	contractID := reportCreate(t, server, token, orgID, "/api/v1/contracts",
		fmt.Sprintf(`{"tenant_id":"%s","property_id":"%s","room_id":"%s","status":"Active","start_date":"2026-01-01T00:00:00Z","end_date":"2026-12-31T00:00:00Z","monthly_rent":500000}`,
			tenantID, propertyID, reportLastRoomID(t, server, token, orgID, propertyID)))
	invoiceID := reportCreate(t, server, token, orgID, "/api/v1/invoices",
		fmt.Sprintf(`{"contract_id":"%s","tenant_id":"%s","amount":500000,"currency":"IDR","status":"Unpaid","due_date":"%s"}`, contractID, tenantID, time.Now().AddDate(0, 0, 30).Format(time.RFC3339)))

	env := &reportTestEnv{server: server, db: db, token: token, orgID: orgID, tenantID: tenantID, contractID: contractID, invoiceID: invoiceID}

	t.Cleanup(func() {
		server.Close()
		db.Close()
	})
	return env
}

// reportLastRoomID returns the most recently created room for a property.
func reportLastRoomID(t *testing.T, server *httptest.Server, token, orgID, propertyID string) string {
	t.Helper()
	resp := reportDo(t, server, "GET", "/api/v1/rooms?property_id="+propertyID, "", token, orgID)
	defer resp.Body.Close()
	var env reportEnvelope
	var list struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		t.Fatalf("decode rooms: %v", err)
	}
	if err := json.Unmarshal(env.Data, &list); err != nil || len(list.Data) == 0 {
		t.Fatalf("no rooms found for property %s", propertyID)
	}
	return list.Data[0].ID
}

func reportTestToken(t *testing.T, server *httptest.Server) string {
	t.Helper()
	email := fmt.Sprintf("test_report_%d@epmp-test.com", time.Now().UnixNano())
	body := fmt.Sprintf(`{"email":"%s","password":"testpass123","name":"Test User"}`, email)
	resp := reportDo(t, server, "POST", "/api/v1/auth/register", body, "", "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		resp2 := reportDo(t, server, "POST", "/api/v1/auth/login", body, "", "")
		defer resp2.Body.Close()
		if resp2.StatusCode != http.StatusOK {
			t.Fatalf("register/login failed: %d/%d", resp.StatusCode, resp2.StatusCode)
		}
		return reportToken(t, resp2.Body)
	}
	return reportToken(t, resp.Body)
}

func reportToken(t *testing.T, body io.Reader) string {
	t.Helper()
	var env reportEnvelope
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

func reportCreateOrg(t *testing.T, server *httptest.Server, token string) string {
	t.Helper()
	body := fmt.Sprintf(`{"name":"Report Org %d","domain":"%d.report.test","is_active":true}`, time.Now().UnixNano(), time.Now().UnixNano())
	resp := reportDo(t, server, "POST", "/api/v1/organizations", body, token, "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create org failed: %d", resp.StatusCode)
	}
	return reportID(t, resp.Body)
}

// reportCreate posts a payload and returns the created entity id.
func reportCreate(t *testing.T, server *httptest.Server, token, orgID, path, body string) string {
	t.Helper()
	resp := reportDo(t, server, "POST", path, body, token, orgID)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("create %s failed: %d, body: %s", path, resp.StatusCode, string(b))
	}
	return reportID(t, resp.Body)
}

func reportID(t *testing.T, body io.Reader) string {
	t.Helper()
	var env reportEnvelope
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

func reportDo(t *testing.T, server *httptest.Server, method, path, body, token, orgID string) *http.Response {
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

// ─── Tests ───────────────────────────────────────────────────────────────────

func TestReportAPI_Occupancy(t *testing.T) {
	env := setupReportTestEnv(t)

	resp := reportDo(t, env.server, "GET", "/api/v1/reports/occupancy", "", env.token, env.orgID)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var env2 reportEnvelope
	var rep occupancyReport
	if err := json.NewDecoder(resp.Body).Decode(&env2); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if err := json.Unmarshal(env2.Data, &rep); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if rep.Totals.TotalRooms < 1 {
		t.Errorf("expected total_rooms>=1, got %d", rep.Totals.TotalRooms)
	}
	// The room is linked to an Active contract, so the sync trigger marks it Reserved.
	if rep.Totals.Reserved+rep.Totals.Occupied+rep.Totals.Available < 1 {
		t.Errorf("expected room counted in some status bucket, got %+v", rep.Totals)
	}
	if len(rep.Properties) == 0 {
		t.Error("expected at least one property row")
	}
}

func TestReportAPI_Occupancy_CrossOrg(t *testing.T) {
	env := setupReportTestEnv(t)

	otherOrg := reportCreateOrg(t, env.server, env.token)
	resp := reportDo(t, env.server, "GET", "/api/v1/reports/occupancy", "", env.token, otherOrg)
	defer resp.Body.Close()

	var env2 reportEnvelope
	var rep occupancyReport
	if err := json.NewDecoder(resp.Body).Decode(&env2); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if err := json.Unmarshal(env2.Data, &rep); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if rep.Totals.TotalRooms != 0 {
		t.Errorf("cross-org report leaked rooms: total=%d", rep.Totals.TotalRooms)
	}
}

func TestReportAPI_Revenue(t *testing.T) {
	env := setupReportTestEnv(t)

	from := time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	to := time.Now().AddDate(0, 2, 0).Format("2006-01-02")
	resp := reportDo(t, env.server, "GET", fmt.Sprintf("/api/v1/reports/revenue?from=%s&to=%s", from, to), "", env.token, env.orgID)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var env2 reportEnvelope
	var rep revenueReport
	if err := json.NewDecoder(resp.Body).Decode(&env2); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if err := json.Unmarshal(env2.Data, &rep); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	found := false
	for _, r := range rep.Rows {
		if r.InvoiceCount >= 1 && r.InvoicedAmount >= 500000 {
			found = true
		}
	}
	if !found {
		t.Error("expected revenue row containing the created invoice")
	}
}

func TestReportAPI_Revenue_InvalidRange(t *testing.T) {
	env := setupReportTestEnv(t)

	resp := reportDo(t, env.server, "GET", "/api/v1/reports/revenue?from=2026-12-01&to=2026-01-01", "", env.token, env.orgID)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestReportAPI_ArAging(t *testing.T) {
	env := setupReportTestEnv(t)

	resp := reportDo(t, env.server, "GET", "/api/v1/reports/ar-aging", "", env.token, env.orgID)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var env2 reportEnvelope
	var rep arAgingReport
	if err := json.NewDecoder(resp.Body).Decode(&env2); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if err := json.Unmarshal(env2.Data, &rep); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	found := false
	for _, inv := range rep.Invoices {
		if inv.InvoiceID == env.invoiceID && inv.Outstanding == 500000 {
			found = true
		}
	}
	if !found {
		t.Error("expected unpaid invoice in ar-aging report")
	}
}

func TestReportAPI_Unauthorized(t *testing.T) {
	env := setupReportTestEnv(t)

	resp := reportDo(t, env.server, "GET", "/api/v1/reports/occupancy", "", "", "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestReportAPI_MissingOrgHeader(t *testing.T) {
	env := setupReportTestEnv(t)

	resp := reportDo(t, env.server, "GET", "/api/v1/reports/occupancy", "", env.token, "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}
