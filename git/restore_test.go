package git

import (
	"os/exec"
	"slices"
	"testing"
)

func TestClient_Restore(t *testing.T) {
	cases := []struct {
		name     string
		paths    []string
		opts     *RestoreOptions
		wantArgs []string
	}{
		{
			name:     "restore single file",
			paths:    []string{"file.txt"},
			opts:     nil,
			wantArgs: []string{"git", "restore", "file.txt"},
		},
		{
			name:     "restore multiple files",
			paths:    []string{"file1.txt", "file2.txt"},
			opts:     nil,
			wantArgs: []string{"git", "restore", "file1.txt", "file2.txt"},
		},
		{
			name:     "restore staged file",
			paths:    []string{"file.txt"},
			opts:     &RestoreOptions{Staged: true},
			wantArgs: []string{"git", "restore", "--staged", "file.txt"},
		},
		{
			name:     "restore from specific commit",
			paths:    []string{"file.txt"},
			opts:     &RestoreOptions{Source: "abc123"},
			wantArgs: []string{"git", "restore", "--source", "abc123", "file.txt"},
		},
		{
			name:     "restore staged from specific commit",
			paths:    []string{"file.txt"},
			opts:     &RestoreOptions{Staged: true, Source: "abc123"},
			wantArgs: []string{"git", "restore", "--staged", "--source", "abc123", "file.txt"},
		},
		{
			name:     "restore with empty source",
			paths:    []string{"file.txt"},
			opts:     &RestoreOptions{Source: ""},
			wantArgs: []string{"git", "restore", "file.txt"},
		},
		{
			name:     "restore all files",
			paths:    []string{"."},
			opts:     nil,
			wantArgs: []string{"git", "restore", "."},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var gotArgs []string
			client := &Client{
				execCommand: func(name string, args ...string) *exec.Cmd {
					gotArgs = append([]string{name}, args...)
					return exec.Command("echo")
				},
			}

			_ = client.Restore(tc.paths, tc.opts)

			if !slices.Equal(gotArgs, tc.wantArgs) {
				t.Errorf("got %v, want %v", gotArgs, tc.wantArgs)
			}
		})
	}
}

func TestClient_RestoreWorkingDir(t *testing.T) {
	cases := []struct {
		name     string
		paths    []string
		wantArgs []string
	}{
		{
			name:     "restore single file from working dir",
			paths:    []string{"file.txt"},
			wantArgs: []string{"git", "restore", "file.txt"},
		},
		{
			name:     "restore multiple files from working dir",
			paths:    []string{"file1.txt", "file2.txt", "dir/file3.txt"},
			wantArgs: []string{"git", "restore", "file1.txt", "file2.txt", "dir/file3.txt"},
		},
		{
			name:     "restore no files",
			paths:    []string{},
			wantArgs: []string{"git", "restore"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var gotArgs []string
			client := &Client{
				execCommand: func(name string, args ...string) *exec.Cmd {
					gotArgs = append([]string{name}, args...)
					return exec.Command("echo")
				},
			}

			_ = client.RestoreWorkingDir(tc.paths...)

			if !slices.Equal(gotArgs, tc.wantArgs) {
				t.Errorf("got %v, want %v", gotArgs, tc.wantArgs)
			}
		})
	}
}

func TestClient_RestoreStaged(t *testing.T) {
	cases := []struct {
		name     string
		paths    []string
		wantArgs []string
	}{
		{
			name:     "unstage single file",
			paths:    []string{"file.txt"},
			wantArgs: []string{"git", "restore", "--staged", "file.txt"},
		},
		{
			name:     "unstage multiple files",
			paths:    []string{"file1.txt", "file2.txt"},
			wantArgs: []string{"git", "restore", "--staged", "file1.txt", "file2.txt"},
		},
		{
			name:     "unstage no files",
			paths:    []string{},
			wantArgs: []string{"git", "restore", "--staged"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var gotArgs []string
			client := &Client{
				execCommand: func(name string, args ...string) *exec.Cmd {
					gotArgs = append([]string{name}, args...)
					return exec.Command("echo")
				},
			}

			_ = client.RestoreStaged(tc.paths...)

			if !slices.Equal(gotArgs, tc.wantArgs) {
				t.Errorf("got %v, want %v", gotArgs, tc.wantArgs)
			}
		})
	}
}

func TestClient_RestoreFromCommit(t *testing.T) {
	cases := []struct {
		name     string
		commit   string
		paths    []string
		wantArgs []string
	}{
		{
			name:     "restore single file from commit",
			commit:   "abc123",
			paths:    []string{"file.txt"},
			wantArgs: []string{"git", "restore", "--source", "abc123", "file.txt"},
		},
		{
			name:     "restore multiple files from commit",
			commit:   "def456",
			paths:    []string{"file1.txt", "file2.txt"},
			wantArgs: []string{"git", "restore", "--source", "def456", "file1.txt", "file2.txt"},
		},
		{
			name:     "restore from HEAD",
			commit:   "HEAD",
			paths:    []string{"file.txt"},
			wantArgs: []string{"git", "restore", "--source", "HEAD", "file.txt"},
		},
		{
			name:     "restore from branch",
			commit:   "feature/branch",
			paths:    []string{"file.txt"},
			wantArgs: []string{"git", "restore", "--source", "feature/branch", "file.txt"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var gotArgs []string
			client := &Client{
				execCommand: func(name string, args ...string) *exec.Cmd {
					gotArgs = append([]string{name}, args...)
					return exec.Command("echo")
				},
			}

			_ = client.RestoreFromCommit(tc.commit, tc.paths...)

			if !slices.Equal(gotArgs, tc.wantArgs) {
				t.Errorf("got %v, want %v", gotArgs, tc.wantArgs)
			}
		})
	}
}

func TestClient_RestoreAll(t *testing.T) {
	var gotArgs []string
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo")
		},
	}

	_ = client.RestoreAll()

	wantArgs := []string{"git", "restore", "."}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("got %v, want %v", gotArgs, wantArgs)
	}
}

func TestClient_RestoreAllStaged(t *testing.T) {
	var gotArgs []string
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo")
		},
	}

	_ = client.RestoreAllStaged()

	wantArgs := []string{"git", "restore", "--staged", "."}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("got %v, want %v", gotArgs, wantArgs)
	}
}
