package interactive

import (
	"testing"

	"github.com/bmf-san/ggc/v8/internal/history"
	kb "github.com/bmf-san/ggc/v8/internal/keybindings"
)

func newRecallState() *UIState {
	return &UIState{
		commands: []CommandInfo{},
		filtered: []CommandInfo{},
		context:  kb.ContextInput,
	}
}

func seededEntries(cmds ...string) []history.Entry {
	out := make([]history.Entry, len(cmds))
	for i, c := range cmds {
		out[i] = history.Entry{Command: c, Raw: c}
	}
	return out
}

func TestBeginHistoryRecall_SnapshotsDraftAndCursor(t *testing.T) {
	s := newRecallState()
	s.input = "wip"
	s.cursorPos = 3

	s.BeginHistoryRecall(seededEntries("status", "commit", "push"))
	if !s.HistoryRecallActive() {
		t.Fatal("recall should be active after BeginHistoryRecall")
	}
	if s.historyDraft != "wip" {
		t.Errorf("draft = %q, want %q", s.historyDraft, "wip")
	}
	if s.historyDraftCursor != 3 {
		t.Errorf("draftCursor = %d, want 3", s.historyDraftCursor)
	}
	if s.historyCursor != 3 {
		t.Errorf("initial cursor should be len(entries)=3, got %d", s.historyCursor)
	}
}

func TestStepHistoryPrev_WalksOlderThenStops(t *testing.T) {
	s := newRecallState()
	s.BeginHistoryRecall(seededEntries("status", "commit", "push"))

	// First Prev → newest entry (push).
	if !s.StepHistoryPrev() || s.input != "push" {
		t.Fatalf("first Prev should set input=push, got %q", s.input)
	}
	// commit
	if !s.StepHistoryPrev() || s.input != "commit" {
		t.Fatalf("second Prev should set input=commit, got %q", s.input)
	}
	// status (oldest)
	if !s.StepHistoryPrev() || s.input != "status" {
		t.Fatalf("third Prev should set input=status, got %q", s.input)
	}
	// Pinned at oldest.
	if s.StepHistoryPrev() {
		t.Fatal("fourth Prev past oldest should be a no-op")
	}
	if s.input != "status" {
		t.Errorf("input should still be status, got %q", s.input)
	}
}

func TestStepHistoryNext_RestoresDraftAtTop(t *testing.T) {
	s := newRecallState()
	s.input = "wip-draft"
	s.cursorPos = 9
	s.BeginHistoryRecall(seededEntries("status", "commit", "push"))

	s.StepHistoryPrev() // push
	s.StepHistoryPrev() // commit

	// Next → push (newer)
	if !s.StepHistoryNext() || s.input != "push" {
		t.Fatalf("Next should walk back to push, got %q", s.input)
	}
	// Next past newest → draft restored, recall ends.
	if !s.StepHistoryNext() {
		t.Fatal("Next past newest should still claim the keystroke")
	}
	if s.HistoryRecallActive() {
		t.Error("recall should end after restoring draft")
	}
	if s.input != "wip-draft" {
		t.Errorf("draft restored as %q, want %q", s.input, "wip-draft")
	}
	if s.cursorPos != 9 {
		t.Errorf("draft cursor restored as %d, want 9", s.cursorPos)
	}
}

func TestStepHistoryNext_NoopWhenInactive(t *testing.T) {
	s := newRecallState()
	if s.StepHistoryNext() {
		t.Error("Next on inactive recall should return false")
	}
}

func TestStepHistoryPrev_NoopWhenInactive(t *testing.T) {
	s := newRecallState()
	if s.StepHistoryPrev() {
		t.Error("Prev on inactive recall should return false")
	}
}

func TestResetHistoryRecall_FiredByUserMutators(t *testing.T) {
	cases := map[string]func(*UIState){
		"AddRune":     func(s *UIState) { s.AddRune('x') },
		"RemoveChar":  func(s *UIState) { s.input = "abc"; s.cursorPos = 3; s.RemoveChar() },
		"ClearInput":  func(s *UIState) { s.ClearInput() },
		"DeleteWord":  func(s *UIState) { s.input = "ab"; s.cursorPos = 2; s.DeleteWord() },
		"DeleteToEnd": func(s *UIState) { s.input = "abcd"; s.cursorPos = 1; s.DeleteToEnd() },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			s := newRecallState()
			s.BeginHistoryRecall(seededEntries("status"))
			s.StepHistoryPrev()
			mutate(s)
			if s.HistoryRecallActive() {
				t.Errorf("%s should reset recall state", name)
			}
		})
	}
}

func TestHandleHistoryRecallKeys_LazyLoadsEntries(t *testing.T) {
	calls := 0
	withHistoryReader(t, historyReaderFunc(func(n int) ([]history.Entry, error) {
		calls++
		return seededEntries("status", "push"), nil
	}))

	h, _, _ := newSelectorHandler("")
	// Walk: Prev (load + step to newest), Prev (step older), Prev (pinned).
	h.invokeHistoryPrev()
	h.invokeHistoryPrev()
	h.invokeHistoryPrev()

	if calls != 1 {
		t.Errorf("entries should be loaded exactly once per recall session, got %d calls", calls)
	}
	if h.ui.state.input != "status" {
		t.Errorf("input = %q, want status (oldest)", h.ui.state.input)
	}
}

func TestHandleHistoryRecallKeys_EmptyHistoryIsSilent(t *testing.T) {
	withHistoryReader(t, fakeHistoryReader{})

	h, stdout, stderr := newSelectorHandler("")
	h.invokeHistoryPrev()

	if h.ui.state.HistoryRecallActive() {
		t.Error("Prev on empty history should not activate recall")
	}
	if stdout.Len() != 0 || stderr.Len() != 0 {
		t.Errorf("Prev on empty history should be silent, got stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
}

func TestHandleHistoryRecallKeys_NextWithoutActiveIsNoop(t *testing.T) {
	h, _, _ := newSelectorHandler("")
	h.invokeHistoryNext()
	if h.ui.state.HistoryRecallActive() {
		t.Error("Next without prior Prev should not activate recall")
	}
}

func TestClampCursor(t *testing.T) {
	cases := []struct {
		pos  int
		s    string
		want int
	}{
		{-1, "abc", 0},
		{0, "abc", 0},
		{2, "abc", 2},
		{3, "abc", 3},
		{99, "abc", 3},
		{0, "あい", 0},
		{5, "あい", 2}, // rune-aware
	}
	for _, c := range cases {
		if got := clampCursor(c.pos, c.s); got != c.want {
			t.Errorf("clampCursor(%d, %q) = %d, want %d", c.pos, c.s, got, c.want)
		}
	}
}
