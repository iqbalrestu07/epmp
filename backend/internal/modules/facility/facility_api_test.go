package facility_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/epmp/backend/internal/testutil"
)

func TestFacilityAPI_CRUD(t *testing.T) {
	env := testutil.Setup(t)
	propertyID, _, _, _ := env.RoomChain(t)
	n := time.Now().UnixNano()

	create := fmt.Sprintf(`{"property_id":"%s","name":"Pool %d","description":"outdoor pool"}`, propertyID, n)
	update := fmt.Sprintf(`{"property_id":"%s","name":"Gym %d","description":"fitness center"}`, propertyID, n)

	testutil.AssertCRUD(t, env, "/api/v1/facilities", create, update)
}

func TestFacilityAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/facilities")
}
