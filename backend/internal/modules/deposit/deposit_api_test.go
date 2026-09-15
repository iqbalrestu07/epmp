package deposit_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/epmp/backend/internal/testutil"
)

func TestDepositAPI_CRUD(t *testing.T) {
	env := testutil.Setup(t)
	contractID, tenantID, _, _ := env.ContractChain(t)
	now := time.Now().Format(time.RFC3339)

	create := fmt.Sprintf(`{"contract_id":"%s","tenant_id":"%s","amount":500000,"status":"Collected","collection_date":"%s","notes":"t"}`, contractID, tenantID, now)
	update := fmt.Sprintf(`{"contract_id":"%s","tenant_id":"%s","amount":450000,"status":"Refunded","collection_date":"%s","notes":"u"}`, contractID, tenantID, now)

	testutil.AssertCRUD(t, env, "/api/v1/deposits", create, update)
}

func TestDepositAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/deposits")
}
