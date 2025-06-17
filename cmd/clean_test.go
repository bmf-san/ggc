package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"os/exec"
	"testing"
)

func TestCleaner_Clean_Files(t *testing.T) {
	called := false
	cleaner := &Cleaner{
		CleanFiles: func() error {
			called = true
			return nil
		},
		CleanDirs: func() error { return nil },
	}
	cleaner.Clean([]string{"files"})
	if !called {
		t.Error("CleanFilesが呼ばれていません")
	}
}

func TestCleaner_Clean_Dirs(t *testing.T) {
	called := false
	cleaner := &Cleaner{
		CleanFiles: func() error { return nil },
		CleanDirs: func() error {
			called = true
			return nil
		},
	}
	cleaner.Clean([]string{"dirs"})
	if !called {
		t.Error("CleanDirsが呼ばれていません")
	}
}

func TestCleaner_Clean_Help(t *testing.T) {
	cleaner := &Cleaner{
		CleanFiles: func() error { return nil },
		CleanDirs:  func() error { return nil },
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cleaner.Clean([]string{"unknown"})

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

func TestCleaner_Clean_Files_Error(t *testing.T) {
	cleaner := &Cleaner{
		CleanFiles: func() error { return errors.New("fail") },
		CleanDirs:  func() error { return nil },
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cleaner.Clean([]string{"files"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stdout = oldStdout

	output := buf.String()
	if output == "" || output[:5] != "Error" {
		t.Errorf("エラー出力がされていません: %s", output)
	}
}

func TestCleaner_CleanInteractive_All(t *testing.T) {
	called := false
	cleaner := &Cleaner{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && len(arg) > 0 && arg[0] == "clean" {
				called = true
				return exec.Command("echo")
			}
			if name == "git" && len(arg) > 0 && arg[0] == "clean" && arg[1] == "-nd" {
				// git clean -nd の場合は候補を返す
				return exec.Command("echo", "Would remove foo.txt\nWould remove bar.txt")
			}
			return exec.Command("echo")
		},
		inputReader: bufio.NewReader(bytes.NewBufferString("all\n")),
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cleaner.CleanInteractive()

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stdout = oldStdout
	if !called {
		t.Errorf("all入力でgit cleanが呼ばれていません")
	}
}
