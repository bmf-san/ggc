package interactive

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v7/pkg/config"
)

func TestWorkflow_AddStep(t *testing.T) {
	workflow := NewWorkflow()

	id1 := workflow.AddStep("add", []string{"."}, "add .")
	id2 := workflow.AddStep("commit", []string{"-m", "test"}, "commit -m test")

	if id1 != 1 {
		t.Errorf("Expected first ID to be 1, got %d", id1)
	}
	if id2 != 2 {
		t.Errorf("Expected second ID to be 2, got %d", id2)
	}

	steps := workflow.GetSteps()
	if len(steps) != 2 {
		t.Errorf("Expected 2 steps, got %d", len(steps))
	}

	if steps[0].Command != "add" {
		t.Errorf("Expected first command 'add', got '%s'", steps[0].Command)
	}
	if steps[1].Command != "commit" {
		t.Errorf("Expected second command 'commit', got '%s'", steps[1].Command)
	}
}

func TestWorkflow_Clear(t *testing.T) {
	workflow := NewWorkflow()

	workflow.AddStep("add", []string{"."}, "add .")
	workflow.AddStep("commit", []string{"-m", "test"}, "commit -m test")

	if workflow.IsEmpty() {
		t.Error("Expected workflow to not be empty")
	}
	if workflow.Size() != 2 {
		t.Errorf("Expected size 2, got %d", workflow.Size())
	}

	workflow.Clear()

	if !workflow.IsEmpty() {
		t.Error("Expected workflow to be empty after clear")
	}
	if workflow.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", workflow.Size())
	}
}

func TestWorkflowStep_String(t *testing.T) {
	// Test with description
	step1 := WorkflowStep{
		ID:          1,
		Command:     "commit",
		Args:        []string{"-m", "test message"},
		Description: "commit -m test message",
	}

	result1 := step1.String()
	expected1 := "[1] commit -m test message"

	if result1 != expected1 {
		t.Errorf("Expected '%s', got '%s'", expected1, result1)
	}

	// Test without description
	step2 := WorkflowStep{
		ID:      2,
		Command: "push",
		Args:    []string{"origin", "main"},
	}

	result2 := step2.String()
	expected2 := "[2] push origin main"

	if result2 != expected2 {
		t.Errorf("Expected '%s', got '%s'", expected2, result2)
	}
}

// Mock Router for testing workflow execution
type mockWorkflowRouter struct {
	executedCommands [][]string
}

func (m *mockWorkflowRouter) Route(args []string) {
	m.executedCommands = append(m.executedCommands, args)
}

func TestWorkflowExecutor_Execute(t *testing.T) {
	mock := &mockWorkflowRouter{}
	executor := NewWorkflowExecutor(mock, nil)
	workflow := NewWorkflow()

	// Test empty workflow
	err := executor.Execute(workflow)
	if err == nil {
		t.Error("Expected error for empty workflow")
	}

	// Add steps to workflow
	workflow.AddStep("add", []string{"."}, "add .")
	workflow.AddStep("commit", []string{"-m", "test"}, "commit -m test")
	workflow.AddStep("push", []string{}, "push")

	// Execute workflow
	err = executor.Execute(workflow)
	if err != nil {
		t.Errorf("Unexpected error executing workflow: %v", err)
	}

	// Check that all commands were routed in order
	expectedCommands := [][]string{
		{"add", "."},
		{"commit", "-m", "test"},
		{"push"},
	}

	if len(mock.executedCommands) != len(expectedCommands) {
		t.Errorf("Expected %d commands executed, got %d", len(expectedCommands), len(mock.executedCommands))
	}

	for i, expectedCmd := range expectedCommands {
		if len(mock.executedCommands[i]) != len(expectedCmd) {
			t.Errorf("Expected command %d to have %d args, got %d", i, len(expectedCmd), len(mock.executedCommands[i]))
			continue
		}
		for j, expectedArg := range expectedCmd {
			if mock.executedCommands[i][j] != expectedArg {
				t.Errorf("Expected command %d arg %d to be '%s', got '%s'", i, j, expectedArg, mock.executedCommands[i][j])
			}
		}
	}
}

func TestWorkflowExecutor_ExecuteCanceled(t *testing.T) {
	colors := NewANSIColors()
	workflowMgr := NewWorkflowManager()
	ui := &UI{
		stdin:       strings.NewReader("\n"),
		stdout:      &bytes.Buffer{},
		stderr:      &bytes.Buffer{},
		colors:      colors,
		workflowMgr: workflowMgr,
		term:        &mockTerminal{shouldFailRaw: true},
	}

	handler := &KeyHandler{ui: ui}
	ui.handler = handler

	workflow := NewWorkflow()
	workflow.AddStep("commit", nil, "commit <message>")

	executor := NewWorkflowExecutor(&mockWorkflowRouter{}, ui)
	err := executor.Execute(workflow)
	if !errors.Is(err, ErrWorkflowCanceled) {
		t.Fatalf("expected workflow cancellation error, got %v", err)
	}
}

func TestInteractiveInputForWorkflowScanner(t *testing.T) {
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	_, _ = w.WriteString("message\n")
	_ = w.Close()
	os.Stdin = r

	inputs, canceled := interactiveInputForWorkflow(nil, []string{"message"})
	if canceled {
		t.Fatal("expected scanner fallback to succeed")
	}
	if got := inputs["message"]; got != "message" {
		t.Fatalf("expected message 'message', got %q", got)
	}

	_ = r.Close()
}

func TestInteractiveInputForWorkflowScannerCanceled(t *testing.T) {
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	_, _ = w.WriteString("\n")
	_ = w.Close()
	os.Stdin = r

	inputs, canceled := interactiveInputForWorkflow(nil, []string{"message"})
	if !canceled {
		t.Fatal("expected cancellation when placeholder input is empty")
	}
	if inputs != nil {
		t.Fatal("expected nil inputs on cancellation")
	}

	_ = r.Close()
}

// TestWorkflowExecutor_ExecuteEmptyWorkflow tests executing an empty workflow
func TestWorkflowExecutor_ExecuteEmptyWorkflow(t *testing.T) {
	// Setup
	mockRouter := &mockRouterNew{}
	executor := NewWorkflowExecutor(mockRouter, nil)
	workflow := NewWorkflow()

	// Execute
	err := executor.Execute(workflow)

	// Verify
	if err == nil {
		t.Error("Expected error when executing empty workflow")
	}
	if err.Error() != "workflow is empty" {
		t.Errorf("Expected 'workflow is empty' error, got '%s'", err.Error())
	}
}

// TestWorkflow_ConcurrentAccess tests concurrent access to workflow
func TestWorkflow_ConcurrentAccess(t *testing.T) {
	workflow := NewWorkflow()

	// Test concurrent adds
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			workflow.AddStep("test", []string{}, fmt.Sprintf("test %d", id))
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all steps were added
	steps := workflow.GetSteps()
	if len(steps) != 10 {
		t.Errorf("Expected 10 steps, got %d", len(steps))
	}
}

// mockRouterNew is a mock implementation of CommandRouter for testing
type mockRouterNew struct {
	routedCommands [][]string
}

func (m *mockRouterNew) Route(args []string) {
	if m.routedCommands == nil {
		m.routedCommands = make([][]string, 0)
	}
	m.routedCommands = append(m.routedCommands, args)
}

func TestWorkflowManagerInitialization(t *testing.T) {
	manager := NewWorkflowManager()

	if got := manager.WorkflowCount(); got != 1 {
		t.Fatalf("expected manager to start with 1 workflow, got %d", got)
	}

	if _, id := manager.GetActiveWorkflow(); id != 1 {
		t.Fatalf("expected initial active workflow ID 1, got %d", id)
	}
}

func TestWorkflowManagerCreateAndDelete(t *testing.T) {
	manager := NewWorkflowManager()

	id2 := manager.CreateWorkflow()
	if id2 == 0 {
		t.Fatal("expected non-zero workflow ID for created workflow")
	}

	if count := manager.WorkflowCount(); count != 2 {
		t.Fatalf("expected workflow count 2, got %d", count)
	}

	if active := manager.GetActiveID(); active != id2 {
		t.Fatalf("expected newly created workflow to be active, got %d", active)
	}

	_, ok := manager.DeleteWorkflow(id2)
	if !ok {
		t.Fatal("expected delete to succeed")
	}

	if count := manager.WorkflowCount(); count != 1 {
		t.Fatalf("expected workflow count 1 after delete, got %d", count)
	}

	if active := manager.GetActiveID(); active == id2 {
		t.Fatalf("expected active workflow to change after delete, still %d", active)
	}
}

func TestWorkflowManagerAddStep(t *testing.T) {
	manager := NewWorkflowManager()
	active, activeID := manager.GetActiveWorkflow()
	if active == nil {
		t.Fatal("expected active workflow to be non-nil")
	}

	stepID, err := manager.AddStep(activeID, "status", nil, "status")
	if err != nil {
		t.Fatalf("unexpected error adding step: %v", err)
	}
	if stepID != 1 {
		t.Fatalf("expected first step ID to be 1, got %d", stepID)
	}

	summaries := manager.ListWorkflows()
	if len(summaries) != 1 {
		t.Fatalf("expected 1 workflow summary, got %d", len(summaries))
	}
	if summaries[0].StepCount != 1 {
		t.Fatalf("expected step count 1, got %d", summaries[0].StepCount)
	}
	if !summaries[0].IsActive {
		t.Fatal("expected summary to mark workflow as active")
	}
	if summaries[0].Source != WorkflowSourceDynamic {
		t.Fatalf("expected workflow source dynamic, got %v", summaries[0].Source)
	}
	if summaries[0].ReadOnly {
		t.Fatal("expected dynamic workflow to be writable")
	}
}

func TestWorkflowManagerReadOnlyWorkflows(t *testing.T) {
	manager := NewWorkflowManager()

	id, err := manager.CreateReadOnlyWorkflow("config-release", []string{"status"})
	if err != nil {
		t.Fatalf("unexpected error creating config workflow: %v", err)
	}

	summaries := manager.ListWorkflows()
	var configSummary *WorkflowSummary
	for i := range summaries {
		if summaries[i].ID == id {
			configSummary = &summaries[i]
			break
		}
	}
	if configSummary == nil {
		t.Fatalf("expected to find config workflow summary for ID %d", id)
	}
	if !configSummary.ReadOnly {
		t.Fatal("expected config workflow to be read-only")
	}
	if configSummary.Source != WorkflowSourceConfig {
		t.Fatalf("expected config workflow source, got %v", configSummary.Source)
	}

	if _, err := manager.AddStep(id, "add", nil, "add"); err == nil {
		t.Fatal("expected add step to read-only workflow to fail")
	}
	if manager.ClearWorkflow(id) {
		t.Fatal("expected clear workflow on read-only workflow to be rejected")
	}
	if _, ok := manager.DeleteWorkflow(id); ok {
		t.Fatal("expected delete to fail for read-only workflow")
	}
}

func TestWorkflowManagerCreateWorkflowFromTemplates(t *testing.T) {
	manager := NewWorkflowManager()
	id, err := manager.CreateWorkflowFromTemplates("cloned", []string{"add .", "commit <message>"})
	if err != nil {
		t.Fatalf("unexpected error creating workflow from templates: %v", err)
	}
	summaries := manager.ListWorkflows()
	var summary *WorkflowSummary
	for i := range summaries {
		if summaries[i].ID == id {
			summary = &summaries[i]
			break
		}
	}
	if summary == nil {
		t.Fatalf("expected summary for workflow %d", id)
	}
	if summary.StepCount != 2 {
		t.Fatalf("expected step count 2, got %d", summary.StepCount)
	}
	if summary.ReadOnly {
		t.Fatal("expected workflow created from templates to be writable")
	}
	if summary.Source != WorkflowSourceDynamic {
		t.Fatalf("expected dynamic source, got %v", summary.Source)
	}
}

func TestUI_bootstrapConfigWorkflows(t *testing.T) {
	ui := &UI{
		workflowMgr: NewWorkflowManager(),
		state:       &UIState{},
		stderr:      &bytes.Buffer{},
		config: &config.Config{
			Workflows: []config.WorkflowConfig{
				{
					Name:  "release",
					Steps: []string{"status", "commit <message>", "push current"},
				},
			},
		},
	}

	ui.bootstrapConfigWorkflows()

	summaries := ui.listWorkflows()
	if len(summaries) != 2 {
		t.Fatalf("expected 2 workflows (default + config), got %d", len(summaries))
	}

	var found bool
	for _, summary := range summaries {
		if summary.Name == "release" {
			found = true
			if summary.Source != WorkflowSourceConfig {
				t.Fatalf("expected config source, got %v", summary.Source)
			}
			if !summary.ReadOnly {
				t.Fatal("expected config workflow to be read-only")
			}
			break
		}
	}
	if !found {
		t.Fatal("expected to find config-defined workflow in summaries")
	}
}
