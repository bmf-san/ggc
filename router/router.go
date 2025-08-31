// Package router provides routing functionality for the ggc CLI tool with alias support.
package router

import (
	"fmt"
	"os"
	"strings"

	"github.com/bmf-san/ggc/v5/cmd"
	"github.com/bmf-san/ggc/v5/config"
)

// Router represents the command router with config support.
type Router struct {
	Executer      cmd.Executer
	ConfigManager *config.Manager
}

// NewRouter creates a new Router with a config manager.
func NewRouter(e cmd.Executer, cm *config.Manager) *Router {
	return &Router{Executer: e, ConfigManager: cm}
}

// Route routes the command to the appropriate handler
func (r *Router) Route(args []string) {
	if len(args) == 0 {
		r.Executer.Interactive()
		return
	}

	cmdName, cmdArgs := args[0], args[1:]

	if r.ConfigManager != nil && r.ConfigManager.GetConfig().IsAlias(cmdName) {
		r.executeAlias(cmdName, cmdArgs)
	} else {
		r.executeCommand(cmdName, cmdArgs)
	}
}

func (r *Router) executeAlias(name string, args []string) {
	cfg := r.ConfigManager.GetConfig()
	cmds, err := cfg.GetAliasCommands(name)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Alias '%s' error: %v\n", name, err)
		return
	}

	alias, err := cfg.ParseAlias(name)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error parsing alias: %v\n", err)
		return
	}

	switch alias.Type {
	case config.SimpleAlias:
		r.executeCommand(cmds[0], args)

	case config.SequenceAlias:
		if len(args) > 0 {
			_, _ = fmt.Fprintf(os.Stderr, "Warning: arguments ignored for sequence alias '%s'\n", name)
		}

		for _, c := range cmds {
			fmt.Printf("Executing: %s\n", c)
			command := strings.Split(c, " ")
			r.executeCommand(command[0], command[1:])
		}
	}
}

func (r *Router) executeCommand(name string, args []string) {
	switch name {
	case "help":
		r.Executer.Help()
	case "add":
		r.Executer.Add(args)
	case "branch":
		r.Executer.Branch(args)
	case "clean":
		r.Executer.Clean(args)
	case "commit":
		r.Executer.Commit(args)
	case "config":
		r.Executer.Config(args)
	case "diff":
		r.Executer.Diff(args)
	case "fetch":
		r.Executer.Fetch(args)
	case "hook":
		r.Executer.Hook(args)
	case "log":
		r.Executer.Log(args)
	case "pull":
		r.Executer.Pull(args)
	case "push":
		r.Executer.Push(args)
	case "rebase":
		r.Executer.Rebase(args)
	case "remote":
		r.Executer.Remote(args)
	case "reset":
		r.Executer.Reset(args)
	case "restore":
		r.Executer.Restore(args)
	case "stash":
		r.Executer.Stash(args)
	case "status":
		r.Executer.Status(args)
	case "tag":
		r.Executer.Tag(args)
	case "version":
		r.Executer.Version(args)
	default:
		r.Executer.Help()
	}
}
