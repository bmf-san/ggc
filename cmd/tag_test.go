package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v5/internal/testutil"
)

func TestTagger_Constructor(t *testing.T) {
	mockClient := testutil.NewMockGitClient()
	tagger := NewTagger(mockClient)

	if tagger == nil {
		t.Fatal("Expected NewTagger to return a non-nil Tagger")
	}
	if tagger.gitClient == nil {
		t.Error("Expected gitClient to be set")
	}
	if tagger.outputWriter == nil {
		t.Error("Expected outputWriter to be set")
	}
	if tagger.helper == nil {
		t.Error("Expected helper to be set")
	}
}

func TestTagger_Tag(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "no args - show help",
			args: []string{},
		},
		{
			name: "list tags",
			args: []string{"list"},
		},
		{
			name: "create tag",
			args: []string{"create", "v1.0.0"},
		},
		{
			name: "create annotated tag",
			args: []string{"create", "v1.0.0", "-m", "Release version 1.0.0"},
		},
		{
			name: "delete tag",
			args: []string{"delete", "v1.0.0"},
		},
		{
			name: "push tag",
			args: []string{"push", "origin", "v1.0.0"},
		},
		{
			name: "push all tags",
			args: []string{"push", "origin", "--all"},
		},
		{
			name: "show tag",
			args: []string{"show", "v1.0.0"},
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

			tagger := &Tagger{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}

			tagger.Tag(tt.args)

			// Verify that the function executed without panic and produced some output
			output := buf.String()
			
			// For all tag commands, we expect some output (help text, results, or error messages)
			if len(output) == 0 {
				t.Errorf("Expected some output for tag command %v, got empty string", tt.args)
			}
			
			// Verify specific output patterns based on command type
			switch {
			case len(tt.args) == 0:
				// No args should show help
				if !containsAny(output, []string{"help", "usage", "tag"}) {
					t.Errorf("Expected help output for no args, got: %s", output)
				}
			case len(tt.args) > 0 && tt.args[0] == "unknown":
				// Unknown commands should show error or help
				if !containsAny(output, []string{"unknown", "help", "usage", "error"}) {
					t.Errorf("Expected error or help for unknown command, got: %s", output)
				}
			default:
				// Other commands should produce relevant output
				if len(output) < 5 { // Minimum reasonable output length
					t.Errorf("Expected meaningful output for command %v, got: %s", tt.args, output)
				}
			}
		})
	}
}

func TestTagger_UtilityMethods(t *testing.T) {
	mockClient := testutil.NewMockGitClient()
	tagger := NewTagger(mockClient)

	// Test GetLatestTag
	tag, err := tagger.GetLatestTag()
	if err != nil {
		t.Errorf("Expected no error from GetLatestTag, got %v", err)
	}
	if tag == "" {
		t.Error("Expected GetLatestTag to return a non-empty string")
	}

	// Test TagExists
	exists := tagger.TagExists("v1.0.0")
	if !exists {
		t.Error("Expected TagExists to return true for mock client")
	}

	// Test GetTagCommit
	commit, err := tagger.GetTagCommit("v1.0.0")
	if err != nil {
		t.Errorf("Expected no error from GetTagCommit, got %v", err)
	}
	if commit == "" {
		t.Error("Expected GetTagCommit to return a non-empty string")
	}
}

// Helper function to check if output contains any of the expected strings
func containsAny(output string, expected []string) bool {
	for _, exp := range expected {
		if strings.Contains(strings.ToLower(output), strings.ToLower(exp)) {
			return true
		}
	}
	return false
}
