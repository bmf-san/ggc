package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Statuseer handles status operations.
type Statuseer struct {
	outputWriter io.Writer
	helper       *Helper
	execCommand  func(string, ...string) *exec.Cmd
}

// NewStatuseer creates a new Statuseer instance.
func NewStatuseer() *Statuseer {
	return &Statuseer{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		execCommand:  exec.Command,
	}
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

	if _, err := exec.LookPath("less"); err != nil {
		// Fallback: If 'less' is not available, direct output to outputWriter
		cmd.Stdout = s.outputWriter
		cmd.Stderr = s.outputWriter
		if err := cmd.Run(); err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error running git status: %v\n", err)
		}
		return
	}

	// Setup 'less' pipeline
	lessCmd := exec.Command("less", "-R")
	gitStdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		_, _ = fmt.Fprintf(s.outputWriter, "Error creating stdout pipe for git: %v\n", err)
		return
	}
	cmd.Stderr = s.outputWriter
	lessCmd.Stdin = gitStdoutPipe
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

	// Wait for both commands to finish
	if err := cmd.Wait(); err != nil {
		_, _ = fmt.Fprintf(s.outputWriter, "Error waiting for git command: %v\n", err)
	}

	if err := lessCmd.Wait(); err != nil {
		_, _ = fmt.Fprintf(s.outputWriter, "Error waiting for less command: %v\n", err)
	}
}
