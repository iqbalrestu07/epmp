package renderer

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
)

// Renderer is the interface for loading and rendering Go templates.
type Renderer interface {
	// Render loads the template identified by name and renders it with data.
	// Returns the rendered content as bytes.
	Render(name string, data map[string]any) ([]byte, error)
}

// TemplateLoader is the interface for loading raw template content by name.
type TemplateLoader interface {
	// Load returns the raw template content for the given name.
	Load(name string) (string, error)
}

// renderer implements Renderer using text/template.
type renderer struct {
	loader   TemplateLoader
	funcMap  template.FuncMap
	leftDelim  string
	rightDelim string
}

// Option configures a renderer.
type Option func(*renderer)

// WithFuncMap adds custom template functions.
func WithFuncMap(fm template.FuncMap) Option {
	return func(r *renderer) {
		for k, v := range fm {
			r.funcMap[k] = v
		}
	}
}

// WithDelims sets custom template delimiters.
func WithDelims(left, right string) Option {
	return func(r *renderer) {
		r.leftDelim = left
		r.rightDelim = right
	}
}

// New creates a Renderer backed by the given TemplateLoader.
func New(loader TemplateLoader, opts ...Option) Renderer {
	r := &renderer{
		loader:   loader,
		funcMap:  defaultFuncMap(),
		leftDelim:  "",
		rightDelim: "",
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *renderer) Render(name string, data map[string]any) ([]byte, error) {
	raw, err := r.loader.Load(name)
	if err != nil {
		return nil, fmt.Errorf("renderer: load template %q: %w", name, err)
	}

	tmpl, err := template.New(name).
		Funcs(r.funcMap).
		Delims(r.leftDelim, r.rightDelim).
		Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("renderer: parse template %q: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("renderer: execute template %q: %w", name, err)
	}

	return buf.Bytes(), nil
}

// FSLoader is a TemplateLoader backed by an fs.FS.
// Template names are resolved as file paths relative to the FS root
// with a .tmpl extension appended if not present.
type FSLoader struct {
	fsys fs.FS
	ext  string
}

// FSLoaderOption configures an FSLoader.
type FSLoaderOption func(*FSLoader)

// WithExtension sets the file extension used when resolving template names.
// Defaults to ".tmpl".
func WithExtension(ext string) FSLoaderOption {
	return func(l *FSLoader) {
		l.ext = ext
	}
}

// NewFSLoader creates a TemplateLoader from an fs.FS.
func NewFSLoader(fsys fs.FS, opts ...FSLoaderOption) TemplateLoader {
	l := &FSLoader{
		fsys: fsys,
		ext:  ".tmpl",
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (l *FSLoader) Load(name string) (string, error) {
	path := name
	if !strings.HasSuffix(path, l.ext) {
		path = filepath.Clean(path + l.ext)
	}

	data, err := fs.ReadFile(l.fsys, path)
	if err != nil {
		return "", fmt.Errorf("fsloader: read %q: %w", path, err)
	}

	return string(data), nil
}

// MapLoader is a TemplateLoader backed by an in-memory map.
// Useful for testing.
type MapLoader struct {
	templates map[string]string
}

// NewMapLoader creates a TemplateLoader from a map of template names to content.
func NewMapLoader(templates map[string]string) TemplateLoader {
	return &MapLoader{
		templates: templates,
	}
}

func (l *MapLoader) Load(name string) (string, error) {
	content, ok := l.templates[name]
	if !ok {
		return "", fmt.Errorf("maploader: template %q not found", name)
	}
	return content, nil
}

func defaultFuncMap() template.FuncMap {
	return template.FuncMap{
		"title": strings.Title,
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"trim":  strings.TrimSpace,
	}
}
