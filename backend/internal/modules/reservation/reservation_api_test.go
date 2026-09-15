package reservation_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/epmp/backend/internal/testutil"
)

func TestReservationAPI_CRUD(t *testing.T) {
	env := testutil.Setup(t)
	tenantID := env.Tenant(t)
	propertyID, _, _, roomID := env.RoomChain(t)
	in := time.Now().Format(time.RFC3339)
	out := time.Now().AddDate(0, 1, 0).Format(time.RFC3339)

	create := fmt.Sprintf(`{"tenant_id":"%s","property_id":"%s","room_id":"%s","status":"Pending","check_in_date":"%s","check_out_date":"%s","booking_fee":100000,"notes":"t"}`, tenantID, propertyID, roomID, in, out)
	update := fmt.Sprintf(`{"tenant_id":"%s","property_id":"%s","room_id":"%s","status":"Confirmed","check_in_date":"%s","check_out_date":"%s","booking_fee":100000,"notes":"u"}`, tenantID, propertyID, roomID, in, out)

	testutil.AssertCRUD(t, env, "/api/v1/reservations", create, update)
}

func TestReservationAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/reservations")
}
