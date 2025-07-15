package cmd

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// mockCompleteGitClient is a mock implementation for complete tests
type mockCompleteGitClient struct {
	listLocalBranchesFunc func() ([]string, error)
}

func (m *mockCompleteGitClient) ListLocalBranches() ([]string, error) {
	if m.listLocalBranchesFunc != nil {
		return m.listLocalBranchesFunc()
	}
	return []string{"main", "feature/test"}, nil
}

// Implement other required methods to satisfy git.Clienter interface
func (m *mockCompleteGitClient) GetCurrentBranch() (string, error)     { return "main", nil }
func (m *mockCompleteGitClient) GetGitStatus() (string, error)         { return "", nil }
func (m *mockCompleteGitClient) GetBranchName() (string, error)        { return "main", nil }
func (m *mockCompleteGitClient) ListRemoteBranches() ([]string, error) { return nil, nil }
func (m *mockCompleteGitClient) AddFiles(_ []string) error             { return nil }
func (m *mockCompleteGitClient) CommitAllowEmpty() error               { return nil }
func (m *mockCompleteGitClient) CommitTmp() error                      { return nil }
func (m *mockCompleteGitClient) Commit(_ string) error                 { return nil }
func (m *mockCompleteGitClient) Push(_ bool) error                     { return nil }
func (m *mockCompleteGitClient) Pull(_ bool) error                     { return nil }
func (m *mockCompleteGitClient) LogSimple() error                      { return nil }
func (m *mockCompleteGitClient) LogGraph() error                       { return nil }
func (m *mockCompleteGitClient) ResetHardAndClean() error              { return nil }
func (m *mockCompleteGitClient) CleanFiles() error                     { return nil }
func (m *mockCompleteGitClient) CleanDirs() error                      { return nil }
func (m *mockCompleteGitClient) CheckoutNewBranch(_ string) error      { return nil }
func (m *mockCompleteGitClient) FetchPrune() error                     { return nil }

func TestCompleter_Complete_Branch(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	c := &Completer{
		gitClient: &mockGitClient{},
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("echo")
		},
	}

	c.Complete([]string{"branch"})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Should suggest subcommands
	if !strings.Contains(output, "current") {
		t.Errorf("expected 'current' in output, got %q", output)
	}
}

func TestCompleter_Complete_Branch_WithSecondArg(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	c := &Completer{
		gitClient: &mockGitClient{},
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("echo")
		},
	}

	c.Complete([]string{"branch", "checkout"})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Should suggest branch names
	if !strings.Contains(output, "main") || !strings.Contains(output, "feature/test") {
		t.Errorf("expected branch names in output, got %q", output)
	}
}

func TestCompleter_Complete_Files(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	c := &Completer{
		gitClient: &mockGitClient{},
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("echo", "file1.go\nfile2.go")
		},
	}

	c.Complete([]string{"files"})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "file1.go") {
		t.Errorf("expected file names in output, got %q", output)
	}
}

func TestCompleter_Complete_NoArgs(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	c := &Completer{
		gitClient:   &mockGitClient{},
		execCommand: exec.Command,
	}

	c.Complete([]string{})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Should not output anything
	if output != "" {
		t.Errorf("expected no output for no args, got %q", output)
	}
}

func TestCompleter_Complete_Unknown(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	c := &Completer{
		gitClient:   &mockGitClient{},
		execCommand: exec.Command,
	}

	c.Complete([]string{"unknown"})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Should not output anything for unknown commands
	if output != "" {
		t.Errorf("expected no output for unknown command, got %q", output)
	}
}

func TestCompleter_Complete_BranchNames(t *testing.T) {
	completer := &Completer{
		gitClient: &mockCompleteGitClient{},
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	completer.Complete([]string{"branch", "foo"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("buf.ReadFrom failed: %v", err)
	}
	os.Stdout = oldStdout

	output := buf.String()
	for _, want := range []string{"main", "feature/test"} {
		if !bytes.Contains([]byte(output), []byte(want)) {
			t.Errorf("local branch candidate %s not found in output: %s", want, output)
		}
	}
}

func TestCompleter_Complete_BranchList_Error(t *testing.T) {
	completer := &Completer{
		gitClient: &mockCompleteGitClient{
			listLocalBranchesFunc: func() ([]string, error) {
				return nil, errors.New("failed to list branches")
			},
		},
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	completer.Complete([]string{"branch", "any"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("buf.ReadFrom failed: %v", err)
	}
	os.Stdout = oldStdout

	output := buf.String()
	if output != "" {
		t.Errorf("no output expected on error: %s", output)
	}
}

func TestCompleter_Complete_Files_Error(t *testing.T) {
	completer := &Completer{
		gitClient: &mockGitClient{},
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("false")
		},
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	completer.Complete([]string{"files"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("buf.ReadFrom failed: %v", err)
	}
	os.Stdout = oldStdout

	output := buf.String()
	if output != "" {
		t.Errorf("no output expected on error: %s", output)
	}
}
