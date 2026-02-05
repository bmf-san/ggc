// Package interactive houses interactive UI types and helpers shared across the application.
package interactive

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"

	"golang.org/x/term"

	commandregistry "github.com/bmf-san/ggc/v7/cmd/command"
	"github.com/bmf-san/ggc/v7/internal/config"
	kb "github.com/bmf-san/ggc/v7/internal/keybindings"
	"github.com/bmf-san/ggc/v7/internal/termio"
	"github.com/bmf-san/ggc/v7/pkg/git"
)

// initialInputCapacity defines the initial capacity for the input rune buffer
// used by the real-time editor. It helps minimize reallocations during typing
// while keeping memory usage modest.
const initialInputCapacity = 64

// CommandInfo contains the name and description of the command
type CommandInfo struct {
	Command     string
	Description string
}

func buildInteractiveCommands() []CommandInfo {
	var list []CommandInfo
	registry := commandregistry.NewRegistry()
	allCommands := registry.All()
	for i := range allCommands {
		cmd := &allCommands[i]
		if cmd.Hidden {
			continue
		}
		if len(cmd.Subcommands) == 0 {
			list = append(list, CommandInfo{Command: cmd.Name, Description: cmd.Summary})
			continue
		}
		for _, sub := range cmd.Subcommands {
			if sub.Hidden {
				continue
			}
			list = append(list, CommandInfo{Command: sub.Name, Description: sub.Summary})
		}
	}
	return list
}

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

var commands = buildInteractiveCommands()

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

// ToggleWorkflowView toggles between search mode and workflow mode.
func (ui *UI) ToggleWorkflowView() {
	if ui == nil || ui.state == nil {
		return
	}
	if ui.state.IsWorkflowMode() {
		ui.enterSearchMode()
		return
	}
	ui.enterWorkflowMode()
}

// enterWorkflowMode switches UI into workflow management mode.
// Workflow mode has no input field - it's a pure management interface.
func (ui *UI) enterWorkflowMode() {
	if ui == nil || ui.state == nil {
		return
	}
	ui.state.SetMode(ModeWorkflow)
	ui.state.SetContext(kb.ContextGlobal)
	ui.ensureWorkflowListSelection()
	ui.updateWorkflowPointer()
}

// enterSearchMode switches UI back to search mode.
func (ui *UI) enterSearchMode() {
	if ui == nil || ui.state == nil {
		return
	}
	ui.state.SetMode(ModeSearch)
	ui.state.FocusInput()
	ui.state.SetContext(kb.ContextGlobal)
}

// AddToWorkflow adds a command to the active workflow.
func (ui *UI) AddToWorkflow(command string, args []string, description string) int {
	if ui.workflowMgr == nil {
		return 0
	}
	activeID := ui.workflowMgr.GetActiveID()
	if id, ok := ui.workflowMgr.AddStep(activeID, command, args, description); ok {
		ui.updateWorkflowPointer()
		return id
	}
	return 0
}

// ApplyContextualKeybindings updates the active keybinding map, satisfying keybindings.ContextualMapApplier.
func (ui *UI) ApplyContextualKeybindings(contextual *kb.ContextualKeyBindingMap) {
	if ui == nil || ui.handler == nil || contextual == nil {
		return
	}
	ui.handler.contextualMap = contextual
}

func (ui *UI) resetToSearchMode() bool {
	if ui == nil || ui.state == nil {
		return false
	}

	state := ui.state
	active := state.HasInput() || state.IsWorkflowMode() || len(state.contextStack) > 0 || state.GetCurrentContext() != kb.ContextGlobal
	state.ClearInput()
	state.selected = 0
	state.contextStack = nil
	state.SetContext(kb.ContextGlobal)
	state.SetMode(ModeSearch)
	state.FocusInput()
	return active
}

// ResetToSearchMode clears the interactive search UI back to its default state.
func (ui *UI) ResetToSearchMode() bool {
	return ui.resetToSearchMode()
}

func (ui *UI) readPlaceholderInput() (string, bool) {
	if ui == nil || ui.handler == nil {
		return "", true
	}
	return ui.handler.getRealTimeInput()
}

// ClearWorkflow removes all steps from the workflow
func (ui *UI) ClearWorkflow() {
	if ui.workflowMgr == nil {
		return
	}
	ui.workflowMgr.ClearWorkflow(ui.workflowMgr.GetActiveID())
	ui.updateWorkflowPointer()
}

// ExecuteWorkflow executes the current workflow
func (ui *UI) ExecuteWorkflow() error {
	if ui.workflowEx == nil {
		return fmt.Errorf("workflow executor not initialized")
	}

	if ui.workflow == nil || ui.workflow.IsEmpty() {
		return fmt.Errorf("workflow is empty")
	}

	return ui.workflowEx.Execute(ui.workflow)
}

func (ui *UI) updateWorkflowPointer() {
	if ui == nil || ui.workflowMgr == nil {
		return
	}
	wf, ok := ui.workflowMgr.GetWorkflow(ui.workflowMgr.GetActiveID())
	if ok {
		ui.workflow = wf
		return
	}
	ui.workflow = nil
}

func (ui *UI) listWorkflows() []WorkflowSummary {
	if ui.workflowMgr == nil {
		return nil
	}
	return ui.workflowMgr.ListWorkflows()
}

func (ui *UI) ensureWorkflowListSelection() {
	if ui == nil || ui.state == nil {
		return
	}
	summaries := ui.listWorkflows()
	activeID := 0
	if ui.workflowMgr != nil {
		activeID = ui.workflowMgr.GetActiveID()
	}
	if activeID != 0 {
		for i, summary := range summaries {
			if summary.ID == activeID {
				ui.state.workflowListIdx = i
				break
			}
		}
	}
	ui.state.SetWorkflowListIndex(ui.state.workflowListIdx, len(summaries))
}

// Run executes the incremental search interactive UI with the provided custom git client,
// and returns the selected command as []string (or nil if nothing is selected).
func Run(gitClient git.StatusInfoReader) []string {
	ui := NewUI(gitClient)
	return ui.Run()
}

// writeError writes an error message to stderr
func (ui *UI) writeError(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(ui.stderr, format+"\n", a...)
}

// write writes a message to stdout
func (ui *UI) write(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(ui.stdout, format, a...)
}

// writeColor writes a colored message to stdout
func (ui *UI) writeColor(text string) {
	_, _ = fmt.Fprint(ui.stdout, text)
}

// writeln writes a message with newline to stdout
func (ui *UI) writeln(format string, a ...interface{}) {
	// Move to line start, clear line, write content, then CRLF
	_, _ = fmt.Fprint(ui.stdout, "\r\x1b[K")
	_, _ = fmt.Fprintf(ui.stdout, format+"\r\n", a...)
}

func (ui *UI) notifySoftCancel() {
	ui.softCancelFlash.Store(true)
}

func (ui *UI) consumeSoftCancelFlash() bool {
	return ui.softCancelFlash.Swap(false)
}

func (ui *UI) notifyWorkflowError(message string, duration time.Duration) {
	if ui == nil {
		return
	}
	ui.workflowNotice = ""
	ui.workflowError = message
	ui.errorExpiresAt = time.Now().Add(duration)
}

func (ui *UI) workflowErrorMessage() string {
	if ui == nil || ui.workflowError == "" {
		return ""
	}
	if time.Now().After(ui.errorExpiresAt) {
		ui.workflowError = ""
		return ""
	}
	return ui.workflowError
}

func (ui *UI) notifyWorkflowSuccess(message string, duration time.Duration) {
	if ui == nil {
		return
	}
	ui.workflowError = ""
	ui.workflowNotice = message
	ui.noticeExpiresAt = time.Now().Add(duration)
}

func (ui *UI) workflowNoticeMessage() string {
	if ui == nil || ui.workflowNotice == "" {
		return ""
	}
	if time.Now().After(ui.noticeExpiresAt) {
		ui.workflowNotice = ""
		return ""
	}
	return ui.workflowNotice
}

// setupTerminal configures terminal raw mode and returns the old state and error status
func (ui *UI) setupTerminal() (*term.State, bool) {
	var oldState *term.State
	if f, ok := ui.stdin.(*os.File); ok {
		// Check if it's a real terminal (TTY) - skip for non-TTY or when term is nil
		if ui.term == nil || !term.IsTerminal(int(f.Fd())) {
			// Not a real terminal, skip raw mode setup for debugging
			return nil, true
		}

		fd := int(f.Fd())
		var err error
		oldState, err = ui.term.MakeRaw(fd)
		if err != nil {
			ui.writeError("Failed to set terminal to raw mode: %v", err)
			return nil, false
		}
	}
	return oldState, true
}

// Run executes the interactive UI
func (ui *UI) Run() []string {
	oldState, reader, isRawMode := ui.initializeTerminal()
	if oldState == nil && isRawMode {
		return nil
	}

	// Set up terminal restoration for raw mode
	if f, ok := ui.stdin.(*os.File); ok && isRawMode {
		fd := int(f.Fd())
		defer func() {
			if err := ui.term.Restore(fd, oldState); err != nil {
				ui.writeError("failed to restore terminal state: %v", err)
			}
		}()
	}

	return ui.runMainLoop(reader, isRawMode, oldState)
}

// initializeTerminal sets up the terminal and returns the old state, reader, and raw mode status
func (ui *UI) initializeTerminal() (*term.State, *bufio.Reader, bool) {
	oldState, ok := ui.setupTerminal()
	if !ok {
		return nil, nil, false
	}

	// Check if we're in raw mode (real terminal) or not
	isRawMode := oldState != nil

	// Set up reader based on mode
	var reader *bufio.Reader
	if !isRawMode {
		reader = bufio.NewReader(ui.stdin)
	}

	return oldState, reader, isRawMode
}

// runMainLoop handles the main input loop
func (ui *UI) runMainLoop(reader *bufio.Reader, isRawMode bool, oldState *term.State) []string {
	if isRawMode {
		if ui.reader == nil {
			ui.reader = bufio.NewReader(ui.stdin)
		}
	} else {
		ui.reader = reader
	}

	for {
		ui.state.UpdateFiltered()
		ui.renderer.Render(ui, ui.state)

		r, err := ui.readNextRune(reader, isRawMode)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			continue // Skip this iteration for other errors
		}

		// Handle key input with rune
		isSingleByte := isRawMode // In raw mode, we read single bytes; in buffered mode, we read full runes
		shouldContinue, result := ui.handler.HandleKey(r, isSingleByte, oldState, reader)
		if !shouldContinue {
			return result
		}
	}
}

// readNextRune reads the next rune from input based on the mode
func (ui *UI) readNextRune(reader *bufio.Reader, isRawMode bool) (rune, error) {
	if isRawMode {
		for {
			// Read single byte directly from stdin in raw mode
			var buf [1]byte
			n, readErr := ui.stdin.Read(buf[:])

			if readErr != nil {
				return 0, readErr
			}

			if n == 0 {
				continue // Try again if no bytes read
			}

			return rune(buf[0]), nil
		}
	}

	// Use buffered reader for non-TTY
	r, _, err := reader.ReadRune()
	return r, err
}

// Extract <...> placeholders from a string
func extractPlaceholders(s string) []string {
	var res []string
	start := -1
	for i, c := range s {
		if c == '<' {
			start = i + 1
		} else if c == '>' && start != -1 {
			res = append(res, s[start:i])
			start = -1
		}
	}
	return res
}
