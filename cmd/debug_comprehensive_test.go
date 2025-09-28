package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// Test NewDebugger constructor
func TestNewDebugger_Coverage(t *testing.T) {
	debugger := NewDebugger()

	t.Run("debugger_creation", func(t *testing.T) {
		if debugger == nil {
			t.Fatal("NewDebugger returned nil")
		}
	})

	t.Run("output_writer_setup", func(t *testing.T) {
		if debugger == nil {
			t.Skip("Skipping due to nil debugger")
		}
		if debugger.outputWriter != os.Stdout {
			t.Fatalf("outputWriter should be set to os.Stdout, got: %v", debugger.outputWriter)
		}
	})

	t.Run("helper_initialization", func(t *testing.T) {
		if debugger == nil {
			t.Skip("Skipping due to nil debugger")
		}
		if debugger.helper == nil {
			t.Fatal("helper should be initialized")
		}
	})
}

// Test DebugKeys method edge cases
func TestDebugger_DebugKeys_EdgeCases(t *testing.T) {
	var buf bytes.Buffer
	debugger := &Debugger{
		outputWriter: &buf,
		helper:       NewHelper(),
	}

	testCases := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name:     "empty args",
			args:     []string{},
			expected: []string{"=== Active Key Bindings ==="},
		},
		{
			name:     "raw subcommand no file",
			args:     []string{"raw"},
			expected: []string{"Error: debug-keys raw mode requires a terminal"},
		},
		{
			name:     "raw subcommand with file",
			args:     []string{"raw", "output.txt"},
			expected: []string{"Error: debug-keys raw mode requires a terminal"},
		},
		{
			name:     "help subcommand",
			args:     []string{"help"},
			expected: []string{"debug-keys - Debug keybinding issues"},
		},
		{
			name:     "short help flag",
			args:     []string{"-h"},
			expected: []string{"debug-keys - Debug keybinding issues"},
		},
		{
			name:     "long help flag",
			args:     []string{"--help"},
			expected: []string{"debug-keys - Debug keybinding issues"},
		},
		{
			name:     "unknown subcommand",
			args:     []string{"invalid"},
			expected: []string{"Unknown subcommand: invalid", "debug-keys - Debug keybinding issues"},
		},
		{
			name:     "multiple args unknown",
			args:     []string{"invalid", "extra", "args"},
			expected: []string{"Unknown subcommand: invalid"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf.Reset()
			debugger.DebugKeys(tc.args)
			output := buf.String()

			for _, expected := range tc.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got: %s", expected, output)
				}
			}
		})
	}
}

// Test showActiveKeybindings with different configurations
func TestDebugger_showActiveKeybindings_Details(t *testing.T) {
	var buf bytes.Buffer
	debugger := &Debugger{
		outputWriter: &buf,
		helper:       NewHelper(),
	}

	debugger.showActiveKeybindings()
	output := buf.String()

	// Test specific keybinding entries
	expectedKeybindings := map[string]string{
		"Ctrl+P":        "Move selection up",
		"Ctrl+N":        "Move selection down",
		"Enter":         "Execute selected command",
		"Ctrl+C":        "Quit interactive mode",
		"Ctrl+U":        "Clear input",
		"Ctrl+W":        "Delete previous word",
		"Ctrl+K":        "Delete to end of line",
		"Ctrl+A":        "Move cursor to start",
		"Ctrl+E":        "Move cursor to end",
		"Backspace":     "Delete character before cursor",
		"Alt+Backspace": "Delete previous word (terminal dependent)",
		"←/→":           "Move cursor",
		"Ctrl+←/→":      "Move by word (terminal dependent)",
	}

	for key, desc := range expectedKeybindings {
		if !strings.Contains(output, key) {
			t.Errorf("Expected keybinding '%s' to be present", key)
		}
		if !strings.Contains(output, desc) {
			t.Errorf("Expected description '%s' to be present", desc)
		}
	}

	// Test section headers
	expectedSections := []string{
		"=== Active Key Bindings ===",
		"Interactive Mode Default Bindings:",
		"Custom keybinding configuration:",
		"Use 'ggc debug-keys raw'",
	}

	for _, section := range expectedSections {
		if !strings.Contains(output, section) {
			t.Errorf("Expected section '%s' to be present", section)
		}
	}
}

// Test captureRawKeySequences error conditions
func TestDebugger_captureRawKeySequences_ErrorHandling(t *testing.T) {
	var buf bytes.Buffer
	debugger := &Debugger{
		outputWriter: &buf,
		helper:       NewHelper(),
	}

	testCases := []struct {
		name       string
		outputFile string
		expected   string
	}{
		{
			name:       "no output file",
			outputFile: "",
			expected:   "Error: debug-keys raw mode requires a terminal",
		},
		{
			name:       "with output file",
			outputFile: "test-output.txt",
			expected:   "Error: debug-keys raw mode requires a terminal",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf.Reset()
			debugger.captureRawKeySequences(tc.outputFile)
			output := buf.String()

			if !strings.Contains(output, tc.expected) {
				t.Errorf("Expected '%s', got: %s", tc.expected, output)
			}
		})
	}
}

// Test showDebugKeysHelp content verification
func TestDebugger_showDebugKeysHelp_ContentVerification(t *testing.T) {
	var buf bytes.Buffer
	debugger := &Debugger{
		outputWriter: &buf,
		helper:       NewHelper(),
	}

	debugger.showDebugKeysHelp()
	output := buf.String()

	// Test all required sections are present
	requiredSections := []string{
		"debug-keys - Debug keybinding issues and capture raw key sequences",
		"USAGE:",
		"ggc debug-keys [SUBCOMMAND] [OPTIONS]",
		"SUBCOMMANDS:",
		"(none)          Show currently active key bindings",
		"raw [file]      Capture raw key sequences and optionally save to file",
		"help            Show this help message",
		"EXAMPLES:",
		"ggc debug-keys                 # Show active keybindings",
		"ggc debug-keys raw             # Capture key sequences interactively",
		"ggc debug-keys raw keys.txt    # Capture and save to keys.txt",
		"DESCRIPTION:",
		"The debug-keys command helps troubleshoot keybinding issues by:",
		"1. Showing currently active key bindings",
		"2. Capturing raw key sequences sent by your terminal",
		"3. Identifying common key sequences (arrows, function keys, etc.)",
		"4. Providing the correct format for custom keybinding configuration",
		"Press Ctrl+C in raw mode to stop capturing and view results.",
	}

	for _, section := range requiredSections {
		if !strings.Contains(output, section) {
			t.Errorf("Help text missing required section: '%s'", section)
		}
	}

	// Test that help is properly formatted
	lines := strings.Split(output, "\n")
	if len(lines) < 10 {
		t.Error("Help text should contain multiple lines")
	}
}

// Benchmark tests for performance
func BenchmarkDebugger_DebugKeys_Default(b *testing.B) {
	var buf bytes.Buffer
	debugger := &Debugger{
		outputWriter: &buf,
		helper:       NewHelper(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		debugger.DebugKeys([]string{})
	}
}

func BenchmarkDebugger_DebugKeys_Help(b *testing.B) {
	var buf bytes.Buffer
	debugger := &Debugger{
		outputWriter: &buf,
		helper:       NewHelper(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		debugger.DebugKeys([]string{"help"})
	}
}

// Test concurrent access safety
func TestDebugger_ConcurrentAccess(t *testing.T) {
	debugger := NewDebugger()

	// Test multiple goroutines calling DebugKeys simultaneously
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			debugger.DebugKeys([]string{})
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// Test with different output writers
func TestDebugger_DifferentOutputWriters(t *testing.T) {
	testCases := []struct {
		name   string
		writer *bytes.Buffer
	}{
		{"empty buffer", &bytes.Buffer{}},
		{"pre-filled buffer", bytes.NewBufferString("existing content\n")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			debugger := &Debugger{
				outputWriter: tc.writer,
				helper:       NewHelper(),
			}

			initialLen := tc.writer.Len()
			debugger.DebugKeys([]string{})

			if tc.writer.Len() <= initialLen {
				t.Error("Expected output to be written to buffer")
			}
		})
	}
}
