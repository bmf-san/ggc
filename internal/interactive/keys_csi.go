package interactive

import (
	"bufio"

	kb "github.com/bmf-san/ggc/v7/internal/keybindings"
)

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
