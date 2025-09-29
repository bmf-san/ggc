package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"

	"golang.org/x/term"

	"github.com/bmf-san/ggc/v6/internal/termio"
	"github.com/bmf-san/ggc/v6/internal/testutil"
)

// mockTerminal mocks terminal operations
type mockTerminal struct {
	makeRawCalled  bool
	restoreCalled  bool
	shouldFailRaw  bool
	shouldFailRest bool
}

func (m *mockTerminal) MakeRaw(_ int) (*term.State, error) {
	m.makeRawCalled = true
	if m.shouldFailRaw {
		return nil, fmt.Errorf("mock makeRaw error")
	}
	return &term.State{}, nil
}

func (m *mockTerminal) Restore(_ int, _ *term.State) error {
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
			mockGitClient := testutil.NewMockGitClient()
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
					term:      mockTerm,
					renderer:  renderer,
					state:     state,
					colors:    colors,
					gitClient: mockGitClient,
					gitStatus: getGitStatus(mockGitClient),
					workflow:  NewWorkflow(),
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
			if !slices.Equal(got, tt.want) {
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

	// Check that all filtered commands fuzzy match 'add'
	for _, cmd := range state.filtered {
		if !fuzzyMatch(strings.ToLower(cmd.Command), "add") {
			t.Errorf("Filtered command '%s' does not fuzzy match 'add'", cmd.Command)
		}
	}
}

// Test fuzzy matching behavior
func TestUIState_UpdateFiltered_FuzzyMatching(t *testing.T) {
	state := &UIState{
		selected:  0,
		input:     "commit",
		cursorPos: 6,
		filtered:  []CommandInfo{},
	}

	state.UpdateFiltered()

	// Should match commands containing "commit"
	expectedMatches := []string{
		"commit <message>",
		"commit allow empty",
		"commit amend",
		"commit amend no-edit",
	}

	// Check that expected commands are found
	for _, expected := range expectedMatches {
		found := false
		for _, filtered := range state.filtered {
			if filtered.Command == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find command containing 'commit': %s", expected)
		}
	}

	// Check that all filtered commands fuzzy match "commit"
	for _, cmd := range state.filtered {
		if !fuzzyMatch(strings.ToLower(cmd.Command), "commit") {
			t.Errorf("Filtered command '%s' should fuzzy match 'commit'", cmd.Command)
		}
	}

	// Note: With substring matching, we now allow commands that contain "commit"
	// even if they don't start with it, so this test section is no longer needed
}

// Test fuzzy matching with non-consecutive characters
func TestUIState_UpdateFiltered_FuzzyNonConsecutive(t *testing.T) {
	state := &UIState{
		selected:  0,
		input:     "bd", // Should match "branch delete"
		cursorPos: 2,
		filtered:  []CommandInfo{},
	}

	state.UpdateFiltered()

	// Should find "branch delete" with fuzzy matching "bd" -> "Branch Delete"
	found := false
	for _, cmd := range state.filtered {
		if strings.Contains(cmd.Command, "branch") && strings.Contains(cmd.Command, "delete") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find 'branch delete' with fuzzy pattern 'bd'")
	}

	// All filtered commands should fuzzy match "bd"
	for _, cmd := range state.filtered {
		if !fuzzyMatch(strings.ToLower(cmd.Command), "bd") {
			t.Errorf("Filtered command '%s' should fuzzy match 'bd'", cmd.Command)
		}
	}
}

// Test fuzzy matching algorithm directly
func TestFuzzyMatch(t *testing.T) {
	testCases := []struct {
		text     string
		pattern  string
		expected bool
	}{
		{"branch delete", "bd", true},
		{"branch delete", "brdel", true},
		{"commit amend", "ca", true},
		{"commit amend", "cmtam", true},
		{"add interactive", "ai", true},
		{"add interactive", "addi", true},
		{"status short", "ss", true},
		{"branch delete", "db", false}, // wrong order
		{"commit", "xyz", false},       // no match
		{"", "", true},                 // empty pattern
		{"test", "", true},             // empty pattern matches anything
	}

	for _, tc := range testCases {
		result := fuzzyMatch(strings.ToLower(tc.text), strings.ToLower(tc.pattern))
		if result != tc.expected {
			t.Errorf("fuzzyMatch(%q, %q) = %v, expected %v", tc.text, tc.pattern, result, tc.expected)
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

func TestUIState_AddRune_ASCII(t *testing.T) {
	state := &UIState{
		selected:  0,
		input:     "",
		cursorPos: 0,
		filtered:  []CommandInfo{},
	}

	state.AddRune('a')
	if state.input != "a" {
		t.Errorf("Expected input to be 'a', got '%s'", state.input)
	}

	state.AddRune('d')
	if state.input != "ad" {
		t.Errorf("Expected input to be 'ad', got '%s'", state.input)
	}
}

func TestUIState_AddRune_MultibyteCharacters(t *testing.T) {
	state := &UIState{
		selected:  0,
		input:     "",
		cursorPos: 0,
		filtered:  []CommandInfo{},
	}

	// Test Japanese hiragana
	state.AddRune('こ')
	if state.input != "こ" {
		t.Errorf("Expected input to be 'こ', got '%s'", state.input)
	}
	if state.cursorPos != 1 {
		t.Errorf("Expected cursor position to be 1, got %d", state.cursorPos)
	}

	// Test Japanese kanji
	state.AddRune('漢')
	if state.input != "こ漢" {
		t.Errorf("Expected input to be 'こ漢', got '%s'", state.input)
	}
	if state.cursorPos != 2 {
		t.Errorf("Expected cursor position to be 2, got %d", state.cursorPos)
	}

	// Test Chinese characters
	state.input = ""
	state.cursorPos = 0
	state.AddRune('中')
	state.AddRune('文')
	if state.input != "中文" {
		t.Errorf("Expected input to be '中文', got '%s'", state.input)
	}

	// Test emoji
	state.input = ""
	state.cursorPos = 0
	state.AddRune('🎉')
	state.AddRune('✨')
	if state.input != "🎉✨" {
		t.Errorf("Expected input to be '🎉✨', got '%s'", state.input)
	}
	if state.cursorPos != 2 {
		t.Errorf("Expected cursor position to be 2, got %d", state.cursorPos)
	}
}

func TestUIState_RemoveChar_Multibyte(t *testing.T) {
	state := &UIState{
		selected:  0,
		input:     "こんにちは",
		cursorPos: 5, // At the end
		filtered:  []CommandInfo{},
	}

	// Remove last character 'は'
	state.RemoveChar()
	if state.input != "こんにち" {
		t.Errorf("Expected input to be 'こんにち', got '%s'", state.input)
	}
	if state.cursorPos != 4 {
		t.Errorf("Expected cursor position to be 4, got %d", state.cursorPos)
	}

	// Remove middle character
	state.cursorPos = 2 // Position after 'ん'
	state.RemoveChar()
	if state.input != "こにち" {
		t.Errorf("Expected input to be 'こにち', got '%s'", state.input)
	}
	if state.cursorPos != 1 {
		t.Errorf("Expected cursor position to be 1, got %d", state.cursorPos)
	}
}

func TestUIState_UpdateFiltered_MultibyteFuzzy(t *testing.T) {
	state := &UIState{
		input: "こ", // Japanese hiragana input
	}

	// This test verifies that multibyte input doesn't crash the UpdateFiltered method
	// and that the fuzzy matching algorithm can handle UTF-8 characters correctly
	state.UpdateFiltered()

	// The test passes if no panic occurs and filtered slice is initialized
	if state.filtered == nil {
		t.Error("Expected filtered slice to be initialized")
	}

	// Test with emoji input
	state.input = "🎉"
	state.UpdateFiltered()

	if state.filtered == nil {
		t.Error("Expected filtered slice to be initialized with emoji input")
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
	if cmd == nil || cmd.Command != "cmd2" {
		if cmd == nil {
			t.Fatal("Expected non-nil command")
		} else {
			t.Errorf("Expected 'cmd2', got '%s'", cmd.Command)
		}
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
		input:     "", // Empty input to show keybind help
		cursorPos: 0,
		filtered:  []CommandInfo{}, // Empty to trigger keybind help
	}

	mockGitClient := testutil.NewMockGitClient()

	ui := &UI{
		stdin:     strings.NewReader(""),
		stdout:    &buf,
		stderr:    &bytes.Buffer{},
		term:      &mockTerminal{},
		renderer:  renderer,
		state:     state,
		colors:    colors,
		gitClient: mockGitClient,
		gitStatus: getGitStatus(mockGitClient),
		workflow:  NewWorkflow(),
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
		"Enter",
		"Ctrl+c",
		"Tab",
		"Ctrl+t",
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
		workflow: NewWorkflow(),
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
	// Create mock git client
	mockClient := testutil.NewMockGitClient()

	status := getGitStatus(mockClient)

	if status == nil || status.Branch != "main" {
		if status == nil {
			t.Fatal("Expected status to be non-nil with mock client")
		} else {
			t.Errorf("Expected branch name to be 'main', got %s", status.Branch)
		}
	}

	// Should have 1 modified and 1 staged file
	if status.Modified != 1 {
		t.Errorf("Expected 1 modified file, got %d", status.Modified)
	}
	if status.Staged != 1 {
		t.Errorf("Expected 1 staged file, got %d", status.Staged)
	}

	// Should have changes
	if !status.HasChanges {
		t.Error("Expected HasChanges to be true")
	}

	// Should have ahead/behind counts
	if status.Ahead != 2 {
		t.Errorf("Expected ahead count to be 2, got %d", status.Ahead)
	}
	if status.Behind != 1 {
		t.Errorf("Expected behind count to be 1, got %d", status.Behind)
	}
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
		workflow:  NewWorkflow(),
		gitStatus: mockStatus,
	}

	renderer.renderGitStatus(ui, mockStatus)
	output := buf.String()

	// Check that git status elements are displayed
	expectedElements := []string{
		"📍",          // Branch icon
		"main",       // Branch name
		"📝",          // Changes icon
		"2 modified", // Modified files
		"1 staged",   // Staged files
		"↑3",         // Ahead count
		"↓1",         // Behind count
	}

	for _, element := range expectedElements {
		if !strings.Contains(output, element) {
			t.Errorf("Expected git status to contain '%s', but it was not found", element)
		}
	}
}

// Test individual render methods
func TestRenderer_RenderHeader(t *testing.T) {
	var buf bytes.Buffer
	colors := NewANSIColors()
	renderer := &Renderer{
		writer: &buf,
		colors: colors,
		width:  80,
		height: 24,
	}

	ui := &UI{
		stdin:     strings.NewReader(""),
		stdout:    &buf,
		stderr:    &bytes.Buffer{},
		term:      &mockTerminal{},
		renderer:  renderer,
		colors:    colors,
		workflow:  NewWorkflow(),
		gitStatus: nil, // No git status
	}

	renderer.renderHeader(ui)
	output := buf.String()

	// Check that header elements are present
	expectedElements := []string{
		"🚀 ggc Interactive Mode",
	}

	for _, element := range expectedElements {
		if !strings.Contains(output, element) {
			t.Errorf("Expected header to contain '%s', but it was not found", element)
		}
	}
}

func TestRenderer_FormatInputWithCursor(t *testing.T) {
	colors := NewANSIColors()
	renderer := &Renderer{
		writer: &bytes.Buffer{},
		colors: colors,
	}

	tests := []struct {
		name      string
		input     string
		cursorPos int
		expected  string
	}{
		{
			name:      "empty input",
			input:     "",
			cursorPos: 0,
			expected:  "█", // Should contain block cursor
		},
		{
			name:      "cursor at end",
			input:     "test",
			cursorPos: 4,
			expected:  "█", // Should contain block cursor at end
		},
		{
			name:      "cursor in middle",
			input:     "test",
			cursorPos: 2,
			expected:  "│", // Should contain line cursor in middle
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &UIState{
				input:     tt.input,
				cursorPos: tt.cursorPos,
			}

			result := renderer.formatInputWithCursor(state)
			if !strings.Contains(result, tt.expected) {
				t.Errorf("Expected cursor format to contain '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestRenderer_RenderCommandItem(t *testing.T) {
	var buf bytes.Buffer
	colors := NewANSIColors()
	renderer := &Renderer{
		writer: &buf,
		colors: colors,
		width:  80,
		height: 24,
	}

	ui := &UI{
		stdin:    strings.NewReader(""),
		stdout:   &buf,
		stderr:   &bytes.Buffer{},
		term:     &mockTerminal{},
		renderer: renderer,
		colors:   colors,
		workflow: NewWorkflow(),
	}

	cmd := CommandInfo{
		Command:     "test command",
		Description: "Test description",
	}

	// Test selected item
	buf.Reset()
	renderer.renderCommandItem(ui, cmd, 0, 0, 20) // index=0, selected=0
	output := buf.String()
	if !strings.Contains(output, "▶") {
		t.Error("Expected selected item to contain '▶' indicator")
	}

	// Test non-selected item
	buf.Reset()
	renderer.renderCommandItem(ui, cmd, 1, 0, 20) // index=1, selected=0
	output = buf.String()
	if strings.Contains(output, "▶") {
		t.Error("Expected non-selected item to NOT contain '▶' indicator")
	}
}

// Test interactive input functionality
func TestKeyHandler_InteractiveInput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	colors := NewANSIColors()
	mockGitClient := testutil.NewMockGitClient()

	// Mock terminal that fails raw mode to trigger fallback
	mockTerm := &mockTerminal{shouldFailRaw: true}

	ui := &UI{
		stdin:     strings.NewReader("test value\n"),
		stdout:    &stdout,
		stderr:    &stderr,
		colors:    colors,
		workflow:  NewWorkflow(),
		term:      mockTerm,
		gitClient: mockGitClient,
	}

	handler := &KeyHandler{ui: ui}
	ui.handler = handler

	placeholders := []string{"message"}
	result, canceled := handler.interactiveInput(placeholders)
	if canceled {
		t.Fatal("interactive input should not be canceled")
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 input, got %d", len(result))
	}

	if result["message"] != "test value" {
		t.Errorf("Expected 'test value', got '%s'", result["message"])
	}

	output := stdout.String()
	if !strings.Contains(output, "message") {
		t.Error("Expected output to contain placeholder name")
	}
	if !strings.Contains(output, "✓") {
		t.Error("Expected output to contain confirmation checkmark")
	}
}

func TestNewUI_WiresContextualResolver(t *testing.T) {
	const overrideEnv = "GGC_KEYBIND_DELETE_WORD"
	t.Setenv(overrideEnv, "ctrl+q")

	mockClient := testutil.NewMockGitClient()
	ui := NewUI(mockClient)

	if ui.handler == nil {
		t.Fatal("handler should be initialized")
	}

	contextual := ui.handler.contextualMap
	if contextual == nil {
		t.Fatal("expected contextual keybinding map to be set")
	}

	globalMap, exists := contextual.GetContext(ContextGlobal)
	if !exists || globalMap == nil {
		t.Fatalf("expected global context keymap, got exists=%v map=%v", exists, globalMap)
	}

	if !globalMap.MatchesKeyStroke("delete_word", NewCtrlKeyStroke('q')) {
		t.Fatal("environment override should be reflected in global keymap")
	}

	currentMap := ui.handler.GetCurrentKeyMap()
	if !currentMap.MatchesKeyStroke("delete_word", NewCtrlKeyStroke('q')) {
		t.Error("handler current map should resolve overrides via contextual map")
	}
}

func TestKeyHandler_ProcessCommand_NoPlaceholders(t *testing.T) {
	var stdout, stderr bytes.Buffer
	colors := NewANSIColors()

	ui := &UI{
		stdin:  strings.NewReader(""),
		stdout: &stdout,
		stderr: &stderr,
		colors: colors,
		term:   &mockTerminal{},
	}

	handler := &KeyHandler{ui: ui}
	ui.handler = handler

	// Command without placeholders should execute immediately
	result, canceled := handler.processCommand("status")
	if canceled {
		t.Fatal("processCommand should not cancel for commands without placeholders")
	}

	expected := []string{"ggc", "status"}
	if len(result) != len(expected) {
		t.Errorf("Expected %d args, got %d", len(expected), len(result))
	}

	for i, arg := range expected {
		if result[i] != arg {
			t.Errorf("Expected arg[%d] to be '%s', got '%s'", i, arg, result[i])
		}
	}
}

func TestKeyHandler_ProcessCommand_WithPlaceholders(t *testing.T) {
	var stdout, stderr bytes.Buffer
	colors := NewANSIColors()

	// Mock terminal that fails raw mode to trigger fallback
	mockTerm := &mockTerminal{shouldFailRaw: true}

	ui := &UI{
		stdin:    strings.NewReader("fix bug\n"),
		stdout:   &stdout,
		stderr:   &stderr,
		colors:   colors,
		workflow: NewWorkflow(),
		term:     mockTerm,
	}

	handler := &KeyHandler{ui: ui}
	ui.handler = handler

	result, canceled := handler.processCommand("commit <message>")
	if canceled {
		t.Fatal("processCommand should not cancel for provided placeholder input")
	}

	expected := []string{"ggc", "commit", "fix", "bug"}
	if len(result) != len(expected) {
		t.Errorf("Expected %d args, got %d", len(expected), len(result))
	}

	for i, arg := range expected {
		if result[i] != arg {
			t.Errorf("Expected arg[%d] to be '%s', got '%s'", i, arg, result[i])
		}
	}
}

func TestKeyHandler_HandleSoftCancelResetsState(t *testing.T) {
	colors := NewANSIColors()
	ui := &UI{
		stdout: &bytes.Buffer{},
		colors: colors,
		state: &UIState{
			input:        "status",
			context:      ContextSearch,
			contextStack: []Context{ContextInput},
			showWorkflow: true,
			selected:     2,
		},
	}

	handler := &KeyHandler{ui: ui}
	ui.handler = handler

	handler.handleSoftCancel(nil)

	if ui.state.HasInput() {
		t.Error("expected input to be cleared after soft cancel")
	}
	if ui.state.showWorkflow {
		t.Error("expected workflow view to be hidden after soft cancel")
	}
	if ui.state.GetCurrentContext() != ContextGlobal {
		t.Errorf("expected context to reset to global, got %s", ui.state.GetCurrentContext())
	}
	if !ui.consumeSoftCancelFlash() {
		t.Error("expected soft cancel flash flag to be set")
	}

	// Second soft cancel with no active state should not flash again
	handler.handleSoftCancel(nil)
	if ui.consumeSoftCancelFlash() {
		t.Error("expected no flash when soft cancel called without active operation")
	}
}

func TestShouldHandleEscapeAsSoftCancel(t *testing.T) {
	ui := &UI{
		stdin:  os.Stdin,
		stdout: &bytes.Buffer{},
		colors: NewANSIColors(),
		state: &UIState{
			context: ContextGlobal,
		},
	}

	handler := &KeyHandler{ui: ui}
	ui.handler = handler

	restoreZero := termio.SetPendingInputFunc(func(uintptr) (int, error) { return 0, nil })
	got := handler.shouldHandleEscapeAsSoftCancel()
	restoreZero()
	if !got {
		t.Fatal("expected soft cancel to trigger when no input is pending")
	}

	ui.reader = bufio.NewReader(strings.NewReader("buffered"))
	_, _ = ui.reader.Peek(1)
	if handler.shouldHandleEscapeAsSoftCancel() {
		t.Fatal("expected buffered reader to prevent soft cancel")
	}
	ui.reader = nil

	restorePending := termio.SetPendingInputFunc(func(uintptr) (int, error) { return 1, nil })
	if handler.shouldHandleEscapeAsSoftCancel() {
		t.Fatal("expected pending input to prevent soft cancel")
	}
	restorePending()
}

func TestRealTimeEditorShouldSoftCancelOnEscape(t *testing.T) {
	colors := NewANSIColors()
	ui := &UI{stdout: &bytes.Buffer{}, colors: colors}
	inputRunes := make([]rune, 0)
	cursor := 0
	editor := &realTimeEditor{
		ui:         ui,
		inputRunes: &inputRunes,
		cursor:     &cursor,
	}

	restoreZero := termio.SetPendingInputFunc(func(uintptr) (int, error) { return 0, nil })
	if !editor.shouldSoftCancelOnEscape(bufio.NewReader(strings.NewReader(""))) {
		t.Fatal("expected soft cancel when no bytes pending")
	}
	restoreZero()

	restorePending := termio.SetPendingInputFunc(func(uintptr) (int, error) { return 1, nil })
	if editor.shouldSoftCancelOnEscape(bufio.NewReader(strings.NewReader(""))) {
		t.Fatal("expected pending input to block soft cancel")
	}
	restorePending()

	bufferedReader := bufio.NewReader(strings.NewReader("buffer"))
	_, _ = bufferedReader.Peek(1)
	if editor.shouldSoftCancelOnEscape(bufferedReader) {
		t.Fatal("expected buffered reader to block soft cancel")
	}
}

func TestKeyHandler_GetLineInput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	colors := NewANSIColors()

	ui := &UI{
		stdin:  strings.NewReader("test input\n"),
		stdout: &stdout,
		stderr: &stderr,
		colors: colors,
		term:   &mockTerminal{},
	}

	handler := &KeyHandler{ui: ui}

	result, canceled := handler.getLineInput()
	if canceled {
		t.Fatal("expected line input to succeed")
	}

	if result != "test input" {
		t.Errorf("Expected 'test input', got '%s'", result)
	}
}

func TestKeyHandler_GetLineInput_EmptyInput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	colors := NewANSIColors()

	// Simulate empty input followed by valid input
	ui := &UI{
		stdin:  strings.NewReader("\nvalid input\n"),
		stdout: &stdout,
		stderr: &stderr,
		colors: colors,
		term:   &mockTerminal{},
	}

	handler := &KeyHandler{ui: ui}

	result, canceled := handler.getLineInput()
	if canceled {
		t.Fatal("expected line input to succeed after retry")
	}

	if result != "valid input" {
		t.Errorf("Expected 'valid input', got '%s'", result)
	}

	// Check that required message was shown
	output := stdout.String()
	if !strings.Contains(output, "(required)") {
		t.Error("Expected output to contain '(required)' message for empty input")
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
		input:     "", // Empty input to show keybind help
		cursorPos: 0,
		filtered:  []CommandInfo{}, // Empty to show keybind help
	}

	ui := &UI{
		stdin:    strings.NewReader(""),
		stdout:   &buf,
		stderr:   &bytes.Buffer{},
		term:     &mockTerminal{},
		renderer: renderer,
		state:    state,
		colors:   colors,
		workflow: NewWorkflow(),
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
		"Ctrl+c",
		"Tab",
		"Ctrl+t",
	}

	for _, keybind := range expectedKeybinds {
		if !strings.Contains(output, keybind) {
			t.Errorf("Expected output to contain keybind '%s', but it was not found", keybind)
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
		workflow: NewWorkflow(),
	}

	handler := &KeyHandler{ui: ui}

	// Test printable character
	shouldContinue, result := handler.HandleKey('a', true, nil)
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
	shouldContinue, _ = handler.HandleKey(127, true, nil)
	if !shouldContinue {
		t.Error("Expected to continue after backspace")
	}
	if ui.state.input != "" {
		t.Errorf("Expected empty input after backspace, got '%s'", ui.state.input)
	}
}

// TestUIState_MultibyteContinuousDeletion tests the bug where multibyte characters
// couldn't be deleted continuously after some deletions
func TestUIState_MultibyteContinuousDeletion(t *testing.T) {
	ui := &UI{
		state: &UIState{
			input:     "",
			cursorPos: 0,
			filtered:  []CommandInfo{},
		},
	}

	// Add multibyte characters
	text := "こんにちは世界"
	for _, r := range text {
		ui.state.AddRune(r)
	}

	expectedRunes := []rune(text)
	if ui.state.input != text {
		t.Errorf("Expected input '%s', got '%s'", text, ui.state.input)
	}
	if ui.state.cursorPos != len(expectedRunes) {
		t.Errorf("Expected cursor position %d, got %d", len(expectedRunes), ui.state.cursorPos)
	}

	// Delete all characters one by one
	for i := len(expectedRunes); i > 0; i-- {
		// Before deletion
		if ui.state.cursorPos != i {
			t.Errorf("Before deletion %d: expected cursor position %d, got %d",
				len(expectedRunes)-i+1, i, ui.state.cursorPos)
		}

		// Perform deletion
		ui.state.RemoveChar()

		// After deletion
		expectedAfter := string(expectedRunes[:i-1])
		if ui.state.input != expectedAfter {
			t.Errorf("After deletion %d: expected input '%s', got '%s'",
				len(expectedRunes)-i+1, expectedAfter, ui.state.input)
		}
		if ui.state.cursorPos != i-1 {
			t.Errorf("After deletion %d: expected cursor position %d, got %d",
				len(expectedRunes)-i+1, i-1, ui.state.cursorPos)
		}
	}

	// Final state should be empty
	if ui.state.input != "" {
		t.Errorf("Expected empty input after all deletions, got '%s'", ui.state.input)
	}
	if ui.state.cursorPos != 0 {
		t.Errorf("Expected cursor position 0 after all deletions, got %d", ui.state.cursorPos)
	}
}

// TestUIState_DeleteWord_Multibyte tests word deletion with multibyte characters
func TestUIState_DeleteWord_Multibyte(t *testing.T) {
	ui := &UI{
		state: &UIState{
			input:     "",
			cursorPos: 0,
			filtered:  []CommandInfo{},
		},
	}

	// Add text with multibyte words
	text := "hello こんにちは world"
	for _, r := range text {
		ui.state.AddRune(r)
	}

	// Position cursor after "world"
	ui.state.MoveToEnd()

	// Delete "world"
	ui.state.DeleteWord()
	expected := "hello こんにちは "
	if ui.state.input != expected {
		t.Errorf("After deleting 'world': expected '%s', got '%s'", expected, ui.state.input)
	}

	// Delete "こんにちは"
	ui.state.DeleteWord()
	expected = "hello "
	if ui.state.input != expected {
		t.Errorf("After deleting 'こんにちは': expected '%s', got '%s'", expected, ui.state.input)
	}

	// Delete "hello"
	ui.state.DeleteWord()
	expected = ""
	if ui.state.input != expected {
		t.Errorf("After deleting 'hello': expected '%s', got '%s'", expected, ui.state.input)
	}
}

// TestUIState_DeleteToEnd_Multibyte tests deleting to end with multibyte characters
func TestUIState_DeleteToEnd_Multibyte(t *testing.T) {
	ui := &UI{
		state: &UIState{
			input:     "",
			cursorPos: 0,
			filtered:  []CommandInfo{},
		},
	}

	// Add multibyte text
	text := "こんにちは世界"
	for _, r := range text {
		ui.state.AddRune(r)
	}

	// Move cursor to middle (after "こん")
	ui.state.cursorPos = 2

	// Delete to end
	ui.state.DeleteToEnd()
	expected := "こん"
	if ui.state.input != expected {
		t.Errorf("After DeleteToEnd: expected '%s', got '%s'", expected, ui.state.input)
	}
	if ui.state.cursorPos != 2 {
		t.Errorf("After DeleteToEnd: expected cursor position 2, got %d", ui.state.cursorPos)
	}
}

// TestHandleInputChar_MultibyteBackspace tests multibyte character deletion in placeholder input
func TestHandleInputChar_MultibyteBackspace(t *testing.T) {
	ui := &UI{
		colors: NewANSIColors(),
		stdout: &strings.Builder{},
	}
	handler := &KeyHandler{ui: ui}

	// Test with Japanese characters
	var input strings.Builder

	// Add some multibyte characters
	text := "こんにちは"
	for _, r := range text {
		done, canceled := handler.handleInputChar(&input, r)
		if done || canceled {
			t.Errorf("Unexpected completion during character input: done=%v, canceled=%v", done, canceled)
		}
	}

	if input.String() != text {
		t.Errorf("Expected input '%s', got '%s'", text, input.String())
	}

	// Test backspace deletion
	runesExpected := []rune(text)
	for i := len(runesExpected); i > 0; i-- {
		// Perform backspace
		done, canceled := handler.handleInputChar(&input, '\b')
		if done || canceled {
			t.Errorf("Unexpected completion during backspace: done=%v, canceled=%v", done, canceled)
		}

		// Check remaining content
		expected := string(runesExpected[:i-1])
		if input.String() != expected {
			t.Errorf("After backspace %d: expected '%s', got '%s'",
				len(runesExpected)-i+1, expected, input.String())
		}
	}

	// Final state should be empty
	if input.String() != "" {
		t.Errorf("Expected empty input after all backspaces, got '%s'", input.String())
	}
}

// TestHandleInputChar_MultibyteDisplay tests multibyte character display in placeholder input
func TestHandleInputChar_MultibyteDisplay(t *testing.T) {
	var output strings.Builder
	ui := &UI{
		colors: NewANSIColors(),
		stdout: &output,
	}
	handler := &KeyHandler{ui: ui}

	// Test various multibyte characters
	testCases := []struct {
		name     string
		char     rune
		expected string
	}{
		{"Japanese Hiragana", 'こ', "こ"},
		{"Japanese Katakana", 'ア', "ア"},
		{"Japanese Kanji", '漢', "漢"},
		{"Chinese Character", '中', "中"},
		{"Emoji", '🚀', "🚀"},
		{"Accented Character", 'é', "é"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var input strings.Builder
			output.Reset()

			// Add the multibyte character
			done, canceled := handler.handleInputChar(&input, tc.char)
			if done || canceled {
				t.Errorf("Unexpected completion during character input: done=%v, canceled=%v", done, canceled)
			}

			// Check input buffer
			if input.String() != tc.expected {
				t.Errorf("Expected input '%s', got '%s'", tc.expected, input.String())
			}

			// Check terminal output contains the character
			terminalOutput := output.String()
			if !strings.Contains(terminalOutput, tc.expected) {
				t.Errorf("Expected terminal output to contain '%s', got '%s'", tc.expected, terminalOutput)
			}
		})
	}
}

// TestKeyHandler_HandleWorkflowKeys tests workflow key handling
func TestKeyHandler_HandleWorkflowKeys(t *testing.T) {
	tests := []struct {
		name           string
		key            rune
		showWorkflow   bool
		hasInput       bool
		selectedCmd    *CommandInfo
		expectedResult bool
		expectedAction string
	}{
		{
			name:           "Tab key adds command to workflow",
			key:            '\t',
			showWorkflow:   false,
			hasInput:       true,
			selectedCmd:    &CommandInfo{Command: "add .", Description: "Add all changes"},
			expectedResult: true,
			expectedAction: "add_to_workflow",
		},
		{
			name:           "Tab key in workflow view returns true but no action",
			key:            '\t',
			showWorkflow:   true,
			hasInput:       true,
			selectedCmd:    &CommandInfo{Command: "add .", Description: "Add all changes"},
			expectedResult: true,
			expectedAction: "no_action",
		},
		{
			name:           "c key clears workflow in workflow view",
			key:            'c',
			showWorkflow:   true,
			hasInput:       false,
			selectedCmd:    nil,
			expectedResult: true,
			expectedAction: "clear_workflow",
		},
		{
			name:           "c key in search view returns true but no action",
			key:            'c',
			showWorkflow:   false,
			hasInput:       true,
			selectedCmd:    nil,
			expectedResult: true,
			expectedAction: "no_action",
		},
		{
			name:           "unhandled key returns false",
			key:            'x',
			showWorkflow:   false,
			hasInput:       true,
			selectedCmd:    nil,
			expectedResult: false,
			expectedAction: "none",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gitClient := &mockGitClient{}
			ui := NewUI(gitClient)
			ui.state.showWorkflow = tt.showWorkflow
			if tt.hasInput {
				ui.state.input = "test"
				ui.state.filtered = []CommandInfo{
					{Command: "add .", Description: "Add all changes"},
				}
				ui.state.selected = 0
			}

			// Mock workflow for clear test
			if tt.expectedAction == "clear_workflow" {
				ui.workflow.AddStep("test", []string{}, "test command")
			}

			// Execute
			result := ui.handler.handleWorkflowKeys(tt.key)

			// Verify
			if result != tt.expectedResult {
				t.Errorf("Expected result %v, got %v", tt.expectedResult, result)
			}

			// Verify side effects
			switch tt.expectedAction {
			case "add_to_workflow":
				if ui.workflow.IsEmpty() {
					t.Error("Expected workflow to have steps after adding command")
				}
			case "clear_workflow":
				if !ui.workflow.IsEmpty() {
					t.Error("Expected workflow to be empty after clearing")
				}
			}
		})
	}
}

// TestKeyHandler_AddCommandToWorkflow tests adding commands to workflow
func TestKeyHandler_AddCommandToWorkflow(t *testing.T) {
	tests := []struct {
		name        string
		cmdTemplate string
		expectSteps int
	}{
		{
			name:        "Add simple command",
			cmdTemplate: "add .",
			expectSteps: 1,
		},
		{
			name:        "Add command with placeholder",
			cmdTemplate: "commit <message>",
			expectSteps: 1,
		},
		{
			name:        "Add complex command",
			cmdTemplate: "branch create <name>",
			expectSteps: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gitClient := &mockGitClient{}
			ui := NewUI(gitClient)

			// Redirect output to avoid test output pollution
			ui.stdout = &bytes.Buffer{}

			// Execute
			ui.handler.addCommandToWorkflow(tt.cmdTemplate)

			// Verify
			steps := ui.workflow.GetSteps()
			if len(steps) != tt.expectSteps {
				t.Errorf("Expected %d steps, got %d", tt.expectSteps, len(steps))
			}

			if len(steps) > 0 {
				if steps[0].Description != tt.cmdTemplate {
					t.Errorf("Expected description '%s', got '%s'", tt.cmdTemplate, steps[0].Description)
				}
			}
		})
	}
}

// TestKeyHandler_ClearWorkflow tests clearing workflow
func TestKeyHandler_ClearWorkflow(t *testing.T) {
	// Setup
	gitClient := &mockGitClient{}
	ui := NewUI(gitClient)

	// Redirect output to avoid test output pollution
	ui.stdout = &bytes.Buffer{}

	// Add some steps
	ui.workflow.AddStep("add", []string{"."}, "add .")
	ui.workflow.AddStep("commit", []string{"-m", "test"}, "commit -m test")

	if ui.workflow.IsEmpty() {
		t.Fatal("Expected workflow to have steps before clearing")
	}

	// Execute
	ui.handler.clearWorkflow()

	// Verify
	if !ui.workflow.IsEmpty() {
		t.Error("Expected workflow to be empty after clearing")
	}
}

// TestUI_ToggleWorkflowView tests workflow view toggling
func TestUI_ToggleWorkflowView(t *testing.T) {
	// Setup
	gitClient := &mockGitClient{}
	ui := NewUI(gitClient)

	// Initial state should be false
	if ui.state.showWorkflow {
		t.Error("Expected initial showWorkflow to be false")
	}

	// Toggle to true
	ui.ToggleWorkflowView()
	if !ui.state.showWorkflow {
		t.Error("Expected showWorkflow to be true after first toggle")
	}

	// Toggle back to false
	ui.ToggleWorkflowView()
	if ui.state.showWorkflow {
		t.Error("Expected showWorkflow to be false after second toggle")
	}
}

// TestUI_WorkflowOperations tests workflow operations
func TestUI_WorkflowOperations(t *testing.T) {
	// Setup
	gitClient := &mockGitClient{}
	ui := NewUI(gitClient)

	// Test AddToWorkflow
	id := ui.AddToWorkflow("add", []string{"."}, "add .")
	if id != 1 {
		t.Errorf("Expected first workflow ID to be 1, got %d", id)
	}

	steps := ui.workflow.GetSteps()
	if len(steps) != 1 {
		t.Errorf("Expected 1 step, got %d", len(steps))
	}

	// Test ClearWorkflow
	ui.ClearWorkflow()
	if !ui.workflow.IsEmpty() {
		t.Error("Expected workflow to be empty after clearing")
	}
}

// TestKeyHandler_ExecuteWorkflow tests workflow execution
func TestKeyHandler_ExecuteWorkflow(t *testing.T) {
	tests := []struct {
		name           string
		workflowSteps  []struct{ cmd, desc string }
		expectError    bool
		expectResult   bool
		expectContinue bool
	}{
		{
			name:           "Execute empty workflow",
			workflowSteps:  []struct{ cmd, desc string }{},
			expectError:    false,
			expectResult:   false, // Empty workflow returns (true, nil)
			expectContinue: true,
		},
		{
			name: "Execute workflow with steps",
			workflowSteps: []struct{ cmd, desc string }{
				{"add .", "add ."},
				{"commit -m test", "commit -m test"},
			},
			expectError:    false,
			expectResult:   true, // Non-empty workflow returns result
			expectContinue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gitClient := &mockGitClient{}
			ui := NewUI(gitClient)

			// Redirect output to avoid test output pollution
			ui.stdout = &bytes.Buffer{}

			// Set up workflow executor with mock router
			mockRouter := &mockRouterForExecute{}
			ui.workflowEx = NewWorkflowExecutor(mockRouter, ui)

			// Add workflow steps
			for _, step := range tt.workflowSteps {
				parts := strings.Fields(step.cmd)
				if len(parts) > 0 {
					ui.workflow.AddStep(parts[0], parts[1:], step.desc)
				}
			}

			// Execute
			shouldContinue, result := ui.handler.executeWorkflow(nil)

			// Verify
			if shouldContinue != tt.expectContinue {
				t.Errorf("Expected shouldContinue %v, got %v", tt.expectContinue, shouldContinue)
			}

			if tt.expectResult {
				if result == nil {
					t.Error("Expected result to not be nil")
				} else if len(result) < 2 || result[1] != InteractiveWorkflowCommand {
					t.Errorf("Expected result to contain workflow command, got %v", result)
				}
			} else {
				if result != nil {
					t.Errorf("Expected result to be nil for empty workflow, got %v", result)
				}
			}
		})
	}
}

// mockRouterForExecute is a mock router for testing workflow execution
type mockRouterForExecute struct {
	routedCommands [][]string
}

func (m *mockRouterForExecute) Route(args []string) {
	if m.routedCommands == nil {
		m.routedCommands = make([][]string, 0)
	}
	m.routedCommands = append(m.routedCommands, args)
}
