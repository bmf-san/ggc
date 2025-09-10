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

func TestClient_Rebase(t *testing.T) {
	tests := []struct {
		name     string
		upstream string
		wantArgs []string
	}{
		{
			name:     "rebase onto main",
			upstream: "main",
			wantArgs: []string{"git", "rebase", "main"},
		},
		{
			name:     "rebase onto origin/main",
			upstream: "origin/main",
			wantArgs: []string{"git", "rebase", "origin/main"},
		},
		{
			name:     "rebase onto feature branch",
			upstream: "origin/feature/test",
			wantArgs: []string{"git", "rebase", "origin/feature/test"},
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

			err := client.Rebase(tt.upstream)
			if err != nil {
				t.Errorf("Rebase() error = %v", err)
			}

			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("Rebase() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestClient_RebaseContinue(t *testing.T) {
	tests := []struct {
		name     string
		wantArgs []string
	}{
		{
			name:     "continue rebase",
			wantArgs: []string{"git", "rebase", "--continue"},
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

			err := client.RebaseContinue()
			if err != nil {
				t.Errorf("RebaseContinue() error = %v", err)
			}

			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("RebaseContinue() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestClient_RebaseAbort(t *testing.T) {
	tests := []struct {
		name     string
		wantArgs []string
	}{
		{
			name:     "abort rebase",
			wantArgs: []string{"git", "rebase", "--abort"},
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

			err := client.RebaseAbort()
			if err != nil {
				t.Errorf("RebaseAbort() error = %v", err)
			}

			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("RebaseAbort() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestClient_RebaseSkip(t *testing.T) {
	tests := []struct {
		name     string
		wantArgs []string
	}{
		{
			name:     "skip rebase",
			wantArgs: []string{"git", "rebase", "--skip"},
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

			err := client.RebaseSkip()
			if err != nil {
				t.Errorf("RebaseSkip() error = %v", err)
			}

			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("RebaseSkip() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

// Error case tests for better coverage
func TestClient_LogOneline_Error(t *testing.T) {
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			return exec.Command("false") // Command that always fails
		},
	}

	_, err := client.LogOneline("HEAD~3", "HEAD")
	if err == nil {
		t.Error("Expected LogOneline to return an error")
	}
}

func TestClient_RebaseInteractive_Error(t *testing.T) {
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			return exec.Command("false") // Command that always fails
		},
	}

	err := client.RebaseInteractive(3)
	if err == nil {
		t.Error("Expected RebaseInteractive to return an error")
	}
}

func TestClient_Rebase_Error(t *testing.T) {
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			return exec.Command("false") // Command that always fails
		},
	}

	err := client.Rebase("main")
	if err == nil {
		t.Error("Expected Rebase to return an error")
	}
}

func TestClient_RebaseContinue_Error(t *testing.T) {
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			return exec.Command("false") // Command that always fails
		},
	}

	err := client.RebaseContinue()
	if err == nil {
		t.Error("Expected RebaseContinue to return an error")
	}
}

func TestClient_RebaseAbort_Error(t *testing.T) {
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			return exec.Command("false") // Command that always fails
		},
	}

	err := client.RebaseAbort()
	if err == nil {
		t.Error("Expected RebaseAbort to return an error")
	}
}

func TestClient_RebaseSkip_Error(t *testing.T) {
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			return exec.Command("false") // Command that always fails
		},
	}

	err := client.RebaseSkip()
	if err == nil {
		t.Error("Expected RebaseSkip to return an error")
	}
}
