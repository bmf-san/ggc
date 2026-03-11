package keybindings

import (
	"testing"
	"time"

	"github.com/bmf-san/ggc/v8/internal/config"
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

// ── ContextManager: SetContext, ForceEnvironment ─────────────────────────────

func TestContextManager_SetContext(t *testing.T) {
	resolver := NewKeyBindingResolver(&config.Config{})
	RegisterBuiltinProfiles(resolver)
	cm := NewContextManager(resolver)

	var transitions [][2]Context
	cm.RegisterContextCallback(ContextResults, func(from, to Context) {
		transitions = append(transitions, [2]Context{from, to})
	})

	// SetContext to new context
	cm.SetContext(ContextResults)
	if cm.GetCurrentContext() != ContextResults {
		t.Errorf("SetContext: current = %v, want %v", cm.GetCurrentContext(), ContextResults)
	}
	if len(transitions) != 1 {
		t.Errorf("expected 1 transition callback, got %d", len(transitions))
	}

	// SetContext to same context should be no-op
	cm.SetContext(ContextResults)
	if len(transitions) != 1 {
		t.Errorf("SetContext same context should not fire callback, got %d transitions", len(transitions))
	}

	// Stack should be unmodified
	if len(cm.GetContextStack()) != 0 {
		t.Errorf("SetContext should not modify stack, got %v", cm.GetContextStack())
	}
}

func TestContextManager_ForceEnvironment(t *testing.T) {
	resolver := NewKeyBindingResolver(&config.Config{})
	RegisterBuiltinProfiles(resolver)
	cm := NewContextManager(resolver)

	// Should not panic
	cm.ForceEnvironment("darwin", "xterm-256color")
}

func TestContextManager_ForceEnvironment_NilCM(t *testing.T) {
	var cm *ContextManager
	// Should not panic
	cm.ForceEnvironment("linux", "xterm")
}

// ── ContextTransitionAnimator ────────────────────────────────────────────────

func TestContextTransitionAnimator_FadeAndSlide(t *testing.T) {
	cta := NewContextTransitionAnimator()
	cta.SetDuration(0) // no sleep in tests

	cta.SetStyle("fade")
	cta.AnimateTransition(ContextGlobal, ContextResults)

	cta.SetStyle("slide")
	cta.AnimateTransition(ContextGlobal, ContextInput)
}

func TestContextTransitionAnimator_Disable(t *testing.T) {
	cta := NewContextTransitionAnimator()
	cta.Disable()
	// Should return early without doing anything
	cta.AnimateTransition(ContextGlobal, ContextResults)
	if cta.enabled {
		t.Error("expected disabled animator")
	}
}

func TestContextTransitionAnimator_Enable(t *testing.T) {
	cta := NewContextTransitionAnimator()
	cta.Disable()
	cta.Enable()
	if !cta.enabled {
		t.Error("expected enabled animator")
	}
}

func TestContextTransitionAnimator_RegisterAnimation(t *testing.T) {
	cta := NewContextTransitionAnimator()
	cta.RegisterAnimation(func(from, to Context) {})
	cta.RegisterAnimation(func(from, to Context) {})
	if len(cta.animations) != 2 {
		t.Errorf("expected 2 registered animations, got %d", len(cta.animations))
	}
}

func TestContextTransitionAnimator_SetDuration(t *testing.T) {
	cta := NewContextTransitionAnimator()
	cta.SetDuration(500 * time.Millisecond)
	if cta.duration != 500*time.Millisecond {
		t.Errorf("duration = %v, want 500ms", cta.duration)
	}
}
