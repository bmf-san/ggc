package interactive

import (
	"bufio"

	"golang.org/x/term"

	kb "github.com/bmf-san/ggc/v7/internal/keybindings"
)

func (h *KeyHandler) handleControlChar(b byte, oldState *term.State, reader *bufio.Reader) (bool, bool, []string) {
	// Handle Ctrl+letter keys (ASCII 1-26)
	if b >= 1 && b <= 26 {
		handled, cont, result := h.handleCtrlLetterKeys(b, oldState)
		if handled {
			return handled, cont, result
		}
	}

	// Handle special control characters
	return h.handleSpecialCtrlChars(b, oldState, reader)
}

// handleCtrlLetterKeys processes Ctrl+A through Ctrl+Z
func (h *KeyHandler) handleCtrlLetterKeys(b byte, oldState *term.State) (bool, bool, []string) {
	km := h.GetCurrentKeyMap()
	ctrlStroke := kb.NewCtrlKeyStroke(rune('a' + b - 1))

	// Workflow mode has different key handling
	if h.ui.state.IsWorkflowMode() {
		if h.handleWorkflowCtrlKeys(km, ctrlStroke, b, oldState) {
			return true, true, nil
		}
		// Ctrl+C needs to fall through to special handler
		if b != 3 {
			return false, true, nil
		}
		return false, false, nil
	}

	// Search mode: full key handling
	return h.handleSearchCtrlKeys(km, ctrlStroke, oldState)
}

// handleWorkflowCtrlKeys handles Ctrl+letter in workflow mode
func (h *KeyHandler) handleWorkflowCtrlKeys(km *kb.KeyBindingMap, stroke kb.KeyStroke, _ byte, oldState *term.State) bool {
	switch {
	case km.MatchesKeyStroke("workflow_delete", stroke):
		h.deleteActiveWorkflow()
		return true
	case km.MatchesKeyStroke("move_down", stroke):
		h.moveWorkflowList(1)
		return true
	case km.MatchesKeyStroke("move_up", stroke):
		h.moveWorkflowList(-1)
		return true
	case km.MatchesKeyStroke("toggle_workflow_view", stroke):
		h.ui.ToggleWorkflowView()
		return true
	case km.MatchesKeyStroke("soft_cancel", stroke):
		h.handleSoftCancel(oldState)
		return true
	}
	return false
}

// handleSearchCtrlKeys handles Ctrl+letter in search mode
func (h *KeyHandler) handleSearchCtrlKeys(km *kb.KeyBindingMap, stroke kb.KeyStroke, oldState *term.State) (bool, bool, []string) {
	// Navigation keys
	if h.handleSearchNavKeys(km, stroke) {
		return true, true, nil
	}

	// Editing keys
	if h.handleSearchEditKeys(km, stroke) {
		return true, true, nil
	}

	// Mode toggle and cancel
	if h.handleSearchModeKeys(km, stroke, oldState) {
		return true, true, nil
	}

	return false, true, nil
}

// handleSearchNavKeys handles navigation Ctrl+keys in search mode
func (h *KeyHandler) handleSearchNavKeys(km *kb.KeyBindingMap, stroke kb.KeyStroke) bool {
	switch {
	case km.MatchesKeyStroke("move_up", stroke):
		h.handleMoveUp()
		return true
	case km.MatchesKeyStroke("move_down", stroke):
		h.handleMoveDown()
		return true
	case km.MatchesKeyStroke("move_to_beginning", stroke):
		h.ui.state.MoveToBeginning()
		return true
	case km.MatchesKeyStroke("move_to_end", stroke):
		h.ui.state.MoveToEnd()
		return true
	}
	return false
}

// handleSearchEditKeys handles editing Ctrl+keys in search mode
func (h *KeyHandler) handleSearchEditKeys(km *kb.KeyBindingMap, stroke kb.KeyStroke) bool {
	switch {
	case km.MatchesKeyStroke("clear_line", stroke):
		h.ui.state.ClearInput()
		return true
	case km.MatchesKeyStroke("delete_word", stroke):
		h.ui.state.DeleteWord()
		return true
	case km.MatchesKeyStroke("delete_to_end", stroke):
		h.ui.state.DeleteToEnd()
		return true
	}
	return false
}

// handleSearchModeKeys handles mode-related Ctrl+keys in search mode
func (h *KeyHandler) handleSearchModeKeys(km *kb.KeyBindingMap, stroke kb.KeyStroke, oldState *term.State) bool {
	switch {
	case km.MatchesKeyStroke("toggle_workflow_view", stroke) && h.ui.state.input == "":
		h.ui.ToggleWorkflowView()
		return true
	case km.MatchesKeyStroke("soft_cancel", stroke):
		h.handleSoftCancel(oldState)
		return true
	}
	return false
}

// handleSpecialCtrlChars handles non-letter control characters
func (h *KeyHandler) handleSpecialCtrlChars(b byte, oldState *term.State, reader *bufio.Reader) (bool, bool, []string) {
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
	case 27: // ESC
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
