package git

import (
	"errors"
	"os/exec"
	"slices"
	"testing"
)

func TestClient_Show(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantArgs []string
	}{
		{
			name:     "no args defaults to HEAD",
			args:     nil,
			wantArgs: []string{"git", "show"},
		},
		{
			name:     "single ref",
			args:     []string{"HEAD~1"},
			wantArgs: []string{"git", "show", "HEAD~1"},
		},
		{
			name:     "stat option with ref",
			args:     []string{"--stat", "abc123"},
			wantArgs: []string{"git", "show", "--stat", "abc123"},
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

			if err := client.Show(tt.args); err != nil {
				t.Errorf("Show() error = %v", err)
			}

			if !slices.Equal(gotArgs, tt.wantArgs) {
				t.Errorf("Show() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestClient_Show_Error(t *testing.T) {
	client := &Client{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return helperCommand(t, "", errors.New("boom"))
		},
	}
	if err := client.Show([]string{"deadbeef"}); err == nil {
		t.Error("expected error, got nil")
	}
}
