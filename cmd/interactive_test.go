package cmd

import (
	"bytes"
	"fmt"
	"reflect"
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
			// Create test UI
			ui := &testUI{
				UI: UI{
					term: &mockTerminal{},
				},
				inputBytes: tt.input,
			}

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
