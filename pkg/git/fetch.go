package git

import "os"

// Fetch fetches from remote repository.
func (c *Client) Fetch(prune bool) error {
	var cmd = c.execCommand("git", "fetch")
	if prune {
		cmd = c.execCommand("git", "fetch", "--prune")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if prune {
			return NewOpError("fetch with prune", "git fetch --prune", err)
		}
		return NewOpError("fetch", "git fetch", err)
	}
	return nil
}
