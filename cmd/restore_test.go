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

			// Basic test - just ensure no panic occurs
			// In a real test, we would check specific outputs based on mock responses
		})
	}
}
