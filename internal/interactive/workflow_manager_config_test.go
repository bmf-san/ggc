package interactive

import (
	"testing"
)

func TestLoadFromConfig_Nil(t *testing.T) {
	mgr := NewWorkflowManager()
	mgr.LoadFromConfig(nil) // must not panic

	summaries := mgr.ListWorkflows()
	if len(summaries) != 1 {
		t.Errorf("expected 1 workflow (scratch only), got %d", len(summaries))
	}
}

func TestLoadFromConfig_Empty(t *testing.T) {
	mgr := NewWorkflowManager()
	mgr.LoadFromConfig(map[string][]string{})

	summaries := mgr.ListWorkflows()
	if len(summaries) != 1 {
		t.Errorf("expected 1 workflow (scratch only), got %d", len(summaries))
	}
}

func TestLoadFromConfig_SingleWorkflow(t *testing.T) {
	mgr := NewWorkflowManager()
	mgr.LoadFromConfig(map[string][]string{
		"deploy": {"add .", "commit <message>", "push current"},
	})

	summaries := mgr.ListWorkflows()
	if len(summaries) != 2 {
		t.Fatalf("expected 2 workflows (scratch + deploy), got %d", len(summaries))
	}

	var deployID int
	for _, s := range summaries {
		if s.Name == "deploy" {
			deployID = s.ID
			break
		}
	}
	if deployID == 0 {
		t.Fatal("expected workflow named 'deploy' to exist")
	}

	wf, ok := mgr.GetWorkflow(deployID)
	if !ok {
		t.Fatal("failed to retrieve deploy workflow by ID")
	}

	steps := wf.GetSteps()
	if len(steps) != 3 {
		t.Fatalf("expected 3 steps, got %d", len(steps))
	}
	if steps[0].Command != "add" {
		t.Errorf("step[0].Command = %q, want %q", steps[0].Command, "add")
	}
	if len(steps[0].Args) != 1 || steps[0].Args[0] != "." {
		t.Errorf("step[0].Args = %v, want [.]", steps[0].Args)
	}
	if steps[1].Command != "commit" {
		t.Errorf("step[1].Command = %q, want %q", steps[1].Command, "commit")
	}
	if len(steps[1].Args) != 1 || steps[1].Args[0] != "<message>" {
		t.Errorf("step[1].Args = %v, want [<message>]", steps[1].Args)
	}
	if steps[2].Command != "push" {
		t.Errorf("step[2].Command = %q, want %q", steps[2].Command, "push")
	}
}

func TestLoadFromConfig_ScratchWorkflowRemainsActive(t *testing.T) {
	mgr := NewWorkflowManager()
	initialActiveID := mgr.GetActiveID()

	mgr.LoadFromConfig(map[string][]string{
		"acp": {"add .", "commit", "push current"},
	})

	if mgr.GetActiveID() != initialActiveID {
		t.Errorf("active ID changed after LoadFromConfig: got %d, want %d",
			mgr.GetActiveID(), initialActiveID)
	}
}

func TestLoadFromConfig_MultipleWorkflows(t *testing.T) {
	mgr := NewWorkflowManager()
	mgr.LoadFromConfig(map[string][]string{
		"acp":          {"add .", "commit", "push current"},
		"fetch-rebase": {"fetch origin", "rebase origin main"},
	})

	summaries := mgr.ListWorkflows()
	// 1 scratch + 2 config-defined
	if len(summaries) != 3 {
		t.Errorf("expected 3 workflows, got %d", len(summaries))
	}
}

func TestLoadFromConfig_StepParsing(t *testing.T) {
	mgr := NewWorkflowManager()
	mgr.LoadFromConfig(map[string][]string{
		"test": {"push origin main"},
	})

	var testID int
	for _, s := range mgr.ListWorkflows() {
		if s.Name == "test" {
			testID = s.ID
			break
		}
	}
	if testID == 0 {
		t.Fatal("workflow 'test' not found")
	}

	wf, _ := mgr.GetWorkflow(testID)
	steps := wf.GetSteps()
	if len(steps) != 1 {
		t.Fatalf("expected 1 step, got %d", len(steps))
	}

	s := steps[0]
	if s.Command != "push" {
		t.Errorf("Command = %q, want %q", s.Command, "push")
	}
	if len(s.Args) != 2 || s.Args[0] != "origin" || s.Args[1] != "main" {
		t.Errorf("Args = %v, want [origin main]", s.Args)
	}
	if s.Description != "push origin main" {
		t.Errorf("Description = %q, want %q", s.Description, "push origin main")
	}
}
