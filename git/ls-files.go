package git

// ListFiles lists all files managed by git.
func (c *Client) ListFiles() (string, error) {
	cmd := c.execCommand("git", "ls-files")
	out, err := cmd.Output()
	if err != nil {
		return "", NewError("list files", "git ls-files", err)
	}
	return string(out), nil
}
