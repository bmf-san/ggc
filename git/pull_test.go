package git

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestClient_Pull(t *testing.T) {
	cases := []struct {
		name     string
		rebase   bool
		wantArgs []string
	}{
		{
			name:     "pull",
			rebase:   false,
			wantArgs: []string{"git", "pull"},
		},
		{
			name:     "pull with rebase",
			rebase:   true,
			wantArgs: []string{"git", "pull", "--rebase"},
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
			_ = client.Pull(tc.rebase)
			if !reflect.DeepEqual(gotArgs, tc.wantArgs) {
				t.Errorf("got %v, want %v", gotArgs, tc.wantArgs)
			}
		})
	}
}
