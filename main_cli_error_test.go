package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v8/internal/git"
)

func TestWriteCLIError_PlainError(t *testing.T) {
	var buf bytes.Buffer
	writeCLIError(&buf, errors.New("boom"), false)
	got := buf.String()
	want := "Error: boom\n"
	if got != want {
		t.Fatalf("plain error: got %q, want %q", got, want)
	}
}

func TestWriteCLIError_OpErrorHidesCommandByDefault(t *testing.T) {
	var buf bytes.Buffer
	err := git.NewOpError("checkout branch", "git checkout main", errors.New("already on main"))
	writeCLIError(&buf, err, false)
	got := buf.String()

	if !strings.Contains(got, "Error: checkout branch failed") {
		t.Errorf("missing op summary: %q", got)
	}
	if !strings.Contains(got, "already on main") {
		t.Errorf("missing underlying message: %q", got)
	}
	if strings.Contains(got, "git checkout main") {
		t.Errorf("raw command leaked without verbose mode: %q", got)
	}
}

func TestWriteCLIError_OpErrorVerboseShowsCommand(t *testing.T) {
	var buf bytes.Buffer
	err := git.NewOpError("checkout branch", "git checkout main", errors.New("already on main"))
	writeCLIError(&buf, err, true)
	got := buf.String()

	if !strings.Contains(got, "detail: git checkout main") {
		t.Errorf("verbose mode should include operation detail: %q", got)
	}
}

func TestWriteCLIError_OpErrorNilUnderlyingError(t *testing.T) {
	var buf bytes.Buffer
	err := git.NewOpError("add files", "git add .", nil)
	writeCLIError(&buf, err, false)
	got := buf.String()

	if !strings.Contains(got, "Error: add files failed") {
		t.Errorf("missing op summary: %q", got)
	}
	if strings.Contains(got, "<nil>") {
		t.Errorf("nil underlying error should not be rendered: %q", got)
	}
	if strings.Contains(got, "git add .") {
		t.Errorf("operation detail leaked without verbose mode: %q", got)
	}
}

func TestWriteCLIError_WrappedOpError(t *testing.T) {
	var buf bytes.Buffer
	inner := git.NewOpError("push", "git push origin main", errors.New("rejected"))
	wrapped := errors.Join(inner, errors.New("post-run cleanup also failed"))
	writeCLIError(&buf, wrapped, false)
	got := buf.String()

	if !strings.Contains(got, "Error: push failed") {
		t.Errorf("errors.As through join should still find OpError: %q", got)
	}
}
