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

// WorkflowSource identifies how a workflow was created.
type WorkflowSource int

const (
	// WorkflowSourceDynamic marks workflows created during the interactive session.
	WorkflowSourceDynamic WorkflowSource = iota
	// WorkflowSourceConfig marks workflows sourced from configuration.
	WorkflowSourceConfig
)

type managedWorkflow struct {
	data     *Workflow
	name     string
	source   WorkflowSource
	readOnly bool
}

// WorkflowSummary describes a workflow for listing/selection purposes.
type WorkflowSummary struct {
	ID        int
	StepCount int
	IsActive  bool
	Name      string
	Source    WorkflowSource
	ReadOnly  bool
}

// WorkflowManager manages multiple workflows and their lifecycle.
type WorkflowManager struct {
	mutex     sync.RWMutex
	workflows map[int]*managedWorkflow
	order     []int
	nextID    int
	activeID  int
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

	mgr.createWorkflowLocked(NewWorkflow(), WorkflowSourceDynamic, "", false, true)

	return mgr
}

func (m *WorkflowManager) createWorkflowLocked(workflow *Workflow, source WorkflowSource, name string, readOnly bool, makeActive bool) int {
	id := m.nextID
	m.nextID++

	m.workflows[id] = &managedWorkflow{
		data:     workflow,
		name:     name,
		source:   source,
		readOnly: readOnly,
	}
	m.order = append(m.order, id)
	if makeActive || m.activeID == 0 {
		m.activeID = id
	}

	return id
}

// CreateWorkflow adds a new empty workflow and makes it active.
func (m *WorkflowManager) CreateWorkflow() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.createWorkflowLocked(NewWorkflow(), WorkflowSourceDynamic, "", false, true)
}

// CreateReadOnlyWorkflow adds a workflow seeded from configuration templates.
func (m *WorkflowManager) CreateReadOnlyWorkflow(name string, templates []string) (int, error) {
	workflow, err := workflowFromTemplates(templates)
	if err != nil {
		return 0, err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.createWorkflowLocked(workflow, WorkflowSourceConfig, name, true, m.activeID == 0), nil
}

// CreateWorkflowFromTemplates creates a dynamic workflow populated with provided templates.
func (m *WorkflowManager) CreateWorkflowFromTemplates(name string, templates []string) (int, error) {
	workflow, err := workflowFromTemplates(templates)
	if err != nil {
		return 0, err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.createWorkflowLocked(workflow, WorkflowSourceDynamic, name, false, true), nil
}

// DeleteWorkflow removes a workflow by ID, returning the new active workflow ID.
func (m *WorkflowManager) DeleteWorkflow(id int) (int, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	mw, exists := m.workflows[id]
	if !exists || mw == nil || mw.readOnly {
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
		newID := m.createWorkflowLocked(NewWorkflow(), WorkflowSourceDynamic, "", false, true)
		return newID, true
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
			Source:    managed.source,
			ReadOnly:  managed.readOnly,
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

// GetActiveID returns the current active workflow ID.
func (m *WorkflowManager) GetActiveID() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.activeID
}

func (m *WorkflowManager) getManagedWorkflow(id int) (*managedWorkflow, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	mw, exists := m.workflows[id]
	return mw, exists
}

// GetWorkflow returns the workflow pointer for a given ID.
func (m *WorkflowManager) GetWorkflow(id int) (*Workflow, bool) {
	mw, exists := m.getManagedWorkflow(id)
	if !exists || mw == nil {
		return nil, false
	}
	return mw.data, true
}

// GetManagedWorkflowMetadata returns metadata for a workflow without exposing internals.
func (m *WorkflowManager) GetManagedWorkflowMetadata(id int) (WorkflowSummary, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	mw, exists := m.workflows[id]
	if !exists || mw == nil {
		return WorkflowSummary{}, false
	}
	return WorkflowSummary{
		ID:        id,
		StepCount: mw.data.Size(),
		IsActive:  id == m.activeID,
		Name:      mw.name,
		Source:    mw.source,
		ReadOnly:  mw.readOnly,
	}, true
}

// GetActiveWorkflow retrieves the currently active workflow.
func (m *WorkflowManager) GetActiveWorkflow() (*Workflow, int) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	mw, exists := m.workflows[m.activeID]
	if !exists || mw == nil {
		return nil, 0
	}
	return mw.data, m.activeID
}

// AddStep appends a step to the specified workflow.
func (m *WorkflowManager) AddStep(id int, command string, args []string, description string) (int, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	mw, exists := m.workflows[id]
	if !exists || mw == nil {
		return 0, fmt.Errorf("workflow %d not found", id)
	}
	if mw.readOnly {
		return 0, fmt.Errorf("workflow %d is read-only", id)
	}

	return mw.data.AddStep(command, args, description), nil
}

// ClearWorkflow removes all steps from the specified workflow.
func (m *WorkflowManager) ClearWorkflow(id int) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	mw, exists := m.workflows[id]
	if !exists || mw == nil || mw.readOnly {
		return false
	}
	mw.data.Clear()
	return true
}

// WorkflowCount returns the number of managed workflows.
func (m *WorkflowManager) WorkflowCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.workflows)
}

func workflowFromTemplates(templates []string) (*Workflow, error) {
	wf := NewWorkflow()
	for i, tmpl := range templates {
		trimmed := strings.TrimSpace(tmpl)
		if trimmed == "" {
			return nil, fmt.Errorf("workflow template at index %d is empty", i)
		}
		parts := strings.Fields(trimmed)
		if len(parts) == 0 {
			return nil, fmt.Errorf("workflow template at index %d is invalid", i)
		}
		wf.AddStep(parts[0], append([]string(nil), parts[1:]...), trimmed)
	}
	return wf, nil
}

func workflowToTemplates(wf *Workflow) []string {
	if wf == nil {
		return nil
	}
	steps := wf.GetSteps()
	templates := make([]string, 0, len(steps))
	for _, step := range steps {
		template := strings.TrimSpace(step.Description)
		if template == "" {
			fields := append([]string{step.Command}, step.Args...)
			template = strings.Join(fields, " ")
		}
		templates = append(templates, template)
	}
	return templates
}

func normalizeTemplates(templates []string) []string {
	if len(templates) == 0 {
		return nil
	}
	normalized := make([]string, len(templates))
	for i, template := range templates {
		normalized[i] = strings.TrimSpace(template)
	}
	return normalized
}

func templatesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if strings.TrimSpace(a[i]) != strings.TrimSpace(b[i]) {
			return false
		}
	}
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

		// Resolve placeholders in the description (which contains the original template)
		finalCmd := step.Description
		placeholders := extractPlaceholders(finalCmd)

		if len(placeholders) > 0 {
			inputs, canceled := interactiveInputForWorkflow(we.ui, placeholders)
			if canceled {
				return ErrWorkflowCanceled
			}

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
