// Package git provides a high-level interface to git commands.
package git

// ResetHardAndClean resets the current branch to the state of origin and cleans the working directory.
func (c *Client) ResetHardAndClean() error {
	branch, err := c.GetCurrentBranch()
	if err != nil {
		return err
	}
	cmd := c.execCommand("git", "reset", "--hard", "origin/"+branch)
	if err := cmd.Run(); err != nil {
		return err
	}
	return c.CleanDirs()
}
