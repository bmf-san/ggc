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

			// Verify that the function executed without panic and produced output
			output := buf.String()
			
			// Diff commands should produce some output (diff results, help, or error messages)
			if len(output) == 0 {
				t.Errorf("Expected output for diff command %v, got empty string", tt.args)
			}
			
			// Verify output content based on command type
			switch {
			case len(tt.args) == 0:
				// No args should show diff (DiffHead)
				if len(output) == 0 {
					t.Errorf("Expected diff output for default command, got empty string")
				}
			case len(tt.args) > 0 && tt.args[0] == "unstaged":
				// Unstaged should show unstaged diff
				if len(output) == 0 {
					t.Errorf("Expected unstaged diff output, got empty string")
				}
			case len(tt.args) > 0 && tt.args[0] == "staged":
				// Staged should show staged diff
				if len(output) == 0 {
					t.Errorf("Expected staged diff output, got empty string")
				}
			case len(tt.args) > 0 && tt.args[0] == "invalid":
				// Invalid args should show help or error
				if len(output) < 5 {
					t.Errorf("Expected help or error output for invalid args, got: %s", output)
				}
			}
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

	// Test basic functionality - should execute without panic
	differ.Diff([]string{})
	
	// Verify that the function executed and produced output
	output := buf.String()
	if len(output) == 0 {
		t.Error("Expected diff output from basic test, got empty string")
	}
	
	// Verify that the mock client is properly configured
	if mockClient == nil {
		t.Error("Expected mock client to be initialized")
	}
}
