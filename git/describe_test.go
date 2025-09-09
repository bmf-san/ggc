package git

import (
	"errors"
	"os/exec"
	"testing"
)

func TestClient_GetVersion(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		err     error
		want    string
		wantErr bool
	}{
		{
			name:    "success_with_tag",
			output:  "v1.2.3",
			err:     nil,
			want:    "v1.2.3",
			wantErr: false,
		},
		{
			name:    "success_with_tag_and_commits",
			output:  "v1.2.3-5-g1234567",
			err:     nil,
			want:    "v1.2.3-5-g1234567",
			wantErr: false,
		},
		{
			name:    "success_with_dirty_flag",
			output:  "v1.2.3-dirty",
			err:     nil,
			want:    "v1.2.3-dirty",
			wantErr: false,
		},
		{
			name:    "success_with_tag_commits_and_dirty",
			output:  "v1.2.3-5-g1234567-dirty",
			err:     nil,
			want:    "v1.2.3-5-g1234567-dirty",
			wantErr: false,
		},
		{
			name:    "success_commit_hash_only",
			output:  "1234567",
			err:     nil,
			want:    "1234567",
			wantErr: false,
		},
		{
			name:    "success_with_whitespace",
			output:  "  v1.2.3  \n",
			err:     nil,
			want:    "v1.2.3",
			wantErr: false,
		},
		{
			name:    "error_no_repository",
			output:  "",
			err:     errors.New("fatal: not a git repository"),
			want:    "dev",
			wantErr: false, // Should return "dev" as fallback
		},
		{
			name:    "error_no_tags",
			output:  "",
			err:     errors.New("fatal: No names found"),
			want:    "dev",
			wantErr: false, // Should return "dev" as fallback
		},
		{
			name:    "error_permission_denied",
			output:  "",
			err:     errors.New("permission denied"),
			want:    "dev",
			wantErr: false, // Should return "dev" as fallback
		},
		{
			name:    "success_empty_output_with_no_error",
			output:  "",
			err:     nil,
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					expectedArgs := []string{"describe", "--tags", "--always", "--dirty"}
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

			got, err := c.GetVersion()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
