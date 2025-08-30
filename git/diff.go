package git

// Diff gets git diff output.
func (c *Client) Diff() (string, error) {
	cmd := c.execCommand("git", "diff")
	out, err := cmd.Output()
	if err != nil {
		return "", NewError("get diff", "git diff", err)
	}
	return string(out), nil
}

// DiffStaged gets git diff --staged output.
func (c *Client) DiffStaged() (string, error) {
	cmd := c.execCommand("git", "diff", "--staged")
	out, err := cmd.Output()
	if err != nil {
		return "", NewError("get diff staged", "git diff --staged", err)
	}
	return string(out), nil
}

// DiffHead gets git diff HEAD output.
func (c *Client) DiffHead() (string, error) {
	cmd := c.execCommand("git", "diff", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", NewError("get diff HEAD", "git diff HEAD", err)
	}
	return string(out), nil
}
