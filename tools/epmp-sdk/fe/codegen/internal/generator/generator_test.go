package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/epmp/sdk/fe-codegen/config"
	"github.com/epmp/sdk/fe-codegen/internal/filesystem"
)

func sampleConfig() *config.GeneratorConfig {
	return &config.GeneratorConfig{
		Version: "1.0",
		Frontend: config.FrontendConfig{
			OutputRoot: "",
			BaseURL:    "/api",
		},
		Domain: config.DomainSpec{
			Name:    "Property",
			Package: "property",
			Table:   "properties",
			Fields: []config.FieldSpec{
				{Name: "id", Type: config.FieldTypeUUID, PrimaryKey: true, Auto: true},
				{Name: "name", Type: config.FieldTypeString, MaxLength: 255, Searchable: true},
				{Name: "address", Type: config.FieldTypeText, Nullable: true},
				{Name: "is_active", Type: config.FieldTypeBool},
			},
			Behaviors: config.BehaviorSpec{
				SoftDelete: true,
				Pagination: true,
				Search:     true,
			},
			REST: config.RESTSpec{
				BasePath: "/api/properties",
				Operations: []config.Operation{
					config.OperationCreate,
					config.OperationRead,
					config.OperationUpdate,
					config.OperationDelete,
					config.OperationList,
				},
			},
		},
	}
}

func TestGenerate_AllArtifacts(t *testing.T) {
	cfg := sampleConfig()
	outputRoot := t.TempDir()
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	err := gen.Generate(&GenerateRequest{
		Config:     cfg,
		OutputRoot: outputRoot,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedFiles := []string{
		"features/property/types/index.ts",
		"features/property/schema/index.ts",
		"features/property/api/index.ts",
		"features/property/hooks/index.ts",
		"features/property/components/PropertyForm.tsx",
		"features/property/components/PropertyTable.tsx",
		"features/property/pages/PropertyListPage.tsx",
		"features/property/pages/PropertyCreatePage.tsx",
		"features/property/pages/PropertyEditPage.tsx",
		"features/property/pages/PropertyDetailPage.tsx",
	}

	for _, f := range expectedFiles {
		path := filepath.Join(outputRoot, f)
		exists, err := fs.Exists(path)
		if err != nil {
			t.Fatalf("stat %s: %v", f, err)
		}
		if !exists {
			t.Errorf("expected file %s to exist", f)
		}
	}
}

func TestGenerate_TypesContent(t *testing.T) {
	cfg := sampleConfig()
	outputRoot := t.TempDir()
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	if err := gen.Generate(&GenerateRequest{Config: cfg, OutputRoot: outputRoot}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := fs.Read(filepath.Join(outputRoot, "features/property/types/index.ts"))
	if err != nil {
		t.Fatalf("read types: %v", err)
	}

	str := string(content)
	if !strings.Contains(str, "export interface Property") {
		t.Error("types should contain Property interface")
	}
	if !strings.Contains(str, "CreatePropertyRequest") {
		t.Error("types should contain CreatePropertyRequest")
	}
	if !strings.Contains(str, "PropertyListResponse") {
		t.Error("types should contain PropertyListResponse")
	}
}

func TestGenerate_SchemaContent(t *testing.T) {
	cfg := sampleConfig()
	outputRoot := t.TempDir()
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	if err := gen.Generate(&GenerateRequest{Config: cfg, OutputRoot: outputRoot}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := fs.Read(filepath.Join(outputRoot, "features/property/schema/index.ts"))
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	str := string(content)
	if !strings.Contains(str, "z.object") {
		t.Error("schema should contain zod object")
	}
	if !strings.Contains(str, "propertySchema") {
		t.Error("schema should contain propertySchema")
	}
}

func TestGenerate_ApiContent(t *testing.T) {
	cfg := sampleConfig()
	outputRoot := t.TempDir()
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	if err := gen.Generate(&GenerateRequest{Config: cfg, OutputRoot: outputRoot}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := fs.Read(filepath.Join(outputRoot, "features/property/api/index.ts"))
	if err != nil {
		t.Fatalf("read api: %v", err)
	}

	str := string(content)
	if !strings.Contains(str, "fetchPropertys") {
		t.Error("api should contain fetchPropertys function")
	}
	if !strings.Contains(str, "fetchPropertyById") {
		t.Error("api should contain fetchPropertyById function")
	}
	if !strings.Contains(str, "createProperty") {
		t.Error("api should contain createProperty function")
	}
}

func TestGenerate_HooksContent(t *testing.T) {
	cfg := sampleConfig()
	outputRoot := t.TempDir()
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	if err := gen.Generate(&GenerateRequest{Config: cfg, OutputRoot: outputRoot}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := fs.Read(filepath.Join(outputRoot, "features/property/hooks/index.ts"))
	if err != nil {
		t.Fatalf("read hooks: %v", err)
	}

	str := string(content)
	if !strings.Contains(str, "usePropertys") {
		t.Error("hooks should contain usePropertys")
	}
	if !strings.Contains(str, "useCreateProperty") {
		t.Error("hooks should contain useCreateProperty")
	}
	if !strings.Contains(str, "useQueryClient") {
		t.Error("hooks should use useQueryClient for cache invalidation")
	}
}

func TestGenerate_FormContent(t *testing.T) {
	cfg := sampleConfig()
	outputRoot := t.TempDir()
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	if err := gen.Generate(&GenerateRequest{Config: cfg, OutputRoot: outputRoot}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := fs.Read(filepath.Join(outputRoot, "features/property/components/PropertyForm.tsx"))
	if err != nil {
		t.Fatalf("read form: %v", err)
	}

	str := string(content)
	if !strings.Contains(str, "PropertyForm") {
		t.Error("form should contain PropertyForm component")
	}
	if !strings.Contains(str, "useForm") {
		t.Error("form should use react-hook-form")
	}
	if !strings.Contains(str, "zodResolver") {
		t.Error("form should use zodResolver")
	}
}

func TestGenerate_TableContent(t *testing.T) {
	cfg := sampleConfig()
	outputRoot := t.TempDir()
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	if err := gen.Generate(&GenerateRequest{Config: cfg, OutputRoot: outputRoot}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := fs.Read(filepath.Join(outputRoot, "features/property/components/PropertyTable.tsx"))
	if err != nil {
		t.Fatalf("read table: %v", err)
	}

	str := string(content)
	if !strings.Contains(str, "PropertyTable") {
		t.Error("table should contain PropertyTable component")
	}
	if !strings.Contains(str, "useReactTable") {
		t.Error("table should use TanStack Table")
	}
}

func TestGenerate_ListPageContent(t *testing.T) {
	cfg := sampleConfig()
	outputRoot := t.TempDir()
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	if err := gen.Generate(&GenerateRequest{Config: cfg, OutputRoot: outputRoot}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := fs.Read(filepath.Join(outputRoot, "features/property/pages/PropertyListPage.tsx"))
	if err != nil {
		t.Fatalf("read list page: %v", err)
	}

	str := string(content)
	if !strings.Contains(str, "PropertyListPage") {
		t.Error("list page should contain PropertyListPage component")
	}
	if !strings.Contains(str, "usePropertys") {
		t.Error("list page should use usePropertys hook")
	}
}

func TestGenerate_NilConfig(t *testing.T) {
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	err := gen.Generate(&GenerateRequest{
		Config:     nil,
		OutputRoot: t.TempDir(),
	})
	if err == nil {
		t.Fatal("expected error for nil config, got nil")
	}
}

func TestGenerate_EmptyOutputRoot(t *testing.T) {
	cfg := sampleConfig()
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	err := gen.Generate(&GenerateRequest{
		Config:     cfg,
		OutputRoot: "",
	})
	if err == nil {
		t.Fatal("expected error for empty output root, got nil")
	}
}

func TestGenerate_CreatesDirectoryStructure(t *testing.T) {
	cfg := sampleConfig()
	outputRoot := t.TempDir()
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	if err := gen.Generate(&GenerateRequest{Config: cfg, OutputRoot: outputRoot}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dirs := []string{
		"features",
		"features/property",
		"features/property/types",
		"features/property/schema",
		"features/property/api",
		"features/property/hooks",
		"features/property/components",
		"features/property/pages",
	}

	for _, d := range dirs {
		path := filepath.Join(outputRoot, d)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("directory %s not created: %v", d, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("expected %s to be a directory", d)
		}
	}
}

func TestGenerate_Idempotent(t *testing.T) {
	cfg := sampleConfig()
	outputRoot := t.TempDir()
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	req := &GenerateRequest{Config: cfg, OutputRoot: outputRoot}

	if err := gen.Generate(req); err != nil {
		t.Fatalf("first generate: %v", err)
	}
	if err := gen.Generate(req); err != nil {
		t.Fatalf("second generate: %v", err)
	}

	path := filepath.Join(outputRoot, "features/property/types/index.ts")
	exists, err := fs.Exists(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if !exists {
		t.Error("file should still exist after second generate")
	}
}

func TestToTsType(t *testing.T) {
	tests := []struct {
		input    config.FieldType
		expected string
	}{
		{config.FieldTypeUUID, "string"},
		{config.FieldTypeString, "string"},
		{config.FieldTypeInt, "number"},
		{config.FieldTypeBool, "boolean"},
		{config.FieldTypeText, "string"},
		{config.FieldTypeEnum, "string"},
		{config.FieldTypeTimestamp, "string"},
		{config.FieldTypeDecimal, "number"},
	}

	for _, tt := range tests {
		got := toTsType(tt.input)
		if got != tt.expected {
			t.Errorf("toTsType(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Property", "property"},
		{"FullName", "fullName"},
		{"IsActive", "isActive"},
	}

	for _, tt := range tests {
		got := toCamelCase(tt.input)
		if got != tt.expected {
			t.Errorf("toCamelCase(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
