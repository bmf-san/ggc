// Package templates provides templates for help messages.
package templates

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"golang.org/x/term"

	commandregistry "github.com/bmf-san/ggc/v6/cmd/command"
)

// HelpData contains data for help message templates.
type HelpData struct {
	Logo        string
	Usage       string
	Description string
	Examples    []string
}

// Templates for help messages.
var (
	mainHelpTemplate = `{{.Logo}}ggc: A Go-based CLI tool to streamline Git operations

Usage:
  ggc <command> [subcommand] [options]

Main Commands:
{{range .Categories}}{{.Name}}:
{{range .Commands}}{{.Display}}
{{end}}
{{end}}Notes:
{{range .Notes}}  - {{.}}
{{end}}`

	commandHelpTemplate = `{{.Logo}}
Usage: {{.Usage}}

Description:
  {{.Description}}

Examples:
{{range .Examples}}  {{.}}
{{end}}
`
)

// MainHelpData contains data for main help message.
type MainHelpData struct {
	Logo       string
	Categories []helpCategory
	Notes      []string
}

type helpCategory struct {
	Name     string
	Commands []helpCommand
}

type helpCommand struct {
	Usage   string
	Summary string
	Display string
}

func selectLogo() string {
	if termWidth, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		if termWidth < 50 {
			return SmallLogo
		}
	}
	return Logo
}

// RenderMainHelp renders the main help message.
func RenderMainHelp() (string, error) {
	tmpl, err := template.New("mainHelp").Parse(mainHelpTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	data := MainHelpData{
		Logo:       selectLogo(),
		Categories: buildMainHelpCategories(),
		Notes: []string{
			"Unified syntax: no option flags (-/--) — use subcommands and words.",
			"To pass a literal that starts with '-', use the '--' separator: ggc commit -- - fix leading dash",
		},
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func buildMainHelpCategories() []helpCategory {
	orderedCats := commandregistry.OrderedCategories()

	categoryCommands := make(map[commandregistry.Category][]helpCommand)

	visibleCommands := commandregistry.VisibleCommands()
	for i := range visibleCommands {
		cmd := &visibleCommands[i]
		categoryCommands[cmd.Category] = append(categoryCommands[cmd.Category], helpCommandsFor(cmd)...)
	}

	var categories []helpCategory
	for _, cat := range orderedCats {
		commands := categoryCommands[cat]
		if len(commands) == 0 {
			continue
		}

		maxUsage := 0
		for _, cmd := range commands {
			if len(cmd.Usage) > maxUsage {
				maxUsage = len(cmd.Usage)
			}
		}

		for i := range commands {
			usage := commands[i].Usage
			summary := commands[i].Summary
			if summary != "" {
				commands[i].Display = fmt.Sprintf("    %-*s %s", maxUsage, usage, summary)
			} else {
				commands[i].Display = fmt.Sprintf("    %s", usage)
			}
		}

		categories = append(categories, helpCategory{
			Name:     string(cat),
			Commands: commands,
		})
	}

	return categories
}

func helpCommandsFor(info *commandregistry.Info) []helpCommand {
	var entries []helpCommand
	if len(info.Subcommands) == 0 {
		usage := firstUsage(info.Usage, "ggc "+info.Name)
		if shouldIncludeUsage(usage) {
			entries = append(entries, helpCommand{Usage: usage, Summary: info.Summary})
		}
		return entries
	}

	for _, sub := range info.Subcommands {
		if sub.Hidden {
			continue
		}
		usage := firstUsage(sub.Usage, "ggc "+sub.Name)
		if !shouldIncludeUsage(usage) {
			continue
		}
		entries = append(entries, helpCommand{Usage: usage, Summary: sub.Summary})
	}

	return dedupeHelpCommands(entries)
}

func firstUsage(usages []string, fallback string) string {
	for _, usage := range usages {
		trimmed := strings.TrimSpace(usage)
		if trimmed != "" {
			return trimmed
		}
	}
	return fallback
}

func shouldIncludeUsage(usage string) bool {
	return strings.HasPrefix(usage, "ggc ")
}

func dedupeHelpCommands(commands []helpCommand) []helpCommand {
	if len(commands) <= 1 {
		return commands
	}
	seen := make(map[string]struct{}, len(commands))
	var result []helpCommand
	for _, cmd := range commands {
		key := cmd.Usage + "\n" + cmd.Summary
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, cmd)
	}
	return result
}

// RenderCommandHelp renders help message for a specific command.
func RenderCommandHelp(data HelpData) (string, error) {
	tmpl, err := template.New("commandHelp").Parse(commandHelpTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
