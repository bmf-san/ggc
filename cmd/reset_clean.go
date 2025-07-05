// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// ResetCleaner handles reset and clean operations.
type ResetCleaner struct {
	outputWriter io.Writer
	helper       *Helper
	execCommand  func(string, ...string) *exec.Cmd
}

// NewResetCleaner creates a new ResetCleaner instance.
func NewResetCleaner() *ResetCleaner {
	return &ResetCleaner{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		execCommand:  exec.Command,
	}
}

// ResetClean executes git reset and clean commands.
func (r *ResetCleaner) ResetClean() {
	// Reset to HEAD
	resetCmd := r.execCommand("git", "reset", "--hard", "HEAD")
	if err := resetCmd.Run(); err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error resetting changes: reset error\n")
		return
	}

	// Clean untracked files
	cleanCmd := r.execCommand("git", "clean", "-fd")
	if err := cleanCmd.Run(); err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error cleaning untracked files: clean error\n")
		return
	}

	_, _ = fmt.Fprintf(r.outputWriter, "operation successful\n")
}
