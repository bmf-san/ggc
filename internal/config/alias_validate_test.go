package config

import (
	"strings"
	"testing"
)

func TestCommandValidator_ValidateCommand(t *testing.T) {
	v := newCommandValidator()

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
			err := v.validateCommand(tc.cmd)
			if tc.wantError {
				if err == nil {
					t.Errorf("validateCommand(%q) expected error, got nil", tc.cmd)
					return
				}
				if tc.errorMsg != "" && !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("validateCommand(%q) error = %v, want to contain %q", tc.cmd, err.Error(), tc.errorMsg)
				}
			} else if err != nil {
				t.Errorf("validateCommand(%q) unexpected error: %v", tc.cmd, err)
			}
		})
	}
}

func TestCommandValidator_IsValidCommand(t *testing.T) {
	v := newCommandValidator()

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
			result := v.isValidCommand(tc.cmdName)
			if result != tc.expected {
				t.Errorf("isValidCommand(%q) = %v, want %v", tc.cmdName, result, tc.expected)
			}
		})
	}
}

func TestValidateAliasName(t *testing.T) {
	cases := []struct {
		name      string
		aliasName string
		wantError bool
	}{
		{name: "valid simple name", aliasName: "myalias", wantError: false},
		{name: "valid with hyphen", aliasName: "my-alias", wantError: false},
		{name: "valid with underscore", aliasName: "my_alias", wantError: false},
		{name: "empty name", aliasName: "", wantError: true},
		{name: "name with space", aliasName: "my alias", wantError: true},
		{name: "only spaces", aliasName: "   ", wantError: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateAliasName(tc.aliasName)
			if tc.wantError && err == nil {
				t.Errorf("validateAliasName(%q) expected error, got nil", tc.aliasName)
			} else if !tc.wantError && err != nil {
				t.Errorf("validateAliasName(%q) unexpected error: %v", tc.aliasName, err)
			}
		})
	}
}

func TestValidateAliasValue(t *testing.T) {
	cases := []struct {
		name      string
		aliasName string
		value     interface{}
		wantError bool
	}{
		{name: "valid string", aliasName: "test", value: "status", wantError: false},
		{name: "valid string with args", aliasName: "test", value: "branch current", wantError: false},
		{name: "empty string", aliasName: "test", value: "", wantError: true},
		{name: "whitespace only", aliasName: "test", value: "   ", wantError: true},
		{name: "valid sequence", aliasName: "test", value: []interface{}{"status", "branch"}, wantError: false},
		{name: "invalid type", aliasName: "test", value: 123, wantError: true},
		{name: "invalid type map", aliasName: "test", value: map[string]string{"key": "value"}, wantError: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateAliasValue(tc.aliasName, tc.value)
			if tc.wantError && err == nil {
				t.Errorf("validateAliasValue(%q, %v) expected error, got nil", tc.aliasName, tc.value)
			} else if !tc.wantError && err != nil {
				t.Errorf("validateAliasValue(%q, %v) unexpected error: %v", tc.aliasName, tc.value, err)
			}
		})
	}
}

func TestValidateAliasSequence(t *testing.T) {
	cases := []struct {
		name      string
		aliasName string
		sequence  []interface{}
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid sequence",
			aliasName: "test",
			sequence:  []interface{}{"status", "branch"},
			wantError: false,
		},
		{
			name:      "empty sequence",
			aliasName: "test",
			sequence:  []interface{}{},
			wantError: true,
			errorMsg:  "cannot be empty",
		},
		{
			name:      "invalid command type",
			aliasName: "test",
			sequence:  []interface{}{"status", 123},
			wantError: true,
			errorMsg:  "must be strings",
		},
		{
			name:      "empty command in sequence",
			aliasName: "test",
			sequence:  []interface{}{"status", ""},
			wantError: true,
			errorMsg:  "cannot be empty",
		},
		{
			name:      "shell injection in sequence",
			aliasName: "test",
			sequence:  []interface{}{"status", "branch; echo pwned"},
			wantError: true,
			errorMsg:  "unsafe shell metacharacters",
		},
		{
			name:      "invalid command in sequence",
			aliasName: "test",
			sequence:  []interface{}{"status", "notacommand"},
			wantError: true,
			errorMsg:  "not a valid ggc command",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateAliasSequence(tc.aliasName, tc.sequence)
			if tc.wantError {
				if err == nil {
					t.Errorf("validateAliasSequence(%q, %v) expected error, got nil", tc.aliasName, tc.sequence)
					return
				}
				if tc.errorMsg != "" && !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("validateAliasSequence(%q, %v) error = %v, want to contain %q", tc.aliasName, tc.sequence, err.Error(), tc.errorMsg)
				}
			} else if err != nil {
				t.Errorf("validateAliasSequence(%q, %v) unexpected error: %v", tc.aliasName, tc.sequence, err)
			}
		})
	}
}

func TestConfig_ValidateAliases(t *testing.T) {
	cases := []struct {
		name      string
		config    *Config
		wantError bool
	}{
		{
			name: "valid simple alias",
			config: &Config{
				Aliases: map[string]interface{}{
					"st": "status",
				},
			},
			wantError: false,
		},
		{
			name: "valid sequence alias",
			config: &Config{
				Aliases: map[string]interface{}{
					"workflow": []interface{}{"add .", "commit -m test"},
				},
			},
			wantError: false,
		},
		{
			name: "invalid alias name with space",
			config: &Config{
				Aliases: map[string]interface{}{
					"my alias": "status",
				},
			},
			wantError: true,
		},
		{
			name: "invalid command with injection",
			config: &Config{
				Aliases: map[string]interface{}{
					"bad": "status; echo pwned",
				},
			},
			wantError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.validateAliases()
			if tc.wantError && err == nil {
				t.Errorf("validateAliases() expected error, got nil")
			} else if !tc.wantError && err != nil {
				t.Errorf("validateAliases() unexpected error: %v", err)
			}
		})
	}
}

func TestCommandValidator_LazyInit(t *testing.T) {
	// Test that multiple validators work independently
	v1 := newCommandValidator()
	v2 := newCommandValidator()

	if !v1.isValidCommand("status") {
		t.Error("v1 should validate status")
	}

	if !v2.isValidCommand("branch") {
		t.Error("v2 should validate branch")
	}
}
