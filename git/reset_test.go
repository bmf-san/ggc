package git

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestResetClean_ExecutesCorrectCommands(t *testing.T) {
	var gotArgs [][]string
	origExecCommand := execCommand
	execCommand = func(name string, args ...string) *exec.Cmd {
		gotArgs = append(gotArgs, append([]string{name}, args...))
		return exec.Command("echo")
	}
	defer func() { execCommand = origExecCommand }()

	_ = ResetClean()
	want := [][]string{
		{"git", "reset", "--hard", "HEAD"},
		{"git", "clean", "-fd"},
	}
	if !reflect.DeepEqual(gotArgs, want) {
		t.Errorf("got %v, want %v", gotArgs, want)
	}
}
