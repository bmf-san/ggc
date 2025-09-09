// Package config provides minimal filesystem abstraction for testability
package config

import (
	"io"
	"os"
)

// FileReader abstracts file reading operations for testing
type FileReader interface {
	ReadFile(filename string) ([]byte, error)
	Stat(name string) (os.FileInfo, error)
}

// FileWriter abstracts file writing operations for testing
type FileWriter interface {
	WriteFile(filename string, data []byte, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	CreateTemp(dir, pattern string) (TempFile, error)
	Remove(name string) error
	Rename(oldpath, newpath string) error
	Chmod(name string, mode os.FileMode) error
}

// TempFile abstracts temporary file operations
type TempFile interface {
	io.WriteCloser
	Name() string
}

// OSFileOps implements both FileReader and FileWriter using real OS operations
type OSFileOps struct{}

// NewOSFileOps creates a new OS file operations implementation
func NewOSFileOps() *OSFileOps {
	return &OSFileOps{}
}

// FileReader methods
func (fs *OSFileOps) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (fs *OSFileOps) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// FileWriter methods
func (fs *OSFileOps) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

func (fs *OSFileOps) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (fs *OSFileOps) Remove(name string) error {
	return os.Remove(name)
}

func (fs *OSFileOps) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func (fs *OSFileOps) CreateTemp(dir, pattern string) (TempFile, error) {
	return os.CreateTemp(dir, pattern)
}

func (fs *OSFileOps) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}
