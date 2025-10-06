package git

import (
	"errors"
	"os/exec"
	"slices"
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
			if !slices.Equal(got, tt.want) {
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
			if !slices.Equal(got, tt.want) {
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

	t.Run("move_branch_empty_commit", func(t *testing.T) {
		c := &Client{execCommand: func(name string, arg ...string) *exec.Cmd {
			t.Fatalf("execCommand should not be called for invalid commit")
			return helperCommand(t, "", nil)
		}}
		if err := c.MoveBranch("feat", "   "); err == nil {
			t.Error("Expected error for empty commit, got nil")
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

	t.Run("set_upstream_empty_remote", func(t *testing.T) {
		c := &Client{execCommand: func(name string, arg ...string) *exec.Cmd {
			t.Fatalf("execCommand should not be called for empty upstream")
			return helperCommand(t, "", nil)
		}}
		if err := c.SetUpstreamBranch("feat", "  "); err == nil {
			t.Error("Expected error for empty upstream, got nil")
		}
	})
}

func TestValidateBranchName(t *testing.T) {
	tests := []struct {
		name    string
		branch  string
		wantErr bool
	}{
		{name: "valid_simple", branch: "feature", wantErr: false},
		{name: "valid_nested", branch: "feature/awesome", wantErr: false},
		{name: "invalid_empty", branch: "", wantErr: true},
		{name: "invalid_with_space", branch: "feature branch", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBranchName(tt.branch)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBranchName(%q) error = %v, wantErr %v", tt.branch, err, tt.wantErr)
			}
		})
	}
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
	if !slices.Equal(got, []string{"b1", "b2", "b3"}) {
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
	if !slices.Equal(got, []string{"main", "feature", "bugfix"}) {
		t.Errorf("unexpected branches: %v", got)
	}
}

func TestClient_CheckoutBranch(t *testing.T) {
	tests := []struct {
		name       string
		branchName string
		err        error
		wantErr    bool
	}{
		{
			name:       "success_checkout_main",
			branchName: "main",
			err:        nil,
			wantErr:    false,
		},
		{
			name:       "success_checkout_feature_branch",
			branchName: "feature/new-feature",
			err:        nil,
			wantErr:    false,
		},
		{
			name:       "error_branch_not_found",
			branchName: "nonexistent-branch",
			err:        errors.New("error: pathspec 'nonexistent-branch' did not match any file(s) known to git"),
			wantErr:    true,
		},
		{
			name:       "error_uncommitted_changes",
			branchName: "develop",
			err:        errors.New("error: Your local changes to the following files would be overwritten by checkout"),
			wantErr:    true,
		},
		{
			name:       "success_checkout_with_special_chars",
			branchName: "feature/user-story_123",
			err:        nil,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					expectedArgs := []string{"checkout", tt.branchName}
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

			err := c.CheckoutBranch(tt.branchName)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckoutBranch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_CheckoutNewBranchFromRemote(t *testing.T) {
	tests := []struct {
		name         string
		localBranch  string
		remoteBranch string
		err          error
		wantErr      bool
		expectExec   bool
	}{
		{
			name:         "success_checkout_from_origin",
			localBranch:  "feature",
			remoteBranch: "origin/feature",
			err:          nil,
			wantErr:      false,
			expectExec:   true,
		},
		{
			name:         "success_checkout_from_upstream",
			localBranch:  "develop",
			remoteBranch: "upstream/develop",
			err:          nil,
			wantErr:      false,
			expectExec:   true,
		},
		{
			name:         "error_remote_branch_not_found",
			localBranch:  "feature",
			remoteBranch: "origin/nonexistent",
			err:          errors.New("fatal: 'origin/nonexistent' is not a commit and a branch 'feature' cannot be created from it"),
			wantErr:      true,
			expectExec:   true,
		},
		{
			name:         "error_local_branch_exists",
			localBranch:  "main",
			remoteBranch: "origin/main",
			err:          errors.New("fatal: A branch named 'main' already exists"),
			wantErr:      true,
			expectExec:   true,
		},
		{
			name:         "success_deep_branch_hierarchy",
			localBranch:  "feature-local",
			remoteBranch: "origin/feature/user/story/implementation",
			err:          nil,
			wantErr:      false,
			expectExec:   true,
		},
		{
			name:         "error_invalid_local_branch",
			localBranch:  "feature branch",
			remoteBranch: "origin/feature-branch",
			err:          nil,
			wantErr:      true,
			expectExec:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executed := false
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					executed = true
					expectedArgs := []string{"checkout", "-b", tt.localBranch, "--track", tt.remoteBranch}
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

			err := c.CheckoutNewBranchFromRemote(tt.localBranch, tt.remoteBranch)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckoutNewBranchFromRemote() error = %v, wantErr %v", err, tt.wantErr)
			}
			if executed != tt.expectExec {
				t.Errorf("CheckoutNewBranchFromRemote() exec invoked = %v, expectExec %v", executed, tt.expectExec)
			}
		})
	}
}

func TestClient_DeleteBranch(t *testing.T) {
	tests := []struct {
		name       string
		branchName string
		err        error
		wantErr    bool
	}{
		{
			name:       "success_delete_merged_branch",
			branchName: "feature-completed",
			err:        nil,
			wantErr:    false,
		},
		{
			name:       "success_delete_feature_branch",
			branchName: "feature/old-feature",
			err:        nil,
			wantErr:    false,
		},
		{
			name:       "error_delete_unmerged_branch",
			branchName: "feature-unmerged",
			err:        errors.New("error: The branch 'feature-unmerged' is not fully merged"),
			wantErr:    true,
		},
		{
			name:       "error_delete_current_branch",
			branchName: "main",
			err:        errors.New("error: Cannot delete branch 'main' checked out at"),
			wantErr:    true,
		},
		{
			name:       "error_branch_not_found",
			branchName: "nonexistent-branch",
			err:        errors.New("error: branch 'nonexistent-branch' not found"),
			wantErr:    true,
		},
		{
			name:       "success_delete_bugfix_branch",
			branchName: "bugfix/issue-123",
			err:        nil,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					expectedArgs := []string{"branch", "-d", tt.branchName}
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

			err := c.DeleteBranch(tt.branchName)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteBranch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_ListMergedBranches(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		err     error
		want    []string
		wantErr bool
	}{
		{
			name:    "success_multiple_merged_branches",
			output:  "  feature/completed\n* main\n  bugfix/issue-123\n  hotfix/security-fix\n",
			err:     nil,
			want:    []string{"feature/completed", "bugfix/issue-123", "hotfix/security-fix"},
			wantErr: false,
		},
		{
			name:    "success_single_merged_branch",
			output:  "  feature/old-feature\n* main\n",
			err:     nil,
			want:    []string{"feature/old-feature"},
			wantErr: false,
		},
		{
			name:    "success_no_merged_branches",
			output:  "* main\n",
			err:     nil,
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "success_empty_output",
			output:  "",
			err:     nil,
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "success_with_extra_whitespace",
			output:  "  feature/merged  \n* main  \n  another/merged  \n",
			err:     nil,
			want:    []string{"feature/merged", "another/merged"},
			wantErr: false,
		},
		{
			name:    "success_exclude_current_branch_variations",
			output:  "* develop\n  feature/merged\n  hotfix/merged\n",
			err:     nil,
			want:    []string{"feature/merged", "hotfix/merged"},
			wantErr: false,
		},
		{
			name:    "error_git_command_failed",
			output:  "",
			err:     errors.New("fatal: not a git repository"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error_permission_denied",
			output:  "",
			err:     errors.New("permission denied"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error_corrupted_repository",
			output:  "",
			err:     errors.New("fatal: bad object HEAD"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "success_mixed_branch_formats",
			output:  "  release/v1.0\n* main\n  feature/user-story_123\n  bugfix-urgent\n",
			err:     nil,
			want:    []string{"release/v1.0", "feature/user-story_123", "bugfix-urgent"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					expectedArgs := []string{"branch", "--merged"}
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

			got, err := c.ListMergedBranches()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListMergedBranches() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("ListMergedBranches() = %v, want %v", got, tt.want)
			}
		})
	}
}
