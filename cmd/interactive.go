// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// CommandInfo contiene información sobre un comando
type CommandInfo struct {
	Command     string
	Description string
}

// Lista unificada de comandos y descripciones
var commandList = []CommandInfo{
	{"add <file>", "Add specific file to index"},
	{"add .", "Add all changes to index"},
	{"add -p", "Add changes interactively"},
	{"branch current", "Show current branch name"},
	{"branch checkout", "Switch to existing branch"},
	{"branch checkout-remote", "Create and checkout local branch from remote"},
	{"branch create", "Create and checkout new branch"},
	{"branch delete", "Delete local branch"},
	{"branch delete-merged", "Delete merged local branches"},
	{"push current", "Push current branch to remote"},
	{"push force", "Force push current branch"},
	{"pull current", "Pull current branch from remote"},
	{"pull rebase", "Pull with rebase"},
	{"log simple", "Show simple history"},
	{"log graph", "Show history with graph"},
	{"commit <message>", "Create commit with message"},
	{"commit allow-empty", "Create empty commit"},
	{"commit tmp", "Create temporary commit"},
	{"commit amend <message>", "Amend previous commit"},
	{"fetch --prune", "Fetch and clean stale references"},
	{"tag list", "List all tags"},
	{"tag annotated <tag> <message>", "Create annotated tag"},
	{"tag delete <tag>", "Delete tag"},
	{"tag show <tag>", "Show tag information"},
	{"tag push", "Push tags to remote"},
	{"tag create <tag>", "Create tag"},
	{"diff", "Show differences"},
	{"diff unstaged", "Show unstaged changes"},
	{"diff staged", "Show staged changes"},
	{"version", "Show current version"},
	{"clean files", "Clean untracked files"},
	{"clean dirs", "Clean untracked directories"},
	{"clean-interactive", "Clean files interactively"},
	{"reset-clean", "Reset and clean"},
	{"commit-push-interactive", "Interactive commit and push"},
	{"stash trash", "Delete stash"},
	{"status", "Show working tree status"},
	{"status short", "Show concise status"},
	{"rebase interactive", "Interactive rebase"},
	{"remote list", "List remote repositories"},
	{"remote add <n> <url>", "Add remote repository"},
	{"remote remove <n>", "Remove remote repository"},
	{"remote set-url <n> <url>", "Change remote repository URL"},
	{"add-commit-push", "Add, commit and push in one operation"},
	{"pull-rebase-push", "Pull with rebase and push"},
	{"stash-pull-pop", "Stash, pull and pop in sequence"},
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

		// Filtrado
		filtered := []CommandInfo{}
		for _, cmd := range commandList {
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

			// Encontrar el comando más largo para alineación
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
				paddingLen := maxCmdLen - len(cmd.Command)
				if paddingLen < 0 {
					paddingLen = 0
				}
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
