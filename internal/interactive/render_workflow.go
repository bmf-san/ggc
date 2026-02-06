package interactive

import (
	"fmt"
	"strings"
)

func workflowLineCounts(summaries []WorkflowSummary, maxStepPreview int) []int {
	lineCounts := make([]int, len(summaries))
	for i, summary := range summaries {
		preview := summary.StepCount
		if preview > maxStepPreview {
			preview = maxStepPreview
		}
		lines := 1 + preview
		if summary.StepCount > maxStepPreview {
			lines++
		}
		lineCounts[i] = lines
	}
	return lineCounts
}

func workflowMaxLines(height int) int {
	maxLines := height - 12
	if maxLines < 6 {
		maxLines = 6
	}
	return maxLines
}

func workflowWindowFromOffset(lineCounts []int, offset, maxLines int) int {
	endIndex := offset
	linesUsed := 0
	for i := offset; i < len(lineCounts); i++ {
		if linesUsed+lineCounts[i] > maxLines && i > offset {
			break
		}
		linesUsed += lineCounts[i]
		endIndex = i
	}
	return endIndex
}

func workflowWindowForSelection(state *UIState, lineCounts []int, maxLines int) (int, int) {
	if len(lineCounts) == 0 {
		return 0, -1
	}

	offset := state.workflowOffset
	if offset < 0 || offset >= len(lineCounts) {
		offset = 0
	}
	if state.workflowListIdx < offset {
		offset = state.workflowListIdx
	}

	endIndex := workflowWindowFromOffset(lineCounts, offset, maxLines)
	for state.workflowListIdx > endIndex && offset < len(lineCounts)-1 {
		offset++
		endIndex = workflowWindowFromOffset(lineCounts, offset, maxLines)
	}

	return offset, endIndex
}

func workflowStepsForSummary(ui *UI, summary WorkflowSummary) []WorkflowStep {
	if ui == nil || ui.workflowMgr == nil {
		return nil
	}
	if wf, ok := ui.workflowMgr.GetWorkflow(summary.ID); ok && wf != nil {
		return wf.GetSteps()
	}
	return nil
}

// renderWorkflowList renders the list of available workflows with selection state.
func (r *Renderer) renderWorkflowList(ui *UI, state *UIState) {
	summaries := ui.listWorkflows()
	ui.ensureWorkflowListSelection()

	r.writeColorln(ui, fmt.Sprintf("%s📋 Workflows%s", r.colors.BrightYellow+r.colors.Bold, r.colors.Reset))

	if len(summaries) == 0 {
		r.writeColorln(ui, fmt.Sprintf("  %sNo workflows yet. Press Ctrl+N to create a workflow.%s",
			r.colors.BrightBlack, r.colors.Reset))
		return
	}

	const maxStepPreview = 3

	lineCounts := workflowLineCounts(summaries, maxStepPreview)
	maxLines := workflowMaxLines(r.height)
	offset, endIndex := workflowWindowForSelection(state, lineCounts, maxLines)
	state.workflowOffset = offset
	if endIndex < offset {
		return
	}

	for i := offset; i <= endIndex; i++ {
		r.renderWorkflowSummary(ui, state, summaries[i], i, maxStepPreview)
	}
}

func (r *Renderer) renderWorkflowSummary(ui *UI, state *UIState, summary WorkflowSummary, index int, maxStepPreview int) {
	displayName := strings.TrimSpace(summary.Name)
	if displayName == "" {
		displayName = fmt.Sprintf("W%d", summary.ID)
	}

	activePrefix := " "
	if summary.IsActive {
		activePrefix = fmt.Sprintf("%s▶%s", r.colors.BrightCyan+r.colors.Bold, r.colors.Reset)
	}

	selectPrefix := " "
	if state.workflowListIdx == index {
		selectPrefix = fmt.Sprintf("%s>%s", r.colors.BrightWhite+r.colors.Bold, r.colors.Reset)
	}

	activeLabel := ""
	if summary.IsActive {
		activeLabel = fmt.Sprintf(" %s[Active]%s", r.colors.BrightCyan, r.colors.Reset)
	}

	line := fmt.Sprintf("%s%s %s%s%s %s(%d step%s)%s%s",
		selectPrefix,
		activePrefix,
		r.colors.BrightWhite+r.colors.Bold,
		displayName,
		r.colors.Reset,
		r.colors.BrightBlack,
		summary.StepCount,
		pluralize(summary.StepCount),
		r.colors.Reset,
		activeLabel,
	)
	r.writeColorln(ui, line)

	steps := workflowStepsForSummary(ui, summary)
	r.renderWorkflowStepPreview(ui, steps, maxStepPreview)
}

func (r *Renderer) renderWorkflowStepPreview(ui *UI, steps []WorkflowStep, maxStepPreview int) {
	if len(steps) == 0 {
		return
	}

	previewCount := len(steps)
	if previewCount > maxStepPreview {
		previewCount = maxStepPreview
	}

	for s := 0; s < previewCount; s++ {
		step := steps[s]
		description := strings.TrimSpace(step.Description)
		if description == "" {
			description = step.Command
			if len(step.Args) > 0 {
				description += " " + strings.Join(step.Args, " ")
			}
		}
		stepLine := fmt.Sprintf("  %s%d.%s %s%s%s",
			r.colors.BrightBlue+r.colors.Bold,
			s+1,
			r.colors.Reset,
			r.colors.BrightGreen,
			description,
			r.colors.Reset)
		r.writeColorln(ui, stepLine)
	}
	if len(steps) > previewCount {
		r.writeColorln(ui, fmt.Sprintf("  %s... +%d more%s",
			r.colors.BrightBlack,
			len(steps)-previewCount,
			r.colors.Reset))
	}
}

// renderWorkflowModeKeybinds renders keybinds available in workflow mode.
// Simplified: no focus-based dimming since there's no input field.
// The *UI and *UIState parameters are accepted for API consistency with other
// renderer methods and reserved for potential future use; they are
// intentionally ignored in this implementation.
func (r *Renderer) renderWorkflowModeKeybinds(_ *UI, _ *UIState) {
	keybinds := []struct{ key, desc string }{
		{"n", "Create new workflow"},
		{"d / Ctrl+D", "Delete active workflow"},
		{"x", "Execute active workflow"},
		{"Ctrl+n/p", "Navigate workflows"},
		{"Ctrl+t", "Return to Search Mode"},
		{"Ctrl+c", "Quit"},
	}

	r.writeColorln(nil, fmt.Sprintf("%s⌨️  %sWorkflow mode keybinds:%s",
		r.colors.BrightBlue, r.colors.BrightWhite+r.colors.Bold, r.colors.Reset))

	for _, kb := range keybinds {
		r.writeColorln(nil, fmt.Sprintf("   %s%s%s  %s%s%s",
			r.colors.BrightGreen+r.colors.Bold,
			kb.key,
			r.colors.Reset,
			r.colors.BrightBlack,
			kb.desc,
			r.colors.Reset))
	}
}

// renderCommandList renders the filtered command list
func (r *Renderer) renderCommandList(ui *UI, state *UIState) {
	// Clamp selection index to valid range
	if state.selected >= len(state.filtered) {
		state.selected = len(state.filtered) - 1
	}
	if state.selected < 0 {
		state.selected = 0
	}

	// Calculate maximum command length for consistent alignment
	maxCmdLen := r.calculateMaxCommandLength(state.filtered)

	for i, cmd := range state.filtered {
		r.renderCommandItem(ui, cmd, i, state.selected, maxCmdLen)
	}
}

// renderCommandItem renders a single command item
func (r *Renderer) renderCommandItem(ui *UI, cmd CommandInfo, index, selected, maxCmdLen int) {
	desc := cmd.Description
	if desc == "" {
		desc = "No description"
	}

	// Calculate padding for consistent command alignment
	paddingLen := maxCmdLen - len(cmd.Command)
	if paddingLen < 0 {
		paddingLen = 0
	}
	padding := strings.Repeat(" ", paddingLen)

	// Calculate available width for description
	usedWidth := 4 + len(cmd.Command) + len(padding) + 3 // prefix + command + padding + separator
	availableDescWidth := r.width - usedWidth
	if availableDescWidth < 10 {
		availableDescWidth = 10
	}

	// Truncate description if needed
	trimmedDesc := ellipsis(desc, availableDescWidth)

	if index == selected {
		// Selected item with modern highlighting
		selectedLine := fmt.Sprintf("%s▶ %s%s%s%s %s│%s %s%s%s",
			r.colors.BrightCyan+r.colors.Bold,
			r.colors.BrightWhite+r.colors.Bold+r.colors.Reverse,
			" "+cmd.Command+" ",
			r.colors.Reset,
			padding,
			r.colors.BrightBlue,
			r.colors.Reset,
			r.colors.BrightWhite,
			trimmedDesc,
			r.colors.Reset)
		r.writeColorln(ui, selectedLine)
	} else {
		// Regular item with improved styling
		regularLine := fmt.Sprintf("  %s%s%s%s %s│%s %s%s%s",
			r.colors.BrightGreen+r.colors.Bold,
			cmd.Command,
			r.colors.Reset,
			padding,
			r.colors.BrightBlack,
			r.colors.Reset,
			r.colors.BrightBlack,
			trimmedDesc,
			r.colors.Reset)
		r.writeColorln(ui, regularLine)
	}
}

// renderWorkflowView renders the detailed workflow view
func (r *Renderer) renderWorkflowView(ui *UI, _ *UIState) {
	if ui == nil || ui.workflow == nil {
		r.writeColorln(ui, fmt.Sprintf("%s📋 Workflow Details (0 steps)%s",
			r.colors.BrightYellow+r.colors.Bold,
			r.colors.Reset))
		r.writeColorln(ui, fmt.Sprintf("%s  No active workflow%s",
			r.colors.BrightBlack,
			r.colors.Reset))
		r.writeColorln(ui, "")
		return
	}
	steps := ui.workflow.GetSteps()

	// Detailed workflow header
	r.writeColorln(ui, fmt.Sprintf("%s📋 Workflow Details (%d steps)%s",
		r.colors.BrightYellow+r.colors.Bold,
		len(steps),
		r.colors.Reset))
	r.writeColorln(ui, "")

	if len(steps) == 0 {
		r.writeColorln(ui, fmt.Sprintf("%s  No steps in workflow%s",
			r.colors.BrightBlack,
			r.colors.Reset))
		r.writeColorln(ui, "")
		return
	}

	// Render all workflow steps
	for i, step := range steps {
		stepLine := fmt.Sprintf("  %s%d.%s %s%s%s",
			r.colors.BrightBlue+r.colors.Bold,
			i+1,
			r.colors.Reset,
			r.colors.BrightGreen+r.colors.Bold,
			step.Description,
			r.colors.Reset)
		r.writeColorln(ui, stepLine)
	}

	r.writeColorln(ui, "")

	// Keybinds rendered elsewhere in workflow mode view
}

// renderGitStatus renders the Git repository status information
