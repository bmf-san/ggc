package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/bmf-san/ggc/git"
)

type mockStashGitClient struct {
	git.Clienter
	stashPullPopCalled bool
	err                error
}

func (m *mockStashGitClient) StashPullPop() error {
	m.stashPullPopCalled = true
	return m.err
}

func TestStasher_Stash_Trash(t *testing.T) {
	mockClient := &mockStashGitClient{}
	var buf bytes.Buffer
	stasher := &Stasher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	stasher.helper.outputWriter = &buf
	stasher.Stash([]string{"trash"})

	if !mockClient.stashPullPopCalled {
		t.Error("StashPullPop should be called")
	}

	output := buf.String()
	if output != "add . → stash done\n" {
		t.Errorf("unexpected output: got %q", output)
	}
}

func TestStasher_Stash_Error(t *testing.T) {
	mockClient := &mockStashGitClient{err: errors.New("stash failed")}
	var buf bytes.Buffer
	stasher := &Stasher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	stasher.helper.outputWriter = &buf
	stasher.Stash([]string{"trash"})

	if !mockClient.stashPullPopCalled {
		t.Error("StashPullPop should be called")
	}

	output := buf.String()
	if output != "Error: stash failed\n" {
		t.Errorf("unexpected output: got %q", output)
	}
}

func TestStasher_Stash_Help(t *testing.T) {
	var buf bytes.Buffer
	stasher := &Stasher{
		gitClient:    &mockStashGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	stasher.helper.outputWriter = &buf
	stasher.Stash([]string{})

	output := buf.String()
	if output == "" || !bytes.Contains(buf.Bytes(), []byte("Usage")) {
		t.Errorf("Usage should be displayed, but got: %s", output)
	}
}
