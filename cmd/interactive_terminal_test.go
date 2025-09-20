package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"golang.org/x/term"
)

func TestUISetupTerminalNonFile(t *testing.T) {
	mockTerm := &mockTerminal{}
	stderr := &bytes.Buffer{}
	ui := &UI{
		stdin:  bytes.NewBuffer(nil),
		term:   mockTerm,
		stderr: stderr,
	}

	oldState, ok := ui.setupTerminal()
	if !ok {
		t.Fatalf("expected setup to succeed when stdin is not a file")
	}
	if oldState != nil {
		t.Fatalf("want nil old state for non-file stdin, got %#v", oldState)
	}
	if mockTerm.makeRawCalled {
		t.Fatalf("makeRaw should not be called for non-file stdin")
	}
	if mockTerm.restoreCalled {
		t.Fatalf("restore should not be called for non-file stdin")
	}
	if stderr.Len() != 0 {
		t.Fatalf("unexpected stderr output: %q", stderr.String())
	}
}

func TestUISetupTerminalMakeRawFailure(t *testing.T) {
	mockTerm := &mockTerminal{shouldFailRaw: true}
	stderr := &bytes.Buffer{}
	ui := &UI{
		stdin:  os.Stdin,
		term:   mockTerm,
		stderr: stderr,
	}

	oldState, ok := ui.setupTerminal()
	if ok {
		t.Fatalf("expected setup to fail when makeRaw errors")
	}
	if oldState != nil {
		t.Fatalf("expected nil state on failure, got %#v", oldState)
	}
	if !mockTerm.makeRawCalled {
		t.Fatalf("expected makeRaw to be called")
	}
	if mockTerm.restoreCalled {
		t.Fatalf("restore should not be called when makeRaw fails")
	}
	if !strings.Contains(stderr.String(), "Failed to set terminal to raw mode") {
		t.Fatalf("expected failure message in stderr, got %q", stderr.String())
	}
}

func TestDefaultTerminalMakeRawInvalidFD(t *testing.T) {
	dt := &defaultTerminal{}
	if _, err := dt.makeRaw(-1); err == nil {
		t.Fatalf("expected error for invalid file descriptor")
	}
}

func TestDefaultTerminalRestoreInvalidFD(t *testing.T) {
	dt := &defaultTerminal{}
	if err := dt.restore(-1, &term.State{}); err == nil {
		t.Fatalf("expected error for invalid file descriptor")
	}
}

func TestUIWriteError(t *testing.T) {
	stderr := &bytes.Buffer{}
	ui := &UI{stderr: stderr}

	ui.writeError("error: %s", "failed")

	if got, want := stderr.String(), "error: failed\n"; got != want {
		t.Fatalf("writeError = %q, want %q", got, want)
	}
}

func TestUIWriteln(t *testing.T) {
	stdout := &bytes.Buffer{}
	ui := &UI{stdout: stdout}

	ui.writeln("line %d", 42)

	if got, want := stdout.String(), "\r\x1b[Kline 42\r\n"; got != want {
		t.Fatalf("writeln output = %q, want %q", got, want)
	}
}
