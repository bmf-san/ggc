package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go.yaml.in/yaml/v3"

	"github.com/bmf-san/ggc/v5/config"
)

// TestKeyStroke tests the KeyStroke struct and its methods
func TestKeyStroke(t *testing.T) {
	t.Run("NewCtrlKeyStroke", func(t *testing.T) {
		ks := NewCtrlKeyStroke('w')
		if ks.Kind != KeyStrokeCtrl {
			t.Errorf("Expected Kind %v, got %v", KeyStrokeCtrl, ks.Kind)
		}
		if ks.Rune != 'w' {
			t.Errorf("Expected Rune 'w', got %c", ks.Rune)
		}
		if ks.String() != "Ctrl+w" {
			t.Errorf("Expected String 'Ctrl+w', got %s", ks.String())
		}
	})

	t.Run("NewAltKeyStroke", func(t *testing.T) {
		ks := NewAltKeyStroke(0, "backspace")
		if ks.Kind != KeyStrokeAlt {
			t.Errorf("Expected Kind %v, got %v", KeyStrokeAlt, ks.Kind)
		}
		if ks.Name != "backspace" {
			t.Errorf("Expected Name 'backspace', got %s", ks.Name)
		}
		if ks.String() != "Alt+backspace" {
			t.Errorf("Expected String 'Alt+backspace', got %s", ks.String())
		}
	})

	t.Run("ToControlByte", func(t *testing.T) {
		ctrlW := NewCtrlKeyStroke('w')
		if b := ctrlW.ToControlByte(); b != 23 {
			t.Errorf("Expected control byte 23, got %d", b)
		}

		altBack := NewAltKeyStroke(0, "backspace")
		if b := altBack.ToControlByte(); b != 0 {
			t.Errorf("Expected control byte 0 for Alt key, got %d", b)
		}
	})

	t.Run("Equals", func(t *testing.T) {
		ks1 := NewCtrlKeyStroke('w')
		ks2 := NewCtrlKeyStroke('w')
		ks3 := NewCtrlKeyStroke('u')

		if !ks1.Equals(ks2) {
			t.Error("Expected equal KeyStrokes to be equal")
		}
		if ks1.Equals(ks3) {
			t.Error("Expected different KeyStrokes to not be equal")
		}
	})
}

// TestParseKeyStroke tests the enhanced parser with additional formats
func TestParseKeyStroke(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected KeyStroke
		wantErr  bool
	}{
		// Legacy formats
		{"ctrl format", "ctrl+w", NewCtrlKeyStroke('w'), false},
		{"caret format", "^w", NewCtrlKeyStroke('w'), false},
		{"emacs format", "C-w", NewCtrlKeyStroke('w'), false},

		// Alt key formats
		{"alt letter", "alt+w", NewAltKeyStroke('w', ""), false},
		{"alt backspace", "alt+backspace", NewAltKeyStroke(0, "backspace"), false},
		{"alt delete", "alt+delete", NewAltKeyStroke(0, "delete"), false},
		{"meta backspace", "meta+backspace", NewAltKeyStroke(0, "backspace"), false},
		{"emacs meta", "M-backspace", NewAltKeyStroke(0, "backspace"), false},

		// Edge cases
		{"empty", "", KeyStroke{}, true},
		{"invalid", "invalid", KeyStroke{}, true},
		{"unsupported alt", "alt+invalid", KeyStroke{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseKeyStroke(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for input %q, got result: %v", tt.input, result)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input %q: %v", tt.input, err)
				}
				if !result.Equals(tt.expected) {
					t.Errorf("ParseKeyStroke(%q) = %v, expected %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestParseKeyStrokes tests array parsing functionality
func TestParseKeyStrokes(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []KeyStroke
		wantErr  bool
	}{
		{
			name:     "single string",
			input:    "ctrl+w",
			expected: []KeyStroke{NewCtrlKeyStroke('w')},
			wantErr:  false,
		},
		{
			name:     "string array",
			input:    []string{"ctrl+w", "alt+backspace"},
			expected: []KeyStroke{NewCtrlKeyStroke('w'), NewAltKeyStroke(0, "backspace")},
			wantErr:  false,
		},
		{
			name:     "interface array",
			input:    []interface{}{"ctrl+w", "^u"},
			expected: []KeyStroke{NewCtrlKeyStroke('w'), NewCtrlKeyStroke('u')},
			wantErr:  false,
		},
		{
			name:  "mixed formats",
			input: []string{"ctrl+w", "^u", "C-k", "alt+backspace", "M-delete"},
			expected: []KeyStroke{
				NewCtrlKeyStroke('w'),
				NewCtrlKeyStroke('u'),
				NewCtrlKeyStroke('k'),
				NewAltKeyStroke(0, "backspace"),
				NewAltKeyStroke(0, "delete"),
			},
			wantErr: false,
		},
		{
			name:    "invalid type",
			input:   123,
			wantErr: true,
		},
		{
			name:    "invalid array element",
			input:   []interface{}{"ctrl+w", 123},
			wantErr: true,
		},
		{
			name:    "invalid string in array",
			input:   []string{"ctrl+w", "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseKeyStrokes(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for input %v, got result: %v", tt.input, result)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input %v: %v", tt.input, err)
				}
				if len(result) != len(tt.expected) {
					t.Errorf("Length mismatch: got %d, expected %d", len(result), len(tt.expected))
				}
				for i, ks := range result {
					if i < len(tt.expected) && !ks.Equals(tt.expected[i]) {
						t.Errorf("KeyStroke %d: got %v, expected %v", i, ks, tt.expected[i])
					}
				}
			}
		})
	}
}

// TestKeyBindingMapExtended tests the updated KeyBindingMap structure
func TestKeyBindingMapExtended(t *testing.T) {
	t.Run("DefaultKeyBindingMap", func(t *testing.T) {
		km := DefaultKeyBindingMap()

		// Check that all fields are []KeyStroke with appropriate defaults
		if len(km.DeleteWord) != 1 || !km.DeleteWord[0].Equals(NewCtrlKeyStroke('w')) {
			t.Errorf("Expected DeleteWord [Ctrl+w], got %v", km.DeleteWord)
		}
		if len(km.ClearLine) != 1 || !km.ClearLine[0].Equals(NewCtrlKeyStroke('u')) {
			t.Errorf("Expected ClearLine [Ctrl+u], got %v", km.ClearLine)
		}
	})

	t.Run("MatchesKeyStroke", func(t *testing.T) {
		km := DefaultKeyBindingMap()

		// Test basic matching
		if !km.MatchesKeyStroke("delete_word", NewCtrlKeyStroke('w')) {
			t.Error("Expected Ctrl+w to match delete_word")
		}
		if km.MatchesKeyStroke("delete_word", NewCtrlKeyStroke('u')) {
			t.Error("Expected Ctrl+u to not match delete_word")
		}

		// Test with multiple bindings
		km.DeleteWord = []KeyStroke{
			NewCtrlKeyStroke('w'),
			NewAltKeyStroke(0, "backspace"),
		}

		if !km.MatchesKeyStroke("delete_word", NewCtrlKeyStroke('w')) {
			t.Error("Expected Ctrl+w to match delete_word with multiple bindings")
		}
		if !km.MatchesKeyStroke("delete_word", NewAltKeyStroke(0, "backspace")) {
			t.Error("Expected Alt+backspace to match delete_word with multiple bindings")
		}
	})

	t.Run("BackwardCompatibility", func(t *testing.T) {
		km := DefaultKeyBindingMap()

		// Test backward compatibility methods
		if km.GetDeleteWordByte() != 23 { // ctrl('w')
			t.Errorf("Expected GetDeleteWordByte() = 23, got %d", km.GetDeleteWordByte())
		}
		if km.GetClearLineByte() != 21 { // ctrl('u')
			t.Errorf("Expected GetClearLineByte() = 21, got %d", km.GetClearLineByte())
		}
	})
}

// TestResolveKeyBindingMap tests resolution with additional features
func TestResolveKeyBindingMap(t *testing.T) {
	t.Run("single string config", func(t *testing.T) {
		cfg := &config.Config{}
		cfg.Interactive.Keybindings.DeleteWord = "alt+backspace"
		cfg.Interactive.Keybindings.MoveUp = "^k"

		km := resolveKeyBindingMapForTest(t, cfg, ContextInput)

		// Should have Alt+backspace for delete_word
		expectedDeleteWord := NewAltKeyStroke(0, "backspace")
		if len(km.DeleteWord) != 1 || !km.DeleteWord[0].Equals(expectedDeleteWord) {
			t.Errorf("Expected DeleteWord [Alt+backspace], got %v", km.DeleteWord)
		}

		// Should have Ctrl+k for move_up
		expectedMoveUp := NewCtrlKeyStroke('k')
		if len(km.MoveUp) != 1 || !km.MoveUp[0].Equals(expectedMoveUp) {
			t.Errorf("Expected MoveUp [Ctrl+k], got %v", km.MoveUp)
		}

		// Other fields should have defaults
		expectedClearLine := NewCtrlKeyStroke('u')
		if len(km.ClearLine) != 1 || !km.ClearLine[0].Equals(expectedClearLine) {
			t.Errorf("Expected ClearLine [Ctrl+u] (default), got %v", km.ClearLine)
		}
	})
}

// TestConflictDetectionExtended tests enhanced conflict detection
func TestConflictDetectionExtended(t *testing.T) {
	t.Run("no conflicts", func(t *testing.T) {
		km := DefaultKeyBindingMap()
		conflicts := detectConflictsV2(km)
		if len(conflicts) > 0 {
			t.Errorf("Expected no conflicts in defaults, got: %v", conflicts)
		}
	})

	t.Run("KeyStroke conflicts", func(t *testing.T) {
		km := DefaultKeyBindingMap()
		// Create conflict: both delete_word and clear_line use Ctrl+w
		km.DeleteWord = []KeyStroke{NewCtrlKeyStroke('w')}
		km.ClearLine = []KeyStroke{NewCtrlKeyStroke('w')}

		conflicts := detectConflictsV2(km)
		if len(conflicts) == 0 {
			t.Error("Expected conflicts to be detected")
		}

		// Should mention both actions
		conflictStr := conflicts[0]
		if !contains(conflictStr, "delete_word") || !contains(conflictStr, "clear_line") {
			t.Errorf("Expected conflict to mention both actions, got: %s", conflictStr)
		}
	})

	t.Run("multiple KeyStroke conflicts", func(t *testing.T) {
		km := DefaultKeyBindingMap()
		// Create multiple bindings with conflicts
		km.DeleteWord = []KeyStroke{NewCtrlKeyStroke('w'), NewAltKeyStroke(0, "backspace")}
		km.ClearLine = []KeyStroke{NewCtrlKeyStroke('u'), NewAltKeyStroke(0, "backspace")} // Alt+backspace conflict

		conflicts := detectConflictsV2(km)
		if len(conflicts) == 0 {
			t.Error("Expected conflicts to be detected")
		}

		// Should detect Alt+backspace conflict
		found := false
		for _, conflict := range conflicts {
			if contains(conflict, "Alt+backspace") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected Alt+backspace conflict to be detected, got: %v", conflicts)
		}
	})
}

// Note: contains() function is defined in keybindings_layers_test.go

// TestProfileTypes tests Profile type constants
func TestProfileTypes(t *testing.T) {
	profiles := []Profile{ProfileDefault, ProfileEmacs, ProfileVi, ProfileReadline}
	expected := []string{"default", "emacs", "vi", "readline"}

	for i, profile := range profiles {
		if string(profile) != expected[i] {
			t.Errorf("Profile constant %d: expected %q, got %q", i, expected[i], string(profile))
		}
	}
}

// TestContextTypes tests Context type constants
func TestContextTypes(t *testing.T) {
	contexts := []Context{ContextGlobal, ContextInput, ContextResults, ContextSearch}
	expected := []string{"global", "input", "results", "search"}

	for i, context := range contexts {
		if string(context) != expected[i] {
			t.Errorf("Context constant %d: expected %q, got %q", i, expected[i], string(context))
		}
	}
}

// TestGetAllContexts tests that all contexts are returned
func TestGetAllContexts(t *testing.T) {
	contexts := GetAllContexts()
	expectedCount := 4

	if len(contexts) != expectedCount {
		t.Errorf("GetAllContexts: expected %d contexts, got %d", expectedCount, len(contexts))
	}

	// Check that all expected contexts are present
	contextMap := make(map[Context]bool)
	for _, ctx := range contexts {
		contextMap[ctx] = true
	}

	expectedContexts := []Context{ContextGlobal, ContextInput, ContextResults, ContextSearch}
	for _, expected := range expectedContexts {
		if !contextMap[expected] {
			t.Errorf("GetAllContexts: missing context %q", expected)
		}
	}
}

// TestKeyBindingProfile tests KeyBindingProfile functionality
func TestKeyBindingProfile(t *testing.T) {
	profile := &KeyBindingProfile{
		Name:        "test",
		Description: "test profile",
		Global: map[string][]KeyStroke{
			"quit": {NewCtrlKeyStroke('c')},
		},
		Contexts: map[Context]map[string][]KeyStroke{
			ContextInput: {
				"delete_word": {NewCtrlKeyStroke('w')},
			},
		},
	}

	// Test GetBinding for existing binding
	if bindings, exists := profile.GetBinding(ContextInput, "delete_word"); !exists {
		t.Error("GetBinding: expected to find delete_word in input context")
	} else if len(bindings) != 1 {
		t.Errorf("GetBinding: expected 1 binding, got %d", len(bindings))
	}

	// Test GetBinding for non-existent binding
	if _, exists := profile.GetBinding(ContextInput, "nonexistent"); exists {
		t.Error("GetBinding: expected not to find nonexistent action")
	}

	// Test GetBinding for non-existent context
	if _, exists := profile.GetBinding(ContextResults, "delete_word"); exists {
		t.Error("GetBinding: expected not to find delete_word in results context")
	}
}

// TestContextualKeyBindingMap tests ContextualKeyBindingMap functionality
func TestContextualKeyBindingMap(t *testing.T) {
	contextual := NewContextualKeyBindingMap(ProfileDefault, "darwin", "iterm")

	// Test initial state
	if contextual.Profile != ProfileDefault {
		t.Errorf("Profile: expected %q, got %q", ProfileDefault, contextual.Profile)
	}
	if contextual.Platform != "darwin" {
		t.Errorf("Platform: expected darwin, got %s", contextual.Platform)
	}
	if contextual.Terminal != "iterm" {
		t.Errorf("Terminal: expected iterm, got %s", contextual.Terminal)
	}

	// Test SetContext and GetContext
	keyMap := DefaultKeyBindingMap()
	contextual.SetContext(ContextInput, keyMap)

	if retrieved, exists := contextual.GetContext(ContextInput); !exists {
		t.Error("GetContext: expected to find input context")
	} else if retrieved != keyMap {
		t.Error("GetContext: retrieved keymap doesn't match set keymap")
	}

	// Test non-existent context
	if _, exists := contextual.GetContext(ContextResults); exists {
		t.Error("GetContext: expected not to find results context")
	}
}

// TestPlatformDetection tests platform detection
func TestPlatformDetection(t *testing.T) {
	platform := DetectPlatform()

	validPlatforms := []string{"darwin", "linux", "windows", "bsd", "unix"}
	valid := false
	for _, p := range validPlatforms {
		if platform == p {
			valid = true
			break
		}
	}

	if !valid {
		t.Errorf("DetectPlatform: returned invalid platform %q", platform)
	}
}

// TestTerminalDetection tests terminal detection
func TestTerminalDetection(t *testing.T) {
	terminal := DetectTerminal()

	// Should return a non-empty string
	if terminal == "" {
		t.Error("DetectTerminal: returned empty string")
	}
}

// TestGetTerminalCapabilities tests terminal capability detection
func TestGetTerminalCapabilities(t *testing.T) {
	capabilities := GetTerminalCapabilities("iterm")

	// iterm should have modern capabilities
	expectedCapabilities := map[string]bool{
		"alt_keys":      true,
		"function_keys": true,
		"mouse":         true,
		"color_256":     true,
		"unicode":       true,
	}

	for capability, expected := range expectedCapabilities {
		if got := capabilities[capability]; got != expected {
			t.Errorf("GetTerminalCapabilities(iterm)[%q]: expected %v, got %v", capability, expected, got)
		}
	}

	// Test dumb terminal
	dumbCaps := GetTerminalCapabilities("dumb")
	for capability, value := range dumbCaps {
		if value {
			t.Errorf("GetTerminalCapabilities(dumb)[%q]: expected false, got true", capability)
		}
	}
}

// TestBuiltinProfiles tests all built-in profile creation
func TestBuiltinProfiles(t *testing.T) {
	profiles := []struct {
		profile Profile
		creator func() *KeyBindingProfile
		name    string
	}{
		{ProfileDefault, CreateDefaultProfile, "Default"},
		{ProfileEmacs, CreateEmacsProfile, "Emacs"},
		{ProfileVi, CreateViProfile, "Vi"},
		{ProfileReadline, CreateReadlineProfile, "Readline"},
	}

	for _, p := range profiles {
		t.Run(string(p.profile), func(t *testing.T) {
			profile := p.creator()

			if profile.Name != p.name {
				t.Errorf("Profile name: expected %q, got %q", p.name, profile.Name)
			}

			if profile.Description == "" {
				t.Error("Profile description should not be empty")
			}

			if profile.Contexts == nil {
				t.Error("Profile contexts should not be nil")
			}

			// Test that input context has basic bindings
			if inputBindings, exists := profile.Contexts[ContextInput]; exists {
				if len(inputBindings) == 0 {
					t.Error("Input context should have keybindings")
				}
			} else {
				t.Error("Profile should have input context bindings")
			}
		})
	}
}

// TestKeyBindingResolver tests the resolver with multiple layers
func TestKeyBindingResolver(t *testing.T) {
	// Create minimal config for testing
	cfg := &config.Config{}
	cfg.Interactive.Profile = "emacs"

	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)

	// Test profile registration
	if profile, exists := resolver.GetProfile(ProfileEmacs); !exists {
		t.Error("Expected emacs profile to be registered")
	} else if profile.Name != "Emacs" {
		t.Errorf("Emacs profile name: expected 'Emacs', got %q", profile.Name)
	}

	// Test resolution for different contexts
	inputMap, err := resolver.Resolve(ProfileEmacs, ContextInput)
	if err != nil {
		t.Errorf("Failed to resolve emacs input context: %v", err)
	}

	if inputMap == nil {
		t.Error("Input map should not be nil")
	}

	// Test contextual resolution
	contextualMap, err := resolver.ResolveContextual(ProfileEmacs)
	if err != nil {
		t.Errorf("Failed to resolve emacs contextual: %v", err)
	}

	if contextualMap == nil {
		t.Error("Contextual map should not be nil")
	}

	// Verify all contexts were resolved
	for _, context := range GetAllContexts() {
		if _, exists := contextualMap.GetContext(context); !exists {
			t.Errorf("Contextual map missing context: %s", context)
		}
	}
}

// TestConfigValidation tests config validation for supported profiles
func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name      string
		profile   string
		shouldErr bool
	}{
		{"valid default profile", "default", false},
		{"valid emacs profile", "emacs", false},
		{"valid vi profile", "vi", false},
		{"valid readline profile", "readline", false},
		{"empty profile", "", false}, // Empty is allowed, defaults to default
		{"invalid profile", "invalid", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := &config.Config{}
			// Set required defaults to make validation pass
			cfg.Default.Branch = "main"
			cfg.Default.Editor = "vim"
			cfg.Behavior.ConfirmDestructive = "simple"
			cfg.Interactive.Profile = test.profile

			err := cfg.Validate()
			if test.shouldErr && err == nil {
				t.Error("Expected validation error but got none")
			}
			if !test.shouldErr && err != nil {
				t.Errorf("Expected no validation error but got: %v", err)
			}
		})
	}
}

// TestContextualKeybindingValidation tests validation of context-specific keybindings
func TestContextualKeybindingValidation(t *testing.T) {
	cfg := &config.Config{}
	// Set required defaults to make validation pass
	cfg.Default.Branch = "main"
	cfg.Default.Editor = "vim"
	cfg.Behavior.ConfirmDestructive = "simple"

	// Test valid context bindings
	cfg.Interactive.Contexts.Input.Keybindings = map[string]interface{}{
		"delete_word": "ctrl+w",
		"clear_line":  []interface{}{"ctrl+u", "ctrl+k"},
	}

	if err := cfg.Validate(); err != nil {
		t.Errorf("Expected valid context bindings to pass validation: %v", err)
	}

	// Test invalid context bindings
	cfg.Interactive.Contexts.Input.Keybindings = map[string]interface{}{
		"delete_word": "invalid_key",
	}

	if err := cfg.Validate(); err == nil {
		t.Error("Expected invalid context bindings to fail validation")
	}
}

// TestProfileDescriptions tests profile description retrieval
func TestProfileDescriptions(t *testing.T) {
	profiles := []Profile{ProfileDefault, ProfileEmacs, ProfileVi, ProfileReadline}

	for _, profile := range profiles {
		description := GetProfileDescription(profile)
		if description == "" {
			t.Errorf("Profile %s should have a description", profile)
		}
		if description == "Unknown profile" {
			t.Errorf("Profile %s should not return 'Unknown profile'", profile)
		}
	}

	// Test unknown profile
	unknownDesc := GetProfileDescription(Profile("unknown"))
	if unknownDesc != "Unknown profile" {
		t.Errorf("Unknown profile should return 'Unknown profile', got %q", unknownDesc)
	}
}

// TestResolverCaching tests that the resolver properly caches results
func TestResolverCaching(t *testing.T) {
	cfg := &config.Config{}
	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)

	// First resolution
	map1, err1 := resolver.Resolve(ProfileDefault, ContextInput)
	if err1 != nil {
		t.Fatalf("First resolve failed: %v", err1)
	}

	// Second resolution (should use cache)
	map2, err2 := resolver.Resolve(ProfileDefault, ContextInput)
	if err2 != nil {
		t.Fatalf("Second resolve failed: %v", err2)
	}

	// Should return the same map (cached)
	if map1 != map2 {
		t.Error("Expected cached result to return same map instance")
	}

	// Test cache clearing
	resolver.ClearCache()

	// Third resolution (after cache clear)
	map3, err3 := resolver.Resolve(ProfileDefault, ContextInput)
	if err3 != nil {
		t.Fatalf("Third resolve failed: %v", err3)
	}

	// Should be a new map instance
	if map1 == map3 {
		t.Error("Expected cache clear to create new map instance")
	}
}

// TestNewRawKeyStroke tests the NewRawKeyStroke constructor
func TestNewRawKeyStroke(t *testing.T) {
	seq := []byte{27, 91, 65} // ESC[A (up arrow)
	ks := NewRawKeyStroke(seq)

	if ks.Kind != KeyStrokeRawSeq {
		t.Errorf("Kind: expected KeyStrokeRawSeq, got %v", ks.Kind)
	}

	if len(ks.Seq) != len(seq) {
		t.Errorf("Sequence length: expected %d, got %d", len(seq), len(ks.Seq))
	}

	for i, b := range seq {
		if ks.Seq[i] != b {
			t.Errorf("Sequence[%d]: expected %d, got %d", i, b, ks.Seq[i])
		}
	}
}

// TestContextManager tests the dynamic context management functionality
func TestContextManager(t *testing.T) {
	cfg := &config.Config{}
	resolver := NewKeyBindingResolver(cfg)
	cm := NewContextManager(resolver)

	// Test initial state
	if cm.GetCurrentContext() != ContextGlobal {
		t.Errorf("Initial context: expected %s, got %s", ContextGlobal, cm.GetCurrentContext())
	}

	// Test entering a context
	cm.EnterContext(ContextInput)
	if cm.GetCurrentContext() != ContextInput {
		t.Errorf("After entering input: expected %s, got %s", ContextInput, cm.GetCurrentContext())
	}

	// Test context stack
	stack := cm.GetContextStack()
	if len(stack) != 1 || stack[0] != ContextGlobal {
		t.Errorf("Stack: expected [%s], got %v", ContextGlobal, stack)
	}

	// Test nested context
	cm.EnterContext(ContextSearch)
	if cm.GetCurrentContext() != ContextSearch {
		t.Errorf("After entering search: expected %s, got %s", ContextSearch, cm.GetCurrentContext())
	}

	// Test exiting context
	returned := cm.ExitContext()
	if returned != ContextInput {
		t.Errorf("Exit context: expected %s, got %s", ContextInput, returned)
	}

	// Test final exit
	cm.ExitContext()
	if cm.GetCurrentContext() != ContextGlobal {
		t.Errorf("Final context: expected %s, got %s", ContextGlobal, cm.GetCurrentContext())
	}

	// Test empty stack exit
	returned = cm.ExitContext()
	if returned != ContextGlobal {
		t.Errorf("Empty stack exit: expected %s, got %s", ContextGlobal, returned)
	}
}

// TestContextCallbacks tests context change callbacks
func TestContextCallbacks(t *testing.T) {
	cfg := &config.Config{}
	resolver := NewKeyBindingResolver(cfg)
	cm := NewContextManager(resolver)

	callbackCalled := false
	var fromContext, toContext Context

	cm.RegisterContextCallback(ContextInput, func(from, to Context) {
		callbackCalled = true
		fromContext = from
		toContext = to
	})

	cm.EnterContext(ContextInput)

	if !callbackCalled {
		t.Error("Callback should have been called")
	}
	if fromContext != ContextGlobal || toContext != ContextInput {
		t.Errorf("Callback args: expected %s->%s, got %s->%s",
			ContextGlobal, ContextInput, fromContext, toContext)
	}
}

// TestPlatformOptimizations tests platform-specific optimizations
func TestPlatformOptimizations(t *testing.T) {
	tests := []struct {
		platform string
		terminal string
		action   string
		expected bool
	}{
		{"darwin", "iterm", "delete_word", true},
		{"linux", "gnome-terminal", "delete_word", true},
		{"windows", "cmd", "paste", true},
		{"freebsd", "xterm", "delete_word", true},
		{"darwin", "terminal", "nonexistent", false},
	}

	for _, test := range tests {
		po := NewPlatformOptimizations(test.platform, test.terminal)
		_, exists := po.GetOptimizedBindings(test.action)

		if exists != test.expected {
			t.Errorf("Platform %s, terminal %s, action %s: expected %v, got %v",
				test.platform, test.terminal, test.action, test.expected, exists)
		}
	}
}

// TestRuntimeProfileSwitcher tests runtime profile switching
func TestRuntimeProfileSwitcher(t *testing.T) {
	cfg := &config.Config{}
	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)

	cm := NewContextManager(resolver)
	rps := NewRuntimeProfileSwitcher(resolver, cm)

	// Test initial profile
	if rps.GetCurrentProfile() != ProfileDefault {
		t.Errorf("Initial profile: expected %s, got %s", ProfileDefault, rps.GetCurrentProfile())
	}

	// Test switching to valid profile
	err := rps.SwitchProfile(ProfileEmacs)
	if err != nil {
		t.Errorf("Switch to emacs: unexpected error %v", err)
	}
	if rps.GetCurrentProfile() != ProfileEmacs {
		t.Errorf("After switch: expected %s, got %s", ProfileEmacs, rps.GetCurrentProfile())
	}

	// Test switching to invalid profile
	err = rps.SwitchProfile(Profile("invalid"))
	if err == nil {
		t.Error("Switch to invalid profile: expected error")
	}

	// Test profile cycling
	if err := rps.SwitchProfile(ProfileDefault); err != nil {
		t.Fatalf("Switch to default: unexpected error %v", err)
	}
	err = rps.CycleProfile()
	if err != nil {
		t.Errorf("Cycle profile: unexpected error %v", err)
	}
	if rps.GetCurrentProfile() != ProfileEmacs {
		t.Errorf("After cycle: expected %s, got %s", ProfileEmacs, rps.GetCurrentProfile())
	}
}

// TestProfileSwitchCallbacks tests profile switch callbacks
func TestProfileSwitchCallbacks(t *testing.T) {
	cfg := &config.Config{}
	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)

	cm := NewContextManager(resolver)
	rps := NewRuntimeProfileSwitcher(resolver, cm)

	callbackCalled := false
	var fromProfile, toProfile Profile

	rps.RegisterSwitchCallback(func(from, to Profile) {
		callbackCalled = true
		fromProfile = from
		toProfile = to
	})

	if err := rps.SwitchProfile(ProfileVi); err != nil {
		t.Fatalf("Switch to vi: unexpected error %v", err)
	}

	if !callbackCalled {
		t.Error("Switch callback should have been called")
	}
	if fromProfile != ProfileDefault || toProfile != ProfileVi {
		t.Errorf("Callback args: expected %s->%s, got %s->%s",
			ProfileDefault, ProfileVi, fromProfile, toProfile)
	}
}

// TestHotConfigReloader tests the hot config reload functionality
func TestHotConfigReloader(t *testing.T) {
	cfg := &config.Config{}
	resolver := NewKeyBindingResolver(cfg)

	// Use a temp file path for testing
	configPath := "/tmp/test_config.yaml"
	hcr := NewHotConfigReloader(configPath, resolver)

	// Test not watching initially
	if hcr.watching {
		t.Error("Should not be watching initially")
	}

	// Test callback registration
	hcr.RegisterReloadCallback(func(*config.Config) {})

	// Test that callback was registered (callbackCalled would be set on reload)
	if len(hcr.reloadCallbacks) != 1 {
		t.Error("Callback should be registered")
	}

	// Note: Full file watching test would require actual file creation
	// This tests the interface and basic functionality
}

// TestContextTransitionAnimator tests context transition animations
func TestContextTransitionAnimator(t *testing.T) {
	cta := NewContextTransitionAnimator()

	// Test initial state
	if !cta.enabled {
		t.Error("Animator should be enabled initially")
	}
	if cta.style != "highlight" {
		t.Errorf("Style: expected 'highlight', got %s", cta.style)
	}
	if cta.duration != 200*time.Millisecond {
		t.Errorf("Duration: expected 200ms, got %v", cta.duration)
	}

	// Test style setting
	cta.SetStyle("fade")
	if cta.style != "fade" {
		t.Errorf("After setting style: expected 'fade', got %s", cta.style)
	}

	// Test duration setting
	newDuration := 500 * time.Millisecond
	cta.SetDuration(newDuration)
	if cta.duration != newDuration {
		t.Errorf("After setting duration: expected %v, got %v", newDuration, cta.duration)
	}

	// Test enable/disable
	cta.Disable()
	if cta.enabled {
		t.Error("Animator should be disabled")
	}

	cta.Enable()
	if !cta.enabled {
		t.Error("Animator should be enabled")
	}

	// Test animation registration
	callbackCalled := false
	cta.RegisterAnimation(func(from, to Context) {
		callbackCalled = true
	})

	if len(cta.animations) != 1 {
		t.Errorf("Animations: expected 1, got %d", len(cta.animations))
	}

	// Use the callback variable to avoid unused warning
	_ = callbackCalled
}

// TestPlatformSpecificBindings tests platform-specific keybinding creation
func TestPlatformSpecificBindings(t *testing.T) {
	// Test macOS optimizations
	po := NewPlatformOptimizations("darwin", "iterm")

	deleteWordBindings, exists := po.GetOptimizedBindings("delete_word")
	if !exists {
		t.Error("macOS should have delete_word optimization")
	}
	if len(deleteWordBindings) < 1 {
		t.Error("macOS delete_word should have at least one binding")
	}

	// Test Linux optimizations
	po = NewPlatformOptimizations("linux", "gnome-terminal")

	deleteWordBindings, exists = po.GetOptimizedBindings("delete_word")
	if !exists {
		t.Error("Linux should have delete_word optimization")
	}
	if len(deleteWordBindings) == 0 {
		t.Error("Linux delete_word should not be empty")
	}

	// Test Windows optimizations
	po = NewPlatformOptimizations("windows", "cmd")

	pasteBindings, exists := po.GetOptimizedBindings("paste")
	if !exists {
		t.Error("Windows should have paste optimization")
	}
	if len(pasteBindings) < 1 {
		t.Error("Windows paste should have at least one binding")
	}

	// Test Unix fallback
	po = NewPlatformOptimizations("freebsd", "xterm")

	clearLineBindings, exists := po.GetOptimizedBindings("clear_line")
	if !exists {
		t.Error("Unix should have clear_line optimization")
	}
	if len(clearLineBindings) < 1 {
		t.Error("Unix clear_line should have at least one binding")
	}
}

// TestIntegratedWorkflow tests the complete extended workflow
func TestIntegratedWorkflow(t *testing.T) {
	// Set up complete system
	cfg := &config.Config{}
	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)

	cm := NewContextManager(resolver)
	rps := NewRuntimeProfileSwitcher(resolver, cm)
	cta := NewContextTransitionAnimator()

	// Test integrated workflow

	// Switch profile
	err := rps.SwitchProfile(ProfileEmacs)
	if err != nil {
		t.Errorf("Profile switch failed: %v", err)
	}

	// Enter context with animation
	cta.AnimateTransition(cm.GetCurrentContext(), ContextInput)
	cm.EnterContext(ContextInput)

	// Test resolution with new profile and context
	keyMap, err := resolver.Resolve(rps.GetCurrentProfile(), cm.GetCurrentContext())
	if err != nil {
		t.Errorf("Resolution failed: %v", err)
	}
	if keyMap == nil {
		t.Error("Resolved keymap should not be nil")
	}

	// Cycle profile
	err = rps.CycleProfile()
	if err != nil {
		t.Errorf("Profile cycle failed: %v", err)
	}

	// Exit context
	cta.AnimateTransition(cm.GetCurrentContext(), ContextGlobal)
	cm.ExitContext()
}

// TestWorkflowPerformance tests performance of the combined workflow
func TestWorkflowPerformance(t *testing.T) {
	cfg := &config.Config{}
	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)

	cm := NewContextManager(resolver)

	// Test context switching performance
	start := time.Now()
	for i := 0; i < 1000; i++ {
		cm.EnterContext(ContextInput)
		cm.ExitContext()
	}
	contextSwitchTime := time.Since(start)

	if contextSwitchTime > 10*time.Millisecond {
		t.Errorf("Context switching too slow: %v > 10ms", contextSwitchTime)
	}

	// Test profile switching performance
	rps := NewRuntimeProfileSwitcher(resolver, cm)

	start = time.Now()
	for i := 0; i < 100; i++ {
		if err := rps.CycleProfile(); err != nil {
			t.Fatalf("profile cycle failed at iteration %d: %v", i, err)
		}
	}
	profileSwitchTime := time.Since(start)

	if profileSwitchTime > 100*time.Millisecond {
		t.Errorf("Profile switching too slow: %v > 100ms", profileSwitchTime)
	}
}

func TestKeybindingExporter_ExportFull(t *testing.T) {
	cfg := &config.Config{}
	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)

	exporter := NewKeybindingExporter(resolver)

	// Test full export with emacs profile
	opts := ExportOptions{
		Profile:     ProfileEmacs,
		DeltaMode:   false,
		IncludeMeta: true,
		Format:      "yaml",
	}

	result, err := exporter.Export(opts)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Verify basic structure
	if result.Profile != "emacs" {
		t.Errorf("Expected profile 'emacs', got '%s'", result.Profile)
	}
	if result.Metadata.Platform == "" {
		t.Error("Platform should be detected and set")
	}
	if result.Metadata.Terminal == "" {
		t.Error("Terminal should be detected and set")
	}

	// Verify some expected keybindings exist
	if _, exists := result.Keybindings["move_to_beginning"]; !exists {
		t.Error("Expected 'move_to_beginning' keybinding not found")
	}
	if _, exists := result.Keybindings["move_to_end"]; !exists {
		t.Error("Expected 'move_to_end' keybinding not found")
	}

	// Verify metadata is populated
	if result.Metadata.ExportedAt.IsZero() {
		t.Error("ExportedAt timestamp should be set")
	}
	if result.Metadata.Version == "" {
		t.Error("Version should be set")
	}
}

func TestKeybindingExporter_ExportDelta(t *testing.T) {
	cfg := &config.Config{}
	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)

	exporter := NewKeybindingExporter(resolver)

	// Test delta export
	opts := ExportOptions{
		Profile:     ProfileEmacs,
		DeltaMode:   true,
		IncludeMeta: true,
		Format:      "yaml",
	}

	result, err := exporter.Export(opts)
	if err != nil {
		t.Fatalf("Delta export failed: %v", err)
	}

	// Delta export should only contain differences
	// Since this is fresh profile, should be minimal
	if len(result.Keybindings) > 10 {
		t.Errorf("Delta export should contain minimal differences, got %d bindings", len(result.Keybindings))
	}

	// Verify delta metadata
	if result.Metadata.DeltaFrom != "emacs" {
		t.Errorf("Expected delta_from 'emacs', got '%s'", result.Metadata.DeltaFrom)
	}
}

func TestKeybindingImporter_ValidateAndPreview(t *testing.T) {
	cfg := &config.Config{}
	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)

	importer := NewKeybindingImporter(resolver)

	// Create test YAML content
	testYAML := `
profile: emacs
keybindings:
  move_to_beginning: "ctrl+a"
  move_to_end: "ctrl+e"
  delete_word: "alt+d"
contexts:
  input:
    keybindings:
      complete: "tab"
metadata:
  version: "5.0.0"
  platform: "darwin"
`

	// Test validation
	opts := ImportOptions{
		Data:       []byte(testYAML),
		DryRun:     true,
		MergeMode:  "merge",
		BackupPath: "",
	}

	err := importer.Import(opts)
	if err != nil {
		t.Fatalf("Valid YAML should pass validation: %v", err)
	}
}

func TestKeybindingImporter_InvalidYAML(t *testing.T) {
	cfg := &config.Config{}
	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)
	importer := NewKeybindingImporter(resolver)

	// Test invalid YAML syntax
	invalidYAML := `
profile: emacs
keybindings:
  move_to_beginning: "ctrl+a"
  - invalid: syntax
`

	opts := ImportOptions{
		Data:       []byte(invalidYAML),
		DryRun:     true,
		MergeMode:  "merge",
		BackupPath: "",
	}

	err := importer.Import(opts)
	if err == nil {
		t.Error("Invalid YAML should fail validation")
	}
}

func TestKeybindingImporter_Import(t *testing.T) {
	cfg := &config.Config{}
	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)

	importer := NewKeybindingImporter(resolver)

	testYAML := `
profile: emacs
keybindings:
  move_to_beginning: "ctrl+home"
  delete_word: "ctrl+backspace"
contexts:
  input:
    keybindings:
      complete: "ctrl+space"
`

	// Test import
	opts := ImportOptions{
		Data:       []byte(testYAML),
		DryRun:     false,
		MergeMode:  "merge",
		BackupPath: "",
	}

	err := importer.Import(opts)
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	// Verify imported bindings exist
	bindings := resolver.GetEffectiveKeybindings(ProfileEmacs, ContextInput)
	if len(bindings) == 0 {
		t.Error("Should have effective bindings after import")
	}
}

func TestShowKeysCommand_Execute(t *testing.T) {
	cfg := &config.Config{}
	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)

	cmd := NewShowKeysCommand(resolver)

	// Test show keys for emacs profile - just verify it doesn't error
	err := cmd.Execute(ProfileEmacs, ContextInput, "full")
	if err != nil {
		t.Fatalf("ShowKeys failed: %v", err)
	}

	// Test compact format
	err = cmd.Execute(ProfileEmacs, ContextInput, "compact")
	if err != nil {
		t.Fatalf("ShowKeys compact failed: %v", err)
	}
}

func TestDebugKeysCommand_FormatKeySequence(t *testing.T) {
	cmd := NewDebugKeysCommand("")

	testCases := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "Simple character",
			input:    []byte{'a'},
			expected: "a (0x61)",
		},
		{
			name:     "Ctrl+A",
			input:    []byte{1},
			expected: "Ctrl+A (0x01)",
		},
		{
			name:     "Escape",
			input:    []byte{27},
			expected: "Esc (0x1b)",
		},
		{
			name:     "Tab",
			input:    []byte{9},
			expected: "Tab (0x09)",
		},
		{
			name:     "Arrow key sequence",
			input:    []byte{27, 91, 65},
			expected: "↑ (0x1b 0x5b 0x41)",
		},
		{
			name:     "Function key F1",
			input:    []byte{27, 79, 80},
			expected: "F1 (0x1b 0x4f 0x50)",
		},
		{
			name:     "Multi-byte sequence",
			input:    []byte{27, 91, 49, 59, 50, 65},
			expected: "Shift+↑ (0x1b 0x5b 0x31 0x3b 0x32 0x41)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := cmd.formatKeySequence(tc.input)
			if !strings.Contains(result, tc.expected) {
				t.Errorf("Expected output to contain '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestExportImportRoundTrip(t *testing.T) {
	cfg := &config.Config{}
	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)

	exporter := NewKeybindingExporter(resolver)

	// Export current configuration
	exportOpts := ExportOptions{
		Profile:     ProfileEmacs,
		DeltaMode:   false,
		IncludeMeta: true,
		Format:      "yaml",
	}

	exported, err := exporter.Export(exportOpts)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Convert to YAML bytes for import
	yamlData, err := yaml.Marshal(exported)
	if err != nil {
		t.Fatalf("Failed to marshal export to YAML: %v", err)
	}

	// Create a new resolver to test import
	newResolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(newResolver)
	newImporter := NewKeybindingImporter(newResolver)

	// Import the exported configuration
	importOpts := ImportOptions{
		Data:       yamlData,
		DryRun:     false,
		MergeMode:  "replace",
		BackupPath: "",
	}

	err = newImporter.Import(importOpts)
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	// Compare original and imported configurations
	originalBindings := resolver.GetEffectiveKeybindings(ProfileEmacs, ContextInput)
	importedBindings := newResolver.GetEffectiveKeybindings(ProfileEmacs, ContextInput)

	// Check that both have bindings
	if len(originalBindings) == 0 {
		t.Error("Original resolver should have bindings")
	}
	if len(importedBindings) == 0 {
		t.Error("Imported resolver should have bindings")
	}
}

func TestFileOperations(t *testing.T) {
	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "ggc_test_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	cfg := &config.Config{}
	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)

	exporter := NewKeybindingExporter(resolver)
	importer := NewKeybindingImporter(resolver)

	// Test export to file
	exportFile := filepath.Join(tempDir, "test_export.yaml")
	exportOpts := ExportOptions{
		Profile:     ProfileEmacs,
		DeltaMode:   false,
		IncludeMeta: true,
		Format:      "yaml",
		OutputFile:  exportFile,
	}

	exported, err := exporter.Export(exportOpts)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Convert to YAML and write to file
	yamlData, err := yaml.Marshal(exported)
	if err != nil {
		t.Fatalf("Failed to marshal export: %v", err)
	}

	err = os.WriteFile(exportFile, yamlData, 0644)
	if err != nil {
		t.Fatalf("Failed to write export file: %v", err)
	}

	// Test import from file
	importData, err := os.ReadFile(exportFile)
	if err != nil {
		t.Fatalf("Failed to read export file: %v", err)
	}

	importOpts := ImportOptions{
		Data:       importData,
		DryRun:     true,
		MergeMode:  "merge",
		BackupPath: "",
	}

	err = importer.Import(importOpts)
	if err != nil {
		t.Fatalf("Imported file should be valid: %v", err)
	}
}

func TestErrorHandling(t *testing.T) {
	cfg := &config.Config{}
	resolver := NewKeyBindingResolver(cfg)

	// Test with unregistered profile
	exporter := NewKeybindingExporter(resolver)

	exportOpts := ExportOptions{
		Profile:     Profile("nonexistent"),
		DeltaMode:   false,
		IncludeMeta: true,
		Format:      "yaml",
	}

	_, err := exporter.Export(exportOpts)
	if err == nil {
		t.Error("Export should fail with nonexistent profile")
	}

	// Test show keys with invalid context
	showCmd := NewShowKeysCommand(resolver)
	err = showCmd.Execute(Profile("nonexistent"), ContextInput, "full")
	if err == nil {
		t.Error("ShowKeys should fail with nonexistent profile")
	}
}

func TestPerformance(t *testing.T) {
	cfg := &config.Config{}
	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)

	exporter := NewKeybindingExporter(resolver)

	// Measure export performance
	start := time.Now()
	exportOpts := ExportOptions{
		Profile:     ProfileEmacs,
		DeltaMode:   false,
		IncludeMeta: true,
		Format:      "yaml",
	}

	_, err := exporter.Export(exportOpts)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Export should complete within reasonable time
	if duration > time.Second {
		t.Errorf("Export took too long: %v", duration)
	}

	// Measure show keys performance
	showCmd := NewShowKeysCommand(resolver)

	start = time.Now()
	err = showCmd.Execute(ProfileEmacs, ContextInput, "full")
	duration = time.Since(start)

	if err != nil {
		t.Fatalf("ShowKeys failed: %v", err)
	}

	if duration > 100*time.Millisecond {
		t.Errorf("ShowKeys took too long: %v", duration)
	}
}
