// Package interactive houses interactive UI types and helpers shared across the application.
package interactive

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"golang.org/x/term"

	kb "github.com/bmf-san/ggc/v7/internal/keybindings"
	"github.com/bmf-san/ggc/v7/internal/termio"
)

// KeyHandler manages keyboard input processing
type KeyHandler struct {
	ui            *UI
	contextualMap *kb.ContextualKeyBindingMap
}

// GetCurrentKeyMap returns the appropriate keybinding map for the current context
func (h *KeyHandler) GetCurrentKeyMap() *kb.KeyBindingMap {
	if h == nil {
		return kb.DefaultKeyBindingMap()
	}
	if h.contextualMap != nil && h.ui != nil && h.ui.state != nil {
		currentContext := h.ui.state.GetCurrentContext()
		if contextMap, exists := h.contextualMap.GetContext(currentContext); exists && contextMap != nil {
			return contextMap
		}
		if contextMap, exists := h.contextualMap.GetContext(kb.ContextGlobal); exists && contextMap != nil {
			return contextMap
		}
	}

	return kb.DefaultKeyBindingMap()
}

// HandleKey processes UTF-8 rune input and returns true if should continue
// This method handles both single-byte (ASCII/control) and multibyte characters
func (h *KeyHandler) HandleKey(r rune, _ bool, oldState *term.State, reader *bufio.Reader) (bool, []string) {
	// Set the reader for consistent access during escape sequence handling
	h.ui.reader = reader
	// Handle workflow-specific keys first (Tab, etc.)
	if handled, cont, result := h.handleWorkflowKeys(r, oldState); handled {
		return cont, result
	}

	// Handle control characters (ASCII range)
	if r < 128 && unicode.IsControl(r) { // ASCII control characters
		if handled, shouldContinue, result := h.handleControlChar(byte(r), oldState, reader); handled {
			return shouldContinue, result
		}
	}

	// Handle printable characters (both ASCII and multibyte)
	// Workflow mode has no input field, so ignore printable characters
	if unicode.IsPrint(r) {
		if !h.ui.state.IsWorkflowMode() {
			h.ui.state.AddRune(r)
		}
	}
	return true, nil
}

// handleWorkflowKeys processes workflow-related key bindings and returns (handled, result)
func (h *KeyHandler) handleWorkflowKeys(r rune, oldState *term.State) (bool, bool, []string) {
	switch h.ui.state.mode {
	case ModeWorkflow:
		return h.handleWorkflowModeKeys(r, oldState)
	case ModeSearch:
		return h.handleSearchModeWorkflowKeys(r)
	default:
		return false, true, nil
	}
}

func (h *KeyHandler) handleSearchModeWorkflowKeys(r rune) (bool, bool, []string) {
	km := h.GetCurrentKeyMap()
	keyStroke := kb.NewCharKeyStroke(r)

	if km.MatchesKeyStroke("add_to_workflow", keyStroke) {
		if h.ui.state.HasInput() {
			if cmd := h.ui.state.GetSelectedCommand(); cmd != nil {
				h.addCommandToWorkflow(cmd.Command)
				h.ui.state.ClearInput()
			}
		}
		return true, true, nil
	}
	return false, true, nil
}

func (h *KeyHandler) handleWorkflowModeKeys(r rune, oldState *term.State) (bool, bool, []string) {
	if handled := h.handleWorkflowModeShortcut(r, oldState); handled {
		return true, true, nil
	}
	if handled := h.handleWorkflowModeBindings(r); handled {
		return true, true, nil
	}
	return false, true, nil
}

func (h *KeyHandler) handleWorkflowModeShortcut(r rune, oldState *term.State) bool {
	switch r {
	case 'x':
		h.executeWorkflow(oldState)
		return true
	case 'n':
		h.createWorkflow()
		return true
	case 'd':
		h.deleteActiveWorkflow()
		return true
	}
	return false
}

func (h *KeyHandler) handleWorkflowModeBindings(r rune) bool {
	keyStroke := kb.NewCharKeyStroke(r)

	if h.handleWorkflowAdd(keyStroke) {
		return true
	}
	if h.handleWorkflowClear(keyStroke) {
		return true
	}
	return false
}

func (h *KeyHandler) handleWorkflowAdd(keyStroke kb.KeyStroke) bool {
	if !h.ui.state.IsInputFocused() || !h.ui.state.HasInput() {
		return false
	}
	km := h.GetCurrentKeyMap()
	if !km.MatchesKeyStroke("add_to_workflow", keyStroke) {
		return false
	}
	if cmd := h.ui.state.GetSelectedCommand(); cmd != nil {
		h.addCommandToWorkflow(cmd.Command)
		h.ui.state.ClearInput()
	}
	return true
}

func (h *KeyHandler) handleWorkflowClear(keyStroke kb.KeyStroke) bool {
	if h.ui.state.IsInputFocused() {
		return false
	}
	km := h.GetCurrentKeyMap()
	if !km.MatchesKeyStroke("clear_workflow", keyStroke) {
		return false
	}
	h.clearWorkflow()
	return true
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
		ctrlStroke := kb.NewCtrlKeyStroke(rune('a' + b - 1))

		// Workflow mode: simplified key handling (no input field)
		if h.ui.state.IsWorkflowMode() {
			if km.MatchesKeyStroke("workflow_delete", ctrlStroke) {
				h.deleteActiveWorkflow()
				return true, true, nil
			}
			// Ctrl+n/p navigate workflow list
			if km.MatchesKeyStroke("move_down", ctrlStroke) {
				h.moveWorkflowList(1)
				return true, true, nil
			}
			if km.MatchesKeyStroke("move_up", ctrlStroke) {
				h.moveWorkflowList(-1)
				return true, true, nil
			}
			// Ctrl+t toggles back to search mode
			if km.MatchesKeyStroke("toggle_workflow_view", ctrlStroke) {
				h.ui.ToggleWorkflowView()
				return true, true, nil
			}
			if km.MatchesKeyStroke("soft_cancel", ctrlStroke) {
				h.handleSoftCancel(oldState)
				return true, true, nil
			}
			// Ctrl+C must still work in workflow mode - fall through to switch
			if b != 3 {
				// Ignore other input-related Ctrl keys in workflow mode
				return false, true, nil
			}
		}

		// Search mode: full key handling
		if km.MatchesKeyStroke("move_up", ctrlStroke) {
			h.handleMoveUp()
			return true, true, nil
		}
		if km.MatchesKeyStroke("move_down", ctrlStroke) {
			h.handleMoveDown()
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
	if km == nil || !km.MatchesKeyStroke("soft_cancel", kb.NewEscapeKeyStroke()) {
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

func (h *KeyHandler) handleMoveUp() {
	switch h.ui.state.mode {
	case ModeWorkflow:
		h.moveWorkflowList(-1)
	default:
		h.ui.state.MoveUp()
	}
}

func (h *KeyHandler) handleMoveDown() {
	switch h.ui.state.mode {
	case ModeWorkflow:
		h.moveWorkflowList(1)
	default:
		h.ui.state.MoveDown()
	}
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

// processCSIFinalByte processes the final byte of a CSI sequence
func (h *KeyHandler) processCSIFinalByte(final byte, params string) {
	isWord := isWordMotionParam(params)

	// Build the full escape sequence for keybinding matching
	seq := h.buildCSISequence(final, params)
	keyStroke := kb.NewRawKeyStroke(seq)
	km := h.GetCurrentKeyMap()

	// Try keybinding-based handling first
	if h.tryArrowKeybinding(km, keyStroke) {
		return
	}

	// Fallback to default cursor movement and word navigation
	h.handleDefaultArrowMovement(final, isWord)
}

// buildCSISequence builds a CSI escape sequence
func (h *KeyHandler) buildCSISequence(final byte, params string) []byte {
	if params == "" {
		return []byte{27, '[', final}
	}
	seq := append([]byte{27, '['}, []byte(params)...)
	return append(seq, final)
}

// tryArrowKeybinding attempts to handle arrow keys via keybindings
func (h *KeyHandler) tryArrowKeybinding(km *kb.KeyBindingMap, keyStroke kb.KeyStroke) bool {
	if km.MatchesKeyStroke("move_up", keyStroke) {
		// Arrow keys don't navigate in workflow mode (use Ctrl+N/P)
		if !h.ui.state.IsWorkflowMode() {
			h.handleMoveUp()
		}
		return true
	}
	if km.MatchesKeyStroke("move_down", keyStroke) {
		// Arrow keys don't navigate in workflow mode (use Ctrl+N/P)
		if !h.ui.state.IsWorkflowMode() {
			h.handleMoveDown()
		}
		return true
	}
	if km.MatchesKeyStroke("move_left", keyStroke) {
		if h.ui.state.IsInputFocused() {
			h.ui.state.MoveLeft()
		}
		return true
	}
	if km.MatchesKeyStroke("move_right", keyStroke) {
		if h.ui.state.IsInputFocused() {
			h.ui.state.MoveRight()
		}
		return true
	}
	return false
}

// handleDefaultArrowMovement handles default arrow key behavior
func (h *KeyHandler) handleDefaultArrowMovement(final byte, isWord bool) {
	if !h.ui.state.IsInputFocused() {
		return
	}
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
	nb, err := h.readNextByte(reader)
	if err != nil {
		return
	}

	// Build the full escape sequence: ESC O <final>
	seq := []byte{27, 'O', nb}
	keyStroke := kb.NewRawKeyStroke(seq)
	km := h.GetCurrentKeyMap()

	// Try keybinding-based handling first
	if h.tryArrowKeybinding(km, keyStroke) {
		return
	}

	// Fallback to default arrow key behavior
	h.handleDefaultAppCursorMovement(nb)
}

// handleDefaultAppCursorMovement handles default application cursor mode arrow keys
func (h *KeyHandler) handleDefaultAppCursorMovement(nb byte) {
	switch nb {
	case 'A':
		// Arrow keys don't navigate in workflow mode (use Ctrl+N/P)
		if !h.ui.state.IsWorkflowMode() {
			h.handleMoveUp()
		}
	case 'B':
		// Arrow keys don't navigate in workflow mode (use Ctrl+N/P)
		if !h.ui.state.IsWorkflowMode() {
			h.handleMoveDown()
		}
	case 'C':
		if h.ui.state.IsInputFocused() {
			h.ui.state.MoveRight()
		}
	case 'D':
		if h.ui.state.IsInputFocused() {
			h.ui.state.MoveLeft()
		}
	}
}

func (h *KeyHandler) moveWorkflowList(delta int) {
	summaries := h.ui.listWorkflows()
	if len(summaries) == 0 {
		return
	}
	h.ui.state.SetWorkflowListIndex(h.ui.state.workflowListIdx+delta, len(summaries))
	idx := h.ui.state.workflowListIdx
	if idx < 0 || idx >= len(summaries) {
		return
	}
	selected := summaries[idx]
	if h.ui.workflowMgr.SetActive(selected.ID) {
		h.ui.updateWorkflowPointer()
	}
}

func (h *KeyHandler) createWorkflow() {
	if h.ui.workflowMgr == nil {
		return
	}
	newID := h.ui.workflowMgr.CreateWorkflow("")
	summaries := h.ui.listWorkflows()
	for i, summary := range summaries {
		if summary.ID == newID {
			h.ui.state.SetWorkflowListIndex(i, len(summaries))
			break
		}
	}
	h.ui.updateWorkflowPointer()
	h.ui.write("%s✨ Created workflow #%d%s\n", h.ui.colors.BrightGreen, newID, h.ui.colors.Reset)
}

func (h *KeyHandler) deleteActiveWorkflow() {
	if h.ui.workflowMgr == nil {
		return
	}
	activeID := h.ui.workflowMgr.GetActiveID()
	if activeID == 0 {
		h.ui.write("%sNo active workflow to delete%s\n", h.ui.colors.BrightYellow, h.ui.colors.Reset)
		return
	}
	newActive, ok := h.ui.workflowMgr.DeleteWorkflow(activeID)
	if !ok {
		h.ui.write("%sUnable to delete workflow #%d%s\n", h.ui.colors.BrightYellow, activeID, h.ui.colors.Reset)
		return
	}
	summaries := h.ui.listWorkflows()
	if newActive == 0 {
		h.ui.state.SetWorkflowListIndex(0, len(summaries))
	} else {
		for i, summary := range summaries {
			if summary.ID == newActive {
				h.ui.state.SetWorkflowListIndex(i, len(summaries))
				break
			}
		}
	}
	h.ui.updateWorkflowPointer()
	h.ui.write("%s🗑  Deleted workflow #%d%s\n", h.ui.colors.BrightYellow, activeID, h.ui.colors.Reset)
}

// readNextByte reads the next byte from either a buffered reader or stdin
func (h *KeyHandler) readNextByte(reader *bufio.Reader) (byte, error) {
	if reader != nil {
		return reader.ReadByte()
	}
	var buf [1]byte
	_, err := h.ui.stdin.Read(buf[:])
	return buf[0], err
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
		if errors.Is(err, ErrWorkflowCanceled) {
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
func (h *KeyHandler) executeWorkflow(oldState *term.State) {
	if h.ui.workflow == nil {
		h.ui.notifyWorkflowError("No active workflow. Press Ctrl+N to create one.", 3*time.Second)
		return
	}
	if h.ui.workflow.IsEmpty() {
		h.ui.notifyWorkflowError("Workflow is empty. Add some steps first!", 3*time.Second)
		return
	}

	// Restore terminal state before execution
	h.restoreTerminalState(oldState)

	// Clear screen and execute workflow
	clearScreen(h.ui.stdout)

	err := h.ui.ExecuteWorkflow()
	if errors.Is(err, ErrWorkflowCanceled) {
		h.handleSoftCancel(oldState)
		h.reenterRawMode(oldState)
		return
	}
	if err != nil {
		h.ui.notifyWorkflowError(fmt.Sprintf("Workflow execution failed: %v", err), 3*time.Second)
		h.reenterRawMode(oldState)
		return
	}

	h.ui.notifyWorkflowSuccess("Workflow preserved for reuse. Press 'Ctrl+t' to view or modify.", 3*time.Second)
	h.reenterRawMode(oldState)

	// Keep workflow for reuse - don't clear it
}
