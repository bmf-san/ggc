package git

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestClient_ResetHardAndClean(t *testing.T) {
	var gotArgs [][]string
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append(gotArgs, append([]string{name}, args...))
			return exec.Command("echo")
		},
		GetCurrentBranchFunc: func() (string, error) {
			return "main", nil
		},
	}

	_ = client.ResetHardAndClean()
	want := [][]string{
		{"git", "reset", "--hard", "origin/main"},
		{"git", "clean", "-fdx"},
	}
	if !reflect.DeepEqual(gotArgs, want) {
		t.Errorf("got %v, want %v", gotArgs, want)
	}
}
