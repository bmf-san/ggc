// Package interactive houses interactive UI types and helpers shared across the application.
package interactive

import (
	"bufio"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/width"

	kb "github.com/bmf-san/ggc/v7/internal/keybindings"
	"github.com/bmf-san/ggc/v7/internal/termio"
)

// AddRune adds a UTF-8 rune to the input at cursor position
func (s *UIState) AddRune(r rune) {
	// Switch to input context when user starts typing
	if s.context != kb.ContextInput {
		s.SetContext(kb.ContextInput)
	}

	// Convert current input to runes for proper cursor positioning
	inputRunes := []rune(s.input)
	if s.cursorPos <= len(inputRunes) {
		// Efficiently insert the rune at the cursor position
		newRunes := make([]rune, len(inputRunes)+1)
		copy(newRunes, inputRunes[:s.cursorPos])
		newRunes[s.cursorPos] = r
		copy(newRunes[s.cursorPos+1:], inputRunes[s.cursorPos:])

		s.input = string(newRunes)
		s.cursorPos++
		s.UpdateFiltered()

		// Switch to search context when actively filtering
		if s.input != "" && s.context != kb.ContextSearch {
			s.SetContext(kb.ContextSearch)
		}
	}
}

// RemoveChar removes character before cursor (backspace)
func (s *UIState) RemoveChar() {
	if s.cursorPos > 0 && s.input != "" {
		// Convert to runes for proper UTF-8 handling
		inputRunes := []rune(s.input)
		if s.cursorPos <= len(inputRunes) {
			inputRunes = append(inputRunes[:s.cursorPos-1], inputRunes[s.cursorPos:]...)

			s.input = string(inputRunes)
			s.cursorPos--
			s.UpdateFiltered()
		}
	}
}

// ClearInput clears all input
func (s *UIState) ClearInput() {
	s.input = ""
	s.cursorPos = 0
	s.UpdateFiltered()
}

// DeleteWord deletes word before cursor (Ctrl+W)
func (s *UIState) DeleteWord() {
	if s.cursorPos == 0 {
		return
	}

	// Convert to runes for proper UTF-8 handling
	inputRunes := []rune(s.input)

	// Find start of current word (skip trailing spaces first)
	pos := s.cursorPos - 1
	for pos >= 0 && inputRunes[pos] == ' ' {
		pos--
	}

	// Find start of word
	for pos >= 0 && inputRunes[pos] != ' ' {
		pos--
	}
	pos++ // Move to first character of word

	// Delete from word start to cursor
	inputRunes = append(inputRunes[:pos], inputRunes[s.cursorPos:]...)
	s.input = string(inputRunes)
	s.cursorPos = pos
	s.UpdateFiltered()
}

// DeleteToEnd deletes from cursor to end of line (Ctrl+K)
func (s *UIState) DeleteToEnd() {
	if s.cursorPos < utf8.RuneCountInString(s.input) {
		inputRunes := []rune(s.input)
		s.input = string(inputRunes[:s.cursorPos])
		s.UpdateFiltered()
	}
}

// MoveToBeginning moves cursor to beginning of line (Ctrl+A)
func (s *UIState) MoveToBeginning() {
	s.cursorPos = 0
}

// MoveToEnd moves cursor to end of line (Ctrl+E)
func (s *UIState) MoveToEnd() {
	s.cursorPos = utf8.RuneCountInString(s.input)
}

// MoveLeft moves cursor one rune left
func (s *UIState) MoveLeft() {
	if s.cursorPos > 0 {
		s.cursorPos--
	}
}

// MoveRight moves cursor one rune right
func (s *UIState) MoveRight() {
	if s.cursorPos < utf8.RuneCountInString(s.input) {
		s.cursorPos++
	}
}

// MoveWordLeft moves cursor to the beginning of the previous word
func (s *UIState) MoveWordLeft() {
	if s.cursorPos == 0 {
		return
	}
	runes := []rune(s.input)
	pos := s.cursorPos - 1
	// Skip any spaces to the left
	for pos >= 0 && unicode.IsSpace(runes[pos]) {
		pos--
	}
	// Skip the word characters to the left
	for pos >= 0 && !unicode.IsSpace(runes[pos]) {
		pos--
	}
	s.cursorPos = pos + 1
}

// MoveWordRight moves cursor to the start of the next word
func (s *UIState) MoveWordRight() {
	runes := []rune(s.input)
	n := len(runes)
	pos := s.cursorPos
	if pos >= n {
		return
	}
	// Skip current word characters
	for pos < n && !unicode.IsSpace(runes[pos]) {
		pos++
	}
	// Skip following spaces
	for pos < n && unicode.IsSpace(runes[pos]) {
		pos++
	}
	s.cursorPos = pos
}

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

// Helper methods for realTimeEditor

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

func (e *realTimeEditor) runeWidth(r rune) int { return runeDisplayWidth(r) }

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

func (e *realTimeEditor) moveLeft(cols int) {
	if cols <= 0 {
		return
	}
	e.ui.write("\x1b[%dD", cols)
}

func (e *realTimeEditor) moveRight(cols int) {
	if cols <= 0 {
		return
	}
	e.ui.write("\x1b[%dC", cols)
}

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
