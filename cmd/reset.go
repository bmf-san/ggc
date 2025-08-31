// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/bmf-san/ggc/v4/git"
)

// Resetter handles reset operations.
type Resetter struct {
	outputWriter io.Writer
	helper       *Helper
	gitClient    git.Clienter
}

// NewResetter creates a new Resetter instance.
func NewResetter() *Resetter {
	return &Resetter{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		gitClient:    getGitClient(),
	}
}

// NewResetterWithClient creates a new Resetter instance with a custom git client.
func NewResetterWithClient(client git.Clienter) *Resetter {
	return &Resetter{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		gitClient:    client,
	}
}

// Reset executes git reset commands.
func (r *Resetter) Reset(args []string) {
	if len(args) == 0 {
		// Default: reset to origin
		branch, err := r.gitClient.GetCurrentBranch()
		if err != nil {
			_, _ = fmt.Fprintf(r.outputWriter, "Error: failed to get current branch: %v\n", err)
			return
		}

		if err := r.gitClient.ResetHardAndClean(); err != nil {
			_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
			return
		}
		_, _ = fmt.Fprintf(r.outputWriter, "Reset to origin/%s successful\n", branch)
		return
	}

	switch args[0] {
	case "hard":
		if len(args) < 2 {
			_, _ = fmt.Fprintf(r.outputWriter, "Error: commit hash required for hard reset\n")
			r.helper.ShowResetHelp()
			return
		}

		commit := args[1]
		if err := r.gitClient.ResetHard(commit); err != nil {
			_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
			return
		}
		_, _ = fmt.Fprintf(r.outputWriter, "Reset to %s successful\n", commit)
		return
	default:
		r.helper.ShowResetHelp()
	}
}
