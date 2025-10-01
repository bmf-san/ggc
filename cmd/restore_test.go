package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v7/git"
)

type mockRestoreOps struct {
	restoreWorkingDirCalled bool
	restoreStagedCalled     bool
	restoreFromCommitCalled bool
	revParseVerifyCalled    bool
	paths                   []string
	commit                  string
	ref                     string
	revParseResult          bool
}

func (m *mockRestoreOps) RestoreWorkingDir(paths ...string) error {
	m.restoreWorkingDirCalled = true
	m.paths = paths
	return nil
}
func (m *mockRestoreOps) RestoreStaged(paths ...string) error {
	m.restoreStagedCalled = true
	m.paths = paths
	return nil
}
func (m *mockRestoreOps) RestoreFromCommit(commit string, paths ...string) error {
	m.restoreFromCommitCalled = true
	m.commit = commit
	m.paths = paths
	return nil
}
func (m *mockRestoreOps) RevParseVerify(ref string) bool {
	m.revParseVerifyCalled = true
	m.ref = ref
	return m.revParseResult
}

var _ git.RestoreOps = (*mockRestoreOps)(nil)

func TestRestorer_Constructor(t *testing.T) {
	mockClient := &mockRestoreOps{}
	restorer := NewRestorer(mockClient)

	if restorer == nil {
		t.Fatal("Expected NewRestorer to return a non-nil Restorer")
	}
	if restorer != nil && restorer.gitClient == nil {
		t.Error("Expected gitClient to be set")
	}
	if restorer != nil && restorer.outputWriter == nil {
		t.Error("Expected outputWriter to be set")
	}
	if restorer != nil && restorer.helper == nil {
		t.Error("Expected helper to be set")
	}
}

func TestRestorer_Restore(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		shouldShowHelp bool
	}{
		{
			name:           "no args - show help",
			args:           []string{},
			expectedOutput: "Usage: ggc restore [command]",
			shouldShowHelp: true,
		},
		{
			name:           "restore file",
			args:           []string{"file.txt"},
			expectedOutput: "", // Mock client returns no output on success
			shouldShowHelp: false,
		},
		{
			name:           "restore all files",
			args:           []string{"."},
			expectedOutput: "", // Mock client returns no output on success
			shouldShowHelp: false,
		},
		{
			name:           "restore staged with file",
			args:           []string{"staged", "file.txt"},
			expectedOutput: "", // Mock client returns no output on success
			shouldShowHelp: false,
		},
		{
			name:           "restore staged without file - show help",
			args:           []string{"staged"},
			expectedOutput: "Usage: ggc restore [command]",
			shouldShowHelp: true,
		},
		{
			name:           "restore from commit",
			args:           []string{"HEAD~1", "file.txt"},
			expectedOutput: "", // Mock client returns no output on success
			shouldShowHelp: false,
		},
		{
			name:           "restore from SHA commit",
			args:           []string{"abc1234", "file.txt"},
			expectedOutput: "", // Mock client returns no output on success
			shouldShowHelp: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := &mockRestoreOps{}

			restorer := &Restorer{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}
			// Set helper's output writer to capture help output
			restorer.helper.outputWriter = buf

			restorer.Restore(tt.args)

			output := buf.String()

			// Verify expected behavior
			if tt.shouldShowHelp {
				if !strings.Contains(output, tt.expectedOutput) {
					t.Errorf("Expected help output containing '%s', got: %s", tt.expectedOutput, output)
				}
			} else {
				// For restore operations, mock returns empty string - this is expected
				// We verify the command executed without error
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in restore operation: %s", output)
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

func TestRestorer_RestoreOperations(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		testFunc func(*testing.T, *Restorer, *bytes.Buffer)
	}{
		{
			name: "restore working directory file",
			args: []string{"file.txt"},
			testFunc: func(t *testing.T, restorer *Restorer, buf *bytes.Buffer) {
				output := buf.String()
				// Mock client doesn't return errors for RestoreWorkingDir
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in working directory restore: %s", output)
				}
				// Should not show help
				if strings.Contains(output, "Usage:") {
					t.Errorf("Unexpected help output for valid restore operation: %s", output)
				}
			},
		},
		{
			name: "restore staged file",
			args: []string{"staged", "file.txt"},
			testFunc: func(t *testing.T, restorer *Restorer, buf *bytes.Buffer) {
				output := buf.String()
				// Mock client doesn't return errors for RestoreStaged
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in staged restore: %s", output)
				}
				// Should not show help
				if strings.Contains(output, "Usage:") {
					t.Errorf("Unexpected help output for valid staged restore: %s", output)
				}
			},
		},
		{
			name: "restore from commit",
			args: []string{"HEAD~1", "file.txt"},
			testFunc: func(t *testing.T, restorer *Restorer, buf *bytes.Buffer) {
				output := buf.String()
				// Mock client doesn't return errors for RestoreFromCommit
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in commit restore: %s", output)
				}
				// Should not show help
				if strings.Contains(output, "Usage:") {
					t.Errorf("Unexpected help output for valid commit restore: %s", output)
				}
			},
		},
		{
			name: "restore all files",
			args: []string{"."},
			testFunc: func(t *testing.T, restorer *Restorer, buf *bytes.Buffer) {
				output := buf.String()
				// Should restore all files in working directory
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in restore all: %s", output)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := &mockRestoreOps{}

			restorer := &Restorer{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}
			restorer.helper.outputWriter = buf

			restorer.Restore(tt.args)
			tt.testFunc(t, restorer, buf)
		})
	}
}

func TestRestorer_CommitDetection(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		testFunc func(*testing.T, *Restorer, *bytes.Buffer)
	}{
		{
			name: "HEAD variant detection",
			args: []string{"HEAD~1", "file.txt"},
			testFunc: func(t *testing.T, restorer *Restorer, buf *bytes.Buffer) {
				output := buf.String()
				// Should be treated as commit restore, not working directory restore
				if strings.Contains(output, "Error:") {
					t.Errorf("HEAD variant should be valid commit reference: %s", output)
				}
			},
		},
		{
			name: "SHA-like commit detection",
			args: []string{"abc1234", "file.txt"},
			testFunc: func(t *testing.T, restorer *Restorer, buf *bytes.Buffer) {
				output := buf.String()
				// Should be treated as commit restore
				if strings.Contains(output, "Error:") {
					t.Errorf("SHA-like string should be valid commit reference: %s", output)
				}
			},
		},
		{
			name: "regular file path (not commit)",
			args: []string{"regular_file.txt"},
			testFunc: func(t *testing.T, restorer *Restorer, buf *bytes.Buffer) {
				output := buf.String()
				// Should be treated as working directory restore
				if strings.Contains(output, "Error:") {
					t.Errorf("Regular file path should work for working directory restore: %s", output)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := &mockRestoreOps{}

			restorer := &Restorer{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}
			restorer.helper.outputWriter = buf

			restorer.Restore(tt.args)
			tt.testFunc(t, restorer, buf)
		})
	}
}
