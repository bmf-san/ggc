// Package git provides a high-level interface to git commands.
package git

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

// CommitWriter provides write operations for commits.
type CommitWriter interface {
	Commit(message string) error
	CommitAmend() error
	CommitAmendNoEdit() error
	CommitAmendWithMessage(message string) error
	CommitAllowEmpty() error
}

// Commit commits with the given message.
func (c *Client) Commit(message string) error {
	if err := validateCommitMessage(message); err != nil {
		return err
	}

	cmd := c.execCommand("git", "commit", "-m", message)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewOpError("commit", "git commit -m "+message, err)
	}
	return nil
}

// CommitAmend amends the last commit.
func (c *Client) CommitAmend() error {
	cmd := c.execCommand("git", "commit", "--amend")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return NewOpError("commit amend", "git commit --amend", err)
	}
	return nil
}

// CommitAmendNoEdit amends the last commit without editing the message.
func (c *Client) CommitAmendNoEdit() error {
	cmd := c.execCommand("git", "commit", "--amend", "--no-edit")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewOpError("commit amend no-edit", "git commit --amend --no-edit", err)
	}
	return nil
}

// CommitAmendWithMessage amends the last commit with a new message.
func (c *Client) CommitAmendWithMessage(message string) error {
	if err := validateCommitMessage(message); err != nil {
		return err
	}

	cmd := c.execCommand("git", "commit", "--amend", "-m", message)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewOpError("commit amend with message", "git commit --amend -m "+message, err)
	}
	return nil
}

func validateCommitMessage(message string) error {
	trimmed := strings.TrimSpace(message)
	if trimmed == "" {
		return fmt.Errorf("commit message cannot be empty or whitespace")
	}
	if !utf8.ValidString(message) {
		return fmt.Errorf("commit message must be valid UTF-8")
	}
	if strings.ContainsRune(message, '\u0000') {
		return fmt.Errorf("commit message cannot contain null characters")
	}
	return nil
}

// CommitAllowEmpty commits with --allow-empty.
func (c *Client) CommitAllowEmpty() error {
	cmd := c.execCommand("git", "commit", "--allow-empty", "-m", "empty commit")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewOpError("commit allow empty", "git commit --allow-empty -m 'empty commit'", err)
	}
	return nil
}
