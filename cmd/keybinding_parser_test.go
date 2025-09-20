package cmd

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantCode    byte
		wantValid   bool
		expectError bool
	}{
		// Valid ctrl+ format
		{"ctrl+a lowercase", "ctrl+a", 1, true, false},
		{"ctrl+z lowercase", "ctrl+z", 26, true, false},
		{"CTRL+A uppercase", "CTRL+A", 1, true, false},
		{"Ctrl+W mixed case", "Ctrl+W", 23, true, false},
		{"ctrl+w with spaces", "  ctrl+w  ", 23, true, false},

		// Valid caret notation
		{"^a caret", "^a", 1, true, false},
		{"^Z caret uppercase", "^Z", 26, true, false},
		{"^w caret", "^w", 23, true, false},

		// Valid emacs notation
		{"C-a emacs", "C-a", 1, true, false},
		{"c-a emacs lowercase", "c-a", 1, true, false},
		{"C-W emacs uppercase", "C-W", 23, true, false},
		{"c-z emacs", "c-z", 26, true, false},

		// Invalid inputs
		{"empty string", "", 0, false, true},
		{"whitespace only", "   ", 0, false, true},
		{"ctrl+ too short", "ctrl+", 0, false, true},
		{"ctrl+ too long", "ctrl+ab", 0, false, true},
		{"ctrl+ invalid char", "ctrl+1", 0, false, true},
		{"ctrl+ invalid char @", "ctrl+@", 0, false, true},
		{"caret too short", "^", 0, false, true},
		{"caret too long", "^ab", 0, false, true},
		{"caret invalid char", "^1", 0, false, true},
		{"emacs too short", "c-", 0, false, true},
		{"emacs too long", "c-ab", 0, false, true},
		{"emacs invalid char", "c-1", 0, false, true},
		{"unsupported format", "alt+a", 0, false, true},
		{"random string", "hello", 0, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseKeyBindingInternal(tt.input)

			// Check error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("parseKeyBindingInternal(%q) expected error, got nil", tt.input)
				}
				return
			}

			// Check no error when not expected
			if err != nil {
				t.Errorf("parseKeyBindingInternal(%q) unexpected error: %v", tt.input, err)
				return
			}

			// Check validity
			if result.IsValid != tt.wantValid {
				t.Errorf("parseKeyBindingInternal(%q).IsValid = %v, want %v", tt.input, result.IsValid, tt.wantValid)
			}

			// Check control code
			if result.ControlCode != tt.wantCode {
				t.Errorf("parseKeyBindingInternal(%q).ControlCode = %d, want %d", tt.input, result.ControlCode, tt.wantCode)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		// Valid cases
		{"valid ctrl+w", "ctrl+w", false},
		{"valid ^w", "^w", false},
		{"valid C-w", "C-w", false},
		{"valid c-w", "c-w", false},

		// Invalid cases
		{"empty", "", true},
		{"invalid format", "alt+w", true},
		{"invalid char", "ctrl+1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateKeyBindingInternal(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("validateKeyBindingInternal(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("validateKeyBindingInternal(%q) unexpected error: %v", tt.input, err)
				}
			}
		})
	}
}

func TestCtrlCode(t *testing.T) {
	tests := []struct {
		name string
		char rune
		want byte
	}{
		{"lowercase a", 'a', 1},
		{"lowercase z", 'z', 26},
		{"lowercase w", 'w', 23},
		{"uppercase A", 'A', 1},
		{"uppercase Z", 'Z', 26},
		{"uppercase W", 'W', 23},
		{"invalid digit", '1', 0},
		{"invalid symbol", '@', 0},
		{"invalid space", ' ', 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := keyBindingCtrlCode(tt.char)
			if got != tt.want {
				t.Errorf("keyBindingCtrlCode(%q) = %d, want %d", tt.char, got, tt.want)
			}
		})
	}
}

func TestParseSynonyms(t *testing.T) {
	// Test that all synonymous formats produce the same control code
	want := byte(23) // ctrl+w = 23
	synonyms := []string{"ctrl+w", "^w", "C-w", "c-w", "CTRL+W", "Ctrl+W"}

	for _, input := range synonyms {
		t.Run(input, func(t *testing.T) {
			result, err := parseKeyBindingInternal(input)
			if err != nil {
				t.Fatalf("parseKeyBindingInternal(%q) error: %v", input, err)
			}
			if result.ControlCode != want {
				t.Errorf("parseKeyBindingInternal(%q) = %d, want %d", input, result.ControlCode, want)
			}
			if !result.IsValid {
				t.Errorf("parseKeyBindingInternal(%q).IsValid = false, want true", input)
			}
		})
	}
}
