package interactive

import (
	"strings"
	"unicode"

	"golang.org/x/text/width"
)

// isCombining reports whether r is a combining mark (zero display width)
func isCombining(r rune) bool {
	return unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Me, r) || unicode.Is(unicode.Mc, r)
}

// isVariationSelector reports whether r is a variation selector (zero width)
func isVariationSelector(r rune) bool {
	// U+FE00..U+FE0F (VS1..VS16) and U+E0100..U+E01EF (IVS)
	return (r >= 0xFE00 && r <= 0xFE0F) || (r >= 0xE0100 && r <= 0xE01EF)
}

// isRegionalIndicator reports whether r is a regional indicator rune (used for flags)
func isRegionalIndicator(r rune) bool { return r >= 0x1F1E6 && r <= 0x1F1FF }

// isZWJ reports whether r is ZERO WIDTH JOINER
func isZWJ(r rune) bool { return r == 0x200D }

// isEmoji reports common emoji ranges that should render as width 2 on most terminals
func isEmoji(r rune) bool {
	return isEmojiRange1(r) || isEmojiRange2(r)
}

// isEmojiRange1 checks the first set of emoji Unicode ranges
func isEmojiRange1(r rune) bool {
	return (r >= 0x1F300 && r <= 0x1F5FF) || // Misc Symbols and Pictographs
		(r >= 0x1F600 && r <= 0x1F64F) || // Emoticons
		(r >= 0x1F680 && r <= 0x1F6FF) || // Transport and Map Symbols
		(r >= 0x1F700 && r <= 0x1F77F) || // Alchemical Symbols
		(r >= 0x1F780 && r <= 0x1F7FF) // Geometric Shapes Extended
}

// isEmojiRange2 checks the second set of emoji Unicode ranges
func isEmojiRange2(r rune) bool {
	return (r >= 0x1F800 && r <= 0x1F8FF) || // Supplemental Arrows-C
		(r >= 0x1F900 && r <= 0x1F9FF) || // Supplemental Symbols and Pictographs
		(r >= 0x1FA00 && r <= 0x1FAFF) || // Symbols and Pictographs Extended-A
		(r >= 0x2600 && r <= 0x26FF) || // Misc symbols
		(r >= 0x2700 && r <= 0x27BF) // Dingbats
}

// runeDisplayWidth returns the number of terminal columns used by r
func runeDisplayWidth(r rune) int {
	// Zero-width characters
	if isCombining(r) || isVariationSelector(r) || isZWJ(r) {
		return 0
	}
	// East Asian wide/fullwidth
	switch width.LookupRune(r).Kind() {
	case width.EastAsianFullwidth, width.EastAsianWide:
		return 2
	}
	// Common emoji are typically 2 columns
	if isEmoji(r) {
		return 2
	}
	return 1
}

// findGraphemeStart finds the start of the grapheme cluster ending at the given position
func (e *realTimeEditor) findGraphemeStart(pos int) int {
	start := pos
	start = e.skipCombiningMarks(start)
	start = e.handleRegionalIndicators(start)
	start = e.handleZWJSequences(start)
	if start < 0 {
		start = 0
	}
	return start
}

// skipCombiningMarks skips any trailing variation selectors or combining marks
func (e *realTimeEditor) skipCombiningMarks(start int) int {
	for start >= 0 && (isCombining((*e.inputRunes)[start]) || isVariationSelector((*e.inputRunes)[start])) {
		start--
	}
	return start
}

// handleRegionalIndicators handles regional indicator pairs (flags)
func (e *realTimeEditor) handleRegionalIndicators(start int) int {
	if start >= 0 && isRegionalIndicator((*e.inputRunes)[start]) {
		if start > 0 && isRegionalIndicator((*e.inputRunes)[start-1]) {
			start--
		}
	}
	return start
}

// handleZWJSequences handles ZWJ sequences by including the joiner and previous rune repeatedly
func (e *realTimeEditor) handleZWJSequences(start int) int {
	for {
		if start > 0 && isZWJ((*e.inputRunes)[start-1]) {
			// Include ZWJ and the previous rune
			start -= 2
			// Also include any combining marks attached to the new base
			start = e.skipCombiningMarks(start)
			continue
		}
		break
	}
	return start
}

// runeWidth returns the display width of a rune
func (e *realTimeEditor) runeWidth(r rune) int { return runeDisplayWidth(r) }

// colsBetween calculates the number of display columns between two positions
func (e *realTimeEditor) colsBetween(from, to int) int {
	if from < 0 {
		from = 0
	}
	if to < 0 {
		to = 0
	}
	if from > to {
		from, to = to, from
	}
	cols := 0
	for i := from; i < to && i < len(*e.inputRunes); i++ {
		cols += e.runeWidth((*e.inputRunes)[i])
	}
	return cols
}

// moveWordLeft moves the cursor to the beginning of the previous word
func (e *realTimeEditor) moveWordLeft() {
	if *e.cursor == 0 {
		return
	}
	i := *e.cursor - 1
	for i >= 0 && unicode.IsSpace((*e.inputRunes)[i]) {
		i--
	}
	for i >= 0 && !unicode.IsSpace((*e.inputRunes)[i]) {
		i--
	}
	newPos := i + 1
	delta := e.colsBetween(newPos, *e.cursor)
	e.moveLeft(delta)
	*e.cursor = newPos
}

// moveWordRight moves the cursor to the start of the next word
func (e *realTimeEditor) moveWordRight() {
	n := len(*e.inputRunes)
	if *e.cursor >= n {
		return
	}
	i := *e.cursor
	for i < n && !unicode.IsSpace((*e.inputRunes)[i]) {
		i++
	}
	for i < n && unicode.IsSpace((*e.inputRunes)[i]) {
		i++
	}
	delta := e.colsBetween(*e.cursor, i)
	e.moveRight(delta)
	*e.cursor = i
}

// deleteWordLeft deletes the word before the cursor and updates the display
func (e *realTimeEditor) deleteWordLeft() {
	if *e.cursor == 0 {
		return
	}
	// Find new cursor position at the beginning of previous word
	i := *e.cursor - 1
	for i >= 0 && unicode.IsSpace((*e.inputRunes)[i]) {
		i--
	}
	for i >= 0 && !unicode.IsSpace((*e.inputRunes)[i]) {
		i--
	}
	newPos := i + 1
	// Compute columns to move left and columns to clear
	moveCols := e.colsBetween(newPos, *e.cursor)
	clearedCols := 0
	for j := newPos; j < *e.cursor; j++ {
		clearedCols += e.runeWidth((*e.inputRunes)[j])
	}
	// Move cursor left to newPos
	e.moveLeft(moveCols)
	// Delete runes in [newPos, cursor)
	*e.inputRunes = append((*e.inputRunes)[:newPos], (*e.inputRunes)[*e.cursor:]...)
	*e.cursor = newPos
	// Redraw tail and clear leftover cells
	e.printTailAndReposition(*e.cursor, clearedCols)
}

// isWordMotionParam reports whether CSI params include a word-motion modifier
// commonly emitted by terminals (e.g., Ctrl/Alt variants use 5/3/9).
func isWordMotionParam(params string) bool {
	return strings.Contains(params, "5") || strings.Contains(params, "3") || strings.Contains(params, "9")
}
