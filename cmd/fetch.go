package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Fetcher handles git fetch operations.
type Fetcher struct {
	outputWriter io.Writer
	helper       *Helper
	execCommand  func(string, ...string) *exec.Cmd
}

// NewFetcher creates a new Fetcher instance.
func NewFetcher() *Fetcher {
	return &Fetcher{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		execCommand:  exec.Command,
	}
}

// Fetch executes git fetch with the given arguments.
func (f *Fetcher) Fetch(args []string) {
	if len(args) == 0 {
		f.helper.ShowFetchHelp()
		return
	}

	var cmd *exec.Cmd
	switch args[0] {
	case "--prune":
		cmd = f.execCommand("git", "fetch", "--prune")
	default:
		f.helper.ShowFetchHelp()
		return
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		_, _ = fmt.Fprintf(f.outputWriter, "Error: %s\n", err)
		return
	}

	_, _ = fmt.Fprintf(f.outputWriter, "%s", output)
}
