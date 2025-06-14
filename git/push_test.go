package git

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestPushCurrentBranch_ExecutesCorrectCommand(t *testing.T) {
	// Mock: getCurrentBranch
	origGetCurrentBranch := getCurrentBranch
	getCurrentBranch = func() (string, error) {
		return "main", nil
	}
	defer func() { getCurrentBranch = origGetCurrentBranch }()

	// Mock: execCommand
	var gotArgs []string
	origExecCommand := execCommand
	execCommand = func(name string, args ...string) *exec.Cmd {
		gotArgs = append([]string{name}, args...)
		return exec.Command("echo") // Actually does nothing
	}
	defer func() { execCommand = origExecCommand }()

	// Test execution
	_ = PushCurrentBranch()

	// Verification
	want := []string{"git", "push", "origin", "main"}
	if !reflect.DeepEqual(gotArgs, want) {
		t.Errorf("got %v, want %v", gotArgs, want)
	}
}
