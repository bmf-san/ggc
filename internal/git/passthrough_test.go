package git

import (
	"errors"
	"os/exec"
	"slices"
	"testing"
)

func TestClient_RunGit(t *testing.T) {
	tests := []struct {
		name     string
		cmdName  string
		args     []string
		wantArgs []string
	}{
		{
			name:     "no args",
			cmdName:  "blame",
			args:     nil,
			wantArgs: []string{"git", "blame"},
		},
		{
			name:     "with args",
			cmdName:  "cherry-pick",
			args:     []string{"-x", "abc123"},
			wantArgs: []string{"git", "cherry-pick", "-x", "abc123"},
		},
		{
			name:     "subcommand with flag",
			cmdName:  "worktree",
			args:     []string{"add", "-b", "feat", "../tree"},
			wantArgs: []string{"git", "worktree", "add", "-b", "feat", "../tree"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotArgs []string
			client := &Client{
				execCommand: func(name string, args ...string) *exec.Cmd {
					gotArgs = append([]string{name}, args...)
					return exec.Command("echo")
				},
			}
			if err := client.RunGit(tt.cmdName, tt.args); err != nil {
				t.Errorf("RunGit() error = %v", err)
			}
			if !slices.Equal(gotArgs, tt.wantArgs) {
				t.Errorf("RunGit() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestClient_RunGit_Error(t *testing.T) {
	client := &Client{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return helperCommand(t, "", errors.New("boom"))
		},
	}
	err := client.RunGit("revert", []string{"HEAD"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
