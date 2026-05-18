package history

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAppendReadSearch(t *testing.T) {
	dir := t.TempDir()
	// override UserConfigDir by setting XDG_CONFIG_HOME and HOME (macOS)
	os.Setenv("XDG_CONFIG_HOME", dir)
	defer os.Unsetenv("XDG_CONFIG_HOME")
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", oldHome)

	if err := AppendCommand([]string{"status"}); err != nil {
		t.Fatalf("append failed: %v", err)
	}
	if err := AppendCommand([]string{"commit", "-m", "test"}); err != nil {
		t.Fatalf("append 2 failed: %v", err)
	}

	all, err := ReadAll()
	if err != nil {
		t.Fatalf("read all failed: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(all))
	}

	last, err := ReadLast(1)
	if err != nil {
		t.Fatalf("read last failed: %v", err)
	}
	if len(last) != 1 {
		t.Fatalf("expected 1 last, got %d", len(last))
	}

	// Search
	s, err := Search("commit")
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if len(s) != 1 {
		t.Fatalf("expected 1 search hit, got %d", len(s))
	}

	// ensure file exists at user config dir
	cfg, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("user config dir: %v", err)
	}
	hf := filepath.Join(cfg, "ggc", "history")
	if _, err := os.Stat(hf); err != nil {
		t.Fatalf("history file missing: %v", err)
	}
}
