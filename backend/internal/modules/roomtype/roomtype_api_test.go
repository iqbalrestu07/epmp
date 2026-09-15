package roomtype_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/epmp/backend/internal/testutil"
)

func TestRoomTypeAPI_CRUD(t *testing.T) {
	env := testutil.Setup(t)
	n := time.Now().UnixNano()

	create := fmt.Sprintf(`{"name":"Standard %d","description":"std room","base_price":400000}`, n)
	update := fmt.Sprintf(`{"name":"Deluxe %d","description":"deluxe room","base_price":750000}`, n)

	testutil.AssertCRUD(t, env, "/api/v1/room_types", create, update)
}

func TestRoomTypeAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/room_types")
}
