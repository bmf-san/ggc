package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v7/internal/config"
)

// TestCmd_Route_DebugKeys verifies that the debug-keys command routes correctly
func TestCmd_Route_DebugKeys(t *testing.T) {
	t.Parallel()

	mockClient := &mockGitClient{}
	cm := config.NewConfigManager(mockClient)
	cmd := NewCmd(mockClient, cm)

	var buf bytes.Buffer
	cmd.outputWriter = &buf
	cmd.debugger.outputWriter = &buf

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Route() should not panic for debug-keys, but got: %v", r)
		}
	}()

	cmd.Route([]string{"debug-keys"})

	out := buf.String()
	if !strings.Contains(out, "=== Active Key Bindings ===") {
		t.Fatalf("expected debug-keys output, got: %q", out)
	}
}

// TestCmd_Route_DebugKeysHelp verifies that debug-keys --help displays help output
func TestCmd_Route_DebugKeysHelp(t *testing.T) {
	t.Parallel()

	mockClient := &mockGitClient{}
	cm := config.NewConfigManager(mockClient)
	cmd := NewCmd(mockClient, cm)

	var buf bytes.Buffer
	cmd.outputWriter = &buf
	cmd.debugger.outputWriter = &buf

	cmd.Route([]string{"debug-keys", "--help"})

	out := buf.String()
	if !strings.Contains(out, "debug-keys - Debug keybinding issues and capture raw key sequences") {
		t.Fatalf("expected debug-keys help output, got: %q", out)
	}
}

// TestCmd_Route_UnknownInputs verifies that previously "legacy-like" inputs
// now fall through to the standard routing behavior (showing help or handling gracefully)
// instead of showing a specific "legacy-like syntax not supported" error.
func TestCmd_Route_UnknownInputs(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name             string
		args             []string
		expectPanic      bool
		shouldContain    []string
		shouldNotContain string
	}{
		{
			name:             "unknown hyphenated command",
			args:             []string{"clean-interactive"},
			expectPanic:      false,
			shouldNotContain: "legacy-like syntax is not supported",
		},
		{
			name:             "unknown hyphenated command 2",
			args:             []string{"add-interactive"},
			expectPanic:      false,
			shouldNotContain: "legacy-like syntax is not supported",
		},
		{
			name:             "flag-like first token",
			args:             []string{"--prune"},
			expectPanic:      false,
			shouldNotContain: "legacy-like syntax is not supported",
		},
		{
			name:             "single dash as first token",
			args:             []string{"-"},
			expectPanic:      false,
			shouldNotContain: "legacy-like syntax is not supported",
		},
		{
			name:             "double dash as first token",
			args:             []string{"--"},
			expectPanic:      false,
			shouldNotContain: "legacy-like syntax is not supported",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockClient := &mockGitClient{}
			cm := config.NewConfigManager(mockClient)
			cmd := NewCmd(mockClient, cm)

			var buf bytes.Buffer
			cmd.outputWriter = &buf

			defer func() {
				r := recover()
				if tc.expectPanic && r == nil {
					t.Fatalf("expected panic for args %v, but did not panic", tc.args)
				}
				if !tc.expectPanic && r != nil {
					t.Fatalf("Route() should not panic for args %v, but got: %v", tc.args, r)
				}
			}()

			cmd.Route(tc.args)

			out := buf.String()

			// Verify that the old legacy-like error is not shown
			if tc.shouldNotContain != "" && strings.Contains(out, tc.shouldNotContain) {
				t.Fatalf("output should not contain %q, got: %q", tc.shouldNotContain, out)
			}

			// Verify expected content is present
			for _, expected := range tc.shouldContain {
				if !strings.Contains(out, expected) {
					t.Fatalf("expected output to contain %q, got: %q", expected, out)
				}
			}
		})
	}
}
