package git

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestClient_Fetch(t *testing.T) {
	tests := []struct {
		name     string
		prune    bool
		wantArgs []string
	}{
		{
			name:     "fetch without prune",
			prune:    false,
			wantArgs: []string{"git", "fetch"},
		},
		{
			name:     "fetch with prune",
			prune:    true,
			wantArgs: []string{"git", "fetch", "--prune"},
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

			err := client.Fetch(tt.prune)
			if err != nil {
				t.Errorf("Fetch() error = %v", err)
			}

			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("Fetch() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}
