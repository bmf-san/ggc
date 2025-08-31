package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v4/git"
)

// mockGitClient for commit_test
type mockCommitGitClient struct {
	commitAllowEmptyCalled       bool
	commitCalled                 bool
	commitAmendCalled            bool
	commitAmendNoEditCalled      bool
	commitAmendWithMessageCalled bool
	commitMessage                string
	amendMessage                 string
	err                          error
}

func (m *mockCommitGitClient) CommitAllowEmpty() error {
	m.commitAllowEmptyCalled = true
	return m.err
}

func (m *mockCommitGitClient) Commit(message string) error {
	m.commitCalled = true
	m.commitMessage = message
	return m.err
}

// Implement all other required methods from git.Clienter interface
func (m *mockCommitGitClient) GetCurrentBranch() (string, error)     { return "main", nil }
func (m *mockCommitGitClient) GetBranchName() (string, error)        { return "main", nil }
func (m *mockCommitGitClient) GetGitStatus() (string, error)         { return "", nil }
func (m *mockCommitGitClient) Status() (string, error)               { return "", nil }
func (m *mockCommitGitClient) StatusShort() (string, error)          { return "", nil }
func (m *mockCommitGitClient) StatusWithColor() (string, error)      { return "", nil }
func (m *mockCommitGitClient) StatusShortWithColor() (string, error) { return "", nil }
func (m *mockCommitGitClient) Add(_ ...string) error                 { return nil }
func (m *mockCommitGitClient) AddInteractive() error                 { return nil }
func (m *mockCommitGitClient) CommitAmend() error {
	m.commitAmendCalled = true
	return m.err
}
func (m *mockCommitGitClient) CommitAmendNoEdit() error {
	m.commitAmendNoEditCalled = true
	return m.err
}
func (m *mockCommitGitClient) CommitAmendWithMessage(message string) error {
	m.commitAmendWithMessageCalled = true
	m.amendMessage = message
	return m.err
}
func (m *mockCommitGitClient) Diff() (string, error)                 { return "", nil }
func (m *mockCommitGitClient) DiffStaged() (string, error)           { return "", nil }
func (m *mockCommitGitClient) DiffHead() (string, error)             { return "", nil }
func (m *mockCommitGitClient) ListLocalBranches() ([]string, error)  { return []string{}, nil }
func (m *mockCommitGitClient) ListRemoteBranches() ([]string, error) { return []string{}, nil }
func (m *mockCommitGitClient) CheckoutNewBranch(_ string) error      { return nil }
func (m *mockCommitGitClient) CheckoutBranch(_ string) error         { return nil }
func (m *mockCommitGitClient) CheckoutNewBranchFromRemote(_, _ string) error {
	return nil
}
func (m *mockCommitGitClient) DeleteBranch(_ string) error            { return nil }
func (m *mockCommitGitClient) ListMergedBranches() ([]string, error)  { return []string{}, nil }
func (m *mockCommitGitClient) Push(_ bool) error                      { return nil }
func (m *mockCommitGitClient) Pull(_ bool) error                      { return nil }
func (m *mockCommitGitClient) Fetch(_ bool) error                     { return nil }
func (m *mockCommitGitClient) RemoteList() error                      { return nil }
func (m *mockCommitGitClient) RemoteAdd(_, _ string) error            { return nil }
func (m *mockCommitGitClient) RemoteRemove(_ string) error            { return nil }
func (m *mockCommitGitClient) RemoteSetURL(_, _ string) error         { return nil }
func (m *mockCommitGitClient) LogSimple() error                       { return nil }
func (m *mockCommitGitClient) LogGraph() error                        { return nil }
func (m *mockCommitGitClient) LogOneline(_, _ string) (string, error) { return "", nil }
func (m *mockCommitGitClient) RebaseInteractive(_ int) error          { return nil }
func (m *mockCommitGitClient) GetUpstreamBranch(_ string) (string, error) {
	return "origin/main", nil
}
func (m *mockCommitGitClient) Stash() error                                  { return nil }
func (m *mockCommitGitClient) StashList() (string, error)                    { return "", nil }
func (m *mockCommitGitClient) StashShow(_ string) error                      { return nil }
func (m *mockCommitGitClient) StashApply(_ string) error                     { return nil }
func (m *mockCommitGitClient) StashPop(_ string) error                       { return nil }
func (m *mockCommitGitClient) StashDrop(_ string) error                      { return nil }
func (m *mockCommitGitClient) StashClear() error                             { return nil }
func (m *mockCommitGitClient) RestoreWorkingDir(_ ...string) error           { return nil }
func (m *mockCommitGitClient) RestoreStaged(_ ...string) error               { return nil }
func (m *mockCommitGitClient) RestoreFromCommit(_ string, _ ...string) error { return nil }
func (m *mockCommitGitClient) RestoreAll() error                             { return nil }
func (m *mockCommitGitClient) RestoreAllStaged() error                       { return nil }
func (m *mockCommitGitClient) ResetHardAndClean() error                      { return nil }
func (m *mockCommitGitClient) ResetHard(_ string) error                      { return nil }
func (m *mockCommitGitClient) CleanFiles() error                             { return nil }
func (m *mockCommitGitClient) CleanDirs() error                              { return nil }
func (m *mockCommitGitClient) CleanDryRun() (string, error)                  { return "", nil }
func (m *mockCommitGitClient) CleanFilesForce(_ []string) error              { return nil }
func (m *mockCommitGitClient) TagList(_ []string) error                      { return nil }
func (m *mockCommitGitClient) TagCreate(_, _ string) error                   { return nil }
func (m *mockCommitGitClient) TagCreateAnnotated(_, _ string) error          { return nil }
func (m *mockCommitGitClient) TagDelete(_ []string) error                    { return nil }
func (m *mockCommitGitClient) TagPush(_, _ string) error                     { return nil }
func (m *mockCommitGitClient) TagPushAll(_ string) error                     { return nil }
func (m *mockCommitGitClient) TagShow(_ string) error                        { return nil }
func (m *mockCommitGitClient) GetLatestTag() (string, error)                 { return "", nil }
func (m *mockCommitGitClient) TagExists(_ string) bool                       { return false }
func (m *mockCommitGitClient) GetTagCommit(_ string) (string, error)         { return "abc123", nil }
func (m *mockCommitGitClient) ListFiles() (string, error)                    { return "", nil }
func (m *mockCommitGitClient) GetUpstreamBranchName(_ string) (string, error) {
	return "origin/main", nil
}
func (m *mockCommitGitClient) GetAheadBehindCount(_, _ string) (string, error) {
	return "0\t0", nil
}
func (m *mockCommitGitClient) GetVersion() (string, error)    { return "test-version", nil }
func (m *mockCommitGitClient) GetCommitHash() (string, error) { return "test-commit", nil }
func (m *mockCommitGitClient) RevParseVerify(_ string) bool   { return true }

// Config Operations
func (m *mockCommitGitClient) ConfigGet(_ string) (string, error)       { return "", nil }
func (m *mockCommitGitClient) ConfigSet(_, _ string) error              { return nil }
func (m *mockCommitGitClient) ConfigGetGlobal(_ string) (string, error) { return "", nil }
func (m *mockCommitGitClient) ConfigSetGlobal(_, _ string) error        { return nil }

// Enhanced Branch Operations
func (m *mockCommitGitClient) RenameBranch(_, _ string) error      { return nil }
func (m *mockCommitGitClient) MoveBranch(_, _ string) error        { return nil }
func (m *mockCommitGitClient) SetUpstreamBranch(_, _ string) error { return nil }
func (m *mockCommitGitClient) GetBranchInfo(branch string) (*git.BranchInfo, error) {
	bi := &git.BranchInfo{Name: branch}
	return bi, nil
}
func (m *mockCommitGitClient) ListBranchesVerbose() ([]git.BranchInfo, error) {
	return []git.BranchInfo{}, nil
}
func (m *mockCommitGitClient) SortBranches(_ string) ([]string, error)       { return []string{}, nil }
func (m *mockCommitGitClient) BranchesContaining(_ string) ([]string, error) { return []string{}, nil }

func TestCommitter_Commit_AllowEmpty(t *testing.T) {
	mockClient := &mockCommitGitClient{}
	var buf bytes.Buffer
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"allow-empty"})
	if !mockClient.commitAllowEmptyCalled {
		t.Error("CommitAllowEmpty should be called")
	}
}

func TestCommitter_Commit_Help(t *testing.T) {
	var buf bytes.Buffer
	c := &Committer{
		gitClient:    &mockCommitGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{})

	output := buf.String()
	if output == "" || !strings.Contains(output, "Usage") {
		t.Errorf("Usage should be displayed, but got: %s", output)
	}
}

func TestCommitter_Commit_AllowEmpty_Error(t *testing.T) {
	var buf bytes.Buffer
	c := &Committer{
		gitClient:    &mockCommitGitClient{err: errors.New("fail")},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"allow-empty"})

	output := buf.String()
	if output != "Error: fail\n" {
		t.Errorf("unexpected output: got %q", output)
	}
}

func TestCommitter_Commit_Normal(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockCommitGitClient{}
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"test message"})

	if !mockClient.commitCalled {
		t.Error("git commit command should be called with the correct message")
	}
	if mockClient.commitMessage != "test message" {
		t.Errorf("expected commit message 'test message', got '%s'", mockClient.commitMessage)
	}
}

func TestCommitter_Commit_Normal_WithBrackets(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockCommitGitClient{}
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"[update]", "test", "message"})

	if !mockClient.commitCalled {
		t.Error("git commit command should be called with the correct message including brackets")
	}
	expectedMessage := "[update] test message"
	if mockClient.commitMessage != expectedMessage {
		t.Errorf("expected commit message '%s', got '%s'", expectedMessage, mockClient.commitMessage)
	}
}

func TestCommitter_Commit_Normal_Error(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockCommitGitClient{err: errors.New("commit failed")}
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"test message"})

	output := buf.String()
	if !strings.Contains(output, "Error:") {
		t.Errorf("Expected error message, got: %s", output)
	}
}

func TestCommitter_Commit_Amend_WithMessage(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockCommitGitClient{}
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"amend", "updated message"})
	if !mockClient.commitAmendWithMessageCalled {
		t.Error("git commit --amend -m command should be called")
	}
}

func TestCommitter_Commit_Amend_NoEdit(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockCommitGitClient{}
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"amend", "--no-edit"})
	if !mockClient.commitAmendNoEditCalled {
		t.Error("git commit --amend --no-edit command should be called")
	}
}

func TestCommitter_Commit_Amend_Error(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockCommitGitClient{
		err: errors.New("commit failed"),
	}
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"--amend", "test", "message"})
	output := buf.String()
	if !strings.Contains(output, "Error:") {
		t.Errorf("Expected error message, got: %s", output)
	}
}

func TestCommitter_Commit_Normal_NoSpace(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockCommitGitClient{}
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"[update]hoge"})

	if !mockClient.commitCalled {
		t.Error("git commit command should be called with the correct message without spaces")
	}
	if mockClient.commitMessage != "[update]hoge" {
		t.Errorf("expected commit message '[update]hoge', got '%s'", mockClient.commitMessage)
	}
}

func TestCommitter_Commit_Amend_WithMultiWordMessage(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockCommitGitClient{}
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"amend", "[update]", "message", "with", "spaces"})
	if !mockClient.commitAmendWithMessageCalled {
		t.Error("git commit --amend -m command should be called with the complete message")
	}
}
