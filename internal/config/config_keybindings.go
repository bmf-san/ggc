package config

import (
	"fmt"
	"strings"
)

// validateKeybindings validates the keybinding configuration
func (c *Config) validateKeybindings() error {
	// Validate profile selection
	if err := c.validateProfile(); err != nil {
		return err
	}

	// Validate global keybindings
	bindings := map[string]string{
		"delete_word":          c.Interactive.Keybindings.DeleteWord,
		"clear_line":           c.Interactive.Keybindings.ClearLine,
		"delete_to_end":        c.Interactive.Keybindings.DeleteToEnd,
		"move_to_beginning":    c.Interactive.Keybindings.MoveToBeginning,
		"move_to_end":          c.Interactive.Keybindings.MoveToEnd,
		"move_up":              c.Interactive.Keybindings.MoveUp,
		"move_down":            c.Interactive.Keybindings.MoveDown,
		"move_left":            c.Interactive.Keybindings.MoveLeft,
		"move_right":           c.Interactive.Keybindings.MoveRight,
		"add_to_workflow":      c.Interactive.Keybindings.AddToWorkflow,
		"toggle_workflow_view": c.Interactive.Keybindings.ToggleWorkflowView,
		"clear_workflow":       c.Interactive.Keybindings.ClearWorkflow,
		"workflow_create":      c.Interactive.Keybindings.WorkflowCreate,
		"workflow_delete":      c.Interactive.Keybindings.WorkflowDelete,
		"soft_cancel":          c.Interactive.Keybindings.SoftCancel,
	}

	for action, keyStr := range bindings {
		// Empty bindings are allowed (will use defaults)
		if keyStr == "" {
			continue
		}
		if err := parseKeyBinding(keyStr); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("interactive.keybindings.%s", action),
				Value:   keyStr,
				Message: err.Error(),
			}
		}
	}

	// Validate context-specific keybindings
	if err := c.validateContextKeybindings(); err != nil {
		return err
	}

	// Validate platform-specific keybindings
	if err := c.validatePlatformKeybindings(); err != nil {
		return err
	}

	return nil
}

// validateProfile validates the profile selection
func (c *Config) validateProfile() error {
	profile := c.Interactive.Profile
	if profile == "" {
		return nil // Empty profile is allowed (defaults to "default")
	}

	validProfiles := map[string]bool{
		"default":  true,
		"emacs":    true,
		"vi":       true,
		"readline": true,
	}

	if !validProfiles[profile] {
		return &ValidationError{
			Field:   "interactive.profile",
			Value:   profile,
			Message: "must be one of: default, emacs, vi, readline",
		}
	}
	return nil
}

// validateContextKeybindings validates context-specific keybindings
func (c *Config) validateContextKeybindings() error {
	contexts := map[string]map[string]interface{}{
		"input":   c.Interactive.Contexts.Input.Keybindings,
		"results": c.Interactive.Contexts.Results.Keybindings,
		"search":  c.Interactive.Contexts.Search.Keybindings,
	}

	nonNil := 0
	for _, bindings := range contexts {
		if bindings != nil {
			nonNil++
		}
	}

	for contextName, bindings := range contexts {
		if bindings == nil {
			if nonNil > 0 {
				return &ValidationError{
					Field:   fmt.Sprintf("interactive.contexts.%s.keybindings", contextName),
					Value:   bindings,
					Message: "keybindings map is missing for this context",
				}
			}
			continue
		}
		for action, value := range bindings {
			if err := validateKeybindingValue(fmt.Sprintf("interactive.contexts.%s.keybindings.%s", contextName, action), value); err != nil {
				return err
			}
		}
	}
	return nil
}

// validatePlatformKeybindings validates platform and terminal specific keybindings
func (c *Config) validatePlatformKeybindings() error {
	platforms := map[string]map[string]interface{}{
		"darwin":  c.Interactive.Darwin.Keybindings,
		"linux":   c.Interactive.Linux.Keybindings,
		"windows": c.Interactive.Windows.Keybindings,
	}

	for platformName, bindings := range platforms {
		if bindings == nil {
			continue
		}
		for action, value := range bindings {
			if err := validateKeybindingValue(fmt.Sprintf("interactive.%s.keybindings.%s", platformName, action), value); err != nil {
				return err
			}
		}
	}

	// Validate terminal-specific keybindings
	if c.Interactive.Terminals != nil {
		for termName, termConfig := range c.Interactive.Terminals {
			if termConfig.Keybindings == nil {
				continue
			}
			for action, value := range termConfig.Keybindings {
				if err := validateKeybindingValue(fmt.Sprintf("interactive.terminals.%s.keybindings.%s", termName, action), value); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// validateKeybindingValue validates a keybinding value (string or array of strings)
func validateKeybindingValue(fieldPath string, value interface{}) error {
	switch v := value.(type) {
	case string:
		if v == "" {
			return nil // Empty is allowed
		}
		if err := parseKeyBinding(v); err != nil {
			return &ValidationError{
				Field:   fieldPath,
				Value:   v,
				Message: err.Error(),
			}
		}
	case []interface{}:
		for i, item := range v {
			itemStr, ok := item.(string)
			if !ok {
				return &ValidationError{
					Field:   fmt.Sprintf("%s[%d]", fieldPath, i),
					Value:   item,
					Message: "keybinding array items must be strings",
				}
			}
			if itemStr != "" {
				if err := parseKeyBinding(itemStr); err != nil {
					return &ValidationError{
						Field:   fmt.Sprintf("%s[%d]", fieldPath, i),
						Value:   itemStr,
						Message: err.Error(),
					}
				}
			}
		}
	default:
		return &ValidationError{
			Field:   fieldPath,
			Value:   value,
			Message: "keybinding must be a string or array of strings",
		}
	}
	return nil
}

// parseKeyBinding validates key binding strings.
// This simple validation is implemented here to avoid a circular import:
// importing the full keybinding parser from the 'cmd' (interactive UI) package
// would cause a circular dependency, since that package depends on 'config'.
func parseKeyBinding(keyStr string) error { //nolint:revive // parsing multiple legacy formats
	s := strings.TrimSpace(keyStr)
	if s == "" {
		return fmt.Errorf("empty key binding")
	}

	// Basic validation - check for supported formats
	sLower := strings.ToLower(s)

	// Accept ctrl+<key>, ^<key>, or c-<key> formats
	if (strings.HasPrefix(sLower, "ctrl+") && len(s) >= 6) ||
		(strings.HasPrefix(s, "^") && len(s) == 2) ||
		(strings.HasPrefix(sLower, "c-") && len(s) == 3) {
		return nil
	}

	return fmt.Errorf("unsupported key binding format: %s (supported: 'ctrl+<key>', '^<key>', 'c-<key>')", keyStr)
}
