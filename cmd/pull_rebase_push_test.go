package cmd

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestPullRebasePusher_PullRebasePush_Success(t *testing.T) {
	calls := []string{}
	prp := &PullRebasePusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			calls = append(calls, name+" "+strings.Join(arg, " "))
			if name == "git" && arg[0] == "rev-parse" {
				return exec.Command("echo", "main")
			}
			return exec.Command("echo")
		},
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	prp.outputWriter = &buf
	prp.PullRebasePush()
	output := buf.String()
	if !strings.Contains(output, "pull→rebase→push completed") {
		t.Errorf("正常系で完了メッセージが出力されていません: %s", output)
	}
	if len(calls) < 4 {
		t.Errorf("コマンド呼び出し数が想定より少ない: %v", calls)
	}
}

func TestPullRebasePusher_PullRebasePush_BranchError(t *testing.T) {
	prp := &PullRebasePusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			return exec.Command("false")
		},
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	prp.outputWriter = &buf
	prp.PullRebasePush()
	output := buf.String()
	if !strings.Contains(output, "Error: Failed to get branch name") {
		t.Errorf("branch取得失敗時のエラー出力が想定と異なります: %s", output)
	}
}

func TestPullRebasePusher_PullRebasePush_PullError(t *testing.T) {
	step := 0
	prp := &PullRebasePusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if step == 0 {
				step++
				return exec.Command("echo", "main") // branch
			}
			if step == 1 {
				step++
				return exec.Command("false") // pull
			}
			return exec.Command("echo")
		},
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	prp.outputWriter = &buf
	prp.PullRebasePush()
	output := buf.String()
	if !strings.Contains(output, "Error: Failed to git pull") {
		t.Errorf("pull失敗時のエラー出力が想定と異なります: %s", output)
	}
}

func TestPullRebasePusher_PullRebasePush_RebaseError(t *testing.T) {
	step := 0
	prp := &PullRebasePusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if step == 0 {
				step++
				return exec.Command("echo", "main") // branch
			}
			if step == 1 {
				step++
				return exec.Command("echo") // pull
			}
			if step == 2 {
				step++
				return exec.Command("false") // rebase
			}
			return exec.Command("echo")
		},
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	prp.outputWriter = &buf
	prp.PullRebasePush()
	output := buf.String()
	if !strings.Contains(output, "Error: Failed to git rebase") {
		t.Errorf("rebase失敗時のエラー出力が想定と異なります: %s", output)
	}
}

func TestPullRebasePusher_PullRebasePush_PushError(t *testing.T) {
	step := 0
	prp := &PullRebasePusher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if step == 0 {
				step++
				return exec.Command("echo", "main") // branch
			}
			if step == 1 {
				step++
				return exec.Command("echo") // pull
			}
			if step == 2 {
				step++
				return exec.Command("echo") // rebase
			}
			if step == 3 {
				step++
				return exec.Command("false") // push
			}
			return exec.Command("echo")
		},
		outputWriter: &bytes.Buffer{},
	}
	var buf bytes.Buffer
	prp.outputWriter = &buf
	prp.PullRebasePush()
	output := buf.String()
	if !strings.Contains(output, "Error: Failed to git push") {
		t.Errorf("push失敗時のエラー出力が想定と異なります: %s", output)
	}
}

func TestNewPullRebasePusher(t *testing.T) {
	prp := NewPullRebasePusher()
	if prp.execCommand == nil {
		t.Error("execCommandがnilです")
	}
	if prp.outputWriter == nil {
		t.Error("outputWriterがnilです")
	}
}
