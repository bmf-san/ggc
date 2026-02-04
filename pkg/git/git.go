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
