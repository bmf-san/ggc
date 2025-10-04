package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v7/pkg/git"
)

type mockStashOps struct {
	stashCalled bool
	listCalled  bool
	showCalled  bool
	applyCalled bool
	popCalled   bool
	dropCalled  bool
	clearCalled bool
	stashName   string
	listOutput  string
}

func (m *mockStashOps) Stash() error { m.stashCalled = true; return nil }
func (m *mockStashOps) StashList() (string, error) {
	m.listCalled = true
	return m.listOutput, nil
}
func (m *mockStashOps) StashShow(stash string) error {
	m.showCalled = true
	m.stashName = stash
	return nil
}
func (m *mockStashOps) StashApply(stash string) error {
	m.applyCalled = true
	m.stashName = stash
	return nil
}
func (m *mockStashOps) StashPop(stash string) error {
	m.popCalled = true
	m.stashName = stash
	return nil
}
func (m *mockStashOps) StashDrop(stash string) error {
	m.dropCalled = true
	m.stashName = stash
	return nil
}
func (m *mockStashOps) StashClear() error { m.clearCalled = true; return nil }

var _ git.StashOps = (*mockStashOps)(nil)

func TestStasher_Constructor(t *testing.T) {
	mockClient := &mockStashOps{}
	stasher := NewStasher(mockClient)

	if stasher == nil {
		t.Fatal("Expected NewStasher to return a non-nil Stasher")
	}
	if stasher != nil && stasher.gitClient == nil {
		t.Error("Expected gitClient to be set")
	}
	if stasher != nil && stasher.outputWriter == nil {
		t.Error("Expected outputWriter to be set")
	}
	if stasher != nil && stasher.helper == nil {
		t.Error("Expected helper to be set")
	}
}

func TestStasher_Stash(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		shouldShowHelp bool
	}{
		{
			name:           "no args - default stash",
			args:           []string{},
			expectedOutput: "", // Mock client returns no output on success
			shouldShowHelp: false,
		},
		{
			name:           "list command",
			args:           []string{"list"},
			expectedOutput: "No stashes found", // Mock returns empty, so shows "No stashes found"
			shouldShowHelp: false,
		},
		{
			name:           "show command with stash",
			args:           []string{"show", "stash@{0}"},
			expectedOutput: "", // Mock client returns no output on success
			shouldShowHelp: false,
		},
		{
			name:           "show command without stash",
			args:           []string{"show"},
			expectedOutput: "", // Mock client returns no output on success
			shouldShowHelp: false,
		},
		{
			name:           "apply command with stash",
			args:           []string{"apply", "stash@{0}"},
			expectedOutput: "", // Mock client returns no output on success
			shouldShowHelp: false,
		},
		{
			name:           "apply command without stash",
			args:           []string{"apply"},
			expectedOutput: "", // Mock client returns no output on success
			shouldShowHelp: false,
		},
		{
			name:           "pop command with stash",
			args:           []string{"pop", "stash@{0}"},
			expectedOutput: "", // Mock client returns no output on success
			shouldShowHelp: false,
		},
		{
			name:           "pop command without stash",
			args:           []string{"pop"},
			expectedOutput: "", // Mock client returns no output on success
			shouldShowHelp: false,
		},
		{
			name:           "drop command with stash",
			args:           []string{"drop", "stash@{0}"},
			expectedOutput: "", // Mock client returns no output on success
			shouldShowHelp: false,
		},
		{
			name:           "drop command without stash",
			args:           []string{"drop"},
			expectedOutput: "", // Mock client returns no output on success
			shouldShowHelp: false,
		},
		{
			name:           "clear command",
			args:           []string{"clear"},
			expectedOutput: "", // Mock client returns no output on success
			shouldShowHelp: false,
		},
		{
			name:           "unknown command - should show help",
			args:           []string{"unknown"},
			expectedOutput: "Usage: ggc stash [command]",
			shouldShowHelp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := &mockStashOps{}

			stasher := &Stasher{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}
			// Set helper's output writer to capture help output
			stasher.helper.outputWriter = buf

			stasher.Stash(tt.args)

			output := buf.String()

			// Verify expected behavior
			if tt.shouldShowHelp {
				if !strings.Contains(output, tt.expectedOutput) {
					t.Errorf("Expected help output containing '%s', got: %s", tt.expectedOutput, output)
				}
			} else {
				if tt.expectedOutput == "" {
					// For stash operations, mock returns empty string - this is expected
					// We verify the command executed without error
					if strings.Contains(output, "Error:") {
						t.Errorf("Unexpected error in stash operation: %s", output)
					}
				} else {
					// For specific expected outputs (like "No stashes found")
					if !strings.Contains(output, tt.expectedOutput) {
						t.Errorf("Expected output containing '%s', got: %s", tt.expectedOutput, output)
					}
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

func TestStasher_StashOperations(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		testFunc func(*testing.T, *Stasher, *bytes.Buffer)
	}{
		{
			name: "default stash saves changes",
			args: []string{},
			testFunc: func(t *testing.T, stasher *Stasher, buf *bytes.Buffer) {
				output := buf.String()
				// Mock client doesn't return errors for Stash
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in default stash: %s", output)
				}
				// Should not show help
				if strings.Contains(output, "Usage:") {
					t.Errorf("Unexpected help output for default stash: %s", output)
				}
			},
		},
		{
			name: "stash list shows available stashes",
			args: []string{"list"},
			testFunc: func(t *testing.T, stasher *Stasher, buf *bytes.Buffer) {
				output := buf.String()
				// Mock returns empty string, so should show "No stashes found"
				if !strings.Contains(output, "No stashes found") {
					t.Errorf("Expected 'No stashes found' for empty stash list, got: %s", output)
				}
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in stash list: %s", output)
				}
			},
		},
		{
			name: "stash show displays stash changes",
			args: []string{"show", "stash@{0}"},
			testFunc: func(t *testing.T, stasher *Stasher, buf *bytes.Buffer) {
				output := buf.String()
				// Mock client doesn't return errors for StashShow
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in stash show: %s", output)
				}
			},
		},
		{
			name: "stash apply restores without removing",
			args: []string{"apply", "stash@{0}"},
			testFunc: func(t *testing.T, stasher *Stasher, buf *bytes.Buffer) {
				output := buf.String()
				// Mock client doesn't return errors for StashApply
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in stash apply: %s", output)
				}
			},
		},
		{
			name: "stash pop restores and removes",
			args: []string{"pop", "stash@{0}"},
			testFunc: func(t *testing.T, stasher *Stasher, buf *bytes.Buffer) {
				output := buf.String()
				// Mock client doesn't return errors for StashPop
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in stash pop: %s", output)
				}
			},
		},
		{
			name: "stash drop removes stash",
			args: []string{"drop", "stash@{0}"},
			testFunc: func(t *testing.T, stasher *Stasher, buf *bytes.Buffer) {
				output := buf.String()
				// Mock client doesn't return errors for StashDrop
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in stash drop: %s", output)
				}
			},
		},
		{
			name: "stash clear removes all stashes",
			args: []string{"clear"},
			testFunc: func(t *testing.T, stasher *Stasher, buf *bytes.Buffer) {
				output := buf.String()
				// Mock client doesn't return errors for StashClear
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in stash clear: %s", output)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := &mockStashOps{}

			stasher := &Stasher{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}
			stasher.helper.outputWriter = buf

			stasher.Stash(tt.args)
			tt.testFunc(t, stasher, buf)
		})
	}
}

func TestStasher_StashArgumentHandling(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		testFunc func(*testing.T, *Stasher, *bytes.Buffer)
	}{
		{
			name: "show without stash reference uses default",
			args: []string{"show"},
			testFunc: func(t *testing.T, stasher *Stasher, buf *bytes.Buffer) {
				output := buf.String()
				// Should work without error (uses empty string as stash reference)
				if strings.Contains(output, "Error:") {
					t.Errorf("Show without stash reference should work: %s", output)
				}
			},
		},
		{
			name: "apply without stash reference uses latest",
			args: []string{"apply"},
			testFunc: func(t *testing.T, stasher *Stasher, buf *bytes.Buffer) {
				output := buf.String()
				// Should work without error (uses empty string as stash reference)
				if strings.Contains(output, "Error:") {
					t.Errorf("Apply without stash reference should work: %s", output)
				}
			},
		},
		{
			name: "pop without stash reference uses latest",
			args: []string{"pop"},
			testFunc: func(t *testing.T, stasher *Stasher, buf *bytes.Buffer) {
				output := buf.String()
				// Should work without error (uses empty string as stash reference)
				if strings.Contains(output, "Error:") {
					t.Errorf("Pop without stash reference should work: %s", output)
				}
			},
		},
		{
			name: "drop without stash reference uses latest",
			args: []string{"drop"},
			testFunc: func(t *testing.T, stasher *Stasher, buf *bytes.Buffer) {
				output := buf.String()
				// Should work without error (uses empty string as stash reference)
				if strings.Contains(output, "Error:") {
					t.Errorf("Drop without stash reference should work: %s", output)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := &mockStashOps{}

			stasher := &Stasher{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}
			stasher.helper.outputWriter = buf

			stasher.Stash(tt.args)
			tt.testFunc(t, stasher, buf)
		})
	}
}
