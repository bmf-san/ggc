package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/bmf-san/ggc/v7/internal/prompt"
)

type mockPrompter struct {
	inputLine     string
	inputCanceled bool
	inputErr      error
}

func (m *mockPrompter) Input(_ string) (string, bool, error) {
	return m.inputLine, m.inputCanceled, m.inputErr
}

func (m *mockPrompter) Select(_ string, _ []string, _ string) (int, bool, error) {
	return 0, false, nil
}

func (m *mockPrompter) Confirm(_ string) (bool, bool, error) {
	return false, false, nil
}

func (m *mockPrompter) WithCancelMessage(_ string) prompt.Prompter {
	return m
}

func TestReadLine_Success(t *testing.T) {
	var buf bytes.Buffer
	p := &mockPrompter{inputLine: "test input"}

	line, ok := ReadLine(p, &buf, "prompt: ")

	if !ok {
		t.Error("ReadLine() ok = false, want true")
	}
	if line != "test input" {
		t.Errorf("ReadLine() line = %q, want %q", line, "test input")
	}
	if buf.Len() != 0 {
		t.Errorf("ReadLine() wrote output %q, want empty", buf.String())
	}
}

func TestReadLine_Canceled(t *testing.T) {
	var buf bytes.Buffer
	p := &mockPrompter{inputCanceled: true}

	line, ok := ReadLine(p, &buf, "prompt: ")

	if ok {
		t.Error("ReadLine() ok = true, want false")
	}
	if line != "" {
		t.Errorf("ReadLine() line = %q, want empty", line)
	}
}

func TestReadLine_Error(t *testing.T) {
	var buf bytes.Buffer
	p := &mockPrompter{inputErr: errors.New("read error")}

	line, ok := ReadLine(p, &buf, "prompt: ")

	if ok {
		t.Error("ReadLine() ok = true, want false")
	}
	if line != "" {
		t.Errorf("ReadLine() line = %q, want empty", line)
	}
	if got := buf.String(); got != "Error: read error\n" {
		t.Errorf("ReadLine() output = %q, want %q", got, "Error: read error\n")
	}
}

func TestReadLine_NilPrompter(t *testing.T) {
	var buf bytes.Buffer

	line, ok := ReadLine(nil, &buf, "prompt: ")

	if ok {
		t.Error("ReadLine() ok = true, want false")
	}
	if line != "" {
		t.Errorf("ReadLine() line = %q, want empty", line)
	}
}
