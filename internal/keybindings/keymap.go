package keybindings

// KeyBindingMap holds resolved key strokes for interactive actions.
// Supports multiple key strokes per action while maintaining backward compatibility.
type KeyBindingMap struct {
	DeleteWord         []KeyStroke // default: [Ctrl+W]
	ClearLine          []KeyStroke // default: [Ctrl+U]
	DeleteToEnd        []KeyStroke // default: [Ctrl+K]
	MoveToBeginning    []KeyStroke // default: [Ctrl+A]
	MoveToEnd          []KeyStroke // default: [Ctrl+E]
	MoveUp             []KeyStroke // default: [Ctrl+P], can add: [up arrow]
	MoveDown           []KeyStroke // default: [Ctrl+N], can add: [down arrow]
	MoveLeft           []KeyStroke // default: [], can add: [left arrow] for cursor movement
	MoveRight          []KeyStroke // default: [], can add: [right arrow] for cursor movement
	AddToWorkflow      []KeyStroke // default: [Tab]
	ToggleWorkflowView []KeyStroke // default: [Ctrl+T]
	ClearWorkflow      []KeyStroke // default: [c]
	WorkflowCreate     []KeyStroke // default: [Ctrl+N]
	WorkflowDelete     []KeyStroke // default: [Ctrl+D]
	SoftCancel         []KeyStroke // default: [Ctrl+G, Esc]
}

// DefaultKeyBindingMap returns the built-in default control bindings.
func DefaultKeyBindingMap() *KeyBindingMap {
	return &KeyBindingMap{
		DeleteWord:         []KeyStroke{NewCtrlKeyStroke('w')},
		ClearLine:          []KeyStroke{NewCtrlKeyStroke('u')},
		DeleteToEnd:        []KeyStroke{NewCtrlKeyStroke('k')},
		MoveToBeginning:    []KeyStroke{NewCtrlKeyStroke('a')},
		MoveToEnd:          []KeyStroke{NewCtrlKeyStroke('e')},
		MoveUp:             []KeyStroke{NewCtrlKeyStroke('p')},
		MoveDown:           []KeyStroke{NewCtrlKeyStroke('n')},
		MoveLeft:           []KeyStroke{}, // Empty by default, users can add left arrow
		MoveRight:          []KeyStroke{}, // Empty by default, users can add right arrow
		AddToWorkflow:      []KeyStroke{NewTabKeyStroke()},
		ToggleWorkflowView: []KeyStroke{NewCtrlKeyStroke('t')},
		ClearWorkflow:      []KeyStroke{NewCharKeyStroke('c')},
		WorkflowCreate:     []KeyStroke{NewCtrlKeyStroke('n')},
		WorkflowDelete:     []KeyStroke{NewCtrlKeyStroke('d')},
		SoftCancel:         []KeyStroke{NewCtrlKeyStroke('g'), NewEscapeKeyStroke()},
	}
}

// Legacy backward-compatibility methods maintain the old byte-based API
// while internally using the new KeyStroke system.

// GetDeleteWordByte returns the primary control byte for DeleteWord (backward compatibility)
func (km *KeyBindingMap) GetDeleteWordByte() byte {
	return km.getFirstControlByte(km.DeleteWord, ctrl('w'))
}

// GetClearLineByte returns the primary control byte for ClearLine (backward compatibility)
func (km *KeyBindingMap) GetClearLineByte() byte {
	return km.getFirstControlByte(km.ClearLine, ctrl('u'))
}

// GetDeleteToEndByte returns the primary control byte for DeleteToEnd (backward compatibility)
func (km *KeyBindingMap) GetDeleteToEndByte() byte {
	return km.getFirstControlByte(km.DeleteToEnd, ctrl('k'))
}

// GetMoveToBeginningByte returns the primary control byte for MoveToBeginning (backward compatibility)
func (km *KeyBindingMap) GetMoveToBeginningByte() byte {
	return km.getFirstControlByte(km.MoveToBeginning, ctrl('a'))
}

// GetMoveToEndByte returns the primary control byte for MoveToEnd (backward compatibility)
func (km *KeyBindingMap) GetMoveToEndByte() byte {
	return km.getFirstControlByte(km.MoveToEnd, ctrl('e'))
}

// GetMoveUpByte returns the primary control byte for MoveUp (backward compatibility)
func (km *KeyBindingMap) GetMoveUpByte() byte {
	return km.getFirstControlByte(km.MoveUp, ctrl('p'))
}

// GetMoveDownByte returns the primary control byte for MoveDown (backward compatibility)
func (km *KeyBindingMap) GetMoveDownByte() byte {
	return km.getFirstControlByte(km.MoveDown, ctrl('n'))
}

// GetAddToWorkflowByte returns the primary control byte for AddToWorkflow (backward compatibility)
func (km *KeyBindingMap) GetAddToWorkflowByte() byte {
	return km.getFirstControlByte(km.AddToWorkflow, 9) // Tab key
}

// GetToggleWorkflowViewByte returns the primary control byte for ToggleWorkflowView (backward compatibility)
func (km *KeyBindingMap) GetToggleWorkflowViewByte() byte {
	return km.getFirstControlByte(km.ToggleWorkflowView, ctrl('t'))
}

// GetClearWorkflowByte returns the primary control byte for ClearWorkflow (backward compatibility)
func (km *KeyBindingMap) GetClearWorkflowByte() byte {
	return km.getFirstControlByte(km.ClearWorkflow, 'c')
}

// getFirstControlByte finds the first Ctrl KeyStroke and returns its control byte,
// or returns the fallback if none found
func (km *KeyBindingMap) getFirstControlByte(keyStrokes []KeyStroke, fallback byte) byte {
	for _, ks := range keyStrokes {
		if b := ks.ToControlByte(); b != 0 {
			return b
		}
	}
	return fallback
}

// MatchesKeyStroke checks if any KeyStroke in the given action matches the input
func (km *KeyBindingMap) MatchesKeyStroke(action string, input KeyStroke) bool {
	actionMap := map[string][]KeyStroke{
		"delete_word":          km.DeleteWord,
		"clear_line":           km.ClearLine,
		"delete_to_end":        km.DeleteToEnd,
		"move_to_beginning":    km.MoveToBeginning,
		"move_to_end":          km.MoveToEnd,
		"move_up":              km.MoveUp,
		"move_down":            km.MoveDown,
		"move_left":            km.MoveLeft,
		"move_right":           km.MoveRight,
		"add_to_workflow":      km.AddToWorkflow,
		"toggle_workflow_view": km.ToggleWorkflowView,
		"clear_workflow":       km.ClearWorkflow,
		"workflow_create":      km.WorkflowCreate,
		"workflow_delete":      km.WorkflowDelete,
		"soft_cancel":          km.SoftCancel,
	}

	keyStrokes, exists := actionMap[action]
	if !exists {
		return false
	}

	for _, ks := range keyStrokes {
		if input.Equals(ks) {
			return true
		}
	}
	return false
}
