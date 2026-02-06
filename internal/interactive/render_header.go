package interactive

import (
	"fmt"
)

func (r *Renderer) renderSoftCancelFlash(ui *UI) {
	if !ui.consumeSoftCancelFlash() {
		return
	}
	alert := fmt.Sprintf("%s⚠️  Operation canceled%s", r.colors.BrightRed+r.colors.Bold, r.colors.Reset)
	r.writeColorln(ui, alert)
	r.writeColorln(ui, "")
}

func (r *Renderer) renderWorkflowError(ui *UI) {
	if ui == nil || ui.state == nil || !ui.state.IsWorkflowMode() {
		return
	}
	message := ui.workflowErrorMessage()
	if message == "" {
		return
	}
	alert := fmt.Sprintf("%s⚠️  %s%s", r.colors.BrightRed+r.colors.Bold, message, r.colors.Reset)
	r.writeColorln(ui, alert)
	r.writeColorln(ui, "")
}

func (r *Renderer) renderWorkflowNotice(ui *UI) {
	if ui == nil || ui.state == nil || !ui.state.IsWorkflowMode() {
		return
	}
	message := ui.workflowNoticeMessage()
	if message == "" {
		return
	}
	notice := fmt.Sprintf("%s%s%s", r.colors.BrightGreen+r.colors.Bold, message, r.colors.Reset)
	r.writeColorln(ui, notice)
	r.writeColorln(ui, "")
}

// renderWorkflowMode renders the workflow management screen.
// Simplified: no input field, just workflow list and keybinds.
func (r *Renderer) renderWorkflowMode(ui *UI, state *UIState) {
	r.writeEmptyLine()
	r.renderWorkflowList(ui, state)
	r.writeEmptyLine()
	r.renderWorkflowView(ui, state)
	r.writeEmptyLine()
	r.renderWorkflowModeKeybinds(ui, state)
}

// renderHeader renders the title, git status, and navigation subtitle
func (r *Renderer) renderHeader(ui *UI) {
	// Modern header with title
	titleText := "🚀 ggc Interactive Mode"
	if ui != nil && ui.state != nil && ui.state.IsWorkflowMode() {
		titleText = "📋 Workflow Mode"
	}
	title := fmt.Sprintf("%s%s%s",
		r.colors.BrightCyan+r.colors.Bold,
		titleText,
		r.colors.Reset)
	r.writeColorln(ui, title)

	// Git status information
	if ui.gitStatus != nil {
		r.renderGitStatus(ui, ui.gitStatus)
	}

	if ui != nil && ui.state != nil && ui.state.IsWorkflowMode() {
		r.renderWorkflowActiveSummary(ui)
	}
}

func (r *Renderer) renderWorkflowActiveSummary(ui *UI) {
	activeID := 0
	stepCount := 0
	if ui.workflowMgr != nil {
		activeID = ui.workflowMgr.GetActiveID()
		if wf, ok := ui.workflowMgr.GetWorkflow(activeID); ok && wf != nil {
			stepCount = wf.Size()
		}
	}

	if activeID == 0 {
		r.writeColorln(ui, fmt.Sprintf("%sActive:%s %s(none)%s",
			r.colors.BrightYellow+r.colors.Bold,
			r.colors.Reset,
			r.colors.BrightBlack,
			r.colors.Reset))
		return
	}

	r.writeColorln(ui, fmt.Sprintf("%sActive:%s %sW%d%s %s(%d step%s)%s",
		r.colors.BrightYellow+r.colors.Bold,
		r.colors.Reset,
		r.colors.BrightWhite+r.colors.Bold,
		activeID,
		r.colors.Reset,
		r.colors.BrightBlack,
		stepCount,
		pluralize(stepCount),
		r.colors.Reset))
}

// renderSearchPrompt renders the search input with cursor
