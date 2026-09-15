package bed_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/epmp/backend/internal/testutil"
)

func TestBedAPI_CRUD(t *testing.T) {
	env := testutil.Setup(t)
	_, _, _, roomID := env.RoomChain(t)
	n := time.Now().UnixNano()

	create := fmt.Sprintf(`{"room_id":"%s","name":"Bed A %d","status":"Available"}`, roomID, n)
	update := fmt.Sprintf(`{"room_id":"%s","name":"Bed B %d","status":"Occupied"}`, roomID, n)

	testutil.AssertCRUD(t, env, "/api/v1/beds", create, update)
}

func TestBedAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/beds")
}
