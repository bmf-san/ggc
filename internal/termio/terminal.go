// Package termio provides small terminal utilities shared across the interactive UI.
package termio

import "golang.org/x/term"

// Terminal abstracts terminal raw mode operations so callers can swap implementations in tests.
type Terminal interface {
	MakeRaw(fd int) (*term.State, error)
	Restore(fd int, state *term.State) error
}

// DefaultTerminal uses golang.org/x/term to manage terminal state.
type DefaultTerminal struct{}

// MakeRaw switches the terminal into raw mode.
func (DefaultTerminal) MakeRaw(fd int) (*term.State, error) {
	return term.MakeRaw(fd)
}

// Restore returns the terminal to its previous state.
func (DefaultTerminal) Restore(fd int, state *term.State) error {
	return term.Restore(fd, state)
}

var pendingInputHook = pendingInput

// PendingInput reports the number of immediately readable bytes for the given descriptor.
func PendingInput(fd uintptr) (int, error) {
	return pendingInputHook(fd)
}

// SetPendingInputFunc overrides the pending-input probe; the returned closure restores the default implementation.
func SetPendingInputFunc(fn func(uintptr) (int, error)) func() {
	prev := pendingInputHook
	pendingInputHook = fn
	return func() { pendingInputHook = prev }
}
