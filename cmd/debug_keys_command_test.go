package cmd

import (
	"os"
	"strings"
	"testing"
	"time"
)

// Test DebugKeysCommand integration with Debugger
func TestDebugKeysCommand_Integration(t *testing.T) {
	testCases := []struct {
		name       string
		outputFile string
	}{
		{
			name:       "without output file",
			outputFile: "",
		},
		{
			name:       "with output file",
			outputFile: "test_sequences.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewDebugKeysCommand(tc.outputFile)

			// Test initial state
			if cmd.IsCapturing() {
				t.Error("Command should not be capturing initially")
			}

			// Test StartCapture
			cmd.StartCapture()
			if !cmd.IsCapturing() {
				t.Error("Command should be capturing after StartCapture")
			}

			// Test CaptureSequence
			testSequence := []byte{27, 91, 65} // Up arrow key
			cmd.CaptureSequence(testSequence)

			// Test StopCapture
			err := cmd.StopCapture()
			if err != nil {
				t.Errorf("StopCapture should not return error: %v", err)
			}

			if cmd.IsCapturing() {
				t.Error("Command should not be capturing after StopCapture")
			}

			// Cleanup test file if created
			if tc.outputFile != "" {
				_ = os.Remove(tc.outputFile)
			}
		})
	}
}

// Test DebugKeysCommand sequence identification
func TestDebugKeysCommand_SequenceIdentification(t *testing.T) {
	cmd := NewDebugKeysCommand("")

	testCases := []struct {
		name     string
		sequence []byte
		expected string
	}{
		{
			name:     "up arrow",
			sequence: []byte{27, 91, 65},
			expected: "↑",
		},
		{
			name:     "down arrow",
			sequence: []byte{27, 91, 66},
			expected: "↓",
		},
		{
			name:     "right arrow",
			sequence: []byte{27, 91, 67},
			expected: "→",
		},
		{
			name:     "left arrow",
			sequence: []byte{27, 91, 68},
			expected: "←",
		},
		{
			name:     "tab key",
			sequence: []byte{9},
			expected: "Tab",
		},
		{
			name:     "enter key",
			sequence: []byte{13},
			expected: "Enter",
		},
		{
			name:     "escape key",
			sequence: []byte{27},
			expected: "Esc",
		},
		{
			name:     "space key",
			sequence: []byte{32},
			expected: "Space",
		},
		{
			name:     "ctrl+a",
			sequence: []byte{1},
			expected: "Ctrl+A",
		},
		{
			name:     "ctrl+z",
			sequence: []byte{26},
			expected: "Ctrl+Z",
		},
		{
			name:     "F1 key",
			sequence: []byte{27, 79, 80},
			expected: "F1",
		},
		{
			name:     "F4 key",
			sequence: []byte{27, 79, 83},
			expected: "F4",
		},
		{
			name:     "shift+up arrow",
			sequence: []byte{27, 91, 49, 59, 50, 65},
			expected: "Shift+↑",
		},
		{
			name:     "shift+left arrow",
			sequence: []byte{27, 91, 49, 59, 50, 68},
			expected: "Shift+←",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd.StartCapture()
			cmd.CaptureSequence(tc.sequence)
			err := cmd.StopCapture()
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// Test DebugKeysCommand format key sequence comprehensive
func TestDebugKeysCommand_FormatKeySequence_Comprehensive(t *testing.T) {
	cmd := NewDebugKeysCommand("")

	testCases := []struct {
		name     string
		sequence []byte
		contains []string // Strings that should be in the formatted output
	}{
		{
			name:     "empty sequence",
			sequence: []byte{},
			contains: []string{"(empty)"},
		},
		{
			name:     "single printable character",
			sequence: []byte{'a'},
			contains: []string{"a", "0x61"},
		},
		{
			name:     "control character",
			sequence: []byte{1},
			contains: []string{"Ctrl+A", "0x01"},
		},
		{
			name:     "escape sequence",
			sequence: []byte{27, 91, 65},
			contains: []string{"↑", "0x1b", "0x5b", "0x41"},
		},
		{
			name:     "multi-byte sequence",
			sequence: []byte{195, 169}, // é in UTF-8
			contains: []string{"0xc3", "0xa9"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			formatted := cmd.formatKeySequence(tc.sequence)

			for _, expected := range tc.contains {
				if !strings.Contains(formatted, expected) {
					t.Errorf("Expected formatted sequence to contain '%s', got: %s", expected, formatted)
				}
			}
		})
	}
}

// Test DebugKeysCommand file output
func TestDebugKeysCommand_FileOutput(t *testing.T) {
	outputFile := "test_debug_output.txt"
	defer func() {
		_ = os.Remove(outputFile) // Cleanup - ignore error as file may not exist
	}()

	cmd := NewDebugKeysCommand(outputFile)

	// Capture some test sequences
	cmd.StartCapture()
	cmd.CaptureSequence([]byte{27, 91, 65}) // Up arrow
	cmd.CaptureSequence([]byte{9})          // Tab
	cmd.CaptureSequence([]byte{13})         // Enter

	err := cmd.StopCapture()
	if err != nil {
		t.Fatalf("StopCapture failed: %v", err)
	}

	// Check that file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Output file was not created")
		return
	}

	// Read and verify file contents
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)

	// Verify expected content
	expectedContent := []string{
		"# Raw Key Sequences Captured by ggc debug-keys",
		"# Total sequences: 3",
		"# Sequence 1",
		"# Raw: [27 91 65]",
		"# Hex: 1b5b41",
		"# Identified: ↑",
		"raw:1b5b41",
		"# Sequence 2",
		"# Raw: [9]",
		"# Hex: 09",
		"# Identified: Tab",
		"raw:09",
		"# Sequence 3",
		"# Raw: [13]",
		"# Hex: 0d",
		"# Identified: Enter",
		"raw:0d",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Expected file to contain '%s'", expected)
		}
	}
}

// Test DebugKeysCommand capture when not started
func TestDebugKeysCommand_CaptureWithoutStart(t *testing.T) {
	cmd := NewDebugKeysCommand("")

	// Try to capture without starting
	cmd.CaptureSequence([]byte{65}) // 'A'

	// Should not be capturing
	if cmd.IsCapturing() {
		t.Error("Should not be capturing without StartCapture")
	}

	// StopCapture should be safe to call
	err := cmd.StopCapture()
	if err != nil {
		t.Errorf("StopCapture should not error when not capturing: %v", err)
	}
}

// Test DebugKeysCommand with invalid file path
func TestDebugKeysCommand_InvalidFilePath(t *testing.T) {
	// Use an invalid path (directory that doesn't exist)
	cmd := NewDebugKeysCommand("/nonexistent/directory/output.txt")

	cmd.StartCapture()
	cmd.CaptureSequence([]byte{65}) // 'A'

	err := cmd.StopCapture()
	if err == nil {
		t.Error("Expected error when writing to invalid file path")
	}
}

// Test DebugKeysCommand multiple start/stop cycles
func TestDebugKeysCommand_MultipleStartStop(t *testing.T) {
	cmd := NewDebugKeysCommand("")

	// First cycle
	cmd.StartCapture()
	if !cmd.IsCapturing() {
		t.Error("Should be capturing after first StartCapture")
	}

	cmd.CaptureSequence([]byte{65}) // 'A'

	err := cmd.StopCapture()
	if err != nil {
		t.Errorf("First StopCapture failed: %v", err)
	}

	if cmd.IsCapturing() {
		t.Error("Should not be capturing after first StopCapture")
	}

	// Second cycle
	cmd.StartCapture()
	if !cmd.IsCapturing() {
		t.Error("Should be capturing after second StartCapture")
	}

	cmd.CaptureSequence([]byte{66}) // 'B'

	err = cmd.StopCapture()
	if err != nil {
		t.Errorf("Second StopCapture failed: %v", err)
	}
}

// Benchmark DebugKeysCommand operations
func BenchmarkDebugKeysCommand_CaptureSequence(b *testing.B) {
	cmd := NewDebugKeysCommand("")
	cmd.StartCapture()
	testSeq := []byte{27, 91, 65} // Up arrow

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd.CaptureSequence(testSeq)
	}
}

func BenchmarkDebugKeysCommand_IdentifySequence(b *testing.B) {
	cmd := NewDebugKeysCommand("")
	testSeq := []byte{27, 91, 65} // Up arrow

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd.identifySequence(testSeq)
	}
}

// Test concurrent access to DebugKeysCommand
func TestDebugKeysCommand_ConcurrentAccess(t *testing.T) {
	cmd := NewDebugKeysCommand("")
	cmd.StartCapture()

	done := make(chan bool, 10)

	// Start multiple goroutines capturing sequences
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()
			seq := []byte{byte(65 + id)} // Different sequence for each goroutine
			cmd.CaptureSequence(seq)
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("Timeout waiting for concurrent captures")
		}
	}

	err := cmd.StopCapture()
	if err != nil {
		t.Errorf("StopCapture failed after concurrent access: %v", err)
	}
}
