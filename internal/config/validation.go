package config

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bmf-san/ggc/v7/internal/aliasvalidator"
)

func (c *Config) validateBranch() error {
	branch := c.Default.Branch
	if strings.TrimSpace(branch) == "" || strings.Contains(branch, " ") {
		return &ValidationError{"default.branch", branch, "must not contain spaces or be empty"}
	}
	return nil
}

func (c *Config) validateEditor() error {
	editor := strings.TrimSpace(c.Default.Editor)
	bin := parseEditorBinary(editor)
	if validEditorPath(bin) {
		return nil
	}
	if _, err := exec.LookPath(bin); err != nil {
		return &ValidationError{"default.editor", editor, "command not found in PATH or invalid path"}
	}
	return nil
}

func parseEditorBinary(editor string) string {
	if editor == "" {
		return ""
	}
	// Support basic quoted paths or first token before whitespace
	if (strings.HasPrefix(editor, "\"") && strings.Count(editor, "\"") >= 2) || (strings.HasPrefix(editor, "'") && strings.Count(editor, "'") >= 2) {
		q := editor[0:1]
		if idx := strings.Index(editor[1:], q); idx >= 0 {
			return editor[1 : 1+idx]
		}
	}
	if i := strings.IndexAny(editor, " \t"); i > 0 {
		return editor[:i]
	}
	return editor
}

func validEditorPath(bin string) bool {
	if !strings.ContainsAny(bin, "/\\") {
		return false
	}
	_, err := os.Stat(bin)
	return err == nil
}

func (c *Config) validateConfirmDestructive() error {
	val := c.Behavior.ConfirmDestructive
	valid := map[string]bool{"simple": true, "always": true, "never": true}
	if !valid[val] {
		return &ValidationError{"behavior.confirm-destructive", val, "must be one of: simple, always, never"}
	}
	return nil
}

// validateGitDefaultRemote validates git default remote name format
func (c *Config) validateGitDefaultRemote() error {
	remote := c.Git.DefaultRemote
	if remote == "" {
		return nil
	}

	if !gitRemoteNameCharsRe.MatchString(remote) || strings.Contains(remote, " ") {
		return &ValidationError{
			Field:   "git.default-remote",
			Value:   remote,
			Message: "Remote may contain letters, digits, ., _, -, and / only",
		}
	}

	// Additional structural checks: no leading/trailing '.' or '/', and no empty/unsafe segments
	if strings.HasPrefix(remote, "/") || strings.HasSuffix(remote, "/") || strings.HasPrefix(remote, ".") || strings.HasSuffix(remote, ".") || strings.Contains(remote, "//") || strings.Contains(remote, "..") {
		return &ValidationError{
			Field:   "git.default-remote",
			Value:   remote,
			Message: "Remote must not start/end with '.' or '/', nor contain '..' or '//'",
		}
	}

	return nil
}

func (c *Config) validateAliases() error {
	for name, value := range c.Aliases {
		if err := validateAliasName(name); err != nil {
			return err
		}
		if err := validateAliasValue(name, value); err != nil {
			return err
		}
	}
	return nil
}

func validateAliasName(name string) error {
	if strings.TrimSpace(name) == "" || strings.Contains(name, " ") {
		return &ValidationError{"aliases." + name, name, "alias names must not contain spaces"}
	}
	return nil
}

func validateAliasValue(name string, value interface{}) error {
	switch v := value.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			return &ValidationError{"aliases." + name, v, "alias command cannot be empty"}
		}
		return nil
	case []interface{}:
		return validateAliasSequence(name, v)
	default:
		return &ValidationError{Field: "aliases." + name, Value: value, Message: "alias must be either a string or array of strings"}
	}
}

func validateAliasSequence(name string, seq []interface{}) error {
	if len(seq) == 0 {
		return &ValidationError{"aliases." + name, seq, "alias sequence cannot be empty"}
	}

	for i, cmd := range seq {
		cmdStr, ok := cmd.(string)
		if !ok {
			return &ValidationError{Field: fmt.Sprintf("aliases.%s[%d]", name, i), Value: cmd, Message: "sequence commands must be strings"}
		}
		if strings.TrimSpace(cmdStr) == "" {
			return &ValidationError{Field: fmt.Sprintf("aliases.%s[%d]", name, i), Value: cmdStr, Message: "command in sequence cannot be empty"}
		}

		// Validate command security
		if err := aliasvalidator.ValidateCommand(cmdStr); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("aliases.%s[%d]", name, i),
				Value:   cmdStr,
				Message: err.Error(),
			}
		}
	}
	return nil
}

// Validate is a function that handles validation operations
func (c *Config) Validate() error {
	if err := c.validateBranch(); err != nil {
		return err
	}
	if err := c.validateEditor(); err != nil {
		return err
	}
	if err := c.validateConfirmDestructive(); err != nil {
		return err
	}
	if err := c.validateGitDefaultRemote(); err != nil {
		return err
	}
	if err := c.validateAliases(); err != nil {
		return err
	}
	if err := c.validateKeybindings(); err != nil {
		return err
	}
	return nil
}
