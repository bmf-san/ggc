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
