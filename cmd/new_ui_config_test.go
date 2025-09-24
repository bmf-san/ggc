package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bmf-san/ggc/v6/internal/testutil"
)

func TestNewUIHonorsConfigProfileAndOverrides(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)

	configContent := `interactive:
  profile: emacs
  darwin:
    keybindings:
      move_down: "Ctrl+J"
`
	configPath := filepath.Join(tempHome, ".ggcconfig.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	gitClient := testutil.NewMockGitClient()
	ui := NewUI(gitClient)

	if ui.profile != ProfileEmacs {
		t.Fatalf("profile = %v, want %v", ui.profile, ProfileEmacs)
	}

	ui.contextMgr.resolver.platform = "darwin"
	ui.contextMgr.resolver.ClearCache()

	keyMap, err := ui.contextMgr.resolver.Resolve(ProfileEmacs, ContextResults)
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}

	if len(keyMap.MoveDown) == 0 || keyMap.MoveDown[0].Rune != 'j' {
		t.Fatalf("darwin override not applied: %#v", keyMap.MoveDown)
	}
}
