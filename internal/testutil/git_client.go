// Package testutil provides testing utilities for the ggc CLI tool.
// This file contains testing utilities for git clients.
//
// This file is excluded from coverage reports as it contains test utilities.
package testutil

import (
	"github.com/bmf-san/ggc/v7/pkg/git"
)

// NewMockGitClient creates a new mock git client for testing
func NewMockGitClient() *testMockGitClient {
	return &testMockGitClient{
		currentBranch: "main",
		gitStatus:     "A  file1.txt\n M file2.txt\n",
		aheadBehind:   "2\t1",
	}
}

// testMockGitClient is a mock git client for testing
type testMockGitClient struct {
	currentBranch string
	gitStatus     string
	aheadBehind   string
}

func (m *testMockGitClient) GetCurrentBranch() (string, error) { return m.currentBranch, nil }
func (m *testMockGitClient) GetBranchName() (string, error)    { return m.currentBranch, nil }
func (m *testMockGitClient) GetGitStatus() (string, error)     { return m.gitStatus, nil }

// Status Operations
func (m *testMockGitClient) Status() (string, error)               { return m.gitStatus, nil }
func (m *testMockGitClient) StatusShort() (string, error)          { return m.gitStatus, nil }
func (m *testMockGitClient) StatusWithColor() (string, error)      { return m.gitStatus, nil }
func (m *testMockGitClient) StatusShortWithColor() (string, error) { return m.gitStatus, nil }

// Staging Operations
func (m *testMockGitClient) Add(_ ...string) error { return nil }
func (m *testMockGitClient) AddInteractive() error { return nil }

// Commit Operations
func (m *testMockGitClient) Commit(_ string) error                 { return nil }
func (m *testMockGitClient) CommitAmend() error                    { return nil }
func (m *testMockGitClient) CommitAmendNoEdit() error              { return nil }
func (m *testMockGitClient) CommitAmendWithMessage(_ string) error { return nil }
func (m *testMockGitClient) CommitAllowEmpty() error               { return nil }

// Diff Operations
func (m *testMockGitClient) Diff() (string, error)       { return "", nil }
func (m *testMockGitClient) DiffStaged() (string, error) { return "", nil }
func (m *testMockGitClient) DiffHead() (string, error)   { return "", nil }
func (m *testMockGitClient) DiffWith(_ []string) (string, error) {
	return "", nil
}

// Branch Operations
func (m *testMockGitClient) ListLocalBranches() ([]string, error) { return []string{"main"}, nil }
func (m *testMockGitClient) ListRemoteBranches() ([]string, error) {
	return []string{"origin/main"}, nil
}
func (m *testMockGitClient) CheckoutNewBranch(_ string) error              { return nil }
func (m *testMockGitClient) CheckoutBranch(_ string) error                 { return nil }
func (m *testMockGitClient) CheckoutNewBranchFromRemote(_, _ string) error { return nil }
func (m *testMockGitClient) DeleteBranch(_ string) error                   { return nil }
func (m *testMockGitClient) ListMergedBranches() ([]string, error)         { return []string{}, nil }
func (m *testMockGitClient) RevParseVerify(_ string) bool                  { return true }

// Remote Operations
func (m *testMockGitClient) Push(_ bool) error              { return nil }
func (m *testMockGitClient) Pull(_ bool) error              { return nil }
func (m *testMockGitClient) Fetch(_ bool) error             { return nil }
func (m *testMockGitClient) RemoteList() error              { return nil }
func (m *testMockGitClient) RemoteAdd(_, _ string) error    { return nil }
func (m *testMockGitClient) RemoteRemove(_ string) error    { return nil }
func (m *testMockGitClient) RemoteSetURL(_, _ string) error { return nil }

// Tag Operations
func (m *testMockGitClient) TagList(_ []string) error              { return nil }
func (m *testMockGitClient) TagCreate(_, _ string) error           { return nil }
func (m *testMockGitClient) TagCreateAnnotated(_, _ string) error  { return nil }
func (m *testMockGitClient) TagDelete(_ []string) error            { return nil }
func (m *testMockGitClient) TagPush(_, _ string) error             { return nil }
func (m *testMockGitClient) TagPushAll(_ string) error             { return nil }
func (m *testMockGitClient) TagShow(_ string) error                { return nil }
func (m *testMockGitClient) GetLatestTag() (string, error)         { return "v1.0.0", nil }
func (m *testMockGitClient) TagExists(_ string) bool               { return true }
func (m *testMockGitClient) GetTagCommit(_ string) (string, error) { return "abc123", nil }

// Log Operations
func (m *testMockGitClient) LogSimple() error                       { return nil }
func (m *testMockGitClient) LogGraph() error                        { return nil }
func (m *testMockGitClient) LogOneline(_, _ string) (string, error) { return "", nil }

// Rebase Operations
func (m *testMockGitClient) RebaseInteractive(_ int) error              { return nil }
func (m *testMockGitClient) Rebase(_ string) error                      { return nil }
func (m *testMockGitClient) RebaseContinue() error                      { return nil }
func (m *testMockGitClient) RebaseAbort() error                         { return nil }
func (m *testMockGitClient) RebaseSkip() error                          { return nil }
func (m *testMockGitClient) GetUpstreamBranch(_ string) (string, error) { return "origin/main", nil }

// Stash Operations
func (m *testMockGitClient) Stash() error               { return nil }
func (m *testMockGitClient) StashList() (string, error) { return "", nil }
func (m *testMockGitClient) StashShow(_ string) error   { return nil }
func (m *testMockGitClient) StashApply(_ string) error  { return nil }
func (m *testMockGitClient) StashPop(_ string) error    { return nil }
func (m *testMockGitClient) StashPush(_ string) error   { return nil }
func (m *testMockGitClient) StashDrop(_ string) error   { return nil }
func (m *testMockGitClient) StashClear() error          { return nil }

// Restore Operations
func (m *testMockGitClient) RestoreWorkingDir(_ ...string) error           { return nil }
func (m *testMockGitClient) RestoreStaged(_ ...string) error               { return nil }
func (m *testMockGitClient) RestoreFromCommit(_ string, _ ...string) error { return nil }
func (m *testMockGitClient) RestoreAll() error                             { return nil }
func (m *testMockGitClient) RestoreAllStaged() error                       { return nil }

// Config Operations
func (m *testMockGitClient) ConfigGet(_ string) (string, error)       { return "", nil }
func (m *testMockGitClient) ConfigSet(_, _ string) error              { return nil }
func (m *testMockGitClient) ConfigGetGlobal(_ string) (string, error) { return "", nil }
func (m *testMockGitClient) ConfigSetGlobal(_, _ string) error        { return nil }

// Reset Operations
func (m *testMockGitClient) ResetHardAndClean() error { return nil }
func (m *testMockGitClient) ResetHard(_ string) error { return nil }

// Clean Operations
func (m *testMockGitClient) CleanFiles() error                { return nil }
func (m *testMockGitClient) CleanDirs() error                 { return nil }
func (m *testMockGitClient) CleanDryRun() (string, error)     { return "", nil }
func (m *testMockGitClient) CleanFilesForce(_ []string) error { return nil }

// Utility Operations
func (m *testMockGitClient) ListFiles() (string, error) { return "", nil }
func (m *testMockGitClient) GetUpstreamBranchName(_ string) (string, error) {
	return "origin/main", nil
}
func (m *testMockGitClient) GetAheadBehindCount(_, _ string) (string, error) {
	return m.aheadBehind, nil
}
func (m *testMockGitClient) GetVersion() (string, error)    { return "test-version", nil }
func (m *testMockGitClient) GetCommitHash() (string, error) { return "test-commit", nil }
func (m *testMockGitClient) BranchesContaining(_ string) ([]string, error) {
	return []string{"main", "develop"}, nil
}
func (m *testMockGitClient) GetBranchInfo(_ string) (*git.BranchInfo, error) {
	return &git.BranchInfo{
		Name:            "main",
		IsCurrentBranch: true,
		Upstream:        "origin/main",
		AheadBehind:     "up to date",
		LastCommitSHA:   "abc123",
		LastCommitMsg:   "test commit",
	}, nil
}
func (m *testMockGitClient) ListBranchesVerbose() ([]git.BranchInfo, error) {
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
func (m *testMockGitClient) MoveBranch(_, _ string) error            { return nil }
func (m *testMockGitClient) RenameBranch(_, _ string) error          { return nil }
func (m *testMockGitClient) SetUpstreamBranch(_, _ string) error     { return nil }
func (m *testMockGitClient) SortBranches(_ string) ([]string, error) { return []string{"main"}, nil }
