package workorder_test

import (
	"fmt"
	"testing"

	"github.com/epmp/backend/internal/testutil"
)

func TestWorkOrderAPI_CRUD(t *testing.T) {
	env := testutil.Setup(t)
	propertyID, _, _, roomID := env.RoomChain(t)

	create := fmt.Sprintf(`{"property_id":"%s","room_id":"%s","description":"leak fix","status":"Open","priority":"High"}`, propertyID, roomID)
	update := fmt.Sprintf(`{"property_id":"%s","room_id":"%s","description":"leak fix done","status":"Done","priority":"Medium"}`, propertyID, roomID)

	testutil.AssertCRUD(t, env, "/api/v1/work_orders", create, update)
}

func TestWorkOrderAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/work_orders")
}
