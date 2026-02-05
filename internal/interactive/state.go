// Package interactive houses interactive UI types and helpers shared across the application.
package interactive

import (
	"sort"
	"strings"

	kb "github.com/bmf-san/ggc/v7/internal/keybindings"
)

// UIMode describes the high-level mode of the interactive UI.
type UIMode int

const (
	// ModeSearch renders the classic search/execute interface.
	ModeSearch UIMode = iota
	// ModeWorkflow renders the dedicated workflow management interface.
	ModeWorkflow
)

// WorkflowFocus indicates which pane in workflow mode has focus.
type WorkflowFocus int

const (
	// FocusInput targets the command input/results pane.
	FocusInput WorkflowFocus = iota
	// FocusWorkflowList targets the workflow list pane.
	FocusWorkflowList
)

// UIState holds the current state of the interactive UI
type UIState struct {
	selected        int
	input           string
	cursorPos       int // Cursor position in input string
	filtered        []CommandInfo
	context         kb.Context   // Current UI context (input/results/search/global)
	contextStack    []kb.Context // Context stack for nested states
	onContextChange func(kb.Context, kb.Context)
	mode            UIMode
	workflowFocus   WorkflowFocus
	workflowListIdx int
	workflowOffset  int
}

// SetMode switches between search and workflow modes.
func (s *UIState) SetMode(mode UIMode) {
	if s.mode == mode {
		return
	}
	s.mode = mode
}

// IsWorkflowMode reports whether the UI is currently in workflow mode.
func (s *UIState) IsWorkflowMode() bool {
	return s.mode == ModeWorkflow
}

// FocusInput moves focus to the command input/results pane.
func (s *UIState) FocusInput() {
	s.workflowFocus = FocusInput
}

// FocusWorkflowList moves focus to the workflow list pane.
func (s *UIState) FocusWorkflowList() {
	s.workflowFocus = FocusWorkflowList
}

// IsInputFocused reports whether the command input/results pane has focus.
func (s *UIState) IsInputFocused() bool {
	return s.workflowFocus == FocusInput
}

// SetWorkflowListIndex sets and clamps the workflow list selection.
func (s *UIState) SetWorkflowListIndex(idx, total int) {
	if total <= 0 {
		s.workflowListIdx = 0
		s.workflowOffset = 0
		return
	}
	if idx < 0 {
		idx = 0
	}
	if idx >= total {
		idx = total - 1
	}
	s.workflowListIdx = idx
}

// UpdateFiltered updates the filtered commands based on current input using fuzzy matching
func (s *UIState) UpdateFiltered() {
	input := strings.ToLower(s.input)
	if input == "" {
		s.filtered = make([]CommandInfo, len(commands))
		copy(s.filtered, commands)
	} else {
		type match struct {
			info  CommandInfo
			score matchScore
		}
		matches := make([]match, 0, len(commands))
		for _, cmd := range commands {
			cmdLower := strings.ToLower(cmd.Command)
			if ok, score := fuzzyMatchScore(cmdLower, input); ok {
				matches = append(matches, match{info: cmd, score: score})
			}
		}
		sort.SliceStable(matches, func(i, j int) bool {
			return matches[i].score.less(matches[j].score)
		})
		s.filtered = make([]CommandInfo, len(matches))
		for i, match := range matches {
			s.filtered[i] = match.info
		}
	}
	// Reset selection if out of bounds
	if s.selected >= len(s.filtered) {
		s.selected = len(s.filtered) - 1
	}
	if s.selected < 0 {
		s.selected = 0
	}
}

// Context Management Methods

// EnterContext pushes the current context onto the stack and switches to the new context
func (s *UIState) EnterContext(newContext kb.Context) {
	if s.context == newContext {
		return
	}
	old := s.context
	s.contextStack = append(s.contextStack, s.context)
	s.context = newContext
	s.notifyContextChange(old, newContext)
}

// ExitContext pops the previous context from the stack
func (s *UIState) ExitContext() {
	if len(s.contextStack) > 0 {
		old := s.context
		s.context = s.contextStack[len(s.contextStack)-1]
		s.contextStack = s.contextStack[:len(s.contextStack)-1]
		s.notifyContextChange(old, s.context)
	} else if s.context != kb.ContextGlobal {
		old := s.context
		s.context = kb.ContextGlobal
		s.notifyContextChange(old, s.context)
	}
}

// GetCurrentContext returns the current UI context
func (s *UIState) GetCurrentContext() kb.Context {
	return s.context
}

// SetContext directly sets the context (use with caution)
func (s *UIState) SetContext(ctx kb.Context) {
	if s.context == ctx {
		return
	}
	old := s.context
	s.context = ctx
	s.notifyContextChange(old, ctx)
}

// notifyContextChange triggers the callback when the active context changes
func (s *UIState) notifyContextChange(oldCtx, newCtx kb.Context) {
	if s.onContextChange != nil && oldCtx != newCtx {
		s.onContextChange(oldCtx, newCtx)
	}
}

// IsInInputMode returns true if currently in input context
func (s *UIState) IsInInputMode() bool {
	return s.context == kb.ContextInput
}

// IsInResultsMode returns true if currently in results context
func (s *UIState) IsInResultsMode() bool {
	return s.context == kb.ContextResults
}

// IsInSearchMode returns true if currently in search context
func (s *UIState) IsInSearchMode() bool {
	return s.context == kb.ContextSearch
}

// MoveUp moves selection up
func (s *UIState) MoveUp() {
	// Switch to results context when navigating
	if s.context != kb.ContextResults && s.context != kb.ContextSearch {
		s.SetContext(kb.ContextResults)
	}

	if s.selected > 0 {
		s.selected--
	}
}

// MoveDown moves selection down
func (s *UIState) MoveDown() {
	// Switch to results context when navigating
	if s.context != kb.ContextResults && s.context != kb.ContextSearch {
		s.SetContext(kb.ContextResults)
	}

	if s.selected < len(s.filtered)-1 {
		s.selected++
	}
}

// GetSelectedCommand returns the currently selected command
func (s *UIState) GetSelectedCommand() *CommandInfo {
	if len(s.filtered) > 0 && s.selected >= 0 && s.selected < len(s.filtered) {
		return &s.filtered[s.selected]
	}
	return nil
}

// HasInput returns true if there is input
func (s *UIState) HasInput() bool {
	return s.input != ""
}

// HasMatches returns true if there are filtered matches
func (s *UIState) HasMatches() bool {
	return len(s.filtered) > 0
}
