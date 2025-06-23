package cmd

import (
	"bytes"
	"errors"
	"testing"
)

type mockResetGitClient struct {
	mockGitClient
	resetHardAndCleanCalled bool
	err                     error
}

func (m *mockResetGitClient) ResetHardAndClean() error {
	m.resetHardAndCleanCalled = true
	return m.err
}

func TestResetter_Reset(t *testing.T) {
	tests := []struct {
		name         string
		resetError   error
		expectOutput string
		expectReset  bool
	}{
		{
			name:         "successful reset",
			resetError:   nil,
			expectOutput: "",
			expectReset:  true,
		},
		{
			name:         "reset with error",
			resetError:   errors.New("reset failed"),
			expectOutput: "Error: reset failed\n",
			expectReset:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockResetGitClient{
				err: tt.resetError,
			}
			var buf bytes.Buffer
			resetter := &Resetter{
				gitClient:    mockClient,
				outputWriter: &buf,
			}

			resetter.Reset()

			if mockClient.resetHardAndCleanCalled != tt.expectReset {
				t.Errorf("Reset called = %v, want %v", mockClient.resetHardAndCleanCalled, tt.expectReset)
			}

			if got := buf.String(); got != tt.expectOutput {
				t.Errorf("Output = %q, want %q", got, tt.expectOutput)
			}
		})
	}
}
