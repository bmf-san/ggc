package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/bmf-san/ggc/v6/git"
)

// Fetcher handles git fetch operations.
type Fetcher struct {
	gitClient    git.FetchOps
	outputWriter io.Writer
	helper       *Helper
}

// NewFetcher creates a new Fetcher instance.
func NewFetcher(client git.FetchOps) *Fetcher {
	return &Fetcher{
		gitClient:    client,
		outputWriter: os.Stdout,
		helper:       NewHelper(),
	}
}

// Fetch executes git fetch with the given arguments.
func (f *Fetcher) Fetch(args []string) {
	if len(args) == 0 {
		f.helper.ShowFetchHelp()
		return
	}

	switch args[0] {
	case "prune":
		if err := f.gitClient.Fetch(true); err != nil {
			_, _ = fmt.Fprintf(f.outputWriter, "Error: %v\n", err)
		}
	default:
		f.helper.ShowFetchHelp()
		return
	}
}
