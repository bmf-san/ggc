package git

import (
	"os/exec"
)

// Client is a git client.
type Client struct {
	execCommand          func(name string, arg ...string) *exec.Cmd
	GetCurrentBranchFunc func() (string, error)
}

// NewClient creates a new Client.
func NewClient() *Client {
	return &Client{
		execCommand: exec.Command,
	}
}

// === Repository Information ===

// BranchInfo contains rich information about a branch.
type BranchInfo struct {
	Name            string
	IsCurrentBranch bool
	Upstream        string
	AheadBehind     string // e.g. "ahead 2, behind 1"
	LastCommitSHA   string
	LastCommitMsg   string
}
