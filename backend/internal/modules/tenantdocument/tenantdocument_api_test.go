package tenantdocument_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/epmp/backend/internal/testutil"
)

func TestTenantDocumentAPI_CRUD(t *testing.T) {
	env := testutil.Setup(t)
	tenantID := env.Tenant(t)
	n := time.Now().UnixNano()

	create := fmt.Sprintf(`{"tenant_id":"%s","document_type":"KTP","file_url":"https://x/doc-%d.pdf"}`, tenantID, n)
	update := fmt.Sprintf(`{"tenant_id":"%s","document_type":"KK","file_url":"https://x/doc2-%d.pdf"}`, tenantID, n)

	testutil.AssertCRUD(t, env, "/api/v1/tenant_documents", create, update)
}

func TestTenantDocumentAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/tenant_documents")
}
