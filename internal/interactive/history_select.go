package interactive

import (
	"bufio"
	"io"
	"strings"

	"golang.org/x/term"

	"github.com/bmf-san/ggc/v8/internal/history"
	"github.com/bmf-san/ggc/v8/internal/ui"
)

// historySelectorLimit caps how many recent entries the in-REPL history
// selector lists. Beyond this the picker stops being useful as a quick
// recall and `ggc history search <pat>` is the better tool.
const historySelectorLimit = 30

// historyReader is the slice of internal/history exercised by the
// interactive history selector. It is a tiny interface so tests can
// inject deterministic entries without touching the global store.
type historyReader interface {
	ReadLast(n int) ([]history.Entry, error)
}

// defaultHistoryReader adapts the package-level history API to the
// historyReader interface used by the selector. Returning a function
// value keeps the indirection cheap and lets tests swap it out.
var defaultHistoryReader historyReader = historyReaderFunc(history.ReadLast)

type historyReaderFunc func(int) ([]history.Entry, error)

func (f historyReaderFunc) ReadLast(n int) ([]history.Entry, error) { return f(n) }

// runHistorySelector replaces the normal `ggc history` execution path
// in the interactive REPL with a numbered picker over recent commands.
// Selecting an entry returns its args ready to feed back through the
// router (canonical command + args, bypassing reparse), so aliases that
// were originally resolved keep resolving the same way.
func (h *KeyHandler) runHistorySelector(oldState *term.State) (bool, []string) {
	entries, err := defaultHistoryReader.ReadLast(historySelectorLimit)
	if err != nil {
		h.ui.writeError("failed to read history: %v", err)
		h.reenterRawMode(oldState)
		return true, nil
	}
	if len(entries) == 0 {
		h.ui.writeln("\nNo history entries yet. Run a few ggc commands first.\n")
		h.reenterRawMode(oldState)
		return true, nil
	}

	// Render numbered list against the cooked terminal so the user
	// can type a multi-digit index without raw-mode key handling.
	items := newestFirstDisplay(entries)
	formatter := ui.NewFormatter(h.ui.stdout)
	loop := ui.NewSelectionLoop(formatter, "Select a command from history:", items)
	loop.Display()

	line, ok := readSelectorLine(h.ui.stdin)
	if !ok {
		h.reenterRawMode(oldState)
		return true, nil
	}

	input, invalid := loop.ParseInput(line)
	// ParseSelectionInput returns SelectionInput{} (zero value, which
	// is SelectionCanceled) when a field fails to parse. Surface that
	// to the user before the regular cancel path so a typo doesn't
	// silently look like a deliberate Esc.
	if invalid != "" {
		h.ui.writeError("invalid selection: %q", invalid)
		h.reenterRawMode(oldState)
		return true, nil
	}
	switch input.Result {
	case ui.SelectionCanceled, ui.SelectionNone:
		h.reenterRawMode(oldState)
		return true, nil
	case ui.SelectionAll:
		// "all" is meaningless here; treat as cancel rather than
		// flooding the executor with N commands.
		h.ui.writeln("(history) 'all' is not supported; canceled.")
		h.reenterRawMode(oldState)
		return true, nil
	case ui.SelectionItems:
		if len(input.Indices) != 1 {
			h.ui.writeln("(history) pick exactly one entry; canceled.")
			h.reenterRawMode(oldState)
			return true, nil
		}
		// Selection list is newest-first, so the index maps back
		// onto the reversed slice we just rendered.
		picked := entries[len(entries)-1-input.Indices[0]]
		return false, entryToArgs(&picked)
	default:
		h.reenterRawMode(oldState)
		return true, nil
	}
}

// newestFirstDisplay produces the display strings the picker shows,
// reversing the chronological order returned by ReadLast so the most
// recent command sits at position 1 — matching how shells render
// `history | tail` style output.
func newestFirstDisplay(entries []history.Entry) []string {
	out := make([]string, len(entries))
	for i := range entries {
		src := entries[len(entries)-1-i]
		out[i] = src.Display()
	}
	return out
}

// entryToArgs reconstructs the argv used to dispatch a history entry
// back through the router. We deliberately use the canonical command +
// args rather than re-tokenizing Raw, so quoted arguments survive a
// round-trip without depending on the shlex tokenizer living in cmd/.
func entryToArgs(e *history.Entry) []string {
	args := make([]string, 0, len(e.Args)+2)
	args = append(args, "ggc", e.Command)
	args = append(args, e.Args...)
	return args
}

// readSelectorLine reads a single line of input from the (cooked-mode)
// terminal. Returns ok=false on EOF so the caller can fall back to
// cancel semantics instead of blocking the REPL.
func readSelectorLine(r io.Reader) (string, bool) {
	br := bufio.NewReader(r)
	line, err := br.ReadString('\n')
	line = strings.TrimRight(line, "\r\n")
	if err != nil && line == "" {
		return "", false
	}
	return line, true
}

// isInteractiveHistoryCommand reports whether the selected command in
// the REPL should be intercepted by the history selector instead of
// dispatched to the regular `ggc history` handler. Sub-forms like
// `history last 20` or `history search foo` keep their normal printing
// behavior since the user has clearly opted into the textual view.
func isInteractiveHistoryCommand(cmd string) bool {
	return strings.TrimSpace(cmd) == "history"
}
