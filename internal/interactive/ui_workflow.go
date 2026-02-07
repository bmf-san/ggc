package interactive

import (
	"fmt"
)

// AddToWorkflow adds a command to the active workflow.
func (ui *UI) AddToWorkflow(command string, args []string, description string) int {
	if ui.workflowMgr == nil {
		return 0
	}
	activeID := ui.workflowMgr.GetActiveID()
	if id, ok := ui.workflowMgr.AddStep(activeID, command, args, description); ok {
		ui.updateWorkflowPointer()
		return id
	}
	return 0
}

// ClearWorkflow removes all steps from the workflow
func (ui *UI) ClearWorkflow() {
	if ui.workflowMgr == nil {
		return
	}
	ui.workflowMgr.ClearWorkflow(ui.workflowMgr.GetActiveID())
	ui.updateWorkflowPointer()
}

// ExecuteWorkflow executes the current workflow
func (ui *UI) ExecuteWorkflow() error {
	if ui.workflowEx == nil {
		return fmt.Errorf("workflow executor not initialized")
	}

	if ui.workflow == nil || ui.workflow.IsEmpty() {
		return fmt.Errorf("workflow is empty")
	}

	return ui.workflowEx.Execute(ui.workflow)
}

// updateWorkflowPointer updates the workflow pointer from the workflow manager
func (ui *UI) updateWorkflowPointer() {
	if ui == nil || ui.workflowMgr == nil {
		return
	}
	wf, ok := ui.workflowMgr.GetWorkflow(ui.workflowMgr.GetActiveID())
	if ok {
		ui.workflow = wf
		return
	}
	ui.workflow = nil
}

// listWorkflows returns a list of all workflows
func (ui *UI) listWorkflows() []WorkflowSummary {
	if ui.workflowMgr == nil {
		return nil
	}
	return ui.workflowMgr.ListWorkflows()
}

// ensureWorkflowListSelection ensures the workflow list selection is valid
func (ui *UI) ensureWorkflowListSelection() {
	if ui == nil || ui.state == nil {
		return
	}
	summaries := ui.listWorkflows()
	activeID := 0
	if ui.workflowMgr != nil {
		activeID = ui.workflowMgr.GetActiveID()
	}
	if activeID != 0 {
		for i, summary := range summaries {
			if summary.ID == activeID {
				ui.state.workflowListIdx = i
				break
			}
		}
	}
	ui.state.SetWorkflowListIndex(ui.state.workflowListIdx, len(summaries))
}
