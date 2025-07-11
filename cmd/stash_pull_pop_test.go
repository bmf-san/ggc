package cmd

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestStashPullPopper_StashPullPop(t *testing.T) {
	cases := []struct {
		name           string
		expectedCmds   []string
		mockOutput     []byte
		mockError      error
		expectedOutput string
	}{
		{
			name:           "successful execution",
			expectedCmds:   []string{"git pull --autostash"},
			mockOutput:     []byte("operation successful with autostash"),
			mockError:      nil,
			expectedOutput: "operation successful with autostash",
		},
		{
			name:           "pull with autostash error",
			expectedCmds:   []string{"git pull --autostash"},
			mockOutput:     nil,
			mockError:      errors.New("exit status 1"),
			expectedOutput: "Error pulling changes with autostash: exit status 1",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			cmdIndex := 0
			s := &StashPullPopper{
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
			s.helper.outputWriter = &buf

			s.StashPullPop()

			output := buf.String()
			if !strings.Contains(output, tc.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tc.expectedOutput, output)
			}
		})
	}
}
