package interactive

import (
	"fmt"
	"unicode/utf8"

	kb "github.com/bmf-san/ggc/v7/internal/keybindings"
)

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
