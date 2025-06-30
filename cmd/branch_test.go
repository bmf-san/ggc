package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"testing"
)

// mockBranchGitClient is a mock implementation of git.Clienter for branch tests
type mockBranchGitClient struct {
	mockGitClient
	getCurrentBranchCalled bool
	currentBranch          string
	err                    error
}

func (m *mockBranchGitClient) GetCurrentBranch() (string, error) {
	m.getCurrentBranchCalled = true
	return m.currentBranch, m.err
}

func (m *mockBranchGitClient) ListLocalBranches() ([]string, error) {
	return []string{"main", "feature/test"}, nil
}

func (m *mockBranchGitClient) LogGraph() error {
	return nil
}

func (m *mockBranchGitClient) StashPullPop() error {
	return nil
}

func TestBrancher_Branch_Current(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		currentBranch: "feature/test",
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
	}
	brancher.Branch([]string{"current"})

	if !mockClient.getCurrentBranchCalled {
		t.Error("GetCurrentBranch should be called")
	}

	output := buf.String()
	if output != "feature/test\n" {
		t.Errorf("unexpected output: got %q, want %q", output, "feature/test\n")
	}
}

func TestBrancher_Branch_Current_Error(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		err: errors.New("failed to get current branch"),
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
	}
	brancher.Branch([]string{"current"})

	if !mockClient.getCurrentBranchCalled {
		t.Error("GetCurrentBranch should be called")
	}

	output := buf.String()
	if output != "Error: failed to get current branch\n" {
		t.Errorf("unexpected output: got %q, want %q", output, "Error: failed to get current branch\n")
	}
}

func TestBrancher_Branch_Checkout(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		execCommand:  func(_ string, _ ...string) *exec.Cmd { return exec.Command("echo") },
		inputReader:  bufio.NewReader(strings.NewReader("1\n")),
	}
	brancher.Branch([]string{"checkout"})

	output := buf.String()
	expected := "Local branches:\n[1] main\n[2] feature/test\nEnter the number to checkout: "
	if !strings.Contains(output, expected) {
		t.Errorf("unexpected output: got %q, want %q", output, expected)
	}
}
