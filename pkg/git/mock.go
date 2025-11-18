//go:build test

package git

// MockClient is a test-only composite mock for selected git operations.
type MockClient struct {
	GetCurrentBranchFunc   func() (string, error)
	ListLocalBranchesFunc  func() ([]string, error)
	ListRemoteBranchesFunc func() ([]string, error)
	CheckoutNewBranchFunc  func(name string) error
	RevParseVerifyFunc     func(ref string) bool
	CleanFilesFunc         func() error
	CleanDirsFunc          func() error
	CommitAllowEmptyFunc   func() error
	FetchPruneFunc         func() error
	LogSimpleFunc          func() error
	LogGraphFunc           func() error
	PullFunc               func(rebase bool) error
	PushFunc               func(force bool) error
	ResetHardAndCleanFunc  func() error
	DeleteBranchFunc       func(branch string) error
	DeleteMergedBranchFunc func() error
	GetGitStatusFunc       func() (string, error)
	GetBranchNameFunc      func() (string, error)
}

// GetCurrentBranch is a mock of GetCurrentBranch.
func (m *MockClient) GetCurrentBranch() (string, error) {
	return m.GetCurrentBranchFunc()
}

// ListLocalBranches is a mock of ListLocalBranches.
func (m *MockClient) ListLocalBranches() ([]string, error) {
	return m.ListLocalBranchesFunc()
}

// ListRemoteBranches is a mock of ListRemoteBranches.
func (m *MockClient) ListRemoteBranches() ([]string, error) {
	return m.ListRemoteBranchesFunc()
}

// CheckoutNewBranch is a mock of CheckoutNewBranch.
func (m *MockClient) CheckoutNewBranch(name string) error {
	return m.CheckoutNewBranchFunc(name)
}

// RevParseVerify is a mock of RevParseVerify.
func (m *MockClient) RevParseVerify(ref string) bool {
	if m.RevParseVerifyFunc != nil {
		return m.RevParseVerifyFunc(ref)
	}
	return false
}

// CleanFiles is a mock of CleanFiles.
func (m *MockClient) CleanFiles() error {
	return m.CleanFilesFunc()
}

// CleanDirs is a mock of CleanDirs.
func (m *MockClient) CleanDirs() error {
	return m.CleanDirsFunc()
}

// CommitAllowEmpty is a mock of CommitAllowEmpty.
func (m *MockClient) CommitAllowEmpty() error {
	return m.CommitAllowEmptyFunc()
}

// FetchPrune is a mock of FetchPrune.
func (m *MockClient) FetchPrune() error {
	return m.FetchPruneFunc()
}

// LogSimple is a mock of LogSimple.
func (m *MockClient) LogSimple() error {
	return m.LogSimpleFunc()
}

// LogGraph is a mock of LogGraph.
func (m *MockClient) LogGraph() error {
	if m.LogGraphFunc != nil {
		return m.LogGraphFunc()
	}
	return nil
}

// Pull is a mock of Pull.
func (m *MockClient) Pull(rebase bool) error {
	return m.PullFunc(rebase)
}

// Push is a mock of Push.
func (m *MockClient) Push(force bool) error {
	return m.PushFunc(force)
}

// ResetHardAndClean is a mock of ResetHardAndClean.
func (m *MockClient) ResetHardAndClean() error {
	return m.ResetHardAndCleanFunc()
}

// DeleteBranch is a mock of DeleteBranch.
func (m *MockClient) DeleteBranch(branch string) error {
	return m.DeleteBranchFunc(branch)
}

// DeleteMergedBranch is a mock of DeleteMergedBranch.
func (m *MockClient) DeleteMergedBranch() error {
	return m.DeleteMergedBranchFunc()
}

// GetGitStatus is a mock of GetGitStatus.
func (m *MockClient) GetGitStatus() (string, error) {
	return m.GetGitStatusFunc()
}

// GetBranchName is a mock of GetBranchName.
func (m *MockClient) GetBranchName() (string, error) {
	return m.GetBranchNameFunc()
}
