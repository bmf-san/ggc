package interactive

import "testing"

func TestBuildSearchKeybindEntriesUsesConfiguredBindings(t *testing.T) {
	renderer := &Renderer{}

	ui := newUIWithKeyMap(&KeyBindingMap{
		ClearLine:       []KeyStroke{NewCtrlKeyStroke('l')},
		DeleteWord:      []KeyStroke{NewAltKeyStroke(0, "backspace")},
		DeleteToEnd:     []KeyStroke{NewCtrlKeyStroke('k')},
		MoveToBeginning: []KeyStroke{NewCtrlKeyStroke('a')},
		MoveToEnd:       []KeyStroke{NewCtrlKeyStroke('e')},
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

	ui := newUIWithKeyMap(&KeyBindingMap{
		DeleteWord: []KeyStroke{
			NewCtrlKeyStroke('w'),
			NewAltKeyStroke(0, "backspace"),
		},
	})

	entries := renderer.buildSearchKeybindEntries(ui)

	if entry, ok := findEntry(entries, "Delete word"); !ok || entry.key != "Ctrl+w, Alt+Backspace" {
		t.Fatalf("expected combined delete word bindings, got %+v", entry)
	}
}

func newUIWithKeyMap(km *KeyBindingMap) *UI {
	state := &UIState{context: ContextSearch}
	ui := &UI{state: state}
	contextMap := NewContextualKeyBindingMap(ProfileDefault, "darwin", "iterm")
	contextMap.SetContext(ContextSearch, km)
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
