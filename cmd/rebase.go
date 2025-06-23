// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Rebaser handles rebase operations.
type Rebaser struct {
	outputWriter io.Writer
	helper       *Helper
	execCommand  func(string, ...string) *exec.Cmd
	inputReader  *bufio.Reader
}

// NewRebaser creates a new Rebaser instance.
func NewRebaser() *Rebaser {
	return &Rebaser{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		execCommand:  exec.Command,
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
	// Get current branch name
	branchCmd := r.execCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchOutput, err := branchCmd.CombinedOutput()
	if err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: failed to get current branch\n")
		return
	}
	currentBranch := strings.TrimSpace(string(branchOutput))

	// Get upstream or main branch
	upstreamCmd := r.execCommand("git", "rev-parse", "--abbrev-ref", fmt.Sprintf("%s@{upstream}", currentBranch))
	upstreamOutput, err := upstreamCmd.CombinedOutput()
	upstream := "main" // default to main if no upstream is set
	if err == nil {
		upstream = strings.TrimSpace(string(upstreamOutput))
	}

	// Get commit history for current branch only (from upstream/main to HEAD)
	cmd := r.execCommand("git", "log", "--oneline", "--reverse", fmt.Sprintf("%s..HEAD", upstream))
	output, err := cmd.CombinedOutput()
	if err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: failed to get git log\n")
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
	cmd = r.execCommand("git", "rebase", "-i", fmt.Sprintf("HEAD~%d", num))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: failed to start rebase: %v\n", err)
		return
	}

	if err := cmd.Wait(); err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: rebase failed\n")
		return
	}

	_, _ = fmt.Fprintf(r.outputWriter, "Rebase successful\n")
}
