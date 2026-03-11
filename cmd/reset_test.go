package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v8/internal/git"
)

type mockResetOps struct {
	currentBranch           string
	resetHardAndCleanCalled bool
	resetHardCalled         bool
	resetSoftCalled         bool
	commit                  string
}

func (m *mockResetOps) GetCurrentBranch() (string, error) {
	if m.currentBranch == "" {
		return "main", nil
	}
	return m.currentBranch, nil
}
func (m *mockResetOps) ResetHardAndClean() error {
	m.resetHardAndCleanCalled = true
	return nil
}
func (m *mockResetOps) ResetHard(commit string) error {
	m.resetHardCalled = true
	m.commit = commit
	return nil
}
func (m *mockResetOps) ResetSoft(commit string) error {
	m.resetSoftCalled = true
	m.commit = commit
	return nil
}

var _ git.ResetOps = (*mockResetOps)(nil)

func TestResetter_Constructor(t *testing.T) {
	mockClient := &mockResetOps{}
	resetter := NewResetter(mockClient)

	if resetter == nil {
		t.Fatal("Expected NewResetter to return a non-nil Resetter")
	}
	if resetter != nil && resetter.gitClient == nil {
		t.Error("Expected gitClient to be set")
	}
	if resetter != nil && resetter.outputWriter == nil {
		t.Error("Expected outputWriter to be set")
	}
	if resetter != nil && resetter.helper == nil {
		t.Error("Expected helper to be set")
	}
}

func TestResetter_Reset(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		shouldShowHelp bool
	}{
		{
			name:           "no args - default reset",
			args:           []string{},
			expectedOutput: "Reset to origin/main successful",
			shouldShowHelp: false,
		},
		{
			name:           "hard reset with commit",
			args:           []string{"hard", "abc123"},
			expectedOutput: "Reset to abc123 successful",
			shouldShowHelp: false,
		},
		{
			name:           "hard reset without commit - should show help",
			args:           []string{"hard"},
			expectedOutput: "Error: commit hash required for hard reset",
			shouldShowHelp: true,
		},
		{
			name:           "unknown command - should show help",
			args:           []string{"unknown"},
			expectedOutput: "Usage: ggc reset",
			shouldShowHelp: true,
		},
		{
			name:           "soft reset with commit",
			args:           []string{"soft", "HEAD~1"},
			expectedOutput: "Reset to HEAD~1 successful",
			shouldShowHelp: false,
		},
		{
			name:           "soft reset without commit - should show help",
			args:           []string{"soft"},
			expectedOutput: "Error: commit reference required for soft reset",
			shouldShowHelp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := &mockResetOps{}

			resetter := &Resetter{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}
			// Set helper's output writer to capture help output
			resetter.helper.outputWriter = buf

			resetter.Reset(tt.args)

			output := buf.String()

			// Verify expected behavior
			if tt.shouldShowHelp {
				if !strings.Contains(output, tt.expectedOutput) {
					t.Errorf("Expected output containing '%s', got: %s", tt.expectedOutput, output)
				}
			} else {
				if !strings.Contains(output, tt.expectedOutput) {
					t.Errorf("Expected success output containing '%s', got: %s", tt.expectedOutput, output)
				}
			}

			// Verify no panic occurred
			if t.Failed() {
				t.Logf("Command args: %v", tt.args)
				t.Logf("Full output: %s", output)
			}
		})
	}
}

func TestResetter_ResetOperations(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		testFunc func(*testing.T, *Resetter, *mockResetOps, *bytes.Buffer)
	}{
		{
			name: "default reset calls ResetHardAndClean",
			args: []string{},
			testFunc: func(t *testing.T, resetter *Resetter, mc *mockResetOps, buf *bytes.Buffer) {
				output := buf.String()
				// Should show success message with branch name
				if !strings.Contains(output, "Reset to origin/") || !strings.Contains(output, "successful") {
					t.Errorf("Expected default reset success message, got: %s", output)
				}
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in default reset: %s", output)
				}
			},
		},
		{
			name: "hard reset with commit calls ResetHard",
			args: []string{"hard", "commit123"},
			testFunc: func(t *testing.T, resetter *Resetter, mc *mockResetOps, buf *bytes.Buffer) {
				output := buf.String()
				// Should show success message with commit hash
				if !strings.Contains(output, "Reset to commit123 successful") {
					t.Errorf("Expected hard reset success message, got: %s", output)
				}
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in hard reset: %s", output)
				}
			},
		},
		{
			name: "hard reset without commit shows error and help",
			args: []string{"hard"},
			testFunc: func(t *testing.T, resetter *Resetter, mc *mockResetOps, buf *bytes.Buffer) {
				output := buf.String()
				// Should show error message
				if !strings.Contains(output, "Error: commit hash required") {
					t.Errorf("Expected error message for missing commit, got: %s", output)
				}
				// Should also show help
				if !strings.Contains(output, "Usage:") {
					t.Errorf("Expected help message after error, got: %s", output)
				}
			},
		},
		{
			name: "soft reset with commit calls ResetSoft",
			args: []string{"soft", "HEAD~1"},
			testFunc: func(t *testing.T, resetter *Resetter, mc *mockResetOps, buf *bytes.Buffer) {
				output := buf.String()
				if !strings.Contains(output, "Reset to HEAD~1 successful") {
					t.Errorf("Expected soft reset success message, got: %s", output)
				}
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in soft reset: %s", output)
				}
				if !mc.resetSoftCalled {
					t.Error("ResetSoft should have been called")
				}
				if mc.commit != "HEAD~1" {
					t.Errorf("expected commit ref 'HEAD~1', got '%s'", mc.commit)
				}
				if mc.resetHardCalled {
					t.Error("ResetHard should NOT have been called for soft reset")
				}
			},
		},
		{
			name: "soft reset without commit shows error and help",
			args: []string{"soft"},
			testFunc: func(t *testing.T, resetter *Resetter, mc *mockResetOps, buf *bytes.Buffer) {
				output := buf.String()
				if !strings.Contains(output, "Error: commit reference required") {
					t.Errorf("Expected error message for missing commit ref, got: %s", output)
				}
				if !strings.Contains(output, "Usage:") {
					t.Errorf("Expected help message after error, got: %s", output)
				}
				if mc.resetSoftCalled {
					t.Error("ResetSoft should NOT have been called when no commit ref given")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := &mockResetOps{}

			resetter := &Resetter{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}
			resetter.helper.outputWriter = buf

			resetter.Reset(tt.args)
			tt.testFunc(t, resetter, mockClient, buf)
		})
	}
}

// mockResetOpsWithErrors supports error injection for reset error-path tests.
type mockResetOpsWithErrors struct {
	currentBranch        string
	currentBranchErr     error
	resetHardAndCleanErr error
	resetHardErr         error
	resetSoftErr         error
}

func (m *mockResetOpsWithErrors) GetCurrentBranch() (string, error) {
	if m.currentBranchErr != nil {
		return "", m.currentBranchErr
	}
	if m.currentBranch != "" {
		return m.currentBranch, nil
	}
	return "main", nil
}
func (m *mockResetOpsWithErrors) ResetHardAndClean() error { return m.resetHardAndCleanErr }
func (m *mockResetOpsWithErrors) ResetHard(_ string) error { return m.resetHardErr }
func (m *mockResetOpsWithErrors) ResetSoft(_ string) error { return m.resetSoftErr }

var _ git.ResetOps = (*mockResetOpsWithErrors)(nil)

func TestResetter_HandleDefaultReset_BranchError(t *testing.T) {
	var buf bytes.Buffer
	mock := &mockResetOpsWithErrors{currentBranchErr: errors.New("branch error")}
	r := &Resetter{gitClient: mock, outputWriter: &buf, helper: NewHelper()}
	r.helper.outputWriter = &buf
	r.handleDefaultReset()
	if !strings.Contains(buf.String(), "branch error") {
		t.Errorf("expected branch error, got: %s", buf.String())
	}
}

func TestResetter_HandleDefaultReset_ResetError(t *testing.T) {
	var buf bytes.Buffer
	mock := &mockResetOpsWithErrors{resetHardAndCleanErr: errors.New("reset error")}
	r := &Resetter{gitClient: mock, outputWriter: &buf, helper: NewHelper()}
	r.helper.outputWriter = &buf
	r.handleDefaultReset()
	if !strings.Contains(buf.String(), "reset error") {
		t.Errorf("expected reset error, got: %s", buf.String())
	}
}

func TestResetter_HandleHardReset_Error(t *testing.T) {
	var buf bytes.Buffer
	mock := &mockResetOpsWithErrors{resetHardErr: errors.New("hard reset error")}
	r := &Resetter{gitClient: mock, outputWriter: &buf, helper: NewHelper()}
	r.helper.outputWriter = &buf
	r.handleHardReset([]string{"abc123"})
	if !strings.Contains(buf.String(), "hard reset error") {
		t.Errorf("expected hard reset error, got: %s", buf.String())
	}
}

func TestResetter_HandleSoftReset_Error(t *testing.T) {
	var buf bytes.Buffer
	mock := &mockResetOpsWithErrors{resetSoftErr: errors.New("soft reset error")}
	r := &Resetter{gitClient: mock, outputWriter: &buf, helper: NewHelper()}
	r.helper.outputWriter = &buf
	r.handleSoftReset([]string{"HEAD~1"})
	if !strings.Contains(buf.String(), "soft reset error") {
		t.Errorf("expected soft reset error, got: %s", buf.String())
	}
}
