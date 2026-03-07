package git

import (
	"errors"
	"os/exec"
	"slices"
	"testing"
)

func TestClient_RevParseVerify(t *testing.T) {
	tests := []struct {
		name string
		ref  string
		err  error
		want bool
	}{
		{
			name: "success_valid_branch",
			ref:  "main",
			err:  nil,
			want: true,
		},
		{
			name: "success_valid_commit_hash",
			ref:  "1234567890abcdef",
			err:  nil,
			want: true,
		},
		{
			name: "success_valid_tag",
			ref:  "v1.0.0",
			err:  nil,
			want: true,
		},
		{
			name: "success_valid_head",
			ref:  "HEAD",
			err:  nil,
			want: true,
		},
		{
			name: "success_valid_relative_ref",
			ref:  "HEAD~1",
			err:  nil,
			want: true,
		},
		{
			name: "failure_invalid_ref",
			ref:  "nonexistent-branch",
			err:  errors.New("exit status 1"),
			want: false,
		},
		{
			name: "failure_invalid_commit_hash",
			ref:  "invalid-hash",
			err:  errors.New("fatal: ambiguous argument"),
			want: false,
		},
		{
			name: "failure_empty_ref",
			ref:  "",
			err:  errors.New("fatal: Needed a single revision"),
			want: false,
		},
		{
			name: "failure_not_git_repository",
			ref:  "main",
			err:  errors.New("fatal: not a git repository"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					expectedArgs := []string{"rev-parse", "--verify", "--quiet", tt.ref}
					if name != "git" || len(arg) != len(expectedArgs) {
						t.Errorf("unexpected command: %s %v", name, arg)
					}
					for i, a := range arg {
						if a != expectedArgs[i] {
							t.Errorf("unexpected arg[%d]: got %s, want %s", i, a, expectedArgs[i])
						}
					}
					return helperCommand(t, "", tt.err)
				},
			}

			got := c.RevParseVerify(tt.ref)
			if got != tt.want {
				t.Errorf("RevParseVerify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetCommitHash(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		err     error
		want    string
		wantErr bool
	}{
		{
			name:    "success_commit_hash",
			output:  "1234567",
			err:     nil,
			want:    "1234567",
			wantErr: false,
		},
		{
			name:    "success_longer_commit_hash",
			output:  "1234567890abcdef",
			err:     nil,
			want:    "1234567890abcdef",
			wantErr: false,
		},
		{
			name:    "success_with_whitespace",
			output:  "  1234567  \n",
			err:     nil,
			want:    "1234567",
			wantErr: false,
		},
		{
			name:    "error_not_git_repository",
			output:  "",
			err:     errors.New("fatal: not a git repository"),
			want:    "unknown",
			wantErr: false, // Should return "unknown" as fallback
		},
		{
			name:    "error_no_commits",
			output:  "",
			err:     errors.New("fatal: ambiguous argument 'HEAD'"),
			want:    "unknown",
			wantErr: false, // Should return "unknown" as fallback
		},
		{
			name:    "error_permission_denied",
			output:  "",
			err:     errors.New("permission denied"),
			want:    "unknown",
			wantErr: false, // Should return "unknown" as fallback
		},
		{
			name:    "success_empty_output_with_no_error",
			output:  "",
			err:     nil,
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					expectedArgs := []string{"rev-parse", "--short", "HEAD"}
					if name != "git" || len(arg) != len(expectedArgs) {
						t.Errorf("unexpected command: %s %v", name, arg)
					}
					for i, a := range arg {
						if a != expectedArgs[i] {
							t.Errorf("unexpected arg[%d]: got %s, want %s", i, a, expectedArgs[i])
						}
					}
					return helperCommand(t, tt.output, tt.err)
				},
			}

			got, err := c.GetCommitHash()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCommitHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCommitHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetUpstreamBranchName(t *testing.T) {
	tests := []struct {
		name           string
		branch         string
		expectedOutput string
		wantArgs       []string
	}{
		{
			name:           "get upstream for main",
			branch:         "main",
			expectedOutput: "origin/main",
			wantArgs:       []string{"git", "rev-parse", "--abbrev-ref", "main@{upstream}"},
		},
		{
			name:           "get upstream for feature branch",
			branch:         "feature/test",
			expectedOutput: "origin/feature/test",
			wantArgs:       []string{"git", "rev-parse", "--abbrev-ref", "feature/test@{upstream}"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotArgs []string
			client := &Client{
				execCommand: func(name string, args ...string) *exec.Cmd {
					gotArgs = append([]string{name}, args...)
					return exec.Command("echo", "-n", tt.expectedOutput)
				},
			}

			result, err := client.GetUpstreamBranchName(tt.branch)
			if err != nil {
				t.Errorf("GetUpstreamBranchName() error = %v", err)
			}

			if !slices.Equal(gotArgs, tt.wantArgs) {
				t.Errorf("GetUpstreamBranchName() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}

			if result != tt.expectedOutput {
				t.Errorf("GetUpstreamBranchName() result = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}

func TestClient_GetUpstreamBranchName_Error(t *testing.T) {
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			return exec.Command("false") // Command that always fails
		},
	}

	_, err := client.GetUpstreamBranchName("main")
	if err == nil {
		t.Error("Expected GetUpstreamBranchName to return an error")
	}
}
