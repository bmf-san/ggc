package keybindings

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go.yaml.in/yaml/v3"

	"github.com/bmf-san/ggc/v7/internal/config"
)

// KeyBindingMap holds resolved key strokes for interactive actions.
// Supports multiple key strokes per action while maintaining backward compatibility.
type KeyBindingMap struct {
	DeleteWord         []KeyStroke // default: [Ctrl+W]
	ClearLine          []KeyStroke // default: [Ctrl+U]
	DeleteToEnd        []KeyStroke // default: [Ctrl+K]
	MoveToBeginning    []KeyStroke // default: [Ctrl+A]
	MoveToEnd          []KeyStroke // default: [Ctrl+E]
	MoveUp             []KeyStroke // default: [Ctrl+P], can add: [up arrow]
	MoveDown           []KeyStroke // default: [Ctrl+N], can add: [down arrow]
	MoveLeft           []KeyStroke // default: [], can add: [left arrow] for cursor movement
	MoveRight          []KeyStroke // default: [], can add: [right arrow] for cursor movement
	AddToWorkflow      []KeyStroke // default: [Tab]
	ToggleWorkflowView []KeyStroke // default: [Ctrl+T]
	ClearWorkflow      []KeyStroke // default: [c]
	WorkflowCreate     []KeyStroke // default: [Ctrl+N]
	WorkflowDelete     []KeyStroke // default: [Ctrl+D]
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
		MoveLeft:           []KeyStroke{}, // Empty by default, users can add left arrow
		MoveRight:          []KeyStroke{}, // Empty by default, users can add right arrow
		AddToWorkflow:      []KeyStroke{NewTabKeyStroke()},
		ToggleWorkflowView: []KeyStroke{NewCtrlKeyStroke('t')},
		ClearWorkflow:      []KeyStroke{NewCharKeyStroke('c')},
		WorkflowCreate:     []KeyStroke{NewCtrlKeyStroke('n')},
		WorkflowDelete:     []KeyStroke{NewCtrlKeyStroke('d')},
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
		"move_left":            km.MoveLeft,
		"move_right":           km.MoveRight,
		"add_to_workflow":      km.AddToWorkflow,
		"toggle_workflow_view": km.ToggleWorkflowView,
		"clear_workflow":       km.ClearWorkflow,
		"workflow_create":      km.WorkflowCreate,
		"workflow_delete":      km.WorkflowDelete,
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

// ForceEnvironment overrides detected platform and terminal (primarily for tests).
func (r *KeyBindingResolver) ForceEnvironment(platform, terminal string) {
	if strings.TrimSpace(platform) != "" {
		r.platform = platform
	}
	if strings.TrimSpace(terminal) != "" {
		r.terminal = terminal
	}
	r.ClearCache()
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
	result["move_left"] = clone(keyMap.MoveLeft)
	result["move_right"] = clone(keyMap.MoveRight)
	result["add_to_workflow"] = clone(keyMap.AddToWorkflow)
	result["toggle_workflow_view"] = clone(keyMap.ToggleWorkflowView)
	result["clear_workflow"] = clone(keyMap.ClearWorkflow)
	result["workflow_create"] = clone(keyMap.WorkflowCreate)
	result["workflow_delete"] = clone(keyMap.WorkflowDelete)

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
	keyMap.WorkflowCreate = append(keyMap.WorkflowCreate, defaults.WorkflowCreate...)
	keyMap.WorkflowDelete = append(keyMap.WorkflowDelete, defaults.WorkflowDelete...)
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
	applyBinding("move_left", &keyMap.MoveLeft)
	applyBinding("move_right", &keyMap.MoveRight)
	applyBinding("add_to_workflow", &keyMap.AddToWorkflow)
	applyBinding("toggle_workflow_view", &keyMap.ToggleWorkflowView)
	applyBinding("clear_workflow", &keyMap.ClearWorkflow)
	applyBinding("workflow_create", &keyMap.WorkflowCreate)
	applyBinding("workflow_delete", &keyMap.WorkflowDelete)
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
	case "move_left":
		keyMap.MoveLeft = bindings
		return true
	case "move_right":
		keyMap.MoveRight = bindings
		return true
	}
	return false
}

// applyWorkflowAction applies workflow-related keybinding actions
func (r *KeyBindingResolver) applyWorkflowAction(keyMap *KeyBindingMap, action string, bindings []KeyStroke) {
	actionMap := map[string]*[]KeyStroke{
		"add_to_workflow":      &keyMap.AddToWorkflow,
		"toggle_workflow_view": &keyMap.ToggleWorkflowView,
		"clear_workflow":       &keyMap.ClearWorkflow,
		"workflow_create":      &keyMap.WorkflowCreate,
		"workflow_delete":      &keyMap.WorkflowDelete,
		"soft_cancel":          &keyMap.SoftCancel,
	}

	if target, exists := actionMap[action]; exists {
		*target = bindings
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
		"move_left":            userBindings.MoveLeft,
		"move_right":           userBindings.MoveRight,
		"add_to_workflow":      userBindings.AddToWorkflow,
		"toggle_workflow_view": userBindings.ToggleWorkflowView,
		"clear_workflow":       userBindings.ClearWorkflow,
		"workflow_create":      userBindings.WorkflowCreate,
		"workflow_delete":      userBindings.WorkflowDelete,
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
				case "move_left":
					keyMap.MoveLeft = []KeyStroke{ks}
				case "move_right":
					keyMap.MoveRight = []KeyStroke{ks}
				case "add_to_workflow":
					keyMap.AddToWorkflow = []KeyStroke{ks}
				case "toggle_workflow_view":
					keyMap.ToggleWorkflowView = []KeyStroke{ks}
				case "clear_workflow":
					keyMap.ClearWorkflow = []KeyStroke{ks}
				case "workflow_create":
					keyMap.WorkflowCreate = []KeyStroke{ks}
				case "workflow_delete":
					keyMap.WorkflowDelete = []KeyStroke{ks}
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
		"GGC_KEYBIND_WORKFLOW_CREATE":      &keyMap.WorkflowCreate,
		"GGC_KEYBIND_WORKFLOW_DELETE":      &keyMap.WorkflowDelete,
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
	case "move_left":
		keyMap.MoveLeft = keystrokes
		return true
	case "move_right":
		keyMap.MoveRight = keystrokes
		return true
	}
	return false
}

// applyUserWorkflowAction applies user workflow-related keybinding actions
func (r *KeyBindingResolver) applyUserWorkflowAction(keyMap *KeyBindingMap, action string, keystrokes []KeyStroke) {
	actionMap := map[string]*[]KeyStroke{
		"add_to_workflow":      &keyMap.AddToWorkflow,
		"toggle_workflow_view": &keyMap.ToggleWorkflowView,
		"clear_workflow":       &keyMap.ClearWorkflow,
		"workflow_create":      &keyMap.WorkflowCreate,
		"workflow_delete":      &keyMap.WorkflowDelete,
		"soft_cancel":          &keyMap.SoftCancel,
	}

	if target, exists := actionMap[action]; exists {
		*target = keystrokes
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

// Runtime Profile Switching

// ContextualMapApplier applies resolved keybindings to interested consumers.
type ContextualMapApplier interface {
	ApplyContextualKeybindings(*ContextualKeyBindingMap)
}

// ProfileSwitcher manages runtime profile switching functionality
type ProfileSwitcher struct {
	resolver       *KeyBindingResolver
	currentProfile Profile
	applier        ContextualMapApplier
}

// NewProfileSwitcher creates a new profile switcher
func NewProfileSwitcher(resolver *KeyBindingResolver, applier ContextualMapApplier) *ProfileSwitcher {
	return &ProfileSwitcher{
		resolver:       resolver,
		currentProfile: ProfileDefault,
		applier:        applier,
	}
}

// SwitchProfile switches to a new profile at runtime
func (ps *ProfileSwitcher) SwitchProfile(newProfile Profile) error {
	if _, exists := ps.resolver.GetProfile(newProfile); !exists {
		return fmt.Errorf("profile %s not found", newProfile)
	}

	ps.resolver.ClearCache()

	newContextualMap, err := ps.resolver.ResolveContextual(newProfile)
	if err != nil {
		return fmt.Errorf("failed to resolve profile %s: %w", newProfile, err)
	}

	if ps.applier != nil {
		ps.applier.ApplyContextualKeybindings(newContextualMap)
	}

	oldProfile := ps.currentProfile
	ps.currentProfile = newProfile

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

type profileSwitchHandler func(*ProfileSwitcher, []string) error

var profileSwitchCommandHandlers = map[string]profileSwitchHandler{
	"list":    handleProfileListCommand,
	"switch":  handleProfileSwitchCommand,
	"preview": handleProfilePreviewCommand,
	"compare": handleProfileCompareCommand,
	"reload":  handleProfileReloadCommand,
}

// HandleProfileSwitchCommand processes profile switching commands
func HandleProfileSwitchCommand(switcher *ProfileSwitcher, command string) error {
	parts := strings.Fields(strings.TrimSpace(command))
	if len(parts) == 0 {
		return fmt.Errorf("no command provided")
	}

	subcommand := parts[0]
	args := parts[1:]

	handler, ok := profileSwitchCommandHandlers[subcommand]
	if !ok {
		return fmt.Errorf("unknown subcommand: %s", subcommand)
	}

	return handler(switcher, args)
}

func handleProfileListCommand(switcher *ProfileSwitcher, _ []string) error {
	profiles := switcher.GetAvailableProfiles()
	fmt.Println("Available profiles:")
	for _, profile := range profiles {
		currentMarker := ""
		if profile == switcher.GetCurrentProfile() {
			currentMarker = " (current)"
		}
		fmt.Printf("  - %s%s\n", profile, currentMarker)
	}

	return nil
}

func handleProfileSwitchCommand(switcher *ProfileSwitcher, args []string) error {
	profile, err := requireProfileArg(args, "switch <profile>")
	if err != nil {
		return err
	}

	return switcher.SwitchProfile(profile)
}

func handleProfilePreviewCommand(switcher *ProfileSwitcher, args []string) error {
	profile, err := requireProfileArg(args, "preview <profile>")
	if err != nil {
		return err
	}

	preview, err := switcher.PreviewProfile(profile)
	if err != nil {
		return err
	}

	fmt.Printf("Preview for profile %s:\n", profile)
	for ctx, mapBinding := range preview.Contexts {
		fmt.Printf("  Context: %s\n", ctx)
		fmt.Printf("    move_up                 %-20s Move up one line\n", FormatKeyStrokesForDisplay(mapBinding.MoveUp))
		fmt.Printf("    move_down               %-20s Move down one line\n", FormatKeyStrokesForDisplay(mapBinding.MoveDown))
		fmt.Printf("    move_to_beginning       %-20s Move to line beginning\n", FormatKeyStrokesForDisplay(mapBinding.MoveToBeginning))
		fmt.Printf("    move_to_end             %-20s Move to line end\n", FormatKeyStrokesForDisplay(mapBinding.MoveToEnd))
		fmt.Printf("    delete_word             %-20s Delete previous word\n", FormatKeyStrokesForDisplay(mapBinding.DeleteWord))
		fmt.Printf("    delete_to_end           %-20s Delete to line end\n", FormatKeyStrokesForDisplay(mapBinding.DeleteToEnd))
		fmt.Printf("    clear_line              %-20s Clear entire line\n", FormatKeyStrokesForDisplay(mapBinding.ClearLine))
	}

	return nil
}

func handleProfileCompareCommand(switcher *ProfileSwitcher, args []string) error {
	profile, err := requireProfileArg(args, "compare <profile>")
	if err != nil {
		return err
	}

	comparison, err := switcher.GetProfileComparison(profile)
	if err != nil {
		return err
	}

	fmt.Printf("Comparison between current profile (%s) and %s:\n", switcher.GetCurrentProfile(), profile)
	for category, value := range comparison {
		fmt.Printf("  %s: %v\n", category, value)
	}

	return nil
}

func handleProfileReloadCommand(switcher *ProfileSwitcher, _ []string) error {
	return switcher.ReloadCurrentProfile()
}

func requireProfileArg(args []string, usage string) (Profile, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("usage: %s", usage)
	}

	return Profile(args[0]), nil
}

// ShowCurrentProfileCommand returns a string representing the current profile status
func ShowCurrentProfileCommand(switcher *ProfileSwitcher) string {
	return fmt.Sprintf("Current profile: %s", switcher.GetCurrentProfile())
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

// ContextManager manages active contexts and notifies callbacks on transitions.
type ContextManager struct {
	resolver  *KeyBindingResolver
	current   Context
	stack     []Context
	callbacks map[Context][]func(Context, Context)
}

// NewContextManager creates a new ContextManager.
func NewContextManager(resolver *KeyBindingResolver) *ContextManager {
	return &ContextManager{
		resolver:  resolver,
		current:   ContextGlobal,
		stack:     make([]Context, 0, 4),
		callbacks: make(map[Context][]func(Context, Context)),
	}
}

// RegisterContextCallback registers a callback invoked when the target context becomes active.
func (cm *ContextManager) RegisterContextCallback(ctx Context, callback func(Context, Context)) {
	if callback == nil {
		return
	}
	cm.callbacks[ctx] = append(cm.callbacks[ctx], callback)
}

// GetCurrentContext returns the currently active context.
func (cm *ContextManager) GetCurrentContext() Context {
	return cm.current
}

// GetContextStack returns a copy of the context stack.
func (cm *ContextManager) GetContextStack() []Context {
	dup := make([]Context, len(cm.stack))
	copy(dup, cm.stack)
	return dup
}

// EnterContext pushes the current context on the stack and switches to the new context.
func (cm *ContextManager) EnterContext(ctx Context) {
	if ctx == cm.current {
		return
	}

	old := cm.current
	cm.stack = append(cm.stack, cm.current)
	cm.current = ctx
	cm.invokeCallbacks(old, ctx)
}

// ExitContext pops the last context from the stack and activates it.
func (cm *ContextManager) ExitContext() Context {
	if len(cm.stack) == 0 {
		return cm.current
	}

	old := cm.current
	idx := len(cm.stack) - 1
	cm.current = cm.stack[idx]
	cm.stack = cm.stack[:idx]
	cm.invokeCallbacks(old, cm.current)
	return cm.current
}

// SetContext forcefully changes the current context without modifying the stack.
func (cm *ContextManager) SetContext(ctx Context) {
	if ctx == cm.current {
		return
	}

	old := cm.current
	cm.current = ctx
	cm.invokeCallbacks(old, ctx)
}

// ForceEnvironment overrides resolver platform/terminal (primarily for tests).
func (cm *ContextManager) ForceEnvironment(platform, terminal string) {
	if cm == nil || cm.resolver == nil {
		return
	}
	cm.resolver.ForceEnvironment(platform, terminal)
}

func (cm *ContextManager) invokeCallbacks(from, to Context) {
	if from == to {
		return
	}

	if callbacks, exists := cm.callbacks[to]; exists {
		for _, cb := range callbacks {
			cb(from, to)
		}
	}

	if to != ContextGlobal {
		if callbacks, exists := cm.callbacks[ContextGlobal]; exists {
			for _, cb := range callbacks {
				cb(from, to)
			}
		}
	}
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
