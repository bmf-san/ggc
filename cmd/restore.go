package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/bmf-san/ggc/v4/git"
)

// Restoreer handles restore operations.
type Restoreer struct {
	outputWriter io.Writer
	helper       *Helper
	execCommand  func(string, ...string) *exec.Cmd
	gitClient    git.Clienter
}

// NewRestoreer creates a new Restoreer instance.
func NewRestoreer() *Restoreer {
	return &Restoreer{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		execCommand:  exec.Command,
		gitClient:    git.NewClient(),
	}
}

// Restore executes git restore commands.
func (r *Restoreer) Restore(args []string) {
	if len(args) == 0 {
		r.helper.ShowRestoreHelp()
		return
	}

	switch args[0] {
	case "staged":
		if len(args) < 2 {
			r.helper.ShowRestoreHelp()
			return
		}

		paths := args[1:]
		if err := r.gitClient.RestoreStaged(paths...); err != nil {
			_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
			return
		}

	default:
		if len(args) >= 2 && isCommitLike(args[0]) {
			// Handle : ggc restore <commit> <file>
			commit := args[0]
			paths := args[1:]
			if err := r.gitClient.RestoreFromCommit(commit, paths...); err != nil {
				_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
				return
			}
		} else {
			// Handle : ggc restore <file> or ggc restore .
			if err := r.gitClient.RestoreWorkingDir(args...); err != nil {
				_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
				return
			}
		}
	}
}

func isCommitLike(s string) bool {
	// Hex-ish object name (short/long SHA)
	if l := len(s); l >= 7 && l <= 40 {
		isHex := true
		for _, r := range s {
			if (r < '0' || r > '9') && (r < 'a' || r > 'f') && (r < 'A' || r > 'F') {
				isHex = false
				break
			}
		}
		if isHex {
			return true
		}
	}
	// Safe prefix checks for ref-ish values
	if strings.HasPrefix(s, "HEAD") { // e.g., HEAD, HEAD~1, HEAD^
		return true
	}
	if strings.HasPrefix(s, "refs/") { // e.g., refs/heads/main
		return true
	}
	if strings.HasPrefix(s, "origin/") { // e.g., origin/main
		return true
	}
	return false
}
