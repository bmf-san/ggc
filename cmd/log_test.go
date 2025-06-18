package cmd

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestLogger_Log_Graph_Success(t *testing.T) {
	logger := &Logger{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			return exec.Command("echo", "graph output")
		},
		logSimple: func() error { return nil },
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger.Log([]string{"graph"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("buf.ReadFrom failed: %v", err)
	}
	os.Stdout = oldStdout

	output := buf.String()
	if !strings.Contains(output, "graph output") {
		t.Errorf("graphサブコマンドの出力が想定と異なります: %s", output)
	}
}

func TestLogger_Log_Graph_Error(t *testing.T) {
	logger := &Logger{
		execCommand: func(name string, arg ...string) *exec.Cmd {
			return exec.Command("false")
		},
		logSimple: func() error { return nil },
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger.Log([]string{"graph"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("buf.ReadFrom failed: %v", err)
	}
	os.Stdout = oldStdout

	output := buf.String()
	if !strings.Contains(output, "Error:") {
		t.Errorf("graphサブコマンドのエラー出力が想定と異なります: %s", output)
	}
}

func TestLogger_Log_Simple_Success(t *testing.T) {
	logger := &Logger{
		execCommand: exec.Command,
		logSimple:   func() error { return nil },
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger.Log([]string{"simple"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("buf.ReadFrom failed: %v", err)
	}
	os.Stdout = oldStdout

	output := buf.String()
	if output != "" {
		t.Errorf("simpleサブコマンドの正常時は出力なし想定: %s", output)
	}
}

func TestLogger_Log_Simple_Error(t *testing.T) {
	logger := &Logger{
		execCommand: exec.Command,
		logSimple:   func() error { return errors.New("dummy error") },
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger.Log([]string{"simple"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("buf.ReadFrom failed: %v", err)
	}
	os.Stdout = oldStdout

	output := buf.String()
	if !strings.Contains(output, "Error:") {
		t.Errorf("simpleサブコマンドのエラー出力が想定と異なります: %s", output)
	}
}

func TestLogger_Log_Help(t *testing.T) {
	logger := &Logger{
		execCommand: exec.Command,
		logSimple:   func() error { return nil },
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger.Log([]string{"unknown"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("buf.ReadFrom failed: %v", err)
	}
	os.Stdout = oldStdout

	output := buf.String()
	if !strings.Contains(output, "Usage: ggc log simple | ggc log graph") {
		t.Errorf("ヘルプ出力が想定と異なります: %s", output)
	}
}
