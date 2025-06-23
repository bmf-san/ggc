package git

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestClient_CommitAllowEmpty(t *testing.T) {
	var gotArgs []string
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo")
		},
	}

	_ = client.CommitAllowEmpty()
	want := []string{"git", "commit", "--allow-empty", "-m", "empty commit"}
	if !reflect.DeepEqual(gotArgs, want) {
		t.Errorf("got %v, want %v", gotArgs, want)
	}
}

func TestClient_CommitTmp(t *testing.T) {
	var gotArgs []string
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo")
		},
	}

	_ = client.CommitTmp()
	want := []string{"git", "commit", "-m", "tmp"}
	if !reflect.DeepEqual(gotArgs, want) {
		t.Errorf("got %v, want %v", gotArgs, want)
	}
}
