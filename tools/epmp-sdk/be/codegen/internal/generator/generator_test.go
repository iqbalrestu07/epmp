package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/epmp/sdk/codegen/config"
	"github.com/epmp/sdk/codegen/internal/filesystem"
)

func TestGenerate_AllArtifacts(t *testing.T) {
	cfg := &config.GeneratorConfig{
		Version: "1",
		Backend: config.BackendConfig{
			OutputRoot: "",
			ModulePath: "github.com/epmp/backend",
		},
		Domain: config.DomainSpec{
			Name:    "Property",
			Package: "property",
			Table:   "properties",
			Fields: []config.FieldSpec{
				{Name: "id", Type: config.FieldTypeUUID, PrimaryKey: true},
				{Name: "name", Type: config.FieldTypeString, Searchable: true},
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

	outputRoot := t.TempDir()
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	err := gen.Generate(&GenerateRequest{
		Config:     cfg,
		OutputRoot: outputRoot,
		ModulePath: "github.com/epmp/backend",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedFiles := []string{
		"modules/property/MODULE.md",
		"modules/property/module.go",
		"modules/property/domain/entity/property.go",
		"modules/property/domain/repository/property_repository.go",
		"modules/property/application/dto/property_dto.go",
		"modules/property/application/service/property_service.go",
		"modules/property/infrastructure/repository/property_repository_impl.go",
		"modules/property/interfaces/http/property_handler.go",
		"modules/property/interfaces/http/property_routes.go",
		"modules/property/property_test.go",
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

func TestGenerate_ModuleContent(t *testing.T) {
	cfg := &config.GeneratorConfig{
		Domain: config.DomainSpec{
			Name:    "Tenant",
			Package: "tenant",
			Table:   "tenants",
			Fields: []config.FieldSpec{
				{Name: "id", Type: config.FieldTypeUUID, PrimaryKey: true},
				{Name: "full_name", Type: config.FieldTypeString},
			},
			REST: config.RESTSpec{
				BasePath:   "/api/tenants",
				Operations: []config.Operation{config.OperationCreate, config.OperationList},
			},
		},
	}

	outputRoot := t.TempDir()
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	if err := gen.Generate(&GenerateRequest{
		Config:     cfg,
		OutputRoot: outputRoot,
		ModulePath: "github.com/epmp/backend",
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := fs.Read(filepath.Join(outputRoot, "modules", "tenant", "MODULE.md"))
	if err != nil {
		t.Fatalf("read MODULE.md: %v", err)
	}

	if !strings.Contains(string(content), "Tenant") {
		t.Error("MODULE.md should contain entity name")
	}
	if !strings.Contains(string(content), "tenant") {
		t.Error("MODULE.md should contain package name")
	}
}

func TestGenerate_EntityContent(t *testing.T) {
	cfg := &config.GeneratorConfig{
		Domain: config.DomainSpec{
			Name:    "Property",
			Package: "property",
			Table:   "properties",
			Fields: []config.FieldSpec{
				{Name: "id", Type: config.FieldTypeUUID, PrimaryKey: true},
				{Name: "name", Type: config.FieldTypeString},
			},
			Behaviors: config.BehaviorSpec{SoftDelete: true},
		},
	}

	outputRoot := t.TempDir()
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	if err := gen.Generate(&GenerateRequest{
		Config:     cfg,
		OutputRoot: outputRoot,
		ModulePath: "github.com/epmp/backend",
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := fs.Read(filepath.Join(outputRoot, "modules", "property", "domain", "entity", "property.go"))
	if err != nil {
		t.Fatalf("read entity: %v", err)
	}

	str := string(content)
	if !strings.Contains(str, "type Property struct") {
		t.Error("entity should contain struct definition")
	}
	if !strings.Contains(str, "DeletedAt") {
		t.Error("entity should contain soft delete field")
	}
}

func TestGenerate_RepositoryContent(t *testing.T) {
	cfg := &config.GeneratorConfig{
		Domain: config.DomainSpec{
			Name:    "Property",
			Package: "property",
			Table:   "properties",
			Fields: []config.FieldSpec{
				{Name: "id", Type: config.FieldTypeUUID, PrimaryKey: true},
			},
		},
	}

	outputRoot := t.TempDir()
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	if err := gen.Generate(&GenerateRequest{
		Config:     cfg,
		OutputRoot: outputRoot,
		ModulePath: "github.com/epmp/backend",
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := fs.Read(filepath.Join(outputRoot, "modules", "property", "domain", "repository", "property_repository.go"))
	if err != nil {
		t.Fatalf("read repository: %v", err)
	}

	str := string(content)
	if !strings.Contains(str, "PropertyRepository interface") {
		t.Error("repository should contain interface definition")
	}
	if !strings.Contains(str, "Save") {
		t.Error("repository should contain Save method")
	}
	if !strings.Contains(str, "FindByID") {
		t.Error("repository should contain FindByID method")
	}
}

func TestGenerate_DTOContent(t *testing.T) {
	cfg := &config.GeneratorConfig{
		Domain: config.DomainSpec{
			Name:    "Property",
			Package: "property",
			Table:   "properties",
			Fields: []config.FieldSpec{
				{Name: "id", Type: config.FieldTypeUUID, PrimaryKey: true},
				{Name: "name", Type: config.FieldTypeString},
			},
		},
	}

	outputRoot := t.TempDir()
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	if err := gen.Generate(&GenerateRequest{
		Config:     cfg,
		OutputRoot: outputRoot,
		ModulePath: "github.com/epmp/backend",
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := fs.Read(filepath.Join(outputRoot, "modules", "property", "application", "dto", "property_dto.go"))
	if err != nil {
		t.Fatalf("read dto: %v", err)
	}

	str := string(content)
	if !strings.Contains(str, "CreatePropertyRequest") {
		t.Error("dto should contain create request")
	}
	if !strings.Contains(str, "UpdatePropertyRequest") {
		t.Error("dto should contain update request")
	}
	if !strings.Contains(str, "PropertyResponse") {
		t.Error("dto should contain response")
	}
}

func TestGenerate_HandlerContent(t *testing.T) {
	cfg := &config.GeneratorConfig{
		Domain: config.DomainSpec{
			Name:    "Property",
			Package: "property",
			Table:   "properties",
			Fields: []config.FieldSpec{
				{Name: "id", Type: config.FieldTypeUUID, PrimaryKey: true},
			},
			REST: config.RESTSpec{
				BasePath:   "/api/properties",
				Operations: []config.Operation{config.OperationCreate, config.OperationRead},
			},
		},
	}

	outputRoot := t.TempDir()
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	if err := gen.Generate(&GenerateRequest{
		Config:     cfg,
		OutputRoot: outputRoot,
		ModulePath: "github.com/epmp/backend",
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := fs.Read(filepath.Join(outputRoot, "modules", "property", "interfaces", "http", "property_handler.go"))
	if err != nil {
		t.Fatalf("read handler: %v", err)
	}

	str := string(content)
	if !strings.Contains(str, "PropertyHandler") {
		t.Error("handler should contain handler struct")
	}
	if !strings.Contains(str, "func (h *PropertyHandler) Create") {
		t.Error("handler should contain Create method")
	}
	if !strings.Contains(str, "func (h *PropertyHandler) GetByID") {
		t.Error("handler should contain GetByID method")
	}
}

func TestGenerate_TestContent(t *testing.T) {
	cfg := &config.GeneratorConfig{
		Domain: config.DomainSpec{
			Name:    "Property",
			Package: "property",
			Table:   "properties",
			Fields: []config.FieldSpec{
				{Name: "id", Type: config.FieldTypeUUID, PrimaryKey: true},
			},
			REST: config.RESTSpec{
				Operations: []config.Operation{config.OperationCreate},
			},
		},
	}

	outputRoot := t.TempDir()
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	if err := gen.Generate(&GenerateRequest{
		Config:     cfg,
		OutputRoot: outputRoot,
		ModulePath: "github.com/epmp/backend",
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := fs.Read(filepath.Join(outputRoot, "modules", "property", "property_test.go"))
	if err != nil {
		t.Fatalf("read test: %v", err)
	}

	str := string(content)
	if !strings.Contains(str, "TestProperty_Placeholder") {
		t.Error("test should contain placeholder test function")
	}
}

func TestGenerate_NilConfig(t *testing.T) {
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	err := gen.Generate(&GenerateRequest{
		Config:     nil,
		OutputRoot: t.TempDir(),
		ModulePath: "github.com/epmp/backend",
	})
	if err == nil {
		t.Fatal("expected error for nil config, got nil")
	}
}

func TestGenerate_EmptyOutputRoot(t *testing.T) {
	cfg := &config.GeneratorConfig{
		Domain: config.DomainSpec{
			Name:    "Property",
			Package: "property",
			Table:   "properties",
		},
	}
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	err := gen.Generate(&GenerateRequest{
		Config:     cfg,
		OutputRoot: "",
		ModulePath: "github.com/epmp/backend",
	})
	if err == nil {
		t.Fatal("expected error for empty output root, got nil")
	}
}

func TestGenerate_EmptyModulePath(t *testing.T) {
	cfg := &config.GeneratorConfig{
		Domain: config.DomainSpec{
			Name:    "Property",
			Package: "property",
			Table:   "properties",
		},
	}
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	err := gen.Generate(&GenerateRequest{
		Config:     cfg,
		OutputRoot: t.TempDir(),
		ModulePath: "",
	})
	if err == nil {
		t.Fatal("expected error for empty module path, got nil")
	}
}

func TestGenerate_Idempotent(t *testing.T) {
	cfg := &config.GeneratorConfig{
		Domain: config.DomainSpec{
			Name:    "Property",
			Package: "property",
			Table:   "properties",
			Fields: []config.FieldSpec{
				{Name: "id", Type: config.FieldTypeUUID, PrimaryKey: true},
			},
		},
	}

	outputRoot := t.TempDir()
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	req := &GenerateRequest{
		Config:     cfg,
		OutputRoot: outputRoot,
		ModulePath: "github.com/epmp/backend",
	}

	if err := gen.Generate(req); err != nil {
		t.Fatalf("first generate: %v", err)
	}
	if err := gen.Generate(req); err != nil {
		t.Fatalf("second generate: %v", err)
	}

	path := filepath.Join(outputRoot, "modules", "property", "MODULE.md")
	exists, err := fs.Exists(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if !exists {
		t.Error("file should still exist after second generate")
	}
}

func TestGenerate_CreatesDirectoryStructure(t *testing.T) {
	cfg := &config.GeneratorConfig{
		Domain: config.DomainSpec{
			Name:    "Contract",
			Package: "contract",
			Table:   "contracts",
			Fields: []config.FieldSpec{
				{Name: "id", Type: config.FieldTypeUUID, PrimaryKey: true},
			},
		},
	}

	outputRoot := t.TempDir()
	fs := filesystem.New()
	gen := New(DefaultRenderer(), fs)

	if err := gen.Generate(&GenerateRequest{
		Config:     cfg,
		OutputRoot: outputRoot,
		ModulePath: "github.com/epmp/backend",
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dirs := []string{
		"modules",
		"modules/contract",
		"modules/contract/domain",
		"modules/contract/domain/entity",
		"modules/contract/domain/repository",
		"modules/contract/application",
		"modules/contract/application/dto",
		"modules/contract/application/service",
		"modules/contract/infrastructure",
		"modules/contract/infrastructure/repository",
		"modules/contract/interfaces",
		"modules/contract/interfaces/http",
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

func TestToGoType(t *testing.T) {
	tests := []struct {
		input    config.FieldType
		expected string
	}{
		{config.FieldTypeUUID, "string"},
		{config.FieldTypeString, "string"},
		{config.FieldTypeInt, "int"},
		{config.FieldTypeBool, "bool"},
		{config.FieldTypeText, "string"},
		{config.FieldTypeEnum, "string"},
		{config.FieldTypeTimestamp, "time.Time"},
		{config.FieldTypeDecimal, "float64"},
	}

	for _, tt := range tests {
		got := toGoType(tt.input)
		if got != tt.expected {
			t.Errorf("toGoType(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"id", "Id"},
		{"full_name", "FullName"},
		{"is_active", "IsActive"},
		{"name", "Name"},
		{"", ""},
	}

	for _, tt := range tests {
		got := toPascalCase(tt.input)
		if got != tt.expected {
			t.Errorf("toPascalCase(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"id", "id"},
		{"fullName", "full_name"},
		{"isActive", "is_active"},
		{"name", "name"},
	}

	for _, tt := range tests {
		got := toSnakeCase(tt.input)
		if got != tt.expected {
			t.Errorf("toSnakeCase(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
