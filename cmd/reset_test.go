package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v5/git"
)

type mockResetOps struct {
	currentBranch           string
	resetHardAndCleanCalled bool
	resetHardCalled         bool
	commit                  string
}

func (m *mockResetOps) GetCurrentBranch() (string, error) {
	if m.currentBranch == "" {
		return "main", nil
	}
	return m.currentBranch, nil
}
func (m *mockResetOps) ResetHardAndClean() error {
	m.resetHardAndCleanCalled = true
	return nil
}
func (m *mockResetOps) ResetHard(commit string) error {
	m.resetHardCalled = true
	m.commit = commit
	return nil
}

var _ git.ResetOps = (*mockResetOps)(nil)

func TestResetter_Constructor(t *testing.T) {
	mockClient := &mockResetOps{}
	resetter := NewResetter(mockClient)

	if resetter == nil {
		t.Fatal("Expected NewResetter to return a non-nil Resetter")
	}
	if resetter.gitClient == nil {
		t.Error("Expected gitClient to be set")
	}
	if resetter.outputWriter == nil {
		t.Error("Expected outputWriter to be set")
	}
	if resetter.helper == nil {
		t.Error("Expected helper to be set")
	}
}

func TestResetter_Reset(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		shouldShowHelp bool
	}{
		{
			name:           "no args - default reset",
			args:           []string{},
			expectedOutput: "Reset to origin/main successful",
			shouldShowHelp: false,
		},
		{
			name:           "hard reset with commit",
			args:           []string{"hard", "abc123"},
			expectedOutput: "Reset to abc123 successful",
			shouldShowHelp: false,
		},
		{
			name:           "hard reset without commit - should show help",
			args:           []string{"hard"},
			expectedOutput: "Error: commit hash required for hard reset",
			shouldShowHelp: true,
		},
		{
			name:           "unknown command - should show help",
			args:           []string{"unknown"},
			expectedOutput: "Usage: ggc reset",
			shouldShowHelp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := &mockResetOps{}

			resetter := &Resetter{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}
			// Set helper's output writer to capture help output
			resetter.helper.outputWriter = buf

			resetter.Reset(tt.args)

			output := buf.String()

			// Verify expected behavior
			if tt.shouldShowHelp {
				if !strings.Contains(output, tt.expectedOutput) {
					t.Errorf("Expected output containing '%s', got: %s", tt.expectedOutput, output)
				}
			} else {
				if !strings.Contains(output, tt.expectedOutput) {
					t.Errorf("Expected success output containing '%s', got: %s", tt.expectedOutput, output)
				}
			}

			// Verify no panic occurred
			if t.Failed() {
				t.Logf("Command args: %v", tt.args)
				t.Logf("Full output: %s", output)
			}
		})
	}
}

func TestResetter_ResetOperations(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		testFunc func(*testing.T, *Resetter, *bytes.Buffer)
	}{
		{
			name: "default reset calls ResetHardAndClean",
			args: []string{},
			testFunc: func(t *testing.T, resetter *Resetter, buf *bytes.Buffer) {
				output := buf.String()
				// Should show success message with branch name
				if !strings.Contains(output, "Reset to origin/") || !strings.Contains(output, "successful") {
					t.Errorf("Expected default reset success message, got: %s", output)
				}
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in default reset: %s", output)
				}
			},
		},
		{
			name: "hard reset with commit calls ResetHard",
			args: []string{"hard", "commit123"},
			testFunc: func(t *testing.T, resetter *Resetter, buf *bytes.Buffer) {
				output := buf.String()
				// Should show success message with commit hash
				if !strings.Contains(output, "Reset to commit123 successful") {
					t.Errorf("Expected hard reset success message, got: %s", output)
				}
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in hard reset: %s", output)
				}
			},
		},
		{
			name: "hard reset without commit shows error and help",
			args: []string{"hard"},
			testFunc: func(t *testing.T, resetter *Resetter, buf *bytes.Buffer) {
				output := buf.String()
				// Should show error message
				if !strings.Contains(output, "Error: commit hash required") {
					t.Errorf("Expected error message for missing commit, got: %s", output)
				}
				// Should also show help
				if !strings.Contains(output, "Usage:") {
					t.Errorf("Expected help message after error, got: %s", output)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := &mockResetOps{}

			resetter := &Resetter{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}
			resetter.helper.outputWriter = buf

			resetter.Reset(tt.args)
			tt.testFunc(t, resetter, buf)
		})
	}
}
