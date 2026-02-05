package ui

import (
	"fmt"
	"io"
)

// Formatter provides common output formatting utilities for CLI commands.
// It encapsulates an io.Writer and provides consistent formatting methods.
type Formatter struct {
	w      io.Writer
	colors *ANSIColors
}

// NewFormatter creates a new Formatter with the given writer.
func NewFormatter(w io.Writer) *Formatter {
	return &Formatter{
		w:      w,
		colors: NewANSIColors(),
	}
}

// Error prints an error message in the standard format: "Error: <message>\n"
func (f *Formatter) Error(err error) {
	_, _ = fmt.Fprintf(f.w, "Error: %v\n", err)
}

// Errorf prints a formatted error message.
func (f *Formatter) Errorf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(f.w, "Error: "+format+"\n", args...)
}

// Print prints a message without formatting.
func (f *Formatter) Print(msg string) {
	_, _ = fmt.Fprint(f.w, msg)
}

// Println prints a message with a newline.
func (f *Formatter) Println(msg string) {
	_, _ = fmt.Fprintln(f.w, msg)
}

// Printf prints a formatted message.
func (f *Formatter) Printf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(f.w, format, args...)
}

// Header prints a header in cyan bold.
func (f *Formatter) Header(msg string) {
	_, _ = fmt.Fprintf(f.w, "%s%s%s%s\n", f.colors.Bold, f.colors.Cyan, msg, f.colors.Reset)
}

// Success prints a success message in green.
func (f *Formatter) Success(msg string) {
	_, _ = fmt.Fprintf(f.w, "%s%s%s\n", f.colors.Green, msg, f.colors.Reset)
}

// Warning prints a warning message in yellow.
func (f *Formatter) Warning(msg string) {
	_, _ = fmt.Fprintf(f.w, "%s%s%s\n", f.colors.Yellow, msg, f.colors.Reset)
}

// ErrorHighlight prints an error message in red bold.
func (f *Formatter) ErrorHighlight(msg string) {
	_, _ = fmt.Fprintf(f.w, "%s%s%s%s\n", f.colors.Bold, f.colors.Red, msg, f.colors.Reset)
}

// NumberedItem prints a numbered list item with highlighted number.
func (f *Formatter) NumberedItem(index int, item string) {
	_, _ = fmt.Fprintf(f.w, "  [%s%s%d%s] %s\n",
		f.colors.Bold, f.colors.Yellow, index, f.colors.Reset, item)
}

// Prompt prints a prompt indicator.
func (f *Formatter) Prompt() {
	_, _ = fmt.Fprint(f.w, "> ")
}

// Colors returns the ANSI color palette for custom formatting.
func (f *Formatter) Colors() *ANSIColors {
	return f.colors
}

// Writer returns the underlying io.Writer.
func (f *Formatter) Writer() io.Writer {
	return f.w
}
