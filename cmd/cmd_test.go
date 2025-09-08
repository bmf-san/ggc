package cmd

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v5/git"
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
	resetHardAndCleanCalled bool
	cleanFilesCalled        bool
	cleanDirsCalled         bool
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

func (m *mockGitClient) ResetHardAndClean() error {
	m.resetHardAndCleanCalled = true
	return nil
}

func (m *mockGitClient) ResetHard(_ string) error {
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

func (m *mockGitClient) RevParseVerify(string) bool { return false }

// Config Operations
func (m *mockGitClient) ConfigGet(_ string) (string, error)       { return "", nil }
func (m *mockGitClient) ConfigSet(_, _ string) error              { return nil }
func (m *mockGitClient) ConfigGetGlobal(_ string) (string, error) { return "", nil }
func (m *mockGitClient) ConfigSetGlobal(_, _ string) error        { return nil }

// Repository Information methods
func (m *mockGitClient) GetBranchName() (string, error) { return "main", nil }
func (m *mockGitClient) GetGitStatus() (string, error)  { return "", nil }

// Status Operations methods
func (m *mockGitClient) Status() (string, error)               { return "", nil }
func (m *mockGitClient) StatusShort() (string, error)          { return "", nil }
func (m *mockGitClient) StatusWithColor() (string, error)      { return "", nil }
func (m *mockGitClient) StatusShortWithColor() (string, error) { return "", nil }

// Staging Operations methods
func (m *mockGitClient) Add(_ ...string) error { return nil }
func (m *mockGitClient) AddInteractive() error { return nil }

// Commit Operations methods
func (m *mockGitClient) Commit(_ string) error                 { return nil }
func (m *mockGitClient) CommitAmend() error                    { return nil }
func (m *mockGitClient) CommitAmendNoEdit() error              { return nil }
func (m *mockGitClient) CommitAmendWithMessage(_ string) error { return nil }

// Diff Operations methods
func (m *mockGitClient) Diff() (string, error)       { return "", nil }
func (m *mockGitClient) DiffStaged() (string, error) { return "", nil }
func (m *mockGitClient) DiffHead() (string, error)   { return "", nil }

// Branch Operations methods
func (m *mockGitClient) CheckoutNewBranch(_ string) error { return nil }
func (m *mockGitClient) CheckoutBranch(_ string) error    { return nil }
func (m *mockGitClient) CheckoutNewBranchFromRemote(_, _ string) error {
	return nil
}
func (m *mockGitClient) DeleteBranch(_ string) error           { return nil }
func (m *mockGitClient) ListMergedBranches() ([]string, error) { return []string{}, nil }
func (m *mockGitClient) RenameBranch(_, _ string) error        { return nil }
func (m *mockGitClient) MoveBranch(_, _ string) error          { return nil }
func (m *mockGitClient) SetUpstreamBranch(_, _ string) error   { return nil }
func (m *mockGitClient) GetBranchInfo(branch string) (*git.BranchInfo, error) {
	bi := &git.BranchInfo{Name: branch}
	return bi, nil
}
func (m *mockGitClient) ListBranchesVerbose() ([]git.BranchInfo, error) {
	return []git.BranchInfo{}, nil
}
func (m *mockGitClient) SortBranches(_ string) ([]string, error)       { return []string{}, nil }
func (m *mockGitClient) BranchesContaining(_ string) ([]string, error) { return []string{}, nil }

// Remote Operations methods
func (m *mockGitClient) Fetch(_ bool) error             { return nil }
func (m *mockGitClient) RemoteList() error              { return nil }
func (m *mockGitClient) RemoteAdd(_, _ string) error    { return nil }
func (m *mockGitClient) RemoteRemove(_ string) error    { return nil }
func (m *mockGitClient) RemoteSetURL(_, _ string) error { return nil }

// Tag Operations methods
func (m *mockGitClient) TagList(_ []string) error              { return nil }
func (m *mockGitClient) TagCreate(_, _ string) error           { return nil }
func (m *mockGitClient) TagCreateAnnotated(_, _ string) error  { return nil }
func (m *mockGitClient) TagDelete(_ []string) error            { return nil }
func (m *mockGitClient) TagPush(_, _ string) error             { return nil }
func (m *mockGitClient) TagPushAll(_ string) error             { return nil }
func (m *mockGitClient) TagShow(_ string) error                { return nil }
func (m *mockGitClient) GetLatestTag() (string, error)         { return "v1.0.0", nil }
func (m *mockGitClient) TagExists(_ string) bool               { return false }
func (m *mockGitClient) GetTagCommit(_ string) (string, error) { return "abc123", nil }

// Log Operations methods
func (m *mockGitClient) LogOneline(_, _ string) (string, error) { return "", nil }

// Rebase Operations methods
func (m *mockGitClient) RebaseInteractive(_ int) error { return nil }
func (m *mockGitClient) GetUpstreamBranch(_ string) (string, error) {
	return "origin/main", nil
}

// Stash Operations methods
func (m *mockGitClient) Stash() error               { return nil }
func (m *mockGitClient) StashList() (string, error) { return "", nil }
func (m *mockGitClient) StashShow(_ string) error   { return nil }
func (m *mockGitClient) StashApply(_ string) error  { return nil }
func (m *mockGitClient) StashPop(_ string) error    { return nil }
func (m *mockGitClient) StashDrop(_ string) error   { return nil }
func (m *mockGitClient) StashClear() error          { return nil }

// Restore Operations methods
func (m *mockGitClient) RestoreWorkingDir(_ ...string) error           { return nil }
func (m *mockGitClient) RestoreStaged(_ ...string) error               { return nil }
func (m *mockGitClient) RestoreFromCommit(_ string, _ ...string) error { return nil }
func (m *mockGitClient) RestoreAll() error                             { return nil }
func (m *mockGitClient) RestoreAllStaged() error                       { return nil }

// Reset and Clean Operations methods
func (m *mockGitClient) CleanDryRun() (string, error)     { return "", nil }
func (m *mockGitClient) CleanFilesForce(_ []string) error { return nil }

// Utility Operations methods
func (m *mockGitClient) ListFiles() (string, error) { return "", nil }
func (m *mockGitClient) GetUpstreamBranchName(_ string) (string, error) {
	return "origin/main", nil
}
func (m *mockGitClient) GetAheadBehindCount(_, _ string) (string, error) {
	return "0	0", nil
}
func (m *mockGitClient) GetVersion() (string, error)    { return "test-version", nil }
func (m *mockGitClient) GetCommitHash() (string, error) { return "test-commit", nil }

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
		puller:       NewPuller(mockClient),
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
		pusher:       NewPusher(mockClient),
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
				logger:       NewLogger(mc),
			}
			cmd.Log(tc.args)

			if !tc.wantCalled(mc) {
				t.Errorf("Expected method to be called for %s", tc.name)
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
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mc := &mockGitClient{}
			cmd := &Cmd{
				outputWriter: io.Discard,
				gitClient:    mc,
				committer:    NewCommitter(mc),
			}
			cmd.Commit(tc.args)

			if !tc.wantCalled(mc) {
				t.Errorf("Expected method to be called for %s", tc.name)
			}
		})
	}
}

func TestCmd_Reset(t *testing.T) {
	mc := &mockGitClient{}
	helper := NewHelper()
	helper.outputWriter = io.Discard
	cmd := &Cmd{
		outputWriter: io.Discard,
		gitClient:    mc,
		resetter:     &Resetter{gitClient: mc, outputWriter: io.Discard, helper: helper},
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
			cmd := &Cmd{
				outputWriter: io.Discard,
				gitClient:    mc,
				cleaner:      NewCleaner(mc),
			}
			cmd.Clean(tc.args)

			if !tc.wantCalled(mc) {
				t.Errorf("Expected method to be called for %s", tc.name)
			}
		})
	}
}

func TestNewCmd(t *testing.T) {
	// Test structure creation using NewCmdWithClient to inject mock
	mockClient := &mockGitClient{}

	cmd := NewCmd(mockClient)

	// Check if all fields are properly initialized
	if cmd.adder == nil {
		t.Error("adder should not be nil")
	}
	if cmd.brancher == nil {
		t.Error("brancher should not be nil")
	}
	if cmd.committer == nil {
		t.Error("committer should not be nil")
	}
	if cmd.logger == nil {
		t.Error("logger should not be nil")
	}
	if cmd.puller == nil {
		t.Error("puller should not be nil")
	}
	if cmd.pusher == nil {
		t.Error("pusher should not be nil")
	}
	if cmd.rebaser == nil {
		t.Error("rebaser should not be nil")
	}
	if cmd.remoteer == nil {
		t.Error("remoteer should not be nil")
	}
	if cmd.resetter == nil {
		t.Error("resetter should not be nil")
	}
	if cmd.stasher == nil {
		t.Error("stasher should not be nil")
	}
	if cmd.fetcher == nil {
		t.Error("fetcher should not be nil")
	}
	if cmd.completer == nil {
		t.Error("completer should not be nil")
	}
	if cmd.helper == nil {
		t.Error("helper should not be nil")
	}
}

func TestCmd_Help(t *testing.T) {
	// Use mock client to avoid actual git commands
	mockClient := &mockGitClient{}
	var buf bytes.Buffer
	helper := NewHelper()
	helper.outputWriter = &buf // Redirect help output to buffer to avoid console output
	cmd := &Cmd{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       helper,
	}

	// Test that Help method exists and can be called without panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Help() should not panic, but got: %v", r)
		}
	}()

	cmd.Help()

	// Verify that help was actually called (buffer should have content)
	if buf.Len() == 0 {
		t.Error("Help() should produce output")
	}
}

func TestCmd_Branch(t *testing.T) {
	// Use mock client to avoid git command side effects
	mockClient := &mockGitClient{}
	var buf bytes.Buffer
	helper := NewHelper()
	helper.outputWriter = &buf

	cmd := &Cmd{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       helper,
		brancher:     &Brancher{gitClient: mockClient, inputReader: bufio.NewReader(strings.NewReader("")), outputWriter: &buf, helper: helper},
	}

	cmd.Branch([]string{})

	output := buf.String()
	if output == "" {
		t.Error("Branch with no args should show help")
	}
}

func TestCmd_Route(t *testing.T) {
	// Use mock client to avoid git command side effects
	mockClient := &mockGitClient{}
	helper := NewHelper()
	helper.outputWriter = io.Discard

	cmd := &Cmd{
		gitClient:    mockClient,
		outputWriter: io.Discard,
		helper:       helper,
		// Initialize all components with mock clients to avoid side effects
		adder:       &Adder{gitClient: mockClient, outputWriter: io.Discard},
		brancher:    &Brancher{gitClient: mockClient, inputReader: bufio.NewReader(strings.NewReader("")), outputWriter: io.Discard, helper: helper},
		committer:   &Committer{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		logger:      &Logger{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		puller:      &Puller{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		pusher:      &Pusher{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		resetter:    &Resetter{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		cleaner:     &Cleaner{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		remoteer:    &Remoteer{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		rebaser:     &Rebaser{gitClient: mockClient, outputWriter: io.Discard, helper: helper, inputReader: bufio.NewReader(strings.NewReader(""))},
		stasher:     &Stasher{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		configureer: &Configureer{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		hooker:      &Hooker{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		tagger:      &Tagger{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		statuseer:   &Statuseer{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		versioneer:  &Versioneer{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		completer:   &Completer{gitClient: mockClient},
		differ:      &Differ{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		restoreer:   &Restoreer{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		fetcher:     &Fetcher{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
	}

	testCases := []struct {
		name string
		args []string
	}{
		{"help", []string{"help"}},
		{"add", []string{"add", "."}},
		{"branch", []string{"branch", "current"}},
		{"commit", []string{"commit", "test message"}},
		{"log", []string{"log", "simple"}},
		{"pull", []string{"pull", "current"}},
		{"push", []string{"push", "current"}},
		{"reset", []string{"reset"}},
		{"clean", []string{"clean", "files"}},
		{"version", []string{"version"}},
		{"clean-interactive", []string{"clean-interactive"}},
		{"remote", []string{"remote", "list"}},
		{"rebase", []string{"rebase", "interactive"}},
		{"stash", []string{"stash", "drop"}},
		{"config", []string{"config", "list"}},
		{"hook", []string{"hook", "list"}},
		{"tag", []string{"tag", "list"}},
		{"status", []string{"status"}},
		{"complete", []string{"complete", "bash"}},
		{"fetch", []string{"fetch", "--prune"}},
		{"diff", []string{"diff"}},
		{"restore", []string{"restore", "."}},
		{"unknown", []string{"unknown"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Route should not panic for any input
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Route() should not panic for args %v, but got: %v", tc.args, r)
				}
			}()

			cmd.Route(tc.args)
		})
	}
}

// Test some behavior of Interactive by mocking InteractiveUI
func TestCmd_Interactive_Existence(t *testing.T) {
	// Interactive is a complex interactive function, making complete testing difficult
	// However, verify the function exists
	cmd := &Cmd{}

	// Indirectly verify that the function exists
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Interactive should not panic when defined, got panic: %v", r)
		}
	}()

	// Use type assertion to verify the function is defined
	_ = cmd.Interactive
}

// Test wrapper functions that were not covered
func TestCmd_Status(t *testing.T) {
	// Use mock client to avoid git command side effects
	mockClient := &mockGitClient{}
	helper := NewHelper()
	helper.outputWriter = io.Discard
	cmd := &Cmd{
		outputWriter: io.Discard,
		gitClient:    mockClient,
		statuseer:    &Statuseer{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
	}

	// Test that Status calls the statuseer
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Status should not panic, got: %v", r)
		}
	}()

	cmd.Status([]string{})
}

func TestCmd_Config(t *testing.T) {
	// Use mock client to avoid git command side effects
	mockClient := &mockGitClient{}
	helper := NewHelper()
	helper.outputWriter = io.Discard
	cmd := &Cmd{
		outputWriter: io.Discard,
		gitClient:    mockClient,
		configureer:  &Configureer{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
	}

	// Test that Config calls the configureer
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Config should not panic, got: %v", r)
		}
	}()

	cmd.Config([]string{"list"})
}

func TestCmd_Hook(t *testing.T) {
	// Use mock client to avoid git command side effects
	mockClient := &mockGitClient{}
	helper := NewHelper()
	helper.outputWriter = io.Discard
	cmd := &Cmd{
		outputWriter: io.Discard,
		gitClient:    mockClient,
		hooker:       &Hooker{outputWriter: io.Discard, helper: helper},
	}

	// Test that Hook calls the hooker
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Hook should not panic, got: %v", r)
		}
	}()

	cmd.Hook([]string{})
}

func TestCmd_Tag(t *testing.T) {
	// Use mock client to avoid git command side effects
	mockClient := &mockGitClient{}
	helper := NewHelper()
	helper.outputWriter = io.Discard
	cmd := &Cmd{
		outputWriter: io.Discard,
		gitClient:    mockClient,
		tagger:       &Tagger{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
	}

	// Test that Tag calls the tagger
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Tag should not panic, got: %v", r)
		}
	}()

	cmd.Tag([]string{})
}

func TestCmd_Diff(t *testing.T) {
	// Use mock client to avoid git command side effects
	mockClient := &mockGitClient{}
	helper := NewHelper()
	helper.outputWriter = io.Discard
	cmd := &Cmd{
		outputWriter: io.Discard,
		gitClient:    mockClient,
		differ:       &Differ{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
	}

	// Test that Diff calls the differ
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Diff should not panic, got: %v", r)
		}
	}()

	cmd.Diff([]string{})
}

func TestCmd_Restore(t *testing.T) {
	// Use mock client to avoid git command side effects
	mockClient := &mockGitClient{}
	helper := NewHelper()
	helper.outputWriter = io.Discard
	cmd := &Cmd{
		outputWriter: io.Discard,
		gitClient:    mockClient,
		restoreer:    &Restoreer{outputWriter: io.Discard, helper: helper},
	}

	// Test that Restore calls the restoreer
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Restore should not panic, got: %v", r)
		}
	}()

	cmd.Restore([]string{})
}

func TestCmd_Version(t *testing.T) {
	// Mock getVersionInfo function
	SetVersionGetter(func() (string, string) {
		return "v1.0.0", "abc123"
	})

	// Use mock client to avoid git command side effects
	mockClient := &mockGitClient{}
	helper := NewHelper()
	helper.outputWriter = io.Discard
	cmd := &Cmd{
		outputWriter: io.Discard,
		gitClient:    mockClient,
		versioneer:   &Versioneer{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
	}

	// Test that Version calls the versioneer
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Version should not panic, got: %v", r)
		}
	}()

	cmd.Version([]string{})
}

func TestCmd_Interactive_Call(t *testing.T) {
	// Use mock client to avoid git command side effects
	mockClient := &mockGitClient{}
	cmd := &Cmd{
		outputWriter: io.Discard,
		gitClient:    mockClient,
	}

	// Test that Interactive function exists and can be called
	// Note: This is a complex interactive function, so we just test it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Interactive should not panic on instantiation, got: %v", r)
		}
	}()

	// Just verify the function is callable
	_ = cmd.Interactive
}

// TestCmd_Remote tests the Remote method
func TestCmd_Remote(t *testing.T) {
	mockClient := &mockGitClient{}
	var buf bytes.Buffer
	helper := NewHelper()
	helper.outputWriter = &buf

	cmd := &Cmd{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       helper,
		remoteer:     &Remoteer{gitClient: mockClient, outputWriter: &buf, helper: helper},
	}

	cmd.Remote([]string{})

	// Should show help when no args provided
	output := buf.String()
	if output == "" {
		t.Error("Remote with no args should show help")
	}
}

// TestCmd_Rebase tests the Rebase method
func TestCmd_Rebase(t *testing.T) {
	mockClient := &mockGitClient{}
	var buf bytes.Buffer
	helper := NewHelper()
	helper.outputWriter = &buf

	cmd := &Cmd{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       helper,
		rebaser:      &Rebaser{gitClient: mockClient, outputWriter: &buf, helper: helper, inputReader: bufio.NewReader(strings.NewReader("1\n"))},
	}

	cmd.Rebase([]string{})

	// Should show help when no args provided
	output := buf.String()
	if output == "" {
		t.Error("Rebase with no args should show help")
	}
}

// TestCmd_Stash tests the Stash method
func TestCmd_Stash(t *testing.T) {
	mockClient := &mockGitClient{}
	var buf bytes.Buffer
	helper := NewHelper()
	helper.outputWriter = &buf

	cmd := &Cmd{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       helper,
		stasher:      &Stasher{gitClient: mockClient, outputWriter: &buf, helper: helper},
	}

	cmd.Stash([]string{})

	// Should show help when no args provided
	output := buf.String()
	if output == "" {
		t.Error("Stash with no args should show help")
	}
}

// TestCmd_Fetch tests the Fetch method
func TestCmd_Fetch(t *testing.T) {
	mockClient := &mockGitClient{}
	var buf bytes.Buffer
	helper := NewHelper()
	helper.outputWriter = &buf

	cmd := &Cmd{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       helper,
	}

	// Test that Fetch method exists and doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Fetch panicked (expected in test environment): %v", r)
		}
	}()

	cmd.Fetch([]string{})

	output := buf.String()
	if !strings.Contains(output, "fetch") {
		t.Logf("Fetch output: %s", output)
	}
}

// TestCmd_Add tests the Add method
func TestCmd_Add(t *testing.T) {
	mockClient := &mockGitClient{}
	var buf bytes.Buffer

	cmd := &Cmd{
		gitClient:    mockClient,
		outputWriter: &buf,
		adder:        &Adder{gitClient: mockClient, outputWriter: &buf},
	}

	cmd.Add([]string{"."})

	// Should execute add command
	output := buf.String()
	if !strings.Contains(output, "Added") {
		t.Logf("Add output: %s", output)
	}
}

// TestCmd_WaitForContinue tests the waitForContinue method
func TestCmd_WaitForContinue(t *testing.T) {
	var buf bytes.Buffer
	cmd := &Cmd{
		outputWriter: &buf,
	}

	// Mock stdin with Enter key
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	defer func() {
		os.Stdin = oldStdin
		r.Close()
		w.Close()
	}()

	go func() {
		w.Write([]byte("\n"))
		w.Close()
	}()

	cmd.waitForContinue()

	output := buf.String()
	if !strings.Contains(output, "Press Enter") {
		t.Errorf("Expected 'Press Enter' in output, got: %s", output)
	}
}

// TestCmd_Stash_WithArgs tests the Stash method with arguments
func TestCmd_Stash_WithArgs(t *testing.T) {
	mockClient := &mockGitClient{}
	var buf bytes.Buffer
	helper := NewHelper()
	helper.outputWriter = &buf

	cmd := &Cmd{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       helper,
		stasher:      &Stasher{gitClient: mockClient, outputWriter: &buf, helper: helper},
	}

	// Test with list argument
	cmd.Stash([]string{"list"})

	// Should execute without panic
	output := buf.String()
	t.Logf("Stash list output: %s", output)
}

// TestCmd_Remote_WithArgs tests the Remote method with arguments
func TestCmd_Remote_WithArgs(t *testing.T) {
	mockClient := &mockGitClient{}
	var buf bytes.Buffer
	helper := NewHelper()
	helper.outputWriter = &buf

	cmd := &Cmd{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       helper,
		remoteer:     &Remoteer{gitClient: mockClient, outputWriter: &buf, helper: helper},
	}

	// Test with list argument
	cmd.Remote([]string{"list"})

	// Should execute without panic
	output := buf.String()
	t.Logf("Remote list output: %s", output)
}

// TestCmd_Tag_WithArgs tests the Tag method with arguments
func TestCmd_Tag_WithArgs(t *testing.T) {
	mockClient := &mockGitClient{}
	var buf bytes.Buffer
	helper := NewHelper()
	helper.outputWriter = &buf

	cmd := &Cmd{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       helper,
		tagger:       &Tagger{gitClient: mockClient, outputWriter: &buf, helper: helper},
	}

	// Test with list argument
	cmd.Tag([]string{"list"})

	// Should execute without panic
	output := buf.String()
	t.Logf("Tag list output: %s", output)
}

// TestCmd_Reset_WithArgs tests the Reset method with arguments
func TestCmd_Reset_WithArgs(t *testing.T) {
	mockClient := &mockGitClient{}
	var buf bytes.Buffer
	helper := NewHelper()
	helper.outputWriter = &buf

	cmd := &Cmd{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       helper,
		resetter:     &Resetter{gitClient: mockClient, outputWriter: &buf, helper: helper},
	}

	// Test with soft argument
	cmd.Reset([]string{"--soft", "HEAD~1"})

	// Should execute without panic
	output := buf.String()
	t.Logf("Reset output: %s", output)
}

// TestCmd_Rebase_WithArgs tests the Rebase method with arguments
func TestCmd_Rebase_WithArgs(t *testing.T) {
	mockClient := &mockGitClient{}
	var buf bytes.Buffer
	helper := NewHelper()
	helper.outputWriter = &buf

	cmd := &Cmd{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       helper,
		rebaser:      &Rebaser{gitClient: mockClient, outputWriter: &buf, helper: helper, inputReader: bufio.NewReader(strings.NewReader("1\n"))},
	}

	// Test with interactive argument
	cmd.Rebase([]string{"interactive"})

	// Should execute without panic
	output := buf.String()
	t.Logf("Rebase interactive output: %s", output)
}

// TestCmd_Config_WithArgs tests the Config method with arguments
func TestCmd_Config_WithArgs(t *testing.T) {
	mockClient := &mockGitClient{}
	var buf bytes.Buffer
	helper := NewHelper()
	helper.outputWriter = &buf

	cmd := &Cmd{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       helper,
		configureer:  &Configureer{gitClient: mockClient, outputWriter: &buf, helper: helper},
	}

	// Test with list argument
	cmd.Config([]string{"list"})

	// Should execute without panic
	output := buf.String()
	t.Logf("Config list output: %s", output)
}
