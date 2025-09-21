package cmd

import (
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v5/config"
)

func resolveKeyBindingMapForTest(t *testing.T, cfg *config.Config, ctx Context) *KeyBindingMap {
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

	keyMap, exists := contextualMap.GetContext(ctx)
	if !exists || keyMap == nil {
		t.Fatalf("context %s not resolved", ctx)
	}

	return cloneKeyBindingMap(keyMap)
}
