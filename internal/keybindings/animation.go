// Package keybindings provides a configurable keybinding system for interactive mode.
// It supports profile-based configuration, platform-specific bindings, context-aware
// key mapping, and runtime profile switching.
package keybindings

import (
	"fmt"
	"time"
)

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
