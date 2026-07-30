package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateDir_NewDirectory(t *testing.T) {
	fs := New()
	dir := filepath.Join(t.TempDir(), "newdir")

	if err := fs.CreateDir(dir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("expected directory, got file")
	}
}

func TestCreateDir_NestedDirectories(t *testing.T) {
	fs := New()
	dir := filepath.Join(t.TempDir(), "a", "b", "c")

	if err := fs.CreateDir(dir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("nested directory not created: %v", err)
	}
}

func TestCreateDir_AlreadyExists(t *testing.T) {
	fs := New()
	dir := t.TempDir()

	if err := fs.CreateDir(dir); err != nil {
		t.Fatalf("expected no error for existing dir, got: %v", err)
	}
}

func TestWriteFile_NewFile(t *testing.T) {
	fs := New()
	path := filepath.Join(t.TempDir(), "test.txt")
	data := []byte("hello world")

	if err := fs.WriteFile(path, data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("file not written: %v", err)
	}
	if string(content) != "hello world" {
		t.Errorf("got %q, want %q", content, data)
	}
}

func TestWriteFile_OverwriteExisting(t *testing.T) {
	fs := New()
	path := filepath.Join(t.TempDir(), "test.txt")

	if err := fs.WriteFile(path, []byte("old")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := fs.WriteFile(path, []byte("new")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("file not written: %v", err)
	}
	if string(content) != "new" {
		t.Errorf("got %q, want %q", content, "new")
	}
}

func TestWriteFile_EmptyData(t *testing.T) {
	fs := New()
	path := filepath.Join(t.TempDir(), "empty.txt")

	if err := fs.WriteFile(path, []byte{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("file not written: %v", err)
	}
	if info.Size() != 0 {
		t.Errorf("expected size 0, got %d", info.Size())
	}
}

func TestExists_True(t *testing.T) {
	fs := New()
	path := filepath.Join(t.TempDir(), "exists.txt")
	if err := os.WriteFile(path, []byte("data"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	exists, err := fs.Exists(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Error("expected true, got false")
	}
}

func TestExists_False(t *testing.T) {
	fs := New()
	path := filepath.Join(t.TempDir(), "nonexistent.txt")

	exists, err := fs.Exists(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exists {
		t.Error("expected false, got true")
	}
}

func TestExists_Directory(t *testing.T) {
	fs := New()
	path := t.TempDir()

	exists, err := fs.Exists(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Error("expected true for directory, got false")
	}
}

func TestCopy_File(t *testing.T) {
	fs := New()
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")

	if err := os.WriteFile(src, []byte("copy me"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := fs.Copy(src, dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("destination file not created: %v", err)
	}
	if string(content) != "copy me" {
		t.Errorf("got %q, want %q", content, "copy me")
	}
}

func TestCopy_OverwriteExisting(t *testing.T) {
	fs := New()
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")

	if err := os.WriteFile(src, []byte("new content"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := os.WriteFile(dst, []byte("old content"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := fs.Copy(src, dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("destination file not read: %v", err)
	}
	if string(content) != "new content" {
		t.Errorf("got %q, want %q", content, "new content")
	}
}

func TestCopy_SourceNotFound(t *testing.T) {
	fs := New()
	dir := t.TempDir()
	src := filepath.Join(dir, "missing.txt")
	dst := filepath.Join(dir, "dst.txt")

	err := fs.Copy(src, dst)
	if err == nil {
		t.Fatal("expected error for missing source, got nil")
	}
}

func TestCopy_PreservesPermissions(t *testing.T) {
	fs := New()
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")

	if err := os.WriteFile(src, []byte("data"), 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := fs.Copy(src, dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	srcInfo, _ := os.Stat(src)
	dstInfo, _ := os.Stat(dst)
	if srcInfo.Mode() != dstInfo.Mode() {
		t.Errorf("mode mismatch: src=%v dst=%v", srcInfo.Mode(), dstInfo.Mode())
	}
}

func TestRead_File(t *testing.T) {
	fs := New()
	path := filepath.Join(t.TempDir(), "read.txt")
	if err := os.WriteFile(path, []byte("read me"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	data, err := fs.Read(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != "read me" {
		t.Errorf("got %q, want %q", data, "read me")
	}
}

func TestRead_FileNotFound(t *testing.T) {
	fs := New()
	path := filepath.Join(t.TempDir(), "missing.txt")

	_, err := fs.Read(path)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestRead_EmptyFile(t *testing.T) {
	fs := New()
	path := filepath.Join(t.TempDir(), "empty.txt")
	if err := os.WriteFile(path, []byte{}, 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	data, err := fs.Read(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) != 0 {
		t.Errorf("expected empty data, got %d bytes", len(data))
	}
}

func TestJoin(t *testing.T) {
	result := Join("a", "b", "c.txt")
	expected := filepath.Join("a", "b", "c.txt")
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestWriteFile_ParentDirMissing(t *testing.T) {
	fs := New()
	path := filepath.Join(t.TempDir(), "nonexistent", "test.txt")

	err := fs.WriteFile(path, []byte("data"))
	if err == nil {
		t.Fatal("expected error for missing parent directory, got nil")
	}
}
