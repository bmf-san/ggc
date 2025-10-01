// Package main generates command documentation from the centralized registry.
package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/bmf-san/ggc/v7/cmd/command"
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
	commands := command.VisibleCommands()

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
