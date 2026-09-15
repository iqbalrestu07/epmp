package communication_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/epmp/backend/internal/testutil"
)

func TestCommunicationAPI_DeviceCRUD(t *testing.T) {
	env := testutil.Setup(t)
	n := time.Now().UnixNano()

	// Create
	deviceID := env.Create(t, "/api/v1/communication/devices", fmt.Sprintf(`{"label":"Dev %d"}`, n))

	// List contains the device and is org-scoped
	devices := env.ListArray(t, "/api/v1/communication/devices", env.OrgID)
	found := false
	for _, d := range devices {
		if d["id"] == deviceID && d["org_id"] == env.OrgID {
			found = true
		}
	}
	if !found {
		t.Fatal("created device not found in org device list")
	}

	// Cross-org: list and delete must not see it
	otherOrg := env.CreateOrg(t)
	for _, d := range env.ListArray(t, "/api/v1/communication/devices", otherOrg) {
		if d["id"] == deviceID {
			t.Fatal("cross-org device list leaked device")
		}
	}
	resp := env.Do(t, "DELETE", "/api/v1/communication/devices/"+deviceID, "", env.Token, otherOrg)
	resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Fatal("cross-org device delete succeeded")
	}

	// Update status + delete in own org
	resp = env.Do(t, "PUT", "/api/v1/communication/devices/"+deviceID+"/status", `{"status":"disconnected"}`, env.Token, env.OrgID)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("update status: %d", resp.StatusCode)
	}
	resp = env.Do(t, "DELETE", "/api/v1/communication/devices/"+deviceID, "", env.Token, env.OrgID)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("delete: %d", resp.StatusCode)
	}
}

func TestCommunicationAPI_TemplateCRUD(t *testing.T) {
	env := testutil.Setup(t)
	n := time.Now().UnixNano()

	tmplID := env.Create(t, "/api/v1/communication/templates",
		fmt.Sprintf(`{"name":"Tmpl %d","content":"Halo {{tenant_name}}","category":"general","variables":["tenant_name"]}`, n))

	templates := env.ListArray(t, "/api/v1/communication/templates", env.OrgID)
	found := false
	for _, x := range templates {
		if x["id"] == tmplID && x["org_id"] == env.OrgID {
			found = true
		}
	}
	if !found {
		t.Fatal("created template not found in org template list")
	}

	otherOrg := env.CreateOrg(t)
	for _, x := range env.ListArray(t, "/api/v1/communication/templates", otherOrg) {
		if x["id"] == tmplID {
			t.Fatal("cross-org template list leaked template")
		}
	}
	resp := env.Do(t, "DELETE", "/api/v1/communication/templates/"+tmplID, "", env.Token, otherOrg)
	resp.Body.Close()
	// Delete is org-scoped soft-delete; a cross-org delete is silently ignored —
	// verify the template is still present in the owner's list.
	for _, x := range env.ListArray(t, "/api/v1/communication/templates", env.OrgID) {
		if x["id"] == tmplID {
			return // still present — pass
		}
	}
	t.Fatal("template was deleted by a cross-org request")
}

func TestCommunicationAPI_BlastOrgScoped(t *testing.T) {
	env := testutil.Setup(t)

	blastID := env.Create(t, "/api/v1/communication/blast",
		`{"title":"Promo","template":"Halo {{tenant_name}}","target_type":"all_tenants"}`)

	blasts := env.ListArray(t, "/api/v1/communication/blast", env.OrgID)
	found := false
	for _, b := range blasts {
		if b["id"] == blastID && b["org_id"] == env.OrgID {
			found = true
		}
	}
	if !found {
		t.Fatal("created blast not found in org blast list")
	}

	otherOrg := env.CreateOrg(t)
	resp := env.Do(t, "GET", "/api/v1/communication/blast/"+blastID, "", env.Token, otherOrg)
	resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Fatal("cross-org blast get succeeded")
	}
}

func TestCommunicationAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/communication/devices")
}
