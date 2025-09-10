package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v5/git"
)

// mockDiffClient implements git.DiffReader with simple call tracking.
type mockDiffClient struct {
	called      string
	headOut     string
	stagedOut   string
	unstagedOut string
}

func (m *mockDiffClient) Diff() (string, error) {
	m.called = "unstaged"
	return m.unstagedOut, nil
}

func (m *mockDiffClient) DiffStaged() (string, error) {
	m.called = "staged"
	return m.stagedOut, nil
}

func (m *mockDiffClient) DiffHead() (string, error) {
	m.called = "head"
	return m.headOut, nil
}

var _ git.DiffReader = (*mockDiffClient)(nil)

func TestDiffer_Constructor(t *testing.T) {
	mockClient := &mockDiffClient{}
	differ := NewDiffer(mockClient)

	if differ == nil {
		t.Fatal("Expected NewDiffer to return a non-nil Differ")
	}
	if differ.gitClient == nil {
		t.Error("Expected gitClient to be set")
	}
	if differ.outputWriter == nil {
		t.Error("Expected outputWriter to be set")
	}
	if differ.helper == nil {
		t.Error("Expected helper to be set")
	}
}

func TestDiffer_Diff(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		shouldShowHelp bool
	}{
		{
			name:           "no args - should call DiffHead",
			args:           []string{},
			expectedOutput: "",
			shouldShowHelp: false,
		},
		{
			name:           "unstaged - should call Diff",
			args:           []string{"unstaged"},
			expectedOutput: "",
			shouldShowHelp: false,
		},
		{
			name:           "staged - should call DiffStaged",
			args:           []string{"staged"},
			expectedOutput: "",
			shouldShowHelp: false,
		},
		{
			name:           "head - should call DiffHead",
			args:           []string{"head"},
			expectedOutput: "",
			shouldShowHelp: false,
		},
		{
			name:           "unknown arg - should show help",
			args:           []string{"unknown"},
			expectedOutput: "Usage: ggc diff [options]",
			shouldShowHelp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := &mockDiffClient{}

			differ := &Differ{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}
			differ.helper.outputWriter = buf

			differ.Diff(tt.args)

			output := buf.String()

			if tt.shouldShowHelp {
				if !strings.Contains(output, tt.expectedOutput) {
					t.Errorf("Expected help output containing '%s', got: %s", tt.expectedOutput, output)
				}
			} else {
				if tt.expectedOutput == "" {
					if strings.Contains(output, "Error:") {
						t.Errorf("Unexpected error in diff operation: %s", output)
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

func TestDiffer_DiffOperations(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		testFunc func(*testing.T, *Differ, *mockDiffClient, *bytes.Buffer)
	}{
		{
			name: "default calls DiffHead",
			args: []string{},
			testFunc: func(t *testing.T, differ *Differ, mock *mockDiffClient, buf *bytes.Buffer) {
				if mock.called != "head" {
					t.Errorf("Expected DiffHead to be called, got %s", mock.called)
				}
			},
		},
		{
			name: "unstaged calls Diff",
			args: []string{"unstaged"},
			testFunc: func(t *testing.T, differ *Differ, mock *mockDiffClient, buf *bytes.Buffer) {
				if mock.called != "unstaged" {
					t.Errorf("Expected Diff to be called, got %s", mock.called)
				}
			},
		},
		{
			name: "staged calls DiffStaged",
			args: []string{"staged"},
			testFunc: func(t *testing.T, differ *Differ, mock *mockDiffClient, buf *bytes.Buffer) {
				if mock.called != "staged" {
					t.Errorf("Expected DiffStaged to be called, got %s", mock.called)
				}
			},
		},
		{
			name: "head calls DiffHead",
			args: []string{"head"},
			testFunc: func(t *testing.T, differ *Differ, mock *mockDiffClient, buf *bytes.Buffer) {
				if mock.called != "head" {
					t.Errorf("Expected DiffHead to be called, got %s", mock.called)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := &mockDiffClient{}

			differ := &Differ{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}

			differ.Diff(tt.args)
			tt.testFunc(t, differ, mockClient, buf)
		})
	}
}

func TestDiffer_Diff_DefaultsToHead(t *testing.T) {
	mc := &mockDiffClient{headOut: "HEAD DIFF"}
	var buf bytes.Buffer

	d := &Differ{gitClient: mc, outputWriter: &buf, helper: NewHelper()}
	d.Diff([]string{})

	if mc.called != "head" {
		t.Fatalf("expected DiffHead to be called, got %q", mc.called)
	}
	if got := buf.String(); got != "HEAD DIFF" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestDiffer_Diff_Unstaged(t *testing.T) {
	mc := &mockDiffClient{unstagedOut: "UNSTAGED DIFF"}
	var buf bytes.Buffer

	d := &Differ{gitClient: mc, outputWriter: &buf, helper: NewHelper()}
	d.Diff([]string{"unstaged"})

	if mc.called != "unstaged" {
		t.Fatalf("expected Diff to be called, got %q", mc.called)
	}
	if got := buf.String(); got != "UNSTAGED DIFF" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestDiffer_Diff_Staged(t *testing.T) {
	mc := &mockDiffClient{stagedOut: "STAGED DIFF"}
	var buf bytes.Buffer

	d := &Differ{gitClient: mc, outputWriter: &buf, helper: NewHelper()}
	d.Diff([]string{"staged"})

	if mc.called != "staged" {
		t.Fatalf("expected DiffStaged to be called, got %q", mc.called)
	}
	if got := buf.String(); got != "STAGED DIFF" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestDiffer_Diff_InvalidArg_ShowsHelp(t *testing.T) {
	mc := &mockDiffClient{}
	var buf bytes.Buffer

	// Use helper with our buffer to capture help output
	h := NewHelper()
	h.outputWriter = &buf

	d := &Differ{gitClient: mc, outputWriter: &buf, helper: h}
	d.Diff([]string{"unknown"})

	if mc.called != "" {
		t.Fatalf("diff methods should not be called for invalid arg, got %q", mc.called)
	}
	if buf.Len() == 0 {
		t.Fatal("expected help output for invalid arg")
	}
}
