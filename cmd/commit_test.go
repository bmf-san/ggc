package cmd

import (
	"bytes"
	"errors"
	"os/exec"
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
		execCommand:  exec.Command,
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
		execCommand:  exec.Command,
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
		execCommand:  exec.Command,
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
		execCommand:  exec.Command,
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"allow-empty"})

	output := buf.String()
	if output != "Error: fail\n" {
		t.Errorf("unexpected output: got %q", output)
	}
}

func TestCommitter_Commit_Normal(t *testing.T) {
	var buf bytes.Buffer
	commandCalled := false
	c := &Committer{
		gitClient:    &mockCommitGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && len(arg) == 3 && arg[0] == "commit" && arg[1] == "-m" && arg[2] == "test message" {
				commandCalled = true
			}
			return exec.Command("echo")
		},
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"test message"})

	if !commandCalled {
		t.Error("git commit command should be called with the correct message")
	}
}

func TestCommitter_Commit_Normal_Error(t *testing.T) {
	var buf bytes.Buffer
	c := &Committer{
		gitClient:    &mockCommitGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, _  ...string) *exec.Cmd {
			return exec.Command("false") // 失敗するコマンド
		},
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"test message"})

	output := buf.String()
	if !strings.Contains(output, "Error:") {
		t.Errorf("Expected error message, got: %s", output)
	}
}

func TestCommitter_Commit_Tmp_Error(t *testing.T) {
	var buf bytes.Buffer
	c := &Committer{
		gitClient:    &mockCommitGitClient{err: errors.New("tmp commit failed")},
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand:  exec.Command,
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"tmp"})

	output := buf.String()
	if output != "Error: tmp commit failed\n" {
		t.Errorf("Expected tmp error message, got: %q", output)
	}
}

func TestShowCommitHelp(t *testing.T) {
	// ShowCommitHelp関数が呼び出せることを確認
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ShowCommitHelp should not panic: %v", r)
		}
	}()

	ShowCommitHelp()
	// 関数が正常に実行されることを確認
}
