package git

import (
	"fmt"
	"os"
	"strings"
)

// RebaseOps provides operations used by the rebase command.
type RebaseOps interface {
	// sequence operations
	RebaseInteractive(commitCount int) error
	Rebase(upstream string) error
	RebaseContinue() error
	RebaseAbort() error
	RebaseSkip() error
	// discovery
	GetCurrentBranch() (string, error)
	GetUpstreamBranch(branch string) (string, error)
	LogOneline(from, to string) (string, error)
	RevParseVerify(ref string) bool
}

// LogOneline gets git log output in oneline format between commits.
func (c *Client) LogOneline(from, to string) (string, error) {
	cmd := c.execCommand("git", "log", "--oneline", "--reverse", fmt.Sprintf("%s..%s", from, to))
	out, err := cmd.Output()
	if err != nil {
		return "", NewOpError("log oneline", fmt.Sprintf("git log --oneline --reverse %s..%s", from, to), err)
	}
	return string(out), nil
}

// RebaseInteractive starts an interactive rebase for the specified number of commits.
func (c *Client) RebaseInteractive(commitCount int) error {
	cmd := c.execCommand("git", "rebase", "-i", fmt.Sprintf("HEAD~%d", commitCount))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewOpError("rebase interactive", fmt.Sprintf("git rebase -i HEAD~%d", commitCount), err)
	}
	return nil
}

// Rebase performs a basic rebase onto the given upstream reference.
func (c *Client) Rebase(upstream string) error {
	cmd := c.execCommand("git", "rebase", upstream)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewOpError("rebase", fmt.Sprintf("git rebase %s", upstream), err)
	}
	return nil
}

// RebaseContinue continues an in-progress rebase.
func (c *Client) RebaseContinue() error {
	cmd := c.execCommand("git", "rebase", "--continue")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewOpError("rebase continue", "git rebase --continue", err)
	}
	return nil
}

// RebaseAbort aborts an in-progress rebase.
func (c *Client) RebaseAbort() error {
	cmd := c.execCommand("git", "rebase", "--abort")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewOpError("rebase abort", "git rebase --abort", err)
	}
	return nil
}

// RebaseSkip skips the current patch and continues rebasing.
func (c *Client) RebaseSkip() error {
	cmd := c.execCommand("git", "rebase", "--skip")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewOpError("rebase skip", "git rebase --skip", err)
	}
	return nil
}

// GetUpstreamBranch gets the upstream branch for the given branch.
func (c *Client) GetUpstreamBranch(branch string) (string, error) {
	cmd := c.execCommand("git", "rev-parse", "--abbrev-ref", fmt.Sprintf("%s@{upstream}", branch))
	out, err := cmd.Output()
	if err != nil {
		// If no upstream is set, return "main" as default
		return "main", nil
	}
	return strings.TrimSpace(string(out)), nil
}
