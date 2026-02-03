package aliasvalidator

import (
	"strings"
	"testing"
)

func TestValidator_ValidateCommand(t *testing.T) {
	v := NewValidator()

	cases := []struct {
		name      string
		cmd       string
		wantError bool
		errorMsg  string
	}{
		// Valid commands
		{name: "simple status", cmd: "status", wantError: false},
		{name: "branch with args", cmd: "branch current", wantError: false},
		{name: "commit with args", cmd: "commit allow empty", wantError: false},
		{name: "diff with args", cmd: "diff staged", wantError: false},
		{name: "command with hyphen", cmd: "branch delete-merged", wantError: false},
		{name: "command with dot", cmd: "restore .", wantError: false},

		// Shell metacharacter injection
		{name: "semicolon injection", cmd: "status; echo pwned", wantError: true, errorMsg: "unsafe shell metacharacters"},
		{name: "pipe injection", cmd: "status | cat", wantError: true, errorMsg: "unsafe shell metacharacters"},
		{name: "ampersand injection", cmd: "status && echo pwned", wantError: true, errorMsg: "unsafe shell metacharacters"},
		{name: "command substitution", cmd: "status $(whoami)", wantError: true, errorMsg: "unsafe shell metacharacters"},
		{name: "backtick injection", cmd: "`whoami`", wantError: true, errorMsg: "unsafe shell metacharacters"},
		{name: "redirection", cmd: "status > /tmp/output", wantError: true, errorMsg: "unsafe shell metacharacters"},
		{name: "stash ref with braces", cmd: "stash show stash@{0}", wantError: true, errorMsg: "unsafe shell metacharacters"},
		{name: "newline injection", cmd: "status\necho pwned", wantError: true, errorMsg: "unsafe shell metacharacters"},

		// Invalid commands
		{name: "invalid command", cmd: "notacommand", wantError: true, errorMsg: "not a valid ggc command"},
		{name: "echo command", cmd: "echo test", wantError: true, errorMsg: "not a valid ggc command"},
		{name: "cat command", cmd: "cat file.txt", wantError: true, errorMsg: "not a valid ggc command"},

		// Edge cases
		{name: "valid command with invalid args metachar", cmd: "branch name; echo", wantError: true, errorMsg: "unsafe shell metacharacters"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := v.ValidateCommand(tc.cmd)
			if tc.wantError {
				if err == nil {
					t.Errorf("ValidateCommand(%q) expected error, got nil", tc.cmd)
					return
				}
				if tc.errorMsg != "" && !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("ValidateCommand(%q) error = %v, want to contain %q", tc.cmd, err.Error(), tc.errorMsg)
				}
			} else if err != nil {
				t.Errorf("ValidateCommand(%q) unexpected error: %v", tc.cmd, err)
			}
		})
	}
}

func TestValidator_IsValidCommand(t *testing.T) {
	v := NewValidator()

	cases := []struct {
		name     string
		cmdName  string
		expected bool
	}{
		// Valid commands
		{name: "status", cmdName: "status", expected: true},
		{name: "branch", cmdName: "branch", expected: true},
		{name: "commit", cmdName: "commit", expected: true},
		{name: "push", cmdName: "push", expected: true},
		{name: "pull", cmdName: "pull", expected: true},
		{name: "log", cmdName: "log", expected: true},
		{name: "diff", cmdName: "diff", expected: true},
		{name: "add", cmdName: "add", expected: true},
		{name: "tag", cmdName: "tag", expected: true},
		{name: "stash", cmdName: "stash", expected: true},
		{name: "help", cmdName: "help", expected: true},

		// Invalid commands
		{name: "invalid command", cmdName: "notacommand", expected: false},
		{name: "echo", cmdName: "echo", expected: false},
		{name: "cat", cmdName: "cat", expected: false},
		{name: "ls", cmdName: "ls", expected: false},
		{name: "rm", cmdName: "rm", expected: false},
		{name: "empty", cmdName: "", expected: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := v.IsValidCommand(tc.cmdName)
			if result != tc.expected {
				t.Errorf("IsValidCommand(%q) = %v, want %v", tc.cmdName, result, tc.expected)
			}
		})
	}
}

func TestPackageFunctions(t *testing.T) {
	// Test package-level convenience functions
	if err := ValidateCommand("status"); err != nil {
		t.Errorf("ValidateCommand(\"status\") unexpected error: %v", err)
	}

	if err := ValidateCommand("status; echo pwned"); err == nil {
		t.Error("ValidateCommand with injection should error")
	}

	if !IsValidCommand("status") {
		t.Error("IsValidCommand(\"status\") should be true")
	}

	if IsValidCommand("notacommand") {
		t.Error("IsValidCommand(\"notacommand\") should be false")
	}
}

func TestValidator_LazyInit(t *testing.T) {
	// Test that multiple validators work independently
	v1 := NewValidator()
	v2 := NewValidator()

	if !v1.IsValidCommand("status") {
		t.Error("v1 should validate status")
	}

	if !v2.IsValidCommand("branch") {
		t.Error("v2 should validate branch")
	}
}
