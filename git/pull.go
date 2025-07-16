// Package git provides a high-level interface to git commands.
package git

import (
	"fmt"
	"strings"
)

// Pull pulls from a remote.
func (c *Client) Pull(rebase bool) error {
	args := []string{"pull"}
	if rebase {
		args = append(args, "--rebase")
	}
	cmd := c.execCommand("git", args...)
	if err := cmd.Run(); err != nil {
		return NewError("pull", fmt.Sprintf("git %s", strings.Join(args, " ")), err)
	}
	return nil
}
