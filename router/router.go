package router

import (
	"github.com/bmf-san/gcl/cmd"
)

func Route(args []string) {
	if len(args) < 2 {
		cmd.ShowHelp()
		return
	}
	switch args[1] {
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
	default:
		cmd.ShowHelp()
	}
}
