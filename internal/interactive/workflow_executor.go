package interactive

import (
	"errors"
	"fmt"
	"strings"
)

// CommandRouter represents an interface for routing commands
type CommandRouter interface {
	Route(args []string) error
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

// uiWrite writes to the UI stdout when the UI is available; otherwise falls back to fmt.Printf.
// This allows WorkflowExecutor to work correctly in tests where the UI may be nil.
func (we *WorkflowExecutor) uiWrite(format string, a ...interface{}) {
	if we.ui != nil {
		we.ui.write(format, a...)
		return
	}
	_, _ = fmt.Printf(format, a...)
}

// Execute runs all steps in the workflow sequentially
func (we *WorkflowExecutor) Execute(workflow *Workflow) error {
	steps := workflow.GetSteps()

	if len(steps) == 0 {
		return fmt.Errorf("workflow is empty")
	}

	we.uiWrite("🚀 Starting workflow execution (%d steps)\n\n", len(steps))

	for i, step := range steps {
		we.uiWrite("📋 Step %d/%d: %s\n", i+1, len(steps), step.String())

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
		we.uiWrite("   → Resolved to: %s\n", strings.Join(parts, " "))

		// Execute the resolved command and propagate any routing error
		if err := we.router.Route(parts); err != nil {
			return fmt.Errorf("step %d/%d failed: %w", i+1, len(steps), err)
		}

		we.uiWrite("✅ Step %d completed successfully\n", i+1)

		// Add separator between steps (except for the last one)
		if i < len(steps)-1 {
			we.uiWrite("─────────────────────────────────────\n")
		}
	}

	we.uiWrite("\n🎉 Workflow completed successfully! (%d steps executed)\n", len(steps))
	return nil
}
