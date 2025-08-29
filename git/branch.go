// Package git provides a high-level interface to git commands.
package git

import (
	"fmt"
	"os"
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

// CheckoutBranch checks out an existing branch.
func (c *Client) CheckoutBranch(name string) error {
	cmd := c.execCommand("git", "checkout", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("checkout branch", "git checkout "+name, err)
	}
	return nil
}

// CheckoutNewBranchFromRemote creates a new local branch tracking a remote branch.
func (c *Client) CheckoutNewBranchFromRemote(localBranch, remoteBranch string) error {
	cmd := c.execCommand("git", "checkout", "-b", localBranch, "--track", remoteBranch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("checkout new branch from remote", fmt.Sprintf("git checkout -b %s --track %s", localBranch, remoteBranch), err)
	}
	return nil
}

// DeleteBranch deletes a branch.
func (c *Client) DeleteBranch(name string) error {
	cmd := c.execCommand("git", "branch", "-d", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("delete branch", "git branch -d "+name, err)
	}
	return nil
}

// ListMergedBranches lists branches that have been merged.
func (c *Client) ListMergedBranches() ([]string, error) {
	cmd := c.execCommand("git", "branch", "--merged")
	out, err := cmd.Output()
	if err != nil {
		return nil, NewError("list merged branches", "git branch --merged", err)
	}

	branches := strings.Split(strings.TrimSpace(string(out)), "\n")
	var result []string
	for _, branch := range branches {
		branch = strings.TrimSpace(branch)
		if branch != "" && !strings.HasPrefix(branch, "*") {
			result = append(result, branch)
		}
	}
	return result, nil
}
