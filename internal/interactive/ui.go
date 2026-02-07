// Package interactive houses interactive UI types and helpers shared across the application.
package interactive

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"

	"github.com/bmf-san/ggc/v7/internal/config"
	kb "github.com/bmf-san/ggc/v7/internal/keybindings"
	"github.com/bmf-san/ggc/v7/internal/termio"
	"github.com/bmf-san/ggc/v7/pkg/git"
)

// initialInputCapacity defines the initial capacity for the input rune buffer
// used by the real-time editor. It helps minimize reallocations during typing
// while keeping memory usage modest.
const initialInputCapacity = 64

// UI represents the interface for terminal UI operations
type UI struct {
	stdin           io.Reader
	stdout          io.Writer
	stderr          io.Writer
	term            termio.Terminal
	renderer        *Renderer
	state           *UIState
	handler         *KeyHandler
	colors          *ANSIColors
	gitStatus       *GitStatus
	gitClient       git.StatusInfoReader
	reader          *bufio.Reader
	contextMgr      *kb.ContextManager
	profile         kb.Profile
	workflowMgr     *WorkflowManager
	workflow        *Workflow
	workflowEx      *WorkflowExecutor
	softCancelFlash atomic.Bool
	workflowError   string
	errorExpiresAt  time.Time
	workflowNotice  string
	noticeExpiresAt time.Time
}

// NewUI creates a new UI with the provided git client and loads keybindings from config
func NewUI(gitClient git.StatusInfoReader, router ...CommandRouter) *UI {
	colors := NewANSIColors()

	renderer := &Renderer{
		writer: os.Stdout,
		colors: colors,
	}
	renderer.updateSize()

	state := &UIState{
		selected:       0,
		input:          "",
		filtered:       []CommandInfo{},
		context:        kb.ContextGlobal, // Start in global context
		contextStack:   []kb.Context{},
		mode:           ModeSearch,
		workflowFocus:  FocusInput,
		workflowOffset: 0,
	}

	// Load config and create resolver
	var cfg *config.Config
	if ops, ok := gitClient.(git.ConfigOps); ok {
		configManager := config.NewConfigManager(ops)
		// Load config - if it fails, we'll use defaults from manager
		_ = configManager.Load()
		cfg = configManager.GetConfig()
	} else {
		// Fallback to empty config (built-in defaults and profiles will be used)
		cfg = &config.Config{}
	}

	// Create KeyBinding resolver and register built-in profiles
	resolver := kb.NewKeyBindingResolver(cfg)
	kb.RegisterBuiltinProfiles(resolver)
	contextManager := kb.NewContextManager(resolver)

	// Determine which profile to use (default to "default" profile)
	profile := kb.ProfileDefault
	if cfg.Interactive.Profile != "" {
		switch kb.Profile(cfg.Interactive.Profile) {
		case kb.ProfileEmacs, kb.ProfileVi, kb.ProfileReadline:
			profile = kb.Profile(cfg.Interactive.Profile)
		default:
			fmt.Fprintf(os.Stderr, "Warning: Unknown profile '%s', using default\n", cfg.Interactive.Profile)
		}
	}

	// Resolve contextual keybindings for all contexts
	contextualMap, err := resolver.ResolveContextual(profile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to resolve keybindings: %v. Using defaults.\n", err)
		// Fallback to legacy defaults
		keyMap := kb.DefaultKeyBindingMap()
		contextualMap = &kb.ContextualKeyBindingMap{
			Profile:  kb.ProfileDefault,
			Platform: kb.DetectPlatform(),
			Terminal: kb.DetectTerminal(),
			Contexts: map[kb.Context]*kb.KeyBindingMap{
				kb.ContextGlobal:  keyMap,
				kb.ContextInput:   keyMap,
				kb.ContextResults: keyMap,
				kb.ContextSearch:  keyMap,
			},
		}
	}

	ui := &UI{
		stdin:       os.Stdin,
		stdout:      os.Stdout,
		stderr:      os.Stderr,
		term:        termio.DefaultTerminal{},
		renderer:    renderer,
		state:       state,
		colors:      colors,
		gitClient:   gitClient,
		gitStatus:   getGitStatus(gitClient),
		contextMgr:  contextManager,
		profile:     profile,
		workflowMgr: NewWorkflowManager(),
	}
	ui.updateWorkflowPointer()
	state.onContextChange = func(_ kb.Context, newCtx kb.Context) {
		contextManager.SetContext(newCtx)
	}

	ui.handler = &KeyHandler{
		ui:            ui,
		contextualMap: contextualMap,
	}

	// Set up workflow executor if router is provided
	if len(router) > 0 && router[0] != nil {
		ui.workflowEx = NewWorkflowExecutor(router[0], ui)
	}

	return ui
}

// Run executes the incremental search interactive UI with the provided custom git client,
// and returns the selected command as []string (or nil if nothing is selected).
func Run(gitClient git.StatusInfoReader) []string {
	ui := NewUI(gitClient)
	return ui.Run()
}
