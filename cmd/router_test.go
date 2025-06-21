package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestComplete(t *testing.T) {
	// Save the original NewCompleter
	originalNewCompleter := NewCompleter
	defer func() {
		// Restore it after the test finishes
		NewCompleter = originalNewCompleter
	}()

	// Mock NewCompleter
	NewCompleter = func() *Completer {
		return &Completer{
			listLocalBranches: func() ([]string, error) {
				return []string{"feature/test-branch", "main"}, nil
			},
		}
	}

	// Capture standard output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute the function under test
	// Call it with "sub" argument for "branch" subcommand
	Complete([]string{"branch", "sub"})

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}

	// Check if the output contains the expected branch name
	expected := "feature/test-branch"
	if !bytes.Contains(buf.Bytes(), []byte(expected)) {
		t.Errorf("expected output to contain %q, but got %q", expected, buf.String())
	}

	expected = "main"
	if !bytes.Contains(buf.Bytes(), []byte(expected)) {
		t.Errorf("expected output to contain %q, but got %q", expected, buf.String())
	}
}
