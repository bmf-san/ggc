// Package history persists ggc command invocations to a per-user file
// under the OS temp directory so that recent commands survive across
// processes within a single boot. The file is deliberately ephemeral
// (cleared on reboot) to keep the security envelope close to shell
// history while still enabling `ggc history` lookups across separate
// CLI invocations and the interactive REPL.
package history

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// DefaultMaxEntries is the cap applied when no explicit limit is
// configured. Entries beyond this are truncated from the head on the
// next write.
const DefaultMaxEntries = 1000

// scannerMaxBuffer bounds individual history lines. Commit messages or
// long argument lists can exceed bufio.Scanner's default 64 KiB, so we
// raise the ceiling to 1 MiB. Lines larger than this are dropped with
// an error rather than panicking.
const scannerMaxBuffer = 1024 * 1024

// envDisable, when set to a truthy value, suppresses every write to the
// history file. Mirrors the bash `HISTCONTROL=ignorespace` escape hatch
// in spirit: "I want this one invocation off the record".
const envDisable = "GGC_NO_HISTORY"

// Entry is one persisted command invocation.
type Entry struct {
	// Timestamp is when AppendCommand was called, in UTC.
	Timestamp time.Time `json:"ts"`
	// Command is the canonical command name resolved by the registry
	// (e.g. "checkout" even when the user typed an alias "co"). It is
	// the right key to feed back into the router for replays.
	Command string `json:"cmd"`
	// Args are the canonical args passed to the command handler. May be
	// nil for commands invoked without arguments.
	Args []string `json:"args,omitempty"`
	// Raw is the exact line the user typed, including aliases or
	// alternate spacing. Empty when the caller could not preserve the
	// original input, in which case Display() falls back to Command +
	// Args.
	Raw string `json:"raw,omitempty"`
}

// Display returns the human-readable form of the entry, preferring the
// original raw input so that aliases survive round-trips through the
// history command. Uses a pointer receiver because Entry carries a
// time.Time which makes value-receiver copies non-trivial.
func (e *Entry) Display() string {
	if e.Raw != "" {
		return e.Raw
	}
	if len(e.Args) == 0 {
		return e.Command
	}
	return e.Command + " " + strings.Join(e.Args, " ")
}

// Store is the persistence layer for history entries. A zero-value Store
// uses defaults; tests can override the path and the entry cap.
type Store struct {
	// Path is the JSONL file backing the store. When empty, DefaultPath
	// is used lazily on the first call.
	Path string
	// MaxEntries is the cap before truncate-rewrite kicks in. Values
	// <= 0 fall back to DefaultMaxEntries.
	MaxEntries int
	// Disabled short-circuits every write. Reads still work so that
	// users can inspect prior history while temporarily off the record.
	Disabled bool
}

// Default returns a Store wired to DefaultPath, DefaultMaxEntries, and
// honoring the GGC_NO_HISTORY env variable.
func Default() *Store {
	return &Store{
		Disabled: envTrue(os.Getenv(envDisable)),
	}
}

// DefaultPath returns the per-user history file location. On Unix-like
// systems it lives under os.TempDir() in a uid-scoped subdirectory; on
// Windows (where Getuid returns -1) it falls back to UserCacheDir so the
// path is still user-private without colliding in a shared /tmp.
func DefaultPath() (string, error) {
	if runtime.GOOS == "windows" {
		base, err := os.UserCacheDir()
		if err != nil {
			return "", fmt.Errorf("locate user cache dir: %w", err)
		}
		dir := filepath.Join(base, "ggc")
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return "", fmt.Errorf("create history dir: %w", err)
		}
		return filepath.Join(dir, "history.jsonl"), nil
	}
	uid := os.Getuid()
	subdir := "ggc-" + strconv.Itoa(uid)
	dir := filepath.Join(os.TempDir(), subdir)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("create history dir: %w", err)
	}
	return filepath.Join(dir, "history.jsonl"), nil
}

func envTrue(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "", "0", "false", "no", "off":
		return false
	}
	return true
}

func (s *Store) path() (string, error) {
	if s.Path != "" {
		// Ensure parent exists for caller-supplied paths (typical in tests).
		if err := os.MkdirAll(filepath.Dir(s.Path), 0o700); err != nil {
			return "", err
		}
		return s.Path, nil
	}
	return DefaultPath()
}

func (s *Store) cap() int {
	if s.MaxEntries > 0 {
		return s.MaxEntries
	}
	return DefaultMaxEntries
}

// Append records a single invocation. It is best-effort: on disabled or
// empty command the call returns nil without touching the disk. Callers
// typically ignore the error and let the user keep working if the FS is
// in a bad state.
func (s *Store) Append(command string, args []string, raw string) error {
	if s.Disabled || strings.TrimSpace(command) == "" {
		return nil
	}
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Command:   command,
		Args:      args,
		Raw:       raw,
	}
	path, err := s.path()
	if err != nil {
		return err
	}
	line, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}
	if _, err := f.Write(append(line, '\n')); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return s.trim()
}

// ReadAll returns every entry in chronological order (oldest first).
// A missing file is not an error; it simply yields an empty slice.
func (s *Store) ReadAll() ([]Entry, error) {
	path, err := s.path()
	if err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	defer func() { _ = f.Close() }()
	return decodeEntries(f)
}

func decodeEntries(r io.Reader) ([]Entry, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), scannerMaxBuffer)
	var out []Entry
	for scanner.Scan() {
		raw := scanner.Bytes()
		if len(raw) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(raw, &e); err != nil {
			// Skip malformed lines (could be a partial write from a
			// concurrent process) rather than failing the whole read.
			continue
		}
		out = append(out, e)
	}
	if err := scanner.Err(); err != nil {
		return out, err
	}
	return out, nil
}

// ReadLast returns up to n most recent entries (oldest first, newest
// last). Non-positive n returns an empty slice; n larger than the
// available history returns everything.
func (s *Store) ReadLast(n int) ([]Entry, error) {
	if n <= 0 {
		return nil, nil
	}
	all, err := s.ReadAll()
	if err != nil {
		return all, err
	}
	if n >= len(all) {
		return all, nil
	}
	return all[len(all)-n:], nil
}

// Search returns entries whose canonical or raw form contains pattern
// (case-insensitive). An empty pattern matches everything.
func (s *Store) Search(pattern string) ([]Entry, error) {
	all, err := s.ReadAll()
	if err != nil {
		return nil, err
	}
	if pattern == "" {
		return all, nil
	}
	low := strings.ToLower(pattern)
	out := make([]Entry, 0, len(all))
	for _, e := range all {
		hay := strings.ToLower(e.Display() + " " + e.Command)
		if strings.Contains(hay, low) {
			out = append(out, e)
		}
	}
	return out, nil
}

// Clear removes every persisted entry by truncating the file. A missing
// file is treated as success.
func (s *Store) Clear() error {
	path, err := s.path()
	if err != nil {
		return err
	}
	if err := os.Truncate(path, 0); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

// trim rewrites the history file so it contains at most s.cap() entries,
// keeping the newest. Called after every successful append. The rewrite
// uses a temp file + rename so a crash mid-trim leaves the previous
// state intact.
func (s *Store) trim() error {
	max := s.cap()
	all, err := s.ReadAll()
	if err != nil || len(all) <= max {
		return err
	}
	path, err := s.path()
	if err != nil {
		return err
	}
	return rewriteKeeping(path, all[len(all)-max:])
}

// rewriteKeeping atomically replaces path with a freshly-written JSONL
// file containing only the supplied entries. On any failure the temp
// file is cleaned up so we never leave half-written rewrites around.
func rewriteKeeping(path string, entries []Entry) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".history-*.jsonl")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	if err := writeEntries(tmp, entries); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	if err := os.Chmod(tmpPath, 0o600); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	return os.Rename(tmpPath, path)
}

// writeEntries marshals each entry as JSONL into w, flushing before
// return. The caller owns w and is responsible for closing it.
func writeEntries(w io.Writer, entries []Entry) error {
	bw := bufio.NewWriter(w)
	for i := range entries {
		b, err := json.Marshal(entries[i])
		if err != nil {
			return err
		}
		if _, err := bw.Write(append(b, '\n')); err != nil {
			return err
		}
	}
	return bw.Flush()
}

// Package-level convenience wrappers using the Default store. These keep
// older call-sites compatible while letting tests inject a Store
// directly when they need control over the path or cap.
var defaultStore = Default()

// SetDefault swaps the package-level store. Tests can use this to point
// the convenience wrappers at a temp directory.
func SetDefault(s *Store) { defaultStore = s }

// AppendCommand records command + args on the default store.
func AppendCommand(command string, args []string, raw string) error {
	return defaultStore.Append(command, args, raw)
}

// ReadAll reads every entry from the default store.
func ReadAll() ([]Entry, error) { return defaultStore.ReadAll() }

// ReadLast reads the last n entries from the default store.
func ReadLast(n int) ([]Entry, error) { return defaultStore.ReadLast(n) }

// Search runs Search on the default store.
func Search(pattern string) ([]Entry, error) { return defaultStore.Search(pattern) }

// Clear truncates the default store.
func Clear() error { return defaultStore.Clear() }
