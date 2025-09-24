package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/bmf-san/ggc/v6/git"
)

// Differ handles git diff operations.
type Differ struct {
	gitClient    git.DiffReader
	outputWriter io.Writer
	helper       *Helper
}

// NewDiffer creates a new Differ instance.
func NewDiffer(client git.DiffReader) *Differ {
	return &Differ{
		gitClient:    client,
		outputWriter: os.Stdout,
		helper:       NewHelper(),
	}
}

// Diff executes git diff with the given arguments.
func (d *Differ) Diff(args []string) {
	var output string
	var err error

	if len(args) == 0 {
		output, err = d.gitClient.DiffHead()
	} else {
		switch args[0] {
		case "unstaged":
			output, err = d.gitClient.Diff()
		case "staged":
			output, err = d.gitClient.DiffStaged()
		case "head":
			output, err = d.gitClient.DiffHead()
		default:
			d.helper.ShowDiffHelp()
			return
		}
	}

	if err != nil {
		_, _ = fmt.Fprintf(d.outputWriter, "Error: %v\n", err)
		return
	}

	_, _ = fmt.Fprint(d.outputWriter, output)
}
