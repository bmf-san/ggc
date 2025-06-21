package router

import (
	"github.com/bmf-san/ggc/cmd"
)

// Router is a struct that holds all the commands.
type Router struct {
	Executer cmd.Executer
}

// NewRouter creates a new Router.
func NewRouter(e cmd.Executer) *Router {
	return &Router{
		Executer: e,
	}
}

// Route routes the command to the appropriate function.
func (r *Router) Route(args []string) {
	if len(args) < 2 {
		r.Executer.ShowHelp()
		return
	}
	switch args[1] {
	case "__complete":
		r.Executer.Complete(args[2:])
		return
	case "branch":
		r.Executer.Branch(args[2:])
	case "push":
		r.Executer.Push(args[2:])
	case "pull":
		r.Executer.Pull(args[2:])
	case "log":
		r.Executer.Log(args[2:])
	case "commit":
		r.Executer.Commit(args[2:])
	case "add":
		r.Executer.Add(args[2:])
	case "fetch":
		r.Executer.Fetch(args[2:])
	case "clean":
		if len(args) > 2 && args[2] == "interactive" {
			r.Executer.CleanInteractive()
		} else {
			r.Executer.Clean(args[2:])
		}
	case "commit-push":
		r.Executer.CommitPushInteractive()
	case "stash":
		r.Executer.Stash(args[2:])
	case "rebase":
		r.Executer.Rebase(args[2:])
	case "remote":
		r.Executer.Remote(args[2:])
	case "add-commit-push":
		r.Executer.AddCommitPush()
	case "prp", "pull-rebase-push":
		r.Executer.PullRebasePush()
	case "sp", "stash-pull-pop":
		r.Executer.StashPullPop()
	case "reset-clean":
		r.Executer.Reset(args[2:])
	default:
		r.Executer.ShowHelp()
	}
}
