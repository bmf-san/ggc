package interactive

import (
	commandregistry "github.com/bmf-san/ggc/v7/cmd/command"
)

// CommandInfo contains the name and description of the command
type CommandInfo struct {
	Command     string
	Description string
}

var commands = buildInteractiveCommands()

// buildInteractiveCommands builds the list of commands available in interactive mode
func buildInteractiveCommands() []CommandInfo {
	var list []CommandInfo
	registry := commandregistry.NewRegistry()
	allCommands := registry.All()
	for i := range allCommands {
		cmd := &allCommands[i]
		if cmd.Hidden {
			continue
		}
		if len(cmd.Subcommands) == 0 {
			list = append(list, CommandInfo{Command: cmd.Name, Description: cmd.Summary})
			continue
		}
		for _, sub := range cmd.Subcommands {
			if sub.Hidden {
				continue
			}
			list = append(list, CommandInfo{Command: sub.Name, Description: sub.Summary})
		}
	}
	return list
}

// extractPlaceholders extracts <...> placeholders from a string
func extractPlaceholders(s string) []string {
	var res []string
	start := -1
	for i, c := range s {
		if c == '<' {
			start = i + 1
		} else if c == '>' && start != -1 {
			res = append(res, s[start:i])
			start = -1
		}
	}
	return res
}
