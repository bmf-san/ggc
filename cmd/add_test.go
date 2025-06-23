package cmd

import (
	"bytes"
	"os"
	"os/exec"
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
		t.Errorf("Usageが出力されていません: %s", output)
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
		t.Error("execCommandが呼ばれていません")
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
		t.Errorf("コマンド名が想定と異なります: got=%s", gotName)
	}
	wantArgs := []string{"add", "foo.txt", "bar.txt"}
	for i, a := range wantArgs {
		if i >= len(gotArgs) || gotArgs[i] != a {
			t.Errorf("引数が想定と異なります: want=%v, got=%v", wantArgs, gotArgs)
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
			cmd := exec.Command("false") // 常にエラーを返すコマンド
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
		t.Errorf("エラー出力がされていません: %s", output)
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
		t.Error("-pオプションでexecCommandが呼ばれていません")
	}
	if gotName != "git" || len(gotArgs) != 2 || gotArgs[0] != "add" || gotArgs[1] != "-p" {
		t.Errorf("-pオプション時のコマンド・引数が想定と異なります: name=%s, args=%v", gotName, gotArgs)
	}
}
