package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
	"testing"
)

// mockGitClient for clean_test
type mockCleanGitClient struct {
	cleanFilesErr     error
	cleanDirsErr      error
	cleanFilesCalled  bool
	cleanDirsCalled   bool
	cleanDryRunResult string
	cleanDryRunErr    error
}

func (m *mockCleanGitClient) CleanFiles() error {
	m.cleanFilesCalled = true
	return m.cleanFilesErr
}

func (m *mockCleanGitClient) CleanDirs() error {
	m.cleanDirsCalled = true
	return m.cleanDirsErr
}

func (m *mockCleanGitClient) CleanDryRun() (string, error) {
	return m.cleanDryRunResult, m.cleanDryRunErr
}

func (m *mockCleanGitClient) CleanFilesForce(_ []string) error {
	return nil
}

// Implement all other required methods from git.Clienter interface
func (m *mockCleanGitClient) GetCurrentBranch() (string, error)     { return "main", nil }
func (m *mockCleanGitClient) GetBranchName() (string, error)        { return "main", nil }
func (m *mockCleanGitClient) GetGitStatus() (string, error)         { return "", nil }
func (m *mockCleanGitClient) Status() (string, error)               { return "", nil }
func (m *mockCleanGitClient) StatusShort() (string, error)          { return "", nil }
func (m *mockCleanGitClient) StatusWithColor() (string, error)      { return "", nil }
func (m *mockCleanGitClient) StatusShortWithColor() (string, error) { return "", nil }
func (m *mockCleanGitClient) Add(_ ...string) error                 { return nil }
func (m *mockCleanGitClient) AddInteractive() error                 { return nil }
func (m *mockCleanGitClient) Commit(_ string) error                 { return nil }
func (m *mockCleanGitClient) CommitAmend() error                    { return nil }
func (m *mockCleanGitClient) CommitAmendNoEdit() error              { return nil }
func (m *mockCleanGitClient) CommitAmendWithMessage(_ string) error { return nil }
func (m *mockCleanGitClient) CommitAllowEmpty() error               { return nil }
func (m *mockCleanGitClient) Diff() (string, error)                 { return "", nil }
func (m *mockCleanGitClient) DiffStaged() (string, error)           { return "", nil }
func (m *mockCleanGitClient) DiffHead() (string, error)             { return "", nil }
func (m *mockCleanGitClient) ListLocalBranches() ([]string, error)  { return []string{}, nil }
func (m *mockCleanGitClient) ListRemoteBranches() ([]string, error) { return []string{}, nil }
func (m *mockCleanGitClient) CheckoutNewBranch(_ string) error      { return nil }
func (m *mockCleanGitClient) CheckoutBranch(_ string) error         { return nil }
func (m *mockCleanGitClient) CheckoutNewBranchFromRemote(_, _ string) error {
	return nil
}
func (m *mockCleanGitClient) DeleteBranch(_ string) error            { return nil }
func (m *mockCleanGitClient) ListMergedBranches() ([]string, error)  { return []string{}, nil }
func (m *mockCleanGitClient) Push(_ bool) error                      { return nil }
func (m *mockCleanGitClient) Pull(_ bool) error                      { return nil }
func (m *mockCleanGitClient) Fetch(_ bool) error                     { return nil }
func (m *mockCleanGitClient) RemoteList() error                      { return nil }
func (m *mockCleanGitClient) RemoteAdd(_, _ string) error            { return nil }
func (m *mockCleanGitClient) RemoteRemove(_ string) error            { return nil }
func (m *mockCleanGitClient) RemoteSetURL(_, _ string) error         { return nil }
func (m *mockCleanGitClient) LogSimple() error                       { return nil }
func (m *mockCleanGitClient) LogGraph() error                        { return nil }
func (m *mockCleanGitClient) LogOneline(_, _ string) (string, error) { return "", nil }
func (m *mockCleanGitClient) RebaseInteractive(_ int) error          { return nil }
func (m *mockCleanGitClient) GetUpstreamBranch(_ string) (string, error) {
	return "origin/main", nil
}
func (m *mockCleanGitClient) Stash() error                                  { return nil }
func (m *mockCleanGitClient) StashList() (string, error)                    { return "", nil }
func (m *mockCleanGitClient) StashShow(_ string) error                      { return nil }
func (m *mockCleanGitClient) StashApply(_ string) error                     { return nil }
func (m *mockCleanGitClient) StashPop(_ string) error                       { return nil }
func (m *mockCleanGitClient) StashDrop(_ string) error                      { return nil }
func (m *mockCleanGitClient) StashClear() error                             { return nil }
func (m *mockCleanGitClient) RestoreWorkingDir(_ ...string) error           { return nil }
func (m *mockCleanGitClient) RestoreStaged(_ ...string) error               { return nil }
func (m *mockCleanGitClient) RestoreFromCommit(_ string, _ ...string) error { return nil }
func (m *mockCleanGitClient) RestoreAll() error                             { return nil }
func (m *mockCleanGitClient) RestoreAllStaged() error                       { return nil }
func (m *mockCleanGitClient) ResetHardAndClean() error                      { return nil }
func (m *mockCleanGitClient) ResetHard(_ string) error                      { return nil }
func (m *mockCleanGitClient) TagList(_ []string) error                      { return nil }
func (m *mockCleanGitClient) TagCreate(_, _ string) error                   { return nil }
func (m *mockCleanGitClient) TagCreateAnnotated(_, _ string) error          { return nil }
func (m *mockCleanGitClient) TagDelete(_ []string) error                    { return nil }
func (m *mockCleanGitClient) TagPush(_, _ string) error                     { return nil }
func (m *mockCleanGitClient) TagPushAll(_ string) error                     { return nil }
func (m *mockCleanGitClient) TagShow(_ string) error                        { return nil }
func (m *mockCleanGitClient) GetLatestTag() (string, error)                 { return "", nil }
func (m *mockCleanGitClient) TagExists(_ string) bool                       { return false }
func (m *mockCleanGitClient) GetTagCommit(_ string) (string, error)         { return "abc123", nil }
func (m *mockCleanGitClient) ListFiles() (string, error)                    { return "", nil }
func (m *mockCleanGitClient) GetUpstreamBranchName(_ string) (string, error) {
	return "origin/main", nil
}
func (m *mockCleanGitClient) GetAheadBehindCount(_, _ string) (string, error) {
	return "0\t0", nil
}
func (m *mockCleanGitClient) GetVersion() (string, error)    { return "test-version", nil }
func (m *mockCleanGitClient) GetCommitHash() (string, error) { return "test-commit", nil }
func (m *mockCleanGitClient) RevParseVerify(_ string) bool   { return true }

// Config Operations
func (m *mockCleanGitClient) ConfigGet(_ string) (string, error)       { return "", nil }
func (m *mockCleanGitClient) ConfigSet(_, _ string) error              { return nil }
func (m *mockCleanGitClient) ConfigGetGlobal(_ string) (string, error) { return "", nil }
func (m *mockCleanGitClient) ConfigSetGlobal(_, _ string) error        { return nil }

func TestCleaner_Clean(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		wantCleanFiles bool
		wantCleanDirs  bool
	}{
		{
			name:           "clean files",
			args:           []string{"files"},
			wantCleanFiles: true,
			wantCleanDirs:  false,
		},
		{
			name:           "clean dirs",
			args:           []string{"dirs"},
			wantCleanFiles: false,
			wantCleanDirs:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockCleanGitClient{}
			var buf bytes.Buffer
			cleaner := NewCleanerWithClient(mockClient)
			cleaner.outputWriter = &buf
			cleaner.Clean(tt.args)

			if mockClient.cleanFilesCalled != tt.wantCleanFiles {
				t.Errorf("CleanFiles called = %v, want %v", mockClient.cleanFilesCalled, tt.wantCleanFiles)
			}
			if mockClient.cleanDirsCalled != tt.wantCleanDirs {
				t.Errorf("CleanDirs called = %v, want %v", mockClient.cleanDirsCalled, tt.wantCleanDirs)
			}
		})
	}
}

func TestCleaner_Clean_Files(t *testing.T) {
	var buf bytes.Buffer
	mock := &mockCleanGitClient{}
	cleaner := &Cleaner{
		gitClient:    mock,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	cleaner.helper.outputWriter = &buf

	cleaner.Clean([]string{"files"})

	if buf.Len() > 0 {
		t.Errorf("Expected no output, got %q", buf.String())
	}
}

func TestCleaner_Clean_Files_Error(t *testing.T) {
	var buf bytes.Buffer
	mock := &mockCleanGitClient{cleanFilesErr: errors.New("failed to clean files")}
	cleaner := &Cleaner{
		gitClient:    mock,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	cleaner.helper.outputWriter = &buf

	cleaner.Clean([]string{"files"})

	expected := "Error: failed to clean files\n"
	if got := buf.String(); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

func TestCleaner_Clean_Dirs(t *testing.T) {
	var buf bytes.Buffer
	mock := &mockCleanGitClient{}
	cleaner := &Cleaner{
		gitClient:    mock,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	cleaner.helper.outputWriter = &buf

	cleaner.Clean([]string{"dirs"})

	if buf.Len() > 0 {
		t.Errorf("Expected no output, got %q", buf.String())
	}
}

func TestCleaner_Clean_Dirs_Error(t *testing.T) {
	var buf bytes.Buffer
	mock := &mockCleanGitClient{cleanDirsErr: errors.New("failed to clean directories")}
	cleaner := &Cleaner{
		gitClient:    mock,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	cleaner.helper.outputWriter = &buf

	cleaner.Clean([]string{"dirs"})

	expected := "Error: failed to clean directories\n"
	if got := buf.String(); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

func TestCleaner_Clean_Help(t *testing.T) {
	var buf bytes.Buffer
	cleaner := &Cleaner{
		gitClient:    &mockCleanGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	cleaner.helper.outputWriter = &buf

	cleaner.Clean([]string{})

	output := buf.String()
	if output == "" || !bytes.Contains(buf.Bytes(), []byte("Usage")) {
		t.Errorf("Usage should be displayed, but got: %s", output)
	}
}

func TestCleaner_CleanInteractive_NoFiles(t *testing.T) {
	var buf bytes.Buffer
	cleaner := &Cleaner{
		gitClient:    &mockCleanGitClient{},
		outputWriter: &buf,

		helper: NewHelper(),
	}
	cleaner.helper.outputWriter = &buf

	cleaner.CleanInteractive()

	expected := "No files to clean.\n"
	if got := buf.String(); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

func TestCleaner_CleanInteractive_WithFiles(t *testing.T) {
	var buf bytes.Buffer
	inputBuf := strings.NewReader("all\n")
	cleaner := &Cleaner{
		gitClient:    &mockCleanGitClient{cleanDryRunResult: "Would remove file1.txt\nWould remove file2.txt\n"},
		outputWriter: &buf,

		inputReader: bufio.NewReader(inputBuf),
		helper:      NewHelper(),
	}
	cleaner.helper.outputWriter = &buf

	cleaner.CleanInteractive()

	if !strings.Contains(buf.String(), "Selected files deleted.") {
		t.Error("Expected output to contain 'Selected files deleted.'")
	}
}

func TestCleaner_CleanInteractive_Error(t *testing.T) {
	var buf bytes.Buffer
	cleaner := &Cleaner{
		gitClient:    &mockCleanGitClient{cleanDryRunErr: errors.New("failed to get candidates with git clean -nd")},
		outputWriter: &buf,

		helper: NewHelper(),
	}
	cleaner.helper.outputWriter = &buf

	cleaner.CleanInteractive()

	expected := "Error: failed to get candidates with git clean -nd"
	if !strings.Contains(buf.String(), expected) {
		t.Errorf("Expected output to contain %q, got %q", expected, buf.String())
	}
}

func TestCleaner_CleanInteractive_Cancel(t *testing.T) {
	var buf bytes.Buffer
	inputBuf := strings.NewReader("\n")
	cleaner := &Cleaner{
		gitClient:    &mockCleanGitClient{cleanDryRunResult: "Would remove file1.txt\nWould remove file2.txt\n"},
		outputWriter: &buf,

		inputReader: bufio.NewReader(inputBuf),
		helper:      NewHelper(),
	}
	cleaner.helper.outputWriter = &buf

	cleaner.CleanInteractive()

	if !strings.Contains(buf.String(), "Cancelled.") {
		t.Error("Expected output to contain 'Cancelled.'")
	}
}

func TestCleaner_CleanInteractive_InvalidNumber(t *testing.T) {
	var buf bytes.Buffer
	inputBuf := strings.NewReader("invalid\nnone\nall\n")
	cleaner := &Cleaner{
		gitClient:    &mockCleanGitClient{cleanDryRunResult: "Would remove file1.txt\nWould remove file2.txt\n"},
		outputWriter: &buf,

		inputReader: bufio.NewReader(inputBuf),
		helper:      NewHelper(),
	}
	cleaner.helper.outputWriter = &buf

	cleaner.CleanInteractive()

	if !strings.Contains(buf.String(), "Invalid number: invalid") {
		t.Error("Expected output to contain 'Invalid number: invalid'")
	}
}

func TestCleaner_CleanInteractive_EmptySelection(t *testing.T) {
	var buf bytes.Buffer
	inputBuf := strings.NewReader("\nall\n")
	cleaner := &Cleaner{
		gitClient:    &mockCleanGitClient{cleanDryRunResult: "Would remove file1.txt\nWould remove file2.txt\n"},
		outputWriter: &buf,

		inputReader: bufio.NewReader(inputBuf),
		helper:      NewHelper(),
	}
	cleaner.helper.outputWriter = &buf

	cleaner.CleanInteractive()

	if !strings.Contains(buf.String(), "Cancelled.") {
		t.Error("Expected output to contain 'Cancelled.' for empty input")
	}
}

func TestCleaner_CleanInteractive_FileRejection(t *testing.T) {
	var buf bytes.Buffer
	inputBuf := strings.NewReader("1\nn\nall\n")
	cleaner := &Cleaner{
		gitClient:    &mockCleanGitClient{cleanDryRunResult: "Would remove file1.txt\nWould remove file2.txt\n"},
		outputWriter: &buf,

		inputReader: bufio.NewReader(inputBuf),
		helper:      NewHelper(),
	}
	cleaner.helper.outputWriter = &buf

	cleaner.CleanInteractive()

	output := buf.String()
	if !strings.Contains(output, "Delete these files? (y/n):") {
		t.Error("Expected output to contain 'Delete these files? (y/n):'")
	}
	if !strings.Contains(output, "Selected files deleted.") {
		t.Error("Expected final deletion to succeed")
	}
}

func TestCleaner_CleanInteractive_NothingSelected(t *testing.T) {
	var buf bytes.Buffer
	// Simulate entering an out-of-range number, which results in no actual selection
	inputBuf := strings.NewReader("10\nall\n")
	cleaner := &Cleaner{
		gitClient:    &mockCleanGitClient{cleanDryRunResult: "Would remove file1.txt\nWould remove file2.txt\n"},
		outputWriter: &buf,

		inputReader: bufio.NewReader(inputBuf),
		helper:      NewHelper(),
	}
	cleaner.helper.outputWriter = &buf

	cleaner.CleanInteractive()

	if !strings.Contains(buf.String(), "Invalid number: 10") {
		t.Error("Expected output to contain 'Invalid number: 10'")
	}
}
