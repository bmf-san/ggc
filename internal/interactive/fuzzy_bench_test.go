package interactive

import "testing"

var benchFuzzyCmds = []string{
	"add",
	"branch checkout",
	"branch delete",
	"branch rename",
	"commit amend",
	"commit amend no-edit",
	"config list",
	"diff",
	"fetch --prune",
	"log graph",
	"push force",
	"pull rebase",
	"rebase interactive",
	"remote set-url",
	"reset hard",
	"restore staged",
	"stash drop",
	"stash pop",
	"status short",
	"tag annotated",
}

func BenchmarkFuzzyMatch(b *testing.B) {
	cases := []struct {
		name    string
		pattern string
	}{
		{"short", "br"},
		{"medium", "push"},
		{"long", "rebase interactive"},
		{"typo_miss", "zzzzz"},
		{"single_char", "a"},
	}
	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, cmd := range benchFuzzyCmds {
					_ = fuzzyMatch(cmd, tc.pattern)
				}
			}
		})
	}
}

func BenchmarkFuzzyMatchScore(b *testing.B) {
	pattern := "br ch"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, cmd := range benchFuzzyCmds {
			_, _ = fuzzyMatchScore(cmd, pattern)
		}
	}
}
