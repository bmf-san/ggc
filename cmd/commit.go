// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/bmf-san/ggc/git"
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

	switch args[0] {
	case "allow-empty":
		if err := c.gitClient.CommitAllowEmpty(); err != nil {
			_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
		}
	case "tmp":
		if err := c.gitClient.CommitTmp(); err != nil {
			_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
		}
	case "amend":
		var cmd *exec.Cmd

		if len(args) == 1 {
			cmd = c.execCommand("git", "commit", "--amend")
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		} else if args[1] == "--no-edit" {
			cmd = c.execCommand("git", "commit", "--amend", "--no-edit")
			cmd.Stdout = c.outputWriter
			cmd.Stderr = c.outputWriter
		} else {
			cmd = c.execCommand("git", "commit", "--amend", "-m", args[1])
			cmd.Stdout = c.outputWriter
			cmd.Stderr = c.outputWriter
		}

		if err := cmd.Run(); err != nil {
			_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
		}
	default:
		// Handle normal commit with message
		cmd := c.execCommand("git", "commit", "-m", args[0])
		cmd.Stdout = c.outputWriter
		cmd.Stderr = c.outputWriter
		if err := cmd.Run(); err != nil {
			_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
		}
	}
}

// ShowCommitHelp displays help for the commit command.
func ShowCommitHelp() {
	fmt.Println("Usage: ggc commit <message> | ggc commit amend <message> | ggc commit allow-empty | ggc commit tmp")
}
