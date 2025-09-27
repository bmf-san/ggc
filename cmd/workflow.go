// Package cmd provides workflow functionality for sequential command execution.
package cmd

import (
	"fmt"
	"strings"
	"sync"
)

// WorkflowStep represents a single step in a workflow
type WorkflowStep struct {
	ID          int      `json:"id"`
	Command     string   `json:"command"`
	Args        []string `json:"args"`
	Description string   `json:"description"`
}

// String returns a string representation of the workflow step
func (ws *WorkflowStep) String() string {
	if ws.Description != "" {
		return fmt.Sprintf("[%d] %s", ws.ID, ws.Description)
	}

	cmdStr := ws.Command
	if len(ws.Args) > 0 {
		cmdStr += " " + strings.Join(ws.Args, " ")
	}
	return fmt.Sprintf("[%d] %s", ws.ID, cmdStr)
}

// Workflow manages a sequence of commands to be executed
type Workflow struct {
	steps  []WorkflowStep
	nextID int
	mutex  sync.RWMutex
}

// NewWorkflow creates a new workflow
func NewWorkflow() *Workflow {
	return &Workflow{
		steps:  make([]WorkflowStep, 0),
		nextID: 1,
	}
}

// AddStep adds a step to the workflow
func (w *Workflow) AddStep(command string, args []string, description string) int {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	step := WorkflowStep{
		ID:          w.nextID,
		Command:     command,
		Args:        args,
		Description: description,
	}

	w.steps = append(w.steps, step)
	id := w.nextID
	w.nextID++

	return id
}

// GetSteps returns a copy of all workflow steps
func (w *Workflow) GetSteps() []WorkflowStep {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	result := make([]WorkflowStep, len(w.steps))
	copy(result, w.steps)
	return result
}

// Clear removes all steps from the workflow
func (w *Workflow) Clear() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.steps = w.steps[:0]
	w.nextID = 1
}

// IsEmpty returns true if the workflow has no steps
func (w *Workflow) IsEmpty() bool {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	return len(w.steps) == 0
}

// Size returns the number of steps in the workflow
func (w *Workflow) Size() int {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	return len(w.steps)
}

// CommandRouter represents an interface for routing commands
type CommandRouter interface {
	Route(args []string)
}

// WorkflowExecutor executes workflow steps sequentially using existing Route mechanism
type WorkflowExecutor struct {
	router CommandRouter
}

// NewWorkflowExecutor creates a new workflow executor
func NewWorkflowExecutor(router CommandRouter) *WorkflowExecutor {
	return &WorkflowExecutor{
		router: router,
	}
}

// Execute runs all steps in the workflow sequentially
func (we *WorkflowExecutor) Execute(workflow *Workflow) error {
	steps := workflow.GetSteps()

	if len(steps) == 0 {
		return fmt.Errorf("workflow is empty")
	}

	fmt.Printf("🚀 Starting workflow execution (%d steps)\n\n", len(steps))

	for i, step := range steps {
		fmt.Printf("📋 Step %d/%d: %s\n", i+1, len(steps), step.String())

		// Resolve placeholders in the description (which contains the original template)
		finalCmd := step.Description
		placeholders := extractPlaceholders(finalCmd)

		if len(placeholders) > 0 {
			// Interactive input for placeholders
			inputs := interactiveInputForWorkflow(placeholders)

			// Placeholder replacement
			for ph, val := range inputs {
				finalCmd = strings.ReplaceAll(finalCmd, "<"+ph+">", val)
			}

			fmt.Printf("   → Resolved to: %s\n", finalCmd)
		}

		// Parse resolved command
		parts := strings.Fields(finalCmd)
		if len(parts) == 0 {
			continue
		}

		// Execute the resolved command using existing Route mechanism
		we.router.Route(parts)

		fmt.Printf("✅ Step %d completed successfully\n", i+1)

		// Add separator between steps (except for the last one)
		if i < len(steps)-1 {
			fmt.Println("─────────────────────────────────────")
		}
	}

	fmt.Printf("\n🎉 Workflow completed successfully! (%d steps executed)\n", len(steps))
	return nil
}

// interactiveInputForWorkflow provides interactive input for placeholders during workflow execution
func interactiveInputForWorkflow(placeholders []string) map[string]string {
	inputs := make(map[string]string)

	for i, ph := range placeholders {
		if len(placeholders) > 1 {
			fmt.Printf("\n[%d/%d] ", i+1, len(placeholders))
		} else {
			fmt.Print("\n")
		}

		fmt.Printf("? %s: ", ph)

		// Simple input reading (not raw mode since we're in workflow execution)
		var value string
		_, err := fmt.Scanln(&value)
		if err != nil {
			fmt.Printf("Input error: %v\n", err)
			return nil
		}

		if value == "" {
			fmt.Printf("Operation canceled\n")
			return nil
		}

		inputs[ph] = value
		fmt.Printf("✓ %s: %s\n", ph, value)
	}

	return inputs
}
