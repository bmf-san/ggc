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
		name           string
		args           []string
		expectedOutput string
		shouldShowHelp bool
	}{
		{
			name:           "no args - list all tags",
			args:           []string{},
			expectedOutput: "", // Mock client returns empty for TagList
			shouldShowHelp: false,
		},
		{
			name:           "list tags",
			args:           []string{"list"},
			expectedOutput: "", // Mock client returns empty for TagList
			shouldShowHelp: false,
		},
		{
			name:           "list tags with alias",
			args:           []string{"l"},
			expectedOutput: "", // Mock client returns empty for TagList
			shouldShowHelp: false,
		},
		{
			name:           "create tag",
			args:           []string{"create", "v1.0.0"},
			expectedOutput: "Tag 'v1.0.0' created",
			shouldShowHelp: false,
		},
		{
			name:           "create tag with alias",
			args:           []string{"c", "v1.0.0"},
			expectedOutput: "Tag 'v1.0.0' created",
			shouldShowHelp: false,
		},
		{
			name:           "create tag with commit",
			args:           []string{"create", "v1.0.0", "abc123"},
			expectedOutput: "Tag 'v1.0.0' created",
			shouldShowHelp: false,
		},
		{
			name:           "create tag without name - should show error",
			args:           []string{"create"},
			expectedOutput: "Error: tag name is required",
			shouldShowHelp: false,
		},
		{
			name:           "delete tag",
			args:           []string{"delete", "v1.0.0"},
			expectedOutput: "Tag 'v1.0.0' deleted",
			shouldShowHelp: false,
		},
		{
			name:           "delete tag with alias",
			args:           []string{"d", "v1.0.0"},
			expectedOutput: "Tag 'v1.0.0' deleted",
			shouldShowHelp: false,
		},
		{
			name:           "delete multiple tags",
			args:           []string{"delete", "v1.0.0", "v1.1.0"},
			expectedOutput: "Tag 'v1.0.0' deleted",
			shouldShowHelp: false,
		},
		{
			name:           "delete tag without name - should show error",
			args:           []string{"delete"},
			expectedOutput: "Error: at least one tag name is required",
			shouldShowHelp: false,
		},
		{
			name:           "push all tags",
			args:           []string{"push"},
			expectedOutput: "All tags pushed to origin",
			shouldShowHelp: false,
		},
		{
			name:           "push specific tag",
			args:           []string{"push", "v1.0.0"},
			expectedOutput: "Tag 'v1.0.0' pushed to origin",
			shouldShowHelp: false,
		},
		{
			name:           "push tag to specific remote",
			args:           []string{"push", "v1.0.0", "upstream"},
			expectedOutput: "Tag 'v1.0.0' pushed to upstream",
			shouldShowHelp: false,
		},
		{
			name:           "show tag",
			args:           []string{"show", "v1.0.0"},
			expectedOutput: "", // Mock client returns empty for TagShow
			shouldShowHelp: false,
		},
		{
			name:           "show tag without name - should show error",
			args:           []string{"show"},
			expectedOutput: "Error: tag name is required",
			shouldShowHelp: false,
		},
		{
			name:           "unknown command - should show help",
			args:           []string{"unknown"},
			expectedOutput: "Usage: ggc tag [command] [options]",
			shouldShowHelp: true,
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
			// Set helper's output writer to capture help output
			tagger.helper.outputWriter = buf

			tagger.Tag(tt.args)

			output := buf.String()

			// Verify expected behavior
			if tt.shouldShowHelp {
				if !strings.Contains(output, tt.expectedOutput) {
					t.Errorf("Expected help output containing '%s', got: %s", tt.expectedOutput, output)
				}
			} else {
				if tt.expectedOutput == "" {
					// For operations that return empty (like list, show), mock returns empty string
					// We verify the command executed without error
					if strings.Contains(output, "Error:") {
						t.Errorf("Unexpected error in tag operation: %s", output)
					}
				} else {
					// For operations with expected output (create, delete, push)
					if !strings.Contains(output, tt.expectedOutput) {
						t.Errorf("Expected output containing '%s', got: %s", tt.expectedOutput, output)
					}
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

func TestTagger_TagOperations(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		testFunc func(*testing.T, *Tagger, *bytes.Buffer)
	}{
		{
			name: "list operation calls TagList",
			args: []string{"list"},
			testFunc: func(t *testing.T, tagger *Tagger, buf *bytes.Buffer) {
				output := buf.String()
				// Mock client doesn't return errors for TagList
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in tag list: %s", output)
				}
			},
		},
		{
			name: "create operation success message",
			args: []string{"create", "v1.0.0"},
			testFunc: func(t *testing.T, tagger *Tagger, buf *bytes.Buffer) {
				output := buf.String()
				// Should show success message
				if !strings.Contains(output, "Tag 'v1.0.0' created") {
					t.Errorf("Expected create success message, got: %s", output)
				}
			},
		},
		{
			name: "delete operation success message",
			args: []string{"delete", "v1.0.0"},
			testFunc: func(t *testing.T, tagger *Tagger, buf *bytes.Buffer) {
				output := buf.String()
				// Should show success message
				if !strings.Contains(output, "Tag 'v1.0.0' deleted") {
					t.Errorf("Expected delete success message, got: %s", output)
				}
			},
		},
		{
			name: "push all tags success message",
			args: []string{"push"},
			testFunc: func(t *testing.T, tagger *Tagger, buf *bytes.Buffer) {
				output := buf.String()
				// Should show success message
				if !strings.Contains(output, "All tags pushed to origin") {
					t.Errorf("Expected push all success message, got: %s", output)
				}
			},
		},
		{
			name: "push specific tag success message",
			args: []string{"push", "v1.0.0"},
			testFunc: func(t *testing.T, tagger *Tagger, buf *bytes.Buffer) {
				output := buf.String()
				// Should show success message
				if !strings.Contains(output, "Tag 'v1.0.0' pushed to origin") {
					t.Errorf("Expected push tag success message, got: %s", output)
				}
			},
		},
		{
			name: "show operation calls TagShow",
			args: []string{"show", "v1.0.0"},
			testFunc: func(t *testing.T, tagger *Tagger, buf *bytes.Buffer) {
				output := buf.String()
				// Mock client doesn't return errors for TagShow
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in tag show: %s", output)
				}
			},
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
			tagger.helper.outputWriter = buf

			tagger.Tag(tt.args)
			tt.testFunc(t, tagger, buf)
		})
	}
}

func TestTagger_ErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		testFunc func(*testing.T, *Tagger, *bytes.Buffer)
	}{
		{
			name: "create without tag name shows error",
			args: []string{"create"},
			testFunc: func(t *testing.T, tagger *Tagger, buf *bytes.Buffer) {
				output := buf.String()
				if !strings.Contains(output, "Error: tag name is required") {
					t.Errorf("Expected error message for missing tag name, got: %s", output)
				}
			},
		},
		{
			name: "delete without tag name shows error",
			args: []string{"delete"},
			testFunc: func(t *testing.T, tagger *Tagger, buf *bytes.Buffer) {
				output := buf.String()
				if !strings.Contains(output, "Error: at least one tag name is required") {
					t.Errorf("Expected error message for missing tag names, got: %s", output)
				}
			},
		},
		{
			name: "show without tag name shows error",
			args: []string{"show"},
			testFunc: func(t *testing.T, tagger *Tagger, buf *bytes.Buffer) {
				output := buf.String()
				if !strings.Contains(output, "Error: tag name is required") {
					t.Errorf("Expected error message for missing tag name, got: %s", output)
				}
			},
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
			tagger.helper.outputWriter = buf

			tagger.Tag(tt.args)
			tt.testFunc(t, tagger, buf)
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
