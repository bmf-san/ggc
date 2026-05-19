package history

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// newTestStore returns a Store rooted at t.TempDir(). Tests should not
// share state across runs, so we never touch the package-level Default.
func newTestStore(t *testing.T) *Store {
	t.Helper()
	return &Store{Path: filepath.Join(t.TempDir(), "history.jsonl")}
}

func TestStore_AppendAndRead(t *testing.T) {
	s := newTestStore(t)

	if err := s.Append("status", nil, "status"); err != nil {
		t.Fatalf("append 1: %v", err)
	}
	if err := s.Append("commit", []string{"-m", "msg"}, `commit -m "msg"`); err != nil {
		t.Fatalf("append 2: %v", err)
	}

	all, err := s.ReadAll()
	if err != nil {
		t.Fatalf("read all: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("want 2 entries, got %d", len(all))
	}
	if all[0].Command != "status" || all[1].Command != "commit" {
		t.Fatalf("unexpected order/contents: %+v", all)
	}
	if all[1].Raw == "" || all[1].Args[0] != "-m" {
		t.Fatalf("entry 2 lost fields: %+v", all[1])
	}

	if _, err := os.Stat(s.Path); err != nil {
		t.Fatalf("history file missing: %v", err)
	}
}

func TestStore_AppendEmptyCommandIsNoop(t *testing.T) {
	s := newTestStore(t)
	if err := s.Append("", nil, ""); err != nil {
		t.Fatalf("append empty: %v", err)
	}
	if err := s.Append("   ", []string{"a"}, "   a"); err != nil {
		t.Fatalf("append whitespace: %v", err)
	}
	all, err := s.ReadAll()
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(all) != 0 {
		t.Fatalf("want 0 entries, got %d", len(all))
	}
}

func TestStore_DisabledSkipsWrite(t *testing.T) {
	s := newTestStore(t)
	s.Disabled = true
	if err := s.Append("status", nil, "status"); err != nil {
		t.Fatalf("append: %v", err)
	}
	if _, err := os.Stat(s.Path); !os.IsNotExist(err) {
		t.Fatalf("disabled store should not create file, stat err = %v", err)
	}
}

func TestEnvDisable(t *testing.T) {
	cases := map[string]bool{
		"":      false,
		"0":     false,
		"false": false,
		"no":    false,
		"off":   false,
		"OFF":   false,
		"1":     true,
		"true":  true,
		"yes":   true,
		"on":    true,
	}
	for in, want := range cases {
		if got := envTrue(in); got != want {
			t.Errorf("envTrue(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestStore_ReadLastBounds(t *testing.T) {
	s := newTestStore(t)
	for _, c := range []string{"a", "b", "c", "d"} {
		if err := s.Append(c, nil, c); err != nil {
			t.Fatalf("append %s: %v", c, err)
		}
	}

	cases := []struct {
		n    int
		want int
	}{
		{n: -1, want: 0},
		{n: 0, want: 0},
		{n: 1, want: 1},
		{n: 3, want: 3},
		{n: 4, want: 4},
		{n: 10, want: 4},
	}
	for _, tc := range cases {
		got, err := s.ReadLast(tc.n)
		if err != nil {
			t.Fatalf("ReadLast(%d): %v", tc.n, err)
		}
		if len(got) != tc.want {
			t.Errorf("ReadLast(%d) = %d entries, want %d", tc.n, len(got), tc.want)
		}
	}
}

func TestStore_SearchCaseAndScope(t *testing.T) {
	s := newTestStore(t)
	_ = s.Append("commit", []string{"-m", "feat: add"}, `commit -m "feat: add"`)
	_ = s.Append("checkout", []string{"main"}, "co main") // alias raw
	_ = s.Append("status", nil, "status")

	// Empty pattern matches everything.
	all, err := s.Search("")
	if err != nil {
		t.Fatalf("search empty: %v", err)
	}
	if len(all) != 3 {
		t.Errorf("empty pattern -> %d, want 3", len(all))
	}

	// Case-insensitive.
	got, err := s.Search("COMMIT")
	if err != nil {
		t.Fatalf("search COMMIT: %v", err)
	}
	if len(got) != 1 || got[0].Command != "commit" {
		t.Errorf("COMMIT match = %+v", got)
	}

	// Matches against raw (alias) too.
	got, err = s.Search("co main")
	if err != nil {
		t.Fatalf("search alias: %v", err)
	}
	if len(got) != 1 || got[0].Command != "checkout" {
		t.Errorf("alias match = %+v", got)
	}

	// Matches against canonical even when raw is something else.
	got, err = s.Search("checkout")
	if err != nil {
		t.Fatalf("search canonical: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("canonical match = %+v", got)
	}

	// No match.
	got, err = s.Search("nonexistent")
	if err != nil {
		t.Fatalf("search miss: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("miss = %+v", got)
	}
}

func TestStore_TrimEnforcesCap(t *testing.T) {
	s := newTestStore(t)
	s.MaxEntries = 3

	for i := 0; i < 7; i++ {
		if err := s.Append("c", []string{}, strings.Repeat("x", i+1)); err != nil {
			t.Fatalf("append %d: %v", i, err)
		}
	}

	all, err := s.ReadAll()
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("trim should keep 3, got %d", len(all))
	}
	// Newest 3 are kept: lengths 5, 6, 7.
	if len(all[0].Raw) != 5 || len(all[2].Raw) != 7 {
		t.Errorf("trim kept wrong entries: %+v", all)
	}

	// File permissions survived the rewrite.
	info, err := os.Stat(s.Path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("after trim perms = %o, want 0600", info.Mode().Perm())
	}
}

func TestStore_Clear(t *testing.T) {
	s := newTestStore(t)
	_ = s.Append("a", nil, "a")
	_ = s.Append("b", nil, "b")

	if err := s.Clear(); err != nil {
		t.Fatalf("clear: %v", err)
	}
	all, err := s.ReadAll()
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(all) != 0 {
		t.Errorf("after clear got %d entries", len(all))
	}

	// Clearing a non-existent file is a no-op.
	s2 := &Store{Path: filepath.Join(t.TempDir(), "missing.jsonl")}
	if err := s2.Clear(); err != nil {
		t.Errorf("clear missing: %v", err)
	}
}

func TestStore_ReadAllSkipsMalformed(t *testing.T) {
	s := newTestStore(t)
	_ = s.Append("ok", nil, "ok")
	// Inject a bad line between two good ones.
	f, err := os.OpenFile(s.Path, os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if _, err := f.WriteString("not-json\n"); err != nil {
		t.Fatalf("write garbage: %v", err)
	}
	_ = f.Close()
	_ = s.Append("ok2", nil, "ok2")

	all, err := s.ReadAll()
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("want 2 good entries, got %d (%v)", len(all), all)
	}
}

func TestStore_LargeLineWithinBuffer(t *testing.T) {
	s := newTestStore(t)
	big := strings.Repeat("x", 200*1024) // 200 KiB, well over scanner default but under cap
	if err := s.Append("commit", []string{"-m", big}, "commit -m "+big); err != nil {
		t.Fatalf("append big: %v", err)
	}
	all, err := s.ReadAll()
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(all) != 1 || len(all[0].Args[1]) != len(big) {
		t.Fatalf("big entry round-trip failed: got %d entries", len(all))
	}
}

func TestEntryDisplayFallback(t *testing.T) {
	cases := []struct {
		name string
		in   Entry
		want string
	}{
		{"raw preferred", Entry{Command: "checkout", Args: []string{"main"}, Raw: "co main"}, "co main"},
		{"raw empty falls back to cmd+args", Entry{Command: "commit", Args: []string{"-m", "x"}}, "commit -m x"},
		{"raw empty no args", Entry{Command: "status"}, "status"},
	}
	for _, tc := range cases {
		if got := tc.in.Display(); got != tc.want {
			t.Errorf("%s: Display = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestDefaultPathIsUserScoped(t *testing.T) {
	p, err := DefaultPath()
	if err != nil {
		t.Fatalf("default path: %v", err)
	}
	if !strings.HasSuffix(p, "history.jsonl") {
		t.Errorf("default path %q should end with history.jsonl", p)
	}
	// Parent directory should exist and be 0700.
	dir := filepath.Dir(p)
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("stat dir: %v", err)
	}
	if info.Mode().Perm()&0o077 != 0 {
		t.Errorf("history dir perms %o should not be group/world accessible", info.Mode().Perm())
	}
}

func TestSetDefaultAndPackageWrappers(t *testing.T) {
	old := defaultStore
	t.Cleanup(func() { SetDefault(old) })

	s := newTestStore(t)
	SetDefault(s)

	if err := AppendCommand("status", nil, "status"); err != nil {
		t.Fatalf("AppendCommand: %v", err)
	}
	all, err := ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("want 1, got %d", len(all))
	}
	last, err := ReadLast(5)
	if err != nil || len(last) != 1 {
		t.Errorf("ReadLast: %v, %v", last, err)
	}
	if got, err := Search("status"); err != nil || len(got) != 1 {
		t.Errorf("Search: %v, %v", got, err)
	}
	if err := Clear(); err != nil {
		t.Errorf("Clear: %v", err)
	}
}
