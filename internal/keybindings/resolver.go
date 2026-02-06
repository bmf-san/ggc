package keybindings

import (
	"fmt"
	"os"
	"strings"

	"github.com/bmf-san/ggc/v7/internal/config"
)

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
