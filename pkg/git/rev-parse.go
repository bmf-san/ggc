package git

import (
	"strings"
)

// GetCurrentBranch gets the current branch name.
func (c *Client) GetCurrentBranch() (string, error) {
	if c.GetCurrentBranchFunc != nil {
		return c.GetCurrentBranchFunc()
	}
	cmd := c.execCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", NewOpError("get current branch", "git rev-parse --abbrev-ref HEAD", err)
	}
	branch := strings.TrimSpace(string(out))
	return branch, nil
}

// GetBranchName gets branch name.
func (c *Client) GetBranchName() (string, error) {
	cmd := c.execCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", NewOpError("get branch name", "git rev-parse --abbrev-ref HEAD", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// RevParseVerify checks whether the given ref resolves to a valid object.
// It runs: git rev-parse --verify --quiet <ref>
func (c *Client) RevParseVerify(ref string) bool {
	cmd := c.execCommand("git", "rev-parse", "--verify", "--quiet", ref)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// GetCommitHash gets the short commit hash
func (c *Client) GetCommitHash() (string, error) {
	cmd := c.execCommand("git", "rev-parse", "--short", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "unknown", nil // Return "unknown" as fallback instead of error
	}
	return strings.TrimSpace(string(out)), nil
}

// GetUpstreamBranchName gets the upstream branch name for a given branch.
func (c *Client) GetUpstreamBranchName(branch string) (string, error) {
	cmd := c.execCommand("git", "rev-parse", "--abbrev-ref", branch+"@{upstream}")
	out, err := cmd.Output()
	if err != nil {
		return "", NewOpError("get upstream branch", "git rev-parse --abbrev-ref "+branch+"@{upstream}", err)
	}
	return strings.TrimSpace(string(out)), nil
}
