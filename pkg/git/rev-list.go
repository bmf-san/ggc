package git

import (
	"strings"
)

// GetAheadBehindCount gets the ahead/behind count between branch and upstream.
func (c *Client) GetAheadBehindCount(branch, upstream string) (string, error) {
	cmd := c.execCommand("git", "rev-list", "--left-right", "--count", branch+"..."+upstream)
	out, err := cmd.Output()
	if err != nil {
		return "", NewError("get ahead behind count", "git rev-list --left-right --count "+branch+"..."+upstream, err)
	}
	return strings.TrimSpace(string(out)), nil
}

// GetTagCommit gets the commit hash for a tag.
func (c *Client) GetTagCommit(name string) (string, error) {
	cmd := c.execCommand("git", "rev-list", "-n", "1", name)
	output, err := cmd.Output()
	if err != nil {
		return "", NewError("get tag commit", "git rev-list -n 1 "+name, err)
	}
	return strings.TrimSpace(string(output)), nil
}
