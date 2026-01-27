package interactive

import (
	"strings"
	"testing"
)

func TestWorkflowManagerDeleteWorkflow(t *testing.T) {
	mgr := NewWorkflowManager()
	active := mgr.GetActiveID()
	if active == 0 {
		t.Fatal("expected initial active workflow")
	}

	newID, ok := mgr.DeleteWorkflow(active)
	if !ok {
		t.Fatal("expected delete to succeed")
	}
	if newID != 0 {
		t.Fatalf("expected no active workflow after deleting last, got %d", newID)
	}
	if mgr.GetActiveID() != 0 {
		t.Fatalf("expected active workflow to be cleared, got %d", mgr.GetActiveID())
	}
	if len(mgr.ListWorkflows()) != 0 {
		t.Fatal("expected no workflows after deleting last")
	}

	if _, ok := mgr.DeleteWorkflow(999); ok {
		t.Fatal("expected delete to fail for missing workflow")
	}
}

func TestWorkflowManagerCloneWorkflow(t *testing.T) {
	mgr := NewWorkflowManager()
	active := mgr.GetActiveID()
	if _, ok := mgr.AddStep(active, "add", []string{"."}, "add ."); !ok {
		t.Fatal("expected add step to succeed")
	}

	cloneID, ok := mgr.CloneWorkflow(active, "")
	if !ok {
		t.Fatal("expected clone to succeed")
	}
	if cloneID == active {
		t.Fatal("expected cloned workflow to have a new ID")
	}

	clone, ok := mgr.GetWorkflow(cloneID)
	if !ok || clone == nil {
		t.Fatal("expected cloned workflow to exist")
	}
	if clone.Size() != 1 {
		t.Fatalf("expected cloned workflow to have 1 step, got %d", clone.Size())
	}
	if step := clone.GetSteps()[0]; step.Description != "add ." {
		t.Fatalf("expected cloned step description to be preserved, got %q", step.Description)
	}

	found := false
	for _, summary := range mgr.ListWorkflows() {
		if summary.ID == cloneID {
			found = true
			if !strings.Contains(summary.Name, "copy") {
				t.Fatalf("expected cloned workflow name to include copy, got %q", summary.Name)
			}
		}
	}
	if !found {
		t.Fatal("expected clone summary to be present")
	}

	if !mgr.SetActive(cloneID) {
		t.Fatal("expected SetActive to succeed for cloned workflow")
	}
	if mgr.GetActiveID() != cloneID {
		t.Fatalf("expected active workflow to be %d, got %d", cloneID, mgr.GetActiveID())
	}
	if mgr.SetActive(12345) {
		t.Fatal("expected SetActive to fail for missing workflow")
	}
}
