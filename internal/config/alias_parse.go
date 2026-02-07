package config

import (
	"fmt"
	"strings"
)

// analyzePlaceholders analyzes a command string for placeholders and returns placeholder info
func analyzePlaceholders(commands []string) (map[string]struct{}, int, error) {
	placeholders := make(map[string]struct{})
	maxPositionalArg := -1

	for _, cmd := range commands {
		matches := aliasPlaceholderPattern.FindAllStringSubmatch(cmd, -1)
		for _, match := range matches {
			if len(match) < 2 {
				continue
			}
			placeholder := match[1]

			// Validate placeholder format
			if err := validatePlaceholder(placeholder); err != nil {
				return nil, -1, fmt.Errorf("invalid placeholder {%s}: %w", placeholder, err)
			}

			placeholders[placeholder] = struct{}{}

			// Check if it's a positional argument
			if len(placeholder) == 1 && placeholder[0] >= '0' && placeholder[0] <= '9' {
				argIndex := int(placeholder[0] - '0')
				if argIndex > maxPositionalArg {
					maxPositionalArg = argIndex
				}
			}
		}
	}

	return placeholders, maxPositionalArg, nil
}

// isValidPlaceholderChar checks if a character is valid in a placeholder
func isValidPlaceholderChar(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9') ||
		char == '_' || char == '-'
}

// validatePlaceholder validates a placeholder name
func validatePlaceholder(placeholder string) error {
	if placeholder == "" {
		return fmt.Errorf("empty placeholder")
	}

	// Check for shell metacharacters.
	// Note: Braces '{}' are included here because they are used as placeholder delimiters.
	// This prevents nested placeholders like {message: {0}}, as braces in the placeholder content are rejected.
	if strings.ContainsAny(placeholder, ";|&$`()[]{}*?<>\"'\\") {
		return fmt.Errorf("placeholder contains unsafe characters")
	}

	// Allow alphanumeric, underscore, and hyphen
	for _, char := range placeholder {
		if !isValidPlaceholderChar(char) {
			return fmt.Errorf("placeholder contains invalid character: %c", char)
		}
	}

	return nil
}

// ParseAlias parses an alias value and returns its type and commands
func (c *Config) ParseAlias(name string) (*ParsedAlias, error) {
	value, exists := c.Aliases[name]
	if !exists {
		return nil, fmt.Errorf("alias '%s' not found", name)
	}

	switch v := value.(type) {
	case string:
		placeholders, maxPositionalArg, err := analyzePlaceholders([]string{v})
		if err != nil {
			return nil, fmt.Errorf("error analyzing placeholders in simple alias '%s': %w", name, err)
		}

		return &ParsedAlias{
			Type:             SimpleAlias,
			Commands:         []string{v},
			Placeholders:     placeholders,
			MaxPositionalArg: maxPositionalArg,
		}, nil

	case []interface{}:
		commands := make([]string, len(v))
		for i, cmd := range v {
			cmdStr, ok := cmd.(string)
			if !ok {
				return nil, fmt.Errorf("invalid command type in alias '%s'", name)
			}
			commands[i] = cmdStr
		}

		placeholders, maxPositionalArg, err := analyzePlaceholders(commands)
		if err != nil {
			return nil, fmt.Errorf("error analyzing placeholders in sequence alias '%s': %w", name, err)
		}

		return &ParsedAlias{
			Type:             SequenceAlias,
			Commands:         commands,
			Placeholders:     placeholders,
			MaxPositionalArg: maxPositionalArg,
		}, nil

	default:
		return nil, fmt.Errorf("invalid alias type for '%s'", name)
	}
}

// IsAlias checks if a given name is an alias
func (c *Config) IsAlias(name string) bool {
	_, exists := c.Aliases[name]
	return exists
}

// GetAliasCommands returns the commands for a given alias
func (c *Config) GetAliasCommands(name string) ([]string, error) {
	parsed, err := c.ParseAlias(name)
	if err != nil {
		return nil, err
	}
	return parsed.Commands, nil
}

// GetAllAliases returns all aliases with their parsed commands
func (c *Config) GetAllAliases() map[string]*ParsedAlias {
	result := make(map[string]*ParsedAlias)
	for name := range c.Aliases {
		if parsed, err := c.ParseAlias(name); err == nil {
			result[name] = parsed
		}
	}
	return result
}
