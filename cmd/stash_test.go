package cmd

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestStasher_Stash(t *testing.T) {
	cases := []struct {
		name           string
		args           []string
		expectedCmds   []string
		mockOutput     []byte
		mockError      error
		expectedOutput string
	}{
		{
			name:           "stash trash success",
			args:           []string{"trash"},
			expectedCmds:   []string{"git stash drop"},
			mockOutput:     []byte("Dropped refs/stash@{0}"),
			mockError:      nil,
			expectedOutput: "Dropped refs/stash@{0}",
		},
		{
			name:           "stash trash error",
			args:           []string{"trash"},
			expectedCmds:   []string{"git stash drop"},
			mockOutput:     nil,
			mockError:      errors.New("no stash found"),
			expectedOutput: "Error: no stash found",
		},
		{
			name:           "no args",
			args:           []string{},
			expectedCmds:   nil,
			mockOutput:     nil,
			mockError:      nil,
			expectedOutput: "Usage: ggc stash [command]",
		},
		{
			name:           "invalid command",
			args:           []string{"invalid"},
			expectedCmds:   nil,
			mockOutput:     nil,
			mockError:      nil,
			expectedOutput: "Usage: ggc stash [command]",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			cmdIndex := 0
			s := &Stasher{
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

			s.Stash(tc.args)

			output := buf.String()
			if !strings.Contains(output, tc.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tc.expectedOutput, output)
			}
		})
	}
}

func TestStasher_Stash_Trash(t *testing.T) {
	var buf bytes.Buffer
	stasher := &Stasher{
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, _  ...string) *exec.Cmd {
			return exec.Command("echo")
		},
	}
	stasher.helper.outputWriter = &buf

	stasher.Stash([]string{"trash"})

	expected := "Dropped refs/stash@{0}\n"
	if got := buf.String(); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

func TestStasher_Stash_Trash_Error(t *testing.T) {
	var buf bytes.Buffer
	stasher := &Stasher{
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, _  ...string) *exec.Cmd {
			return exec.Command("false")
		},
	}
	stasher.helper.outputWriter = &buf

	stasher.Stash([]string{"trash"})

	expected := "Error: no stash found\n"
	if got := buf.String(); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

func TestStasher_Stash_Help(t *testing.T) {
	var buf bytes.Buffer
	stasher := &Stasher{
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, _  ...string) *exec.Cmd {
			return exec.Command("echo")
		},
	}
	stasher.helper.outputWriter = &buf
	stasher.Stash([]string{})
	output := buf.String()
	if !strings.Contains(output, "Usage") {
		t.Errorf("Usage should be displayed, but got: %s", output)
	}
}

func TestStasher_Stash_Unknown(t *testing.T) {
	var buf bytes.Buffer
	stasher := &Stasher{
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, _  ...string) *exec.Cmd {
			return exec.Command("echo")
		},
	}
	stasher.helper.outputWriter = &buf
	stasher.Stash([]string{"unknown"})
	output := buf.String()
	if !strings.Contains(output, "Usage") {
		t.Errorf("Usage should be displayed for unknown command, but got: %s", output)
	}
}
