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
				// History recall / search are also bound globally so
				// they fire on the fresh empty prompt (initial state
				// is ContextGlobal until the user transitions into
				// ContextInput by typing a rune). Without this overlay
				// Ctrl+P/N/R would silently be claimed by move_up /
				// move_down instead.
				"history_prev":   {NewCtrlKeyStroke('p')},
				"history_next":   {NewCtrlKeyStroke('n')},
				"history_search": {NewCtrlKeyStroke('r')},
			},
			ContextInput: {
				"delete_word":       {NewCtrlKeyStroke('w')},
				"clear_line":        {NewCtrlKeyStroke('u')},
				"delete_to_end":     {NewCtrlKeyStroke('k')},
				"move_to_beginning": {NewCtrlKeyStroke('a')},
				"move_to_end":       {NewCtrlKeyStroke('e')},
				// Ctrl+P/N walk the persisted command history while the
				// user is still editing the input buffer. These bindings
				// only exist in ContextInput so that the same chord can
				// continue to drive list navigation once the user starts
				// filtering (ContextSearch/ContextResults).
				"history_prev":   {NewCtrlKeyStroke('p')},
				"history_next":   {NewCtrlKeyStroke('n')},
				"history_search": {NewCtrlKeyStroke('r')},
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
				// Ctrl+R remains available while filtering so the user
				// can promote a partial query into a reverse history
				// search without having to clear the buffer first.
				"history_search": {NewCtrlKeyStroke('r')},
			},
		},
	}
}
