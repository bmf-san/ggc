package keybindings

// applyDefaults sets up default keybindings
func (r *KeyBindingResolver) applyDefaults(keyMap *KeyBindingMap) {
	// Apply hardcoded defaults (legacy compatibility)
	defaults := DefaultKeyBindingMap()
	keyMap.DeleteWord = append(keyMap.DeleteWord, defaults.DeleteWord...)
	keyMap.ClearLine = append(keyMap.ClearLine, defaults.ClearLine...)
	keyMap.DeleteToEnd = append(keyMap.DeleteToEnd, defaults.DeleteToEnd...)
	keyMap.MoveToBeginning = append(keyMap.MoveToBeginning, defaults.MoveToBeginning...)
	keyMap.MoveToEnd = append(keyMap.MoveToEnd, defaults.MoveToEnd...)
	keyMap.MoveUp = append(keyMap.MoveUp, defaults.MoveUp...)
	keyMap.MoveDown = append(keyMap.MoveDown, defaults.MoveDown...)
	keyMap.AddToWorkflow = append(keyMap.AddToWorkflow, defaults.AddToWorkflow...)
	keyMap.ToggleWorkflowView = append(keyMap.ToggleWorkflowView, defaults.ToggleWorkflowView...)
	keyMap.ClearWorkflow = append(keyMap.ClearWorkflow, defaults.ClearWorkflow...)
	keyMap.WorkflowCreate = append(keyMap.WorkflowCreate, defaults.WorkflowCreate...)
	keyMap.WorkflowDelete = append(keyMap.WorkflowDelete, defaults.WorkflowDelete...)
	keyMap.SoftCancel = append(keyMap.SoftCancel, defaults.SoftCancel...)
}

func (r *KeyBindingResolver) applyProfile(keyMap *KeyBindingMap, profile *KeyBindingProfile, context Context) {
	// Helper function to apply bindings from profile
	applyBinding := func(action string, target *[]KeyStroke) {
		if keystrokes, exists := profile.GetBinding(context, action); exists {
			*target = keystrokes // Replace, don't append (profile overrides defaults)
		}
	}

	applyBinding("delete_word", &keyMap.DeleteWord)
	applyBinding("clear_line", &keyMap.ClearLine)
	applyBinding("delete_to_end", &keyMap.DeleteToEnd)
	applyBinding("move_to_beginning", &keyMap.MoveToBeginning)
	applyBinding("move_to_end", &keyMap.MoveToEnd)
	applyBinding("move_up", &keyMap.MoveUp)
	applyBinding("move_down", &keyMap.MoveDown)
	applyBinding("move_left", &keyMap.MoveLeft)
	applyBinding("move_right", &keyMap.MoveRight)
	applyBinding("add_to_workflow", &keyMap.AddToWorkflow)
	applyBinding("toggle_workflow_view", &keyMap.ToggleWorkflowView)
	applyBinding("clear_workflow", &keyMap.ClearWorkflow)
	applyBinding("workflow_create", &keyMap.WorkflowCreate)
	applyBinding("workflow_delete", &keyMap.WorkflowDelete)
	applyBinding("soft_cancel", &keyMap.SoftCancel)
}

func (r *KeyBindingResolver) applyPlatformLayer(keyMap *KeyBindingMap) {
	platformBindings := GetPlatformSpecificKeyBindings(r.platform)

	// Apply platform-specific overrides
	if bindings, exists := platformBindings["delete_word"]; exists {
		keyMap.DeleteWord = bindings // Platform overrides profile
	}
}

func (r *KeyBindingResolver) applyTerminalLayer(keyMap *KeyBindingMap) {
	terminalBindings := GetTerminalSpecificKeyBindings(r.terminal)

	// Apply terminal-specific overrides with explicit action handling
	for action, bindings := range terminalBindings {
		r.applyTerminalBinding(keyMap, action, bindings)
	}
}

// applyTerminalBinding applies a single terminal binding to reduce cyclomatic complexity
func (r *KeyBindingResolver) applyTerminalBinding(keyMap *KeyBindingMap, action string, bindings []KeyStroke) {
	// Apply editing actions
	if r.applyEditingAction(keyMap, action, bindings) {
		return
	}

	// Apply navigation actions
	if r.applyNavigationAction(keyMap, action, bindings) {
		return
	}

	// Apply workflow actions
	r.applyWorkflowAction(keyMap, action, bindings)
}

// applyEditingAction applies editing-related keybinding actions
func (r *KeyBindingResolver) applyEditingAction(keyMap *KeyBindingMap, action string, bindings []KeyStroke) bool {
	switch action {
	case "delete_word":
		keyMap.DeleteWord = bindings
		return true
	case "clear_line":
		keyMap.ClearLine = bindings
		return true
	case "delete_to_end":
		keyMap.DeleteToEnd = bindings
		return true
	}
	return false
}

// applyNavigationAction applies navigation-related keybinding actions
func (r *KeyBindingResolver) applyNavigationAction(keyMap *KeyBindingMap, action string, bindings []KeyStroke) bool {
	switch action {
	case "move_to_beginning":
		keyMap.MoveToBeginning = bindings
		return true
	case "move_to_end":
		keyMap.MoveToEnd = bindings
		return true
	case "move_up":
		keyMap.MoveUp = bindings
		return true
	case "move_down":
		keyMap.MoveDown = bindings
		return true
	case "move_left":
		keyMap.MoveLeft = bindings
		return true
	case "move_right":
		keyMap.MoveRight = bindings
		return true
	}
	return false
}

// applyWorkflowAction applies workflow-related keybinding actions
func (r *KeyBindingResolver) applyWorkflowAction(keyMap *KeyBindingMap, action string, bindings []KeyStroke) {
	actionMap := map[string]*[]KeyStroke{
		"add_to_workflow":      &keyMap.AddToWorkflow,
		"toggle_workflow_view": &keyMap.ToggleWorkflowView,
		"clear_workflow":       &keyMap.ClearWorkflow,
		"workflow_create":      &keyMap.WorkflowCreate,
		"workflow_delete":      &keyMap.WorkflowDelete,
		"soft_cancel":          &keyMap.SoftCancel,
	}

	if target, exists := actionMap[action]; exists {
		*target = bindings
	}
}
