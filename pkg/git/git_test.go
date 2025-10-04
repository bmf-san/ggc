package git

import (
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient()

	if client == nil {
		t.Error("NewClient() should return a non-nil client")
		return
	}

	if client.execCommand == nil {
		t.Error("NewClient() should set execCommand field")
	}
}

func TestClient_GetGitStatus(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		err     error
		want    string
		wantErr bool
	}{
		{
			name:    "success_no_changes",
			output:  "",
			err:     nil,
			want:    "",
			wantErr: false,
		},
		{
			name:    "success_with_changes",
			output:  " M file.go\n?? new.go\n",
			err:     nil,
			want:    " M file.go\n?? new.go\n",
			wantErr: false,
		},
		{
			name:    "success_multiple_change_types",
			output:  " M modified.go\n A added.go\n D deleted.go\n?? untracked.go\nR  renamed.go -> new_name.go\n",
			err:     nil,
			want:    " M modified.go\n A added.go\n D deleted.go\n?? untracked.go\nR  renamed.go -> new_name.go\n",
			wantErr: false,
		},
		{
			name:    "error_git_command_failure",
			output:  "",
			err:     errors.New("not a git repository"),
			want:    "",
			wantErr: true,
		},
		{
			name:    "error_permission_denied",
			output:  "",
			err:     errors.New("permission denied"),
			want:    "",
			wantErr: true,
		},
		{
			name:    "success_staging_area_changes",
			output:  "M  staged.go\n M unstaged.go\n",
			err:     nil,
			want:    "M  staged.go\n M unstaged.go\n",
			wantErr: false,
		},
		{
			name:    "success_empty_repository",
			output:  "",
			err:     nil,
			want:    "",
			wantErr: false,
		},
		{
			name:    "error_repository_corruption",
			output:  "",
			err:     errors.New("fatal: bad object refs/heads/main"),
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					if name != "git" || !strings.Contains(strings.Join(arg, " "), "status --porcelain") {
						t.Errorf("unexpected command: %s %v", name, arg)
					}
					return helperCommand(t, tt.output, tt.err)
				},
			}

			got, err := c.GetGitStatus()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGitStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetGitStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetBranchName(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		err     error
		want    string
		wantErr bool
	}{
		{
			name:    "success_main_branch",
			output:  "main\n",
			err:     nil,
			want:    "main",
			wantErr: false,
		},
		{
			name:    "success_feature_branch",
			output:  "feature/test\n",
			err:     nil,
			want:    "feature/test",
			wantErr: false,
		},
		{
			name:    "success_deep_branch_hierarchy",
			output:  "feature/user/story/implementation\n",
			err:     nil,
			want:    "feature/user/story/implementation",
			wantErr: false,
		},
		{
			name:    "success_trim_whitespace",
			output:  "  main  \n\n",
			err:     nil,
			want:    "main",
			wantErr: false,
		},
		{
			name:    "error_git_command_failure",
			output:  "",
			err:     errors.New("not a git repository"),
			want:    "",
			wantErr: true,
		},
		{
			name:    "error_head_not_found",
			output:  "",
			err:     errors.New("fatal: ref HEAD is not a symbolic ref"),
			want:    "",
			wantErr: true,
		},
		{
			name:    "success_empty_output",
			output:  "",
			err:     nil,
			want:    "",
			wantErr: false,
		},
		{
			name:    "error_detached_head",
			output:  "",
			err:     errors.New("fatal: HEAD does not point to a branch"),
			want:    "",
			wantErr: true,
		},
		{
			name:    "success_branch_with_numbers",
			output:  "release/v1.2.3\n",
			err:     nil,
			want:    "release/v1.2.3",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					if name != "git" || !strings.Contains(strings.Join(arg, " "), "rev-parse --abbrev-ref HEAD") {
						t.Errorf("unexpected command: %s %v", name, arg)
					}
					return helperCommand(t, tt.output, tt.err)
				},
			}

			got, err := c.GetBranchName()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBranchName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBranchName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_CheckoutNewBranch(t *testing.T) {
	tests := []struct {
		name    string
		branch  string
		err     error
		wantErr bool
	}{
		{
			name:    "success_new_branch",
			branch:  "feature/test",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "success_deep_hierarchy_branch",
			branch:  "feature/user/story/implementation",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "error_existing_branch",
			branch:  "main",
			err:     errors.New("fatal: A branch named 'main' already exists"),
			wantErr: true,
		},
		{
			name:    "error_invalid_branch_name",
			branch:  "invalid..branch",
			err:     errors.New("fatal: 'invalid..branch' is not a valid branch name"),
			wantErr: true,
		},
		{
			name:    "error_empty_branch_name",
			branch:  "",
			err:     errors.New("fatal: branch name required"),
			wantErr: true,
		},
		{
			name:    "error_permission_denied",
			branch:  "feature/test",
			err:     errors.New("permission denied"),
			wantErr: true,
		},
		{
			name:    "error_disk_space_full",
			branch:  "feature/test",
			err:     errors.New("fatal: unable to write new index file"),
			wantErr: true,
		},
		{
			name:    "success_branch_with_special_chars",
			branch:  "feature/user-story_123",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "error_branch_name_too_long",
			branch:  strings.Repeat("very-long-branch-name", 20),
			err:     errors.New("fatal: branch name too long"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					if name != "git" || !strings.Contains(strings.Join(arg, " "), "checkout -b "+tt.branch) {
						t.Errorf("unexpected command: %s %v", name, arg)
					}
					return helperCommand(t, "", tt.err)
				},
			}

			if err := c.CheckoutNewBranch(tt.branch); (err != nil) != tt.wantErr {
				t.Errorf("CheckoutNewBranch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_GetCurrentBranch(t *testing.T) {
	tests := []struct {
		name        string
		output      string
		err         error
		mockFunc    func() (string, error)
		want        string
		wantErr     bool
		useMockFunc bool
	}{
		{
			name:    "success_main_branch",
			output:  "main\n",
			err:     nil,
			want:    "main",
			wantErr: false,
		},
		{
			name:    "success_feature_branch",
			output:  "feature/test\n",
			err:     nil,
			want:    "feature/test",
			wantErr: false,
		},
		{
			name:    "success_trim_whitespace",
			output:  "  develop  \n\n",
			err:     nil,
			want:    "develop",
			wantErr: false,
		},
		{
			name:    "error_git_command_failure",
			output:  "",
			err:     errors.New("not a git repository"),
			want:    "",
			wantErr: true,
		},
		{
			name:    "error_detached_head",
			output:  "",
			err:     errors.New("fatal: ref HEAD is not a symbolic ref"),
			want:    "",
			wantErr: true,
		},
		{
			name:        "success_with_mock_function",
			mockFunc:    func() (string, error) { return "mocked-branch", nil },
			want:        "mocked-branch",
			wantErr:     false,
			useMockFunc: true,
		},
		{
			name:        "error_mock_function_failure",
			mockFunc:    func() (string, error) { return "", errors.New("mock error") },
			want:        "",
			wantErr:     true,
			useMockFunc: true,
		},
		{
			name:    "success_release_branch",
			output:  "release/v2.1.0\n",
			err:     nil,
			want:    "release/v2.1.0",
			wantErr: false,
		},
		{
			name:    "error_corrupted_head",
			output:  "",
			err:     errors.New("fatal: bad object HEAD"),
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					if name != "git" || !strings.Contains(strings.Join(arg, " "), "rev-parse --abbrev-ref HEAD") {
						t.Errorf("unexpected command: %s %v", name, arg)
					}
					return helperCommand(t, tt.output, tt.err)
				},
			}

			if tt.useMockFunc {
				c.GetCurrentBranchFunc = tt.mockFunc
			}

			got, err := c.GetCurrentBranch()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCurrentBranch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCurrentBranch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_LogGraph(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr bool
	}{
		{
			name:    "success_log_display",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "error_git_command_failure",
			err:     errors.New("not a git repository"),
			wantErr: true,
		},
		{
			name:    "error_permission_denied",
			err:     errors.New("permission denied"),
			wantErr: true,
		},
		{
			name:    "error_empty_repository",
			err:     errors.New("fatal: your current branch does not have any commits yet"),
			wantErr: true,
		},
		{
			name:    "error_corrupted_repository",
			err:     errors.New("fatal: bad object"),
			wantErr: true,
		},
		{
			name:    "error_no_commits",
			err:     errors.New("fatal: bad default revision 'HEAD'"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					if name != "git" || !strings.Contains(strings.Join(arg, " "), "log --graph --oneline --decorate --all") {
						t.Errorf("unexpected command: %s %v", name, arg)
					}
					return helperCommand(t, "", tt.err)
				},
			}

			if err := c.LogGraph(); (err != nil) != tt.wantErr {
				t.Errorf("LogGraph() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestClient_ErrorHandling tests error handling with various git command failures
func TestClient_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		expectError bool
		errorType   string
	}{
		{
			name:        "GetGitStatus_with_command_failure",
			method:      "GetGitStatus",
			expectError: true,
			errorType:   "repository_error",
		},
		{
			name:        "GetBranchName_with_command_failure",
			method:      "GetBranchName",
			expectError: true,
			errorType:   "branch_error",
		},
		{
			name:        "GetCurrentBranch_with_command_failure",
			method:      "GetCurrentBranch",
			expectError: true,
			errorType:   "branch_error",
		},
		{
			name:        "CheckoutNewBranch_with_command_failure",
			method:      "CheckoutNewBranch",
			expectError: true,
			errorType:   "checkout_error",
		},
		{
			name:        "LogGraph_with_command_failure",
			method:      "LogGraph",
			expectError: true,
			errorType:   "log_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(_ string, _ ...string) *exec.Cmd {
					return helperCommand(t, "", errors.New("command failed: "+tt.errorType))
				},
			}

			var err error
			switch tt.method {
			case "GetGitStatus":
				_, err = c.GetGitStatus()
			case "GetBranchName":
				_, err = c.GetBranchName()
			case "GetCurrentBranch":
				_, err = c.GetCurrentBranch()
			case "CheckoutNewBranch":
				err = c.CheckoutNewBranch("test-branch")
			case "LogGraph":
				err = c.LogGraph()
			}

			if tt.expectError && err == nil {
				t.Errorf("%s() expected error but got nil", tt.method)
			}
			if !tt.expectError && err != nil {
				t.Errorf("%s() expected no error but got: %v", tt.method, err)
			}
		})
	}
}

// TestClient_CommandValidation tests that the correct git commands are called
func TestClient_CommandValidation(t *testing.T) {
	tests := []struct {
		name            string
		method          string
		args            []any
		expectedCommand string
	}{
		{
			name:            "GetGitStatus_command_validation",
			method:          "GetGitStatus",
			expectedCommand: "git status --porcelain",
		},
		{
			name:            "GetBranchName_command_validation",
			method:          "GetBranchName",
			expectedCommand: "git rev-parse --abbrev-ref HEAD",
		},
		{
			name:            "GetCurrentBranch_command_validation",
			method:          "GetCurrentBranch",
			expectedCommand: "git rev-parse --abbrev-ref HEAD",
		},
		{
			name:            "CheckoutNewBranch_command_validation",
			method:          "CheckoutNewBranch",
			args:            []any{"test-branch"},
			expectedCommand: "git checkout -b test-branch",
		},
		{
			name:            "LogGraph_command_validation",
			method:          "LogGraph",
			expectedCommand: "git log --graph --oneline --decorate --all",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commandCalled := ""
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					commandCalled = name + " " + strings.Join(arg, " ")
					return helperCommand(t, "test-output", nil)
				},
			}

			switch tt.method {
			case "GetGitStatus":
				_, _ = c.GetGitStatus()
			case "GetBranchName":
				_, _ = c.GetBranchName()
			case "GetCurrentBranch":
				_, _ = c.GetCurrentBranch()
			case "CheckoutNewBranch":
				if len(tt.args) > 0 {
					_ = c.CheckoutNewBranch(tt.args[0].(string))
				}
			case "LogGraph":
				_ = c.LogGraph()
			}

			if commandCalled != tt.expectedCommand {
				t.Errorf("%s() called command %q, expected %q", tt.method, commandCalled, tt.expectedCommand)
			}
		})
	}
}

// TestClient_EdgeCases tests edge cases and boundary conditions
func TestClient_EdgeCases(t *testing.T) {
	t.Run("empty_branch_name_in_CheckoutNewBranch", func(t *testing.T) {
		c := &Client{
			execCommand: func(_ string, _ ...string) *exec.Cmd {
				return helperCommand(t, "", errors.New("invalid branch name"))
			},
		}
		err := c.CheckoutNewBranch("")
		if err == nil {
			t.Error("CheckoutNewBranch('') should return error")
		}
	})

	t.Run("very_long_branch_name", func(t *testing.T) {
		longBranch := strings.Repeat("a", 255)
		c := &Client{
			execCommand: func(_ string, _ ...string) *exec.Cmd {
				return helperCommand(t, "", nil)
			},
		}
		err := c.CheckoutNewBranch(longBranch)
		if err != nil {
			t.Errorf("CheckoutNewBranch with long name failed: %v", err)
		}
	})

	t.Run("special_characters_in_branch_name", func(t *testing.T) {
		specialBranch := "feature/user-story_123"
		c := &Client{
			execCommand: func(_ string, _ ...string) *exec.Cmd {
				return helperCommand(t, "", nil)
			},
		}
		err := c.CheckoutNewBranch(specialBranch)
		if err != nil {
			t.Errorf("CheckoutNewBranch with special characters failed: %v", err)
		}
	})

	t.Run("unicode_in_branch_name", func(t *testing.T) {
		unicodeBranch := "feature/test-branch"
		c := &Client{
			execCommand: func(_ string, _ ...string) *exec.Cmd {
				return helperCommand(t, "", nil)
			},
		}
		err := c.CheckoutNewBranch(unicodeBranch)
		if err != nil {
			t.Errorf("CheckoutNewBranch with unicode characters failed: %v", err)
		}
	})

	t.Run("whitespace_only_output", func(t *testing.T) {
		c := &Client{
			execCommand: func(_ string, _ ...string) *exec.Cmd {
				return helperCommand(t, "   \n\t\n   ", nil)
			},
		}
		branch, err := c.GetBranchName()
		if err != nil {
			t.Errorf("GetBranchName with whitespace failed: %v", err)
		}
		if branch != "" {
			t.Errorf("Expected empty string after trimming whitespace, got: %q", branch)
		}
	})
}

// TestClient_MockFunctionBehavior tests the behavior of mock functions
func TestClient_MockFunctionBehavior(t *testing.T) {
	t.Run("mock_function_overrides_command", func(t *testing.T) {
		c := &Client{
			execCommand: func(_ string, _ ...string) *exec.Cmd {
				t.Error("execCommand should not be called when mock function is set")
				return helperCommand(t, "", errors.New("should not be called"))
			},
			GetCurrentBranchFunc: func() (string, error) {
				return "mocked-branch", nil
			},
		}

		branch, err := c.GetCurrentBranch()
		if err != nil {
			t.Errorf("GetCurrentBranch with mock failed: %v", err)
		}
		if branch != "mocked-branch" {
			t.Errorf("Expected 'mocked-branch', got: %q", branch)
		}
	})

	t.Run("mock_function_nil_check", func(t *testing.T) {
		c := &Client{
			execCommand: func(_ string, _ ...string) *exec.Cmd {
				return helperCommand(t, "real-branch", nil)
			},
			GetCurrentBranchFunc: nil,
		}

		branch, err := c.GetCurrentBranch()
		if err != nil {
			t.Errorf("GetCurrentBranch without mock failed: %v", err)
		}
		if branch != "real-branch" {
			t.Errorf("Expected 'real-branch', got: %q", branch)
		}
	})
}

// helperCommand creates a mock command for testing
func helperCommand(t *testing.T, output string, err error) *exec.Cmd {
	if t != nil {
		t.Helper()
	}
	if err != nil {
		return exec.Command("false")
	}
	if output == "" {
		return exec.Command("true")
	}
	return exec.Command("echo", "-n", output)
}

// fakeExecCommand creates a mock command that returns specified output
func fakeExecCommand(output string) *exec.Cmd {
	return exec.Command("echo", "-n", output)
}
