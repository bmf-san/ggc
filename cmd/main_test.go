package cmd

import (
	"os"
	"testing"

	"github.com/bmf-san/ggc/v8/internal/history"
)

// TestMain swaps the package-level history store for a disabled one so
// that running `go test ./cmd/...` never touches the real per-user
// history file. Without this, every test that builds a Cmd and invokes
// Route() would append to /tmp/ggc-<uid>/history.jsonl and shared state
// would leak between unrelated packages.
func TestMain(m *testing.M) {
	prev := history.Default()
	history.SetDefault(&history.Store{Disabled: true})
	code := m.Run()
	history.SetDefault(prev)
	os.Exit(code)
}
