package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
)

type mockCompleteGitClient struct {
	mockGitClient
	branches []string
	err      error
}

func (m *mockCompleteGitClient) GetBranchName() (string, error) {
	return "main", nil
}

func (m *mockCompleteGitClient) GetGitStatus() (string, error) {
	return "", nil
}

func (m *mockCompleteGitClient) ListLocalBranches() ([]string, error) {
	return m.branches, m.err
}

func TestCompleter_Complete_BranchSubcommands(t *testing.T) {
	completer := &Completer{
		gitClient: &mockCompleteGitClient{},
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	completer.Complete([]string{"branch"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("buf.ReadFrom failed: %v", err)
	}
	os.Stdout = oldStdout

	output := buf.String()
	for _, want := range []string{"current", "checkout", "checkout-remote", "create", "delete", "delete-merged"} {
		if !bytes.Contains([]byte(output), []byte(want)) {
			t.Errorf("branchサブコマンド候補 %s が出力されていません: %s", want, output)
		}
	}
}

func TestCompleter_Complete_BranchNames(t *testing.T) {
	completer := &Completer{
		gitClient: &mockCompleteGitClient{
			branches: []string{"main", "feature/test"},
		},
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	completer.Complete([]string{"branch", "foo"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("buf.ReadFrom failed: %v", err)
	}
	os.Stdout = oldStdout

	output := buf.String()
	for _, want := range []string{"main", "feature/test"} {
		if !bytes.Contains([]byte(output), []byte(want)) {
			t.Errorf("ローカルブランチ候補 %s が出力されていません: %s", want, output)
		}
	}
}

func TestCompleter_Complete_Files(t *testing.T) {
	completer := &Completer{
		gitClient: &mockCompleteGitClient{},
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("echo", "foo.txt\nbar.txt")
		},
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	completer.Complete([]string{"files"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("buf.ReadFrom failed: %v", err)
	}
	os.Stdout = oldStdout

	output := buf.String()
	for _, want := range []string{"foo.txt", "bar.txt"} {
		if !bytes.Contains([]byte(output), []byte(want)) {
			t.Errorf("filesサブコマンドの出力が想定と異なります: %s", output)
		}
	}
}

func TestCompleter_Complete_BranchList_Error(t *testing.T) {
	completer := &Completer{
		gitClient: &mockCompleteGitClient{
			err: os.ErrNotExist,
		},
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	completer.Complete([]string{"branch", "any"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("buf.ReadFrom failed: %v", err)
	}
	os.Stdout = oldStdout

	output := buf.String()
	if output != "" {
		t.Errorf("エラー時は出力されない想定: %s", output)
	}
}

func TestCompleter_Complete_Files_Error(t *testing.T) {
	completer := &Completer{
		gitClient: &mockCompleteGitClient{},
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("false")
		},
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	completer.Complete([]string{"files"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("buf.ReadFrom failed: %v", err)
	}
	os.Stdout = oldStdout

	output := buf.String()
	if output != "" {
		t.Errorf("エラー時は出力されない想定: %s", output)
	}
}

func TestCompleter_Complete_Default(t *testing.T) {
	completer := &Completer{
		gitClient: &mockCompleteGitClient{},
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	completer.Complete([]string{"unknown"})

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() failed: %v", err)
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("buf.ReadFrom failed: %v", err)
	}
	os.Stdout = oldStdout

	output := buf.String()
	if output != "" {
		t.Errorf("未対応分岐は出力なし想定: %s", output)
	}
}
