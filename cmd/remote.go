// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Remoteer provides functionality for the remote command.
type Remoteer struct {
	execCommand  func(name string, arg ...string) *exec.Cmd
	outputWriter io.Writer
	helper       *Helper
}

// NewRemoteer creates a new Remoteer.
func NewRemoteer() *Remoteer {
	r := &Remoteer{
		execCommand:  exec.Command,
		outputWriter: os.Stdout,
		helper:       NewHelper(),
	}
	r.helper.outputWriter = r.outputWriter
	return r
}

// Remote executes the remote command with the given arguments.
func (r *Remoteer) Remote(args []string) {
	if len(args) == 0 {
		r.helper.ShowRemoteHelp()
		return
	}

	switch args[0] {
	case "list":
		r.remoteList()
	case "add":
		if len(args) != 3 {
			r.helper.ShowRemoteHelp()
			return
		}
		r.remoteAdd(args[1], args[2])
	case "remove":
		if len(args) != 2 {
			r.helper.ShowRemoteHelp()
			return
		}
		r.remoteRemove(args[1])
	case "set-url":
		if len(args) != 3 {
			r.helper.ShowRemoteHelp()
			return
		}
		r.remoteSetURL(args[1], args[2])
	default:
		r.helper.ShowRemoteHelp()
	}
}

func (r *Remoteer) remoteList() {
	cmd := r.execCommand("git", "remote", "-v")
	cmd.Stdout = r.outputWriter
	cmd.Stderr = r.outputWriter
	if err := cmd.Run(); err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: failed to list remotes: %v\n", err)
		return
	}
}

func (r *Remoteer) remoteAdd(name, url string) {
	cmd := r.execCommand("git", "remote", "add", name, url)
	cmd.Stdout = r.outputWriter
	cmd.Stderr = r.outputWriter
	if err := cmd.Run(); err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: failed to add remote: %v\n", err)
		return
	}
	_, _ = fmt.Fprintf(r.outputWriter, "Remote '%s' added\n", name)
}

func (r *Remoteer) remoteRemove(name string) {
	cmd := r.execCommand("git", "remote", "remove", name)
	cmd.Stdout = r.outputWriter
	cmd.Stderr = r.outputWriter
	if err := cmd.Run(); err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: failed to remove remote: %v\n", err)
		return
	}
	_, _ = fmt.Fprintf(r.outputWriter, "Remote '%s' removed\n", name)
}

func (r *Remoteer) remoteSetURL(name, url string) {
	cmd := r.execCommand("git", "remote", "set-url", name, url)
	cmd.Stdout = r.outputWriter
	cmd.Stderr = r.outputWriter
	if err := cmd.Run(); err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: failed to set remote URL: %v\n", err)
		return
	}
	_, _ = fmt.Fprintf(r.outputWriter, "Remote '%s' URL updated\n", name)
}
