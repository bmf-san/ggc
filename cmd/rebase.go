// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/bmf-san/ggc/v4/git"
)

// Rebaser handles rebase operations.
type Rebaser struct {
	gitClient    git.Clienter
	outputWriter io.Writer
	helper       *Helper
	inputReader  *bufio.Reader
}

// NewRebaser creates a new Rebaser instance.
func NewRebaser(client git.Clienter) *Rebaser {
	return &Rebaser{
		gitClient:    client,
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		inputReader:  bufio.NewReader(os.Stdin),
	}
}

// Rebase executes git rebase commands.
func (r *Rebaser) Rebase(args []string) {
	if len(args) == 0 {
		r.helper.ShowRebaseHelp()
		return
	}

	switch args[0] {
	case "-i", "--interactive":
		r.RebaseInteractive()
	case "interactive":
		_, _ = fmt.Fprintln(r.outputWriter, "Error: 'ggc rebase interactive' is no longer supported. Use 'ggc rebase -i' or 'ggc rebase --interactive'.")
		r.helper.ShowRebaseHelp()
	default:
		r.helper.ShowRebaseHelp()
	}
}

// RebaseInteractive executes interactive rebase.
func (r *Rebaser) RebaseInteractive() {
	// Get current branch name
	currentBranch, err := r.gitClient.GetCurrentBranch()
	if err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
		return
	}

	// Get upstream or main branch
	upstream, err := r.gitClient.GetUpstreamBranch(currentBranch)
	if err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
		return
	}

	// Get commit history for current branch only (from upstream/main to HEAD)
	output, err := r.gitClient.LogOneline(upstream, "HEAD")
	if err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
		return
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: no commit history found\n")
		return
	}

	_, _ = fmt.Fprintf(r.outputWriter, "Current branch: %s\n", currentBranch)
	_, _ = fmt.Fprintln(r.outputWriter, "Select number of commits to rebase (commits are shown from oldest to newest):")
	for i, line := range lines {
		_, _ = fmt.Fprintf(r.outputWriter, "  [%d] %s\n", i+1, line)
	}
	_, _ = fmt.Fprint(r.outputWriter, "> ")

	input, err := r.inputReader.ReadString('\n')
	if err != nil || strings.TrimSpace(input) == "" {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: operation cancelled\n")
		return
	}

	num, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || num < 1 || num > len(lines) {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: invalid number\n")
		return
	}

	// Start interactive rebase using HEAD~n format
	if err := r.gitClient.RebaseInteractive(num); err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
		return
	}
	_, _ = fmt.Fprintf(r.outputWriter, "Rebase successful\n")
}
