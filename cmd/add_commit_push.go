// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// AddCommitPusher provides functionality for add-commit-push command.
type AddCommitPusher struct {
	execCommand  func(name string, arg ...string) *exec.Cmd
	inputReader  *bufio.Reader
	outputWriter io.Writer
}

// NewAddCommitPusher creates a new AddCommitPusher.
func NewAddCommitPusher() *AddCommitPusher {
	return &AddCommitPusher{
		execCommand:  exec.Command,
		inputReader:  bufio.NewReader(os.Stdin),
		outputWriter: os.Stdout,
	}
}

// AddCommitPush executes add, commit, and push in sequence.
func (a *AddCommitPusher) AddCommitPush() {
	// git add .
	addCmd := a.execCommand("git", "add", ".")
	addCmd.Stdout = a.outputWriter
	addCmd.Stderr = a.outputWriter
	if err := addCmd.Run(); err != nil {
		if _, err := fmt.Fprintf(a.outputWriter, "Error: failed to add all files: %v\n", err); err != nil {
			_ = err
		}
		return
	}
	// Enter commit message
	if _, err := fmt.Fprint(a.outputWriter, "\n\r"); err != nil {
		_ = err
	}
	if _, err := fmt.Fprint(a.outputWriter, "Enter commit message: "); err != nil {
		_ = err
	}
	msg, _ := a.inputReader.ReadString('\n')
	msg = strings.TrimSpace(msg)
	if msg == "" {
		if _, err := fmt.Fprintln(a.outputWriter, "Cancelled."); err != nil {
			_ = err
		}
		return
	}
	// git commit
	commitCmd := a.execCommand("git", "commit", "-m", msg)
	commitCmd.Stdout = a.outputWriter
	commitCmd.Stderr = a.outputWriter
	if err := commitCmd.Run(); err != nil {
		if _, err := fmt.Fprintf(a.outputWriter, "Error: failed to commit: %v\n", err); err != nil {
			_ = err
		}
		return
	}
	// Get current branch name
	branchCmd := a.execCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchOut, err := branchCmd.Output()
	if err != nil {
		if _, err := fmt.Fprintf(a.outputWriter, "Error: failed to get branch name: %v\n", err); err != nil {
			_ = err
		}
		return
	}
	branch := strings.TrimSpace(string(branchOut))
	// git push
	pushCmd := a.execCommand("git", "push", "origin", branch)
	pushCmd.Stdout = a.outputWriter
	pushCmd.Stderr = a.outputWriter
	if err := pushCmd.Run(); err != nil {
		if _, err := fmt.Fprintf(a.outputWriter, "Error: failed to push: %v\n", err); err != nil {
			_ = err
		}
		return
	}
	if _, err := fmt.Fprintln(a.outputWriter, "add→commit→push done"); err != nil {
		_ = err
	}
}
