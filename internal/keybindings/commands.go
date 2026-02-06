package keybindings

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// ShowKeysCommand displays effective keybindings
type ShowKeysCommand struct {
	resolver *KeyBindingResolver
	platform string
	terminal string
}

// NewShowKeysCommand creates a new show keys command
func NewShowKeysCommand(resolver *KeyBindingResolver) *ShowKeysCommand {
	return &ShowKeysCommand{
		resolver: resolver,
		platform: DetectPlatform(),
		terminal: DetectTerminal(),
	}
}

// Execute runs the show keys command
func (skc *ShowKeysCommand) Execute(profile Profile, context Context, format string) error { //nolint:revive // rich output grouped by sections
	fmt.Printf("ggc Interactive Mode - Effective Keybindings\n")
	fmt.Printf("=============================================\n\n")

	// Get profile info
	prof, exists := skc.resolver.GetProfile(profile)
	if !exists {
		return fmt.Errorf("profile '%s' not found", profile)
	}

	fmt.Printf("Profile: %s", prof.Name)
	if prof.Description != "" {
		fmt.Printf(" (%s)", prof.Description)
	}
	fmt.Printf("\n")

	fmt.Printf("Platform: %s/%s\n", skc.platform, skc.terminal)
	fmt.Printf("Context: %s\n\n", context)

	// Get effective keybindings
	keyMap, err := skc.resolver.Resolve(profile, context)
	if err != nil {
		return fmt.Errorf("failed to resolve keybindings: %w", err)
	}

	// Display keybindings by category
	fmt.Printf("Core Actions:\n")
	fmt.Printf("  Navigation:\n")
	if len(keyMap.MoveUp) > 0 {
		fmt.Printf("    move_up                 %-20s Move up one line\n", FormatKeyStrokesForDisplay(keyMap.MoveUp))
	}
	if len(keyMap.MoveDown) > 0 {
		fmt.Printf("    move_down               %-20s Move down one line\n", FormatKeyStrokesForDisplay(keyMap.MoveDown))
	}
	if len(keyMap.MoveToBeginning) > 0 {
		fmt.Printf("    move_to_beginning       %-20s Move to line beginning\n", FormatKeyStrokesForDisplay(keyMap.MoveToBeginning))
	}
	if len(keyMap.MoveToEnd) > 0 {
		fmt.Printf("    move_to_end             %-20s Move to line end\n", FormatKeyStrokesForDisplay(keyMap.MoveToEnd))
	}

	fmt.Printf("\n  Editing:\n")
	if len(keyMap.DeleteWord) > 0 {
		fmt.Printf("    delete_word             %-20s Delete previous word\n", FormatKeyStrokesForDisplay(keyMap.DeleteWord))
	}
	if len(keyMap.DeleteToEnd) > 0 {
		fmt.Printf("    delete_to_end           %-20s Delete to line end\n", FormatKeyStrokesForDisplay(keyMap.DeleteToEnd))
	}
	if len(keyMap.ClearLine) > 0 {
		fmt.Printf("    clear_line              %-20s Clear entire line\n", FormatKeyStrokesForDisplay(keyMap.ClearLine))
	}

	fmt.Printf("\nQuick Reference:\n")
	fmt.Printf("  quit                    %-20s Exit to shell\n", "Ctrl+C")

	// Show resolution layers
	fmt.Printf("\nResolution Layers Applied:\n")
	fmt.Printf("  1. Base Profile: %s\n", profile)
	fmt.Printf("  2. Platform: %s\n", skc.platform)
	fmt.Printf("  3. Terminal: %s\n", skc.terminal)
	fmt.Printf("  4. User Config: (if configured)\n")

	fmt.Printf("\nTips:\n")
	fmt.Printf("  • Use 'ggc config keybindings --export' to backup your settings\n")
	fmt.Printf("  • Profile switching: set 'interactive.profile' in config\n")

	return nil
}

// DebugKeysCommand captures and displays raw key sequences
type DebugKeysCommand struct {
	capturing  bool
	sequences  [][]byte
	outputFile string
}

// NewDebugKeysCommand creates a new debug keys command
func NewDebugKeysCommand(outputFile string) *DebugKeysCommand {
	return &DebugKeysCommand{
		capturing:  false,
		sequences:  make([][]byte, 0),
		outputFile: outputFile,
	}
}

// StartCapture begins capturing raw key sequences
func (dkc *DebugKeysCommand) StartCapture() {
	dkc.capturing = true
	dkc.sequences = make([][]byte, 0)

	fmt.Printf("=== Debug Keys Mode ===\n")
	fmt.Printf("Raw key sequence capture started.\n")
	fmt.Printf("Press keys to see their sequences.\n")
	fmt.Printf("Press Ctrl+C to stop and view results.\n\n")
}

// CaptureSequence captures a raw key sequence
func (dkc *DebugKeysCommand) CaptureSequence(seq []byte) {
	if !dkc.capturing {
		return
	}

	// Make a copy of the sequence
	captured := make([]byte, len(seq))
	copy(captured, seq)
	dkc.sequences = append(dkc.sequences, captured)

	// Display immediately
	fmt.Printf("Captured: %v (hex: %x) (chars: %q)\n", seq, seq, seq)
}

// StopCapture stops capturing and shows results
func (dkc *DebugKeysCommand) StopCapture() error {
	if !dkc.capturing {
		return nil
	}

	dkc.capturing = false

	fmt.Printf("\n=== Capture Results ===\n")
	fmt.Printf("Total sequences captured: %d\n\n", len(dkc.sequences))

	if len(dkc.sequences) == 0 {
		fmt.Printf("No sequences captured.\n")
		return nil
	}

	// Display all captured sequences
	for i, seq := range dkc.sequences {
		fmt.Printf("%d. %v (hex: %x)\n", i+1, seq, seq)

		// Try to identify common sequences
		if identified := dkc.identifySequence(seq); identified != "" {
			fmt.Printf("   → Identified as: %s\n", identified)
		}

		// Show binding format
		fmt.Printf("   → Config format: \"raw:%x\"\n", seq)
	}

	// Save to file if requested
	if dkc.outputFile != "" {
		if err := dkc.saveToFile(); err != nil {
			return fmt.Errorf("failed to save to file: %w", err)
		}
		fmt.Printf("\nSequences saved to: %s\n", dkc.outputFile)
	}

	fmt.Printf("\nTip: Use the 'raw:' format in your config to bind these sequences.\n")

	return nil
}

func (dkc *DebugKeysCommand) formatKeySequence(seq []byte) string { //nolint:revive // classification supports many key types
	if len(seq) == 0 {
		return "(empty)"
	}

	label := dkc.identifySequence(seq)
	if label == "" {
		if len(seq) == 1 {
			b := seq[0]
			switch {
			case b >= 32 && b <= 126:
				label = string([]byte{b})
			case b >= 1 && b <= 26:
				label = fmt.Sprintf("Ctrl+%c", 'A'+b-1)
			case b == 27:
				label = "Esc"
			default:
				label = fmt.Sprintf("0x%02x", b)
			}
		} else {
			label = fmt.Sprintf("%v", seq)
		}
	}

	hexParts := make([]string, len(seq))
	for i, b := range seq {
		hexParts[i] = fmt.Sprintf("0x%02x", b)
	}

	return fmt.Sprintf("%s (%s)", label, strings.Join(hexParts, " "))
}

// identifySequence tries to identify common key sequences
func (dkc *DebugKeysCommand) identifySequence(seq []byte) string { //nolint:revive // identifies many terminal escape sequences
	if len(seq) == 1 {
		switch seq[0] {
		case 9:
			return "Tab"
		case 13:
			return "Enter"
		case 27:
			return "Esc"
		case 32:
			return "Space"
		}
		if seq[0] >= 1 && seq[0] <= 26 {
			return fmt.Sprintf("Ctrl+%c", 'A'+seq[0]-1)
		}
	}

	if len(seq) == 3 && seq[0] == 27 && seq[1] == 91 {
		switch seq[2] {
		case 65:
			return "↑"
		case 66:
			return "↓"
		case 67:
			return "→"
		case 68:
			return "←"
		}
	}

	// Shift-modified arrow keys (CSI 1;2X sequences)
	if len(seq) == 6 && seq[0] == 27 && seq[1] == 91 && seq[2] == 49 && seq[3] == 59 {
		if seq[4] == 50 {
			switch seq[5] {
			case 65:
				return "Shift+↑"
			case 66:
				return "Shift+↓"
			case 67:
				return "Shift+→"
			case 68:
				return "Shift+←"
			}
		}
	}

	// Function keys
	if len(seq) >= 3 && seq[0] == 27 && seq[1] == 79 {
		switch seq[2] {
		case 80:
			return "F1"
		case 81:
			return "F2"
		case 82:
			return "F3"
		case 83:
			return "F4"
		}
	}

	return ""
}

// saveToFile saves captured sequences to a file
func (dkc *DebugKeysCommand) saveToFile() error {
	var content strings.Builder

	content.WriteString("# Raw Key Sequences Captured by ggc debug-keys\n")
	content.WriteString(fmt.Sprintf("# Captured on: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	content.WriteString(fmt.Sprintf("# Total sequences: %d\n\n", len(dkc.sequences)))

	for i, seq := range dkc.sequences {
		content.WriteString(fmt.Sprintf("# Sequence %d\n", i+1))
		content.WriteString(fmt.Sprintf("# Raw: %v\n", seq))
		content.WriteString(fmt.Sprintf("# Hex: %x\n", seq))
		if identified := dkc.identifySequence(seq); identified != "" {
			content.WriteString(fmt.Sprintf("# Identified: %s\n", identified))
		}
		content.WriteString(fmt.Sprintf("raw:%x\n\n", seq))
	}

	if err := os.WriteFile(dkc.outputFile, []byte(content.String()), 0600); err != nil {
		return err
	}

	fmt.Printf("Saved to %s:\n%s", dkc.outputFile, content.String())

	return nil
}

// IsCapturing returns whether debug capture is active
func (dkc *DebugKeysCommand) IsCapturing() bool {
	return dkc.capturing
}
