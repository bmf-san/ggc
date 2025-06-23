package git

import (
	"errors"
	"fmt"
)

// Error represents a git operation error.
type Error struct {
	Op      string // Operation that failed
	Command string // Git command that was executed
	Err     error  // Underlying error
}

func (e *Error) Error() string {
	if e.Command != "" {
		return fmt.Sprintf("git: %s failed: %s (command: %s)", e.Op, e.Err, e.Command)
	}
	return fmt.Sprintf("git: %s failed: %s", e.Op, e.Err)
}

// NewError creates a new Error.
func NewError(op string, command string, err error) error {
	return &Error{
		Op:      op,
		Command: command,
		Err:     err,
	}
}

var (
	// ErrCleanFiles is returned when cleaning files fails.
	ErrCleanFiles = errors.New("failed to clean files")
	// ErrCleanDirs is returned when cleaning directories fails.
	ErrCleanDirs = errors.New("failed to clean directories")
)
