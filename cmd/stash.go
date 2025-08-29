// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bmf-san/ggc/v4/git"
)

// Stasher handles stash operations.
type Stasher struct {
	gitClient    git.Clienter
	outputWriter io.Writer
	helper       *Helper
}

// NewStasher creates a new Stasher instance.
func NewStasher() *Stasher {
	return &Stasher{
		gitClient:    git.NewClient(),
		outputWriter: os.Stdout,
		helper:       NewHelper(),
	}
}

// Stash executes git stash commands.
func (s *Stasher) Stash(args []string) {
	if len(args) == 0 {
		// Default stash operation - stash current changes
		if err := s.gitClient.Stash(); err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: %v\n", err)
			return
		}
		return
	}

	switch args[0] {
	case "list":
		// List all stashes
		output, err := s.gitClient.StashList()
		if err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: %v\n", err)
			return
		}
		if len(strings.TrimSpace(output)) == 0 {
			_, _ = fmt.Fprintf(s.outputWriter, "No stashes found\n")
			return
		}
		_, _ = fmt.Fprintf(s.outputWriter, "%s", output)

	case "show":
		// Show the changes recorded in the stash
		var stash string
		if len(args) > 1 {
			stash = args[1]
		}
		if err := s.gitClient.StashShow(stash); err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: %v\n", err)
			return
		}

	case "apply":
		// Apply the stash without removing it
		var stash string
		if len(args) > 1 {
			stash = args[1]
		}
		if err := s.gitClient.StashApply(stash); err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: %v\n", err)
			return
		}

	case "pop":
		// Apply and remove the latest stash
		var stash string
		if len(args) > 1 {
			stash = args[1]
		}
		if err := s.gitClient.StashPop(stash); err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: %v\n", err)
			return
		}

	case "drop":
		// Drop the specified stash
		var stash string
		if len(args) > 1 {
			stash = args[1]
		}
		if err := s.gitClient.StashDrop(stash); err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: %v\n", err)
			return
		}

	case "clear":
		// Remove all stashes
		if err := s.gitClient.StashClear(); err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: %v\n", err)
			return
		}

	default:
		s.helper.ShowStashHelp()
		return
	}
}
