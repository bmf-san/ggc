package router

import (
	"github.com/bmf-san/ggc/cmd"
)

func Route(args []string) {
	if len(args) < 2 {
		cmd.ShowHelp()
		return
	}
	switch args[1] {
	case "__complete":
		cmd.Complete(args[2:])
		return
	case "branch":
		cmd.Branch(args[2:])
	case "push":
		cmd.Push(args[2:])
	case "pull":
		cmd.Pull(args[2:])
	case "log":
		cmd.Log(args[2:])
	case "commit":
		cmd.Commit(args[2:])
	case "add":
		cmd.Add(args[2:])
	case "fetch":
		cmd.Fetch(args[2:])
	case "clean":
		if len(args) > 2 && args[2] == "interactive" {
			cmd.CleanInteractive()
		} else {
			cmd.Clean(args[2:])
		}
	case "commit-push":
		cmd.CommitPushInteractive()
	case "stash":
		cmd.Stash(args[2:])
	case "rebase":
		cmd.Rebase(args[2:])
	case "remote":
		cmd.Remote(args[2:])
	case "add-commit-push":
		cmd.AddCommitPush()
	case "pull-rebase-push":
		cmd.PullRebasePush()
	case "stash-pull-pop":
		cmd.StashPullPop()
	case "reset-clean":
		cmd.ResetClean()
	default:
		cmd.ShowHelp()
	}
}
