package interactive

import (
	"bufio"

	kb "github.com/bmf-san/ggc/v7/internal/keybindings"
)

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
