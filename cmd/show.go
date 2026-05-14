package cmd

import (
	"io"
	"os"

	"github.com/bmf-san/ggc/v8/internal/git"
)

// Shower handles git show operations.
type Shower struct {
	gitClient    git.ShowOps
	outputWriter io.Writer
	helper       *Helper
}

// NewShower creates a new Shower instance.
func NewShower(client git.ShowOps) *Shower {
	return &Shower{
		gitClient:    client,
		outputWriter: os.Stdout,
		helper:       NewHelper(),
	}
}

// Show executes git show with the given arguments. With no arguments,
// it shows the HEAD commit. The first argument may be "help" to print
// usage information without invoking git.
func (s *Shower) Show(args []string) {
	if len(args) > 0 && args[0] == "help" {
		s.helper.ShowShowHelp()
		return
	}
	if err := s.gitClient.Show(args); err != nil {
		WriteError(s.outputWriter, err)
	}
}
