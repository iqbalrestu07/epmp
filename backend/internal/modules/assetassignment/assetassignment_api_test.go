package assetassignment_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/epmp/backend/internal/testutil"
)

func TestAssetAssignmentAPI_CRUD(t *testing.T) {
	env := testutil.Setup(t)
	propertyID, _, _, roomID := env.RoomChain(t)
	n := time.Now().UnixNano()
	assetID := env.Create(t, "/api/v1/assets",
		fmt.Sprintf(`{"property_id":"%s","name":"AA %d","category":"Furniture","status":"Available","purchase_price":100000}`, propertyID, n))
	now := time.Now().Format(time.RFC3339)

	create := fmt.Sprintf(`{"asset_id":"%s","room_id":"%s","assigned_date":"%s"}`, assetID, roomID, now)
	update := fmt.Sprintf(`{"asset_id":"%s","room_id":"%s","assigned_date":"%s"}`, assetID, roomID, now)

	testutil.AssertCRUD(t, env, "/api/v1/asset_assignments", create, update)
}

func TestAssetAssignmentAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/asset_assignments")
}
