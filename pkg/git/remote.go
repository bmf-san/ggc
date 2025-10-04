package git

import (
	"os"
)

// RemoteList lists all remotes.
func (c *Client) RemoteList() error {
	cmd := c.execCommand("git", "remote", "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("remote list", "git remote -v", err)
	}
	return nil
}

// RemoteAdd adds a new remote.
func (c *Client) RemoteAdd(name, url string) error {
	cmd := c.execCommand("git", "remote", "add", name, url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("remote add", "git remote add "+name+" "+url, err)
	}
	return nil
}

// RemoteRemove removes a remote.
func (c *Client) RemoteRemove(name string) error {
	cmd := c.execCommand("git", "remote", "remove", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("remote remove", "git remote remove "+name, err)
	}
	return nil
}

// RemoteSetURL sets the URL for a remote.
func (c *Client) RemoteSetURL(name, url string) error {
	cmd := c.execCommand("git", "remote", "set-url", name, url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("remote set-url", "git remote set-url "+name+" "+url, err)
	}
	return nil
}
