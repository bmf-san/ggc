package git

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestClient_ListFiles(t *testing.T) {
	var gotArgs []string
	expectedOutput := "file1.go\nfile2.go\nREADME.md"

	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo", "-n", expectedOutput)
		},
	}

	result, err := client.ListFiles()
	if err != nil {
		t.Errorf("ListFiles() error = %v", err)
	}

	wantArgs := []string{"git", "ls-files"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("ListFiles() gotArgs = %v, want %v", gotArgs, wantArgs)
	}

	if result != expectedOutput {
		t.Errorf("ListFiles() result = %v, want %v", result, expectedOutput)
	}
}

func TestClient_GetUpstreamBranchName(t *testing.T) {
	var gotArgs []string
	expectedOutput := "origin/main"

	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo", "-n", expectedOutput)
		},
	}

	result, err := client.GetUpstreamBranchName("main")
	if err != nil {
		t.Errorf("GetUpstreamBranchName() error = %v", err)
	}

	wantArgs := []string{"git", "rev-parse", "--abbrev-ref", "main@{upstream}"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("GetUpstreamBranchName() gotArgs = %v, want %v", gotArgs, wantArgs)
	}

	if result != expectedOutput {
		t.Errorf("GetUpstreamBranchName() result = %v, want %v", result, expectedOutput)
	}
}

func TestClient_GetAheadBehindCount(t *testing.T) {
	var gotArgs []string
	expectedOutput := "2	1"

	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo", "-n", expectedOutput)
		},
	}

	result, err := client.GetAheadBehindCount("main", "origin/main")
	if err != nil {
		t.Errorf("GetAheadBehindCount() error = %v", err)
	}

	wantArgs := []string{"git", "rev-list", "--left-right", "--count", "main...origin/main"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("GetAheadBehindCount() gotArgs = %v, want %v", gotArgs, wantArgs)
	}

	if result != expectedOutput {
		t.Errorf("GetAheadBehindCount() result = %v, want %v", result, expectedOutput)
	}
}
