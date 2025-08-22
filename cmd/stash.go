// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Stasher handles stash operations.
type Stasher struct {
	outputWriter io.Writer
	helper       *Helper
	execCommand  func(string, ...string) *exec.Cmd
}

// NewStasher creates a new Stasher instance.
func NewStasher() *Stasher {
	return &Stasher{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		execCommand:  exec.Command,
	}
}

// Stash executes git stash commands.
func (s *Stasher) Stash(args []string) {
	if len(args) == 0 {
		// Default stash operation - stash current changes
		cmd := s.execCommand("git", "stash")
		if err := cmd.Run(); err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: no changes to stash\n")
			return
		}
		_, _ = fmt.Fprintf(s.outputWriter, "Saved working directory and index state\n")
		return
	}

	var cmd *exec.Cmd
	switch args[0] {
	case "list":
		// List all stashes
		cmd = s.execCommand("git", "stash", "list")
		output, err := cmd.Output()
		if err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: failed to list stashes\n")
			return
		}
		if len(output) == 0 {
			_, _ = fmt.Fprintf(s.outputWriter, "No stashes found\n")
			return
		}
		_, _ = fmt.Fprintf(s.outputWriter, "%s", output)

	case "show":
		// Show the changes recorded in the stash
		if len(args) > 1 {
			cmd = s.execCommand("git", "stash", "show", args[1])
		} else {
			cmd = s.execCommand("git", "stash", "show")
		}
		output, err := cmd.Output()
		if err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: no stash found\n")
			return
		}
		_, _ = fmt.Fprintf(s.outputWriter, "%s", output)

	case "apply":
		// Apply the stash without removing it
		if len(args) > 1 {
			cmd = s.execCommand("git", "stash", "apply", args[1])
		} else {
			cmd = s.execCommand("git", "stash", "apply")
		}
		if err := cmd.Run(); err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: no stash found\n")
			return
		}
		_, _ = fmt.Fprintf(s.outputWriter, "Applied stash\n")

	case "pop":
		// Apply and remove the latest stash
		if len(args) > 1 {
			cmd = s.execCommand("git", "stash", "pop", args[1])
		} else {
			cmd = s.execCommand("git", "stash", "pop")
		}
		if err := cmd.Run(); err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: no stash found\n")
			return
		}
		_, _ = fmt.Fprintf(s.outputWriter, "Applied and dropped stash\n")

	case "drop":
		// Drop the specified stash
		if len(args) > 1 {
			cmd = s.execCommand("git", "stash", "drop", args[1])
		} else {
			cmd = s.execCommand("git", "stash", "drop")
		}
		if err := cmd.Run(); err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: no stash found\n")
			return
		}
		_, _ = fmt.Fprintf(s.outputWriter, "Dropped stash\n")

	case "branch":
		// Create a new branch from a stash
		if len(args) < 2 {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: branch name required\n")
			return
		}
		if len(args) > 2 {
			cmd = s.execCommand("git", "stash", "branch", args[1], args[2])
		} else {
			cmd = s.execCommand("git", "stash", "branch", args[1])
		}
		if err := cmd.Run(); err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: failed to create branch from stash\n")
			return
		}
		_, _ = fmt.Fprintf(s.outputWriter, "Created branch '%s' from stash\n", args[1])

	case "push":
		// Save changes to a new stash (same as save)
		gitArgs := []string{"stash", "push"}
		if len(args) > 1 {
			gitArgs = append(gitArgs, args[1:]...)
		}
		cmd = s.execCommand("git", gitArgs...)
		if err := cmd.Run(); err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: no changes to stash\n")
			return
		}
		_, _ = fmt.Fprintf(s.outputWriter, "Saved working directory and index state\n")

	case "save":
		// Save changes to a new stash with message
		gitArgs := []string{"stash", "save"}
		if len(args) > 1 {
			gitArgs = append(gitArgs, args[1:]...)
		}
		cmd = s.execCommand("git", gitArgs...)
		if err := cmd.Run(); err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: no changes to stash\n")
			return
		}
		_, _ = fmt.Fprintf(s.outputWriter, "Saved working directory and index state\n")

	case "clear":
		// Remove all stashes
		cmd = s.execCommand("git", "stash", "clear")
		if err := cmd.Run(); err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: failed to clear stashes\n")
			return
		}
		_, _ = fmt.Fprintf(s.outputWriter, "Removed all stashes\n")

	case "create":
		// Create a stash entry and return its object name
		gitArgs := []string{"stash", "create"}
		if len(args) > 1 {
			gitArgs = append(gitArgs, args[1:]...)
		}
		cmd = s.execCommand("git", gitArgs...)
		output, err := cmd.Output()
		if err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: failed to create stash\n")
			return
		}
		_, _ = fmt.Fprintf(s.outputWriter, "%s", output)

	case "store":
		// Store the stash created by git stash create
		if len(args) < 2 {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: stash object required\n")
			return
		}
		gitArgs := []string{"stash", "store"}
		gitArgs = append(gitArgs, args[1:]...)
		cmd = s.execCommand("git", gitArgs...)
		if err := cmd.Run(); err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: failed to store stash\n")
			return
		}
		_, _ = fmt.Fprintf(s.outputWriter, "Stored stash\n")

	default:
		s.helper.ShowStashHelp()
		return
	}
}
