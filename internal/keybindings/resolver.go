// Package keybindings provides a configurable keybinding system for interactive mode.
// It supports profile-based configuration, platform-specific bindings, context-aware
// key mapping, and runtime profile switching.
package keybindings

import (
	"fmt"
	"strings"

	"github.com/bmf-san/ggc/v7/internal/config"
)

// ContextualKeyBindingMap holds resolved keybindings for all contexts

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
