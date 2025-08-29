package cmd

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v3/git"
)

// mockGitClient for commit_test
type mockCommitGitClient struct {
	git.Clienter
	commitAllowEmptyCalled bool
	err                    error
}

func (m *mockCommitGitClient) CommitAllowEmpty() error {
	m.commitAllowEmptyCalled = true
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

func TestCommitter_Commit_Normal_WithBrackets(t *testing.T) {
	var buf bytes.Buffer
	commandCalled := false
	c := &Committer{
		gitClient:    &mockCommitGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && len(arg) == 3 && arg[0] == "commit" && arg[1] == "-m" && arg[2] == "[update] test message" {
				commandCalled = true
			}
			return exec.Command("echo")
		},
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"[update]", "test", "message"})

	if !commandCalled {
		t.Error("git commit command should be called with the correct message including brackets")
	}
}

func TestCommitter_Commit_Normal_Error(t *testing.T) {
	var buf bytes.Buffer
	c := &Committer{
		gitClient:    &mockCommitGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("false") // command that fails
		},
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
	commandCalled := false
	c := &Committer{
		gitClient:    &mockCommitGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && len(arg) == 4 && arg[0] == "commit" &&
				arg[1] == "--amend" && arg[2] == "-m" && arg[3] == "updated message" {
				commandCalled = true
			}
			return exec.Command("echo")
		},
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"amend", "updated message"})
	if !commandCalled {
		t.Error("git commit --amend -m command should be called")
	}
}

func TestCommitter_Commit_Amend_NoEdit(t *testing.T) {
	var buf bytes.Buffer
	commandCalled := false
	c := &Committer{
		gitClient:    &mockCommitGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && len(arg) == 3 && arg[0] == "commit" && arg[1] == "--amend" && arg[2] == "--no-edit" {
				commandCalled = true
			}
			return exec.Command("echo")
		},
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"amend", "--no-edit"})
	if !commandCalled {
		t.Error("git commit --amend --no-edit command should be called")
	}
}

func TestCommitter_Commit_Amend_Error(t *testing.T) {
	var buf bytes.Buffer
	c := &Committer{
		gitClient:    &mockCommitGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("false") // command that fails
		},
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"amend", "test message"})
	output := buf.String()
	if !strings.Contains(output, "Error:") {
		t.Errorf("Expected error message, got: %s", output)
	}
}

func TestCommitter_Commit_Normal_NoSpace(t *testing.T) {
	var buf bytes.Buffer
	commandCalled := false
	c := &Committer{
		gitClient:    &mockCommitGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && len(arg) == 3 && arg[0] == "commit" && arg[1] == "-m" && arg[2] == "[update]hoge" {
				commandCalled = true
			}
			return exec.Command("echo")
		},
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"[update]hoge"})

	if !commandCalled {
		t.Error("git commit command should be called with the correct message without spaces")
	}
}

func TestCommitter_Commit_Amend_WithMultiWordMessage(t *testing.T) {
	var buf bytes.Buffer
	commandCalled := false
	c := &Committer{
		gitClient:    &mockCommitGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(name string, arg ...string) *exec.Cmd {
			if name == "git" && len(arg) == 4 && arg[0] == "commit" &&
				arg[1] == "--amend" && arg[2] == "-m" && arg[3] == "[update] message with spaces" {
				commandCalled = true
			}
			return exec.Command("echo")
		},
	}
	c.helper.outputWriter = &buf
	c.Commit([]string{"amend", "[update]", "message", "with", "spaces"})
	if !commandCalled {
		t.Error("git commit --amend -m command should be called with the complete message")
	}
}
