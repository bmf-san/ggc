// Package aliasvalidator provides alias command validation utilities for ggc.
package aliasvalidator

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

// Validator validates alias commands against security policies.
// It maintains a whitelist of valid commands and checks for shell metacharacters.
type Validator struct {
	validCommands map[string]struct{}
	once          sync.Once
}

// NewValidator creates a new command validator.
func NewValidator() *Validator {
	return &Validator{}
}

// initCommands lazily initializes the valid command set.
func (v *Validator) initCommands() {
	v.once.Do(func() {
		v.validCommands = make(map[string]struct{})
		allCommands := command.All()
		for i := range allCommands {
			v.validCommands[allCommands[i].Name] = struct{}{}
		}
	})
}

// ValidateCommand validates a single alias command string for security.
// It checks for both shell metacharacters and valid command names.
//
// Returns an error if:
//   - The command contains shell metacharacters
//   - The command name is not a valid ggc command
func (v *Validator) ValidateCommand(cmd string) error {
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

// IsValidCommand checks if a command name is in the valid command registry.
func (v *Validator) IsValidCommand(cmdName string) bool {
	v.initCommands()
	_, valid := v.validCommands[cmdName]
	return valid
}

// defaultValidator is the package-level validator instance.
var defaultValidator = NewValidator()

// ValidateCommand validates a command using the default validator.
// This is a convenience function for one-off validations.
func ValidateCommand(cmd string) error {
	return defaultValidator.ValidateCommand(cmd)
}

// IsValidCommand checks if a command is valid using the default validator.
// This is a convenience function for one-off checks.
func IsValidCommand(cmdName string) bool {
	return defaultValidator.IsValidCommand(cmdName)
}
