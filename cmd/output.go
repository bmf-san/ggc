// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
)

// WriteError writes an error message to the writer
func WriteError(w io.Writer, err error) {
	_, _ = fmt.Fprintf(w, "Error: %v\n", err)
}

// WriteErrorf writes a formatted error message to the writer
func WriteErrorf(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, "Error: "+format+"\n", args...)
}

// WriteLine writes a line to the writer
func WriteLine(w io.Writer, msg string) {
	_, _ = fmt.Fprintln(w, msg)
}

// WriteLinef writes a formatted line to the writer
func WriteLinef(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format+"\n", args...)
}
