// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"unicode"
	"unicode/utf8"

	"golang.org/x/term"
	"golang.org/x/text/width"

	commandregistry "github.com/bmf-san/ggc/v7/cmd/command"
	"github.com/bmf-san/ggc/v7/config"
	"github.com/bmf-san/ggc/v7/git"
	"github.com/bmf-san/ggc/v7/internal/termio"
)

// initialInputCapacity defines the initial capacity for the input rune buffer
// used by the real-time editor. It helps minimize reallocations during typing
// while keeping memory usage modest.
const initialInputCapacity = 64

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
func getGitStatus(gitClient git.StatusInfoReader) *GitStatus {
	status := &GitStatus{}

	// Get current branch name
	if branch := getGitBranch(gitClient); branch != "" {
		status.Branch = branch
	} else {
		return nil // Not in a git repository
	}

	// Get working directory status
	modified, staged := getGitWorkingStatus(gitClient)
	status.Modified = modified
	status.Staged = staged
	status.HasChanges = modified > 0 || staged > 0

	// Get remote tracking status
	ahead, behind := getGitRemoteStatus(gitClient)
	status.Ahead = ahead
	status.Behind = behind

	return status
}

// getGitBranch gets the current branch name
func getGitBranch(gitClient git.StatusInfoReader) string {
	branch, err := gitClient.GetCurrentBranch()
	if err != nil {
		return ""
	}
	return branch
}

// getGitWorkingStatus gets the number of modified and staged files
func getGitWorkingStatus(gitClient git.StatusInfoReader) (modified, staged int) {
	output, err := gitClient.StatusShortWithColor()
	if err != nil {
		return 0, 0
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
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
func getGitRemoteStatus(gitClient git.StatusInfoReader) (ahead, behind int) {
	output, err := gitClient.GetAheadBehindCount("HEAD", "@{upstream}")
	if err != nil {
		return 0, 0 // No upstream or other error
	}

	parts := strings.Fields(strings.TrimSpace(output))
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

func buildInteractiveCommands() []CommandInfo {
	var list []CommandInfo
	allCommands := commandregistry.All()
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

// UI represents the interface for terminal UI operations
type UI struct {
	stdin           io.Reader
	stdout          io.Writer
	stderr          io.Writer
	term            termio.Terminal
	renderer        *Renderer
	state           *UIState
	handler         *KeyHandler
	colors          *ANSIColors
	gitStatus       *GitStatus
	gitClient       git.StatusInfoReader
	reader          *bufio.Reader
	contextMgr      *ContextManager
	profile         Profile
	workflow        *Workflow
	workflowEx      *WorkflowExecutor
	softCancelFlash atomic.Bool
}

// UIState holds the current state of the interactive UI
type UIState struct {
	selected        int
	input           string
	cursorPos       int // Cursor position in input string
	filtered        []CommandInfo
	context         Context   // Current UI context (input/results/search/global)
	contextStack    []Context // Context stack for nested states
	onContextChange func(Context, Context)
	showWorkflow    bool // Whether to show the workflow builder
}

// UpdateFiltered updates the filtered commands based on current input using fuzzy matching
func (s *UIState) UpdateFiltered() {
	s.filtered = []CommandInfo{}
	input := strings.ToLower(s.input)
	for _, cmd := range commands {
		cmdLower := strings.ToLower(cmd.Command)
		if fuzzyMatch(cmdLower, input) {
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

// Context Management Methods

// EnterContext pushes the current context onto the stack and switches to the new context
func (s *UIState) EnterContext(newContext Context) {
	if s.context == newContext {
		return
	}
	old := s.context
	s.contextStack = append(s.contextStack, s.context)
	s.context = newContext
	s.notifyContextChange(old, newContext)
}

// ExitContext pops the previous context from the stack
func (s *UIState) ExitContext() {
	if len(s.contextStack) > 0 {
		old := s.context
		s.context = s.contextStack[len(s.contextStack)-1]
		s.contextStack = s.contextStack[:len(s.contextStack)-1]
		s.notifyContextChange(old, s.context)
	} else if s.context != ContextGlobal {
		old := s.context
		s.context = ContextGlobal
		s.notifyContextChange(old, s.context)
	}
}

// GetCurrentContext returns the current UI context
func (s *UIState) GetCurrentContext() Context {
	return s.context
}

// SetContext directly sets the context (use with caution)
func (s *UIState) SetContext(ctx Context) {
	if s.context == ctx {
		return
	}
	old := s.context
	s.context = ctx
	s.notifyContextChange(old, ctx)
}

// notifyContextChange triggers the callback when the active context changes
func (s *UIState) notifyContextChange(oldCtx, newCtx Context) {
	if s.onContextChange != nil && oldCtx != newCtx {
		s.onContextChange(oldCtx, newCtx)
	}
}

// IsInInputMode returns true if currently in input context
func (s *UIState) IsInInputMode() bool {
	return s.context == ContextInput
}

// IsInResultsMode returns true if currently in results context
func (s *UIState) IsInResultsMode() bool {
	return s.context == ContextResults
}

// IsInSearchMode returns true if currently in search context
func (s *UIState) IsInSearchMode() bool {
	return s.context == ContextSearch
}

// fuzzyMatch performs fuzzy matching between text and pattern
// Returns true if all characters in pattern appear in text in order (but not necessarily consecutive)
func fuzzyMatch(text, pattern string) bool {
	if pattern == "" {
		return true
	}

	textRunes := []rune(text)
	patternRunes := []rune(pattern)

	textIdx := 0
	patternIdx := 0

	for textIdx < len(textRunes) && patternIdx < len(patternRunes) {
		if textRunes[textIdx] == patternRunes[patternIdx] {
			patternIdx++
		}
		textIdx++
	}

	return patternIdx == len(patternRunes)
}

// MoveUp moves selection up
func (s *UIState) MoveUp() {
	// Switch to results context when navigating
	if s.context != ContextResults && s.context != ContextSearch {
		s.SetContext(ContextResults)
	}

	if s.selected > 0 {
		s.selected--
	}
}

// MoveDown moves selection down
func (s *UIState) MoveDown() {
	// Switch to results context when navigating
	if s.context != ContextResults && s.context != ContextSearch {
		s.SetContext(ContextResults)
	}

	if s.selected < len(s.filtered)-1 {
		s.selected++
	}
}

// AddRune adds a UTF-8 rune to the input at cursor position
func (s *UIState) AddRune(r rune) {
	// Switch to input context when user starts typing
	if s.context != ContextInput {
		s.SetContext(ContextInput)
	}

	// Convert current input to runes for proper cursor positioning
	inputRunes := []rune(s.input)
	if s.cursorPos <= len(inputRunes) {
		// Efficiently insert the rune at the cursor position
		newRunes := make([]rune, len(inputRunes)+1)
		copy(newRunes, inputRunes[:s.cursorPos])
		newRunes[s.cursorPos] = r
		copy(newRunes[s.cursorPos+1:], inputRunes[s.cursorPos:])

		s.input = string(newRunes)
		s.cursorPos++
		s.UpdateFiltered()

		// Switch to search context when actively filtering
		if s.input != "" && s.context != ContextSearch {
			s.SetContext(ContextSearch)
		}
	}
}

// RemoveChar removes character before cursor (backspace)
func (s *UIState) RemoveChar() {
	if s.cursorPos > 0 && s.input != "" {
		// Convert to runes for proper UTF-8 handling
		inputRunes := []rune(s.input)
		if s.cursorPos <= len(inputRunes) {
			inputRunes = append(inputRunes[:s.cursorPos-1], inputRunes[s.cursorPos:]...)

			s.input = string(inputRunes)
			s.cursorPos--
			s.UpdateFiltered()
		}
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

	// Convert to runes for proper UTF-8 handling
	inputRunes := []rune(s.input)

	// Find start of current word (skip trailing spaces first)
	pos := s.cursorPos - 1
	for pos >= 0 && inputRunes[pos] == ' ' {
		pos--
	}

	// Find start of word
	for pos >= 0 && inputRunes[pos] != ' ' {
		pos--
	}
	pos++ // Move to first character of word

	// Delete from word start to cursor
	inputRunes = append(inputRunes[:pos], inputRunes[s.cursorPos:]...)
	s.input = string(inputRunes)
	s.cursorPos = pos
	s.UpdateFiltered()
}

// DeleteToEnd deletes from cursor to end of line (Ctrl+K)
func (s *UIState) DeleteToEnd() {
	if s.cursorPos < utf8.RuneCountInString(s.input) {
		inputRunes := []rune(s.input)
		s.input = string(inputRunes[:s.cursorPos])
		s.UpdateFiltered()
	}
}

// MoveToBeginning moves cursor to beginning of line (Ctrl+A)
func (s *UIState) MoveToBeginning() {
	s.cursorPos = 0
}

// MoveToEnd moves cursor to end of line (Ctrl+E)
func (s *UIState) MoveToEnd() {
	s.cursorPos = utf8.RuneCountInString(s.input)
}

// MoveLeft moves cursor one rune left
func (s *UIState) MoveLeft() {
	if s.cursorPos > 0 {
		s.cursorPos--
	}
}

// MoveRight moves cursor one rune right
func (s *UIState) MoveRight() {
	if s.cursorPos < utf8.RuneCountInString(s.input) {
		s.cursorPos++
	}
}

// MoveWordLeft moves cursor to the beginning of the previous word
func (s *UIState) MoveWordLeft() {
	if s.cursorPos == 0 {
		return
	}
	runes := []rune(s.input)
	pos := s.cursorPos - 1
	// Skip any spaces to the left
	for pos >= 0 && unicode.IsSpace(runes[pos]) {
		pos--
	}
	// Skip the word characters to the left
	for pos >= 0 && !unicode.IsSpace(runes[pos]) {
		pos--
	}
	s.cursorPos = pos + 1
}

// MoveWordRight moves cursor to the start of the next word
func (s *UIState) MoveWordRight() {
	runes := []rune(s.input)
	n := len(runes)
	pos := s.cursorPos
	if pos >= n {
		return
	}
	// Skip current word characters
	for pos < n && !unicode.IsSpace(runes[pos]) {
		pos++
	}
	// Skip following spaces
	for pos < n && unicode.IsSpace(runes[pos]) {
		pos++
	}
	s.cursorPos = pos
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

type keybindHelpEntry struct {
	key  string
	desc string
}

// KeyHandler manages keyboard input processing
type KeyHandler struct {
	ui            *UI
	contextualMap *ContextualKeyBindingMap
}

// GetCurrentKeyMap returns the appropriate keybinding map for the current context
func (h *KeyHandler) GetCurrentKeyMap() *KeyBindingMap {
	if h == nil {
		return DefaultKeyBindingMap()
	}
	if h.contextualMap != nil && h.ui != nil && h.ui.state != nil {
		currentContext := h.ui.state.GetCurrentContext()
		if contextMap, exists := h.contextualMap.GetContext(currentContext); exists && contextMap != nil {
			return contextMap
		}
		if contextMap, exists := h.contextualMap.GetContext(ContextGlobal); exists && contextMap != nil {
			return contextMap
		}
	}

	return DefaultKeyBindingMap()
}

// HandleKey processes UTF-8 rune input and returns true if should continue
// This method handles both single-byte (ASCII/control) and multibyte characters
func (h *KeyHandler) HandleKey(r rune, _ bool, oldState *term.State, reader *bufio.Reader) (bool, []string) {
	// Set the reader for consistent access during escape sequence handling
	h.ui.reader = reader
	// Handle workflow-specific keys first (Tab, etc.)
	if handled := h.handleWorkflowKeys(r); handled {
		return true, nil
	}

	// Handle control characters (ASCII range)
	if r < 128 && unicode.IsControl(r) { // ASCII control characters
		if handled, shouldContinue, result := h.handleControlChar(byte(r), oldState, reader); handled {
			return shouldContinue, result
		}
	}

	// Handle printable characters (both ASCII and multibyte)
	if unicode.IsPrint(r) {
		// Don't accept text input in workflow view
		if !h.ui.state.showWorkflow {
			h.ui.state.AddRune(r)
		}
	}
	return true, nil
}

// handleWorkflowKeys processes workflow-related key bindings and returns (handled, result)
func (h *KeyHandler) handleWorkflowKeys(r rune) bool {
	km := h.GetCurrentKeyMap()

	// Create KeyStroke for this character
	keyStroke := NewCharKeyStroke(r)

	// Check add to workflow
	if km.MatchesKeyStroke("add_to_workflow", keyStroke) {
		// Add current selection to workflow
		if !h.ui.state.showWorkflow && h.ui.state.HasInput() {
			selectedCmd := h.ui.state.GetSelectedCommand()
			if selectedCmd != nil {
				h.addCommandToWorkflow(selectedCmd.Command)
				// Clear input after adding to workflow
				h.ui.state.ClearInput()
				return true
			}
		}
		return true
	}

	// Check clear workflow
	if km.MatchesKeyStroke("clear_workflow", keyStroke) {
		// Clear workflow
		if h.ui.state.showWorkflow {
			h.clearWorkflow()
			return true
		}
		return true
	}

	return false
}

// handleControlChar processes control characters and returns (handled, shouldContinue, result)
// Enhanced to support KeyStroke matching while maintaining backward compatibility
//
//nolint:revive // Control character handling inherently requires many cases
func (h *KeyHandler) handleControlChar(b byte, oldState *term.State, reader *bufio.Reader) (bool, bool, []string) {
	// Get the appropriate keybinding map for current context
	km := h.GetCurrentKeyMap()

	// Create KeyStroke for this control character
	if b >= 1 && b <= 26 {
		// Control character: convert back to letter
		ctrlStroke := NewCtrlKeyStroke(rune('a' + b - 1))

		// Check each action using new KeyStroke matching
		if km.MatchesKeyStroke("move_up", ctrlStroke) {
			// No navigation in workflow view - just normal command navigation
			if !h.ui.state.showWorkflow {
				h.ui.state.MoveUp()
			}
			return true, true, nil
		}
		if km.MatchesKeyStroke("move_down", ctrlStroke) {
			// No navigation in workflow view - just normal command navigation
			if !h.ui.state.showWorkflow {
				h.ui.state.MoveDown()
			}
			return true, true, nil
		}
		if km.MatchesKeyStroke("clear_line", ctrlStroke) {
			h.ui.state.ClearInput()
			return true, true, nil
		}
		if km.MatchesKeyStroke("delete_word", ctrlStroke) {
			h.ui.state.DeleteWord()
			return true, true, nil
		}
		if km.MatchesKeyStroke("delete_to_end", ctrlStroke) {
			h.ui.state.DeleteToEnd()
			return true, true, nil
		}
		if km.MatchesKeyStroke("move_to_beginning", ctrlStroke) {
			h.ui.state.MoveToBeginning()
			return true, true, nil
		}
		if km.MatchesKeyStroke("move_to_end", ctrlStroke) {
			h.ui.state.MoveToEnd()
			return true, true, nil
		}

		// Check for workflow toggle
		if km.MatchesKeyStroke("toggle_workflow_view", ctrlStroke) && h.ui.state.input == "" {
			h.ui.ToggleWorkflowView()
			return true, true, nil
		}
		if km.MatchesKeyStroke("soft_cancel", ctrlStroke) {
			h.handleSoftCancel(oldState)
			return true, true, nil
		}
	}

	// Handle special cases that are not Ctrl+letter
	switch b {
	case 3: // Ctrl+C
		h.handleCtrlC(oldState)
		return true, false, nil
	case 13: // Enter
		shouldContinue, result := h.handleEnter(oldState)
		return true, shouldContinue, result
	case 127, 8: // Backspace
		h.ui.state.RemoveChar()
		return true, true, nil
	case 27: // ESC: arrow keys and Option/Alt modifiers
		if h.shouldHandleEscapeAsSoftCancel() {
			h.handleSoftCancel(oldState)
			return true, true, nil
		}
		h.handleEscapeSequence(reader)
		return true, true, nil
	default:
		return false, true, nil
	}
}

// handleEscapeSequence parses common ESC sequences for arrow and word navigation.
// Supports:
// - Arrow keys: ESC [ C/D (right/left), ESC O C/D (application mode)
// - Ctrl+Arrow: ESC [ 1;5 C/D or ESC [ 5 C/D
// - Alt/Option+Arrow: ESC [ 1;3 C/D, ESC [ 1;9 C/D (varies by terminal)
// - macOS Option word nav: ESC b / ESC f
func (h *KeyHandler) handleEscapeSequence(reader *bufio.Reader) {
	if h.ui == nil {
		return
	}

	// Read next byte after ESC
	var b byte
	var err error

	if reader != nil {
		// Use provided buffered reader (non-raw mode)
		b, err = reader.ReadByte()
	} else {
		// Raw mode: read directly from stdin
		var buf [1]byte
		_, err = h.ui.stdin.Read(buf[:])
		b = buf[0]
	}

	if err != nil {
		return
	}

	switch b {
	case '[':
		h.handleCSISequence(reader)
	case 'O':
		h.handleApplicationCursorMode(reader)
	case 'b':
		h.ui.state.MoveWordLeft()
	case 'f':
		h.ui.state.MoveWordRight()
	case 127, 8:
		// Meta-Backspace (Option+Backspace): delete word left
		h.ui.state.DeleteWord()
	}
}

func (h *KeyHandler) handleSoftCancel(_ *term.State) {
	if h == nil || h.ui == nil {
		return
	}

	if h.ui.resetToSearchMode() {
		h.ui.notifySoftCancel()
	}
}

func (h *KeyHandler) shouldHandleEscapeAsSoftCancel() bool {
	km := h.GetCurrentKeyMap()
	if km == nil || !km.MatchesKeyStroke("soft_cancel", NewEscapeKeyStroke()) {
		return false
	}

	if h.ui == nil {
		return false
	}

	if h.ui.reader != nil && h.ui.reader.Buffered() > 0 {
		return false
	}

	if file, ok := h.ui.stdin.(*os.File); ok {
		if pending, err := termio.PendingInput(file.Fd()); err == nil {
			return pending == 0
		}
	}

	return false
}

// handleCSISequence handles CSI (Control Sequence Introducer) sequences
func (h *KeyHandler) handleCSISequence(reader *bufio.Reader) {
	var params []byte
	for {
		var nb byte
		var err error

		if reader != nil {
			// Use provided buffered reader (non-raw mode)
			nb, err = reader.ReadByte()
		} else {
			// Raw mode: read directly from stdin
			var buf [1]byte
			_, err = h.ui.stdin.Read(buf[:])
			nb = buf[0]
		}

		if err != nil {
			return
		}
		if (nb >= 'A' && nb <= 'Z') || nb == '~' {
			h.processCSIFinalByte(nb, string(params))
			return
		}
		params = append(params, nb)
	}
}

// isWordMotionParam reports whether CSI params include a word-motion modifier
// commonly emitted by terminals (e.g., Ctrl/Alt variants use 5/3/9).
func isWordMotionParam(params string) bool {
	return strings.Contains(params, "5") || strings.Contains(params, "3") || strings.Contains(params, "9")
}

// processCSIFinalByte processes the final byte of a CSI sequence
func (h *KeyHandler) processCSIFinalByte(final byte, params string) {
	isWord := isWordMotionParam(params)
	switch final {
	case 'C': // Right
		if isWord {
			h.ui.state.MoveWordRight()
		} else {
			h.ui.state.MoveRight()
		}
	case 'D': // Left
		if isWord {
			h.ui.state.MoveWordLeft()
		} else {
			h.ui.state.MoveLeft()
		}
	}
}

// handleApplicationCursorMode handles application cursor mode sequences
func (h *KeyHandler) handleApplicationCursorMode(reader *bufio.Reader) {
	var nb byte
	var err error

	if reader != nil {
		// Use provided buffered reader (non-raw mode)
		nb, err = reader.ReadByte()
	} else {
		// Raw mode: read directly from stdin
		var buf [1]byte
		_, err = h.ui.stdin.Read(buf[:])
		nb = buf[0]
	}

	if err != nil {
		return
	}
	switch nb {
	case 'C':
		h.ui.state.MoveRight()
	case 'D':
		h.ui.state.MoveLeft()
	}
}

// handleCtrlC handles Ctrl+C key press
func (h *KeyHandler) handleCtrlC(oldState *term.State) {
	if oldState != nil {
		if f, ok := h.ui.stdin.(*os.File); ok {
			if err := h.ui.term.Restore(int(f.Fd()), oldState); err != nil {
				h.ui.writeError("failed to restore terminal state: %v", err)
			}
		}
	}
	h.ui.writeln("\nExiting...")
	os.Exit(0)
}

// restoreTerminalState restores the terminal from raw mode to cooked mode
func (h *KeyHandler) restoreTerminalState(oldState *term.State) {
	if oldState == nil {
		return
	}
	if f, ok := h.ui.stdin.(*os.File); ok {
		if err := h.ui.term.Restore(int(f.Fd()), oldState); err != nil {
			h.ui.writeError("failed to restore terminal state: %v", err)
		}
	}
}

// reenterRawMode re-enters raw mode after being restored
func (h *KeyHandler) reenterRawMode(oldState *term.State) {
	if oldState == nil {
		return
	}
	if f, ok := h.ui.stdin.(*os.File); ok {
		fd := int(f.Fd())
		if _, err := h.ui.term.MakeRaw(fd); err != nil {
			h.ui.writeError("failed to set terminal to raw mode: %v", err)
		}
	}
}

// handleEnter handles Enter key press
func (h *KeyHandler) handleEnter(oldState *term.State) (bool, []string) {
	// Handle workflow mode
	if h.ui.state.showWorkflow {
		// Execute workflow
		return h.executeWorkflow(oldState)
	}

	if !h.ui.state.HasInput() {
		return true, nil
	}

	selectedCmd := h.ui.state.GetSelectedCommand()
	if selectedCmd == nil {
		return true, nil
	}

	// Restore terminal state BEFORE showing Execute message
	h.restoreTerminalState(oldState)

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
	args, canceled := h.processCommand(selectedCmd.Command)
	if canceled {
		// Re-enter raw mode before returning to main loop
		h.reenterRawMode(oldState)
		return true, nil
	}
	return false, args
}

// processCommand processes the command with placeholder replacement
func (h *KeyHandler) processCommand(cmdTemplate string) ([]string, bool) {
	placeholders := extractPlaceholders(cmdTemplate)

	if len(placeholders) == 0 {
		// No placeholders - execute immediately
		args := []string{"ggc"}
		args = append(args, strings.Fields(cmdTemplate)...)
		return args, false
	}

	// Interactive input for placeholders
	inputs, canceled := h.interactiveInput(placeholders)
	if canceled {
		h.handleSoftCancel(nil)
		return nil, true
	}

	// Placeholder replacement
	finalCmd := cmdTemplate
	for ph, val := range inputs {
		finalCmd = strings.ReplaceAll(finalCmd, "<"+ph+">", val)
	}

	args := []string{"ggc"}
	args = append(args, strings.Fields(finalCmd)...)
	return args, false
}

// interactiveInput provides real-time interactive input for placeholders
func (h *KeyHandler) interactiveInput(placeholders []string) (map[string]string, bool) {
	inputs := make(map[string]string)

	for i, ph := range placeholders {
		h.ui.write("\n")

		// Show progress and prompt
		if len(placeholders) > 1 {
			h.ui.write("%s[%d/%d]%s ",
				h.ui.colors.BrightBlue+h.ui.colors.Bold,
				i+1, len(placeholders),
				h.ui.colors.Reset)
		}

		h.ui.write("%s? %s%s%s: ",
			h.ui.colors.BrightGreen,
			h.ui.colors.BrightWhite+h.ui.colors.Bold,
			ph,
			h.ui.colors.Reset)

		// Get input with real-time feedback
		value, canceled := h.ui.readPlaceholderInput()
		if canceled {
			return nil, true
		}
		inputs[ph] = value

		// Show confirmation
		h.ui.write("%s✓ %s%s: %s%s%s\n",
			h.ui.colors.BrightGreen,
			h.ui.colors.BrightBlue,
			ph,
			h.ui.colors.BrightYellow+h.ui.colors.Bold,
			value,
			h.ui.colors.Reset)
	}

	return inputs, false
}

// getRealTimeInput gets user input with real-time display using raw terminal mode
func (h *KeyHandler) getRealTimeInput() (string, bool) {
	fd := int(os.Stdin.Fd())
	oldState, err := h.ui.term.MakeRaw(fd)
	if err != nil {
		return h.getLineInput()
	}
	defer func() { _ = h.ui.term.Restore(fd, oldState) }()

	return h.processRealTimeInput()
}

// processRealTimeInput handles the main input processing loop
func (h *KeyHandler) processRealTimeInput() (string, bool) {
	reader := bufio.NewReader(os.Stdin)
	inputRunes := make([]rune, 0, initialInputCapacity)
	cursor := 0

	editor := &realTimeEditor{
		ui:         h.ui,
		inputRunes: &inputRunes,
		cursor:     &cursor,
	}

	for {
		r, _, err := reader.ReadRune()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			continue
		}

		result := editor.handleInput(r, reader)
		if result.done {
			return result.text, false
		}
		if result.canceled {
			return "", true
		}
	}
	return string(inputRunes), false
}

// inputResult represents the result of handling input
type inputResult struct {
	done     bool
	canceled bool
	text     string
}

// realTimeEditor handles real-time input editing
type realTimeEditor struct {
	ui         *UI
	inputRunes *[]rune
	cursor     *int
}

// handleInput processes a single input rune
func (e *realTimeEditor) handleInput(r rune, reader *bufio.Reader) inputResult {
	switch r {
	case '\n', '\r':
		return e.handleEnter()
	case 3: // Ctrl+C
		return e.handleCtrlC()
	case 7: // Ctrl+G
		return e.handleSoftCancel()
	case 127, '\b': // Backspace
		e.handleBackspace()
		return inputResult{}
	case 27: // ESC sequences
		if e.shouldSoftCancelOnEscape(reader) {
			return e.handleSoftCancel()
		}
		e.handleEscape(reader)
		return inputResult{}
	default:
		if unicode.IsPrint(r) {
			e.handlePrintableChar(r)
		}
		return inputResult{}
	}
}

// handleEnter processes Enter key
func (e *realTimeEditor) handleEnter() inputResult {
	if len(*e.inputRunes) > 0 {
		e.ui.write("\r\n")
		return inputResult{done: true, text: string(*e.inputRunes)}
	}
	e.ui.write(" %s(required)%s", e.ui.colors.BrightRed, e.ui.colors.Reset)
	return inputResult{}
}

// handleCtrlC processes Ctrl+C
func (e *realTimeEditor) handleCtrlC() inputResult {
	e.ui.write("\r\n%sOperation canceled%s\r\n", e.ui.colors.BrightRed, e.ui.colors.Reset)
	return inputResult{canceled: true}
}

func (e *realTimeEditor) handleSoftCancel() inputResult {
	e.ui.write("\r\n")
	return inputResult{canceled: true}
}

// handleBackspace processes backspace key
func (e *realTimeEditor) handleBackspace() {
	if *e.cursor == 0 {
		return
	}
	start := e.findGraphemeStart(*e.cursor - 1)
	// Compute columns to move left/clear for the removed cluster
	moveCols := e.colsBetween(start, *e.cursor)
	clearedCols := 0
	for i := start; i < *e.cursor; i++ {
		clearedCols += e.runeWidth((*e.inputRunes)[i])
	}
	// Move cursor left, remove runes, and redraw tail
	e.moveLeft(moveCols)
	*e.inputRunes = append((*e.inputRunes)[:start], (*e.inputRunes)[*e.cursor:]...)
	*e.cursor = start
	e.printTailAndReposition(*e.cursor, clearedCols)
}

// findGraphemeStart finds the start of the grapheme cluster ending at the given position
func (e *realTimeEditor) findGraphemeStart(pos int) int {
	start := pos
	start = e.skipCombiningMarks(start)
	start = e.handleRegionalIndicators(start)
	start = e.handleZWJSequences(start)
	if start < 0 {
		start = 0
	}
	return start
}

// skipCombiningMarks skips any trailing variation selectors or combining marks
func (e *realTimeEditor) skipCombiningMarks(start int) int {
	for start >= 0 && (isCombining((*e.inputRunes)[start]) || isVariationSelector((*e.inputRunes)[start])) {
		start--
	}
	return start
}

// handleRegionalIndicators handles regional indicator pairs (flags)
func (e *realTimeEditor) handleRegionalIndicators(start int) int {
	if start >= 0 && isRegionalIndicator((*e.inputRunes)[start]) {
		if start > 0 && isRegionalIndicator((*e.inputRunes)[start-1]) {
			start--
		}
	}
	return start
}

// handleZWJSequences handles ZWJ sequences by including the joiner and previous rune repeatedly
func (e *realTimeEditor) handleZWJSequences(start int) int {
	for {
		if start > 0 && isZWJ((*e.inputRunes)[start-1]) {
			// Include ZWJ and the previous rune
			start -= 2
			// Also include any combining marks attached to the new base
			start = e.skipCombiningMarks(start)
			continue
		}
		break
	}
	return start
}

// handlePrintableChar processes printable characters
func (e *realTimeEditor) handlePrintableChar(r rune) {
	if *e.cursor == len(*e.inputRunes) {
		*e.inputRunes = append(*e.inputRunes, r)
	} else {
		*e.inputRunes = append((*e.inputRunes)[:*e.cursor], append([]rune{r}, (*e.inputRunes)[*e.cursor:]...)...)
	}
	e.ui.write("%s", string(r))
	*e.cursor++
	e.printTailAndReposition(*e.cursor, 0)
}

// handleEscape processes escape sequences for real-time input
func (e *realTimeEditor) handleEscape(reader *bufio.Reader) {
	b, err := reader.ReadByte()
	if err != nil {
		return
	}
	switch b {
	case '[':
		e.handleCSIEscape(reader)
	case 'O':
		e.handleApplicationEscape(reader)
	case 'b':
		e.moveWordLeft()
	case 'f':
		e.moveWordRight()
	case 127, '\b':
		// Option+Backspace: delete previous word
		e.deleteWordLeft()
	}
}

func (e *realTimeEditor) shouldSoftCancelOnEscape(reader *bufio.Reader) bool {
	if reader != nil && reader.Buffered() > 0 {
		return false
	}

	if os.Stdin != nil {
		if pending, err := termio.PendingInput(os.Stdin.Fd()); err == nil {
			return pending == 0
		}
	}

	return false
}

// handleCSIEscape processes CSI escape sequences for real-time input
func (e *realTimeEditor) handleCSIEscape(reader *bufio.Reader) {
	var params []byte
	for {
		nb, err := reader.ReadByte()
		if err != nil {
			return
		}
		if (nb >= 'A' && nb <= 'Z') || nb == '~' {
			e.processCSIEscape(nb, string(params))
			return
		}
		params = append(params, nb)
	}
}

// processCSIEscape handles CSI final byte for real-time input
func (e *realTimeEditor) processCSIEscape(final byte, params string) {
	isWord := isWordMotionParam(params)
	switch final {
	case 'C': // Right
		if isWord {
			e.moveWordRight()
		} else if *e.cursor < len(*e.inputRunes) {
			e.moveRight(e.runeWidth((*e.inputRunes)[*e.cursor]))
			*e.cursor++
		}
	case 'D': // Left
		if isWord {
			e.moveWordLeft()
		} else if *e.cursor > 0 {
			e.moveLeft(e.runeWidth((*e.inputRunes)[*e.cursor-1]))
			*e.cursor--
		}
	}
}

// handleApplicationEscape processes application mode escape sequences
func (e *realTimeEditor) handleApplicationEscape(reader *bufio.Reader) {
	nb, err := reader.ReadByte()
	if err != nil {
		return
	}
	switch nb {
	case 'C':
		if *e.cursor < len(*e.inputRunes) {
			e.moveRight(e.runeWidth((*e.inputRunes)[*e.cursor]))
			*e.cursor++
		}
	case 'D':
		if *e.cursor > 0 {
			e.moveLeft(e.runeWidth((*e.inputRunes)[*e.cursor-1]))
			*e.cursor--
		}
	}
}

// Helper methods for realTimeEditor

// isCombining reports whether r is a combining mark (zero display width)
func isCombining(r rune) bool {
	return unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Me, r) || unicode.Is(unicode.Mc, r)
}

// isVariationSelector reports whether r is a variation selector (zero width)
func isVariationSelector(r rune) bool {
	// U+FE00..U+FE0F (VS1..VS16) and U+E0100..U+E01EF (IVS)
	return (r >= 0xFE00 && r <= 0xFE0F) || (r >= 0xE0100 && r <= 0xE01EF)
}

// isRegionalIndicator reports whether r is a regional indicator rune (used for flags)
func isRegionalIndicator(r rune) bool { return r >= 0x1F1E6 && r <= 0x1F1FF }

// isZWJ reports whether r is ZERO WIDTH JOINER
func isZWJ(r rune) bool { return r == 0x200D }

// isEmoji reports common emoji ranges that should render as width 2 on most terminals
func isEmoji(r rune) bool {
	return isEmojiRange1(r) || isEmojiRange2(r)
}

// isEmojiRange1 checks the first set of emoji Unicode ranges
func isEmojiRange1(r rune) bool {
	return (r >= 0x1F300 && r <= 0x1F5FF) || // Misc Symbols and Pictographs
		(r >= 0x1F600 && r <= 0x1F64F) || // Emoticons
		(r >= 0x1F680 && r <= 0x1F6FF) || // Transport and Map Symbols
		(r >= 0x1F700 && r <= 0x1F77F) || // Alchemical Symbols
		(r >= 0x1F780 && r <= 0x1F7FF) // Geometric Shapes Extended
}

// isEmojiRange2 checks the second set of emoji Unicode ranges
func isEmojiRange2(r rune) bool {
	return (r >= 0x1F800 && r <= 0x1F8FF) || // Supplemental Arrows-C
		(r >= 0x1F900 && r <= 0x1F9FF) || // Supplemental Symbols and Pictographs
		(r >= 0x1FA00 && r <= 0x1FAFF) || // Symbols and Pictographs Extended-A
		(r >= 0x2600 && r <= 0x26FF) || // Misc symbols
		(r >= 0x2700 && r <= 0x27BF) // Dingbats
}

// runeDisplayWidth returns the number of terminal columns used by r
func runeDisplayWidth(r rune) int {
	// Zero-width characters
	if isCombining(r) || isVariationSelector(r) || isZWJ(r) {
		return 0
	}
	// East Asian wide/fullwidth
	switch width.LookupRune(r).Kind() {
	case width.EastAsianFullwidth, width.EastAsianWide:
		return 2
	}
	// Common emoji are typically 2 columns
	if isEmoji(r) {
		return 2
	}
	return 1
}

func (e *realTimeEditor) runeWidth(r rune) int { return runeDisplayWidth(r) }

func (e *realTimeEditor) colsBetween(from, to int) int {
	if from < 0 {
		from = 0
	}
	if to < 0 {
		to = 0
	}
	if from > to {
		from, to = to, from
	}
	cols := 0
	for i := from; i < to && i < len(*e.inputRunes); i++ {
		cols += e.runeWidth((*e.inputRunes)[i])
	}
	return cols
}

func (e *realTimeEditor) moveLeft(cols int) {
	if cols <= 0 {
		return
	}
	e.ui.write("\x1b[%dD", cols)
}

func (e *realTimeEditor) moveRight(cols int) {
	if cols <= 0 {
		return
	}
	e.ui.write("\x1b[%dC", cols)
}

func (e *realTimeEditor) printTailAndReposition(from int, clearedCols int) {
	tailCols := 0
	if from < len(*e.inputRunes) {
		tail := string((*e.inputRunes)[from:])
		e.ui.write("%s", tail)
		for _, rr := range (*e.inputRunes)[from:] {
			tailCols += e.runeWidth(rr)
		}
	}
	if clearedCols > 0 {
		e.ui.write("%s", strings.Repeat(" ", clearedCols))
	}
	e.moveLeft(tailCols + clearedCols)
}

func (e *realTimeEditor) moveWordLeft() {
	if *e.cursor == 0 {
		return
	}
	i := *e.cursor - 1
	for i >= 0 && unicode.IsSpace((*e.inputRunes)[i]) {
		i--
	}
	for i >= 0 && !unicode.IsSpace((*e.inputRunes)[i]) {
		i--
	}
	newPos := i + 1
	delta := e.colsBetween(newPos, *e.cursor)
	e.moveLeft(delta)
	*e.cursor = newPos
}

func (e *realTimeEditor) moveWordRight() {
	n := len(*e.inputRunes)
	if *e.cursor >= n {
		return
	}
	i := *e.cursor
	for i < n && !unicode.IsSpace((*e.inputRunes)[i]) {
		i++
	}
	for i < n && unicode.IsSpace((*e.inputRunes)[i]) {
		i++
	}
	delta := e.colsBetween(*e.cursor, i)
	e.moveRight(delta)
	*e.cursor = i
}

// deleteWordLeft deletes the word before the cursor and updates the display
func (e *realTimeEditor) deleteWordLeft() {
	if *e.cursor == 0 {
		return
	}
	// Find new cursor position at the beginning of previous word
	i := *e.cursor - 1
	for i >= 0 && unicode.IsSpace((*e.inputRunes)[i]) {
		i--
	}
	for i >= 0 && !unicode.IsSpace((*e.inputRunes)[i]) {
		i--
	}
	newPos := i + 1
	// Compute columns to move left and columns to clear
	moveCols := e.colsBetween(newPos, *e.cursor)
	clearedCols := 0
	for j := newPos; j < *e.cursor; j++ {
		clearedCols += e.runeWidth((*e.inputRunes)[j])
	}
	// Move cursor left to newPos
	e.moveLeft(moveCols)
	// Delete runes in [newPos, cursor)
	*e.inputRunes = append((*e.inputRunes)[:newPos], (*e.inputRunes)[*e.cursor:]...)
	*e.cursor = newPos
	// Redraw tail and clear leftover cells
	e.printTailAndReposition(*e.cursor, clearedCols)
}

//nolint:revive // Input character handling inherently requires multiple cases
func (h *KeyHandler) handleInputChar(input *strings.Builder, char rune) (done bool, canceled bool) {
	switch char {
	case '\n', '\r':
		if input.Len() > 0 {
			h.ui.write("\r\n")
			return true, false
		}
		h.ui.write(" %s(required)%s", h.ui.colors.BrightRed, h.ui.colors.Reset)
		return false, false
	case '\b', 127:
		if input.Len() == 0 {
			return false, false
		}
		runes := []rune(input.String())
		if len(runes) == 0 {
			return false, false
		}
		// Identify start of previous grapheme-like cluster
		end := len(runes)
		start := end - 1
		for start >= 0 && (isCombining(runes[start]) || isVariationSelector(runes[start])) {
			start--
		}
		if start >= 0 && isRegionalIndicator(runes[start]) {
			if start > 0 && isRegionalIndicator(runes[start-1]) {
				start--
			}
		}
		for {
			if start > 0 && isZWJ(runes[start-1]) {
				start -= 2
				for start >= 0 && (isCombining(runes[start]) || isVariationSelector(runes[start])) {
					start--
				}
				continue
			}
			break
		}
		if start < 0 {
			start = 0
		}
		// Calculate total columns to clear
		cols := 0
		for i := start; i < end; i++ {
			cols += runeDisplayWidth(runes[i])
		}
		// Update input
		input.Reset()
		input.WriteString(string(runes[:start]))
		// Clear terminal cells
		for i := 0; i < cols; i++ {
			h.ui.write("\b \b")
		}
		return false, false
	case 3: // Ctrl+C
		h.ui.write("\r\n%sOperation canceled%s\r\n", h.ui.colors.BrightRed, h.ui.colors.Reset)
		return true, true
	default:
		// Accept all printable characters including multibyte
		if unicode.IsPrint(char) {
			input.WriteRune(char)
			h.ui.write("%s", string(char))
		}
		return false, false
	}
}

// getLineInput provides fallback line-based input when raw mode is not available
func (h *KeyHandler) getLineInput() (string, bool) {
	reader := bufio.NewReader(h.ui.stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", true
		}
		line = strings.TrimSpace(line)
		if line != "" {
			return line, false
		}
		h.ui.write("%s(required)%s ",
			h.ui.colors.BrightRed,
			h.ui.colors.Reset)
	}
}

// NewUI creates a new UI with the provided git client and loads keybindings from config
func NewUI(gitClient git.StatusInfoReader, router ...CommandRouter) *UI {
	colors := NewANSIColors()

	renderer := &Renderer{
		writer: os.Stdout,
		colors: colors,
	}
	renderer.updateSize()

	state := &UIState{
		selected:     0,
		input:        "",
		filtered:     []CommandInfo{},
		context:      ContextGlobal, // Start in global context
		contextStack: []Context{},
	}

	// Load config and create resolver
	var cfg *config.Config
	if ops, ok := gitClient.(git.ConfigOps); ok {
		configManager := config.NewConfigManager(ops)
		// Load config - if it fails, we'll use defaults from manager
		_ = configManager.Load()
		cfg = configManager.GetConfig()
	} else {
		// Fallback to empty config (built-in defaults and profiles will be used)
		cfg = &config.Config{}
	}

	// Create KeyBinding resolver and register built-in profiles
	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)
	contextManager := NewContextManager(resolver)

	// Determine which profile to use (default to "default" profile)
	profile := ProfileDefault
	if cfg.Interactive.Profile != "" {
		switch Profile(cfg.Interactive.Profile) {
		case ProfileEmacs, ProfileVi, ProfileReadline:
			profile = Profile(cfg.Interactive.Profile)
		default:
			fmt.Fprintf(os.Stderr, "Warning: Unknown profile '%s', using default\n", cfg.Interactive.Profile)
		}
	}

	// Resolve contextual keybindings for all contexts
	contextualMap, err := resolver.ResolveContextual(profile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to resolve keybindings: %v. Using defaults.\n", err)
		// Fallback to legacy defaults
		keyMap := DefaultKeyBindingMap()
		contextualMap = &ContextualKeyBindingMap{
			Profile:  ProfileDefault,
			Platform: DetectPlatform(),
			Terminal: DetectTerminal(),
			Contexts: map[Context]*KeyBindingMap{
				ContextGlobal:  keyMap,
				ContextInput:   keyMap,
				ContextResults: keyMap,
				ContextSearch:  keyMap,
			},
		}
	}

	ui := &UI{
		stdin:      os.Stdin,
		stdout:     os.Stdout,
		stderr:     os.Stderr,
		term:       termio.DefaultTerminal{},
		renderer:   renderer,
		state:      state,
		colors:     colors,
		gitClient:  gitClient,
		gitStatus:  getGitStatus(gitClient),
		contextMgr: contextManager,
		profile:    profile,
		workflow:   NewWorkflow(),
	}
	state.onContextChange = func(_ Context, newCtx Context) {
		contextManager.SetContext(newCtx)
	}

	ui.handler = &KeyHandler{
		ui:            ui,
		contextualMap: contextualMap,
	}

	// Set up workflow executor if router is provided
	if len(router) > 0 && router[0] != nil {
		ui.workflowEx = NewWorkflowExecutor(router[0], ui)
	}

	return ui
}

// ToggleWorkflowView toggles between normal command view and workflow view
func (ui *UI) ToggleWorkflowView() {
	ui.state.showWorkflow = !ui.state.showWorkflow
}

// AddToWorkflow adds a command to the workflow
func (ui *UI) AddToWorkflow(command string, args []string, description string) int {
	return ui.workflow.AddStep(command, args, description)
}

func (ui *UI) resetToSearchMode() bool {
	if ui == nil || ui.state == nil {
		return false
	}

	state := ui.state
	active := state.HasInput() || state.showWorkflow || len(state.contextStack) > 0 || state.GetCurrentContext() != ContextGlobal
	state.ClearInput()
	state.selected = 0
	state.contextStack = nil
	state.SetContext(ContextGlobal)
	state.showWorkflow = false
	return active
}

func (ui *UI) readPlaceholderInput() (string, bool) {
	if ui == nil || ui.handler == nil {
		return "", true
	}
	return ui.handler.getRealTimeInput()
}

// ClearWorkflow removes all steps from the workflow
func (ui *UI) ClearWorkflow() {
	ui.workflow.Clear()
}

// ExecuteWorkflow executes the current workflow
func (ui *UI) ExecuteWorkflow() error {
	if ui.workflowEx == nil {
		return fmt.Errorf("workflow executor not initialized")
	}

	if ui.workflow.IsEmpty() {
		return fmt.Errorf("workflow is empty")
	}

	return ui.workflowEx.Execute(ui.workflow)
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

var commands = buildInteractiveCommands()

// InteractiveUI provides an incremental search interactive UI with custom git client.
// Returns the selected command as []string (nil if nothing selected)
func InteractiveUI(gitClient git.StatusInfoReader) []string {
	ui := NewUI(gitClient)
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

func (ui *UI) notifySoftCancel() {
	ui.softCancelFlash.Store(true)
}

func (ui *UI) consumeSoftCancelFlash() bool {
	return ui.softCancelFlash.Swap(false)
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
	var restoreCursor func()
	defer func() {
		_, _ = fmt.Fprint(r.writer, "\x1b[?7h")
		if restoreCursor != nil {
			restoreCursor()
		}
		showCursor(r.writer)
	}()

	// Update terminal size
	r.updateSize()

	// Render each section
	r.renderHeader(ui)
	r.renderSoftCancelFlash(ui)

	// Render workflow status if it has steps and no search input
	if !ui.workflow.IsEmpty() && !state.showWorkflow && state.input == "" {
		r.renderWorkflowStatus(ui)
	}

	// Render content based on state
	if state.showWorkflow {
		r.renderWorkflowView(ui, state)
	} else {
		r.renderSearchPrompt(ui, state)
		restoreCursor = r.saveCursorAtSearchPrompt(state)

		switch {
		case state.input == "":
			if ui.workflow.IsEmpty() {
				r.renderEmptyState(ui)
			} else {
				r.renderEmptyStateWithWorkflow(ui)
			}
			r.writeEmptyLine()
			r.renderSearchKeybinds(ui)
		case len(state.filtered) == 0:
			r.renderNoMatches(ui, state)
		default:
			r.renderCommandList(ui, state)
		}
	}

}

func (r *Renderer) renderSoftCancelFlash(ui *UI) {
	if !ui.consumeSoftCancelFlash() {
		return
	}
	alert := fmt.Sprintf("%s⚠️  Operation canceled%s", r.colors.BrightRed+r.colors.Bold, r.colors.Reset)
	r.writeColorln(ui, alert)
	r.writeColorln(ui, "")
}

// renderHeader renders the title, git status, and navigation subtitle
func (r *Renderer) renderHeader(ui *UI) {
	// Modern header with title
	title := fmt.Sprintf("%s%s🚀 ggc Interactive Mode%s",
		r.colors.BrightCyan+r.colors.Bold,
		r.colors.Reset,
		r.colors.Reset)
	r.writeColorln(ui, title)

	// Git status information
	if ui.gitStatus != nil {
		r.renderGitStatus(ui, ui.gitStatus)
	}

}

// renderSearchPrompt renders the search input with cursor
func (r *Renderer) renderSearchPrompt(ui *UI, state *UIState) {
	inputWithCursor := r.formatInputWithCursor(state)

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
}

func (r *Renderer) saveCursorAtSearchPrompt(state *UIState) func() {
	linesUp := 2
	if state.input != "" {
		linesUp++
	}
	_, _ = fmt.Fprintf(r.writer, "\x1b[%dA", linesUp)
	const prefix = "┌─ Search: "
	// Compute display width (columns) of the prefix using runeDisplayWidth
	prefixCols := 0
	for _, pr := range prefix {
		prefixCols += runeDisplayWidth(pr)
	}
	// Compute display width up to the logical cursor position
	runes := []rune(state.input)
	cursorPos := state.cursorPos
	if cursorPos > len(runes) {
		cursorPos = len(runes)
	}
	cursorWidth := 0
	for _, rr := range runes[:cursorPos] {
		cursorWidth += runeDisplayWidth(rr)
	}
	column := prefixCols + cursorWidth + 1
	if column < 1 {
		column = 1
	}
	_, _ = fmt.Fprintf(r.writer, "\x1b[%dG", column)
	_, _ = fmt.Fprint(r.writer, "\x1b[s")
	_, _ = fmt.Fprintf(r.writer, "\x1b[%dB", linesUp)
	return func() {
		_, _ = fmt.Fprint(r.writer, "\x1b[u")
	}
}

// formatInputWithCursor formats the input string with cursor position
func (r *Renderer) formatInputWithCursor(state *UIState) string {
	if state.input == "" {
		return fmt.Sprintf("%s█%s", r.colors.BrightWhite+r.colors.Bold, r.colors.Reset)
	}

	inputRunes := []rune(state.input)
	beforeCursor := string(inputRunes[:state.cursorPos])
	afterCursor := string(inputRunes[state.cursorPos:])
	cursor := "│"
	if state.cursorPos >= utf8.RuneCountInString(state.input) {
		cursor = "█"
	}

	return fmt.Sprintf("%s%s%s%s%s%s%s",
		r.colors.BrightYellow,
		beforeCursor,
		r.colors.BrightWhite+r.colors.Bold,
		cursor,
		r.colors.Reset+r.colors.BrightYellow,
		afterCursor,
		r.colors.Reset)
}

// renderEmptyState renders the empty input state
func (r *Renderer) renderEmptyState(ui *UI) {
	r.writeColorln(ui, fmt.Sprintf("%s💭 %sStart typing to search commands...%s",
		r.colors.BrightBlue, r.colors.BrightBlack, r.colors.Reset))
}

func (r *Renderer) buildSearchKeybindEntries(ui *UI) []keybindHelpEntry {
	entries := []keybindHelpEntry{
		{key: "←/→", desc: "Move cursor"},
		{key: "Ctrl+←/→", desc: "Move by word"},
		{key: "Option+←/→", desc: "Move by word (macOS)"},
	}
	// Future: extend this helper for additional contexts such as workflow views.

	var km *KeyBindingMap
	if ui != nil && ui.handler != nil {
		km = ui.handler.GetCurrentKeyMap()
	}
	if km == nil {
		km = DefaultKeyBindingMap()
	}

	defaultMap := DefaultKeyBindingMap()

	appendDynamic := func(primary []KeyStroke, fallback []KeyStroke, desc string) {
		keys := primary
		if len(keys) == 0 {
			keys = fallback
		}
		if len(keys) == 0 {
			return
		}
		formatted := FormatKeyStrokesForDisplay(keys)
		if formatted == "" || formatted == "none" {
			return
		}
		entries = append(entries, keybindHelpEntry{key: formatted, desc: desc})
	}

	appendDynamic(km.ClearLine, defaultMap.ClearLine, "Clear all input")
	appendDynamic(km.DeleteWord, defaultMap.DeleteWord, "Delete word")
	appendDynamic(km.DeleteToEnd, defaultMap.DeleteToEnd, "Delete to end")
	appendDynamic(km.MoveToBeginning, defaultMap.MoveToBeginning, "Move to beginning")
	appendDynamic(km.MoveToEnd, defaultMap.MoveToEnd, "Move to end")

	entries = append(entries, keybindHelpEntry{key: "Backspace", desc: "Delete character"})
	entries = append(entries, keybindHelpEntry{key: "Enter", desc: "Execute selected command"})

	appendDynamic(km.AddToWorkflow, defaultMap.AddToWorkflow, "Add to workflow")
	appendDynamic(km.ToggleWorkflowView, defaultMap.ToggleWorkflowView, "Toggle workflow view")

	entries = append(entries, keybindHelpEntry{key: "Ctrl+c", desc: "Quit"})

	return entries
}

func (r *Renderer) renderKeybindEntries(ui *UI, entries []keybindHelpEntry) {
	if len(entries) == 0 {
		return
	}

	r.writeColorln(ui, fmt.Sprintf("%s⌨️  %sAvailable keybinds:%s",
		r.colors.BrightBlue, r.colors.BrightWhite+r.colors.Bold, r.colors.Reset))

	for _, entry := range entries {
		r.writeColorln(ui, fmt.Sprintf("   %s%s%s  %s%s%s",
			r.colors.BrightGreen+r.colors.Bold,
			entry.key,
			r.colors.Reset,
			r.colors.BrightBlack,
			entry.desc,
			r.colors.Reset))
	}
}

// renderNoMatches renders the no matches found state with keybind help
func (r *Renderer) renderNoMatches(ui *UI, state *UIState) {
	// No matches message
	r.writeColorln(ui, fmt.Sprintf("%s🔍 %sNo commands found for '%s%s%s'%s",
		r.colors.BrightYellow,
		r.colors.BrightWhite,
		r.colors.BrightYellow+r.colors.Bold,
		state.input,
		r.colors.Reset+r.colors.BrightWhite,
		r.colors.Reset))

	r.writeEmptyLine()
	r.renderKeybindEntries(ui, r.buildSearchKeybindEntries(ui))
}

// renderSearchKeybinds renders keybinds available in search UI
func (r *Renderer) renderSearchKeybinds(ui *UI) {
	r.renderKeybindEntries(ui, r.buildSearchKeybindEntries(ui))
}

// renderWorkflowKeybinds renders keybinds available in workflow UI
func (r *Renderer) renderWorkflowKeybinds(ui *UI) {
	// Workflow-specific keybinds
	keybinds := []struct{ key, desc string }{
		{"Enter", "Execute workflow"},
		{"c", "Clear workflow"},
		{"Ctrl+t", "Back to search view"},
		{"Ctrl+c", "Quit"},
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
}

// renderCommandList renders the filtered command list
func (r *Renderer) renderCommandList(ui *UI, state *UIState) {
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
		r.renderCommandItem(ui, cmd, i, state.selected, maxCmdLen)
	}
}

// renderCommandItem renders a single command item
func (r *Renderer) renderCommandItem(ui *UI, cmd CommandInfo, index, selected, maxCmdLen int) {
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

	if index == selected {
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

// writeColorln writes a colored line to the terminal
func (r *Renderer) writeColorln(_ *UI, text string) {
	// Move to line start, clear line, write content, then CRLF
	_, _ = fmt.Fprint(r.writer, "\r\x1b[K")
	_, _ = fmt.Fprint(r.writer, text+"\r\n")
}

// renderEmptyStateWithWorkflow renders the empty state with workflow info
func (r *Renderer) renderEmptyStateWithWorkflow(ui *UI) {
	r.writeColorln(ui, fmt.Sprintf("%s📝 Start typing to search commands...%s",
		r.colors.BrightBlue,
		r.colors.Reset))
}

// renderWorkflowStatus renders workflow information at the top of the UI
func (r *Renderer) renderWorkflowStatus(ui *UI) {
	steps := ui.workflow.GetSteps()
	if len(steps) == 0 {
		return
	}

	// Workflow status bar
	statusText := fmt.Sprintf("%s📋 Workflow Ready (%d steps):%s",
		r.colors.BrightYellow+r.colors.Bold,
		len(steps),
		r.colors.Reset)

	// Show first few steps inline
	stepTexts := make([]string, 0, min(3, len(steps)))
	for i, step := range steps[:min(3, len(steps))] {
		stepText := fmt.Sprintf("%s%d.%s %s%s%s",
			r.colors.BrightBlue+r.colors.Bold,
			i+1,
			r.colors.Reset,
			r.colors.BrightGreen,
			step.Description,
			r.colors.Reset)
		stepTexts = append(stepTexts, stepText)
	}

	if len(steps) > 3 {
		stepTexts = append(stepTexts, fmt.Sprintf("%s...+%d more%s",
			r.colors.BrightBlack, len(steps)-3, r.colors.Reset))
	}

	r.writeColorln(ui, statusText+" "+strings.Join(stepTexts, " → "))
	r.writeColorln(ui, "")
}

// renderWorkflowView renders the detailed workflow view
func (r *Renderer) renderWorkflowView(ui *UI, _ *UIState) {
	steps := ui.workflow.GetSteps()

	// Detailed workflow header
	r.writeColorln(ui, fmt.Sprintf("%s📋 Workflow Details (%d steps)%s",
		r.colors.BrightYellow+r.colors.Bold,
		len(steps),
		r.colors.Reset))
	r.writeColorln(ui, "")

	if len(steps) == 0 {
		r.writeColorln(ui, fmt.Sprintf("%s  No steps in workflow%s",
			r.colors.BrightBlack,
			r.colors.Reset))
		r.writeColorln(ui, "")

		// Render workflow keybinds even for empty workflow
		r.renderWorkflowKeybinds(ui)
		return
	}

	// Render all workflow steps
	for i, step := range steps {
		stepLine := fmt.Sprintf("  %s%d.%s %s%s%s",
			r.colors.BrightBlue+r.colors.Bold,
			i+1,
			r.colors.Reset,
			r.colors.BrightGreen+r.colors.Bold,
			step.Description,
			r.colors.Reset)
		r.writeColorln(ui, stepLine)
	}

	r.writeColorln(ui, "")

	// Render workflow keybinds
	r.renderWorkflowKeybinds(ui)
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
	if len(filtered) == 0 {
		return 0
	}

	maxLen := 0
	for _, cmd := range filtered {
		if len(cmd.Command) > maxLen {
			maxLen = len(cmd.Command)
		}
	}
	return maxLen
}

// setupTerminal configures terminal raw mode and returns the old state and error status
func (ui *UI) setupTerminal() (*term.State, bool) {
	var oldState *term.State
	if f, ok := ui.stdin.(*os.File); ok {
		// Check if it's a real terminal (TTY) - skip for non-TTY or when term is nil
		if ui.term == nil || !term.IsTerminal(int(f.Fd())) {
			// Not a real terminal, skip raw mode setup for debugging
			return nil, true
		}

		fd := int(f.Fd())
		var err error
		oldState, err = ui.term.MakeRaw(fd)
		if err != nil {
			ui.writeError("Failed to set terminal to raw mode: %v", err)
			return nil, false
		}
	}
	return oldState, true
}

// Run executes the interactive UI
func (ui *UI) Run() []string {
	oldState, reader, isRawMode := ui.initializeTerminal()
	if oldState == nil && isRawMode {
		return nil
	}

	// Set up terminal restoration for raw mode
	if f, ok := ui.stdin.(*os.File); ok && isRawMode {
		fd := int(f.Fd())
		defer func() {
			if err := ui.term.Restore(fd, oldState); err != nil {
				ui.writeError("failed to restore terminal state: %v", err)
			}
		}()
	}

	return ui.runMainLoop(reader, isRawMode, oldState)
}

// initializeTerminal sets up the terminal and returns the old state, reader, and raw mode status
func (ui *UI) initializeTerminal() (*term.State, *bufio.Reader, bool) {
	oldState, ok := ui.setupTerminal()
	if !ok {
		return nil, nil, false
	}

	// Check if we're in raw mode (real terminal) or not
	isRawMode := oldState != nil

	// Set up reader based on mode
	var reader *bufio.Reader
	if !isRawMode {
		reader = bufio.NewReader(ui.stdin)
	}

	return oldState, reader, isRawMode
}

// runMainLoop handles the main input loop
func (ui *UI) runMainLoop(reader *bufio.Reader, isRawMode bool, oldState *term.State) []string {
	if isRawMode {
		if ui.reader == nil {
			ui.reader = bufio.NewReader(ui.stdin)
		}
	} else {
		ui.reader = reader
	}

	for {
		ui.state.UpdateFiltered()
		ui.renderer.Render(ui, ui.state)

		r, err := ui.readNextRune(reader, isRawMode)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			continue // Skip this iteration for other errors
		}

		// Handle key input with rune
		isSingleByte := isRawMode // In raw mode, we read single bytes; in buffered mode, we read full runes
		shouldContinue, result := ui.handler.HandleKey(r, isSingleByte, oldState, reader)
		if !shouldContinue {
			return result
		}
	}
}

// readNextRune reads the next rune from input based on the mode
func (ui *UI) readNextRune(reader *bufio.Reader, isRawMode bool) (rune, error) {
	if isRawMode {
		for {
			// Read single byte directly from stdin in raw mode
			var buf [1]byte
			n, readErr := ui.stdin.Read(buf[:])

			if readErr != nil {
				return 0, readErr
			}

			if n == 0 {
				continue // Try again if no bytes read
			}

			return rune(buf[0]), nil
		}
	}

	// Use buffered reader for non-TTY
	r, _, err := reader.ReadRune()
	return r, err
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

// addCommandToWorkflow adds a command to the workflow (preserving placeholders)
func (h *KeyHandler) addCommandToWorkflow(cmdTemplate string) {
	// Don't process placeholders here - save the template as-is
	// Placeholders will be resolved during workflow execution

	// Parse command and arguments from template
	parts := strings.Fields(cmdTemplate)
	if len(parts) == 0 {
		return
	}

	command := parts[0]
	args := parts[1:]

	// Add template to workflow (with placeholders intact)
	id := h.ui.AddToWorkflow(command, args, cmdTemplate)

	// Show success message
	placeholders := extractPlaceholders(cmdTemplate)
	if len(placeholders) > 0 {
		h.ui.write("\n%s🎯 Added to workflow!%s\n",
			h.ui.colors.BrightGreen+h.ui.colors.Bold, h.ui.colors.Reset)
		h.ui.write("%s  Step %d: %s%s%s %s(will prompt for: %v)%s\n",
			h.ui.colors.BrightCyan, id, h.ui.colors.BrightWhite+h.ui.colors.Bold, cmdTemplate, h.ui.colors.Reset,
			h.ui.colors.BrightYellow, placeholders, h.ui.colors.Reset)
	} else {
		h.ui.write("\n%s🎯 Added to workflow!%s\n",
			h.ui.colors.BrightGreen+h.ui.colors.Bold, h.ui.colors.Reset)
		h.ui.write("%s  Step %d: %s%s%s\n",
			h.ui.colors.BrightCyan, id, h.ui.colors.BrightWhite+h.ui.colors.Bold, cmdTemplate, h.ui.colors.Reset)
	}
	h.ui.write("%s  Press 'Ctrl+t' to view workflow, or continue adding more commands%s\n\n",
		h.ui.colors.BrightBlack, h.ui.colors.Reset)
}

// clearWorkflow clears all steps from workflow
func (h *KeyHandler) clearWorkflow() {
	h.ui.ClearWorkflow()
	h.ui.write("%s🧹 Workflow cleared%s\n", h.ui.colors.BrightYellow, h.ui.colors.Reset)
}

// executeWorkflow executes the current workflow
func (h *KeyHandler) executeWorkflow(oldState *term.State) (bool, []string) {
	if h.ui.workflow.IsEmpty() {
		h.ui.write("%sWorkflow is empty. Add some steps first!%s\n",
			h.ui.colors.BrightRed, h.ui.colors.Reset)
		return true, nil
	}

	// Restore terminal state before execution
	if oldState != nil {
		if f, ok := h.ui.stdin.(*os.File); ok {
			if err := h.ui.term.Restore(int(f.Fd()), oldState); err != nil {
				h.ui.writeError("failed to restore terminal state: %v", err)
			}
		}
	}

	// Clear screen and execute workflow
	clearScreen(h.ui.stdout)

	err := h.ui.ExecuteWorkflow()
	if errors.Is(err, ErrWorkflowCanceled) {
		h.handleSoftCancel(oldState)
		return true, nil
	}
	if err != nil {
		fmt.Printf("\n❌ Workflow execution failed: %v\n", err)
	} else {
		fmt.Printf("\n✨ Workflow preserved for reuse. Press 'Ctrl+t' to view or modify.\n")
	}

	// Keep workflow for reuse - don't clear it
	return false, []string{"ggc", InteractiveWorkflowCommand}
}
