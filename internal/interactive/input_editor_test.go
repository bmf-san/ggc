package interactive

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

// makeEditor creates a realTimeEditor for unit testing.
// Returns the editor, pointer to the rune slice, and pointer to the cursor.
func makeEditor(runes []rune, cursorPos int) (*realTimeEditor, *[]rune, *int) {
	r := make([]rune, len(runes))
	copy(r, runes)
	c := cursorPos
	var out strings.Builder
	ui := &UI{colors: NewANSIColors(), stdout: &out}
	e := &realTimeEditor{ui: ui, inputRunes: &r, cursor: &c}
	return e, &r, &c
}

// --- runeDisplayWidth (missing East Asian wide branch) ---

func TestRuneDisplayWidth_EastAsianWide(t *testing.T) {
	// '中' is U+4E2D, an East Asian Wide character
	got := runeDisplayWidth('中')
	if got != 2 {
		t.Errorf("'中' width: got %d, want 2", got)
	}
}

func TestRuneDisplayWidth_ZeroWidth(t *testing.T) {
	// U+0301 combining acute accent is zero width
	got := runeDisplayWidth(0x0301)
	if got != 0 {
		t.Errorf("combining mark width: got %d, want 0", got)
	}
}

func TestRuneDisplayWidth_Normal(t *testing.T) {
	got := runeDisplayWidth('A')
	if got != 1 {
		t.Errorf("'A' width: got %d, want 1", got)
	}
}

// --- colsBetween edge cases (currently 70%) ---

func TestColsBetween_NegativeFrom(t *testing.T) {
	e, _, _ := makeEditor([]rune{'a', 'b', 'c'}, 0)
	got := e.colsBetween(-1, 2)
	if got != 2 {
		t.Errorf("colsBetween(-1, 2): got %d, want 2", got)
	}
}

func TestColsBetween_NegativeTo(t *testing.T) {
	e, _, _ := makeEditor([]rune{'a', 'b', 'c'}, 0)
	// to=-1 is clamped to 0; from=1 > to=0 so they swap → count rune[0]='a' = 1
	got := e.colsBetween(1, -1)
	if got != 1 {
		t.Errorf("colsBetween(1, -1): got %d, want 1", got)
	}
}

func TestColsBetween_Reversed(t *testing.T) {
	e, _, _ := makeEditor([]rune{'a', 'b', 'c'}, 0)
	// from > to: should swap and compute 2 columns
	got := e.colsBetween(2, 0)
	if got != 2 {
		t.Errorf("colsBetween(2, 0): got %d, want 2", got)
	}
}

func TestColsBetween_Empty(t *testing.T) {
	e, _, _ := makeEditor([]rune{}, 0)
	got := e.colsBetween(0, 0)
	if got != 0 {
		t.Errorf("colsBetween(0, 0) on empty: got %d, want 0", got)
	}
}

// --- findGraphemeStart ---

func TestFindGraphemeStart_Simple(t *testing.T) {
	// Simple ASCII: each rune is its own grapheme
	e, _, _ := makeEditor([]rune{'h', 'e', 'l', 'l', 'o'}, 4)
	got := e.findGraphemeStart(4)
	if got != 4 {
		t.Errorf("ASCII findGraphemeStart(4): got %d, want 4", got)
	}
}

func TestFindGraphemeStart_WithCombiningMark(t *testing.T) {
	// 'e' + combining acute (0x0301): grapheme starts at 'e' (index 0)
	e, _, _ := makeEditor([]rune{'e', 0x0301}, 1)
	got := e.findGraphemeStart(1)
	if got != 0 {
		t.Errorf("combining mark: findGraphemeStart(1) = %d, want 0", got)
	}
}

func TestFindGraphemeStart_ZeroBounds(t *testing.T) {
	// Position 0: should return 0 even with edge cases
	e, _, _ := makeEditor([]rune{'a'}, 0)
	got := e.findGraphemeStart(0)
	if got != 0 {
		t.Errorf("findGraphemeStart(0) = %d, want 0", got)
	}
}

// --- skipCombiningMarks ---

func TestSkipCombiningMarks_NoCombining(t *testing.T) {
	e, _, _ := makeEditor([]rune{'a', 'b', 'c'}, 0)
	got := e.skipCombiningMarks(2)
	if got != 2 {
		t.Errorf("no combining marks: got %d, want 2", got)
	}
}

func TestSkipCombiningMarks_WithCombining(t *testing.T) {
	// 'e' at 0, combining acute at 1, combining grave at 2
	e, _, _ := makeEditor([]rune{'e', 0x0301, 0x0300}, 0)
	start := e.skipCombiningMarks(2)
	if start != 0 {
		t.Errorf("skip two combining marks from 2: got %d, want 0", start)
	}
}

func TestSkipCombiningMarks_NegativeStart(t *testing.T) {
	e, _, _ := makeEditor([]rune{'a'}, 0)
	got := e.skipCombiningMarks(-1)
	if got != -1 {
		t.Errorf("negative start: got %d, want -1", got)
	}
}

// --- handleRegionalIndicators ---

func TestHandleRegionalIndicators_Pair(t *testing.T) {
	// Two regional indicators: 🇺 (U+1F1FA) 🇸 (U+1F1F8) — "US" flag
	ri1 := rune(0x1F1FA)
	ri2 := rune(0x1F1F8)
	e, _, _ := makeEditor([]rune{ri1, ri2}, 0)
	// Pointing at index 1 (second RI): should back up to index 0
	got := e.handleRegionalIndicators(1)
	if got != 0 {
		t.Errorf("RI pair from 1: got %d, want 0", got)
	}
}

func TestHandleRegionalIndicators_Single(t *testing.T) {
	// Only one RI (no preceding RI): stays at the same position
	ri := rune(0x1F1FA)
	e, _, _ := makeEditor([]rune{'x', ri}, 0)
	got := e.handleRegionalIndicators(1)
	if got != 1 {
		t.Errorf("single RI at 1: got %d, want 1", got)
	}
}

func TestHandleRegionalIndicators_NonRI(t *testing.T) {
	e, _, _ := makeEditor([]rune{'a', 'b'}, 0)
	got := e.handleRegionalIndicators(1)
	if got != 1 {
		t.Errorf("non-RI: got %d, want 1", got)
	}
}

// --- handleZWJSequences ---

func TestHandleZWJSequences_NoZWJ(t *testing.T) {
	e, _, _ := makeEditor([]rune{'a', 'b', 'c'}, 0)
	got := e.handleZWJSequences(2)
	if got != 2 {
		t.Errorf("no ZWJ: got %d, want 2", got)
	}
}

func TestHandleZWJSequences_WithZWJ(t *testing.T) {
	// woman (👩=U+1F469), ZWJ (U+200D), rocket (🚀=U+1F680)
	// positions: [0]=👩, [1]=ZWJ, [2]=🚀
	// pointing at 2 (🚀): ZWJ is at 1, so skip ZWJ+base → position 0
	woman := rune(0x1F469)
	zwj := rune(0x200D)
	rocket := rune(0x1F680)
	e, _, _ := makeEditor([]rune{woman, zwj, rocket}, 0)
	got := e.handleZWJSequences(2)
	if got != 0 {
		t.Errorf("ZWJ sequence: got %d, want 0", got)
	}
}

// --- moveWordLeft ---

func TestMoveWordLeft_AtStart(t *testing.T) {
	e, _, c := makeEditor([]rune{'h', 'i'}, 0)
	e.moveWordLeft()
	if *c != 0 {
		t.Errorf("moveWordLeft at start: cursor should stay 0, got %d", *c)
	}
}

func TestMoveWordLeft_SimpleWord(t *testing.T) {
	// "hello world", cursor at end (11)
	runes := []rune("hello world")
	e, _, c := makeEditor(runes, len(runes))
	e.moveWordLeft()
	// Should stop at start of "world" = index 6
	if *c != 6 {
		t.Errorf("moveWordLeft: cursor = %d, want 6", *c)
	}
}

func TestMoveWordLeft_SkipsTrailingSpace(t *testing.T) {
	// "foo   ", cursor at 6
	runes := []rune("foo   ")
	e, _, c := makeEditor(runes, len(runes))
	e.moveWordLeft()
	// Skip trailing spaces, then back over "foo" → cursor at 0
	if *c != 0 {
		t.Errorf("moveWordLeft with trailing spaces: cursor = %d, want 0", *c)
	}
}

// --- deleteWordLeft ---

func TestDeleteWordLeft_AtStart(t *testing.T) {
	e, r, c := makeEditor([]rune{'a', 'b'}, 0)
	e.deleteWordLeft()
	// Nothing should change
	if *c != 0 {
		t.Errorf("deleteWordLeft at start: cursor = %d, want 0", *c)
	}
	if string(*r) != "ab" {
		t.Errorf("deleteWordLeft at start: runes = %q, want %q", string(*r), "ab")
	}
}

func TestDeleteWordLeft_DeletesWord(t *testing.T) {
	// "hello world", cursor at 11 (end)
	runes := []rune("hello world")
	e, r, c := makeEditor(runes, len(runes))
	e.deleteWordLeft()
	// Should delete "world" (positions 6-10), leaving "hello " with cursor at 6
	if *c != 6 {
		t.Errorf("deleteWordLeft: cursor = %d, want 6", *c)
	}
	if string(*r) != "hello " {
		t.Errorf("deleteWordLeft: runes = %q, want %q", string(*r), "hello ")
	}
}

// --- moveLeft / moveRight zero and negative ---

func TestMoveLeft_ZeroDoesNothing(t *testing.T) {
	var out strings.Builder
	ui := &UI{colors: NewANSIColors(), stdout: &out}
	runes := []rune("ab")
	cursor := 1
	e := &realTimeEditor{ui: ui, inputRunes: &runes, cursor: &cursor}
	e.moveLeft(0)
	if out.Len() != 0 {
		t.Errorf("moveLeft(0) should write nothing, got %q", out.String())
	}
}

func TestMoveLeft_NegativeDoesNothing(t *testing.T) {
	var out strings.Builder
	ui := &UI{colors: NewANSIColors(), stdout: &out}
	runes := []rune("ab")
	cursor := 1
	e := &realTimeEditor{ui: ui, inputRunes: &runes, cursor: &cursor}
	e.moveLeft(-3)
	if out.Len() != 0 {
		t.Errorf("moveLeft(-3) should write nothing, got %q", out.String())
	}
}

func TestMoveRight_ZeroDoesNothing(t *testing.T) {
	var out strings.Builder
	ui := &UI{colors: NewANSIColors(), stdout: &out}
	runes := []rune("ab")
	cursor := 0
	e := &realTimeEditor{ui: ui, inputRunes: &runes, cursor: &cursor}
	e.moveRight(0)
	if out.Len() != 0 {
		t.Errorf("moveRight(0) should write nothing, got %q", out.String())
	}
}

func TestMoveRight_NegativeDoesNothing(t *testing.T) {
	var out strings.Builder
	ui := &UI{colors: NewANSIColors(), stdout: &out}
	runes := []rune("ab")
	cursor := 0
	e := &realTimeEditor{ui: ui, inputRunes: &runes, cursor: &cursor}
	e.moveRight(-1)
	if out.Len() != 0 {
		t.Errorf("moveRight(-1) should write nothing, got %q", out.String())
	}
}

// --- printTailAndReposition edge cases ---

func TestPrintTailAndReposition_EmptyTail(t *testing.T) {
	// from == len(runes): no tail, no spaces → nothing meaningful written
	var out strings.Builder
	ui := &UI{colors: NewANSIColors(), stdout: &out}
	runes := []rune("ab")
	cursor := 2
	e := &realTimeEditor{ui: ui, inputRunes: &runes, cursor: &cursor}
	e.printTailAndReposition(2, 0) // from == len(runes), clearedCols == 0
	// moveLeft(0) is a no-op, nothing should be written
	if out.Len() != 0 {
		t.Errorf("printTailAndReposition with empty tail wrote %q, want empty", out.String())
	}
}

func TestPrintTailAndReposition_NoClearedCols(t *testing.T) {
	// from < len: tail is printed, no extra spaces
	var out strings.Builder
	ui := &UI{colors: NewANSIColors(), stdout: &out}
	runes := []rune("ab")
	cursor := 0
	e := &realTimeEditor{ui: ui, inputRunes: &runes, cursor: &cursor}
	e.printTailAndReposition(0, 0) // clearedCols == 0 → no spaces appended
	written := out.String()
	if !strings.Contains(written, "ab") {
		t.Errorf("printTailAndReposition: expected tail 'ab' in output, got %q", written)
	}
}

// --- handleEscape uncovered branches (b, f, 127/\b) ---

func TestHandleEscape_MoveWordLeft(t *testing.T) {
	// ESC b → moveWordLeft
	runes := []rune("hello world")
	e, _, c := makeEditor(runes, len(runes))
	reader := bufio.NewReader(bytes.NewBufferString("b"))
	e.handleEscape(reader)
	if *c != 6 {
		t.Errorf("ESC b: cursor = %d, want 6", *c)
	}
}

func TestHandleEscape_MoveWordRight(t *testing.T) {
	// ESC f → moveWordRight: skips "hello" then the space → cursor at 6
	runes := []rune("hello world")
	e, _, c := makeEditor(runes, 0)
	reader := bufio.NewReader(bytes.NewBufferString("f"))
	e.handleEscape(reader)
	if *c != 6 {
		t.Errorf("ESC f: cursor = %d, want 6", *c)
	}
}

func TestHandleEscape_DeleteWordLeft(t *testing.T) {
	// ESC 127 → deleteWordLeft (Option+Backspace)
	runes := []rune("hello world")
	e, r, c := makeEditor(runes, len(runes))
	reader := bufio.NewReader(bytes.NewReader([]byte{127}))
	e.handleEscape(reader)
	if *c != 6 {
		t.Errorf("ESC 127: cursor = %d, want 6", *c)
	}
	if string(*r) != "hello " {
		t.Errorf("ESC 127: runes = %q, want 'hello '", string(*r))
	}
}

func TestHandleEscape_EOF(t *testing.T) {
	// Empty reader → ReadByte returns error → returns without panic
	runes := []rune("ab")
	e, _, c := makeEditor(runes, 1)
	reader := bufio.NewReader(bytes.NewBufferString(""))
	e.handleEscape(reader) // must not panic
	if *c != 1 {
		t.Errorf("ESC EOF: cursor changed unexpectedly to %d", *c)
	}
}

// --- handleApplicationEscape uncovered branches ---

func TestHandleApplicationEscape_MoveRight(t *testing.T) {
	// OC → move cursor right
	runes := []rune("ab")
	e, _, c := makeEditor(runes, 0)
	reader := bufio.NewReader(bytes.NewBufferString("C"))
	e.handleApplicationEscape(reader)
	if *c != 1 {
		t.Errorf("app escape C: cursor = %d, want 1", *c)
	}
}

func TestHandleApplicationEscape_MoveRight_AtEnd(t *testing.T) {
	// OC with cursor at end → no movement
	runes := []rune("ab")
	e, _, c := makeEditor(runes, 2)
	reader := bufio.NewReader(bytes.NewBufferString("C"))
	e.handleApplicationEscape(reader)
	if *c != 2 {
		t.Errorf("app escape C at end: cursor = %d, want 2", *c)
	}
}

func TestHandleApplicationEscape_MoveLeft_AtStart(t *testing.T) {
	// OD with cursor at start → no movement
	runes := []rune("ab")
	e, _, c := makeEditor(runes, 0)
	reader := bufio.NewReader(bytes.NewBufferString("D"))
	e.handleApplicationEscape(reader)
	if *c != 0 {
		t.Errorf("app escape D at start: cursor = %d, want 0", *c)
	}
}

func TestHandleApplicationEscape_EOF(t *testing.T) {
	// Empty reader → ReadByte returns error → returns without panic
	runes := []rune("ab")
	e, _, c := makeEditor(runes, 1)
	reader := bufio.NewReader(bytes.NewBufferString(""))
	e.handleApplicationEscape(reader) // must not panic
	if *c != 1 {
		t.Errorf("app escape EOF: cursor changed to %d", *c)
	}
}
