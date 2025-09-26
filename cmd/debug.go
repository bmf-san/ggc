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
		"↑/k":      "Navigate up",
		"↓/j":      "Navigate down",
		"Enter":    "Execute selected command",
		"q/Ctrl+C": "Quit interactive mode",
		"r":        "Refresh status",
		"?":        "Show help",
		"Space":    "Toggle item selection (where applicable)",
		"Tab":      "Switch between sections",
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
	// Check if we're in a terminal
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		_, _ = fmt.Fprintln(d.outputWriter, "Error: debug-keys raw mode requires a terminal")
		return
	}

	debugCmd := NewDebugKeysCommand(outputFile)

	// Set up signal handling for graceful exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Save original terminal state
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		_, _ = fmt.Fprintf(d.outputWriter, "Error setting terminal to raw mode: %v\n", err)
		return
	}

	defer func() {
		signal.Reset(os.Interrupt, syscall.SIGTERM)
		_ = term.Restore(int(os.Stdin.Fd()), oldState)
	}()

	// Start capture
	debugCmd.StartCapture()

	// Handle graceful shutdown
	go func() {
		<-sigChan
		_, _ = fmt.Fprintln(d.outputWriter, "\n\nReceived interrupt signal, stopping capture...")
		if err := debugCmd.StopCapture(); err != nil {
			_, _ = fmt.Fprintf(d.outputWriter, "Error stopping capture: %v\n", err)
		}
		_ = term.Restore(int(os.Stdin.Fd()), oldState)
		os.Exit(0)
	}()

	// Read input continuously
	buffer := make([]byte, 64)
	for debugCmd.IsCapturing() {
		n, err := os.Stdin.Read(buffer)
		if err != nil {
			_, _ = fmt.Fprintf(d.outputWriter, "Error reading input: %v\n", err)
			break
		}

		if n > 0 {
			sequence := buffer[:n]

			// Check for Ctrl+C in the sequence
			for _, b := range sequence {
				if b == 3 { // Ctrl+C
					_, _ = fmt.Fprintln(d.outputWriter, "\nCapture stopped by user")
					if err := debugCmd.StopCapture(); err != nil {
						_, _ = fmt.Fprintf(d.outputWriter, "Error stopping capture: %v\n", err)
					}
					return
				}
			}

			debugCmd.CaptureSequence(sequence)
		}
	}
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
