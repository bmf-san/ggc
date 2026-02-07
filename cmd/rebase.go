// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/bmf-san/ggc/v7/internal/prompt"
	"github.com/bmf-san/ggc/v7/pkg/git"
)

// Rebaser handles rebase operations.
type Rebaser struct {
	gitClient    git.RebaseOps
	outputWriter io.Writer
	helper       *Helper
	prompter     prompt.Prompter
}

// NewRebaser creates a new Rebaser instance.
func NewRebaser(client git.RebaseOps) *Rebaser {
	output := os.Stdout
	helper := NewHelper()
	helper.outputWriter = output
	return &Rebaser{
		gitClient:    client,
		outputWriter: output,
		helper:       helper,
		prompter:     prompt.New(os.Stdin, output),
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
	case "continue":
		r.handleRebaseContinue()
	case "abort":
		r.handleRebaseAbort()
	case "skip":
		r.handleRebaseSkip()
	default:
		r.handleStandardRebase(args[0])
	}
}

func (r *Rebaser) handleRebaseContinue() {
	if err := r.gitClient.RebaseContinue(); err != nil {
		WriteError(r.outputWriter, err)
		return
	}
	WriteLine(r.outputWriter, "Rebase successful")
}

func (r *Rebaser) handleRebaseAbort() {
	if err := r.gitClient.RebaseAbort(); err != nil {
		WriteError(r.outputWriter, err)
		return
	}
	WriteLine(r.outputWriter, "Rebase aborted")
}

func (r *Rebaser) handleRebaseSkip() {
	if err := r.gitClient.RebaseSkip(); err != nil {
		WriteError(r.outputWriter, err)
		return
	}
	WriteLine(r.outputWriter, "Rebase successful")
}

func (r *Rebaser) handleStandardRebase(ref string) {
	upstream := r.resolveUpstream(ref)
	if upstream == "" {
		return
	}
	if err := r.gitClient.Rebase(upstream); err != nil {
		WriteError(r.outputWriter, err)
		return
	}
	WriteLine(r.outputWriter, "Rebase successful")
}

func (r *Rebaser) resolveUpstream(ref string) string {
	if r.gitClient.RevParseVerify(ref) {
		return ref
	}
	try := "origin/" + ref
	if r.gitClient.RevParseVerify(try) {
		return try
	}
	WriteErrorf(r.outputWriter, "unknown ref '%s'", ref)
	return ""
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
		WriteError(r.outputWriter, err)
		return
	}
	WriteLine(r.outputWriter, "Rebase successful")
}

type rebaseCtx struct {
	currentBranch string
	upstream      string
	lines         []string
}

func (r *Rebaser) prepareRebaseContext() (rebaseCtx, bool) {
	currentBranch, err := r.gitClient.GetCurrentBranch()
	if err != nil {
		WriteError(r.outputWriter, err)
		return rebaseCtx{}, false
	}
	upstream, err := r.gitClient.GetUpstreamBranch(currentBranch)
	if err != nil {
		WriteError(r.outputWriter, err)
		return rebaseCtx{}, false
	}
	output, err := r.gitClient.LogOneline(upstream, "HEAD")
	if err != nil {
		WriteError(r.outputWriter, err)
		return rebaseCtx{}, false
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		WriteErrorf(r.outputWriter, "no commit history found")
		return rebaseCtx{}, false
	}
	return rebaseCtx{currentBranch: currentBranch, upstream: upstream, lines: lines}, true
}

func (r *Rebaser) printCommitChoices(currentBranch string, lines []string) {
	WriteLinef(r.outputWriter, "Current branch: %s", currentBranch)
	WriteLine(r.outputWriter, "Select number of commits to rebase (commits are shown from oldest to newest):")
	for i, line := range lines {
		WriteLinef(r.outputWriter, "  [%d] %s", i+1, line)
	}
}

func (r *Rebaser) promptRebaseCount(max int) (int, bool) {
	input, ok := ReadLine(r.prompter, r.outputWriter, "> ")
	if !ok || strings.TrimSpace(input) == "" {
		WriteErrorf(r.outputWriter, "operation canceled")
		return 0, false
	}
	num, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || num < 1 || num > max {
		WriteErrorf(r.outputWriter, "invalid number")
		return 0, false
	}
	return num, true
}
