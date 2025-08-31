// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/bmf-san/ggc/v5/git"
)

// Pusher provides functionality for the push command.
type Pusher struct {
	gitClient    git.Clienter
	outputWriter io.Writer
	helper       *Helper
}

// NewPusher creates a new Pusher.
func NewPusher(client git.Clienter) *Pusher {
	p := &Pusher{
		gitClient:    client,
		outputWriter: os.Stdout,
		helper:       NewHelper(),
	}
	p.helper.outputWriter = p.outputWriter
	return p
}

// Push executes the push command with the given arguments.
func (p *Pusher) Push(args []string) {
	if len(args) == 0 {
		p.helper.ShowPushHelp()
		return
	}

	switch args[0] {
	case "current":
		if err := p.gitClient.Push(false); err != nil {
			_, _ = fmt.Fprintf(p.outputWriter, "Error: %v\n", err)
		}
	case "force":
		if err := p.gitClient.Push(true); err != nil {
			_, _ = fmt.Fprintf(p.outputWriter, "Error: %v\n", err)
		}
	default:
		p.helper.ShowPushHelp()
	}
}
