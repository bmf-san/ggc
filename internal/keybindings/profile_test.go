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

func TestWorkflowCreateKeybindingsAvoidNavigationConflicts(t *testing.T) {
	profiles := []struct {
		name    string
		profile *KeyBindingProfile
	}{
		{"default", CreateDefaultProfile()},
		{"emacs", CreateEmacsProfile()},
		{"vi", CreateViProfile()},
		{"readline", CreateReadlineProfile()},
	}

	for _, tc := range profiles {
		bindings := tc.profile.Contexts[ContextWorkflowView]
		moveDown := FormatKeyStrokesForDisplay(bindings["move_down"])
		create := FormatKeyStrokesForDisplay(bindings["workflow_create"])
		if moveDown != "" && create != "" && moveDown == create {
			t.Fatalf("profile %s has conflicting keybindings for move_down (%s) and workflow_create (%s)", tc.name, moveDown, create)
		}
	}
}
