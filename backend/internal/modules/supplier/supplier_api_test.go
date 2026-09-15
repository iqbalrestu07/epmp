package supplier_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/epmp/backend/internal/testutil"
)

func TestSupplierAPI_CRUD(t *testing.T) {
	env := testutil.Setup(t)
	n := time.Now().UnixNano()

	create := fmt.Sprintf(`{"name":"Supplier %d","contact_person":"John","phone":"081%d","service_type":"Cleaning"}`, n, n%1000000)
	update := fmt.Sprintf(`{"name":"Supplier %d","contact_person":"Jane","phone":"082%d","service_type":"Laundry"}`, n, n%1000000)

	testutil.AssertCRUD(t, env, "/api/v1/vendors", create, update)
}

func TestSupplierAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/vendors")
}
