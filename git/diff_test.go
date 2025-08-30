package git

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestClient_Diff(t *testing.T) {
	var gotArgs []string
	expectedOutput := "diff --git a/file.go b/file.go\nindex 1234567..abcdefg 100644"

	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo", "-n", expectedOutput)
		},
	}

	result, err := client.Diff()
	if err != nil {
		t.Errorf("Diff() error = %v", err)
	}

	wantArgs := []string{"git", "diff"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("Diff() gotArgs = %v, want %v", gotArgs, wantArgs)
	}

	if result != expectedOutput {
		t.Errorf("Diff() result = %v, want %v", result, expectedOutput)
	}
}

func TestClient_DiffStaged(t *testing.T) {
	var gotArgs []string
	expectedOutput := "diff --git a/staged.go b/staged.go\nindex 1234567..abcdefg 100644"

	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo", "-n", expectedOutput)
		},
	}

	result, err := client.DiffStaged()
	if err != nil {
		t.Errorf("DiffStaged() error = %v", err)
	}

	wantArgs := []string{"git", "diff", "--staged"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("DiffStaged() gotArgs = %v, want %v", gotArgs, wantArgs)
	}

	if result != expectedOutput {
		t.Errorf("DiffStaged() result = %v, want %v", result, expectedOutput)
	}
}

func TestClient_DiffHead(t *testing.T) {
	var gotArgs []string
	expectedOutput := "diff --git a/head.go b/head.go\nindex 1234567..abcdefg 100644"

	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo", "-n", expectedOutput)
		},
	}

	result, err := client.DiffHead()
	if err != nil {
		t.Errorf("DiffHead() error = %v", err)
	}

	wantArgs := []string{"git", "diff", "HEAD"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("DiffHead() gotArgs = %v, want %v", gotArgs, wantArgs)
	}

	if result != expectedOutput {
		t.Errorf("DiffHead() result = %v, want %v", result, expectedOutput)
	}
}
