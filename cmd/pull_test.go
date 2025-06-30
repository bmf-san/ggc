package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/bmf-san/ggc/git"
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
