// Package testutil holds shared test helpers.
package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

// InstallEditorStubsOnPath places minimal executables named after the editors
// the default configuration references (vim, cat) at the front of PATH, so
// config validation by exec.LookPath succeeds on every OS — Windows and
// minimal container images ship neither command. Tests that load or save the
// real config (which validates the default editor) need this.
func InstallEditorStubsOnPath(t *testing.T, names ...string) {
	t.Helper()
	dir := t.TempDir()
	for _, name := range names {
		path := filepath.Join(dir, name+".exe")
		// Content is irrelevant; only that the file exists and LookPath
		// finds it. The 0o755 mode satisfies LookPath's exec-bit check on Unix.
		if err := os.WriteFile(path, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
			t.Fatalf("write editor stub %s: %v", name, err)
		}
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}
