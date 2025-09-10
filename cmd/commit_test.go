package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// mockGitClient for commit_test (minimal CommitWriter)
type mockCommitGitClient struct {
	commitAllowEmptyCalled       bool
	commitCalled                 bool
	commitAmendCalled            bool
	commitAmendNoEditCalled      bool
	commitAmendWithMessageCalled bool
	commitMessage                string
	amendMessage                 string
	err                          error
}

func (m *mockCommitGitClient) CommitAllowEmpty() error {
	m.commitAllowEmptyCalled = true
	return m.err
}

func (m *mockCommitGitClient) Commit(message string) error {
	m.commitCalled = true
	m.commitMessage = message
	return m.err
}

func (m *mockCommitGitClient) CommitAmend() error {
	m.commitAmendCalled = true
	return m.err
}
func (m *mockCommitGitClient) CommitAmendNoEdit() error {
	m.commitAmendNoEditCalled = true
	return m.err
}
func (m *mockCommitGitClient) CommitAmendWithMessage(message string) error {
	m.commitAmendWithMessageCalled = true
	m.amendMessage = message
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
	c.Commit([]string{"allow", "empty"})
	if !mockClient.commitAllowEmptyCalled {
		t.Error("CommitAllowEmpty should be called")
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
	c.Commit([]string{})

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
	c.Commit([]string{"allow", "empty"})

	output := buf.String()
	if output != "Error: fail\n" {
		t.Errorf("unexpected output: got %q", output)
	}
}

func TestCommitter_Commit_Normal(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockCommitGitClient{}
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"test message"})

	if !mockClient.commitCalled {
		t.Error("git commit command should be called with the correct message")
	}
	if mockClient.commitMessage != "test message" {
		t.Errorf("expected commit message 'test message', got '%s'", mockClient.commitMessage)
	}
}

func TestCommitter_Commit_Normal_WithBrackets(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockCommitGitClient{}
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"[update]", "test", "message"})

	if !mockClient.commitCalled {
		t.Error("git commit command should be called with the correct message including brackets")
	}
	expectedMessage := "[update] test message"
	if mockClient.commitMessage != expectedMessage {
		t.Errorf("expected commit message '%s', got '%s'", expectedMessage, mockClient.commitMessage)
	}
}

func TestCommitter_Commit_Normal_Error(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockCommitGitClient{err: errors.New("commit failed")}
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"test message"})

	output := buf.String()
	if !strings.Contains(output, "Error:") {
		t.Errorf("Expected error message, got: %s", output)
	}
}

func TestCommitter_Commit_Amend_WithMessage(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockCommitGitClient{}
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"amend", "updated message"})
	if !mockClient.commitAmendWithMessageCalled {
		t.Error("git commit --amend -m command should be called")
	}
}

func TestCommitter_Commit_Amend_NoEdit(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockCommitGitClient{}
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"amend", "no-edit"})
	if !mockClient.commitAmendNoEditCalled {
		t.Error("git commit --amend --no-edit command should be called")
	}
}

func TestCommitter_Commit_Amend_Error(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockCommitGitClient{
		err: errors.New("commit failed"),
	}
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"amend", "test", "message"})
	output := buf.String()
	if !strings.Contains(output, "Error:") {
		t.Errorf("Expected error message, got: %s", output)
	}
}

func TestCommitter_Commit_Normal_NoSpace(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockCommitGitClient{}
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"[update]hoge"})

	if !mockClient.commitCalled {
		t.Error("Expected commit invocation to manage spacing.")
	}
	if mockClient.commitMessage != "[update]hoge" {
		t.Errorf("expected commit message '[update]hoge', got '%s'", mockClient.commitMessage)
	}
}

func TestCommitter_Commit_Amend_WithMultiWordMessage(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockCommitGitClient{}
	c := &Committer{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"amend", "[update]", "message", "with", "spaces"})
	if !mockClient.commitAmendWithMessageCalled {
		t.Error("git commit --amend -m command should be called with the complete message")
	}
}
