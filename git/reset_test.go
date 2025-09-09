package git

import (
	"errors"
	"os/exec"
	"reflect"
	"testing"
)

func TestClient_ResetHardAndClean(t *testing.T) {
	var gotArgs [][]string
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append(gotArgs, append([]string{name}, args...))
			return exec.Command("echo")
		},
		GetCurrentBranchFunc: func() (string, error) {
			return "main", nil
		},
	}

	_ = client.ResetHardAndClean()
	want := [][]string{
		{"git", "reset", "--hard", "origin/main"},
		{"git", "clean", "-fdx"},
	}
	if !reflect.DeepEqual(gotArgs, want) {
		t.Errorf("got %v, want %v", gotArgs, want)
	}
}

func TestClient_ResetHard(t *testing.T) {
	tests := []struct {
		name    string
		commit  string
		err     error
		wantErr bool
	}{
		{
			name:    "success_reset_to_commit_hash",
			commit:  "abc123def456",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "success_reset_to_head",
			commit:  "HEAD",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "success_reset_to_relative_commit",
			commit:  "HEAD~3",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "success_reset_to_branch",
			commit:  "main",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "success_reset_to_tag",
			commit:  "v1.0.0",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "error_invalid_commit",
			commit:  "nonexistent-commit",
			err:     errors.New("fatal: ambiguous argument 'nonexistent-commit': unknown revision or path not in the working tree"),
			wantErr: true,
		},
		{
			name:    "error_not_git_repository",
			commit:  "HEAD",
			err:     errors.New("fatal: not a git repository (or any of the parent directories): .git"),
			wantErr: true,
		},
		{
			name:    "success_reset_to_remote_branch",
			commit:  "origin/develop",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "error_corrupted_repository",
			commit:  "HEAD",
			err:     errors.New("fatal: bad object HEAD"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					expectedArgs := []string{"reset", "--hard", tt.commit}
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

			err := c.ResetHard(tt.commit)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResetHard() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
