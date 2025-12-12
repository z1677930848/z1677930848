package files

import (
	"io/fs"
	"os"
	"path/filepath"
)

// File wraps a file system path with helpers similar to TeaGo/files.
type File struct {
	path string
	info fs.FileInfo
}

// NewFile creates a File wrapper.
func NewFile(path string) *File {
	return &File{path: path}
}

// Path returns the full path.
func (f *File) Path() string {
	return f.path
}

// Name returns the base name.
func (f *File) Name() string {
	if f.info != nil {
		return f.info.Name()
	}
	return filepath.Base(f.path)
}

// Exists reports whether the file or directory exists.
func (f *File) Exists() bool {
	_, err := os.Stat(f.path)
	return err == nil
}

// Mkdir creates the directory (all parents).
func (f *File) Mkdir() error {
	return os.MkdirAll(f.path, 0755)
}

// List lists files in the directory.
func (f *File) List() []*File {
	entries, err := os.ReadDir(f.path)
	if err != nil {
		return nil
	}
	result := make([]*File, 0, len(entries))
	for _, e := range entries {
		info, _ := e.Info()
		result = append(result, &File{
			path: filepath.Join(f.path, e.Name()),
			info: info,
		})
	}
	return result
}

// Delete removes the file or directory.
func (f *File) Delete() error {
	return os.RemoveAll(f.path)
}

// Move moves/renames the file.
func (f *File) Move(target string) error {
	return os.Rename(f.path, target)
}
