package cmd

import (
	"bytes"
	"os"
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
