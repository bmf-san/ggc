package keybindings

// CreateDefaultProfile returns the default keybinding profile compatible with legacy behavior.
func CreateDefaultProfile() *KeyBindingProfile {
	return &KeyBindingProfile{
		Name:        "Default",
		Description: "Default keybindings compatible with legacy behavior",
		Global:      make(map[string][]KeyStroke),
		Contexts: map[Context]map[string][]KeyStroke{
			ContextGlobal: {
				"soft_cancel": {NewCtrlKeyStroke('g'), NewEscapeKeyStroke()},
			},
			ContextInput: {
				"delete_word":       {NewCtrlKeyStroke('w')},
				"clear_line":        {NewCtrlKeyStroke('u')},
				"delete_to_end":     {NewCtrlKeyStroke('k')},
				"move_to_beginning": {NewCtrlKeyStroke('a')},
				"move_to_end":       {NewCtrlKeyStroke('e')},
			},
			ContextResults: {
				"move_up":              {NewCtrlKeyStroke('p')},
				"move_down":            {NewCtrlKeyStroke('n')},
				"add_to_workflow":      {NewTabKeyStroke()},
				"toggle_workflow_view": {NewCtrlKeyStroke('t')},
				"clear_workflow":       {NewCharKeyStroke('c')},
			},
			ContextSearch: {
				"move_up":              {NewCtrlKeyStroke('p')},
				"move_down":            {NewCtrlKeyStroke('n')},
				"add_to_workflow":      {NewTabKeyStroke()},
				"toggle_workflow_view": {NewCtrlKeyStroke('t')},
				"clear_workflow":       {NewCharKeyStroke('c')},
			},
		},
	}
}
