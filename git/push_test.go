package git

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestPushCurrentBranch_ExecutesCorrectCommand(t *testing.T) {
	// モック: getCurrentBranch
	origGetCurrentBranch := getCurrentBranch
	getCurrentBranch = func() (string, error) {
		return "main", nil
	}
	defer func() { getCurrentBranch = origGetCurrentBranch }()

	// モック: execCommand
	var gotArgs []string
	origExecCommand := execCommand
	execCommand = func(name string, args ...string) *exec.Cmd {
		gotArgs = append([]string{name}, args...)
		return exec.Command("echo") // 実際には何もしない
	}
	defer func() { execCommand = origExecCommand }()

	// テスト実行
	_ = PushCurrentBranch()

	// 検証
	want := []string{"git", "push", "origin", "main"}
	if !reflect.DeepEqual(gotArgs, want) {
		t.Errorf("got %v, want %v", gotArgs, want)
	}
}
