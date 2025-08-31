// Package cmd provides command implementations for the ggc CLI tool.
// This file contains testing utilities for git clients.
package cmd

import (
	"github.com/bmf-san/ggc/v4/git"
)

// NewMockGitClient creates a new mock git client for testing
func NewMockGitClient() git.Clienter {
	return &TestMockGitClient{
		currentBranch: "main",
		gitStatus:     "",
		aheadBehind:   "0\t0",
	}
}

// TestMockGitClient is a mock git client for testing
type TestMockGitClient struct {
	currentBranch string
	gitStatus     string
	aheadBehind   string
}

// Repository Information
func (m *TestMockGitClient) GetCurrentBranch() (string, error) { return m.currentBranch, nil }
func (m *TestMockGitClient) GetBranchName() (string, error)    { return m.currentBranch, nil }
func (m *TestMockGitClient) GetGitStatus() (string, error)     { return m.gitStatus, nil }

// Status Operations
func (m *TestMockGitClient) Status() (string, error)               { return "", nil }
func (m *TestMockGitClient) StatusShort() (string, error)          { return "", nil }
func (m *TestMockGitClient) StatusWithColor() (string, error)      { return "", nil }
func (m *TestMockGitClient) StatusShortWithColor() (string, error) { return "", nil }

// Staging Operations
func (m *TestMockGitClient) Add(_ ...string) error { return nil }
func (m *TestMockGitClient) AddInteractive() error { return nil }

// Commit Operations
func (m *TestMockGitClient) Commit(_ string) error                 { return nil }
func (m *TestMockGitClient) CommitAmend() error                    { return nil }
func (m *TestMockGitClient) CommitAmendNoEdit() error              { return nil }
func (m *TestMockGitClient) CommitAmendWithMessage(_ string) error { return nil }
func (m *TestMockGitClient) CommitAllowEmpty() error               { return nil }

// Diff Operations
func (m *TestMockGitClient) Diff() (string, error)       { return "", nil }
func (m *TestMockGitClient) DiffStaged() (string, error) { return "", nil }
func (m *TestMockGitClient) DiffHead() (string, error)   { return "", nil }

// Branch Operations
func (m *TestMockGitClient) ListLocalBranches() ([]string, error) { return []string{"main"}, nil }
func (m *TestMockGitClient) ListRemoteBranches() ([]string, error) {
	return []string{"origin/main"}, nil
}
func (m *TestMockGitClient) CheckoutNewBranch(_ string) error              { return nil }
func (m *TestMockGitClient) CheckoutBranch(_ string) error                 { return nil }
func (m *TestMockGitClient) CheckoutNewBranchFromRemote(_, _ string) error { return nil }
func (m *TestMockGitClient) DeleteBranch(_ string) error                   { return nil }
func (m *TestMockGitClient) ListMergedBranches() ([]string, error)         { return []string{}, nil }
func (m *TestMockGitClient) RevParseVerify(_ string) bool                  { return true }

// Remote Operations
func (m *TestMockGitClient) Push(_ bool) error              { return nil }
func (m *TestMockGitClient) Pull(_ bool) error              { return nil }
func (m *TestMockGitClient) Fetch(_ bool) error             { return nil }
func (m *TestMockGitClient) RemoteList() error              { return nil }
func (m *TestMockGitClient) RemoteAdd(_, _ string) error    { return nil }
func (m *TestMockGitClient) RemoteRemove(_ string) error    { return nil }
func (m *TestMockGitClient) RemoteSetURL(_, _ string) error { return nil }

// Tag Operations
func (m *TestMockGitClient) TagList(_ []string) error              { return nil }
func (m *TestMockGitClient) TagCreate(_, _ string) error           { return nil }
func (m *TestMockGitClient) TagCreateAnnotated(_, _ string) error  { return nil }
func (m *TestMockGitClient) TagDelete(_ []string) error            { return nil }
func (m *TestMockGitClient) TagPush(_, _ string) error             { return nil }
func (m *TestMockGitClient) TagPushAll(_ string) error             { return nil }
func (m *TestMockGitClient) TagShow(_ string) error                { return nil }
func (m *TestMockGitClient) GetLatestTag() (string, error)         { return "v1.0.0", nil }
func (m *TestMockGitClient) TagExists(_ string) bool               { return true }
func (m *TestMockGitClient) GetTagCommit(_ string) (string, error) { return "abc123", nil }

// Log Operations
func (m *TestMockGitClient) LogSimple() error                       { return nil }
func (m *TestMockGitClient) LogGraph() error                        { return nil }
func (m *TestMockGitClient) LogOneline(_, _ string) (string, error) { return "", nil }

// Rebase Operations
func (m *TestMockGitClient) RebaseInteractive(_ int) error              { return nil }
func (m *TestMockGitClient) GetUpstreamBranch(_ string) (string, error) { return "origin/main", nil }

// Stash Operations
func (m *TestMockGitClient) Stash() error               { return nil }
func (m *TestMockGitClient) StashList() (string, error) { return "", nil }
func (m *TestMockGitClient) StashShow(_ string) error   { return nil }
func (m *TestMockGitClient) StashApply(_ string) error  { return nil }
func (m *TestMockGitClient) StashPop(_ string) error    { return nil }
func (m *TestMockGitClient) StashDrop(_ string) error   { return nil }
func (m *TestMockGitClient) StashClear() error          { return nil }

// Restore Operations
func (m *TestMockGitClient) RestoreWorkingDir(_ ...string) error           { return nil }
func (m *TestMockGitClient) RestoreStaged(_ ...string) error               { return nil }
func (m *TestMockGitClient) RestoreFromCommit(_ string, _ ...string) error { return nil }
func (m *TestMockGitClient) RestoreAll() error                             { return nil }
func (m *TestMockGitClient) RestoreAllStaged() error                       { return nil }

// Config Operations
func (m *TestMockGitClient) ConfigGet(_ string) (string, error)       { return "", nil }
func (m *TestMockGitClient) ConfigSet(_, _ string) error              { return nil }
func (m *TestMockGitClient) ConfigGetGlobal(_ string) (string, error) { return "", nil }
func (m *TestMockGitClient) ConfigSetGlobal(_, _ string) error        { return nil }

// Reset Operations
func (m *TestMockGitClient) ResetHardAndClean() error { return nil }
func (m *TestMockGitClient) ResetHard(_ string) error { return nil }

// Clean Operations
func (m *TestMockGitClient) CleanFiles() error                { return nil }
func (m *TestMockGitClient) CleanDirs() error                 { return nil }
func (m *TestMockGitClient) CleanDryRun() (string, error)     { return "", nil }
func (m *TestMockGitClient) CleanFilesForce(_ []string) error { return nil }

// Utility Operations
func (m *TestMockGitClient) ListFiles() (string, error) { return "", nil }
func (m *TestMockGitClient) GetUpstreamBranchName(_ string) (string, error) {
	return "origin/main", nil
}
func (m *TestMockGitClient) GetAheadBehindCount(_, _ string) (string, error) {
	return m.aheadBehind, nil
}
func (m *TestMockGitClient) GetVersion() (string, error)    { return "test-version", nil }
func (m *TestMockGitClient) GetCommitHash() (string, error) { return "test-commit", nil }
