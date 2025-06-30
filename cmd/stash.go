// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/bmf-san/ggc/git"
)

// Stasher provides functionality for the stash command.
type Stasher struct {
	gitClient    git.Clienter
	outputWriter io.Writer
	helper       *Helper
}

// NewStasher creates a new Stasher.
func NewStasher() *Stasher {
	return NewStasherWithClient(git.NewClient())
}

// NewStasherWithClient creates a new Stasher with the specified git client.
func NewStasherWithClient(client git.Clienter) *Stasher {
	s := &Stasher{
		gitClient:    client,
		outputWriter: os.Stdout,
		helper:       NewHelper(),
	}
	s.helper.outputWriter = s.outputWriter
	return s
}

// Stash executes the stash command with the given arguments.
func (s *Stasher) Stash(args []string) {
	if len(args) == 0 {
		s.helper.ShowStashHelp()
		return
	}

	switch args[0] {
	case "trash":
		if err := s.gitClient.StashPullPop(); err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: %v\n", err)
			return
		}
		_, _ = fmt.Fprintln(s.outputWriter, "add . → stash done")
	default:
		s.helper.ShowStashHelp()
	}
}
