package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/bmf-san/ggc/v5/git"
)

type mockPushGitClient struct {
	git.Clienter
	pushCalled bool
	pushForce  bool
	err        error
}

func (m *mockPushGitClient) Push(force bool) error {
	m.pushCalled = true
	m.pushForce = force
	return m.err
}

// Implement all required methods from git.Clienter interface
func (m *mockPushGitClient) GetCurrentBranch() (string, error)               { return "main", nil }
func (m *mockPushGitClient) GetBranchName() (string, error)                  { return "main", nil }
func (m *mockPushGitClient) GetGitStatus() (string, error)                   { return "", nil }
func (m *mockPushGitClient) Status() (string, error)                         { return "", nil }
func (m *mockPushGitClient) StatusShort() (string, error)                    { return "", nil }
func (m *mockPushGitClient) StatusWithColor() (string, error)                { return "", nil }
func (m *mockPushGitClient) StatusShortWithColor() (string, error)           { return "", nil }
func (m *mockPushGitClient) Add(_ ...string) error                           { return nil }
func (m *mockPushGitClient) AddInteractive() error                           { return nil }
func (m *mockPushGitClient) Commit(_ string) error                           { return nil }
func (m *mockPushGitClient) CommitAmend() error                              { return nil }
func (m *mockPushGitClient) CommitAmendNoEdit() error                        { return nil }
func (m *mockPushGitClient) CommitAmendWithMessage(_ string) error           { return nil }
func (m *mockPushGitClient) CommitAllowEmpty() error                         { return nil }
func (m *mockPushGitClient) Diff() (string, error)                           { return "", nil }
func (m *mockPushGitClient) DiffStaged() (string, error)                     { return "", nil }
func (m *mockPushGitClient) DiffHead() (string, error)                       { return "", nil }
func (m *mockPushGitClient) ListLocalBranches() ([]string, error)            { return nil, nil }
func (m *mockPushGitClient) ListRemoteBranches() ([]string, error)           { return nil, nil }
func (m *mockPushGitClient) CheckoutNewBranch(_ string) error                { return nil }
func (m *mockPushGitClient) CheckoutBranch(_ string) error                   { return nil }
func (m *mockPushGitClient) CheckoutNewBranchFromRemote(_, _ string) error   { return nil }
func (m *mockPushGitClient) DeleteBranch(_ string) error                     { return nil }
func (m *mockPushGitClient) ListMergedBranches() ([]string, error)           { return nil, nil }
func (m *mockPushGitClient) RevParseVerify(_ string) bool                    { return false }
func (m *mockPushGitClient) Pull(_ bool) error                               { return nil }
func (m *mockPushGitClient) Fetch(_ bool) error                              { return nil }
func (m *mockPushGitClient) RemoteList() error                               { return nil }
func (m *mockPushGitClient) RemoteAdd(_, _ string) error                     { return nil }
func (m *mockPushGitClient) RemoteRemove(_ string) error                     { return nil }
func (m *mockPushGitClient) RemoteSetURL(_, _ string) error                  { return nil }
func (m *mockPushGitClient) TagList(_ []string) error                        { return nil }
func (m *mockPushGitClient) TagCreate(_, _ string) error                     { return nil }
func (m *mockPushGitClient) TagCreateAnnotated(_, _ string) error            { return nil }
func (m *mockPushGitClient) TagDelete(_ []string) error                      { return nil }
func (m *mockPushGitClient) TagPush(_, _ string) error                       { return nil }
func (m *mockPushGitClient) TagPushAll(_ string) error                       { return nil }
func (m *mockPushGitClient) TagShow(_ string) error                          { return nil }
func (m *mockPushGitClient) GetLatestTag() (string, error)                   { return "", nil }
func (m *mockPushGitClient) TagExists(_ string) bool                         { return false }
func (m *mockPushGitClient) GetTagCommit(_ string) (string, error)           { return "", nil }
func (m *mockPushGitClient) LogSimple() error                                { return nil }
func (m *mockPushGitClient) LogGraph() error                                 { return nil }
func (m *mockPushGitClient) LogOneline(_, _ string) (string, error)          { return "", nil }
func (m *mockPushGitClient) RebaseInteractive(_ int) error                   { return nil }
func (m *mockPushGitClient) GetUpstreamBranch(_ string) (string, error)      { return "", nil }
func (m *mockPushGitClient) Stash() error                                    { return nil }
func (m *mockPushGitClient) StashList() (string, error)                      { return "", nil }
func (m *mockPushGitClient) StashShow(_ string) error                        { return nil }
func (m *mockPushGitClient) StashApply(_ string) error                       { return nil }
func (m *mockPushGitClient) StashPop(_ string) error                         { return nil }
func (m *mockPushGitClient) StashDrop(_ string) error                        { return nil }
func (m *mockPushGitClient) StashClear() error                               { return nil }
func (m *mockPushGitClient) RestoreWorkingDir(_ ...string) error             { return nil }
func (m *mockPushGitClient) RestoreStaged(_ ...string) error                 { return nil }
func (m *mockPushGitClient) RestoreFromCommit(_ string, _ ...string) error   { return nil }
func (m *mockPushGitClient) RestoreAll() error                               { return nil }
func (m *mockPushGitClient) RestoreAllStaged() error                         { return nil }
func (m *mockPushGitClient) ConfigGet(_ string) (string, error)              { return "", nil }
func (m *mockPushGitClient) ConfigSet(_, _ string) error                     { return nil }
func (m *mockPushGitClient) ConfigGetGlobal(_ string) (string, error)        { return "", nil }
func (m *mockPushGitClient) ConfigSetGlobal(_, _ string) error               { return nil }
func (m *mockPushGitClient) ResetHardAndClean() error                        { return nil }
func (m *mockPushGitClient) ResetHard(_ string) error                        { return nil }
func (m *mockPushGitClient) CleanFiles() error                               { return nil }
func (m *mockPushGitClient) CleanDirs() error                                { return nil }
func (m *mockPushGitClient) CleanDryRun() (string, error)                    { return "", nil }
func (m *mockPushGitClient) CleanFilesForce(_ []string) error                { return nil }
func (m *mockPushGitClient) ListFiles() (string, error)                      { return "", nil }
func (m *mockPushGitClient) GetUpstreamBranchName(_ string) (string, error)  { return "", nil }
func (m *mockPushGitClient) GetAheadBehindCount(_, _ string) (string, error) { return "", nil }
func (m *mockPushGitClient) GetVersion() (string, error)                     { return "test-version", nil }
func (m *mockPushGitClient) GetCommitHash() (string, error)                  { return "test-commit", nil }

func TestPusher_Push(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantPush bool
		wantErr  bool
		err      error
	}{
		{
			name:     "normal_push",
			args:     []string{"current"},
			wantPush: true,
			wantErr:  false,
		},
		{
			name:     "push_with_error",
			args:     []string{"current"},
			wantPush: true,
			wantErr:  true,
			err:      errors.New("push failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockPushGitClient{err: tt.err}
			var buf bytes.Buffer
			pusher := &Pusher{
				gitClient:    mockClient,
				outputWriter: &buf,
				helper:       NewHelper(),
			}
			pusher.helper.outputWriter = &buf
			pusher.Push(tt.args)

			if mockClient.pushCalled != tt.wantPush {
				t.Errorf("Push called = %v, want %v", mockClient.pushCalled, tt.wantPush)
			}

			if tt.wantErr {
				output := buf.String()
				if output != "Error: push failed\n" {
					t.Errorf("Output = %q, want %q", output, "Error: push failed\n")
				}
			}
		})
	}
}

func TestPusher_Push_Help(t *testing.T) {
	var buf bytes.Buffer
	pusher := &Pusher{
		gitClient:    &mockPushGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	pusher.helper.outputWriter = &buf
	pusher.Push([]string{})

	output := buf.String()
	if output == "" || !bytes.Contains(buf.Bytes(), []byte("Usage")) {
		t.Errorf("Usage should be displayed, but got: %s", output)
	}
}

func TestPusher_Push_Force(t *testing.T) {
	mockClient := &mockPushGitClient{}
	var buf bytes.Buffer
	pusher := &Pusher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	pusher.helper.outputWriter = &buf
	pusher.Push([]string{"force"})

	if !mockClient.pushCalled {
		t.Error("Push should be called")
	}
	if !mockClient.pushForce {
		t.Error("Push should be called with force=true")
	}
}

func TestPusher_Push_ForceError(t *testing.T) {
	mockClient := &mockPushGitClient{err: errors.New("force push failed")}
	var buf bytes.Buffer
	pusher := &Pusher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	pusher.helper.outputWriter = &buf
	pusher.Push([]string{"force"})

	if !mockClient.pushCalled {
		t.Error("Push should be called")
	}
	if !mockClient.pushForce {
		t.Error("Push should be called with force=true")
	}

	output := buf.String()
	if output != "Error: force push failed\n" {
		t.Errorf("Expected error message, got: %s", output)
	}
}

func TestPusher_Push_UnknownCommand(t *testing.T) {
	var buf bytes.Buffer
	pusher := &Pusher{
		gitClient:    &mockPushGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	pusher.helper.outputWriter = &buf
	pusher.Push([]string{"unknown"})

	output := buf.String()
	if output == "" || !bytes.Contains(buf.Bytes(), []byte("Usage")) {
		t.Errorf("Usage should be displayed for unknown command, but got: %s", output)
	}
}
