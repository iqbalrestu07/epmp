package charge_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/epmp/backend/internal/testutil"
)

func TestChargeAPI_CRUD(t *testing.T) {
	env := testutil.Setup(t)
	invoiceID, tenantID, contractID := env.InvoiceChain(t)
	_ = tenantID
	now := time.Now().Format(time.RFC3339)

	create := fmt.Sprintf(`{"contract_id":"%s","invoice_id":"%s","charge_type":"Utility","amount":150000,"status":"Unbilled","charge_date":"%s","notes":"t"}`, contractID, invoiceID, now)
	update := fmt.Sprintf(`{"contract_id":"%s","invoice_id":"%s","charge_type":"Utility","amount":175000,"status":"Billed","charge_date":"%s","notes":"u"}`, contractID, invoiceID, now)

	testutil.AssertCRUD(t, env, "/api/v1/charges", create, update)
}

func TestChargeAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/charges")
}
