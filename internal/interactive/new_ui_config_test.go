package interactive

import (
	"os"
	"path/filepath"
	"testing"

	kb "github.com/bmf-san/ggc/v7/internal/keybindings"
	"github.com/bmf-san/ggc/v7/internal/testutil"
)

func TestNewUIHonorsConfigProfileAndOverrides(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)

	configContent := `interactive:
  profile: emacs
  contexts:
    results:
      keybindings:
        move_down: "Ctrl+J"
`
	configPath := filepath.Join(tempHome, ".ggcconfig.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	gitClient := testutil.NewMockGitClient()
	ui := NewUI(gitClient)

	if ui.profile != kb.ProfileEmacs {
		t.Fatalf("profile = %v, want %v", ui.profile, kb.ProfileEmacs)
	}

	contextual := ui.handler.contextualMap
	if contextual == nil {
		t.Fatal("expected contextual keybinding map to be initialized")
	}

	resultsMap, exists := contextual.GetContext(kb.ContextResults)
	if !exists || resultsMap == nil {
		t.Fatal("expected results context map")
	}

	if len(resultsMap.MoveDown) == 0 {
		t.Fatalf("config override not applied: %#v", resultsMap.MoveDown)
	}
	stroke := resultsMap.MoveDown[0]
	if stroke.Kind != kb.KeyStrokeCtrl || stroke.Rune != 'j' {
		t.Fatalf("config override not applied: %#v", resultsMap.MoveDown)
	}
}
