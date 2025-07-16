// Package git provides a high-level interface to git commands.
package git

import (
	"fmt"
	"strings"
)

// Push pushes to a remote.
func (c *Client) Push(force bool) error {
	branch, err := c.GetCurrentBranch()
	if err != nil {
		return NewError("push", "get current branch", err)
	}
	args := []string{"push", "origin", branch}
	if force {
		args = append(args, "--force-with-lease")
	}
	cmd := c.execCommand("git", args...)
	if err := cmd.Run(); err != nil {
		return NewError("push", fmt.Sprintf("git %s", strings.Join(args, " ")), err)
	}
	return nil
}
