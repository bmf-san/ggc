package git

// StatusReader provides read-only status output with color support.
type StatusReader interface {
	StatusWithColor() (string, error)
	StatusShortWithColor() (string, error)
}

// BranchUpstreamReader provides information about the current branch and its upstream.
type BranchUpstreamReader interface {
	GetCurrentBranch() (string, error)
	GetUpstreamBranchName(branch string) (string, error)
	GetAheadBehindCount(branch, upstream string) (string, error)
}

// StatusInfoReader is a pragmatic composite for the status command dependencies.
// It avoids pulling in an overly broad client surface area.
type StatusInfoReader interface {
	StatusReader
	BranchUpstreamReader
}

// Status gets git status output.
func (c *Client) Status() (string, error) {
	cmd := c.execCommand("git", "status")
	out, err := cmd.Output()
	if err != nil {
		return "", NewOpError("get status", "git status", err)
	}
	return string(out), nil
}

// StatusShort gets git status --short output.
func (c *Client) StatusShort() (string, error) {
	cmd := c.execCommand("git", "status", "--short")
	out, err := cmd.Output()
	if err != nil {
		return "", NewOpError("get status short", "git status --short", err)
	}
	return string(out), nil
}

// StatusWithColor gets git status output with color.
func (c *Client) StatusWithColor() (string, error) {
	cmd := c.execCommand("git", "-c", "color.status=always", "status")
	out, err := cmd.Output()
	if err != nil {
		return "", NewOpError("get status with color", "git -c color.status=always status", err)
	}
	return string(out), nil
}

// StatusShortWithColor gets git status --short output with color.
func (c *Client) StatusShortWithColor() (string, error) {
	cmd := c.execCommand("git", "-c", "color.status=always", "status", "--short")
	out, err := cmd.Output()
	if err != nil {
		return "", NewOpError("get status short with color", "git -c color.status=always status --short", err)
	}
	return string(out), nil
}
