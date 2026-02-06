// Package command provides centralized command registry and metadata management.
package command

import (
	"fmt"
	"strings"
)

// Registry manages command metadata.
type Registry struct {
	commands []Info
}

// NewRegistry creates a new Registry with default commands.
func NewRegistry() *Registry {
	return &Registry{
		commands: defaultCommands(),
	}
}

// NewRegistryWith creates a Registry with custom commands (for testing).
func NewRegistryWith(commands []Info) *Registry {
	return &Registry{commands: commands}
}

// All returns a defensive copy of all commands.
func (r *Registry) All() []Info {
	out := make([]Info, len(r.commands))
	for i := range r.commands {
		out[i] = (&r.commands[i]).clone()
	}
	return out
}

// Find returns the command metadata by name.
func (r *Registry) Find(name string) (Info, bool) {
	for i := range r.commands {
		if strings.EqualFold(r.commands[i].Name, name) {
			return (&r.commands[i]).clone(), true
		}
	}
	return Info{}, false
}

// VisibleCommands returns non-hidden commands.
func (r *Registry) VisibleCommands() []Info {
	var out []Info
	for i := range r.commands {
		if r.commands[i].Hidden {
			continue
		}
		out = append(out, (&r.commands[i]).clone())
	}
	return out
}

// Validate ensures registry consistency.
func (r *Registry) Validate() error {
	return Validate(r.commands)
}

// defaultCommands returns the default command set.
func defaultCommands() []Info {
	var commands []Info
	commands = append(commands, basics()...)
	commands = append(commands, branch()...)
	commands = append(commands, remote()...)
	commands = append(commands, commit()...)
	commands = append(commands, tag()...)
	commands = append(commands, config()...)
	commands = append(commands, hook()...)
	commands = append(commands, diff()...)
	commands = append(commands, utility()...)
	commands = append(commands, cleanup()...)
	commands = append(commands, stash()...)
	commands = append(commands, status()...)
	commands = append(commands, rebase()...)
	return commands
}

// Validate ensures the provided command metadata is internally consistent.
func Validate(commands []Info) error {
	seen := make(map[string]struct{})
	for i := range commands {
		cmd := &commands[i]
		if err := validateCommand(cmd, seen); err != nil {
			return err
		}
	}
	return nil
}

func validateCommand(cmd *Info, seen map[string]struct{}) error {
	if strings.TrimSpace(cmd.Name) == "" {
		return fmt.Errorf("command name cannot be empty")
	}
	key := strings.ToLower(cmd.Name)
	if _, ok := seen[key]; ok {
		return fmt.Errorf("duplicate command name: %s", cmd.Name)
	}
	seen[key] = struct{}{}
	if strings.TrimSpace(cmd.Summary) == "" {
		return fmt.Errorf("command summary missing for %s", cmd.Name)
	}
	if !cmd.Hidden && strings.TrimSpace(cmd.HandlerID) == "" {
		return fmt.Errorf("handler ID missing for %s", cmd.Name)
	}

	return validateSubcommands(cmd)
}

func validateSubcommands(cmd *Info) error {
	subSeen := make(map[string]struct{})
	for _, sub := range cmd.Subcommands {
		if strings.TrimSpace(sub.Name) == "" {
			return fmt.Errorf("subcommand name cannot be empty for %s", cmd.Name)
		}
		subKey := strings.ToLower(sub.Name)
		if _, ok := subSeen[subKey]; ok {
			return fmt.Errorf("duplicate subcommand %s under %s", sub.Name, cmd.Name)
		}
		subSeen[subKey] = struct{}{}
		if strings.TrimSpace(sub.Summary) == "" {
			return fmt.Errorf("subcommand summary missing for %s -> %s", cmd.Name, sub.Name)
		}
	}
	return nil
}
