package git

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestClient_RemoteList(t *testing.T) {
	var gotArgs []string
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo")
		},
	}

	err := client.RemoteList()
	if err != nil {
		t.Errorf("RemoteList() error = %v", err)
	}

	wantArgs := []string{"git", "remote", "-v"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("RemoteList() gotArgs = %v, want %v", gotArgs, wantArgs)
	}
}

func TestClient_RemoteAdd(t *testing.T) {
	var gotArgs []string
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo")
		},
	}

	err := client.RemoteAdd("origin", "https://github.com/user/repo.git")
	if err != nil {
		t.Errorf("RemoteAdd() error = %v", err)
	}

	wantArgs := []string{"git", "remote", "add", "origin", "https://github.com/user/repo.git"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("RemoteAdd() gotArgs = %v, want %v", gotArgs, wantArgs)
	}
}

func TestClient_RemoteRemove(t *testing.T) {
	var gotArgs []string
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo")
		},
	}

	err := client.RemoteRemove("origin")
	if err != nil {
		t.Errorf("RemoteRemove() error = %v", err)
	}

	wantArgs := []string{"git", "remote", "remove", "origin"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("RemoteRemove() gotArgs = %v, want %v", gotArgs, wantArgs)
	}
}

func TestClient_RemoteSetURL(t *testing.T) {
	var gotArgs []string
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo")
		},
	}

	err := client.RemoteSetURL("origin", "https://github.com/user/new-repo.git")
	if err != nil {
		t.Errorf("RemoteSetURL() error = %v", err)
	}

	wantArgs := []string{"git", "remote", "set-url", "origin", "https://github.com/user/new-repo.git"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("RemoteSetURL() gotArgs = %v, want %v", gotArgs, wantArgs)
	}
}
