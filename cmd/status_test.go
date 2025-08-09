package cmd

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func mockExecCommand(command string, args ...string) *exec.Cmd {
	if command == "git" && len(args) >= 3 {
		if args[0] == "rev-parse" && args[1] == "--abbrev-ref" && args[2] == "HEAD" {
			return exec.Command("echo", "main")
		} else if args[0] == "rev-parse" && args[1] == "--abbrev-ref" && strings.Contains(args[2], "@{upstream}") {
			return exec.Command("echo", "origin/main")
		} else if args[0] == "rev-list" && args[1] == "--left-right" && args[2] == "--count" {
			return exec.Command("echo", "0 0")
		}
	}

	return exec.Command("echo", "mock output")
}

func TestStatuseer_getCurrentBranch(t *testing.T) {
	tests := []struct {
		name        string
		execCommand func(string, ...string) *exec.Cmd
		expected    string
		expectError bool
	}{
		{
			name: "successful branch retrieval",
			execCommand: func(command string, args ...string) *exec.Cmd {
				if command == "git" && len(args) >= 3 && args[0] == "rev-parse" && args[1] == "--abbrev-ref" && args[2] == "HEAD" {
					return exec.Command("echo", "main")
				}
				return exec.Command("echo", "mock")
			},
			expected:    "main",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Statuseer{
				execCommand: tt.execCommand,
				gitClient:   &mockGitClient{},
			}

			branch, err := s.gitClient.GetCurrentBranch()

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if branch != tt.expected {
				t.Errorf("Expected branch '%s', got '%s'", tt.expected, branch)
			}
		})
	}
}

func TestStatuseer_getUpstreamStatus(t *testing.T) {
	tests := []struct {
		name        string
		branch      string
		execCommand func(string, ...string) *exec.Cmd
		expected    string
	}{
		{
			name:   "up to date branch",
			branch: "main",
			execCommand: func(command string, args ...string) *exec.Cmd {
				if command == "git" && len(args) >= 3 {
					if args[0] == "rev-parse" && args[1] == "--abbrev-ref" && strings.Contains(args[2], "@{upstream}") {
						return exec.Command("echo", "origin/main")
					}
					if args[0] == "rev-list" && args[1] == "--left-right" && args[2] == "--count" {
						return exec.Command("echo", "0 0")
					}
				}
				return exec.Command("echo", "mock")
			},
			expected: "Your branch is up to date with 'origin/main'",
		},
		{
			name:   "branch ahead of upstream",
			branch: "feature",
			execCommand: func(command string, args ...string) *exec.Cmd {
				if command == "git" && len(args) >= 3 {
					if args[0] == "rev-parse" && args[1] == "--abbrev-ref" && strings.Contains(args[2], "@{upstream}") {
						return exec.Command("echo", "origin/feature")
					}
					if args[0] == "rev-list" && args[1] == "--left-right" && args[2] == "--count" {
						return exec.Command("echo", "2 0")
					}
				}
				return exec.Command("echo", "mock")
			},
			expected: "Your branch is ahead of 'origin/feature' by 2 commit(s)",
		},
		{
			name:   "branch behind upstream",
			branch: "feature",
			execCommand: func(command string, args ...string) *exec.Cmd {
				if command == "git" && len(args) >= 3 {
					if args[0] == "rev-parse" && args[1] == "--abbrev-ref" && strings.Contains(args[2], "@{upstream}") {
						return exec.Command("echo", "origin/feature")
					}
					if args[0] == "rev-list" && args[1] == "--left-right" && args[2] == "--count" {
						return exec.Command("echo", "0 3")
					}
				}
				return exec.Command("echo", "mock")
			},
			expected: "Your branch is behind 'origin/feature' by 3 commit(s)",
		},
		{
			name:   "diverged branches",
			branch: "feature",
			execCommand: func(command string, args ...string) *exec.Cmd {
				if command == "git" && len(args) >= 3 {
					if args[0] == "rev-parse" && args[1] == "--abbrev-ref" && strings.Contains(args[2], "@{upstream}") {
						return exec.Command("echo", "origin/feature")
					}
					if args[0] == "rev-list" && args[1] == "--left-right" && args[2] == "--count" {
						return exec.Command("echo", "2 1")
					}
				}
				return exec.Command("echo", "mock")
			},
			expected: "Your branch and 'origin/feature' have diverged",
		},
		{
			name:   "no upstream branch",
			branch: "local-branch",
			execCommand: func(command string, args ...string) *exec.Cmd {
				if command == "git" && len(args) >= 3 {
					if args[0] == "rev-parse" && args[1] == "--abbrev-ref" && strings.Contains(args[2], "@{upstream}") {
						return exec.Command("false") // Command that returns error
					}
				}
				return exec.Command("echo", "mock")
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Statuseer{
				execCommand: tt.execCommand,
				gitClient:   &mockGitClient{},
			}

			status := s.getUpstreamStatus(tt.branch)

			if tt.expected == "" && status != "" {
				t.Errorf("Expected empty status, got '%s'", status)
			}
			if tt.expected != "" && status != tt.expected {
				t.Errorf("Expected status '%s', got '%s'", tt.expected, status)
			}
		})
	}
}

func TestStatuseer_Status(t *testing.T) {
	cases := []struct {
		name           string
		args           []string
		expectedCmds   []string
		mockOutput     []byte
		mockError      error
		expectedOutput string
	}{
		{
			name:           "status no args",
			args:           []string{},
			expectedCmds:   []string{"git -c color.status=always status", "git rev-parse --abbrev-ref main@{upstream}", "git rev-list --left-right --count main...On branch main\nChanges not staged for commit:\n  modified:   modified_file.go\n\nUntracked files:\n  untracked_file.go"},
			mockOutput:     []byte("On branch main\nChanges not staged for commit:\n  modified:   modified_file.go\n\nUntracked files:\n  untracked_file.go\n"),
			mockError:      nil,
			expectedOutput: "On branch main",
		},
		{
			name:           "status short",
			args:           []string{"short"},
			expectedCmds:   []string{"git -c color.status=always status --short", "git rev-parse --abbrev-ref main@{upstream}", "git rev-list --left-right --count main...M  modified_file.go\n?? untracked_file.go"},
			mockOutput:     []byte("M  modified_file.go\n?? untracked_file.go\n"),
			mockError:      nil,
			expectedOutput: "M  modified_file.go",
		},
		{
			name:           "invalid command",
			args:           []string{"invalid"},
			expectedCmds:   nil,
			mockOutput:     nil,
			mockError:      nil,
			expectedOutput: "Usage: ggc status",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			cmdIndex := 0
			s := &Statuseer{
				outputWriter: &buf,
				helper:       NewHelper(),
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
				gitClient: &mockGitClient{},
			}
			s.helper.outputWriter = &buf
			s.Status(tc.args)
			output := buf.String()
			if !strings.Contains(output, tc.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tc.expectedOutput, output)
			}
		})
	}
}

func TestStatuseer_StatusWithBranchInfo(t *testing.T) {
	var output bytes.Buffer
	s := &Statuseer{
		outputWriter: &output,
		helper:       NewHelper(),
		execCommand:  mockExecCommand,
		gitClient:    &mockGitClient{},
	}

	s.Status([]string{})

	outputStr := output.String()

	if !strings.Contains(outputStr, "On branch main") {
		t.Error("Expected 'On branch main' in output")
	}

	if !strings.Contains(outputStr, "Your branch is up to date") {
		t.Error("Expected upstream status in output")
	}
}
