package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
)

func TestRemoteer_Remote_List(t *testing.T) {
	called := false
	remoteer := &Remoteer{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			called = true
			cmd := exec.Command("echo", "origin\thttps://example.com (fetch)")
			return cmd
		},
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	remoteer.Remote([]string{"list"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stdout = oldStdout

	output := buf.String()
	if !called || output == "" || output[:6] != "origin" {
		t.Errorf("listサブコマンドの出力が想定と異なります: %s", output)
	}
}

func TestRemoteer_Remote_Add(t *testing.T) {
	called := false
	remoteer := &Remoteer{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			called = true
			return exec.Command("echo")
		},
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	remoteer.Remote([]string{"add", "origin", "https://example.com"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stdout = oldStdout
	if !called {
		t.Error("addサブコマンドでexecCommandが呼ばれていません")
	}
}

func TestRemoteer_Remote_Remove(t *testing.T) {
	called := false
	remoteer := &Remoteer{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			called = true
			return exec.Command("echo")
		},
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	remoteer.Remote([]string{"remove", "origin"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stdout = oldStdout
	if !called {
		t.Error("removeサブコマンドでexecCommandが呼ばれていません")
	}
}

func TestRemoteer_Remote_SetURL(t *testing.T) {
	called := false
	remoteer := &Remoteer{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			called = true
			return exec.Command("echo")
		},
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	remoteer.Remote([]string{"set-url", "origin", "https://example.com"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stdout = oldStdout
	if !called {
		t.Error("set-urlサブコマンドでexecCommandが呼ばれていません")
	}
}

func TestRemoteer_Remote_Help(t *testing.T) {
	remoteer := &Remoteer{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			return exec.Command("echo")
		},
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	remoteer.Remote([]string{"unknown"})

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
