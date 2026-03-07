package git

import (
	"os/exec"
	"slices"
	"testing"
)

func TestClient_Stash(t *testing.T) {
	var gotArgs []string
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo")
		},
	}

	err := client.Stash()
	if err != nil {
		t.Errorf("Stash() error = %v", err)
	}

	wantArgs := []string{"git", "stash"}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("Stash() gotArgs = %v, want %v", gotArgs, wantArgs)
	}
}

func TestClient_StashList(t *testing.T) {
	var gotArgs []string
	expectedOutput := "stash@{0}: WIP on main: 1234567 commit message\nstash@{1}: WIP on feature: abcdefg another commit"

	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo", "-n", expectedOutput)
		},
	}

	result, err := client.StashList()
	if err != nil {
		t.Errorf("StashList() error = %v", err)
	}

	wantArgs := []string{"git", "stash", "list"}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("StashList() gotArgs = %v, want %v", gotArgs, wantArgs)
	}

	if result != expectedOutput {
		t.Errorf("StashList() result = %v, want %v", result, expectedOutput)
	}
}

func TestClient_StashShow(t *testing.T) {
	tests := []struct {
		name     string
		stash    string
		wantArgs []string
	}{
		{
			name:     "show default stash",
			stash:    "",
			wantArgs: []string{"git", "stash", "show"},
		},
		{
			name:     "show specific stash",
			stash:    "stash@{1}",
			wantArgs: []string{"git", "stash", "show", "stash@{1}"},
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

			err := client.StashShow(tt.stash)
			if err != nil {
				t.Errorf("StashShow() error = %v", err)
			}

			if !slices.Equal(gotArgs, tt.wantArgs) {
				t.Errorf("StashShow() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestClient_StashApply(t *testing.T) {
	tests := []struct {
		name     string
		stash    string
		wantArgs []string
	}{
		{
			name:     "apply default stash",
			stash:    "",
			wantArgs: []string{"git", "stash", "apply"},
		},
		{
			name:     "apply specific stash",
			stash:    "stash@{1}",
			wantArgs: []string{"git", "stash", "apply", "stash@{1}"},
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

			err := client.StashApply(tt.stash)
			if err != nil {
				t.Errorf("StashApply() error = %v", err)
			}

			if !slices.Equal(gotArgs, tt.wantArgs) {
				t.Errorf("StashApply() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestClient_StashPop(t *testing.T) {
	tests := []struct {
		name     string
		stash    string
		wantArgs []string
	}{
		{
			name:     "pop default stash",
			stash:    "",
			wantArgs: []string{"git", "stash", "pop"},
		},
		{
			name:     "pop specific stash",
			stash:    "stash@{1}",
			wantArgs: []string{"git", "stash", "pop", "stash@{1}"},
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

			err := client.StashPop(tt.stash)
			if err != nil {
				t.Errorf("StashPop() error = %v", err)
			}

			if !slices.Equal(gotArgs, tt.wantArgs) {
				t.Errorf("StashPop() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestClient_StashPush(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		wantArgs []string
	}{
		{
			name:     "push without message",
			message:  "",
			wantArgs: []string{"git", "stash", "push"},
		},
		{
			name:     "push with message",
			message:  "update docs",
			wantArgs: []string{"git", "stash", "push", "-m", "update docs"},
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

			err := client.StashPush(tt.message)
			if err != nil {
				t.Errorf("StashPush() error = %v", err)
			}

			if !slices.Equal(gotArgs, tt.wantArgs) {
				t.Errorf("StashPush() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestClient_StashDrop(t *testing.T) {
	tests := []struct {
		name     string
		stash    string
		wantArgs []string
	}{
		{
			name:     "drop default stash",
			stash:    "",
			wantArgs: []string{"git", "stash", "drop"},
		},
		{
			name:     "drop specific stash",
			stash:    "stash@{1}",
			wantArgs: []string{"git", "stash", "drop", "stash@{1}"},
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

			err := client.StashDrop(tt.stash)
			if err != nil {
				t.Errorf("StashDrop() error = %v", err)
			}

			if !slices.Equal(gotArgs, tt.wantArgs) {
				t.Errorf("StashDrop() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestClient_StashClear(t *testing.T) {
	var gotArgs []string
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo")
		},
	}

	err := client.StashClear()
	if err != nil {
		t.Errorf("StashClear() error = %v", err)
	}

	wantArgs := []string{"git", "stash", "clear"}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("StashClear() gotArgs = %v, want %v", gotArgs, wantArgs)
	}
}

// Error case tests for better coverage
func TestClient_Stash_Error(t *testing.T) {
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			return exec.Command("false") // Command that always fails
		},
	}

	err := client.Stash()
	if err == nil {
		t.Error("Expected Stash to return an error")
	}
}

func TestClient_StashList_Error(t *testing.T) {
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			return exec.Command("false") // Command that always fails
		},
	}

	_, err := client.StashList()
	if err == nil {
		t.Error("Expected StashList to return an error")
	}
}

func TestClient_StashShow_Error(t *testing.T) {
	tests := []struct {
		name  string
		stash string
	}{
		{
			name:  "error with default stash",
			stash: "",
		},
		{
			name:  "error with specific stash",
			stash: "stash@{1}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				execCommand: func(name string, args ...string) *exec.Cmd {
					return exec.Command("false") // Command that always fails
				},
			}

			err := client.StashShow(tt.stash)
			if err == nil {
				t.Error("Expected StashShow to return an error")
			}
		})
	}
}

func TestClient_StashApply_Error(t *testing.T) {
	tests := []struct {
		name  string
		stash string
	}{
		{
			name:  "error with default stash",
			stash: "",
		},
		{
			name:  "error with specific stash",
			stash: "stash@{1}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				execCommand: func(name string, args ...string) *exec.Cmd {
					return exec.Command("false") // Command that always fails
				},
			}

			err := client.StashApply(tt.stash)
			if err == nil {
				t.Error("Expected StashApply to return an error")
			}
		})
	}
}

func TestClient_StashPop_Error(t *testing.T) {
	tests := []struct {
		name  string
		stash string
	}{
		{
			name:  "error with default stash",
			stash: "",
		},
		{
			name:  "error with specific stash",
			stash: "stash@{1}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				execCommand: func(name string, args ...string) *exec.Cmd {
					return exec.Command("false") // Command that always fails
				},
			}

			err := client.StashPop(tt.stash)
			if err == nil {
				t.Error("Expected StashPop to return an error")
			}
		})
	}
}

func TestClient_StashDrop_Error(t *testing.T) {
	tests := []struct {
		name  string
		stash string
	}{
		{
			name:  "error with default stash",
			stash: "",
		},
		{
			name:  "error with specific stash",
			stash: "stash@{1}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				execCommand: func(name string, args ...string) *exec.Cmd {
					return exec.Command("false") // Command that always fails
				},
			}

			err := client.StashDrop(tt.stash)
			if err == nil {
				t.Error("Expected StashDrop to return an error")
			}
		})
	}
}

func TestClient_StashClear_Error(t *testing.T) {
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			return exec.Command("false") // Command that always fails
		},
	}

	err := client.StashClear()
	if err == nil {
		t.Error("Expected StashClear to return an error")
	}
}

func TestClient_StashPush_Error(t *testing.T) {
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			return exec.Command("false") // Command that always fails
		},
	}

	err := client.StashPush("message")
	if err == nil {
		t.Error("Expected StashPush to return an error")
	}
}
