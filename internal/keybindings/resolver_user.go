package keybindings

import "os"

func (r *KeyBindingResolver) applyUserConfig(keyMap *KeyBindingMap, context Context) { //nolint:revive // layered override logic retained for clarity
	// Apply user global keybindings first
	userBindings := r.userConfig.Interactive.Keybindings

	userValues := map[string]string{
		"delete_word":          userBindings.DeleteWord,
		"clear_line":           userBindings.ClearLine,
		"delete_to_end":        userBindings.DeleteToEnd,
		"move_to_beginning":    userBindings.MoveToBeginning,
		"move_to_end":          userBindings.MoveToEnd,
		"move_up":              userBindings.MoveUp,
		"move_down":            userBindings.MoveDown,
		"move_left":            userBindings.MoveLeft,
		"move_right":           userBindings.MoveRight,
		"add_to_workflow":      userBindings.AddToWorkflow,
		"toggle_workflow_view": userBindings.ToggleWorkflowView,
		"clear_workflow":       userBindings.ClearWorkflow,
		"workflow_create":      userBindings.WorkflowCreate,
		"workflow_delete":      userBindings.WorkflowDelete,
		"soft_cancel":          userBindings.SoftCancel,
	}

	// Apply non-empty user overrides
	for action, keyStr := range userValues {
		if keyStr != "" {
			if ks, err := ParseKeyStroke(keyStr); err == nil {
				switch action {
				case "delete_word":
					keyMap.DeleteWord = []KeyStroke{ks}
				case "clear_line":
					keyMap.ClearLine = []KeyStroke{ks}
				case "delete_to_end":
					keyMap.DeleteToEnd = []KeyStroke{ks}
				case "move_to_beginning":
					keyMap.MoveToBeginning = []KeyStroke{ks}
				case "move_to_end":
					keyMap.MoveToEnd = []KeyStroke{ks}
				case "move_up":
					keyMap.MoveUp = []KeyStroke{ks}
				case "move_down":
					keyMap.MoveDown = []KeyStroke{ks}
				case "move_left":
					keyMap.MoveLeft = []KeyStroke{ks}
				case "move_right":
					keyMap.MoveRight = []KeyStroke{ks}
				case "add_to_workflow":
					keyMap.AddToWorkflow = []KeyStroke{ks}
				case "toggle_workflow_view":
					keyMap.ToggleWorkflowView = []KeyStroke{ks}
				case "clear_workflow":
					keyMap.ClearWorkflow = []KeyStroke{ks}
				case "workflow_create":
					keyMap.WorkflowCreate = []KeyStroke{ks}
				case "workflow_delete":
					keyMap.WorkflowDelete = []KeyStroke{ks}
				case "soft_cancel":
					keyMap.SoftCancel = []KeyStroke{ks}
				}
			}
		}
	}

	// Apply context-specific user bindings
	r.applyUserContextBindings(keyMap, context)

	// Apply platform-specific user bindings
	r.applyUserPlatformBindings(keyMap)

	// Apply terminal-specific user bindings
	r.applyUserTerminalBindings(keyMap)
}

func (r *KeyBindingResolver) applyEnvironmentOverrides(keyMap *KeyBindingMap) {
	// Check for environment variable overrides
	envOverrides := map[string]*[]KeyStroke{
		"GGC_KEYBIND_DELETE_WORD":          &keyMap.DeleteWord,
		"GGC_KEYBIND_CLEAR_LINE":           &keyMap.ClearLine,
		"GGC_KEYBIND_DELETE_TO_END":        &keyMap.DeleteToEnd,
		"GGC_KEYBIND_MOVE_TO_BEGINNING":    &keyMap.MoveToBeginning,
		"GGC_KEYBIND_MOVE_TO_END":          &keyMap.MoveToEnd,
		"GGC_KEYBIND_MOVE_UP":              &keyMap.MoveUp,
		"GGC_KEYBIND_MOVE_DOWN":            &keyMap.MoveDown,
		"GGC_KEYBIND_ADD_TO_WORKFLOW":      &keyMap.AddToWorkflow,
		"GGC_KEYBIND_TOGGLE_WORKFLOW_VIEW": &keyMap.ToggleWorkflowView,
		"GGC_KEYBIND_CLEAR_WORKFLOW":       &keyMap.ClearWorkflow,
		"GGC_KEYBIND_WORKFLOW_CREATE":      &keyMap.WorkflowCreate,
		"GGC_KEYBIND_WORKFLOW_DELETE":      &keyMap.WorkflowDelete,
		"GGC_KEYBIND_SOFT_CANCEL":          &keyMap.SoftCancel,
	}

	for envVar, target := range envOverrides {
		if keyStr := os.Getenv(envVar); keyStr != "" {
			if ks, err := ParseKeyStroke(keyStr); err == nil {
				*target = []KeyStroke{ks}
			}
		}
	}
}

func (r *KeyBindingResolver) applyUserContextBindings(keyMap *KeyBindingMap, context Context) {
	// Apply context-specific user bindings if they exist
	var contextBindings map[string]interface{}

	switch context {
	case ContextInput:
		contextBindings = r.userConfig.Interactive.Contexts.Input.Keybindings
	case ContextResults:
		contextBindings = r.userConfig.Interactive.Contexts.Results.Keybindings
	case ContextSearch:
		contextBindings = r.userConfig.Interactive.Contexts.Search.Keybindings
	}

	if contextBindings != nil {
		r.applyUserBindings(keyMap, contextBindings)
	}
}

func (r *KeyBindingResolver) applyUserPlatformBindings(keyMap *KeyBindingMap) {
	var platformBindings map[string]interface{}

	switch r.platform {
	case "darwin":
		platformBindings = r.userConfig.Interactive.Darwin.Keybindings
	case "linux":
		platformBindings = r.userConfig.Interactive.Linux.Keybindings
	case "windows":
		platformBindings = r.userConfig.Interactive.Windows.Keybindings
	}

	if platformBindings != nil {
		r.applyUserBindings(keyMap, platformBindings)
	}
}

func (r *KeyBindingResolver) applyUserTerminalBindings(keyMap *KeyBindingMap) {
	if r.userConfig.Interactive.Terminals != nil {
		if termConfig, exists := r.userConfig.Interactive.Terminals[r.terminal]; exists {
			if termConfig.Keybindings != nil {
				r.applyUserBindings(keyMap, termConfig.Keybindings)
			}
		}
	}
}

func (r *KeyBindingResolver) applyUserBindings(keyMap *KeyBindingMap, bindings map[string]interface{}) {
	for action, value := range bindings {
		keystrokes := r.parseUserBindingValue(value)
		if len(keystrokes) > 0 {
			r.applyUserBinding(keyMap, action, keystrokes)
		}
	}
}

// applyUserBinding applies a single user binding to reduce cyclomatic complexity
func (r *KeyBindingResolver) applyUserBinding(keyMap *KeyBindingMap, action string, keystrokes []KeyStroke) {
	// Apply editing actions
	if r.applyUserEditingAction(keyMap, action, keystrokes) {
		return
	}

	// Apply navigation actions
	if r.applyUserNavigationAction(keyMap, action, keystrokes) {
		return
	}

	// Apply workflow actions
	r.applyUserWorkflowAction(keyMap, action, keystrokes)
}

// applyUserEditingAction applies user editing-related keybinding actions
func (r *KeyBindingResolver) applyUserEditingAction(keyMap *KeyBindingMap, action string, keystrokes []KeyStroke) bool {
	switch action {
	case "delete_word":
		keyMap.DeleteWord = keystrokes
		return true
	case "clear_line":
		keyMap.ClearLine = keystrokes
		return true
	case "delete_to_end":
		keyMap.DeleteToEnd = keystrokes
		return true
	}
	return false
}

// applyUserNavigationAction applies user navigation-related keybinding actions
func (r *KeyBindingResolver) applyUserNavigationAction(keyMap *KeyBindingMap, action string, keystrokes []KeyStroke) bool {
	switch action {
	case "move_to_beginning":
		keyMap.MoveToBeginning = keystrokes
		return true
	case "move_to_end":
		keyMap.MoveToEnd = keystrokes
		return true
	case "move_up":
		keyMap.MoveUp = keystrokes
		return true
	case "move_down":
		keyMap.MoveDown = keystrokes
		return true
	case "move_left":
		keyMap.MoveLeft = keystrokes
		return true
	case "move_right":
		keyMap.MoveRight = keystrokes
		return true
	}
	return false
}

// applyUserWorkflowAction applies user workflow-related keybinding actions
func (r *KeyBindingResolver) applyUserWorkflowAction(keyMap *KeyBindingMap, action string, keystrokes []KeyStroke) {
	actionMap := map[string]*[]KeyStroke{
		"add_to_workflow":      &keyMap.AddToWorkflow,
		"toggle_workflow_view": &keyMap.ToggleWorkflowView,
		"clear_workflow":       &keyMap.ClearWorkflow,
		"workflow_create":      &keyMap.WorkflowCreate,
		"workflow_delete":      &keyMap.WorkflowDelete,
		"soft_cancel":          &keyMap.SoftCancel,
	}

	if target, exists := actionMap[action]; exists {
		*target = keystrokes
	}
}

func (r *KeyBindingResolver) parseUserBindingValue(value interface{}) []KeyStroke {
	switch v := value.(type) {
	case string:
		if v == "" {
			return []KeyStroke{}
		}
		if ks, err := ParseKeyStroke(v); err == nil {
			return []KeyStroke{ks}
		}
	case []interface{}:
		var keystrokes []KeyStroke
		for _, item := range v {
			if itemStr, ok := item.(string); ok && itemStr != "" {
				if ks, err := ParseKeyStroke(itemStr); err == nil {
					keystrokes = append(keystrokes, ks)
				}
			}
		}
		return keystrokes
	}
	return []KeyStroke{}
}
