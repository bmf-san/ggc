package cmd

import (
	"bytes"
	"errors"
	"testing"
)

func TestWriteError(t *testing.T) {
	var buf bytes.Buffer
	err := errors.New("something went wrong")
	WriteError(&buf, err)

	want := "Error: something went wrong\n"
	if got := buf.String(); got != want {
		t.Errorf("WriteError() = %q, want %q", got, want)
	}
}

func TestWriteErrorf(t *testing.T) {
	var buf bytes.Buffer
	WriteErrorf(&buf, "failed to open %s", "file.txt")

	want := "Error: failed to open file.txt\n"
	if got := buf.String(); got != want {
		t.Errorf("WriteErrorf() = %q, want %q", got, want)
	}
}

func TestWriteLine(t *testing.T) {
	var buf bytes.Buffer
	WriteLine(&buf, "hello world")

	want := "hello world\n"
	if got := buf.String(); got != want {
		t.Errorf("WriteLine() = %q, want %q", got, want)
	}
}

func TestWriteLinef(t *testing.T) {
	var buf bytes.Buffer
	WriteLinef(&buf, "count: %d", 42)

	want := "count: 42\n"
	if got := buf.String(); got != want {
		t.Errorf("WriteLinef() = %q, want %q", got, want)
	}
}
