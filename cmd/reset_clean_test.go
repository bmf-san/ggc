package cmd

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestResetCleaner_ResetClean(t *testing.T) {
	cases := []struct {
		name           string
		expectedCmds   []string
		mockOutput     []byte
		mockError      error
		expectedOutput string
	}{
		{
			name:           "successful execution",
			expectedCmds:   []string{"git reset --hard HEAD", "git clean -fd"},
			mockOutput:     []byte("operation successful"),
			mockError:      nil,
			expectedOutput: "operation successful",
		},
		{
			name:           "reset error",
			expectedCmds:   []string{"git reset --hard HEAD"},
			mockOutput:     nil,
			mockError:      errors.New("reset error"),
			expectedOutput: "Error resetting changes: reset error",
		},
		{
			name:           "clean error",
			expectedCmds:   []string{"git reset --hard HEAD", "git clean -fd"},
			mockOutput:     nil,
			mockError:      errors.New("clean error"),
			expectedOutput: "Error cleaning untracked files: clean error",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			cmdIndex := 0
			r := &ResetCleaner{
				outputWriter: &buf,
				helper:       NewHelper(),
				execCommand: func(_ string, args ...string) *exec.Cmd {
					if cmdIndex < len(tc.expectedCmds) {
						gotCmd := strings.Join(append([]string{"git"}, args...), " ")
						if gotCmd != tc.expectedCmds[cmdIndex] {
							t.Errorf("expected command %q, got %q", tc.expectedCmds[cmdIndex], gotCmd)
						}
					}
					cmdIndex++
					if tc.mockError != nil && cmdIndex == len(tc.expectedCmds) {
						return exec.Command("false")
					}
					return exec.Command("echo", string(tc.mockOutput))
				},
			}
			r.helper.outputWriter = &buf

			r.ResetClean()

			output := buf.String()
			if !strings.Contains(output, tc.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tc.expectedOutput, output)
			}
		})
	}
}
