// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"bufio"
	"fmt"
	"io"
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

// UI represents the interface for terminal UI operations
type UI struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	term   terminal
}

// terminal represents terminal operations
type terminal interface {
	makeRaw(fd int) (*term.State, error)
	restore(fd int, state *term.State) error
}

type defaultTerminal struct{}

func (t *defaultTerminal) makeRaw(fd int) (*term.State, error) {
	return term.MakeRaw(fd)
}

func (t *defaultTerminal) restore(fd int, state *term.State) error {
	return term.Restore(fd, state)
}

// NewUI creates a new UI with default settings
func NewUI() *UI {
	return &UI{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
		term:   &defaultTerminal{},
	}
}

var commands = []CommandInfo{
	{"help", "Show help message"},
	{"add <file>", "Add a specific file to the index"},
	{"add .", "Add all changes to index"},
	{"add interactive", "Add changes interactively"},
	{"add -p", "Add changes interactively (patch mode)"},
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
	{"commit amend", "Amend previous commit (editor)"},
	{"commit amend --no-edit", "Amend without editing commit message"},
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
	{"stash", "Stash current changes"},
	{"stash list", "List all stashes"},
	{"stash show", "Show changes in stash"},
	{"stash show <stash>", "Show changes in specific stash"},
	{"stash apply", "Apply stash without removing it"},
	{"stash apply <stash>", "Apply specific stash without removing it"},
	{"stash pop", "Apply and remove the latest stash"},
	{"stash pop <stash>", "Apply and remove specific stash"},
	{"stash drop", "Remove the latest stash"},
	{"stash drop <stash>", "Remove specific stash"},
	{"stash branch <branch>", "Create branch from stash"},
	{"stash branch <branch> <stash>", "Create branch from specific stash"},
	{"stash push", "Save changes to new stash"},
	{"stash push -m <message>", "Save changes to new stash with message"},
	{"stash save <message>", "Save changes to new stash with message"},
	{"stash clear", "Remove all stashes"},
	{"stash create", "Create stash and return object name"},
	{"stash store <object>", "Store stash object"},
	{"status", "Show working tree status"},
	{"status short", "Show concise status (porcelain format)"},
	{"rebase -i", "Interactive rebase"},
	{"rebase --interactive", "Interactive rebase"},
	{"remote list", "List all remote repositories"},
	{"remote add <name> <url>", "Add remote repository"},
	{"remote remove <name>", "Remove remote repository"},
	{"remote set-url <name> <url>", "Change remote URL"},
	{"restore <file>", "Restore file in working directory from index"},
	{"restore .", "Restore all files in working directory from index"},
	{"restore staged <file>", "Unstage file (restore from HEAD to index)"},
	{"restore staged .", "Unstage all files"},
	{"restore <commit> <file>", "Restore file from specific commit"},
	{"quit", "Exit interactive mode"},
}

// InteractiveUI provides an incremental search interactive UI for command selection.
// Returns the selected command as []string (nil if nothing selected)
func InteractiveUI() []string {
	ui := NewUI()
	return ui.Run()
}

// writeError writes an error message to stderr
func (ui *UI) writeError(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(ui.stderr, format+"\n", a...)
}

// write writes a message to stdout
func (ui *UI) write(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(ui.stdout, format, a...)
}

// writeln writes a message with newline to stdout
func (ui *UI) writeln(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(ui.stdout, format+"\n", a...)
}

// Run executes the interactive UI
func (ui *UI) Run() []string {
	// Only set raw mode if stdin is a terminal
	var oldState *term.State
	if f, ok := ui.stdin.(*os.File); ok {
		fd := int(f.Fd())
		var err error
		oldState, err = ui.term.makeRaw(fd)
		if err != nil {
			ui.writeError("Failed to set terminal to raw mode: %v", err)
			return nil
		}
		defer func() {
			if err := ui.term.restore(fd, oldState); err != nil {
				ui.writeError("failed to restore terminal state: %v", err)
			}
		}()
	}

	reader := bufio.NewReader(ui.stdin)
	selected := 0
	input := ""

	for {
		ui.write("\033[H\033[2J\033[H") // Clear screen
		ui.writeln("Select a command (incremental search: type to filter, ctrl+n: down, ctrl+p: up, enter: execute, ctrl+c: quit)")
		ui.writeln("\rSearch: %s\n", input)

		// Filtering
		filtered := []CommandInfo{}
		for _, cmd := range commands {
			if strings.Contains(cmd.Command, input) {
				filtered = append(filtered, cmd)
			}
		}
		if input == "" {
			ui.writeln("(Type to filter commands...)")
		} else {
			if len(filtered) == 0 {
				ui.writeln("  (No matching command)")
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
					ui.writeln("\r> %s%s  %s", cmd.Command, padding, desc)
				} else {
					ui.writeln("\r  %s%s  %s", cmd.Command, padding, desc)
				}
			}
		}
		ui.write("\n\r") // Ensure next output starts at left edge

		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			continue
		}
		if b == 3 { // Ctrl+C in raw mode
			if oldState != nil {
				if f, ok := ui.stdin.(*os.File); ok {
					if err := ui.term.restore(int(f.Fd()), oldState); err != nil {
						ui.writeError("failed to restore terminal state: %v", err)
					}
				}
			}
			ui.writeln("\nExiting...")
			os.Exit(0)
		} else if b == 13 { // Enter
			if input == "" {
				continue
			}
			if len(filtered) > 0 {
				ui.writeln("\nExecute: %s", filtered[selected].Command)
				if oldState != nil {
					if f, ok := ui.stdin.(*os.File); ok {
						if err := ui.term.restore(int(f.Fd()), oldState); err != nil {
							ui.writeError("failed to restore terminal state: %v", err)
						}
					}
				}
				// Placeholder detection
				cmdTemplate := filtered[selected].Command
				placeholders := extractPlaceholders(cmdTemplate)
				inputs := make(map[string]string)
				readerStdin := bufio.NewReader(ui.stdin)
				for _, ph := range placeholders {
					ui.write("\n\r") // Newline + carriage return
					ui.write("Enter value for %s: ", ph)
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
