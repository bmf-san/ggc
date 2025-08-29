package git

import (
	"strings"
)

// ListFiles lists all files managed by git.
func (c *Client) ListFiles() (string, error) {
	cmd := c.execCommand("git", "ls-files")
	out, err := cmd.Output()
	if err != nil {
		return "", NewError("list files", "git ls-files", err)
	}
	return string(out), nil
}

// GetUpstreamBranchName gets the upstream branch name for a given branch.
func (c *Client) GetUpstreamBranchName(branch string) (string, error) {
	cmd := c.execCommand("git", "rev-parse", "--abbrev-ref", branch+"@{upstream}")
	out, err := cmd.Output()
	if err != nil {
		return "", NewError("get upstream branch", "git rev-parse --abbrev-ref "+branch+"@{upstream}", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// GetAheadBehindCount gets the ahead/behind count between branch and upstream.
func (c *Client) GetAheadBehindCount(branch, upstream string) (string, error) {
	cmd := c.execCommand("git", "rev-list", "--left-right", "--count", branch+"..."+upstream)
	out, err := cmd.Output()
	if err != nil {
		return "", NewError("get ahead behind count", "git rev-list --left-right --count "+branch+"..."+upstream, err)
	}
	return strings.TrimSpace(string(out)), nil
}
