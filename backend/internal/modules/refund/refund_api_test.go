package refund_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/epmp/backend/internal/testutil"
)

func TestRefundAPI_CRUD(t *testing.T) {
	env := testutil.Setup(t)
	paymentID, _, tenantID := env.PaymentChain(t)
	now := time.Now().Format(time.RFC3339)

	create := fmt.Sprintf(`{"payment_id":"%s","tenant_id":"%s","amount":50000,"status":"Pending","refund_date":"%s","reason":"overpayment"}`, paymentID, tenantID, now)
	update := fmt.Sprintf(`{"payment_id":"%s","tenant_id":"%s","amount":50000,"status":"Approved","refund_date":"%s","reason":"overpayment confirmed"}`, paymentID, tenantID, now)

	testutil.AssertCRUD(t, env, "/api/v1/refunds", create, update)
}

func TestRefundAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/refunds")
}
