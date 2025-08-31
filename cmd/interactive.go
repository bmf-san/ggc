// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"bufio"
	"fmt"
	"io"
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
	stdin    io.Reader
	stdout   io.Writer
	stderr   io.Writer
	term     terminal
	renderer *Renderer
	state    *UIState
	handler  *KeyHandler
}

// UIState holds the current state of the interactive UI
type UIState struct {
	selected int
	input    string
	filtered []CommandInfo
}

// UpdateFiltered updates the filtered commands based on current input
func (s *UIState) UpdateFiltered() {
	s.filtered = []CommandInfo{}
	for _, cmd := range commands {
		if strings.Contains(cmd.Command, s.input) {
			s.filtered = append(s.filtered, cmd)
		}
	}
	// Reset selection if out of bounds
	if s.selected >= len(s.filtered) {
		s.selected = len(s.filtered) - 1
	}
	if s.selected < 0 {
		s.selected = 0
	}
}

// MoveUp moves selection up
func (s *UIState) MoveUp() {
	if s.selected > 0 {
		s.selected--
	}
}

// MoveDown moves selection down
func (s *UIState) MoveDown() {
	if s.selected < len(s.filtered)-1 {
		s.selected++
	}
}

// AddChar adds a character to the input
func (s *UIState) AddChar(c byte) {
	s.input += string(c)
	s.UpdateFiltered()
}

// RemoveChar removes the last character from input
func (s *UIState) RemoveChar() {
	if len(s.input) > 0 {
		s.input = s.input[:len(s.input)-1]
		s.UpdateFiltered()
	}
}

// GetSelectedCommand returns the currently selected command
func (s *UIState) GetSelectedCommand() *CommandInfo {
	if len(s.filtered) > 0 && s.selected >= 0 && s.selected < len(s.filtered) {
		return &s.filtered[s.selected]
	}
	return nil
}

// HasInput returns true if there is input
func (s *UIState) HasInput() bool {
	return s.input != ""
}

// HasMatches returns true if there are filtered matches
func (s *UIState) HasMatches() bool {
	return len(s.filtered) > 0
}

// Renderer handles all terminal rendering operations
type Renderer struct {
	writer io.Writer
	width  int
	height int
}

// KeyHandler manages keyboard input processing
type KeyHandler struct {
	ui *UI
}

// HandleKey processes a single key input and returns true if should continue
func (h *KeyHandler) HandleKey(b byte, oldState *term.State) (bool, []string) {
	switch b {
	case 3: // Ctrl+C
		h.handleCtrlC(oldState)
		return false, nil
	case 13: // Enter
		return h.handleEnter(oldState)
	case 16: // Ctrl+P (up)
		h.ui.state.MoveUp()
		return true, nil
	case 14: // Ctrl+N (down)
		h.ui.state.MoveDown()
		return true, nil
	case 127, 8: // Backspace
		h.ui.state.RemoveChar()
		return true, nil
	default:
		if b >= 32 && b <= 126 { // Printable ASCII
			h.ui.state.AddChar(b)
		}
		return true, nil
	}
}

// handleCtrlC handles Ctrl+C key press
func (h *KeyHandler) handleCtrlC(oldState *term.State) {
	if oldState != nil {
		if f, ok := h.ui.stdin.(*os.File); ok {
			if err := h.ui.term.restore(int(f.Fd()), oldState); err != nil {
				h.ui.writeError("failed to restore terminal state: %v", err)
			}
		}
	}
	h.ui.writeln("\nExiting...")
	os.Exit(0)
}

// handleEnter handles Enter key press
func (h *KeyHandler) handleEnter(oldState *term.State) (bool, []string) {
	if !h.ui.state.HasInput() {
		return true, nil
	}

	selectedCmd := h.ui.state.GetSelectedCommand()
	if selectedCmd == nil {
		return true, nil
	}

	// Restore terminal state BEFORE showing Execute message
	if oldState != nil {
		if f, ok := h.ui.stdin.(*os.File); ok {
			if err := h.ui.term.restore(int(f.Fd()), oldState); err != nil {
				h.ui.writeError("failed to restore terminal state: %v", err)
			}
		}
	}

	h.ui.write("Execute: %s\n", selectedCmd.Command)

	// Handle placeholders
	return false, h.processCommand(selectedCmd.Command)
}

// processCommand processes the command with placeholder replacement
func (h *KeyHandler) processCommand(cmdTemplate string) []string {
	placeholders := extractPlaceholders(cmdTemplate)
	inputs := make(map[string]string)
	readerStdin := bufio.NewReader(h.ui.stdin)

	for _, ph := range placeholders {
		h.ui.write("\n") // Newline
		h.ui.write("Enter value for %s: ", ph)
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
	renderer := &Renderer{writer: os.Stdout}
	renderer.updateSize()

	state := &UIState{
		selected: 0,
		input:    "",
		filtered: []CommandInfo{},
	}

	ui := &UI{
		stdin:    os.Stdin,
		stdout:   os.Stdout,
		stderr:   os.Stderr,
		term:     &defaultTerminal{},
		renderer: renderer,
		state:    state,
	}

	ui.handler = &KeyHandler{ui: ui}

	return ui
}

// updateSize updates the terminal dimensions
func (r *Renderer) updateSize() {
	if f, ok := r.writer.(*os.File); ok {
		if w, h, err := term.GetSize(int(f.Fd())); err == nil && w > 0 && h > 0 {
			r.width, r.height = w, h
			return
		}
	}
	r.width, r.height = 80, 24 // Default fallback
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
	{"branch rename <old> <new>", "Rename a branch"},
	{"branch move <branch> <commit>", "Move branch to specified commit"},
	{"branch set-upstream <branch> <upstream>", "Set upstream for a branch"},
	{"branch info <branch>", "Show detailed branch information"},
	{"branch list --verbose", "Show detailed branch listing"},
	{"branch sort [date|name]", "List branches sorted by date or name"},
	{"branch contains <commit>", "Show branches containing a commit"},
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
	// Move to line start, clear line, write content, then CRLF
	_, _ = fmt.Fprint(ui.stdout, "\r\x1b[K")
	_, _ = fmt.Fprintf(ui.stdout, format+"\r\n", a...)
}

// clearScreen clears the entire screen and hides cursor
func clearScreen(w io.Writer) {
	// Clear screen, move cursor to home, hide cursor
	_, _ = fmt.Fprint(w, "\x1b[2J\x1b[H\x1b[?25l")
}

// showCursor shows the terminal cursor
func showCursor(w io.Writer) {
	_, _ = fmt.Fprint(w, "\x1b[?25h")
}

// ellipsis truncates string and adds ellipsis if it exceeds maxLen (ASCII only)
func ellipsis(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 1 {
		return "…"
	}
	return s[:maxLen-1] + "…"
}

// Render displays the command list with proper terminal handling
func (r *Renderer) Render(ui *UI, state *UIState) {
	clearScreen(r.writer)
	// Disable line wrapping during rendering, restore at end
	_, _ = fmt.Fprint(r.writer, "\x1b[?7l")
	defer func() {
		_, _ = fmt.Fprint(r.writer, "\x1b[?7h")
		showCursor(r.writer)
	}()

	// Update terminal size
	r.updateSize()

	// Header
	r.writeln(ui, "Select a command (incremental search: type to filter, Ctrl+N: down, Ctrl+P: up, Enter: execute, Ctrl+C: quit)")
	r.writeln(ui, "Search: %s", state.input)
	r.writeln(ui, "") // Empty line

	if state.input == "" {
		r.writeln(ui, "(Type to filter commands...)")
		return
	}

	if len(state.filtered) == 0 {
		r.writeln(ui, "  (No matching command)")
		return
	}

	// Clamp selection index to valid range
	if state.selected >= len(state.filtered) {
		state.selected = len(state.filtered) - 1
	}
	if state.selected < 0 {
		state.selected = 0
	}

	// Calculate maximum command length for consistent alignment
	maxCmdLen := r.calculateMaxCommandLength(state.filtered)

	for i, cmd := range state.filtered {
		desc := cmd.Description
		if desc == "" {
			desc = "No description"
		}

		// Calculate padding for consistent command alignment
		paddingLen := maxCmdLen - len(cmd.Command)
		if paddingLen < 0 {
			paddingLen = 0
		}
		padding := strings.Repeat(" ", paddingLen)

		// Calculate line layout for truncation
		prefix := "  "
		if i == state.selected {
			prefix = "> "
		}

		// Calculate available width for description
		usedWidth := len(prefix) + len(cmd.Command) + len(padding) + 2 // 2 spaces separator
		availableDescWidth := r.width - usedWidth
		if availableDescWidth < 10 {
			availableDescWidth = 10
		}

		// Truncate description if needed
		trimmedDesc := ellipsis(desc, availableDescWidth)

		if i == state.selected {
			r.writeln(ui, "> %s%s  %s", cmd.Command, padding, trimmedDesc)
		} else {
			r.writeln(ui, "  %s%s  %s", cmd.Command, padding, trimmedDesc)
		}
	}
}

// writeln writes a message with newline through the renderer
func (r *Renderer) writeln(_ *UI, format string, a ...interface{}) {
	// Move to line start, clear line, write content, then CRLF
	_, _ = fmt.Fprint(r.writer, "\r\x1b[K")
	_, _ = fmt.Fprintf(r.writer, format+"\r\n", a...)
}

// calculateMaxCommandLength calculates the maximum command length for alignment
func (r *Renderer) calculateMaxCommandLength(filtered []CommandInfo) int {
	maxCmdLen := 0
	for _, cmd := range filtered {
		if len(cmd.Command) > maxCmdLen {
			maxCmdLen = len(cmd.Command)
		}
	}
	return maxCmdLen
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

	for {
		// Update filtered commands
		ui.state.UpdateFiltered()

		// Render the UI
		ui.renderer.Render(ui, ui.state)

		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			continue
		}

		// Handle key input
		shouldContinue, result := ui.handler.HandleKey(b, oldState)
		if !shouldContinue {
			return result
		}
	}
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
