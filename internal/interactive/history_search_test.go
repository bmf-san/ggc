package interactive

import (
	"testing"

	"github.com/bmf-san/ggc/v8/internal/history"
	kb "github.com/bmf-san/ggc/v8/internal/keybindings"
)

func newSearchState() *UIState {
	return &UIState{
		commands: []CommandInfo{
			{Command: "status", Description: "show status"},
			{Command: "commit", Description: "commit changes"},
		},
		filtered: []CommandInfo{},
		context:  kb.ContextInput,
	}
}

func TestEnterHistorySearch_SwapsCommandsNewestFirst(t *testing.T) {
	s := newSearchState()
	entries := []history.Entry{
		{Command: "status", Raw: "status"},
		{Command: "commit", Args: []string{"-m", "wip"}, Raw: "commit -m wip"},
		{Command: "push", Raw: "push"},
	}
	s.EnterHistorySearch(entries)

	if !s.IsHistorySearch() {
		t.Fatal("EnterHistorySearch should activate the mode")
	}
	if len(s.commands) != 3 {
		t.Fatalf("commands should hold 3 entries, got %d", len(s.commands))
	}
	if s.commands[0].Command != "push" {
		t.Errorf("newest entry should be first; got %q", s.commands[0].Command)
	}
	if s.commands[2].Command != "status" {
		t.Errorf("oldest entry should be last; got %q", s.commands[2].Command)
	}
	if s.GetCurrentContext() != kb.ContextSearch {
		t.Errorf("context should switch to ContextSearch, got %v", s.GetCurrentContext())
	}
}

func TestEnterHistorySearch_DeduplicatesByDisplay(t *testing.T) {
	s := newSearchState()
	// Same raw line repeated; we want one entry in the picker so the
	// user isn't asked to scroll past N copies of the same command.
	s.EnterHistorySearch([]history.Entry{
		{Command: "status", Raw: "status"},
		{Command: "status", Raw: "status"},
		{Command: "status", Raw: "status"},
	})
	if len(s.commands) != 1 {
		t.Errorf("duplicate displays should collapse; got %d entries", len(s.commands))
	}
}

func TestEnterHistorySearch_SkipsEmptyDisplay(t *testing.T) {
	s := newSearchState()
	// An entry with no Command and no Raw would render as "" — skip
	// rather than poison the picker with a blank row.
	s.EnterHistorySearch([]history.Entry{{}})
	if len(s.commands) != 0 {
		t.Errorf("blank display should be skipped; got %d entries", len(s.commands))
	}
}

func TestEnterHistorySearch_IsIdempotent(t *testing.T) {
	s := newSearchState()
	orig := s.commands
	s.EnterHistorySearch([]history.Entry{{Command: "status", Raw: "status"}})
	// Second call must not stomp the original backup, otherwise the
	// user could never return to the real commands list.
	s.EnterHistorySearch([]history.Entry{{Command: "push", Raw: "push"}})
	s.ExitHistorySearch()
	if len(s.commands) != len(orig) {
		t.Errorf("backup commands not restored; got %d want %d", len(s.commands), len(orig))
	}
}

func TestExitHistorySearch_RestoresCommands(t *testing.T) {
	s := newSearchState()
	origLen := len(s.commands)
	s.EnterHistorySearch([]history.Entry{{Command: "push", Raw: "push"}})
	s.ExitHistorySearch()

	if s.IsHistorySearch() {
		t.Error("Exit should clear active flag")
	}
	if len(s.commands) != origLen {
		t.Errorf("commands not restored; got %d want %d", len(s.commands), origLen)
	}
	if s.input != "" {
		t.Errorf("input should be cleared on exit, got %q", s.input)
	}
}

func TestExitHistorySearch_NoopWhenInactive(t *testing.T) {
	s := newSearchState()
	origCmds := s.commands
	s.ExitHistorySearch()
	if len(s.commands) != len(origCmds) {
		t.Error("ExitHistorySearch on inactive state should not touch commands")
	}
}

func TestHistorySearchEntryFor_LookupByDisplay(t *testing.T) {
	s := newSearchState()
	s.EnterHistorySearch([]history.Entry{
		{Command: "commit", Args: []string{"-m", "msg"}, Raw: "commit -m msg"},
	})
	got, ok := s.HistorySearchEntryFor("commit -m msg")
	if !ok {
		t.Fatal("expected lookup to succeed for known display string")
	}
	if got.Command != "commit" || len(got.Args) != 2 {
		t.Errorf("recovered entry mismatched: %+v", got)
	}
	if _, ok := s.HistorySearchEntryFor("nonexistent"); ok {
		t.Error("lookup for unknown display should miss")
	}
}

func TestHandleHistorySearchTrigger_EntersOnCtrlR(t *testing.T) {
	withHistoryReader(t, fakeHistoryReader{entries: []history.Entry{
		{Command: "status", Raw: "status"},
	}})
	h, _, _ := newSelectorHandler("")

	km := &kb.KeyBindingMap{HistorySearch: []kb.KeyStroke{kb.NewCtrlKeyStroke('r')}}

	claimed := h.handleHistorySearchTrigger(km, kb.NewCtrlKeyStroke('r'))
	if !claimed {
		t.Fatal("Ctrl+R should be claimed by the trigger")
	}
	if !h.ui.state.IsHistorySearch() {
		t.Error("Ctrl+R should activate history search")
	}
}

func TestHandleHistorySearchTrigger_EmptyHistorySilent(t *testing.T) {
	withHistoryReader(t, fakeHistoryReader{})
	h, stdout, stderr := newSelectorHandler("")

	km := &kb.KeyBindingMap{HistorySearch: []kb.KeyStroke{kb.NewCtrlKeyStroke('r')}}

	// Trigger is still claimed (so move_up/down don't run) but no
	// search is entered when there is nothing to search.
	if !h.handleHistorySearchTrigger(km, kb.NewCtrlKeyStroke('r')) {
		t.Fatal("Ctrl+R on empty history should still claim the chord")
	}
	if h.ui.state.IsHistorySearch() {
		t.Error("empty history must not enter search mode")
	}
	if stdout.Len() != 0 || stderr.Len() != 0 {
		t.Errorf("empty-history Ctrl+R should be silent; stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
}

func TestHandleHistorySearchTrigger_NoMatchReturnsFalse(t *testing.T) {
	h, _, _ := newSelectorHandler("")
	km := &kb.KeyBindingMap{}
	// no history_search binding installed
	if h.handleHistorySearchTrigger(km, kb.NewCtrlKeyStroke('r')) {
		t.Error("unmatched chord should not be claimed")
	}
}

func TestResetToSearchMode_ExitsHistorySearch(t *testing.T) {
	ui := &UI{state: newSearchState()}
	ui.state.EnterHistorySearch([]history.Entry{{Command: "status", Raw: "status"}})
	if !ui.state.IsHistorySearch() {
		t.Fatal("setup precondition failed")
	}
	ui.resetToSearchMode()
	if ui.state.IsHistorySearch() {
		t.Error("soft cancel should drop out of history search")
	}
}

func TestHandleSpecialCtrlChars_CtrlCInHistorySearchCancelsInsteadOfQuitting(t *testing.T) {
	// While reverse-i-search is active, Ctrl+C must back out of the
	// overlay rather than tear down ggc. We assert both the
	// shouldContinue=true return (REPL keeps running) and that the
	// history-search flag is cleared.
	h, _, _ := newSelectorHandler("")
	h.ui.state.EnterHistorySearch([]history.Entry{{Command: "status", Raw: "status"}})
	if !h.ui.state.IsHistorySearch() {
		t.Fatal("setup precondition failed: history search not active")
	}

	handled, shouldContinue, result := h.handleSpecialCtrlChars(3, nil, nil)
	if !handled {
		t.Fatal("Ctrl+C must always be handled")
	}
	if !shouldContinue {
		t.Error("Ctrl+C in history search must keep the REPL alive (shouldContinue=true)")
	}
	if result != nil {
		t.Errorf("Ctrl+C in history search should not dispatch args, got %v", result)
	}
	if h.ui.state.IsHistorySearch() {
		t.Error("Ctrl+C in history search must exit the overlay")
	}
}

func TestHandleSpecialCtrlChars_CtrlCOutsideHistorySearchStillQuits(t *testing.T) {
	// Outside the overlay Ctrl+C keeps its global "quit ggc" meaning;
	// shouldContinue=false signals the Run loop to exit.
	h, _, _ := newSelectorHandler("")

	handled, shouldContinue, result := h.handleSpecialCtrlChars(3, nil, nil)
	if !handled {
		t.Fatal("Ctrl+C must always be handled")
	}
	if shouldContinue {
		t.Error("Ctrl+C outside history search must quit (shouldContinue=false)")
	}
	if result != nil {
		t.Errorf("quit path should not return args, got %v", result)
	}
}
