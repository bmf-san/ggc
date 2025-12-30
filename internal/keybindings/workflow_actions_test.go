package keybindings

import (
	"testing"

	"github.com/bmf-san/ggc/v7/pkg/config"
)

func TestApplyWorkflowActions(t *testing.T) {
	resolver := NewKeyBindingResolver(&config.Config{})
	keyMap := DefaultKeyBindingMap()

	create := []KeyStroke{NewCtrlKeyStroke('n')}
	resolver.applyWorkflowAction(keyMap, "workflow_create", create)
	if len(keyMap.WorkflowCreate) != 1 || keyMap.WorkflowCreate[0].Kind != KeyStrokeCtrl || keyMap.WorkflowCreate[0].Rune != 'n' {
		t.Fatalf("expected workflow_create to be applied, got %#v", keyMap.WorkflowCreate)
	}

	deleteKeys := []KeyStroke{NewCtrlKeyStroke('d')}
	resolver.applyUserWorkflowAction(keyMap, "workflow_delete", deleteKeys)
	if len(keyMap.WorkflowDelete) != 1 || keyMap.WorkflowDelete[0].Kind != KeyStrokeCtrl || keyMap.WorkflowDelete[0].Rune != 'd' {
		t.Fatalf("expected workflow_delete to be applied, got %#v", keyMap.WorkflowDelete)
	}
}
