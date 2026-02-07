package keybindings

// GetProfileStatistics returns statistics about a keybinding profile
func GetProfileStatistics(profile *KeyBindingProfile) map[string]interface{} {
	stats := make(map[string]interface{})

	if profile == nil {
		return stats
	}

	// Count total bindings
	totalBindings := 0
	contextStats := make(map[Context]int)

	for context, bindings := range profile.Contexts {
		count := len(bindings)
		contextStats[context] = count
		totalBindings += count
	}

	// Count global bindings
	globalBindings := 0
	if profile.Global != nil {
		globalBindings = len(profile.Global)
	}

	stats["profile_name"] = profile.Name
	stats["description"] = profile.Description
	stats["total_context_bindings"] = totalBindings
	stats["global_bindings"] = globalBindings
	stats["context_breakdown"] = contextStats
	stats["contexts_defined"] = len(profile.Contexts)

	// Calculate keystroke type distribution
	keystrokeTypes := make(map[KeyStrokeKind]int)
	for _, bindings := range profile.Contexts {
		for _, keystrokes := range bindings {
			for _, ks := range keystrokes {
				keystrokeTypes[ks.Kind]++
			}
		}
	}
	stats["keystroke_types"] = keystrokeTypes

	return stats
}

// CompareProfiles compares two profiles and returns differences
func CompareProfiles(profile1, profile2 *KeyBindingProfile) map[string]interface{} { //nolint:revive // comparison builds rich analysis report
	comparison := make(map[string]interface{})

	if profile1 == nil || profile2 == nil {
		comparison["error"] = "one or both profiles are nil"
		return comparison
	}

	comparison["profile1_name"] = profile1.Name
	comparison["profile2_name"] = profile2.Name

	// Compare contexts
	contexts1 := make(map[Context]bool)
	contexts2 := make(map[Context]bool)

	for ctx := range profile1.Contexts {
		contexts1[ctx] = true
	}
	for ctx := range profile2.Contexts {
		contexts2[ctx] = true
	}

	var uniqueToProfile1, uniqueToProfile2, sharedContexts []Context
	for ctx := range contexts1 {
		if contexts2[ctx] {
			sharedContexts = append(sharedContexts, ctx)
		} else {
			uniqueToProfile1 = append(uniqueToProfile1, ctx)
		}
	}
	for ctx := range contexts2 {
		if !contexts1[ctx] {
			uniqueToProfile2 = append(uniqueToProfile2, ctx)
		}
	}

	comparison["unique_to_profile1"] = uniqueToProfile1
	comparison["unique_to_profile2"] = uniqueToProfile2
	comparison["shared_contexts"] = sharedContexts

	// Compare action coverage in shared contexts
	actionComparison := make(map[Context]map[string]interface{})
	for _, ctx := range sharedContexts {
		bindings1 := profile1.Contexts[ctx]
		bindings2 := profile2.Contexts[ctx]

		actions1 := make(map[string]bool)
		actions2 := make(map[string]bool)

		for action := range bindings1 {
			actions1[action] = true
		}
		for action := range bindings2 {
			actions2[action] = true
		}

		var uniqueActions1, uniqueActions2, sharedActions []string
		for action := range actions1 {
			if actions2[action] {
				sharedActions = append(sharedActions, action)
			} else {
				uniqueActions1 = append(uniqueActions1, action)
			}
		}
		for action := range actions2 {
			if !actions1[action] {
				uniqueActions2 = append(uniqueActions2, action)
			}
		}

		actionComparison[ctx] = map[string]interface{}{
			"unique_to_profile1": uniqueActions1,
			"unique_to_profile2": uniqueActions2,
			"shared_actions":     sharedActions,
		}
	}

	comparison["action_comparison"] = actionComparison
	return comparison
}
