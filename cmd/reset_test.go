package cmd

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestResetter_Reset(t *testing.T) {
	cases := []struct {
		name           string
		args           []string
		expectedCmds   []string
		mockOutput     []byte
		mockError      error
		expectedOutput string
	}{
		{
			name:           "reset clean success",
			args:           []string{"clean"},
			expectedCmds:   []string{"git reset --hard HEAD", "git clean -fd"},
			mockOutput:     []byte("Reset and clean successful"),
			mockError:      nil,
			expectedOutput: "Reset and clean successful",
		},
		{
			name:           "reset clean error on reset",
			args:           []string{"clean"},
			expectedCmds:   []string{"git reset --hard HEAD"},
			mockOutput:     nil,
			mockError:      errors.New("reset failed"),
			expectedOutput: "Error resetting changes: reset failed",
		},
		{
			name:           "reset clean error on clean",
			args:           []string{"clean"},
			expectedCmds:   []string{"git reset --hard HEAD", "git clean -fd"},
			mockOutput:     nil,
			mockError:      errors.New("clean failed"),
			expectedOutput: "Error cleaning untracked files: clean failed",
		},
		{
			name:           "no args",
			args:           []string{},
			expectedCmds:   nil,
			mockOutput:     nil,
			mockError:      nil,
			expectedOutput: "Usage: ggc reset",
		},
		{
			name:           "invalid command",
			args:           []string{"invalid"},
			expectedCmds:   nil,
			mockOutput:     nil,
			mockError:      nil,
			expectedOutput: "Usage: ggc reset",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			cmdIndex := 0
			r := &Resetter{
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

			r.Reset(tc.args)

			output := buf.String()
			if !strings.Contains(output, tc.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tc.expectedOutput, output)
			}
		})
	}
}
