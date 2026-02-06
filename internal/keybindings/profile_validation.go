package keybindings

import "fmt"

func ValidateProfile(profile *KeyBindingProfile) error { //nolint:revive // performs exhaustive validation checks
	if profile == nil {
		return fmt.Errorf("profile is nil")
	}

	if profile.Name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}

	if profile.Description == "" {
		return fmt.Errorf("profile description cannot be empty")
	}

	if profile.Contexts == nil {
		return fmt.Errorf("profile contexts cannot be nil")
	}

	// Validate that profile has required contexts
	requiredContexts := []Context{ContextGlobal, ContextInput, ContextResults, ContextSearch}
	for _, requiredCtx := range requiredContexts {
		if _, exists := profile.Contexts[requiredCtx]; !exists {
			return fmt.Errorf("profile missing required context: %s", requiredCtx)
		}
	}

	// Validate that each context has at least basic navigation bindings
	if inputBindings, exists := profile.Contexts[ContextInput]; exists {
		requiredInputActions := []string{"move_to_beginning", "move_to_end", "delete_word", "clear_line"}
		for _, action := range requiredInputActions {
			if _, hasAction := inputBindings[action]; !hasAction {
				return fmt.Errorf("profile input context missing required action: %s", action)
			}
		}
	}

	if resultsBindings, exists := profile.Contexts[ContextResults]; exists {
		requiredResultsActions := []string{"move_up", "move_down"}
		for _, action := range requiredResultsActions {
			if _, hasAction := resultsBindings[action]; !hasAction {
				return fmt.Errorf("profile results context missing required action: %s", action)
			}
		}
	}

	// Validate KeyStroke consistency
	for contextName, contextBindings := range profile.Contexts {
		for action, keystrokes := range contextBindings {
			if len(keystrokes) == 0 {
				return fmt.Errorf("profile %s context %s action %s has no keystrokes", profile.Name, contextName, action)
			}
			for i, ks := range keystrokes {
				if err := validateKeyStroke(ks); err != nil {
					return fmt.Errorf("profile %s context %s action %s keystroke %d invalid: %w", profile.Name, contextName, action, i, err)
				}
			}
		}
	}

	return nil
}

// ValidateAllBuiltinProfiles validates all built-in profiles
func ValidateAllBuiltinProfiles() error {
	profiles := map[Profile]func() *KeyBindingProfile{
		ProfileDefault:  CreateDefaultProfile,
		ProfileEmacs:    CreateEmacsProfile,
		ProfileVi:       CreateViProfile,
		ProfileReadline: CreateReadlineProfile,
	}

	for profileName, creator := range profiles {
		profile := creator()
		if err := ValidateProfile(profile); err != nil {
			return fmt.Errorf("built-in profile %s validation failed: %w", profileName, err)
		}
	}

	return nil
}
