package cmd

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"testing"
)

func TestBrancher_Branch_Current(t *testing.T) {
	brancher := &Brancher{
		GetCurrentBranch: func() (string, error) {
			return "feature/test", nil
		},
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	brancher.Branch([]string{"current"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stdout = oldStdout

	output := buf.String()
	if output == "" || output[:7] != "feature" {
		t.Errorf("currentサブコマンドの出力が想定と異なります: %s", output)
	}
}

func TestBrancher_Branch_Checkout(t *testing.T) {
	called := false
	brancher := &Brancher{
		ListLocalBranches: func() ([]string, error) {
			return []string{"main", "feature/test"}, nil
		},
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && arg[0] == "checkout" && arg[1] == "feature/test" {
				called = true
			}
			return exec.Command("echo")
		},
		inputReader: bufio.NewReader(bytes.NewBufferString("2\n")),
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	brancher.Branch([]string{"checkout"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stdout = oldStdout
	if !called {
		t.Errorf("checkoutで正しいブランチがcheckoutされていません")
	}
}
