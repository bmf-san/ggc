// Package git provides a high-level interface to git commands.
package git

// Push pushes to a remote.
func (c *Client) Push(force bool) error {
	branch, err := c.GetCurrentBranch()
	if err != nil {
		return err
	}
	args := []string{"push", "origin", branch}
	if force {
		args = append(args, "--force-with-lease")
	}
	cmd := c.execCommand("git", args...)
	return cmd.Run()
}
