package git

import (
	"errors"
	"os/exec"
	"slices"
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
	if !slices.Equal(gotArgs, want) {
		t.Errorf("got %v, want %v", gotArgs, want)
	}
}

func TestClient_Commit(t *testing.T) {
	tests := []struct {
		name    string
		message string
		err     error
		wantErr bool
	}{
		{
			name:    "success_simple_message",
			message: "Add new feature",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "success_multiline_message",
			message: "Fix critical bug\n\nThis fixes the issue with user authentication",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "success_empty_message",
			message: "",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "success_message_with_special_chars",
			message: "Update: fix issue #123 & add tests (v2.0)",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "error_nothing_to_commit",
			message: "Attempt to commit",
			err:     errors.New("nothing to commit, working tree clean"),
			wantErr: true,
		},
		{
			name:    "error_pre_commit_hook_failed",
			message: "Add new code",
			err:     errors.New("pre-commit hook failed"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					expectedArgs := []string{"commit", "-m", tt.message}
					if name != "git" || len(arg) != len(expectedArgs) {
						t.Errorf("unexpected command: %s %v", name, arg)
					}
					for i, a := range arg {
						if a != expectedArgs[i] {
							t.Errorf("unexpected arg[%d]: got %s, want %s", i, a, expectedArgs[i])
						}
					}
					return helperCommand(t, "", tt.err)
				},
			}

			err := c.Commit(tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Commit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_CommitAmend(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr bool
	}{
		{
			name:    "success_amend_commit",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "error_no_commits",
			err:     errors.New("fatal: --amend: no previous commit"),
			wantErr: true,
		},
		{
			name:    "error_nothing_to_amend",
			err:     errors.New("nothing to commit, working tree clean"),
			wantErr: true,
		},
		{
			name:    "error_editor_failed",
			err:     errors.New("error: There was a problem with the editor"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					expectedArgs := []string{"commit", "--amend"}
					if name != "git" || len(arg) != len(expectedArgs) {
						t.Errorf("unexpected command: %s %v", name, arg)
					}
					for i, a := range arg {
						if a != expectedArgs[i] {
							t.Errorf("unexpected arg[%d]: got %s, want %s", i, a, expectedArgs[i])
						}
					}
					return helperCommand(t, "", tt.err)
				},
			}

			err := c.CommitAmend()
			if (err != nil) != tt.wantErr {
				t.Errorf("CommitAmend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_CommitAmendNoEdit(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr bool
	}{
		{
			name:    "success_amend_no_edit",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "error_no_commits",
			err:     errors.New("fatal: --amend: no previous commit"),
			wantErr: true,
		},
		{
			name:    "error_nothing_to_amend",
			err:     errors.New("nothing to commit, working tree clean"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					expectedArgs := []string{"commit", "--amend", "--no-edit"}
					if name != "git" || len(arg) != len(expectedArgs) {
						t.Errorf("unexpected command: %s %v", name, arg)
					}
					for i, a := range arg {
						if a != expectedArgs[i] {
							t.Errorf("unexpected arg[%d]: got %s, want %s", i, a, expectedArgs[i])
						}
					}
					return helperCommand(t, "", tt.err)
				},
			}

			err := c.CommitAmendNoEdit()
			if (err != nil) != tt.wantErr {
				t.Errorf("CommitAmendNoEdit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_CommitAmendWithMessage(t *testing.T) {
	tests := []struct {
		name    string
		message string
		err     error
		wantErr bool
	}{
		{
			name:    "success_amend_with_new_message",
			message: "Updated commit message",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "success_amend_with_multiline_message",
			message: "Fix bug\n\nThis addresses the critical issue found in production",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "success_amend_with_empty_message",
			message: "",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "error_no_commits",
			message: "New message",
			err:     errors.New("fatal: --amend: no previous commit"),
			wantErr: true,
		},
		{
			name:    "error_nothing_to_amend",
			message: "Update message",
			err:     errors.New("nothing to commit, working tree clean"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					expectedArgs := []string{"commit", "--amend", "-m", tt.message}
					if name != "git" || len(arg) != len(expectedArgs) {
						t.Errorf("unexpected command: %s %v", name, arg)
					}
					for i, a := range arg {
						if a != expectedArgs[i] {
							t.Errorf("unexpected arg[%d]: got %s, want %s", i, a, expectedArgs[i])
						}
					}
					return helperCommand(t, "", tt.err)
				},
			}

			err := c.CommitAmendWithMessage(tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommitAmendWithMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
