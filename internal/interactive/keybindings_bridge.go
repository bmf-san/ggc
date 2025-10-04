package interactive

import kb "github.com/bmf-san/ggc/v7/internal/keybindings"

type (
	Profile                 = kb.Profile
	Context                 = kb.Context
	KeyStroke               = kb.KeyStroke
	KeyBindingMap           = kb.KeyBindingMap
	ContextualKeyBindingMap = kb.ContextualKeyBindingMap
	KeyBindingResolver      = kb.KeyBindingResolver
	ContextManager          = kb.ContextManager
)

var (
	NewKeyBindingResolver      = kb.NewKeyBindingResolver
	RegisterBuiltinProfiles    = kb.RegisterBuiltinProfiles
	DefaultKeyBindingMap       = kb.DefaultKeyBindingMap
	DetectPlatform             = kb.DetectPlatform
	DetectTerminal             = kb.DetectTerminal
	NewContextManager          = kb.NewContextManager
	NewCtrlKeyStroke           = kb.NewCtrlKeyStroke
	NewCharKeyStroke           = kb.NewCharKeyStroke
	NewRawKeyStroke            = kb.NewRawKeyStroke
	NewEscapeKeyStroke         = kb.NewEscapeKeyStroke
	NewAltKeyStroke            = kb.NewAltKeyStroke
	NewContextualKeyBindingMap = kb.NewContextualKeyBindingMap
	FormatKeyStrokesForDisplay = kb.FormatKeyStrokesForDisplay
)

const (
	ProfileDefault  = kb.ProfileDefault
	ProfileEmacs    = kb.ProfileEmacs
	ProfileVi       = kb.ProfileVi
	ProfileReadline = kb.ProfileReadline

	ContextGlobal  = kb.ContextGlobal
	ContextInput   = kb.ContextInput
	ContextResults = kb.ContextResults
	ContextSearch  = kb.ContextSearch
)
