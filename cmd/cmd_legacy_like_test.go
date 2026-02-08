package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v7/internal/config"
)

func TestCmd_Route_LegacyLikeError_Extended(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
	}{
		{name: "top-level hyphenated clean-interactive", args: []string{"clean-interactive"}},
		{name: "top-level hyphenated add-interactive", args: []string{"add-interactive"}},
		{name: "short flag after command", args: []string{"rebase", "-i"}},
		{name: "long flag after command", args: []string{"fetch", "--prune"}},
		{name: "multiple flags after command", args: []string{"push", "-f", "--force-with-lease"}},
		{name: "single dash after command", args: []string{"commit", "-"}},
		{name: "double dash as command token", args: []string{"--"}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockClient := &mockGitClient{}
			cm := config.NewConfigManager(mockClient)
			cmd := NewCmd(mockClient, cm)

			var buf bytes.Buffer
			cmd.outputWriter = &buf // capture legacy-like error output

			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("Route() should not panic for args %v, but got: %v", tc.args, r)
				}
			}()

			cmd.Route(tc.args)

			out := buf.String()
			if !strings.Contains(out, "legacy-like syntax is not supported") {
				t.Fatalf("expected legacy-like error message, got: %q", out)
			}
		})
	}
}

func TestCmd_Route_LegacyLike_AllowsRegisteredHyphenatedCommand(t *testing.T) {
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
	if strings.Contains(out, "legacy-like syntax is not supported") {
		t.Fatalf("did not expect legacy-like error for debug-keys, got: %q", out)
	}
	if !strings.Contains(out, "=== Active Key Bindings ===") {
		t.Fatalf("expected debug-keys output, got: %q", out)
	}
}

func TestCmd_Route_LegacyLike_AllowsDebugKeysHelpFlags(t *testing.T) {
	t.Parallel()

	mockClient := &mockGitClient{}
	cm := config.NewConfigManager(mockClient)
	cmd := NewCmd(mockClient, cm)

	var buf bytes.Buffer
	cmd.outputWriter = &buf
	cmd.debugger.outputWriter = &buf

	cmd.Route([]string{"debug-keys", "--help"})

	out := buf.String()
	if strings.Contains(out, "legacy-like syntax is not supported") {
		t.Fatalf("did not expect legacy-like error for debug-keys --help, got: %q", out)
	}
	if !strings.Contains(out, "debug-keys - Debug keybinding issues and capture raw key sequences") {
		t.Fatalf("expected debug-keys help output, got: %q", out)
	}
}
