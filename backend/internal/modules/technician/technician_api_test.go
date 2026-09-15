package technician_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/epmp/backend/internal/testutil"
)

func TestTechnicianAPI_CRUD(t *testing.T) {
	env := testutil.Setup(t)
	n := time.Now().UnixNano()

	create := fmt.Sprintf(`{"name":"Tech %d","phone":"081%d","specialty":"Plumbing"}`, n, n%1000000)
	update := fmt.Sprintf(`{"name":"Tech %d","phone":"082%d","specialty":"Electrical"}`, n, n%1000000)

	testutil.AssertCRUD(t, env, "/api/v1/technicians", create, update)
}

func TestTechnicianAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/technicians")
}
