package git

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestCleanFiles_ExecutesCorrectCommand(t *testing.T) {
	var gotArgs []string
	origExecCommand := execCommand
	execCommand = func(name string, args ...string) *exec.Cmd {
		gotArgs = append([]string{name}, args...)
		return exec.Command("echo")
	}
	defer func() { execCommand = origExecCommand }()

	_ = CleanFiles()
	want := []string{"git", "clean", "-f"}
	if !reflect.DeepEqual(gotArgs, want) {
		t.Errorf("got %v, want %v", gotArgs, want)
	}
}

func TestCleanDirs_ExecutesCorrectCommand(t *testing.T) {
	var gotArgs []string
	origExecCommand := execCommand
	execCommand = func(name string, args ...string) *exec.Cmd {
		gotArgs = append([]string{name}, args...)
		return exec.Command("echo")
	}
	defer func() { execCommand = origExecCommand }()

	_ = CleanDirs()
	want := []string{"git", "clean", "-d"}
	if !reflect.DeepEqual(gotArgs, want) {
		t.Errorf("got %v, want %v", gotArgs, want)
	}
}
