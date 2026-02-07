package interactive

import (
	"errors"
	"fmt"
	"strings"
)

// CommandRouter represents an interface for routing commands
type CommandRouter interface {
	Route(args []string)
}

// WorkflowExecutor executes workflow steps sequentially using existing Route mechanism
type WorkflowExecutor struct {
	router CommandRouter
	ui     *UI
}

// ErrWorkflowCanceled indicates the workflow was aborted by the user via soft cancel.
var ErrWorkflowCanceled = errors.New("workflow canceled")

// NewWorkflowExecutor creates a new workflow executor
func NewWorkflowExecutor(router CommandRouter, ui *UI) *WorkflowExecutor {
	return &WorkflowExecutor{
		router: router,
		ui:     ui,
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

		// Resolve placeholders in each argument individually to preserve multiword values
		resolvedArgs, canceled := resolveStepPlaceholders(we.ui, step)
		if canceled {
			return ErrWorkflowCanceled
		}

		// Build parts array: command + resolved args
		parts := append([]string{step.Command}, resolvedArgs...)

		if parts[0] == "" {
			continue
		}

		// Show resolved command
		fmt.Printf("   → Resolved to: %s\n", strings.Join(parts, " "))

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
