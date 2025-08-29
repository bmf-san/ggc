package git

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestClient_TagList(t *testing.T) {
	tests := []struct {
		name     string
		pattern  []string
		wantArgs []string
	}{
		{
			name:     "list all tags",
			pattern:  []string{},
			wantArgs: []string{"git", "tag", "--sort=-version:refname"},
		},
		{
			name:     "list tags with pattern",
			pattern:  []string{"v*"},
			wantArgs: []string{"git", "tag", "--sort=-version:refname", "-l", "v*"},
		},
		{
			name:     "list tags with multiple patterns",
			pattern:  []string{"v1.*", "v2.*"},
			wantArgs: []string{"git", "tag", "--sort=-version:refname", "-l", "v1.*", "v2.*"},
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

			err := client.TagList(tt.pattern)
			if err != nil {
				t.Errorf("TagList() error = %v", err)
			}

			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("TagList() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestClient_TagCreate(t *testing.T) {
	tests := []struct {
		name     string
		tagName  string
		commit   string
		wantArgs []string
	}{
		{
			name:     "create tag at HEAD",
			tagName:  "v1.0.0",
			commit:   "",
			wantArgs: []string{"git", "tag", "v1.0.0"},
		},
		{
			name:     "create tag at specific commit",
			tagName:  "v1.0.1",
			commit:   "abc1234",
			wantArgs: []string{"git", "tag", "v1.0.1", "abc1234"},
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

			err := client.TagCreate(tt.tagName, tt.commit)
			if err != nil {
				t.Errorf("TagCreate() error = %v", err)
			}

			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("TagCreate() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestClient_TagCreateAnnotated(t *testing.T) {
	tests := []struct {
		name     string
		tagName  string
		message  string
		wantArgs []string
	}{
		{
			name:     "create annotated tag without message",
			tagName:  "v1.0.0",
			message:  "",
			wantArgs: []string{"git", "tag", "-a", "v1.0.0"},
		},
		{
			name:     "create annotated tag with message",
			tagName:  "v1.0.1",
			message:  "Release version 1.0.1",
			wantArgs: []string{"git", "tag", "-a", "v1.0.1", "-m", "Release version 1.0.1"},
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

			err := client.TagCreateAnnotated(tt.tagName, tt.message)
			if err != nil {
				t.Errorf("TagCreateAnnotated() error = %v", err)
			}

			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("TagCreateAnnotated() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestClient_TagDelete(t *testing.T) {
	tests := []struct {
		name      string
		tagNames  []string
		wantCalls int
	}{
		{
			name:      "delete single tag",
			tagNames:  []string{"v1.0.0"},
			wantCalls: 1,
		},
		{
			name:      "delete multiple tags",
			tagNames:  []string{"v1.0.0", "v1.0.1", "v2.0.0"},
			wantCalls: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var callCount int
			var lastArgs []string
			client := &Client{
				execCommand: func(name string, args ...string) *exec.Cmd {
					callCount++
					lastArgs = append([]string{name}, args...)
					return exec.Command("echo")
				},
			}

			err := client.TagDelete(tt.tagNames)
			if err != nil {
				t.Errorf("TagDelete() error = %v", err)
			}

			if callCount != tt.wantCalls {
				t.Errorf("TagDelete() callCount = %v, want %v", callCount, tt.wantCalls)
			}

			// Check the last call for the pattern
			if len(tt.tagNames) > 0 {
				expectedLastArgs := []string{"git", "tag", "-d", tt.tagNames[len(tt.tagNames)-1]}
				if !reflect.DeepEqual(lastArgs, expectedLastArgs) {
					t.Errorf("TagDelete() lastArgs = %v, want %v", lastArgs, expectedLastArgs)
				}
			}
		})
	}
}

func TestClient_TagPush(t *testing.T) {
	var gotArgs []string
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo")
		},
	}

	err := client.TagPush("origin", "v1.0.0")
	if err != nil {
		t.Errorf("TagPush() error = %v", err)
	}

	wantArgs := []string{"git", "push", "origin", "v1.0.0"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("TagPush() gotArgs = %v, want %v", gotArgs, wantArgs)
	}
}

func TestClient_TagPushAll(t *testing.T) {
	var gotArgs []string
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo")
		},
	}

	err := client.TagPushAll("origin")
	if err != nil {
		t.Errorf("TagPushAll() error = %v", err)
	}

	wantArgs := []string{"git", "push", "origin", "--tags"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("TagPushAll() gotArgs = %v, want %v", gotArgs, wantArgs)
	}
}

func TestClient_TagShow(t *testing.T) {
	var gotArgs []string
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo")
		},
	}

	err := client.TagShow("v1.0.0")
	if err != nil {
		t.Errorf("TagShow() error = %v", err)
	}

	wantArgs := []string{"git", "show", "v1.0.0"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("TagShow() gotArgs = %v, want %v", gotArgs, wantArgs)
	}
}

func TestClient_GetLatestTag(t *testing.T) {
	var gotArgs []string
	expectedOutput := "v1.2.3"

	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo", "-n", expectedOutput)
		},
	}

	result, err := client.GetLatestTag()
	if err != nil {
		t.Errorf("GetLatestTag() error = %v", err)
	}

	wantArgs := []string{"git", "describe", "--tags", "--abbrev=0"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("GetLatestTag() gotArgs = %v, want %v", gotArgs, wantArgs)
	}

	if result != expectedOutput {
		t.Errorf("GetLatestTag() result = %v, want %v", result, expectedOutput)
	}
}

func TestClient_TagExists(t *testing.T) {
	tests := []struct {
		name       string
		tagName    string
		mockOutput string
		mockError  bool
		expected   bool
	}{
		{
			name:       "tag exists",
			tagName:    "v1.0.0",
			mockOutput: "v1.0.0",
			mockError:  false,
			expected:   true,
		},
		{
			name:       "tag does not exist",
			tagName:    "v2.0.0",
			mockOutput: "",
			mockError:  false,
			expected:   false,
		},
		{
			name:       "command error",
			tagName:    "v3.0.0",
			mockOutput: "",
			mockError:  true,
			expected:   false,
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

			result := client.TagExists(tt.tagName)

			wantArgs := []string{"git", "tag", "-l", tt.tagName}
			if !reflect.DeepEqual(gotArgs, wantArgs) {
				t.Errorf("TagExists() gotArgs = %v, want %v", gotArgs, wantArgs)
			}

			if result != tt.expected {
				t.Errorf("TagExists() result = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestClient_GetTagCommit(t *testing.T) {
	var gotArgs []string
	expectedOutput := "abc1234567890def"

	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo", "-n", expectedOutput)
		},
	}

	result, err := client.GetTagCommit("v1.0.0")
	if err != nil {
		t.Errorf("GetTagCommit() error = %v", err)
	}

	wantArgs := []string{"git", "rev-list", "-n", "1", "v1.0.0"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("GetTagCommit() gotArgs = %v, want %v", gotArgs, wantArgs)
	}

	if result != expectedOutput {
		t.Errorf("GetTagCommit() result = %v, want %v", result, expectedOutput)
	}
}
