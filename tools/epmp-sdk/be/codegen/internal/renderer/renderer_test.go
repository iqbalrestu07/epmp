package renderer

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
	"text/template"
)

func TestNewMapLoader_Render(t *testing.T) {
	templates := map[string]string{
		"greeting": "Hello, {{.Name}}!",
	}
	r := New(NewMapLoader(templates))

	out, err := r.Render("greeting", map[string]any{"Name": "World"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []byte("Hello, World!")
	if !bytes.Equal(out, expected) {
		t.Errorf("got %q, want %q", out, expected)
	}
}

func TestRender_MissingTemplate(t *testing.T) {
	r := New(NewMapLoader(map[string]string{}))

	_, err := r.Render("nonexistent", map[string]any{})
	if err == nil {
		t.Fatal("expected error for missing template, got nil")
	}
}

func TestRender_InvalidTemplate(t *testing.T) {
	templates := map[string]string{
		"bad": "Hello, {{.Name}",
	}
	r := New(NewMapLoader(templates))

	_, err := r.Render("bad", map[string]any{})
	if err == nil {
		t.Fatal("expected error for invalid template syntax, got nil")
	}
}

type testPerson struct {
	Name string
}

func TestRender_ExecutionError(t *testing.T) {
	templates := map[string]string{
		"err": "{{.User.NonExistent}}",
	}
	r := New(NewMapLoader(templates))

	_, err := r.Render("err", map[string]any{
		"User": testPerson{Name: "test"},
	})
	if err == nil {
		t.Fatal("expected error for execution failure, got nil")
	}
}

func TestRender_NilData(t *testing.T) {
	templates := map[string]string{
		"static": "static content",
	}
	r := New(NewMapLoader(templates))

	out, err := r.Render("static", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []byte("static content")
	if !bytes.Equal(out, expected) {
		t.Errorf("got %q, want %q", out, expected)
	}
}

func TestRender_EmptyData(t *testing.T) {
	templates := map[string]string{
		"static": "static content",
	}
	r := New(NewMapLoader(templates))

	out, err := r.Render("static", map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []byte("static content")
	if !bytes.Equal(out, expected) {
		t.Errorf("got %q, want %q", out, expected)
	}
}

func TestRender_WithFuncMap(t *testing.T) {
	templates := map[string]string{
		"upper": "{{upper .Name}}",
	}
	r := New(
		NewMapLoader(templates),
		WithFuncMap(template.FuncMap{
			"upper": func(s string) string { return s },
		}),
	)

	out, err := r.Render("upper", map[string]any{"Name": "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []byte("test")
	if !bytes.Equal(out, expected) {
		t.Errorf("got %q, want %q", out, expected)
	}
}

func TestRender_DefaultFuncMap(t *testing.T) {
	templates := map[string]string{
		"funcs": "{{upper .Name}}|{{lower .Name}}|{{trim .Name}}",
	}
	r := New(NewMapLoader(templates))

	out, err := r.Render("funcs", map[string]any{"Name": "  Hello  "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []byte("  HELLO  |  hello  |Hello")
	if !bytes.Equal(out, expected) {
		t.Errorf("got %q, want %q", out, expected)
	}
}

func TestRender_WithCustomDelims(t *testing.T) {
	templates := map[string]string{
		"custom": "<<.Name>>",
	}
	r := New(
		NewMapLoader(templates),
		WithDelims("<<", ">>"),
	)

	out, err := r.Render("custom", map[string]any{"Name": "World"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []byte("World")
	if !bytes.Equal(out, expected) {
		t.Errorf("got %q, want %q", out, expected)
	}
}

func TestFSLoader_Load(t *testing.T) {
	fsys := fstest.MapFS{
		"greeting.tmpl": &fstest.MapFile{
			Data: []byte("Hello, {{.Name}}!"),
		},
	}
	r := New(NewFSLoader(fsys))

	out, err := r.Render("greeting", map[string]any{"Name": "World"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []byte("Hello, World!")
	if !bytes.Equal(out, expected) {
		t.Errorf("got %q, want %q", out, expected)
	}
}

func TestFSLoader_LoadWithExtension(t *testing.T) {
	fsys := fstest.MapFS{
		"entity.tpl": &fstest.MapFile{
			Data: []byte("type {{.Name}} struct{}"),
		},
	}
	r := New(NewFSLoader(fsys, WithExtension(".tpl")))

	out, err := r.Render("entity", map[string]any{"Name": "Property"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []byte("type Property struct{}")
	if !bytes.Equal(out, expected) {
		t.Errorf("got %q, want %q", out, expected)
	}
}

func TestFSLoader_LoadNotFound(t *testing.T) {
	fsys := fstest.MapFS{}
	loader := NewFSLoader(fsys)

	_, err := loader.Load("missing")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestFSLoader_LoadFromDisk(t *testing.T) {
	dir := t.TempDir()
	tmplFile := filepath.Join(dir, "hello.tmpl")
	if err := os.WriteFile(tmplFile, []byte("Hi {{.Name}}"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	fsys := os.DirFS(dir)
	r := New(NewFSLoader(fsys))

	out, err := r.Render("hello", map[string]any{"Name": "EPMP"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []byte("Hi EPMP")
	if !bytes.Equal(out, expected) {
		t.Errorf("got %q, want %q", out, expected)
	}
}

func TestMapLoader_LoadNotFound(t *testing.T) {
	loader := NewMapLoader(map[string]string{})

	_, err := loader.Load("missing")
	if err == nil {
		t.Fatal("expected error for missing template, got nil")
	}
}

func TestRender_ComplexTemplate(t *testing.T) {
	templates := map[string]string{
		"entity": `// Code generated by epmp-sdk. DO NOT EDIT.
package {{.Package}}

type {{.Name}} struct {
{{- range .Fields}}
	{{.Name}} {{.Type}} ` + "`json:\"{{.JSONName}}\"`" + `
{{- end}}
}`,
	}
	r := New(NewMapLoader(templates))

	data := map[string]any{
		"Package": "property",
		"Name":    "Property",
		"Fields": []map[string]any{
			{"Name": "ID", "Type": "string", "JSONName": "id"},
			{"Name": "Name", "Type": "string", "JSONName": "name"},
		},
	}

	out, err := r.Render("entity", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `// Code generated by epmp-sdk. DO NOT EDIT.
package property

type Property struct {
	ID string ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}`
	if string(out) != expected {
		t.Errorf("got:\n%s\nwant:\n%s", out, expected)
	}
}

func TestRender_ReturnsBytes(t *testing.T) {
	templates := map[string]string{
		"test": "content",
	}
	r := New(NewMapLoader(templates))

	out, err := r.Render("test", map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out == nil {
		t.Fatal("expected non-nil bytes")
	}
	if len(out) == 0 {
		t.Fatal("expected non-empty bytes")
	}
}

func TestRender_DoesNotWriteFiles(t *testing.T) {
	dir := t.TempDir()
	templates := map[string]string{
		"test": "content",
	}
	r := New(NewMapLoader(templates))

	_, err := r.Render("test", map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("unexpected error reading dir: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("renderer wrote files to disk: %d entries found", len(entries))
	}
}

func TestRender_LoaderErrorWrapped(t *testing.T) {
	loader := &failingLoader{}
	r := New(loader)

	_, err := r.Render("test", map[string]any{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var target errNotFound
	if !errors.As(err, &target) {
		t.Errorf("expected wrapped error to contain errNotFound, got: %v", err)
	}
}

type errNotFound struct{}

func (errNotFound) Error() string { return "not found" }

type failingLoader struct{}

func (failingLoader) Load(name string) (string, error) {
	return "", errNotFound{}
}
