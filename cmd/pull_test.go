package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/bmf-san/ggc/v3/git"
)

type mockPullGitClient struct {
	git.Clienter
	pullCalled bool
	pullRebase bool
	err        error
}

func (m *mockPullGitClient) Pull(rebase bool) error {
	m.pullCalled = true
	m.pullRebase = rebase
	return m.err
}

func TestPuller_Pull(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantPull bool
		wantErr  bool
		err      error
	}{
		{
			name:     "normal_pull",
			args:     []string{"current"},
			wantPull: true,
			wantErr:  false,
		},
		{
			name:     "pull_with_error",
			args:     []string{"current"},
			wantPull: true,
			wantErr:  true,
			err:      errors.New("pull failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockPullGitClient{err: tt.err}
			var buf bytes.Buffer
			puller := &Puller{
				gitClient:    mockClient,
				outputWriter: &buf,
				helper:       NewHelper(),
			}
			puller.helper.outputWriter = &buf
			puller.Pull(tt.args)

			if mockClient.pullCalled != tt.wantPull {
				t.Errorf("Pull called = %v, want %v", mockClient.pullCalled, tt.wantPull)
			}

			if tt.wantErr {
				output := buf.String()
				if output != "Error: pull failed\n" {
					t.Errorf("Output = %q, want %q", output, "Error: pull failed\n")
				}
			}
		})
	}
}

func TestPuller_Pull_Help(t *testing.T) {
	var buf bytes.Buffer
	puller := &Puller{
		gitClient:    &mockPullGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	puller.helper.outputWriter = &buf
	puller.Pull([]string{})

	output := buf.String()
	if output == "" || !bytes.Contains(buf.Bytes(), []byte("Usage")) {
		t.Errorf("Usage should be displayed, but got: %s", output)
	}
}

func TestPuller_Pull_Rebase(t *testing.T) {
	mockClient := &mockPullGitClient{}
	var buf bytes.Buffer
	puller := &Puller{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	puller.helper.outputWriter = &buf
	puller.Pull([]string{"rebase"})

	if !mockClient.pullCalled {
		t.Error("Pull should be called")
	}
	if !mockClient.pullRebase {
		t.Error("Pull should be called with rebase=true")
	}
}

func TestPuller_Pull_RebaseError(t *testing.T) {
	mockClient := &mockPullGitClient{err: errors.New("rebase failed")}
	var buf bytes.Buffer
	puller := &Puller{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	puller.helper.outputWriter = &buf
	puller.Pull([]string{"rebase"})

	if !mockClient.pullCalled {
		t.Error("Pull should be called")
	}
	if !mockClient.pullRebase {
		t.Error("Pull should be called with rebase=true")
	}

	output := buf.String()
	if output != "Error: rebase failed\n" {
		t.Errorf("Expected error message, got: %s", output)
	}
}

func TestPuller_Pull_UnknownCommand(t *testing.T) {
	var buf bytes.Buffer
	puller := &Puller{
		gitClient:    &mockPullGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	puller.helper.outputWriter = &buf
	puller.Pull([]string{"unknown"})

	output := buf.String()
	if output == "" || !bytes.Contains(buf.Bytes(), []byte("Usage")) {
		t.Errorf("Usage should be displayed for unknown command, but got: %s", output)
	}
}
