package interactive

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/term"

	kb "github.com/bmf-san/ggc/v7/internal/keybindings"
)

func (h *KeyHandler) handleWorkflowKeys(r rune, oldState *term.State) (bool, bool, []string) {
	switch h.ui.state.mode {
	case ModeWorkflow:
		return h.handleWorkflowModeKeys(r, oldState)
	case ModeSearch:
		return h.handleSearchModeWorkflowKeys(r)
	default:
		return false, true, nil
	}
}

func (h *KeyHandler) handleSearchModeWorkflowKeys(r rune) (bool, bool, []string) {
	km := h.GetCurrentKeyMap()
	keyStroke := kb.NewCharKeyStroke(r)

	if km.MatchesKeyStroke("add_to_workflow", keyStroke) {
		if h.ui.state.HasInput() {
			if cmd := h.ui.state.GetSelectedCommand(); cmd != nil {
				h.addCommandToWorkflow(cmd.Command)
				h.ui.state.ClearInput()
			}
		}
		return true, true, nil
	}
	return false, true, nil
}

func (h *KeyHandler) handleWorkflowModeKeys(r rune, oldState *term.State) (bool, bool, []string) {
	if handled := h.handleWorkflowModeShortcut(r, oldState); handled {
		return true, true, nil
	}
	if handled := h.handleWorkflowModeBindings(r); handled {
		return true, true, nil
	}
	return false, true, nil
}

func (h *KeyHandler) handleWorkflowModeShortcut(r rune, oldState *term.State) bool {
	switch r {
	case 'x':
		h.executeWorkflow(oldState)
		return true
	case 'n':
		h.createWorkflow()
		return true
	case 'd':
		h.deleteActiveWorkflow()
		return true
	}
	return false
}

func (h *KeyHandler) handleWorkflowModeBindings(r rune) bool {
	keyStroke := kb.NewCharKeyStroke(r)

	if h.handleWorkflowAdd(keyStroke) {
		return true
	}
	if h.handleWorkflowClear(keyStroke) {
		return true
	}
	return false
}

func (h *KeyHandler) handleWorkflowAdd(keyStroke kb.KeyStroke) bool {
	if !h.ui.state.IsInputFocused() || !h.ui.state.HasInput() {
		return false
	}
	km := h.GetCurrentKeyMap()
	if !km.MatchesKeyStroke("add_to_workflow", keyStroke) {
		return false
	}
	if cmd := h.ui.state.GetSelectedCommand(); cmd != nil {
		h.addCommandToWorkflow(cmd.Command)
		h.ui.state.ClearInput()
	}
	return true
}

func (h *KeyHandler) handleWorkflowClear(keyStroke kb.KeyStroke) bool {
	if h.ui.state.IsInputFocused() {
		return false
	}
	km := h.GetCurrentKeyMap()
	if !km.MatchesKeyStroke("clear_workflow", keyStroke) {
		return false
	}
	h.clearWorkflow()
	return true
}

// handleControlChar processes control characters and returns (handled, shouldContinue, result)
// Enhanced to support KeyStroke matching while maintaining backward compatibility
//
//nolint:revive // Control character handling inherently requires many cases
func (h *KeyHandler) moveWorkflowList(delta int) {
	summaries := h.ui.listWorkflows()
	if len(summaries) == 0 {
		return
	}
	h.ui.state.SetWorkflowListIndex(h.ui.state.workflowListIdx+delta, len(summaries))
	idx := h.ui.state.workflowListIdx
	if idx < 0 || idx >= len(summaries) {
		return
	}
	selected := summaries[idx]
	if h.ui.workflowMgr.SetActive(selected.ID) {
		h.ui.updateWorkflowPointer()
	}
}

func (h *KeyHandler) createWorkflow() {
	if h.ui.workflowMgr == nil {
		return
	}
	newID := h.ui.workflowMgr.CreateWorkflow("")
	summaries := h.ui.listWorkflows()
	for i, summary := range summaries {
		if summary.ID == newID {
			h.ui.state.SetWorkflowListIndex(i, len(summaries))
			break
		}
	}
	h.ui.updateWorkflowPointer()
	h.ui.write("%s✨ Created workflow #%d%s\n", h.ui.colors.BrightGreen, newID, h.ui.colors.Reset)
}

func (h *KeyHandler) deleteActiveWorkflow() {
	if h.ui.workflowMgr == nil {
		return
	}
	activeID := h.ui.workflowMgr.GetActiveID()
	if activeID == 0 {
		h.ui.write("%sNo active workflow to delete%s\n", h.ui.colors.BrightYellow, h.ui.colors.Reset)
		return
	}
	newActive, ok := h.ui.workflowMgr.DeleteWorkflow(activeID)
	if !ok {
		h.ui.write("%sUnable to delete workflow #%d%s\n", h.ui.colors.BrightYellow, activeID, h.ui.colors.Reset)
		return
	}
	summaries := h.ui.listWorkflows()
	if newActive == 0 {
		h.ui.state.SetWorkflowListIndex(0, len(summaries))
	} else {
		for i, summary := range summaries {
			if summary.ID == newActive {
				h.ui.state.SetWorkflowListIndex(i, len(summaries))
				break
			}
		}
	}
	h.ui.updateWorkflowPointer()
	h.ui.write("%s🗑  Deleted workflow #%d%s\n", h.ui.colors.BrightYellow, activeID, h.ui.colors.Reset)
}

// readNextByte reads the next byte from either a buffered reader or stdin
func (h *KeyHandler) addCommandToWorkflow(cmdTemplate string) {
	// Don't process placeholders here - save the template as-is
	// Placeholders will be resolved during workflow execution

	// Parse command and arguments from template
	parts := strings.Fields(cmdTemplate)
	if len(parts) == 0 {
		return
	}

	command := parts[0]
	args := parts[1:]

	// Add template to workflow (with placeholders intact)
	id := h.ui.AddToWorkflow(command, args, cmdTemplate)

	// Show success message
	placeholders := extractPlaceholders(cmdTemplate)
	if len(placeholders) > 0 {
		h.ui.write("\n%s🎯 Added to workflow!%s\n",
			h.ui.colors.BrightGreen+h.ui.colors.Bold, h.ui.colors.Reset)
		h.ui.write("%s  Step %d: %s%s%s %s(will prompt for: %v)%s\n",
			h.ui.colors.BrightCyan, id, h.ui.colors.BrightWhite+h.ui.colors.Bold, cmdTemplate, h.ui.colors.Reset,
			h.ui.colors.BrightYellow, placeholders, h.ui.colors.Reset)
	} else {
		h.ui.write("\n%s🎯 Added to workflow!%s\n",
			h.ui.colors.BrightGreen+h.ui.colors.Bold, h.ui.colors.Reset)
		h.ui.write("%s  Step %d: %s%s%s\n",
			h.ui.colors.BrightCyan, id, h.ui.colors.BrightWhite+h.ui.colors.Bold, cmdTemplate, h.ui.colors.Reset)
	}
	h.ui.write("%s  Press 'Ctrl+t' to view workflow, or continue adding more commands%s\n\n",
		h.ui.colors.BrightBlack, h.ui.colors.Reset)
}

// clearWorkflow clears all steps from workflow
func (h *KeyHandler) clearWorkflow() {
	h.ui.ClearWorkflow()
	h.ui.write("%s🧹 Workflow cleared%s\n", h.ui.colors.BrightYellow, h.ui.colors.Reset)
}

// executeWorkflow executes the current workflow
func (h *KeyHandler) executeWorkflow(oldState *term.State) {
	if h.ui.workflow == nil {
		h.ui.notifyWorkflowError("No active workflow. Press Ctrl+N to create one.", 3*time.Second)
		return
	}
	if h.ui.workflow.IsEmpty() {
		h.ui.notifyWorkflowError("Workflow is empty. Add some steps first!", 3*time.Second)
		return
	}

	// Restore terminal state before execution
	h.restoreTerminalState(oldState)

	// Clear screen and execute workflow
	clearScreen(h.ui.stdout)

	err := h.ui.ExecuteWorkflow()
	if errors.Is(err, ErrWorkflowCanceled) {
		h.handleSoftCancel(oldState)
		h.reenterRawMode(oldState)
		return
	}
	if err != nil {
		h.ui.notifyWorkflowError(fmt.Sprintf("Workflow execution failed: %v", err), 3*time.Second)
		h.reenterRawMode(oldState)
		return
	}

	h.ui.notifyWorkflowSuccess("Workflow preserved for reuse. Press 'Ctrl+t' to view or modify.", 3*time.Second)
	h.reenterRawMode(oldState)

	// Keep workflow for reuse - don't clear it
}
