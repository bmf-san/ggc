package interactive

import (
	"fmt"
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
	items, picked := newestFirstUniqueDisplay(entries)
	formatter := ui.NewFormatter(h.ui.stdout)
	loop := ui.NewSelectionLoop(formatter, "Select a command from history:", items)
	loop.Display()

	// Re-enter raw mode for the line read so Ctrl+C is delivered as
	// the 0x03 byte and treated as a local cancel. Otherwise cooked
	// mode forwards SIGINT to the interactive REPL's signal handler,
	// which tears down the whole ggc process.
	h.reenterRawMode(oldState)
	line, ok := readSelectorLine(h.ui.stdin, h.ui.stdout)
	h.restoreTerminalState(oldState)
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
	return h.dispatchHistorySelection(input, picked, oldState)
}

// dispatchHistorySelection turns a parsed selection input into either
// (a) a dispatched argv that the REPL should execute, or (b) a cancel
// outcome that returns control to the prompt. Keeping the switch in
// its own function trims runHistorySelector below the cyclomatic
// complexity ceiling without changing observable behavior.
func (h *KeyHandler) dispatchHistorySelection(input ui.SelectionInput, picked []history.Entry, oldState *term.State) (bool, []string) {
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
		return h.dispatchHistoryItem(input.Indices, picked, oldState)
	default:
		h.reenterRawMode(oldState)
		return true, nil
	}
}

// dispatchHistoryItem validates the single-item selection contract and
// either returns the canonical argv for the picked entry or cancels
// with a user-visible reason. Extracted from dispatchHistorySelection
// to keep each function under the cyclomatic budget.
func (h *KeyHandler) dispatchHistoryItem(indices []int, picked []history.Entry, oldState *term.State) (bool, []string) {
	if len(indices) != 1 {
		h.ui.writeln("(history) pick exactly one entry; canceled.")
		h.reenterRawMode(oldState)
		return true, nil
	}
	idx := indices[0]
	if idx < 0 || idx >= len(picked) {
		h.ui.writeError("invalid selection: out of range")
		h.reenterRawMode(oldState)
		return true, nil
	}
	// `picked` is index-aligned with `items` (which is the
	// deduplicated newest-first list), so the user-typed index
	// maps straight onto the displayed row.
	entry := picked[idx]
	return false, entryToArgs(&entry)
}

// newestFirstUniqueDisplay produces the display strings the picker
// shows along with the originating entries, in newest-first order and
// deduplicated by display string. Deduplication keeps the most recent
// occurrence, matching the Ctrl+R reverse-i-search behavior so the two
// views never disagree about "how many times did I run X?". Returning
// the picked entries alongside the display strings keeps the index
// mapping trivial: items[i] and picked[i] reference the same row.
func newestFirstUniqueDisplay(entries []history.Entry) ([]string, []history.Entry) {
	items := make([]string, 0, len(entries))
	picked := make([]history.Entry, 0, len(entries))
	seen := make(map[string]struct{}, len(entries))
	for i := len(entries) - 1; i >= 0; i-- {
		disp := entries[i].Display()
		if disp == "" {
			continue
		}
		if _, dup := seen[disp]; dup {
			continue
		}
		seen[disp] = struct{}{}
		items = append(items, disp)
		picked = append(picked, entries[i])
	}
	return items, picked
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

// readSelectorLine reads one line of selector input one byte at a
// time so we can interpret control characters locally. We deliberately
// avoid cooked-mode line buffering: Ctrl+C in cooked mode is converted
// to SIGINT by the tty driver, and the interactive REPL's signal
// handler treats SIGINT as "quit ggc", which is the wrong outcome
// when the user just wants to back out of the history picker.
//
//   - Enter (\r or \n)   -> return the accumulated line, ok=true
//   - Ctrl+C / Ctrl+G / Esc -> cancel, ok=false
//   - Backspace (DEL or BS) -> erase one byte and redraw
//   - printable ASCII       -> echo and append
//
// Multi-byte UTF-8 is not expected here (the picker accepts digits and
// the literal word "all"), so we keep the reader byte-oriented.
func readSelectorLine(r io.Reader, w io.Writer) (string, bool) {
	buf := make([]byte, 0, 8)
	one := make([]byte, 1)
	for {
		n, err := r.Read(one)
		if err != nil || n == 0 {
			if len(buf) > 0 {
				return string(buf), true
			}
			return "", false
		}
		newBuf, done, accepted := stepSelectorByte(one[0], buf, w)
		buf = newBuf
		if done {
			return string(buf), accepted
		}
	}
}

// stepSelectorByte advances the readSelectorLine state machine by one
// byte. Splitting the per-byte switch out keeps the loop function
// under the cyclomatic complexity budget without losing the inline
// control flow of a single switch. Returns the (possibly modified)
// buffer, a done flag that ends the loop, and the accepted flag that
// distinguishes Enter (true) from cancel (false).
func stepSelectorByte(b byte, buf []byte, w io.Writer) ([]byte, bool, bool) {
	switch b {
	case '\r', '\n':
		_, _ = fmt.Fprint(w, "\r\n")
		return buf, true, true
	case 0x03, 0x07, 0x1b: // Ctrl+C, Ctrl+G, Esc
		_, _ = fmt.Fprint(w, "\r\n")
		return buf[:0], true, false
	case 0x7f, 0x08: // DEL / Backspace
		if len(buf) > 0 {
			buf = buf[:len(buf)-1]
			_, _ = fmt.Fprint(w, "\b \b")
		}
		return buf, false, false
	default:
		if b >= 0x20 && b < 0x7f {
			buf = append(buf, b)
			_, _ = w.Write([]byte{b})
		}
		return buf, false, false
	}
}

// isInteractiveHistoryCommand reports whether the selected command in
// the REPL should be intercepted by the history selector instead of
// dispatched to the regular `ggc history` handler. Sub-forms like
// `history last 20` or `history search foo` keep their normal printing
// behavior since the user has clearly opted into the textual view.
func isInteractiveHistoryCommand(cmd string) bool {
	return strings.TrimSpace(cmd) == "history"
}
