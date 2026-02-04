package interactive

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
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
	ui := &UI{
		stdin:    strings.NewReader("\n"),
		stdout:   &bytes.Buffer{},
		stderr:   &bytes.Buffer{},
		colors:   colors,
		workflow: NewWorkflow(),
		term:     &mockTerminal{shouldFailRaw: true},
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

// TestWorkflowExecutor_MultiwordPlaceholder tests that multiword placeholder values are preserved as single arguments
func TestWorkflowExecutor_MultiwordPlaceholder(t *testing.T) {
	// Setup stdin with multiword input
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	// User enters "my commit message" for the <message> placeholder
	_, _ = w.WriteString("my commit message\n")
	_ = w.Close()
	os.Stdin = r

	mock := &mockWorkflowRouter{}
	executor := NewWorkflowExecutor(mock, nil)
	workflow := NewWorkflow()

	// Add step with placeholder in Args - this is the proper way to preserve multiword values
	workflow.AddStep("commit", []string{"-m", "<message>"}, "commit -m <message>")

	err = executor.Execute(workflow)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that "my commit message" was preserved as a single argument
	if len(mock.executedCommands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(mock.executedCommands))
	}

	cmd := mock.executedCommands[0]
	// Expected: ["commit", "-m", "my commit message"]
	if len(cmd) != 3 {
		t.Fatalf("expected 3 args, got %d: %v", len(cmd), cmd)
	}
	if cmd[0] != "commit" {
		t.Errorf("expected command 'commit', got '%s'", cmd[0])
	}
	if cmd[1] != "-m" {
		t.Errorf("expected arg 1 '-m', got '%s'", cmd[1])
	}
	if cmd[2] != "my commit message" {
		t.Errorf("expected arg 2 'my commit message', got '%s'", cmd[2])
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
