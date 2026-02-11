package cmd

import (
	"bytes"
	"errors"
	"testing"
)

func TestWriteError(t *testing.T) {
	var buf bytes.Buffer
	err := errors.New("test error")

	WriteError(&buf, err)

	got := buf.String()
	want := "Error: test error\n"
	if got != want {
		t.Errorf("WriteError() = %q, want %q", got, want)
	}
}

func TestWriteErrorf(t *testing.T) {
	var buf bytes.Buffer

	WriteErrorf(&buf, "failed to %s: %d", "process", 42)

	got := buf.String()
	want := "Error: failed to process: 42\n"
	if got != want {
		t.Errorf("WriteErrorf() = %q, want %q", got, want)
	}
}

func TestWriteLine(t *testing.T) {
	var buf bytes.Buffer

	WriteLine(&buf, "test message")

	got := buf.String()
	want := "test message\n"
	if got != want {
		t.Errorf("WriteLine() = %q, want %q", got, want)
	}
}

func TestWriteLinef(t *testing.T) {
	var buf bytes.Buffer

	WriteLinef(&buf, "count: %d", 5)

	got := buf.String()
	want := "count: 5\n"
	if got != want {
		t.Errorf("WriteLinef() = %q, want %q", got, want)
	}
}
