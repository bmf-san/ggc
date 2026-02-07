package interactive

import (
	"bufio"
	"errors"
	"io"
	"os"

	"golang.org/x/term"
)

// setupTerminal configures terminal raw mode and returns the old state and error status
func (ui *UI) setupTerminal() (*term.State, bool) {
	var oldState *term.State
	if f, ok := ui.stdin.(*os.File); ok {
		// Check if it's a real terminal (TTY) - skip for non-TTY or when term is nil
		if ui.term == nil || !term.IsTerminal(int(f.Fd())) {
			// Not a real terminal, skip raw mode setup for debugging
			return nil, true
		}

		fd := int(f.Fd())
		var err error
		oldState, err = ui.term.MakeRaw(fd)
		if err != nil {
			ui.writeError("failed to set terminal to raw mode: %v", err)
			return nil, false
		}
	}
	return oldState, true
}

// Run executes the interactive UI
func (ui *UI) Run() []string {
	oldState, reader, isRawMode := ui.initializeTerminal()
	// If terminal initialization failed in non-raw mode, abort cleanly.
	if !isRawMode && reader == nil {
		return nil
	}

	// Set up terminal restoration for raw mode
	if f, ok := ui.stdin.(*os.File); ok && isRawMode {
		fd := int(f.Fd())
		defer func() {
			if err := ui.term.Restore(fd, oldState); err != nil {
				ui.writeError("failed to restore terminal state: %v", err)
			}
		}()
	}

	return ui.runMainLoop(reader, isRawMode, oldState)
}

// initializeTerminal sets up the terminal and returns the old state, reader, and raw mode status
func (ui *UI) initializeTerminal() (*term.State, *bufio.Reader, bool) {
	oldState, ok := ui.setupTerminal()
	if !ok {
		return nil, nil, false
	}

	// Check if we're in raw mode (real terminal) or not
	isRawMode := oldState != nil

	// Set up reader based on mode
	var reader *bufio.Reader
	if !isRawMode {
		reader = bufio.NewReader(ui.stdin)
	}

	return oldState, reader, isRawMode
}

// runMainLoop handles the main input loop
func (ui *UI) runMainLoop(reader *bufio.Reader, isRawMode bool, oldState *term.State) []string {
	if isRawMode {
		if ui.reader == nil {
			ui.reader = bufio.NewReader(ui.stdin)
		}
	} else {
		ui.reader = reader
	}

	for {
		ui.state.UpdateFiltered()
		ui.renderer.Render(ui, ui.state)

		r, err := ui.readNextRune(reader, isRawMode)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			continue // Skip this iteration for other errors
		}

		// Handle key input with rune
		isSingleByte := isRawMode // In raw mode, we read single bytes; in buffered mode, we read full runes
		shouldContinue, result := ui.handler.HandleKey(r, isSingleByte, oldState, reader)
		if !shouldContinue {
			return result
		}
	}
}

// readNextRune reads the next rune from input based on the mode
func (ui *UI) readNextRune(reader *bufio.Reader, isRawMode bool) (rune, error) {
	if isRawMode {
		for {
			// Read single byte directly from stdin in raw mode
			var buf [1]byte
			n, readErr := ui.stdin.Read(buf[:])

			if readErr != nil {
				return 0, readErr
			}

			if n == 0 {
				continue // Try again if no bytes read
			}

			return rune(buf[0]), nil
		}
	}

	// Use buffered reader for non-TTY
	r, _, err := reader.ReadRune()
	return r, err
}
