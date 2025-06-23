// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/bmf-san/ggc/git"
)

// PullRebasePusher provides functionality for the pull-rebase-push command.
type PullRebasePusher struct {
	gitClient    git.Clienter
	outputWriter io.Writer
}

// NewPullRebasePusher creates a new PullRebasePusher.
func NewPullRebasePusher() *PullRebasePusher {
	return NewPullRebasePusherWithClient(git.NewClient())
}

// NewPullRebasePusherWithClient creates a new PullRebasePusher with the specified git client.
func NewPullRebasePusherWithClient(client git.Clienter) *PullRebasePusher {
	return &PullRebasePusher{
		gitClient:    client,
		outputWriter: os.Stdout,
	}
}

// PullRebasePush executes the pull-rebase-push command.
func (p *PullRebasePusher) PullRebasePush() {
	// Pull with rebase
	if err := p.gitClient.Pull(true); err != nil {
		_, _ = fmt.Fprintf(p.outputWriter, "Error: %v\n", err)
		return
	}

	// Push to remote
	if err := p.gitClient.Push(false); err != nil {
		_, _ = fmt.Fprintf(p.outputWriter, "Error: %v\n", err)
		return
	}

	_, _ = fmt.Fprintln(p.outputWriter, "pull→rebase→push completed")
}
