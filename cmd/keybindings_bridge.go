package cmd

import "github.com/bmf-san/ggc/v7/internal/keybindings"

// DebugKeysCommand exposes the keybindings debug command to the cmd package.
type DebugKeysCommand = keybindings.DebugKeysCommand

var (
	// NewDebugKeysCommand builds a command that emits debug information about keybindings.
	NewDebugKeysCommand = keybindings.NewDebugKeysCommand
)
