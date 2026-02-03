package git

import (
	"fmt"
	"strings"
)

// ConfigGet retrieves a git configuration value from local repository
func (c *Client) ConfigGet(key string) (string, error) {
	cmd := c.execCommand("git", "config", key)
	out, err := cmd.Output()
	if err != nil {
		return "", NewOpError("config get", fmt.Sprintf("git config %s", key), err)
	}
	return strings.TrimSpace(string(out)), nil
}

// ConfigSet sets a git configuration value in local repository
func (c *Client) ConfigSet(key, value string) error {
	cmd := c.execCommand("git", "config", key, value)
	if err := cmd.Run(); err != nil {
		return NewOpError("config set", fmt.Sprintf("git config %s %s", key, value), err)
	}
	return nil
}

// ConfigGetGlobal retrieves a git configuration value from global config
func (c *Client) ConfigGetGlobal(key string) (string, error) {
	cmd := c.execCommand("git", "config", "--global", key)
	out, err := cmd.Output()
	if err != nil {
		return "", NewOpError("config get global", fmt.Sprintf("git config --global %s", key), err)
	}
	return strings.TrimSpace(string(out)), nil
}

// ConfigSetGlobal sets a git configuration value in global config
func (c *Client) ConfigSetGlobal(key, value string) error {
	cmd := c.execCommand("git", "config", "--global", key, value)
	if err := cmd.Run(); err != nil {
		return NewOpError("config set global", fmt.Sprintf("git config --global %s %s", key, value), err)
	}
	return nil
}
