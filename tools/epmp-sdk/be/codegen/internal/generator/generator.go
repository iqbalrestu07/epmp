package generator

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/epmp/sdk/codegen/config"
	"github.com/epmp/sdk/codegen/internal/filesystem"
	"github.com/epmp/sdk/codegen/internal/renderer"
)

// Generator generates module artifacts using Renderer and Filesystem.
type Generator struct {
	renderer   renderer.Renderer
	filesystem filesystem.Filesystem
}

// New creates a Generator with the given Renderer and Filesystem.
func New(r renderer.Renderer, fs filesystem.Filesystem) *Generator {
	return &Generator{
		renderer:   r,
		filesystem: fs,
	}
}

// GenerateRequest contains everything needed to generate a module.
type GenerateRequest struct {
	Config     *config.GeneratorConfig
	OutputRoot string
	ModulePath string
}

// Generate generates all module artifacts.
func (g *Generator) Generate(req *GenerateRequest) error {
	if req.Config == nil {
		return fmt.Errorf("generator: config is required")
	}
	if req.OutputRoot == "" {
		return fmt.Errorf("generator: output root is required")
	}
	if req.ModulePath == "" {
		return fmt.Errorf("generator: module path is required")
	}

	domain := req.Config.Domain
	data := buildTemplateData(&domain, req.ModulePath)

	nameLower := strings.ToLower(domain.Name)
	moduleDir := filepath.Join(req.OutputRoot, "modules", domain.Package)

	artifacts := []artifact{
		{template: "module.md", path: filepath.Join(moduleDir, "MODULE.md")},
		{template: "module.go", path: filepath.Join(moduleDir, "module.go")},
		{template: "entity.go", path: filepath.Join(moduleDir, "domain", "entity", nameLower+".go")},
		{template: "repository.go", path: filepath.Join(moduleDir, "domain", "repository", nameLower+"_repository.go")},
		{template: "dto.go", path: filepath.Join(moduleDir, "application", "dto", nameLower+"_dto.go")},
		{template: "service.go", path: filepath.Join(moduleDir, "application", "service", nameLower+"_service.go")},
		{template: "repository_impl.go", path: filepath.Join(moduleDir, "infrastructure", "repository", nameLower+"_repository_impl.go")},
		{template: "handler.go", path: filepath.Join(moduleDir, "interfaces", "http", nameLower+"_handler.go")},
		{template: "routes.go", path: filepath.Join(moduleDir, "interfaces", "http", nameLower+"_routes.go")},
		{template: "test.go", path: filepath.Join(moduleDir, nameLower+"_test.go")},
	}

	for _, a := range artifacts {
		content, err := g.renderer.Render(a.template, data)
		if err != nil {
			return fmt.Errorf("generator: render %s: %w", a.template, err)
		}

		dir := filepath.Dir(a.path)
		if err := g.filesystem.CreateDir(dir); err != nil {
			return fmt.Errorf("generator: create dir for %s: %w", a.template, err)
		}

		if err := g.filesystem.WriteFile(a.path, content); err != nil {
			return fmt.Errorf("generator: write %s: %w", a.template, err)
		}
	}

	return nil
}

// DryRun renders all templates and prints the output paths without writing files.
func (g *Generator) DryRun(req *GenerateRequest) error {
	if req.Config == nil {
		return fmt.Errorf("generator: config is required")
	}
	if req.OutputRoot == "" {
		return fmt.Errorf("generator: output root is required")
	}
	if req.ModulePath == "" {
		return fmt.Errorf("generator: module path is required")
	}

	domain := req.Config.Domain
	data := buildTemplateData(&domain, req.ModulePath)

	artifacts := g.artifacts(req.OutputRoot, domain)

	for _, a := range artifacts {
		content, err := g.renderer.Render(a.template, data)
		if err != nil {
			return fmt.Errorf("generator: render %s: %w", a.template, err)
		}
		fmt.Printf("[dry-run] %s (%d bytes)\n", a.path, len(content))
	}

	return nil
}

func (g *Generator) artifacts(outputRoot string, domain config.DomainSpec) []artifact {
	nameLower := strings.ToLower(domain.Name)
	moduleDir := filepath.Join(outputRoot, "modules", domain.Package)
	return []artifact{
		{template: "module.md", path: filepath.Join(moduleDir, "MODULE.md")},
		{template: "module.go", path: filepath.Join(moduleDir, "module.go")},
		{template: "entity.go", path: filepath.Join(moduleDir, "domain", "entity", nameLower+".go")},
		{template: "repository.go", path: filepath.Join(moduleDir, "domain", "repository", nameLower+"_repository.go")},
		{template: "dto.go", path: filepath.Join(moduleDir, "application", "dto", nameLower+"_dto.go")},
		{template: "service.go", path: filepath.Join(moduleDir, "application", "service", nameLower+"_service.go")},
		{template: "repository_impl.go", path: filepath.Join(moduleDir, "infrastructure", "repository", nameLower+"_repository_impl.go")},
		{template: "handler.go", path: filepath.Join(moduleDir, "interfaces", "http", nameLower+"_handler.go")},
		{template: "routes.go", path: filepath.Join(moduleDir, "interfaces", "http", nameLower+"_routes.go")},
		{template: "test.go", path: filepath.Join(moduleDir, nameLower+"_test.go")},
	}
}

type artifact struct {
	template string
	path     string
}

// templateData is the data model passed to templates.
type templateData struct {
	Name           string      // PascalCase entity name (e.g. Property)
	Package        string      // lowercase package name (e.g. property)
	Table          string      // database table name
	ModulePath     string      // Go module path
	BoundedContext string      // bounded context name
	Fields         []fieldData // field list
	HasSoftDelete  bool        // soft delete enabled
	HasPagination  bool        // pagination enabled
	HasSearch      bool        // search enabled
	RESTBasePath   string      // REST base path (e.g. /api/properties)
	Operations     []string    // enabled operations
}

type fieldData struct {
	Name       string // Go field name (PascalCase)
	GoType     string // Go type (string, int, etc.)
	JSONName   string // JSON tag
	ColumnName string // database column name
	PrimaryKey bool
	Nullable   bool
	Searchable bool
}

func buildTemplateData(domain *config.DomainSpec, modulePath string) map[string]any {
	operations := make([]string, 0, len(domain.REST.Operations))
	for _, op := range domain.REST.Operations {
		operations = append(operations, string(op))
	}

	fields := make([]fieldData, 0, len(domain.Fields))
	for _, f := range domain.Fields {
		fields = append(fields, fieldData{
			Name:       toPascalCase(f.Name),
			GoType:     toGoType(f.Type),
			JSONName:   strings.ToLower(f.Name),
			ColumnName: toSnakeCase(f.Name),
			PrimaryKey: f.PrimaryKey,
			Nullable:   f.Nullable,
			Searchable: f.Searchable,
		})
	}

	return map[string]any{
		"Name":           domain.Name,
		"Package":        domain.Package,
		"Table":          domain.Table,
		"ModulePath":     modulePath,
		"BoundedContext": domain.BoundedContext,
		"Fields":         fields,
		"HasSoftDelete":  domain.Behaviors.SoftDelete,
		"HasPagination":  domain.Behaviors.Pagination,
		"HasSearch":      domain.Behaviors.Search,
		"RESTBasePath":   domain.REST.BasePath,
		"Operations":     operations,
	}
}

func toGoType(t config.FieldType) string {
	switch t {
	case config.FieldTypeUUID:
		return "string"
	case config.FieldTypeString:
		return "string"
	case config.FieldTypeInt:
		return "int"
	case config.FieldTypeBool:
		return "bool"
	case config.FieldTypeText:
		return "string"
	case config.FieldTypeEnum:
		return "string"
	case config.FieldTypeTimestamp:
		return "time.Time"
	case config.FieldTypeDecimal:
		return "float64"
	default:
		return "string"
	}
}

func toPascalCase(s string) string {
	if s == "" {
		return s
	}
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, "")
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(toLowerRune(r))
	}
	return result.String()
}

func toLowerRune(r rune) rune {
	if r >= 'A' && r <= 'Z' {
		return r + 32
	}
	return r
}
