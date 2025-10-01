package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/bmf-san/ggc/v7/git"
)

// Restorer handles restore operations.
type Restorer struct {
	outputWriter io.Writer
	helper       *Helper
	execCommand  func(string, ...string) *exec.Cmd
	gitClient    git.RestoreOps
}

// NewRestorer creates a new Restorer instance.
func NewRestorer(client git.RestoreOps) *Restorer {
	return &Restorer{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		execCommand:  exec.Command,
		gitClient:    client,
	}
}

// Restore executes git restore commands.
func (r *Restorer) Restore(args []string) {
	if len(args) == 0 {
		r.helper.ShowRestoreHelp()
		return
	}
	if args[0] == "staged" {
		if len(args) < 2 {
			r.helper.ShowRestoreHelp()
			return
		}
		r.restoreStaged(args[1:])
		return
	}

	r.restoreCommitOrWorking(args)
}

func (r *Restorer) restoreStaged(paths []string) {
	if len(paths) < 1 {
		r.helper.ShowRestoreHelp()
		return
	}
	if err := r.gitClient.RestoreStaged(paths...); err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
	}
}

func (r *Restorer) restoreCommitOrWorking(args []string) {
	if len(args) >= 2 && (r.gitClient.RevParseVerify(args[0]) || isCommitLikeStrict(args[0])) {
		commit := args[0]
		paths := args[1:]
		if err := r.gitClient.RestoreFromCommit(commit, paths...); err != nil {
			_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
		}
		return
	}
	if err := r.gitClient.RestoreWorkingDir(args...); err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
	}
}

// isCommitLikeStrict performs cheap, defensive checks without panicking.
// It intentionally narrows matches to avoid false positives and defers to
// RevParseVerify when available for authoritative validation.
func isCommitLikeStrict(s string) bool {
	return isHEADVariant(s) || isExplicitRef(s) || isRemoteRef(s) || isHexObjectName(s)
}

// isHEADVariant checks for tight HEAD variants
func isHEADVariant(s string) bool {
	return s == "HEAD" || strings.HasPrefix(s, "HEAD^") || strings.HasPrefix(s, "HEAD~") || strings.HasPrefix(s, "HEAD@{")
}

// isExplicitRef checks for explicit ref namespaces
func isExplicitRef(s string) bool {
	return strings.HasPrefix(s, "refs/") // e.g., refs/heads/main
}

// isRemoteRef checks for common remote ref format
func isRemoteRef(s string) bool {
	return strings.HasPrefix(s, "origin/") // e.g., origin/main
}

// isHexObjectName checks for hex-ish object name (short/long SHA)
func isHexObjectName(s string) bool {
	l := len(s)
	if l < 7 || l > 64 { // Allow up to 64 for SHA-256 compatibility
		return false
	}

	for _, r := range s {
		if !isHexChar(r) {
			return false
		}
	}
	return true
}

// isHexChar checks if a rune is a valid hex character
func isHexChar(r rune) bool {
	return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}
