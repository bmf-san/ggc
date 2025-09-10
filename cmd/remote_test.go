package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v5/internal/testutil"
)

func TestRemoter_Constructor(t *testing.T) {
	mockClient := testutil.NewMockGitClient()
	remoter := NewRemoter(mockClient)

	if remoter == nil {
		t.Fatal("Expected NewRemoter to return a non-nil Remoter")
	}
	if remoter.gitClient == nil {
		t.Error("Expected gitClient to be set")
	}
	if remoter.outputWriter == nil {
		t.Error("Expected outputWriter to be set")
	}
	if remoter.helper == nil {
		t.Error("Expected helper to be set")
	}
}

func TestRemoter_Remote(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		shouldShowHelp bool
	}{
		{
			name:           "no args - should show help",
			args:           []string{},
			expectedOutput: "Usage: ggc remote <command>",
			shouldShowHelp: true,
		},
		{
			name:           "list command",
			args:           []string{"list"},
			expectedOutput: "",
			shouldShowHelp: false,
		},
		{
			name:           "add command with correct args",
			args:           []string{"add", "origin", "https://github.com/user/repo.git"},
			expectedOutput: "Remote 'origin' added",
			shouldShowHelp: false,
		},
		{
			name:           "add command with incorrect args",
			args:           []string{"add", "origin"},
			expectedOutput: "Usage: ggc remote <command>",
			shouldShowHelp: true,
		},
		{
			name:           "remove command with correct args",
			args:           []string{"remove", "origin"},
			expectedOutput: "Remote 'origin' removed",
			shouldShowHelp: false,
		},
		{
			name:           "remove command with incorrect args",
			args:           []string{"remove"},
			expectedOutput: "Usage: ggc remote <command>",
			shouldShowHelp: true,
		},
		{
			name:           "set-url command with correct args",
			args:           []string{"set-url", "origin", "https://github.com/user/newrepo.git"},
			expectedOutput: "Remote 'origin' URL updated",
			shouldShowHelp: false,
		},
		{
			name:           "set-url command with incorrect args",
			args:           []string{"set-url", "origin"},
			expectedOutput: "Usage: ggc remote <command>",
			shouldShowHelp: true,
		},
		{
			name:           "unknown command",
			args:           []string{"unknown"},
			expectedOutput: "Usage: ggc remote <command>",
			shouldShowHelp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := testutil.NewMockGitClient()

			remoter := &Remoter{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}
			// Set helper's output writer to capture help output
			remoter.helper.outputWriter = buf

			remoter.Remote(tt.args)

			output := buf.String()

			// Verify expected behavior
			if tt.shouldShowHelp {
				if !strings.Contains(output, tt.expectedOutput) {
					t.Errorf("Expected help output containing '%s', got: %s", tt.expectedOutput, output)
				}
			} else if tt.expectedOutput != "" {
				if !strings.Contains(output, tt.expectedOutput) {
					t.Errorf("Expected output containing '%s', got: %s", tt.expectedOutput, output)
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

func TestRemoter_RemoteOperations(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		testFunc func(*testing.T, *Remoter, *bytes.Buffer)
	}{
		{
			name: "list operation calls git client",
			args: []string{"list"},
			testFunc: func(t *testing.T, remoter *Remoter, buf *bytes.Buffer) {
				// Test that RemoteList is called (mock doesn't return error)
				// In a more sophisticated test, we'd verify the git client method was called
				if buf.String() != "" && strings.Contains(buf.String(), "Error:") {
					t.Errorf("Unexpected error in list operation: %s", buf.String())
				}
			},
		},
		{
			name: "add operation with success",
			args: []string{"add", "upstream", "https://github.com/upstream/repo.git"},
			testFunc: func(t *testing.T, remoter *Remoter, buf *bytes.Buffer) {
				output := buf.String()
				if !strings.Contains(output, "Remote 'upstream' added") {
					t.Errorf("Expected success message for add operation, got: %s", output)
				}
			},
		},
		{
			name: "remove operation with success",
			args: []string{"remove", "upstream"},
			testFunc: func(t *testing.T, remoter *Remoter, buf *bytes.Buffer) {
				output := buf.String()
				if !strings.Contains(output, "Remote 'upstream' removed") {
					t.Errorf("Expected success message for remove operation, got: %s", output)
				}
			},
		},
		{
			name: "set-url operation with success",
			args: []string{"set-url", "origin", "https://github.com/newowner/repo.git"},
			testFunc: func(t *testing.T, remoter *Remoter, buf *bytes.Buffer) {
				output := buf.String()
				if !strings.Contains(output, "Remote 'origin' URL updated") {
					t.Errorf("Expected success message for set-url operation, got: %s", output)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := testutil.NewMockGitClient()

			remoter := &Remoter{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}
			remoter.helper.outputWriter = buf

			remoter.Remote(tt.args)
			tt.testFunc(t, remoter, buf)
		})
	}
}
