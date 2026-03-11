package keybindings

import (
	"testing"

	"github.com/bmf-san/ggc/v8/internal/config"
)

func TestKeyBindingResolverLayering(t *testing.T) {
	cfg := &config.Config{}
	cfg.Interactive.Keybindings.DeleteWord = "Ctrl+X"
	cfg.Interactive.Contexts.Results.Keybindings = map[string]interface{}{
		"move_up": []interface{}{"Alt+F"},
	}
	cfg.Interactive.Darwin.Keybindings = map[string]interface{}{
		"move_down": "Ctrl+J",
	}
	cfg.Interactive.Terminals = map[string]config.KeybindingsConfig{
		"wezterm": {
			Keybindings: map[string]interface{}{
				"move_to_end": "Ctrl+L",
			},
		},
	}

	resolver := NewKeyBindingResolver(cfg)
	resolver.platform = "darwin"
	resolver.terminal = "wezterm"

	const testProfile = Profile("custom")
	profile := NewKeyBindingProfile("custom", "test profile")
	profile.SetGlobalBinding("delete_word", []KeyStroke{NewCtrlKeyStroke('d')})
	profile.SetContextBinding(ContextResults, "move_up", []KeyStroke{NewCtrlKeyStroke('p')})
	profile.SetContextBinding(ContextResults, "move_down", []KeyStroke{NewCtrlKeyStroke('n')})
	profile.SetContextBinding(ContextInput, "move_to_end", []KeyStroke{NewCtrlKeyStroke('e')})
	profile.SetContextBinding(ContextInput, "move_to_beginning", []KeyStroke{NewCtrlKeyStroke('a')})
	profile.Contexts[ContextGlobal] = map[string][]KeyStroke{}
	profile.Contexts[ContextSearch] = map[string][]KeyStroke{}
	resolver.RegisterProfile(testProfile, profile)

	keyMap, err := resolver.Resolve(testProfile, ContextResults)
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}

	if len(keyMap.DeleteWord) != 1 || keyMap.DeleteWord[0].Kind != KeyStrokeCtrl || keyMap.DeleteWord[0].Rune != 'x' {
		t.Fatalf("user override for delete_word not applied: %#v", keyMap.DeleteWord)
	}

	if len(keyMap.MoveUp) != 1 || keyMap.MoveUp[0].Kind != KeyStrokeAlt || keyMap.MoveUp[0].Rune != 'f' {
		t.Fatalf("context user override for move_up not applied: %#v", keyMap.MoveUp)
	}

	if len(keyMap.MoveDown) != 1 || keyMap.MoveDown[0].Rune != 'j' {
		t.Fatalf("platform override for move_down not applied: %#v", keyMap.MoveDown)
	}

	inputMap, err := resolver.Resolve(testProfile, ContextInput)
	if err != nil {
		t.Fatalf("Resolve input returned error: %v", err)
	}

	if len(inputMap.MoveToEnd) != 1 || inputMap.MoveToEnd[0].Rune != 'l' {
		t.Fatalf("terminal override for move_to_end not applied: %#v", inputMap.MoveToEnd)
	}

	cachedMap, err := resolver.Resolve(testProfile, ContextResults)
	if err != nil {
		t.Fatalf("Resolve cached returned error: %v", err)
	}
	if cachedMap != keyMap {
		t.Fatalf("expected resolver to return cached pointer for repeated calls")
	}
}

func TestKeyBindingResolverUserBindingParsing(t *testing.T) {
	cfg := &config.Config{}
	cfg.Interactive.Contexts.Input.Keybindings = map[string]interface{}{
		"move_up": []interface{}{"Ctrl+P", "Ctrl+N"},
	}

	resolver := NewKeyBindingResolver(cfg)
	resolver.platform = "linux"
	resolver.terminal = "xterm"

	profile := NewKeyBindingProfile("minimal", "desc")
	profile.SetContextBinding(ContextInput, "move_up", []KeyStroke{NewCtrlKeyStroke('p')})
	profile.SetContextBinding(ContextInput, "move_down", []KeyStroke{NewCtrlKeyStroke('n')})
	profile.Contexts[ContextGlobal] = map[string][]KeyStroke{}
	profile.Contexts[ContextResults] = map[string][]KeyStroke{}
	profile.Contexts[ContextSearch] = map[string][]KeyStroke{}
	resolver.RegisterProfile(Profile("minimal"), profile)

	keyMap, err := resolver.Resolve(Profile("minimal"), ContextInput)
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}

	if len(keyMap.MoveUp) != 2 {
		t.Fatalf("expected two move_up bindings, got %d", len(keyMap.MoveUp))
	}
}

func TestKeyBindingResolverLayerPrecedence(t *testing.T) {
	cfg := &config.Config{}
	cfg.Interactive.Keybindings.MoveUp = "Ctrl+Q"
	cfg.Interactive.Contexts.Input.Keybindings = map[string]interface{}{
		"move_up": []interface{}{"Ctrl+R"},
	}
	cfg.Interactive.Darwin.Keybindings = map[string]interface{}{
		"move_up": []interface{}{"Ctrl+S"},
	}
	cfg.Interactive.Terminals = map[string]config.KeybindingsConfig{
		"wezterm": {Keybindings: map[string]interface{}{"move_up": "Ctrl+T"}},
	}

	t.Setenv("GGC_KEYBIND_MOVE_UP", "Ctrl+U")

	resolver := NewKeyBindingResolver(cfg)
	resolver.platform = "darwin"
	resolver.terminal = "wezterm"
	RegisterBuiltinProfiles(resolver)

	keyMap, err := resolver.Resolve(ProfileDefault, ContextInput)
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}

	if len(keyMap.MoveUp) != 1 || keyMap.MoveUp[0].Rune != 'u' {
		t.Fatalf("environment override should win, got %#v", keyMap.MoveUp)
	}
}

func TestResolveContextualAppliesOverridesPerContext(t *testing.T) {
	cfg := &config.Config{}
	cfg.Interactive.Contexts.Input.Keybindings = map[string]interface{}{
		"move_up": "Ctrl+P",
	}
	cfg.Interactive.Contexts.Results.Keybindings = map[string]interface{}{
		"move_up": "Ctrl+R",
	}
	cfg.Interactive.Contexts.Search.Keybindings = map[string]interface{}{
		"move_up": []interface{}{"Ctrl+S", "Ctrl+T"},
	}

	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)

	contextual, err := resolver.ResolveContextual(ProfileDefault)
	if err != nil {
		t.Fatalf("ResolveContextual returned error: %v", err)
	}

	ctxs := map[Context]struct {
		expected rune
	}{
		ContextInput:   {expected: 'p'},
		ContextResults: {expected: 'r'},
	}

	for ctx, want := range ctxs {
		keyMap, exists := contextual.GetContext(ctx)
		if !exists {
			t.Fatalf("missing context %v", ctx)
		}
		if len(keyMap.MoveUp) != 1 || keyMap.MoveUp[0].Rune != want.expected {
			t.Fatalf("context %v move_up = %#v", ctx, keyMap.MoveUp)
		}
	}

	searchMap, exists := contextual.GetContext(ContextSearch)
	if !exists {
		t.Fatalf("missing search context")
	}
	if len(searchMap.MoveUp) != 2 || searchMap.MoveUp[0].Rune != 's' || searchMap.MoveUp[1].Rune != 't' {
		t.Fatalf("search context bindings unexpected: %#v", searchMap.MoveUp)
	}
}

// ── resolver_user: editing/navigation actions, parseUserBindingValue ──────────

func testResolver() *KeyBindingResolver {
	cfg := &config.Config{}
	cfg.Interactive.Profile = "emacs"
	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)
	return resolver
}

func TestApplyUserEditingAction_ClearLine(t *testing.T) {
	r := testResolver()
	km := DefaultKeyBindingMap()
	ks := []KeyStroke{NewTabKeyStroke()}
	if !r.applyUserEditingAction(km, "clear_line", ks) {
		t.Error("expected applyUserEditingAction to return true for clear_line")
	}
	if len(km.ClearLine) == 0 {
		t.Error("expected ClearLine to be set")
	}
}

func TestApplyUserEditingAction_DeleteToEnd(t *testing.T) {
	r := testResolver()
	km := DefaultKeyBindingMap()
	ks := []KeyStroke{NewTabKeyStroke()}
	if !r.applyUserEditingAction(km, "delete_to_end", ks) {
		t.Error("expected applyUserEditingAction to return true for delete_to_end")
	}
	if len(km.DeleteToEnd) == 0 {
		t.Error("expected DeleteToEnd to be set")
	}
}

func TestApplyUserNavigationAction_MoveLeft(t *testing.T) {
	r := testResolver()
	km := DefaultKeyBindingMap()
	ks := []KeyStroke{NewTabKeyStroke()}
	if !r.applyUserNavigationAction(km, "move_left", ks) {
		t.Error("expected applyUserNavigationAction to return true for move_left")
	}
	if len(km.MoveLeft) == 0 {
		t.Error("expected MoveLeft to be set")
	}
}

func TestApplyUserNavigationAction_MoveRight(t *testing.T) {
	r := testResolver()
	km := DefaultKeyBindingMap()
	ks := []KeyStroke{NewTabKeyStroke()}
	if !r.applyUserNavigationAction(km, "move_right", ks) {
		t.Error("expected applyUserNavigationAction to return true for move_right")
	}
	if len(km.MoveRight) == 0 {
		t.Error("expected MoveRight to be set")
	}
}

func TestParseUserBindingValue_Slice(t *testing.T) {
	r := testResolver()
	got := r.parseUserBindingValue([]interface{}{"ctrl+w", "ctrl+u"})
	if len(got) != 2 {
		t.Errorf("expected 2 keystrokes from slice input, got %d", len(got))
	}
}
