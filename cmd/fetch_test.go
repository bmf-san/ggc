package cmd

import (
	"bytes"
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
			expectedOutput: "",
			mockOutput:     []byte(""),
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
			mockClient := &mockAddGitClient{}
			f := &Fetcher{
				gitClient:    mockClient,
				outputWriter: &buf,
				helper:       NewHelper(),
			}
			f.helper.outputWriter = &buf

			f.Fetch(tc.args)

			output := buf.String()
			if tc.expectedOutput == "" {
				// For empty expected output, check that no error message is present
				if strings.Contains(output, "Error:") {
					t.Errorf("expected no error output, got %q", output)
				}
			} else if !strings.Contains(output, tc.expectedOutput) {
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
	}
	f.helper.outputWriter = &buf

	f.Fetch([]string{})

	output := buf.String()
	// When no arguments are provided, it shows help message instead of error
	if !strings.Contains(output, "Usage: ggc fetch") {
		t.Errorf("Expected help message, got: %s", output)
	}
}
