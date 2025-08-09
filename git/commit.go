// Package git provides a high-level interface to git commands.
package git

import (
	"os"
)

// CommitAllowEmpty commits with --allow-empty.
func (c *Client) CommitAllowEmpty() error {
	cmd := c.execCommand("git", "commit", "--allow-empty", "-m", "empty commit")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("commit allow empty", "git commit --allow-empty -m 'empty commit'", err)
	}
	return nil
}
