// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// StashPullPopper handles stash-pull-pop operations.
type StashPullPopper struct {
	outputWriter io.Writer
	helper       *Helper
	execCommand  func(string, ...string) *exec.Cmd
}

// NewStashPullPopper creates a new StashPullPopper instance.
func NewStashPullPopper() *StashPullPopper {
	return &StashPullPopper{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		execCommand:  exec.Command,
	}
}

// StashPullPop executes stash, pull, and pop commands in sequence.
func (s *StashPullPopper) StashPullPop() {
	// Stash changes
	stashCmd := s.execCommand("git", "stash")
	if err := stashCmd.Run(); err != nil {
		_, _ = fmt.Fprintf(s.outputWriter, "Error stashing changes: stash error\n")
		return
	}

	// Pull changes
	pullCmd := s.execCommand("git", "pull")
	if err := pullCmd.Run(); err != nil {
		_, _ = fmt.Fprintf(s.outputWriter, "Error pulling changes: pull error\n")
		return
	}

	// Pop stashed changes
	popCmd := s.execCommand("git", "stash", "pop")
	if err := popCmd.Run(); err != nil {
		_, _ = fmt.Fprintf(s.outputWriter, "Error popping stashed changes: pop error\n")
		return
	}

	_, _ = fmt.Fprintf(s.outputWriter, "operation successful\n")
}
