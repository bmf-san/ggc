// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/bmf-san/ggc/v6/git"
)

// Logger provides functionality for the log command.
type Logger struct {
	gitClient    git.LogReader
	outputWriter io.Writer
	execCommand  func(name string, arg ...string) *exec.Cmd
	helper       *Helper
}

// NewLogger creates a new Logger.
func NewLogger(client git.LogReader) *Logger {
	l := &Logger{
		gitClient:    client,
		outputWriter: os.Stdout,
		execCommand:  exec.Command,
		helper:       NewHelper(),
	}
	l.helper.outputWriter = l.outputWriter
	return l
}

// Log executes the log command with the given arguments.
func (l *Logger) Log(args []string) {
	if len(args) == 0 {
		l.helper.ShowLogHelp()
		return
	}

	switch args[0] {
	case "simple":
		if err := l.gitClient.LogSimple(); err != nil {
			_, _ = fmt.Fprintf(l.outputWriter, "Error: %v\n", err)
		}
	case "graph":
		if err := l.gitClient.LogGraph(); err != nil {
			_, _ = fmt.Fprintf(l.outputWriter, "Error: %v\n", err)
		}
	default:
		l.helper.ShowLogHelp()
	}
}
