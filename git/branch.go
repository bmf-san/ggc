// Package git provides a high-level interface to git commands.
package git

import (
	"strings"
)

// ListLocalBranches lists local branches.
func (c *Client) ListLocalBranches() ([]string, error) {
	cmd := c.execCommand("git", "branch", "--format", "%(refname:short)")
	out, err := cmd.Output()
	if err != nil {
		return nil, NewError("list local branches", "git branch --format %(refname:short)", err)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	return lines, nil
}

// ListRemoteBranches lists remote branches.
func (c *Client) ListRemoteBranches() ([]string, error) {
	cmd := c.execCommand("git", "branch", "-r", "--format", "%(refname:short)")
	out, err := cmd.Output()
	if err != nil {
		return nil, NewError("list remote branches", "git branch -r --format %(refname:short)", err)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	// Exclude HEAD references such as origin/HEAD -> origin/main
	filtered := []string{}
	for _, l := range lines {
		if strings.Contains(l, "->") {
			continue
		}
		filtered = append(filtered, strings.TrimSpace(l))
	}
	return filtered, nil
}
