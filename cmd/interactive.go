// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/term"
)

// GitStatus represents the current Git repository status
type GitStatus struct {
	Branch     string
	Modified   int
	Staged     int
	Ahead      int
	Behind     int
	HasChanges bool
}

// ANSIColors defines color codes for terminal output
type ANSIColors struct {
	// Basic colors (0-7)
	Black   string
	Red     string
	Green   string
	Yellow  string
	Blue    string
	Magenta string
	Cyan    string
	White   string

	// Bright colors (8-15)
	BrightBlack   string // Gray
	BrightRed     string
	BrightGreen   string
	BrightYellow  string
	BrightBlue    string
	BrightMagenta string
	BrightCyan    string
	BrightWhite   string

	// Text attributes
	Bold      string
	Underline string
	Reverse   string
	Reset     string
}

// NewANSIColors creates a new ANSIColors instance
func NewANSIColors() *ANSIColors {
	return &ANSIColors{
		// Basic colors
		Black:   "\033[30m",
		Red:     "\033[31m",
		Green:   "\033[32m",
		Yellow:  "\033[33m",
		Blue:    "\033[34m",
		Magenta: "\033[35m",
		Cyan:    "\033[36m",
		White:   "\033[37m",

		// Bright colors
		BrightBlack:   "\033[90m",
		BrightRed:     "\033[91m",
		BrightGreen:   "\033[92m",
		BrightYellow:  "\033[93m",
		BrightBlue:    "\033[94m",
		BrightMagenta: "\033[95m",
		BrightCyan:    "\033[96m",
		BrightWhite:   "\033[97m",

		// Text attributes
		Bold:      "\033[1m",
		Underline: "\033[4m",
		Reverse:   "\033[7m",
		Reset:     "\033[0m",
	}
}

// getGitStatus retrieves the current Git repository status
func getGitStatus() *GitStatus {
	status := &GitStatus{}

	// Get current branch name
	if branch := getGitBranch(); branch != "" {
		status.Branch = branch
	} else {
		return nil // Not in a git repository
	}

	// Get working directory status
	modified, staged := getGitWorkingStatus()
	status.Modified = modified
	status.Staged = staged
	status.HasChanges = modified > 0 || staged > 0

	// Get remote tracking status
	ahead, behind := getGitRemoteStatus()
	status.Ahead = ahead
	status.Behind = behind

	return status
}

// getGitBranch gets the current branch name
func getGitBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// getGitWorkingStatus gets the number of modified and staged files
func getGitWorkingStatus() (modified, staged int) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return 0, 0
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if len(line) < 2 {
			continue
		}
		
		// First character: staged status
		// Second character: working tree status
		if line[0] != ' ' && line[0] != '?' {
			staged++
		}
		if line[1] != ' ' && line[1] != '?' {
			modified++
		}
	}
	return modified, staged
}

// getGitRemoteStatus gets ahead/behind count compared to remote
func getGitRemoteStatus() (ahead, behind int) {
	cmd := exec.Command("git", "rev-list", "--count", "--left-right", "HEAD...@{upstream}")
	output, err := cmd.Output()
	if err != nil {
		return 0, 0 // No upstream or other error
	}

	parts := strings.Fields(strings.TrimSpace(string(output)))
	if len(parts) != 2 {
		return 0, 0
	}

	ahead, _ = strconv.Atoi(parts[0])
	behind, _ = strconv.Atoi(parts[1])
	return ahead, behind
}

// CommandInfo contains the name and description of the command
type CommandInfo struct {
	Command     string
	Description string
}

// UI represents the interface for terminal UI operations
type UI struct {
	stdin     io.Reader
	stdout    io.Writer
	stderr    io.Writer
	term      terminal
	renderer  *Renderer
	state     *UIState
	handler   *KeyHandler
	colors    *ANSIColors
	gitStatus *GitStatus
}

// UIState holds the current state of the interactive UI
type UIState struct {
	selected  int
	input     string
	cursorPos int // Cursor position in input string
	filtered  []CommandInfo
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

// AddChar adds a character to the input at cursor position
func (s *UIState) AddChar(c byte) {
	if s.cursorPos <= len(s.input) {
		s.input = s.input[:s.cursorPos] + string(c) + s.input[s.cursorPos:]
		s.cursorPos++
		s.UpdateFiltered()
	}
}

// RemoveChar removes character before cursor (backspace)
func (s *UIState) RemoveChar() {
	if s.cursorPos > 0 && len(s.input) > 0 {
		s.input = s.input[:s.cursorPos-1] + s.input[s.cursorPos:]
		s.cursorPos--
		s.UpdateFiltered()
	}
}

// ClearInput clears all input
func (s *UIState) ClearInput() {
	s.input = ""
	s.cursorPos = 0
	s.UpdateFiltered()
}

// DeleteWord deletes word before cursor (Ctrl+W)
func (s *UIState) DeleteWord() {
	if s.cursorPos == 0 {
		return
	}

	// Find start of current word (skip trailing spaces first)
	pos := s.cursorPos - 1
	for pos >= 0 && s.input[pos] == ' ' {
		pos--
	}

	// Find start of word
	for pos >= 0 && s.input[pos] != ' ' {
		pos--
	}
	pos++ // Move to first character of word

	// Delete from word start to cursor
	s.input = s.input[:pos] + s.input[s.cursorPos:]
	s.cursorPos = pos
	s.UpdateFiltered()
}

// DeleteToEnd deletes from cursor to end of line (Ctrl+K)
func (s *UIState) DeleteToEnd() {
	if s.cursorPos < len(s.input) {
		s.input = s.input[:s.cursorPos]
		s.UpdateFiltered()
	}
}

// MoveToBeginning moves cursor to beginning of line (Ctrl+A)
func (s *UIState) MoveToBeginning() {
	s.cursorPos = 0
}

// MoveToEnd moves cursor to end of line (Ctrl+E)
func (s *UIState) MoveToEnd() {
	s.cursorPos = len(s.input)
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
	colors *ANSIColors
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
	case 21: // Ctrl+U (clear line)
		h.ui.state.ClearInput()
		return true, nil
	case 23: // Ctrl+W (delete word)
		h.ui.state.DeleteWord()
		return true, nil
	case 11: // Ctrl+K (delete to end)
		h.ui.state.DeleteToEnd()
		return true, nil
	case 1: // Ctrl+A (beginning of line)
		h.ui.state.MoveToBeginning()
		return true, nil
	case 5: // Ctrl+E (end of line)
		h.ui.state.MoveToEnd()
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

	// Clear screen and show execution message
	clearScreen(h.ui.stdout)
	executeMsg := fmt.Sprintf("%s🚀 %sExecuting:%s %s%s%s\n\n",
		h.ui.colors.BrightGreen,
		h.ui.colors.BrightWhite+h.ui.colors.Bold,
		h.ui.colors.Reset,
		h.ui.colors.BrightCyan+h.ui.colors.Bold,
		selectedCmd.Command,
		h.ui.colors.Reset)
	h.ui.writeColor(executeMsg)

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
	colors := NewANSIColors()

	renderer := &Renderer{
		writer: os.Stdout,
		colors: colors,
	}
	renderer.updateSize()

	state := &UIState{
		selected: 0,
		input:    "",
		filtered: []CommandInfo{},
	}

	ui := &UI{
		stdin:     os.Stdin,
		stdout:    os.Stdout,
		stderr:    os.Stderr,
		term:      &defaultTerminal{},
		renderer:  renderer,
		state:     state,
		colors:    colors,
		gitStatus: getGitStatus(),
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

// writeColor writes a colored message to stdout
func (ui *UI) writeColor(text string) {
	_, _ = fmt.Fprint(ui.stdout, text)
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

	// Modern header with title and subtitle
	title := fmt.Sprintf("%s%s🚀 ggc Interactive Mode%s",
		r.colors.BrightCyan+r.colors.Bold,
		r.colors.Reset,
		r.colors.Reset)
	r.writeColorln(ui, title)

	// Git status information
	if ui.gitStatus != nil {
		r.renderGitStatus(ui, ui.gitStatus)
	}

	subtitle := fmt.Sprintf("%sType to search • %sCtrl+n/p%s navigate • %sCtrl+a/e%s move • %sEnter%s execute • %sCtrl+c%s quit%s",
		r.colors.BrightBlack,
		r.colors.BrightGreen+r.colors.Bold,
		r.colors.BrightBlack,
		r.colors.BrightCyan+r.colors.Bold,
		r.colors.BrightBlack,
		r.colors.BrightGreen+r.colors.Bold,
		r.colors.BrightBlack,
		r.colors.BrightGreen+r.colors.Bold,
		r.colors.BrightBlack,
		r.colors.Reset)
	r.writeColorln(ui, subtitle)
	r.writeEmptyLine()

	// Enhanced search prompt with cursor at correct position
	var inputWithCursor string
	if len(state.input) == 0 {
		inputWithCursor = fmt.Sprintf("%s█%s", r.colors.BrightWhite+r.colors.Bold, r.colors.Reset)
	} else {
		beforeCursor := state.input[:state.cursorPos]
		afterCursor := state.input[state.cursorPos:]
		cursor := "│"
		if state.cursorPos >= len(state.input) {
			cursor = "█"
		}
		inputWithCursor = fmt.Sprintf("%s%s%s%s%s%s%s",
			r.colors.BrightYellow,
			beforeCursor,
			r.colors.BrightWhite+r.colors.Bold,
			cursor,
			r.colors.Reset+r.colors.BrightYellow,
			afterCursor,
			r.colors.Reset)
	}

	searchPrompt := fmt.Sprintf("%s┌─ %sSearch:%s %s",
		r.colors.BrightBlue,
		r.colors.BrightGreen+r.colors.Bold,
		r.colors.Reset,
		inputWithCursor)
	r.writeColorln(ui, searchPrompt)

	// Results separator
	if state.input != "" {
		separator := fmt.Sprintf("%s└─ %sResults:%s",
			r.colors.BrightBlue,
			r.colors.BrightMagenta+r.colors.Bold,
			r.colors.Reset)
		r.writeColorln(ui, separator)
	}
	r.writeEmptyLine()

	if state.input == "" {
		// Empty state - simple message
		r.writeColorln(ui, fmt.Sprintf("%s💭 %sStart typing to search commands...%s",
			r.colors.BrightBlue, r.colors.BrightBlack, r.colors.Reset))
		return
	}

	if len(state.filtered) == 0 {
		// No matches found - helpful message
		r.writeColorln(ui, fmt.Sprintf("%s🔍 %sNo commands found for '%s%s%s'%s",
			r.colors.BrightYellow,
			r.colors.BrightWhite,
			r.colors.BrightYellow+r.colors.Bold,
			state.input,
			r.colors.Reset+r.colors.BrightWhite,
			r.colors.Reset))
		r.writeEmptyLine()

		// Available keybinds
		keybinds := []struct{ key, desc string }{
			{"Ctrl+u", "Clear all input"},
			{"Ctrl+w", "Delete word"},
			{"Ctrl+k", "Delete to end"},
			{"Ctrl+a", "Move to beginning"},
			{"Ctrl+e", "Move to end"},
			{"Backspace", "Delete character"},
		}

		r.writeColorln(ui, fmt.Sprintf("%s⌨️  %sAvailable keybinds:%s",
			r.colors.BrightBlue, r.colors.BrightWhite+r.colors.Bold, r.colors.Reset))

		for _, kb := range keybinds {
			r.writeColorln(ui, fmt.Sprintf("   %s%s%s  %s%s%s",
				r.colors.BrightGreen+r.colors.Bold,
				kb.key,
				r.colors.Reset,
				r.colors.BrightBlack,
				kb.desc,
				r.colors.Reset))
		}
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

		// Calculate available width for description
		usedWidth := 4 + len(cmd.Command) + len(padding) + 3 // prefix + command + padding + separator
		availableDescWidth := r.width - usedWidth
		if availableDescWidth < 10 {
			availableDescWidth = 10
		}

		// Truncate description if needed
		trimmedDesc := ellipsis(desc, availableDescWidth)

		if i == state.selected {
			// Selected item with modern highlighting
			selectedLine := fmt.Sprintf("%s▶ %s%s%s%s %s│%s %s%s%s",
				r.colors.BrightCyan+r.colors.Bold,
				r.colors.BrightWhite+r.colors.Bold+r.colors.Reverse,
				" "+cmd.Command+" ",
				r.colors.Reset,
				padding,
				r.colors.BrightBlue,
				r.colors.Reset,
				r.colors.BrightWhite,
				trimmedDesc,
				r.colors.Reset)
			r.writeColorln(ui, selectedLine)
		} else {
			// Regular item with improved styling
			regularLine := fmt.Sprintf("  %s%s%s%s %s│%s %s%s%s",
				r.colors.BrightGreen+r.colors.Bold,
				cmd.Command,
				r.colors.Reset,
				padding,
				r.colors.BrightBlack,
				r.colors.Reset,
				r.colors.BrightBlack,
				trimmedDesc,
				r.colors.Reset)
			r.writeColorln(ui, regularLine)
		}
	}

	// Footer with navigation hints
	r.writeEmptyLine()
	footer := fmt.Sprintf("%s%sCtrl+n/p%s Navigate  %sCtrl+a/e%s Move  %sCtrl+u/w/k%s Edit  %sEnter%s Execute  %sCtrl+c%s Exit%s",
		r.colors.BrightBlack,
		r.colors.BrightGreen+r.colors.Bold,
		r.colors.BrightBlack,
		r.colors.BrightCyan+r.colors.Bold,
		r.colors.BrightBlack,
		r.colors.BrightYellow+r.colors.Bold,
		r.colors.BrightBlack,
		r.colors.BrightGreen+r.colors.Bold,
		r.colors.BrightBlack,
		r.colors.BrightGreen+r.colors.Bold,
		r.colors.BrightBlack,
		r.colors.Reset)
	r.writeColorln(ui, footer)
}

// writeColorln writes a colored line to the terminal
func (r *Renderer) writeColorln(_ *UI, text string) {
	// Move to line start, clear line, write content, then CRLF
	_, _ = fmt.Fprint(r.writer, "\r\x1b[K")
	_, _ = fmt.Fprint(r.writer, text+"\r\n")
}

// renderGitStatus renders the Git repository status information
func (r *Renderer) renderGitStatus(ui *UI, status *GitStatus) {
	var parts []string

	// Branch name
	branchPart := fmt.Sprintf("%s📍 %s%s%s",
		r.colors.BrightBlue,
		r.colors.BrightWhite+r.colors.Bold,
		status.Branch,
		r.colors.Reset)
	parts = append(parts, branchPart)

	// Working directory status
	if status.HasChanges {
		var statusParts []string
		if status.Modified > 0 {
			statusParts = append(statusParts, fmt.Sprintf("%d modified", status.Modified))
		}
		if status.Staged > 0 {
			statusParts = append(statusParts, fmt.Sprintf("%d staged", status.Staged))
		}
		
		workingPart := fmt.Sprintf("%s📝 %s%s%s",
			r.colors.BrightYellow,
			r.colors.BrightWhite+r.colors.Bold,
			strings.Join(statusParts, ", "),
			r.colors.Reset)
		parts = append(parts, workingPart)
	}

	// Remote tracking status
	if status.Ahead > 0 || status.Behind > 0 {
		var remoteParts []string
		if status.Ahead > 0 {
			remoteParts = append(remoteParts, fmt.Sprintf("↑%d", status.Ahead))
		}
		if status.Behind > 0 {
			remoteParts = append(remoteParts, fmt.Sprintf("↓%d", status.Behind))
		}
		
		remotePart := fmt.Sprintf("%s%s%s",
			r.colors.BrightMagenta+r.colors.Bold,
			strings.Join(remoteParts, " "),
			r.colors.Reset)
		parts = append(parts, remotePart)
	}

	// Render the status line
	statusLine := strings.Join(parts, "  ")
	r.writeColorln(ui, statusLine)
}

// writeEmptyLine writes an empty line
func (r *Renderer) writeEmptyLine() {
	_, _ = fmt.Fprint(r.writer, "\r\x1b[K\r\n")
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
