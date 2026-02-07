package interactive

import (
	kb "github.com/bmf-san/ggc/v7/internal/keybindings"
)

// ToggleWorkflowView toggles between search mode and workflow mode.
func (ui *UI) ToggleWorkflowView() {
	if ui == nil || ui.state == nil {
		return
	}
	if ui.state.IsWorkflowMode() {
		ui.enterSearchMode()
		return
	}
	ui.enterWorkflowMode()
}

// enterWorkflowMode switches UI into workflow management mode.
// Workflow mode has no input field - it's a pure management interface.
func (ui *UI) enterWorkflowMode() {
	if ui == nil || ui.state == nil {
		return
	}
	ui.state.SetMode(ModeWorkflow)
	ui.state.SetContext(kb.ContextGlobal)
	ui.ensureWorkflowListSelection()
	ui.updateWorkflowPointer()
}

// enterSearchMode switches UI back to search mode.
func (ui *UI) enterSearchMode() {
	if ui == nil || ui.state == nil {
		return
	}
	ui.state.SetMode(ModeSearch)
	ui.state.FocusInput()
	ui.state.SetContext(kb.ContextGlobal)
}

// resetToSearchMode clears the interactive search UI back to its default state.
func (ui *UI) resetToSearchMode() bool {
	if ui == nil || ui.state == nil {
		return false
	}

	state := ui.state
	active := state.HasInput() || state.IsWorkflowMode() || len(state.contextStack) > 0 || state.GetCurrentContext() != kb.ContextGlobal
	state.ClearInput()
	state.selected = 0
	state.contextStack = nil
	state.SetContext(kb.ContextGlobal)
	state.SetMode(ModeSearch)
	state.FocusInput()
	return active
}

// ResetToSearchMode clears the interactive search UI back to its default state.
func (ui *UI) ResetToSearchMode() bool {
	return ui.resetToSearchMode()
}

// readPlaceholderInput reads input for placeholder replacement
func (ui *UI) readPlaceholderInput() (string, bool) {
	if ui == nil || ui.handler == nil {
		return "", true
	}
	return ui.handler.getRealTimeInput()
}

// ApplyContextualKeybindings updates the active keybinding map, satisfying keybindings.ContextualMapApplier.
func (ui *UI) ApplyContextualKeybindings(contextual *kb.ContextualKeyBindingMap) {
	if ui == nil || ui.handler == nil || contextual == nil {
		return
	}
	ui.handler.contextualMap = contextual
}
