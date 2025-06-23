// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/term"
)

var commands = []string{
	"add <file>",
	"add .",
	"add -p",
	"branch current",
	"branch checkout",
	"branch checkout-remote",
	"branch delete",
	"branch delete-merged",
	"push current",
	"push force",
	"pull current",
	"pull rebase",
	"log simple",
	"log graph",
	"commit <message>",
	"commit allow-empty",
	"commit tmp",
	"fetch --prune",
	"clean files",
	"clean dirs",
	"clean-interactive",
	"reset-clean",
	"commit-push-interactive",
	"stash trash",
	"rebase interactive",
	"remote list",
	"remote add <n> <url>",
	"remote remove <n>",
	"remote set-url <n> <url>",
	"add-commit-push",
	"pull-rebase-push",
	"stash-pull-pop",
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

	// For Ctrl+C exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if err := term.Restore(fd, oldState); err != nil {
			fmt.Fprintln(os.Stderr, "failed to restore terminal state:", err)
		}
		os.Exit(0)
	}()

	reader := bufio.NewReader(os.Stdin)
	selected := 0
	input := ""

	for {
		if _, err := os.Stdout.Write([]byte("\033[H\033[2J\033[H")); err != nil {
			fmt.Fprintln(os.Stderr, "failed to write clear screen sequence:", err)
		}
		fmt.Printf("Select a command (incremental search: type to filter, ctrl+n: down, ctrl+p: up, Enter: execute, Ctrl+C: quit)\n")
		fmt.Printf("\rSearch: %s\n\n", input)

		// Filtering
		filtered := []string{}
		for _, cmd := range commands {
			if strings.Contains(cmd, input) {
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
			for i, cmd := range filtered {
				if i == selected {
					fmt.Printf("\r> %s\n", cmd)
				} else {
					fmt.Printf("\r  %s\n", cmd)
				}
			}
		}
		fmt.Print("\n\r") // Ensure next output starts at left edge

		b, err := reader.ReadByte()
		if err != nil {
			continue
		}
		if b == 13 { // Enter
			if len(filtered) > 0 {
				fmt.Printf("\nExecute: %s\n", filtered[selected])
				if err := term.Restore(fd, oldState); err != nil {
					fmt.Fprintln(os.Stderr, "failed to restore terminal state:", err)
				}
				// Placeholder detection
				cmdTemplate := filtered[selected]
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
