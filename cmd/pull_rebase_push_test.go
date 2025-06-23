package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/bmf-san/ggc/git"
)

type mockPullRebasePushGitClient struct {
	git.Clienter
	pullCalled bool
	pullRebase bool
	pushCalled bool
	pushForce  bool
	err        error
}

func (m *mockPullRebasePushGitClient) Pull(rebase bool) error {
	m.pullCalled = true
	m.pullRebase = rebase
	return m.err
}

func (m *mockPullRebasePushGitClient) Push(force bool) error {
	m.pushCalled = true
	m.pushForce = force
	return m.err
}

func TestPullRebasePusher_PullRebasePush(t *testing.T) {
	mockClient := &mockPullRebasePushGitClient{}
	var buf bytes.Buffer
	prp := NewPullRebasePusherWithClient(mockClient)
	prp.outputWriter = &buf
	prp.PullRebasePush()

	if !mockClient.pullCalled {
		t.Error("Pull should be called")
	}
	if !mockClient.pullRebase {
		t.Error("Pull should be called with rebase=true")
	}
	if !mockClient.pushCalled {
		t.Error("Push should be called")
	}
	if mockClient.pushForce {
		t.Error("Push should be called with force=false")
	}

	output := buf.String()
	if output != "pull→rebase→push completed\n" {
		t.Errorf("unexpected output: got %q", output)
	}
}

func TestPullRebasePusher_PullRebasePush_Error(t *testing.T) {
	mockClient := &mockPullRebasePushGitClient{err: errors.New("fail")}
	var buf bytes.Buffer
	prp := NewPullRebasePusherWithClient(mockClient)
	prp.outputWriter = &buf
	prp.PullRebasePush()

	output := buf.String()
	if output != "Error: fail\n" {
		t.Errorf("unexpected output: got %q", output)
	}
}
