package interactive

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"

	"golang.org/x/term"
)

func (h *KeyHandler) readNextByte(reader *bufio.Reader) (byte, error) {
	if reader != nil {
		return reader.ReadByte()
	}
	var buf [1]byte
	_, err := h.ui.stdin.Read(buf[:])
	return buf[0], err
}

// handleCtrlC handles Ctrl+C key press
func (h *KeyHandler) handleCtrlC(oldState *term.State) {
	if oldState != nil {
		if f, ok := h.ui.stdin.(*os.File); ok {
			if err := h.ui.term.Restore(int(f.Fd()), oldState); err != nil {
				h.ui.writeError("failed to restore terminal state: %v", err)
			}
		}
	}
	h.ui.writeln("\nExiting...")
	os.Exit(0)
}

// restoreTerminalState restores the terminal from raw mode to cooked mode
func (h *KeyHandler) restoreTerminalState(oldState *term.State) {
	if oldState == nil {
		return
	}
	if f, ok := h.ui.stdin.(*os.File); ok {
		if err := h.ui.term.Restore(int(f.Fd()), oldState); err != nil {
			h.ui.writeError("failed to restore terminal state: %v", err)
		}
	}
}

// reenterRawMode re-enters raw mode after being restored
func (h *KeyHandler) reenterRawMode(oldState *term.State) {
	if oldState == nil {
		return
	}
	if f, ok := h.ui.stdin.(*os.File); ok {
		fd := int(f.Fd())
		if _, err := h.ui.term.MakeRaw(fd); err != nil {
			h.ui.writeError("failed to set terminal to raw mode: %v", err)
		}
	}
}

// handleEnter handles Enter key press
func (h *KeyHandler) handleEnter(oldState *term.State) (bool, []string) {
	if !h.ui.state.HasInput() {
		return true, nil
	}

	selectedCmd := h.ui.state.GetSelectedCommand()
	if selectedCmd == nil {
		return true, nil
	}

	// Restore terminal state BEFORE showing Execute message
	h.restoreTerminalState(oldState)

	// Clear screen and show execution message
	clearScreen(h.ui.stdout)
	executeMsg := fmt.Sprintf("%s🚀 %sExecuting:%s %s%s%s\n\n",
		h.ui.colors.BrightGreen,
		h.ui.colors.BrightWhite+h.ui.colors.Bold,
		h.ui.colors.Reset,
		h.ui.colors.BrightCyan+h.ui.colors.Bold,
		selectedCmd.Command,
		h.ui.colors.Reset)
	h.ui.writeColor(executeMsg)

	// Handle placeholders
	args, canceled := h.processCommand(selectedCmd.Command)
	if canceled {
		// Re-enter raw mode before returning to main loop
		h.reenterRawMode(oldState)
		return true, nil
	}
	return false, args
}

// processCommand processes the command with placeholder replacement
func (h *KeyHandler) processCommand(cmdTemplate string) ([]string, bool) {
	placeholders := extractPlaceholders(cmdTemplate)

	if len(placeholders) == 0 {
		// No placeholders - execute immediately
		args := []string{"ggc"}
		args = append(args, strings.Fields(cmdTemplate)...)
		return args, false
	}

	// Interactive input for placeholders
	inputs, canceled := h.interactiveInput(placeholders)
	if canceled {
		h.handleSoftCancel(nil)
		return nil, true
	}

	// Placeholder replacement
	finalCmd := cmdTemplate
	for ph, val := range inputs {
		finalCmd = strings.ReplaceAll(finalCmd, "<"+ph+">", val)
	}

	args := []string{"ggc"}
	args = append(args, strings.Fields(finalCmd)...)
	return args, false
}

// interactiveInput provides real-time interactive input for placeholders
func (h *KeyHandler) interactiveInput(placeholders []string) (map[string]string, bool) {
	inputs := make(map[string]string)

	for i, ph := range placeholders {
		h.ui.write("\n")

		// Show progress and prompt
		if len(placeholders) > 1 {
			h.ui.write("%s[%d/%d]%s ",
				h.ui.colors.BrightBlue+h.ui.colors.Bold,
				i+1, len(placeholders),
				h.ui.colors.Reset)
		}

		h.ui.write("%s? %s%s%s: ",
			h.ui.colors.BrightGreen,
			h.ui.colors.BrightWhite+h.ui.colors.Bold,
			ph,
			h.ui.colors.Reset)

		// Get input with real-time feedback
		value, canceled := h.ui.readPlaceholderInput()
		if canceled {
			return nil, true
		}
		inputs[ph] = value

		// Show confirmation
		h.ui.write("%s✓ %s%s: %s%s%s\n",
			h.ui.colors.BrightGreen,
			h.ui.colors.BrightBlue,
			ph,
			h.ui.colors.BrightYellow+h.ui.colors.Bold,
			value,
			h.ui.colors.Reset)
	}

	return inputs, false
}

// getRealTimeInput gets user input with real-time display using raw terminal mode
func (h *KeyHandler) getRealTimeInput() (string, bool) {
	fd := int(os.Stdin.Fd())
	oldState, err := h.ui.term.MakeRaw(fd)
	if err != nil {
		return h.getLineInput()
	}
	defer func() { _ = h.ui.term.Restore(fd, oldState) }()

	return h.processRealTimeInput()
}

// processRealTimeInput handles the main input processing loop
func (h *KeyHandler) processRealTimeInput() (string, bool) {
	reader := bufio.NewReader(os.Stdin)
	inputRunes := make([]rune, 0, initialInputCapacity)
	cursor := 0

	editor := &realTimeEditor{
		ui:         h.ui,
		inputRunes: &inputRunes,
		cursor:     &cursor,
	}

	for {
		r, _, err := reader.ReadRune()
		if errors.Is(err, ErrWorkflowCanceled) {
			break
		}
		if err != nil {
			continue
		}

		result := editor.handleInput(r, reader)
		if result.done {
			return result.text, false
		}
		if result.canceled {
			return "", true
		}
	}
	return string(inputRunes), false
}

//nolint:revive // Input character handling inherently requires multiple cases
func (h *KeyHandler) handleInputChar(input *strings.Builder, char rune) (done bool, canceled bool) {
	switch char {
	case '\n', '\r':
		if input.Len() > 0 {
			h.ui.write("\r\n")
			return true, false
		}
		h.ui.write(" %s(required)%s", h.ui.colors.BrightRed, h.ui.colors.Reset)
		return false, false
	case '\b', 127:
		if input.Len() == 0 {
			return false, false
		}
		runes := []rune(input.String())
		if len(runes) == 0 {
			return false, false
		}
		// Identify start of previous grapheme-like cluster
		end := len(runes)
		start := end - 1
		for start >= 0 && (isCombining(runes[start]) || isVariationSelector(runes[start])) {
			start--
		}
		if start >= 0 && isRegionalIndicator(runes[start]) {
			if start > 0 && isRegionalIndicator(runes[start-1]) {
				start--
			}
		}
		for {
			if start > 0 && isZWJ(runes[start-1]) {
				start -= 2
				for start >= 0 && (isCombining(runes[start]) || isVariationSelector(runes[start])) {
					start--
				}
				continue
			}
			break
		}
		if start < 0 {
			start = 0
		}
		// Calculate total columns to clear
		cols := 0
		for i := start; i < end; i++ {
			cols += runeDisplayWidth(runes[i])
		}
		// Update input
		input.Reset()
		input.WriteString(string(runes[:start]))
		// Clear terminal cells
		for i := 0; i < cols; i++ {
			h.ui.write("\b \b")
		}
		return false, false
	case 3: // Ctrl+C
		h.ui.write("\r\n%sOperation canceled%s\r\n", h.ui.colors.BrightRed, h.ui.colors.Reset)
		return true, true
	default:
		// Accept all printable characters including multibyte
		if unicode.IsPrint(char) {
			input.WriteRune(char)
			h.ui.write("%s", string(char))
		}
		return false, false
	}
}

// getLineInput provides fallback line-based input when raw mode is not available
func (h *KeyHandler) getLineInput() (string, bool) {
	reader := bufio.NewReader(h.ui.stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", true
		}
		line = strings.TrimSpace(line)
		if line != "" {
			return line, false
		}
		h.ui.write("%s(required)%s ",
			h.ui.colors.BrightRed,
			h.ui.colors.Reset)
	}
}

// addCommandToWorkflow adds a command to the workflow (preserving placeholders)
