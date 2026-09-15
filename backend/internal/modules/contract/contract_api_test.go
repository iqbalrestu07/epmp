package contract_test

import (
	"fmt"
	"testing"

	"github.com/epmp/backend/internal/testutil"
)

func TestContractAPI_CRUD(t *testing.T) {
	env := testutil.Setup(t)
	tenantID := env.Tenant(t)
	propertyID, _, _, roomID := env.RoomChain(t)

	create := fmt.Sprintf(`{"tenant_id":"%s","property_id":"%s","room_id":"%s","status":"Active","start_date":"2026-01-01T00:00:00Z","end_date":"2026-12-31T00:00:00Z","monthly_rent":500000,"deposit_amount":500000,"terms":"std"}`, tenantID, propertyID, roomID)
	update := fmt.Sprintf(`{"tenant_id":"%s","property_id":"%s","room_id":"%s","status":"Terminated","start_date":"2026-01-01T00:00:00Z","end_date":"2026-12-31T00:00:00Z","monthly_rent":550000,"deposit_amount":500000,"terms":"updated"}`, tenantID, propertyID, roomID)

	testutil.AssertCRUD(t, env, "/api/v1/contracts", create, update)
}

func TestContractAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/contracts")
}
