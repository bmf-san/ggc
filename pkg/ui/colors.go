// Package ui provides shared terminal rendering utilities for the CLI.
package ui

// ANSIColors defines terminal color escape sequences for both base and bright palettes,
// plus common text attributes. All fields contain raw ANSI escape codes suitable for
// writing directly to an io.Writer.
type ANSIColors struct {
	Black   string
	Red     string
	Green   string
	Yellow  string
	Blue    string
	Magenta string
	Cyan    string
	White   string

	BrightBlack   string
	BrightRed     string
	BrightGreen   string
	BrightYellow  string
	BrightBlue    string
	BrightMagenta string
	BrightCyan    string
	BrightWhite   string

	Bold      string
	Underline string
	Reverse   string
	Reset     string
}

// NewANSIColors returns a palette initialized with the standard ANSI escape codes.
func NewANSIColors() *ANSIColors {
	return &ANSIColors{
		Black:   "\033[30m",
		Red:     "\033[31m",
		Green:   "\033[32m",
		Yellow:  "\033[33m",
		Blue:    "\033[34m",
		Magenta: "\033[35m",
		Cyan:    "\033[36m",
		White:   "\033[37m",

		BrightBlack:   "\033[90m",
		BrightRed:     "\033[91m",
		BrightGreen:   "\033[92m",
		BrightYellow:  "\033[93m",
		BrightBlue:    "\033[94m",
		BrightMagenta: "\033[95m",
		BrightCyan:    "\033[96m",
		BrightWhite:   "\033[97m",

		Bold:      "\033[1m",
		Underline: "\033[4m",
		Reverse:   "\033[7m",
		Reset:     "\033[0m",
	}
}
