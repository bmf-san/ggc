package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bmf-san/ggc/v8/internal/history"
)

// defaultHistoryShow is the number of entries returned by `ggc history`
// when no count is supplied. Mirrors a typical `history | tail -n 20`.
const defaultHistoryShow = 20

// historyTimeFormat is the timestamp format used in the human-readable
// listing. Local timezone so users can correlate with their own clock,
// fixed-width so columns line up when piped through `column -t`.
const historyTimeFormat = "2006-01-02 15:04:05"

// History displays, searches, or clears the persisted ggc command
// history. It is the user-facing front for the internal/history package.
//
// Subcommand grammar (the first arg that matches wins):
//
//	(no args)                 — show the most recent 20 entries
//	<positive integer N>      — shortcut for `history last N`
//	last <N>                  — show the most recent N entries
//	search <pattern>          — case-insensitive substring filter
//	clear                     — delete every entry
//
// Anything else falls through to a Usage error so typos surface quickly
// instead of being silently treated as "show defaults".
func (c *Cmd) History(args []string) {
	if len(args) == 0 {
		c.showLast(defaultHistoryShow)
		return
	}
	switch args[0] {
	case "last":
		c.handleHistoryLast(args[1:])
	case "search":
		c.handleHistorySearch(args[1:])
	case "clear":
		c.handleHistoryClear()
	default:
		// `ggc history 10` as a shortcut for `ggc history last 10`.
		if n, err := strconv.Atoi(args[0]); err == nil && n > 0 {
			c.showLast(n)
			return
		}
		_, _ = fmt.Fprintln(c.outputWriter, "Usage: ggc history [N | last <N> | search <pattern> | clear]")
	}
}

func (c *Cmd) handleHistoryLast(rest []string) {
	if len(rest) == 0 {
		c.showLast(defaultHistoryShow)
		return
	}
	n, err := strconv.Atoi(rest[0])
	if err != nil || n <= 0 {
		_, _ = fmt.Fprintln(c.outputWriter, "Error: 'last' requires a positive integer")
		return
	}
	c.showLast(n)
}

func (c *Cmd) handleHistorySearch(rest []string) {
	if len(rest) == 0 {
		_, _ = fmt.Fprintln(c.outputWriter, "Usage: ggc history search <pattern>")
		return
	}
	// Join so multi-word patterns like `search feat: add` work
	// without forcing quoting at the shell.
	pattern := strings.Join(rest, " ")
	entries, err := history.Search(pattern)
	if err != nil {
		_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
		return
	}
	c.printEntries(entries)
}

func (c *Cmd) handleHistoryClear() {
	if err := history.Clear(); err != nil {
		_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
		return
	}
	_, _ = fmt.Fprintln(c.outputWriter, "History cleared.")
}

func (c *Cmd) showLast(n int) {
	entries, err := history.ReadLast(n)
	if err != nil {
		_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
		return
	}
	c.printEntries(entries)
}

// printEntries renders entries one per line in chronological order
// using local-timezone timestamps. The format is intentionally simple
// (`YYYY-MM-DD HH:MM:SS<TAB>command`) so it stays greppable from shell
// pipelines: `ggc history | grep commit`.
func (c *Cmd) printEntries(entries []history.Entry) {
	for i := range entries {
		ts := entries[i].Timestamp.Local().Format(historyTimeFormat)
		_, _ = fmt.Fprintf(c.outputWriter, "%s\t%s\n", ts, entries[i].Display())
	}
}
