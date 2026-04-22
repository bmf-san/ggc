// Package testutil provides testing utilities for the ggc CLI tool.
// This file contains testing utilities for git clients.
//
// This file is excluded from coverage reports as it contains test utilities.
package testutil

import (
	"github.com/bmf-san/ggc/v8/internal/git"
)

// NewMockGitClient creates a new mock git client for testing.
// For custom behavior, embed MockGitClient in a wrapper struct and override
// the specific methods needed; Go's method resolution will prefer the
// outer struct's methods. Example:
//
//	type myMock struct { *testutil.MockGitClient; calls int }
//	func (m *myMock) Push(force bool) error { m.calls++; return nil }
func NewMockGitClient() *MockGitClient {
	return &MockGitClient{
		currentBranch: "main",
		gitStatus:     "A  file1.txt\n M file2.txt\n",
		aheadBehind:   "2\t1",
	}
}

// MockGitClient is a mock git client for testing. It implements the full git
// client surface with no-op / default-value behavior. Tests may embed this
// type to override only the methods that matter to them.
type MockGitClient struct {
	currentBranch string
	gitStatus     string
	aheadBehind   string
}

func (m *MockGitClient) GetCurrentBranch() (string, error) { return m.currentBranch, nil }
func (m *MockGitClient) GetBranchName() (string, error)    { return m.currentBranch, nil }

// Status Operations
func (m *MockGitClient) Status() (string, error)               { return m.gitStatus, nil }
func (m *MockGitClient) StatusShort() (string, error)          { return m.gitStatus, nil }
func (m *MockGitClient) StatusWithColor() (string, error)      { return m.gitStatus, nil }
func (m *MockGitClient) StatusShortWithColor() (string, error) { return m.gitStatus, nil }

// Staging Operations
func (m *MockGitClient) Add(_ ...string) error { return nil }
func (m *MockGitClient) AddInteractive() error { return nil }

// Commit Operations
func (m *MockGitClient) Commit(_ string) error                 { return nil }
func (m *MockGitClient) CommitAmend() error                    { return nil }
func (m *MockGitClient) CommitAmendNoEdit() error              { return nil }
func (m *MockGitClient) CommitAmendWithMessage(_ string) error { return nil }
func (m *MockGitClient) CommitAllowEmpty() error               { return nil }
func (m *MockGitClient) CommitFixup(_ string) error            { return nil }

// Diff Operations
func (m *MockGitClient) Diff() (string, error)       { return "", nil }
func (m *MockGitClient) DiffStaged() (string, error) { return "", nil }
func (m *MockGitClient) DiffHead() (string, error)   { return "", nil }
func (m *MockGitClient) DiffWith(_ []string) (string, error) {
	return "", nil
}

// Branch Operations
func (m *MockGitClient) ListLocalBranches() ([]string, error) { return []string{"main"}, nil }
func (m *MockGitClient) ListRemoteBranches() ([]string, error) {
	return []string{"origin/main"}, nil
}
func (m *MockGitClient) CheckoutNewBranch(_ string) error              { return nil }
func (m *MockGitClient) CheckoutBranch(_ string) error                 { return nil }
func (m *MockGitClient) CheckoutNewBranchFromRemote(_, _ string) error { return nil }
func (m *MockGitClient) DeleteBranch(_ string) error                   { return nil }
func (m *MockGitClient) ListMergedBranches() ([]string, error)         { return []string{}, nil }
func (m *MockGitClient) RevParseVerify(_ string) bool                  { return true }

// Remote Operations
func (m *MockGitClient) Push(_ bool) error              { return nil }
func (m *MockGitClient) Pull(_ bool) error              { return nil }
func (m *MockGitClient) Fetch(_ bool) error             { return nil }
func (m *MockGitClient) RemoteList() error              { return nil }
func (m *MockGitClient) RemoteAdd(_, _ string) error    { return nil }
func (m *MockGitClient) RemoteRemove(_ string) error    { return nil }
func (m *MockGitClient) RemoteSetURL(_, _ string) error { return nil }

// Tag Operations
func (m *MockGitClient) TagList(_ []string) error              { return nil }
func (m *MockGitClient) TagCreate(_, _ string) error           { return nil }
func (m *MockGitClient) TagCreateAnnotated(_, _ string) error  { return nil }
func (m *MockGitClient) TagDelete(_ []string) error            { return nil }
func (m *MockGitClient) TagPush(_, _ string) error             { return nil }
func (m *MockGitClient) TagPushAll(_ string) error             { return nil }
func (m *MockGitClient) TagShow(_ string) error                { return nil }
func (m *MockGitClient) GetLatestTag() (string, error)         { return "v1.0.0", nil }
func (m *MockGitClient) TagExists(_ string) bool               { return true }
func (m *MockGitClient) GetTagCommit(_ string) (string, error) { return "abc123", nil }

// Log Operations
func (m *MockGitClient) LogSimple() error                       { return nil }
func (m *MockGitClient) LogGraph() error                        { return nil }
func (m *MockGitClient) LogOneline(_, _ string) (string, error) { return "", nil }

// Rebase Operations
func (m *MockGitClient) RebaseInteractive(_ int) error              { return nil }
func (m *MockGitClient) RebaseInteractiveAutosquash(_ int) error    { return nil }
func (m *MockGitClient) Rebase(_ string) error                      { return nil }
func (m *MockGitClient) RebaseContinue() error                      { return nil }
func (m *MockGitClient) RebaseAbort() error                         { return nil }
func (m *MockGitClient) RebaseSkip() error                          { return nil }
func (m *MockGitClient) GetUpstreamBranch(_ string) (string, error) { return "origin/main", nil }

// Stash Operations
func (m *MockGitClient) Stash() error               { return nil }
func (m *MockGitClient) StashList() (string, error) { return "", nil }
func (m *MockGitClient) StashShow(_ string) error   { return nil }
func (m *MockGitClient) StashApply(_ string) error  { return nil }
func (m *MockGitClient) StashPop(_ string) error    { return nil }
func (m *MockGitClient) StashPush(_ string) error   { return nil }
func (m *MockGitClient) StashDrop(_ string) error   { return nil }
func (m *MockGitClient) StashClear() error          { return nil }

// Restore Operations
func (m *MockGitClient) RestoreWorkingDir(_ ...string) error           { return nil }
func (m *MockGitClient) RestoreStaged(_ ...string) error               { return nil }
func (m *MockGitClient) RestoreFromCommit(_ string, _ ...string) error { return nil }
func (m *MockGitClient) RestoreAll() error                             { return nil }
func (m *MockGitClient) RestoreAllStaged() error                       { return nil }

// Config Operations
func (m *MockGitClient) ConfigGet(_ string) (string, error)       { return "", nil }
func (m *MockGitClient) ConfigSet(_, _ string) error              { return nil }
func (m *MockGitClient) ConfigGetGlobal(_ string) (string, error) { return "", nil }
func (m *MockGitClient) ConfigSetGlobal(_, _ string) error        { return nil }

// Reset Operations
func (m *MockGitClient) ResetHardAndClean() error { return nil }
func (m *MockGitClient) ResetHard(_ string) error { return nil }
func (m *MockGitClient) ResetSoft(_ string) error { return nil }

// Clean Operations
func (m *MockGitClient) CleanFiles() error                { return nil }
func (m *MockGitClient) CleanDirs() error                 { return nil }
func (m *MockGitClient) CleanDryRun() (string, error)     { return "", nil }
func (m *MockGitClient) CleanFilesForce(_ []string) error { return nil }

// Utility Operations
func (m *MockGitClient) ListFiles() (string, error) { return "", nil }
func (m *MockGitClient) GetUpstreamBranchName(_ string) (string, error) {
	return "origin/main", nil
}
func (m *MockGitClient) GetAheadBehindCount(_, _ string) (string, error) {
	return m.aheadBehind, nil
}
func (m *MockGitClient) GetVersion() (string, error)    { return "test-version", nil }
func (m *MockGitClient) GetCommitHash() (string, error) { return "test-commit", nil }
func (m *MockGitClient) BranchesContaining(_ string) ([]string, error) {
	return []string{"main", "develop"}, nil
}
func (m *MockGitClient) GetBranchInfo(_ string) (*git.BranchInfo, error) {
	return &git.BranchInfo{
		Name:            "main",
		IsCurrentBranch: true,
		Upstream:        "origin/main",
		AheadBehind:     "up to date",
		LastCommitSHA:   "abc123",
		LastCommitMsg:   "test commit",
	}, nil
}
func (m *MockGitClient) ListBranchesVerbose() ([]git.BranchInfo, error) {
	return []git.BranchInfo{
		{
			Name:            "main",
			IsCurrentBranch: true,
			Upstream:        "origin/main",
			AheadBehind:     "up to date",
			LastCommitSHA:   "abc123",
			LastCommitMsg:   "test commit",
		},
	}, nil
}

// Additional missing methods
func (m *MockGitClient) MoveBranch(_, _ string) error            { return nil }
func (m *MockGitClient) RenameBranch(_, _ string) error          { return nil }
func (m *MockGitClient) SetUpstreamBranch(_, _ string) error     { return nil }
func (m *MockGitClient) SortBranches(_ string) ([]string, error) { return []string{"main"}, nil }
func (m *MockGitClient) ValidateBranchName(_ string) error       { return nil }
