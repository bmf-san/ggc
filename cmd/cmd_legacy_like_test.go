package cmd

import (
	"bytes"
	"strings"
	"testing"
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
			cmd := NewCmd(mockClient)

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
