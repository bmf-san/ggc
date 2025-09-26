package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v6/git"
)

type mockTagOps struct {
	listCalled      bool
	showCalled      bool
	createCalled    bool
	createAnnCalled bool
	deleteCalled    bool
	pushCalled      bool
	pushAllCalled   bool

	listPattern   []string
	showName      string
	createName    string
	createCommit  string
	createAnnName string
	createAnnMsg  string
	deleteNames   []string
	pushRemote    string
	pushName      string

	errList      error
	errShow      error
	errCreate    error
	errCreateAnn error
	errDelete    error
	errPush      error
	errPushAll   error

	latestTag string
	tagExists bool
	tagCommit string
}

func (m *mockTagOps) TagList(pattern []string) error {
	m.listCalled = true
	m.listPattern = pattern
	return m.errList
}
func (m *mockTagOps) TagShow(name string) error {
	m.showCalled = true
	m.showName = name
	return m.errShow
}
func (m *mockTagOps) TagCreate(name, commit string) error {
	m.createCalled = true
	m.createName, m.createCommit = name, commit
	return m.errCreate
}
func (m *mockTagOps) TagCreateAnnotated(name, message string) error {
	m.createAnnCalled = true
	m.createAnnName, m.createAnnMsg = name, message
	return m.errCreateAnn
}
func (m *mockTagOps) TagDelete(names []string) error {
	m.deleteCalled = true
	m.deleteNames = names
	return m.errDelete
}
func (m *mockTagOps) TagPush(remote, name string) error {
	m.pushCalled = true
	m.pushRemote, m.pushName = remote, name
	return m.errPush
}
func (m *mockTagOps) TagPushAll(remote string) error {
	m.pushAllCalled = true
	m.pushRemote = remote
	return m.errPushAll
}
func (m *mockTagOps) GetLatestTag() (string, error)       { return m.latestTag, nil }
func (m *mockTagOps) TagExists(string) bool               { return m.tagExists }
func (m *mockTagOps) GetTagCommit(string) (string, error) { return m.tagCommit, nil }

func (m *mockTagOps) ConfigGetGlobal(string) (string, error) { return "", nil }
func (m *mockTagOps) ConfigSetGlobal(string, string) error   { return nil }
func (m *mockTagOps) GetVersion() (string, error)            { return "test-version", nil }
func (m *mockTagOps) GetCommitHash() (string, error)         { return "test-commit", nil }

var _ git.TagOps = (*mockTagOps)(nil)

func TestTagger_Constructor(t *testing.T) {
	mockClient := &mockTagOps{}
	tagger := NewTagger(mockClient)

	if tagger == nil {
		t.Fatal("Expected NewTagger to return a non-nil Tagger")
	}
	if tagger != nil && tagger.gitClient == nil {
		t.Error("Expected gitClient to be set")
	}
	if tagger != nil && tagger.outputWriter == nil {
		t.Error("Expected outputWriter to be set")
	}
	if tagger != nil && tagger.helper == nil {
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
			expectedOutput: "",
			shouldShowHelp: false,
		},
		{
			name:           "list tags",
			args:           []string{"list"},
			expectedOutput: "",
			shouldShowHelp: false,
		},
		{
			name:           "list tags with alias",
			args:           []string{"l"},
			expectedOutput: "",
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
			expectedOutput: "",
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
			mockClient := &mockTagOps{}

			tagger := &Tagger{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}
			tagger.helper.outputWriter = buf

			tagger.Tag(tt.args)

			output := buf.String()

			if tt.shouldShowHelp {
				if !strings.Contains(output, tt.expectedOutput) {
					t.Errorf("Expected help output containing '%s', got: %s", tt.expectedOutput, output)
				}
			} else {
				if tt.expectedOutput == "" {
					if strings.Contains(output, "Error:") {
						t.Errorf("Unexpected error in tag operation: %s", output)
					}
				} else {
					if !strings.Contains(output, tt.expectedOutput) {
						t.Errorf("Expected output containing '%s', got: %s", tt.expectedOutput, output)
					}
				}
			}

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
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in tag show: %s", output)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := &mockTagOps{}

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
			mockClient := &mockTagOps{}

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
	mockClient := &mockTagOps{
		latestTag: "v2.0.0",
		tagExists: true,
		tagCommit: "abc123",
	}
	tagger := NewTagger(mockClient)

	tag, err := tagger.GetLatestTag()
	if err != nil {
		t.Errorf("Expected no error from GetLatestTag, got %v", err)
	}
	if tag != "v2.0.0" {
		t.Errorf("Expected GetLatestTag to return 'v2.0.0', got %s", tag)
	}

	exists := tagger.TagExists("v1.0.0")
	if !exists {
		t.Error("Expected TagExists to return true for mock client")
	}

	commit, err := tagger.GetTagCommit("v1.0.0")
	if err != nil {
		t.Errorf("Expected no error from GetTagCommit, got %v", err)
	}
	if commit != "abc123" {
		t.Errorf("Expected GetTagCommit to return 'abc123', got %s", commit)
	}
}

func TestTagger_List_UsesTagOps(t *testing.T) {
	m := &mockTagOps{}
	var buf bytes.Buffer
	tg := &Tagger{gitClient: m, outputWriter: &buf, helper: NewHelper()}
	tg.helper.outputWriter = &buf

	tg.Tag([]string{"list"})

	if !m.listCalled {
		t.Fatal("expected TagList to be called")
	}
}

func TestTagger_Create_UsesTagOps(t *testing.T) {
	m := &mockTagOps{}
	var buf bytes.Buffer
	tg := &Tagger{gitClient: m, outputWriter: &buf, helper: NewHelper()}
	tg.helper.outputWriter = &buf

	tg.Tag([]string{"create", "v1.2.3"})

	if !m.createCalled || m.createName != "v1.2.3" || m.createCommit != "" {
		t.Fatalf("unexpected create call: called=%v name=%q commit=%q", m.createCalled, m.createName, m.createCommit)
	}
}

// Tests for CreateAnnotatedTag function (0% coverage before)
func TestTagger_CreateAnnotatedTag(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedOut string
		expectedMsg string
		shouldError bool
	}{
		{
			name:        "no tag name provided",
			args:        []string{},
			expectedOut: "Error: tag name is required",
			expectedMsg: "",
			shouldError: false,
		},
		{
			name:        "tag with message",
			args:        []string{"v1.0.0", "Release", "version", "1.0.0"},
			expectedOut: "Annotated tag 'v1.0.0' created",
			expectedMsg: "Release version 1.0.0",
			shouldError: false,
		},
		{
			name:        "tag without message (editor mode)",
			args:        []string{"v1.0.0"},
			expectedOut: "Annotated tag 'v1.0.0' created",
			expectedMsg: "",
			shouldError: false,
		},
		{
			name:        "error during tag creation",
			args:        []string{"invalid-tag"},
			expectedOut: "Error: failed to create annotated tag",
			expectedMsg: "",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			mockClient := &mockTagOps{}

			if tt.shouldError {
				mockClient.errCreateAnn = errors.New("failed to create annotated tag")
			}

			tagger := &Tagger{
				gitClient:    mockClient,
				outputWriter: &buf,
			}

			tagger.CreateAnnotatedTag(tt.args)

			output := buf.String()
			if !strings.Contains(output, tt.expectedOut) {
				t.Errorf("Expected %q in output, got: %s", tt.expectedOut, output)
			}

			// Verify mock was called appropriately
			if len(tt.args) > 0 && !tt.shouldError {
				if !mockClient.createAnnCalled {
					t.Error("Expected TagCreateAnnotated to be called")
				}
				if mockClient.createAnnName != tt.args[0] {
					t.Errorf("Expected tag name %q, got %q", tt.args[0], mockClient.createAnnName)
				}
				if mockClient.createAnnMsg != tt.expectedMsg {
					t.Errorf("Expected message %q, got %q", tt.expectedMsg, mockClient.createAnnMsg)
				}
			} else if tt.shouldError && len(tt.args) > 0 {
				if !mockClient.createAnnCalled {
					t.Error("Expected TagCreateAnnotated to be called even when error occurs")
				}
			}
		})
	}
}

func TestTagger_Push_All_UsesTagOps(t *testing.T) {
	m := &mockTagOps{}
	var buf bytes.Buffer
	tg := &Tagger{gitClient: m, outputWriter: &buf, helper: NewHelper()}
	tg.helper.outputWriter = &buf

	tg.Tag([]string{"push"})

	if !m.pushAllCalled || m.pushRemote != "origin" {
		t.Fatalf("expected TagPushAll to origin, got called=%v remote=%q", m.pushAllCalled, m.pushRemote)
	}
}

func TestTagger_Show_ErrorOutputs(t *testing.T) {
	m := &mockTagOps{errShow: errors.New("boom")}
	var buf bytes.Buffer
	tg := &Tagger{gitClient: m, outputWriter: &buf, helper: NewHelper()}
	tg.helper.outputWriter = &buf

	tg.Tag([]string{"show", "v0.1.0"})

	if got := buf.String(); got == "" || got[:6] != "Error:"[:6] {
		t.Fatalf("expected error output, got: %q", got)
	}
}
