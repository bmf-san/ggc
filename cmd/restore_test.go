package cmd

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestRestoreer_Restore(t *testing.T) {
	cases := []struct {
		name           string
		args           []string
		expectedCmds   []string
		mockOutput     []byte
		mockError      error
		expectedOutput string
	}{
		{
			name:           "restore single file",
			args:           []string{"file.txt"},
			expectedCmds:   []string{"git restore file.txt"},
			mockOutput:     []byte("Restored file.txt"),
			mockError:      nil,
			expectedOutput: "",
		},
		{
			name:           "restore all files",
			args:           []string{"."},
			expectedCmds:   []string{"git restore ."},
			mockOutput:     []byte("Restored all files"),
			mockError:      nil,
			expectedOutput: "",
		},
		{
			name:           "restore multiple files",
			args:           []string{"file1.txt", "file2.txt"},
			expectedCmds:   []string{"git restore file1.txt file2.txt"},
			mockOutput:     []byte("Restored files"),
			mockError:      nil,
			expectedOutput: "",
		},
		{
			name:           "restore staged file",
			args:           []string{"staged", "file.txt"},
			expectedCmds:   []string{"git restore --staged file.txt"},
			mockOutput:     []byte("Unstaged file.txt"),
			mockError:      nil,
			expectedOutput: "",
		},
		{
			name:           "restore staged all files",
			args:           []string{"staged", "."},
			expectedCmds:   []string{"git restore --staged ."},
			mockOutput:     []byte("Unstaged all files"),
			mockError:      nil,
			expectedOutput: "",
		},
		{
			name:           "restore staged multiple files",
			args:           []string{"staged", "file1.txt", "file2.txt"},
			expectedCmds:   []string{"git restore --staged file1.txt file2.txt"},
			mockOutput:     []byte("Unstaged files"),
			mockError:      nil,
			expectedOutput: "",
		},
		{
			name:           "restore from commit",
			args:           []string{"HEAD~1", "file.txt"},
			expectedCmds:   []string{"git restore --source HEAD~1 file.txt"},
			mockOutput:     []byte("Restored from commit"),
			mockError:      nil,
			expectedOutput: "",
		},
		{
			name:           "restore from commit hash",
			args:           []string{"abc123f", "file.txt"},
			expectedCmds:   []string{"git restore --source abc123f file.txt"},
			mockOutput:     []byte("Restored from commit"),
			mockError:      nil,
			expectedOutput: "",
		},
		{
			name:           "restore file error",
			args:           []string{"file.txt"},
			expectedCmds:   []string{"git restore file.txt"},
			mockOutput:     nil,
			mockError:      errors.New("restore failed"),
			expectedOutput: "Error:",
		},
		{
			name:           "restore staged error",
			args:           []string{"staged", "file.txt"},
			expectedCmds:   []string{"git restore --staged file.txt"},
			mockOutput:     nil,
			mockError:      errors.New("staged restore failed"),
			expectedOutput: "Error:",
		},
		{
			name:           "restore from commit error",
			args:           []string{"HEAD~1", "file.txt"},
			expectedCmds:   []string{"git restore --source HEAD~1 file.txt"},
			mockOutput:     nil,
			mockError:      errors.New("commit restore failed"),
			expectedOutput: "Error:",
		},
		{
			name:           "no args",
			args:           []string{},
			expectedCmds:   nil,
			mockOutput:     nil,
			mockError:      nil,
			expectedOutput: "Usage: ggc restore",
		},
		{
			name:           "staged without file",
			args:           []string{"staged"},
			expectedCmds:   nil,
			mockOutput:     nil,
			mockError:      nil,
			expectedOutput: "Usage: ggc restore",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			cmdIndex := 0

			MockGitClient := &MockGitClient{
				execCommand: func(_ string, args ...string) *exec.Cmd {
					if cmdIndex < len(tc.expectedCmds) {
						gotCmd := strings.Join(append([]string{"git"}, args...), " ")
						if gotCmd != tc.expectedCmds[cmdIndex] {
							t.Errorf("expected command %q, got %q", tc.expectedCmds[cmdIndex], gotCmd)
						}
					}
					cmdIndex++
					if tc.mockError != nil {
						return exec.Command("false")
					}
					return exec.Command("echo", string(tc.mockOutput))
				},
			}

			r := &Restoreer{
				outputWriter: &buf,
				helper:       NewHelper(),
				execCommand:  MockGitClient.execCommand,
				gitClient:    MockGitClient,
			}
			r.helper.outputWriter = &buf

			r.Restore(tc.args)

			output := buf.String()
			if !strings.Contains(output, tc.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tc.expectedOutput, output)
			}
		})
	}
}

// MockGitClient that uses execCommand for testing
type MockGitClient struct {
	execCommand func(string, ...string) *exec.Cmd
}

func (m *MockGitClient) RestoreWorkingDir(paths ...string) error {
	args := append([]string{"restore"}, paths...)
	cmd := m.execCommand("git", args...)
	return cmd.Run()
}

func (m *MockGitClient) RestoreStaged(paths ...string) error {
	args := append([]string{"restore", "--staged"}, paths...)
	cmd := m.execCommand("git", args...)
	return cmd.Run()
}

func (m *MockGitClient) RestoreFromCommit(commit string, paths ...string) error {
	args := append([]string{"restore", "--source", commit}, paths...)
	cmd := m.execCommand("git", args...)
	return cmd.Run()
}

func (m *MockGitClient) GetCurrentBranch() (string, error) {
	return "main", nil
}

func (m *MockGitClient) ListLocalBranches() ([]string, error) {
	return []string{"main", "feature/test"}, nil
}

func (m *MockGitClient) ListRemoteBranches() ([]string, error) {
	return []string{"origin/main", "origin/feature/test"}, nil
}

func (m *MockGitClient) Push(bool) error {
	return nil
}

func (m *MockGitClient) Pull(bool) error {
	return nil
}

func (m *MockGitClient) FetchPrune() error {
	return nil
}

func (m *MockGitClient) LogSimple() error {
	return nil
}

func (m *MockGitClient) LogGraph() error {
	return nil
}

func (m *MockGitClient) CommitAllowEmpty() error {
	return nil
}

func (m *MockGitClient) ResetHardAndClean() error {
	return nil
}

func (m *MockGitClient) CleanFiles() error {
	return nil
}

func (m *MockGitClient) CleanDirs() error {
	return nil
}

func (m *MockGitClient) GetBranchName() (string, error) {
	return "main", nil
}

func (m *MockGitClient) GetGitStatus() (string, error) {
	return "", nil
}

func (m *MockGitClient) CheckoutNewBranch(_ string) error {
	return nil
}

func (m *MockGitClient) RestoreAll() error {
	return nil
}

func (m *MockGitClient) RestoreAllStaged() error {
	return nil
}
