package git

import (
	"os"
	"strings"
)

// PassthroughOps provides a generic mechanism to run an arbitrary git
// subcommand and stream its output to the current process's stdout/stderr.
// It is used by lightweight ggc wrappers that simply forward arguments to git
// (e.g. ggc cherry-pick, ggc revert, ggc blame).
type PassthroughOps interface {
	RunGit(name string, args []string) error
}

// RunGit invokes `git <name> [args...]`, wiring stdout/stderr to the host
// process so the user sees normal git output (including pagers when the
// terminal supports them).
func (c *Client) RunGit(name string, args []string) error {
	gitArgs := append([]string{name}, args...)
	cmd := c.execCommand("git", gitArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewOpError(name, "git "+name+joinArgs(args), err)
	}
	return nil
}

func joinArgs(args []string) string {
	if len(args) == 0 {
		return ""
	}
	return " " + strings.Join(args, " ")
}
