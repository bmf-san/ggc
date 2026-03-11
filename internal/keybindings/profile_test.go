package keybindings

import "testing"

func TestProfileValidation(t *testing.T) {
	valid := []Profile{ProfileDefault, ProfileEmacs, ProfileVi, ProfileReadline}
	for _, profile := range valid {
		if !profile.IsValid() {
			t.Fatalf("expected profile %q to be valid", profile)
		}
		if profile.String() == "" {
			t.Fatalf("expected profile %q to return non-empty string", profile)
		}
	}

	invalid := []Profile{"", "custom", "DEFAULT"}
	for _, profile := range invalid {
		if Profile(profile).IsValid() {
			t.Fatalf("expected profile %q to be invalid", profile)
		}
	}
}

func TestContextValidation(t *testing.T) {
	valid := []Context{ContextGlobal, ContextInput, ContextResults, ContextSearch}
	for _, ctx := range valid {
		if !ctx.IsValid() {
			t.Fatalf("expected context %q to be valid", ctx)
		}
		if ctx.String() == "" {
			t.Fatalf("expected context %q to return non-empty string", ctx)
		}
	}

	invalid := []Context{"", "listing"}
	for _, ctx := range invalid {
		if Context(ctx).IsValid() {
			t.Fatalf("expected context %q to be invalid", ctx)
		}
	}
}

// ── GetAllProfiles / GetAllContexts ──────────────────────────────────────────

func TestGetAllProfiles(t *testing.T) {
	profiles := GetAllProfiles()
	if len(profiles) != 4 {
		t.Fatalf("GetAllProfiles() returned %d profiles, want 4", len(profiles))
	}
}

func TestGetAllContexts_Count(t *testing.T) {
	contexts := GetAllContexts()
	if len(contexts) != 4 {
		t.Fatalf("GetAllContexts() returned %d contexts, want 4", len(contexts))
	}
}

func TestProfile_IsValid(t *testing.T) {
	for _, p := range GetAllProfiles() {
		if !p.IsValid() {
			t.Errorf("Profile(%q).IsValid() = false, want true", p)
		}
	}
	if Profile("invalid").IsValid() {
		t.Error("Profile(\"invalid\").IsValid() = true, want false")
	}
}

func TestContext_IsValid(t *testing.T) {
	for _, c := range GetAllContexts() {
		if !c.IsValid() {
			t.Errorf("Context(%q).IsValid() = false, want true", c)
		}
	}
	if Context("invalid").IsValid() {
		t.Error("Context(\"invalid\").IsValid() = true, want false")
	}
}

// ── KeyBindingProfile methods ─────────────────────────────────────────────────

func TestKeyBindingProfile_GetAllActions(t *testing.T) {
	p := NewKeyBindingProfile("test", "Test profile")
	p.SetGlobalBinding("action1", []KeyStroke{NewCtrlKeyStroke('a')})
	p.SetContextBinding(ContextInput, "action2", []KeyStroke{NewCtrlKeyStroke('b')})

	actions := p.GetAllActions()
	if len(actions) != 2 {
		t.Errorf("GetAllActions() returned %d actions, want 2: %v", len(actions), actions)
	}
}

func TestKeyBindingProfile_Clone(t *testing.T) {
	p := NewKeyBindingProfile("orig", "Original")
	p.SetGlobalBinding("action1", []KeyStroke{NewCtrlKeyStroke('a')})
	p.SetContextBinding(ContextInput, "action2", []KeyStroke{NewCtrlKeyStroke('b')})

	clone := p.Clone()
	if clone.Name != p.Name {
		t.Errorf("Clone().Name = %q, want %q", clone.Name, p.Name)
	}
	if len(clone.Global) != len(p.Global) {
		t.Errorf("Clone().Global len = %d, want %d", len(clone.Global), len(p.Global))
	}
	// Mutations to clone should not affect original
	clone.Global["action1"] = []KeyStroke{NewCtrlKeyStroke('z')}
	if clone.Global["action1"][0].Rune == p.Global["action1"][0].Rune {
		t.Error("Clone is not a deep copy: modifying clone affected original")
	}
}

// ── Builtin profiles ──────────────────────────────────────────────────────────

func TestGetAllProfilesBuiltin(t *testing.T) {
	profiles := GetAllProfilesBuiltin()
	if len(profiles) != 4 {
		t.Fatalf("GetAllProfilesBuiltin() returned %d profiles, want 4", len(profiles))
	}
}

func TestValidateAllBuiltinProfiles(t *testing.T) {
	if err := ValidateAllBuiltinProfiles(); err != nil {
		t.Fatalf("ValidateAllBuiltinProfiles() error: %v", err)
	}
}

func TestValidateProfile_Nil(t *testing.T) {
	if err := ValidateProfile(nil); err == nil {
		t.Error("expected error for nil profile")
	}
}

func TestValidateProfile_EmptyName(t *testing.T) {
	p := NewKeyBindingProfile("", "desc")
	if err := ValidateProfile(p); err == nil {
		t.Error("expected error for empty profile name")
	}
}

func TestValidateProfile_EmptyDescription(t *testing.T) {
	p := NewKeyBindingProfile("test", "")
	if err := ValidateProfile(p); err == nil {
		t.Error("expected error for empty description")
	}
}

func TestValidateProfile_NilContexts(t *testing.T) {
	p := &KeyBindingProfile{Name: "test", Description: "desc", Contexts: nil}
	if err := ValidateProfile(p); err == nil {
		t.Error("expected error for nil contexts")
	}
}

// ── GetProfileStatistics / CompareProfiles ────────────────────────────────────

func TestGetProfileStatistics_Nil(t *testing.T) {
	stats := GetProfileStatistics(nil)
	if len(stats) != 0 {
		t.Errorf("GetProfileStatistics(nil) returned non-empty map: %v", stats)
	}
}

func TestGetProfileStatistics_WithProfile(t *testing.T) {
	p := NewKeyBindingProfile("test", "Test")
	p.SetContextBinding(ContextInput, "move_up", []KeyStroke{NewCtrlKeyStroke('p')})
	stats := GetProfileStatistics(p)
	if stats["profile_name"] != "test" {
		t.Errorf("stats[profile_name] = %v, want %q", stats["profile_name"], "test")
	}
	if stats["total_context_bindings"].(int) != 1 {
		t.Errorf("stats[total_context_bindings] = %v, want 1", stats["total_context_bindings"])
	}
}

func TestCompareProfiles_NilInputs(t *testing.T) {
	result := CompareProfiles(nil, nil)
	if result["error"] == nil {
		t.Error("expected error key in comparison result for nil inputs")
	}
}

func TestCompareProfiles_TwoProfiles(t *testing.T) {
	p1 := NewKeyBindingProfile("p1", "Profile 1")
	p1.SetContextBinding(ContextInput, "move_up", []KeyStroke{NewCtrlKeyStroke('p')})
	p2 := NewKeyBindingProfile("p2", "Profile 2")
	p2.SetContextBinding(ContextResults, "move_down", []KeyStroke{NewCtrlKeyStroke('n')})

	result := CompareProfiles(p1, p2)
	if result["profile1_name"] != "p1" {
		t.Errorf("profile1_name = %v, want p1", result["profile1_name"])
	}
}
