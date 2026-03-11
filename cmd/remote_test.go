package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v8/internal/git"
)

type mockRemoteManager struct {
	listCalled   bool
	addCalled    bool
	removeCalled bool
	setURLCalled bool
	addName      string
	addURL       string
	removeName   string
	setName      string
	setURL       string
}

func (m *mockRemoteManager) RemoteList() error { m.listCalled = true; return nil }
func (m *mockRemoteManager) RemoteAdd(name, url string) error {
	m.addCalled = true
	m.addName = name
	m.addURL = url
	return nil
}
func (m *mockRemoteManager) RemoteRemove(name string) error {
	m.removeCalled = true
	m.removeName = name
	return nil
}
func (m *mockRemoteManager) RemoteSetURL(name, url string) error {
	m.setURLCalled = true
	m.setName = name
	m.setURL = url
	return nil
}

var _ git.RemoteManager = (*mockRemoteManager)(nil)

func TestRemoter_Constructor(t *testing.T) {
	mockClient := &mockRemoteManager{}
	remoter := NewRemoter(mockClient)

	if remoter == nil {
		t.Fatal("Expected NewRemoter to return a non-nil Remoter")
	}
	if remoter != nil && remoter.gitClient == nil {
		t.Error("Expected gitClient to be set")
	}
	if remoter != nil && remoter.outputWriter == nil {
		t.Error("Expected outputWriter to be set")
	}
	if remoter != nil && remoter.helper == nil {
		t.Error("Expected helper to be set")
	}
}

func TestRemoter_Remote(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		shouldShowHelp bool
	}{
		{
			name:           "no args - should show help",
			args:           []string{},
			expectedOutput: "Usage: ggc remote <command>",
			shouldShowHelp: true,
		},
		{
			name:           "list command",
			args:           []string{"list"},
			expectedOutput: "",
			shouldShowHelp: false,
		},
		{
			name:           "add command with correct args",
			args:           []string{"add", "origin", "https://github.com/user/repo.git"},
			expectedOutput: "Remote 'origin' added",
			shouldShowHelp: false,
		},
		{
			name:           "add command with incorrect args",
			args:           []string{"add", "origin"},
			expectedOutput: "Usage: ggc remote <command>",
			shouldShowHelp: true,
		},
		{
			name:           "remove command with correct args",
			args:           []string{"remove", "origin"},
			expectedOutput: "Remote 'origin' removed",
			shouldShowHelp: false,
		},
		{
			name:           "remove command with incorrect args",
			args:           []string{"remove"},
			expectedOutput: "Usage: ggc remote <command>",
			shouldShowHelp: true,
		},
		{
			name:           "set-url command with correct args",
			args:           []string{"set-url", "origin", "https://github.com/user/newrepo.git"},
			expectedOutput: "Remote 'origin' URL updated",
			shouldShowHelp: false,
		},
		{
			name:           "set-url command with incorrect args",
			args:           []string{"set-url", "origin"},
			expectedOutput: "Usage: ggc remote <command>",
			shouldShowHelp: true,
		},
		{
			name:           "unknown command",
			args:           []string{"unknown"},
			expectedOutput: "Usage: ggc remote <command>",
			shouldShowHelp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := &mockRemoteManager{}

			remoter := &Remoter{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}
			remoter.helper.outputWriter = buf

			remoter.Remote(tt.args)

			output := buf.String()

			if tt.shouldShowHelp {
				if !strings.Contains(output, tt.expectedOutput) {
					t.Errorf("Expected help output containing '%s', got: %s", tt.expectedOutput, output)
				}
			} else if tt.expectedOutput != "" {
				if !strings.Contains(output, tt.expectedOutput) {
					t.Errorf("Expected output containing '%s', got: %s", tt.expectedOutput, output)
				}
			}

			if t.Failed() {
				t.Logf("Command args: %v", tt.args)
				t.Logf("Full output: %s", output)
			}
		})
	}
}

func TestRemoter_RemoteOperations(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		testFunc func(*testing.T, *Remoter, *bytes.Buffer)
	}{
		{
			name: "list operation calls git client",
			args: []string{"list"},
			testFunc: func(t *testing.T, remoter *Remoter, buf *bytes.Buffer) {
				if buf.String() != "" && strings.Contains(buf.String(), "Error:") {
					t.Errorf("Unexpected error in list operation: %s", buf.String())
				}
			},
		},
		{
			name: "add operation with success",
			args: []string{"add", "upstream", "https://github.com/upstream/repo.git"},
			testFunc: func(t *testing.T, remoter *Remoter, buf *bytes.Buffer) {
				output := buf.String()
				if !strings.Contains(output, "Remote 'upstream' added") {
					t.Errorf("Expected success message for add operation, got: %s", output)
				}
			},
		},
		{
			name: "remove operation with success",
			args: []string{"remove", "upstream"},
			testFunc: func(t *testing.T, remoter *Remoter, buf *bytes.Buffer) {
				output := buf.String()
				if !strings.Contains(output, "Remote 'upstream' removed") {
					t.Errorf("Expected success message for remove operation, got: %s", output)
				}
			},
		},
		{
			name: "set-url operation with success",
			args: []string{"set-url", "origin", "https://github.com/newowner/repo.git"},
			testFunc: func(t *testing.T, remoter *Remoter, buf *bytes.Buffer) {
				output := buf.String()
				if !strings.Contains(output, "Remote 'origin' URL updated") {
					t.Errorf("Expected success message for set-url operation, got: %s", output)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := &mockRemoteManager{}

			remoter := &Remoter{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}
			remoter.helper.outputWriter = buf

			remoter.Remote(tt.args)
			tt.testFunc(t, remoter, buf)
		})
	}
}

func TestRemoter_List_Add_Remove_SetURL(t *testing.T) {
	m := &mockRemoteManager{}
	var buf bytes.Buffer
	r := &Remoter{gitClient: m, outputWriter: &buf, helper: NewHelper()}
	r.helper.outputWriter = &buf

	r.Remote([]string{"list"})
	if !m.listCalled {
		t.Fatal("expected RemoteList to be called")
	}

	buf.Reset()
	r.Remote([]string{"add", "origin", "https://example.com/repo.git"})
	if !m.addCalled || m.addName != "origin" || m.addURL != "https://example.com/repo.git" {
		t.Fatalf("unexpected add state: %+v", m)
	}
	if got := buf.String(); got == "" {
		t.Fatal("expected confirmation output for add")
	}

	buf.Reset()
	r.Remote([]string{"remove", "origin"})
	if !m.removeCalled || m.removeName != "origin" {
		t.Fatalf("expected RemoteRemove to be called with 'origin', got: %+v", m)
	}
	if got := buf.String(); got == "" {
		t.Fatal("expected confirmation output for remove")
	}

	buf.Reset()
	r.Remote([]string{"set-url", "origin", "https://example.com/new.git"})
	if !m.setURLCalled || m.setName != "origin" || m.setURL != "https://example.com/new.git" {
		t.Fatalf("unexpected set-url state: %+v", m)
	}
	if got := buf.String(); got == "" {
		t.Fatal("expected confirmation output for set-url")
	}
}

// mockRemoteManagerWithErrors allows injecting errors for error-path testing.
type mockRemoteManagerWithErrors struct {
	listErr   error
	addErr    error
	removeErr error
	setURLErr error
}

func (m *mockRemoteManagerWithErrors) RemoteList() error              { return m.listErr }
func (m *mockRemoteManagerWithErrors) RemoteAdd(_, _ string) error    { return m.addErr }
func (m *mockRemoteManagerWithErrors) RemoteRemove(_ string) error    { return m.removeErr }
func (m *mockRemoteManagerWithErrors) RemoteSetURL(_, _ string) error { return m.setURLErr }

var _ git.RemoteManager = (*mockRemoteManagerWithErrors)(nil)

func TestRemoter_ErrorPaths(t *testing.T) {
	errMsg := "remote error"
	cases := []struct {
		name string
		args []string
		mock *mockRemoteManagerWithErrors
	}{
		{
			name: "list error",
			args: []string{"list"},
			mock: &mockRemoteManagerWithErrors{listErr: errors.New(errMsg)},
		},
		{
			name: "add error",
			args: []string{"add", "origin", "https://example.com/r.git"},
			mock: &mockRemoteManagerWithErrors{addErr: errors.New(errMsg)},
		},
		{
			name: "remove error",
			args: []string{"remove", "origin"},
			mock: &mockRemoteManagerWithErrors{removeErr: errors.New(errMsg)},
		},
		{
			name: "set-url error",
			args: []string{"set-url", "origin", "https://example.com/n.git"},
			mock: &mockRemoteManagerWithErrors{setURLErr: errors.New(errMsg)},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			r := &Remoter{gitClient: tc.mock, outputWriter: &buf, helper: NewHelper()}
			r.helper.outputWriter = &buf
			r.Remote(tc.args)
			if !strings.Contains(buf.String(), errMsg) {
				t.Errorf("expected error %q in output, got: %q", errMsg, buf.String())
			}
		})
	}
}
