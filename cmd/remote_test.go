package cmd

import (
	"bytes"
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
		name string
		args []string
	}{
		{
			name: "no args - should show help",
			args: []string{},
		},
		{
			name: "list command",
			args: []string{"list"},
		},
		{
			name: "add command with correct args",
			args: []string{"add", "origin", "https://github.com/user/repo.git"},
		},
		{
			name: "add command with incorrect args",
			args: []string{"add", "origin"},
		},
		{
			name: "remove command with correct args",
			args: []string{"remove", "origin"},
		},
		{
			name: "remove command with incorrect args",
			args: []string{"remove"},
		},
		{
			name: "set-url command with correct args",
			args: []string{"set-url", "origin", "https://github.com/user/newrepo.git"},
		},
		{
			name: "set-url command with incorrect args",
			args: []string{"set-url", "origin"},
		},
		{
			name: "unknown command",
			args: []string{"unknown"},
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

			remoter.Remote(tt.args)

			// Verify that the function executed without panic and produced output
			output := buf.String()

			// Remote commands should produce some output (help, results, or error messages)
			if len(output) == 0 {
				t.Errorf("Expected output for remote command %v, got empty string", tt.args)
			}

			// Verify output content based on command and arguments
			switch {
			case len(tt.args) == 0:
				// No args should show help
				if len(output) < 10 {
					t.Errorf("Expected help output for no args, got: %s", output)
				}
			case len(tt.args) > 0 && tt.args[0] == "list":
				// List command should show remote information
				if len(output) == 0 {
					t.Errorf("Expected list output, got empty string")
				}
			case len(tt.args) > 0 && tt.args[0] == "unknown":
				// Unknown commands should show error or help
				if len(output) < 5 {
					t.Errorf("Expected error output for unknown command, got: %s", output)
				}
			case len(tt.args) > 0 && (tt.args[0] == "add" || tt.args[0] == "remove" || tt.args[0] == "set-url"):
				// Commands with arguments should show result or error
				if len(output) == 0 {
					t.Errorf("Expected output for %s command, got empty string", tt.args[0])
				}
				// Commands with insufficient args should show error or help
				if (tt.args[0] == "add" && len(tt.args) < 3) ||
					(tt.args[0] == "remove" && len(tt.args) < 2) ||
					(tt.args[0] == "set-url" && len(tt.args) < 3) {
					// These should produce error or help output
					if len(output) < 5 {
						t.Errorf("Expected error or help for insufficient args in %s, got: %s", tt.args[0], output)
					}
				}
			}
		})
	}
}
