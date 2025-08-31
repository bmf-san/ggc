package cmd

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"
)

// mockCompleteGitClient is a mock implementation for complete tests
type mockCompleteGitClient struct {
	listLocalBranchesFunc func() ([]string, error)
	getTrackedFilesFunc   func() ([]string, error)
}

func (m *mockCompleteGitClient) ListLocalBranches() ([]string, error) {
	if m.listLocalBranchesFunc != nil {
		return m.listLocalBranchesFunc()
	}
	return []string{"main", "feature/test"}, nil
}

func (m *mockCompleteGitClient) GetTrackedFiles() ([]string, error) {
	if m.getTrackedFilesFunc != nil {
		return m.getTrackedFilesFunc()
	}
	return []string{"file1.go", "file2.go", "README.md"}, nil
}

func (m *mockCompleteGitClient) ListFiles() (string, error) {
	if m.getTrackedFilesFunc != nil {
		files, err := m.getTrackedFilesFunc()
		if err != nil {
			return "", err
		}
		return strings.Join(files, "\n"), nil
	}
	return "file1.go\nfile2.go\nREADME.md", nil
}

// Implement other required methods to satisfy git Clienter interface
func (m *mockCompleteGitClient) GetCurrentBranch() (string, error)         { return "main", nil }
func (m *mockCompleteGitClient) GetGitStatus() (string, error)             { return "", nil }
func (m *mockCompleteGitClient) GetBranchName() (string, error)            { return "main", nil }
func (m *mockCompleteGitClient) ListRemoteBranches() ([]string, error)     { return nil, nil }
func (m *mockCompleteGitClient) AddFiles(_ []string) error                 { return nil }
func (m *mockCompleteGitClient) CommitAllowEmpty() error                   { return nil }
func (m *mockCompleteGitClient) Commit(_ string) error                     { return nil }
func (m *mockCompleteGitClient) Push(_ bool) error                         { return nil }
func (m *mockCompleteGitClient) Pull(_ bool) error                         { return nil }
func (m *mockCompleteGitClient) LogSimple() error                          { return nil }
func (m *mockCompleteGitClient) LogGraph() error                           { return nil }
func (m *mockCompleteGitClient) ResetHardAndClean() error                  { return nil }
func (m *mockCompleteGitClient) CleanFiles() error                         { return nil }
func (m *mockCompleteGitClient) CleanDirs() error                          { return nil }
func (m *mockCompleteGitClient) CheckoutNewBranch(_ string) error          { return nil }
func (m *mockCompleteGitClient) FetchPrune() error                         { return nil }
func (m *mockCompleteGitClient) RestoreAll() error                         { return nil }
func (m *mockCompleteGitClient) RestoreAllStaged() error                   { return nil }
func (m *mockCompleteGitClient) RestoreStaged(...string) error             { return nil }
func (m *mockCompleteGitClient) RestoreWorkingDir(...string) error         { return nil }
func (m *mockCompleteGitClient) RestoreFromCommit(string, ...string) error { return nil }
func (m *mockCompleteGitClient) RevParseVerify(string) bool                { return false }

// Config Operations
func (m *mockCompleteGitClient) ConfigGet(_ string) (string, error)       { return "", nil }
func (m *mockCompleteGitClient) ConfigSet(_, _ string) error              { return nil }
func (m *mockCompleteGitClient) ConfigGetGlobal(_ string) (string, error) { return "", nil }
func (m *mockCompleteGitClient) ConfigSetGlobal(_, _ string) error        { return nil }

// Add missing methods to satisfy git.Clienter interface
func (m *mockCompleteGitClient) Add(_ ...string) error                 { return nil }
func (m *mockCompleteGitClient) AddInteractive() error                 { return nil }
func (m *mockCompleteGitClient) Status() (string, error)               { return "", nil }
func (m *mockCompleteGitClient) StatusShort() (string, error)          { return "", nil }
func (m *mockCompleteGitClient) StatusWithColor() (string, error)      { return "", nil }
func (m *mockCompleteGitClient) StatusShortWithColor() (string, error) { return "", nil }
func (m *mockCompleteGitClient) CommitAmend() error                    { return nil }
func (m *mockCompleteGitClient) CommitAmendNoEdit() error              { return nil }
func (m *mockCompleteGitClient) CommitAmendWithMessage(_ string) error { return nil }
func (m *mockCompleteGitClient) Diff() (string, error)                 { return "", nil }
func (m *mockCompleteGitClient) DiffStaged() (string, error)           { return "", nil }
func (m *mockCompleteGitClient) DiffHead() (string, error)             { return "", nil }
func (m *mockCompleteGitClient) CheckoutBranch(_ string) error         { return nil }
func (m *mockCompleteGitClient) CheckoutNewBranchFromRemote(_, _ string) error {
	return nil
}
func (m *mockCompleteGitClient) DeleteBranch(_ string) error            { return nil }
func (m *mockCompleteGitClient) ListMergedBranches() ([]string, error)  { return []string{}, nil }
func (m *mockCompleteGitClient) Fetch(_ bool) error                     { return nil }
func (m *mockCompleteGitClient) RemoteList() error                      { return nil }
func (m *mockCompleteGitClient) RemoteAdd(_, _ string) error            { return nil }
func (m *mockCompleteGitClient) RemoteRemove(_ string) error            { return nil }
func (m *mockCompleteGitClient) RemoteSetURL(_, _ string) error         { return nil }
func (m *mockCompleteGitClient) TagList(_ []string) error               { return nil }
func (m *mockCompleteGitClient) TagCreate(_, _ string) error            { return nil }
func (m *mockCompleteGitClient) TagCreateAnnotated(_, _ string) error   { return nil }
func (m *mockCompleteGitClient) TagDelete(_ []string) error             { return nil }
func (m *mockCompleteGitClient) TagPush(_, _ string) error              { return nil }
func (m *mockCompleteGitClient) TagPushAll(_ string) error              { return nil }
func (m *mockCompleteGitClient) TagShow(_ string) error                 { return nil }
func (m *mockCompleteGitClient) GetLatestTag() (string, error)          { return "v1.0.0", nil }
func (m *mockCompleteGitClient) TagExists(_ string) bool                { return false }
func (m *mockCompleteGitClient) GetTagCommit(_ string) (string, error)  { return "abc123", nil }
func (m *mockCompleteGitClient) LogOneline(_, _ string) (string, error) { return "", nil }
func (m *mockCompleteGitClient) RebaseInteractive(_ int) error          { return nil }
func (m *mockCompleteGitClient) GetUpstreamBranch(_ string) (string, error) {
	return "origin/main", nil
}
func (m *mockCompleteGitClient) Stash() error                     { return nil }
func (m *mockCompleteGitClient) StashList() (string, error)       { return "", nil }
func (m *mockCompleteGitClient) StashShow(_ string) error         { return nil }
func (m *mockCompleteGitClient) StashApply(_ string) error        { return nil }
func (m *mockCompleteGitClient) StashPop(_ string) error          { return nil }
func (m *mockCompleteGitClient) StashDrop(_ string) error         { return nil }
func (m *mockCompleteGitClient) StashClear() error                { return nil }
func (m *mockCompleteGitClient) ResetHard(_ string) error         { return nil }
func (m *mockCompleteGitClient) CleanDryRun() (string, error)     { return "", nil }
func (m *mockCompleteGitClient) CleanFilesForce(_ []string) error { return nil }

func (m *mockCompleteGitClient) GetUpstreamBranchName(_ string) (string, error) {
	return "origin/main", nil
}
func (m *mockCompleteGitClient) GetAheadBehindCount(_, _ string) (string, error) {
	return "0	0", nil
}
func (m *mockCompleteGitClient) GetVersion() (string, error)    { return "test-version", nil }
func (m *mockCompleteGitClient) GetCommitHash() (string, error) { return "test-commit", nil }

func TestCompleter_Complete_Branch(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	c := &Completer{
		gitClient: &mockCompleteGitClient{},
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
		gitClient: &mockCompleteGitClient{},
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
		gitClient: &mockCompleteGitClient{},
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
		gitClient: &mockGitClient{},
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
		gitClient: &mockGitClient{},
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
		gitClient: &mockCompleteGitClient{
			getTrackedFilesFunc: func() ([]string, error) {
				return nil, errors.New("git error")
			},
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
