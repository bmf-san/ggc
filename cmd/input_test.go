package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/bmf-san/ggc/v7/internal/prompt"
)

type mockPrompter struct {
	input    string
	canceled bool
	err      error
}

func (m *mockPrompter) Input(_ string) (string, bool, error) {
	return m.input, m.canceled, m.err
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
	p := &mockPrompter{input: "test input"}

	result, ok := ReadLine(p, &buf, "Enter:")
	if !ok {
		t.Error("ReadLine() ok = false, want true")
	}
	if result != "test input" {
		t.Errorf("ReadLine() = %q, want %q", result, "test input")
	}
	if buf.Len() != 0 {
		t.Errorf("ReadLine() wrote output %q, want empty", buf.String())
	}
}

func TestReadLine_Canceled(t *testing.T) {
	var buf bytes.Buffer
	p := &mockPrompter{canceled: true}

	result, ok := ReadLine(p, &buf, "Enter:")
	if ok {
		t.Error("ReadLine() ok = true, want false")
	}
	if result != "" {
		t.Errorf("ReadLine() = %q, want empty string", result)
	}
}

func TestReadLine_Error(t *testing.T) {
	var buf bytes.Buffer
	p := &mockPrompter{err: errors.New("input error")}

	result, ok := ReadLine(p, &buf, "Enter:")
	if ok {
		t.Error("ReadLine() ok = true, want false")
	}
	if result != "" {
		t.Errorf("ReadLine() = %q, want empty string", result)
	}
	// Verify error was written
	if buf.String() != "Error: input error\n" {
		t.Errorf("Error output = %q, want %q", buf.String(), "Error: input error\n")
	}
}

func TestReadLine_NilPrompter(t *testing.T) {
	var buf bytes.Buffer

	result, ok := ReadLine(nil, &buf, "Enter:")
	if ok {
		t.Error("ReadLine() ok = true, want false")
	}
	if result != "" {
		t.Errorf("ReadLine() = %q, want empty string", result)
	}
}
