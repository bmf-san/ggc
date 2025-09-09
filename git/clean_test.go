package git

import (
	"errors"
	"os/exec"
	"reflect"
	"testing"
)

func TestClient_CleanFiles(t *testing.T) {
	var gotArgs []string
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo")
		},
	}

	_ = client.CleanFiles()
	want := []string{"git", "clean", "-fd"}
	if !reflect.DeepEqual(gotArgs, want) {
		t.Errorf("got %v, want %v", gotArgs, want)
	}
}

func TestClient_CleanDirs(t *testing.T) {
	var gotArgs []string
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo")
		},
	}

	_ = client.CleanDirs()
	want := []string{"git", "clean", "-fdx"}
	if !reflect.DeepEqual(gotArgs, want) {
		t.Errorf("got %v, want %v", gotArgs, want)
	}
}

func TestClient_CleanDryRun(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		err     error
		want    string
		wantErr bool
	}{
		{
			name:    "success_with_files_to_clean",
			output:  "Would remove file1.txt\nWould remove file2.log\n",
			err:     nil,
			want:    "Would remove file1.txt\nWould remove file2.log\n",
			wantErr: false,
		},
		{
			name:    "success_no_files_to_clean",
			output:  "",
			err:     nil,
			want:    "",
			wantErr: false,
		},
		{
			name:    "success_with_directories",
			output:  "Would remove dir1/\nWould remove dir2/file.txt\n",
			err:     nil,
			want:    "Would remove dir1/\nWould remove dir2/file.txt\n",
			wantErr: false,
		},
		{
			name:    "error_not_git_repository",
			output:  "",
			err:     errors.New("fatal: not a git repository"),
			want:    "",
			wantErr: true,
		},
		{
			name:    "error_permission_denied",
			output:  "",
			err:     errors.New("permission denied"),
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					expectedArgs := []string{"clean", "-nd"}
					if name != "git" || len(arg) != len(expectedArgs) {
						t.Errorf("unexpected command: %s %v", name, arg)
					}
					for i, a := range arg {
						if a != expectedArgs[i] {
							t.Errorf("unexpected arg[%d]: got %s, want %s", i, a, expectedArgs[i])
						}
					}
					return helperCommand(t, tt.output, tt.err)
				},
			}

			got, err := c.CleanDryRun()
			if (err != nil) != tt.wantErr {
				t.Errorf("CleanDryRun() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CleanDryRun() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_CleanFilesForce(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		err      error
		wantArgs []string
		wantErr  bool
	}{
		{
			name:     "success_single_file",
			files:    []string{"file1.txt"},
			err:      nil,
			wantArgs: []string{"git", "clean", "-f", "--", "file1.txt"},
			wantErr:  false,
		},
		{
			name:     "success_multiple_files",
			files:    []string{"file1.txt", "file2.log", "dir/file3.tmp"},
			err:      nil,
			wantArgs: []string{"git", "clean", "-f", "--", "file1.txt", "file2.log", "dir/file3.tmp"},
			wantErr:  false,
		},
		{
			name:     "success_no_files",
			files:    []string{},
			err:      nil,
			wantArgs: nil, // No command should be executed
			wantErr:  false,
		},
		{
			name:     "error_permission_denied",
			files:    []string{"file1.txt"},
			err:      errors.New("permission denied"),
			wantArgs: []string{"git", "clean", "-f", "--", "file1.txt"},
			wantErr:  true,
		},
		{
			name:     "error_file_not_found",
			files:    []string{"nonexistent.txt"},
			err:      errors.New("pathspec 'nonexistent.txt' did not match any files"),
			wantArgs: []string{"git", "clean", "-f", "--", "nonexistent.txt"},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotArgs []string
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					gotArgs = append([]string{name}, arg...)
					if tt.wantArgs != nil {
						if name != "git" || len(arg) != len(tt.wantArgs)-1 {
							t.Errorf("unexpected command: %s %v", name, arg)
						}
						for i, a := range arg {
							if a != tt.wantArgs[i+1] {
								t.Errorf("unexpected arg[%d]: got %s, want %s", i, a, tt.wantArgs[i+1])
							}
						}
					}
					return helperCommand(t, "", tt.err)
				},
			}

			err := c.CleanFilesForce(tt.files)
			if (err != nil) != tt.wantErr {
				t.Errorf("CleanFilesForce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantArgs == nil && gotArgs != nil {
				t.Errorf("CleanFilesForce() should not execute command for empty files, but got: %v", gotArgs)
			} else if tt.wantArgs != nil && !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("CleanFilesForce() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}
