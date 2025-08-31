package cmd

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v5/internal/testutil"
)

func TestVersioneer_Version(t *testing.T) {
	cases := []struct {
		name           string
		args           []string
		expectedOutput []string
	}{
		{
			name: "version no args with default values",
			args: []string{},
			expectedOutput: []string{
				"ggc version",
				"commit:",
				"built:",
				"os/arch:",
			},
		},
		{
			name: "version no args with custom version info",
			args: []string{},
			expectedOutput: []string{
				"ggc version",
				"commit:",
				"built:",
				"os/arch:",
			},
		},
		{
			name: "version with args shows help",
			args: []string{"help"},
			expectedOutput: []string{
				"Usage:",
			},
		},
		{
			name: "version with multiple args shows help",
			args: []string{"invalid", "args"},
			expectedOutput: []string{
				"Usage:",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock getVersionInfo function
			SetVersionGetter(func() (string, string) {
				return "v1.0.0", "abc123"
			})

			var buf bytes.Buffer
			v := &Versioneer{
				gitClient:    testutil.NewMockGitClient(),
				outputWriter: &buf,
				helper:       NewHelper(),
				execCommand:  exec.Command,
			}

			v.helper.outputWriter = &buf
			v.Version(tc.args)

			output := buf.String()
			for _, expected := range tc.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("expected output to contain %q, got %q", expected, output)
				}
			}
		})
	}
}
