package git

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestCommitAllowEmpty_ExecutesCorrectCommand(t *testing.T) {
	var gotArgs []string
	origExecCommand := execCommand
	execCommand = func(name string, args ...string) *exec.Cmd {
		gotArgs = append([]string{name}, args...)
		return exec.Command("echo")
	}
	defer func() { execCommand = origExecCommand }()

	_ = CommitAllowEmpty()
	want := []string{"git", "commit", "--allow-empty", "-m", "empty commit"}
	if !reflect.DeepEqual(gotArgs, want) {
		t.Errorf("got %v, want %v", gotArgs, want)
	}
}

func TestCommitTmp_ExecutesCorrectCommand(t *testing.T) {
	var gotArgs []string
	origExecCommand := execCommand
	execCommand = func(name string, args ...string) *exec.Cmd {
		gotArgs = append([]string{name}, args...)
		return exec.Command("echo")
	}
	defer func() { execCommand = origExecCommand }()

	_ = CommitTmp()
	want := []string{"git", "commit", "-m", "tmp"}
	if !reflect.DeepEqual(gotArgs, want) {
		t.Errorf("got %v, want %v", gotArgs, want)
	}
}
