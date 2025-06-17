package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
)

func TestStasher_Stash_Trash(t *testing.T) {
	calls := []string{}
	stasher := &Stasher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			calls = append(calls, name+" "+arg[0])
			return exec.Command("echo")
		},
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	stasher.Stash([]string{"trash"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stdout = oldStdout

	output := buf.String()
	if output == "" || output[:15] != "add . → stash" {
		t.Errorf("add . → stash doneが出力されていません: %s", output)
	}
	if len(calls) != 2 || calls[0] != "git add" || calls[1] != "git stash" {
		t.Errorf("コマンド呼び出しが想定と異なります: %v", calls)
	}
}

func TestStasher_Stash_Help(t *testing.T) {
	stasher := &Stasher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			return exec.Command("echo")
		},
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	stasher.Stash([]string{"unknown"})

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

func TestStasher_Stash_AddError(t *testing.T) {
	stasher := &Stasher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if arg[0] == "add" {
				return exec.Command("false")
			}
			return exec.Command("echo")
		},
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	stasher.Stash([]string{"trash"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stdout = oldStdout

	output := buf.String()
	if output == "" || output[:5] != "Error" {
		t.Errorf("addコマンド失敗時のエラー出力がされていません: %s", output)
	}
}

func TestStasher_Stash_StashError(t *testing.T) {
	stasher := &Stasher{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if arg[0] == "stash" {
				return exec.Command("false")
			}
			return exec.Command("echo")
		},
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	stasher.Stash([]string{"trash"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stdout = oldStdout

	output := buf.String()
	if output == "" || output[:5] != "Error" {
		t.Errorf("stashコマンド失敗時のエラー出力がされていません: %s", output)
	}
}
