package interactive

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func runeIndex(haystack, needle string) int {
	idx := strings.Index(haystack, needle)
	if idx < 0 {
		return -1
	}
	return utf8.RuneCountInString(haystack[:idx])
}

func TestUIStateContextManagement(t *testing.T) {
	t.Run("enter and exit context stack", func(t *testing.T) {
		state := &UIState{
			context:      ContextGlobal,
			contextStack: []Context{},
		}

		var transitions [][2]Context
		state.onContextChange = func(oldCtx, newCtx Context) {
			transitions = append(transitions, [2]Context{oldCtx, newCtx})
		}

		state.EnterContext(ContextInput)
		if state.context != ContextInput {
			t.Fatalf("want context %v, got %v", ContextInput, state.context)
		}
		if len(state.contextStack) != 1 || state.contextStack[0] != ContextGlobal {
			t.Fatalf("unexpected stack contents: %#v", state.contextStack)
		}

		state.EnterContext(ContextInput)
		if len(state.contextStack) != 1 {
			t.Fatalf("expected stack to remain unchanged, got %d", len(state.contextStack))
		}

		state.EnterContext(ContextResults)
		if got := state.GetCurrentContext(); got != ContextResults {
			t.Fatalf("want current context %v, got %v", ContextResults, got)
		}
		if len(state.contextStack) != 2 {
			t.Fatalf("want stack size 2, got %d", len(state.contextStack))
		}

		state.ExitContext()
		if state.context != ContextInput {
			t.Fatalf("want context %v after exit, got %v", ContextInput, state.context)
		}
		state.ExitContext()
		if state.context != ContextGlobal {
			t.Fatalf("want context %v after exit, got %v", ContextGlobal, state.context)
		}
		state.ExitContext()
		if state.context != ContextGlobal {
			t.Fatalf("context should remain %v when stack empty, got %v", ContextGlobal, state.context)
		}

		wantTransitions := [][2]Context{
			{ContextGlobal, ContextInput},
			{ContextInput, ContextResults},
			{ContextResults, ContextInput},
			{ContextInput, ContextGlobal},
		}
		if len(transitions) != len(wantTransitions) {
			t.Fatalf("want %d transitions, got %d", len(wantTransitions), len(transitions))
		}
		for i, want := range wantTransitions {
			if transitions[i] != want {
				t.Fatalf("transition %d: want %v, got %v", i, want, transitions[i])
			}
		}
	})

	t.Run("fallback to global when stack empty", func(t *testing.T) {
		state := &UIState{
			context:      ContextSearch,
			contextStack: nil,
		}

		transitioned := false
		state.onContextChange = func(oldCtx, newCtx Context) {
			transitioned = true
			if oldCtx != ContextSearch || newCtx != ContextGlobal {
				t.Fatalf("unexpected transition %v -> %v", oldCtx, newCtx)
			}
		}

		state.ExitContext()
		if state.context != ContextGlobal {
			t.Fatalf("want context %v, got %v", ContextGlobal, state.context)
		}
		if !transitioned {
			t.Fatalf("expected transition callback to run")
		}
	})
}

func TestUIStateModeHelpers(t *testing.T) {
	state := &UIState{context: ContextInput}
	if !state.IsInInputMode() {
		t.Fatalf("expected input mode")
	}
	state.context = ContextResults
	if !state.IsInResultsMode() {
		t.Fatalf("expected results mode")
	}
	state.context = ContextSearch
	if !state.IsInSearchMode() {
		t.Fatalf("expected search mode")
	}
}

func TestUIStateHasMatches(t *testing.T) {
	state := &UIState{}
	if state.HasMatches() {
		t.Fatalf("expected no matches initially")
	}
	state.filtered = []CommandInfo{{Command: "status"}}
	if !state.HasMatches() {
		t.Fatalf("expected matches after adding filtered command")
	}
}

func TestUIStateMoveLeftRight(t *testing.T) {
	input := "a界😊"
	state := &UIState{
		input:     input,
		cursorPos: utf8.RuneCountInString(input),
	}

	state.MoveLeft()
	if state.cursorPos != 2 {
		t.Fatalf("want cursor 2 after MoveLeft, got %d", state.cursorPos)
	}
	state.MoveLeft()
	if state.cursorPos != 1 {
		t.Fatalf("want cursor 1 after MoveLeft, got %d", state.cursorPos)
	}
	state.MoveLeft()
	if state.cursorPos != 0 {
		t.Fatalf("cursor should not move past beginning, got %d", state.cursorPos)
	}
	state.MoveLeft()
	if state.cursorPos != 0 {
		t.Fatalf("cursor should remain at beginning, got %d", state.cursorPos)
	}

	state.MoveRight()
	if state.cursorPos != 1 {
		t.Fatalf("want cursor 1 after MoveRight, got %d", state.cursorPos)
	}
	state.MoveRight()
	if state.cursorPos != 2 {
		t.Fatalf("want cursor 2 after MoveRight, got %d", state.cursorPos)
	}
	state.MoveRight()
	if state.cursorPos != 3 {
		t.Fatalf("want cursor 3 after MoveRight, got %d", state.cursorPos)
	}
	state.MoveRight()
	if state.cursorPos != 3 {
		t.Fatalf("cursor should not exceed rune length, got %d", state.cursorPos)
	}
}

func TestUIStateMoveWordBoundaries(t *testing.T) {
	text := "alpha  beta \tgamma"
	totalRunes := utf8.RuneCountInString(text)
	gammaIndex := runeIndex(text, "gamma")
	betaIndex := runeIndex(text, "beta")
	alphaIndex := runeIndex(text, "alpha")

	state := &UIState{
		input:     text,
		cursorPos: totalRunes,
	}

	state.MoveWordLeft()
	if state.cursorPos != gammaIndex {
		t.Fatalf("want cursor at gamma (%d), got %d", gammaIndex, state.cursorPos)
	}
	state.MoveWordLeft()
	if state.cursorPos != betaIndex {
		t.Fatalf("want cursor at beta (%d), got %d", betaIndex, state.cursorPos)
	}
	state.MoveWordLeft()
	if state.cursorPos != alphaIndex {
		t.Fatalf("want cursor at alpha (%d), got %d", alphaIndex, state.cursorPos)
	}
	state.MoveWordLeft()
	if state.cursorPos != alphaIndex {
		t.Fatalf("cursor should remain at first word, got %d", state.cursorPos)
	}

	state.cursorPos = 0
	state.MoveWordRight()
	if state.cursorPos != betaIndex {
		t.Fatalf("want cursor at beta (%d), got %d", betaIndex, state.cursorPos)
	}
	state.MoveWordRight()
	if state.cursorPos != gammaIndex {
		t.Fatalf("want cursor at gamma (%d), got %d", gammaIndex, state.cursorPos)
	}
	state.MoveWordRight()
	if state.cursorPos != totalRunes {
		t.Fatalf("cursor should move to end after last word, got %d", state.cursorPos)
	}
	state.MoveWordRight()
	if state.cursorPos != totalRunes {
		t.Fatalf("cursor should remain at end, got %d", state.cursorPos)
	}
}
