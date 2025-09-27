package git

import (
	"os/exec"
	"slices"
	"testing"
)

func TestClient_Push(t *testing.T) {
	cases := []struct {
		name     string
		force    bool
		wantArgs []string
	}{
		{
			name:     "push",
			force:    false,
			wantArgs: []string{"git", "push", "origin", "main"},
		},
		{
			name:     "force push",
			force:    true,
			wantArgs: []string{"git", "push", "origin", "main", "--force-with-lease"},
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
				GetCurrentBranchFunc: func() (string, error) {
					return "main", nil
				},
			}

			_ = client.Push(tc.force)
			if !slices.Equal(gotArgs, tc.wantArgs) {
				t.Errorf("got %v, want %v", gotArgs, tc.wantArgs)
			}
		})
	}
}
