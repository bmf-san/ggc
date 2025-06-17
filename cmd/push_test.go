package cmd

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestPusher_Push_Current(t *testing.T) {
	called := false
	p := &Pusher{
		PushCurrentBranch: func() error {
			called = true
			return nil
		},
		PushForceCurrentBranch: func() error { return nil },
	}
	p.Push([]string{"current"})
	if !called {
		t.Error("PushCurrentBranchが呼ばれていません")
	}
}

func TestPusher_Push_Force(t *testing.T) {
	called := false
	p := &Pusher{
		PushCurrentBranch: func() error { return nil },
		PushForceCurrentBranch: func() error {
			called = true
			return nil
		},
	}
	p.Push([]string{"force"})
	if !called {
		t.Error("PushForceCurrentBranchが呼ばれていません")
	}
}

func TestPusher_Push_Help(t *testing.T) {
	p := &Pusher{
		PushCurrentBranch:      func() error { return nil },
		PushForceCurrentBranch: func() error { return nil },
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	p.Push([]string{"unknown"})

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

func TestPusher_Push_Current_Error(t *testing.T) {
	p := &Pusher{
		PushCurrentBranch:      func() error { return errors.New("fail") },
		PushForceCurrentBranch: func() error { return nil },
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	p.Push([]string{"current"})

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
