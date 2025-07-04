package cmd

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestCommitPusher_CommitPushInteractive_AllSuccess(t *testing.T) {
	calls := []string{}
	cp := &CommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			calls = append(calls, name+" "+strings.Join(arg, " "))
			if name == "git" && arg[0] == "status" {
				return exec.Command("echo", " M foo.txt\n M bar.txt")
			}
			if name == "git" && arg[0] == "rev-parse" {
				return exec.Command("echo", "main")
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("all\ntestmsg\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	cp.outputWriter = &buf
	cp.CommitPushInteractive()
	output := buf.String()
	if !strings.Contains(output, "Done!") {
		t.Errorf("success message should be displayed, but not: %s", output)
	}
	if len(calls) < 5 {
		t.Errorf("number of command calls is less than expected: %v", calls)
	}
}

func TestCommitPusher_CommitPushInteractive_Cancel(t *testing.T) {
	cp := &CommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && arg[0] == "status" {
				return exec.Command("echo", " M foo.txt\n M bar.txt")
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	cp.outputWriter = &buf
	cp.CommitPushInteractive()
	output := buf.String()
	if !strings.Contains(output, "Cancelled.") {
		t.Errorf("cancel message should be displayed, but not: %s", output)
	}
}

func TestCommitPusher_CommitPushInteractive_AddError(t *testing.T) {
	cp := &CommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && arg[0] == "status" {
				return exec.Command("echo", " M foo.txt\n M bar.txt")
			}
			if name == "git" && arg[0] == "add" {
				return exec.Command("false")
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("all\ntestmsg\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	cp.outputWriter = &buf
	cp.CommitPushInteractive()
	output := buf.String()
	if !strings.Contains(output, "Error: failed to add files") {
		t.Errorf("error message on add failure should be displayed, but not: %s", output)
	}
}

func TestCommitPusher_CommitPushInteractive_CommitError(t *testing.T) {
	cp := &CommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && arg[0] == "status" {
				return exec.Command("echo", " M foo.txt\n M bar.txt")
			}
			if name == "git" && arg[0] == "add" {
				return exec.Command("echo")
			}
			if name == "git" && arg[0] == "commit" {
				return exec.Command("false")
			}
			if name == "git" && arg[0] == "rev-parse" {
				return exec.Command("echo", "main")
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("all\ntestmsg\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	cp.outputWriter = &buf
	cp.CommitPushInteractive()
	output := buf.String()
	if !strings.Contains(output, "Error: failed to commit") {
		t.Errorf("error message on commit failure should be displayed, but not: %s", output)
	}
}

func TestCommitPusher_CommitPushInteractive_BranchError(t *testing.T) {
	cp := &CommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && arg[0] == "status" {
				return exec.Command("echo", " M foo.txt\n M bar.txt")
			}
			if name == "git" && arg[0] == "add" {
				return exec.Command("echo")
			}
			if name == "git" && arg[0] == "commit" {
				return exec.Command("echo")
			}
			if name == "git" && arg[0] == "rev-parse" {
				return exec.Command("false")
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("all\ntestmsg\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	cp.outputWriter = &buf
	cp.CommitPushInteractive()
	output := buf.String()
	if !strings.Contains(output, "Error: failed to get branch name") {
		t.Errorf("error message on branch name fetch failure should be displayed, but not: %s", output)
	}
}

func TestCommitPusher_CommitPushInteractive_PushError(t *testing.T) {
	cp := &CommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && arg[0] == "status" {
				return exec.Command("echo", " M foo.txt\n M bar.txt")
			}
			if name == "git" && arg[0] == "add" {
				return exec.Command("echo")
			}
			if name == "git" && arg[0] == "commit" {
				return exec.Command("echo")
			}
			if name == "git" && arg[0] == "rev-parse" {
				return exec.Command("echo", "main")
			}
			if name == "git" && arg[0] == "push" {
				return exec.Command("false")
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("all\ntestmsg\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	cp.outputWriter = &buf
	cp.CommitPushInteractive()
	output := buf.String()
	if !strings.Contains(output, "Error: failed to push") {
		t.Errorf("error message on push failure should be displayed, but not: %s", output)
	}
}

func TestCommitPusher_CommitPushInteractive_StatusError(t *testing.T) {
	cp := &CommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && arg[0] == "status" {
				return exec.Command("false")
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	cp.outputWriter = &buf
	cp.CommitPushInteractive()
	output := buf.String()
	if !strings.Contains(output, "Error: failed to get git status") {
		t.Errorf("status error message should be displayed, but not: %s", output)
	}
}

func TestCommitPusher_CommitPushInteractive_NoChanges(t *testing.T) {
	cp := &CommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && arg[0] == "status" {
				return exec.Command("echo", "")
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	cp.outputWriter = &buf
	cp.CommitPushInteractive()
	output := buf.String()
	if !strings.Contains(output, "No changed files.") {
		t.Errorf("no changes message should be displayed, but not: %s", output)
	}
}

func TestCommitPusher_CommitPushInteractive_NoFilesToStage(t *testing.T) {
	cp := &CommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && arg[0] == "status" {
				return exec.Command("echo", "?? \n")
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	cp.outputWriter = &buf
	cp.CommitPushInteractive()
	output := buf.String()
	if !strings.Contains(output, "No files to stage.") {
		t.Errorf("no files to stage message should be displayed, but not: %s", output)
	}
}

func TestCommitPusher_CommitPushInteractive_SelectNone(t *testing.T) {
	cp := &CommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && arg[0] == "status" {
				return exec.Command("echo", " M foo.txt\n M bar.txt")
			}
			if name == "git" && arg[0] == "add" {
				return exec.Command("echo")
			}
			if name == "git" && arg[0] == "commit" {
				return exec.Command("echo")
			}
			if name == "git" && arg[0] == "rev-parse" {
				return exec.Command("echo", "main")
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("none\nall\ntestmsg\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	cp.outputWriter = &buf
	cp.CommitPushInteractive()
	output := buf.String()
	if !strings.Contains(output, "Done!") {
		t.Errorf("success message should be displayed, but not: %s", output)
	}
}

func TestCommitPusher_CommitPushInteractive_SelectSpecificFiles(t *testing.T) {
	cp := &CommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && arg[0] == "status" {
				return exec.Command("echo", " M foo.txt\n M bar.txt\n M baz.txt")
			}
			if name == "git" && arg[0] == "add" {
				return exec.Command("echo")
			}
			if name == "git" && arg[0] == "commit" {
				return exec.Command("echo")
			}
			if name == "git" && arg[0] == "rev-parse" {
				return exec.Command("echo", "main")
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("1 3\ny\ntestmsg\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	cp.outputWriter = &buf
	cp.CommitPushInteractive()
	output := buf.String()
	if !strings.Contains(output, "Done!") {
		t.Errorf("success message should be displayed, but not: %s", output)
	}
	if !strings.Contains(output, "Selected files: [foo.txt baz.txt]") {
		t.Errorf("selected files should be displayed, but not: %s", output)
	}
}

func TestCommitPusher_CommitPushInteractive_InvalidNumber(t *testing.T) {
	cp := &CommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && arg[0] == "status" {
				return exec.Command("echo", " M foo.txt\n M bar.txt")
			}
			if name == "git" && arg[0] == "add" {
				return exec.Command("echo")
			}
			if name == "git" && arg[0] == "commit" {
				return exec.Command("echo")
			}
			if name == "git" && arg[0] == "rev-parse" {
				return exec.Command("echo", "main")
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("abc\nall\ntestmsg\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	cp.outputWriter = &buf
	cp.CommitPushInteractive()
	output := buf.String()
	if !strings.Contains(output, "Invalid number: abc") {
		t.Errorf("invalid number message should be displayed, but not: %s", output)
	}
	if !strings.Contains(output, "Done!") {
		t.Errorf("success message should be displayed after recovery, but not: %s", output)
	}
}

func TestCommitPusher_CommitPushInteractive_OutOfRangeNumber(t *testing.T) {
	cp := &CommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && arg[0] == "status" {
				return exec.Command("echo", " M foo.txt\n M bar.txt")
			}
			if name == "git" && arg[0] == "add" {
				return exec.Command("echo")
			}
			if name == "git" && arg[0] == "commit" {
				return exec.Command("echo")
			}
			if name == "git" && arg[0] == "rev-parse" {
				return exec.Command("echo", "main")
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("5\nall\ntestmsg\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	cp.outputWriter = &buf
	cp.CommitPushInteractive()
	output := buf.String()
	if !strings.Contains(output, "Invalid number: 5") {
		t.Errorf("invalid number message should be displayed, but not: %s", output)
	}
}

func TestCommitPusher_CommitPushInteractive_NothingSelected(t *testing.T) {
	cp := &CommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && arg[0] == "status" {
				return exec.Command("echo", " M foo.txt\n M bar.txt")
			}
			if name == "git" && arg[0] == "add" {
				return exec.Command("echo")
			}
			if name == "git" && arg[0] == "commit" {
				return exec.Command("echo")
			}
			if name == "git" && arg[0] == "rev-parse" {
				return exec.Command("echo", "main")
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("0\nall\ntestmsg\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	cp.outputWriter = &buf
	cp.CommitPushInteractive()
	output := buf.String()
	if !strings.Contains(output, "Invalid number: 0") {
		t.Errorf("invalid number message should be displayed, but not: %s", output)
	}
}

func TestCommitPusher_CommitPushInteractive_RejectFiles(t *testing.T) {
	cp := &CommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && arg[0] == "status" {
				return exec.Command("echo", " M foo.txt\n M bar.txt")
			}
			if name == "git" && arg[0] == "add" {
				return exec.Command("echo")
			}
			if name == "git" && arg[0] == "commit" {
				return exec.Command("echo")
			}
			if name == "git" && arg[0] == "rev-parse" {
				return exec.Command("echo", "main")
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("1 2\nn\nall\ntestmsg\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	cp.outputWriter = &buf
	cp.CommitPushInteractive()
	output := buf.String()
	if !strings.Contains(output, "Done!") {
		t.Errorf("success message should be displayed, but not: %s", output)
	}
}

func TestCommitPusher_CommitPushInteractive_EmptyCommitMessage(t *testing.T) {
	cp := &CommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && arg[0] == "status" {
				return exec.Command("echo", " M foo.txt\n M bar.txt")
			}
			if name == "git" && arg[0] == "add" {
				return exec.Command("echo")
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("all\n\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	cp.outputWriter = &buf
	cp.CommitPushInteractive()
	output := buf.String()
	if !strings.Contains(output, "Cancelled.") {
		t.Errorf("cancelled message should be displayed for empty commit message, but not: %s", output)
	}
}

func TestCommitPusher_CommitPushInteractive_ShortStatusLine(t *testing.T) {
	cp := &CommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && arg[0] == "status" {
				return exec.Command("echo", "??\n M foo.txt")
			}
			if name == "git" && arg[0] == "add" {
				return exec.Command("echo")
			}
			if name == "git" && arg[0] == "commit" {
				return exec.Command("echo")
			}
			if name == "git" && arg[0] == "rev-parse" {
				return exec.Command("echo", "main")
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("all\ntestmsg\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	cp.outputWriter = &buf
	cp.CommitPushInteractive()
	output := buf.String()
	if !strings.Contains(output, "Done!") {
		t.Errorf("success message should be displayed, but not: %s", output)
	}
}

func TestCommitPusher_CommitPushInteractive_EmptySelection(t *testing.T) {
	cp := &CommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && arg[0] == "status" {
				return exec.Command("echo", " M foo.txt\n M bar.txt")
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	cp.outputWriter = &buf
	cp.CommitPushInteractive()
	output := buf.String()
	if !strings.Contains(output, "Cancelled.") {
		t.Errorf("cancelled message should be displayed for empty input, but not: %s", output)
	}
}

func TestCommitPusher_CommitPushInteractive_Success(t *testing.T) {
	calls := []string{}
	cp := &CommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			calls = append(calls, name+" "+strings.Join(arg, " "))
			if name == "git" && arg[0] == "status" {
				return exec.Command("echo", " M test.go")
			}
			if name == "git" && arg[0] == "rev-parse" {
				return exec.Command("echo", "main")
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("all\ntest message\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	cp.outputWriter = &buf

	cp.CommitPushInteractive()

	output := buf.String()
	if !strings.Contains(output, "Done!") {
		t.Errorf("success message should be displayed, but not: %s", output)
	}
	if len(calls) < 3 {
		t.Errorf("number of command calls is less than expected: %v", calls)
	}
}

func TestCommitPusher_CommitPushInteractive_Cancel_EmptyMessage(t *testing.T) {
	cp := &CommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && arg[0] == "status" {
				return exec.Command("echo", " M test.go")
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("all\n\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	cp.outputWriter = &buf

	cp.CommitPushInteractive()

	output := buf.String()
	if !strings.Contains(output, "Cancelled.") {
		t.Errorf("cancel message should be displayed, but not: %s", output)
	}
}

func TestCommitPusher_CommitPushInteractive_CommitFailure(t *testing.T) {
	cp := &CommitPusher{
		execCommand: func(name string, args ...string) *exec.Cmd {
			if name == "git" && args[0] == "status" {
				return exec.Command("echo", " M test.go")
			}
			if name == "git" && args[0] == "commit" {
				return exec.Command("false")
			}
			return exec.Command("echo", "main")
		},
		inputReader:  bufio.NewReader(strings.NewReader("all\ntest message\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	cp.outputWriter = &buf

	cp.CommitPushInteractive()

	output := buf.String()
	if !strings.Contains(output, "Error: failed to commit") {
		t.Errorf("error message on commit failure should be displayed, but not: %s", output)
	}
}

func TestCommitPusher_CommitPushInteractive_PushFailure(t *testing.T) {
	cp := &CommitPusher{
		execCommand: func(name string, args ...string) *exec.Cmd {
			if name == "git" && args[0] == "status" {
				return exec.Command("echo", " M test.go")
			}
			if name == "git" && args[0] == "push" {
				return exec.Command("false")
			}
			if name == "git" && args[0] == "rev-parse" {
				return exec.Command("echo", "main")
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("all\ntest message\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	cp.outputWriter = &buf

	cp.CommitPushInteractive()

	output := buf.String()
	if !strings.Contains(output, "Error: failed to push") {
		t.Errorf("error message on push failure should be displayed, but not: %s", output)
	}
}
