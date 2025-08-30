package git

import (
	"fmt"
	"os"
	"strings"
)

// LogOneline gets git log output in oneline format between commits.
func (c *Client) LogOneline(from, to string) (string, error) {
	cmd := c.execCommand("git", "log", "--oneline", "--reverse", fmt.Sprintf("%s..%s", from, to))
	out, err := cmd.Output()
	if err != nil {
		return "", NewError("log oneline", fmt.Sprintf("git log --oneline --reverse %s..%s", from, to), err)
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
		return NewError("rebase interactive", fmt.Sprintf("git rebase -i HEAD~%d", commitCount), err)
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
