package ui

import (
	"bytes"
	"strings"
	"testing"
)

// ── terminal.go ──────────────────────────────────────────────────────────────

func TestClearScreen(t *testing.T) {
	var buf bytes.Buffer
	ClearScreen(&buf)
	got := buf.String()
	if !strings.Contains(got, "\x1b[2J") {
		t.Errorf("ClearScreen() output %q does not contain escape sequence \\x1b[2J", got)
	}
}

func TestHideCursor(t *testing.T) {
	var buf bytes.Buffer
	HideCursor(&buf)
	if buf.String() != escHideCursor {
		t.Errorf("HideCursor() = %q, want %q", buf.String(), escHideCursor)
	}
}

func TestShowCursor(t *testing.T) {
	var buf bytes.Buffer
	ShowCursor(&buf)
	if buf.String() != escShowCursor {
		t.Errorf("ShowCursor() = %q, want %q", buf.String(), escShowCursor)
	}
}

func TestDisableWrap(t *testing.T) {
	var buf bytes.Buffer
	DisableWrap(&buf)
	if buf.String() != escDisableWrap {
		t.Errorf("DisableWrap() = %q, want %q", buf.String(), escDisableWrap)
	}
}

func TestEnableWrap(t *testing.T) {
	var buf bytes.Buffer
	EnableWrap(&buf)
	if buf.String() != escEnableWrap {
		t.Errorf("EnableWrap() = %q, want %q", buf.String(), escEnableWrap)
	}
}

func TestDimensions_Fallback(t *testing.T) {
	// bytes.Buffer is not an *os.File, so Dimensions uses the fallback.
	var buf bytes.Buffer
	w, h := Dimensions(&buf, 120, 40)
	if w != 120 {
		t.Errorf("Dimensions width = %d, want 120", w)
	}
	if h != 40 {
		t.Errorf("Dimensions height = %d, want 40", h)
	}
}

func TestDimensions_ZeroFallback(t *testing.T) {
	// Zero/negative fallbacks should be replaced by safe defaults (80×24).
	var buf bytes.Buffer
	w, h := Dimensions(&buf, 0, -1)
	if w != 80 {
		t.Errorf("Dimensions width with zero fallback = %d, want 80", w)
	}
	if h != 24 {
		t.Errorf("Dimensions height with negative fallback = %d, want 24", h)
	}
}

// ── text.go ──────────────────────────────────────────────────────────────────

func TestEllipsis(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		maxLen int
		want   string
	}{
		{"empty string", "", 10, ""},
		{"fits exactly", "hello", 5, "hello"},
		{"fits shorter", "hi", 5, "hi"},
		{"truncated", "hello world", 8, "hello w…"},
		{"maxLen 1", "abc", 1, "…"},
		{"maxLen 0", "abc", 0, ""},
		{"maxLen negative", "abc", -1, ""},
		{"maxLen 2", "abc", 2, "a…"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Ellipsis(tt.s, tt.maxLen)
			if got != tt.want {
				t.Errorf("Ellipsis(%q, %d) = %q, want %q", tt.s, tt.maxLen, got, tt.want)
			}
		})
	}
}
