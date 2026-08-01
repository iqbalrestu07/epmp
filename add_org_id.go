package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	schemaDir := "/Users/macbookpro/pjc/personal/epmp/tools/epmp-sdk/schemas"
	files, err := os.ReadDir(schemaDir)
	if err != nil {
		panic(err)
	}

	tables := []string{}

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" && f.Name() != "organization.yaml" {
			path := filepath.Join(schemaDir, f.Name())
			contentBytes, err := os.ReadFile(path)
			if err != nil {
				panic(err)
			}
			content := string(contentBytes)
			
			// Only add if not already there
			if !strings.Contains(content, "name: organization_id") {
				orgField := `
    - name: organization_id
      type: uuid
      searchable: true`
				
				content = strings.Replace(content, "fields:", "fields:"+orgField, 1)
				os.WriteFile(path, []byte(content), 0644)
			}
			
			// Extract table name to alter
			lines := strings.Split(content, "\n")
			for _, line := range lines {
				if strings.HasPrefix(strings.TrimSpace(line), "table:") {
					tableName := strings.TrimSpace(strings.Split(line, ":")[1])
					tables = append(tables, tableName)
					break
				}
			}
		}
	}

	// 000029_create_organizations
	up29 := `CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    domain VARCHAR(100),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON organizations FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();`
	down29 := `DROP TRIGGER IF EXISTS set_updated_at ON organizations; DROP TABLE IF EXISTS organizations;`
	os.WriteFile("/Users/macbookpro/pjc/personal/epmp/backend/migrations/000029_create_organizations.up.sql", []byte(up29), 0644)
	os.WriteFile("/Users/macbookpro/pjc/personal/epmp/backend/migrations/000029_create_organizations.down.sql", []byte(down29), 0644)

	// 000030_add_org_id_to_all
	up30 := "ALTER TABLE users ADD COLUMN IF NOT EXISTS organization_id UUID REFERENCES organizations(id);\n"
	up30 += "ALTER TABLE roles ADD COLUMN IF NOT EXISTS organization_id UUID REFERENCES organizations(id);\n"
	for _, t := range tables {
		up30 += fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS organization_id UUID REFERENCES organizations(id);\n", t)
	}
	down30 := "ALTER TABLE users DROP COLUMN IF EXISTS organization_id;\n"
	down30 += "ALTER TABLE roles DROP COLUMN IF EXISTS organization_id;\n"
	for _, t := range tables {
		down30 += fmt.Sprintf("ALTER TABLE %s DROP COLUMN IF EXISTS organization_id;\n", t)
	}
	os.WriteFile("/Users/macbookpro/pjc/personal/epmp/backend/migrations/000030_add_org_id_to_all.up.sql", []byte(up30), 0644)
	os.WriteFile("/Users/macbookpro/pjc/personal/epmp/backend/migrations/000030_add_org_id_to_all.down.sql", []byte(down30), 0644)

	// Generate script for codegen
	shContent := "#!/bin/bash\nset -e\ncd /Users/macbookpro/pjc/personal/epmp/tools/epmp-sdk/be/codegen\n"
	shContent += "./epmp-codegen --config ../../schemas/organization.yaml --output ../../../../backend/internal\n"
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" {
			shContent += fmt.Sprintf("./epmp-codegen --config ../../schemas/%s --output ../../../../backend/internal\n", f.Name())
		}
	}
	shContent += "\ncd /Users/macbookpro/pjc/personal/epmp/tools/epmp-sdk/fe/codegen\n"
	shContent += "./epmp-fe-codegen --config ../../schemas/organization.yaml --output ../../../../frontend/src\n"
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" {
			shContent += fmt.Sprintf("./epmp-fe-codegen --config ../../schemas/%s --output ../../../../frontend/src\n", f.Name())
		}
	}
	os.WriteFile("/Users/macbookpro/pjc/personal/epmp/regen_all.sh", []byte(shContent), 0755)

	fmt.Println("Added organization_id to all schemas and generated migrations.")
}
