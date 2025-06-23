// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Stasher handles stash operations.
type Stasher struct {
	outputWriter io.Writer
	helper       *Helper
	execCommand  func(string, ...string) *exec.Cmd
}

// NewStasher creates a new Stasher instance.
func NewStasher() *Stasher {
	return &Stasher{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		execCommand:  exec.Command,
	}
}

// Stash executes git stash commands.
func (s *Stasher) Stash(args []string) {
	if len(args) == 0 {
		s.helper.ShowStashHelp()
		return
	}

	var cmd *exec.Cmd
	switch args[0] {
	case "trash":
		// Drop the latest stash
		cmd = s.execCommand("git", "stash", "drop")
	default:
		s.helper.ShowStashHelp()
		return
	}

	if err := cmd.Run(); err != nil {
		_, _ = fmt.Fprintf(s.outputWriter, "Error: no stash found\n")
		return
	}

	_, _ = fmt.Fprintf(s.outputWriter, "Dropped refs/stash@{0}\n")
}
