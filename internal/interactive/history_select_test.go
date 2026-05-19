package interactive

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v8/internal/history"
	"github.com/bmf-san/ggc/v8/internal/testutil"
)

// fakeHistoryReader lets each test seed the picker with deterministic
// entries (and optionally an error) without touching the on-disk store
// or relying on environment isolation.
type fakeHistoryReader struct {
	entries []history.Entry
	err     error
}

func (f fakeHistoryReader) ReadLast(_ int) ([]history.Entry, error) {
	if f.err != nil {
		return nil, f.err
	}
	// Defensive copy: the selector reverses the slice, and we don't
	// want a test mutation to leak into subsequent assertions.
	out := make([]history.Entry, len(f.entries))
	copy(out, f.entries)
	return out, nil
}

// withHistoryReader swaps the package-level reader for the duration of
// a single test and restores it via t.Cleanup so parallel-safe tests
// don't bleed state.
func withHistoryReader(t *testing.T, r historyReader) {
	t.Helper()
	prev := defaultHistoryReader
	defaultHistoryReader = r
	t.Cleanup(func() { defaultHistoryReader = prev })
}

func newSelectorHandler(stdin string) (*KeyHandler, *bytes.Buffer, *bytes.Buffer) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	ui := &UI{
		stdin:     strings.NewReader(stdin),
		stdout:    stdout,
		stderr:    stderr,
		colors:    NewANSIColors(),
		term:      &mockTerminal{shouldFailRaw: true},
		gitClient: testutil.NewMockGitClient(),
	}
	h := &KeyHandler{ui: ui}
	ui.handler = h
	return h, stdout, stderr
}

func TestIsInteractiveHistoryCommand(t *testing.T) {
	cases := map[string]bool{
		"history":           true,
		" history ":         true,
		"history last 5":    false,
		"history search x":  false,
		"history-something": false,
		"":                  false,
	}
	for in, want := range cases {
		if got := isInteractiveHistoryCommand(in); got != want {
			t.Errorf("isInteractiveHistoryCommand(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestRunHistorySelector_EmptyShowsNotice(t *testing.T) {
	withHistoryReader(t, fakeHistoryReader{})
	h, stdout, _ := newSelectorHandler("")

	cont, args := h.runHistorySelector(nil)
	if !cont {
		t.Fatalf("expected REPL to keep looping when history is empty, got cont=%v args=%v", cont, args)
	}
	if args != nil {
		t.Fatalf("expected nil args, got %v", args)
	}
	if !strings.Contains(stdout.String(), "No history entries yet") {
		t.Errorf("expected empty-history notice, got %q", stdout.String())
	}
}

func TestRunHistorySelector_ReadErrorReportedAndContinues(t *testing.T) {
	withHistoryReader(t, fakeHistoryReader{err: errors.New("boom")})
	h, _, stderr := newSelectorHandler("")

	cont, args := h.runHistorySelector(nil)
	if !cont || args != nil {
		t.Fatalf("expected (true, nil) on read error, got (%v, %v)", cont, args)
	}
	if !strings.Contains(stderr.String(), "failed to read history") {
		t.Errorf("expected read error in stderr, got %q", stderr.String())
	}
}

func TestRunHistorySelector_PicksNewestFirst(t *testing.T) {
	// ReadLast returns oldest-first; the picker must reverse so the
	// most recent command is at position 1.
	withHistoryReader(t, fakeHistoryReader{entries: []history.Entry{
		{Command: "status", Raw: "status"},
		{Command: "commit", Args: []string{"tmp"}, Raw: "commit tmp"},
		{Command: "push", Raw: "push"},
	}})
	h, stdout, _ := newSelectorHandler("1\n")

	cont, args := h.runHistorySelector(nil)
	if cont {
		t.Fatalf("expected selector to dispatch args, got cont=true")
	}
	wantArgs := []string{"ggc", "push"}
	if !equalStrings(args, wantArgs) {
		t.Errorf("args = %v, want %v", args, wantArgs)
	}
	out := stdout.String()
	// Position 1 must be the newest entry (`push`).
	pushIdx := strings.Index(out, "push")
	statusIdx := strings.Index(out, "status")
	if pushIdx < 0 || statusIdx < 0 || pushIdx > statusIdx {
		t.Errorf("expected push to appear before status in picker output, got %q", out)
	}
}

func TestRunHistorySelector_PreservesCanonicalArgs(t *testing.T) {
	// Raw includes an alias spelling, but we replay through the
	// canonical command so the executor gets a deterministic argv.
	withHistoryReader(t, fakeHistoryReader{entries: []history.Entry{
		{Command: "checkout", Args: []string{"feat/x"}, Raw: "co feat/x"},
	}})
	h, _, _ := newSelectorHandler("1\n")

	_, args := h.runHistorySelector(nil)
	want := []string{"ggc", "checkout", "feat/x"}
	if !equalStrings(args, want) {
		t.Errorf("args = %v, want %v", args, want)
	}
}

func TestRunHistorySelector_CancelOnBlankInput(t *testing.T) {
	withHistoryReader(t, fakeHistoryReader{entries: []history.Entry{
		{Command: "status", Raw: "status"},
	}})
	h, _, _ := newSelectorHandler("\n")

	cont, args := h.runHistorySelector(nil)
	if !cont || args != nil {
		t.Fatalf("blank input should cancel; got (%v, %v)", cont, args)
	}
}

func TestRunHistorySelector_RejectsAll(t *testing.T) {
	withHistoryReader(t, fakeHistoryReader{entries: []history.Entry{
		{Command: "a", Raw: "a"}, {Command: "b", Raw: "b"},
	}})
	h, stdout, _ := newSelectorHandler("all\n")

	cont, args := h.runHistorySelector(nil)
	if !cont || args != nil {
		t.Fatalf("'all' should cancel; got (%v, %v)", cont, args)
	}
	if !strings.Contains(stdout.String(), "'all' is not supported") {
		t.Errorf("expected 'all' rejection notice, got %q", stdout.String())
	}
}

func TestRunHistorySelector_RejectsMultiPick(t *testing.T) {
	withHistoryReader(t, fakeHistoryReader{entries: []history.Entry{
		{Command: "a", Raw: "a"}, {Command: "b", Raw: "b"},
	}})
	h, stdout, _ := newSelectorHandler("1 2\n")

	cont, args := h.runHistorySelector(nil)
	if !cont || args != nil {
		t.Fatalf("multi-pick should cancel; got (%v, %v)", cont, args)
	}
	if !strings.Contains(stdout.String(), "pick exactly one") {
		t.Errorf("expected multi-pick rejection notice, got %q", stdout.String())
	}
}

func TestRunHistorySelector_InvalidIndexReportsError(t *testing.T) {
	withHistoryReader(t, fakeHistoryReader{entries: []history.Entry{
		{Command: "a", Raw: "a"},
	}})
	h, _, stderr := newSelectorHandler("99\n")

	cont, args := h.runHistorySelector(nil)
	if !cont || args != nil {
		t.Fatalf("invalid index should cancel; got (%v, %v)", cont, args)
	}
	if !strings.Contains(stderr.String(), "invalid selection") {
		t.Errorf("expected invalid-selection error, got %q", stderr.String())
	}
}

func TestEntryToArgs_NoArgs(t *testing.T) {
	got := entryToArgs(&history.Entry{Command: "status"})
	want := []string{"ggc", "status"}
	if !equalStrings(got, want) {
		t.Errorf("entryToArgs = %v, want %v", got, want)
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
