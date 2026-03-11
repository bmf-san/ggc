package keybindings

import (
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v8/internal/config"
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

type mockMapApplier struct {
	last *ContextualKeyBindingMap
}

func (m *mockMapApplier) ApplyContextualKeybindings(km *ContextualKeyBindingMap) {
	m.last = km
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

	defaultInput, exists := defaultContextual.GetContext(ContextInput)
	if !exists || defaultInput == nil {
		t.Fatalf("expected default input keymap, exists=%v", exists)
	}
	if !defaultInput.MatchesKeyStroke("delete_word", NewCtrlKeyStroke('w')) {
		t.Fatalf("default profile should use Ctrl+W for delete_word")
	}

	applier := &mockMapApplier{}
	switcher := NewProfileSwitcher(resolver, applier)
	if err := switcher.SwitchProfile(ProfileEmacs); err != nil {
		t.Fatalf("SwitchProfile(ProfileEmacs) error: %v", err)
	}

	if applier.last == nil {
		t.Fatal("expected applier to receive contextual map")
	}
	if applier.last.Profile != ProfileEmacs {
		t.Fatalf("contextual profile = %s, want %s", applier.last.Profile, ProfileEmacs)
	}

	emacsInput, exists := applier.last.GetContext(ContextInput)
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
}

// ── ProfileSwitcher ──────────────────────────────────────────────────────────

func newTestSwitcher() *ProfileSwitcher {
	resolver := NewKeyBindingResolver(nil)
	RegisterBuiltinProfiles(resolver)
	return NewProfileSwitcher(resolver, nil)
}

func TestProfileSwitcher_GetCurrentProfile(t *testing.T) {
	ps := newTestSwitcher()
	if ps.GetCurrentProfile() != ProfileDefault {
		t.Errorf("GetCurrentProfile() = %v, want %v", ps.GetCurrentProfile(), ProfileDefault)
	}
}

func TestProfileSwitcher_GetAvailableProfiles(t *testing.T) {
	ps := newTestSwitcher()
	profiles := ps.GetAvailableProfiles()
	if len(profiles) != 4 {
		t.Fatalf("GetAvailableProfiles() returned %d profiles, want 4", len(profiles))
	}
}

func TestProfileSwitcher_CanSwitchTo_Valid(t *testing.T) {
	ps := newTestSwitcher()
	ok, err := ps.CanSwitchTo(ProfileEmacs)
	if err != nil {
		t.Fatalf("CanSwitchTo(emacs) error: %v", err)
	}
	if !ok {
		t.Error("CanSwitchTo(emacs) = false, want true")
	}
}

func TestProfileSwitcher_CanSwitchTo_NotRegistered(t *testing.T) {
	resolver := NewKeyBindingResolver(nil) // no profiles registered
	ps := NewProfileSwitcher(resolver, nil)
	_, err := ps.CanSwitchTo(ProfileEmacs)
	if err == nil {
		t.Error("expected error for unregistered profile")
	}
}

func TestProfileSwitcher_GetProfileComparison(t *testing.T) {
	ps := newTestSwitcher()
	result, err := ps.GetProfileComparison(ProfileEmacs)
	if err != nil {
		t.Fatalf("GetProfileComparison(emacs) error: %v", err)
	}
	if result["profile1_name"] == nil {
		t.Error("expected profile1_name in comparison result")
	}
}

func TestProfileSwitcher_ShowCurrentProfileCommand(t *testing.T) {
	ps := newTestSwitcher()
	got := ShowCurrentProfileCommand(ps)
	if !strings.Contains(got, "default") {
		t.Errorf("ShowCurrentProfileCommand() = %q, want to contain 'default'", got)
	}
}

// ── HandleProfileSwitchCommand ───────────────────────────────────────────────

func TestHandleProfileSwitchCommand_List(t *testing.T) {
	ps := newTestSwitcher()
	if err := HandleProfileSwitchCommand(ps, "list"); err != nil {
		t.Fatalf("HandleProfileSwitchCommand(list) error: %v", err)
	}
}

func TestHandleProfileSwitchCommand_Unknown(t *testing.T) {
	ps := newTestSwitcher()
	if err := HandleProfileSwitchCommand(ps, "bogus"); err == nil {
		t.Error("expected error for unknown subcommand")
	}
}

func TestHandleProfileSwitchCommand_Empty(t *testing.T) {
	ps := newTestSwitcher()
	if err := HandleProfileSwitchCommand(ps, ""); err == nil {
		t.Error("expected error for empty command")
	}
}

func TestHandleProfileSwitchCommand_Preview(t *testing.T) {
	ps := newTestSwitcher()
	if err := HandleProfileSwitchCommand(ps, "preview emacs"); err != nil {
		t.Fatalf("HandleProfileSwitchCommand(preview emacs) error: %v", err)
	}
}

func TestHandleProfileSwitchCommand_Compare(t *testing.T) {
	ps := newTestSwitcher()
	if err := HandleProfileSwitchCommand(ps, "compare emacs"); err != nil {
		t.Fatalf("HandleProfileSwitchCommand(compare emacs) error: %v", err)
	}
}

func TestHandleProfileSwitchCommand_PreviewNoArg(t *testing.T) {
	ps := newTestSwitcher()
	if err := HandleProfileSwitchCommand(ps, "preview"); err == nil {
		t.Error("expected error for preview without arg")
	}
}
