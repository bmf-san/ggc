// Package config provides filesystem abstraction for testability
package config

import (
	"io"
	"os"
)

// FileSystem abstracts file system operations for testing
type FileSystem interface {
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Remove(name string) error
	Rename(oldpath, newpath string) error
	Stat(name string) (os.FileInfo, error)
	CreateTemp(dir, pattern string) (File, error)
	Chmod(name string, mode os.FileMode) error
}

// File abstracts file operations
type File interface {
	io.WriteCloser
	Name() string
}

// OSFileSystem implements FileSystem using real OS operations
type OSFileSystem struct{}

// NewOSFileSystem creates a new OS filesystem implementation
func NewOSFileSystem() *OSFileSystem {
	return &OSFileSystem{}
}

func (fs *OSFileSystem) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (fs *OSFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

func (fs *OSFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (fs *OSFileSystem) Remove(name string) error {
	return os.Remove(name)
}

func (fs *OSFileSystem) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func (fs *OSFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (fs *OSFileSystem) CreateTemp(dir, pattern string) (File, error) {
	return os.CreateTemp(dir, pattern)
}

func (fs *OSFileSystem) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}
