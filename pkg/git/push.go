// Package git provides a high-level interface to git commands.
package git

import (
	"fmt"
	"os"
	"strings"
)

// Push pushes to a remote.
func (c *Client) Push(force bool) error {
	branch, err := c.GetCurrentBranch()
	if err != nil {
		return NewOpError("push", "get current branch", err)
	}
	args := []string{"push", "origin", branch}
	if force {
		args = append(args, "--force-with-lease")
	}
	cmd := c.execCommand("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewOpError("push", fmt.Sprintf("git %s", strings.Join(args, " ")), err)
	}
	return nil
}
