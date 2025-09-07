// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/bmf-san/ggc/v5/git"
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
	case "interactive":
		r.RebaseInteractive()
	default:
		r.helper.ShowRebaseHelp()
	}
}

// RebaseInteractive executes interactive rebase.
func (r *Rebaser) RebaseInteractive() {
	ctx, ok := r.prepareRebaseContext()
	if !ok {
		return
	}
	r.printCommitChoices(ctx.currentBranch, ctx.lines)
	num, ok := r.promptRebaseCount(len(ctx.lines))
	if !ok {
		return
	}
	if err := r.gitClient.RebaseInteractive(num); err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
		return
	}
	_, _ = fmt.Fprintf(r.outputWriter, "Rebase successful\n")
}

type rebaseCtx struct {
	currentBranch string
	upstream      string
	lines         []string
}

func (r *Rebaser) prepareRebaseContext() (rebaseCtx, bool) {
	currentBranch, err := r.gitClient.GetCurrentBranch()
	if err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
		return rebaseCtx{}, false
	}
	upstream, err := r.gitClient.GetUpstreamBranch(currentBranch)
	if err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
		return rebaseCtx{}, false
	}
	output, err := r.gitClient.LogOneline(upstream, "HEAD")
	if err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
		return rebaseCtx{}, false
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: no commit history found\n")
		return rebaseCtx{}, false
	}
	return rebaseCtx{currentBranch: currentBranch, upstream: upstream, lines: lines}, true
}

func (r *Rebaser) printCommitChoices(currentBranch string, lines []string) {
	_, _ = fmt.Fprintf(r.outputWriter, "Current branch: %s\n", currentBranch)
	_, _ = fmt.Fprintln(r.outputWriter, "Select number of commits to rebase (commits are shown from oldest to newest):")
	for i, line := range lines {
		_, _ = fmt.Fprintf(r.outputWriter, "  [%d] %s\n", i+1, line)
	}
}

func (r *Rebaser) promptRebaseCount(max int) (int, bool) {
	_, _ = fmt.Fprint(r.outputWriter, "> ")
	input, err := r.inputReader.ReadString('\n')
	if err != nil || strings.TrimSpace(input) == "" {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: operation canceled\n")
		return 0, false
	}
	num, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || num < 1 || num > max {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: invalid number\n")
		return 0, false
	}
	return num, true
}
