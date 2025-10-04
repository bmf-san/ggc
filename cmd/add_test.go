package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v7/pkg/git"
)

func TestAdder_Add_NoArgs_PrintsUsage(t *testing.T) {
	mockClient := &mockAddGitClient{}
	var buf bytes.Buffer
	adder := &Adder{
		gitClient:    mockClient,
		outputWriter: &buf,
	}

	adder.Add([]string{})

	output := buf.String()
	if output == "" || output[:5] != "Usage" {
		t.Errorf("Usage not output: %s", output)
	}
}

// mockGitClient for testing
type mockAddGitClient struct {
	addCalled             bool
	addInteractiveCalled  bool
	addFiles              []string
	addError              error
	addInteractiveError   error
	GetCurrentBranchFunc  func() (string, error)
	LogOnelineFunc        func(from, to string) (string, error)
	RebaseInteractiveFunc func(commitCount int) error
	RebaseCalled          bool
	RebaseUpstream        string
	RebaseErr             error
	RebaseContinueCalled  bool
	RebaseContinueErr     error
	RebaseAbortCalled     bool
	RebaseAbortErr        error
	RebaseSkipCalled      bool
	RebaseSkipErr         error
	RevParseVerifyFunc    func(ref string) bool
}

func (m *mockAddGitClient) Add(files ...string) error {
	m.addCalled = true
	m.addFiles = files
	return m.addError
}

func (m *mockAddGitClient) AddInteractive() error {
	m.addInteractiveCalled = true
	return m.addInteractiveError
}

// Repository Information methods
func (m *mockAddGitClient) GetCurrentBranch() (string, error) {
	if m.GetCurrentBranchFunc != nil {
		return m.GetCurrentBranchFunc()
	}
	return "feature/test", nil
}
func (m *mockAddGitClient) GetBranchName() (string, error) { return "main", nil }
func (m *mockAddGitClient) GetGitStatus() (string, error)  { return "", nil }

// Status Operations methods
func (m *mockAddGitClient) Status() (string, error)               { return "", nil }
func (m *mockAddGitClient) StatusShort() (string, error)          { return "", nil }
func (m *mockAddGitClient) StatusWithColor() (string, error)      { return "", nil }
func (m *mockAddGitClient) StatusShortWithColor() (string, error) { return "", nil }

// Commit Operations methods
func (m *mockAddGitClient) Commit(_ string) error                 { return nil }
func (m *mockAddGitClient) CommitAmend() error                    { return nil }
func (m *mockAddGitClient) CommitAmendNoEdit() error              { return nil }
func (m *mockAddGitClient) CommitAmendWithMessage(_ string) error { return nil }
func (m *mockAddGitClient) CommitAllowEmpty() error               { return nil }

// Diff Operations methods
func (m *mockAddGitClient) Diff() (string, error)       { return "", nil }
func (m *mockAddGitClient) DiffStaged() (string, error) { return "", nil }
func (m *mockAddGitClient) DiffHead() (string, error)   { return "", nil }
func (m *mockAddGitClient) DiffWith(_ []string) (string, error) {
	return "", nil
}

// Branch Operations methods
func (m *mockAddGitClient) ListLocalBranches() ([]string, error) { return []string{"main"}, nil }
func (m *mockAddGitClient) ListRemoteBranches() ([]string, error) {
	return []string{"origin/main"}, nil
}
func (m *mockAddGitClient) CheckoutNewBranch(_ string) error { return nil }
func (m *mockAddGitClient) CheckoutBranch(_ string) error    { return nil }
func (m *mockAddGitClient) CheckoutNewBranchFromRemote(_, _ string) error {
	return nil
}
func (m *mockAddGitClient) DeleteBranch(_ string) error           { return nil }
func (m *mockAddGitClient) ListMergedBranches() ([]string, error) { return []string{}, nil }

// Remote Operations methods
func (m *mockAddGitClient) Push(_ bool) error              { return nil }
func (m *mockAddGitClient) Pull(_ bool) error              { return nil }
func (m *mockAddGitClient) Fetch(_ bool) error             { return nil }
func (m *mockAddGitClient) RemoteList() error              { return nil }
func (m *mockAddGitClient) RemoteAdd(_, _ string) error    { return nil }
func (m *mockAddGitClient) RemoteRemove(_ string) error    { return nil }
func (m *mockAddGitClient) RemoteSetURL(_, _ string) error { return nil }

// Tag Operations methods
func (m *mockAddGitClient) TagList(_ []string) error              { return nil }
func (m *mockAddGitClient) TagCreate(_, _ string) error           { return nil }
func (m *mockAddGitClient) TagCreateAnnotated(_, _ string) error  { return nil }
func (m *mockAddGitClient) TagDelete(_ []string) error            { return nil }
func (m *mockAddGitClient) TagPush(_, _ string) error             { return nil }
func (m *mockAddGitClient) TagPushAll(_ string) error             { return nil }
func (m *mockAddGitClient) TagShow(_ string) error                { return nil }
func (m *mockAddGitClient) GetLatestTag() (string, error)         { return "v1.0.0", nil }
func (m *mockAddGitClient) TagExists(_ string) bool               { return false }
func (m *mockAddGitClient) GetTagCommit(_ string) (string, error) { return "abc123", nil }

// Log Operations methods
func (m *mockAddGitClient) LogSimple() error { return nil }
func (m *mockAddGitClient) LogGraph() error  { return nil }
func (m *mockAddGitClient) LogOneline(from, to string) (string, error) {
	if m.LogOnelineFunc != nil {
		return m.LogOnelineFunc(from, to)
	}
	return "abc123 First commit\ndef456 Second commit\nghi789 Third commit", nil
}

// Rebase Operations methods
func (m *mockAddGitClient) RebaseInteractive(commitCount int) error {
	if m.RebaseInteractiveFunc != nil {
		return m.RebaseInteractiveFunc(commitCount)
	}
	return nil
}
func (m *mockAddGitClient) Rebase(upstream string) error {
	m.RebaseCalled = true
	m.RebaseUpstream = upstream
	return m.RebaseErr
}
func (m *mockAddGitClient) RebaseContinue() error {
	m.RebaseContinueCalled = true
	return m.RebaseContinueErr
}
func (m *mockAddGitClient) RebaseAbort() error {
	m.RebaseAbortCalled = true
	return m.RebaseAbortErr
}
func (m *mockAddGitClient) RebaseSkip() error {
	m.RebaseSkipCalled = true
	return m.RebaseSkipErr
}
func (m *mockAddGitClient) GetUpstreamBranch(_ string) (string, error) {
	return "origin/main", nil
}

// Stash Operations methods
func (m *mockAddGitClient) Stash() error               { return nil }
func (m *mockAddGitClient) StashList() (string, error) { return "", nil }
func (m *mockAddGitClient) StashShow(_ string) error   { return nil }
func (m *mockAddGitClient) StashApply(_ string) error  { return nil }
func (m *mockAddGitClient) StashPop(_ string) error    { return nil }
func (m *mockAddGitClient) StashDrop(_ string) error   { return nil }
func (m *mockAddGitClient) StashClear() error          { return nil }

// Restore Operations methods
func (m *mockAddGitClient) RestoreWorkingDir(_ ...string) error           { return nil }
func (m *mockAddGitClient) RestoreStaged(_ ...string) error               { return nil }
func (m *mockAddGitClient) RestoreFromCommit(_ string, _ ...string) error { return nil }
func (m *mockAddGitClient) RestoreAll() error                             { return nil }
func (m *mockAddGitClient) RestoreAllStaged() error                       { return nil }

// Reset and Clean Operations methods
func (m *mockAddGitClient) ResetHardAndClean() error         { return nil }
func (m *mockAddGitClient) ResetHard(_ string) error         { return nil }
func (m *mockAddGitClient) CleanFiles() error                { return nil }
func (m *mockAddGitClient) CleanDirs() error                 { return nil }
func (m *mockAddGitClient) CleanDryRun() (string, error)     { return "", nil }
func (m *mockAddGitClient) CleanFilesForce(_ []string) error { return nil }

// Utility Operations methods
func (m *mockAddGitClient) ListFiles() (string, error) { return "", nil }
func (m *mockAddGitClient) GetUpstreamBranchName(_ string) (string, error) {
	return "origin/main", nil
}
func (m *mockAddGitClient) GetAheadBehindCount(_, _ string) (string, error) {
	return "0	0", nil
}
func (m *mockAddGitClient) GetVersion() (string, error)    { return "test-version", nil }
func (m *mockAddGitClient) GetCommitHash() (string, error) { return "test-commit", nil }
func (m *mockAddGitClient) RevParseVerify(ref string) bool {
	if m.RevParseVerifyFunc != nil {
		return m.RevParseVerifyFunc(ref)
	}
	return true
}

// Config Operations
func (m *mockAddGitClient) ConfigGet(_ string) (string, error)       { return "", nil }
func (m *mockAddGitClient) ConfigSet(_, _ string) error              { return nil }
func (m *mockAddGitClient) ConfigGetGlobal(_ string) (string, error) { return "", nil }
func (m *mockAddGitClient) ConfigSetGlobal(_, _ string) error        { return nil }

// Enhanced Branch Operations
func (m *mockAddGitClient) RenameBranch(_, _ string) error      { return nil }
func (m *mockAddGitClient) MoveBranch(_, _ string) error        { return nil }
func (m *mockAddGitClient) SetUpstreamBranch(_, _ string) error { return nil }
func (m *mockAddGitClient) GetBranchInfo(branch string) (*git.BranchInfo, error) {
	bi := &git.BranchInfo{Name: branch}
	return bi, nil
}
func (m *mockAddGitClient) ListBranchesVerbose() ([]git.BranchInfo, error) {
	return []git.BranchInfo{}, nil
}
func (m *mockAddGitClient) SortBranches(_ string) ([]string, error)       { return []string{}, nil }
func (m *mockAddGitClient) BranchesContaining(_ string) ([]string, error) { return []string{}, nil }

func TestAdder_Add_GitAddCalled(t *testing.T) {
	mockClient := &mockAddGitClient{}
	adder := &Adder{
		gitClient:    mockClient,
		outputWriter: &bytes.Buffer{},
	}
	adder.Add([]string{"hoge.txt"})
	if !mockClient.addCalled {
		t.Error("Add was not called")
	}
	if len(mockClient.addFiles) != 1 || mockClient.addFiles[0] != "hoge.txt" {
		t.Errorf("Expected files [hoge.txt], got %v", mockClient.addFiles)
	}
}

func TestAdder_Add_GitAddArgs(t *testing.T) {
	mockClient := &mockAddGitClient{}
	adder := &Adder{
		gitClient:    mockClient,
		outputWriter: &bytes.Buffer{},
	}
	adder.Add([]string{"foo.txt", "bar.txt"})

	if !mockClient.addCalled {
		t.Error("Add was not called")
	}

	wantFiles := []string{"foo.txt", "bar.txt"}
	if len(mockClient.addFiles) != len(wantFiles) {
		t.Errorf("Expected %d files, got %d", len(wantFiles), len(mockClient.addFiles))
		return
	}

	for i, expected := range wantFiles {
		if mockClient.addFiles[i] != expected {
			t.Errorf("Expected file %s at index %d, got %s", expected, i, mockClient.addFiles[i])
		}
	}
}

func TestAdder_Add_RunError_PrintsError(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockAddGitClient{
		addError: errors.New("git add failed"),
	}
	adder := &Adder{
		gitClient:    mockClient,
		outputWriter: &buf,
	}
	adder.Add([]string{"foo.txt"})

	output := buf.String()
	if !strings.Contains(output, "Error") {
		t.Errorf("Error message not output: %s", output)
	}
	if !strings.Contains(output, "git add failed") {
		t.Errorf("Expected error message not found: %s", output)
	}
}

func TestAdder_Add_PatchSubcommand_CallsInteractive(t *testing.T) {
	mockClient := &mockAddGitClient{}
	adder := &Adder{
		gitClient:    mockClient,
		outputWriter: &bytes.Buffer{},
	}
	adder.Add([]string{"patch"})

	if !mockClient.addInteractiveCalled {
		t.Error("AddInteractive was not called")
	}
	if mockClient.addCalled {
		t.Error("Add should not be called for patch subcommand")
	}
}

func TestAdder_Add_PatchSubcommand_Error(t *testing.T) {
	mockClient := &mockAddGitClient{addInteractiveError: errors.New("interactive add failed")}
	var buf bytes.Buffer
	adder := &Adder{
		gitClient:    mockClient,
		outputWriter: &buf,
	}
	adder.Add([]string{"patch"})

	output := buf.String()
	if output == "" || output[:5] != "Error" {
		t.Errorf("Error output not generated with patch subcommand: %s", output)
	}
}

func TestAdder_Add_Interactive(t *testing.T) {
	mockClient := &mockAddGitClient{}
	var buf bytes.Buffer
	adder := &Adder{
		gitClient:    mockClient,
		outputWriter: &buf,
	}

	adder.Add([]string{"patch"})

	// Check that AddInteractive was called
	if !mockClient.addInteractiveCalled {
		t.Error("AddInteractive should be called for patch subcommand")
	}
}

func TestAdder_Add(t *testing.T) {
	cases := []struct {
		name        string
		args        []string
		expectedCmd string
		expectError bool
	}{
		{
			name:        "add all files",
			args:        []string{"."},
			expectedCmd: "git add .",
			expectError: false,
		},
		{
			name:        "add specific file",
			args:        []string{"file.txt"},
			expectedCmd: "git add file.txt",
			expectError: false,
		},
		{
			name:        "add multiple files",
			args:        []string{"file1.txt", "file2.txt"},
			expectedCmd: "git add file1.txt file2.txt",
			expectError: false,
		},
		{
			name:        "no args",
			args:        []string{},
			expectedCmd: "",
			expectError: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &mockAddGitClient{}
			var buf bytes.Buffer
			a := &Adder{
				gitClient:    mockClient,
				outputWriter: &buf,
			}

			a.Add(tc.args)

			// Check output for no args case
			if len(tc.args) == 0 {
				output := buf.String()
				if !strings.Contains(output, "Usage:") {
					t.Errorf("expected usage message, got %q", output)
				}
			}
		})
	}
}

func TestAdder_Add_Error(t *testing.T) {
	mockClient := &mockAddGitClient{addError: errors.New("git add failed")}
	var buf bytes.Buffer
	a := &Adder{
		gitClient:    mockClient,
		outputWriter: &buf,
	}

	a.Add([]string{"file.txt"})

	output := buf.String()
	if !strings.Contains(output, "Error:") {
		t.Errorf("expected error message, got %q", output)
	}
}

func TestAdder_Add_InteractiveSubcommand_CallsInteractive(t *testing.T) {
	mockClient := &mockAddGitClient{}
	var buf bytes.Buffer
	adder := &Adder{
		gitClient:    mockClient,
		outputWriter: &buf,
	}

	adder.Add([]string{"interactive"})

	if !mockClient.addInteractiveCalled {
		t.Error("AddInteractive was not called for 'interactive' subcommand")
	}
	if mockClient.addCalled {
		t.Error("Add should not be called for 'interactive' subcommand")
	}
}
