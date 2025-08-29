package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/bmf-san/ggc/v4/git"
)

// Statuseer handles status operations.
type Statuseer struct {
	outputWriter io.Writer
	helper       *Helper
	execCommand  func(string, ...string) *exec.Cmd
	gitClient    git.Clienter
}

// NewStatuseer creates a new Statuseer instance.
func NewStatuseer() *Statuseer {
	return &Statuseer{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		execCommand:  exec.Command,
		gitClient:    git.NewClient(),
	}
}

// getUpstreamStatus gets the upstream tracking status
func (s *Statuseer) getUpstreamStatus(branch string) string {
	// Check if upstream exists
	cmd := s.execCommand("git", "rev-parse", "--abbrev-ref", branch+"@{upstream}")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	upstream := strings.TrimSpace(string(output))

	// Get ahead/behind count
	cmd = s.execCommand("git", "rev-list", "--left-right", "--count", branch+"..."+upstream)
	output, err = cmd.Output()
	if err != nil {
		return fmt.Sprintf("Your branch is up to date with '%s'", upstream)
	}

	counts := strings.Fields(strings.TrimSpace(string(output)))
	if len(counts) == 2 {
		ahead := counts[0]
		behind := counts[1]

		if ahead == "0" && behind == "0" {
			return fmt.Sprintf("Your branch is up to date with '%s'", upstream)
		}
		if ahead != "0" && behind == "0" {
			return fmt.Sprintf("Your branch is ahead of '%s' by %s commit(s)", upstream, ahead)
		}
		if ahead == "0" && behind != "0" {
			return fmt.Sprintf("Your branch is behind '%s' by %s commit(s)", upstream, behind)
		}
		return fmt.Sprintf("Your branch and '%s' have diverged", upstream)
	}

	return fmt.Sprintf("Your branch is up to date with '%s'", upstream)
}

// Status executes git status with the given arguments.
func (s *Statuseer) Status(args []string) {
	var cmd *exec.Cmd
	if len(args) == 0 {
		// Add '-c color.status=always' to ensure colour showing up in 'less'
		cmd = s.execCommand("git", "-c", "color.status=always", "status")
	} else {
		switch args[0] {
		case "short":
			cmd = s.execCommand("git", "-c", "color.status=always", "status", "--short")
		default:
			s.helper.ShowStatusHelp()
			return
		}
	}

	branch, err := s.gitClient.GetCurrentBranch()
	if err != nil {
		_, _ = fmt.Fprintf(s.outputWriter, "Error getting current branch: %v\n", err)
		return
	}
	upstreamStatus := s.getUpstreamStatus(branch)

	if _, err := exec.LookPath("less"); err != nil {
		// Fallback: If 'less' is not available, direct output to outputWriter
		_, _ = fmt.Fprintf(s.outputWriter, "On branch %s\n", branch)
		if upstreamStatus != "" {
			_, _ = fmt.Fprintf(s.outputWriter, "%s\n", upstreamStatus)
		}
		_, _ = fmt.Fprintf(s.outputWriter, "\n")

		cmd.Stdout = s.outputWriter
		cmd.Stderr = s.outputWriter
		if err := cmd.Run(); err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error running git status: %v\n", err)
		}
		return
	}

	// Setup 'less' pipeline with branch info prepended
	lessCmd := exec.Command("less", "-R")
	gitStdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		_, _ = fmt.Fprintf(s.outputWriter, "Error creating stdout pipe for git: %v\n", err)
		return
	}
	cmd.Stderr = s.outputWriter

	// Create a pipe to combine branch info with git output
	lessStdinPipe, err := lessCmd.StdinPipe()
	if err != nil {
		_, _ = fmt.Fprintf(s.outputWriter, "Error creating stdin pipe for less: %v\n", err)
		return
	}

	lessCmd.Stdout = s.outputWriter
	lessCmd.Stderr = s.outputWriter

	// Start both commands
	if err := cmd.Start(); err != nil {
		_, _ = fmt.Fprintf(s.outputWriter, "Error starting git command: %v\n", err)
		return
	}
	if err := lessCmd.Start(); err != nil {
		_, _ = fmt.Fprintf(s.outputWriter, "Error starting less command: %v\n", err)
		// Drain output to avoid deadlocking
		if _, err := io.Copy(io.Discard, gitStdoutPipe); err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error discarding output: %v\n", err)
		}
		if err := cmd.Wait(); err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error waiting for git command: %v\n", err)
		}
		return
	}

	_, _ = fmt.Fprintf(lessStdinPipe, "On branch %s\n", branch)
	if upstreamStatus != "" {
		_, _ = fmt.Fprintf(lessStdinPipe, "%s\n", upstreamStatus)
	}
	_, _ = fmt.Fprintf(lessStdinPipe, "\n")

	go func() {
		defer func() {
			if err := lessStdinPipe.Close(); err != nil {
				_, _ = fmt.Fprintf(s.outputWriter, "Error closing lessStdinPipe: %v\n", err)
			}
		}()
		if _, err := io.Copy(lessStdinPipe, gitStdoutPipe); err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error copying git output: %v\n", err)
		}
	}()

	// Wait for both commands to finish
	if err := cmd.Wait(); err != nil {
		_, _ = fmt.Fprintf(s.outputWriter, "Error waiting for git command: %v\n", err)
	}
	if err := lessCmd.Wait(); err != nil {
		_, _ = fmt.Fprintf(s.outputWriter, "Error waiting for less command: %v\n", err)
	}
}
