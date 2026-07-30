package generator

import (
	"embed"
	"io/fs"
	"text/template"

	"github.com/epmp/sdk/codegen/internal/renderer"
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
			"contains": contains,
		}),
	)
}

// contains checks if a slice of strings contains the given value.
func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
