package cmd

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestResetter_Reset_Clean(t *testing.T) {
	called := false
	resetter := &Resetter{
		ResetClean: func() error {
			called = true
			return nil
		},
	}
	resetter.Reset([]string{"clean"})
	if !called {
		t.Error("ResetCleanが呼ばれていません")
	}
}

func TestResetter_Reset_Help(t *testing.T) {
	resetter := &Resetter{
		ResetClean: func() error { return nil },
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	resetter.Reset([]string{"unknown"})

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

func TestResetter_Reset_Clean_Error(t *testing.T) {
	resetter := &Resetter{
		ResetClean: func() error { return errors.New("fail") },
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	resetter.Reset([]string{"clean"})

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
