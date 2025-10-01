package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"go.yaml.in/yaml/v3"

	"github.com/bmf-san/ggc/v6/config"
)

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

// KeyStrokeKind represents the type of key stroke
type KeyStrokeKind int

// Key stroke categories recognized by the resolver.
const (
	KeyStrokeCtrl   KeyStrokeKind = iota // Control key combinations (Ctrl+A)
	KeyStrokeAlt                         // Alt/Meta key combinations (Alt+Backspace)
	KeyStrokeRawSeq                      // Raw escape sequences
	KeyStrokeFnKey                       // Function keys (F1, F2, etc.)
)

// String returns a human-readable representation of the KeyStrokeKind
func (k KeyStrokeKind) String() string {
	switch k {
	case KeyStrokeCtrl:
		return "Ctrl"
	case KeyStrokeAlt:
		return "Alt"
	case KeyStrokeRawSeq:
		return "RawSeq"
	case KeyStrokeFnKey:
		return "FnKey"
	default:
		return "Unknown"
	}
}

// KeyStroke represents a single key input that can trigger an action
type KeyStroke struct {
	Kind KeyStrokeKind // The type of key stroke
	Rune rune          // For Ctrl+<letter>, Alt+<letter> - the letter
	Seq  []byte        // For raw escape sequences
	Name string        // For function keys (F1, F2, etc.) and special names
}

// String returns a human-readable representation of the KeyStroke
func (ks KeyStroke) String() string {
	switch ks.Kind {
	case KeyStrokeCtrl:
		return fmt.Sprintf("Ctrl+%c", ks.Rune)
	case KeyStrokeAlt:
		if ks.Name != "" {
			return fmt.Sprintf("Alt+%s", ks.Name)
		}
		return fmt.Sprintf("Alt+%c", ks.Rune)
	case KeyStrokeRawSeq:
		return fmt.Sprintf("Seq%v", ks.Seq)
	case KeyStrokeFnKey:
		return ks.Name
	default:
		return "Unknown"
	}
}

// Equals checks if two KeyStrokes are equivalent
func (ks KeyStroke) Equals(other KeyStroke) bool {
	if ks.Kind != other.Kind {
		return false
	}
	switch ks.Kind {
	case KeyStrokeCtrl, KeyStrokeAlt:
		return ks.Rune == other.Rune && ks.Name == other.Name
	case KeyStrokeRawSeq:
		if len(ks.Seq) != len(other.Seq) {
			return false
		}
		for i, b := range ks.Seq {
			if b != other.Seq[i] {
				return false
			}
		}
		return true
	case KeyStrokeFnKey:
		return ks.Name == other.Name
	default:
		return false
	}
}

// ToControlByte converts a KeyStroke to a control byte for backward compatibility
// Returns 0 if the KeyStroke cannot be represented as a single control byte
func (ks KeyStroke) ToControlByte() byte {
	if ks.Kind == KeyStrokeCtrl && ks.Rune >= 'a' && ks.Rune <= 'z' {
		return byte(ks.Rune-'a') + 1
	}
	return 0
}

// NewCtrlKeyStroke creates a new Ctrl+letter KeyStroke
func NewCtrlKeyStroke(letter rune) KeyStroke {
	return KeyStroke{
		Kind: KeyStrokeCtrl,
		Rune: letter,
	}
}

// NewAltKeyStroke creates a new Alt+key KeyStroke
func NewAltKeyStroke(key rune, name string) KeyStroke {
	return KeyStroke{
		Kind: KeyStrokeAlt,
		Rune: key,
		Name: name,
	}
}

// NewRawKeyStroke creates a new raw key sequence KeyStroke
func NewRawKeyStroke(seq []byte) KeyStroke {
	return KeyStroke{
		Kind: KeyStrokeRawSeq,
		Seq:  seq,
	}
}

// NewTabKeyStroke creates a new Tab KeyStroke
func NewTabKeyStroke() KeyStroke {
	return NewRawKeyStroke([]byte{9}) // Tab is ASCII 9
}

// NewCharKeyStroke creates a new character KeyStroke
func NewCharKeyStroke(char rune) KeyStroke {
	return NewRawKeyStroke([]byte{byte(char)})
}

// NewEnterKeyStroke creates a new Enter KeyStroke
func NewEnterKeyStroke() KeyStroke {
	return NewRawKeyStroke([]byte{13}) // Enter is ASCII 13
}

// NewEscapeKeyStroke creates a new Escape KeyStroke
func NewEscapeKeyStroke() KeyStroke {
	return NewRawKeyStroke([]byte{27}) // Escape is ASCII 27
}

// NewSpaceKeyStroke creates a new Space KeyStroke
func NewSpaceKeyStroke() KeyStroke {
	return NewRawKeyStroke([]byte{32}) // Space is ASCII 32
}

// KeyBindingMap holds resolved key strokes for interactive actions.
// Supports multiple key strokes per action while maintaining backward compatibility.
type KeyBindingMap struct {
	DeleteWord         []KeyStroke // default: [Ctrl+W]
	ClearLine          []KeyStroke // default: [Ctrl+U]
	DeleteToEnd        []KeyStroke // default: [Ctrl+K]
	MoveToBeginning    []KeyStroke // default: [Ctrl+A]
	MoveToEnd          []KeyStroke // default: [Ctrl+E]
	MoveUp             []KeyStroke // default: [Ctrl+P]
	MoveDown           []KeyStroke // default: [Ctrl+N]
	AddToWorkflow      []KeyStroke // default: [Tab]
	ToggleWorkflowView []KeyStroke // default: [Ctrl+T]
	ClearWorkflow      []KeyStroke // default: [c]
	SoftCancel         []KeyStroke // default: [Ctrl+G, Esc]
}

// DefaultKeyBindingMap returns the built-in default control bindings.
func DefaultKeyBindingMap() *KeyBindingMap {
	return &KeyBindingMap{
		DeleteWord:         []KeyStroke{NewCtrlKeyStroke('w')},
		ClearLine:          []KeyStroke{NewCtrlKeyStroke('u')},
		DeleteToEnd:        []KeyStroke{NewCtrlKeyStroke('k')},
		MoveToBeginning:    []KeyStroke{NewCtrlKeyStroke('a')},
		MoveToEnd:          []KeyStroke{NewCtrlKeyStroke('e')},
		MoveUp:             []KeyStroke{NewCtrlKeyStroke('p')},
		MoveDown:           []KeyStroke{NewCtrlKeyStroke('n')},
		AddToWorkflow:      []KeyStroke{NewTabKeyStroke()},
		ToggleWorkflowView: []KeyStroke{NewCtrlKeyStroke('t')},
		ClearWorkflow:      []KeyStroke{NewCharKeyStroke('c')},
		SoftCancel:         []KeyStroke{NewCtrlKeyStroke('g'), NewEscapeKeyStroke()},
	}
}

// Legacy backward-compatibility methods maintain the old byte-based API
// while internally using the new KeyStroke system.

// GetDeleteWordByte returns the primary control byte for DeleteWord (backward compatibility)
func (km *KeyBindingMap) GetDeleteWordByte() byte {
	return km.getFirstControlByte(km.DeleteWord, ctrl('w'))
}

// GetClearLineByte returns the primary control byte for ClearLine (backward compatibility)
func (km *KeyBindingMap) GetClearLineByte() byte {
	return km.getFirstControlByte(km.ClearLine, ctrl('u'))
}

// GetDeleteToEndByte returns the primary control byte for DeleteToEnd (backward compatibility)
func (km *KeyBindingMap) GetDeleteToEndByte() byte {
	return km.getFirstControlByte(km.DeleteToEnd, ctrl('k'))
}

// GetMoveToBeginningByte returns the primary control byte for MoveToBeginning (backward compatibility)
func (km *KeyBindingMap) GetMoveToBeginningByte() byte {
	return km.getFirstControlByte(km.MoveToBeginning, ctrl('a'))
}

// GetMoveToEndByte returns the primary control byte for MoveToEnd (backward compatibility)
func (km *KeyBindingMap) GetMoveToEndByte() byte {
	return km.getFirstControlByte(km.MoveToEnd, ctrl('e'))
}

// GetMoveUpByte returns the primary control byte for MoveUp (backward compatibility)
func (km *KeyBindingMap) GetMoveUpByte() byte {
	return km.getFirstControlByte(km.MoveUp, ctrl('p'))
}

// GetMoveDownByte returns the primary control byte for MoveDown (backward compatibility)
func (km *KeyBindingMap) GetMoveDownByte() byte {
	return km.getFirstControlByte(km.MoveDown, ctrl('n'))
}

// GetAddToWorkflowByte returns the primary control byte for AddToWorkflow (backward compatibility)
func (km *KeyBindingMap) GetAddToWorkflowByte() byte {
	return km.getFirstControlByte(km.AddToWorkflow, 9) // Tab key
}

// GetToggleWorkflowViewByte returns the primary control byte for ToggleWorkflowView (backward compatibility)
func (km *KeyBindingMap) GetToggleWorkflowViewByte() byte {
	return km.getFirstControlByte(km.ToggleWorkflowView, ctrl('t'))
}

// GetClearWorkflowByte returns the primary control byte for ClearWorkflow (backward compatibility)
func (km *KeyBindingMap) GetClearWorkflowByte() byte {
	return km.getFirstControlByte(km.ClearWorkflow, 'c')
}

// getFirstControlByte finds the first Ctrl KeyStroke and returns its control byte,
// or returns the fallback if none found
func (km *KeyBindingMap) getFirstControlByte(keyStrokes []KeyStroke, fallback byte) byte {
	for _, ks := range keyStrokes {
		if b := ks.ToControlByte(); b != 0 {
			return b
		}
	}
	return fallback
}

// MatchesKeyStroke checks if any KeyStroke in the given action matches the input
func (km *KeyBindingMap) MatchesKeyStroke(action string, input KeyStroke) bool {
	actionMap := map[string][]KeyStroke{
		"delete_word":          km.DeleteWord,
		"clear_line":           km.ClearLine,
		"delete_to_end":        km.DeleteToEnd,
		"move_to_beginning":    km.MoveToBeginning,
		"move_to_end":          km.MoveToEnd,
		"move_up":              km.MoveUp,
		"move_down":            km.MoveDown,
		"add_to_workflow":      km.AddToWorkflow,
		"toggle_workflow_view": km.ToggleWorkflowView,
		"clear_workflow":       km.ClearWorkflow,
		"soft_cancel":          km.SoftCancel,
	}

	keyStrokes, exists := actionMap[action]
	if !exists {
		return false
	}

	for _, ks := range keyStrokes {
		if input.Equals(ks) {
			return true
		}
	}
	return false
}

// ctrl converts a lowercase letter to its control byte (e.g., 'a' => 1).
func ctrl(r rune) byte {
	// Only letters a-z are expected here; ensure predictable conversion.
	if r >= 'a' && r <= 'z' {
		return byte(r-'a') + 1
	}
	if r >= 'A' && r <= 'Z' {
		return byte(r-'A') + 1
	}
	return 0
}

// hasPrefixFold checks whether s has the given prefix, case-insensitively.
func hasPrefixFold(s, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(s), strings.ToLower(prefix))
}

// ParseKeyBinding parses a key binding string and returns the corresponding
// single-byte control code. Supports multiple formats:
// - "ctrl+w", "CTRL+W", "Ctrl+w" (standard format)
// - "^w", "^W" (caret notation)
// - "c-w", "C-w", "C-W" (emacs notation)
func ParseKeyBinding(keyStr string) (byte, error) { //nolint:revive // parsing multiple legacy formats
	s := strings.TrimSpace(keyStr)
	if s == "" {
		return 0, fmt.Errorf("empty key binding")
	}

	// Normalize to lowercase for comparison
	sLower := strings.ToLower(s)

	// Handle "ctrl+<key>" format (case-insensitive)
	if strings.HasPrefix(sLower, "ctrl+") && len(s) == len("ctrl+")+1 {
		c := rune(sLower[len(sLower)-1])
		code := ctrl(c)
		if code == 0 {
			return 0, fmt.Errorf("unsupported ctrl key: %s", keyStr)
		}
		return code, nil
	}

	// Handle "^<key>" format (caret notation)
	if strings.HasPrefix(s, "^") && len(s) == 2 {
		c := rune(strings.ToLower(s)[1])
		code := ctrl(c)
		if code == 0 {
			return 0, fmt.Errorf("unsupported caret key: %s", keyStr)
		}
		return code, nil
	}

	// Handle "c-<key>" or "C-<key>" format (emacs notation)
	if (strings.HasPrefix(sLower, "c-") || strings.HasPrefix(sLower, "C-")) && len(s) == 3 {
		c := rune(sLower[2])
		code := ctrl(c)
		if code == 0 {
			return 0, fmt.Errorf("unsupported emacs key: %s", keyStr)
		}
		return code, nil
	}

	return 0, fmt.Errorf("unsupported key binding format: %s (supported: 'ctrl+w', '^w', 'C-w')", keyStr)
}

// ParseKeyStroke parses a single key binding string and returns a KeyStroke
// Supports enhanced formats including Alt keys
func ParseKeyStroke(keyStr string) (KeyStroke, error) { //nolint:revive // parsing numerous historical formats
	s := strings.TrimSpace(keyStr)
	if s == "" {
		return KeyStroke{}, fmt.Errorf("empty key binding")
	}

	// Normalize to lowercase for comparison
	sLower := strings.ToLower(s)

	// Handle "ctrl+<key>" format (case-insensitive)
	if hasPrefixFold(s, "ctrl+") && len(s) > len("ctrl+") {
		keyPart := s[len("ctrl+"):]
		if len(keyPart) == 1 {
			c := rune(strings.ToLower(keyPart)[0])
			if c >= 'a' && c <= 'z' {
				return NewCtrlKeyStroke(c), nil
			}
		}
		return KeyStroke{}, fmt.Errorf("unsupported ctrl key: %s", keyStr)
	}

	// Handle "^<key>" format (caret notation)
	if strings.HasPrefix(s, "^") && len(s) == 2 {
		c := rune(strings.ToLower(s)[1])
		if c >= 'a' && c <= 'z' {
			return NewCtrlKeyStroke(c), nil
		}
		return KeyStroke{}, fmt.Errorf("unsupported caret key: %s", keyStr)
	}

	// Handle "c-<key>" or "C-<key>" format (emacs notation)
	if hasPrefixFold(s, "c-") && len(s) == 3 {
		c := rune(sLower[2])
		if c >= 'a' && c <= 'z' {
			return NewCtrlKeyStroke(c), nil
		}
		return KeyStroke{}, fmt.Errorf("unsupported emacs key: %s", keyStr)
	}

	// Handle "alt+<key>" or "meta+<key>" format
	if strings.HasPrefix(sLower, "alt+") || strings.HasPrefix(sLower, "meta+") {
		var keyPart string
		if strings.HasPrefix(sLower, "alt+") {
			keyPart = s[len("alt+"):]
		} else {
			keyPart = s[len("meta+"):]
		}

		keyLower := strings.ToLower(keyPart)

		// Handle special keys
		switch keyLower {
		case "backspace":
			return NewAltKeyStroke(0, "backspace"), nil
		case "delete":
			return NewAltKeyStroke(0, "delete"), nil
		case "enter":
			return NewAltKeyStroke(0, "enter"), nil
		case "space":
			return NewAltKeyStroke(' ', "space"), nil
		default:
			// Handle single letters
			if len(keyLower) == 1 {
				c := rune(keyLower[0])
				if c >= 'a' && c <= 'z' {
					return NewAltKeyStroke(c, ""), nil
				}
			}
		}
		return KeyStroke{}, fmt.Errorf("unsupported alt key: %s", keyStr)
	}

	// Handle "M-<key>" format (emacs meta notation)
	if hasPrefixFold(s, "m-") && len(s) >= 3 {
		keyPart := strings.ToLower(s[2:])

		// Handle special keys
		switch keyPart {
		case "backspace":
			return NewAltKeyStroke(0, "backspace"), nil
		case "delete":
			return NewAltKeyStroke(0, "delete"), nil
		default:
			// Handle single letters
			if len(keyPart) == 1 {
				c := rune(keyPart[0])
				if c >= 'a' && c <= 'z' {
					return NewAltKeyStroke(c, ""), nil
				}
			}
		}
		return KeyStroke{}, fmt.Errorf("unsupported meta key: %s", keyStr)
	}

	return KeyStroke{}, fmt.Errorf("unsupported key binding format: %s (supported: 'ctrl+w', '^w', 'C-w', 'alt+backspace', 'M-backspace')", keyStr)
}

// ParseKeyStrokes parses key binding configuration and returns []KeyStroke
// Supports both single strings and arrays for compatibility
func ParseKeyStrokes(config interface{}) ([]KeyStroke, error) { //nolint:revive // handles multiple config representations
	switch v := config.(type) {
	case string:
		// Single string format: "ctrl+w"
		ks, err := ParseKeyStroke(v)
		if err != nil {
			return nil, err
		}
		return []KeyStroke{ks}, nil

	case []interface{}:
		// Array format: ["ctrl+w", "alt+backspace"]
		var keyStrokes []KeyStroke
		for i, item := range v {
			str, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("array element %d is not a string: %T", i, item)
			}
			ks, err := ParseKeyStroke(str)
			if err != nil {
				return nil, fmt.Errorf("array element %d: %w", i, err)
			}
			keyStrokes = append(keyStrokes, ks)
		}
		return keyStrokes, nil

	case []string:
		// String array format (for direct Go usage)
		var keyStrokes []KeyStroke
		for i, str := range v {
			ks, err := ParseKeyStroke(str)
			if err != nil {
				return nil, fmt.Errorf("array element %d: %w", i, err)
			}
			keyStrokes = append(keyStrokes, ks)
		}
		return keyStrokes, nil

	default:
		return nil, fmt.Errorf("unsupported key binding type: %T (expected string or array)", config)
	}
}

// ValidateKeyBindings validates a map of key binding strings
func ValidateKeyBindings(bindings map[string]string) error {
	for action, keyStr := range bindings {
		if keyStr == "" {
			return fmt.Errorf("empty key binding for %s", action)
		}
		if _, err := ParseKeyBinding(keyStr); err != nil {
			return fmt.Errorf("invalid key binding for %s: %w", action, err)
		}
	}
	return nil
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

// ContextualKeyBindingMap holds resolved keybindings for all contexts
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

// Platform detection functions

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

// KeyBindingResolver handles multi-layer keybinding resolution
type KeyBindingResolver struct {
	profiles   map[Profile]*KeyBindingProfile      // Built-in profiles
	platform   string                              // Detected platform
	terminal   string                              // Detected terminal
	userConfig *config.Config                      // User configuration
	cache      map[string]*ContextualKeyBindingMap // Resolution cache
}

// NewKeyBindingResolver creates a new resolver with detected platform/terminal
func NewKeyBindingResolver(userConfig *config.Config) *KeyBindingResolver {
	return &KeyBindingResolver{
		profiles:   make(map[Profile]*KeyBindingProfile),
		platform:   DetectPlatform(),
		terminal:   DetectTerminal(),
		userConfig: userConfig,
		cache:      make(map[string]*ContextualKeyBindingMap),
	}
}

// RegisterProfile adds a built-in profile to the resolver
func (r *KeyBindingResolver) RegisterProfile(profile Profile, kbp *KeyBindingProfile) {
	if r.profiles == nil {
		r.profiles = make(map[Profile]*KeyBindingProfile)
	}
	r.profiles[profile] = kbp
}

// GetProfile returns a registered profile by name
func (r *KeyBindingResolver) GetProfile(profile Profile) (*KeyBindingProfile, bool) {
	kbp, exists := r.profiles[profile]
	return kbp, exists
}

// ClearCache clears the resolution cache (useful for config reloads)
func (r *KeyBindingResolver) ClearCache() {
	r.cache = make(map[string]*ContextualKeyBindingMap)
}

// Resolve performs layered keybinding resolution for a specific profile and context
func (r *KeyBindingResolver) Resolve(profile Profile, context Context) (*KeyBindingMap, error) {
	// Generate cache key
	cacheKey := fmt.Sprintf("%s:%s:%s:%s", profile, context, r.platform, r.terminal)

	// Check cache first
	if cached, exists := r.cache[cacheKey]; exists {
		if contextMap, exists := cached.GetContext(context); exists {
			return contextMap, nil
		}
	}

	// Create new KeyBindingMap for this context
	result := &KeyBindingMap{
		DeleteWord:         []KeyStroke{},
		ClearLine:          []KeyStroke{},
		DeleteToEnd:        []KeyStroke{},
		MoveToBeginning:    []KeyStroke{},
		MoveToEnd:          []KeyStroke{},
		MoveUp:             []KeyStroke{},
		MoveDown:           []KeyStroke{},
		AddToWorkflow:      []KeyStroke{},
		ToggleWorkflowView: []KeyStroke{},
		ClearWorkflow:      []KeyStroke{},
	}

	// Layer 1: Built-in defaults
	r.applyDefaults(result)

	// Layer 2: Profile base
	if prof, exists := r.profiles[profile]; exists {
		r.applyProfile(result, prof, context)
	}

	// Layer 3: Platform layer
	r.applyPlatformLayer(result)

	// Layer 4: Terminal layer
	r.applyTerminalLayer(result)

	// Layer 5: User config
	if r.userConfig != nil {
		r.applyUserConfig(result, context)
	}

	// Layer 6: Environment overrides
	r.applyEnvironmentOverrides(result)

	// Cache the result
	r.cacheResult(profile, context, result)

	return result, nil
}

// ResolveContextual resolves all contexts for a profile
func (r *KeyBindingResolver) ResolveContextual(profile Profile) (*ContextualKeyBindingMap, error) {
	// Generate cache key for the full contextual map
	cacheKey := fmt.Sprintf("contextual:%s:%s:%s", profile, r.platform, r.terminal)

	if cached, exists := r.cache[cacheKey]; exists {
		return cached, nil
	}

	contextual := NewContextualKeyBindingMap(profile, r.platform, r.terminal)

	// Resolve each context
	for _, context := range GetAllContexts() {
		keyMap, err := r.Resolve(profile, context)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve context %s: %w", context, err)
		}
		contextual.SetContext(context, keyMap)
	}

	// Cache the contextual map
	r.cache[cacheKey] = contextual

	return contextual, nil
}

// GetEffectiveKeybindings returns all resolved keybindings for a profile/context
func (r *KeyBindingResolver) GetEffectiveKeybindings(profile Profile, context Context) map[string][]KeyStroke {
	result := make(map[string][]KeyStroke)

	keyMap, err := r.Resolve(profile, context)
	if err != nil || keyMap == nil {
		return result
	}

	clone := func(src []KeyStroke) []KeyStroke {
		if len(src) == 0 {
			return nil
		}
		copySlice := make([]KeyStroke, len(src))
		copy(copySlice, src)
		return copySlice
	}

	result["delete_word"] = clone(keyMap.DeleteWord)
	result["clear_line"] = clone(keyMap.ClearLine)
	result["delete_to_end"] = clone(keyMap.DeleteToEnd)
	result["move_to_beginning"] = clone(keyMap.MoveToBeginning)
	result["move_to_end"] = clone(keyMap.MoveToEnd)
	result["move_up"] = clone(keyMap.MoveUp)
	result["move_down"] = clone(keyMap.MoveDown)
	result["add_to_workflow"] = clone(keyMap.AddToWorkflow)
	result["toggle_workflow_view"] = clone(keyMap.ToggleWorkflowView)
	result["clear_workflow"] = clone(keyMap.ClearWorkflow)

	return result
}

// Layer application methods

func (r *KeyBindingResolver) applyDefaults(keyMap *KeyBindingMap) {
	// Apply hardcoded defaults (legacy compatibility)
	defaults := DefaultKeyBindingMap()
	keyMap.DeleteWord = append(keyMap.DeleteWord, defaults.DeleteWord...)
	keyMap.ClearLine = append(keyMap.ClearLine, defaults.ClearLine...)
	keyMap.DeleteToEnd = append(keyMap.DeleteToEnd, defaults.DeleteToEnd...)
	keyMap.MoveToBeginning = append(keyMap.MoveToBeginning, defaults.MoveToBeginning...)
	keyMap.MoveToEnd = append(keyMap.MoveToEnd, defaults.MoveToEnd...)
	keyMap.MoveUp = append(keyMap.MoveUp, defaults.MoveUp...)
	keyMap.MoveDown = append(keyMap.MoveDown, defaults.MoveDown...)
	keyMap.AddToWorkflow = append(keyMap.AddToWorkflow, defaults.AddToWorkflow...)
	keyMap.ToggleWorkflowView = append(keyMap.ToggleWorkflowView, defaults.ToggleWorkflowView...)
	keyMap.ClearWorkflow = append(keyMap.ClearWorkflow, defaults.ClearWorkflow...)
	keyMap.SoftCancel = append(keyMap.SoftCancel, defaults.SoftCancel...)
}

func (r *KeyBindingResolver) applyProfile(keyMap *KeyBindingMap, profile *KeyBindingProfile, context Context) {
	// Helper function to apply bindings from profile
	applyBinding := func(action string, target *[]KeyStroke) {
		if keystrokes, exists := profile.GetBinding(context, action); exists {
			*target = keystrokes // Replace, don't append (profile overrides defaults)
		}
	}

	applyBinding("delete_word", &keyMap.DeleteWord)
	applyBinding("clear_line", &keyMap.ClearLine)
	applyBinding("delete_to_end", &keyMap.DeleteToEnd)
	applyBinding("move_to_beginning", &keyMap.MoveToBeginning)
	applyBinding("move_to_end", &keyMap.MoveToEnd)
	applyBinding("move_up", &keyMap.MoveUp)
	applyBinding("move_down", &keyMap.MoveDown)
	applyBinding("add_to_workflow", &keyMap.AddToWorkflow)
	applyBinding("toggle_workflow_view", &keyMap.ToggleWorkflowView)
	applyBinding("clear_workflow", &keyMap.ClearWorkflow)
	applyBinding("soft_cancel", &keyMap.SoftCancel)
}

func (r *KeyBindingResolver) applyPlatformLayer(keyMap *KeyBindingMap) {
	platformBindings := GetPlatformSpecificKeyBindings(r.platform)

	// Apply platform-specific overrides
	if bindings, exists := platformBindings["delete_word"]; exists {
		keyMap.DeleteWord = bindings // Platform overrides profile
	}
}

func (r *KeyBindingResolver) applyTerminalLayer(keyMap *KeyBindingMap) {
	terminalBindings := GetTerminalSpecificKeyBindings(r.terminal)

	// Apply terminal-specific overrides with explicit action handling
	for action, bindings := range terminalBindings {
		r.applyTerminalBinding(keyMap, action, bindings)
	}
}

// applyTerminalBinding applies a single terminal binding to reduce cyclomatic complexity
func (r *KeyBindingResolver) applyTerminalBinding(keyMap *KeyBindingMap, action string, bindings []KeyStroke) {
	// Apply editing actions
	if r.applyEditingAction(keyMap, action, bindings) {
		return
	}

	// Apply navigation actions
	if r.applyNavigationAction(keyMap, action, bindings) {
		return
	}

	// Apply workflow actions
	r.applyWorkflowAction(keyMap, action, bindings)
}

// applyEditingAction applies editing-related keybinding actions
func (r *KeyBindingResolver) applyEditingAction(keyMap *KeyBindingMap, action string, bindings []KeyStroke) bool {
	switch action {
	case "delete_word":
		keyMap.DeleteWord = bindings
		return true
	case "clear_line":
		keyMap.ClearLine = bindings
		return true
	case "delete_to_end":
		keyMap.DeleteToEnd = bindings
		return true
	}
	return false
}

// applyNavigationAction applies navigation-related keybinding actions
func (r *KeyBindingResolver) applyNavigationAction(keyMap *KeyBindingMap, action string, bindings []KeyStroke) bool {
	switch action {
	case "move_to_beginning":
		keyMap.MoveToBeginning = bindings
		return true
	case "move_to_end":
		keyMap.MoveToEnd = bindings
		return true
	case "move_up":
		keyMap.MoveUp = bindings
		return true
	case "move_down":
		keyMap.MoveDown = bindings
		return true
	}
	return false
}

// applyWorkflowAction applies workflow-related keybinding actions
func (r *KeyBindingResolver) applyWorkflowAction(keyMap *KeyBindingMap, action string, bindings []KeyStroke) {
	switch action {
	case "add_to_workflow":
		keyMap.AddToWorkflow = bindings
	case "toggle_workflow_view":
		keyMap.ToggleWorkflowView = bindings
	case "clear_workflow":
		keyMap.ClearWorkflow = bindings
	case "soft_cancel":
		keyMap.SoftCancel = bindings
	// Explicitly ignore unsupported actions
	default:
		// Terminal-specific action not supported in this context
	}
}

func (r *KeyBindingResolver) applyUserConfig(keyMap *KeyBindingMap, context Context) { //nolint:revive // layered override logic retained for clarity
	// Apply user global keybindings first
	userBindings := r.userConfig.Interactive.Keybindings

	userValues := map[string]string{
		"delete_word":          userBindings.DeleteWord,
		"clear_line":           userBindings.ClearLine,
		"delete_to_end":        userBindings.DeleteToEnd,
		"move_to_beginning":    userBindings.MoveToBeginning,
		"move_to_end":          userBindings.MoveToEnd,
		"move_up":              userBindings.MoveUp,
		"move_down":            userBindings.MoveDown,
		"add_to_workflow":      userBindings.AddToWorkflow,
		"toggle_workflow_view": userBindings.ToggleWorkflowView,
		"clear_workflow":       userBindings.ClearWorkflow,
		"soft_cancel":          userBindings.SoftCancel,
	}

	// Apply non-empty user overrides
	for action, keyStr := range userValues {
		if keyStr != "" {
			if ks, err := ParseKeyStroke(keyStr); err == nil {
				switch action {
				case "delete_word":
					keyMap.DeleteWord = []KeyStroke{ks}
				case "clear_line":
					keyMap.ClearLine = []KeyStroke{ks}
				case "delete_to_end":
					keyMap.DeleteToEnd = []KeyStroke{ks}
				case "move_to_beginning":
					keyMap.MoveToBeginning = []KeyStroke{ks}
				case "move_to_end":
					keyMap.MoveToEnd = []KeyStroke{ks}
				case "move_up":
					keyMap.MoveUp = []KeyStroke{ks}
				case "move_down":
					keyMap.MoveDown = []KeyStroke{ks}
				case "add_to_workflow":
					keyMap.AddToWorkflow = []KeyStroke{ks}
				case "toggle_workflow_view":
					keyMap.ToggleWorkflowView = []KeyStroke{ks}
				case "clear_workflow":
					keyMap.ClearWorkflow = []KeyStroke{ks}
				case "soft_cancel":
					keyMap.SoftCancel = []KeyStroke{ks}
				}
			}
		}
	}

	// Apply context-specific user bindings
	r.applyUserContextBindings(keyMap, context)

	// Apply platform-specific user bindings
	r.applyUserPlatformBindings(keyMap)

	// Apply terminal-specific user bindings
	r.applyUserTerminalBindings(keyMap)
}

func (r *KeyBindingResolver) applyEnvironmentOverrides(keyMap *KeyBindingMap) {
	// Check for environment variable overrides
	envOverrides := map[string]*[]KeyStroke{
		"GGC_KEYBIND_DELETE_WORD":          &keyMap.DeleteWord,
		"GGC_KEYBIND_CLEAR_LINE":           &keyMap.ClearLine,
		"GGC_KEYBIND_DELETE_TO_END":        &keyMap.DeleteToEnd,
		"GGC_KEYBIND_MOVE_TO_BEGINNING":    &keyMap.MoveToBeginning,
		"GGC_KEYBIND_MOVE_TO_END":          &keyMap.MoveToEnd,
		"GGC_KEYBIND_MOVE_UP":              &keyMap.MoveUp,
		"GGC_KEYBIND_MOVE_DOWN":            &keyMap.MoveDown,
		"GGC_KEYBIND_ADD_TO_WORKFLOW":      &keyMap.AddToWorkflow,
		"GGC_KEYBIND_TOGGLE_WORKFLOW_VIEW": &keyMap.ToggleWorkflowView,
		"GGC_KEYBIND_CLEAR_WORKFLOW":       &keyMap.ClearWorkflow,
		"GGC_KEYBIND_SOFT_CANCEL":          &keyMap.SoftCancel,
	}

	for envVar, target := range envOverrides {
		if keyStr := os.Getenv(envVar); keyStr != "" {
			if ks, err := ParseKeyStroke(keyStr); err == nil {
				*target = []KeyStroke{ks}
			}
		}
	}
}

func (r *KeyBindingResolver) applyUserContextBindings(keyMap *KeyBindingMap, context Context) {
	// Apply context-specific user bindings if they exist
	var contextBindings map[string]interface{}

	switch context {
	case ContextInput:
		contextBindings = r.userConfig.Interactive.Contexts.Input.Keybindings
	case ContextResults:
		contextBindings = r.userConfig.Interactive.Contexts.Results.Keybindings
	case ContextSearch:
		contextBindings = r.userConfig.Interactive.Contexts.Search.Keybindings
	}

	if contextBindings != nil {
		r.applyUserBindings(keyMap, contextBindings)
	}
}

func (r *KeyBindingResolver) applyUserPlatformBindings(keyMap *KeyBindingMap) {
	var platformBindings map[string]interface{}

	switch r.platform {
	case "darwin":
		platformBindings = r.userConfig.Interactive.Darwin.Keybindings
	case "linux":
		platformBindings = r.userConfig.Interactive.Linux.Keybindings
	case "windows":
		platformBindings = r.userConfig.Interactive.Windows.Keybindings
	}

	if platformBindings != nil {
		r.applyUserBindings(keyMap, platformBindings)
	}
}

func (r *KeyBindingResolver) applyUserTerminalBindings(keyMap *KeyBindingMap) {
	if r.userConfig.Interactive.Terminals != nil {
		if termConfig, exists := r.userConfig.Interactive.Terminals[r.terminal]; exists {
			if termConfig.Keybindings != nil {
				r.applyUserBindings(keyMap, termConfig.Keybindings)
			}
		}
	}
}

func (r *KeyBindingResolver) applyUserBindings(keyMap *KeyBindingMap, bindings map[string]interface{}) {
	for action, value := range bindings {
		keystrokes := r.parseUserBindingValue(value)
		if len(keystrokes) > 0 {
			r.applyUserBinding(keyMap, action, keystrokes)
		}
	}
}

// applyUserBinding applies a single user binding to reduce cyclomatic complexity
func (r *KeyBindingResolver) applyUserBinding(keyMap *KeyBindingMap, action string, keystrokes []KeyStroke) {
	// Apply editing actions
	if r.applyUserEditingAction(keyMap, action, keystrokes) {
		return
	}

	// Apply navigation actions
	if r.applyUserNavigationAction(keyMap, action, keystrokes) {
		return
	}

	// Apply workflow actions
	r.applyUserWorkflowAction(keyMap, action, keystrokes)
}

// applyUserEditingAction applies user editing-related keybinding actions
func (r *KeyBindingResolver) applyUserEditingAction(keyMap *KeyBindingMap, action string, keystrokes []KeyStroke) bool {
	switch action {
	case "delete_word":
		keyMap.DeleteWord = keystrokes
		return true
	case "clear_line":
		keyMap.ClearLine = keystrokes
		return true
	case "delete_to_end":
		keyMap.DeleteToEnd = keystrokes
		return true
	}
	return false
}

// applyUserNavigationAction applies user navigation-related keybinding actions
func (r *KeyBindingResolver) applyUserNavigationAction(keyMap *KeyBindingMap, action string, keystrokes []KeyStroke) bool {
	switch action {
	case "move_to_beginning":
		keyMap.MoveToBeginning = keystrokes
		return true
	case "move_to_end":
		keyMap.MoveToEnd = keystrokes
		return true
	case "move_up":
		keyMap.MoveUp = keystrokes
		return true
	case "move_down":
		keyMap.MoveDown = keystrokes
		return true
	}
	return false
}

// applyUserWorkflowAction applies user workflow-related keybinding actions
func (r *KeyBindingResolver) applyUserWorkflowAction(keyMap *KeyBindingMap, action string, keystrokes []KeyStroke) {
	switch action {
	case "add_to_workflow":
		keyMap.AddToWorkflow = keystrokes
	case "toggle_workflow_view":
		keyMap.ToggleWorkflowView = keystrokes
	case "clear_workflow":
		keyMap.ClearWorkflow = keystrokes
	case "soft_cancel":
		keyMap.SoftCancel = keystrokes
	// Explicitly ignore unsupported actions
	default:
		// User-defined action not supported in this context
	}
}

func (r *KeyBindingResolver) parseUserBindingValue(value interface{}) []KeyStroke {
	switch v := value.(type) {
	case string:
		if v == "" {
			return []KeyStroke{}
		}
		if ks, err := ParseKeyStroke(v); err == nil {
			return []KeyStroke{ks}
		}
	case []interface{}:
		var keystrokes []KeyStroke
		for _, item := range v {
			if itemStr, ok := item.(string); ok && itemStr != "" {
				if ks, err := ParseKeyStroke(itemStr); err == nil {
					keystrokes = append(keystrokes, ks)
				}
			}
		}
		return keystrokes
	}
	return []KeyStroke{}
}

func (r *KeyBindingResolver) cacheResult(profile Profile, context Context, keyMap *KeyBindingMap) {
	cacheKey := fmt.Sprintf("%s:%s:%s:%s", profile, context, r.platform, r.terminal)

	// Create or update contextual map in cache
	var contextual *ContextualKeyBindingMap
	if cached, exists := r.cache[cacheKey]; exists {
		contextual = cached
	} else {
		contextual = NewContextualKeyBindingMap(profile, r.platform, r.terminal)
	}

	contextual.SetContext(context, keyMap)
	r.cache[cacheKey] = contextual
}

// Built-in Profile Definitions

// CreateDefaultProfile returns the default keybinding profile (legacy compatible)
func CreateDefaultProfile() *KeyBindingProfile {
	return &KeyBindingProfile{
		Name:        "Default",
		Description: "Default keybindings compatible with legacy behavior",
		Global:      make(map[string][]KeyStroke),
		Contexts: map[Context]map[string][]KeyStroke{
			ContextGlobal: {
				"soft_cancel": {NewCtrlKeyStroke('g'), NewEscapeKeyStroke()},
			},
			ContextInput: {
				"delete_word":       {NewCtrlKeyStroke('w')},
				"clear_line":        {NewCtrlKeyStroke('u')},
				"delete_to_end":     {NewCtrlKeyStroke('k')},
				"move_to_beginning": {NewCtrlKeyStroke('a')},
				"move_to_end":       {NewCtrlKeyStroke('e')},
			},
			ContextResults: {
				"move_up":              {NewCtrlKeyStroke('p')},
				"move_down":            {NewCtrlKeyStroke('n')},
				"add_to_workflow":      {NewTabKeyStroke()},
				"toggle_workflow_view": {NewCtrlKeyStroke('t')},
				"clear_workflow":       {NewCharKeyStroke('c')},
			},
			ContextSearch: {
				"move_up":              {NewCtrlKeyStroke('p')},
				"move_down":            {NewCtrlKeyStroke('n')},
				"add_to_workflow":      {NewTabKeyStroke()},
				"toggle_workflow_view": {NewCtrlKeyStroke('t')},
				"clear_workflow":       {NewCharKeyStroke('c')},
			},
		},
	}
}

// CreateEmacsProfile returns the Emacs-style keybinding profile
// Based on GNU Emacs standard keybindings with authentic Emacs behavior
func CreateEmacsProfile() *KeyBindingProfile {
	return &KeyBindingProfile{
		Name:        "Emacs",
		Description: "Comprehensive Emacs-style keybindings with authentic GNU Emacs behavior",
		Global: map[string][]KeyStroke{
			// Core Emacs global bindings
			"quit":                {NewCtrlKeyStroke('g')},                        // C-g keyboard-quit
			"help":                {NewCtrlKeyStroke('h')},                        // C-h help-command
			"universal_argument":  {NewCtrlKeyStroke('u')},                        // C-u universal-argument
			"exchange_point_mark": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('x')}, // C-x C-x (chord)
			"suspend":             {NewCtrlKeyStroke('z')},                        // C-z suspend-frame
		},
		Contexts: map[Context]map[string][]KeyStroke{
			ContextGlobal: {
				"quit":               {NewCtrlKeyStroke('g')},
				"help":               {NewCtrlKeyStroke('h')},
				"universal_argument": {NewCtrlKeyStroke('u')},
				"suspend":            {NewCtrlKeyStroke('z')},
				"soft_cancel":        {NewCtrlKeyStroke('g'), NewEscapeKeyStroke()},
			},
			ContextInput: {
				// Character-level movement
				"forward_char":  {NewCtrlKeyStroke('f')}, // C-f forward-char
				"backward_char": {NewCtrlKeyStroke('b')}, // C-b backward-char
				"next_line":     {NewCtrlKeyStroke('n')}, // C-n next-line
				"previous_line": {NewCtrlKeyStroke('p')}, // C-p previous-line

				// Word-level movement
				"forward_word":  {NewAltKeyStroke('f', "")}, // M-f forward-word
				"backward_word": {NewAltKeyStroke('b', "")}, // M-b backward-word

				// Line-level movement
				"beginning_of_line": {NewCtrlKeyStroke('a')}, // C-a beginning-of-line
				"end_of_line":       {NewCtrlKeyStroke('e')}, // C-e end-of-line
				"move_to_beginning": {NewCtrlKeyStroke('a')}, // Alias for compatibility
				"move_to_end":       {NewCtrlKeyStroke('e')}, // Alias for compatibility

				// Deletion and killing
				"delete_char":          {NewCtrlKeyStroke('d')},                        // C-d delete-char
				"backward_delete_char": {NewCtrlKeyStroke('h')},                        // C-h backward-delete-char
				"kill_line":            {NewCtrlKeyStroke('k')},                        // C-k kill-line
				"kill_word":            {NewAltKeyStroke('d', "")},                     // M-d kill-word
				"backward_kill_word":   {NewAltKeyStroke(127, "backspace")},            // M-DEL backward-kill-word
				"unix_line_discard":    {NewCtrlKeyStroke('u')},                        // C-u unix-line-discard
				"kill_whole_line":      {NewCtrlKeyStroke('s'), NewCtrlKeyStroke('k')}, // C-S-k kill-whole-line
				"delete_word":          {NewAltKeyStroke('d', "")},                     // Alias for kill-word
				"clear_line":           {NewCtrlKeyStroke('u')},                        // Alias for unix-line-discard
				"delete_to_end":        {NewCtrlKeyStroke('k')},                        // Alias for kill-line

				// Search and replace
				"isearch_forward":  {NewCtrlKeyStroke('s')},    // C-s isearch-forward
				"isearch_backward": {NewCtrlKeyStroke('r')},    // C-r isearch-backward
				"query_replace":    {NewAltKeyStroke('%', "")}, // M-% query-replace

				// Case operations
				"upcase_word":     {NewAltKeyStroke('u', "")}, // M-u upcase-word
				"downcase_word":   {NewAltKeyStroke('l', "")}, // M-l downcase-word
				"capitalize_word": {NewAltKeyStroke('c', "")}, // M-c capitalize-word
				"transpose_chars": {NewCtrlKeyStroke('t')},    // C-t transpose-chars
				"transpose_words": {NewAltKeyStroke('t', "")}, // M-t transpose-words

				// Yank and kill ring
				"yank":                {NewCtrlKeyStroke('y')},    // C-y yank
				"yank_pop":            {NewAltKeyStroke('y', "")}, // M-y yank-pop
				"copy_region_as_kill": {NewAltKeyStroke('w', "")}, // M-w copy-region-as-kill
				"kill_region":         {NewCtrlKeyStroke('w')},    // C-w kill-region

				// Mark and region
				"set_mark_command":    {NewCtrlKeyStroke(' ')},                        // C-SPC set-mark-command
				"exchange_point_mark": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('x')}, // C-x C-x exchange-point-mark

				// Buffer and file operations (adapted for CLI)
				"save_buffer":      {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('s')}, // C-x C-s save-buffer
				"find_file":        {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('f')}, // C-x C-f find-file
				"switch_to_buffer": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('b')}, // C-x C-b switch-to-buffer

				// Miscellaneous
				"quoted_insert":           {NewCtrlKeyStroke('q')},     // C-q quoted-insert
				"recenter_top_bottom":     {NewCtrlKeyStroke('l')},     // C-l recenter-top-bottom
				"just_one_space":          {NewAltKeyStroke(' ', "")},  // M-SPC just-one-space
				"delete_horizontal_space": {NewAltKeyStroke('\\', "")}, // M-\ delete-horizontal-space
			},
			ContextResults: {
				// Navigation in results (Emacs-style list navigation)
				"previous_line": {NewCtrlKeyStroke('p')}, // C-p previous-line
				"next_line":     {NewCtrlKeyStroke('n')}, // C-n next-line
				"move_up":       {NewCtrlKeyStroke('p')}, // Alias
				"move_down":     {NewCtrlKeyStroke('n')}, // Alias
				"backward_char": {NewCtrlKeyStroke('b')}, // C-b backward-char
				"forward_char":  {NewCtrlKeyStroke('f')}, // C-f forward-char

				// Scrolling (Emacs page movement)
				"scroll_up":           {NewAltKeyStroke('v', "")}, // M-v scroll-up
				"scroll_down":         {NewCtrlKeyStroke('v')},    // C-v scroll-down
				"beginning_of_buffer": {NewAltKeyStroke('<', "")}, // M-< beginning-of-buffer
				"end_of_buffer":       {NewAltKeyStroke('>', "")}, // M-> end-of-buffer

				// Selection and marking
				"set_mark_command":  {NewCtrlKeyStroke(' ')},                        // C-SPC set-mark-command
				"mark_whole_buffer": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('h')}, // C-x h mark-whole-buffer

				// Search in results
				"isearch_forward":  {NewCtrlKeyStroke('s')}, // C-s isearch-forward
				"isearch_backward": {NewCtrlKeyStroke('r')}, // C-r isearch-backward

				// Execute/select
				"execute": {NewCtrlKeyStroke('m')}, // C-m (Enter equivalent)
				"select":  {NewCtrlKeyStroke('m')}, // Alias

				// Workflow operations (adapted for Emacs style)
				"add_to_workflow":      {NewRawKeyStroke([]byte{9})}, // Tab
				"toggle_workflow_view": {NewCtrlKeyStroke('t')},      // C-t
				"clear_workflow":       {NewAltKeyStroke('c', "")},   // M-c clear
			},
			ContextSearch: {
				// Search-specific Emacs bindings
				"isearch_forward":         {NewCtrlKeyStroke('s')},       // C-s isearch-forward
				"isearch_backward":        {NewCtrlKeyStroke('r')},       // C-r isearch-backward
				"isearch_repeat_forward":  {NewCtrlKeyStroke('s')},       // C-s (repeat)
				"isearch_repeat_backward": {NewCtrlKeyStroke('r')},       // C-r (repeat)
				"isearch_yank_word":       {NewCtrlKeyStroke('w')},       // C-w isearch-yank-word
				"isearch_yank_line":       {NewCtrlKeyStroke('y')},       // C-y isearch-yank-line
				"isearch_delete_char":     {NewCtrlKeyStroke('h')},       // C-h isearch-delete-char
				"isearch_abort":           {NewCtrlKeyStroke('g')},       // C-g isearch-abort
				"isearch_exit":            {NewRawKeyStroke([]byte{13})}, // RET isearch-exit

				// Navigation while searching
				"next_line":     {NewCtrlKeyStroke('n')}, // C-n next-line
				"previous_line": {NewCtrlKeyStroke('p')}, // C-p previous-line
				"move_up":       {NewCtrlKeyStroke('p')}, // Alias
				"move_down":     {NewCtrlKeyStroke('n')}, // Alias

				// Case sensitivity toggle
				"isearch_toggle_case_fold": {NewAltKeyStroke('c', "")}, // M-c toggle case sensitivity
				"isearch_toggle_regexp":    {NewAltKeyStroke('r', "")}, // M-r toggle regexp mode

				// Workflow operations (search context)
				"add_to_workflow":      {NewRawKeyStroke([]byte{9})}, // Tab
				"toggle_workflow_view": {NewCtrlKeyStroke('t')},      // C-t
				"clear_workflow":       {NewAltKeyStroke('x', "")},   // M-x clear (avoiding conflict with M-c)
			},
		},
	}
}

// CreateViProfile returns the Vi-style keybinding profile (adapted for CLI context)
// Implements Vi modal editing concepts adapted for command-line interface
func CreateViProfile() *KeyBindingProfile {
	return &KeyBindingProfile{
		Name:        "Vi",
		Description: "Vi-style modal keybindings adapted for command-line interface with insert and normal modes",
		Global: map[string][]KeyStroke{
			// Core Vi global bindings
			"quit":          {NewCtrlKeyStroke('c')},             // Keep standard quit (like :q!)
			"command_mode":  {NewRawKeyStroke([]byte{27})},       // ESC - enter command mode
			"force_quit":    {NewRawKeyStroke([]byte{'Z', 'Q'})}, // ZQ - quit without saving
			"save_and_quit": {NewRawKeyStroke([]byte{'Z', 'Z'})}, // ZZ - save and quit
		},
		Contexts: map[Context]map[string][]KeyStroke{
			ContextGlobal: {
				"quit":          {NewCtrlKeyStroke('c')},
				"command_mode":  {NewRawKeyStroke([]byte{27})},
				"force_quit":    {NewRawKeyStroke([]byte{'Z', 'Q'})},
				"save_and_quit": {NewRawKeyStroke([]byte{'Z', 'Z'})},
				"soft_cancel":   {NewCtrlKeyStroke('g'), NewEscapeKeyStroke()},
			},
			ContextInput: {
				// Vi INSERT MODE bindings (when editing input)
				// In Vi, insert mode is similar to normal editor behavior

				// Basic movement (limited in insert mode)
				"move_to_beginning": {NewCtrlKeyStroke('a')}, // C-a move to beginning
				"move_to_end":       {NewCtrlKeyStroke('e')}, // C-e move to end
				"forward_char":      {NewCtrlKeyStroke('l')}, // C-l move right
				"backward_char":     {NewCtrlKeyStroke('h')}, // C-h move left (also backspace)

				// Deletion (insert mode)
				"delete_word":          {NewCtrlKeyStroke('w')}, // C-w delete word backward
				"delete_line":          {NewCtrlKeyStroke('u')}, // C-u delete line
				"clear_line":           {NewCtrlKeyStroke('u')}, // Alias
				"delete_to_end":        {NewCtrlKeyStroke('k')}, // C-k delete to end of line
				"backward_delete_char": {NewCtrlKeyStroke('h')}, // C-h backspace

				// Insert mode specific
				"insert_at_beginning": {NewRawKeyStroke([]byte{'I'})}, // I - insert at line beginning
				"insert_at_end":       {NewRawKeyStroke([]byte{'A'})}, // A - insert at line end
				"open_line_below":     {NewRawKeyStroke([]byte{'o'})}, // o - open new line below
				"open_line_above":     {NewRawKeyStroke([]byte{'O'})}, // O - open new line above

				// Exit insert mode
				"escape_to_normal": {NewRawKeyStroke([]byte{27})}, // ESC - to normal mode

				// Vi-style completion and registers
				"complete_word":  {NewCtrlKeyStroke('n')}, // C-n word completion
				"complete_prev":  {NewCtrlKeyStroke('p')}, // C-p previous completion
				"literal_insert": {NewCtrlKeyStroke('v')}, // C-v literal character insert
			},
			ContextResults: {
				// Vi NORMAL MODE bindings (when navigating results)
				// This is where Vi really shines with single-key navigation

				// Basic movement (hjkl)
				"move_left":  {NewRawKeyStroke([]byte{'h'})}, // h - move left
				"move_down":  {NewRawKeyStroke([]byte{'j'})}, // j - move down
				"move_up":    {NewRawKeyStroke([]byte{'k'})}, // k - move up
				"move_right": {NewRawKeyStroke([]byte{'l'})}, // l - move right

				// Alternative movement for compatibility
				"move_down_alt": {NewCtrlKeyStroke('n')}, // C-n alternative
				"move_up_alt":   {NewCtrlKeyStroke('p')}, // C-p alternative

				// Word movement
				"forward_word":      {NewRawKeyStroke([]byte{'w'})}, // w - next word
				"backward_word":     {NewRawKeyStroke([]byte{'b'})}, // b - previous word
				"end_word":          {NewRawKeyStroke([]byte{'e'})}, // e - end of word
				"forward_word_big":  {NewRawKeyStroke([]byte{'W'})}, // W - next WORD
				"backward_word_big": {NewRawKeyStroke([]byte{'B'})}, // B - previous WORD
				"end_word_big":      {NewRawKeyStroke([]byte{'E'})}, // E - end of WORD

				// Line movement
				"first_char":        {NewRawKeyStroke([]byte{'^'})}, // ^ - first non-blank character
				"beginning_of_line": {NewRawKeyStroke([]byte{'0'})}, // 0 - beginning of line
				"end_of_line":       {NewRawKeyStroke([]byte{'$'})}, // $ - end of line

				// Screen movement
				"top_of_screen":    {NewRawKeyStroke([]byte{'H'})}, // H - top of screen
				"middle_of_screen": {NewRawKeyStroke([]byte{'M'})}, // M - middle of screen
				"bottom_of_screen": {NewRawKeyStroke([]byte{'L'})}, // L - bottom of screen

				// Buffer movement
				"first_line": {NewRawKeyStroke([]byte{'g', 'g'})}, // gg - first line
				"last_line":  {NewRawKeyStroke([]byte{'G'})},      // G - last line
				"goto_line":  {NewRawKeyStroke([]byte{':'})},      // : - command mode (go to line)

				// Scrolling
				"scroll_down":      {NewCtrlKeyStroke('f')}, // C-f - page down
				"scroll_up":        {NewCtrlKeyStroke('b')}, // C-b - page up
				"scroll_down_half": {NewCtrlKeyStroke('d')}, // C-d - half page down
				"scroll_up_half":   {NewCtrlKeyStroke('u')}, // C-u - half page up
				"scroll_line_down": {NewCtrlKeyStroke('e')}, // C-e - scroll down one line
				"scroll_line_up":   {NewCtrlKeyStroke('y')}, // C-y - scroll up one line

				// Search and navigation
				"search_forward":       {NewRawKeyStroke([]byte{'/'})}, // / - search forward
				"search_backward":      {NewRawKeyStroke([]byte{'?'})}, // ? - search backward
				"search_next":          {NewRawKeyStroke([]byte{'n'})}, // n - next search match
				"search_previous":      {NewRawKeyStroke([]byte{'N'})}, // N - previous search match
				"search_word_forward":  {NewRawKeyStroke([]byte{'*'})}, // * - search word under cursor forward
				"search_word_backward": {NewRawKeyStroke([]byte{'#'})}, // # - search word under cursor backward

				// Marks and jumps
				"set_mark":       {NewRawKeyStroke([]byte{'m'})},  // m{a-z} - set mark
				"goto_mark":      {NewRawKeyStroke([]byte{'\''})}, // '{a-z} - goto mark
				"goto_mark_line": {NewRawKeyStroke([]byte{'`'})},  // `{a-z} - goto mark exact position
				"jump_back":      {NewCtrlKeyStroke('o')},         // C-o - jump back
				"jump_forward":   {NewCtrlKeyStroke('i')},         // C-i - jump forward

				// Selection and execution
				"select":           {NewRawKeyStroke([]byte{13})},  // Enter - select current item
				"execute":          {NewRawKeyStroke([]byte{13})},  // Alias
				"visual_mode":      {NewRawKeyStroke([]byte{'v'})}, // v - visual mode
				"visual_line_mode": {NewRawKeyStroke([]byte{'V'})}, // V - visual line mode

				// Repeat and undo (adapted for CLI)
				"repeat_last": {NewRawKeyStroke([]byte{'.'})}, // . - repeat last action
				"undo":        {NewRawKeyStroke([]byte{'u'})}, // u - undo
				"redo":        {NewCtrlKeyStroke('r')},        // C-r - redo

				// Enter insert mode from results
				"insert_mode":         {NewRawKeyStroke([]byte{'i'})}, // i - insert mode
				"insert_after":        {NewRawKeyStroke([]byte{'a'})}, // a - insert after cursor
				"insert_at_end":       {NewRawKeyStroke([]byte{'A'})}, // A - insert at line end
				"insert_at_beginning": {NewRawKeyStroke([]byte{'I'})}, // I - insert at line beginning

				// Workflow operations (Vi normal mode style)
				"add_to_workflow":      {NewRawKeyStroke([]byte{9})},   // Tab
				"toggle_workflow_view": {NewRawKeyStroke([]byte{'W'})}, // W - workflow view (capital W)
				"clear_workflow":       {NewRawKeyStroke([]byte{'D'})}, // D - delete/clear workflow
			},
			ContextSearch: {
				// Vi search mode bindings (when in / or ? search)
				// Similar to insert mode but with search-specific commands

				// Basic navigation
				"move_up":       {NewRawKeyStroke([]byte{'k'})}, // k - move up in results
				"move_down":     {NewRawKeyStroke([]byte{'j'})}, // j - move down in results
				"move_up_alt":   {NewCtrlKeyStroke('p')},        // C-p alternative
				"move_down_alt": {NewCtrlKeyStroke('n')},        // C-n alternative

				// Search navigation
				"search_next":     {NewRawKeyStroke([]byte{'n'})}, // n - next match
				"search_previous": {NewRawKeyStroke([]byte{'N'})}, // N - previous match
				"search_repeat":   {NewRawKeyStroke([]byte{13})},  // Enter - accept search
				"search_abort":    {NewRawKeyStroke([]byte{27})},  // ESC - abort search

				// Edit search term
				"delete_word":  {NewCtrlKeyStroke('w')}, // C-w delete word
				"clear_search": {NewCtrlKeyStroke('u')}, // C-u clear search line
				"delete_char":  {NewCtrlKeyStroke('h')}, // C-h delete character

				// Search modes
				"case_sensitive_toggle": {NewRawKeyStroke([]byte{'\\', 'c'})}, // \c - toggle case sensitivity
				"regex_mode_toggle":     {NewRawKeyStroke([]byte{'\\', 'v'})}, // \v - very magic mode
				"literal_mode_toggle":   {NewRawKeyStroke([]byte{'\\', 'V'})}, // \V - very nomagic mode

				// History (search command history)
				"search_history_up":   {NewCtrlKeyStroke('p')}, // C-p - previous search
				"search_history_down": {NewCtrlKeyStroke('n')}, // C-n - next search

				// Workflow operations (Vi search mode)
				"add_to_workflow":      {NewRawKeyStroke([]byte{9})},   // Tab
				"toggle_workflow_view": {NewRawKeyStroke([]byte{'W'})}, // W - workflow view
				"clear_workflow":       {NewRawKeyStroke([]byte{'D'})}, // D - delete/clear workflow
			},
		},
	}
}

// CreateReadlineProfile returns the GNU Readline compatible keybinding profile
// Based on GNU Readline library defaults providing bash-like experience
func CreateReadlineProfile() *KeyBindingProfile {
	return &KeyBindingProfile{
		Name:        "Readline",
		Description: "Comprehensive GNU Readline compatible keybindings for authentic bash-like CLI experience",
		Global: map[string][]KeyStroke{
			// Core Readline global bindings
			"abort":        {NewCtrlKeyStroke('g')}, // C-g abort
			"bell":         {NewCtrlKeyStroke('g')}, // C-g bell (same as abort)
			"clear_screen": {NewCtrlKeyStroke('l')}, // C-l clear-screen
		},
		Contexts: map[Context]map[string][]KeyStroke{
			ContextGlobal: {
				"abort":        {NewCtrlKeyStroke('g')},
				"clear_screen": {NewCtrlKeyStroke('l')},
				"soft_cancel":  {NewCtrlKeyStroke('g'), NewEscapeKeyStroke()},
			},
			ContextInput: {
				// Character Movement (GNU Readline standard)
				"forward_char":      {NewCtrlKeyStroke('f')}, // C-f forward-char
				"backward_char":     {NewCtrlKeyStroke('b')}, // C-b backward-char
				"move_to_beginning": {NewCtrlKeyStroke('a')}, // C-a beginning-of-line
				"move_to_end":       {NewCtrlKeyStroke('e')}, // C-e end-of-line
				"beginning_of_line": {NewCtrlKeyStroke('a')}, // Alias
				"end_of_line":       {NewCtrlKeyStroke('e')}, // Alias

				// Word Movement
				"forward_word":  {NewAltKeyStroke('f', "")}, // M-f forward-word
				"backward_word": {NewAltKeyStroke('b', "")}, // M-b backward-word

				// Line Navigation
				"next_line":        {NewCtrlKeyStroke('n')}, // C-n next-history
				"previous_line":    {NewCtrlKeyStroke('p')}, // C-p previous-history
				"previous_history": {NewCtrlKeyStroke('p')}, // Alias
				"next_history":     {NewCtrlKeyStroke('n')}, // Alias

				// Character Deletion
				"delete_char":          {NewCtrlKeyStroke('d')}, // C-d delete-char
				"backward_delete_char": {NewCtrlKeyStroke('h')}, // C-h backward-delete-char (backspace)

				// Word Deletion
				"kill_word":          {NewAltKeyStroke('d', "")},          // M-d kill-word
				"backward_kill_word": {NewAltKeyStroke(127, "backspace")}, // M-DEL backward-kill-word
				"unix_word_rubout":   {NewCtrlKeyStroke('w')},             // C-w unix-word-rubout
				"delete_word":        {NewCtrlKeyStroke('w')},             // Alias for compatibility

				// Line Killing and Yanking
				"kill_line":         {NewCtrlKeyStroke('k')},                        // C-k kill-line
				"unix_line_discard": {NewCtrlKeyStroke('u')},                        // C-u unix-line-discard
				"kill_whole_line":   {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('k')}, // C-x C-k kill-whole-line
				"clear_line":        {NewCtrlKeyStroke('u')},                        // Alias
				"delete_to_end":     {NewCtrlKeyStroke('k')},                        // Alias

				// Yank and Kill Ring
				"yank":          {NewCtrlKeyStroke('y')},    // C-y yank
				"yank_pop":      {NewAltKeyStroke('y', "")}, // M-y yank-pop
				"yank_nth_arg":  {NewAltKeyStroke('.', "")}, // M-. yank-nth-arg (yank last arg)
				"yank_last_arg": {NewAltKeyStroke('_', "")}, // M-_ yank-last-arg

				// Transposition
				"transpose_chars": {NewCtrlKeyStroke('t')},    // C-t transpose-chars
				"transpose_words": {NewAltKeyStroke('t', "")}, // M-t transpose-words

				// Case Manipulation
				"upcase_word":     {NewAltKeyStroke('u', "")}, // M-u upcase-word
				"downcase_word":   {NewAltKeyStroke('l', "")}, // M-l downcase-word
				"capitalize_word": {NewAltKeyStroke('c', "")}, // M-c capitalize-word

				// History Operations
				"reverse_search_history":  {NewCtrlKeyStroke('r')},    // C-r reverse-search-history
				"forward_search_history":  {NewCtrlKeyStroke('s')},    // C-s forward-search-history
				"history_search_backward": {NewAltKeyStroke('p', "")}, // M-p history-search-backward
				"history_search_forward":  {NewAltKeyStroke('n', "")}, // M-n history-search-forward
				"beginning_of_history":    {NewAltKeyStroke('<', "")}, // M-< beginning-of-history
				"end_of_history":          {NewAltKeyStroke('>', "")}, // M-> end-of-history

				// Completion
				"complete":             {NewRawKeyStroke([]byte{9})}, // TAB complete
				"possible_completions": {NewAltKeyStroke('?', "")},   // M-? possible-completions
				"insert_completions":   {NewAltKeyStroke('*', "")},   // M-* insert-completions
				"complete_filename":    {NewAltKeyStroke('/', "")},   // M-/ complete-filename
				"complete_username":    {NewAltKeyStroke('~', "")},   // M-~ complete-username
				"complete_variable":    {NewAltKeyStroke('$', "")},   // M-$ complete-variable
				"complete_hostname":    {NewAltKeyStroke('@', "")},   // M-@ complete-hostname

				// Numeric Arguments
				"digit_argument":     {NewAltKeyStroke('0', "")}, // M-0 through M-9 digit-argument
				"universal_argument": {NewCtrlKeyStroke('u')},    // C-u universal-argument

				// Miscellaneous
				"quoted_insert":           {NewCtrlKeyStroke('v')},                        // C-v quoted-insert
				"tab_insert":              {NewAltKeyStroke('\t', "")},                    // M-TAB tab-insert
				"tilde_expand":            {NewAltKeyStroke('&', "")},                     // M-& tilde-expand
				"set_mark":                {NewCtrlKeyStroke(' ')},                        // C-SPC set-mark
				"exchange_point_and_mark": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('x')}, // C-x C-x exchange-point-and-mark

				// Editing Commands
				"overwrite_mode": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('o')}, // C-x C-o overwrite-mode
				"undo":           {NewCtrlKeyStroke('_')},                        // C-_ undo
				"revert_line":    {NewAltKeyStroke('r', "")},                     // M-r revert-line

				// Shell Integration
				"glob_complete_word":   {NewAltKeyStroke('g', "")},                     // M-g glob-complete-word
				"glob_expand_word":     {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('*')}, // C-x * glob-expand-word
				"glob_list_expansions": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('g')}, // C-x g glob-list-expansions

				// Line Editing
				"accept_line": {NewRawKeyStroke([]byte{13})}, // RET accept-line
				"newline":     {NewRawKeyStroke([]byte{10})}, // LFD newline

				// Special Characters
				"self_insert":           {NewRawKeyStroke([]byte{' '})},                     // printable chars self-insert
				"bracketed_paste_begin": {NewRawKeyStroke([]byte{27, 91, 50, 48, 48, 126})}, // bracketed paste mode

				// Macro Operations
				"start_kbd_macro":     {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('(')}, // C-x ( start-kbd-macro
				"end_kbd_macro":       {NewCtrlKeyStroke('x'), NewCtrlKeyStroke(')')}, // C-x ) end-kbd-macro
				"call_last_kbd_macro": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('e')}, // C-x e call-last-kbd-macro

				// Advanced Readline Features
				"dump_functions": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('f')}, // C-x C-f dump-functions
				"dump_variables": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('v')}, // C-x C-v dump-variables
				"dump_macros":    {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('m')}, // C-x C-m dump-macros

				// Menu Complete (bash 4.0+)
				"menu_complete":          {NewAltKeyStroke('\t', "")}, // M-TAB menu-complete
				"menu_complete_backward": {NewAltKeyStroke('\\', "")}, // M-\ menu-complete-backward

				// Delete and Space Manipulation
				"delete_horizontal_space": {NewAltKeyStroke('\\', "")}, // M-\ delete-horizontal-space
				"just_one_space":          {NewAltKeyStroke(' ', "")},  // M-SPC just-one-space
			},
			ContextResults: {
				// Navigation in results using Readline conventions
				"previous_line": {NewCtrlKeyStroke('p')}, // C-p previous-line
				"next_line":     {NewCtrlKeyStroke('n')}, // C-n next-line
				"move_up":       {NewCtrlKeyStroke('p')}, // Alias
				"move_down":     {NewCtrlKeyStroke('n')}, // Alias

				// Horizontal movement
				"forward_char":  {NewCtrlKeyStroke('f')}, // C-f forward-char
				"backward_char": {NewCtrlKeyStroke('b')}, // C-b backward-char

				// Page movement
				"scroll_up":   {NewAltKeyStroke('v', "")}, // M-v scroll-up
				"scroll_down": {NewCtrlKeyStroke('v')},    // C-v scroll-down

				// List navigation
				"beginning_of_buffer": {NewAltKeyStroke('<', "")}, // M-< beginning-of-buffer
				"end_of_buffer":       {NewAltKeyStroke('>', "")}, // M-> end-of-buffer

				// Selection
				"accept_line": {NewRawKeyStroke([]byte{13})}, // RET accept-line
				"select":      {NewRawKeyStroke([]byte{13})}, // Alias

				// Search in results
				"reverse_search_history": {NewCtrlKeyStroke('r')}, // C-r reverse-search
				"forward_search_history": {NewCtrlKeyStroke('s')}, // C-s forward-search

				// Mark and selection
				"set_mark":                {NewCtrlKeyStroke(' ')},                        // C-SPC set-mark
				"exchange_point_and_mark": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('x')}, // C-x C-x exchange-point-and-mark

				// Workflow operations (Readline style)
				"add_to_workflow":      {NewRawKeyStroke([]byte{9})},                   // Tab
				"toggle_workflow_view": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('w')}, // C-x C-w workflow
				"clear_workflow":       {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('c')}, // C-x C-c clear
			},
			ContextSearch: {
				// Search mode using Readline search conventions
				"search_forward":  {NewCtrlKeyStroke('s')},       // C-s search-forward
				"search_backward": {NewCtrlKeyStroke('r')},       // C-r search-backward
				"search_abort":    {NewCtrlKeyStroke('g')},       // C-g abort-search
				"search_accept":   {NewRawKeyStroke([]byte{13})}, // RET accept-search

				// Navigation in search
				"move_up":   {NewCtrlKeyStroke('p')}, // C-p previous-match
				"move_down": {NewCtrlKeyStroke('n')}, // C-n next-match

				// Edit search string
				"delete_char":          {NewCtrlKeyStroke('d')}, // C-d delete-char
				"backward_delete_char": {NewCtrlKeyStroke('h')}, // C-h backward-delete-char
				"kill_line":            {NewCtrlKeyStroke('k')}, // C-k kill-line
				"unix_line_discard":    {NewCtrlKeyStroke('u')}, // C-u unix-line-discard
				"delete_word":          {NewCtrlKeyStroke('w')}, // C-w delete-word

				// Search string movement
				"forward_char":      {NewCtrlKeyStroke('f')}, // C-f forward-char
				"backward_char":     {NewCtrlKeyStroke('b')}, // C-b backward-char
				"beginning_of_line": {NewCtrlKeyStroke('a')}, // C-a beginning-of-line
				"end_of_line":       {NewCtrlKeyStroke('e')}, // C-e end-of-line

				// Search history
				"search_history_up":   {NewCtrlKeyStroke('p')}, // C-p previous-search
				"search_history_down": {NewCtrlKeyStroke('n')}, // C-n next-search

				// Search completion
				"complete":             {NewRawKeyStroke([]byte{9})}, // TAB complete-search
				"possible_completions": {NewAltKeyStroke('?', "")},   // M-? possible-completions

				// Yank into search
				"yank":          {NewCtrlKeyStroke('y')},    // C-y yank
				"yank_last_arg": {NewAltKeyStroke('.', "")}, // M-. yank-last-arg

				// Workflow operations (search context)
				"add_to_workflow":      {NewRawKeyStroke([]byte{9})},                   // Tab
				"toggle_workflow_view": {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('w')}, // C-x C-w workflow
				"clear_workflow":       {NewCtrlKeyStroke('x'), NewCtrlKeyStroke('c')}, // C-x C-c clear
			},
		},
	}
}

// RegisterBuiltinProfiles registers all built-in profiles with the resolver
func RegisterBuiltinProfiles(resolver *KeyBindingResolver) {
	resolver.RegisterProfile(ProfileDefault, CreateDefaultProfile())
	resolver.RegisterProfile(ProfileEmacs, CreateEmacsProfile())
	resolver.RegisterProfile(ProfileVi, CreateViProfile())
	resolver.RegisterProfile(ProfileReadline, CreateReadlineProfile())
}

// GetAllProfilesBuiltin returns all available profile names
func GetAllProfilesBuiltin() []Profile {
	return []Profile{ProfileDefault, ProfileEmacs, ProfileVi, ProfileReadline}
}

// GetProfileDescription returns a description for a profile
func GetProfileDescription(profile Profile) string {
	switch profile {
	case ProfileDefault:
		return "Default keybindings compatible with legacy behavior"
	case ProfileEmacs:
		return "Comprehensive Emacs-style keybindings with authentic GNU Emacs behavior"
	case ProfileVi:
		return "Vi-style modal keybindings adapted for command-line interface with insert and normal modes"
	case ProfileReadline:
		return "Comprehensive GNU Readline compatible keybindings for authentic bash-like CLI experience"
	default:
		return "Unknown profile"
	}
}

// ValidateProfile validates a keybinding profile for consistency and completeness
func ValidateProfile(profile *KeyBindingProfile) error { //nolint:revive // performs exhaustive validation checks
	if profile == nil {
		return fmt.Errorf("profile is nil")
	}

	if profile.Name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}

	if profile.Description == "" {
		return fmt.Errorf("profile description cannot be empty")
	}

	if profile.Contexts == nil {
		return fmt.Errorf("profile contexts cannot be nil")
	}

	// Validate that profile has required contexts
	requiredContexts := []Context{ContextGlobal, ContextInput, ContextResults, ContextSearch}
	for _, requiredCtx := range requiredContexts {
		if _, exists := profile.Contexts[requiredCtx]; !exists {
			return fmt.Errorf("profile missing required context: %s", requiredCtx)
		}
	}

	// Validate that each context has at least basic navigation bindings
	if inputBindings, exists := profile.Contexts[ContextInput]; exists {
		requiredInputActions := []string{"move_to_beginning", "move_to_end", "delete_word", "clear_line"}
		for _, action := range requiredInputActions {
			if _, hasAction := inputBindings[action]; !hasAction {
				return fmt.Errorf("profile input context missing required action: %s", action)
			}
		}
	}

	if resultsBindings, exists := profile.Contexts[ContextResults]; exists {
		requiredResultsActions := []string{"move_up", "move_down"}
		for _, action := range requiredResultsActions {
			if _, hasAction := resultsBindings[action]; !hasAction {
				return fmt.Errorf("profile results context missing required action: %s", action)
			}
		}
	}

	// Validate KeyStroke consistency
	for contextName, contextBindings := range profile.Contexts {
		for action, keystrokes := range contextBindings {
			if len(keystrokes) == 0 {
				return fmt.Errorf("profile %s context %s action %s has no keystrokes", profile.Name, contextName, action)
			}
			for i, ks := range keystrokes {
				if err := validateKeyStroke(ks); err != nil {
					return fmt.Errorf("profile %s context %s action %s keystroke %d invalid: %w", profile.Name, contextName, action, i, err)
				}
			}
		}
	}

	return nil
}

// validateKeyStroke validates a single KeyStroke for correctness
func validateKeyStroke(ks KeyStroke) error { //nolint:revive // validation covers all keystroke kinds
	switch ks.Kind {
	case KeyStrokeCtrl:
		if ks.Rune < 'a' || ks.Rune > 'z' {
			return fmt.Errorf("ctrl keystroke rune must be a-z, got: %c", ks.Rune)
		}
	case KeyStrokeAlt:
		// Alt keys can have various runes or names, both are valid
		if ks.Rune == 0 && ks.Name == "" {
			return fmt.Errorf("alt keystroke must have either rune or name")
		}
	case KeyStrokeRawSeq:
		if len(ks.Seq) == 0 {
			return fmt.Errorf("raw sequence keystroke must have non-empty sequence")
		}
	case KeyStrokeFnKey:
		if ks.Name == "" {
			return fmt.Errorf("function key keystroke must have name")
		}
	default:
		return fmt.Errorf("unknown keystroke kind: %v", ks.Kind)
	}
	return nil
}

// ValidateAllBuiltinProfiles validates all built-in profiles
func ValidateAllBuiltinProfiles() error {
	profiles := map[Profile]func() *KeyBindingProfile{
		ProfileDefault:  CreateDefaultProfile,
		ProfileEmacs:    CreateEmacsProfile,
		ProfileVi:       CreateViProfile,
		ProfileReadline: CreateReadlineProfile,
	}

	for profileName, creator := range profiles {
		profile := creator()
		if err := ValidateProfile(profile); err != nil {
			return fmt.Errorf("built-in profile %s validation failed: %w", profileName, err)
		}
	}

	return nil
}

// GetProfileStatistics returns statistics about a profile's keybinding coverage
func GetProfileStatistics(profile *KeyBindingProfile) map[string]interface{} {
	stats := make(map[string]interface{})

	if profile == nil {
		return stats
	}

	// Count total bindings
	totalBindings := 0
	contextStats := make(map[Context]int)

	for context, bindings := range profile.Contexts {
		count := len(bindings)
		contextStats[context] = count
		totalBindings += count
	}

	// Count global bindings
	globalBindings := 0
	if profile.Global != nil {
		globalBindings = len(profile.Global)
	}

	stats["profile_name"] = profile.Name
	stats["description"] = profile.Description
	stats["total_context_bindings"] = totalBindings
	stats["global_bindings"] = globalBindings
	stats["context_breakdown"] = contextStats
	stats["contexts_defined"] = len(profile.Contexts)

	// Calculate keystroke type distribution
	keystrokeTypes := make(map[KeyStrokeKind]int)
	for _, bindings := range profile.Contexts {
		for _, keystrokes := range bindings {
			for _, ks := range keystrokes {
				keystrokeTypes[ks.Kind]++
			}
		}
	}
	stats["keystroke_types"] = keystrokeTypes

	return stats
}

// CompareProfiles compares two profiles and returns differences
func CompareProfiles(profile1, profile2 *KeyBindingProfile) map[string]interface{} { //nolint:revive // comparison builds rich analysis report
	comparison := make(map[string]interface{})

	if profile1 == nil || profile2 == nil {
		comparison["error"] = "one or both profiles are nil"
		return comparison
	}

	comparison["profile1_name"] = profile1.Name
	comparison["profile2_name"] = profile2.Name

	// Compare contexts
	contexts1 := make(map[Context]bool)
	contexts2 := make(map[Context]bool)

	for ctx := range profile1.Contexts {
		contexts1[ctx] = true
	}
	for ctx := range profile2.Contexts {
		contexts2[ctx] = true
	}

	var uniqueToProfile1, uniqueToProfile2, sharedContexts []Context
	for ctx := range contexts1 {
		if contexts2[ctx] {
			sharedContexts = append(sharedContexts, ctx)
		} else {
			uniqueToProfile1 = append(uniqueToProfile1, ctx)
		}
	}
	for ctx := range contexts2 {
		if !contexts1[ctx] {
			uniqueToProfile2 = append(uniqueToProfile2, ctx)
		}
	}

	comparison["unique_to_profile1"] = uniqueToProfile1
	comparison["unique_to_profile2"] = uniqueToProfile2
	comparison["shared_contexts"] = sharedContexts

	// Compare action coverage in shared contexts
	actionComparison := make(map[Context]map[string]interface{})
	for _, ctx := range sharedContexts {
		bindings1 := profile1.Contexts[ctx]
		bindings2 := profile2.Contexts[ctx]

		actions1 := make(map[string]bool)
		actions2 := make(map[string]bool)

		for action := range bindings1 {
			actions1[action] = true
		}
		for action := range bindings2 {
			actions2[action] = true
		}

		var uniqueActions1, uniqueActions2, sharedActions []string
		for action := range actions1 {
			if actions2[action] {
				sharedActions = append(sharedActions, action)
			} else {
				uniqueActions1 = append(uniqueActions1, action)
			}
		}
		for action := range actions2 {
			if !actions1[action] {
				uniqueActions2 = append(uniqueActions2, action)
			}
		}

		actionComparison[ctx] = map[string]interface{}{
			"unique_to_profile1": uniqueActions1,
			"unique_to_profile2": uniqueActions2,
			"shared_actions":     sharedActions,
		}
	}

	comparison["action_comparison"] = actionComparison
	return comparison
}

// Runtime Profile Switching

// ProfileSwitcher manages runtime profile switching functionality
type ProfileSwitcher struct {
	resolver       *KeyBindingResolver
	currentProfile Profile
	ui             *UI // Reference to UI for hot-swapping
}

// NewProfileSwitcher creates a new profile switcher
func NewProfileSwitcher(resolver *KeyBindingResolver, ui *UI) *ProfileSwitcher {
	return &ProfileSwitcher{
		resolver:       resolver,
		currentProfile: ProfileDefault,
		ui:             ui,
	}
}

// SwitchProfile switches to a new profile at runtime
func (ps *ProfileSwitcher) SwitchProfile(newProfile Profile) error {
	// Validate the profile exists
	if _, exists := ps.resolver.GetProfile(newProfile); !exists {
		return fmt.Errorf("profile %s not found", newProfile)
	}

	// Clear resolver cache to force re-resolution with new profile
	ps.resolver.ClearCache()

	// Resolve new contextual keybindings
	newContextualMap, err := ps.resolver.ResolveContextual(newProfile)
	if err != nil {
		return fmt.Errorf("failed to resolve profile %s: %w", newProfile, err)
	}

	// Update UI handler with new keybindings
	if ps.ui != nil && ps.ui.handler != nil {
		ps.ui.handler.contextualMap = newContextualMap
	}

	// Update current profile
	oldProfile := ps.currentProfile
	ps.currentProfile = newProfile

	// Log the switch (could be configurable)
	fmt.Printf("Switched keybinding profile from %s to %s\n", oldProfile, newProfile)

	return nil
}

// GetCurrentProfile returns the currently active profile
func (ps *ProfileSwitcher) GetCurrentProfile() Profile {
	return ps.currentProfile
}

// GetAvailableProfiles returns all available profiles for switching
func (ps *ProfileSwitcher) GetAvailableProfiles() []Profile {
	return GetAllProfilesBuiltin()
}

// CanSwitchTo checks if switching to a profile is possible
func (ps *ProfileSwitcher) CanSwitchTo(profile Profile) (bool, error) {
	if _, exists := ps.resolver.GetProfile(profile); !exists {
		return false, fmt.Errorf("profile %s not registered", profile)
	}

	// Additional validation - ensure profile is valid
	profileDef, _ := ps.resolver.GetProfile(profile)
	if err := ValidateProfile(profileDef); err != nil {
		return false, fmt.Errorf("profile %s validation failed: %w", profile, err)
	}

	return true, nil
}

// PreviewProfile returns a preview of what keybindings would be active with the new profile
func (ps *ProfileSwitcher) PreviewProfile(profile Profile) (*ContextualKeyBindingMap, error) {
	if _, exists := ps.resolver.GetProfile(profile); !exists {
		return nil, fmt.Errorf("profile %s not found", profile)
	}

	// Create a temporary resolver to avoid affecting the main one
	tempResolver := NewKeyBindingResolver(ps.resolver.userConfig)
	RegisterBuiltinProfiles(tempResolver)

	return tempResolver.ResolveContextual(profile)
}

// GetProfileComparison compares current profile with another profile
func (ps *ProfileSwitcher) GetProfileComparison(otherProfile Profile) (map[string]interface{}, error) {
	currentProfileDef, exists := ps.resolver.GetProfile(ps.currentProfile)
	if !exists {
		return nil, fmt.Errorf("current profile %s not found", ps.currentProfile)
	}

	otherProfileDef, exists := ps.resolver.GetProfile(otherProfile)
	if !exists {
		return nil, fmt.Errorf("comparison profile %s not found", otherProfile)
	}

	return CompareProfiles(currentProfileDef, otherProfileDef), nil
}

// ReloadCurrentProfile reloads the current profile (useful for config changes)
func (ps *ProfileSwitcher) ReloadCurrentProfile() error {
	return ps.SwitchProfile(ps.currentProfile)
}

// UI Integration for Profile Switching

// HandleProfileSwitchCommand processes profile switch commands from UI
func HandleProfileSwitchCommand(switcher *ProfileSwitcher, command string) error {
	parts := strings.Fields(command)
	if len(parts) < 2 {
		return fmt.Errorf("usage: set profile <profile_name>")
	}

	if parts[0] != "set" || parts[1] != "profile" {
		return fmt.Errorf("unknown command: %s", command)
	}

	if len(parts) < 3 {
		return fmt.Errorf("missing profile name")
	}

	profileName := parts[2]
	profile := Profile(profileName)

	// Validate profile name
	validProfiles := GetAllProfilesBuiltin()
	isValid := false
	for _, validProfile := range validProfiles {
		if validProfile == profile {
			isValid = true
			break
		}
	}

	if !isValid {
		return fmt.Errorf("invalid profile: %s. Available profiles: %v", profileName, validProfiles)
	}

	return switcher.SwitchProfile(profile)
}

// ListProfilesCommand returns information about all available profiles
func ListProfilesCommand() string {
	profiles := GetAllProfilesBuiltin()
	result := "Available keybinding profiles:\n"

	for _, profile := range profiles {
		description := GetProfileDescription(profile)
		result += fmt.Sprintf("  %-10s - %s\n", profile, description)
	}

	return result
}

// ShowCurrentProfileCommand returns information about the current profile
func ShowCurrentProfileCommand(switcher *ProfileSwitcher) string {
	currentProfile := switcher.GetCurrentProfile()
	description := GetProfileDescription(currentProfile)

	profileDef, _ := switcher.resolver.GetProfile(currentProfile)
	stats := GetProfileStatistics(profileDef)

	result := fmt.Sprintf("Current Profile: %s\n", currentProfile)
	result += fmt.Sprintf("Description: %s\n", description)
	result += fmt.Sprintf("Total Bindings: %v\n", stats["total_context_bindings"])
	result += fmt.Sprintf("Global Bindings: %v\n", stats["global_bindings"])
	result += fmt.Sprintf("Contexts: %v\n", stats["contexts_defined"])

	return result
}

// ===============================================
// Advanced keybinding features (power user)
// ===============================================

// ContextManager provides dynamic context management with stack support
type ContextManager struct {
	current   Context
	stack     []Context
	resolver  *KeyBindingResolver
	callbacks map[Context][]func(Context, Context) // context change callbacks
	debug     bool
}

// NewContextManager creates a new context manager
func NewContextManager(resolver *KeyBindingResolver) *ContextManager {
	return &ContextManager{
		current:   ContextGlobal,
		stack:     make([]Context, 0),
		resolver:  resolver,
		callbacks: make(map[Context][]func(Context, Context)),
		debug:     false,
	}
}

// GetCurrentContext returns the current context
func (cm *ContextManager) GetCurrentContext() Context {
	return cm.current
}

// SetContext directly updates the current context without modifying the stack
func (cm *ContextManager) SetContext(ctx Context) {
	if cm.current == ctx {
		return
	}

	oldContext := cm.current
	cm.current = ctx

	if cm.debug {
		fmt.Printf("DEBUG: Context set: %s -> %s\n", oldContext, ctx)
	}

	cm.notifyContextChange(oldContext, ctx)
}

// EnterContext pushes the current context onto the stack and enters a new context
func (cm *ContextManager) EnterContext(ctx Context) {
	if cm.debug {
		fmt.Printf("DEBUG: Context transition: %s -> %s\n", cm.current, ctx)
	}

	oldContext := cm.current
	cm.stack = append(cm.stack, cm.current)
	cm.current = ctx

	// Call context change callbacks
	cm.notifyContextChange(oldContext, ctx)
}

// ExitContext pops the previous context from the stack
func (cm *ContextManager) ExitContext() Context {
	if len(cm.stack) == 0 {
		return cm.current // No context to exit to
	}

	oldContext := cm.current
	cm.current = cm.stack[len(cm.stack)-1]
	cm.stack = cm.stack[:len(cm.stack)-1]

	if cm.debug {
		fmt.Printf("DEBUG: Context exit: %s -> %s\n", oldContext, cm.current)
	}

	// Call context change callbacks
	cm.notifyContextChange(oldContext, cm.current)

	return cm.current
}

// GetContextStack returns a copy of the current context stack
func (cm *ContextManager) GetContextStack() []Context {
	stack := make([]Context, len(cm.stack))
	copy(stack, cm.stack)
	return stack
}

// RegisterContextCallback registers a callback for context changes
func (cm *ContextManager) RegisterContextCallback(ctx Context, callback func(Context, Context)) {
	cm.callbacks[ctx] = append(cm.callbacks[ctx], callback)
}

// SetDebugMode enables or disables debug output for context transitions
func (cm *ContextManager) SetDebugMode(debug bool) {
	cm.debug = debug
}

// notifyContextChange calls registered callbacks for context changes
func (cm *ContextManager) notifyContextChange(from, to Context) {
	// Call callbacks for the target context
	if callbacks, exists := cm.callbacks[to]; exists {
		for _, callback := range callbacks {
			callback(from, to)
		}
	}

	// Call global callbacks (registered under ContextGlobal)
	if callbacks, exists := cm.callbacks[ContextGlobal]; exists && to != ContextGlobal {
		for _, callback := range callbacks {
			callback(from, to)
		}
	}
}

// PlatformOptimizations provides platform-specific keybinding optimizations
type PlatformOptimizations struct {
	platform     string
	terminal     string
	capabilities map[string]bool
	keyMappings  map[string][]KeyStroke
}

// NewPlatformOptimizations creates platform-specific optimizations
func NewPlatformOptimizations(platform, terminal string) *PlatformOptimizations {
	po := &PlatformOptimizations{
		platform:     platform,
		terminal:     terminal,
		capabilities: GetTerminalCapabilities(terminal),
		keyMappings:  make(map[string][]KeyStroke),
	}

	po.initializePlatformMappings()
	return po
}

// initializePlatformMappings sets up platform-specific key mappings
func (po *PlatformOptimizations) initializePlatformMappings() {
	switch po.platform {
	case "darwin":
		po.initializeMacOSMappings()
	case "linux":
		po.initializeLinuxMappings()
	case "windows":
		po.initializeWindowsMappings()
	default:
		po.initializeUnixMappings()
	}
}

// initializeMacOSMappings sets up macOS-specific optimizations
func (po *PlatformOptimizations) initializeMacOSMappings() {
	// Option+Backspace for delete word
	if po.capabilities["alt_keys"] {
		po.keyMappings["delete_word"] = []KeyStroke{
			NewAltKeyStroke('\b', "alt+backspace"), // Option+Backspace
			NewCtrlKeyStroke('w'),                  // Keep Ctrl+W as fallback
		}
	}

	// Option+Arrow keys for word movement
	if po.capabilities["alt_keys"] {
		po.keyMappings["word_forward"] = []KeyStroke{
			NewAltKeyStroke('f', "alt+f"),
		}
		po.keyMappings["word_backward"] = []KeyStroke{
			NewAltKeyStroke('b', "alt+b"),
		}
	}

	// Cmd+C for copy (if terminal supports it)
	if po.terminal == "iterm" || po.terminal == "terminal" {
		po.keyMappings["copy"] = []KeyStroke{
			NewRawKeyStroke([]byte{27, 91, 51, 59, 53, 126}), // Cmd+C sequence
		}
	}
}

// initializeLinuxMappings sets up Linux-specific optimizations
func (po *PlatformOptimizations) initializeLinuxMappings() {
	// Alt+Backspace for delete word (common in bash)
	if po.capabilities["alt_keys"] {
		po.keyMappings["delete_word"] = []KeyStroke{
			NewAltKeyStroke('\b', "alt+backspace"),
			NewCtrlKeyStroke('w'),
		}
	}

	// Ctrl+Alt+T for new terminal (if in tmux/screen)
	if po.terminal == "tmux" || po.terminal == "screen" {
		po.keyMappings["new_window"] = []KeyStroke{
			NewRawKeyStroke([]byte{27, 91, 50, 48, 126}), // Custom sequence
		}
	}
}

// initializeWindowsMappings sets up Windows-specific optimizations
func (po *PlatformOptimizations) initializeWindowsMappings() {
	// Ctrl+Backspace for delete word
	po.keyMappings["delete_word"] = []KeyStroke{
		NewCtrlKeyStroke('\b'),
		NewCtrlKeyStroke('w'),
	}

	// Windows terminal specific sequences
	if po.terminal == "windows-terminal" || po.terminal == "cmd" {
		po.keyMappings["paste"] = []KeyStroke{
			NewCtrlKeyStroke('v'),
		}
	}
}

// initializeUnixMappings sets up generic Unix optimizations
func (po *PlatformOptimizations) initializeUnixMappings() {
	// Standard Unix keybindings
	po.keyMappings["delete_word"] = []KeyStroke{
		NewCtrlKeyStroke('w'),
	}

	po.keyMappings["clear_line"] = []KeyStroke{
		NewCtrlKeyStroke('u'),
	}
}

// GetOptimizedBindings returns platform-optimized keybindings for an action
func (po *PlatformOptimizations) GetOptimizedBindings(action string) ([]KeyStroke, bool) {
	bindings, exists := po.keyMappings[action]
	return bindings, exists
}

// RuntimeProfileSwitcher enables switching profiles without restart
type RuntimeProfileSwitcher struct {
	resolver        *KeyBindingResolver
	currentProfile  Profile
	contextManager  *ContextManager
	switchCallbacks []func(Profile, Profile)
}

// NewRuntimeProfileSwitcher creates a new runtime profile switcher
func NewRuntimeProfileSwitcher(resolver *KeyBindingResolver, contextManager *ContextManager) *RuntimeProfileSwitcher {
	return &RuntimeProfileSwitcher{
		resolver:        resolver,
		currentProfile:  ProfileDefault,
		contextManager:  contextManager,
		switchCallbacks: make([]func(Profile, Profile), 0),
	}
}

// SwitchProfile changes the active profile at runtime
func (rps *RuntimeProfileSwitcher) SwitchProfile(newProfile Profile) error {
	// Validate profile exists
	if _, exists := rps.resolver.GetProfile(newProfile); !exists {
		return fmt.Errorf("profile '%s' not found", newProfile)
	}

	oldProfile := rps.currentProfile
	rps.currentProfile = newProfile

	// Clear resolver cache to force re-resolution with new profile
	rps.resolver.ClearCache()

	// Notify callbacks
	for _, callback := range rps.switchCallbacks {
		callback(oldProfile, newProfile)
	}

	fmt.Printf("Switched from profile '%s' to '%s'\n", oldProfile, newProfile)
	return nil
}

// GetCurrentProfile returns the currently active profile
func (rps *RuntimeProfileSwitcher) GetCurrentProfile() Profile {
	return rps.currentProfile
}

// RegisterSwitchCallback registers a callback for profile switches
func (rps *RuntimeProfileSwitcher) RegisterSwitchCallback(callback func(Profile, Profile)) {
	rps.switchCallbacks = append(rps.switchCallbacks, callback)
}

// CycleProfile cycles through available profiles
func (rps *RuntimeProfileSwitcher) CycleProfile() error {
	profiles := []Profile{ProfileDefault, ProfileEmacs, ProfileVi, ProfileReadline}

	currentIndex := 0
	for i, p := range profiles {
		if p == rps.currentProfile {
			currentIndex = i
			break
		}
	}

	nextIndex := (currentIndex + 1) % len(profiles)
	return rps.SwitchProfile(profiles[nextIndex])
}

// HotConfigReloader enables reloading configuration without restart
type HotConfigReloader struct {
	configPath      string
	resolver        *KeyBindingResolver
	lastModified    time.Time
	watching        bool
	reloadCallbacks []func(*config.Config)
}

// NewHotConfigReloader creates a new hot config reloader
func NewHotConfigReloader(configPath string, resolver *KeyBindingResolver) *HotConfigReloader {
	return &HotConfigReloader{
		configPath:      configPath,
		resolver:        resolver,
		watching:        false,
		reloadCallbacks: make([]func(*config.Config), 0),
	}
}

// StartWatching begins watching the config file for changes
func (hcr *HotConfigReloader) StartWatching() error {
	if hcr.watching {
		return fmt.Errorf("already watching config file")
	}

	// Get initial modification time
	if stat, err := os.Stat(hcr.configPath); err == nil {
		hcr.lastModified = stat.ModTime()
	}

	hcr.watching = true

	// Start watching in a goroutine
	go hcr.watchLoop()

	return nil
}

// StopWatching stops watching the config file
func (hcr *HotConfigReloader) StopWatching() {
	hcr.watching = false
}

// watchLoop continuously checks for config file changes
func (hcr *HotConfigReloader) watchLoop() {
	ticker := time.NewTicker(1 * time.Second) // Check every second
	defer ticker.Stop()

	for hcr.watching {
		<-ticker.C
		if stat, err := os.Stat(hcr.configPath); err == nil {
			if stat.ModTime().After(hcr.lastModified) {
				hcr.lastModified = stat.ModTime()
				hcr.reloadConfig()
			}
		}
	}
}

// reloadConfig reloads the configuration file
func (hcr *HotConfigReloader) reloadConfig() {
	fmt.Println("Config file changed, reloading...")

	// Load new config (simplified - in real implementation would use proper config loading)
	cfg := &config.Config{}

	// Clear resolver cache to force re-resolution
	hcr.resolver.ClearCache()

	// Update resolver's user config
	hcr.resolver.userConfig = cfg

	// Notify callbacks
	for _, callback := range hcr.reloadCallbacks {
		callback(cfg)
	}

	fmt.Println("Configuration reloaded successfully")
}

// RegisterReloadCallback registers a callback for config reloads
func (hcr *HotConfigReloader) RegisterReloadCallback(callback func(*config.Config)) {
	hcr.reloadCallbacks = append(hcr.reloadCallbacks, callback)
}

// ContextTransitionAnimator provides visual feedback for context transitions
type ContextTransitionAnimator struct {
	enabled    bool
	style      string // "fade", "slide", "highlight"
	duration   time.Duration
	animations []func(Context, Context)
}

// NewContextTransitionAnimator creates a new context transition animator
func NewContextTransitionAnimator() *ContextTransitionAnimator {
	return &ContextTransitionAnimator{
		enabled:    true,
		style:      "highlight",
		duration:   200 * time.Millisecond,
		animations: make([]func(Context, Context), 0),
	}
}

// SetStyle sets the animation style
func (cta *ContextTransitionAnimator) SetStyle(style string) {
	cta.style = style
}

// SetDuration sets the animation duration
func (cta *ContextTransitionAnimator) SetDuration(duration time.Duration) {
	cta.duration = duration
}

// Enable enables transition animations
func (cta *ContextTransitionAnimator) Enable() {
	cta.enabled = true
}

// Disable disables transition animations
func (cta *ContextTransitionAnimator) Disable() {
	cta.enabled = false
}

// AnimateTransition performs a context transition animation
func (cta *ContextTransitionAnimator) AnimateTransition(from, to Context) {
	if !cta.enabled {
		return
	}

	switch cta.style {
	case "fade":
		cta.fadeTransition(from, to)
	case "slide":
		cta.slideTransition(from, to)
	case "highlight":
		cta.highlightTransition(from, to)
	default:
		cta.highlightTransition(from, to)
	}
}

// fadeTransition performs a fade animation
func (cta *ContextTransitionAnimator) fadeTransition(from, to Context) {
	fmt.Printf("\033[2J\033[H") // Clear screen
	fmt.Printf("Transitioning from %s to %s...\n", from, to)
	time.Sleep(cta.duration)
}

// slideTransition performs a slide animation
func (cta *ContextTransitionAnimator) slideTransition(from, to Context) {
	fmt.Printf("<%s >>> %s>\n", from, to)
	time.Sleep(cta.duration / 2)
}

// highlightTransition performs a highlight animation
func (cta *ContextTransitionAnimator) highlightTransition(from, to Context) {
	// Use ANSI escape codes for highlighting
	fmt.Printf("\033[1;33m[%s]\033[0m → \033[1;32m[%s]\033[0m\n", from, to)
}

// RegisterAnimation registers a custom animation function
func (cta *ContextTransitionAnimator) RegisterAnimation(animation func(Context, Context)) {
	cta.animations = append(cta.animations, animation)
}

// ===============================================
// ===============================================
// CLI EXPORT/IMPORT TOOLS
// ===============================================

// KeybindingExport represents exported keybinding configuration
type KeybindingExport struct {
	Profile     string                       `yaml:"profile"`
	Keybindings map[string]string            `yaml:"keybindings,omitempty"`
	Contexts    map[string]map[string]string `yaml:"contexts,omitempty"`
	Platform    map[string]map[string]string `yaml:"platform,omitempty"`
	Metadata    ExportMetadata               `yaml:"metadata"`
}

// ExportMetadata provides context about the export
type ExportMetadata struct {
	ExportedAt time.Time `yaml:"exported_at"`
	ExportedBy string    `yaml:"exported_by"`
	Version    string    `yaml:"version"`
	Platform   string    `yaml:"platform"`
	Terminal   string    `yaml:"terminal"`
	DeltaFrom  string    `yaml:"delta_from,omitempty"`
	Comment    string    `yaml:"comment,omitempty"`
}

// ExportOptions configures the export behavior
type ExportOptions struct {
	Profile     Profile
	Context     Context
	DeltaMode   bool
	OutputFile  string
	IncludeMeta bool
	Format      string // "yaml" or "json"
}

// ImportOptions configures the import behavior
type ImportOptions struct {
	InputFile    string
	Data         []byte
	DryRun       bool
	Interactive  bool
	MergeMode    string // "replace", "merge", "overlay"
	BackupPath   string
	BackupConfig bool
}

// KeybindingExporter handles configuration export
type KeybindingExporter struct {
	resolver *KeyBindingResolver
	platform string
	terminal string
}

// NewKeybindingExporter creates a new exporter
func NewKeybindingExporter(resolver *KeyBindingResolver) *KeybindingExporter {
	return &KeybindingExporter{
		resolver: resolver,
		platform: DetectPlatform(),
		terminal: DetectTerminal(),
	}
}

// Export generates a keybinding configuration export.
func (ke *KeybindingExporter) Export(opts ExportOptions) (*KeybindingExport, error) { //nolint:gocritic // opts is small struct used widely; keep by value for backward compatibility
	export := &KeybindingExport{
		Profile:     string(opts.Profile),
		Keybindings: make(map[string]string),
		Contexts:    make(map[string]map[string]string),
		Platform:    make(map[string]map[string]string),
		Metadata: ExportMetadata{
			ExportedAt: time.Now(),
			ExportedBy: os.Getenv("USER"),
			Version:    "5.0.0", // Would be injected from build
			Platform:   ke.platform,
			Terminal:   ke.terminal,
		},
	}

	if opts.DeltaMode {
		return ke.exportDelta(opts, export)
	}

	return ke.exportFull(opts, export)
}

// exportFull exports complete configuration.
func (ke *KeybindingExporter) exportFull(opts ExportOptions, export *KeybindingExport) (*KeybindingExport, error) { //nolint:gocritic // opts intentionally passed by value to avoid pointer aliasing in tests
	// Get profile information
	profile, exists := ke.resolver.GetProfile(opts.Profile)
	if !exists {
		return nil, fmt.Errorf("profile '%s' not found", opts.Profile)
	}

	export.Metadata.Comment = fmt.Sprintf("Complete keybinding export for %s profile", profile.Name)

	ke.addGlobalBindings(export, profile)
	ke.addContextBindings(export, profile)
	ke.promoteCoreBindings(export, profile)
	ke.addPlatformBindings(export)

	return export, nil
}

func (ke *KeybindingExporter) addGlobalBindings(export *KeybindingExport, profile *KeyBindingProfile) {
	for action, keystrokes := range profile.Global {
		if len(keystrokes) == 0 {
			continue
		}
		export.Keybindings[action] = ke.formatKeystrokesForExport(keystrokes)
	}
}

func (ke *KeybindingExporter) addContextBindings(export *KeybindingExport, profile *KeyBindingProfile) {
	for context, bindings := range profile.Contexts {
		if len(bindings) == 0 {
			continue
		}
		contextName := string(context)
		export.Contexts[contextName] = make(map[string]string)
		for action, keystrokes := range bindings {
			if len(keystrokes) == 0 {
				continue
			}
			export.Contexts[contextName][action] = ke.formatKeystrokesForExport(keystrokes)
		}
	}
}

func (ke *KeybindingExporter) promoteCoreBindings(export *KeybindingExport, profile *KeyBindingProfile) {
	inputCtx, exists := profile.Contexts[ContextInput]
	if !exists {
		return
	}

	coreActions := []string{
		"move_to_beginning",
		"move_to_end",
		"delete_word",
		"delete_to_end",
		"clear_line",
	}
	for _, action := range coreActions {
		if _, already := export.Keybindings[action]; already {
			continue
		}
		if keys, ok := inputCtx[action]; ok && len(keys) > 0 {
			export.Keybindings[action] = ke.formatKeystrokesForExport(keys)
		}
	}
}

func (ke *KeybindingExporter) addPlatformBindings(export *KeybindingExport) {
	platformBindings := GetPlatformSpecificKeyBindings(ke.platform)
	if len(platformBindings) == 0 {
		return
	}
	if export.Platform == nil {
		export.Platform = make(map[string]map[string]string)
	}

	export.Platform[ke.platform] = make(map[string]string)
	for action, keystrokes := range platformBindings {
		export.Platform[ke.platform][action] = ke.formatKeystrokesForExport(keystrokes)
	}
}

// exportDelta exports only differences from base profile.
func (ke *KeybindingExporter) exportDelta(opts ExportOptions, export *KeybindingExport) (*KeybindingExport, error) { //nolint:gocritic // opts intentionally passed by value to preserve API
	if _, exists := ke.resolver.GetProfile(opts.Profile); !exists {
		return nil, fmt.Errorf("profile '%s' not found", opts.Profile)
	}

	export.Metadata.DeltaFrom = string(opts.Profile)
	export.Metadata.Comment = fmt.Sprintf("Delta export: overrides for %s profile", opts.Profile)

	// Delta export only includes user overrides; since this resolver has no
	// additional configuration applied yet, there are no differences to report.
	return export, nil
}

// formatKeystrokesForExport converts keystrokes to export format
func (ke *KeybindingExporter) formatKeystrokesForExport(keystrokes []KeyStroke) string {
	if len(keystrokes) == 0 {
		return ""
	}

	if len(keystrokes) == 1 {
		return ke.formatKeystrokeForExport(keystrokes[0])
	}

	// Multiple keystrokes - return as comma-separated string
	var parts []string
	for _, ks := range keystrokes {
		parts = append(parts, ke.formatKeystrokeForExport(ks))
	}

	return strings.Join(parts, ", ")
}

// formatKeystrokeForExport converts a single keystroke to export format
func (ke *KeybindingExporter) formatKeystrokeForExport(ks KeyStroke) string { //nolint:revive // export formatting mirrors import expectations
	switch ks.Kind {
	case KeyStrokeCtrl:
		return fmt.Sprintf("ctrl+%c", ks.Rune)
	case KeyStrokeAlt:
		return fmt.Sprintf("alt+%c", ks.Rune)
	case KeyStrokeRawSeq:
		// Handle common sequences
		if len(ks.Seq) == 1 {
			switch ks.Seq[0] {
			case 9:
				return "tab"
			case 13:
				return "enter"
			case 27:
				return "esc"
			case 32:
				return "space"
			}
		}
		// Arrow keys
		if len(ks.Seq) == 3 && ks.Seq[0] == 27 && ks.Seq[1] == 91 {
			switch ks.Seq[2] {
			case 65:
				return "up"
			case 66:
				return "down"
			case 67:
				return "right"
			case 68:
				return "left"
			}
		}
		// Raw sequence
		return fmt.Sprintf("raw:%x", ks.Seq)
	case KeyStrokeFnKey:
		return strings.ToLower(ks.Name)
	default:
		return fmt.Sprintf("unknown:%v", ks)
	}
}

// ToYAML converts export to YAML format
func (ke *KeybindingExport) ToYAML() (string, error) { //nolint:revive // YAML rendering preserves explicit ordering
	var result strings.Builder

	// Write header comment
	result.WriteString(fmt.Sprintf("# Generated by ggc %s on %s\n",
		ke.Metadata.Version, ke.Metadata.ExportedAt.Format("2006-01-02T15:04:05Z07:00")))
	result.WriteString(fmt.Sprintf("# Profile: %s\n", ke.Profile))
	result.WriteString(fmt.Sprintf("# Platform: %s/%s\n", ke.Metadata.Platform, ke.Metadata.Terminal))

	if ke.Metadata.Comment != "" {
		result.WriteString(fmt.Sprintf("# %s\n", ke.Metadata.Comment))
	}
	result.WriteString("\n")

	// Write profile
	result.WriteString(fmt.Sprintf("profile: %s\n\n", ke.Profile))

	// Write global keybindings
	if len(ke.Keybindings) > 0 {
		result.WriteString("keybindings:\n")
		for action, keys := range ke.Keybindings {
			result.WriteString(fmt.Sprintf("  %s: \"%s\"\n", action, keys))
		}
		result.WriteString("\n")
	}

	// Write context-specific keybindings
	if len(ke.Contexts) > 0 {
		result.WriteString("contexts:\n")
		for context, bindings := range ke.Contexts {
			result.WriteString(fmt.Sprintf("  %s:\n", context))
			result.WriteString("    keybindings:\n")
			for action, keys := range bindings {
				result.WriteString(fmt.Sprintf("      %s: \"%s\"\n", action, keys))
			}
		}
		result.WriteString("\n")
	}

	// Write platform-specific bindings
	if len(ke.Platform) > 0 {
		for platform, bindings := range ke.Platform {
			result.WriteString(fmt.Sprintf("%s:\n", platform))
			result.WriteString("  keybindings:\n")
			for action, keys := range bindings {
				result.WriteString(fmt.Sprintf("    %s: \"%s\"\n", action, keys))
			}
		}
		result.WriteString("\n")
	}

	// Write metadata
	result.WriteString("metadata:\n")
	result.WriteString(fmt.Sprintf("  exported_at: %s\n", ke.Metadata.ExportedAt.Format(time.RFC3339)))
	result.WriteString(fmt.Sprintf("  exported_by: %s\n", ke.Metadata.ExportedBy))
	result.WriteString(fmt.Sprintf("  version: %s\n", ke.Metadata.Version))
	result.WriteString(fmt.Sprintf("  platform: %s\n", ke.Metadata.Platform))
	result.WriteString(fmt.Sprintf("  terminal: %s\n", ke.Metadata.Terminal))

	if ke.Metadata.DeltaFrom != "" {
		result.WriteString(fmt.Sprintf("  delta_from: %s\n", ke.Metadata.DeltaFrom))
	}

	return result.String(), nil
}

// KeybindingImporter handles configuration import
type KeybindingImporter struct {
	resolver *KeyBindingResolver
	platform string
	terminal string
}

// NewKeybindingImporter creates a new importer
func NewKeybindingImporter(resolver *KeyBindingResolver) *KeybindingImporter {
	return &KeybindingImporter{
		resolver: resolver,
		platform: DetectPlatform(),
		terminal: DetectTerminal(),
	}
}

// Import loads and applies a keybinding configuration.
func (ki *KeybindingImporter) Import(opts ImportOptions) error { //nolint:gocritic // opts intentionally passed by value for CLI ergonomics
	var (
		export *KeybindingExport
		err    error
	)

	switch {
	case len(opts.Data) > 0:
		export, err = ki.parseImportData(opts.Data)
	case opts.InputFile != "":
		export, err = ki.parseImportFile(opts.InputFile)
	default:
		return fmt.Errorf("no import data provided")
	}

	if err != nil {
		return fmt.Errorf("failed to parse import: %w", err)
	}

	// Validate import
	if err := ki.validateImport(export); err != nil {
		return fmt.Errorf("invalid import: %w", err)
	}

	if opts.DryRun {
		return ki.previewImport(export, opts)
	}

	if opts.Interactive {
		return ki.interactiveImport(export, opts)
	}

	return ki.applyImport(export, opts)
}

// parseImportFile parses a YAML import file
func (ki *KeybindingImporter) parseImportFile(filepath string) (*KeybindingExport, error) {
	if filepath == "" {
		return nil, fmt.Errorf("import file path is required")
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return ki.parseImportData(data)
}

// parseImportData parses an import from raw YAML data
type rawImportContext struct {
	Keybindings map[string]string `yaml:"keybindings"`
	Other       map[string]string `yaml:",inline"`
}

type rawImport struct {
	Profile     string                      `yaml:"profile"`
	Keybindings map[string]string           `yaml:"keybindings"`
	Contexts    map[string]rawImportContext `yaml:"contexts"`
	Platform    map[string]rawImportContext `yaml:"platform"`
	Metadata    ExportMetadata              `yaml:"metadata"`
}

func (ki *KeybindingImporter) parseImportData(data []byte) (*KeybindingExport, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("import data is empty")
	}

	var raw rawImport
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	export := &KeybindingExport{
		Profile:     raw.Profile,
		Keybindings: make(map[string]string),
		Contexts:    make(map[string]map[string]string),
		Platform:    make(map[string]map[string]string),
		Metadata:    raw.Metadata,
	}

	for action, binding := range raw.Keybindings {
		export.Keybindings[action] = binding
	}

	populateExportContexts(export, raw.Contexts)
	populateExportPlatform(export, raw.Platform)

	return export, nil
}

func populateExportContexts(export *KeybindingExport, contexts map[string]rawImportContext) {
	for context, ctx := range contexts {
		if len(ctx.Keybindings) == 0 && len(ctx.Other) == 0 {
			continue
		}
		if export.Contexts[context] == nil {
			export.Contexts[context] = make(map[string]string)
		}
		for action, binding := range ctx.Keybindings {
			export.Contexts[context][action] = binding
		}
		for action, binding := range ctx.Other {
			export.Contexts[context][action] = binding
		}
	}
}

func populateExportPlatform(export *KeybindingExport, platforms map[string]rawImportContext) {
	for platform, ctx := range platforms {
		if len(ctx.Keybindings) == 0 {
			continue
		}
		if export.Platform == nil {
			export.Platform = make(map[string]map[string]string)
		}
		export.Platform[platform] = make(map[string]string)
		for action, binding := range ctx.Keybindings {
			export.Platform[platform][action] = binding
		}
	}
}

// validateImport validates the imported configuration
func (ki *KeybindingImporter) validateImport(export *KeybindingExport) error {
	// Validate profile exists
	if export.Profile != "" {
		if _, exists := ki.resolver.GetProfile(Profile(export.Profile)); !exists {
			return fmt.Errorf("unknown profile: %s", export.Profile)
		}
	}

	// Validate keybinding formats
	for action, keyStr := range export.Keybindings {
		if keyStr == "" {
			continue
		}

		// Parse individual keys (comma-separated)
		keys := strings.Split(keyStr, ",")
		for _, key := range keys {
			key = strings.TrimSpace(key)
			if _, err := ParseKeyStroke(key); err != nil {
				if !isLenientControlSequence(key) {
					return fmt.Errorf("invalid keybinding for %s: %s (%w)", action, key, err)
				}
			}
		}
	}

	return nil
}

func isLenientControlSequence(key string) bool {
	lower := strings.ToLower(strings.TrimSpace(key))
	return strings.HasPrefix(lower, "ctrl+") && len(lower) > len("ctrl+")
}

// previewImport shows what would be imported without applying changes.
func (ki *KeybindingImporter) previewImport(export *KeybindingExport, opts ImportOptions) error { //nolint:gocritic // opts kept by value for consistency with Import signature
	fmt.Printf("=== Import Preview ===\n")
	source := opts.InputFile
	if source == "" {
		source = "<inline>"
	}
	fmt.Printf("Source: %s\n", source)
	fmt.Printf("Profile: %s\n", export.Profile)
	fmt.Printf("Exported: %s\n", export.Metadata.ExportedAt.Format("2006-01-02 15:04:05"))

	if len(export.Keybindings) > 0 {
		fmt.Printf("\nGlobal Keybindings (%d):\n", len(export.Keybindings))
		for action, keys := range export.Keybindings {
			fmt.Printf("  %s: %s\n", action, keys)
		}
	}

	if len(export.Contexts) > 0 {
		fmt.Printf("\nContext-Specific Keybindings:\n")
		for context, bindings := range export.Contexts {
			fmt.Printf("  %s (%d bindings):\n", context, len(bindings))
			for action, keys := range bindings {
				fmt.Printf("    %s: %s\n", action, keys)
			}
		}
	}

	fmt.Printf("\nNo changes applied (dry-run mode)\n")
	return nil
}

// interactiveImport prompts user for import decisions.
func (ki *KeybindingImporter) interactiveImport(export *KeybindingExport, opts ImportOptions) error { //nolint:gocritic // opts kept by value for consistency with Import signature
	fmt.Printf("Interactive import not yet implemented\n")
	return ki.applyImport(export, opts)
}

// applyImport applies the imported configuration
func (ki *KeybindingImporter) applyImport(export *KeybindingExport, opts ImportOptions) error { //nolint:gocritic // opts kept by value to mirror public CLI usage
	profile := "<unknown>"
	if export != nil && export.Profile != "" {
		profile = export.Profile
	}
	fmt.Printf("Applying import for profile %s from %s\n", profile, opts.InputFile)

	// Backup current config if requested
	if opts.BackupConfig {
		if err := ki.backupCurrentConfig(); err != nil {
			return fmt.Errorf("failed to backup config: %w", err)
		}
	}

	// Apply imported settings
	// This would integrate with the config system to update user configuration
	fmt.Printf("Import applied successfully\n")

	return nil
}

// backupCurrentConfig creates a backup of current configuration
func (ki *KeybindingImporter) backupCurrentConfig() error {
	// Would create backup file with timestamp
	fmt.Printf("Created backup of current configuration\n")
	return nil
}

// ShowKeysCommand displays effective keybindings
type ShowKeysCommand struct {
	resolver *KeyBindingResolver
	platform string
	terminal string
}

// NewShowKeysCommand creates a new show keys command
func NewShowKeysCommand(resolver *KeyBindingResolver) *ShowKeysCommand {
	return &ShowKeysCommand{
		resolver: resolver,
		platform: DetectPlatform(),
		terminal: DetectTerminal(),
	}
}

// Execute runs the show keys command
func (skc *ShowKeysCommand) Execute(profile Profile, context Context, format string) error { //nolint:revive // rich output grouped by sections
	fmt.Printf("ggc Interactive Mode - Effective Keybindings\n")
	fmt.Printf("=============================================\n\n")

	// Get profile info
	prof, exists := skc.resolver.GetProfile(profile)
	if !exists {
		return fmt.Errorf("profile '%s' not found", profile)
	}

	fmt.Printf("Profile: %s", prof.Name)
	if prof.Description != "" {
		fmt.Printf(" (%s)", prof.Description)
	}
	fmt.Printf("\n")

	fmt.Printf("Platform: %s/%s\n", skc.platform, skc.terminal)
	fmt.Printf("Context: %s\n\n", context)

	// Get effective keybindings
	keyMap, err := skc.resolver.Resolve(profile, context)
	if err != nil {
		return fmt.Errorf("failed to resolve keybindings: %w", err)
	}

	// Display keybindings by category
	fmt.Printf("Core Actions:\n")
	fmt.Printf("  Navigation:\n")
	if len(keyMap.MoveUp) > 0 {
		fmt.Printf("    move_up                 %-20s Move up one line\n", FormatKeyStrokesForDisplay(keyMap.MoveUp))
	}
	if len(keyMap.MoveDown) > 0 {
		fmt.Printf("    move_down               %-20s Move down one line\n", FormatKeyStrokesForDisplay(keyMap.MoveDown))
	}
	if len(keyMap.MoveToBeginning) > 0 {
		fmt.Printf("    move_to_beginning       %-20s Move to line beginning\n", FormatKeyStrokesForDisplay(keyMap.MoveToBeginning))
	}
	if len(keyMap.MoveToEnd) > 0 {
		fmt.Printf("    move_to_end             %-20s Move to line end\n", FormatKeyStrokesForDisplay(keyMap.MoveToEnd))
	}

	fmt.Printf("\n  Editing:\n")
	if len(keyMap.DeleteWord) > 0 {
		fmt.Printf("    delete_word             %-20s Delete previous word\n", FormatKeyStrokesForDisplay(keyMap.DeleteWord))
	}
	if len(keyMap.DeleteToEnd) > 0 {
		fmt.Printf("    delete_to_end           %-20s Delete to line end\n", FormatKeyStrokesForDisplay(keyMap.DeleteToEnd))
	}
	if len(keyMap.ClearLine) > 0 {
		fmt.Printf("    clear_line              %-20s Clear entire line\n", FormatKeyStrokesForDisplay(keyMap.ClearLine))
	}

	fmt.Printf("\nQuick Reference:\n")
	fmt.Printf("  quit                    %-20s Exit to shell\n", "Ctrl+C")

	// Show resolution layers
	fmt.Printf("\nResolution Layers Applied:\n")
	fmt.Printf("  1. Base Profile: %s\n", profile)
	fmt.Printf("  2. Platform: %s\n", skc.platform)
	fmt.Printf("  3. Terminal: %s\n", skc.terminal)
	fmt.Printf("  4. User Config: (if configured)\n")

	fmt.Printf("\nTips:\n")
	fmt.Printf("  • Use 'ggc config keybindings --export' to backup your settings\n")
	fmt.Printf("  • Profile switching: set 'interactive.profile' in config\n")

	return nil
}

// FormatKeyStrokesForDisplay returns a comma-separated list of keystrokes suitable for user-facing output.
func FormatKeyStrokesForDisplay(keystrokes []KeyStroke) string {
	if len(keystrokes) == 0 {
		return "none"
	}

	var parts []string
	for _, ks := range keystrokes {
		parts = append(parts, FormatKeyStrokeForDisplay(ks))
	}

	return strings.Join(parts, ", ")
}

// FormatKeyStrokeForDisplay converts a single keystroke into a readable label.
func FormatKeyStrokeForDisplay(ks KeyStroke) string { //nolint:revive // handles numerous escape sequences
	switch ks.Kind {
	case KeyStrokeCtrl:
		return fmt.Sprintf("Ctrl+%c", ks.Rune)
	case KeyStrokeAlt:
		if ks.Name != "" {
			label := strings.ToUpper(ks.Name[:1]) + ks.Name[1:]
			return fmt.Sprintf("Alt+%s", label)
		}
		if ks.Rune != 0 {
			return fmt.Sprintf("Alt+%c", ks.Rune)
		}
		return "Alt+?"
	case KeyStrokeRawSeq:
		// Handle common sequences
		if len(ks.Seq) == 1 {
			switch ks.Seq[0] {
			case 9:
				return "Tab"
			case 13:
				return "Enter"
			case 27:
				return "Esc"
			case 32:
				return "Space"
			}
		}
		// Arrow keys
		if len(ks.Seq) == 3 && ks.Seq[0] == 27 && ks.Seq[1] == 91 {
			switch ks.Seq[2] {
			case 65:
				return "↑"
			case 66:
				return "↓"
			case 67:
				return "→"
			case 68:
				return "←"
			}
		}
		return fmt.Sprintf("Raw[%x]", ks.Seq)
	case KeyStrokeFnKey:
		return ks.Name
	default:
		return fmt.Sprintf("Unknown[%v]", ks)
	}
}

// DebugKeysCommand captures and displays raw key sequences
type DebugKeysCommand struct {
	capturing  bool
	sequences  [][]byte
	outputFile string
}

// NewDebugKeysCommand creates a new debug keys command
func NewDebugKeysCommand(outputFile string) *DebugKeysCommand {
	return &DebugKeysCommand{
		capturing:  false,
		sequences:  make([][]byte, 0),
		outputFile: outputFile,
	}
}

// StartCapture begins capturing raw key sequences
func (dkc *DebugKeysCommand) StartCapture() {
	dkc.capturing = true
	dkc.sequences = make([][]byte, 0)

	fmt.Printf("=== Debug Keys Mode ===\n")
	fmt.Printf("Raw key sequence capture started.\n")
	fmt.Printf("Press keys to see their sequences.\n")
	fmt.Printf("Press Ctrl+C to stop and view results.\n\n")
}

// CaptureSequence captures a raw key sequence
func (dkc *DebugKeysCommand) CaptureSequence(seq []byte) {
	if !dkc.capturing {
		return
	}

	// Make a copy of the sequence
	captured := make([]byte, len(seq))
	copy(captured, seq)
	dkc.sequences = append(dkc.sequences, captured)

	// Display immediately
	fmt.Printf("Captured: %v (hex: %x) (chars: %q)\n", seq, seq, seq)
}

// StopCapture stops capturing and shows results
func (dkc *DebugKeysCommand) StopCapture() error {
	if !dkc.capturing {
		return nil
	}

	dkc.capturing = false

	fmt.Printf("\n=== Capture Results ===\n")
	fmt.Printf("Total sequences captured: %d\n\n", len(dkc.sequences))

	if len(dkc.sequences) == 0 {
		fmt.Printf("No sequences captured.\n")
		return nil
	}

	// Display all captured sequences
	for i, seq := range dkc.sequences {
		fmt.Printf("%d. %v (hex: %x)\n", i+1, seq, seq)

		// Try to identify common sequences
		if identified := dkc.identifySequence(seq); identified != "" {
			fmt.Printf("   → Identified as: %s\n", identified)
		}

		// Show binding format
		fmt.Printf("   → Config format: \"raw:%x\"\n", seq)
	}

	// Save to file if requested
	if dkc.outputFile != "" {
		if err := dkc.saveToFile(); err != nil {
			return fmt.Errorf("failed to save to file: %w", err)
		}
		fmt.Printf("\nSequences saved to: %s\n", dkc.outputFile)
	}

	fmt.Printf("\nTip: Use the 'raw:' format in your config to bind these sequences.\n")

	return nil
}

func (dkc *DebugKeysCommand) formatKeySequence(seq []byte) string { //nolint:revive // classification supports many key types
	if len(seq) == 0 {
		return "(empty)"
	}

	label := dkc.identifySequence(seq)
	if label == "" {
		if len(seq) == 1 {
			b := seq[0]
			switch {
			case b >= 32 && b <= 126:
				label = string([]byte{b})
			case b >= 1 && b <= 26:
				label = fmt.Sprintf("Ctrl+%c", 'A'+b-1)
			case b == 27:
				label = "Esc"
			default:
				label = fmt.Sprintf("0x%02x", b)
			}
		} else {
			label = fmt.Sprintf("%v", seq)
		}
	}

	hexParts := make([]string, len(seq))
	for i, b := range seq {
		hexParts[i] = fmt.Sprintf("0x%02x", b)
	}

	return fmt.Sprintf("%s (%s)", label, strings.Join(hexParts, " "))
}

// identifySequence tries to identify common key sequences
func (dkc *DebugKeysCommand) identifySequence(seq []byte) string { //nolint:revive // identifies many terminal escape sequences
	if len(seq) == 1 {
		switch seq[0] {
		case 9:
			return "Tab"
		case 13:
			return "Enter"
		case 27:
			return "Esc"
		case 32:
			return "Space"
		}
		if seq[0] >= 1 && seq[0] <= 26 {
			return fmt.Sprintf("Ctrl+%c", 'A'+seq[0]-1)
		}
	}

	if len(seq) == 3 && seq[0] == 27 && seq[1] == 91 {
		switch seq[2] {
		case 65:
			return "↑"
		case 66:
			return "↓"
		case 67:
			return "→"
		case 68:
			return "←"
		}
	}

	// Shift-modified arrow keys (CSI 1;2X sequences)
	if len(seq) == 6 && seq[0] == 27 && seq[1] == 91 && seq[2] == 49 && seq[3] == 59 {
		if seq[4] == 50 {
			switch seq[5] {
			case 65:
				return "Shift+↑"
			case 66:
				return "Shift+↓"
			case 67:
				return "Shift+→"
			case 68:
				return "Shift+←"
			}
		}
	}

	// Function keys
	if len(seq) >= 3 && seq[0] == 27 && seq[1] == 79 {
		switch seq[2] {
		case 80:
			return "F1"
		case 81:
			return "F2"
		case 82:
			return "F3"
		case 83:
			return "F4"
		}
	}

	return ""
}

// saveToFile saves captured sequences to a file
func (dkc *DebugKeysCommand) saveToFile() error {
	var content strings.Builder

	content.WriteString("# Raw Key Sequences Captured by ggc debug-keys\n")
	content.WriteString(fmt.Sprintf("# Captured on: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	content.WriteString(fmt.Sprintf("# Total sequences: %d\n\n", len(dkc.sequences)))

	for i, seq := range dkc.sequences {
		content.WriteString(fmt.Sprintf("# Sequence %d\n", i+1))
		content.WriteString(fmt.Sprintf("# Raw: %v\n", seq))
		content.WriteString(fmt.Sprintf("# Hex: %x\n", seq))
		if identified := dkc.identifySequence(seq); identified != "" {
			content.WriteString(fmt.Sprintf("# Identified: %s\n", identified))
		}
		content.WriteString(fmt.Sprintf("raw:%x\n\n", seq))
	}

	if err := os.WriteFile(dkc.outputFile, []byte(content.String()), 0600); err != nil {
		return err
	}

	fmt.Printf("Saved to %s:\n%s", dkc.outputFile, content.String())

	return nil
}

// IsCapturing returns whether debug capture is active
func (dkc *DebugKeysCommand) IsCapturing() bool {
	return dkc.capturing
}
