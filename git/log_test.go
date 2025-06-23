package git

import (
	"os/exec"
	"strings"
	"testing"
)

func TestClient_LogSimple(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr bool
	}{
		{
			name:    "正常系：シンプルログ表示",
			err:     nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					if name != "git" || !strings.Contains(strings.Join(arg, " "), "log --oneline --graph --decorate -10") {
						t.Errorf("unexpected command: %s %v", name, arg)
					}
					return helperCommand(t, "", tt.err)
				},
			}

			if err := c.LogSimple(); (err != nil) != tt.wantErr {
				t.Errorf("LogSimple() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
