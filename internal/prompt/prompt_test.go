package prompt

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"syscall"
	"testing"
)

func TestNewDefault(t *testing.T) {
	p := NewDefault()
	if p == nil {
		t.Fatal("NewDefault() returned nil")
	}
}

func TestNew_NilReader(t *testing.T) {
	p := New(nil, &bytes.Buffer{})
	if p == nil {
		t.Fatal("New(nil, writer) returned nil")
	}
}

func TestNew_NilWriter(t *testing.T) {
	p := New(strings.NewReader(""), nil)
	if p == nil {
		t.Fatal("New(reader, nil) returned nil")
	}
}

func TestWithCancelMessage_Override(t *testing.T) {
	var buf bytes.Buffer
	p := New(strings.NewReader(""), &buf).(*StandardPrompter)
	p2 := p.WithCancelMessage("custom cancel").(*StandardPrompter)
	if p2.cancelMessage != "custom cancel" {
		t.Errorf("cancelMessage = %q, want %q", p2.cancelMessage, "custom cancel")
	}
}

func TestWithCancelMessage_EmptyKeepsDefault(t *testing.T) {
	var buf bytes.Buffer
	p := New(strings.NewReader(""), &buf).(*StandardPrompter)
	p.WithCancelMessage("")
	if p.cancelMessage != defaultCancelMessage {
		t.Errorf("empty override changed cancelMessage to %q", p.cancelMessage)
	}
}

func TestInputWithBuffered_Success(t *testing.T) {
	var buf bytes.Buffer
	p := New(strings.NewReader("hello world\n"), &buf).(*StandardPrompter)
	line, canceled, err := p.inputWithBuffered("input: ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if canceled {
		t.Fatal("expected not canceled")
	}
	if line != "hello world" {
		t.Errorf("got %q, want %q", line, "hello world")
	}
	if !strings.Contains(buf.String(), "input: ") {
		t.Errorf("prompt not written to output, got: %q", buf.String())
	}
}

func TestInputWithBuffered_CRLFNewline(t *testing.T) {
	var buf bytes.Buffer
	p := New(strings.NewReader("hello\r\n"), &buf).(*StandardPrompter)
	line, canceled, err := p.inputWithBuffered("> ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if canceled {
		t.Fatal("expected not canceled")
	}
	if line != "hello" {
		t.Errorf("got %q, want %q", line, "hello")
	}
}

func TestInputWithBuffered_EOFEmpty(t *testing.T) {
	var buf bytes.Buffer
	p := New(strings.NewReader(""), &buf).(*StandardPrompter)
	_, _, err := p.inputWithBuffered("input: ")
	if !errors.Is(err, io.EOF) {
		t.Fatalf("expected io.EOF, got: %v", err)
	}
}

func TestInputWithBuffered_EOFWithContent(t *testing.T) {
	var buf bytes.Buffer
	// No trailing newline — bufio will return content with io.EOF
	p := New(strings.NewReader("partial"), &buf).(*StandardPrompter)
	line, canceled, err := p.inputWithBuffered("input: ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if canceled {
		t.Fatal("expected not canceled")
	}
	if line != "partial" {
		t.Errorf("got %q, want %q", line, "partial")
	}
}

func TestSelect_Success(t *testing.T) {
	var buf bytes.Buffer
	p := New(strings.NewReader("2\n"), &buf)
	idx, canceled, err := p.Select("Choose:", []string{"foo", "bar", "baz"}, "> ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if canceled {
		t.Fatal("expected not canceled")
	}
	if idx != 1 {
		t.Errorf("got index %d, want 1", idx)
	}
}

func TestSelect_NoItems(t *testing.T) {
	var buf bytes.Buffer
	p := New(strings.NewReader(""), &buf)
	_, _, err := p.Select("Choose:", nil, "> ")
	if err == nil {
		t.Fatal("expected error for empty items list")
	}
}

func TestSelect_InvalidSelection(t *testing.T) {
	var buf bytes.Buffer
	p := New(strings.NewReader("99\n"), &buf)
	_, _, err := p.Select("Choose:", []string{"foo"}, "> ")
	if !errors.Is(err, ErrInvalidSelection) {
		t.Fatalf("expected ErrInvalidSelection, got %v", err)
	}
}

func TestSelect_NonNumericInput(t *testing.T) {
	var buf bytes.Buffer
	p := New(strings.NewReader("abc\n"), &buf)
	_, _, err := p.Select("Choose:", []string{"foo"}, "> ")
	if !errors.Is(err, ErrInvalidSelection) {
		t.Fatalf("expected ErrInvalidSelection, got %v", err)
	}
}

func TestSelect_NoTitle(t *testing.T) {
	var buf bytes.Buffer
	p := New(strings.NewReader("1\n"), &buf)
	idx, canceled, err := p.Select("", []string{"only"}, "> ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if canceled {
		t.Fatal("expected not canceled")
	}
	if idx != 0 {
		t.Errorf("got index %d, want 0", idx)
	}
}

func TestConfirm_Yes(t *testing.T) {
	cases := []string{"y\n", "yes\n", "Y\n", "YES\n", "Yes\n"}
	for _, input := range cases {
		var buf bytes.Buffer
		p := New(strings.NewReader(input), &buf)
		got, canceled, err := p.Confirm("ok? ")
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", input, err)
		}
		if canceled {
			t.Fatalf("input %q: expected not canceled", input)
		}
		if !got {
			t.Errorf("input %q: expected true, got false", input)
		}
	}
}

func TestConfirm_No(t *testing.T) {
	cases := []string{"n\n", "no\n", "N\n", "NO\n", "\n"}
	for _, input := range cases {
		var buf bytes.Buffer
		p := New(strings.NewReader(input), &buf)
		got, canceled, err := p.Confirm("ok? ")
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", input, err)
		}
		if canceled {
			t.Fatalf("input %q: expected not canceled", input)
		}
		if got {
			t.Errorf("input %q: expected false, got true", input)
		}
	}
}

func TestConfirm_InvalidAnswer(t *testing.T) {
	var buf bytes.Buffer
	p := New(strings.NewReader("maybe\n"), &buf)
	_, canceled, err := p.Confirm("ok? ")
	if err == nil {
		t.Fatal("expected error for invalid answer")
	}
	if canceled {
		t.Fatal("expected not canceled")
	}
	if !errors.Is(err, ErrInvalidConfirmation) {
		t.Errorf("expected ErrInvalidConfirmation, got %v", err)
	}
}

func TestDisplaySelection_WithTitle(t *testing.T) {
	var buf bytes.Buffer
	p := New(strings.NewReader(""), &buf).(*StandardPrompter)
	err := p.displaySelection("Pick one:", []string{"alpha", "beta"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "Pick one:") {
		t.Errorf("missing title in output: %q", got)
	}
	if !strings.Contains(got, "[1] alpha") {
		t.Errorf("missing item 1 in output: %q", got)
	}
	if !strings.Contains(got, "[2] beta") {
		t.Errorf("missing item 2 in output: %q", got)
	}
}

func TestDisplaySelection_NoTitle(t *testing.T) {
	var buf bytes.Buffer
	p := New(strings.NewReader(""), &buf).(*StandardPrompter)
	err := p.displaySelection("", []string{"one"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[1] one") {
		t.Errorf("missing item in output: %q", buf.String())
	}
}

func TestParseSelection_Valid(t *testing.T) {
	p := &StandardPrompter{}
	idx, _, err := p.parseSelection("2", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idx != 1 {
		t.Errorf("got %d, want 1", idx)
	}
}

func TestParseSelection_OutOfRange(t *testing.T) {
	p := &StandardPrompter{}
	_, _, err := p.parseSelection("5", 3)
	if !errors.Is(err, ErrInvalidSelection) {
		t.Fatalf("expected ErrInvalidSelection, got %v", err)
	}
}

func TestParseSelection_Zero(t *testing.T) {
	p := &StandardPrompter{}
	_, _, err := p.parseSelection("0", 3)
	if !errors.Is(err, ErrInvalidSelection) {
		t.Fatalf("expected ErrInvalidSelection for 0, got %v", err)
	}
}

func TestParseSelection_NotANumber(t *testing.T) {
	p := &StandardPrompter{}
	_, _, err := p.parseSelection("abc", 3)
	if !errors.Is(err, ErrInvalidSelection) {
		t.Fatalf("expected ErrInvalidSelection, got %v", err)
	}
}

func TestParseSelection_WithSpaces(t *testing.T) {
	p := &StandardPrompter{}
	idx, _, err := p.parseSelection("  1  ", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idx != 0 {
		t.Errorf("got %d, want 0", idx)
	}
}

func TestTrimNewline(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"hello\n", "hello"},
		{"hello\r\n", "hello"},
		{"hello\r", "hello"},
		{"hello", "hello"},
		{"", ""},
		{"\n", ""},
		{"\r\n", ""},
	}
	for _, tt := range tests {
		got := trimNewline(tt.in)
		if got != tt.want {
			t.Errorf("trimNewline(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestIsInterrupted(t *testing.T) {
	if isInterrupted(nil) {
		t.Error("nil should return false")
	}
	if !isInterrupted(syscall.EINTR) {
		t.Error("syscall.EINTR should return true")
	}
	if !isInterrupted(errors.New("interrupted")) {
		t.Error("error containing 'interrupted' should return true")
	}
	if isInterrupted(errors.New("connection refused")) {
		t.Error("unrelated error should return false")
	}
}

func TestNilPrompter_Input(t *testing.T) {
	var p *StandardPrompter
	line, canceled, err := p.Input("prompt: ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !canceled {
		t.Fatal("nil prompter Input should return canceled=true")
	}
	if line != "" {
		t.Errorf("nil prompter Input should return empty line, got %q", line)
	}
}

func TestNilPrompter_Select(t *testing.T) {
	var p *StandardPrompter
	idx, canceled, err := p.Select("title", []string{"a"}, "> ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !canceled {
		t.Fatal("nil prompter Select should return canceled=true")
	}
	if idx != -1 {
		t.Errorf("nil prompter Select should return -1, got %d", idx)
	}
}

func TestNilPrompter_Confirm(t *testing.T) {
	var p *StandardPrompter
	got, canceled, err := p.Confirm("ok? ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !canceled {
		t.Fatal("nil prompter Confirm should return canceled=true")
	}
	if got {
		t.Error("nil prompter Confirm should return false")
	}
}

func TestTerminalReadWriter_NilReader(t *testing.T) {
	rw := &terminalReadWriter{reader: nil, writer: &bytes.Buffer{}}
	buf := make([]byte, 4)
	_, err := rw.Read(buf)
	if !errors.Is(err, errNilTerminalReader) {
		t.Errorf("expected errNilTerminalReader, got %v", err)
	}
}

func TestTerminalReadWriter_NilWriter(t *testing.T) {
	rw := &terminalReadWriter{writer: nil}
	_, err := rw.Write([]byte("x"))
	if !errors.Is(err, errNilTerminalWriter) {
		t.Errorf("expected errNilTerminalWriter, got %v", err)
	}
}

func TestTerminalReadWriter_NilStruct(t *testing.T) {
	var rw *terminalReadWriter
	buf := make([]byte, 4)
	_, err := rw.Read(buf)
	if !errors.Is(err, errNilTerminalReader) {
		t.Errorf("nil Read: expected errNilTerminalReader, got %v", err)
	}
	_, err = rw.Write([]byte("x"))
	if !errors.Is(err, errNilTerminalWriter) {
		t.Errorf("nil Write: expected errNilTerminalWriter, got %v", err)
	}
}
