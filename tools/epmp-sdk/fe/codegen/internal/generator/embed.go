package generator

import (
	"embed"
	"io/fs"
	"text/template"

	"github.com/epmp/sdk/fe-codegen/internal/renderer"
)

//go:embed templates/*
var templateFS embed.FS

// DefaultRenderer returns a Renderer backed by the embedded templates.
func DefaultRenderer() renderer.Renderer {
	sub, err := fs.Sub(templateFS, "templates")
	if err != nil {
		panic("generator: cannot create sub FS for templates: " + err.Error())
	}
	return renderer.New(
		renderer.NewFSLoader(sub, renderer.WithExtension(".tmpl")),
		renderer.WithFuncMap(template.FuncMap{
			"contains":  contains,
			"tsZodType": tsZodType,
		}),
	)
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}

// tsZodType maps a TypeScript type to the corresponding Zod schema function.
func tsZodType(tsType string, nullable bool) string {
	switch tsType {
	case "string":
		return "string"
	case "number":
		return "number"
	case "boolean":
		return "boolean"
	default:
		return "string"
	}
}
