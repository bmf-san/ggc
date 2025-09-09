package cmd

import (
	"bytes"
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
		shouldCallHelp bool
	}{
		{
			name:           "no args - should call DiffHead",
			args:           []string{},
			expectedOutput: "diff output from DiffHead",
			shouldCallHelp: false,
		},
		{
			name:           "unstaged - should call Diff",
			args:           []string{"unstaged"},
			expectedOutput: "diff output from Diff",
			shouldCallHelp: false,
		},
		{
			name:           "staged - should call DiffStaged",
			args:           []string{"staged"},
			expectedOutput: "diff output from DiffStaged",
			shouldCallHelp: false,
		},
		{
			name:           "invalid arg - should show help",
			args:           []string{"invalid"},
			expectedOutput: "",
			shouldCallHelp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create buffer to capture output
			buf := &bytes.Buffer{}

			// Use testutil mock client
			mockClient := testutil.NewMockGitClient()

			differ := &Differ{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}

			differ.Diff(tt.args)

			// Basic test - just ensure no panic occurs
			// In a real test, we would check specific outputs based on mock responses
		})
	}
}

// Basic error test using testutil mock client
func TestDiffer_DiffBasic(t *testing.T) {
	buf := &bytes.Buffer{}
	mockClient := testutil.NewMockGitClient()

	differ := &Differ{
		gitClient:    mockClient,
		outputWriter: buf,
		helper:       NewHelper(),
	}

	// Test basic functionality
	differ.Diff([]string{})
	// Should not panic and should produce some output or call methods
}
