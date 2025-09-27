package cmd

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"
)

// Debugger handles debug operations.
type Debugger struct {
	outputWriter io.Writer
	helper       *Helper
}

// NewDebugger creates a new Debugger instance.
func NewDebugger() *Debugger {
	return &Debugger{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
	}
}

// DebugKeys handles the debug-keys command with subcommand support
func (d *Debugger) DebugKeys(args []string) {
	if len(args) == 0 {
		d.showActiveKeybindings()
		return
	}

	switch args[0] {
	case "raw":
		outputFile := ""
		if len(args) > 1 {
			outputFile = args[1]
		}
		d.captureRawKeySequences(outputFile)
	case "help", "-h", "--help":
		d.showDebugKeysHelp()
	default:
		_, _ = fmt.Fprintf(d.outputWriter, "Unknown subcommand: %s\n", args[0])
		d.showDebugKeysHelp()
	}
}

// showActiveKeybindings displays currently active key bindings
func (d *Debugger) showActiveKeybindings() {
	_, _ = fmt.Fprintln(d.outputWriter, "=== Active Key Bindings ===")
	_, _ = fmt.Fprintln(d.outputWriter, "")

	// Show default keybindings for interactive mode
	defaultBindings := map[string]string{
		"Ctrl+P":        "Move selection up",
		"Ctrl+N":        "Move selection down",
		"Enter":         "Execute selected command",
		"Ctrl+C":        "Quit interactive mode",
		"Ctrl+U":        "Clear input",
		"Ctrl+W":        "Delete previous word",
		"Ctrl+K":        "Delete to end of line",
		"Ctrl+A":        "Move cursor to start",
		"Ctrl+E":        "Move cursor to end",
		"Backspace":     "Delete character before cursor",
		"Alt+Backspace": "Delete previous word (terminal dependent)",
		"←/→":           "Move cursor",
		"Ctrl+←/→":      "Move by word (terminal dependent)",
	}

	_, _ = fmt.Fprintln(d.outputWriter, "Interactive Mode Default Bindings:")
	for key, desc := range defaultBindings {
		_, _ = fmt.Fprintf(d.outputWriter, "  %-12s %s\n", key, desc)
	}

	_, _ = fmt.Fprintln(d.outputWriter, "")
	_, _ = fmt.Fprintln(d.outputWriter, "Custom keybinding configuration:")
	_, _ = fmt.Fprintln(d.outputWriter, "  Use 'ggc debug-keys raw' to capture key sequences for custom bindings")
	_, _ = fmt.Fprintln(d.outputWriter, "  Add them to your ggc config with the 'raw:' prefix")
}

// captureRawKeySequences captures and displays raw key sequences
func (d *Debugger) captureRawKeySequences(outputFile string) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		_, _ = fmt.Fprintln(d.outputWriter, "Error: debug-keys raw mode requires a terminal")
		return
	}

	debugCmd := NewDebugKeysCommand(outputFile)
	oldState, err := d.setupTerminalRawMode()
	if err != nil {
		_, _ = fmt.Fprintf(d.outputWriter, "Error setting terminal to raw mode: %v\n", err)
		return
	}

	defer d.restoreTerminal(oldState)

	debugCmd.StartCapture()
	d.handleSignals(debugCmd, oldState)
	d.processInput(debugCmd)
}

// setupTerminalRawMode configures the terminal for raw input
func (d *Debugger) setupTerminalRawMode() (*term.State, error) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	return term.MakeRaw(int(os.Stdin.Fd()))
}

// restoreTerminal restores the terminal to its original state
func (d *Debugger) restoreTerminal(oldState *term.State) {
	signal.Reset(os.Interrupt, syscall.SIGTERM)
	if err := term.Restore(int(os.Stdin.Fd()), oldState); err != nil {
		_, _ = fmt.Fprintf(d.outputWriter, "Error restoring terminal: %v\n", err)
	}
}

// handleSignals sets up signal handling for graceful shutdown
func (d *Debugger) handleSignals(debugCmd *DebugKeysCommand, oldState *term.State) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		_, _ = fmt.Fprintln(d.outputWriter, "\n\nReceived interrupt signal, stopping capture...")
		if err := debugCmd.StopCapture(); err != nil {
			_, _ = fmt.Fprintf(d.outputWriter, "Error stopping capture: %v\n", err)
		}
		if err := term.Restore(int(os.Stdin.Fd()), oldState); err != nil {
			_, _ = fmt.Fprintf(d.outputWriter, "Error restoring terminal: %v\n", err)
		}
		os.Exit(0)
	}()
}

// processInput reads and processes keyboard input
func (d *Debugger) processInput(debugCmd *DebugKeysCommand) {
	buffer := make([]byte, 64)
	for debugCmd.IsCapturing() {
		n, err := os.Stdin.Read(buffer)
		if err != nil {
			_, _ = fmt.Fprintf(d.outputWriter, "Error reading input: %v\n", err)
			break
		}

		if n > 0 {
			sequence := buffer[:n]
			if d.checkForCtrlC(sequence, debugCmd) {
				return
			}
			debugCmd.CaptureSequence(sequence)
		}
	}
}

// checkForCtrlC checks if Ctrl+C was pressed and stops capture if so
func (d *Debugger) checkForCtrlC(sequence []byte, debugCmd *DebugKeysCommand) bool {
	for _, b := range sequence {
		if b == 3 { // Ctrl+C
			_, _ = fmt.Fprintln(d.outputWriter, "\nCapture stopped by user")
			if err := debugCmd.StopCapture(); err != nil {
				_, _ = fmt.Fprintf(d.outputWriter, "Error stopping capture: %v\n", err)
			}
			return true
		}
	}
	return false
}

// showDebugKeysHelp displays help for the debug-keys command
func (d *Debugger) showDebugKeysHelp() {
	help := `debug-keys - Debug keybinding issues and capture raw key sequences

USAGE:
    ggc debug-keys [SUBCOMMAND] [OPTIONS]

SUBCOMMANDS:
    (none)          Show currently active key bindings
    raw [file]      Capture raw key sequences and optionally save to file
    help            Show this help message

EXAMPLES:
    ggc debug-keys                 # Show active keybindings
    ggc debug-keys raw             # Capture key sequences interactively
    ggc debug-keys raw keys.txt    # Capture and save to keys.txt

DESCRIPTION:
    The debug-keys command helps troubleshoot keybinding issues by:
    1. Showing currently active key bindings
    2. Capturing raw key sequences sent by your terminal
    3. Identifying common key sequences (arrows, function keys, etc.)
    4. Providing the correct format for custom keybinding configuration

    Use 'raw' mode to see exactly what key codes your terminal sends when
    you press specific keys. This is especially useful for:
    - Custom key combinations that don't work as expected
    - Terminal-specific key sequences (tmux, different terminal emulators)
    - Function keys, Alt/Meta combinations, and special keys

    Press Ctrl+C in raw mode to stop capturing and view results.
`
	_, _ = fmt.Fprint(d.outputWriter, help)
}
