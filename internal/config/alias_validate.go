package config

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bmf-san/ggc/v7/cmd/command"
)

// shellMetacharacters are characters that could enable command injection attacks.
// Blocked characters: ;|&<>(){}[]$`\n\r
// These characters are commonly used in shell command injection:
//   - ; | && || - Command chaining
//   - < > >> - I/O redirection
//   - () - Subshells
//   - {} [] - Brace/bracket expansion
//   - $ - Variable expansion and command substitution
//   - ` - Command substitution (backticks)
//   - \n \r - Newline injection
const shellMetacharacters = ";|&<>(){}[]$`\n\r"

// commandValidator validates alias commands against security policies.
// It maintains a whitelist of valid commands and checks for shell metacharacters.
type commandValidator struct {
	validCommands map[string]struct{}
	once          sync.Once
}

// newCommandValidator creates a new command validator.
func newCommandValidator() *commandValidator {
	return &commandValidator{}
}

// initCommands lazily initializes the valid command set.
func (v *commandValidator) initCommands() {
	v.once.Do(func() {
		v.validCommands = make(map[string]struct{})
		registry := command.NewRegistry()
		allCommands := registry.All()
		for i := range allCommands {
			v.validCommands[allCommands[i].Name] = struct{}{}
		}
	})
}

// validateCommand validates a single alias command string for security.
// It checks for both shell metacharacters and valid command names.
//
// Returns an error if:
//   - The command contains shell metacharacters
//   - The command name is not a valid ggc command
func (v *commandValidator) validateCommand(cmd string) error {
	// Check for shell metacharacters first (fast path)
	if strings.ContainsAny(cmd, shellMetacharacters) {
		return fmt.Errorf("command '%s' contains unsafe shell metacharacters (;|&<>(){}[]$`)", cmd)
	}

	// Ensure command set is initialized
	v.initCommands()

	// Extract command name (first word)
	cmdParts := strings.SplitN(cmd, " ", 2)
	cmdName := cmdParts[0]

	// Check if command is in the whitelist
	if _, valid := v.validCommands[cmdName]; !valid {
		return fmt.Errorf("'%s' is not a valid ggc command", cmdName)
	}

	return nil
}

// isValidCommand checks if a command name is in the valid command registry.
func (v *commandValidator) isValidCommand(cmdName string) bool {
	v.initCommands()
	_, valid := v.validCommands[cmdName]
	return valid
}

// defaultValidator is the package-level validator instance.
var defaultValidator = newCommandValidator()

// validateAliases validates all aliases in the configuration.
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

// validateAliasName validates an alias name format.
func validateAliasName(name string) error {
	if strings.TrimSpace(name) == "" || strings.Contains(name, " ") {
		return &ValidationError{"aliases." + name, name, "alias names must not contain spaces"}
	}
	return nil
}

// validateAliasValue validates an alias value (string or sequence).
func validateAliasValue(name string, value interface{}) error {
	switch v := value.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			return &ValidationError{"aliases." + name, v, "alias command cannot be empty"}
		}
		// Validate command security for simple string aliases
		if err := defaultValidator.validateCommand(v); err != nil {
			return &ValidationError{
				Field:   "aliases." + name,
				Value:   v,
				Message: err.Error(),
			}
		}
		return nil
	case []interface{}:
		return validateAliasSequence(name, v)
	default:
		return &ValidationError{Field: "aliases." + name, Value: value, Message: "alias must be either a string or array of strings"}
	}
}

// validateAliasSequence validates an alias sequence (array of commands).
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
		if err := defaultValidator.validateCommand(cmdStr); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("aliases.%s[%d]", name, i),
				Value:   cmdStr,
				Message: err.Error(),
			}
		}
	}
	return nil
}
