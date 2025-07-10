// Package router provides routing functionality for the ggc CLI tool.
package router

import (
	"github.com/bmf-san/ggc/cmd"
)

// Router represents the command router.
type Router struct {
	Executer cmd.Executer
}

// NewRouter creates a new Router.
func NewRouter(e cmd.Executer) *Router {
	return &Router{
		Executer: e,
	}
}

// Route routes the command to the appropriate handler.
func (r *Router) Route(args []string) {
	if len(args) == 0 {
		r.Executer.Interactive()
		return
	}

	switch args[0] {
	case "help":
		r.Executer.Help()
	case "branch":
		r.Executer.Branch(args[1:])
	case "commit":
		r.Executer.Commit(args[1:])
	case "log":
		r.Executer.Log(args[1:])
	case "pull":
		r.Executer.Pull(args[1:])
	case "push":
		r.Executer.Push(args[1:])
	case "reset":
		r.Executer.Reset(args[1:])
    case "diff":
		r.Executer.Diff(args[1:])
	case "tag":
		r.Executer.Tag(args[1:])
    case "status":
		r.Executer.Status(args[1:])
	case "clean":
		r.Executer.Clean(args[1:])
	case "pull-rebase-push":
		r.Executer.PullRebasePush()
	default:
		r.Executer.Help()
	}
}
