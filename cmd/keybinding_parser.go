package cmd

import (
	"fmt"
	"strings"
)

// KeyBindingParseResult contains the result of parsing a key binding
type KeyBindingParseResult struct {
	ControlCode byte
	IsValid     bool
}

// parseKeyBindingInternal parses a key binding string and returns the control code and validity
func parseKeyBindingInternal(keyStr string) (KeyBindingParseResult, error) {
	s := strings.TrimSpace(keyStr)
	if s == "" {
		return KeyBindingParseResult{}, fmt.Errorf("empty key binding")
	}

	// Normalize to lowercase for comparison
	sLower := strings.ToLower(s)

	// Try different format parsers
	if result, err := parseCtrlFormat(s, sLower, keyStr); err == nil {
		return result, nil
	}
	if result, err := parseCaretFormat(s, keyStr); err == nil {
		return result, nil
	}
	if result, err := parseEmacsFormat(s, sLower, keyStr); err == nil {
		return result, nil
	}

	return KeyBindingParseResult{}, fmt.Errorf("unsupported key binding format: %s (supported: 'ctrl+w', '^w', 'C-w')", keyStr)
}

// parseCtrlFormat handles "ctrl+<key>" format (case-insensitive)
func parseCtrlFormat(s, sLower, keyStr string) (KeyBindingParseResult, error) {
	if !strings.HasPrefix(sLower, "ctrl+") || len(s) != len("ctrl+")+1 {
		return KeyBindingParseResult{}, fmt.Errorf("not ctrl format")
	}

	c := rune(sLower[len(sLower)-1])
	code := keyBindingCtrlCode(c)
	if code == 0 {
		return KeyBindingParseResult{}, fmt.Errorf("unsupported ctrl key: %s", keyStr)
	}
	return KeyBindingParseResult{ControlCode: code, IsValid: true}, nil
}

// parseCaretFormat handles "^<key>" format (caret notation)
func parseCaretFormat(s, keyStr string) (KeyBindingParseResult, error) {
	if !strings.HasPrefix(s, "^") || len(s) != 2 {
		return KeyBindingParseResult{}, fmt.Errorf("not caret format")
	}

	c := rune(strings.ToLower(s)[1])
	code := keyBindingCtrlCode(c)
	if code == 0 {
		return KeyBindingParseResult{}, fmt.Errorf("unsupported caret key: %s", keyStr)
	}
	return KeyBindingParseResult{ControlCode: code, IsValid: true}, nil
}

// parseEmacsFormat handles "c-<key>" or "C-<key>" format (emacs notation)
func parseEmacsFormat(s, sLower, keyStr string) (KeyBindingParseResult, error) {
	if !strings.HasPrefix(sLower, "c-") || len(s) != 3 {
		return KeyBindingParseResult{}, fmt.Errorf("not emacs format")
	}

	c := rune(sLower[2])
	code := keyBindingCtrlCode(c)
	if code == 0 {
		return KeyBindingParseResult{}, fmt.Errorf("unsupported emacs key: %s", keyStr)
	}
	return KeyBindingParseResult{ControlCode: code, IsValid: true}, nil
}

// validateKeyBindingInternal validates a key binding string without returning the control code
func validateKeyBindingInternal(keyStr string) error {
	_, err := parseKeyBindingInternal(keyStr)
	return err
}

// keyBindingCtrlCode converts a letter to its control byte (e.g., 'a' => 1).
// Handles both uppercase and lowercase letters for compatibility.
func keyBindingCtrlCode(r rune) byte {
	// Handle lowercase letters a-z
	if r >= 'a' && r <= 'z' {
		return byte(r-'a') + 1
	}
	// Handle uppercase letters A-Z
	if r >= 'A' && r <= 'Z' {
		return byte(r-'A') + 1
	}
	return 0
}
