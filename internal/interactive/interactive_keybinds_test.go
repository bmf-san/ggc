package interactive

import (
	"testing"

	kb "github.com/bmf-san/ggc/v7/internal/keybindings"
)

func TestBuildSearchKeybindEntriesUsesConfiguredBindings(t *testing.T) {
	renderer := &Renderer{}

	ui := newUIWithKeyMap(&kb.KeyBindingMap{
		ClearLine:       []kb.KeyStroke{kb.NewCtrlKeyStroke('l')},
		DeleteWord:      []kb.KeyStroke{kb.NewAltKeyStroke(0, "backspace")},
		DeleteToEnd:     []kb.KeyStroke{kb.NewCtrlKeyStroke('k')},
		MoveToBeginning: []kb.KeyStroke{kb.NewCtrlKeyStroke('a')},
		MoveToEnd:       []kb.KeyStroke{kb.NewCtrlKeyStroke('e')},
	})

	entries := renderer.buildSearchKeybindEntries(ui)

	if entry, ok := findEntry(entries, "Clear all input"); !ok || entry.key != "Ctrl+l" {
		t.Fatalf("expected clear line key to be Ctrl+l, got %+v", entry)
	}

	if entry, ok := findEntry(entries, "Delete word"); !ok || entry.key != "Alt+Backspace" {
		t.Fatalf("expected delete word key to be Alt+Backspace, got %+v", entry)
	}
}

func TestBuildSearchKeybindEntriesFallsBackToDefaults(t *testing.T) {
	renderer := &Renderer{}

	entries := renderer.buildSearchKeybindEntries(&UI{})

	if entry, ok := findEntry(entries, "Delete word"); !ok || entry.key != "Ctrl+w" {
		t.Fatalf("expected default delete word to be Ctrl+w, got %+v", entry)
	}

	if _, ok := findEntry(entries, "Delete character"); !ok {
		t.Fatalf("expected static Backspace hint to be present")
	}
}

func TestBuildSearchKeybindEntriesFormatsMultipleKeys(t *testing.T) {
	renderer := &Renderer{}

	ui := newUIWithKeyMap(&kb.KeyBindingMap{
		DeleteWord: []kb.KeyStroke{
			kb.NewCtrlKeyStroke('w'),
			kb.NewAltKeyStroke(0, "backspace"),
		},
	})

	entries := renderer.buildSearchKeybindEntries(ui)

	if entry, ok := findEntry(entries, "Delete word"); !ok || entry.key != "Ctrl+w, Alt+Backspace" {
		t.Fatalf("expected combined delete word bindings, got %+v", entry)
	}
}

func newUIWithKeyMap(km *kb.KeyBindingMap) *UI {
	state := &UIState{context: kb.ContextSearch}
	ui := &UI{state: state}
	contextMap := kb.NewContextualKeyBindingMap(kb.ProfileDefault, "darwin", "iterm")
	contextMap.SetContext(kb.ContextSearch, km)
	handler := &KeyHandler{contextualMap: contextMap}
	handler.ui = ui
	ui.handler = handler
	return ui
}

func findEntry(entries []keybindHelpEntry, desc string) (keybindHelpEntry, bool) {
	for _, entry := range entries {
		if entry.desc == desc {
			return entry, true
		}
	}
	return keybindHelpEntry{}, false
}
