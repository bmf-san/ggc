package git

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestClient_StashPullPop(t *testing.T) {
	var gotCommands [][]string
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotCommands = append(gotCommands, append([]string{name}, args...))
			return exec.Command("echo")
		},
	}

	err := client.StashPullPop()
	if err != nil {
		t.Errorf("StashPullPop() error = %v, want nil", err)
	}

	want := [][]string{
		{"git", "stash"},
		{"git", "pull"},
		{"git", "stash", "pop"},
	}
	if !reflect.DeepEqual(gotCommands, want) {
		t.Errorf("StashPullPop() commands = %v, want %v", gotCommands, want)
	}
}

func TestStashPullPop_Error(t *testing.T) {
	client := &Client{
		execCommand: func(_ string, _  ...string) *exec.Cmd {
			return exec.Command("false") // Always fails
		},
	}

	err := client.StashPullPop()
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestStashPullPop_PullError(t *testing.T) {
	var cmdCount int
	client := &Client{
		execCommand: func(_ string, args ...string) *exec.Cmd {
			cmdCount++
			// stashは成功、pullで失敗
			if cmdCount == 1 && len(args) > 0 && args[0] == "stash" && len(args) == 1 {
				return exec.Command("echo") // stash成功
			}
			if cmdCount == 2 && len(args) > 0 && args[0] == "pull" {
				return exec.Command("false") // pull失敗
			}
			return exec.Command("echo")
		},
	}

	err := client.StashPullPop()
	if err == nil {
		t.Error("Expected error when pull fails")
	}
}

func TestStashPullPop_PopError(t *testing.T) {
	var cmdCount int
	client := &Client{
		execCommand: func(_ string, args ...string) *exec.Cmd {
			cmdCount++
			// stashとpullは成功、popで失敗
			if cmdCount == 1 && len(args) > 0 && args[0] == "stash" && len(args) == 1 {
				return exec.Command("echo") // stash成功
			}
			if cmdCount == 2 && len(args) > 0 && args[0] == "pull" {
				return exec.Command("echo") // pull成功
			}
			if cmdCount == 3 && len(args) > 1 && args[0] == "stash" && args[1] == "pop" {
				return exec.Command("false") // pop失敗
			}
			return exec.Command("echo")
		},
	}

	err := client.StashPullPop()
	if err == nil {
		t.Error("Expected error when pop fails")
	}
}
