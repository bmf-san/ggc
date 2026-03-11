package keybindings

import "testing"

// ── KeyStrokeKind.String ─────────────────────────────────────────────────────

func TestKeyStrokeKind_String(t *testing.T) {
	tests := []struct {
		kind KeyStrokeKind
		want string
	}{
		{KeyStrokeCtrl, "Ctrl"},
		{KeyStrokeAlt, "Alt"},
		{KeyStrokeRawSeq, "RawSeq"},
		{KeyStrokeFnKey, "FnKey"},
		{KeyStrokeKind(99), "Unknown"},
	}
	for _, tt := range tests {
		got := tt.kind.String()
		if got != tt.want {
			t.Errorf("KeyStrokeKind(%d).String() = %q, want %q", tt.kind, got, tt.want)
		}
	}
}

// ── KeyStroke constructors ───────────────────────────────────────────────────

func TestNewEnterKeyStroke(t *testing.T) {
	ks := NewEnterKeyStroke()
	if ks.Kind != KeyStrokeRawSeq {
		t.Errorf("NewEnterKeyStroke().Kind = %v, want KeyStrokeRawSeq", ks.Kind)
	}
	if len(ks.Seq) != 1 || ks.Seq[0] != 13 {
		t.Errorf("NewEnterKeyStroke().Seq = %v, want [13]", ks.Seq)
	}
}

func TestNewSpaceKeyStroke(t *testing.T) {
	ks := NewSpaceKeyStroke()
	if ks.Kind != KeyStrokeRawSeq {
		t.Errorf("NewSpaceKeyStroke().Kind = %v, want KeyStrokeRawSeq", ks.Kind)
	}
	if len(ks.Seq) != 1 || ks.Seq[0] != 32 {
		t.Errorf("NewSpaceKeyStroke().Seq = %v, want [32]", ks.Seq)
	}
}

// ── KeyStroke.String ─────────────────────────────────────────────────────────

func TestKeyStroke_String_Ctrl(t *testing.T) {
	ks := NewCtrlKeyStroke('a')
	got := ks.String()
	if got != "Ctrl+a" {
		t.Errorf("KeyStroke.String() for ctrl = %q, want %q", got, "Ctrl+a")
	}
}

func TestKeyStroke_String_Alt_WithName(t *testing.T) {
	ks := KeyStroke{Kind: KeyStrokeAlt, Name: "backspace"}
	got := ks.String()
	if got != "Alt+backspace" {
		t.Errorf("KeyStroke.String() for Alt+name = %q, want %q", got, "Alt+backspace")
	}
}

func TestKeyStroke_String_Alt_WithRune(t *testing.T) {
	ks := KeyStroke{Kind: KeyStrokeAlt, Rune: 'x'}
	got := ks.String()
	if got != "Alt+x" {
		t.Errorf("KeyStroke.String() for Alt+rune = %q, want %q", got, "Alt+x")
	}
}

func TestKeyStroke_String_FnKey(t *testing.T) {
	ks := KeyStroke{Kind: KeyStrokeFnKey, Name: "F1"}
	got := ks.String()
	if got != "F1" {
		t.Errorf("KeyStroke.String() for FnKey = %q, want 'F1'", got)
	}
}

func TestKeyStroke_String_Unknown(t *testing.T) {
	ks := KeyStroke{Kind: KeyStrokeKind(999)}
	got := ks.String()
	if got != "Unknown" {
		t.Errorf("KeyStroke.String() for unknown = %q, want 'Unknown'", got)
	}
}

// ── validateKeyStroke ────────────────────────────────────────────────────────

func TestValidateKeyStroke_Ctrl_Valid(t *testing.T) {
	ks := NewCtrlKeyStroke('w')
	if err := validateKeyStroke(ks); err != nil {
		t.Errorf("unexpected error for valid ctrl keystroke: %v", err)
	}
}

func TestValidateKeyStroke_Ctrl_InvalidRune(t *testing.T) {
	ks := KeyStroke{Kind: KeyStrokeCtrl, Rune: 0}
	if err := validateKeyStroke(ks); err == nil {
		t.Error("expected error for ctrl with zero rune")
	}
}

func TestValidateKeyStroke_Alt_ValidRune(t *testing.T) {
	ks := KeyStroke{Kind: KeyStrokeAlt, Rune: 'b'}
	if err := validateKeyStroke(ks); err != nil {
		t.Errorf("unexpected error for valid alt keystroke: %v", err)
	}
}

func TestValidateKeyStroke_Alt_Empty(t *testing.T) {
	ks := KeyStroke{Kind: KeyStrokeAlt}
	if err := validateKeyStroke(ks); err == nil {
		t.Error("expected error for alt keystroke with no rune or name")
	}
}

func TestValidateKeyStroke_RawSeq_Valid(t *testing.T) {
	ks := NewEnterKeyStroke()
	if err := validateKeyStroke(ks); err != nil {
		t.Errorf("unexpected error for valid raw seq: %v", err)
	}
}

func TestValidateKeyStroke_RawSeq_Empty(t *testing.T) {
	ks := KeyStroke{Kind: KeyStrokeRawSeq, Seq: []byte{}}
	if err := validateKeyStroke(ks); err == nil {
		t.Error("expected error for raw seq with empty sequence")
	}
}

func TestValidateKeyStroke_FnKey_Valid(t *testing.T) {
	ks := KeyStroke{Kind: KeyStrokeFnKey, Name: "F1"}
	if err := validateKeyStroke(ks); err != nil {
		t.Errorf("unexpected error for valid fn key: %v", err)
	}
}

func TestValidateKeyStroke_FnKey_NoName(t *testing.T) {
	ks := KeyStroke{Kind: KeyStrokeFnKey}
	if err := validateKeyStroke(ks); err == nil {
		t.Error("expected error for fn key with no name")
	}
}

func TestValidateKeyStroke_Unknown(t *testing.T) {
	ks := KeyStroke{Kind: KeyStrokeKind(999)}
	if err := validateKeyStroke(ks); err == nil {
		t.Error("expected error for unknown keystroke kind")
	}
}
