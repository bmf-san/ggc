package cmd

import (
	"bytes"
	"testing"

	"github.com/bmf-san/ggc/v5/internal/testutil"
)

func TestResetter_Constructor(t *testing.T) {
	mockClient := testutil.NewMockGitClient()
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
		name string
		args []string
	}{
		{
			name: "no args - default reset",
			args: []string{},
		},
		{
			name: "hard reset with commit",
			args: []string{"hard", "abc123"},
		},
		{
			name: "hard reset without commit",
			args: []string{"hard"},
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

			resetter := &Resetter{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}

			resetter.Reset(tt.args)

			// Verify that the function executed without panic and produced output
			output := buf.String()

			// Reset commands should produce some output (results, help, or error messages)
			if len(output) == 0 {
				t.Errorf("Expected output for reset command %v, got empty string", tt.args)
			}

			// Verify output content based on command type
			switch {
			case len(tt.args) == 0:
				// Default reset should produce some output
				if len(output) == 0 {
					t.Errorf("Expected default reset output, got empty string")
				}
			case len(tt.args) > 0 && tt.args[0] == "hard":
				// Hard reset should show confirmation or result
				if len(output) == 0 {
					t.Errorf("Expected hard reset output, got empty string")
				}
			case len(tt.args) > 0 && tt.args[0] == "unknown":
				// Unknown commands should show error or help
				if len(output) < 5 {
					t.Errorf("Expected error or help output for unknown command, got: %s", output)
				}
			}
		})
	}
}
