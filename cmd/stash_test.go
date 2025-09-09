package cmd

import (
	"bytes"
	"testing"

	"github.com/bmf-san/ggc/v5/internal/testutil"
)

func TestStasher_Constructor(t *testing.T) {
	mockClient := testutil.NewMockGitClient()
	stasher := NewStasher(mockClient)

	if stasher == nil {
		t.Fatal("Expected NewStasher to return a non-nil Stasher")
	}
	if stasher.gitClient == nil {
		t.Error("Expected gitClient to be set")
	}
	if stasher.outputWriter == nil {
		t.Error("Expected outputWriter to be set")
	}
	if stasher.helper == nil {
		t.Error("Expected helper to be set")
	}
}

func TestStasher_Stash(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "no args - default stash",
			args: []string{},
		},
		{
			name: "list command",
			args: []string{"list"},
		},
		{
			name: "show command",
			args: []string{"show", "stash@{0}"},
		},
		{
			name: "apply command",
			args: []string{"apply", "stash@{0}"},
		},
		{
			name: "pop command",
			args: []string{"pop", "stash@{0}"},
		},
		{
			name: "drop command",
			args: []string{"drop", "stash@{0}"},
		},
		{
			name: "clear command",
			args: []string{"clear"},
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

			stasher := &Stasher{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}

			stasher.Stash(tt.args)

			// Verify that the function executed without panic and produced output
			output := buf.String()
			
			// Stash commands should produce some output (results, help, or error messages)
			if len(output) == 0 {
				t.Errorf("Expected output for stash command %v, got empty string", tt.args)
			}
			
			// Verify output content based on command type
			switch {
			case len(tt.args) == 0:
				// Default stash should create a stash and show confirmation
				if len(output) == 0 {
					t.Errorf("Expected stash output for default command, got empty string")
				}
			case len(tt.args) > 0:
				command := tt.args[0]
				switch command {
				case "list", "show", "apply", "pop", "drop", "clear":
					// Valid stash commands should produce output
					if len(output) == 0 {
						t.Errorf("Expected output for stash %s command, got empty string", command)
					}
				case "unknown":
					// Unknown commands should show error or help
					if len(output) < 5 {
						t.Errorf("Expected error or help output for unknown command, got: %s", output)
					}
				}
			}
		})
	}
}
