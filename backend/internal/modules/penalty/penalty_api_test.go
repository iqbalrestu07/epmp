package penalty_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/epmp/backend/internal/testutil"
)

func TestPenaltyAPI_CRUD(t *testing.T) {
	env := testutil.Setup(t)
	invoiceID, _, _ := env.InvoiceChain(t)
	now := time.Now().Format(time.RFC3339)

	create := fmt.Sprintf(`{"invoice_id":"%s","amount":25000,"status":"Unpaid","penalty_date":"%s","description":"late fee"}`, invoiceID, now)
	update := fmt.Sprintf(`{"invoice_id":"%s","amount":30000,"status":"Paid","penalty_date":"%s","description":"late fee updated"}`, invoiceID, now)

	testutil.AssertCRUD(t, env, "/api/v1/penalties", create, update)
}

func TestPenaltyAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/penalties")
}
