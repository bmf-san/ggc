package cmd

import (
	"strings"
	"testing"
)

// Test backspace handling with combining marks (e + \u0301)
func TestHandleInputChar_CombiningMarkBackspace(t *testing.T) {
	var output strings.Builder
	ui := &UI{
		colors: NewANSIColors(),
		stdout: &output,
	}
	handler := &KeyHandler{ui: ui}

	var input strings.Builder

	// Type 'e' + combining acute accent
	seq := []rune{'e', 0x0301}
	for _, r := range seq {
		done, canceled := handler.handleInputChar(&input, r)
		if done || canceled {
			t.Fatalf("unexpected completion during input: done=%v canceled=%v", done, canceled)
		}
	}
	if input.String() != string(seq) {
		t.Fatalf("unexpected buffer: got %q want %q", input.String(), string(seq))
	}

	// Backspace once should remove the entire grapheme cluster
	output.Reset()
	done, canceled := handler.handleInputChar(&input, '\b')
	if done || canceled {
		t.Fatalf("unexpected completion on backspace: done=%v canceled=%v", done, canceled)
	}
	if got := input.String(); got != "" {
		t.Fatalf("after backspace, want empty string; got %q", got)
	}
	// Should have cleared exactly one cell (base 'e' = width 1)
	if strings.Count(output.String(), "\b \b") != 1 {
		t.Fatalf("expected one cell cleared, got output: %q", output.String())
	}
}

// Test backspace handling for emoji (width 2)
func TestHandleInputChar_EmojiBackspace(t *testing.T) {
	var output strings.Builder
	ui := &UI{
		colors: NewANSIColors(),
		stdout: &output,
	}
	handler := &KeyHandler{ui: ui}

	var input strings.Builder

	// Type rocket emoji
	rocket := '🚀'
	done, canceled := handler.handleInputChar(&input, rocket)
	if done || canceled {
		t.Fatalf("unexpected completion during input: done=%v canceled=%v", done, canceled)
	}
	if input.String() != string(rocket) {
		t.Fatalf("unexpected buffer: got %q want %q", input.String(), string(rocket))
	}

	// Backspace should clear two cells
	output.Reset()
	done, canceled = handler.handleInputChar(&input, '\b')
	if done || canceled {
		t.Fatalf("unexpected completion on backspace: done=%v canceled=%v", done, canceled)
	}
	if got := input.String(); got != "" {
		t.Fatalf("after backspace, want empty string; got %q", got)
	}
	if strings.Count(output.String(), "\b \b") != 2 {
		t.Fatalf("expected two cells cleared, got output: %q", output.String())
	}
}

// Test backspace handling for a simple ZWJ emoji sequence (woman + ZWJ + rocket)
// Note: We intentionally do not test ZWJ sequences here because the current
// input handler filters out non-printable runes like ZWJ from the buffer.
