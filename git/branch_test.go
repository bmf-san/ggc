package git

import (
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

func TestClient_ListLocalBranches(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		want    []string
		wantErr bool
	}{
		{
			name:    "success_multiple_branches",
			output:  "main\nfeature/test\ndevelop",
			want:    []string{"main", "feature/test", "develop"},
			wantErr: false,
		},
		{
			name:    "success_single_branch",
			output:  "main",
			want:    []string{"main"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					if name != "git" || !strings.Contains(strings.Join(arg, " "), "branch --format %(refname:short)") {
						t.Errorf("unexpected command: %s %v", name, arg)
					}
					return fakeExecCommand(tt.output)
				},
			}

			got, err := c.ListLocalBranches()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListLocalBranches() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListLocalBranches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ListRemoteBranches(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		want    []string
		wantErr bool
	}{
		{
			name:    "success_exclude_head",
			output:  "origin/main\norigin/HEAD -> origin/main\norigin/feature/test",
			want:    []string{"origin/main", "origin/feature/test"},
			wantErr: false,
		},
		{
			name:    "success_single_remote_branch",
			output:  "origin/main",
			want:    []string{"origin/main"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					if name != "git" || !strings.Contains(strings.Join(arg, " "), "branch -r --format %(refname:short)") {
						t.Errorf("unexpected command: %s %v", name, arg)
					}
					return fakeExecCommand(tt.output)
				},
			}

			got, err := c.ListRemoteBranches()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListRemoteBranches() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListRemoteBranches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListLocalBranches_Error(t *testing.T) {
	client := &Client{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("false") // Always fails
		},
	}

	_, err := client.ListLocalBranches()
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestListRemoteBranches_Error(t *testing.T) {
	client := &Client{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("false") // Always fails
		},
	}

	_, err := client.ListRemoteBranches()
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
