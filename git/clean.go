// Package git provides a high-level interface to git commands.
package git

import (
	"os"
)

// CleanFiles cleans untracked files.
func (c *Client) CleanFiles() error {
	cmd := c.execCommand("git", "clean", "-fd")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("clean files", "git clean -fd", err)
	}
	return nil
}

// CleanDirs cleans untracked directories.
func (c *Client) CleanDirs() error {
	cmd := c.execCommand("git", "clean", "-fdx")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("clean directories", "git clean -fdx", err)
	}
	return nil
}
