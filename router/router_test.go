package router

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v7/config"
	"github.com/bmf-san/ggc/v7/internal/testutil"
)

type mockExecuter struct {
	helpCalled        bool
	branchCalled      bool
	branchArgs        []string
	commitCalled      bool
	commitArgs        []string
	logCalled         bool
	logArgs           []string
	diffCalled        bool
	diffArgs          []string
	statusCalled      bool
	statusArgs        []string
	pullCalled        bool
	pullArgs          []string
	pushCalled        bool
	pushArgs          []string
	resetCalled       bool
	resetArgs         []string
	tagCalled         bool
	tagArgs           []string
	versionCalled     bool
	versionArgs       []string
	cleanCalled       bool
	cleanArgs         []string
	configCalled      bool
	configArgs        []string
	hookerCalled      bool
	hookerArgs        []string
	restoreCalled     bool
	restoreArgs       []string
	addCalled         bool
	addArgs           []string
	remoteCalled      bool
	remoteArgs        []string
	rebaseCalled      bool
	rebaseArgs        []string
	stashCalled       bool
	stashArgs         []string
	fetchCalled       bool
	fetchArgs         []string
	debugKeysCalled   bool
	debugKeysArgs     []string
	interactiveCalled bool
}

func (m *mockExecuter) Help() {
	m.helpCalled = true
}

func (m *mockExecuter) Branch(args []string) {
	m.branchCalled = true
	m.branchArgs = args
}

func (m *mockExecuter) Commit(args []string) {
	m.commitCalled = true
	m.commitArgs = args
}

func (m *mockExecuter) Log(args []string) {
	m.logCalled = true
	m.logArgs = args
}

func (m *mockExecuter) Config(args []string) {
	m.configCalled = true
	m.configArgs = args
}

func (m *mockExecuter) Hook(args []string) {
	m.hookerCalled = true
	m.hookerArgs = args
}

func (m *mockExecuter) Status(args []string) {
	m.statusCalled = true
	m.statusArgs = args
}

func (m *mockExecuter) Version(args []string) {
	m.versionCalled = true
	m.versionArgs = args
}

func (m *mockExecuter) Diff(args []string) {
	m.diffCalled = true
	m.diffArgs = args
}

func (m *mockExecuter) Add(args []string) {
	m.addCalled = true
	m.addArgs = args
}

func (m *mockExecuter) Remote(args []string) {
	m.remoteCalled = true
	m.remoteArgs = args
}

func (m *mockExecuter) Rebase(args []string) {
	m.rebaseCalled = true
	m.rebaseArgs = args
}

func (m *mockExecuter) Stash(args []string) {
	m.stashCalled = true
	m.stashArgs = args
}

func (m *mockExecuter) Fetch(args []string) {
	m.fetchCalled = true
	m.fetchArgs = args
}

func (m *mockExecuter) Restore(args []string) {
	m.restoreCalled = true
	m.restoreArgs = args
}

func (m *mockExecuter) Tag(args []string) {
	m.tagCalled = true
	m.tagArgs = args
}

func (m *mockExecuter) Pull(args []string) {
	m.pullCalled = true
	m.pullArgs = args
}

func (m *mockExecuter) Push(args []string) {
	m.pushCalled = true
	m.pushArgs = args
}

func (m *mockExecuter) Reset(args []string) {
	m.resetCalled = true
	m.resetArgs = args
}

func (m *mockExecuter) Clean(args []string) {
	m.cleanCalled = true
	m.cleanArgs = args
}

func (m *mockExecuter) DebugKeys(args []string) {
	m.debugKeysCalled = true
	m.debugKeysArgs = args
}

func (m *mockExecuter) Interactive() {
	m.interactiveCalled = true
}

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()

	orig := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	os.Stderr = w

	var buf bytes.Buffer
	done := make(chan struct{})
	go func() {
		_, _ = io.Copy(&buf, r)
		close(done)
	}()

	fn()

	_ = w.Close()
	os.Stderr = orig
	<-done
	_ = r.Close()

	return buf.String()
}

func TestRouter(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		validate func(t *testing.T, m *mockExecuter)
	}{
		{
			name: "help",
			args: []string{"help"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.helpCalled {
					t.Error("Help should be called")
				}
			},
		},
		{
			name: "branch",
			args: []string{"branch", "current"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.branchCalled {
					t.Error("Branch should be called")
				}
				if len(m.branchArgs) != 1 || m.branchArgs[0] != "current" {
					t.Errorf("unexpected branch args: got %v", m.branchArgs)
				}
			},
		},
		{
			name: "commit",
			args: []string{"commit", "allow-empty"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.commitCalled {
					t.Error("Commit should be called")
				}
				if len(m.commitArgs) != 1 || m.commitArgs[0] != "allow-empty" {
					t.Errorf("unexpected commit args: got %v", m.commitArgs)
				}
			},
		},
		{
			name: "log",
			args: []string{"log", "simple"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.logCalled {
					t.Error("Log should be called")
				}
				if len(m.logArgs) != 1 || m.logArgs[0] != "simple" {
					t.Errorf("unexpected log args: got %v", m.logArgs)
				}
			},
		},
		{
			name: "pull",
			args: []string{"pull", "rebase"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.pullCalled {
					t.Error("Pull should be called")
				}
				if len(m.pullArgs) != 1 || m.pullArgs[0] != "rebase" {
					t.Errorf("unexpected pull args: got %v", m.pullArgs)
				}
			},
		},
		{
			name: "push",
			args: []string{"push", "force"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.pushCalled {
					t.Error("Push should be called")
				}
				if len(m.pushArgs) != 1 || m.pushArgs[0] != "force" {
					t.Errorf("unexpected push args: got %v", m.pushArgs)
				}
			},
		},
		{
			name: "reset",
			args: []string{"reset"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.resetCalled {
					t.Error("Reset should be called")
				}
			},
		},
		{
			name: "hooker no args",
			args: []string{"hook"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.hookerCalled {
					t.Error("Hooker should be called")
				}
			},
		},
		{
			name: "hooker list",
			args: []string{"hook", "list"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.hookerCalled {
					t.Error("Hooker should be called")
				}
			},
		},
		{
			name: "restore no args",
			args: []string{"restore"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.restoreCalled {
					t.Error("restore should be called")
				}
			},
		},
		{
			name: "restore all",
			args: []string{"restore", "."},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.restoreCalled {
					t.Error("restore should be called")
				}
			},
		},
		{
			name: "config no args",
			args: []string{"config"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.configCalled {
					t.Error("config should be called")
				}
			},
		},
		{
			name: "config with list",
			args: []string{"config", "list"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.configCalled {
					t.Error("config should be called")
				}
			},
		},
		{
			name: "status no args",
			args: []string{"status"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.statusCalled {
					t.Error("Status should be called")
				}
				if len(m.statusArgs) != 0 {
					t.Errorf("unexpected status args: got %v, expected empty", m.statusArgs)
				}
			},
		},
		{
			name: "status with short arg",
			args: []string{"status", "short"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.statusCalled {
					t.Error("Status should be called")
				}
				if len(m.statusArgs) != 1 || m.statusArgs[0] != "short" {
					t.Errorf("unexpected status args: got %v, expected [short]", m.statusArgs)
				}
			},
		},
		{
			name: "tag no args",
			args: []string{"tag"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.tagCalled {
					t.Error("Tag should be called")
				}
				if len(m.tagArgs) != 0 {
					t.Errorf("unexpected status args: got %v, expected empty", m.tagArgs)
				}
			},
		},
		{
			name: "tag with arg",
			args: []string{"tag", "list"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.tagCalled {
					t.Error("Tag should be called")
				}
				if len(m.tagArgs) != 1 || m.tagArgs[0] != "list" {
					t.Errorf("unexpected tag args: got %v, expected [list]", m.tagArgs)
				}
			},
		},
		{
			name: "version",
			args: []string{"version"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.versionCalled {
					t.Error("Version should be called")
				}
			},
		},
		{
			name: "diff no args",
			args: []string{"diff"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.diffCalled {
					t.Error("Diff should be called")
				}
				if len(m.diffArgs) != 0 {
					t.Errorf("unexpected diff args: got %v, expected empty", m.diffArgs)
				}
			},
		},
		{
			name: "diff unstaged",
			args: []string{"diff", "unstaged"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.diffCalled {
					t.Error("Diff should be called")
				}
				if len(m.diffArgs) != 1 || m.diffArgs[0] != "unstaged" {
					t.Errorf("unexpected diff args: got %v, expected [unstaged]", m.diffArgs)
				}
			},
		},
		{
			name: "diff staged",
			args: []string{"diff", "staged"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.diffCalled {
					t.Error("Diff should be called")
				}
				if len(m.diffArgs) != 1 || m.diffArgs[0] != "staged" {
					t.Errorf("unexpected diff args: got %v, expected [staged]", m.diffArgs)
				}
			},
		},
		{
			name: "clean",
			args: []string{"clean", "files"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.cleanCalled {
					t.Error("Clean should be called")
				}
				if len(m.cleanArgs) != 1 || m.cleanArgs[0] != "files" {
					t.Errorf("unexpected clean args: got %v", m.cleanArgs)
				}
			},
		},
		{
			name: "add",
			args: []string{"add", "."},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.addCalled {
					t.Error("Add should be called")
				}
				if len(m.addArgs) != 1 || m.addArgs[0] != "." {
					t.Errorf("unexpected add args: got %v", m.addArgs)
				}
			},
		},
		{
			name: "remote",
			args: []string{"remote"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.remoteCalled {
					t.Error("remote should be called")
				}
			},
		},
		{
			name: "rebase",
			args: []string{"rebase"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.rebaseCalled {
					t.Error("Rebase should be called")
				}
			},
		},
		{
			name: "stash",
			args: []string{"stash"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.stashCalled {
					t.Error("stash should be called")
				}
			},
		},
		{
			name: "fetch",
			args: []string{"fetch"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.fetchCalled {
					t.Error("fetch should be called")
				}
			},
		},
		{
			name: "unknown",
			args: []string{"unknown"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.helpCalled {
					t.Error("Help should be called")
				}
			},
		},
		{
			name: "empty",
			args: []string{},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.interactiveCalled {
					t.Error("Interactive should be called")
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := &mockExecuter{}
			r := NewRouter(m, config.NewConfigManager(testutil.NewMockGitClient()))
			r.Route(tc.args)
			tc.validate(t, m)
		})
	}
}

func TestRouter_WithAliases(t *testing.T) {
	cases := []struct {
		name             string
		aliases          map[string]interface{}
		args             []string
		expectExit       bool
		expectedExitCode int
		validate         func(t *testing.T, m *mockExecuter)
	}{
		{
			name: "simple alias",
			aliases: map[string]interface{}{
				"st": "status",
			},
			args: []string{"st"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.statusCalled {
					t.Error("Status should be called for simple alias")
				}
			},
		},
		{
			name: "simple alias with args",
			aliases: map[string]interface{}{
				"br": "branch",
			},
			args: []string{"br", "new-branch"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.branchCalled {
					t.Error("Branch should be called for simple alias")
				}
				if len(m.branchArgs) != 1 || m.branchArgs[0] != "new-branch" {
					t.Errorf("unexpected branch args: got %v", m.branchArgs)
				}
			},
		},
		{
			name: "sequence alias",
			aliases: map[string]interface{}{
				"sync": []interface{}{"pull current", "push current"},
			},
			args: []string{"sync"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.pullCalled {
					t.Error("Pull should be called for sequence alias")
				}
				if !m.pushCalled {
					t.Error("Push should be called for sequence alias")
				}
			},
		},
		{
			name: "sequence alias with args (should be ignored)",
			aliases: map[string]interface{}{
				"sync": []interface{}{"pull current", "push current"},
			},
			args:             []string{"sync", "ignored"},
			expectExit:       true,
			expectedExitCode: 1,
			validate: func(t *testing.T, m *mockExecuter) {
				if m.pullCalled || m.pushCalled {
					t.Error("Sequence alias should not run when arguments are provided")
				}
			},
		},
		{
			name: "non-alias command",
			aliases: map[string]interface{}{
				"st": "status",
			},
			args: []string{"commit", "test"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.commitCalled {
					t.Error("Commit should be called for non-alias command")
				}
				if len(m.commitArgs) != 1 || m.commitArgs[0] != "test" {
					t.Errorf("unexpected commit args: got %v", m.commitArgs)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a real config manager and manually set aliases
			mockClient := testutil.NewMockGitClient()
			configManager := config.NewConfigManager(mockClient)
			configManager.LoadConfig()

			// Manually set aliases in the config
			cfg := configManager.GetConfig()
			cfg.Aliases = tc.aliases

			m := &mockExecuter{}
			r := NewRouter(m, configManager)
			exitCalled := false
			exitCode := 0
			r.SetExitFunc(func(code int) {
				exitCalled = true
				exitCode = code
			})
			r.Route(tc.args)
			tc.validate(t, m)
			if tc.expectExit {
				if !exitCalled {
					t.Fatal("expected router to exit but it did not")
				}
				if exitCode != tc.expectedExitCode {
					t.Fatalf("expected exit code %d, got %d", tc.expectedExitCode, exitCode)
				}
			} else if exitCalled {
				t.Fatalf("unexpected exit with code %d", exitCode)
			}
		})
	}
}

func TestRouter_SequenceAliasRejectsArguments(t *testing.T) {
	mockClient := testutil.NewMockGitClient()
	configManager := config.NewConfigManager(mockClient)
	configManager.LoadConfig()

	cfg := configManager.GetConfig()
	cfg.Aliases = map[string]interface{}{
		"deploy": []interface{}{"status"},
	}

	m := &mockExecuter{}
	r := NewRouter(m, configManager)
	exitCalled := false
	exitCode := 0
	r.SetExitFunc(func(code int) {
		exitCalled = true
		exitCode = code
	})

	output := captureStderr(t, func() {
		r.Route([]string{"deploy", "production"})
	})

	if !exitCalled {
		t.Fatal("sequence alias should exit when arguments are provided")
	}
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
	if m.statusCalled {
		t.Fatal("sequence commands should not execute when arguments are rejected")
	}
	if !strings.Contains(output, "sequence alias 'deploy' does not accept arguments") {
		t.Fatalf("expected rejection message, got %q", output)
	}
	if !strings.Contains(output, "production") {
		t.Fatalf("expected error to list offending arguments, got %q", output)
	}
}

func TestRouter_ConfigManagerNil(t *testing.T) {
	m := &mockExecuter{}
	r := NewRouter(m, nil)

	// Should not panic and should execute normal command
	r.Route([]string{"status"})

	if !m.statusCalled {
		t.Error("Status should be called when ConfigManager is nil")
	}
}

func TestRouter_AliasErrors(t *testing.T) {
	cases := []struct {
		name     string
		aliases  map[string]interface{}
		args     []string
		validate func(t *testing.T, m *mockExecuter)
	}{
		{
			name: "invalid alias format",
			aliases: map[string]interface{}{
				"invalid": 123, // Invalid format - should be string or []interface{}
			},
			args: []string{"invalid"},
			validate: func(t *testing.T, m *mockExecuter) {
				// Should not call any command due to error
				if m.statusCalled || m.branchCalled || m.commitCalled {
					t.Error("No command should be called for invalid alias")
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a real config manager and manually set aliases
			mockClient := testutil.NewMockGitClient()
			configManager := config.NewConfigManager(mockClient)
			configManager.LoadConfig()

			// Manually set aliases in the config
			cfg := configManager.GetConfig()
			cfg.Aliases = tc.aliases

			m := &mockExecuter{}
			r := NewRouter(m, configManager)
			r.Route(tc.args)
			tc.validate(t, m)
		})
	}
}

func TestRouter_AliasSequenceExecution(t *testing.T) {
	cases := []struct {
		name     string
		aliases  map[string]interface{}
		args     []string
		validate func(t *testing.T, m *mockExecuter)
	}{
		{
			name: "sequence alias with valid commands executes all",
			aliases: map[string]interface{}{
				"valid-seq": []interface{}{"status", "branch current", "log simple"},
			},
			args: []string{"valid-seq"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.statusCalled {
					t.Error("Status should be called for valid sequence")
				}
				if !m.branchCalled {
					t.Error("Branch should be called for valid sequence")
				}
				if !m.logCalled {
					t.Error("Log should be called for valid sequence")
				}
			},
		},
		{
			name: "sequence alias with valid command and arguments",
			aliases: map[string]interface{}{
				"valid-args": []interface{}{"branch current", "diff staged", "status short"},
			},
			args: []string{"valid-args"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.branchCalled {
					t.Error("Branch should be called")
				}
				if len(m.branchArgs) != 1 || m.branchArgs[0] != "current" {
					t.Errorf("Branch args incorrect: got %v, want [current]", m.branchArgs)
				}
				if !m.diffCalled {
					t.Error("Diff should be called")
				}
				if len(m.diffArgs) != 1 || m.diffArgs[0] != "staged" {
					t.Errorf("Diff args incorrect: got %v, want [staged]", m.diffArgs)
				}
				if !m.statusCalled {
					t.Error("Status should be called")
				}
				if len(m.statusArgs) != 1 || m.statusArgs[0] != "short" {
					t.Errorf("Status args incorrect: got %v, want [short]", m.statusArgs)
				}
			},
		},
		{
			name: "simple alias delegates to executeCommand",
			aliases: map[string]interface{}{
				"st": "status",
			},
			args: []string{"st"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.statusCalled {
					t.Error("Status should be called for simple alias")
				}
			},
		},
		{
			name: "sequence alias with multiple valid commands",
			aliases: map[string]interface{}{
				"multi-cmd": []interface{}{"status", "branch", "diff", "log", "tag"},
			},
			args: []string{"multi-cmd"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.statusCalled {
					t.Error("Status should be called")
				}
				if !m.branchCalled {
					t.Error("Branch should be called")
				}
				if !m.diffCalled {
					t.Error("Diff should be called")
				}
				if !m.logCalled {
					t.Error("Log should be called")
				}
				if !m.tagCalled {
					t.Error("Tag should be called")
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := testutil.NewMockGitClient()
			configManager := config.NewConfigManager(mockClient)
			configManager.LoadConfig()

			cfg := configManager.GetConfig()
			cfg.Aliases = tc.aliases

			m := &mockExecuter{}
			r := NewRouter(m, configManager)
			r.Route(tc.args)
			tc.validate(t, m)
		})
	}
}
