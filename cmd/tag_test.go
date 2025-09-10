package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v5/git"
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

// Test CreateAnnotatedTag function to improve coverage from 0.0%
func TestTagger_CreateAnnotatedTag(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "create annotated tag with message",
			args:           []string{"v1.0.0", "Release", "version", "1.0.0"},
			expectedOutput: "Annotated tag 'v1.0.0' created",
			expectError:    false,
		},
		{
			name:           "create annotated tag without message (editor)",
			args:           []string{"v1.0.0"},
			expectedOutput: "Annotated tag 'v1.0.0' created",
			expectError:    false,
		},
		{
			name:           "create annotated tag with single word message",
			args:           []string{"v2.0.0", "Major"},
			expectedOutput: "Annotated tag 'v2.0.0' created",
			expectError:    false,
		},
		{
			name:           "create annotated tag without name - should show error",
			args:           []string{},
			expectedOutput: "Error: tag name is required",
			expectError:    true,
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

			tagger.CreateAnnotatedTag(tt.args)

			output := buf.String()
			if !strings.Contains(output, tt.expectedOutput) {
				t.Errorf("Expected output containing '%s', got: %s", tt.expectedOutput, output)
			}
		})
	}
}

// mockTagErrorClient wraps the regular mock client and overrides tag methods to return errors
type mockTagErrorClient struct {
	git.Clienter
}

func newMockTagErrorClient() *mockTagErrorClient {
	return &mockTagErrorClient{
		Clienter: testutil.NewMockGitClient(),
	}
}

func (m *mockTagErrorClient) TagList(args []string) error {
	return fmt.Errorf("mock tag list error")
}

func (m *mockTagErrorClient) TagCreate(name, commit string) error {
	return fmt.Errorf("mock tag create error")
}

func (m *mockTagErrorClient) TagDelete(names []string) error {
	return fmt.Errorf("mock tag delete error")
}

func (m *mockTagErrorClient) TagPushAll(remote string) error {
	return fmt.Errorf("mock tag push all error")
}

func (m *mockTagErrorClient) TagPush(remote, name string) error {
	return fmt.Errorf("mock tag push error")
}

func (m *mockTagErrorClient) TagShow(name string) error {
	return fmt.Errorf("mock tag show error")
}

func (m *mockTagErrorClient) TagCreateAnnotated(name, message string) error {
	return fmt.Errorf("mock annotated tag create error")
}

// Test error cases in git operations to improve coverage
func TestTagger_GitOperationErrors(t *testing.T) {
	tests := []struct {
		name           string
		operation      string
		args           []string
		expectedOutput string
	}{
		{
			name:           "list tags with git error",
			operation:      "list",
			args:           []string{"list"},
			expectedOutput: "Error:",
		},
		{
			name:           "create tag with git error",
			operation:      "create",
			args:           []string{"create", "v1.0.0"},
			expectedOutput: "Error:",
		},
		{
			name:           "create tag with commit and git error",
			operation:      "create",
			args:           []string{"create", "v1.0.0", "abc123"},
			expectedOutput: "Error:",
		},
		{
			name:           "delete tag with git error",
			operation:      "delete",
			args:           []string{"delete", "v1.0.0"},
			expectedOutput: "Error:",
		},
		{
			name:           "push all tags with git error",
			operation:      "push",
			args:           []string{"push"},
			expectedOutput: "Error:",
		},
		{
			name:           "push specific tag with git error",
			operation:      "push",
			args:           []string{"push", "v1.0.0"},
			expectedOutput: "Error:",
		},
		{
			name:           "show tag with git error",
			operation:      "show",
			args:           []string{"show", "v1.0.0"},
			expectedOutput: "Error:",
		},
		{
			name:           "create annotated tag with git error",
			operation:      "annotated",
			args:           []string{"v1.0.0", "message"},
			expectedOutput: "Error:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := newMockTagErrorClient()

			tagger := &Tagger{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}

			switch tt.operation {
			case "list":
				tagger.Tag(tt.args)
			case "create":
				tagger.Tag(tt.args)
			case "delete":
				tagger.Tag(tt.args)
			case "push":
				tagger.Tag(tt.args)
			case "show":
				tagger.Tag(tt.args)
			case "annotated":
				tagger.CreateAnnotatedTag(tt.args)
			}

			output := buf.String()
			if !strings.Contains(output, tt.expectedOutput) {
				t.Errorf("Expected error output containing '%s', got: %s", tt.expectedOutput, output)
			}
		})
	}
}

// Test listTags with various arguments to improve coverage from 50.0%
func TestTagger_listTags_Coverage(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		mockError bool
	}{
		{
			name:      "list with pattern",
			args:      []string{"v1.*"},
			mockError: false,
		},
		{
			name:      "list with multiple patterns",
			args:      []string{"v1.*", "v2.*"},
			mockError: false,
		},
		{
			name:      "list with git error",
			args:      []string{},
			mockError: true,
		},
		{
			name:      "list with pattern and git error",
			args:      []string{"v*"},
			mockError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			var mockClient git.Clienter
			
			if tt.mockError {
				mockClient = newMockTagErrorClient()
			} else {
				mockClient = testutil.NewMockGitClient()
			}

			tagger := &Tagger{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}

			tagger.listTags(tt.args)

			output := buf.String()
			if tt.mockError {
				if !strings.Contains(output, "Error:") {
					t.Errorf("Expected error output, got: %s", output)
				}
			}
		})
	}
}

// Test edge cases and additional scenarios
func TestTagger_EdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		testFunc       func(*testing.T, *Tagger, *bytes.Buffer)
		expectedOutput string
	}{
		{
			name: "no args should call TagList",
			testFunc: func(t *testing.T, tagger *Tagger, buf *bytes.Buffer) {
				tagger.Tag([]string{})
			},
			expectedOutput: "", // No output expected from successful TagList
		},
		{
			name: "no args with git error should show error",
			testFunc: func(t *testing.T, tagger *Tagger, buf *bytes.Buffer) {
				// Use error mock client
				tagger.gitClient = newMockTagErrorClient()
				tagger.Tag([]string{})
			},
			expectedOutput: "Error:",
		},
		{
			name: "delete multiple tags success",
			testFunc: func(t *testing.T, tagger *Tagger, buf *bytes.Buffer) {
				tagger.deleteTags([]string{"v1.0.0", "v1.1.0", "v2.0.0"})
			},
			expectedOutput: "Tag 'v1.0.0' deleted",
		},
		{
			name: "push tag with custom remote",
			testFunc: func(t *testing.T, tagger *Tagger, buf *bytes.Buffer) {
				tagger.pushTags([]string{"v1.0.0", "upstream"})
			},
			expectedOutput: "Tag 'v1.0.0' pushed to upstream",
		},
		{
			name: "create tag with long commit hash",
			testFunc: func(t *testing.T, tagger *Tagger, buf *bytes.Buffer) {
				tagger.createTag([]string{"v1.0.0", "abcdef1234567890abcdef1234567890abcdef12"})
			},
			expectedOutput: "Tag 'v1.0.0' created",
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

			tt.testFunc(t, tagger, buf)

			output := buf.String()
			if tt.expectedOutput != "" {
				if !strings.Contains(output, tt.expectedOutput) {
					t.Errorf("Expected output containing '%s', got: %s", tt.expectedOutput, output)
				}
			}
		})
	}
}

// Test comprehensive tag command routing
func TestTagger_CommandRouting(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
	}{
		{
			name:           "list command with l alias",
			args:           []string{"l", "v*"},
			expectedOutput: "", // Successful list operation
		},
		{
			name:           "create command with c alias",
			args:           []string{"c", "v3.0.0"},
			expectedOutput: "Tag 'v3.0.0' created",
		},
		{
			name:           "delete command with d alias",
			args:           []string{"d", "v3.0.0"},
			expectedOutput: "Tag 'v3.0.0' deleted",
		},
		{
			name:           "invalid command shows help",
			args:           []string{"invalid"},
			expectedOutput: "Usage:",
		},
		{
			name:           "empty string command shows help",
			args:           []string{""},
			expectedOutput: "Usage:",
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

			output := buf.String()
			if tt.expectedOutput != "" {
				if !strings.Contains(output, tt.expectedOutput) {
					t.Errorf("Expected output containing '%s', got: %s", tt.expectedOutput, output)
				}
			}
		})
	}
}
