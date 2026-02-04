package interactive

import (
	"bufio"
	"errors"
	"fmt"
	"os"
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
	m.mutex.RUnlock()
	if !exists || src == nil {
		return 0, false
	}

	clone := NewWorkflow()
	for _, step := range src.data.GetSteps() {
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

		if len(parts) == 0 || parts[0] == "" {
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

// deriveArgsFromDescription extracts arguments from a description string.
// Returns the args portion (everything after the command name).
func deriveArgsFromDescription(description string) []string {
	parts := strings.Fields(description)
	if len(parts) > 1 {
		return parts[1:]
	}
	return nil
}

// collectPlaceholders extracts unique placeholders from a list of arguments.
func collectPlaceholders(args []string) []string {
	var placeholders []string
	seen := make(map[string]bool)
	for _, arg := range args {
		for _, ph := range extractPlaceholders(arg) {
			if !seen[ph] {
				seen[ph] = true
				placeholders = append(placeholders, ph)
			}
		}
	}
	return placeholders
}

// replacePlaceholdersInArgs replaces placeholders in each argument with their values.
func replacePlaceholdersInArgs(args []string, inputs map[string]string) []string {
	resolvedArgs := make([]string, len(args))
	for i, arg := range args {
		resolved := arg
		for ph, val := range inputs {
			resolved = strings.ReplaceAll(resolved, "<"+ph+">", val)
		}
		resolvedArgs[i] = resolved
	}
	return resolvedArgs
}

// resolveStepPlaceholders resolves placeholders in a workflow step's arguments.
// Each argument is processed individually, preserving multiword placeholder values as single arguments.
func resolveStepPlaceholders(ui *UI, step WorkflowStep) ([]string, bool) {
	// If Args is empty, derive from Description
	args := step.Args
	if len(args) == 0 {
		args = deriveArgsFromDescription(step.Description)
	}

	// Extract unique placeholders from all args
	placeholders := collectPlaceholders(args)
	if len(placeholders) == 0 {
		return args, false
	}

	// Get user input for each placeholder
	inputs, canceled := interactiveInputForWorkflow(ui, placeholders)
	if canceled {
		return nil, true
	}

	return replacePlaceholdersInArgs(args, inputs), false
}

// interactiveInputForWorkflow provides interactive input for placeholders during workflow execution
func interactiveInputForWorkflow(ui *UI, placeholders []string) (map[string]string, bool) {
	if ui != nil && ui.handler != nil {
		return interactiveInputForWorkflowUI(ui, placeholders)
	}
	scanner := bufio.NewScanner(os.Stdin)
	return interactiveInputForWorkflowScanner(scanner, placeholders)
}

func interactiveInputForWorkflowUI(ui *UI, placeholders []string) (map[string]string, bool) {
	inputs := make(map[string]string)
	for i, ph := range placeholders {
		ui.write("\n")
		if len(placeholders) > 1 {
			ui.write("%s[%d/%d]%s ",
				ui.colors.BrightBlue+ui.colors.Bold,
				i+1, len(placeholders),
				ui.colors.Reset)
		}
		ui.write("%s? %s%s%s: ",
			ui.colors.BrightGreen,
			ui.colors.BrightWhite+ui.colors.Bold,
			ph,
			ui.colors.Reset)

		value, canceled := ui.readPlaceholderInput()
		if canceled {
			return nil, true
		}

		inputs[ph] = value
		ui.write("%s✓ %s%s: %s%s%s\n",
			ui.colors.BrightGreen,
			ui.colors.BrightBlue,
			ph,
			ui.colors.BrightYellow+ui.colors.Bold,
			value,
			ui.colors.Reset)
	}
	return inputs, false
}

func interactiveInputForWorkflowScanner(scanner *bufio.Scanner, placeholders []string) (map[string]string, bool) {
	inputs := make(map[string]string)
	for i, ph := range placeholders {
		if len(placeholders) > 1 {
			fmt.Printf("\n[%d/%d] ", i+1, len(placeholders))
		} else {
			fmt.Print("\n")
		}

		fmt.Printf("? %s: ", ph)

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Printf("Input error: %v\n", err)
			}
			return nil, true
		}
		value := strings.TrimSpace(scanner.Text())

		if value == "" {
			fmt.Printf("Operation canceled\n")
			return nil, true
		}

		inputs[ph] = value
		fmt.Printf("✓ %s: %s\n", ph, value)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Scanner error: %v\n", err)
		return nil, true
	}

	return inputs, false
}
