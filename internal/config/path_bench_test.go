package config

import (
	"testing"

	"github.com/bmf-san/ggc/v8/internal/testutil"
)

func BenchmarkManagerGet(b *testing.B) {
	mock := testutil.NewMockGitClient()
	cm := NewConfigManager(mock)
	cm.config = getDefaultConfig(mock)
	keys := []string{
		"ui.color",
		"default.branch",
		"default.editor",
		"behavior.auto_fetch",
		"git.default_remote",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, k := range keys {
			_, _ = cm.Get(k)
		}
	}
}

func BenchmarkSanitizeConfigPath(b *testing.B) {
	paths := []string{
		"ui.color",
		"default.branch",
		"interactive.keybindings.global.quit",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, p := range paths {
			_, _ = sanitizeConfigPath(p)
		}
	}
}
