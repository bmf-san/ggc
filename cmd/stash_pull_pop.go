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

// StashPullPop executes stash, pull, and pop commands in sequence using --autostash.
func (s *StashPullPopper) StashPullPop() {
	pullCmd := s.execCommand("git", "pull", "--autostash")
	if err := pullCmd.Run(); err != nil {
		_, _ = fmt.Fprintf(s.outputWriter, "Error pulling changes with autostash: %v\n", err)
		return
	}

	_, _ = fmt.Fprintf(s.outputWriter, "operation successful with autostash\n")
}
