// Package router provides routing functionality for the ggc CLI tool with alias support.
package router

import (
	"fmt"
	"strings"

	"github.com/bmf-san/ggc/v7/cmd"
	"github.com/bmf-san/ggc/v7/internal/config"
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
func (r *Router) Route(args []string) error {
	if len(args) == 0 {
		r.Executer.Interactive()
		return nil
	}

	cmdName, cmdArgs := args[0], args[1:]

	if r.ConfigManager != nil && r.ConfigManager.GetConfig().IsAlias(cmdName) {
		return r.executeAlias(cmdName, cmdArgs)
	}
	r.executeCommand(cmdName, cmdArgs)
	return nil
}

func (r *Router) executeAlias(name string, args []string) error {
	cfg := r.ConfigManager.GetConfig()
	alias, err := cfg.ParseAlias(name)
	if err != nil {
		return fmt.Errorf("error parsing alias: %w", err)
	}

	switch alias.Type {
	case config.SimpleAlias:
		// For simple aliases, process placeholders if any exist
		processedCommands, err := r.processPlaceholders(alias, args, name)
		if err != nil {
			return err
		}

		// Note: Using strings.Split may incorrectly split quoted arguments.
		// For example, "commit -m 'fix bug'" becomes ["commit", "-m", "'fix", "bug'"]
		// instead of ["commit", "-m", "'fix bug'"]. This is a known limitation.
		command := strings.Split(processedCommands[0], " ")
		if len(alias.Placeholders) == 0 {
			// No placeholders, forward user arguments
			r.executeCommand(command[0], args)
		} else {
			// Placeholders were processed, use the processed command
			r.executeCommand(command[0], command[1:])
		}

	case config.SequenceAlias:
		// Process placeholders for sequence aliases
		processedCommands, err := r.processPlaceholders(alias, args, name)
		if err != nil {
			return err
		}

		for _, c := range processedCommands {
			fmt.Printf("Executing: %s\n", c)
			// Note: Using strings.Split may incorrectly split quoted arguments.
			// This is a known limitation that affects commands with quoted parameters.
			command := strings.Split(c, " ")
			r.executeCommand(command[0], command[1:])
		}
	}
	return nil
}

// processPlaceholders processes placeholder replacement in alias commands
func (r *Router) processPlaceholders(alias *config.ParsedAlias, args []string, aliasName string) ([]string, error) {
	// If no placeholders are used, handle arguments appropriately
	if len(alias.Placeholders) == 0 {
		if alias.Type == config.SequenceAlias && len(args) > 0 {
			return nil, fmt.Errorf("sequence alias '%s' does not accept arguments (got %s)", aliasName, strings.Join(args, " "))
		}
		// For simple aliases without placeholders, arguments are forwarded as usual
		return alias.Commands, nil
	}

	// Validate that we have enough arguments for positional placeholders.
	// Note: MaxPositionalArg is 0-indexed (the highest placeholder index used),
	// so if MaxPositionalArg = 0, we need at least 1 argument (for {0}).
	if alias.MaxPositionalArg >= 0 && len(args) <= alias.MaxPositionalArg {
		return nil, fmt.Errorf("alias '%s' requires at least %d argument(s), got %d",
			aliasName, alias.MaxPositionalArg+1, len(args))
	}

	// Process each command
	processedCommands := make([]string, len(alias.Commands))
	for i, cmd := range alias.Commands {
		processed := cmd

		// Replace positional placeholders {0}, {1}, etc.
		for j := 0; j <= alias.MaxPositionalArg; j++ {
			placeholder := fmt.Sprintf("{%d}", j)
			if strings.Contains(processed, placeholder) {
				if j < len(args) {
					processed = strings.ReplaceAll(processed, placeholder, args[j])
				}
			}
		}

		// TODO: Support named placeholders in the future
		// For now, we only support positional placeholders

		processedCommands[i] = processed
	}

	return processedCommands, nil
}

func (r *Router) executeCommand(name string, args []string) {
	allArgs := append([]string{name}, args...)
	r.Executer.Route(allArgs)
}
