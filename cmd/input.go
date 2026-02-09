// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"io"

	"github.com/bmf-san/ggc/v7/internal/prompt"
)

// ReadLine prompts for and reads a single line of input.
//
// It returns the entered line and true on success.
// It returns an empty string and false if the prompter is nil, if the input
// is canceled by the user, or if an error occurs while reading input.
// In the case of a non-cancellation error, the error is also written to w
// via WriteError before returning.
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
