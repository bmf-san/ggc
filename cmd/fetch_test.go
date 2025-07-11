package cmd

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestFetcher_Fetch(t *testing.T) {
	cases := []struct {
		name           string
		args           []string
		expectedCmd    string
		expectedOutput string
		mockOutput     []byte
		mockError      error
	}{
		{
			name:           "fetch with prune",
			args:           []string{"--prune"},
			expectedCmd:    "git fetch --prune",
			expectedOutput: "mock output",
			mockOutput:     []byte("mock output"),
			mockError:      nil,
		},
		{
			name:           "fetch with no args",
			args:           []string{},
			expectedOutput: "Usage: ggc fetch [options]",
		},
		{
			name:           "fetch with invalid arg",
			args:           []string{"invalid"},
			expectedOutput: "Usage: ggc fetch [options]",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			f := &Fetcher{
				outputWriter: &buf,
				helper:       NewHelper(),
				execCommand: func(name string, args ...string) *exec.Cmd {
					if tc.expectedCmd != "" {
						gotCmd := strings.Join(append([]string{name}, args...), " ")
						if gotCmd != tc.expectedCmd {
							t.Errorf("expected command %q, got %q", tc.expectedCmd, gotCmd)
						}
					}
					return exec.Command("echo", string(tc.mockOutput))
				},
			}
			f.helper.outputWriter = &buf

			f.Fetch(tc.args)

			output := buf.String()
			if !strings.Contains(output, tc.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tc.expectedOutput, output)
			}
		})
	}
}

func TestFetcher_Fetch_Error(t *testing.T) {
	var buf bytes.Buffer
	f := &Fetcher{
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("false") // Command that fails
		},
	}
	f.helper.outputWriter = &buf

	f.Fetch([]string{})

	output := buf.String()
	// When no arguments are provided, it shows help message instead of error
	if !strings.Contains(output, "Usage: ggc fetch") {
		t.Errorf("Expected help message, got: %s", output)
	}
}
