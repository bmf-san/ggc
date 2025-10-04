package git

import (
	"os/exec"
	"slices"
	"testing"
)

func TestClient_Add(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		wantArgs []string
		wantErr  bool
	}{
		{
			name:     "add single file",
			files:    []string{"file1.go"},
			wantArgs: []string{"git", "add", "file1.go"},
			wantErr:  false,
		},
		{
			name:     "add multiple files",
			files:    []string{"file1.go", "file2.go"},
			wantArgs: []string{"git", "add", "file1.go", "file2.go"},
			wantErr:  false,
		},
		{
			name:     "add all files",
			files:    []string{"."},
			wantArgs: []string{"git", "add", "."},
			wantErr:  false,
		},
		{
			name:    "no files provided",
			files:   []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotArgs []string
			client := &Client{
				execCommand: func(name string, args ...string) *exec.Cmd {
					gotArgs = append([]string{name}, args...)
					return exec.Command("echo")
				},
			}

			err := client.Add(tt.files...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !slices.Equal(gotArgs, tt.wantArgs) {
				t.Errorf("Add() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestClient_AddInteractive(t *testing.T) {
	var gotArgs []string
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo")
		},
	}

	err := client.AddInteractive()
	if err != nil {
		t.Errorf("AddInteractive() error = %v", err)
	}

	wantArgs := []string{"git", "add", "-p"}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("AddInteractive() gotArgs = %v, want %v", gotArgs, wantArgs)
	}
}
