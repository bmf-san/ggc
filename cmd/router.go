package cmd

import (
	"fmt"
	"sort"
	"strings"

	commandregistry "github.com/bmf-san/ggc/v8/cmd/command"
)

// commandRouter dispatches a command name (plus its args) to the matching
// handler on Cmd. It is initialized once at NewCmd time and then consulted
// for every command the user types, both in scripted and interactive mode.
type commandRouter struct {
	registry *commandregistry.Registry
	handlers map[string]func([]string)
}

// newCommandRouter builds the handler map and validates that every
// non-hidden command in the registry has a handler. The map-key is the
// canonical command name as exposed by the registry; this keeps shell
// completions, README.md and the router in strict lockstep.
func newCommandRouter(cmd *Cmd) (*commandRouter, error) {
	if err := cmd.registry.Validate(); err != nil {
		return nil, fmt.Errorf("command registry validation failed: %w", err)
	}

	handlers := map[string]func([]string){
		"help":       func(args []string) { cmd.Help(args) },
		"add":        func(args []string) { cmd.Add(args) },
		"branch":     func(args []string) { cmd.Branch(args) },
		"commit":     func(args []string) { cmd.Commit(args) },
		"log":        func(args []string) { cmd.Log(args) },
		"pull":       func(args []string) { cmd.Pull(args) },
		"push":       func(args []string) { cmd.Push(args) },
		"reset":      func(args []string) { cmd.Reset(args) },
		"clean":      func(args []string) { cmd.Clean(args) },
		"version":    func(args []string) { cmd.Version(args) },
		"remote":     func(args []string) { cmd.Remote(args) },
		"rebase":     func(args []string) { cmd.Rebase(args) },
		"stash":      func(args []string) { cmd.Stash(args) },
		"config":     func(args []string) { cmd.Config(args) },
		"hook":       func(args []string) { cmd.Hook(args) },
		"tag":        func(args []string) { cmd.Tag(args) },
		"status":     func(args []string) { cmd.Status(args) },
		"fetch":      func(args []string) { cmd.Fetch(args) },
		"diff":       func(args []string) { cmd.Diff(args) },
		"restore":    func(args []string) { cmd.Restore(args) },
		"debug-keys": func(args []string) { cmd.DebugKeys(args) },
		interactiveQuitCommand: func([]string) {
			_, _ = fmt.Fprintln(cmd.outputWriter, "The 'quit' command is only available in interactive mode.")
		},
	}

	available := make(map[string]struct{}, len(handlers))
	for key := range handlers {
		available[key] = struct{}{}
	}

	missing := missingHandlers(cmd.registry, available)
	if len(missing) > 0 {
		sort.Strings(missing)
		return nil, fmt.Errorf("no handler registered for commands: %s", strings.Join(missing, ", "))
	}

	return &commandRouter{registry: cmd.registry, handlers: handlers}, nil
}

// route looks up cmd in the registry (which handles aliases and canonical
// names), then dispatches to the registered handler. It returns false when
// the command is unknown to the registry, or when no handler is registered
// for the canonical name, so the caller can fall back to alias or error
// paths.
func (r *commandRouter) route(cmd string, args []string) bool {
	info, ok := r.registry.Find(cmd)
	if !ok {
		return false
	}
	// Use the canonical command name from the registry as the handler key.
	handler, ok := r.handlers[info.Name]
	if !ok {
		return false
	}
	handler(args)
	return true
}

// missingHandlers returns every non-hidden registry command that has no
// matching handler in available. It is used at startup to turn a registry
// drift into a loud construction error instead of a silent "unknown command"
// at runtime.
func missingHandlers(registry *commandregistry.Registry, available map[string]struct{}) []string {
	var missing []string
	allCommands := registry.All()
	for i := range allCommands {
		info := &allCommands[i]
		if info.Hidden {
			continue
		}
		if _, ok := available[info.Name]; !ok {
			missing = append(missing, info.Name)
		}
	}
	return missing
}
