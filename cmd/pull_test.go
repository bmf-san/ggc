package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/bmf-san/ggc/v4/git"
)

type mockPullGitClient struct {
	git.Clienter
	pullCalled bool
	pullRebase bool
	err        error
}

func (m *mockPullGitClient) Pull(rebase bool) error {
	m.pullCalled = true
	m.pullRebase = rebase
	return m.err
}

// Implement all required methods from git.Clienter interface
func (m *mockPullGitClient) GetCurrentBranch() (string, error)               { return "main", nil }
func (m *mockPullGitClient) GetBranchName() (string, error)                  { return "main", nil }
func (m *mockPullGitClient) GetGitStatus() (string, error)                   { return "", nil }
func (m *mockPullGitClient) Status() (string, error)                         { return "", nil }
func (m *mockPullGitClient) StatusShort() (string, error)                    { return "", nil }
func (m *mockPullGitClient) StatusWithColor() (string, error)                { return "", nil }
func (m *mockPullGitClient) StatusShortWithColor() (string, error)           { return "", nil }
func (m *mockPullGitClient) Add(_ ...string) error                           { return nil }
func (m *mockPullGitClient) AddInteractive() error                           { return nil }
func (m *mockPullGitClient) Commit(_ string) error                           { return nil }
func (m *mockPullGitClient) CommitAmend() error                              { return nil }
func (m *mockPullGitClient) CommitAmendNoEdit() error                        { return nil }
func (m *mockPullGitClient) CommitAmendWithMessage(_ string) error           { return nil }
func (m *mockPullGitClient) CommitAllowEmpty() error                         { return nil }
func (m *mockPullGitClient) Diff() (string, error)                           { return "", nil }
func (m *mockPullGitClient) DiffStaged() (string, error)                     { return "", nil }
func (m *mockPullGitClient) DiffHead() (string, error)                       { return "", nil }
func (m *mockPullGitClient) ListLocalBranches() ([]string, error)            { return nil, nil }
func (m *mockPullGitClient) ListRemoteBranches() ([]string, error)           { return nil, nil }
func (m *mockPullGitClient) CheckoutNewBranch(_ string) error                { return nil }
func (m *mockPullGitClient) CheckoutBranch(_ string) error                   { return nil }
func (m *mockPullGitClient) CheckoutNewBranchFromRemote(_, _ string) error   { return nil }
func (m *mockPullGitClient) DeleteBranch(_ string) error                     { return nil }
func (m *mockPullGitClient) ListMergedBranches() ([]string, error)           { return nil, nil }
func (m *mockPullGitClient) RevParseVerify(_ string) bool                    { return false }
func (m *mockPullGitClient) Push(_ bool) error                               { return nil }
func (m *mockPullGitClient) Fetch(_ bool) error                              { return nil }
func (m *mockPullGitClient) RemoteList() error                               { return nil }
func (m *mockPullGitClient) RemoteAdd(_, _ string) error                     { return nil }
func (m *mockPullGitClient) RemoteRemove(_ string) error                     { return nil }
func (m *mockPullGitClient) RemoteSetURL(_, _ string) error                  { return nil }
func (m *mockPullGitClient) TagList(_ []string) error                        { return nil }
func (m *mockPullGitClient) TagCreate(_, _ string) error                     { return nil }
func (m *mockPullGitClient) TagCreateAnnotated(_, _ string) error            { return nil }
func (m *mockPullGitClient) TagDelete(_ []string) error                      { return nil }
func (m *mockPullGitClient) TagPush(_, _ string) error                       { return nil }
func (m *mockPullGitClient) TagPushAll(_ string) error                       { return nil }
func (m *mockPullGitClient) TagShow(_ string) error                          { return nil }
func (m *mockPullGitClient) GetLatestTag() (string, error)                   { return "", nil }
func (m *mockPullGitClient) TagExists(_ string) bool                         { return false }
func (m *mockPullGitClient) GetTagCommit(_ string) (string, error)           { return "", nil }
func (m *mockPullGitClient) LogSimple() error                                { return nil }
func (m *mockPullGitClient) LogGraph() error                                 { return nil }
func (m *mockPullGitClient) LogOneline(_, _ string) (string, error)          { return "", nil }
func (m *mockPullGitClient) RebaseInteractive(_ int) error                   { return nil }
func (m *mockPullGitClient) GetUpstreamBranch(_ string) (string, error)      { return "", nil }
func (m *mockPullGitClient) Stash() error                                    { return nil }
func (m *mockPullGitClient) StashList() (string, error)                      { return "", nil }
func (m *mockPullGitClient) StashShow(_ string) error                        { return nil }
func (m *mockPullGitClient) StashApply(_ string) error                       { return nil }
func (m *mockPullGitClient) StashPop(_ string) error                         { return nil }
func (m *mockPullGitClient) StashDrop(_ string) error                        { return nil }
func (m *mockPullGitClient) StashClear() error                               { return nil }
func (m *mockPullGitClient) RestoreWorkingDir(_ ...string) error             { return nil }
func (m *mockPullGitClient) RestoreStaged(_ ...string) error                 { return nil }
func (m *mockPullGitClient) RestoreFromCommit(_ string, _ ...string) error   { return nil }
func (m *mockPullGitClient) RestoreAll() error                               { return nil }
func (m *mockPullGitClient) RestoreAllStaged() error                         { return nil }
func (m *mockPullGitClient) ConfigGet(_ string) (string, error)              { return "", nil }
func (m *mockPullGitClient) ConfigSet(_, _ string) error                     { return nil }
func (m *mockPullGitClient) ConfigGetGlobal(_ string) (string, error)        { return "", nil }
func (m *mockPullGitClient) ConfigSetGlobal(_, _ string) error               { return nil }
func (m *mockPullGitClient) ResetHardAndClean() error                        { return nil }
func (m *mockPullGitClient) ResetHard(_ string) error                        { return nil }
func (m *mockPullGitClient) CleanFiles() error                               { return nil }
func (m *mockPullGitClient) CleanDirs() error                                { return nil }
func (m *mockPullGitClient) CleanDryRun() (string, error)                    { return "", nil }
func (m *mockPullGitClient) CleanFilesForce(_ []string) error                { return nil }
func (m *mockPullGitClient) ListFiles() (string, error)                      { return "", nil }
func (m *mockPullGitClient) GetUpstreamBranchName(_ string) (string, error)  { return "", nil }
func (m *mockPullGitClient) GetAheadBehindCount(_, _ string) (string, error) { return "", nil }
func (m *mockPullGitClient) GetVersion() (string, error)                     { return "test-version", nil }
func (m *mockPullGitClient) GetCommitHash() (string, error)                  { return "test-commit", nil }

func TestPuller_Pull(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantPull bool
		wantErr  bool
		err      error
	}{
		{
			name:     "normal_pull",
			args:     []string{"current"},
			wantPull: true,
			wantErr:  false,
		},
		{
			name:     "pull_with_error",
			args:     []string{"current"},
			wantPull: true,
			wantErr:  true,
			err:      errors.New("pull failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockPullGitClient{err: tt.err}
			var buf bytes.Buffer
			puller := &Puller{
				gitClient:    mockClient,
				outputWriter: &buf,
				helper:       NewHelper(),
			}
			puller.helper.outputWriter = &buf
			puller.Pull(tt.args)

			if mockClient.pullCalled != tt.wantPull {
				t.Errorf("Pull called = %v, want %v", mockClient.pullCalled, tt.wantPull)
			}

			if tt.wantErr {
				output := buf.String()
				if output != "Error: pull failed\n" {
					t.Errorf("Output = %q, want %q", output, "Error: pull failed\n")
				}
			}
		})
	}
}

func TestPuller_Pull_Help(t *testing.T) {
	var buf bytes.Buffer
	puller := &Puller{
		gitClient:    &mockPullGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	puller.helper.outputWriter = &buf
	puller.Pull([]string{})

	output := buf.String()
	if output == "" || !bytes.Contains(buf.Bytes(), []byte("Usage")) {
		t.Errorf("Usage should be displayed, but got: %s", output)
	}
}

func TestPuller_Pull_Rebase(t *testing.T) {
	mockClient := &mockPullGitClient{}
	var buf bytes.Buffer
	puller := &Puller{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	puller.helper.outputWriter = &buf
	puller.Pull([]string{"rebase"})

	if !mockClient.pullCalled {
		t.Error("Pull should be called")
	}
	if !mockClient.pullRebase {
		t.Error("Pull should be called with rebase=true")
	}
}

func TestPuller_Pull_RebaseError(t *testing.T) {
	mockClient := &mockPullGitClient{err: errors.New("rebase failed")}
	var buf bytes.Buffer
	puller := &Puller{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	puller.helper.outputWriter = &buf
	puller.Pull([]string{"rebase"})

	if !mockClient.pullCalled {
		t.Error("Pull should be called")
	}
	if !mockClient.pullRebase {
		t.Error("Pull should be called with rebase=true")
	}

	output := buf.String()
	if output != "Error: rebase failed\n" {
		t.Errorf("Expected error message, got: %s", output)
	}
}

func TestPuller_Pull_UnknownCommand(t *testing.T) {
	var buf bytes.Buffer
	puller := &Puller{
		gitClient:    &mockPullGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	puller.helper.outputWriter = &buf
	puller.Pull([]string{"unknown"})

	output := buf.String()
	if output == "" || !bytes.Contains(buf.Bytes(), []byte("Usage")) {
		t.Errorf("Usage should be displayed for unknown command, but got: %s", output)
	}
}
