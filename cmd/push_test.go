package cmd

import (
	"bytes"
	"errors"
	"testing"
)

type mockPushGitClient struct {
	pushCalled bool
	pushForce  bool
	err        error
}

func (m *mockPushGitClient) Push(force bool) error {
	m.pushCalled = true
	m.pushForce = force
	return m.err
}

// All other methods intentionally omitted

func TestPusher_Push(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantPush bool
		wantErr  bool
		err      error
	}{
		{
			name:     "normal_push",
			args:     []string{"current"},
			wantPush: true,
			wantErr:  false,
		},
		{
			name:     "push_with_error",
			args:     []string{"current"},
			wantPush: true,
			wantErr:  true,
			err:      errors.New("push failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockPushGitClient{err: tt.err}
			var buf bytes.Buffer
			pusher := &Pusher{
				gitClient:    mockClient,
				outputWriter: &buf,
				helper:       NewHelper(),
			}
			pusher.helper.outputWriter = &buf
			pusher.Push(tt.args)

			if mockClient.pushCalled != tt.wantPush {
				t.Errorf("Push called = %v, want %v", mockClient.pushCalled, tt.wantPush)
			}

			if tt.wantErr {
				output := buf.String()
				if output != "Error: push failed\n" {
					t.Errorf("Output = %q, want %q", output, "Error: push failed\n")
				}
			}
		})
	}
}

func TestPusher_Push_Help(t *testing.T) {
	var buf bytes.Buffer
	pusher := &Pusher{
		gitClient:    &mockPushGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	pusher.helper.outputWriter = &buf
	pusher.Push([]string{})

	output := buf.String()
	if output == "" || !bytes.Contains(buf.Bytes(), []byte("Usage")) {
		t.Errorf("Usage should be displayed, but got: %s", output)
	}
}

func TestPusher_Push_Force(t *testing.T) {
	mockClient := &mockPushGitClient{}
	var buf bytes.Buffer
	pusher := &Pusher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	pusher.helper.outputWriter = &buf
	pusher.Push([]string{"force"})

	if !mockClient.pushCalled {
		t.Error("Push should be called")
	}
	if !mockClient.pushForce {
		t.Error("Push should be called with force=true")
	}
}

func TestPusher_Push_ForceError(t *testing.T) {
	mockClient := &mockPushGitClient{err: errors.New("force push failed")}
	var buf bytes.Buffer
	pusher := &Pusher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	pusher.helper.outputWriter = &buf
	pusher.Push([]string{"force"})

	if !mockClient.pushCalled {
		t.Error("Push should be called")
	}
	if !mockClient.pushForce {
		t.Error("Push should be called with force=true")
	}

	output := buf.String()
	if output != "Error: force push failed\n" {
		t.Errorf("Expected error message, got: %s", output)
	}
}

func TestPusher_Push_UnknownCommand(t *testing.T) {
	var buf bytes.Buffer
	pusher := &Pusher{
		gitClient:    &mockPushGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	pusher.helper.outputWriter = &buf
	pusher.Push([]string{"unknown"})

	output := buf.String()
	if output == "" || !bytes.Contains(buf.Bytes(), []byte("Usage")) {
		t.Errorf("Usage should be displayed for unknown command, but got: %s", output)
	}
}
