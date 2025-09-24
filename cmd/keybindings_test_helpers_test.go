package cmd

import (
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v6/config"
)

// resolveKeyBindingMapForTest resolves the keybinding map for testing purposes,
// with platform and terminal overrides disabled. It takes a testing.T and an optional
// config.Config, and returns a KeyBindingMap for the default or specified profile.
// This ensures that platform-specific keybinding overrides do not affect test results.
func resolveKeyBindingMapForTest(t *testing.T, cfg *config.Config) *KeyBindingMap {
	t.Helper()

	effectiveCfg := cfg
	if effectiveCfg == nil {
		effectiveCfg = &config.Config{}
	}

	resolver := NewKeyBindingResolver(effectiveCfg)
	RegisterBuiltinProfiles(resolver)

	resolver.platform = ""
	resolver.terminal = ""

	profile := ProfileDefault
	if name := strings.TrimSpace(effectiveCfg.Interactive.Profile); name != "" {
		candidate := Profile(name)
		if candidate.IsValid() {
			profile = candidate
		}
	}

	contextualMap, err := resolver.ResolveContextual(profile)
	if err != nil {
		t.Fatalf("ResolveContextual(%s) error = %v", profile, err)
	}

	keyMap, exists := contextualMap.GetContext(ContextInput)
	if !exists || keyMap == nil {
		t.Fatalf("context %s not resolved", ContextInput)
	}

	return cloneKeyBindingMap(keyMap)
}

// cloneKeyBindingMap creates a deep copy of a KeyBindingMap to prevent test interference.
func cloneKeyBindingMap(src *KeyBindingMap) *KeyBindingMap {
	if src == nil {
		return nil
	}

	return &KeyBindingMap{
		DeleteWord:      cloneKeyStrokes(src.DeleteWord),
		ClearLine:       cloneKeyStrokes(src.ClearLine),
		DeleteToEnd:     cloneKeyStrokes(src.DeleteToEnd),
		MoveToBeginning: cloneKeyStrokes(src.MoveToBeginning),
		MoveToEnd:       cloneKeyStrokes(src.MoveToEnd),
		MoveUp:          cloneKeyStrokes(src.MoveUp),
		MoveDown:        cloneKeyStrokes(src.MoveDown),
	}
}

// cloneKeyStrokes creates and returns a copy of the provided KeyStroke slice.
func cloneKeyStrokes(src []KeyStroke) []KeyStroke {
	if len(src) == 0 {
		return nil
	}
	copySlice := make([]KeyStroke, len(src))
	copy(copySlice, src)
	return copySlice
}

func TestResolveKeyBindingMapForTest_DisablesPlatformOverrides(t *testing.T) {
	cfg := &config.Config{}
	helperMap := resolveKeyBindingMapForTest(t, cfg)

	if containsAltBackspace(helperMap.DeleteWord) {
		t.Fatal("helper map should not include platform-specific alt+backspace override")
	}

	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)
	resolver.platform = "darwin"

	directMap, err := resolver.Resolve(ProfileDefault, ContextInput)
	if err != nil {
		t.Fatalf("Resolve default profile: %v", err)
	}

	if !containsAltBackspace(directMap.DeleteWord) {
		t.Fatal("expected platform overrides to include alt+backspace for delete_word")
	}
}

func containsAltBackspace(strokes []KeyStroke) bool {
	for _, ks := range strokes {
		if ks.Kind == KeyStrokeAlt && strings.EqualFold(ks.Name, "backspace") {
			return true
		}
	}
	return false
}
