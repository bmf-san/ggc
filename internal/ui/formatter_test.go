package ui

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestFormatter_Error(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf)

	f.Error(errors.New("something went wrong"))

	got := buf.String()
	want := "Error: something went wrong\n"
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestFormatter_Errorf(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf)

	f.Errorf("failed to %s: %v", "connect", "timeout")

	got := buf.String()
	want := "Error: failed to connect: timeout\n"
	if got != want {
		t.Errorf("Errorf() = %q, want %q", got, want)
	}
}

func TestFormatter_Print(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf)

	f.Print("hello")

	got := buf.String()
	want := "hello"
	if got != want {
		t.Errorf("Print() = %q, want %q", got, want)
	}
}

func TestFormatter_Println(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf)

	f.Println("hello")

	got := buf.String()
	want := "hello\n"
	if got != want {
		t.Errorf("Println() = %q, want %q", got, want)
	}
}

func TestFormatter_Printf(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf)

	f.Printf("count: %d", 42)

	got := buf.String()
	want := "count: 42"
	if got != want {
		t.Errorf("Printf() = %q, want %q", got, want)
	}
}

func TestFormatter_Header(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf)

	f.Header("Select items:")

	got := buf.String()
	// Should contain the message and ANSI codes
	if !strings.Contains(got, "Select items:") {
		t.Errorf("Header() should contain message, got %q", got)
	}
	if !strings.Contains(got, "\033[") {
		t.Errorf("Header() should contain ANSI codes, got %q", got)
	}
}

func TestFormatter_NumberedItem(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf)

	f.NumberedItem(1, "first item")

	got := buf.String()
	if !strings.Contains(got, "1") {
		t.Errorf("NumberedItem() should contain index, got %q", got)
	}
	if !strings.Contains(got, "first item") {
		t.Errorf("NumberedItem() should contain item text, got %q", got)
	}
}

func TestFormatter_Prompt(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf)

	f.Prompt()

	got := buf.String()
	want := "> "
	if got != want {
		t.Errorf("Prompt() = %q, want %q", got, want)
	}
}

func TestFormatter_Colors(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf)

	colors := f.Colors()
	if colors == nil {
		t.Fatal("Colors() should not return nil")
	}
	if colors.Reset != "\033[0m" {
		t.Errorf("Colors().Reset = %q, want %q", colors.Reset, "\033[0m")
	}
}

func TestFormatter_Writer(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf)

	w := f.Writer()
	if w != &buf {
		t.Error("Writer() should return the underlying writer")
	}
}

func TestFormatter_Success(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf)

	f.Success("Operation completed")

	got := buf.String()
	if !strings.Contains(got, "Operation completed") {
		t.Errorf("Success() should contain message, got %q", got)
	}
}

func TestFormatter_Warning(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf)

	f.Warning("Be careful")

	got := buf.String()
	if !strings.Contains(got, "Be careful") {
		t.Errorf("Warning() should contain message, got %q", got)
	}
}

func TestFormatter_ErrorHighlight(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf)

	f.ErrorHighlight("Critical failure")

	got := buf.String()
	if !strings.Contains(got, "Critical failure") {
		t.Errorf("ErrorHighlight() should contain message, got %q", got)
	}
}
