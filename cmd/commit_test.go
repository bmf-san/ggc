package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/git"
)

// mockGitClient for commit_test
type mockCommitGitClient struct {
	git.Clienter
	commitAllowEmptyCalled bool
	commitTmpCalled        bool
	err                    error
}

func (m *mockCommitGitClient) CommitAllowEmpty() error {
	m.commitAllowEmptyCalled = true
	return m.err
}

func (m *mockCommitGitClient) CommitTmp() error {
	m.commitTmpCalled = true
	return m.err
}

func TestCommitter_Commit_AllowEmpty(t *testing.T) {
	mockClient := &mockCommitGitClient{}
	var buf bytes.Buffer
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"allow-empty"})
	if !mockClient.commitAllowEmptyCalled {
		t.Error("CommitAllowEmpty should be called")
	}
}

func TestCommitter_Commit_Tmp(t *testing.T) {
	mockClient := &mockCommitGitClient{}
	var buf bytes.Buffer
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"tmp"})
	if !mockClient.commitTmpCalled {
		t.Error("CommitTmp should be called")
	}
}

func TestCommitter_Commit_Help(t *testing.T) {
	var buf bytes.Buffer
	c := &Committer{
		gitClient:    &mockCommitGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"unknown"})

	output := buf.String()
	if output == "" || !strings.Contains(output, "Usage") {
		t.Errorf("Usage should be displayed, but got: %s", output)
	}
}

func TestCommitter_Commit_AllowEmpty_Error(t *testing.T) {
	var buf bytes.Buffer
	c := &Committer{
		gitClient:    &mockCommitGitClient{err: errors.New("fail")},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"allow-empty"})

	output := buf.String()
	if output != "Error: fail\n" {
		t.Errorf("unexpected output: got %q", output)
	}
}
