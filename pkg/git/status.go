package git

// Status gets git status output.
func (c *Client) Status() (string, error) {
	cmd := c.execCommand("git", "status")
	out, err := cmd.Output()
	if err != nil {
		return "", NewError("get status", "git status", err)
	}
	return string(out), nil
}

// StatusShort gets git status --short output.
func (c *Client) StatusShort() (string, error) {
	cmd := c.execCommand("git", "status", "--short")
	out, err := cmd.Output()
	if err != nil {
		return "", NewError("get status short", "git status --short", err)
	}
	return string(out), nil
}

// StatusWithColor gets git status output with color.
func (c *Client) StatusWithColor() (string, error) {
	cmd := c.execCommand("git", "-c", "color.status=always", "status")
	out, err := cmd.Output()
	if err != nil {
		return "", NewError("get status with color", "git -c color.status=always status", err)
	}
	return string(out), nil
}

// StatusShortWithColor gets git status --short output with color.
func (c *Client) StatusShortWithColor() (string, error) {
	cmd := c.execCommand("git", "-c", "color.status=always", "status", "--short")
	out, err := cmd.Output()
	if err != nil {
		return "", NewError("get status short with color", "git -c color.status=always status --short", err)
	}
	return string(out), nil
}
