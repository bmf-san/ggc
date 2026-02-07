package interactive

import (
	"bufio"
	"strings"
	"unicode"
)

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

// handleSoftCancel processes soft cancel (Ctrl+G or ESC)
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

// moveLeft moves the cursor left by the specified number of columns
func (e *realTimeEditor) moveLeft(cols int) {
	if cols <= 0 {
		return
	}
	e.ui.write("\x1b[%dD", cols)
}

// moveRight moves the cursor right by the specified number of columns
func (e *realTimeEditor) moveRight(cols int) {
	if cols <= 0 {
		return
	}
	e.ui.write("\x1b[%dC", cols)
}

// printTailAndReposition prints the tail of the input and repositions the cursor
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
