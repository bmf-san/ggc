package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/git"
)

type mockLogGitClient struct {
	git.Clienter
	logSimpleCalled bool
	logGraphCalled  bool
	err             error
}

func (m *mockLogGitClient) LogSimple() error {
	m.logSimpleCalled = true
	return m.err
}

func (m *mockLogGitClient) LogGraph() error {
	m.logGraphCalled = true
	return m.err
}

func TestLogger_Log_Simple(t *testing.T) {
	mockClient := &mockLogGitClient{}
	var buf bytes.Buffer
	l := &Logger{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	l.helper.outputWriter = &buf
	l.Log([]string{"simple"})
	if !mockClient.logSimpleCalled {
		t.Error("LogSimple should be called")
	}
}

func TestLogger_Log_Graph(t *testing.T) {
	mockClient := &mockLogGitClient{}
	var buf bytes.Buffer
	l := &Logger{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	l.helper.outputWriter = &buf
	l.Log([]string{"graph"})
	if !mockClient.logGraphCalled {
		t.Error("LogGraph should be called")
	}
}

func TestLogger_Log_Help(t *testing.T) {
	var buf bytes.Buffer
	l := &Logger{
		gitClient:    &mockLogGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	l.helper.outputWriter = &buf
	l.Log([]string{"unknown"})

	output := buf.String()
	if output == "" || !strings.Contains(output, "Usage") {
		t.Errorf("Usage should be displayed, but got: %s", output)
	}
}

func TestLogger_Log_Simple_Error(t *testing.T) {
	var buf bytes.Buffer
	l := &Logger{
		gitClient:    &mockLogGitClient{err: errors.New("fail")},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	l.helper.outputWriter = &buf
	l.Log([]string{"simple"})

	output := buf.String()
	if output != "Error: fail\n" {
		t.Errorf("unexpected output: got %q", output)
	}
}

func TestLogger_Log_Graph_Error(t *testing.T) {
	var buf bytes.Buffer
	l := &Logger{
		gitClient:    &mockLogGitClient{err: errors.New("fail")},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	l.helper.outputWriter = &buf
	l.Log([]string{"graph"})

	output := buf.String()
	if output != "Error: fail\n" {
		t.Errorf("unexpected output: got %q", output)
	}
}
