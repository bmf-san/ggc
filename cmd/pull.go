// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/bmf-san/ggc/git"
)

// Puller provides functionality for the pull command.
type Puller struct {
	gitClient    git.Clienter
	outputWriter io.Writer
	helper       *Helper
}

// NewPuller creates a new Puller.
func NewPuller() *Puller {
	return NewPullerWithClient(git.NewClient())
}

// NewPullerWithClient creates a new Puller with the specified git client.
func NewPullerWithClient(client git.Clienter) *Puller {
	p := &Puller{
		gitClient:    client,
		outputWriter: os.Stdout,
		helper:       NewHelper(),
	}
	p.helper.outputWriter = p.outputWriter
	return p
}

// Pull executes the pull command with the given arguments.
func (p *Puller) Pull(args []string) {
	if len(args) == 0 {
		p.helper.ShowPullHelp()
		return
	}

	switch args[0] {
	case "current":
		if err := p.gitClient.Pull(false); err != nil {
			_, _ = fmt.Fprintf(p.outputWriter, "Error: %v\n", err)
		}
	case "rebase":
		if err := p.gitClient.Pull(true); err != nil {
			_, _ = fmt.Fprintf(p.outputWriter, "Error: %v\n", err)
		}
	default:
		p.helper.ShowPullHelp()
	}
}
