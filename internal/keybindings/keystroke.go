package keybindings

import (
	"fmt"
	"strings"
)

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

// NewUpArrowKeyStroke creates a new Up Arrow KeyStroke (CSI A)
// Can be used to rebind up arrow for list navigation
func NewUpArrowKeyStroke() KeyStroke {
	return NewRawKeyStroke([]byte{27, '[', 'A'}) // ESC [ A
}

// NewDownArrowKeyStroke creates a new Down Arrow KeyStroke (CSI B)
// Can be used to rebind down arrow for list navigation
func NewDownArrowKeyStroke() KeyStroke {
	return NewRawKeyStroke([]byte{27, '[', 'B'}) // ESC [ B
}

// NewLeftArrowKeyStroke creates a new Left Arrow KeyStroke (CSI D)
// Can be used to rebind left arrow for cursor movement
func NewLeftArrowKeyStroke() KeyStroke {
	return NewRawKeyStroke([]byte{27, '[', 'D'}) // ESC [ D
}

// NewRightArrowKeyStroke creates a new Right Arrow KeyStroke (CSI C)
// Can be used to rebind right arrow for cursor movement
func NewRightArrowKeyStroke() KeyStroke {
	return NewRawKeyStroke([]byte{27, '[', 'C'}) // ESC [ C
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

	// Handle arrow keys - all four directions are now rebindable
	switch sLower {
	case "up", "arrow-up", "arrowup":
		return NewUpArrowKeyStroke(), nil
	case "down", "arrow-down", "arrowdown":
		return NewDownArrowKeyStroke(), nil
	case "left", "arrow-left", "arrowleft":
		return NewLeftArrowKeyStroke(), nil
	case "right", "arrow-right", "arrowright":
		return NewRightArrowKeyStroke(), nil
	}

	return KeyStroke{}, fmt.Errorf("unsupported key binding format: %s (supported: 'ctrl+w', '^w', 'C-w', 'alt+backspace', 'M-backspace', 'up', 'down', 'left', 'right')", keyStr)
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
