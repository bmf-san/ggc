// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/bmf-san/ggc/v7/git"
)

// Remoter provides functionality for the remote command.
type Remoter struct {
	gitClient    git.RemoteManager
	outputWriter io.Writer
	helper       *Helper
}

// NewRemoter creates a new Remoter.
func NewRemoter(client git.RemoteManager) *Remoter {
	r := &Remoter{
		gitClient:    client,
		outputWriter: os.Stdout,
		helper:       NewHelper(),
	}
	r.helper.outputWriter = r.outputWriter
	return r
}

// Remote executes the remote command with the given arguments.
func (r *Remoter) Remote(args []string) {
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

func (r *Remoter) remoteList() {
	if err := r.gitClient.RemoteList(); err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
	}
}

func (r *Remoter) remoteAdd(name, url string) {
	if err := r.gitClient.RemoteAdd(name, url); err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
		return
	}
	_, _ = fmt.Fprintf(r.outputWriter, "Remote '%s' added\n", name)
}

func (r *Remoter) remoteRemove(name string) {
	if err := r.gitClient.RemoteRemove(name); err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
		return
	}
	_, _ = fmt.Fprintf(r.outputWriter, "Remote '%s' removed\n", name)
}

func (r *Remoter) remoteSetURL(name, url string) {
	if err := r.gitClient.RemoteSetURL(name, url); err != nil {
		_, _ = fmt.Fprintf(r.outputWriter, "Error: %v\n", err)
		return
	}
	_, _ = fmt.Fprintf(r.outputWriter, "Remote '%s' URL updated\n", name)
}
