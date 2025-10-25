// Package router provides routing functionality for the ggc CLI tool with alias support.
package router

import (
	"fmt"
	"os"
	"strings"

	"github.com/bmf-san/ggc/v7/cmd"
	"github.com/bmf-san/ggc/v7/pkg/config"
)

// Router represents the command router with config support.
type Router struct {
	Executer      cmd.Executer
	ConfigManager *config.Manager
}

// NewRouter creates a new Router with a config manager.
func NewRouter(e cmd.Executer, cm *config.Manager) *Router {
	return &Router{
		Executer:      e,
		ConfigManager: cm,
	}
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
	allArgs := append([]string{name}, args...)
	r.Executer.Route(allArgs)
}
