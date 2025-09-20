package cmd

import (
	"testing"

	"github.com/bmf-san/ggc/v5/config"
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
