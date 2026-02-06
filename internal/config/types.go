package config

import (
	"fmt"
	"os"
)

// TempFile interface for temporary file operations
type TempFile interface {
	Write([]byte) (int, error)
	Close() error
	Name() string
}

// FileOps interface for file operations (for testability)
type FileOps interface {
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
	Stat(name string) (os.FileInfo, error)
	MkdirAll(path string, perm os.FileMode) error
	CreateTemp(dir, pattern string) (TempFile, error)
	Remove(name string) error
	Rename(oldpath, newpath string) error
	Chmod(name string, mode os.FileMode) error
}

// OSFileOps implements FileOps using real OS operations
type OSFileOps struct{}

// ReadFile reads a file from the filesystem
func (OSFileOps) ReadFile(filename string) ([]byte, error) { return os.ReadFile(filename) }

// WriteFile writes data to a file
func (OSFileOps) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

// Stat returns file information
func (OSFileOps) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }

// MkdirAll creates directories recursively
func (OSFileOps) MkdirAll(path string, perm os.FileMode) error { return os.MkdirAll(path, perm) }

// CreateTemp creates a temporary file
func (OSFileOps) CreateTemp(dir, pattern string) (TempFile, error) {
	return os.CreateTemp(dir, pattern)
}

// Remove removes a file
func (OSFileOps) Remove(name string) error { return os.Remove(name) }

// Rename renames a file
func (OSFileOps) Rename(oldpath, newpath string) error { return os.Rename(oldpath, newpath) }

// Chmod changes file permissions
func (OSFileOps) Chmod(name string, mode os.FileMode) error { return os.Chmod(name, mode) }

// KeybindingsConfig represents a configuration section with keybindings
type KeybindingsConfig struct {
	Keybindings map[string]interface{} `yaml:"keybindings,omitempty"`
}

// AliasType represents the type of alias
type AliasType int

const (
	// SimpleAlias represents a simple string alias
	SimpleAlias AliasType = iota
	// SequenceAlias represents an array string alias
	SequenceAlias
)

// ParsedAlias represents a parsed alias with its type and commands
type ParsedAlias struct {
	Type             AliasType
	Commands         []string
	Placeholders     map[string]struct{} // Track which placeholders are used
	MaxPositionalArg int                 // Highest positional argument index (-1 if none)
}

// ValidationError creates a new error manager for validation operations
type ValidationError struct {
	Field   string
	Value   any
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid value for '%s': %v (%s)", e.Field, e.Value, e.Message)
}
