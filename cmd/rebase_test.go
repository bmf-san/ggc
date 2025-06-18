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
	var gotArgs [][]string
	rebaser := &Rebaser{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			gotArgs = append(gotArgs, append([]string{name}, arg...))
			if name == "git" && arg[0] == "log" {
				return exec.Command("echo", "a1 first\nb2 second\nc3 third")
			}
			return exec.Command("echo")
		},
		inputReader: bufioReaderWithInput("2\n"),
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rebaser.RebaseInteractive()

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("buf.ReadFrom failed: %v", err)
	}
	os.Stdout = oldStdout

	output := buf.String()
	if !strings.Contains(output, "Where do you want to rebase up to?") {
		t.Errorf("プロンプトが出力されていません: %s", output)
	}
	if len(gotArgs) < 2 || gotArgs[1][0] != "git" || gotArgs[1][1] != "rebase" {
		t.Errorf("rebaseコマンドが呼ばれていません: %+v", gotArgs)
	}
}

func TestRebaser_RebaseInteractive_Cancel(t *testing.T) {
	rebaser := &Rebaser{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && arg[0] == "log" {
				return exec.Command("echo", "a1 first\nb2 second\nc3 third")
			}
			return exec.Command("echo")
		},
		inputReader: bufioReaderWithInput("\n"),
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rebaser.RebaseInteractive()

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("buf.ReadFrom failed: %v", err)
	}
	os.Stdout = oldStdout

	output := buf.String()
	if !strings.Contains(output, "Cancelled") {
		t.Errorf("キャンセル時の出力が想定と異なります: %s", output)
	}
}

func TestRebaser_RebaseInteractive_InvalidNumber(t *testing.T) {
	rebaser := &Rebaser{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && arg[0] == "log" {
				return exec.Command("echo", "a1 first\nb2 second\nc3 third")
			}
			return exec.Command("echo")
		},
		inputReader: bufioReaderWithInput("abc\n"),
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rebaser.RebaseInteractive()

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("buf.ReadFrom failed: %v", err)
	}
	os.Stdout = oldStdout

	output := buf.String()
	if !strings.Contains(output, "Invalid number") {
		t.Errorf("不正な番号入力時の出力が想定と異なります: %s", output)
	}
}

func TestRebaser_RebaseInteractive_NoHistory(t *testing.T) {
	rebaser := &Rebaser{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			return exec.Command("echo", "")
		},
		inputReader: bufioReaderWithInput("1\n"),
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rebaser.RebaseInteractive()

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("buf.ReadFrom failed: %v", err)
	}
	os.Stdout = oldStdout

	output := buf.String()
	if !strings.Contains(output, "No commit history found") {
		t.Errorf("履歴なし時の出力が想定と異なります: %s", output)
	}
}

func TestRebaser_RebaseInteractive_LogError(t *testing.T) {
	rebaser := &Rebaser{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			return exec.Command("false")
		},
		inputReader: bufioReaderWithInput("1\n"),
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rebaser.RebaseInteractive()

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("buf.ReadFrom failed: %v", err)
	}
	os.Stdout = oldStdout

	output := buf.String()
	if !strings.Contains(output, "error: failed to get git log") {
		t.Errorf("git logエラー時の出力が想定と異なります: %s", output)
	}
}

// テスト用: 任意の入力を返すbufio.Readerを生成
func bufioReaderWithInput(s string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(s))
}
