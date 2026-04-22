package config

import (
	"strings"
	"testing"
	"unicode/utf8"
)

// FuzzValidatePlaceholder ensures placeholder validation never panics on
// arbitrary strings and rejects anything containing characters outside the
// documented alphabet (alnum plus `_-`).
func FuzzValidatePlaceholder(f *testing.F) {
	seeds := []string{
		"",
		"0",
		"branch",
		"kebab-case",
		"snake_case",
		"MixedCase123",
		"with space",
		"dot.invalid",
		"emoji🙂",
		"-leading",
		"trailing-",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		err := validatePlaceholder(s)
		if err != nil {
			return
		}
		// When validation passes the input must consist of only valid bytes.
		if s == "" {
			t.Errorf("validatePlaceholder accepted empty input")
			return
		}
		if !utf8.ValidString(s) {
			t.Errorf("validatePlaceholder accepted invalid UTF-8: %q", s)
			return
		}
		for _, r := range s {
			if !isValidPlaceholderChar(r) {
				t.Errorf("validatePlaceholder accepted %q which contains invalid rune %q", s, r)
				return
			}
		}
	})
}

// FuzzAnalyzePlaceholders drives the full placeholder extractor to catch
// panics on malformed alias commands.
func FuzzAnalyzePlaceholders(f *testing.F) {
	seeds := []string{
		"",
		"ggc status",
		"ggc commit -m {0}",
		"ggc branch checkout {branch}",
		"ggc commit -m {message} {0} {1}",
		"ggc {} empty",
		"ggc {emoji🙂} weird",
		"ggc {a-b-c} dashes",
		strings.Repeat("{x}", 100),
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		_, _, _ = analyzePlaceholders([]string{s})
	})
}
