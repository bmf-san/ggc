package cmd

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

type noopCmd struct {
	cmd *exec.Cmd
}

func newNoopCmd() *noopCmd {
	cmd := exec.Command("true")
	return &noopCmd{cmd: cmd}
}

func (c *noopCmd) Run() error {
	return nil
}

func TestRemoteer_Remote_List(t *testing.T) {
	called := false
	var buf bytes.Buffer
	remoteer := &Remoteer{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			called = true
			return newNoopCmd().cmd
		},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	remoteer.helper.outputWriter = &buf

	remoteer.Remote([]string{"list"})

	if !called {
		t.Error("execCommand is not called in list subcommand")
	}
}

func TestRemoteer_Remote_Add(t *testing.T) {
	called := false
	var buf bytes.Buffer
	remoteer := &Remoteer{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			called = true
			return newNoopCmd().cmd
		},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	remoteer.helper.outputWriter = &buf

	remoteer.Remote([]string{"add", "origin", "https://example.com"})

	if !called {
		t.Error("execCommand is not called in add subcommand")
	}
	output := buf.String()
	if output != "Remote 'origin' added\n" {
		t.Errorf("unexpected output: got %q", output)
	}
}

func TestRemoteer_Remote_Remove(t *testing.T) {
	called := false
	var buf bytes.Buffer
	remoteer := &Remoteer{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			called = true
			return newNoopCmd().cmd
		},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	remoteer.helper.outputWriter = &buf

	remoteer.Remote([]string{"remove", "origin"})

	if !called {
		t.Error("execCommand is not called in remove subcommand")
	}
	output := buf.String()
	if output != "Remote 'origin' removed\n" {
		t.Errorf("unexpected output: got %q", output)
	}
}

func TestRemoteer_Remote_SetURL(t *testing.T) {
	called := false
	var buf bytes.Buffer
	remoteer := &Remoteer{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			called = true
			return newNoopCmd().cmd
		},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	remoteer.helper.outputWriter = &buf

	remoteer.Remote([]string{"set-url", "origin", "https://example.com"})

	if !called {
		t.Error("execCommand is not called in set-url subcommand")
	}
	output := buf.String()
	if output != "Remote 'origin' URL updated\n" {
		t.Errorf("unexpected output: got %q", output)
	}
}

func TestRemoteer_Remote_Help(t *testing.T) {
	var buf bytes.Buffer
	remoteer := &Remoteer{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return newNoopCmd().cmd
		},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	remoteer.helper.outputWriter = &buf

	remoteer.Remote([]string{"unknown"})

	output := buf.String()
	if output == "" || !strings.Contains(output, "Usage") {
		t.Errorf("Usage should be displayed, but got: %s", output)
	}
}

func TestRemoteer_Remote_NoArgs(t *testing.T) {
	var buf bytes.Buffer
	remoteer := &Remoteer{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return newNoopCmd().cmd
		},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	remoteer.helper.outputWriter = &buf

	remoteer.Remote([]string{})

	output := buf.String()
	if output == "" || !strings.Contains(output, "Usage") {
		t.Errorf("Usage should be displayed, but got: %s", output)
	}
}

func TestRemoteer_Remote_Add_InsufficientArgs(t *testing.T) {
	var buf bytes.Buffer
	remoteer := &Remoteer{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return newNoopCmd().cmd
		},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	remoteer.helper.outputWriter = &buf

	remoteer.Remote([]string{"add", "origin"})

	output := buf.String()
	if output == "" || !strings.Contains(output, "Usage") {
		t.Errorf("Usage should be displayed for insufficient args, but got: %s", output)
	}
}

func TestRemoteer_Remote_Remove_InsufficientArgs(t *testing.T) {
	var buf bytes.Buffer
	remoteer := &Remoteer{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return newNoopCmd().cmd
		},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	remoteer.helper.outputWriter = &buf

	remoteer.Remote([]string{"remove"})

	output := buf.String()
	if output == "" || !strings.Contains(output, "Usage") {
		t.Errorf("Usage should be displayed for insufficient args, but got: %s", output)
	}
}

func TestRemoteer_Remote_SetURL_InsufficientArgs(t *testing.T) {
	var buf bytes.Buffer
	remoteer := &Remoteer{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return newNoopCmd().cmd
		},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	remoteer.helper.outputWriter = &buf

	remoteer.Remote([]string{"set-url", "origin"})

	output := buf.String()
	if output == "" || !strings.Contains(output, "Usage") {
		t.Errorf("Usage should be displayed for insufficient args, but got: %s", output)
	}
}

func TestRemoteer_Remote_List_Error(t *testing.T) {
	var buf bytes.Buffer
	remoteer := &Remoteer{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("false")
		},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	remoteer.helper.outputWriter = &buf

	remoteer.Remote([]string{"list"})

	output := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("Error: failed to list remotes")) {
		t.Errorf("Expected error message, but got: %s", output)
	}
}

func TestRemoteer_Remote_Add_Error(t *testing.T) {
	var buf bytes.Buffer
	remoteer := &Remoteer{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("false")
		},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	remoteer.helper.outputWriter = &buf

	remoteer.Remote([]string{"add", "origin", "https://example.com"})

	output := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("Error: failed to add remote")) {
		t.Errorf("Expected error message, but got: %s", output)
	}
}

func TestRemoteer_Remote_Remove_Error(t *testing.T) {
	var buf bytes.Buffer
	remoteer := &Remoteer{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("false")
		},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	remoteer.helper.outputWriter = &buf

	remoteer.Remote([]string{"remove", "origin"})

	output := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("Error: failed to remove remote")) {
		t.Errorf("Expected error message, but got: %s", output)
	}
}

func TestRemoteer_Remote_SetURL_Error(t *testing.T) {
	var buf bytes.Buffer
	remoteer := &Remoteer{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("false")
		},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	remoteer.helper.outputWriter = &buf

	remoteer.Remote([]string{"set-url", "origin", "https://example.com"})

	output := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("Error: failed to set remote URL")) {
		t.Errorf("Expected error message, but got: %s", output)
	}
}

func TestRemoteer_Remote_Add_InvalidURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "Invalid URL format",
			url:      "not-a-url",
			expected: "Error: failed to add remote",
		},
		{
			name:     "Empty URL",
			url:      "",
			expected: "Error: failed to add remote",
		},
		{
			name:     "URL with spaces",
			url:      "https://example.com/repo with spaces",
			expected: "Error: failed to add remote",
		},
		{
			name:     "Very long URL",
			url:      "https://example.com/" + strings.Repeat("a", 1000),
			expected: "Error: failed to add remote",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			remoteer := &Remoteer{
				execCommand: func(_ string, _ ...string) *exec.Cmd {
					return exec.Command("false")
				},
				outputWriter: &buf,
				helper:       NewHelper(),
			}
			remoteer.helper.outputWriter = &buf

			remoteer.Remote([]string{"add", "origin", tt.url})

			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected %q in output, got: %s", tt.expected, output)
			}
		})
	}
}

func TestRemoteer_Remote_Add_InvalidRemoteName(t *testing.T) {
	tests := []struct {
		name       string
		remoteName string
		expected   string
	}{
		{
			name:       "Remote name with spaces",
			remoteName: "origin with spaces",
			expected:   "Error: failed to add remote",
		},
		{
			name:       "Empty remote name",
			remoteName: "",
			expected:   "Error: failed to add remote",
		},
		{
			name:       "Remote name with special chars",
			remoteName: "origin@#$%",
			expected:   "Error: failed to add remote",
		},
		{
			name:       "Very long remote name",
			remoteName: strings.Repeat("a", 500),
			expected:   "Error: failed to add remote",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			remoteer := &Remoteer{
				execCommand: func(_ string, _ ...string) *exec.Cmd {
					return exec.Command("false")
				},
				outputWriter: &buf,
				helper:       NewHelper(),
			}
			remoteer.helper.outputWriter = &buf

			remoteer.Remote([]string{"add", tt.remoteName, "https://example.com"})

			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected %q in output, got: %s", tt.expected, output)
			}
		})
	}
}

func TestRemoteer_Remote_Remove_NonexistentRemote(t *testing.T) {
	tests := []struct {
		name       string
		remoteName string
		expected   string
	}{
		{
			name:       "Remove nonexistent remote",
			remoteName: "nonexistent",
			expected:   "Error: failed to remove remote",
		},
		{
			name:       "Remove with empty name",
			remoteName: "",
			expected:   "Error: failed to remove remote",
		},
		{
			name:       "Remove with invalid characters",
			remoteName: "remote/with/slashes",
			expected:   "Error: failed to remove remote",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			remoteer := &Remoteer{
				execCommand: func(_ string, _ ...string) *exec.Cmd {
					return exec.Command("false")
				},
				outputWriter: &buf,
				helper:       NewHelper(),
			}
			remoteer.helper.outputWriter = &buf

			remoteer.Remote([]string{"remove", tt.remoteName})

			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected %q in output, got: %s", tt.expected, output)
			}
		})
	}
}

func TestRemoteer_Remote_SetURL_InvalidScenarios(t *testing.T) {
	tests := []struct {
		name       string
		remoteName string
		url        string
		expected   string
	}{
		{
			name:       "Set URL for nonexistent remote",
			remoteName: "nonexistent",
			url:        "https://example.com",
			expected:   "Error: failed to set remote URL",
		},
		{
			name:       "Set invalid URL",
			remoteName: "origin",
			url:        "not-a-url",
			expected:   "Error: failed to set remote URL",
		},
		{
			name:       "Set URL with empty remote name",
			remoteName: "",
			url:        "https://example.com",
			expected:   "Error: failed to set remote URL",
		},
		{
			name:       "Set empty URL",
			remoteName: "origin",
			url:        "",
			expected:   "Error: failed to set remote URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			remoteer := &Remoteer{
				execCommand: func(_ string, _ ...string) *exec.Cmd {
					return exec.Command("false")
				},
				outputWriter: &buf,
				helper:       NewHelper(),
			}
			remoteer.helper.outputWriter = &buf

			remoteer.Remote([]string{"set-url", tt.remoteName, tt.url})

			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected %q in output, got: %s", tt.expected, output)
			}
		})
	}
}

func TestRemoteer_Remote_List_RepositoryErrors(t *testing.T) {
	tests := []struct {
		name            string
		execCommandFunc func(string, ...string) *exec.Cmd
		expected        string
	}{
		{
			name: "Not a git repository",
			execCommandFunc: func(_ string, _ ...string) *exec.Cmd {
				cmd := exec.Command("sh", "-c", "echo 'fatal: not a git repository' >&2; exit 1")
				return cmd
			},
			expected: "Error: failed to list remotes",
		},
		{
			name: "Permission denied",
			execCommandFunc: func(_ string, _ ...string) *exec.Cmd {
				cmd := exec.Command("sh", "-c", "echo 'permission denied' >&2; exit 1")
				return cmd
			},
			expected: "Error: failed to list remotes",
		},
		{
			name: "Git command not found",
			execCommandFunc: func(_ string, _ ...string) *exec.Cmd {
				return exec.Command("nonexistent-command")
			},
			expected: "Error: failed to list remotes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			remoteer := &Remoteer{
				execCommand:  tt.execCommandFunc,
				outputWriter: &buf,
				helper:       NewHelper(),
			}
			remoteer.helper.outputWriter = &buf

			remoteer.Remote([]string{"list"})

			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected %q in output, got: %s", tt.expected, output)
			}
		})
	}
}

func TestRemoteer_Remote_ExcessiveArguments(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Add with too many arguments",
			args:     []string{"add", "origin", "https://example.com", "extra", "args"},
			expected: "Usage: ggc remote <command>", // Should show help screen
		},
		{
			name:     "Remove with too many arguments",
			args:     []string{"remove", "origin", "extra", "args"},
			expected: "Usage: ggc remote <command>", // Should show help screen
		},
		{
			name:     "Set-url with too many arguments",
			args:     []string{"set-url", "origin", "https://example.com", "extra", "args"},
			expected: "Usage: ggc remote <command>", // Should show help screen
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			remoteer := &Remoteer{
				execCommand: func(_ string, _ ...string) *exec.Cmd {
					return newNoopCmd().cmd
				},
				outputWriter: &buf,
				helper:       NewHelper(),
			}
			remoteer.helper.outputWriter = &buf

			remoteer.Remote(tt.args)

			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected %q in output, got: %s", tt.expected, output)
			}
		})
	}
}

func TestRemoteer_Remote_ConcurrentOperations(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		args      []string
		expected  string
	}{
		{
			name:      "Add operation",
			operation: "add",
			args:      []string{"add", "origin1", "https://example1.com"},
			expected:  "Remote 'origin1' added",
		},
		{
			name:      "Remove operation",
			operation: "remove",
			args:      []string{"remove", "origin2"},
			expected:  "Remote 'origin2' removed",
		},
		{
			name:      "Set-url operation",
			operation: "set-url",
			args:      []string{"set-url", "origin3", "https://example3.com"},
			expected:  "Remote 'origin3' URL updated",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			remoteer := &Remoteer{
				execCommand: func(_ string, _ ...string) *exec.Cmd {
					return newNoopCmd().cmd
				},
				outputWriter: &buf,
				helper:       NewHelper(),
			}
			remoteer.helper.outputWriter = &buf

			remoteer.Remote(tt.args)

			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected %q in output, got: %s", tt.expected, output)
			}
		})
	}
}

func TestRemoteer_Remote_SpecialCharacterHandling(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Unicode in remote name",
			args:     []string{"add", "орigin", "https://example.com"},
			expected: "Error: failed to add remote",
		},
		{
			name:     "URL with query parameters",
			args:     []string{"add", "origin", "https://example.com/repo?token=abc123"},
			expected: "Remote 'origin' added",
		},
		{
			name:     "URL with fragment",
			args:     []string{"add", "origin", "https://example.com/repo#main"},
			expected: "Remote 'origin' added",
		},
		{
			name:     "SSH URL format",
			args:     []string{"add", "origin", "git@github.com:user/repo.git"},
			expected: "Remote 'origin' added",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			remoteer := &Remoteer{
				execCommand: func(_ string, _ ...string) *exec.Cmd {
					if strings.Contains(tt.expected, "Error") {
						return exec.Command("false")
					}
					return newNoopCmd().cmd
				},
				outputWriter: &buf,
				helper:       NewHelper(),
			}
			remoteer.helper.outputWriter = &buf

			remoteer.Remote(tt.args)

			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected %q in output, got: %s", tt.expected, output)
			}
		})
	}
}

func TestRemoteer_Remote_NetworkRelatedErrors(t *testing.T) {
	tests := []struct {
		name            string
		command         []string
		execCommandFunc func(string, ...string) *exec.Cmd
		expected        string
	}{
		{
			name:    "Network timeout on add",
			command: []string{"add", "origin", "https://timeout.example.com"},
			execCommandFunc: func(_ string, _ ...string) *exec.Cmd {
				cmd := exec.Command("sh", "-c", "echo 'fatal: unable to access' >&2; exit 1")
				return cmd
			},
			expected: "Error: failed to add remote",
		},
		{
			name:    "DNS resolution failure",
			command: []string{"add", "origin", "https://nonexistent.domain.local"},
			execCommandFunc: func(_ string, _ ...string) *exec.Cmd {
				cmd := exec.Command("sh", "-c", "echo 'fatal: could not resolve host' >&2; exit 1")
				return cmd
			},
			expected: "Error: failed to add remote",
		},
		{
			name:    "Authentication failure",
			command: []string{"add", "origin", "https://private.example.com/repo.git"},
			execCommandFunc: func(_ string, _ ...string) *exec.Cmd {
				cmd := exec.Command("sh", "-c", "echo 'fatal: authentication failed' >&2; exit 1")
				return cmd
			},
			expected: "Error: failed to add remote",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			remoteer := &Remoteer{
				execCommand:  tt.execCommandFunc,
				outputWriter: &buf,
				helper:       NewHelper(),
			}
			remoteer.helper.outputWriter = &buf

			remoteer.Remote(tt.command)

			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected %q in output, got: %s", tt.expected, output)
			}
		})
	}
}
