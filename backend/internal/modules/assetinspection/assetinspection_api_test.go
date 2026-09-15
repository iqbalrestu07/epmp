package assetinspection_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/epmp/backend/internal/testutil"
)

func TestAssetInspectionAPI_CRUD(t *testing.T) {
	env := testutil.Setup(t)
	propertyID, _, _, _ := env.RoomChain(t)
	n := time.Now().UnixNano()
	assetID := env.Create(t, "/api/v1/assets",
		fmt.Sprintf(`{"property_id":"%s","name":"AI %d","category":"Furniture","status":"Available","purchase_price":100000}`, propertyID, n))
	now := time.Now().Format(time.RFC3339)

	create := fmt.Sprintf(`{"asset_id":"%s","inspection_date":"%s","condition":"Good","notes":"ok"}`, assetID, now)
	update := fmt.Sprintf(`{"asset_id":"%s","inspection_date":"%s","condition":"Fair","notes":"scratches"}`, assetID, now)

	testutil.AssertCRUD(t, env, "/api/v1/asset_inspections", create, update)
}

func TestAssetInspectionAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/asset_inspections")
}
