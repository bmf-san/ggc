package interactive

import (
	"fmt"
	"sync"
)

// WorkflowSummary describes a workflow for listing/selection purposes.
type WorkflowSummary struct {
	ID        int
	StepCount int
	IsActive  bool
	Name      string
}

type managedWorkflow struct {
	data *Workflow
	name string
}

// WorkflowManager manages multiple workflows and their lifecycle.
type WorkflowManager struct {
	mutex     sync.RWMutex
	workflows map[int]*managedWorkflow
	order     []int
	activeID  int
	nextID    int
}

// NewWorkflowManager constructs a manager with an initial empty workflow.
func NewWorkflowManager() *WorkflowManager {
	mgr := &WorkflowManager{
		workflows: make(map[int]*managedWorkflow),
		order:     make([]int, 0, 4),
		nextID:    1,
	}

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	mgr.createWorkflowLocked(NewWorkflow(), "")

	return mgr
}

func (m *WorkflowManager) createWorkflowLocked(workflow *Workflow, name string) int {
	id := m.nextID
	m.nextID++

	m.workflows[id] = &managedWorkflow{
		data: workflow,
		name: name,
	}
	m.order = append(m.order, id)
	m.activeID = id

	return id
}

// CreateWorkflow adds a new empty workflow and makes it active.
func (m *WorkflowManager) CreateWorkflow(name string) int {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.createWorkflowLocked(NewWorkflow(), name)
}

// DeleteWorkflow removes a workflow by ID, returning the new active workflow ID.
func (m *WorkflowManager) DeleteWorkflow(id int) (int, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.workflows[id]; !exists {
		return m.activeID, false
	}

	delete(m.workflows, id)

	removedIndex := -1
	for i, existing := range m.order {
		if existing == id {
			removedIndex = i
			m.order = append(m.order[:i], m.order[i+1:]...)
			break
		}
	}

	if len(m.workflows) == 0 {
		m.activeID = 0
		return 0, true
	}

	if m.activeID == id {
		if removedIndex >= len(m.order) {
			removedIndex = len(m.order) - 1
		}
		if removedIndex < 0 {
			removedIndex = 0
		}
		m.activeID = m.order[removedIndex]
	}

	return m.activeID, true
}

// ListWorkflows returns ordered summaries of all workflows.
func (m *WorkflowManager) ListWorkflows() []WorkflowSummary {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	summaries := make([]WorkflowSummary, 0, len(m.order))
	for _, id := range m.order {
		managed := m.workflows[id]
		if managed == nil {
			continue
		}
		stepCount := managed.data.Size()
		summaries = append(summaries, WorkflowSummary{
			ID:        id,
			StepCount: stepCount,
			IsActive:  id == m.activeID,
			Name:      managed.name,
		})
	}
	return summaries
}

// SetActive designates the workflow with the provided ID as active.
func (m *WorkflowManager) SetActive(id int) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.workflows[id]; !exists {
		return false
	}
	m.activeID = id
	return true
}

// CycleActive moves the active workflow by delta in the ordered list.
func (m *WorkflowManager) CycleActive(delta int) int {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.order) == 0 {
		m.activeID = 0
		return 0
	}

	idx := 0
	for i, id := range m.order {
		if id == m.activeID {
			idx = i
			break
		}
	}

	idx = (idx + delta) % len(m.order)
	if idx < 0 {
		idx += len(m.order)
	}
	m.activeID = m.order[idx]
	return m.activeID
}

// GetActiveID returns the current active workflow ID.
func (m *WorkflowManager) GetActiveID() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.activeID
}

// GetWorkflow returns the workflow and whether it exists.
func (m *WorkflowManager) GetWorkflow(id int) (*Workflow, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	w, exists := m.workflows[id]
	if !exists || w == nil {
		return nil, false
	}
	return w.data, true
}

// AddStep adds a step to the specified workflow.
func (m *WorkflowManager) AddStep(id int, command string, args []string, description string) (int, bool) {
	m.mutex.RLock()
	w, exists := m.workflows[id]
	m.mutex.RUnlock()
	if !exists || w == nil {
		return 0, false
	}
	return w.data.AddStep(command, args, description), true
}

// CloneWorkflow duplicates an existing workflow into a new dynamic workflow.
func (m *WorkflowManager) CloneWorkflow(id int, name string) (int, bool) {
	m.mutex.RLock()
	src, exists := m.workflows[id]
	if !exists || src == nil {
		m.mutex.RUnlock()
		return 0, false
	}
	// Take a snapshot of the steps while holding the read lock so that
	// we only clone workflows that are still managed.
	steps := src.data.GetSteps()
	m.mutex.RUnlock()

	clone := NewWorkflow()
	for _, step := range steps {
		clone.AddStep(step.Command, append([]string(nil), step.Args...), step.Description)
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	if name == "" {
		name = fmt.Sprintf("Workflow #%d (copy)", id)
	}
	newID := m.createWorkflowLocked(clone, name)
	return newID, true
}

// ClearWorkflow removes all steps from the specified workflow.
func (m *WorkflowManager) ClearWorkflow(id int) bool {
	m.mutex.RLock()
	w, exists := m.workflows[id]
	m.mutex.RUnlock()
	if !exists || w == nil {
		return false
	}
	w.data.Clear()
	return true
}
