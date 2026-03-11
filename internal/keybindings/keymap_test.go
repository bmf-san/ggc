package keybindings

import "testing"

// ── KeyBindingMap getter methods ────────────────────────────────────────────

func TestKeyBindingMap_GetDeleteToEndByte(t *testing.T) {
	km := DefaultKeyBindingMap()
	got := km.GetDeleteToEndByte()
	want := ctrl('k')
	if got != want {
		t.Errorf("GetDeleteToEndByte() = %d, want %d", got, want)
	}
}

func TestKeyBindingMap_GetMoveToBeginningByte(t *testing.T) {
	km := DefaultKeyBindingMap()
	got := km.GetMoveToBeginningByte()
	want := ctrl('a')
	if got != want {
		t.Errorf("GetMoveToBeginningByte() = %d, want %d", got, want)
	}
}

func TestKeyBindingMap_GetMoveToEndByte(t *testing.T) {
	km := DefaultKeyBindingMap()
	got := km.GetMoveToEndByte()
	want := ctrl('e')
	if got != want {
		t.Errorf("GetMoveToEndByte() = %d, want %d", got, want)
	}
}

func TestKeyBindingMap_GetMoveUpByte(t *testing.T) {
	km := DefaultKeyBindingMap()
	got := km.GetMoveUpByte()
	want := ctrl('p')
	if got != want {
		t.Errorf("GetMoveUpByte() = %d, want %d", got, want)
	}
}

func TestKeyBindingMap_GetMoveDownByte(t *testing.T) {
	km := DefaultKeyBindingMap()
	got := km.GetMoveDownByte()
	want := ctrl('n')
	if got != want {
		t.Errorf("GetMoveDownByte() = %d, want %d", got, want)
	}
}

func TestKeyBindingMap_GetAddToWorkflowByte(t *testing.T) {
	km := DefaultKeyBindingMap()
	got := km.GetAddToWorkflowByte()
	// Tab is ASCII 9; getFirstControlByte won't find a ctrl keystroke in a raw seq,
	// so it returns the fallback (9).
	if got != 9 {
		t.Errorf("GetAddToWorkflowByte() = %d, want 9 (Tab)", got)
	}
}

func TestKeyBindingMap_GetToggleWorkflowViewByte(t *testing.T) {
	km := DefaultKeyBindingMap()
	got := km.GetToggleWorkflowViewByte()
	want := ctrl('t')
	if got != want {
		t.Errorf("GetToggleWorkflowViewByte() = %d, want %d", got, want)
	}
}

func TestKeyBindingMap_GetClearWorkflowByte(t *testing.T) {
	km := DefaultKeyBindingMap()
	got := km.GetClearWorkflowByte()
	// ClearWorkflow uses NewCharKeyStroke('c') which is a raw seq, not ctrl.
	// So it falls back to 'c'.
	if got != 'c' {
		t.Errorf("GetClearWorkflowByte() = %d (%q), want %d (%q)", got, got, byte('c'), 'c')
	}
}

func TestKeyBindingMap_GetDeleteToEndByte_EmptySlice(t *testing.T) {
	km := &KeyBindingMap{DeleteToEnd: []KeyStroke{}}
	got := km.GetDeleteToEndByte()
	want := ctrl('k') // fallback
	if got != want {
		t.Errorf("GetDeleteToEndByte (empty) = %d, want %d", got, want)
	}
}
