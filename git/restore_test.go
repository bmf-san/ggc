package git

import (
	"os/exec"
	"reflect"
	"strings"
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

			if !reflect.DeepEqual(gotArgs, tc.wantArgs) {
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

			if !reflect.DeepEqual(gotArgs, tc.wantArgs) {
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

			if !reflect.DeepEqual(gotArgs, tc.wantArgs) {
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

			if !reflect.DeepEqual(gotArgs, tc.wantArgs) {
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
	if !reflect.DeepEqual(gotArgs, wantArgs) {
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
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("got %v, want %v", gotArgs, wantArgs)
	}
}

// Test error handling in Restore function to improve coverage from 92.3% to higher
func TestClient_Restore_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		paths       []string
		opts        *RestoreOptions
		commandErr  bool
		expectedErr string
	}{
		{
			name:        "command execution error",
			paths:       []string{"file.txt"},
			opts:        nil,
			commandErr:  true,
			expectedErr: "restore",
		},
		{
			name:        "command execution error with staged option",
			paths:       []string{"file.txt"},
			opts:        &RestoreOptions{Staged: true},
			commandErr:  true,
			expectedErr: "restore",
		},
		{
			name:        "command execution error with source option",
			paths:       []string{"file.txt"},
			opts:        &RestoreOptions{Source: "abc123"},
			commandErr:  true,
			expectedErr: "restore",
		},
		{
			name:        "command execution error with both options",
			paths:       []string{"file.txt"},
			opts:        &RestoreOptions{Staged: true, Source: "HEAD~1"},
			commandErr:  true,
			expectedErr: "restore",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				execCommand: func(name string, args ...string) *exec.Cmd {
					if tt.commandErr {
						// Return a command that will fail
						return exec.Command("false")
					}
					return exec.Command("echo")
				},
			}

			err := client.Restore(tt.paths, tt.opts)

			if tt.commandErr {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				if !strings.Contains(err.Error(), tt.expectedErr) {
					t.Errorf("Expected error to contain %q, got: %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

// Test error handling in wrapper functions to ensure they propagate errors correctly
func TestClient_RestoreWrappers_ErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*Client) error
	}{
		{
			name: "RestoreWorkingDir error",
			testFunc: func(c *Client) error {
				return c.RestoreWorkingDir("file.txt")
			},
		},
		{
			name: "RestoreStaged error",
			testFunc: func(c *Client) error {
				return c.RestoreStaged("file.txt")
			},
		},
		{
			name: "RestoreFromCommit error",
			testFunc: func(c *Client) error {
				return c.RestoreFromCommit("abc123", "file.txt")
			},
		},
		{
			name: "RestoreAll error",
			testFunc: func(c *Client) error {
				return c.RestoreAll()
			},
		},
		{
			name: "RestoreAllStaged error",
			testFunc: func(c *Client) error {
				return c.RestoreAllStaged()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				execCommand: func(name string, args ...string) *exec.Cmd {
					// Return a command that will fail
					return exec.Command("false")
				},
			}

			err := tt.testFunc(client)

			if err == nil {
				t.Error("Expected error but got nil")
			}
			if !strings.Contains(err.Error(), "restore") {
				t.Errorf("Expected error to contain 'restore', got: %v", err)
			}
		})
	}
}

// Test RestoreOptions struct variations to ensure all combinations work
func TestClient_RestoreOptions_Combinations(t *testing.T) {
	tests := []struct {
		name     string
		opts     *RestoreOptions
		wantArgs []string
	}{
		{
			name:     "nil options",
			opts:     nil,
			wantArgs: []string{"git", "restore", "file.txt"},
		},
		{
			name:     "empty options struct",
			opts:     &RestoreOptions{},
			wantArgs: []string{"git", "restore", "file.txt"},
		},
		{
			name:     "only staged true",
			opts:     &RestoreOptions{Staged: true},
			wantArgs: []string{"git", "restore", "--staged", "file.txt"},
		},
		{
			name:     "only staged false",
			opts:     &RestoreOptions{Staged: false},
			wantArgs: []string{"git", "restore", "file.txt"},
		},
		{
			name:     "only source set",
			opts:     &RestoreOptions{Source: "HEAD~1"},
			wantArgs: []string{"git", "restore", "--source", "HEAD~1", "file.txt"},
		},
		{
			name:     "source empty string",
			opts:     &RestoreOptions{Source: ""},
			wantArgs: []string{"git", "restore", "file.txt"},
		},
		{
			name:     "staged true and source set",
			opts:     &RestoreOptions{Staged: true, Source: "main"},
			wantArgs: []string{"git", "restore", "--staged", "--source", "main", "file.txt"},
		},
		{
			name:     "staged false and source set",
			opts:     &RestoreOptions{Staged: false, Source: "develop"},
			wantArgs: []string{"git", "restore", "--source", "develop", "file.txt"},
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

			_ = client.Restore([]string{"file.txt"}, tt.opts)

			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("got %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

// Test edge cases with various path patterns
func TestClient_Restore_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		paths    []string
		opts     *RestoreOptions
		wantArgs []string
	}{
		{
			name:     "empty paths slice",
			paths:    []string{},
			opts:     nil,
			wantArgs: []string{"git", "restore"},
		},
		{
			name:     "single dot path",
			paths:    []string{"."},
			opts:     nil,
			wantArgs: []string{"git", "restore", "."},
		},
		{
			name:     "multiple dot paths",
			paths:    []string{".", ".."},
			opts:     nil,
			wantArgs: []string{"git", "restore", ".", ".."},
		},
		{
			name:     "paths with spaces",
			paths:    []string{"file with spaces.txt", "another file.txt"},
			opts:     nil,
			wantArgs: []string{"git", "restore", "file with spaces.txt", "another file.txt"},
		},
		{
			name:     "paths with special characters",
			paths:    []string{"file@#$.txt", "file[].txt", "file().txt"},
			opts:     nil,
			wantArgs: []string{"git", "restore", "file@#$.txt", "file[].txt", "file().txt"},
		},
		{
			name:     "deeply nested paths",
			paths:    []string{"dir1/dir2/dir3/file.txt", "another/deep/path/file.txt"},
			opts:     nil,
			wantArgs: []string{"git", "restore", "dir1/dir2/dir3/file.txt", "another/deep/path/file.txt"},
		},
		{
			name:     "wildcard patterns",
			paths:    []string{"*.txt", "src/**/*.go"},
			opts:     nil,
			wantArgs: []string{"git", "restore", "*.txt", "src/**/*.go"},
		},
		{
			name:     "mixed path types with staged option",
			paths:    []string{".", "file.txt", "dir/", "*.md"},
			opts:     &RestoreOptions{Staged: true},
			wantArgs: []string{"git", "restore", "--staged", ".", "file.txt", "dir/", "*.md"},
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

			_ = client.Restore(tt.paths, tt.opts)

			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("got %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

// Test various source commit formats
func TestClient_Restore_SourceFormats(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		wantArgs []string
	}{
		{
			name:     "short commit hash",
			source:   "abc123",
			wantArgs: []string{"git", "restore", "--source", "abc123", "file.txt"},
		},
		{
			name:     "full commit hash",
			source:   "abcdef1234567890abcdef1234567890abcdef12",
			wantArgs: []string{"git", "restore", "--source", "abcdef1234567890abcdef1234567890abcdef12", "file.txt"},
		},
		{
			name:     "HEAD reference",
			source:   "HEAD",
			wantArgs: []string{"git", "restore", "--source", "HEAD", "file.txt"},
		},
		{
			name:     "HEAD with tilde",
			source:   "HEAD~1",
			wantArgs: []string{"git", "restore", "--source", "HEAD~1", "file.txt"},
		},
		{
			name:     "HEAD with caret",
			source:   "HEAD^",
			wantArgs: []string{"git", "restore", "--source", "HEAD^", "file.txt"},
		},
		{
			name:     "branch name",
			source:   "main",
			wantArgs: []string{"git", "restore", "--source", "main", "file.txt"},
		},
		{
			name:     "feature branch",
			source:   "feature/new-feature",
			wantArgs: []string{"git", "restore", "--source", "feature/new-feature", "file.txt"},
		},
		{
			name:     "remote branch",
			source:   "origin/main",
			wantArgs: []string{"git", "restore", "--source", "origin/main", "file.txt"},
		},
		{
			name:     "tag reference",
			source:   "v1.0.0",
			wantArgs: []string{"git", "restore", "--source", "v1.0.0", "file.txt"},
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

			_ = client.RestoreFromCommit(tt.source, "file.txt")

			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("got %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}
