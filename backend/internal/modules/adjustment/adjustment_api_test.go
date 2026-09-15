package adjustment_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/epmp/backend/internal/testutil"
)

func TestAdjustmentAPI_CRUD(t *testing.T) {
	env := testutil.Setup(t)
	invoiceID, _, _ := env.InvoiceChain(t)
	now := time.Now().Format(time.RFC3339)

	create := fmt.Sprintf(`{"invoice_id":"%s","adjustment_type":"Discount","amount":50000,"adjustment_date":"%s","reason":"promo"}`, invoiceID, now)
	update := fmt.Sprintf(`{"invoice_id":"%s","adjustment_type":"Correction","amount":75000,"adjustment_date":"%s","reason":"price fix"}`, invoiceID, now)

	testutil.AssertCRUD(t, env, "/api/v1/adjustments", create, update)
}

func TestAdjustmentAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/adjustments")
}
