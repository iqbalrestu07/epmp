package main

import (
	"fmt"
	"os"
	"strings"
)

type Entity struct {
	Name           string
	Package        string
	Table          string
	BoundedContext string
	Fields         string
	SQLFields      string
}

func main() {
	entities := []Entity{
		{
			Name:           "Building",
			Package:        "building",
			Table:          "buildings",
			BoundedContext: "property",
			Fields: `
    - name: id
      type: uuid
      primary_key: true
      auto: true
    - name: property_id
      type: uuid
      searchable: true
    - name: name
      type: string
      max_length: 100
      searchable: true
    - name: total_floors
      type: int
      default: 1`,
			SQLFields: `
    property_id UUID NOT NULL REFERENCES properties(id),
    name VARCHAR(100) NOT NULL,
    total_floors INT NOT NULL DEFAULT 1`,
		},
		{
			Name:           "Zone",
			Package:        "zone",
			Table:          "zones",
			BoundedContext: "property",
			Fields: `
    - name: id
      type: uuid
      primary_key: true
      auto: true
    - name: building_id
      type: uuid
      searchable: true
    - name: floor
      type: int
    - name: name
      type: string
      max_length: 100
      searchable: true`,
			SQLFields: `
    building_id UUID NOT NULL REFERENCES buildings(id),
    floor INT NOT NULL,
    name VARCHAR(100) NOT NULL`,
		},
		{
			Name:           "Bed",
			Package:        "bed",
			Table:          "beds",
			BoundedContext: "property",
			Fields: `
    - name: id
      type: uuid
      primary_key: true
      auto: true
    - name: room_id
      type: uuid
      searchable: true
    - name: name
      type: string
      max_length: 100
      searchable: true
    - name: status
      type: enum
      enum_values:
        - Available
        - Occupied
      default: "Available"`,
			SQLFields: `
    room_id UUID NOT NULL REFERENCES rooms(id),
    name VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'Available'`,
		},
		{
			Name:           "Facility",
			Package:        "facility",
			Table:          "facilities",
			BoundedContext: "property",
			Fields: `
    - name: id
      type: uuid
      primary_key: true
      auto: true
    - name: property_id
      type: uuid
      searchable: true
    - name: name
      type: string
      max_length: 100
      searchable: true
    - name: description
      type: text
      nullable: true`,
			SQLFields: `
    property_id UUID NOT NULL REFERENCES properties(id),
    name VARCHAR(100) NOT NULL,
    description TEXT`,
		},
		{
			Name:           "RoomType",
			Package:        "roomtype",
			Table:          "room_types",
			BoundedContext: "property",
			Fields: `
    - name: id
      type: uuid
      primary_key: true
      auto: true
    - name: name
      type: string
      max_length: 100
      searchable: true
    - name: description
      type: text
      nullable: true
    - name: base_price
      type: decimal
      default: 0`,
			SQLFields: `
    name VARCHAR(100) NOT NULL,
    description TEXT,
    base_price DECIMAL(19, 4) NOT NULL DEFAULT 0`,
		},
		{
			Name:           "TenantIdentity",
			Package:        "tenantidentity",
			Table:          "tenant_identities",
			BoundedContext: "tenant",
			Fields: `
    - name: id
      type: uuid
      primary_key: true
      auto: true
    - name: tenant_id
      type: uuid
      searchable: true
    - name: identity_type
      type: enum
      enum_values:
        - KTP
        - Passport
      default: "KTP"
    - name: identity_number
      type: string
      max_length: 100
      searchable: true
    - name: file_url
      type: string
      nullable: true`,
			SQLFields: `
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    identity_type VARCHAR(50) NOT NULL DEFAULT 'KTP',
    identity_number VARCHAR(100) NOT NULL,
    file_url VARCHAR(255)`,
		},
		{
			Name:           "TenantContact",
			Package:        "tenantcontact",
			Table:          "tenant_contacts",
			BoundedContext: "tenant",
			Fields: `
    - name: id
      type: uuid
      primary_key: true
      auto: true
    - name: tenant_id
      type: uuid
      searchable: true
    - name: contact_type
      type: enum
      enum_values:
        - Email
        - Phone
        - Emergency
      default: "Phone"
    - name: contact_value
      type: string
      max_length: 100
      searchable: true
    - name: is_primary
      type: bool
      default: false`,
			SQLFields: `
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    contact_type VARCHAR(50) NOT NULL DEFAULT 'Phone',
    contact_value VARCHAR(100) NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT false`,
		},
		{
			Name:           "TenantDocument",
			Package:        "tenantdocument",
			Table:          "tenant_documents",
			BoundedContext: "tenant",
			Fields: `
    - name: id
      type: uuid
      primary_key: true
      auto: true
    - name: tenant_id
      type: uuid
      searchable: true
    - name: document_type
      type: string
      max_length: 100
    - name: file_url
      type: string
      max_length: 255`,
			SQLFields: `
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    document_type VARCHAR(100) NOT NULL,
    file_url VARCHAR(255) NOT NULL`,
		},
		{
			Name:           "Asset",
			Package:        "asset",
			Table:          "assets",
			BoundedContext: "asset",
			Fields: `
    - name: id
      type: uuid
      primary_key: true
      auto: true
    - name: property_id
      type: uuid
      searchable: true
    - name: name
      type: string
      max_length: 100
      searchable: true
    - name: category
      type: string
      max_length: 100
    - name: status
      type: enum
      enum_values:
        - Available
        - Assigned
        - Maintenance
        - Disposed
      default: "Available"
    - name: purchase_price
      type: decimal
      default: 0`,
			SQLFields: `
    property_id UUID NOT NULL REFERENCES properties(id),
    name VARCHAR(100) NOT NULL,
    category VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'Available',
    purchase_price DECIMAL(19, 4) NOT NULL DEFAULT 0`,
		},
		{
			Name:           "AssetAssignment",
			Package:        "assetassignment",
			Table:          "asset_assignments",
			BoundedContext: "asset",
			Fields: `
    - name: id
      type: uuid
      primary_key: true
      auto: true
    - name: asset_id
      type: uuid
      searchable: true
    - name: room_id
      type: uuid
      searchable: true
    - name: assigned_date
      type: timestamp`,
			SQLFields: `
    asset_id UUID NOT NULL REFERENCES assets(id),
    room_id UUID NOT NULL REFERENCES rooms(id),
    assigned_date TIMESTAMP WITH TIME ZONE NOT NULL`,
		},
		{
			Name:           "AssetInspection",
			Package:        "assetinspection",
			Table:          "asset_inspections",
			BoundedContext: "asset",
			Fields: `
    - name: id
      type: uuid
      primary_key: true
      auto: true
    - name: asset_id
      type: uuid
      searchable: true
    - name: inspection_date
      type: timestamp
    - name: condition
      type: string
      max_length: 100
    - name: notes
      type: text
      nullable: true`,
			SQLFields: `
    asset_id UUID NOT NULL REFERENCES assets(id),
    inspection_date TIMESTAMP WITH TIME ZONE NOT NULL,
    condition VARCHAR(100) NOT NULL,
    notes TEXT`,
		},
		{
			Name:           "WorkOrder",
			Package:        "workorder",
			Table:          "work_orders",
			BoundedContext: "maintenance",
			Fields: `
    - name: id
      type: uuid
      primary_key: true
      auto: true
    - name: property_id
      type: uuid
      searchable: true
    - name: room_id
      type: uuid
      nullable: true
      searchable: true
    - name: description
      type: text
    - name: status
      type: enum
      enum_values:
        - Open
        - InProgress
        - Resolved
      default: "Open"
    - name: priority
      type: string
      max_length: 50
      default: "Medium"`,
			SQLFields: `
    property_id UUID NOT NULL REFERENCES properties(id),
    room_id UUID REFERENCES rooms(id),
    description TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'Open',
    priority VARCHAR(50) NOT NULL DEFAULT 'Medium'`,
		},
		{
			Name:           "Technician",
			Package:        "technician",
			Table:          "technicians",
			BoundedContext: "maintenance",
			Fields: `
    - name: id
      type: uuid
      primary_key: true
      auto: true
    - name: name
      type: string
      max_length: 100
      searchable: true
    - name: phone
      type: string
      max_length: 50
    - name: specialty
      type: string
      max_length: 100`,
			SQLFields: `
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    specialty VARCHAR(100) NOT NULL`,
		},
		{
			Name:           "Vendor",
			Package:        "vendor",
			Table:          "vendors",
			BoundedContext: "maintenance",
			Fields: `
    - name: id
      type: uuid
      primary_key: true
      auto: true
    - name: name
      type: string
      max_length: 100
      searchable: true
    - name: contact_person
      type: string
      max_length: 100
    - name: phone
      type: string
      max_length: 50
    - name: service_type
      type: string
      max_length: 100`,
			SQLFields: `
    name VARCHAR(100) NOT NULL,
    contact_person VARCHAR(100) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    service_type VARCHAR(100) NOT NULL`,
		},
	}

	yamlTemplate := `version: "1.0"
dry_run: false

backend:
  output_root: "./internal"
  module_path: "github.com/epmp/backend"
  artifacts:
    - dto
    - repository
    - rest
    - test

frontend:
  output_root: "./src"
  base_url: "/api"

domain:
  name: %s
  package: %s
  table: %s
  bounded_context: %s

  fields:%s

  behaviors:
    soft_delete: true
    pagination: true
    search: true
    sort_by:
      - created_at

  rest:
    base_path: /api/%s
    auth_required: true
    operations:
      - create
      - read
      - update
      - delete
      - list

  events:
    - %sCreated
`

	sqlUpTemplate := `CREATE TABLE %s (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),%s,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON %s
FOR EACH ROW
EXECUTE FUNCTION trigger_set_updated_at();
`

	sqlDownTemplate := `DROP TRIGGER IF EXISTS set_updated_at ON %s;
DROP TABLE IF EXISTS %s;
`

	baseMigrationIndex := 15
	
	// Create run_codegen.sh script
	shContent := "#!/bin/bash\nset -e\n"
	shContent += "cd /Users/macbookpro/pjc/personal/epmp/tools/epmp-sdk/be/codegen\n"

	for i, e := range entities {
		// Generate YAML
		yamlPath := fmt.Sprintf("/Users/macbookpro/pjc/personal/epmp/tools/epmp-sdk/schemas/%s.yaml", e.Package)
		yamlContent := fmt.Sprintf(yamlTemplate, e.Name, e.Package, e.Table, e.BoundedContext, e.Fields, e.Table, e.Name)
		os.WriteFile(yamlPath, []byte(yamlContent), 0644)

		// Generate SQL Migrations
		upName := fmt.Sprintf("/Users/macbookpro/pjc/personal/epmp/backend/migrations/%06d_create_%s.up.sql", baseMigrationIndex+i, e.Table)
		downName := fmt.Sprintf("/Users/macbookpro/pjc/personal/epmp/backend/migrations/%06d_create_%s.down.sql", baseMigrationIndex+i, e.Table)
		
		upContent := fmt.Sprintf(sqlUpTemplate, e.Table, e.SQLFields, e.Table)
		downContent := fmt.Sprintf(sqlDownTemplate, e.Table, e.Table)
		
		os.WriteFile(upName, []byte(upContent), 0644)
		os.WriteFile(downName, []byte(downContent), 0644)

		shContent += fmt.Sprintf("./epmp-codegen --config ../../schemas/%s.yaml --output ../../../../backend/internal\n", e.Package)
	}

	shContent += "\ncd /Users/macbookpro/pjc/personal/epmp/tools/epmp-sdk/fe/codegen\n"
	for _, e := range entities {
		shContent += fmt.Sprintf("./epmp-fe-codegen --config ../../schemas/%s.yaml --output ../../../../frontend/src\n", e.Package)
	}

	os.WriteFile("/Users/macbookpro/pjc/personal/epmp/generate_batch.sh", []byte(shContent), 0755)

	// Append to modules.go
	modulesGoPath := "/Users/macbookpro/pjc/personal/epmp/backend/internal/modules/modules.go"
	b, _ := os.ReadFile(modulesGoPath)
	content := string(b)

	// Inject imports
	imports := ""
	registers := ""
	for _, e := range entities {
		imports += fmt.Sprintf("\t\"github.com/epmp/backend/internal/modules/%s\"\n", e.Package)
		registers += fmt.Sprintf("\t%s.NewModule(db, log).RegisterRoutes(protected)\n", e.Package)
	}

	content = strings.Replace(content, "import (", "import (\n"+imports, 1)
	content = strings.Replace(content, "return nil", registers+"\n\treturn nil", 1)

	os.WriteFile(modulesGoPath, []byte(content), 0644)
	fmt.Println("Generation complete. Run generate_batch.sh now.")
}
