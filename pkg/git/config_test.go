package git

import (
	"errors"
	"os/exec"
	"testing"
)

func TestClient_ConfigGet(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		output  string
		err     error
		want    string
		wantErr bool
	}{
		{
			name:    "success_user_name",
			key:     "user.name",
			output:  "John Doe",
			err:     nil,
			want:    "John Doe",
			wantErr: false,
		},
		{
			name:    "success_user_email",
			key:     "user.email",
			output:  "john.doe@example.com",
			err:     nil,
			want:    "john.doe@example.com",
			wantErr: false,
		},
		{
			name:    "success_with_whitespace",
			key:     "core.editor",
			output:  "  vim  \n",
			err:     nil,
			want:    "vim",
			wantErr: false,
		},
		{
			name:    "error_key_not_found",
			key:     "nonexistent.key",
			output:  "",
			err:     errors.New("exit status 1"),
			want:    "",
			wantErr: true,
		},
		{
			name:    "error_invalid_key",
			key:     "invalid..key",
			output:  "",
			err:     errors.New("error: invalid key"),
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					expectedArgs := []string{"config", tt.key}
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

			got, err := c.ConfigGet(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConfigGet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ConfigSet(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		value   string
		err     error
		wantErr bool
	}{
		{
			name:    "success_set_user_name",
			key:     "user.name",
			value:   "John Doe",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "success_set_user_email",
			key:     "user.email",
			value:   "john.doe@example.com",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "success_set_empty_value",
			key:     "test.key",
			value:   "",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "error_invalid_key",
			key:     "invalid..key",
			value:   "value",
			err:     errors.New("error: invalid key"),
			wantErr: true,
		},
		{
			name:    "error_permission_denied",
			key:     "user.name",
			value:   "John Doe",
			err:     errors.New("permission denied"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					expectedArgs := []string{"config", tt.key, tt.value}
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

			if err := c.ConfigSet(tt.key, tt.value); (err != nil) != tt.wantErr {
				t.Errorf("ConfigSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_ConfigGetGlobal(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		output  string
		err     error
		want    string
		wantErr bool
	}{
		{
			name:    "success_global_user_name",
			key:     "user.name",
			output:  "Global User",
			err:     nil,
			want:    "Global User",
			wantErr: false,
		},
		{
			name:    "success_global_user_email",
			key:     "user.email",
			output:  "global@example.com",
			err:     nil,
			want:    "global@example.com",
			wantErr: false,
		},
		{
			name:    "error_key_not_found",
			key:     "nonexistent.key",
			output:  "",
			err:     errors.New("exit status 1"),
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					expectedArgs := []string{"config", "--global", tt.key}
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

			got, err := c.ConfigGetGlobal(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigGetGlobal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConfigGetGlobal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ConfigSetGlobal(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		value   string
		err     error
		wantErr bool
	}{
		{
			name:    "success_set_global_user_name",
			key:     "user.name",
			value:   "Global User",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "success_set_global_user_email",
			key:     "user.email",
			value:   "global@example.com",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "error_invalid_key",
			key:     "invalid..key",
			value:   "value",
			err:     errors.New("error: invalid key"),
			wantErr: true,
		},
		{
			name:    "error_permission_denied",
			key:     "user.name",
			value:   "Global User",
			err:     errors.New("permission denied"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					expectedArgs := []string{"config", "--global", tt.key, tt.value}
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

			if err := c.ConfigSetGlobal(tt.key, tt.value); (err != nil) != tt.wantErr {
				t.Errorf("ConfigSetGlobal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
