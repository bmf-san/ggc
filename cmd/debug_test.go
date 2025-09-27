package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewDebugger(t *testing.T) {
	debugger := NewDebugger()
	if debugger == nil {
		t.Fatal("NewDebugger() returned nil")
	}
	if debugger.outputWriter == nil {
		t.Error("outputWriter should not be nil")
	}
	if debugger.helper == nil {
		t.Error("helper should not be nil")
	}
}

func TestDebugger_DebugKeys_NoArgs(t *testing.T) {
	var buf bytes.Buffer
	debugger := &Debugger{
		outputWriter: &buf,
		helper:       NewHelper(),
	}

	debugger.DebugKeys([]string{})

	output := buf.String()
	if !strings.Contains(output, "=== Active Key Bindings ===") {
		t.Error("Expected active keybindings output when no args provided")
	}
	if !strings.Contains(output, "Interactive Mode Default Bindings:") {
		t.Error("Expected default bindings section in output")
	}
}

func TestDebugger_DebugKeys_Help(t *testing.T) {
	var buf bytes.Buffer
	debugger := &Debugger{
		outputWriter: &buf,
		helper:       NewHelper(),
	}

	testCases := []string{"help", "-h", "--help"}
	for _, arg := range testCases {
		buf.Reset()
		debugger.DebugKeys([]string{arg})

		output := buf.String()
		if !strings.Contains(output, "debug-keys - Debug keybinding issues") {
			t.Errorf("Expected help text for arg '%s'", arg)
		}
		if !strings.Contains(output, "USAGE:") {
			t.Errorf("Expected USAGE section for arg '%s'", arg)
		}
		if !strings.Contains(output, "SUBCOMMANDS:") {
			t.Errorf("Expected SUBCOMMANDS section for arg '%s'", arg)
		}
	}
}

func TestDebugger_DebugKeys_UnknownSubcommand(t *testing.T) {
	var buf bytes.Buffer
	debugger := &Debugger{
		outputWriter: &buf,
		helper:       NewHelper(),
	}

	debugger.DebugKeys([]string{"unknown"})

	output := buf.String()
	if !strings.Contains(output, "Unknown subcommand: unknown") {
		t.Error("Expected unknown subcommand error message")
	}
	if !strings.Contains(output, "debug-keys - Debug keybinding issues") {
		t.Error("Expected help text to be shown after unknown subcommand")
	}
}

func TestDebugger_showActiveKeybindings(t *testing.T) {
	var buf bytes.Buffer
	debugger := &Debugger{
		outputWriter: &buf,
		helper:       NewHelper(),
	}

	debugger.showActiveKeybindings()

	output := buf.String()
	expectedContent := []string{
		"=== Active Key Bindings ===",
		"Interactive Mode Default Bindings:",
		"Ctrl+P",
		"Move selection up",
		"Ctrl+N",
		"Move selection down",
		"Ctrl+U",
		"Clear input",
		"Ctrl+W",
		"Delete previous word",
		"Ctrl+K",
		"Delete to end of line",
		"Ctrl+A",
		"Move cursor to start",
		"Ctrl+E",
		"Move cursor to end",
		"Backspace",
		"Delete character before cursor",
		"Alt+Backspace",
		"Delete previous word (terminal dependent)",
		"←/→",
		"Move cursor",
		"Ctrl+←/→",
		"Move by word (terminal dependent)",
		"Custom keybinding configuration:",
		"Use 'ggc debug-keys raw'",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain '%s'", expected)
		}
	}
}

func TestDebugger_showDebugKeysHelp(t *testing.T) {
	var buf bytes.Buffer
	debugger := &Debugger{
		outputWriter: &buf,
		helper:       NewHelper(),
	}

	debugger.showDebugKeysHelp()

	output := buf.String()
	expectedContent := []string{
		"debug-keys - Debug keybinding issues and capture raw key sequences",
		"USAGE:",
		"SUBCOMMANDS:",
		"EXAMPLES:",
		"DESCRIPTION:",
		"ggc debug-keys",
		"ggc debug-keys raw",
		"ggc debug-keys raw keys.txt",
		"Press Ctrl+C in raw mode to stop capturing",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected help output to contain '%s'", expected)
		}
	}
}

// Test that the raw mode would fail gracefully in non-terminal environment
func TestDebugger_captureRawKeySequences_NonTerminal(t *testing.T) {
	var buf bytes.Buffer
	debugger := &Debugger{
		outputWriter: &buf,
		helper:       NewHelper(),
	}

	// In test environment, stdin is typically not a terminal
	debugger.captureRawKeySequences("")

	output := buf.String()
	if !strings.Contains(output, "Error: debug-keys raw mode requires a terminal") {
		t.Error("Expected error message for non-terminal environment")
	}
}
