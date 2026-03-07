// Package git provides a high-level interface to git commands.
package git

import (
	"os"
	"strings"
)

// CleanOps provides operations used by the clean command.
type CleanOps interface {
	CleanFiles() error
	CleanDirs() error
	CleanDryRun() (string, error)
	CleanFilesForce(files []string) error
}

// CleanFiles cleans untracked files.
func (c *Client) CleanFiles() error {
	cmd := c.execCommand("git", "clean", "-fd")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewOpError("clean files", "git clean -fd", err)
	}
	return nil
}

// CleanDirs cleans untracked directories.
func (c *Client) CleanDirs() error {
	cmd := c.execCommand("git", "clean", "-fdx")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewOpError("clean directories", "git clean -fdx", err)
	}
	return nil
}

// CleanDryRun shows what would be cleaned without actually cleaning.
func (c *Client) CleanDryRun() (string, error) {
	cmd := c.execCommand("git", "clean", "-nd")
	out, err := cmd.Output()
	if err != nil {
		return "", NewOpError("clean dry run", "git clean -nd", err)
	}
	return string(out), nil
}

// CleanFilesForce removes specific files forcefully.
func (c *Client) CleanFilesForce(files []string) error {
	if len(files) == 0 {
		return nil
	}

	args := append([]string{"clean", "-f", "--"}, files...)
	cmd := c.execCommand("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewOpError("clean files force", "git clean -f -- "+strings.Join(files, " "), err)
	}
	return nil
}
