// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"io"

	"github.com/bmf-san/ggc/v7/internal/prompt"
)

// ReadLine reads a line from the prompter
func ReadLine(p prompt.Prompter, w io.Writer, promptText string) (string, bool) {
	if p == nil {
		return "", false
	}
	line, canceled, err := p.Input(promptText)
	if canceled {
		return "", false
	}
	if err != nil {
		WriteError(w, err)
		return "", false
	}
	return line, true
}
