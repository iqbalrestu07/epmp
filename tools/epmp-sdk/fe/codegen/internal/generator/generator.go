package generator

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/epmp/sdk/fe-codegen/config"
	"github.com/epmp/sdk/fe-codegen/internal/filesystem"
	"github.com/epmp/sdk/fe-codegen/internal/renderer"
)

// Generator generates frontend module artifacts using Renderer and Filesystem.
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

// GenerateRequest contains everything needed to generate a frontend module.
type GenerateRequest struct {
	Config     *config.GeneratorConfig
	OutputRoot string
	BaseURL    string
}

// Generate generates all frontend module artifacts.
func (g *Generator) Generate(req *GenerateRequest) error {
	if req.Config == nil {
		return fmt.Errorf("generator: config is required")
	}
	if req.OutputRoot == "" {
		return fmt.Errorf("generator: output root is required")
	}

	domain := req.Config.Domain
	baseURL := req.BaseURL
	if baseURL == "" {
		baseURL = req.Config.Frontend.BaseURL
	}
	data := buildTemplateData(&domain, baseURL)

	featureDir := filepath.Join(req.OutputRoot, "features", domain.Package)

	artifacts := []artifact{
		// types
		{template: "types.ts", path: filepath.Join(featureDir, "types", "index.ts")},
		// schema (Zod validation)
		{template: "schema.ts", path: filepath.Join(featureDir, "schema", "index.ts")},
		// api client
		{template: "api.ts", path: filepath.Join(featureDir, "api", "index.ts")},
		// hooks (TanStack Query)
		{template: "hooks.ts", path: filepath.Join(featureDir, "hooks", "index.ts")},
		// form component
		{template: "form.tsx", path: filepath.Join(featureDir, "components", domain.Name+"Form.tsx")},
		// table component
		{template: "table.tsx", path: filepath.Join(featureDir, "components", domain.Name+"Table.tsx")},
		// list page
		{template: "list-page.tsx", path: filepath.Join(featureDir, "pages", domain.Name+"ListPage.tsx")},
		// create page
		{template: "create-page.tsx", path: filepath.Join(featureDir, "pages", domain.Name+"CreatePage.tsx")},
		// edit page
		{template: "edit-page.tsx", path: filepath.Join(featureDir, "pages", domain.Name+"EditPage.tsx")},
		// detail page
		{template: "detail-page.tsx", path: filepath.Join(featureDir, "pages", domain.Name+"DetailPage.tsx")},
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

	domain := req.Config.Domain
	baseURL := req.BaseURL
	if baseURL == "" {
		baseURL = req.Config.Frontend.BaseURL
	}
	data := buildTemplateData(&domain, baseURL)

	featureDir := filepath.Join(req.OutputRoot, "features", domain.Package)

	artifacts := g.artifacts(featureDir, domain)

	for _, a := range artifacts {
		content, err := g.renderer.Render(a.template, data)
		if err != nil {
			return fmt.Errorf("generator: render %s: %w", a.template, err)
		}
		fmt.Printf("[dry-run] %s (%d bytes)\n", a.path, len(content))
	}

	return nil
}

func (g *Generator) artifacts(featureDir string, domain config.DomainSpec) []artifact {
	return []artifact{
		{template: "types.ts", path: filepath.Join(featureDir, "types", "index.ts")},
		{template: "schema.ts", path: filepath.Join(featureDir, "schema", "index.ts")},
		{template: "api.ts", path: filepath.Join(featureDir, "api", "index.ts")},
		{template: "hooks.ts", path: filepath.Join(featureDir, "hooks", "index.ts")},
		{template: "form.tsx", path: filepath.Join(featureDir, "components", domain.Name+"Form.tsx")},
		{template: "table.tsx", path: filepath.Join(featureDir, "components", domain.Name+"Table.tsx")},
		{template: "list-page.tsx", path: filepath.Join(featureDir, "pages", domain.Name+"ListPage.tsx")},
		{template: "create-page.tsx", path: filepath.Join(featureDir, "pages", domain.Name+"CreatePage.tsx")},
		{template: "edit-page.tsx", path: filepath.Join(featureDir, "pages", domain.Name+"EditPage.tsx")},
		{template: "detail-page.tsx", path: filepath.Join(featureDir, "pages", domain.Name+"DetailPage.tsx")},
	}
}

type artifact struct {
	template string
	path     string
}

type fieldData struct {
	Name        string // PascalCase (e.g. FullName)
	TsType      string // TypeScript type (string, number, boolean, etc.)
	JSONName    string // JSON key (e.g. full_name)
	ColumnName  string // snake_case
	PrimaryKey  bool
	Nullable    bool
	Searchable  bool
	IsAuto      bool
	EnumValues  []string
	HasEnum     bool
	MaxLength   int
	HasMaxLen   bool
}

type templateData struct {
	Name           string   // PascalCase (e.g. Property)
	NameLower      string   // camelCase (e.g. property)
	Package        string   // lowercase (e.g. property)
	Table          string
	BoundedContext string
	BaseURL        string   // API base URL
	RESTBasePath   string   // e.g. /api/properties
	Fields         []fieldData
	HasSoftDelete  bool
	HasPagination  bool
	HasSearch      bool
	Operations     []string
	HasCreate      bool
	HasRead        bool
	HasUpdate      bool
	HasDelete      bool
	HasList        bool
}

func buildTemplateData(domain *config.DomainSpec, baseURL string) map[string]any {
	operations := make([]string, 0, len(domain.REST.Operations))
	hasCreate, hasRead, hasUpdate, hasDelete, hasList := false, false, false, false, false
	for _, op := range domain.REST.Operations {
		operations = append(operations, string(op))
		switch op {
		case config.OperationCreate:
			hasCreate = true
		case config.OperationRead:
			hasRead = true
		case config.OperationUpdate:
			hasUpdate = true
		case config.OperationDelete:
			hasDelete = true
		case config.OperationList:
			hasList = true
		}
	}

	fields := make([]fieldData, 0, len(domain.Fields))
	for _, f := range domain.Fields {
		fd := fieldData{
			Name:       toPascalCase(f.Name),
			TsType:     toTsType(f.Type),
			JSONName:   strings.ToLower(f.Name),
			ColumnName: toSnakeCase(f.Name),
			PrimaryKey: f.PrimaryKey,
			Nullable:   f.Nullable,
			Searchable: f.Searchable,
			IsAuto:     f.Auto,
			EnumValues: f.EnumValues,
			HasEnum:    f.Type == config.FieldTypeEnum,
			MaxLength:  f.MaxLength,
			HasMaxLen:  f.MaxLength > 0,
		}
		fields = append(fields, fd)
	}

	return map[string]any{
		"Name":           domain.Name,
		"NameLower":      toCamelCase(domain.Name),
		"Package":        domain.Package,
		"Table":          domain.Table,
		"BoundedContext": domain.BoundedContext,
		"BaseURL":        baseURL,
		"RESTBasePath":   domain.REST.BasePath,
		"Fields":         fields,
		"HasSoftDelete":  domain.Behaviors.SoftDelete,
		"HasPagination":  domain.Behaviors.Pagination,
		"HasSearch":      domain.Behaviors.Search,
		"Operations":     operations,
		"HasCreate":      hasCreate,
		"HasRead":        hasRead,
		"HasUpdate":      hasUpdate,
		"HasDelete":      hasDelete,
		"HasList":        hasList,
	}
}

func toTsType(t config.FieldType) string {
	switch t {
	case config.FieldTypeUUID:
		return "string"
	case config.FieldTypeString:
		return "string"
	case config.FieldTypeInt:
		return "number"
	case config.FieldTypeBool:
		return "boolean"
	case config.FieldTypeText:
		return "string"
	case config.FieldTypeEnum:
		return "string"
	case config.FieldTypeTimestamp:
		return "string"
	case config.FieldTypeDecimal:
		return "number"
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

func toCamelCase(s string) string {
	pascal := toPascalCase(s)
	if pascal == "" {
		return pascal
	}
	return strings.ToLower(pascal[:1]) + pascal[1:]
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
