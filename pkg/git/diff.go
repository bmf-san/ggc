package git

import "strings"

// DiffReader provides read-only diff operations.
// Implemented by Client and any compatible mock in tests.
type DiffReader interface {
	Diff() (string, error)
	DiffStaged() (string, error)
	DiffHead() (string, error)
	DiffWith(args []string) (string, error)
}

// Diff gets git diff output.
func (c *Client) Diff() (string, error) {
	return c.DiffWith(nil)
}

// DiffStaged gets git diff --staged output.
func (c *Client) DiffStaged() (string, error) {
	return c.DiffWith([]string{"--staged"})
}

// DiffHead gets git diff HEAD output.
func (c *Client) DiffHead() (string, error) {
	return c.DiffWith([]string{"HEAD"})
}

// DiffWith executes git diff with custom arguments.
func (c *Client) DiffWith(args []string) (string, error) {
	cmdArgs := append([]string{"diff"}, args...)
	cmd := c.execCommand("git", cmdArgs...)
	out, err := cmd.Output()
	if err != nil {
		command := strings.Join(append([]string{"git"}, cmdArgs...), " ")
		return "", NewOpError("get diff", command, err)
	}
	return string(out), nil
}
