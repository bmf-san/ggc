package cmd

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/term"
)

// mockTerminal mocks terminal operations
type mockTerminal struct {
	makeRawCalled  bool
	restoreCalled  bool
	shouldFailRaw  bool
	shouldFailRest bool
}

func (m *mockTerminal) makeRaw(_ int) (*term.State, error) {
	m.makeRawCalled = true
	if m.shouldFailRaw {
		return nil, fmt.Errorf("mock makeRaw error")
	}
	return &term.State{}, nil
}

func (m *mockTerminal) restore(_ int, _ *term.State) error {
	m.restoreCalled = true
	if m.shouldFailRest {
		return fmt.Errorf("mock restore error")
	}
	return nil
}

// testUI is a test structure for UI
type testUI struct {
	UI
	inputBytes []byte
}

func (ui *testUI) Run() []string {
	// Simulate standard input
	ui.stdin = bytes.NewReader(ui.inputBytes)
	ui.stdout = &bytes.Buffer{}
	ui.stderr = &bytes.Buffer{}

	// Update renderer writer to use test buffer
	ui.renderer.writer = ui.stdout

	return ui.UI.Run()
}

func TestUI_Run(t *testing.T) {
	tests := []struct {
		name         string
		input        []byte
		expectedArgs []string
		expectNil    bool
	}{
		{
			name:      "press Enter with empty input",
			input:     []byte{13}, // Enter key
			expectNil: true,
		},
		{
			name:         "type 'help' and press Enter",
			input:        []byte{'h', 'e', 'l', 'p', 13}, // 'h','e','l','p' + Enter
			expectedArgs: []string{"ggc", "help"},
			expectNil:    false,
		},
		{
			name:      "type non-existent command",
			input:     []byte{'x', 'y', 'z', 13}, // 'x','y','z' + Enter
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test UI with new design
			mockTerm := &mockTerminal{}
			colors := NewANSIColors()
			renderer := &Renderer{
				writer: &bytes.Buffer{},
				colors: colors,
			}
			state := &UIState{
				selected:  0,
				input:     "",
				cursorPos: 0,
				filtered:  []CommandInfo{},
			}

			ui := &testUI{
				UI: UI{
					term:     mockTerm,
					renderer: renderer,
					state:    state,
					colors:   colors,
				},
				inputBytes: tt.input,
			}

			// Set up handler after UI is created
			ui.handler = &KeyHandler{ui: &ui.UI}

			// Execute test
			result := ui.Run()

			// Verify results
			if tt.expectNil && result != nil {
				t.Errorf("expected: nil, got: %v", result)
			}

			if !tt.expectNil {
				if result == nil {
					t.Error("expected: not nil, got: nil")
					return
				}
				if len(result) != len(tt.expectedArgs) {
					t.Errorf("expected length: %d, got length: %d", len(tt.expectedArgs), len(result))
					return
				}
				for i, arg := range tt.expectedArgs {
					if result[i] != arg {
						t.Errorf("expected[%d]: %s, got[%d]: %s", i, arg, i, result[i])
					}
				}
			}
		})
	}
}

func TestExtractPlaceholders(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "no placeholders",
			input: "simple command",
			want:  []string{},
		},
		{
			name:  "single placeholder",
			input: "add <file>",
			want:  []string{"file"},
		},
		{
			name:  "multiple placeholders",
			input: "remote add <name> <url>",
			want:  []string{"name", "url"},
		},
		{
			name:  "empty placeholder",
			input: "command <>",
			want:  []string{""},
		},
		{
			name:  "invalid format",
			input: "command <incomplete",
			want:  []string{},
		},
		{
			name:  "nested placeholder",
			input: "command <<nested>>",
			want:  []string{"nested"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractPlaceholders(tt.input)
			// Compare without distinguishing between nil and empty slice
			if len(tt.want) == 0 && len(got) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractPlaceholders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandDescriptions(t *testing.T) {
	// Ensure that all commands have a description (catch in tests)
	for _, cmd := range commands {
		description := cmd.Description
		if description == "" {
			t.Errorf("Command '%s' has no description", cmd.Command)
		}
	}
}

func TestCommandDescriptionsContent(t *testing.T) {
	// Test description content updated + showing correctly
	expectedDescriptions := map[string]string{
		"add <file>":       "Add a specific file to the index",
		"status":           "Show working tree status",
		"commit <message>": "Create commit with a message",
		"quit":             "Exit interactive mode",
	}

	for cmdStr, expectedDesc := range expectedDescriptions {
		found := false
		for _, cmd := range commands {
			if cmd.Command == cmdStr {
				if cmd.Description != expectedDesc {
					t.Errorf("Description for '%s' is '%s', expected '%s'", cmdStr, cmd.Description, expectedDesc)
				}
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Command '%s' not found in commandList", cmdStr)
		}
	}
}

// Test UIState functionality
func TestUIState_UpdateFiltered(t *testing.T) {
	state := &UIState{
		selected:  0,
		input:     "add",
		cursorPos: 3,
		filtered:  []CommandInfo{},
	}

	state.UpdateFiltered()

	if len(state.filtered) == 0 {
		t.Error("Expected filtered commands for 'add' input")
	}

	// Check that all filtered commands contain 'add'
	for _, cmd := range state.filtered {
		if !strings.Contains(cmd.Command, "add") {
			t.Errorf("Filtered command '%s' does not contain 'add'", cmd.Command)
		}
	}
}

func TestUIState_MoveUp(t *testing.T) {
	state := &UIState{
		selected:  2,
		input:     "",
		cursorPos: 0,
		filtered: []CommandInfo{
			{"cmd1", "desc1"},
			{"cmd2", "desc2"},
			{"cmd3", "desc3"},
		},
	}

	state.MoveUp()
	if state.selected != 1 {
		t.Errorf("Expected selected to be 1, got %d", state.selected)
	}

	state.selected = 0
	state.MoveUp()
	if state.selected != 0 {
		t.Errorf("Expected selected to stay 0 when at top, got %d", state.selected)
	}
}

func TestUIState_MoveDown(t *testing.T) {
	state := &UIState{
		selected:  0,
		input:     "",
		cursorPos: 0,
		filtered: []CommandInfo{
			{"cmd1", "desc1"},
			{"cmd2", "desc2"},
			{"cmd3", "desc3"},
		},
	}

	state.MoveDown()
	if state.selected != 1 {
		t.Errorf("Expected selected to be 1, got %d", state.selected)
	}

	state.selected = 2
	state.MoveDown()
	if state.selected != 2 {
		t.Errorf("Expected selected to stay 2 when at bottom, got %d", state.selected)
	}
}

func TestUIState_AddChar(t *testing.T) {
	state := &UIState{
		selected:  0,
		input:     "",
		cursorPos: 0,
		filtered:  []CommandInfo{},
	}

	state.AddChar('a')
	if state.input != "a" {
		t.Errorf("Expected input to be 'a', got '%s'", state.input)
	}

	state.AddChar('d')
	if state.input != "ad" {
		t.Errorf("Expected input to be 'ad', got '%s'", state.input)
	}
}

func TestUIState_RemoveChar(t *testing.T) {
	state := &UIState{
		selected:  0,
		input:     "test",
		cursorPos: 4,
		filtered:  []CommandInfo{},
	}

	state.RemoveChar()
	if state.input != "tes" {
		t.Errorf("Expected input to be 'tes', got '%s'", state.input)
	}

	// Test removing from empty string
	state.input = ""
	state.RemoveChar()
	if state.input != "" {
		t.Errorf("Expected input to remain empty, got '%s'", state.input)
	}
}

func TestUIState_ClearInput(t *testing.T) {
	state := &UIState{
		selected:  0,
		input:     "test input",
		cursorPos: 5,
		filtered:  []CommandInfo{},
	}

	state.ClearInput()
	if state.input != "" {
		t.Errorf("Expected input to be empty after clear, got '%s'", state.input)
	}
	if state.cursorPos != 0 {
		t.Errorf("Expected cursor position to be 0 after clear, got %d", state.cursorPos)
	}

	// Test clearing empty input
	state.ClearInput()
	if state.input != "" {
		t.Errorf("Expected input to remain empty, got '%s'", state.input)
	}
}

func TestUIState_DeleteWord(t *testing.T) {
	state := &UIState{
		selected:  0,
		input:     "hello world test",
		cursorPos: 16, // At end
		filtered:  []CommandInfo{},
	}

	state.DeleteWord()
	if state.input != "hello world " {
		t.Errorf("Expected 'hello world ', got '%s'", state.input)
	}
	if state.cursorPos != 12 {
		t.Errorf("Expected cursor at 12, got %d", state.cursorPos)
	}

	// Delete another word
	state.DeleteWord()
	if state.input != "hello " {
		t.Errorf("Expected 'hello ', got '%s'", state.input)
	}
}

func TestUIState_DeleteToEnd(t *testing.T) {
	state := &UIState{
		selected:  0,
		input:     "hello world",
		cursorPos: 5, // After "hello"
		filtered:  []CommandInfo{},
	}

	state.DeleteToEnd()
	if state.input != "hello" {
		t.Errorf("Expected 'hello', got '%s'", state.input)
	}
}

func TestUIState_MoveToBeginning(t *testing.T) {
	state := &UIState{
		selected:  0,
		input:     "test",
		cursorPos: 4,
		filtered:  []CommandInfo{},
	}

	state.MoveToBeginning()
	if state.cursorPos != 0 {
		t.Errorf("Expected cursor at 0, got %d", state.cursorPos)
	}
}

func TestUIState_MoveToEnd(t *testing.T) {
	state := &UIState{
		selected:  0,
		input:     "test",
		cursorPos: 0,
		filtered:  []CommandInfo{},
	}

	state.MoveToEnd()
	if state.cursorPos != 4 {
		t.Errorf("Expected cursor at 4, got %d", state.cursorPos)
	}
}

func TestUIState_GetSelectedCommand(t *testing.T) {
	state := &UIState{
		selected:  1,
		input:     "",
		cursorPos: 0,
		filtered: []CommandInfo{
			{"cmd1", "desc1"},
			{"cmd2", "desc2"},
		},
	}

	cmd := state.GetSelectedCommand()
	if cmd == nil {
		t.Fatal("Expected non-nil command")
	}
	if cmd.Command != "cmd2" {
		t.Errorf("Expected 'cmd2', got '%s'", cmd.Command)
	}

	// Test out of bounds
	state.selected = 10
	cmd = state.GetSelectedCommand()
	if cmd != nil {
		t.Error("Expected nil for out of bounds selection")
	}
}

// Test Renderer functionality
func TestRenderer_UpdateSize(t *testing.T) {
	var buf bytes.Buffer
	colors := NewANSIColors()
	renderer := &Renderer{
		writer: &buf,
		colors: colors,
	}

	renderer.updateSize()

	// Should have default values when not a file
	if renderer.width != 80 || renderer.height != 24 {
		t.Errorf("Expected default size 80x24, got %dx%d", renderer.width, renderer.height)
	}
}

func TestRenderer_CalculateMaxCommandLength(t *testing.T) {
	renderer := &Renderer{}
	commands := []CommandInfo{
		{"short", "desc"},
		{"very long command", "desc"},
		{"medium", "desc"},
	}

	maxLen := renderer.calculateMaxCommandLength(commands)
	expected := len("very long command")
	if maxLen != expected {
		t.Errorf("Expected max length %d, got %d", expected, maxLen)
	}
}

// Test Renderer keybind help functionality
func TestRenderer_KeybindHelp(t *testing.T) {
	var buf bytes.Buffer
	colors := NewANSIColors()
	renderer := &Renderer{
		writer: &buf,
		colors: colors,
		width:  80,
		height: 24,
	}

	state := &UIState{
		selected:  0,
		input:     "nonexistent",
		cursorPos: 11,
		filtered:  []CommandInfo{}, // Empty to trigger keybind help
	}

	ui := &UI{
		stdin:    strings.NewReader(""),
		stdout:   &buf,
		stderr:   &bytes.Buffer{},
		term:     &mockTerminal{},
		renderer: renderer,
		state:    state,
		colors:   colors,
	}

	renderer.Render(ui, state)
	output := buf.String()

	// Check that keybind help is displayed
	expectedKeybinds := []string{
		"Available keybinds:",
		"Ctrl+u",
		"Ctrl+w",
		"Ctrl+k",
		"Ctrl+a",
		"Ctrl+e",
		"Backspace",
	}

	for _, keybind := range expectedKeybinds {
		if !strings.Contains(output, keybind) {
			t.Errorf("Expected keybind help to contain '%s', but it was not found", keybind)
		}
	}
}

// Test Renderer empty state display
func TestRenderer_EmptyState(t *testing.T) {
	var buf bytes.Buffer
	colors := NewANSIColors()
	renderer := &Renderer{
		writer: &buf,
		colors: colors,
		width:  80,
		height: 24,
	}

	state := &UIState{
		selected:  0,
		input:     "", // Empty input
		cursorPos: 0,
		filtered:  []CommandInfo{},
	}

	ui := &UI{
		stdin:    strings.NewReader(""),
		stdout:   &buf,
		stderr:   &bytes.Buffer{},
		term:     &mockTerminal{},
		renderer: renderer,
		state:    state,
		colors:   colors,
	}

	renderer.Render(ui, state)
	output := buf.String()

	// Check that simple message is displayed
	if !strings.Contains(output, "Start typing to search commands...") {
		t.Error("Expected empty state to show simple search message")
	}

	// Check that popular commands are NOT displayed
	unwantedTexts := []string{
		"Popular commands to get started:",
		"status",
		"add .",
		"commit",
		"push",
		"pull",
	}

	for _, unwanted := range unwantedTexts {
		if strings.Contains(output, unwanted) {
			t.Errorf("Expected empty state to NOT contain '%s', but it was found", unwanted)
		}
	}
}

// Test Git status functionality
func TestGetGitStatus(t *testing.T) {
	// This test may fail in environments without git or outside git repos
	status := getGitStatus()
	
	// If we're in a git repository, status should not be nil
	if status != nil {
		// Branch name should not be empty
		if status.Branch == "" {
			t.Error("Expected branch name to be non-empty when in git repository")
		}
		
		// Values should be non-negative
		if status.Modified < 0 || status.Staged < 0 || status.Ahead < 0 || status.Behind < 0 {
			t.Error("Git status values should be non-negative")
		}
	}
	// If status is nil, we're probably not in a git repository, which is fine for tests
}

// Test Git status rendering
func TestRenderer_RenderGitStatus(t *testing.T) {
	var buf bytes.Buffer
	colors := NewANSIColors()
	renderer := &Renderer{
		writer: &buf,
		colors: colors,
		width:  80,
		height: 24,
	}

	// Test with mock git status
	mockStatus := &GitStatus{
		Branch:     "main",
		Modified:   2,
		Staged:     1,
		Ahead:      3,
		Behind:     1,
		HasChanges: true,
	}

	ui := &UI{
		stdin:     strings.NewReader(""),
		stdout:    &buf,
		stderr:    &bytes.Buffer{},
		term:      &mockTerminal{},
		renderer:  renderer,
		colors:    colors,
		gitStatus: mockStatus,
	}

	renderer.renderGitStatus(ui, mockStatus)
	output := buf.String()

	// Check that git status elements are displayed
	expectedElements := []string{
		"📍", // Branch icon
		"main", // Branch name
		"📝", // Changes icon
		"2 modified", // Modified files
		"1 staged", // Staged files
		"↑3", // Ahead count
		"↓1", // Behind count
	}

	for _, element := range expectedElements {
		if !strings.Contains(output, element) {
			t.Errorf("Expected git status to contain '%s', but it was not found", element)
		}
	}
}

// Test Renderer header/footer keybind display
func TestRenderer_KeybindDisplay(t *testing.T) {
	var buf bytes.Buffer
	colors := NewANSIColors()
	renderer := &Renderer{
		writer: &buf,
		colors: colors,
		width:  80,
		height: 24,
	}

	state := &UIState{
		selected:  0,
		input:     "test",
		cursorPos: 4,
		filtered: []CommandInfo{
			{"test command", "test description"},
		},
	}

	ui := &UI{
		stdin:    strings.NewReader(""),
		stdout:   &buf,
		stderr:   &bytes.Buffer{},
		term:     &mockTerminal{},
		renderer: renderer,
		state:    state,
		colors:   colors,
	}

	renderer.Render(ui, state)
	output := buf.String()

	// Check that lowercase keybind notation is used
	expectedKeybinds := []string{
		"Ctrl+n/p",
		"Ctrl+a/e",
		"Ctrl+u/w/k",
		"Ctrl+c",
	}

	for _, keybind := range expectedKeybinds {
		if !strings.Contains(output, keybind) {
			t.Errorf("Expected output to contain lowercase keybind '%s', but it was not found", keybind)
		}
	}

	// Check that uppercase versions are NOT used
	uppercaseKeybinds := []string{
		"Ctrl+N/P",
		"Ctrl+A/E",
		"Ctrl+U/W/K",
		"Ctrl+C",
	}

	for _, keybind := range uppercaseKeybinds {
		if strings.Contains(output, keybind) {
			t.Errorf("Expected output to NOT contain uppercase keybind '%s', but it was found", keybind)
		}
	}
}

// Test KeyHandler functionality
func TestKeyHandler_HandleKey(t *testing.T) {
	var stdout, stderr bytes.Buffer
	colors := NewANSIColors()
	renderer := &Renderer{
		writer: &stdout,
		colors: colors,
	}
	state := &UIState{
		selected:  0,
		input:     "",
		cursorPos: 0,
		filtered:  []CommandInfo{},
	}

	ui := &UI{
		stdin:    strings.NewReader(""),
		stdout:   &stdout,
		stderr:   &stderr,
		term:     &mockTerminal{},
		renderer: renderer,
		state:    state,
		colors:   colors,
	}

	handler := &KeyHandler{ui: ui}

	// Test printable character
	shouldContinue, result := handler.HandleKey('a', nil)
	if !shouldContinue {
		t.Error("Expected to continue after printable character")
	}
	if result != nil {
		t.Error("Expected nil result for printable character")
	}
	if ui.state.input != "a" {
		t.Errorf("Expected input 'a', got '%s'", ui.state.input)
	}

	// Test backspace
	shouldContinue, _ = handler.HandleKey(127, nil)
	if !shouldContinue {
		t.Error("Expected to continue after backspace")
	}
	if ui.state.input != "" {
		t.Errorf("Expected empty input after backspace, got '%s'", ui.state.input)
	}
}
