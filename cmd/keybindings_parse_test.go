package cmd

import (
	"strings"
	"testing"
)

func TestParseKeyStrokeVariants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		wantKind KeyStrokeKind
		wantRune rune
		wantName string
	}{
		{name: "ctrl uppercase", input: "Ctrl+A", wantKind: KeyStrokeCtrl, wantRune: 'a'},
		{name: "caret notation", input: "^Z", wantKind: KeyStrokeCtrl, wantRune: 'z'},
		{name: "emacs notation", input: "C-k", wantKind: KeyStrokeCtrl, wantRune: 'k'},
		{name: "alt special", input: "Alt+Backspace", wantKind: KeyStrokeAlt, wantName: "backspace"},
		{name: "meta letter", input: "M-b", wantKind: KeyStrokeAlt, wantRune: 'b'},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ks, err := ParseKeyStroke(tt.input)
			if err != nil {
				t.Fatalf("ParseKeyStroke(%q) returned error: %v", tt.input, err)
			}
			if ks.Kind != tt.wantKind {
				t.Fatalf("ParseKeyStroke(%q) kind = %v, want %v", tt.input, ks.Kind, tt.wantKind)
			}
			if ks.Rune != tt.wantRune {
				t.Fatalf("ParseKeyStroke(%q) rune = %q, want %q", tt.input, ks.Rune, tt.wantRune)
			}
			if ks.Name != tt.wantName {
				t.Fatalf("ParseKeyStroke(%q) name = %q, want %q", tt.input, ks.Name, tt.wantName)
			}
		})
	}
}

func TestParseKeyStrokeInvalid(t *testing.T) {
	t.Parallel()

	invalid := []string{"", "Alt+1", "Ctrl+", "Shift+A", "meta+unknown"}
	for _, input := range invalid {
		input := input
		t.Run(input, func(t *testing.T) {
			if _, err := ParseKeyStroke(input); err == nil {
				t.Fatalf("ParseKeyStroke(%q) expected error", input)
			}
		})
	}
}

func TestParseKeyStrokesContainers(t *testing.T) {
	t.Parallel()

	single, err := ParseKeyStrokes("Ctrl+U")
	if err != nil {
		t.Fatalf("ParseKeyStrokes single returned error: %v", err)
	}
	if len(single) != 1 || single[0].Kind != KeyStrokeCtrl || single[0].Rune != 'u' {
		t.Fatalf("unexpected single ParseKeyStrokes result: %#v", single)
	}

	ifaceSlice, err := ParseKeyStrokes([]interface{}{"Ctrl+N", "Alt+Backspace"})
	if err != nil {
		t.Fatalf("ParseKeyStrokes interface slice error: %v", err)
	}
	if len(ifaceSlice) != 2 {
		t.Fatalf("expected 2 key strokes, got %d", len(ifaceSlice))
	}

	strSlice, err := ParseKeyStrokes([]string{"Ctrl+P", "M-f"})
	if err != nil {
		t.Fatalf("ParseKeyStrokes string slice error: %v", err)
	}
	if len(strSlice) != 2 || strSlice[1].Kind != KeyStrokeAlt {
		t.Fatalf("unexpected string slice ParseKeyStrokes result: %#v", strSlice)
	}

	if _, err := ParseKeyStrokes(42); err == nil {
		t.Fatalf("ParseKeyStrokes expected error for unsupported type")
	}

	_, err = ParseKeyStrokes([]interface{}{123})
	if err == nil {
		t.Fatalf("ParseKeyStrokes expected error for non-string array element")
	}
}

func TestKeyStrokeEqualsAndControlByte(t *testing.T) {
	t.Parallel()

	ctrlA := NewCtrlKeyStroke('a')
	ctrlACopy := NewCtrlKeyStroke('a')
	if !ctrlA.Equals(ctrlACopy) {
		t.Fatalf("expected ctrlA to equal copy")
	}
	if ctrlA.ToControlByte() == 0 {
		t.Fatalf("expected ctrlA to map to control byte")
	}

	altB := NewAltKeyStroke('b', "")
	if altB.Equals(ctrlA) {
		t.Fatalf("expected altB and ctrlA to differ")
	}

	rawSeq := NewRawKeyStroke([]byte{27, 91})
	rawSeqCopy := NewRawKeyStroke([]byte{27, 91})
	if !rawSeq.Equals(rawSeqCopy) {
		t.Fatalf("expected raw sequences to match")
	}

	rawDifferent := NewRawKeyStroke([]byte{27, 92})
	if rawSeq.Equals(rawDifferent) {
		t.Fatalf("expected different raw sequences to differ")
	}

	fnKey := KeyStroke{Kind: KeyStrokeFnKey, Name: "F1"}
	fnKeyCopy := KeyStroke{Kind: KeyStrokeFnKey, Name: "F1"}
	if !fnKey.Equals(fnKeyCopy) {
		t.Fatalf("expected function keys to compare equal")
	}

	fnKeyOther := KeyStroke{Kind: KeyStrokeFnKey, Name: "F2"}
	if fnKey.Equals(fnKeyOther) {
		t.Fatalf("expected different function keys to differ")
	}
}

func TestParseKeyStrokeErrorMessages(t *testing.T) {
	t.Parallel()

	_, err := ParseKeyStroke("Alt+Entertain")
	if err == nil {
		t.Fatalf("expected alt error")
	}
	if !strings.Contains(err.Error(), "unsupported alt key") {
		t.Fatalf("unexpected error message: %v", err)
	}
}
