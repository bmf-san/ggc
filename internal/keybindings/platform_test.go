package keybindings

import "testing"

// ── DetectTerminal ───────────────────────────────────────────────────────────

func TestDetectTerminal_TermProgram(t *testing.T) {
	t.Setenv("TERM", "")

	cases := []struct {
		prog string
		want string
	}{
		{"iTerm.app", "iterm"},
		{"Apple_Terminal", "terminal"},
		{"vscode", "vscode"},
		{"Hyper", "hyper"},
	}
	for _, c := range cases {
		t.Setenv("TERM_PROGRAM", c.prog)
		if got := DetectTerminal(); got != c.want {
			t.Errorf("TERM_PROGRAM=%q: want %q, got %q", c.prog, c.want, got)
		}
	}
}

func TestDetectTerminal_TERM(t *testing.T) {
	t.Setenv("TERM_PROGRAM", "")

	cases := []struct {
		term string
		want string
	}{
		{"tmux-256color", "tmux"},
		{"screen-256color", "screen"},
		{"xterm-256color", "xterm"},
		{"alacritty", "alacritty"},
		{"kitty", "kitty"},
		{"wezterm", "wezterm"},
		{"konsole-256color", "konsole"},
		{"gnome-256color", "gnome-terminal"},
		{"rxvt-unicode", "rxvt"},
		{"dumb", "dumb"},
		{"", "generic"},
	}
	for _, c := range cases {
		t.Setenv("TERM", c.term)
		if got := DetectTerminal(); got != c.want {
			t.Errorf("TERM=%q: want %q, got %q", c.term, c.want, got)
		}
	}
}

// ── GetTerminalCapabilities ──────────────────────────────────────────────────

func TestGetTerminalCapabilities_StandardTerminals(t *testing.T) {
	for _, terminal := range []string{"xterm", "gnome-terminal", "konsole"} {
		caps := GetTerminalCapabilities(terminal)
		if !caps["alt_keys"] {
			t.Errorf("%s: expected alt_keys=true", terminal)
		}
		if caps["mouse"] {
			t.Errorf("%s: expected mouse=false", terminal)
		}
	}
}

func TestGetTerminalCapabilities_Multiplexers(t *testing.T) {
	for _, terminal := range []string{"tmux", "screen"} {
		caps := GetTerminalCapabilities(terminal)
		if !caps["alt_keys"] {
			t.Errorf("%s: expected alt_keys=true", terminal)
		}
		if caps["mouse"] {
			t.Errorf("%s: expected mouse=false", terminal)
		}
	}
}

func TestGetTerminalCapabilities_MacTerminal(t *testing.T) {
	caps := GetTerminalCapabilities("terminal")
	if !caps["alt_keys"] {
		t.Error("terminal: expected alt_keys=true")
	}
	if caps["mouse"] {
		t.Error("terminal: expected mouse=false")
	}
}

func TestGetTerminalCapabilities_Generic(t *testing.T) {
	caps := GetTerminalCapabilities("generic")
	if !caps["alt_keys"] {
		t.Error("generic: expected alt_keys=true")
	}
	if caps["function_keys"] {
		t.Error("generic: expected function_keys=false")
	}
	if caps["mouse"] {
		t.Error("generic: expected mouse=false")
	}
}

func TestGetTerminalCapabilities_Unknown(t *testing.T) {
	caps := GetTerminalCapabilities("unknown-terminal-xyz")
	// default branch — same as generic
	if caps["mouse"] {
		t.Error("unknown terminal: expected mouse=false")
	}
}

// ── GetPlatformSpecificKeyBindings ───────────────────────────────────────────

func TestGetPlatformSpecificKeyBindings_Linux(t *testing.T) {
	bindings := GetPlatformSpecificKeyBindings("linux")
	ks, ok := bindings["delete_word"]
	if !ok {
		t.Fatal("linux: expected delete_word binding")
	}
	if len(ks) == 0 {
		t.Error("linux: expected at least one delete_word keystroke")
	}
}

func TestGetPlatformSpecificKeyBindings_BSD(t *testing.T) {
	bindings := GetPlatformSpecificKeyBindings("bsd")
	if _, ok := bindings["delete_word"]; !ok {
		t.Error("bsd: expected delete_word binding")
	}
}

func TestGetPlatformSpecificKeyBindings_Unix(t *testing.T) {
	bindings := GetPlatformSpecificKeyBindings("unix")
	if _, ok := bindings["delete_word"]; !ok {
		t.Error("unix: expected delete_word binding")
	}
}

func TestGetPlatformSpecificKeyBindings_Windows(t *testing.T) {
	// Windows branch exists but currently returns empty map
	bindings := GetPlatformSpecificKeyBindings("windows")
	_ = bindings // no panic = pass
}

func TestGetPlatformSpecificKeyBindings_Default(t *testing.T) {
	bindings := GetPlatformSpecificKeyBindings("unknown-platform")
	if len(bindings) != 0 {
		t.Error("unknown platform: expected empty bindings")
	}
}

// ── GetTerminalSpecificKeyBindings ───────────────────────────────────────────

func TestGetTerminalSpecificKeyBindings_All(t *testing.T) {
	for _, terminal := range []string{"tmux", "screen", "iterm", "alacritty", "kitty", "wezterm", "other"} {
		bindings := GetTerminalSpecificKeyBindings(terminal)
		_ = bindings // all return empty map — no panic = pass
	}
}
