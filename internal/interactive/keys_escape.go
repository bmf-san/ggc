package interactive

import (
	"bufio"
	"os"

	kb "github.com/bmf-san/ggc/v7/internal/keybindings"
	"github.com/bmf-san/ggc/v7/internal/termio"
	"golang.org/x/term"
)

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
