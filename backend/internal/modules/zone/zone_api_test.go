package zone_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/epmp/backend/internal/testutil"
)

func TestZoneAPI_CRUD(t *testing.T) {
	env := testutil.Setup(t)
	_, buildingID, _, _ := env.RoomChain(t)
	n := time.Now().UnixNano()

	create := fmt.Sprintf(`{"building_id":"%s","floor":1,"name":"Zone A %d"}`, buildingID, n)
	update := fmt.Sprintf(`{"building_id":"%s","floor":2,"name":"Zone B %d"}`, buildingID, n)

	testutil.AssertCRUD(t, env, "/api/v1/zones", create, update)
}

func TestZoneAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/zones")
}
