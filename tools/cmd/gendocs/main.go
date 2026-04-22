// Package main generates command documentation from the centralized registry.
package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/bmf-san/ggc/v8/cmd/command"
)

func main() {
	lines, err := readREADME()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading README: %v\n", err)
		os.Exit(1)
	}

	startIdx, endIdx, err := findCommandSection(lines)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding command section: %v\n", err)
		os.Exit(1)
	}

	tableStartIdx, tableEndIdx, err := findCommandTable(lines, startIdx, endIdx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding command table: %v\n", err)
		os.Exit(1)
	}

	newTable := generateCommandTable()

	if err := writeUpdatedREADME(lines, tableStartIdx, tableEndIdx, newTable); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing README: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("README.md command table updated successfully")

	if err := writeCommandsReference("docs/guide/commands.md"); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing commands reference: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("docs/guide/commands.md regenerated from registry")
}

func readREADME() ([]string, error) {
	file, err := os.Open("README.md")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func findCommandSection(lines []string) (int, int, error) {
	startIdx := -1
	endIdx := -1

	for i, line := range lines {
		if strings.Contains(line, "### Available Commands") {
			startIdx = i
		}
		if startIdx != -1 && startIdx < i && strings.HasPrefix(line, "### ") && !strings.Contains(line, "Available Commands") {
			endIdx = i
			break
		}
	}

	if startIdx == -1 {
		return 0, 0, fmt.Errorf("could not find '### Available Commands' section")
	}

	return startIdx, endIdx, nil
}

func findCommandTable(lines []string, startIdx, endIdx int) (int, int, error) {
	tableStartIdx := findTableStart(lines, startIdx, endIdx)
	if tableStartIdx == -1 {
		return 0, 0, fmt.Errorf("could not find command table")
	}

	tableEndIdx := findTableEnd(lines, tableStartIdx, endIdx)
	return tableStartIdx, tableEndIdx, nil
}

func findTableStart(lines []string, startIdx, endIdx int) int {
	for i := startIdx; i < len(lines) && (endIdx == -1 || i < endIdx); i++ {
		if strings.HasPrefix(lines[i], "| Command | Description |") {
			return i
		}
	}
	return -1
}

func findTableEnd(lines []string, tableStartIdx, endIdx int) int {
	for i := tableStartIdx + 2; i < len(lines); i++ { // Skip header and separator
		if (endIdx != -1 && i >= endIdx) ||
			(!strings.HasPrefix(lines[i], "|") && strings.TrimSpace(lines[i]) != "") {
			return i
		}
	}
	return len(lines)
}

func writeUpdatedREADME(lines []string, tableStartIdx, tableEndIdx int, newTable []string) error {
	var newLines []string
	newLines = append(newLines, lines[:tableStartIdx]...)
	newLines = append(newLines, newTable...)
	newLines = append(newLines, lines[tableEndIdx:]...)

	output, err := os.Create("README.md")
	if err != nil {
		return err
	}
	defer func() {
		_ = output.Close()
	}()

	for _, line := range newLines {
		if _, err := fmt.Fprintln(output, line); err != nil {
			return err
		}
	}

	return nil
}

func generateCommandTable() []string {
	registry := command.NewRegistry()
	commands := registry.VisibleCommands()

	// Sort commands by category, then by name
	sort.Slice(commands, func(i, j int) bool {
		if commands[i].Category != commands[j].Category {
			return command.CategoryOrder(commands[i].Category) < command.CategoryOrder(commands[j].Category)
		}
		return commands[i].Name < commands[j].Name
	})

	var table []string
	table = append(table, "| Command | Description |")
	table = append(table, "|--------|-------------|")

	for i := range commands {
		cmd := &commands[i]
		if len(cmd.Subcommands) == 0 {
			table = append(table, fmt.Sprintf("| `%s` | %s |", cmd.Name, cmd.Summary))
		} else {
			// Sort subcommands by name
			subcommands := make([]command.SubcommandInfo, 0, len(cmd.Subcommands))
			for _, sub := range cmd.Subcommands {
				if !sub.Hidden {
					subcommands = append(subcommands, sub)
				}
			}
			sort.Slice(subcommands, func(i, j int) bool {
				return subcommands[i].Name < subcommands[j].Name
			})

			for _, sub := range subcommands {
				table = append(table, fmt.Sprintf("| `%s` | %s |", sub.Name, sub.Summary))
			}
		}
	}

	return table
}

func writeCommandsReference(path string) error {
	registry := command.NewRegistry()
	commands := registry.VisibleCommands()
	sort.Slice(commands, func(i, j int) bool {
		if commands[i].Category != commands[j].Category {
			return command.CategoryOrder(commands[i].Category) < command.CategoryOrder(commands[j].Category)
		}
		return commands[i].Name < commands[j].Name
	})

	var b strings.Builder
	b.WriteString("# Commands\n\n")
	b.WriteString("This reference is auto-generated from the command registry in [`cmd/command/`](https://github.com/bmf-san/ggc/tree/main/cmd/command). Do not edit this file by hand; run `make docs`.\n\n")
	b.WriteString("For quick lookup, `ggc help` lists every command and `ggc help <command>` shows the same detail in your terminal.\n\n")
	b.WriteString("## Table of contents\n\n")

	byCategory := make(map[command.Category][]command.Info)
	for _, c := range commands {
		byCategory[c.Category] = append(byCategory[c.Category], c)
	}
	for _, cat := range command.OrderedCategories() {
		list := byCategory[cat]
		if len(list) == 0 {
			continue
		}
		anchor := strings.ToLower(string(cat))
		fmt.Fprintf(&b, "- [%s](#%s)\n", cat, anchor)
	}
	b.WriteString("\n")

	for _, cat := range command.OrderedCategories() {
		list := byCategory[cat]
		if len(list) == 0 {
			continue
		}
		fmt.Fprintf(&b, "## %s\n\n", cat)
		for i := range list {
			writeCommandSection(&b, &list[i])
		}
	}

	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func writeCommandSection(b *strings.Builder, c *command.Info) {
	fmt.Fprintf(b, "### `ggc %s`\n\n", c.Name)
	if c.Summary != "" {
		fmt.Fprintf(b, "%s.\n\n", strings.TrimSuffix(c.Summary, "."))
	}
	if len(c.Aliases) > 0 {
		quoted := make([]string, 0, len(c.Aliases))
		for _, a := range c.Aliases {
			quoted = append(quoted, "`"+a+"`")
		}
		fmt.Fprintf(b, "**Aliases:** %s\n\n", strings.Join(quoted, ", "))
	}
	if len(c.Usage) > 0 {
		b.WriteString("**Usage:**\n\n```bash\n")
		for _, u := range c.Usage {
			fmt.Fprintf(b, "%s\n", u)
		}
		b.WriteString("```\n\n")
	}
	subs := visibleSubs(c.Subcommands)
	if len(subs) > 0 {
		b.WriteString("**Subcommands:**\n\n")
		b.WriteString("| Subcommand | Description |\n")
		b.WriteString("|---|---|\n")
		sort.Slice(subs, func(i, j int) bool { return subs[i].Name < subs[j].Name })
		for _, s := range subs {
			fmt.Fprintf(b, "| `%s` | %s |\n", s.Name, s.Summary)
		}
		b.WriteString("\n")
		for _, s := range subs {
			if len(s.Examples) == 0 {
				continue
			}
			fmt.Fprintf(b, "_Examples for `%s`:_\n\n```bash\n", s.Name)
			for _, ex := range s.Examples {
				fmt.Fprintf(b, "%s\n", ex)
			}
			b.WriteString("```\n\n")
		}
	}
	if len(c.Examples) > 0 {
		b.WriteString("**Examples:**\n\n```bash\n")
		for _, ex := range c.Examples {
			fmt.Fprintf(b, "%s\n", ex)
		}
		b.WriteString("```\n\n")
	}
}

func visibleSubs(subs []command.SubcommandInfo) []command.SubcommandInfo {
	var out []command.SubcommandInfo
	for _, s := range subs {
		if !s.Hidden {
			out = append(out, s)
		}
	}
	return out
}
