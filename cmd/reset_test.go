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

			// Basic test - just ensure no panic occurs
			// In a real test, we would check specific outputs based on mock responses
		})
	}
}
