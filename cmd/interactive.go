// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"

	"golang.org/x/term"
)

// CommandInfo contains the name and description of the command
type CommandInfo struct {
	Command     string
	Description string
}

var commands = []CommandInfo{
	{"add <file>", "Add a specific file to the index"},
	{"add .", "Add all changes to index"},
	{"add -p", "Add changes interactively"},
	{"branch current", "Show current branch name"},
	{"branch checkout", "Switch to an existing branch"},
	{"branch checkout-remote", "Create and checkout a local branch from the remote"},
	{"branch create", "Create and checkout new branch"},
	{"branch delete", "Delete local branch"},
	{"branch delete-merged", "Delete local merged branch"},
	{"push current", "Push current branch from remote repository"},
	{"push force", "Force push current branch"},
	{"pull current", "Pull current branch from remote repository"},
	{"pull rebase", "Pull and rebase"},
	{"log simple", "Show simple historical log"},
	{"log graph", "Show log with graph"},
	{"commit <message>", "Create commit with a message"},
	{"commit allow-empty", "Create an empty commit"},
	{"commit tmp", "Create a temporary commit"},
	{"commit amend <message>", "Amend a previous commit"},
	{"fetch --prune", "Fetch and clean stale references"},
	{"tag list", "List all tags"},
	{"tag annotated <tag> <message>", "Create annotated tag"},
	{"tag delete <tag>", "Delete tag"},
	{"tag show <tag>", "Show tag information"},
	{"tag push", "Push tags to remote"},
	{"tag create <tag>", "Create tag"},
	{"config list", "List all configuration"},
	{"config get <key>", "Get a specific config value"},
	{"config set <key> <value>", "Set a configuration value"},
	{"hook list", "List all hooks"},
	{"hook install <hook>", "Install a hook"},
	{"hook enable <hook>", "Enable/Turn on a hook"},
	{"hook disable <hook>", "Disable/Turn off a hook"},
	{"hook uninstall <hook>", "Uninstall an existing hook"},
	{"hook edit <hook>", "Edit a hook's contents"},
	{"diff", "Show changes (git diff HEAD)"},
	{"diff unstaged", "Show unstaged changes"},
	{"diff staged", "Show staged changes"},
	{"version", "Show current version"},
	{"clean files", "Clean untracked files"},
	{"clean dirs", "Clean untracked directories"},
	{"clean-interactive", "Clean files interactively"},
	{"stash trash", "Delete stash"},
	{"status", "Show working tree status"},
	{"status short", "Show concise status (porcelain format)"},
	{"rebase interactive", "Interactive rebase"},
	{"remote list", "List all remote repositories"},
	{"remote add <name> <url>", "Add remote repository"},
	{"remote remove <name>", "Remove remote repository"},
	{"remote set-url <name> <url>", "Change remote URL"},
	{"quit", "Exit interactive mode"},
}

// InteractiveUI provides an incremental search interactive UI for command selection.
// Returns the selected command as []string (nil if nothing selected)
func InteractiveUI() []string {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Println("Failed to set terminal to raw mode:", err)
		return nil
	}
	defer func() {
		if err := term.Restore(fd, oldState); err != nil {
			fmt.Fprintln(os.Stderr, "failed to restore terminal state:", err)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	selected := 0
	input := ""

	for {
		if _, err := os.Stdout.Write([]byte("\033[H\033[2J\033[H")); err != nil {
			fmt.Fprintln(os.Stderr, "failed to write clear screen sequence:", err)
		}
		fmt.Printf("Select a command (incremental search: type to filter, ctrl+n: down, ctrl+p: up, enter: execute, ctrl+c: quit)\n")
		fmt.Printf("\rSearch: %s\n\n", input)

		// Filtering
		filtered := []CommandInfo{}
		for _, cmd := range commands {
			if strings.Contains(cmd.Command, input) {
				filtered = append(filtered, cmd)
			}
		}
		if input == "" {
			fmt.Println("(Type to filter commands...)")
		} else {
			if len(filtered) == 0 {
				fmt.Println("  (No matching command)")
			}
			if selected >= len(filtered) {
				selected = len(filtered) - 1
			}
			if selected < 0 {
				selected = 0
			}
			maxCmdLen := 0
			for _, cmd := range filtered {
				if len(cmd.Command) > maxCmdLen {
					maxCmdLen = len(cmd.Command)
				}
			}

			for i, cmd := range filtered {
				desc := cmd.Description
				if desc == "" {
					desc = "No description"
				}
				paddingLen := int(math.Max(0, float64(maxCmdLen-len(cmd.Command))))
				padding := strings.Repeat(" ", paddingLen)
				if i == selected {
					fmt.Printf("\r> %s%s  %s\n", cmd.Command, padding, desc)
				} else {
					fmt.Printf("\r  %s%s  %s\n", cmd.Command, padding, desc)
				}
			}
		}
		fmt.Print("\n\r") // Ensure next output starts at left edge

		b, err := reader.ReadByte()
		if err != nil {
			continue
		}
		if b == 3 { // Ctrl+C in raw mode
			if err := term.Restore(fd, oldState); err != nil {
				fmt.Fprintln(os.Stderr, "failed to restore terminal state:", err)
			}
			fmt.Println("\nExiting...")
			os.Exit(0)
		} else if b == 13 { // Enter
			if input == "" {
				continue
			}
			if len(filtered) > 0 {
				fmt.Printf("\nExecute: %s\n", filtered[selected].Command)
				if err := term.Restore(fd, oldState); err != nil {
					fmt.Fprintln(os.Stderr, "failed to restore terminal state:", err)
				}
				// Placeholder detection
				cmdTemplate := filtered[selected].Command
				placeholders := extractPlaceholders(cmdTemplate)
				inputs := make(map[string]string)
				readerStdin := bufio.NewReader(os.Stdin)
				for _, ph := range placeholders {
					fmt.Print("\n\r") // Newline + carriage return
					fmt.Printf("Enter value for %s: ", ph)
					val, _ := readerStdin.ReadString('\n')
					val = strings.TrimSpace(val)
					inputs[ph] = val
				}
				// Placeholder replacement
				finalCmd := cmdTemplate
				for ph, val := range inputs {
					finalCmd = strings.ReplaceAll(finalCmd, "<"+ph+">", val)
				}
				args := []string{"ggc"}
				args = append(args, strings.Fields(finalCmd)...)
				return args
			}
			break
		} else if b == 16 { // Ctrl+p
			if selected > 0 {
				selected--
			}
		} else if b == 14 { // Ctrl+n
			if selected < len(filtered)-1 {
				selected++
			}
		} else if b == 127 || b == 8 { // Backspace
			if len(input) > 0 {
				input = input[:len(input)-1]
			}
		} else if b >= 32 && b <= 126 { // Printable ASCII
			input += string(b)
		}
	}
	return nil
}

// Extract <...> placeholders from a string
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
