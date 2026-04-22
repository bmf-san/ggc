package cmd

import (
	"bytes"
	"io"
	"strings"
	"testing"

	commandregistry "github.com/bmf-san/ggc/v8/cmd/command"
	"github.com/bmf-san/ggc/v8/internal/config"
	"github.com/bmf-san/ggc/v8/internal/interactive"
	"github.com/bmf-san/ggc/v8/internal/prompt"
	"github.com/bmf-san/ggc/v8/internal/testutil"
)

// mockGitClient embeds testutil.MockGitClient for no-op defaults and only
// overrides the operations whose invocation needs to be observed.
type mockGitClient struct {
	testutil.MockGitClient
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

func (m *mockGitClient) CleanFiles() error {
	m.cleanFilesCalled = true
	return nil
}

func (m *mockGitClient) CleanDirs() error {
	m.cleanDirsCalled = true
	return nil
}

type mockCmdGitClient struct {
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
		gitClient:    nil,
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
		gitClient:    nil,
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
			name: "allow empty",
			args: []string{"allow", "empty"},
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
	cm := config.NewConfigManager(mockClient)

	cmd, _ := NewCmd(mockClient, cm)

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
	if cmd.remoter == nil {
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

	cmd.Help(nil)

	// Verify that help was actually called (buffer should have content)
	if buf.Len() == 0 {
		t.Error("Help() should produce output")
	}

	buf.Reset()
	cmd.Help([]string{"branch"})
	if buf.Len() == 0 {
		t.Error("Help(command) should produce output")
	}
	if !strings.Contains(buf.String(), "ggc branch") {
		t.Error("Help(command) should render branch usage")
	}
}

func TestCmd_Branch(t *testing.T) {
	// Use mock client to avoid git command side effects
	var buf bytes.Buffer
	helper := NewHelper()
	helper.outputWriter = &buf

	cmd := &Cmd{
		gitClient:    nil,
		outputWriter: &buf,
		helper:       helper,
		// help path does not access gitClient; pass nil to minimize dependencies
		brancher: &Brancher{gitClient: nil, prompter: prompt.New(strings.NewReader(""), &buf), outputWriter: &buf, helper: helper},
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
		registry:     commandregistry.NewRegistry(),
		// Initialize all components with mock clients to avoid side effects
		adder:      &Adder{gitClient: mockClient, outputWriter: io.Discard},
		brancher:   &Brancher{gitClient: mockClient, prompter: prompt.New(strings.NewReader(""), io.Discard), outputWriter: io.Discard, helper: helper},
		committer:  &Committer{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		logger:     &Logger{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		puller:     &Puller{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		pusher:     &Pusher{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		resetter:   &Resetter{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		cleaner:    &Cleaner{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		remoter:    &Remoter{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		rebaser:    &Rebaser{gitClient: mockClient, outputWriter: io.Discard, helper: helper, prompter: prompt.New(strings.NewReader(""), io.Discard)},
		stasher:    &Stasher{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		configurer: &Configurer{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		hooker:     &Hooker{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		tagger:     &Tagger{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		statuser:   &Statuser{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		versioner:  &Versioner{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		differ:     &Differ{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		restorer:   &Restorer{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		fetcher:    &Fetcher{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
	}
	cmd.cmdRouter, _ = newCommandRouter(cmd)

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

			_ = cmd.Route(tc.args)
		})
	}
}

func TestCmd_Route_SeparatorAllowsHyphenValues(t *testing.T) {
	// Use mock client to avoid git command side effects
	mockClient := &mockGitClient{}
	helper := NewHelper()
	helper.outputWriter = io.Discard

	var buf bytes.Buffer
	cmd := &Cmd{
		gitClient:    mockClient,
		outputWriter: &buf, // capture legacy-like error output if any
		helper:       helper,
		registry:     commandregistry.NewRegistry(),
		adder:        &Adder{gitClient: mockClient, outputWriter: io.Discard},
		brancher:     &Brancher{gitClient: mockClient, prompter: prompt.New(strings.NewReader(""), io.Discard), outputWriter: io.Discard, helper: helper},
		committer:    &Committer{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		logger:       &Logger{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		puller:       &Puller{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		pusher:       &Pusher{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		resetter:     &Resetter{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		cleaner:      &Cleaner{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		remoter:      &Remoter{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		rebaser:      &Rebaser{gitClient: mockClient, outputWriter: io.Discard, helper: helper, prompter: prompt.New(strings.NewReader(""), io.Discard)},
		stasher:      &Stasher{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		configurer:   &Configurer{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		hooker:       &Hooker{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		tagger:       &Tagger{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		statuser:     &Statuser{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		versioner:    &Versioner{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		differ:       &Differ{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		restorer:     &Restorer{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
		fetcher:      &Fetcher{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
	}
	cmd.cmdRouter, _ = newCommandRouter(cmd)

	// Using "--" should allow a value starting with '-' to pass through
	// without triggering the legacy-like error.
	args := []string{"commit", "--", "-leading dash message"}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Route() should not panic for args %v, but got: %v", args, r)
		}
	}()

	_ = cmd.Route(args)

	if strings.Contains(buf.String(), "legacy-like syntax is not supported") {
		t.Fatalf("did not expect legacy-like error when using '--' separator, got: %q", buf.String())
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
		statuser:     &Statuser{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
	}

	// Test that Status calls the statuser
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
		configurer:   &Configurer{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
	}

	// Test that Config calls the configurer
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
		restorer:     &Restorer{outputWriter: io.Discard, helper: helper},
	}

	// Test that Restore calls the restorer
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
		versioner:    &Versioner{gitClient: mockClient, outputWriter: io.Discard, helper: helper},
	}

	// Test that Version calls the versioner
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

func TestCmd_waitForContinue(t *testing.T) {
	// Test that waitForContinue function exists and can be called
	cmd := &Cmd{}

	// Just verify the function exists - we can't easily test the interactive part
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("waitForContinue should not panic, got panic: %v", r)
		}
	}()

	// We can't actually test the interactive input in unit tests
	// but we can verify the function is callable
	_ = cmd.waitForContinue
}

// TestCmd_InteractiveWorkflowIntegration tests workflow integration with Cmd
func TestCmd_InteractiveWorkflowIntegration(t *testing.T) {
	// Setup
	mockClient := &mockGitClient{}
	cm := config.NewConfigManager(mockClient)
	cmd, _ := NewCmd(mockClient, cm)

	ui := interactive.NewUI(mockClient, nil, nil, cmd)
	if ui == nil {
		t.Fatal("interactive.NewUI returned nil")
	}

	// Test workflow operations using exported methods
	if got := ui.AddToWorkflow("add", []string{"."}, "add ."); got != 1 {
		t.Errorf("AddToWorkflow first ID = %d, want 1", got)
	}
	if got := ui.AddToWorkflow("commit", []string{"-m", "msg"}, "commit"); got != 2 {
		t.Errorf("AddToWorkflow second ID = %d, want 2", got)
	}

	ui.ClearWorkflow()
	if got := ui.AddToWorkflow("status", nil, "status"); got != 1 {
		t.Errorf("AddToWorkflow after ClearWorkflow ID = %d, want 1", got)
	}
}
