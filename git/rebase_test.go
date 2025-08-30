package git

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestClient_LogOneline(t *testing.T) {
	tests := []struct {
		name           string
		from           string
		to             string
		expectedOutput string
		wantArgs       []string
	}{
		{
			name:           "log between commits",
			from:           "HEAD~3",
			to:             "HEAD",
			expectedOutput: "abc1234 commit 1\ndef5678 commit 2\n",
			wantArgs:       []string{"git", "log", "--oneline", "--reverse", "HEAD~3..HEAD"},
		},
		{
			name:           "log between branches",
			from:           "main",
			to:             "feature",
			expectedOutput: "123abcd feature commit 1\n456efgh feature commit 2\n",
			wantArgs:       []string{"git", "log", "--oneline", "--reverse", "main..feature"},
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

			result, err := client.LogOneline(tt.from, tt.to)
			if err != nil {
				t.Errorf("LogOneline() error = %v", err)
			}

			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("LogOneline() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}

			if result != tt.expectedOutput {
				t.Errorf("LogOneline() result = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}

func TestClient_RebaseInteractive(t *testing.T) {
	tests := []struct {
		name        string
		commitCount int
		wantArgs    []string
	}{
		{
			name:        "rebase last 3 commits",
			commitCount: 3,
			wantArgs:    []string{"git", "rebase", "-i", "HEAD~3"},
		},
		{
			name:        "rebase last commit",
			commitCount: 1,
			wantArgs:    []string{"git", "rebase", "-i", "HEAD~1"},
		},
		{
			name:        "rebase last 5 commits",
			commitCount: 5,
			wantArgs:    []string{"git", "rebase", "-i", "HEAD~5"},
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

			err := client.RebaseInteractive(tt.commitCount)
			if err != nil {
				t.Errorf("RebaseInteractive() error = %v", err)
			}

			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("RebaseInteractive() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestClient_GetUpstreamBranch(t *testing.T) {
	tests := []struct {
		name           string
		branch         string
		mockOutput     string
		mockError      bool
		expectedResult string
		wantArgs       []string
	}{
		{
			name:           "get upstream for main branch",
			branch:         "main",
			mockOutput:     "origin/main",
			mockError:      false,
			expectedResult: "origin/main",
			wantArgs:       []string{"git", "rev-parse", "--abbrev-ref", "main@{upstream}"},
		},
		{
			name:           "get upstream for feature branch",
			branch:         "feature/test",
			mockOutput:     "origin/feature/test",
			mockError:      false,
			expectedResult: "origin/feature/test",
			wantArgs:       []string{"git", "rev-parse", "--abbrev-ref", "feature/test@{upstream}"},
		},
		{
			name:           "no upstream set - returns default",
			branch:         "local-branch",
			mockOutput:     "",
			mockError:      true,
			expectedResult: "main",
			wantArgs:       []string{"git", "rev-parse", "--abbrev-ref", "local-branch@{upstream}"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotArgs []string
			client := &Client{
				execCommand: func(name string, args ...string) *exec.Cmd {
					gotArgs = append([]string{name}, args...)
					if tt.mockError {
						return exec.Command("false")
					}
					return exec.Command("echo", "-n", tt.mockOutput)
				},
			}

			result, err := client.GetUpstreamBranch(tt.branch)
			if err != nil {
				t.Errorf("GetUpstreamBranch() error = %v", err)
			}

			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("GetUpstreamBranch() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}

			if result != tt.expectedResult {
				t.Errorf("GetUpstreamBranch() result = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}
