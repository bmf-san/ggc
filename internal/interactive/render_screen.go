// Package interactive houses interactive UI types and helpers shared across the application.
package interactive

import (
	"fmt"
	"strings"
	"unicode/utf8"

	kb "github.com/bmf-san/ggc/v7/internal/keybindings"
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
func (r *Renderer) renderSearchPrompt(ui *UI, state *UIState) {
	inputWithCursor := r.formatInputWithCursor(state)

	searchPrompt := fmt.Sprintf("%s┌─ %sSearch:%s %s",
		r.colors.BrightBlue,
		r.colors.BrightGreen+r.colors.Bold,
		r.colors.Reset,
		inputWithCursor)
	r.writeColorln(ui, searchPrompt)

	// Results separator
	if state.input != "" {
		separator := fmt.Sprintf("%s└─ %sResults:%s",
			r.colors.BrightBlue,
			r.colors.BrightMagenta+r.colors.Bold,
			r.colors.Reset)
		r.writeColorln(ui, separator)
	}
	r.writeEmptyLine()
}

func (r *Renderer) saveCursorAtSearchPrompt(state *UIState) func() {
	linesUp := 2
	if state.input != "" {
		linesUp++
	}
	_, _ = fmt.Fprintf(r.writer, "\x1b[%dA", linesUp)
	const prefix = "┌─ Search: "
	// Compute display width (columns) of the prefix using runeDisplayWidth
	prefixCols := 0
	for _, pr := range prefix {
		prefixCols += runeDisplayWidth(pr)
	}
	// Compute display width up to the logical cursor position
	runes := []rune(state.input)
	cursorPos := state.cursorPos
	if cursorPos > len(runes) {
		cursorPos = len(runes)
	}
	cursorWidth := 0
	for _, rr := range runes[:cursorPos] {
		cursorWidth += runeDisplayWidth(rr)
	}
	column := prefixCols + cursorWidth + 1
	if column < 1 {
		column = 1
	}
	_, _ = fmt.Fprintf(r.writer, "\x1b[%dG", column)
	_, _ = fmt.Fprint(r.writer, "\x1b[s")
	_, _ = fmt.Fprintf(r.writer, "\x1b[%dB", linesUp)
	return func() {
		_, _ = fmt.Fprint(r.writer, "\x1b[u")
	}
}

// formatInputWithCursor formats the input string with cursor position
func (r *Renderer) formatInputWithCursor(state *UIState) string {
	if state.input == "" {
		return fmt.Sprintf("%s█%s", r.colors.BrightWhite+r.colors.Bold, r.colors.Reset)
	}

	inputRunes := []rune(state.input)
	beforeCursor := string(inputRunes[:state.cursorPos])
	afterCursor := string(inputRunes[state.cursorPos:])
	cursor := "│"
	if state.cursorPos >= utf8.RuneCountInString(state.input) {
		cursor = "█"
	}

	return fmt.Sprintf("%s%s%s%s%s%s%s",
		r.colors.BrightYellow,
		beforeCursor,
		r.colors.BrightWhite+r.colors.Bold,
		cursor,
		r.colors.Reset+r.colors.BrightYellow,
		afterCursor,
		r.colors.Reset)
}

// renderEmptyState renders the empty input state
func (r *Renderer) renderEmptyState(ui *UI) {
	r.writeColorln(ui, fmt.Sprintf("%s💭 %sStart typing to search commands...%s",
		r.colors.BrightBlue, r.colors.BrightBlack, r.colors.Reset))
}

func (r *Renderer) buildSearchKeybindEntries(ui *UI) []keybindHelpEntry {
	entries := []keybindHelpEntry{
		{key: "←/→", desc: "Move cursor"},
		{key: "Ctrl+←/→", desc: "Move by word"},
		{key: "Option+←/→", desc: "Move by word (macOS)"},
	}
	// Future: extend this helper for additional contexts such as workflow views.

	var km *kb.KeyBindingMap
	if ui != nil && ui.handler != nil {
		km = ui.handler.GetCurrentKeyMap()
	}
	if km == nil {
		km = kb.DefaultKeyBindingMap()
	}

	defaultMap := kb.DefaultKeyBindingMap()

	appendDynamic := func(primary []kb.KeyStroke, fallback []kb.KeyStroke, desc string) {
		keys := primary
		if len(keys) == 0 {
			keys = fallback
		}
		if len(keys) == 0 {
			return
		}
		formatted := kb.FormatKeyStrokesForDisplay(keys)
		if formatted == "" || formatted == "none" {
			return
		}
		entries = append(entries, keybindHelpEntry{key: formatted, desc: desc})
	}

	appendDynamic(km.MoveUp, defaultMap.MoveUp, "Navigate up")
	appendDynamic(km.MoveDown, defaultMap.MoveDown, "Navigate down")
	appendDynamic(km.ClearLine, defaultMap.ClearLine, "Clear all input")
	appendDynamic(km.DeleteWord, defaultMap.DeleteWord, "Delete word")
	appendDynamic(km.DeleteToEnd, defaultMap.DeleteToEnd, "Delete to end")
	appendDynamic(km.MoveToBeginning, defaultMap.MoveToBeginning, "Move to beginning")
	appendDynamic(km.MoveToEnd, defaultMap.MoveToEnd, "Move to end")

	entries = append(entries, keybindHelpEntry{key: "Backspace", desc: "Delete character"})
	entries = append(entries, keybindHelpEntry{key: "Enter", desc: "Execute selected command"})

	appendDynamic(km.AddToWorkflow, defaultMap.AddToWorkflow, "Add to workflow")
	appendDynamic(km.ToggleWorkflowView, defaultMap.ToggleWorkflowView, "Toggle workflow view")

	entries = append(entries, keybindHelpEntry{key: "Ctrl+c", desc: "Quit"})

	return entries
}

func (r *Renderer) renderKeybindEntries(ui *UI, entries []keybindHelpEntry) {
	if len(entries) == 0 {
		return
	}

	r.writeColorln(ui, fmt.Sprintf("%s⌨️  %sAvailable keybinds:%s",
		r.colors.BrightBlue, r.colors.BrightWhite+r.colors.Bold, r.colors.Reset))

	for _, entry := range entries {
		r.writeColorln(ui, fmt.Sprintf("   %s%s%s  %s%s%s",
			r.colors.BrightGreen+r.colors.Bold,
			entry.key,
			r.colors.Reset,
			r.colors.BrightBlack,
			entry.desc,
			r.colors.Reset))
	}
}

// renderNoMatches renders the no matches found state with keybind help
func (r *Renderer) renderNoMatches(ui *UI, state *UIState) {
	// No matches message
	r.writeColorln(ui, fmt.Sprintf("%s🔍 %sNo commands found for '%s%s%s'%s",
		r.colors.BrightYellow,
		r.colors.BrightWhite,
		r.colors.BrightYellow+r.colors.Bold,
		state.input,
		r.colors.Reset+r.colors.BrightWhite,
		r.colors.Reset))

	r.writeEmptyLine()
	r.renderKeybindEntries(ui, r.buildSearchKeybindEntries(ui))
}

// renderSearchKeybinds renders keybinds available in search UI
func (r *Renderer) renderSearchKeybinds(ui *UI) {
	r.renderKeybindEntries(ui, r.buildSearchKeybindEntries(ui))
}

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
func (r *Renderer) renderGitStatus(ui *UI, status *GitStatus) {
	var parts []string

	// Branch name
	branchPart := fmt.Sprintf("%s📍 %s%s%s",
		r.colors.BrightBlue,
		r.colors.BrightWhite+r.colors.Bold,
		status.Branch,
		r.colors.Reset)
	parts = append(parts, branchPart)

	// Working directory status
	if status.HasChanges {
		var statusParts []string
		if status.Modified > 0 {
			statusParts = append(statusParts, fmt.Sprintf("%d modified", status.Modified))
		}
		if status.Staged > 0 {
			statusParts = append(statusParts, fmt.Sprintf("%d staged", status.Staged))
		}

		workingPart := fmt.Sprintf("%s📝 %s%s%s",
			r.colors.BrightYellow,
			r.colors.BrightWhite+r.colors.Bold,
			strings.Join(statusParts, ", "),
			r.colors.Reset)
		parts = append(parts, workingPart)
	}

	// Remote tracking status
	if status.Ahead > 0 || status.Behind > 0 {
		var remoteParts []string
		if status.Ahead > 0 {
			remoteParts = append(remoteParts, fmt.Sprintf("↑%d", status.Ahead))
		}
		if status.Behind > 0 {
			remoteParts = append(remoteParts, fmt.Sprintf("↓%d", status.Behind))
		}

		remotePart := fmt.Sprintf("%s%s%s",
			r.colors.BrightMagenta+r.colors.Bold,
			strings.Join(remoteParts, " "),
			r.colors.Reset)
		parts = append(parts, remotePart)
	}

	// Render the status line
	statusLine := strings.Join(parts, "  ")
	r.writeColorln(ui, statusLine)
}
