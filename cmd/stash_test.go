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
			name:           "stash default success",
			args:           []string{},
			expectedCmds:   []string{"git stash"},
			mockOutput:     []byte("Saved working directory and index state"),
			mockError:      nil,
			expectedOutput: "Saved working directory and index state",
		},
		{
			name:           "stash default error",
			args:           []string{},
			expectedCmds:   []string{"git stash"},
			mockOutput:     nil,
			mockError:      errors.New("no changes to stash"),
			expectedOutput: "Error: no changes to stash",
		},
		{
			name:           "stash pop success",
			args:           []string{"pop"},
			expectedCmds:   []string{"git stash pop"},
			mockOutput:     []byte("Applied and dropped stash"),
			mockError:      nil,
			expectedOutput: "Applied and dropped stash",
		},
		{
			name:           "stash pop error",
			args:           []string{"pop"},
			expectedCmds:   []string{"git stash pop"},
			mockOutput:     nil,
			mockError:      errors.New("no stash found"),
			expectedOutput: "Error: no stash found",
		},
		{
			name:           "stash drop success",
			args:           []string{"drop"},
			expectedCmds:   []string{"git stash drop"},
			mockOutput:     []byte("Dropped stash"),
			mockError:      nil,
			expectedOutput: "Dropped stash",
		},
		{
			name:           "stash drop error",
			args:           []string{"drop"},
			expectedCmds:   []string{"git stash drop"},
			mockOutput:     nil,
			mockError:      errors.New("no stash found"),
			expectedOutput: "Error: no stash found",
		},
		{
			name:           "stash list success",
			args:           []string{"list"},
			expectedCmds:   []string{"git stash list"},
			mockOutput:     []byte("stash@{0}: WIP on main: 1234567 commit message"),
			mockError:      nil,
			expectedOutput: "stash@{0}: WIP on main: 1234567 commit message",
		},
		{
			name:           "stash list empty",
			args:           []string{"list"},
			expectedCmds:   []string{"git stash list"},
			mockOutput:     []byte(""),
			mockError:      nil,
			expectedOutput: "No stashes found",
		},
		{
			name:           "stash show success",
			args:           []string{"show"},
			expectedCmds:   []string{"git stash show"},
			mockOutput:     []byte(" file.txt | 2 +-\n 1 file changed, 1 insertion(+), 1 deletion(-)"),
			mockError:      nil,
			expectedOutput: "file.txt | 2 +-",
		},
		{
			name:           "stash show with ref",
			args:           []string{"show", "stash@{1}"},
			expectedCmds:   []string{"git stash show stash@{1}"},
			mockOutput:     []byte(" file.txt | 2 +-\n 1 file changed, 1 insertion(+), 1 deletion(-)"),
			mockError:      nil,
			expectedOutput: "file.txt | 2 +-",
		},
		{
			name:           "stash apply success",
			args:           []string{"apply"},
			expectedCmds:   []string{"git stash apply"},
			mockOutput:     []byte("Applied stash"),
			mockError:      nil,
			expectedOutput: "Applied stash",
		},
		{
			name:           "stash apply with ref",
			args:           []string{"apply", "stash@{1}"},
			expectedCmds:   []string{"git stash apply stash@{1}"},
			mockOutput:     []byte("Applied stash"),
			mockError:      nil,
			expectedOutput: "Applied stash",
		},
		{
			name:           "stash branch success",
			args:           []string{"branch", "feature-branch"},
			expectedCmds:   []string{"git stash branch feature-branch"},
			mockOutput:     []byte("Created branch"),
			mockError:      nil,
			expectedOutput: "Created branch 'feature-branch' from stash",
		},
		{
			name:           "stash branch with ref",
			args:           []string{"branch", "feature-branch", "stash@{1}"},
			expectedCmds:   []string{"git stash branch feature-branch stash@{1}"},
			mockOutput:     []byte("Created branch"),
			mockError:      nil,
			expectedOutput: "Created branch 'feature-branch' from stash",
		},
		{
			name:           "stash push success",
			args:           []string{"push", "-m", "message"},
			expectedCmds:   []string{"git stash push -m message"},
			mockOutput:     []byte("Saved working directory"),
			mockError:      nil,
			expectedOutput: "Saved working directory and index state",
		},
		{
			name:           "stash save success",
			args:           []string{"save", "message"},
			expectedCmds:   []string{"git stash save message"},
			mockOutput:     []byte("Saved working directory"),
			mockError:      nil,
			expectedOutput: "Saved working directory and index state",
		},
		{
			name:           "stash clear success",
			args:           []string{"clear"},
			expectedCmds:   []string{"git stash clear"},
			mockOutput:     []byte(""),
			mockError:      nil,
			expectedOutput: "Removed all stashes",
		},
		{
			name:           "stash create success",
			args:           []string{"create"},
			expectedCmds:   []string{"git stash create"},
			mockOutput:     []byte("abc123def456"),
			mockError:      nil,
			expectedOutput: "abc123def456",
		},
		{
			name:           "stash store success",
			args:           []string{"store", "abc123def456"},
			expectedCmds:   []string{"git stash store abc123def456"},
			mockOutput:     []byte(""),
			mockError:      nil,
			expectedOutput: "Stored stash",
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

					// Handle commands that need Output() vs Run()
					switch args[0] {
					case "stash":
						if len(args) > 1 && (args[1] == "list" || args[1] == "show" || args[1] == "create") {
							if tc.mockError != nil {
								return exec.Command("false")
							}
							return exec.Command("echo", "-n", string(tc.mockOutput))
						}
					}

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

func TestStasher_Stash_Drop(t *testing.T) {
	var buf bytes.Buffer
	stasher := &Stasher{
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("echo")
		},
	}
	stasher.helper.outputWriter = &buf

	stasher.Stash([]string{"drop"})

	expected := "Dropped stash\n"
	if got := buf.String(); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

func TestStasher_Stash_Drop_Error(t *testing.T) {
	var buf bytes.Buffer
	stasher := &Stasher{
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("false")
		},
	}
	stasher.helper.outputWriter = &buf

	stasher.Stash([]string{"drop"})

	expected := "Error: no stash found\n"
	if got := buf.String(); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

func TestStasher_Stash_Pop(t *testing.T) {
	var buf bytes.Buffer
	stasher := &Stasher{
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("echo")
		},
	}
	stasher.helper.outputWriter = &buf

	stasher.Stash([]string{"pop"})

	expected := "Applied and dropped stash\n"
	if got := buf.String(); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

func TestStasher_Stash_Pop_Error(t *testing.T) {
	var buf bytes.Buffer
	stasher := &Stasher{
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("false")
		},
	}
	stasher.helper.outputWriter = &buf

	stasher.Stash([]string{"pop"})

	expected := "Error: no stash found\n"
	if got := buf.String(); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

func TestStasher_Stash_Default(t *testing.T) {
	var buf bytes.Buffer
	stasher := &Stasher{
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("echo")
		},
	}
	stasher.helper.outputWriter = &buf

	stasher.Stash([]string{})

	expected := "Saved working directory and index state\n"
	if got := buf.String(); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

func TestStasher_Stash_Default_Error(t *testing.T) {
	var buf bytes.Buffer
	stasher := &Stasher{
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("false")
		},
	}
	stasher.helper.outputWriter = &buf

	stasher.Stash([]string{})

	expected := "Error: no changes to stash\n"
	if got := buf.String(); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

func TestStasher_Stash_Unknown(t *testing.T) {
	var buf bytes.Buffer
	stasher := &Stasher{
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, _ ...string) *exec.Cmd {
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

func TestStasher_Stash_BranchMissingName(t *testing.T) {
	var buf bytes.Buffer
	stasher := &Stasher{
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("echo")
		},
	}
	stasher.helper.outputWriter = &buf

	stasher.Stash([]string{"branch"})

	expected := "Error: branch name required\n"
	if got := buf.String(); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

func TestStasher_Stash_StoreMissingObject(t *testing.T) {
	var buf bytes.Buffer
	stasher := &Stasher{
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("echo")
		},
	}
	stasher.helper.outputWriter = &buf

	stasher.Stash([]string{"store"})

	expected := "Error: stash object required\n"
	if got := buf.String(); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}
