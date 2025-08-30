// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/bmf-san/ggc/v4/git"
)

// Committer provides functionality for the commit command.
type Committer struct {
	gitClient    git.Clienter
	outputWriter io.Writer
	helper       *Helper
	execCommand  func(name string, arg ...string) *exec.Cmd
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
		execCommand:  exec.Command,
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

	// Flag-style handling
	// --allow-empty
	for _, a := range args {
		if a == "--allow-empty" {
			if err := c.gitClient.CommitAllowEmpty(); err != nil {
				_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
			}
			return
		}
	}

	// --amend [--no-edit] [message...]
	amend := false
	noEdit := false
	msgParts := []string{}
	for _, a := range args {
		switch a {
		case "--amend":
			amend = true
		case "--no-edit":
			noEdit = true
		default:
			// treat as part of message if not a flag
			if !strings.HasPrefix(a, "-") {
				msgParts = append(msgParts, a)
			}
		}
	}

	if amend {
		var cmd *exec.Cmd
		switch {
		case noEdit:
			cmd = c.execCommand("git", "commit", "--amend", "--no-edit")
			cmd.Stdout = c.outputWriter
			cmd.Stderr = c.outputWriter
		case len(msgParts) > 0:
			cmd = c.execCommand("git", "commit", "--amend", "-m", strings.Join(msgParts, " "))
			cmd.Stdout = c.outputWriter
			cmd.Stderr = c.outputWriter
		default:
			cmd = c.execCommand("git", "commit", "--amend")
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		if err := cmd.Run(); err != nil {
			_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
		}
		return
	}

	// Default: normal commit with message
	msg := strings.Join(args, " ")
	cmd := c.execCommand("git", "commit", "-m", msg)
	cmd.Stdout = c.outputWriter
	cmd.Stderr = c.outputWriter
	if err := cmd.Run(); err != nil {
		_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
	}
}
