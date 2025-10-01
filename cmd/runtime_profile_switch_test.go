package cmd

import (
	"testing"

	"github.com/bmf-san/ggc/v7/config"
)

func TestRuntimeProfileSwitcherSwitchProfile(t *testing.T) {
	cfg := &config.Config{}
	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)

	contextMgr := NewContextManager(resolver)
	switcher := NewRuntimeProfileSwitcher(resolver, contextMgr)

	// Seed cache to ensure switches clear it
	resolver.cache["dummy"] = NewContextualKeyBindingMap(ProfileDefault, "darwin", "wezterm")

	var callbackOld, callbackNew Profile
	switcher.RegisterSwitchCallback(func(oldProfile, newProfile Profile) {
		callbackOld = oldProfile
		callbackNew = newProfile
	})

	if err := switcher.SwitchProfile(ProfileEmacs); err != nil {
		t.Fatalf("SwitchProfile returned error: %v", err)
	}

	if callbackOld != ProfileDefault || callbackNew != ProfileEmacs {
		t.Fatalf("callback received (%v, %v)", callbackOld, callbackNew)
	}

	if _, exists := resolver.cache["dummy"]; exists {
		t.Fatalf("expected resolver cache to be cleared on switch")
	}

	if got := switcher.GetCurrentProfile(); got != ProfileEmacs {
		t.Fatalf("current profile = %v, want %v", got, ProfileEmacs)
	}

	if err := switcher.SwitchProfile(Profile("missing")); err == nil {
		t.Fatalf("expected error switching to unknown profile")
	}

	if err := switcher.CycleProfile(); err != nil {
		t.Fatalf("CycleProfile returned error: %v", err)
	}
}

func TestProfileSwitcherUpdatesHandler(t *testing.T) {
	cfg := &config.Config{}
	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)

	resolver.platform = ""
	resolver.terminal = ""

	defaultContextual, err := resolver.ResolveContextual(ProfileDefault)
	if err != nil {
		t.Fatalf("ResolveContextual default profile: %v", err)
	}

	state := &UIState{context: ContextInput}
	ui := &UI{state: state}
	ui.handler = &KeyHandler{ui: ui, contextualMap: defaultContextual}

	defaultInput, exists := defaultContextual.GetContext(ContextInput)
	if !exists || defaultInput == nil {
		t.Fatalf("expected default input keymap, exists=%v", exists)
	}
	if !defaultInput.MatchesKeyStroke("delete_word", NewCtrlKeyStroke('w')) {
		t.Fatalf("default profile should use Ctrl+W for delete_word")
	}

	switcher := NewProfileSwitcher(resolver, ui)
	if err := switcher.SwitchProfile(ProfileEmacs); err != nil {
		t.Fatalf("SwitchProfile(ProfileEmacs) error: %v", err)
	}

	emacsContextual := ui.handler.contextualMap
	if emacsContextual == nil {
		t.Fatal("expected handler contextual map to be updated")
	}
	if emacsContextual != nil && emacsContextual.Profile != ProfileEmacs {
		t.Fatalf("contextual profile = %s, want %s", emacsContextual.Profile, ProfileEmacs)
	}

	emacsInput, exists := emacsContextual.GetContext(ContextInput)
	if !exists || emacsInput == nil {
		t.Fatalf("expected emacs input keymap, exists=%v", exists)
	}
	if emacsInput == defaultInput {
		t.Error("expected emacs context to replace default keymap reference")
	}
	if !emacsInput.MatchesKeyStroke("delete_word", NewAltKeyStroke('d', "")) {
		t.Logf("emacs delete_word bindings: %#v", emacsInput.DeleteWord)
		t.Error("emacs profile should use Alt+d for delete_word")
	}

	currentMap := ui.handler.GetCurrentKeyMap()
	if !currentMap.MatchesKeyStroke("delete_word", NewAltKeyStroke('d', "")) {
		t.Logf("current delete_word bindings: %#v", currentMap.DeleteWord)
		t.Error("handler current map should reflect switched profile bindings")
	}
}
