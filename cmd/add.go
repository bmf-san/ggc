// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/bmf-san/ggc/v4/git"
)

// Adder provides functionality for the add command.
type Adder struct {
	gitClient    git.Clienter
	outputWriter io.Writer
}

// NewAdder creates a new Adder.
func NewAdder() *Adder {
	return &Adder{
		gitClient:    getGitClient(),
		outputWriter: os.Stdout,
	}
}

// Add executes the add command with the given arguments.
func (a *Adder) Add(args []string) {
	if len(args) == 0 {
		_, _ = fmt.Fprintf(a.outputWriter, "Usage: ggc add <file> | ggc add interactive | ggc add -p\n")
		return
	}

	if len(args) == 1 && args[0] == "interactive" {
		if err := a.gitClient.AddInteractive(); err != nil {
			_, _ = fmt.Fprintf(a.outputWriter, "Error: %v\n", err)
		}
		return
	}

	if len(args) == 1 && args[0] == "-p" {
		if err := a.gitClient.AddInteractive(); err != nil {
			_, _ = fmt.Fprintf(a.outputWriter, "Error: %v\n", err)
		}
		return
	}

	if err := a.gitClient.Add(args...); err != nil {
		_, _ = fmt.Fprintf(a.outputWriter, "Error: %v\n", err)
	}
}
