package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Differ handles git diff operations.
type Differ struct {
	outputWriter io.Writer
	helper       *Helper
	execCommand  func(string, ...string) *exec.Cmd
}

// NewDiffer creates a new Differ instance.
func NewDiffer() *Differ {
	return &Differ{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		execCommand:  exec.Command,
	}
}

// Diff executes git diff with the given arguments.
func (d *Differ) Diff(args []string) {
	var cmd *exec.Cmd
	if len(args) == 0 {
		cmd = d.execCommand("git", "diff", "HEAD")
	} else {
		switch args[0] {
		case "unstaged":
			cmd = d.execCommand("git", "diff")
		case "staged":
			cmd = d.execCommand("git", "diff", "--staged")
		default:
			d.helper.ShowDiffHelp()
			return
		}
	}

	cmd.Stdout = d.outputWriter
	cmd.Stderr = d.outputWriter
	if err := cmd.Run(); err != nil {
		_, _ = fmt.Fprintf(d.outputWriter, "Error: %v\n", err)
	}
}
