package cmd

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestCommitter_Commit_AllowEmpty(t *testing.T) {
	called := false
	c := &Committer{
		CommitAllowEmpty: func() error {
			called = true
			return nil
		},
		CommitTmp: func() error {
			return nil
		},
	}
	c.Commit([]string{"allow-empty"})
	if !called {
		t.Error("CommitAllowEmpty should be called")
	}
}

func TestCommitter_Commit_Tmp(t *testing.T) {
	called := false
	c := &Committer{
		CommitAllowEmpty: func() error {
			return nil
		},
		CommitTmp: func() error {
			called = true
			return nil
		},
	}
	c.Commit([]string{"tmp"})
	if !called {
		t.Error("CommitTmp should be called")
	}
}

func TestCommitter_Commit_Help(t *testing.T) {
	c := &Committer{
		CommitAllowEmpty: func() error { return nil },
		CommitTmp:        func() error { return nil },
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	c.Commit([]string{"unknown"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stdout = oldStdout

	output := buf.String()
	if output == "" || output[:5] != "Usage" {
		t.Errorf("Usage should be displayed, but got: %s", output)
	}
}

func TestCommitter_Commit_AllowEmpty_Error(t *testing.T) {
	c := &Committer{
		CommitAllowEmpty: func() error { return errors.New("fail") },
		CommitTmp:        func() error { return nil },
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	c.Commit([]string{"allow-empty"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stdout = oldStdout

	output := buf.String()
	if output == "" || output[:5] != "Error" {
		t.Errorf("Error should be displayed, but got: %s", output)
	}
}
