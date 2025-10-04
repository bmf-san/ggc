package git

// No imports needed for this file

// GetGitStatus gets git status.
func (c *Client) GetGitStatus() (string, error) {
	cmd := c.execCommand("git", "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return "", NewError("get status", "git status --porcelain", err)
	}
	return string(out), nil
}
