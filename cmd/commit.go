// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bmf-san/ggc/v4/git"
)

// Committer provides functionality for the commit command.
type Committer struct {
	gitClient    git.Clienter
	outputWriter io.Writer
	helper       *Helper
}

// NewCommitter creates a new Committer.
func NewCommitter() *Committer {
	return NewCommitterWithClient(git.NewClient())
}

// NewCommitterWithClient creates a new Committer with the specified git client.
func NewCommitterWithClient(client git.Clienter) *Committer {
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
		if err := c.gitClient.CommitAllowEmpty(); err != nil {
			_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
		}
	case "amend":
		if len(args) == 1 {
			if err := c.gitClient.CommitAmend(); err != nil {
				_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
			}
		} else if args[1] == "--no-edit" {
			if err := c.gitClient.CommitAmendNoEdit(); err != nil {
				_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
			}
		} else {
			// Join all arguments after "amend" as the commit message
			msg := strings.Join(args[1:], " ")
			if err := c.gitClient.CommitAmendWithMessage(msg); err != nil {
				_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
			}
		}
	default:
		// Handle normal commit with message
		msg := strings.Join(args, " ")
		if err := c.gitClient.Commit(msg); err != nil {
			_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
		}
	}
}
