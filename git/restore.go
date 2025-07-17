// Package git provides a high-level interface to git commands.
package git

import (
	"fmt"
	"strings"
)

// RestoreOptions holds options for git restore command
type RestoreOptions struct {
	Staged bool   //  (from HEAD to index)
	Source string // (from specific commit)
}

// Restore runs `git restore` with optional paths and options.
func (c *Client) Restore(paths []string, opts *RestoreOptions) error {
	args := []string{"restore"}

	if opts != nil {
		if opts.Staged {
			args = append(args, "--staged")
		}
		if opts.Source != "" {
			args = append(args, "--source", opts.Source)
		}
	}

	args = append(args, paths...)
	cmd := c.execCommand("git", args...)
	if err := cmd.Run(); err != nil {
		return NewError("restore", fmt.Sprintf("git %s", strings.Join(args, " ")), err)
	}
	return nil
}

// RestoreWorkingDir restores files in working directory from index
func (c *Client) RestoreWorkingDir(paths ...string) error {
	return c.Restore(paths, nil)
}

// RestoreStaged unstages files (restores from HEAD to index)
func (c *Client) RestoreStaged(paths ...string) error {
	return c.Restore(paths, &RestoreOptions{Staged: true})
}

// RestoreFromCommit restores files from a specific commit
func (c *Client) RestoreFromCommit(commit string, paths ...string) error {
	return c.Restore(paths, &RestoreOptions{Source: commit})
}

// RestoreAll restores all files in working directory from index
func (c *Client) RestoreAll() error {
	return c.Restore([]string{"."}, nil)
}

// RestoreAllStaged unstages all files
func (c *Client) RestoreAllStaged() error {
	return c.Restore([]string{"."}, &RestoreOptions{Staged: true})
}
