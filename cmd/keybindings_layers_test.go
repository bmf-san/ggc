package cmd

import (
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v5/config"
)

// TestKeyBindingResolution tests the layering and resolution logic
func TestKeyBindingResolution(t *testing.T) {
	tests := []struct {
		name     string
		config   *config.Config
		expected *KeyBindingMap
	}{
		{
			name:   "defaults only",
			config: &config.Config{},
			expected: &KeyBindingMap{
				DeleteWord:      []KeyStroke{NewCtrlKeyStroke('w')},
				ClearLine:       []KeyStroke{NewCtrlKeyStroke('u')},
				DeleteToEnd:     []KeyStroke{NewCtrlKeyStroke('k')},
				MoveToBeginning: []KeyStroke{NewCtrlKeyStroke('a')},
				MoveToEnd:       []KeyStroke{NewCtrlKeyStroke('e')},
				MoveUp:          []KeyStroke{NewCtrlKeyStroke('p')},
				MoveDown:        []KeyStroke{NewCtrlKeyStroke('n')},
			},
		},
		{
			name: "user overrides single binding",
			config: func() *config.Config {
				cfg := &config.Config{}
				cfg.Interactive.Keybindings.DeleteWord = "ctrl+o"
				return cfg
			}(),
			expected: &KeyBindingMap{
				DeleteWord:      []KeyStroke{NewCtrlKeyStroke('o')}, // overridden
				ClearLine:       []KeyStroke{NewCtrlKeyStroke('u')}, // default
				DeleteToEnd:     []KeyStroke{NewCtrlKeyStroke('k')}, // default
				MoveToBeginning: []KeyStroke{NewCtrlKeyStroke('a')}, // default
				MoveToEnd:       []KeyStroke{NewCtrlKeyStroke('e')}, // default
				MoveUp:          []KeyStroke{NewCtrlKeyStroke('p')}, // default
				MoveDown:        []KeyStroke{NewCtrlKeyStroke('n')}, // default
			},
		},
		{
			name: "multiple user overrides",
			config: func() *config.Config {
				cfg := &config.Config{}
				cfg.Interactive.Keybindings.DeleteWord = "ctrl+o"
				cfg.Interactive.Keybindings.MoveUp = "ctrl+l"
				cfg.Interactive.Keybindings.MoveDown = "ctrl+j"
				return cfg
			}(),
			expected: &KeyBindingMap{
				DeleteWord:      []KeyStroke{NewCtrlKeyStroke('o')}, // overridden
				ClearLine:       []KeyStroke{NewCtrlKeyStroke('u')}, // default
				DeleteToEnd:     []KeyStroke{NewCtrlKeyStroke('k')}, // default
				MoveToBeginning: []KeyStroke{NewCtrlKeyStroke('a')}, // default
				MoveToEnd:       []KeyStroke{NewCtrlKeyStroke('e')}, // default
				MoveUp:          []KeyStroke{NewCtrlKeyStroke('l')}, // overridden
				MoveDown:        []KeyStroke{NewCtrlKeyStroke('j')}, // overridden
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolveKeyBindingMapForTest(t, tt.config)

			if !keyBindingMapsEqual(result, tt.expected) {
				t.Errorf("resolved map = %+v, expected %+v", result, tt.expected)
			}
		})
	}
}

// TestKeyBindingConflictDetection tests conflict detection in key bindings
func TestKeyBindingConflictDetection(t *testing.T) {
	tests := []struct {
		name            string
		config          *config.Config
		expectConflicts bool
	}{
		{
			name:            "no conflicts",
			config:          &config.Config{},
			expectConflicts: false,
		},
		{
			name: "conflict detected emits warning",
			config: func() *config.Config {
				cfg := &config.Config{}
				cfg.Interactive.Keybindings.DeleteWord = "ctrl+k"
				cfg.Interactive.Keybindings.DeleteToEnd = "ctrl+k"
				return cfg
			}(),
			expectConflicts: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyMap := resolveKeyBindingMapForTest(t, tt.config)

			conflicts := detectConflicts(keyMap)
			if tt.expectConflicts && len(conflicts) == 0 {
				t.Errorf("Expected conflicts, but none detected")
			}
			if !tt.expectConflicts && len(conflicts) > 0 {
				t.Errorf("Expected no conflicts, but detected: %v", conflicts)
			}
		})
	}
}

// TestParseKeyBindingExtended tests extended parsing capabilities
func TestParseKeyBindingExtended(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected byte
		wantErr  bool
	}{
		// Current supported formats
		{"simple ctrl", "ctrl+w", ctrl('w'), false},
		{"uppercase", "CTRL+W", ctrl('w'), false},
		{"mixed case", "Ctrl+W", ctrl('w'), false},

		// Additional supported formats
		{"caret notation", "^W", ctrl('w'), false},
		{"emacs notation", "C-w", ctrl('w'), false},

		// Invalid inputs
		{"empty", "", 0, true},
		{"invalid key", "ctrl+@", 0, true},
		{"unsupported", "alt+w", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseKeyBinding(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseKeyBinding(%q) expected error, got result: %d", tt.input, result)
				}
			} else {
				if err != nil {
					t.Errorf("ParseKeyBinding(%q) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ParseKeyBinding(%q) = %d, expected %d", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestKeyBindingValidation tests validation of key binding values
func TestKeyBindingValidation(t *testing.T) {
	tests := []struct {
		name        string
		bindings    map[string]string
		expectError bool
		errorMsg    string
	}{
		{
			name: "all valid",
			bindings: map[string]string{
				"delete_word": "ctrl+w",
				"clear_line":  "ctrl+u",
			},
			expectError: false,
		},
		{
			name: "invalid binding format",
			bindings: map[string]string{
				"delete_word": "invalid",
			},
			expectError: true,
			errorMsg:    "unsupported key binding",
		},
		{
			name: "empty binding",
			bindings: map[string]string{
				"delete_word": "",
			},
			expectError: true,
			errorMsg:    "empty key binding",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateKeyBindings(tt.bindings)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error containing '%s', got nil", tt.errorMsg)
				} else if !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing '%s', got: %v", tt.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// keyBindingMapsEqual compares two KeyBindingMaps for equality
func keyBindingMapsEqual(a, b *KeyBindingMap) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	return keyStrokesEqual(a.DeleteWord, b.DeleteWord) &&
		keyStrokesEqual(a.ClearLine, b.ClearLine) &&
		keyStrokesEqual(a.DeleteToEnd, b.DeleteToEnd) &&
		keyStrokesEqual(a.MoveToBeginning, b.MoveToBeginning) &&
		keyStrokesEqual(a.MoveToEnd, b.MoveToEnd) &&
		keyStrokesEqual(a.MoveUp, b.MoveUp) &&
		keyStrokesEqual(a.MoveDown, b.MoveDown)
}

// keyStrokesEqual compares two KeyStroke slices for equality
func keyStrokesEqual(a, b []KeyStroke) bool {
	if len(a) != len(b) {
		return false
	}
	for i, ks := range a {
		if !ks.Equals(b[i]) {
			return false
		}
	}
	return true
}
