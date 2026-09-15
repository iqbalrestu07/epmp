package occupancy_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/epmp/backend/internal/testutil"
)

func TestOccupancyAPI_CRUD(t *testing.T) {
	env := testutil.Setup(t)
	contractID, tenantID, _, roomID := env.ContractChain(t)
	now := time.Now().Format(time.RFC3339)
	out := time.Now().AddDate(0, 1, 0).Format(time.RFC3339)

	create := fmt.Sprintf(`{"contract_id":"%s","room_id":"%s","tenant_id":"%s","status":"CheckedIn","check_in_time":"%s","check_out_time":"%s","notes":"t"}`, contractID, roomID, tenantID, now, out)
	update := fmt.Sprintf(`{"contract_id":"%s","room_id":"%s","tenant_id":"%s","status":"CheckedOut","check_in_time":"%s","check_out_time":"%s","notes":"u"}`, contractID, roomID, tenantID, now, out)

	testutil.AssertCRUD(t, env, "/api/v1/occupancies", create, update)
}

func TestOccupancyAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/occupancies")
}
