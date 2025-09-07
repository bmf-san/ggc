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
	gitClient    git.Clienter
	outputWriter io.Writer
	helper       *Helper
}

// NewCommitter creates a new Committer.
func NewCommitter(client git.Clienter) *Committer {
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
	case "allow-empty":
		c.handleAllowEmpty()
	case "amend":
		c.handleAmend(args[1:])
	default:
		c.handleCommitMessage(args)
	}
}

// handleAllowEmpty creates an empty commit
func (c *Committer) handleAllowEmpty() {
	if err := c.gitClient.CommitAllowEmpty(); err != nil {
		_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
	}
}

// handleAmend processes amend variations
func (c *Committer) handleAmend(args []string) {
	switch {
	case len(args) == 0:
		if err := c.gitClient.CommitAmend(); err != nil {
			_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
		}
	case args[0] == "--no-edit":
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

// handleCommitMessage creates a normal commit with message
func (c *Committer) handleCommitMessage(args []string) {
	msg := strings.Join(args, " ")
	if err := c.gitClient.Commit(msg); err != nil {
		_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
	}
}
