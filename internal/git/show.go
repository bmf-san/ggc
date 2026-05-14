package git

import "os"

// ShowOps provides access to the git show command.
type ShowOps interface {
	Show(args []string) error
}

// Show runs `git show` with the supplied arguments, streaming output to stdout.
// When args is empty, `git show` (with no arguments) shows the HEAD commit.
func (c *Client) Show(args []string) error {
	gitArgs := append([]string{"show"}, args...)
	cmd := c.execCommand("git", gitArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		full := "git show"
		for _, a := range args {
			full += " " + a
		}
		return NewOpError("show", full, err)
	}
	return nil
}
