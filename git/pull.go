// Package git provides a high-level interface to git commands.
package git

// Pull pulls from a remote.
func (c *Client) Pull(rebase bool) error {
	args := []string{"pull"}
	if rebase {
		args = append(args, "--rebase")
	}
	cmd := c.execCommand("git", args...)
	return cmd.Run()
}
