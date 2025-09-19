package cmd

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
