// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bmf-san/ggc/v6/git"
)

// Stasher handles stash operations.
type Stasher struct {
	gitClient    git.StashOps
	outputWriter io.Writer
	helper       *Helper
}

// NewStasher creates a new Stasher instance.
func NewStasher(client git.StashOps) *Stasher {
	return &Stasher{
		gitClient:    client,
		outputWriter: os.Stdout,
		helper:       NewHelper(),
	}
}

// Stash executes git stash commands.
func (s *Stasher) Stash(args []string) {
	if len(args) == 0 {
		s.stashDefault()
		return
	}

	switch args[0] {
	case "list":
		s.stashList()
	case "show":
		s.stashShow(args)
	case "apply":
		s.stashApply(args)
	case "pop":
		s.stashPop(args)
	case "drop":
		s.stashDrop(args)
	case "clear":
		s.stashClear()
	default:
		s.helper.ShowStashHelp()
	}
}

// stashDefault performs default stash operation - stash current changes
func (s *Stasher) stashDefault() {
	if err := s.gitClient.Stash(); err != nil {
		_, _ = fmt.Fprintf(s.outputWriter, "Error: %v\n", err)
	}
}

// stashList lists all stashes
func (s *Stasher) stashList() {
	output, err := s.gitClient.StashList()
	if err != nil {
		_, _ = fmt.Fprintf(s.outputWriter, "Error: %v\n", err)
		return
	}
	if strings.TrimSpace(output) == "" {
		_, _ = fmt.Fprintf(s.outputWriter, "No stashes found\n")
		return
	}
	_, _ = fmt.Fprintf(s.outputWriter, "%s", output)
}

// stashShow shows the changes recorded in the stash
func (s *Stasher) stashShow(args []string) {
	var stash string
	if len(args) > 1 {
		stash = args[1]
	}
	if err := s.gitClient.StashShow(stash); err != nil {
		_, _ = fmt.Fprintf(s.outputWriter, "Error: %v\n", err)
	}
}

// stashApply applies the stash without removing it
func (s *Stasher) stashApply(args []string) {
	var stash string
	if len(args) > 1 {
		stash = args[1]
	}
	if err := s.gitClient.StashApply(stash); err != nil {
		_, _ = fmt.Fprintf(s.outputWriter, "Error: %v\n", err)
	}
}

// stashPop applies and removes the latest stash
func (s *Stasher) stashPop(args []string) {
	var stash string
	if len(args) > 1 {
		stash = args[1]
	}
	if err := s.gitClient.StashPop(stash); err != nil {
		_, _ = fmt.Fprintf(s.outputWriter, "Error: %v\n", err)
	}
}

// stashDrop drops the specified stash
func (s *Stasher) stashDrop(args []string) {
	var stash string
	if len(args) > 1 {
		stash = args[1]
	}
	if err := s.gitClient.StashDrop(stash); err != nil {
		_, _ = fmt.Fprintf(s.outputWriter, "Error: %v\n", err)
	}
}

// stashClear removes all stashes
func (s *Stasher) stashClear() {
	if err := s.gitClient.StashClear(); err != nil {
		_, _ = fmt.Fprintf(s.outputWriter, "Error: %v\n", err)
	}
}
