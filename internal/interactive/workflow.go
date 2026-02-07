package interactive

import (
	"sync"
)

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
