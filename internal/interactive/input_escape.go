package interactive

import (
	"bufio"
	"os"

	"github.com/bmf-san/ggc/v7/internal/termio"
)

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

// shouldSoftCancelOnEscape checks if ESC key should trigger soft cancel
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
