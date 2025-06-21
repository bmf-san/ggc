package cmd

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestPuller_Pull_Current(t *testing.T) {
	called := false
	p := &Puller{
		PullCurrentBranch: func() error {
			called = true
			return nil
		},
		PullRebaseCurrentBranch: func() error { return nil },
	}
	p.Pull([]string{"current"})
	if !called {
		t.Error("PullCurrentBranch should be called")
	}
}

func TestPuller_Pull_Rebase(t *testing.T) {
	called := false
	p := &Puller{
		PullCurrentBranch: func() error { return nil },
		PullRebaseCurrentBranch: func() error {
			called = true
			return nil
		},
	}
	p.Pull([]string{"rebase"})
	if !called {
		t.Error("PullRebaseCurrentBranch should be called")
	}
}

func TestPuller_Pull_Help(t *testing.T) {
	p := &Puller{
		PullCurrentBranch:       func() error { return nil },
		PullRebaseCurrentBranch: func() error { return nil },
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	p.Pull([]string{"unknown"})

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

func TestPuller_Pull_Current_Error(t *testing.T) {
	p := &Puller{
		PullCurrentBranch:       func() error { return errors.New("fail") },
		PullRebaseCurrentBranch: func() error { return nil },
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	p.Pull([]string{"current"})

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
