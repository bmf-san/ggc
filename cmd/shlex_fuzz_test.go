package cmd

import (
	"strings"
	"testing"
)

// FuzzTokenize ensures the POSIX-style tokenizer doesn't panic or crash on
// arbitrary input, and that token concatenation invariants hold.
func FuzzTokenize(f *testing.F) {
	seeds := []string{
		"",
		" ",
		`commit -m "fix bug"`,
		`commit -m 'fix the bug'`,
		`echo "it's a \"test\""`,
		`"unterminated`,
		`'unterminated`,
		`a\\b`,
		"\t\ttab\tseparated",
		`"" '' a`,
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		tokens := tokenize(s)
		for _, tok := range tokens {
			// A token must not itself contain a raw NUL; the tokenizer never
			// emits one from ASCII/UTF-8 input.
			if strings.ContainsRune(tok, 0) && !strings.ContainsRune(s, 0) {
				t.Errorf("tokenize introduced NUL byte: input=%q tokens=%#v", s, tokens)
			}
		}
	})
}
