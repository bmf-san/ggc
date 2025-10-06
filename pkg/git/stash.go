package git

import (
	"os"
)

// Stash creates a new stash.
func (c *Client) Stash() error {
	cmd := c.execCommand("git", "stash")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("stash", "git stash", err)
	}
	return nil
}

// StashList lists all stashes.
func (c *Client) StashList() (string, error) {
	cmd := c.execCommand("git", "stash", "list")
	out, err := cmd.Output()
	if err != nil {
		return "", NewError("stash list", "git stash list", err)
	}
	return string(out), nil
}

// StashShow shows a stash.
func (c *Client) StashShow(stash string) error {
	var cmd = c.execCommand("git", "stash", "show")
	if stash != "" {
		cmd = c.execCommand("git", "stash", "show", stash)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		cmdStr := "git stash show"
		if stash != "" {
			cmdStr = "git stash show " + stash
		}
		return NewError("stash show", cmdStr, err)
	}
	return nil
}

// StashApply applies a stash.
func (c *Client) StashApply(stash string) error {
	var cmd = c.execCommand("git", "stash", "apply")
	if stash != "" {
		cmd = c.execCommand("git", "stash", "apply", stash)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		cmdStr := "git stash apply"
		if stash != "" {
			cmdStr = "git stash apply " + stash
		}
		return NewError("stash apply", cmdStr, err)
	}
	return nil
}

// StashPop pops a stash.
func (c *Client) StashPop(stash string) error {
	var cmd = c.execCommand("git", "stash", "pop")
	if stash != "" {
		cmd = c.execCommand("git", "stash", "pop", stash)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		cmdStr := "git stash pop"
		if stash != "" {
			cmdStr = "git stash pop " + stash
		}
		return NewError("stash pop", cmdStr, err)
	}
	return nil
}

// StashDrop drops a stash.
func (c *Client) StashDrop(stash string) error {
	var cmd = c.execCommand("git", "stash", "drop")
	if stash != "" {
		cmd = c.execCommand("git", "stash", "drop", stash)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		cmdStr := "git stash drop"
		if stash != "" {
			cmdStr = "git stash drop " + stash
		}
		return NewError("stash drop", cmdStr, err)
	}
	return nil
}

// StashClear clears all stashes.
func (c *Client) StashClear() error {
	cmd := c.execCommand("git", "stash", "clear")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("stash clear", "git stash clear", err)
	}
	return nil
}
