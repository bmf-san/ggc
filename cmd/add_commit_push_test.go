package cmd

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestAddCommitPusher_AddCommitPush_Success(t *testing.T) {
	calls := []string{}
	acp := &AddCommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			calls = append(calls, name+" "+strings.Join(arg, " "))
			if name == "git" && arg[0] == "rev-parse" {
				return exec.Command("echo", "main")
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("test message\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	acp.outputWriter = &buf
	acp.AddCommitPush()
	output := buf.String()
	if !strings.Contains(output, "add→commit→push done") {
		t.Errorf("success message should be displayed, but not: %s", output)
	}
	if len(calls) < 4 {
		t.Errorf("number of command calls is less than expected: %v", calls)
	}
}

func TestAddCommitPusher_AddCommitPush_Cancel(t *testing.T) {
	acp := &AddCommitPusher{
		execCommand:  exec.Command,
		inputReader:  bufio.NewReader(strings.NewReader("\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	acp.outputWriter = &buf
	acp.AddCommitPush()
	output := buf.String()
	if !strings.Contains(output, "Cancelled.") {
		t.Errorf("cancel message should be displayed, but not: %s", output)
	}
}

func TestAddCommitPusher_AddCommitPush_AddError(t *testing.T) {
	acp := &AddCommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			return exec.Command("false")
		},
		inputReader:  bufio.NewReader(strings.NewReader("test\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	acp.outputWriter = &buf
	acp.AddCommitPush()
	output := buf.String()
	if !strings.Contains(output, "Error: failed to add all files") {
		t.Errorf("error message on add failure should be displayed, but not: %s", output)
	}
}

func TestAddCommitPusher_AddCommitPush_CommitError(t *testing.T) {
	step := 0
	acp := &AddCommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if step == 0 {
				step++
				return exec.Command("echo") // add
			}
			if step == 1 {
				step++
				return exec.Command("false") // commit
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("msg\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	acp.outputWriter = &buf
	acp.AddCommitPush()
	output := buf.String()
	if !strings.Contains(output, "Error: failed to commit") {
		t.Errorf("error message on commit failure should be displayed, but not: %s", output)
	}
}

func TestAddCommitPusher_AddCommitPush_BranchError(t *testing.T) {
	step := 0
	acp := &AddCommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if step == 0 {
				step++
				return exec.Command("echo") // add
			}
			if step == 1 {
				step++
				return exec.Command("echo") // commit
			}
			if step == 2 {
				step++
				return exec.Command("false") // branch
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("msg\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	acp.outputWriter = &buf
	acp.AddCommitPush()
	output := buf.String()
	if !strings.Contains(output, "Error: failed to get branch name") {
		t.Errorf("error message on branch name fetch failure should be displayed, but not: %s", output)
	}
}

func TestAddCommitPusher_AddCommitPush_PushError(t *testing.T) {
	step := 0
	acp := &AddCommitPusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if step == 0 {
				step++
				return exec.Command("echo") // add
			}
			if step == 1 {
				step++
				return exec.Command("echo") // commit
			}
			if step == 2 {
				step++
				return exec.Command("echo", "main") // branch
			}
			if step == 3 {
				step++
				return exec.Command("false") // push
			}
			return exec.Command("echo")
		},
		inputReader:  bufio.NewReader(strings.NewReader("msg\n")),
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	acp.outputWriter = &buf
	acp.AddCommitPush()
	output := buf.String()
	if !strings.Contains(output, "Error: failed to push") {
		t.Errorf("error message on push failure should be displayed, but not: %s", output)
	}
}
