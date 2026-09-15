package dashboard_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/epmp/backend/internal/testutil"
)

func TestDashboardAPI_Summary(t *testing.T) {
	env := testutil.Setup(t)

	// Seed a room so the summary has data.
	env.RoomChain(t)

	resp := env.Do(t, "GET", "/api/v1/dashboard/summary", "", env.Token, env.OrgID)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var envJson struct {
		Success bool            `json:"success"`
		Data    json.RawMessage `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envJson); err != nil {
		t.Fatalf("decode: %v", err)
	}
	var summary struct {
		Metrics struct {
			TotalProperties int64 `json:"total_properties"`
			TotalRooms      int64 `json:"total_rooms"`
			ActiveTenants   int64 `json:"active_tenants"`
		} `json:"metrics"`
	}
	if err := json.Unmarshal(envJson.Data, &summary); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if summary.Metrics.TotalProperties < 1 || summary.Metrics.TotalRooms < 1 {
		t.Errorf("expected seeded data in summary, got %+v", summary)
	}
}

func TestDashboardAPI_Summary_CrossOrg(t *testing.T) {
	env := testutil.Setup(t)
	env.RoomChain(t)

	otherOrg := env.CreateOrg(t)
	resp := env.Do(t, "GET", "/api/v1/dashboard/summary", "", env.Token, otherOrg)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var envJson struct {
		Data json.RawMessage `json:"data"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&envJson)
	var summary struct {
		Metrics struct {
			TotalRooms int64 `json:"total_rooms"`
		} `json:"metrics"`
	}
	_ = json.Unmarshal(envJson.Data, &summary)
	if summary.Metrics.TotalRooms != 0 {
		t.Errorf("cross-org summary leaked rooms: %d", summary.Metrics.TotalRooms)
	}
}

func TestDashboardAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/dashboard/summary")
}
