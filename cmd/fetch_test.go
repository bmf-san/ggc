package cmd

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestFetcher_Fetch_Prune(t *testing.T) {
	called := false
	f := &Fetcher{
		FetchPrune: func() error {
			called = true
			return nil
		},
	}
	f.Fetch([]string{"--prune"})
	if !called {
		t.Error("FetchPrune should be called")
	}
}

func TestFetcher_Fetch_Help(t *testing.T) {
	f := &Fetcher{
		FetchPrune: func() error { return nil },
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f.Fetch([]string{"unknown"})

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

func TestFetcher_Fetch_Prune_Error(t *testing.T) {
	f := &Fetcher{
		FetchPrune: func() error { return errors.New("fail") },
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f.Fetch([]string{"--prune"})

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
