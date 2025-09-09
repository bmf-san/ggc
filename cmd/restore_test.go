package cmd

import (
	"bytes"
	"testing"

	"github.com/bmf-san/ggc/v5/internal/testutil"
)

func TestRestorer_Constructor(t *testing.T) {
	mockClient := testutil.NewMockGitClient()
	restorer := NewRestorer(mockClient)

	if restorer == nil {
		t.Fatal("Expected NewRestorer to return a non-nil Restorer")
	}
	if restorer.gitClient == nil {
		t.Error("Expected gitClient to be set")
	}
	if restorer.outputWriter == nil {
		t.Error("Expected outputWriter to be set")
	}
	if restorer.helper == nil {
		t.Error("Expected helper to be set")
	}
}

func TestRestorer_Restore(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "no args - show help",
			args: []string{},
		},
		{
			name: "restore file",
			args: []string{"file.txt"},
		},
		{
			name: "restore all files",
			args: []string{"."},
		},
		{
			name: "restore staged",
			args: []string{"--staged", "file.txt"},
		},
		{
			name: "restore from commit",
			args: []string{"--source", "HEAD~1", "file.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := testutil.NewMockGitClient()

			restorer := &Restorer{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}

			restorer.Restore(tt.args)

			// Verify that the function executed without panic and produced output
			output := buf.String()

			// Restore commands should produce some output (results, help, or error messages)
			if len(output) == 0 {
				t.Errorf("Expected output for restore command %v, got empty string", tt.args)
			}

			// Verify output content based on command type
			switch {
			case len(tt.args) == 0:
				// No args should show help
				if len(output) < 10 {
					t.Errorf("Expected help output for no args, got: %s", output)
				}
			case len(tt.args) > 0:
				// Commands with arguments should show result or confirmation
				if len(tt.args) == 1 {
					// Single argument commands (file or directory)
					if len(output) == 0 {
						t.Errorf("Expected restore output for %v, got empty string", tt.args)
					}
				} else if len(tt.args) > 1 {
					// Commands with flags should produce appropriate output
					if len(output) == 0 {
						t.Errorf("Expected output for restore command with flags %v, got empty string", tt.args)
					}
				}
			}
		})
	}
}
