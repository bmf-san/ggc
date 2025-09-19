package cmd

import "testing"

func TestParseKeyBinding_Synonyms(t *testing.T) {
	want := ctrl('w')
	cases := []string{"ctrl+w", "^w", "C-w", "c-w", "CTRL+W"}
	for _, in := range cases {
		got, err := ParseKeyBinding(in)
		if err != nil {
			t.Fatalf("ParseKeyBinding(%q) error: %v", in, err)
		}
		if got != want {
			t.Fatalf("ParseKeyBinding(%q) = %d, want %d", in, got, want)
		}
	}
}

func TestParseKeyBinding_Unsupported(t *testing.T) {
	if _, err := ParseKeyBinding("alt+backspace"); err == nil {
		t.Fatalf("expected error for alt+backspace, got nil")
	}
}
