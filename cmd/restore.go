package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/bmf-san/ggc/v5/git"
)

// Restoreer handles restore operations.
type Restoreer struct {
	outputWriter io.Writer
	helper       *Helper
	execCommand  func(string, ...string) *exec.Cmd
	gitClient    git.Clienter
}

// NewRestoreer creates a new Restoreer instance.
func NewRestoreer(client git.Clienter) *Restoreer {
	return &Restoreer{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		execCommand:  exec.Command,
		gitClient:    client,
	}
}

// Restore executes git restore commands.
func (r *Restoreer) Restore(args []string) {
	if len(args) == 0 {
		r.helper.ShowRestoreHelp()
		return
	}

	switch args[0] {
	case "--staged":
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
		// Prefer git to validate commit-ish to avoid false positives
		if len(args) >= 2 && (r.gitClient.RevParseVerify(args[0]) || isCommitLikeStrict(args[0])) {
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

// isCommitLikeStrict performs cheap, defensive checks without panicking.
// It intentionally narrows matches to avoid false positives and defers to
// RevParseVerify when available for authoritative validation.
func isCommitLikeStrict(s string) bool {
	// Tight HEAD variants
	if s == "HEAD" || strings.HasPrefix(s, "HEAD^") || strings.HasPrefix(s, "HEAD~") || strings.HasPrefix(s, "HEAD@{") {
		return true
	}
	// Explicit ref namespaces
	if strings.HasPrefix(s, "refs/") { // e.g., refs/heads/main
		return true
	}
	// Common remote ref format. Note: other remotes will be caught by rev-parse.
	if strings.HasPrefix(s, "origin/") { // e.g., origin/main
		return true
	}
	// Hex-ish object name (short/long SHA). Allow up to 64 for SHA-256 compatibility.
	if l := len(s); l >= 7 && l <= 64 {
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
	return false
}
