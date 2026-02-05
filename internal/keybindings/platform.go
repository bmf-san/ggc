package keybindings

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// DetectPlatform identifies the current operating system platform
func DetectPlatform() string {
	switch runtime.GOOS {
	case "darwin":
		return "darwin"
	case "linux":
		return "linux"
	case "windows":
		return "windows"
	case "freebsd", "openbsd", "netbsd":
		return "bsd"
	default:
		return "unix"
	}
}

// DetectTerminal identifies the current terminal type from environment variables
func DetectTerminal() string { //nolint:revive // terminal detection relies on heuristics
	term := os.Getenv("TERM")
	termProgram := os.Getenv("TERM_PROGRAM")

	// Check TERM_PROGRAM first (more specific)
	switch termProgram {
	case "iTerm.app":
		return "iterm"
	case "Apple_Terminal":
		return "terminal"
	case "vscode":
		return "vscode"
	case "Hyper":
		return "hyper"
	}

	// Check TERM environment variable
	switch {
	case strings.Contains(term, "tmux"):
		return "tmux"
	case strings.Contains(term, "screen"):
		return "screen"
	case strings.HasPrefix(term, "xterm"):
		return "xterm"
	case strings.Contains(term, "alacritty"):
		return "alacritty"
	case strings.Contains(term, "kitty"):
		return "kitty"
	case strings.Contains(term, "wezterm"):
		return "wezterm"
	case strings.Contains(term, "konsole"):
		return "konsole"
	case strings.Contains(term, "gnome"):
		return "gnome-terminal"
	case strings.Contains(term, "rxvt"):
		return "rxvt"
	case term == "dumb":
		return "dumb"
	default:
		return "generic"
	}
}

// GetTerminalCapabilities returns a set of capabilities for the detected terminal
func GetTerminalCapabilities(terminal string) map[string]bool {
	capabilities := make(map[string]bool)

	switch terminal {
	case "iterm", "alacritty", "kitty", "wezterm":
		// Modern terminals with full capability
		capabilities["alt_keys"] = true
		capabilities["function_keys"] = true
		capabilities["mouse"] = true
		capabilities["color_256"] = true
		capabilities["unicode"] = true

	case "xterm", "gnome-terminal", "konsole":
		// Standard terminals
		capabilities["alt_keys"] = true
		capabilities["function_keys"] = true
		capabilities["mouse"] = false
		capabilities["color_256"] = true
		capabilities["unicode"] = true

	case "tmux", "screen":
		// Terminal multiplexers
		capabilities["alt_keys"] = true // may need prefix
		capabilities["function_keys"] = true
		capabilities["mouse"] = false
		capabilities["color_256"] = true
		capabilities["unicode"] = true

	case "terminal": // macOS Terminal
		// macOS Terminal specifics
		capabilities["alt_keys"] = true
		capabilities["function_keys"] = true
		capabilities["mouse"] = false
		capabilities["color_256"] = true
		capabilities["unicode"] = true

	case "dumb":
		// Minimal terminal
		capabilities["alt_keys"] = false
		capabilities["function_keys"] = false
		capabilities["mouse"] = false
		capabilities["color_256"] = false
		capabilities["unicode"] = false

	default:
		// Generic terminal - assume basic capabilities
		capabilities["alt_keys"] = true
		capabilities["function_keys"] = false
		capabilities["mouse"] = false
		capabilities["color_256"] = false
		capabilities["unicode"] = true
	}

	return capabilities
}

// GetPlatformSpecificKeyBindings returns platform-specific keybinding adjustments
func GetPlatformSpecificKeyBindings(platform string) map[string][]KeyStroke {
	platformBindings := make(map[string][]KeyStroke)

	switch platform {
	case "darwin":
		// macOS specific bindings
		// Option+Backspace for delete word (common macOS behavior)
		platformBindings["delete_word"] = []KeyStroke{NewAltKeyStroke(0, "backspace")}
		// Command key handling would go here if we supported it

	case "windows":
		// Windows specific bindings
		// Windows typically uses Ctrl+Backspace for delete word
		// NOTE: Ctrl+Backspace is not supported by NewCtrlKeyStroke; omitting until proper encoding is supported.

	case "linux", "bsd", "unix":
		// Unix-like systems - typically follow readline conventions
		// Most Linux terminals use Alt+Backspace or Ctrl+W
		platformBindings["delete_word"] = []KeyStroke{
			NewCtrlKeyStroke('w'),
			NewAltKeyStroke(0, "backspace"),
		}

	default:
		// No platform-specific adjustments
	}

	return platformBindings
}

// GetTerminalSpecificKeyBindings returns terminal-specific keybinding adjustments
func GetTerminalSpecificKeyBindings(terminal string) map[string][]KeyStroke {
	terminalBindings := make(map[string][]KeyStroke)

	switch terminal {
	case "tmux":
		// tmux prefix handling - these would need special handling
		// For now, just document that some keys might need prefix
		break

	case "screen":
		// GNU Screen specific adjustments
		break

	case "iterm":
		// iTerm2 specific features
		break

	case "alacritty", "kitty", "wezterm":
		// Modern terminal features
		break

	default:
		// No terminal-specific adjustments
	}

	return terminalBindings
}

// detectConflicts finds duplicate key assignments in a KeyBindingMap (legacy compatibility)
func detectConflicts(keyMap *KeyBindingMap) []string {
	// Convert to extended format and use newer conflict detection
	return detectConflictsV2(keyMap)
}

// detectConflictsV2 finds duplicate KeyStroke assignments in a KeyBindingMap (extended)
func detectConflictsV2(keyMap *KeyBindingMap) []string {
	var conflicts []string

	// Build a map of KeyStrokes to actions
	keystrokeToActions := make(map[string][]string)

	// Helper function to add KeyStrokes to conflict map
	addKeyStrokes := func(keyStrokes []KeyStroke, action string) {
		for _, ks := range keyStrokes {
			key := ks.String()
			keystrokeToActions[key] = append(keystrokeToActions[key], action)
		}
	}

	// Add all actions
	addKeyStrokes(keyMap.DeleteWord, "delete_word")
	addKeyStrokes(keyMap.ClearLine, "clear_line")
	addKeyStrokes(keyMap.DeleteToEnd, "delete_to_end")
	addKeyStrokes(keyMap.MoveToBeginning, "move_to_beginning")
	addKeyStrokes(keyMap.MoveToEnd, "move_to_end")
	addKeyStrokes(keyMap.MoveUp, "move_up")
	addKeyStrokes(keyMap.MoveDown, "move_down")
	addKeyStrokes(keyMap.MoveLeft, "move_left")
	addKeyStrokes(keyMap.MoveRight, "move_right")
	addKeyStrokes(keyMap.AddToWorkflow, "add_to_workflow")
	addKeyStrokes(keyMap.ToggleWorkflowView, "toggle_workflow_view")
	addKeyStrokes(keyMap.ClearWorkflow, "clear_workflow")

	// Find conflicts (multiple actions for same keystroke)
	for keystroke, actions := range keystrokeToActions {
		if len(actions) > 1 {
			conflicts = append(conflicts, fmt.Sprintf("keystroke %s assigned to: %v", keystroke, actions))
		}
	}

	return conflicts
}

// PlatformOptimizations provides platform-specific keybinding recommendations.
type PlatformOptimizations struct {
	platform string
	terminal string
	keyMap   map[string][]KeyStroke
}

// NewPlatformOptimizations builds platform-aware keybinding suggestions.
func NewPlatformOptimizations(platform, terminal string) *PlatformOptimizations {
	po := &PlatformOptimizations{
		platform: platform,
		terminal: terminal,
		keyMap:   make(map[string][]KeyStroke),
	}

	po.initialize()
	return po
}

func (po *PlatformOptimizations) initialize() {
	switch po.platform {
	case "darwin":
		po.keyMap["delete_word"] = []KeyStroke{
			NewAltKeyStroke(0, "alt+backspace"),
			NewCtrlKeyStroke('w'),
		}
	case "linux":
		po.keyMap["delete_word"] = []KeyStroke{NewCtrlKeyStroke('w')}
		po.keyMap["clear_line"] = []KeyStroke{NewCtrlKeyStroke('u')}
	case "windows":
		po.keyMap["paste"] = []KeyStroke{NewCtrlKeyStroke('v')}
		po.keyMap["clear_line"] = []KeyStroke{NewCtrlKeyStroke('u')}
	default:
		po.keyMap["delete_word"] = []KeyStroke{NewCtrlKeyStroke('w')}
		po.keyMap["clear_line"] = []KeyStroke{NewCtrlKeyStroke('u')}
	}

	if po.platform == "darwin" {
		po.keyMap["word_forward"] = []KeyStroke{NewAltKeyStroke('f', "alt+f")}
		po.keyMap["word_backward"] = []KeyStroke{NewAltKeyStroke('b', "alt+b")}
	}
}

// GetOptimizedBindings returns platform-aware bindings for the given action.
func (po *PlatformOptimizations) GetOptimizedBindings(action string) ([]KeyStroke, bool) {
	bindings, ok := po.keyMap[action]
	return bindings, ok
}
