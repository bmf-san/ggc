package git

import (
	"fmt"
)

// OpError represents a git operation error.
// It follows the same pattern as net.OpError in the standard library.
type OpError struct {
	Op      string // Operation that failed (e.g., "checkout branch")
	Command string // Git command that was executed (e.g., "git checkout main")
	Err     error  // Underlying error
}

func (e *OpError) Error() string {
	if e.Command != "" {
		return fmt.Sprintf("git: %s failed: %s (command: %s)", e.Op, e.Err, e.Command)
	}
	return fmt.Sprintf("git: %s failed: %s", e.Op, e.Err)
}

// NewOpError creates a new OpError.
func NewOpError(op string, command string, err error) error {
	return &OpError{
		Op:      op,
		Command: command,
		Err:     err,
	}
}
