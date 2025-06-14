package git

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestFetchPrune_ExecutesCorrectCommand(t *testing.T) {
	var gotArgs []string
	origExecCommand := execCommand
	execCommand = func(name string, args ...string) *exec.Cmd {
		gotArgs = append([]string{name}, args...)
		return exec.Command("echo")
	}
	defer func() { execCommand = origExecCommand }()

	_ = FetchPrune()
	want := []string{"git", "fetch", "--prune"}
	if !reflect.DeepEqual(gotArgs, want) {
		t.Errorf("got %v, want %v", gotArgs, want)
	}
}
