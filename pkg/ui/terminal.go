package ui

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

const (
	escDisableWrap = "\x1b[?7l"
	escEnableWrap  = "\x1b[?7h"
	escHideCursor  = "\x1b[?25l"
	escShowCursor  = "\x1b[?25h"
	escClearScreen = "\x1b[2J\x1b[H"
)

// ClearScreen clears the terminal and positions the cursor at the top-left corner.
func ClearScreen(w io.Writer) {
	_, _ = fmt.Fprint(w, escClearScreen)
}

// HideCursor hides the terminal cursor until ShowCursor is called.
func HideCursor(w io.Writer) {
	_, _ = fmt.Fprint(w, escHideCursor)
}

// ShowCursor makes the terminal cursor visible.
func ShowCursor(w io.Writer) {
	_, _ = fmt.Fprint(w, escShowCursor)
}

// DisableWrap turns off automatic line wrapping.
func DisableWrap(w io.Writer) {
	_, _ = fmt.Fprint(w, escDisableWrap)
}

// EnableWrap reenables automatic line wrapping.
func EnableWrap(w io.Writer) {
	_, _ = fmt.Fprint(w, escEnableWrap)
}

// Dimensions attempts to determine the terminal size for the provided writer. If the
// writer is not backed by an *os.File or the lookup fails, it returns a safe default.
func Dimensions(w io.Writer, fallbackWidth, fallbackHeight int) (width, height int) {
	if f, ok := w.(*os.File); ok {
		if fw, fh, err := term.GetSize(int(f.Fd())); err == nil && fw > 0 && fh > 0 {
			return fw, fh
		}
	}
	if fallbackWidth <= 0 {
		fallbackWidth = 80
	}
	if fallbackHeight <= 0 {
		fallbackHeight = 24
	}
	return fallbackWidth, fallbackHeight
}
