package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// Client is a git client.
type Client struct {
	execCommand          func(name string, arg ...string) *exec.Cmd
	GetCurrentBranchFunc func() (string, error)
}

// Clienter is an interface for a git client.
type Clienter interface {
	GetCurrentBranch() (string, error)
	ListLocalBranches() ([]string, error)
	ListRemoteBranches() ([]string, error)
	CheckoutNewBranch(name string) error
	Push(force bool) error
	Pull(rebase bool) error
	LogSimple() error
	LogGraph() error
	CommitAllowEmpty() error
	ResetHardAndClean() error
	CleanFiles() error
	CleanDirs() error
	GetGitStatus() (string, error)
	GetBranchName() (string, error)
	RestoreWorkingDir(paths ...string) error
	RestoreStaged(paths ...string) error
	RestoreFromCommit(commit string, paths ...string) error
	RestoreAll() error
	RestoreAllStaged() error
}

// NewClient creates a new Client.
func NewClient() *Client {
	return &Client{
		execCommand: exec.Command,
	}
}

// GetGitStatus gets git status.
func (c *Client) GetGitStatus() (string, error) {
	cmd := c.execCommand("git", "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return "", NewError("get status", "git status --porcelain", err)
	}
	return string(out), nil
}

// GetBranchName gets branch name.
func (c *Client) GetBranchName() (string, error) {
	cmd := c.execCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", NewError("get branch name", "git rev-parse --abbrev-ref HEAD", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// CheckoutNewBranch creates a new branch and checks it out.
func (c *Client) CheckoutNewBranch(name string) error {
	cmd := c.execCommand("git", "checkout", "-b", name)
	if err := cmd.Run(); err != nil {
		return NewError("checkout new branch", fmt.Sprintf("git checkout -b %s", name), err)
	}
	return nil
}

// GetCurrentBranch gets the current branch name.
func (c *Client) GetCurrentBranch() (string, error) {
	if c.GetCurrentBranchFunc != nil {
		return c.GetCurrentBranchFunc()
	}
	cmd := c.execCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", NewError("get current branch", "git rev-parse --abbrev-ref HEAD", err)
	}
	branch := strings.TrimSpace(string(out))
	return branch, nil
}
