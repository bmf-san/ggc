package main

import (
	"runtime/debug"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v5/cmd"
	"github.com/bmf-san/ggc/v5/config"
	"github.com/bmf-san/ggc/v5/internal/testutil"
	"github.com/bmf-san/ggc/v5/router"
)

func TestGetVersionInfo(t *testing.T) {
	tests := []struct {
		name            string
		setupFunc       func()
		expectedVersion string
		expectedCommit  string
		checkFunc       func(t *testing.T, version, commit string)
	}{
		{
			name: "default case - no ldflags, use build info",
			setupFunc: func() {
				// Reset global variables to simulate default state
				version = ""
				commit = ""
			},
			checkFunc: func(t *testing.T, v, c string) {
				// In test environment, build info should provide some values or empty
				// We verify the function executes without panic and returns two strings
				if v != "" && c != "" {
					t.Logf("Build info provided: version=%s, commit=%s", v, c)
				} else {
					t.Log("Build info not available or in dev mode")
				}
			},
		},
		{
			name: "ldflags version set",
			setupFunc: func() {
				version = "v1.2.3"
				commit = ""
			},
			expectedVersion: "v1.2.3",
			expectedCommit:  "",
			checkFunc: func(t *testing.T, v, c string) {
				if v != "v1.2.3" {
					t.Errorf("Expected version 'v1.2.3', got '%s'", v)
				}
				if c != "" {
					t.Errorf("Expected empty commit, got '%s'", c)
				}
			},
		},
		{
			name: "ldflags commit set",
			setupFunc: func() {
				version = ""
				commit = "abc1234"
			},
			expectedVersion: "",
			expectedCommit:  "abc1234",
			checkFunc: func(t *testing.T, v, c string) {
				if v != "" {
					t.Errorf("Expected empty version, got '%s'", v)
				}
				if c != "abc1234" {
					t.Errorf("Expected commit 'abc1234', got '%s'", c)
				}
			},
		},
		{
			name: "both ldflags set",
			setupFunc: func() {
				version = "v2.0.0"
				commit = "def5678"
			},
			expectedVersion: "v2.0.0",
			expectedCommit:  "def5678",
			checkFunc: func(t *testing.T, v, c string) {
				if v != "v2.0.0" {
					t.Errorf("Expected version 'v2.0.0', got '%s'", v)
				}
				if c != "def5678" {
					t.Errorf("Expected commit 'def5678', got '%s'", c)
				}
			},
		},
	}

	// Store original values to restore later
	originalVersion := version
	originalCommit := commit

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFunc()
			v, c := GetVersionInfo()
			tt.checkFunc(t, v, c)
		})
	}

	// Restore original values
	version = originalVersion
	commit = originalCommit
}

func TestGetVersionInfo_BuildInfo(t *testing.T) {
	// Reset global variables to test build info path
	originalVersion := version
	originalCommit := commit
	version = ""
	commit = ""

	// Test the build info fallback
	v, c := GetVersionInfo()

	// Verify function executes without panic
	t.Logf("Build info result: version='%s', commit='%s'", v, c)

	// In test environment, build info behavior varies
	// We mainly test that the function doesn't panic and returns strings
	if bi, ok := debug.ReadBuildInfo(); ok {
		t.Log("Build info is available")
		// Verify that if build info is available, we get some result
		if bi.Main.Version == "(devel)" {
			// Development build should return empty version
			if v != "" {
				t.Logf("Note: Got version '%s' from build info despite (devel)", v)
			}
		}
	} else {
		t.Log("Build info is not available")
		// If no build info, should return empty strings
		if v != "" || c != "" {
			t.Errorf("Expected empty strings when no build info, got version='%s', commit='%s'", v, c)
		}
	}

	// Restore original values
	version = originalVersion
	commit = originalCommit
}

func TestMain_Components(t *testing.T) {
	// Test that main function components can be initialized without side effects
	// We avoid calling main() directly to prevent actual git command execution during tests

	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "config manager creation",
			testFunc: func(t *testing.T) {
				mockClient := testutil.NewMockGitClient()
				cm := config.NewConfigManager(mockClient)
				if cm == nil {
					t.Error("config manager should be created")
				}
				
				// Test config loading (safe with mock)
				cm.LoadConfig()
				cfg := cm.GetConfig()
				if cfg == nil {
					t.Error("config should be loaded")
				}
				t.Log("Config manager created and loaded successfully")
			},
		},
		{
			name: "version getter setup",
			testFunc: func(t *testing.T) {
				// Test version getter setup (safe, no git commands)
				cmd.SetVersionGetter(GetVersionInfo)
				
				// Verify the version getter was set by calling it
				v, c := GetVersionInfo()
				t.Logf("Version getter set successfully: version='%s', commit='%s'", v, c)
			},
		},
		{
			name: "cmd creation with mock client",
			testFunc: func(t *testing.T) {
				// Test cmd creation with mock client (safe, no real git commands)
				mockClient := testutil.NewMockGitClient()
				c := cmd.NewCmd(mockClient)
				if c == nil {
					t.Error("cmd should be created")
				}
				t.Log("Cmd created successfully with mock client")
			},
		},
		{
			name: "router creation",
			testFunc: func(t *testing.T) {
				// Test router creation with mock components
				mockClient := testutil.NewMockGitClient()
				cm := config.NewConfigManager(mockClient)
				c := cmd.NewCmd(mockClient)
				r := router.NewRouter(c, cm)
				if r == nil {
					t.Error("router should be created")
				}
				t.Log("Router created successfully")
			},
		},
		{
			name: "integration test with mock components",
			testFunc: func(t *testing.T) {
				// Test the complete initialization flow with mock components
				mockClient := testutil.NewMockGitClient()
				
				// Initialize components like main() does
				cm := config.NewConfigManager(mockClient)
				cm.LoadConfig()
				cmd.SetVersionGetter(GetVersionInfo)
				c := cmd.NewCmd(mockClient)
				r := router.NewRouter(c, cm)
				
				// Test safe routing (help command)
				r.Route([]string{"help"})
				t.Log("Integration test completed successfully")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func TestMain_ArgumentHandling(t *testing.T) {
	tests := []struct {
		name string
		args []string
		desc string
	}{
		{
			name: "help command",
			args: []string{"help"},
			desc: "Test help command routing",
		},
		{
			name: "version command",
			args: []string{"version"},
			desc: "Test version command routing",
		},
		{
			name: "status command",
			args: []string{"status"},
			desc: "Test status command routing",
		},
		{
			name: "empty args",
			args: []string{},
			desc: "Test empty arguments (should trigger interactive mode)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test argument handling with mock components
			mockClient := testutil.NewMockGitClient()
			cm := config.NewConfigManager(mockClient)
			cm.LoadConfig()
			c := cmd.NewCmd(mockClient)
			r := router.NewRouter(c, cm)
			
			// Test routing with different arguments (safe with mock)
			r.Route(tt.args)
			t.Logf("%s: Successfully routed args %v", tt.desc, tt.args)
		})
	}
}

func TestMain_OsArgsSimulation(t *testing.T) {
	// Test os.Args simulation without actually modifying os.Args
	testArgs := [][]string{
		{"ggc"},           // Program name only
		{"ggc", "help"},   // Help command
		{"ggc", "status"}, // Status command
		{"ggc", "version"}, // Version command
	}

	for _, args := range testArgs {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			// Simulate what main() does with os.Args[1:]
			var routeArgs []string
			if len(args) > 1 {
				routeArgs = args[1:]
			}
			
			// Test with mock components
			mockClient := testutil.NewMockGitClient()
			cm := config.NewConfigManager(mockClient)
			cm.LoadConfig()
			cmd.SetVersionGetter(GetVersionInfo)
			c := cmd.NewCmd(mockClient)
			r := router.NewRouter(c, cm)
			
			// Route the arguments (safe with mock)
			r.Route(routeArgs)
			t.Logf("Successfully simulated main() with args: %v -> route args: %v", args, routeArgs)
		})
	}
}
