package git

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestClient_GetAheadBehindCount(t *testing.T) {
	tests := []struct {
		name           string
		branch         string
		upstream       string
		expectedOutput string
		wantArgs       []string
	}{
		{
			name:           "ahead behind count main vs origin/main",
			branch:         "main",
			upstream:       "origin/main",
			expectedOutput: "2\t1",
			wantArgs:       []string{"git", "rev-list", "--left-right", "--count", "main...origin/main"},
		},
		{
			name:           "ahead behind count feature vs main",
			branch:         "feature",
			upstream:       "main",
			expectedOutput: "3\t0",
			wantArgs:       []string{"git", "rev-list", "--left-right", "--count", "feature...main"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotArgs []string
			client := &Client{
				execCommand: func(name string, args ...string) *exec.Cmd {
					gotArgs = append([]string{name}, args...)
					return exec.Command("echo", "-n", tt.expectedOutput)
				},
			}

			result, err := client.GetAheadBehindCount(tt.branch, tt.upstream)
			if err != nil {
				t.Errorf("GetAheadBehindCount() error = %v", err)
			}

			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("GetAheadBehindCount() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}

			if result != tt.expectedOutput {
				t.Errorf("GetAheadBehindCount() result = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}

func TestClient_GetTagCommit(t *testing.T) {
	tests := []struct {
		name           string
		tagName        string
		expectedOutput string
		wantArgs       []string
	}{
		{
			name:           "get commit for tag v1.0.0",
			tagName:        "v1.0.0",
			expectedOutput: "abc123def456",
			wantArgs:       []string{"git", "rev-list", "-n", "1", "v1.0.0"},
		},
		{
			name:           "get commit for tag v2.1.0",
			tagName:        "v2.1.0",
			expectedOutput: "def456abc123",
			wantArgs:       []string{"git", "rev-list", "-n", "1", "v2.1.0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotArgs []string
			client := &Client{
				execCommand: func(name string, args ...string) *exec.Cmd {
					gotArgs = append([]string{name}, args...)
					return exec.Command("echo", "-n", tt.expectedOutput)
				},
			}

			result, err := client.GetTagCommit(tt.tagName)
			if err != nil {
				t.Errorf("GetTagCommit() error = %v", err)
			}

			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("GetTagCommit() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}

			if result != tt.expectedOutput {
				t.Errorf("GetTagCommit() result = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}

// Error case tests for better coverage
func TestClient_GetAheadBehindCount_Error(t *testing.T) {
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			return exec.Command("false") // Command that always fails
		},
	}

	_, err := client.GetAheadBehindCount("main", "origin/main")
	if err == nil {
		t.Error("Expected GetAheadBehindCount to return an error")
	}
}

func TestClient_GetTagCommit_Error(t *testing.T) {
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			return exec.Command("false") // Command that always fails
		},
	}

	_, err := client.GetTagCommit("v1.0.0")
	if err == nil {
		t.Error("Expected GetTagCommit to return an error")
	}
}
