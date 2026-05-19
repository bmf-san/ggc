package interactive

import (
	"unicode/utf8"

	"github.com/bmf-san/ggc/v8/internal/history"
	kb "github.com/bmf-san/ggc/v8/internal/keybindings"
)

// historyRecallLimit caps how many entries Ctrl+P / Ctrl+N can walk
// back through. We use the same ceiling as the picker so behavior is
// predictable across both UX surfaces.
const historyRecallLimit = historySelectorLimit

// BeginHistoryRecall snapshots the current input as the "draft" the
// user will return to via Ctrl+N past newest, and seeds the recall
// state with the supplied entries (oldest-first, the order the history
// store returns). historyCursor is parked at len(entries) which is the
// "above newest" sentinel — the first Ctrl+P then lands on the newest
// entry (entries[len-1]).
func (s *UIState) BeginHistoryRecall(entries []history.Entry) {
	s.historyRecallActive = true
	s.historyDraft = s.input
	s.historyDraftCursor = s.cursorPos
	s.historyEntries = make([]historyRecallEntry, len(entries))
	for i := range entries {
		s.historyEntries[i] = historyRecallEntry{display: entries[i].Display()}
	}
	s.historyCursor = len(entries) // before first Prev: above newest
}

// resetHistoryRecall clears every field touched by BeginHistoryRecall.
// Called from user-driven mutators in input_state.go so that any direct
// edit (typing, backspace, Ctrl+U/W/K) exits recall cleanly. Safe to
// call when recall is already inactive.
func (s *UIState) resetHistoryRecall() {
	if !s.historyRecallActive {
		return
	}
	s.historyRecallActive = false
	s.historyDraft = ""
	s.historyDraftCursor = 0
	s.historyCursor = 0
	s.historyEntries = nil
}

// HistoryRecallActive reports whether the user is currently walking
// the persisted history via Ctrl+P / Ctrl+N. Exposed for tests and the
// renderer (which may want to hint at the state in a future change).
func (s *UIState) HistoryRecallActive() bool {
	return s.historyRecallActive
}

// StepHistoryPrev moves one entry older. Returns false when recall is
// not yet active (the caller is expected to seed it with entries
// first) or when the cursor is already pinned at the oldest entry.
func (s *UIState) StepHistoryPrev() bool {
	if !s.historyRecallActive || len(s.historyEntries) == 0 {
		return false
	}
	if s.historyCursor <= 0 {
		return false
	}
	s.historyCursor--
	s.setRecallInput(s.historyEntries[s.historyCursor].display)
	return true
}

// StepHistoryNext moves one entry newer. When the cursor walks past
// the newest entry, the original draft is restored and recall exits.
// Returns false if recall is not active.
func (s *UIState) StepHistoryNext() bool {
	if !s.historyRecallActive {
		return false
	}
	if s.historyCursor >= len(s.historyEntries)-1 {
		// Past the newest entry → restore the user's draft and end
		// recall, mirroring how readline / fish behave at the bottom
		// of the history stack.
		draft := s.historyDraft
		draftCursor := s.historyDraftCursor
		s.resetHistoryRecall()
		s.input = draft
		s.cursorPos = clampCursor(draftCursor, draft)
		s.UpdateFiltered()
		return true
	}
	s.historyCursor++
	s.setRecallInput(s.historyEntries[s.historyCursor].display)
	return true
}

// setRecallInput replaces the input buffer with a recalled entry,
// leaving the cursor at the end. It deliberately does NOT call
// resetHistoryRecall, since the whole point is that recall stays
// active until the user types or cancels.
func (s *UIState) setRecallInput(text string) {
	s.input = text
	s.cursorPos = utf8.RuneCountInString(text)
	s.UpdateFiltered()
}

// clampCursor keeps a cursor position inside the rune length of s.
// Used when restoring the draft so a stale cursor (e.g. saved from a
// longer draft that the user shortened mid-recall, which currently is
// not possible but is cheap insurance) never points off the end.
func clampCursor(pos int, s string) int {
	max := utf8.RuneCountInString(s)
	if pos < 0 {
		return 0
	}
	if pos > max {
		return max
	}
	return pos
}

// handleHistoryRecallKeys dispatches Ctrl+P / Ctrl+N to the recall
// state machine. Returns true when the chord was claimed so the
// caller stops walking the rest of the Ctrl-key handlers.
//
// First Ctrl+P seeds recall by snapshotting the entries via
// defaultHistoryReader and saving the current input as the draft. All
// subsequent Prev/Next presses just walk the cached snapshot, so the
// disk is touched at most once per recall session.
func (h *KeyHandler) handleHistoryRecallKeys(km *kb.KeyBindingMap, stroke kb.KeyStroke) bool {
	switch {
	case km.MatchesKeyStroke("history_prev", stroke):
		h.invokeHistoryPrev()
		return true
	case km.MatchesKeyStroke("history_next", stroke):
		h.invokeHistoryNext()
		return true
	}
	return false
}

func (h *KeyHandler) invokeHistoryPrev() {
	if !h.ui.state.HistoryRecallActive() {
		entries, err := defaultHistoryReader.ReadLast(historyRecallLimit)
		if err != nil || len(entries) == 0 {
			// Nothing to recall — stay silent so Ctrl+P on an
			// empty history doesn't spam the UI.
			return
		}
		h.ui.state.BeginHistoryRecall(entries)
	}
	h.ui.state.StepHistoryPrev()
}

func (h *KeyHandler) invokeHistoryNext() {
	// Ctrl+N is only meaningful while recall is already active; pressing
	// it from a fresh prompt would have nothing to step toward.
	if !h.ui.state.HistoryRecallActive() {
		return
	}
	h.ui.state.StepHistoryNext()
}
