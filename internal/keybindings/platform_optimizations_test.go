package keybindings

import "testing"

func TestPlatformOptimizationsMappings(t *testing.T) {
	po := NewPlatformOptimizations("darwin", "wezterm")
	bindings, ok := po.GetOptimizedBindings("delete_word")
	if !ok || len(bindings) == 0 {
		t.Fatalf("expected delete_word bindings for macOS")
	}
	if bindings[0].Kind != KeyStrokeAlt {
		t.Fatalf("expected alt-based delete word binding, got %#v", bindings[0])
	}

	linuxPO := NewPlatformOptimizations("linux", "tmux")
	if linuxBindings, ok := linuxPO.GetOptimizedBindings("delete_word"); !ok || len(linuxBindings) == 0 {
		t.Fatalf("expected linux delete_word override")
	}

	unixPO := NewPlatformOptimizations("unix", "generic")
	clearBindings, ok := unixPO.GetOptimizedBindings("clear_line")
	if !ok || len(clearBindings) == 0 || clearBindings[0].Rune != 'u' {
		t.Fatalf("expected unix clear_line fallback, got %#v", clearBindings)
	}
}
