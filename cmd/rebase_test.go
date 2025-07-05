package cmd

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestRebaser_RebaseInteractive_SelectValid(t *testing.T) {
	var buf bytes.Buffer
	commandCount := 0
	r := &Rebaser{
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, args ...string) *exec.Cmd {
			commandCount++
			if len(args) > 0 {
				switch args[0] {
				case "rev-parse":
					if strings.Contains(args[len(args)-1], "@{upstream}") {
						return exec.Command("echo", "origin/main")
					}
					return exec.Command("echo", "feature/test")
				case "log":
					return exec.Command("echo", "a1 first\nb2 second\nc3 third")
				}
			}
			cmd := exec.Command("echo")
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd
		},
		inputReader: bufio.NewReader(strings.NewReader("2\n")),
	}
	r.helper.outputWriter = &buf

	r.RebaseInteractive()

	output := buf.String()
	if !strings.Contains(output, "Current branch: feature/test") {
		t.Errorf("expected branch name, got: %s", output)
	}
	if !strings.Contains(output, "Select number of commits to rebase") {
		t.Errorf("expected prompt, got: %s", output)
	}
	if !strings.Contains(output, "Rebase successful") {
		t.Errorf("expected success message, got: %s", output)
	}
}

func TestRebaser_RebaseInteractive_BranchError(t *testing.T) {
	var buf bytes.Buffer
	r := &Rebaser{
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("false")
		},
		inputReader: bufio.NewReader(strings.NewReader("1\n")),
	}
	r.helper.outputWriter = &buf

	r.RebaseInteractive()

	output := buf.String()
	if !strings.Contains(output, "Error: failed to get current branch") {
		t.Errorf("expected branch error message, got: %s", output)
	}
}

func TestRebaser_RebaseInteractive_Cancel(t *testing.T) {
	var buf bytes.Buffer
	r := &Rebaser{
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, args ...string) *exec.Cmd {
			if len(args) > 0 {
				switch args[0] {
				case "rev-parse":
					if strings.Contains(args[len(args)-1], "@{upstream}") {
						return exec.Command("echo", "origin/main")
					}
					return exec.Command("echo", "feature/test")
				case "log":
					return exec.Command("echo", "a1 first\nb2 second\nc3 third")
				}
			}
			return exec.Command("echo")
		},
		inputReader: bufio.NewReader(strings.NewReader("\n")),
	}
	r.helper.outputWriter = &buf

	r.RebaseInteractive()

	output := buf.String()
	if !strings.Contains(output, "Error: operation cancelled") {
		t.Errorf("expected cancellation message, got: %s", output)
	}
}

func TestRebaser_RebaseInteractive_InvalidNumber(t *testing.T) {
	var buf bytes.Buffer
	r := &Rebaser{
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, args ...string) *exec.Cmd {
			if len(args) > 0 {
				switch args[0] {
				case "rev-parse":
					if strings.Contains(args[len(args)-1], "@{upstream}") {
						return exec.Command("echo", "origin/main")
					}
					return exec.Command("echo", "feature/test")
				case "log":
					return exec.Command("echo", "a1 first\nb2 second\nc3 third")
				}
			}
			return exec.Command("echo")
		},
		inputReader: bufio.NewReader(strings.NewReader("abc\n")),
	}
	r.helper.outputWriter = &buf

	r.RebaseInteractive()

	output := buf.String()
	if !strings.Contains(output, "Error: invalid number") {
		t.Errorf("expected invalid number message, got: %s", output)
	}
}

func TestRebaser_RebaseInteractive_NoHistory(t *testing.T) {
	var buf bytes.Buffer
	r := &Rebaser{
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, args ...string) *exec.Cmd {
			if len(args) > 0 {
				switch args[0] {
				case "rev-parse":
					if strings.Contains(args[len(args)-1], "@{upstream}") {
						return exec.Command("echo", "origin/main")
					}
					return exec.Command("echo", "feature/test")
				case "log":
					return exec.Command("echo", "")
				}
			}
			return exec.Command("echo")
		},
		inputReader: bufio.NewReader(strings.NewReader("1\n")),
	}
	r.helper.outputWriter = &buf

	r.RebaseInteractive()

	output := buf.String()
	if !strings.Contains(output, "Error: no commit history found") {
		t.Errorf("expected no history message, got: %s", output)
	}
}

func TestRebaser_RebaseInteractive_LogError(t *testing.T) {
	var buf bytes.Buffer
	commandCount := 0
	r := &Rebaser{
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, args ...string) *exec.Cmd {
			commandCount++
			if len(args) > 0 {
				switch args[0] {
				case "rev-parse":
					if strings.Contains(args[len(args)-1], "@{upstream}") {
						return exec.Command("echo", "origin/main")
					}
					return exec.Command("echo", "feature/test")
				}
			}
			return exec.Command("false")
		},
		inputReader: bufio.NewReader(strings.NewReader("1\n")),
	}
	r.helper.outputWriter = &buf

	r.RebaseInteractive()

	output := buf.String()
	if !strings.Contains(output, "Error: failed to get git log") {
		t.Errorf("expected log error message, got: %s", output)
	}
}

func TestRebaser_RebaseInteractive_RebaseError(t *testing.T) {
	var buf bytes.Buffer
	commandCount := 0
	r := &Rebaser{
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, args ...string) *exec.Cmd {
			commandCount++
			if len(args) > 0 {
				switch args[0] {
				case "rev-parse":
					if strings.Contains(args[len(args)-1], "@{upstream}") {
						return exec.Command("echo", "origin/main")
					}
					return exec.Command("echo", "feature/test")
				case "log":
					return exec.Command("echo", "a1 first")
				case "rebase":
					return exec.Command("false")
				}
			}
			return exec.Command("echo")
		},
		inputReader: bufio.NewReader(strings.NewReader("1\n")),
	}
	r.helper.outputWriter = &buf

	r.RebaseInteractive()

	output := buf.String()
	if !strings.Contains(output, "Error: rebase failed") {
		t.Errorf("expected rebase error message, got: %s", output)
	}
}

func TestRebaser_Rebase(t *testing.T) {
	cases := []struct {
		name           string
		args           []string
		expectedOutput string
		mockInput      string
	}{
		{
			name:           "interactive rebase",
			args:           []string{"interactive"},
			expectedOutput: "Current branch:",
			mockInput:      "1\n",
		},
		{
			name:           "no args",
			args:           []string{},
			expectedOutput: "Usage: ggc rebase",
		},
		{
			name:           "invalid command",
			args:           []string{"invalid"},
			expectedOutput: "Usage: ggc rebase",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			r := &Rebaser{
				outputWriter: &buf,
				helper:       NewHelper(),
				execCommand: func(_ string, args ...string) *exec.Cmd {
					if len(args) > 0 {
						switch args[0] {
						case "rev-parse":
							if strings.Contains(args[len(args)-1], "@{upstream}") {
								return exec.Command("echo", "origin/main")
							}
							return exec.Command("echo", "feature/test")
						case "log":
							return exec.Command("echo", "abc123 commit message")
						}
					}
					return exec.Command("echo", "rebase output")
				},
				inputReader: bufio.NewReader(strings.NewReader(tc.mockInput)),
			}
			r.helper.outputWriter = &buf

			r.Rebase(tc.args)

			output := buf.String()
			if !strings.Contains(output, tc.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tc.expectedOutput, output)
			}
		})
	}
}

func TestRebaser_Rebase_InteractiveCancel(t *testing.T) {
	var buf bytes.Buffer
	r := &Rebaser{
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, args ...string) *exec.Cmd {
			if len(args) > 0 {
				switch args[0] {
				case "rev-parse":
					if strings.Contains(args[len(args)-1], "@{upstream}") {
						return exec.Command("echo", "origin/main")
					}
					return exec.Command("echo", "feature/test")
				case "log":
					return exec.Command("echo", "")
				}
			}
			return exec.Command("echo")
		},
		inputReader: bufio.NewReader(strings.NewReader("\n")),
	}
	r.helper.outputWriter = &buf

	r.Rebase([]string{"interactive"})

	output := buf.String()
	if !strings.Contains(output, "Error: no commit history found") {
		t.Errorf("expected no history message, got %q", output)
	}
}

func TestRebaser_Rebase_Error(t *testing.T) {
	var buf bytes.Buffer
	r := &Rebaser{
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("false") // Command fails
		},
		inputReader: bufio.NewReader(strings.NewReader("y\n")),
	}
	r.helper.outputWriter = &buf

	r.Rebase([]string{"interactive"})

	output := buf.String()
	if !strings.Contains(output, "Error:") {
		t.Errorf("expected error message, got %q", output)
	}
}
