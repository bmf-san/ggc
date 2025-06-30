// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/bmf-san/ggc/git"
)

// Resetter provides functionality for the reset command.
type Resetter struct {
	gitClient    git.Clienter
	outputWriter io.Writer
	helper       *Helper
}

// NewResetter creates a new Resetter.
func NewResetter() *Resetter {
	return NewResetterWithClient(git.NewClient())
}

// NewResetterWithClient creates a new Resetter with the specified git client.
func NewResetterWithClient(client git.Clienter) *Resetter {
	return &Resetter{
		gitClient:    client,
		outputWriter: os.Stdout,
		helper:       NewHelper(),
	}
}

// Reset executes the reset command.
func (r *Resetter) Reset() {
	if err := r.gitClient.ResetHardAndClean(); err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
	}
}
