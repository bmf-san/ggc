package keybindings

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
