// Package interactive houses interactive UI types and helpers shared across the application.
package interactive

import (
	"bufio"
	"unicode"

	"golang.org/x/term"

	kb "github.com/bmf-san/ggc/v7/internal/keybindings"
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
