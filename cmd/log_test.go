package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v4/git"
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

func TestLogger_Log_NoArgs(t *testing.T) {
	var buf bytes.Buffer
	l := &Logger{
		gitClient:    &mockLogGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	l.helper.outputWriter = &buf
	l.Log([]string{})

	output := buf.String()
	if output == "" || !strings.Contains(output, "Usage") {
		t.Errorf("Usage should be displayed when no args provided, but got: %s", output)
	}
}

func TestLogger_Log_Simple_OutputFormat(t *testing.T) {
	tests := []struct {
		name        string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "Successful simple log",
			shouldError: false,
		},
		{
			name:        "Error during simple log",
			shouldError: true,
			errorMsg:    "git error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			mockClient := &mockLogGitClient{}

			if tt.shouldError {
				mockClient.err = errors.New(tt.errorMsg)
			}

			l := &Logger{
				gitClient:    mockClient,
				outputWriter: &buf,
				helper:       NewHelper(),
			}
			l.helper.outputWriter = &buf

			l.Log([]string{"simple"})

			if tt.shouldError {
				output := buf.String()
				if !strings.Contains(output, "Error:") {
					t.Errorf("Expected error message, got: %s", output)
				}
				if !strings.Contains(output, tt.errorMsg) {
					t.Errorf("Expected error message to contain %q, got: %s", tt.errorMsg, output)
				}
			} else {
				if !mockClient.logSimpleCalled {
					t.Error("LogSimple should be called")
				}
			}
		})
	}
}

func TestLogger_Log_Graph_OutputFormat(t *testing.T) {
	tests := []struct {
		name        string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "Successful graph log",
			shouldError: false,
		},
		{
			name:        "Error during graph log",
			shouldError: true,
			errorMsg:    "git error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			mockClient := &mockLogGitClient{}

			if tt.shouldError {
				mockClient.err = errors.New(tt.errorMsg)
			}

			l := &Logger{
				gitClient:    mockClient,
				outputWriter: &buf,
				helper:       NewHelper(),
			}
			l.helper.outputWriter = &buf

			l.Log([]string{"graph"})

			if tt.shouldError {
				output := buf.String()
				if !strings.Contains(output, "Error:") {
					t.Errorf("Expected error message, got: %s", output)
				}
				if !strings.Contains(output, tt.errorMsg) {
					t.Errorf("Expected error message to contain %q, got: %s", tt.errorMsg, output)
				}
			} else {
				if !mockClient.logGraphCalled {
					t.Error("LogGraph should be called")
				}
			}
		})
	}
}
