package git

import (
	"os"
	"strings"
)

// Add adds files to the staging area.
func (c *Client) Add(files ...string) error {
	if len(files) == 0 {
		return NewError("add files", "git add", nil)
	}

	args := append([]string{"add"}, files...)
	cmd := c.execCommand("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("add files", "git "+strings.Join(args, " "), err)
	}
	return nil
}

// AddInteractive starts interactive staging.
func (c *Client) AddInteractive() error {
	cmd := c.execCommand("git", "add", "-p")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return NewError("interactive add", "git add -p", err)
	}
	return nil
}
