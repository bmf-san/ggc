// Package git provides a high-level interface to git commands.
package git

import (
	"os"
)

// StashPullPop stashes changes, pulls from remote, and pops stashed changes.
func (c *Client) StashPullPop() error {
	// Stash changes
	stashCmd := c.execCommand("git", "stash")
	stashCmd.Stdout = os.Stdout
	stashCmd.Stderr = os.Stderr
	if err := stashCmd.Run(); err != nil {
		return NewError("stash", "git stash", err)
	}

	// Pull changes
	if err := c.Pull(false); err != nil {
		return err
	}

	// Pop stashed changes
	popCmd := c.execCommand("git", "stash", "pop")
	popCmd.Stdout = os.Stdout
	popCmd.Stderr = os.Stderr
	if err := popCmd.Run(); err != nil {
		return NewError("stash pop", "git stash pop", err)
	}

	return nil
}
