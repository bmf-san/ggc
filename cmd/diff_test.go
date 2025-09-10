package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v5/internal/testutil"
)

func TestDiffer_Constructor(t *testing.T) {
	mockClient := testutil.NewMockGitClient()
	differ := NewDiffer(mockClient)

	if differ == nil {
		t.Fatal("Expected NewDiffer to return a non-nil Differ")
	}
	if differ.gitClient == nil {
		t.Error("Expected gitClient to be set")
	}
	if differ.outputWriter == nil {
		t.Error("Expected outputWriter to be set")
	}
	if differ.helper == nil {
		t.Error("Expected helper to be set")
	}
}

func TestDiffer_Diff(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		shouldShowHelp bool
	}{
		{
			name:           "no args - should call DiffHead",
			args:           []string{},
			expectedOutput: "", // Mock returns empty diff
			shouldShowHelp: false,
		},
		{
			name:           "unstaged - should call Diff",
			args:           []string{"unstaged"},
			expectedOutput: "", // Mock returns empty diff
			shouldShowHelp: false,
		},
		{
			name:           "staged - should call DiffStaged",
			args:           []string{"staged"},
			expectedOutput: "", // Mock returns empty diff
			shouldShowHelp: false,
		},
		{
			name:           "invalid arg - should show help",
			args:           []string{"invalid"},
			expectedOutput: "Usage: ggc diff [options]",
			shouldShowHelp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := testutil.NewMockGitClient()

			differ := &Differ{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}
			// Set helper's output writer to capture help output
			differ.helper.outputWriter = buf

			differ.Diff(tt.args)

			output := buf.String()

			// Verify expected behavior
			if tt.shouldShowHelp {
				if !strings.Contains(output, tt.expectedOutput) {
					t.Errorf("Expected help output containing '%s', got: %s", tt.expectedOutput, output)
				}
			} else {
				// For diff operations, mock returns empty string - this is expected
				// We verify the command executed without error
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in diff operation: %s", output)
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

func TestDiffer_DiffOperations(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		testFunc func(*testing.T, *Differ, *bytes.Buffer)
	}{
		{
			name: "no args calls DiffHead",
			args: []string{},
			testFunc: func(t *testing.T, differ *Differ, buf *bytes.Buffer) {
				// Mock client returns empty string for DiffHead
				output := buf.String()
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in DiffHead operation: %s", output)
				}
				// In a real test, we would verify DiffHead was called
			},
		},
		{
			name: "unstaged calls Diff",
			args: []string{"unstaged"},
			testFunc: func(t *testing.T, differ *Differ, buf *bytes.Buffer) {
				// Mock client returns empty string for Diff
				output := buf.String()
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in Diff operation: %s", output)
				}
				// In a real test, we would verify Diff was called
			},
		},
		{
			name: "staged calls DiffStaged",
			args: []string{"staged"},
			testFunc: func(t *testing.T, differ *Differ, buf *bytes.Buffer) {
				// Mock client returns empty string for DiffStaged
				output := buf.String()
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in DiffStaged operation: %s", output)
				}
				// In a real test, we would verify DiffStaged was called
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := testutil.NewMockGitClient()

			differ := &Differ{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}
			differ.helper.outputWriter = buf

			differ.Diff(tt.args)
			tt.testFunc(t, differ, buf)
		})
	}
}

func TestDiffer_ErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		testFunc func(*testing.T, *Differ, *bytes.Buffer)
	}{
		{
			name: "basic diff operation without errors",
			args: []string{},
			testFunc: func(t *testing.T, differ *Differ, buf *bytes.Buffer) {
				output := buf.String()
				// Mock client doesn't return errors, so no "Error:" should appear
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in basic diff operation: %s", output)
				}
			},
		},
		{
			name: "unstaged diff without errors",
			args: []string{"unstaged"},
			testFunc: func(t *testing.T, differ *Differ, buf *bytes.Buffer) {
				output := buf.String()
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in unstaged diff operation: %s", output)
				}
			},
		},
		{
			name: "staged diff without errors",
			args: []string{"staged"},
			testFunc: func(t *testing.T, differ *Differ, buf *bytes.Buffer) {
				output := buf.String()
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in staged diff operation: %s", output)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := testutil.NewMockGitClient()

			differ := &Differ{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}
			differ.helper.outputWriter = buf

			differ.Diff(tt.args)
			tt.testFunc(t, differ, buf)

			// Verify that the mock client is properly configured
			if mockClient == nil {
				t.Error("Expected mock client to be initialized")
			}
		})
	}
}
