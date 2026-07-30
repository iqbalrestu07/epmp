package filesystem

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Filesystem is the interface for filesystem operations.
type Filesystem interface {
	// CreateDir creates a directory and all parent directories.
	// It is idempotent — no error if the directory already exists.
	CreateDir(path string) error

	// WriteFile writes data to the given path.
	// Parent directories must already exist.
	// It overwrites the file if it already exists.
	WriteFile(path string, data []byte) error

	// Exists returns true if the given path exists.
	Exists(path string) (bool, error)

	// Copy copies a file from src to dst.
	// Parent directories of dst must already exist.
	// If dst already exists, it will be overwritten.
	Copy(src, dst string) error

	// Read reads the entire content of a file.
	Read(path string) ([]byte, error)
}

// osFilesystem implements Filesystem using the operating system.
type osFilesystem struct{}

// New creates a Filesystem backed by the OS.
func New() Filesystem {
	return &osFilesystem{}
}

func (f *osFilesystem) CreateDir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("filesystem: create dir %q: %w", path, err)
	}
	return nil
}

func (f *osFilesystem) WriteFile(path string, data []byte) error {
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("filesystem: write file %q: %w", path, err)
	}
	return nil
}

func (f *osFilesystem) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("filesystem: stat %q: %w", path, err)
}

func (f *osFilesystem) Copy(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("filesystem: open src %q: %w", src, err)
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("filesystem: stat src %q: %w", src, err)
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("filesystem: open dst %q: %w", dst, err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("filesystem: copy %q to %q: %w", src, dst, err)
	}
	return nil
}

func (f *osFilesystem) Read(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("filesystem: read %q: %w", path, err)
	}
	return data, nil
}

// Join joins path elements using the OS-specific separator.
func Join(elements ...string) string {
	return filepath.Join(elements...)
}
