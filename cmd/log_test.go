package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v4/git"
)

type mockLogGitClient struct {
	git.Clienter
	logSimpleCalled bool
	logGraphCalled  bool
	err             error
}

func (m *mockLogGitClient) LogSimple() error {
	m.logSimpleCalled = true
	return m.err
}

func (m *mockLogGitClient) LogGraph() error {
	m.logGraphCalled = true
	return m.err
}

// Implement all required methods from git.Clienter interface
func (m *mockLogGitClient) GetCurrentBranch() (string, error)               { return "main", nil }
func (m *mockLogGitClient) GetBranchName() (string, error)                  { return "main", nil }
func (m *mockLogGitClient) GetGitStatus() (string, error)                   { return "", nil }
func (m *mockLogGitClient) Status() (string, error)                         { return "", nil }
func (m *mockLogGitClient) StatusShort() (string, error)                    { return "", nil }
func (m *mockLogGitClient) StatusWithColor() (string, error)                { return "", nil }
func (m *mockLogGitClient) StatusShortWithColor() (string, error)           { return "", nil }
func (m *mockLogGitClient) Add(_ ...string) error                           { return nil }
func (m *mockLogGitClient) AddInteractive() error                           { return nil }
func (m *mockLogGitClient) Commit(_ string) error                           { return nil }
func (m *mockLogGitClient) CommitAmend() error                              { return nil }
func (m *mockLogGitClient) CommitAmendNoEdit() error                        { return nil }
func (m *mockLogGitClient) CommitAmendWithMessage(_ string) error           { return nil }
func (m *mockLogGitClient) CommitAllowEmpty() error                         { return nil }
func (m *mockLogGitClient) Diff() (string, error)                           { return "", nil }
func (m *mockLogGitClient) DiffStaged() (string, error)                     { return "", nil }
func (m *mockLogGitClient) DiffHead() (string, error)                       { return "", nil }
func (m *mockLogGitClient) ListLocalBranches() ([]string, error)            { return nil, nil }
func (m *mockLogGitClient) ListRemoteBranches() ([]string, error)           { return nil, nil }
func (m *mockLogGitClient) CheckoutNewBranch(_ string) error                { return nil }
func (m *mockLogGitClient) CheckoutBranch(_ string) error                   { return nil }
func (m *mockLogGitClient) CheckoutNewBranchFromRemote(_, _ string) error   { return nil }
func (m *mockLogGitClient) DeleteBranch(_ string) error                     { return nil }
func (m *mockLogGitClient) ListMergedBranches() ([]string, error)           { return nil, nil }
func (m *mockLogGitClient) RevParseVerify(_ string) bool                    { return false }
func (m *mockLogGitClient) Push(_ bool) error                               { return nil }
func (m *mockLogGitClient) Pull(_ bool) error                               { return nil }
func (m *mockLogGitClient) Fetch(_ bool) error                              { return nil }
func (m *mockLogGitClient) RemoteList() error                               { return nil }
func (m *mockLogGitClient) RemoteAdd(_, _ string) error                     { return nil }
func (m *mockLogGitClient) RemoteRemove(_ string) error                     { return nil }
func (m *mockLogGitClient) RemoteSetURL(_, _ string) error                  { return nil }
func (m *mockLogGitClient) TagList(_ []string) error                        { return nil }
func (m *mockLogGitClient) TagCreate(_, _ string) error                     { return nil }
func (m *mockLogGitClient) TagCreateAnnotated(_, _ string) error            { return nil }
func (m *mockLogGitClient) TagDelete(_ []string) error                      { return nil }
func (m *mockLogGitClient) TagPush(_, _ string) error                       { return nil }
func (m *mockLogGitClient) TagPushAll(_ string) error                       { return nil }
func (m *mockLogGitClient) TagShow(_ string) error                          { return nil }
func (m *mockLogGitClient) GetLatestTag() (string, error)                   { return "", nil }
func (m *mockLogGitClient) TagExists(_ string) bool                         { return false }
func (m *mockLogGitClient) GetTagCommit(_ string) (string, error)           { return "", nil }
func (m *mockLogGitClient) LogOneline(_, _ string) (string, error)          { return "", nil }
func (m *mockLogGitClient) RebaseInteractive(_ int) error                   { return nil }
func (m *mockLogGitClient) GetUpstreamBranch(_ string) (string, error)      { return "", nil }
func (m *mockLogGitClient) Stash() error                                    { return nil }
func (m *mockLogGitClient) StashList() (string, error)                      { return "", nil }
func (m *mockLogGitClient) StashShow(_ string) error                        { return nil }
func (m *mockLogGitClient) StashApply(_ string) error                       { return nil }
func (m *mockLogGitClient) StashPop(_ string) error                         { return nil }
func (m *mockLogGitClient) StashDrop(_ string) error                        { return nil }
func (m *mockLogGitClient) StashClear() error                               { return nil }
func (m *mockLogGitClient) RestoreWorkingDir(_ ...string) error             { return nil }
func (m *mockLogGitClient) RestoreStaged(_ ...string) error                 { return nil }
func (m *mockLogGitClient) RestoreFromCommit(_ string, _ ...string) error   { return nil }
func (m *mockLogGitClient) RestoreAll() error                               { return nil }
func (m *mockLogGitClient) RestoreAllStaged() error                         { return nil }
func (m *mockLogGitClient) ConfigGet(_ string) (string, error)              { return "", nil }
func (m *mockLogGitClient) ConfigSet(_, _ string) error                     { return nil }
func (m *mockLogGitClient) ConfigGetGlobal(_ string) (string, error)        { return "", nil }
func (m *mockLogGitClient) ConfigSetGlobal(_, _ string) error               { return nil }
func (m *mockLogGitClient) ResetHardAndClean() error                        { return nil }
func (m *mockLogGitClient) ResetHard(_ string) error                        { return nil }
func (m *mockLogGitClient) CleanFiles() error                               { return nil }
func (m *mockLogGitClient) CleanDirs() error                                { return nil }
func (m *mockLogGitClient) CleanDryRun() (string, error)                    { return "", nil }
func (m *mockLogGitClient) CleanFilesForce(_ []string) error                { return nil }
func (m *mockLogGitClient) ListFiles() (string, error)                      { return "", nil }
func (m *mockLogGitClient) GetUpstreamBranchName(_ string) (string, error)  { return "", nil }
func (m *mockLogGitClient) GetAheadBehindCount(_, _ string) (string, error) { return "", nil }
func (m *mockLogGitClient) GetVersion() (string, error)                     { return "test-version", nil }
func (m *mockLogGitClient) GetCommitHash() (string, error)                  { return "test-commit", nil }

func TestLogger_Log_Simple(t *testing.T) {
	mockClient := &mockLogGitClient{}
	var buf bytes.Buffer
	l := &Logger{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	l.helper.outputWriter = &buf
	l.Log([]string{"simple"})
	if !mockClient.logSimpleCalled {
		t.Error("LogSimple should be called")
	}
}

func TestLogger_Log_Graph(t *testing.T) {
	mockClient := &mockLogGitClient{}
	var buf bytes.Buffer
	l := &Logger{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	l.helper.outputWriter = &buf
	l.Log([]string{"graph"})
	if !mockClient.logGraphCalled {
		t.Error("LogGraph should be called")
	}
}

func TestLogger_Log_Help(t *testing.T) {
	var buf bytes.Buffer
	l := &Logger{
		gitClient:    &mockLogGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	l.helper.outputWriter = &buf
	l.Log([]string{"unknown"})

	output := buf.String()
	if output == "" || !strings.Contains(output, "Usage") {
		t.Errorf("Usage should be displayed, but got: %s", output)
	}
}

func TestLogger_Log_Simple_Error(t *testing.T) {
	var buf bytes.Buffer
	l := &Logger{
		gitClient:    &mockLogGitClient{err: errors.New("fail")},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	l.helper.outputWriter = &buf
	l.Log([]string{"simple"})

	output := buf.String()
	if output != "Error: fail\n" {
		t.Errorf("unexpected output: got %q", output)
	}
}

func TestLogger_Log_Graph_Error(t *testing.T) {
	var buf bytes.Buffer
	l := &Logger{
		gitClient:    &mockLogGitClient{err: errors.New("fail")},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	l.helper.outputWriter = &buf
	l.Log([]string{"graph"})

	output := buf.String()
	if output != "Error: fail\n" {
		t.Errorf("unexpected output: got %q", output)
	}
}

func TestLogger_Log_NoArgs(t *testing.T) {
	var buf bytes.Buffer
	l := &Logger{
		gitClient:    &mockLogGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	l.helper.outputWriter = &buf
	l.Log([]string{})

	output := buf.String()
	if output == "" || !strings.Contains(output, "Usage") {
		t.Errorf("Usage should be displayed when no args provided, but got: %s", output)
	}
}

func TestLogger_Log_Simple_OutputFormat(t *testing.T) {
	tests := []struct {
		name        string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "Successful simple log",
			shouldError: false,
		},
		{
			name:        "Error during simple log",
			shouldError: true,
			errorMsg:    "git error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			mockClient := &mockLogGitClient{}

			if tt.shouldError {
				mockClient.err = errors.New(tt.errorMsg)
			}

			l := &Logger{
				gitClient:    mockClient,
				outputWriter: &buf,
				helper:       NewHelper(),
			}
			l.helper.outputWriter = &buf

			l.Log([]string{"simple"})

			if tt.shouldError {
				output := buf.String()
				if !strings.Contains(output, "Error:") {
					t.Errorf("Expected error message, got: %s", output)
				}
				if !strings.Contains(output, tt.errorMsg) {
					t.Errorf("Expected error message to contain %q, got: %s", tt.errorMsg, output)
				}
			} else if !mockClient.logSimpleCalled {
				t.Error("LogSimple should be called")
			}
		})
	}
}

func TestLogger_Log_Graph_OutputFormat(t *testing.T) {
	tests := []struct {
		name        string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "Successful graph log",
			shouldError: false,
		},
		{
			name:        "Error during graph log",
			shouldError: true,
			errorMsg:    "git error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			mockClient := &mockLogGitClient{}

			if tt.shouldError {
				mockClient.err = errors.New(tt.errorMsg)
			}

			l := &Logger{
				gitClient:    mockClient,
				outputWriter: &buf,
				helper:       NewHelper(),
			}
			l.helper.outputWriter = &buf

			l.Log([]string{"graph"})

			if tt.shouldError {
				output := buf.String()
				if !strings.Contains(output, "Error:") {
					t.Errorf("Expected error message, got: %s", output)
				}
				if !strings.Contains(output, tt.errorMsg) {
					t.Errorf("Expected error message to contain %q, got: %s", tt.errorMsg, output)
				}
			} else if !mockClient.logGraphCalled {
				t.Error("LogGraph should be called")
			}
		})
	}
}
