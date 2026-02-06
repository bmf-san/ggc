package keybindings

// ContextualKeyBindingMap holds keybinding maps for multiple contexts
type ContextualKeyBindingMap struct {
	Profile  Profile                    // The resolved profile
	Platform string                     // Platform (darwin/linux/windows)
	Terminal string                     // Terminal type (xterm/tmux/etc)
	Contexts map[Context]*KeyBindingMap // Resolved keybindings per context
}

// NewContextualKeyBindingMap creates a new contextual map
func NewContextualKeyBindingMap(profile Profile, platform, terminal string) *ContextualKeyBindingMap {
	return &ContextualKeyBindingMap{
		Profile:  profile,
		Platform: platform,
		Terminal: terminal,
		Contexts: make(map[Context]*KeyBindingMap),
	}
}

// GetContext returns the KeyBindingMap for a specific context
func (ckm *ContextualKeyBindingMap) GetContext(context Context) (*KeyBindingMap, bool) {
	keyMap, exists := ckm.Contexts[context]
	return keyMap, exists
}

// SetContext sets the KeyBindingMap for a specific context
func (ckm *ContextualKeyBindingMap) SetContext(context Context, keyMap *KeyBindingMap) {
	if ckm.Contexts == nil {
		ckm.Contexts = make(map[Context]*KeyBindingMap)
	}
	ckm.Contexts[context] = keyMap
}
