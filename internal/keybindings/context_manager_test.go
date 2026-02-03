package keybindings

import (
	"testing"

	"github.com/bmf-san/ggc/v7/internal/config"
)

func TestContextManagerStackAndCallbacks(t *testing.T) {
	resolver := NewKeyBindingResolver(&config.Config{})
	RegisterBuiltinProfiles(resolver)
	cm := NewContextManager(resolver)

	var targetTransitions [][2]Context
	cm.RegisterContextCallback(ContextResults, func(oldCtx, newCtx Context) {
		targetTransitions = append(targetTransitions, [2]Context{oldCtx, newCtx})
	})

	globalCalls := 0
	cm.RegisterContextCallback(ContextGlobal, func(oldCtx, newCtx Context) {
		globalCalls++
	})

	if cm.GetCurrentContext() != ContextGlobal {
		t.Fatalf("initial context = %v, want %v", cm.GetCurrentContext(), ContextGlobal)
	}

	cm.EnterContext(ContextInput)
	if cm.GetCurrentContext() != ContextInput {
		t.Fatalf("after enter input got %v", cm.GetCurrentContext())
	}

	if stack := cm.GetContextStack(); len(stack) != 1 || stack[0] != ContextGlobal {
		t.Fatalf("unexpected stack after first enter: %#v", stack)
	}

	cm.EnterContext(ContextResults)
	if cm.GetCurrentContext() != ContextResults {
		t.Fatalf("after enter results got %v", cm.GetCurrentContext())
	}

	if len(targetTransitions) != 1 {
		t.Fatalf("expected 1 results transition, got %d", len(targetTransitions))
	}

	if globalCalls != 2 { // ContextGlobal callback fires for input and results
		t.Fatalf("expected 2 global callbacks, got %d", globalCalls)
	}

	if got := cm.ExitContext(); got != ContextInput {
		t.Fatalf("exit context returned %v, want %v", got, ContextInput)
	}

	if got := cm.ExitContext(); got != ContextGlobal {
		t.Fatalf("stack exit returned %v, want %v", got, ContextGlobal)
	}

	// Exiting with empty stack keeps current context and should not panic
	if got := cm.ExitContext(); got != ContextGlobal {
		t.Fatalf("exit on empty stack returned %v", got)
	}
}
