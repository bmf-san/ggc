package interactive

import (
	"bufio"

	"golang.org/x/term"

	kb "github.com/bmf-san/ggc/v7/internal/keybindings"
)

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
