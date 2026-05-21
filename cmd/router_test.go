package cmd

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/bmf-san/ggc/v8/internal/config"
	"github.com/bmf-san/ggc/v8/internal/history"
)

// installIsolatedHistory points the package-level history store at a
// throwaway path and returns it so tests can inspect persisted writes
// from the router. The previous store is restored via t.Cleanup so we
// do not leak state into other tests in this package.
func installIsolatedHistory(t *testing.T) *history.Store {
	t.Helper()
	store := &history.Store{Path: filepath.Join(t.TempDir(), "history.jsonl")}
	prev := history.Default()
	history.SetDefault(store)
	t.Cleanup(func() { history.SetDefault(prev) })
	return store
}

func newRouterCmd(t *testing.T) *Cmd {
	t.Helper()
	mockClient := &mockGitClient{}
	cm := config.NewConfigManager(mockClient)
	cmd, err := NewCmd(mockClient, cm)
	if err != nil {
		t.Fatalf("NewCmd: %v", err)
	}
	cmd.outputWriter = &bytes.Buffer{}
	return cmd
}

func TestRouter_RecordsCommand(t *testing.T) {
	store := installIsolatedHistory(t)
	cmd := newRouterCmd(t)

	if err := cmd.Route([]string{"version"}); err != nil {
		t.Fatalf("Route: %v", err)
	}

	all, err := store.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(all) != 1 || all[0].Command != "version" {
		t.Fatalf("expected one 'version' entry, got %+v", all)
	}
}

func TestRouter_SkipsHistoryItself(t *testing.T) {
	store := installIsolatedHistory(t)
	cmd := newRouterCmd(t)

	// Invoking `ggc history` should not pollute the history with a
	// "history" entry — otherwise every `history search` query would
	// match the searches you just ran.
	if err := cmd.Route([]string{"history"}); err != nil {
		t.Fatalf("Route: %v", err)
	}
	all, err := store.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(all) != 0 {
		t.Fatalf("history command should not record itself, got %+v", all)
	}
}

func TestRouter_DisabledStoreNoOp(t *testing.T) {
	store := installIsolatedHistory(t)
	store.Disabled = true
	cmd := newRouterCmd(t)

	if err := cmd.Route([]string{"version"}); err != nil {
		t.Fatalf("Route: %v", err)
	}
	all, err := store.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(all) != 0 {
		t.Fatalf("disabled store should not record, got %+v", all)
	}
}

func TestRouter_PreservesAliasInRaw(t *testing.T) {
	// `version` has no alias today, so we simulate the canonical/typed
	// distinction by feeding the canonical name and asserting that the
	// stored raw line reflects what was actually typed (plus args).
	store := installIsolatedHistory(t)
	cmd := newRouterCmd(t)

	if err := cmd.Route([]string{"version", "json"}); err != nil {
		t.Fatalf("Route: %v", err)
	}
	all, err := store.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("want 1 entry, got %d", len(all))
	}
	if all[0].Raw != "version json" {
		t.Errorf("raw = %q, want %q", all[0].Raw, "version json")
	}
	if len(all[0].Args) != 1 || all[0].Args[0] != "json" {
		t.Errorf("args = %v", all[0].Args)
	}
}
