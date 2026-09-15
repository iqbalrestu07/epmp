package tenantidentity_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/epmp/backend/internal/testutil"
)

func TestTenantIdentityAPI_CRUD(t *testing.T) {
	env := testutil.Setup(t)
	tenantID := env.Tenant(t)
	n := time.Now().UnixNano()

	create := fmt.Sprintf(`{"tenant_id":"%s","identity_type":"KTP","identity_number":"ID-%d","file_url":"https://x/ktp.pdf"}`, tenantID, n)
	update := fmt.Sprintf(`{"tenant_id":"%s","identity_type":"SIM","identity_number":"ID2-%d","file_url":"https://x/sim.pdf"}`, tenantID, n)

	testutil.AssertCRUD(t, env, "/api/v1/tenant_identities", create, update)
}

func TestTenantIdentityAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/tenant_identities")
}
