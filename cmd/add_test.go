package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestAdder_Add_NoArgs_PrintsUsage(t *testing.T) {
	adder := NewAdder()
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	adder.Add([]string{})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stdout = oldStdout

	output := buf.String()
	if output == "" || output[:5] != "Usage" {
		t.Errorf("Usage not output: %s", output)
	}
}

func TestAdder_Add_GitAddCalled(t *testing.T) {
	called := false
	adder := &Adder{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			called = true
			return exec.Command("echo")
		},
	}
	adder.Add([]string{"hoge.txt"})
	if !called {
		t.Error("execCommand was not called")
	}
}

func TestAdder_Add_GitAddArgs(t *testing.T) {
	var gotName string
	var gotArgs []string
	adder := &Adder{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			gotName = name
			gotArgs = arg
			return exec.Command("echo")
		},
	}
	adder.Add([]string{"foo.txt", "bar.txt"})
	if gotName != "git" {
		t.Errorf("Command name differs from expected: got=%s", gotName)
	}
	wantArgs := []string{"add", "foo.txt", "bar.txt"}
	for i, a := range wantArgs {
		if i >= len(gotArgs) || gotArgs[i] != a {
			t.Errorf("Arguments differ from expected: want=%v, got=%v", wantArgs, gotArgs)
			break
		}
	}
}

func TestAdder_Add_RunError_PrintsError(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	adder := &Adder{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			cmd := exec.Command("false") // command that always returns error
			return cmd
		},
	}
	adder.Add([]string{"foo.txt"})
	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stdout = oldStdout
	output := buf.String()
	if output == "" || output[:5] != "error" {
		t.Errorf("Error output not generated: %s", output)
	}
}

func TestAdder_Add_POption_CallsGitAddP(t *testing.T) {
	called := false
	var gotName string
	var gotArgs []string
	adder := &Adder{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			called = true
			gotName = name
			gotArgs = arg
			return exec.Command("echo")
		},
	}
	adder.Add([]string{"-p"})
	if !called {
		t.Error("execCommand not called with -p option")
	}
	if gotName != "git" || len(gotArgs) != 2 || gotArgs[0] != "add" || gotArgs[1] != "-p" {
		t.Errorf("Command/arguments differ from expected for -p option: name=%s, args=%v", gotName, gotArgs)
	}
}

func TestAdder_Add_POption_Error(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	adder := &Adder{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			cmd := exec.Command("false") // command that always returns error
			return cmd
		},
	}
	adder.Add([]string{"-p"})
	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stdout = oldStdout
	output := buf.String()
	if output == "" || output[:5] != "error" {
		t.Errorf("Error output not generated with -p option: %s", output)
	}
}

func TestAdder_Add_Interactive(t *testing.T) {
	var buf bytes.Buffer
	adder := &Adder{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			cmd := exec.Command("echo", "interactive add")
			return cmd
		},
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	adder.Add([]string{"-p"})

	_ = w.Close()
	os.Stdout = oldStdout

	_, _ = buf.ReadFrom(r)
	output := buf.String()
	if !strings.Contains(output, "interactive add") {
		t.Errorf("expected interactive add output, got %q", output)
	}
}

func TestAdder_Add(t *testing.T) {
	cases := []struct {
		name        string
		args        []string
		expectedCmd string
		expectError bool
	}{
		{
			name:        "add all files",
			args:        []string{"."},
			expectedCmd: "git add .",
			expectError: false,
		},
		{
			name:        "add specific file",
			args:        []string{"file.txt"},
			expectedCmd: "git add file.txt",
			expectError: false,
		},
		{
			name:        "add multiple files",
			args:        []string{"file1.txt", "file2.txt"},
			expectedCmd: "git add file1.txt file2.txt",
			expectError: false,
		},
		{
			name:        "no args",
			args:        []string{},
			expectedCmd: "",
			expectError: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var actualCmd string
			a := &Adder{
				execCommand: func(name string, args ...string) *exec.Cmd {
					if tc.expectedCmd != "" {
						actualCmd = strings.Join(append([]string{name}, args...), " ")
						if actualCmd != tc.expectedCmd {
							t.Errorf("expected command %q, got %q", tc.expectedCmd, actualCmd)
						}
					}
					return exec.Command("echo")
				},
			}

			// Capture stdout for no args case
			if len(tc.args) == 0 {
				oldStdout := os.Stdout
				r, w, _ := os.Pipe()
				os.Stdout = w

				a.Add(tc.args)

				_ = w.Close()
				os.Stdout = oldStdout

				var buf bytes.Buffer
				_, _ = buf.ReadFrom(r)
				output := buf.String()
				if !strings.Contains(output, "Usage:") {
					t.Errorf("expected usage message, got %q", output)
				}
			} else {
				a.Add(tc.args)
			}
		})
	}
}

func TestAdder_Add_Error(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	a := &Adder{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("false") // Command fails
		},
	}

	a.Add([]string{"file.txt"})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()
	if !strings.Contains(output, "error:") {
		t.Errorf("expected error message, got %q", output)
	}
}
