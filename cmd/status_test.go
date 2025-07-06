package cmd

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestStatuseer_Status(t *testing.T) {
	cases := []struct {
		name           string
		args           []string
		expectedCmds   []string
		mockOutput     []byte
		mockError      error
		expectedOutput string
	}{
		{
			name:           "status no args",
			args:           []string{},
			expectedCmds:   []string{"git -c color.status=always status"},
			mockOutput:     []byte("On branch main\nChanges not staged for commit:\n  modified:   modified_file.go\n\nUntracked files:\n  untracked_file.go\n"),
			mockError:      nil,
			expectedOutput: "On branch main",
		},
		{
			name:           "status short",
			args:           []string{"short"},
			expectedCmds:   []string{"git -c color.status=always status --short"},
			mockOutput:     []byte("M  modified_file.go\n?? untracked_file.go\n"),
			mockError:      nil,
			expectedOutput: "M  modified_file.go",
		},
		{
			name:           "invalid command",
			args:           []string{"invalid"},
			expectedCmds:   nil,
			mockOutput:     nil,
			mockError:      nil,
			expectedOutput: "Usage: ggc status",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			cmdIndex := 0
			s := &Statuseer{
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
					if tc.mockError != nil {
						return exec.Command("false")
					}
					return exec.Command("echo", string(tc.mockOutput))
				},
			}
			s.helper.outputWriter = &buf
			s.Status(tc.args)

			output := buf.String()
			if !strings.Contains(output, tc.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tc.expectedOutput, output)
			}
		})
	}
}
