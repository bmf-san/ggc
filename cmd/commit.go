// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/bmf-san/ggc/git"
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
	case "tmp":
		if err := c.gitClient.CommitTmp(); err != nil {
			_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
		}
	default:
		c.helper.ShowCommitHelp()
	}
}

// ShowCommitHelp displays help for the commit command.
func ShowCommitHelp() {
	fmt.Println("Usage: ggc commit allow-empty | ggc commit tmp")
}
