package tenantcontact_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/epmp/backend/internal/testutil"
)

func TestTenantContactAPI_CRUD(t *testing.T) {
	env := testutil.Setup(t)
	tenantID := env.Tenant(t)
	n := time.Now().UnixNano()

	create := fmt.Sprintf(`{"tenant_id":"%s","contact_type":"phone","contact_value":"081%d","is_primary":true}`, tenantID, n%1000000)
	update := fmt.Sprintf(`{"tenant_id":"%s","contact_type":"email","contact_value":"c%d@example.com","is_primary":false}`, tenantID, n)

	testutil.AssertCRUD(t, env, "/api/v1/tenant_contacts", create, update)
}

func TestTenantContactAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/tenant_contacts")
}
