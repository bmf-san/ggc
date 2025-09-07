package cmd

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v5/internal/testutil"
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
		"commit allow-empty",
		"commit amend",
		"commit amend --no-edit",
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

// Test multibyte character input support
func TestUIState_AddRune_MultibyteCcharacters(t *testing.T) {
	state := &UIState{
		selected:  0,
		input:     "",
		cursorPos: 0,
		filtered:  []CommandInfo{},
	}

	// Test Japanese characters
	testCases := []struct {
		name     string
		rune     rune
		expected string
	}{
		{"Hiragana", 'こ', "こ"},
		{"Katakana", 'テ', "テ"},
		{"Kanji", '機', "機"},
		{"Emoji", '🚀', "🚀"},
		{"Chinese", '中', "中"},
		{"Korean", '한', "한"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset state
			state.input = ""
			state.cursorPos = 0

			// Add the multibyte character
			state.AddRune(tc.rune)

			if state.input != tc.expected {
				t.Errorf("Expected input '%s', got '%s'", tc.expected, state.input)
			}

			if state.cursorPos != 1 {
				t.Errorf("Expected cursor position 1, got %d", state.cursorPos)
			}
		})
	}
}

// Test multibyte character removal
func TestUIState_RemoveChar_Multibyte(t *testing.T) {
	state := &UIState{
		selected:  0,
		input:     "こんにちは", // "Hello" in Japanese
		cursorPos: 5,       // At the end
		filtered:  []CommandInfo{},
	}

	// Remove last character (は)
	state.RemoveChar()

	expected := "こんにち"
	if state.input != expected {
		t.Errorf("Expected input '%s', got '%s'", expected, state.input)
	}

	if state.cursorPos != 4 {
		t.Errorf("Expected cursor position 4, got %d", state.cursorPos)
	}

	// Remove another character (ち)
	state.RemoveChar()

	expected = "こんに"
	if state.input != expected {
		t.Errorf("Expected input '%s', got '%s'", expected, state.input)
	}

	if state.cursorPos != 3 {
		t.Errorf("Expected cursor position 3, got %d", state.cursorPos)
	}
}

// Test fuzzy matching with multibyte characters
func TestUIState_UpdateFiltered_MultibyteFuzzy(t *testing.T) {
	// Add a mock command with multibyte characters for testing
	originalCommands := commands
	defer func() { commands = originalCommands }()

	commands = []CommandInfo{
		{"commit 機能追加", "Add feature commit"},
		{"branch テスト", "Test branch"},
		{"add ファイル", "Add file"},
		{"commit message", "Regular commit"},
	}

	state := &UIState{
		selected:  0,
		input:     "機能", // Should match "commit 機能追加"
		cursorPos: 2,
		filtered:  []CommandInfo{},
	}

	state.UpdateFiltered()

	// Should find the command containing "機能"
	found := false
	for _, cmd := range state.filtered {
		if strings.Contains(cmd.Command, "機能追加") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find command containing '機能追加' with fuzzy pattern '機能'")
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
	//nolint:staticcheck // SA5011: false positive - this nil check is intentional
	if cmd == nil {
		t.Fatal("Expected non-nil command")
	}
	//nolint:staticcheck // SA5011: false positive after t.Fatal
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
	// Create mock git client
	mockClient := testutil.NewMockGitClient()

	status := getGitStatus(mockClient)

	//nolint:staticcheck // SA5011: false positive - this nil check is intentional
	if status == nil {
		t.Fatal("Expected status to be non-nil with mock client")
	}

	// Branch name should match mock
	//nolint:staticcheck // SA5011: false positive after t.Fatal
	if status.Branch != "main" {
		t.Errorf("Expected branch name to be 'main', got %s", status.Branch)
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
		gitStatus: nil, // No git status
	}

	renderer.renderHeader(ui)
	output := buf.String()

	// Check that header elements are present
	expectedElements := []string{
		"🚀 ggc Interactive Mode",
		"Type to search",
		"Ctrl+n/p",
		"navigate",
		"Enter",
		"execute",
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
		term:      mockTerm,
		gitClient: mockGitClient,
	}

	handler := &KeyHandler{ui: ui}

	placeholders := []string{"message"}
	result := handler.interactiveInput(placeholders)

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

	// Command without placeholders should execute immediately
	result := handler.processCommand("status")

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
		stdin:  strings.NewReader("fix bug\n"),
		stdout: &stdout,
		stderr: &stderr,
		colors: colors,
		term:   mockTerm,
	}

	handler := &KeyHandler{ui: ui}

	result := handler.processCommand("commit <message>")

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

	result := handler.getLineInput()

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

	result := handler.getLineInput()

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
