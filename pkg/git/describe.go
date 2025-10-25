package git

import (
	"strings"
)

// GetVersion gets the git version/tag information
func (c *Client) GetVersion() (string, error) {
	cmd := c.execCommand("git", "describe", "--tags", "--always", "--dirty")
	out, err := cmd.Output()
	if err != nil {
		return "dev", nil // Return "dev" as fallback instead of error
	}
	return strings.TrimSpace(string(out)), nil
}
