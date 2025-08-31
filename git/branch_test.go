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

func TestClient_RenameMoveSetUpstream(t *testing.T) {
	// Validate commands are constructed correctly and errors handled
	t.Run("rename_branch_command", func(t *testing.T) {
		c := &Client{execCommand: func(name string, arg ...string) *exec.Cmd {
			if name != "git" || strings.Join(arg, " ") != "branch -m old new" {
				t.Errorf("unexpected command: %s %v", name, arg)
			}
			return helperCommand(t, "", nil)
		}}
		if err := c.RenameBranch("old", "new"); err != nil {
			t.Errorf("RenameBranch() error = %v", err)
		}
	})

	t.Run("move_branch_command", func(t *testing.T) {
		c := &Client{execCommand: func(name string, arg ...string) *exec.Cmd {
			if name != "git" || strings.Join(arg, " ") != "branch -f feat abc123" {
				t.Errorf("unexpected command: %s %v", name, arg)
			}
			return helperCommand(t, "", nil)
		}}
		if err := c.MoveBranch("feat", "abc123"); err != nil {
			t.Errorf("MoveBranch() error = %v", err)
		}
	})

	t.Run("set_upstream_command", func(t *testing.T) {
		c := &Client{execCommand: func(name string, arg ...string) *exec.Cmd {
			if name != "git" || strings.Join(arg, " ") != "branch -u origin/main feat" {
				t.Errorf("unexpected command: %s %v", name, arg)
			}
			return helperCommand(t, "", nil)
		}}
		if err := c.SetUpstreamBranch("feat", "origin/main"); err != nil {
			t.Errorf("SetUpstreamBranch() error = %v", err)
		}
	})
}

func TestClient_ListBranchesVerbose_Parse(t *testing.T) {
	output := `* main    1a2b3c4 [origin/main: ahead 2, behind 1] Update README
  feature 5d6e7f8 [origin/feature] Implement feature
  local    abcdef0 Local only commit`
	c := &Client{execCommand: func(name string, arg ...string) *exec.Cmd {
		if name != "git" || !strings.Contains(strings.Join(arg, " "), "branch -vv") {
			t.Errorf("unexpected command: %s %v", name, arg)
		}
		return fakeExecCommand(output)
	}}
	infos, err := c.ListBranchesVerbose()
	if err != nil {
		t.Fatalf("ListBranchesVerbose error: %v", err)
	}
	if len(infos) != 3 {
		t.Fatalf("expected 3 branches, got %d", len(infos))
	}
	if !infos[0].IsCurrentBranch || infos[0].Name != "main" || infos[0].Upstream != "origin/main" || infos[0].AheadBehind == "" {
		t.Errorf("unexpected main info: %+v", infos[0])
	}
	if infos[1].Upstream != "origin/feature" || infos[1].AheadBehind != "" {
		t.Errorf("unexpected feature info: %+v", infos[1])
	}
	if infos[2].Upstream != "" || infos[2].AheadBehind != "" {
		t.Errorf("unexpected local info: %+v", infos[2])
	}
}

func TestClient_GetBranchInfo_Fallback(t *testing.T) {
	// ListBranchesVerbose returns empty -> fallback path
	step := 0
	c := &Client{execCommand: func(_ string, _ ...string) *exec.Cmd {
		step++
		switch step {
		case 1:
			// branch -vv
			return helperCommand(t, "", nil)
		case 2:
			// rev-parse --short branch
			return fakeExecCommand("abc1234")
		case 3:
			// log -1 --pretty=%s branch
			return fakeExecCommand("Commit message")
		default:
			return fakeExecCommand("")
		}
	}}
	info, err := c.GetBranchInfo("test")
	if err != nil {
		t.Fatalf("GetBranchInfo error: %v", err)
	}
	if info == nil || info.LastCommitSHA != "abc1234" || info.LastCommitMsg != "Commit message" {
		t.Errorf("unexpected info: %+v", info)
	}
}

func TestClient_SortBranches(t *testing.T) {
	c := &Client{execCommand: func(name string, arg ...string) *exec.Cmd {
		if name != "git" || !strings.Contains(strings.Join(arg, " "), "branch --sort=-committerdate") {
			t.Errorf("unexpected command: %s %v", name, arg)
		}
		return fakeExecCommand("b1\nb2\nb3")
	}}
	got, err := c.SortBranches("date")
	if err != nil {
		t.Fatalf("SortBranches error: %v", err)
	}
	if !reflect.DeepEqual(got, []string{"b1", "b2", "b3"}) {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestClient_BranchesContaining(t *testing.T) {
	c := &Client{execCommand: func(name string, arg ...string) *exec.Cmd {
		if name != "git" || !strings.Contains(strings.Join(arg, " "), "branch --contains abc123") {
			t.Errorf("unexpected command: %s %v", name, arg)
		}
		return fakeExecCommand("* main\n  feature\n  bugfix")
	}}
	got, err := c.BranchesContaining("abc123")
	if err != nil {
		t.Fatalf("BranchesContaining error: %v", err)
	}
	if !reflect.DeepEqual(got, []string{"main", "feature", "bugfix"}) {
		t.Errorf("unexpected branches: %v", got)
	}
}
