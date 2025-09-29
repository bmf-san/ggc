package cmd

import "testing"

func TestFormatKeyStrokeForDisplay(t *testing.T) {
	tests := []struct {
		name     string
		stroke   KeyStroke
		expected string
	}{
		{
			name:     "ctrl letter",
			stroke:   NewCtrlKeyStroke('q'),
			expected: "Ctrl+q",
		},
		{
			name:     "alt named key",
			stroke:   NewAltKeyStroke(0, "backspace"),
			expected: "Alt+Backspace",
		},
		{
			name:     "alt rune",
			stroke:   NewAltKeyStroke('f', ""),
			expected: "Alt+f",
		},
		{
			name:     "raw tab",
			stroke:   NewRawKeyStroke([]byte{9}),
			expected: "Tab",
		},
		{
			name:     "raw arrow",
			stroke:   NewRawKeyStroke([]byte{27, 91, 67}),
			expected: "→",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := FormatKeyStrokeForDisplay(tc.stroke); got != tc.expected {
				t.Fatalf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestFormatKeyStrokesForDisplayJoinsMultiple(t *testing.T) {
	strokes := []KeyStroke{
		NewCtrlKeyStroke('w'),
		NewAltKeyStroke(0, "backspace"),
	}

	expected := "Ctrl+w, Alt+Backspace"
	if got := FormatKeyStrokesForDisplay(strokes); got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}
