// Package git provides a high-level interface to git commands.
package git

// ResetHardAndClean resets the current branch to the state of origin and cleans the working directory.
func (c *Client) ResetHardAndClean() error {
	branch, err := c.GetCurrentBranch()
	if err != nil {
		return NewError("reset hard and clean", "get current branch", err)
	}
	cmd := c.execCommand("git", "reset", "--hard", "origin/"+branch)
	if err := cmd.Run(); err != nil {
		return NewError("reset hard and clean", "git reset --hard origin/"+branch, err)
	}
	if err := c.CleanDirs(); err != nil {
		return NewError("reset hard and clean", "clean directories", err)
	}
	return nil
}

// ResetHard resets to the specified commit.
func (c *Client) ResetHard(commit string) error {
	cmd := c.execCommand("git", "reset", "--hard", commit)
	if err := cmd.Run(); err != nil {
		return NewError("reset hard", "git reset --hard "+commit, err)
	}
	return nil
}
