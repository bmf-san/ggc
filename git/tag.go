package git

import (
	"os"
	"strings"
)

// TagList lists tags, optionally filtered by pattern.
func (c *Client) TagList(pattern []string) error {
	var cmd = c.execCommand("git", "tag", "--sort=-version:refname")
	if len(pattern) > 0 {
		args := append([]string{"tag", "--sort=-version:refname", "-l"}, pattern...)
		cmd = c.execCommand("git", args...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("tag list", "git tag --sort=-version:refname", err)
	}
	return nil
}

// TagCreate creates a lightweight tag.
func (c *Client) TagCreate(name string, commit string) error {
	var cmd = c.execCommand("git", "tag", name)
	if commit != "" {
		cmd = c.execCommand("git", "tag", name, commit)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("tag create", "git tag "+name, err)
	}
	return nil
}

// TagCreateAnnotated creates an annotated tag.
func (c *Client) TagCreateAnnotated(name, message string) error {
	var cmd = c.execCommand("git", "tag", "-a", name)
	if message != "" {
		cmd = c.execCommand("git", "tag", "-a", name, "-m", message)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("tag create annotated", "git tag -a "+name, err)
	}
	return nil
}

// TagDelete deletes tags.
func (c *Client) TagDelete(names []string) error {
	for _, name := range names {
		cmd := c.execCommand("git", "tag", "-d", name)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return NewError("tag delete", "git tag -d "+name, err)
		}
	}
	return nil
}

// TagPush pushes a specific tag to remote.
func (c *Client) TagPush(remote, name string) error {
	cmd := c.execCommand("git", "push", remote, name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("tag push", "git push "+remote+" "+name, err)
	}
	return nil
}

// TagPushAll pushes all tags to remote.
func (c *Client) TagPushAll(remote string) error {
	cmd := c.execCommand("git", "push", remote, "--tags")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("tag push all", "git push "+remote+" --tags", err)
	}
	return nil
}

// TagShow shows information about a tag.
func (c *Client) TagShow(name string) error {
	cmd := c.execCommand("git", "show", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("tag show", "git show "+name, err)
	}
	return nil
}

// GetLatestTag gets the latest tag.
func (c *Client) GetLatestTag() (string, error) {
	cmd := c.execCommand("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		return "", NewError("get latest tag", "git describe --tags --abbrev=0", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// TagExists checks if a tag exists.
func (c *Client) TagExists(name string) bool {
	cmd := c.execCommand("git", "tag", "-l", name)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) != ""
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
