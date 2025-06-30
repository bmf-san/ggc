package cmd

import (
	"bytes"
	"io"
	"testing"

	"github.com/bmf-san/ggc/git"
)

// mockGitClient is a mock of git.Client.
type mockGitClient struct {
	pushCalled              bool
	pushForce               bool
	pullCalled              bool
	pullRebase              bool
	fetchPruneCalled        bool
	logSimpleCalled         bool
	logGraphCalled          bool
	commitAllowEmptyCalled  bool
	commitTmpCalled         bool
	resetHardAndCleanCalled bool
	cleanFilesCalled        bool
	cleanDirsCalled         bool
	stashPullPopCalled      bool
	stashPullPopErr         error
}

func (m *mockGitClient) GetCurrentBranch() (string, error) {
	return "main", nil
}

func (m *mockGitClient) ListLocalBranches() ([]string, error) {
	return []string{"main", "feature/test"}, nil
}

func (m *mockGitClient) ListRemoteBranches() ([]string, error) {
	return []string{"origin/main", "origin/feature/test"}, nil
}

func (m *mockGitClient) Push(force bool) error {
	m.pushCalled = true
	m.pushForce = force
	return nil
}

func (m *mockGitClient) Pull(rebase bool) error {
	m.pullCalled = true
	m.pullRebase = rebase
	return nil
}

func (m *mockGitClient) FetchPrune() error {
	m.fetchPruneCalled = true
	return nil
}

func (m *mockGitClient) LogSimple() error {
	m.logSimpleCalled = true
	return nil
}

func (m *mockGitClient) LogGraph() error {
	m.logGraphCalled = true
	return nil
}

func (m *mockGitClient) CommitAllowEmpty() error {
	m.commitAllowEmptyCalled = true
	return nil
}

func (m *mockGitClient) CommitTmp() error {
	m.commitTmpCalled = true
	return nil
}

func (m *mockGitClient) ResetHardAndClean() error {
	m.resetHardAndCleanCalled = true
	return nil
}

func (m *mockGitClient) CleanFiles() error {
	m.cleanFilesCalled = true
	return nil
}

func (m *mockGitClient) CleanDirs() error {
	m.cleanDirsCalled = true
	return nil
}

func (m *mockGitClient) GetBranchName() (string, error) {
	return "main", nil
}

func (m *mockGitClient) GetGitStatus() (string, error) {
	return "", nil
}

func (m *mockGitClient) CheckoutNewBranch(_ string) error {
	return nil
}

func (m *mockGitClient) StashPullPop() error {
	m.stashPullPopCalled = true
	return m.stashPullPopErr
}

type mockCmdGitClient struct {
	git.Clienter
	pullCalled bool
	pullRebase bool
	pushCalled bool
	pushForce  bool
}

func (m *mockCmdGitClient) Pull(rebase bool) error {
	m.pullCalled = true
	m.pullRebase = rebase
	return nil
}

func (m *mockCmdGitClient) Push(force bool) error {
	m.pushCalled = true
	m.pushForce = force
	return nil
}

func TestCmd_Pull(t *testing.T) {
	mockClient := &mockCmdGitClient{}
	var buf bytes.Buffer
	cmd := &Cmd{
		gitClient:    mockClient,
		outputWriter: &buf,
		puller:       NewPullerWithClient(mockClient),
	}

	cmd.Pull([]string{"current"})

	if !mockClient.pullCalled {
		t.Error("gitClient.Pull should be called")
	}
	if mockClient.pullRebase {
		t.Error("gitClient.Pull should be called with rebase=false")
	}
}

func TestCmd_Push(t *testing.T) {
	mockClient := &mockCmdGitClient{}
	var buf bytes.Buffer
	cmd := &Cmd{
		gitClient:    mockClient,
		outputWriter: &buf,
		pusher:       NewPusherWithClient(mockClient),
	}

	cmd.Push([]string{"current"})

	if !mockClient.pushCalled {
		t.Error("gitClient.Push should be called")
	}
	if mockClient.pushForce {
		t.Error("gitClient.Push should be called with force=false")
	}
}

func TestCmd_Log(t *testing.T) {
	cases := []struct {
		name       string
		args       []string
		wantCalled func(mc *mockGitClient) bool
	}{
		{
			name: "simple",
			args: []string{"simple"},
			wantCalled: func(mc *mockGitClient) bool {
				return mc.logSimpleCalled
			},
		},
		{
			name: "graph",
			args: []string{"graph"},
			wantCalled: func(mc *mockGitClient) bool {
				return mc.logGraphCalled
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mc := &mockGitClient{}
			cmd := &Cmd{
				outputWriter: io.Discard,
				gitClient:    mc,
				logger:       NewLoggerWithClient(mc),
			}
			cmd.Log(tc.args)
			if !tc.wantCalled(mc) {
				t.Errorf("gitClient.%s should be called", tc.name)
			}
		})
	}
}

func TestCmd_Commit(t *testing.T) {
	cases := []struct {
		name       string
		args       []string
		wantCalled func(mc *mockGitClient) bool
	}{
		{
			name: "allow-empty",
			args: []string{"allow-empty"},
			wantCalled: func(mc *mockGitClient) bool {
				return mc.commitAllowEmptyCalled
			},
		},
		{
			name: "tmp",
			args: []string{"tmp"},
			wantCalled: func(mc *mockGitClient) bool {
				return mc.commitTmpCalled
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mc := &mockGitClient{}
			cmd := &Cmd{
				outputWriter: io.Discard,
				gitClient:    mc,
				committer:    NewCommitterWithClient(mc),
			}
			cmd.Commit(tc.args)
			if !tc.wantCalled(mc) {
				t.Errorf("gitClient.%s should be called", tc.name)
			}
		})
	}
}

func TestCmd_Reset(t *testing.T) {
	mc := &mockGitClient{}
	cmd := &Cmd{
		outputWriter: io.Discard,
		gitClient:    mc,
		resetter:     NewResetterWithClient(mc),
	}
	cmd.Reset([]string{})

	if !mc.resetHardAndCleanCalled {
		t.Error("gitClient.ResetHardAndClean should be called")
	}
}

func TestCmd_Clean(t *testing.T) {
	cases := []struct {
		name       string
		args       []string
		wantCalled func(mc *mockGitClient) bool
	}{
		{
			name: "files",
			args: []string{"files"},
			wantCalled: func(mc *mockGitClient) bool {
				return mc.cleanFilesCalled
			},
		},
		{
			name: "dirs",
			args: []string{"dirs"},
			wantCalled: func(mc *mockGitClient) bool {
				return mc.cleanDirsCalled
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mc := &mockGitClient{}
			cleaner := NewCleaner()
			cleaner.gitClient = mc
			cmd := &Cmd{
				outputWriter: io.Discard,
				gitClient:    mc,
				cleaner:      cleaner,
			}
			cmd.Clean(tc.args)
			if !tc.wantCalled(mc) {
				t.Errorf("gitClient.CleanFiles or gitClient.CleanDirs should be called")
			}
		})
	}
}

func TestCmd_PullRebasePush(t *testing.T) {
	mc := &mockGitClient{}
	cmd := &Cmd{
		outputWriter:     io.Discard,
		gitClient:        mc,
		pullRebasePusher: NewPullRebasePusherWithClient(mc),
	}
	cmd.PullRebasePush()

	if !mc.pullCalled {
		t.Error("gitClient.Pull should be called")
	}
	if !mc.pullRebase {
		t.Error("gitClient.Pull should be called with rebase=true")
	}
	if !mc.pushCalled {
		t.Error("gitClient.Push should be called")
	}
	if mc.pushForce {
		t.Error("gitClient.Push should be called with force=false")
	}
}
