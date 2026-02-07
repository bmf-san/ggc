// Package interactive houses interactive UI types and helpers shared across the application.
package interactive

import (
	"fmt"
	"io"

	uiutil "github.com/bmf-san/ggc/v7/internal/ui"
)

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

// updateSize updates the terminal dimensions
func (r *Renderer) updateSize() {
	w, h := uiutil.Dimensions(r.writer, 80, 24)
	r.width, r.height = w, h
}

// Render displays the command list with proper terminal handling
func (r *Renderer) Render(ui *UI, state *UIState) {
	clearScreen(r.writer)
	// Disable line wrapping during rendering, restore at end
	uiutil.DisableWrap(r.writer)
	var restoreCursor func()
	defer func() {
		uiutil.EnableWrap(r.writer)
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
	r.renderWorkflowError(ui)
	r.renderWorkflowNotice(ui)

	switch state.mode {
	case ModeWorkflow:
		// Workflow mode: no search prompt, just workflow management
		r.renderWorkflowMode(ui, state)
	default:
		r.renderSearchPrompt(ui, state)
		restoreCursor = r.saveCursorAtSearchPrompt(state)

		switch {
		case state.input == "":
			r.renderEmptyState(ui)
			r.writeEmptyLine()
			r.renderSearchKeybinds(ui)
		case len(state.filtered) == 0:
			r.renderNoMatches(ui, state)
		default:
			r.renderCommandList(ui, state)
		}
	}
}

// clearScreen clears the entire screen and hides cursor
func clearScreen(w io.Writer) {
	uiutil.ClearScreen(w)
	uiutil.HideCursor(w)
}

// showCursor shows the terminal cursor
func showCursor(w io.Writer) {
	uiutil.ShowCursor(w)
}

// ellipsis truncates string and adds ellipsis if it exceeds maxLen (ASCII only)
func ellipsis(s string, maxLen int) string {
	return uiutil.Ellipsis(s, maxLen)
}

func pluralize(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// writeColorln writes a colored line to the terminal.
// The *UI parameter is intentionally unused but kept in the signature
// to stay consistent with other rendering helpers and to allow future
// UI-dependent styling without changing the public API.
func (r *Renderer) writeColorln(_ *UI, text string) {
	// Move to line start, clear line, write content, then CRLF
	_, _ = fmt.Fprint(r.writer, "\r\x1b[K")
	_, _ = fmt.Fprint(r.writer, text+"\r\n")
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
