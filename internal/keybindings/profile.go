package keybindings

// Profile represents built-in keybinding profiles that users can select
type Profile string

// Supported keybinding profiles.
const (
	ProfileDefault  Profile = "default"  // Current default behavior (backward compatible)
	ProfileEmacs    Profile = "emacs"    // Emacs-style bindings (Ctrl-based, modeless)
	ProfileVi       Profile = "vi"       // Vi-style bindings (modal concepts adapted for CLI)
	ProfileReadline Profile = "readline" // GNU Readline standard bindings
)

// String returns the string representation of a Profile
func (p Profile) String() string {
	return string(p)
}

// IsValid checks if a Profile value is valid
func (p Profile) IsValid() bool {
	switch p {
	case ProfileDefault, ProfileEmacs, ProfileVi, ProfileReadline:
		return true
	default:
		return false
	}
}

// Context represents different UI states that can have specific keybindings
type Context string

// Available contexts for interactive UI states.
const (
	ContextGlobal  Context = "global"  // Always active (reserved keys like Ctrl+C)
	ContextInput   Context = "input"   // When typing/editing the search query
	ContextResults Context = "results" // When navigating through filtered results
	ContextSearch  Context = "search"  // When fuzzy search is active (combines input + results)
)

// String returns the string representation of a Context
func (c Context) String() string {
	return string(c)
}

// IsValid checks if a Context value is valid
func (c Context) IsValid() bool {
	switch c {
	case ContextGlobal, ContextInput, ContextResults, ContextSearch:
		return true
	default:
		return false
	}
}

// GetAllProfiles returns a list of all valid profiles
func GetAllProfiles() []Profile {
	return []Profile{ProfileDefault, ProfileEmacs, ProfileVi, ProfileReadline}
}

// GetAllContexts returns a list of all valid contexts
func GetAllContexts() []Context {
	return []Context{ContextGlobal, ContextInput, ContextResults, ContextSearch}
}

// KeyBindingProfile defines keybindings for a complete profile
type KeyBindingProfile struct {
	Name        string                             // Profile name (e.g., "emacs")
	Description string                             // Human-readable description
	Global      map[string][]KeyStroke             // Global keybindings (always active)
	Contexts    map[Context]map[string][]KeyStroke // Context-specific keybindings
}

// NewKeyBindingProfile creates a new profile with initialized maps
func NewKeyBindingProfile(name, description string) *KeyBindingProfile {
	return &KeyBindingProfile{
		Name:        name,
		Description: description,
		Global:      make(map[string][]KeyStroke),
		Contexts:    make(map[Context]map[string][]KeyStroke),
	}
}

// SetGlobalBinding sets a global keybinding (active in all contexts)
func (kbp *KeyBindingProfile) SetGlobalBinding(action string, keystrokes []KeyStroke) {
	if kbp.Global == nil {
		kbp.Global = make(map[string][]KeyStroke)
	}
	kbp.Global[action] = keystrokes
}

// SetContextBinding sets a context-specific keybinding
func (kbp *KeyBindingProfile) SetContextBinding(context Context, action string, keystrokes []KeyStroke) {
	if kbp.Contexts == nil {
		kbp.Contexts = make(map[Context]map[string][]KeyStroke)
	}
	if kbp.Contexts[context] == nil {
		kbp.Contexts[context] = make(map[string][]KeyStroke)
	}
	kbp.Contexts[context][action] = keystrokes
}

// GetBinding returns the keybinding for an action in a specific context
// Falls back to global bindings if not found in context
func (kbp *KeyBindingProfile) GetBinding(context Context, action string) ([]KeyStroke, bool) {
	// Try context-specific first
	if contextMap, exists := kbp.Contexts[context]; exists {
		if keystrokes, exists := contextMap[action]; exists {
			return keystrokes, true
		}
	}

	// Fall back to global
	if keystrokes, exists := kbp.Global[action]; exists {
		return keystrokes, true
	}

	return nil, false
}

// GetAllActions returns all action names defined in this profile
func (kbp *KeyBindingProfile) GetAllActions() []string {
	actionSet := make(map[string]bool)

	// Add global actions
	for action := range kbp.Global {
		actionSet[action] = true
	}

	// Add context-specific actions
	for _, contextMap := range kbp.Contexts {
		for action := range contextMap {
			actionSet[action] = true
		}
	}

	// Convert to slice
	actions := make([]string, 0, len(actionSet))
	for action := range actionSet {
		actions = append(actions, action)
	}

	return actions
}

// Clone creates a deep copy of the profile
func (kbp *KeyBindingProfile) Clone() *KeyBindingProfile {
	clone := NewKeyBindingProfile(kbp.Name, kbp.Description)

	// Clone global bindings
	for action, keystrokes := range kbp.Global {
		clonedKeystrokes := make([]KeyStroke, len(keystrokes))
		copy(clonedKeystrokes, keystrokes)
		clone.Global[action] = clonedKeystrokes
	}

	// Clone context bindings
	for context, contextMap := range kbp.Contexts {
		clone.Contexts[context] = make(map[string][]KeyStroke)
		for action, keystrokes := range contextMap {
			clonedKeystrokes := make([]KeyStroke, len(keystrokes))
			copy(clonedKeystrokes, keystrokes)
			clone.Contexts[context][action] = clonedKeystrokes
		}
	}

	return clone
}
