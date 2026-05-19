package interactive

import (
	"github.com/bmf-san/ggc/v8/internal/history"
	kb "github.com/bmf-san/ggc/v8/internal/keybindings"
)

// historySearchLimit caps how many entries the Ctrl+R reverse-search
// session pulls into the fuzzy filter. Larger windows mostly make the
// match list noisier; if the user wants something older they can fall
// back to `ggc history search <pat>` from the CLI.
const historySearchLimit = 200

// EnterHistorySearch flips the interactive UI into reverse-history-
// search mode. It snapshots the current commands list so the search
// is run over history entries instead, but every existing handler
// (typing, fuzzy filter, up/down nav, Enter) keeps working as-is.
//
// We deliberately reuse ContextSearch and the standard renderer here
// — this is the smallest change that gets the reverse-i-search UX in
// front of users, and the prompt label can be themed later without
// disturbing the dispatcher.
func (s *UIState) EnterHistorySearch(entries []history.Entry) {
	if s.historySearchActive {
		return
	}
	s.historySearchActive = true
	s.historySearchBackup = s.commands
	s.historySearchEntries = make(map[string]history.Entry, len(entries))

	// Build the displayable command list newest-first so the first
	// hit when the user types is intuitively the most recent match.
	commands := make([]CommandInfo, 0, len(entries))
	for i := len(entries) - 1; i >= 0; i-- {
		e := entries[i]
		disp := e.Display()
		if disp == "" {
			continue
		}
		// Map back to the entry by display string. Duplicates collapse
		// to the most recent occurrence, which is what we want when
		// the user just presses Enter on a familiar line.
		if _, dup := s.historySearchEntries[disp]; !dup {
			commands = append(commands, CommandInfo{Command: disp})
			s.historySearchEntries[disp] = e
		}
	}
	s.commands = commands
	s.input = ""
	s.cursorPos = 0
	s.selected = 0
	s.UpdateFiltered()
	s.SetContext(kb.ContextSearch)
}

// ExitHistorySearch restores the pre-search command list and clears
// the input buffer. Safe to call when search isn't active so it can
// be hooked into the global soft-cancel path without a guard.
func (s *UIState) ExitHistorySearch() {
	if !s.historySearchActive {
		return
	}
	s.historySearchActive = false
	s.commands = s.historySearchBackup
	s.historySearchBackup = nil
	s.historySearchEntries = nil
	s.input = ""
	s.cursorPos = 0
	s.selected = 0
	s.UpdateFiltered()
}

// IsHistorySearch reports whether the UI is currently running a Ctrl+R
// reverse-history-search session.
func (s *UIState) IsHistorySearch() bool {
	return s.historySearchActive
}

// HistorySearchEntryFor returns the underlying history entry for a
// recalled display string. Used by handleEnter to dispatch the right
// argv instead of re-tokenizing the displayed text.
func (s *UIState) HistorySearchEntryFor(display string) (history.Entry, bool) {
	e, ok := s.historySearchEntries[display]
	return e, ok
}

// handleHistorySearchTrigger enters reverse-history-search mode on
// Ctrl+R. Called from the Ctrl-key dispatcher. Returns true when the
// chord was claimed.
func (h *KeyHandler) handleHistorySearchTrigger(km *kb.KeyBindingMap, stroke kb.KeyStroke) bool {
	if !km.MatchesKeyStroke("history_search", stroke) {
		return false
	}
	if h.ui.state.IsHistorySearch() {
		// Repeated Ctrl+R inside an active search is a no-op for
		// now. A future change can advance to the next match here,
		// matching readline's reverse-i-search behavior.
		return true
	}
	entries, err := defaultHistoryReader.ReadLast(historySearchLimit)
	if err != nil || len(entries) == 0 {
		return true
	}
	h.ui.state.EnterHistorySearch(entries)
	return true
}
