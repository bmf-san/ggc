package git

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestPullCurrentBranch_ExecutesCorrectCommand(t *testing.T) {
	origGetCurrentBranch := getCurrentBranch
	getCurrentBranch = func() (string, error) { return "main", nil }
	defer func() { getCurrentBranch = origGetCurrentBranch }()

	var gotArgs []string
	origExecCommand := execCommand
	execCommand = func(name string, args ...string) *exec.Cmd {
		gotArgs = append([]string{name}, args...)
		return exec.Command("echo")
	}
	defer func() { execCommand = origExecCommand }()

	_ = PullCurrentBranch()
	want := []string{"git", "pull", "origin", "main"}
	if !reflect.DeepEqual(gotArgs, want) {
		t.Errorf("got %v, want %v", gotArgs, want)
	}
}

func TestPullRebaseCurrentBranch_ExecutesCorrectCommand(t *testing.T) {
	origGetCurrentBranch := getCurrentBranch
	getCurrentBranch = func() (string, error) { return "main", nil }
	defer func() { getCurrentBranch = origGetCurrentBranch }()

	var gotArgs []string
	origExecCommand := execCommand
	execCommand = func(name string, args ...string) *exec.Cmd {
		gotArgs = append([]string{name}, args...)
		return exec.Command("echo")
	}
	defer func() { execCommand = origExecCommand }()

	_ = PullRebaseCurrentBranch()
	want := []string{"git", "pull", "--rebase", "origin", "main"}
	if !reflect.DeepEqual(gotArgs, want) {
		t.Errorf("got %v, want %v", gotArgs, want)
	}
}
