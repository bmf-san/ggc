// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bmf-san/ggc/v5/git"
)

// Committer provides functionality for the commit command.
type Committer struct {
	gitClient    git.CommitWriter
	outputWriter io.Writer
	helper       *Helper
}

// NewCommitter creates a new Committer.
func NewCommitter(client git.CommitWriter) *Committer {
	c := &Committer{
		gitClient:    client,
		outputWriter: os.Stdout,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = c.outputWriter
	return c
}

// Commit executes the commit command with the given arguments.
func (c *Committer) Commit(args []string) {
	if len(args) == 0 {
		c.helper.ShowCommitHelp()
		return
	}

	switch args[0] {
	case "allow":
		c.handleAllowCommand(args[1:])
	case "amend":
		c.handleAmendCommand(args[1:])
	default:
		c.handleDefaultCommit(args)
	}
}

// handleAllowCommand handles the "allow" subcommand
func (c *Committer) handleAllowCommand(args []string) {
	if len(args) >= 1 && args[0] == "empty" {
		if err := c.gitClient.CommitAllowEmpty(); err != nil {
			_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
		}
		return
	}
	c.helper.ShowCommitHelp()
}

// handleAmendCommand handles the "amend" subcommand
func (c *Committer) handleAmendCommand(args []string) {
	switch {
	case len(args) == 0:
		if err := c.gitClient.CommitAmend(); err != nil {
			_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
		}
	case args[0] == "no-edit":
		if err := c.gitClient.CommitAmendNoEdit(); err != nil {
			_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
		}
	default:
		msg := strings.Join(args, " ")
		if err := c.gitClient.CommitAmendWithMessage(msg); err != nil {
			_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
		}
	}
}

// handleDefaultCommit handles regular commit with message
func (c *Committer) handleDefaultCommit(args []string) {
	msg := strings.Join(args, " ")
	if err := c.gitClient.Commit(msg); err != nil {
		_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
	}
}
