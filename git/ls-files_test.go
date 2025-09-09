package git

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestClient_ListFiles(t *testing.T) {
	var gotArgs []string
	expectedOutput := "file1.go\nfile2.go\nREADME.md\n"

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

// Error case test for better coverage
func TestClient_ListFiles_Error(t *testing.T) {
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			return exec.Command("false") // Command that always fails
		},
	}

	_, err := client.ListFiles()
	if err == nil {
		t.Error("Expected ListFiles to return an error")
	}
}
