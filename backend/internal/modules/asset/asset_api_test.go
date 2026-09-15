package asset_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/epmp/backend/internal/testutil"
)

func TestAssetAPI_CRUD(t *testing.T) {
	env := testutil.Setup(t)
	propertyID, _, _, _ := env.RoomChain(t)
	n := time.Now().UnixNano()

	create := fmt.Sprintf(`{"property_id":"%s","name":"AC %d","category":"Electronics","status":"Available","purchase_price":2500000}`, propertyID, n)
	update := fmt.Sprintf(`{"property_id":"%s","name":"AC %d","category":"Electronics","status":"InUse","purchase_price":2500000}`, propertyID, n)

	testutil.AssertCRUD(t, env, "/api/v1/assets", create, update)
}

func TestAssetAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/assets")
}
