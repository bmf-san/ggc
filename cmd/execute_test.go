package cmd

import (
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v7/internal/config"
	"github.com/bmf-san/ggc/v7/internal/testutil"
)

func TestExecute_BasicCommands(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{name: "help", args: []string{"help"}},
		{name: "version", args: []string{"version"}},
		{name: "status", args: []string{"status"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := testutil.NewMockGitClient()
			cm := config.NewConfigManager(mockClient)
			_ = cm.LoadConfig()
			c := NewCmd(mockClient, cm)
			err := c.Execute(tc.args)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestExecute_WithSimpleAlias(t *testing.T) {
	mockClient := testutil.NewMockGitClient()
	configManager := config.NewConfigManager(mockClient)
	_ = configManager.LoadConfig()

	cfg := configManager.GetConfig()
	cfg.Aliases = map[string]interface{}{
		"st": "status",
	}

	c := NewCmd(mockClient, configManager)
	err := c.Execute([]string{"st"})

	if err != nil {
		t.Errorf("simple alias should not return error: %v", err)
	}
}

func TestExecute_WithSimpleAliasAndArgs(t *testing.T) {
	mockClient := testutil.NewMockGitClient()
	configManager := config.NewConfigManager(mockClient)
	_ = configManager.LoadConfig()

	cfg := configManager.GetConfig()
	cfg.Aliases = map[string]interface{}{
		"br": "branch",
	}

	c := NewCmd(mockClient, configManager)
	err := c.Execute([]string{"br", "current"})

	if err != nil {
		t.Errorf("simple alias with args should not return error: %v", err)
	}
}

func TestExecute_WithSequenceAlias(t *testing.T) {
	mockClient := testutil.NewMockGitClient()
	configManager := config.NewConfigManager(mockClient)
	_ = configManager.LoadConfig()

	cfg := configManager.GetConfig()
	cfg.Aliases = map[string]interface{}{
		"sync": []interface{}{"status", "log simple"},
	}

	c := NewCmd(mockClient, configManager)
	err := c.Execute([]string{"sync"})

	if err != nil {
		t.Errorf("sequence alias should not return error: %v", err)
	}
}

func TestExecute_SequenceAliasRejectsArguments(t *testing.T) {
	mockClient := testutil.NewMockGitClient()
	configManager := config.NewConfigManager(mockClient)
	_ = configManager.LoadConfig()

	cfg := configManager.GetConfig()
	cfg.Aliases = map[string]interface{}{
		"deploy": []interface{}{"status"},
	}

	c := NewCmd(mockClient, configManager)
	err := c.Execute([]string{"deploy", "production"})

	if err == nil {
		t.Fatal("sequence alias should return error when arguments are provided")
	}
	if !strings.Contains(err.Error(), "sequence alias 'deploy' does not accept arguments") {
		t.Fatalf("expected rejection message, got %q", err.Error())
	}
	if !strings.Contains(err.Error(), "production") {
		t.Fatalf("expected error to list offending arguments, got %q", err.Error())
	}
}

func TestExecute_ConfigManagerNil(t *testing.T) {
	mockClient := testutil.NewMockGitClient()
	c := NewCmd(mockClient, nil)

	// Should not panic and should execute normal command
	err := c.Execute([]string{"help"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestExecute_NonAliasCommand(t *testing.T) {
	mockClient := testutil.NewMockGitClient()
	configManager := config.NewConfigManager(mockClient)
	_ = configManager.LoadConfig()

	cfg := configManager.GetConfig()
	cfg.Aliases = map[string]interface{}{
		"st": "status",
	}

	c := NewCmd(mockClient, configManager)
	// "commit" is not an alias, should be routed directly
	err := c.Execute([]string{"commit", "test"})

	if err != nil {
		t.Errorf("non-alias command should not return error: %v", err)
	}
}

func TestExecute_PlaceholderProcessing(t *testing.T) {
	tests := []struct {
		name        string
		aliases     map[string]interface{}
		args        []string
		expectError bool
		errorSubstr string
	}{
		{
			name: "simple alias with placeholder",
			aliases: map[string]interface{}{
				"commit-msg": "commit -m '{0}'",
			},
			args:        []string{"commit-msg", "fix bug"},
			expectError: false,
		},
		{
			name: "sequence alias with placeholders",
			aliases: map[string]interface{}{
				"deploy": []interface{}{"branch checkout {0}", "log simple"},
			},
			args:        []string{"deploy", "production"},
			expectError: false,
		},
		{
			name: "simple alias with placeholder - missing argument",
			aliases: map[string]interface{}{
				"commit-msg": "commit -m '{0}'",
			},
			args:        []string{"commit-msg"},
			expectError: true,
			errorSubstr: "requires at least 1 argument",
		},
		{
			name: "simple alias without placeholders - arguments forwarded",
			aliases: map[string]interface{}{
				"st": "status",
			},
			args:        []string{"st", "short"},
			expectError: false,
		},
		{
			name: "sequence alias without placeholders - arguments rejected",
			aliases: map[string]interface{}{
				"sync": []interface{}{"status", "log simple"},
			},
			args:        []string{"sync", "unwanted"},
			expectError: true,
			errorSubstr: "does not accept arguments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := testutil.NewMockGitClient()
			configManager := config.NewConfigManager(mockClient)

			cfg := configManager.GetConfig()
			cfg.Aliases = tt.aliases

			c := NewCmd(mockClient, configManager)
			err := c.Execute(tt.args)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if tt.errorSubstr != "" && !strings.Contains(err.Error(), tt.errorSubstr) {
					t.Errorf("expected error containing %q, got %q", tt.errorSubstr, err.Error())
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestExecute_PlaceholderEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		aliases     map[string]interface{}
		args        []string
		expectError bool
		errorSubstr string
	}{
		{
			name: "duplicate placeholders in same command",
			aliases: map[string]interface{}{
				"duplicate": []interface{}{"branch checkout {0}", "commit -m '{0} - {0}'"},
			},
			args:        []string{"duplicate", "main"},
			expectError: false,
		},
		{
			name: "excess arguments beyond placeholders",
			aliases: map[string]interface{}{
				"single-placeholder": "branch checkout {0}",
			},
			args:        []string{"single-placeholder", "main", "extra1", "extra2"},
			expectError: false,
		},
		{
			name: "multiple placeholders in mixed order",
			aliases: map[string]interface{}{
				"mixed-order": "branch checkout-from {1} {0}",
			},
			args:        []string{"mixed-order", "feature/test", "main"},
			expectError: false,
		},
		{
			name: "sequence alias with no placeholders gets extra args rejected",
			aliases: map[string]interface{}{
				"no-placeholders-seq": []interface{}{"status", "branch current"},
			},
			args:        []string{"no-placeholders-seq", "unwanted"},
			expectError: true,
			errorSubstr: "does not accept arguments",
		},
		{
			name: "simple alias with placeholders but no args provided",
			aliases: map[string]interface{}{
				"needs-arg": "commit -m '{0}'",
			},
			args:        []string{"needs-arg"},
			expectError: true,
			errorSubstr: "requires at least 1 argument",
		},
		{
			name: "sequence alias with multiple placeholders - insufficient arguments",
			aliases: map[string]interface{}{
				"feature": []interface{}{"branch checkout-from {0} feature/{1}", "commit -m 'Start {1} from {0}'"},
			},
			args:        []string{"feature", "main"},
			expectError: true,
			errorSubstr: "requires at least 2 argument",
		},
		{
			name: "sequence alias with multiple placeholders - sufficient arguments",
			aliases: map[string]interface{}{
				"feature": []interface{}{"branch checkout-from {0} feature/{1}", "log simple"},
			},
			args:        []string{"feature", "main", "user-auth"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := testutil.NewMockGitClient()
			configManager := config.NewConfigManager(mockClient)

			cfg := configManager.GetConfig()
			cfg.Aliases = tt.aliases

			c := NewCmd(mockClient, configManager)
			err := c.Execute(tt.args)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if tt.errorSubstr != "" && !strings.Contains(err.Error(), tt.errorSubstr) {
					t.Errorf("expected error containing %q, got %q", tt.errorSubstr, err.Error())
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestExecute_InvalidAliasFormat(t *testing.T) {
	mockClient := testutil.NewMockGitClient()
	configManager := config.NewConfigManager(mockClient)
	_ = configManager.LoadConfig()

	cfg := configManager.GetConfig()
	cfg.Aliases = map[string]interface{}{
		"invalid": 123, // Invalid format - should be string or []interface{}
	}

	c := NewCmd(mockClient, configManager)
	err := c.Execute([]string{"invalid"})

	if err == nil {
		t.Fatal("invalid alias format should return error")
	}
}
