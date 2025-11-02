package interactive

import kb "github.com/bmf-san/ggc/v7/internal/keybindings"

type (
	// Profile exposes the kb.Profile type to interactive consumers.
	Profile = kb.Profile
	// Context exposes the kb.Context type to interactive consumers.
	Context = kb.Context
	// KeyStroke exposes the kb.KeyStroke type to interactive consumers.
	KeyStroke = kb.KeyStroke
	// KeyBindingMap exposes the kb.KeyBindingMap type to interactive consumers.
	KeyBindingMap = kb.KeyBindingMap
	// ContextualKeyBindingMap exposes the kb.ContextualKeyBindingMap type to interactive consumers.
	ContextualKeyBindingMap = kb.ContextualKeyBindingMap
	// KeyBindingResolver exposes the kb.KeyBindingResolver type to interactive consumers.
	KeyBindingResolver = kb.KeyBindingResolver
	// ContextManager exposes the kb.ContextManager type to interactive consumers.
	ContextManager = kb.ContextManager
)

var (
	// NewKeyBindingResolver constructs a key binding resolver using the keybindings package implementation.
	NewKeyBindingResolver = kb.NewKeyBindingResolver
	// RegisterBuiltinProfiles registers builtin profiles in the underlying keybindings package.
	RegisterBuiltinProfiles = kb.RegisterBuiltinProfiles
	// DefaultKeyBindingMap provides the default key binding map from the keybindings package.
	DefaultKeyBindingMap = kb.DefaultKeyBindingMap
	// DetectPlatform returns the inferred platform from the keybindings package.
	DetectPlatform = kb.DetectPlatform
	// DetectTerminal returns the inferred terminal from the keybindings package.
	DetectTerminal = kb.DetectTerminal
	// NewContextManager constructs a context manager using the keybindings package.
	NewContextManager = kb.NewContextManager
	// NewCtrlKeyStroke creates a control key stroke using the keybindings package implementation.
	NewCtrlKeyStroke = kb.NewCtrlKeyStroke
	// NewCharKeyStroke creates a character key stroke using the keybindings package implementation.
	NewCharKeyStroke = kb.NewCharKeyStroke
	// NewRawKeyStroke creates a raw key stroke using the keybindings package implementation.
	NewRawKeyStroke = kb.NewRawKeyStroke
	// NewEscapeKeyStroke creates an escape key stroke using the keybindings package implementation.
	NewEscapeKeyStroke = kb.NewEscapeKeyStroke
	// NewAltKeyStroke creates an alt key stroke using the keybindings package implementation.
	NewAltKeyStroke = kb.NewAltKeyStroke
	// NewContextualKeyBindingMap builds a contextual key binding map using the keybindings package implementation.
	NewContextualKeyBindingMap = kb.NewContextualKeyBindingMap
	// FormatKeyStrokesForDisplay formats key strokes for human-readable display.
	FormatKeyStrokesForDisplay = kb.FormatKeyStrokesForDisplay
)

const (
	// ProfileDefault exposes the default key binding profile identifier.
	ProfileDefault = kb.ProfileDefault
	// ProfileEmacs exposes the emacs key binding profile identifier.
	ProfileEmacs = kb.ProfileEmacs
	// ProfileVi exposes the vi key binding profile identifier.
	ProfileVi = kb.ProfileVi
	// ProfileReadline exposes the readline key binding profile identifier.
	ProfileReadline = kb.ProfileReadline

	// ContextGlobal exposes the global key binding context identifier.
	ContextGlobal = kb.ContextGlobal
	// ContextInput exposes the input key binding context identifier.
	ContextInput = kb.ContextInput
	// ContextResults exposes the results key binding context identifier.
	ContextResults = kb.ContextResults
	// ContextSearch exposes the search key binding context identifier.
	ContextSearch = kb.ContextSearch
	// ContextWorkflowView exposes the workflow management context identifier.
	ContextWorkflowView = kb.ContextWorkflowView
	// ContextWorkflowSelection exposes the workflow selection context identifier.
	ContextWorkflowSelection = kb.ContextWorkflowSelection
)
