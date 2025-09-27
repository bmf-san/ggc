package router

import (
	"testing"
)

// Test debug-keys command routing integration
func TestRouter_DebugKeysIntegration(t *testing.T) {
	testCases := []struct {
		name string
		args []string
	}{
		{
			name: "debug-keys no args",
			args: []string{"debug-keys"},
		},
		{
			name: "debug-keys help",
			args: []string{"debug-keys", "help"},
		},
		{
			name: "debug-keys raw",
			args: []string{"debug-keys", "raw"},
		},
		{
			name: "debug-keys raw with file",
			args: []string{"debug-keys", "raw", "output.txt"},
		},
		{
			name: "debug-keys unknown subcommand",
			args: []string{"debug-keys", "unknown"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &mockExecuter{}
			r := NewRouter(m, nil)

			r.Route(tc.args)

			if !m.debugKeysCalled {
				t.Error("Expected DebugKeys to be called")
			}

			expectedArgs := tc.args[1:] // Remove "debug-keys" from args
			if len(m.debugKeysArgs) != len(expectedArgs) {
				t.Errorf("Expected %d args, got %d", len(expectedArgs), len(m.debugKeysArgs))
			}

			for i, arg := range expectedArgs {
				if i >= len(m.debugKeysArgs) || m.debugKeysArgs[i] != arg {
					t.Errorf("Expected arg %d to be '%s', got '%s'", i, arg, m.debugKeysArgs[i])
				}
			}
		})
	}
}

// Test that debug-keys command exists in command handlers map
func TestRouter_DebugKeysCommandHandlers(t *testing.T) {
	m := &mockExecuter{}
	r := NewRouter(m, nil)

	handlers := r.getCommandHandlers()

	if _, exists := handlers["debug-keys"]; !exists {
		t.Error("debug-keys command should exist in command handlers")
	}

	// Test that the handler actually calls the right method
	handlers["debug-keys"]([]string{"test", "args"})

	if !m.debugKeysCalled {
		t.Error("debug-keys handler should call DebugKeys method")
	}

	expectedArgs := []string{"test", "args"}
	if len(m.debugKeysArgs) != len(expectedArgs) {
		t.Errorf("Expected %d args, got %d", len(expectedArgs), len(m.debugKeysArgs))
	}
}

// Test debug-keys with nil config manager
func TestRouter_DebugKeysWithNilConfig(t *testing.T) {
	m := &mockExecuter{}
	r := NewRouter(m, nil)

	// Test that routing works with nil config manager
	r.Route([]string{"debug-keys", "help"})

	if !m.debugKeysCalled {
		t.Error("Expected DebugKeys to be called with nil config")
	}

	expectedArgs := []string{"help"}
	if len(m.debugKeysArgs) != len(expectedArgs) {
		t.Errorf("Expected %d args, got %d", len(expectedArgs), len(m.debugKeysArgs))
	}
}

// Test error handling in routing
func TestRouter_DebugKeysErrorHandling(t *testing.T) {
	testCases := []struct {
		name                  string
		args                  []string
		expectDebugKeysCalled bool
	}{
		{
			name:                  "normal debug-keys command",
			args:                  []string{"debug-keys", "help"},
			expectDebugKeysCalled: true,
		},
		{
			name:                  "debug-keys with args that look like flags",
			args:                  []string{"debug-keys", "--help"},
			expectDebugKeysCalled: true, // This is passed as args to debug-keys, not treated as legacy
		},
		{
			name:                  "legacy-like top level command",
			args:                  []string{"debug-interactive", "help"}, // This would trigger legacy detection
			expectDebugKeysCalled: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &mockExecuter{}
			r := NewRouter(m, nil)

			r.Route(tc.args)

			if tc.expectDebugKeysCalled && !m.debugKeysCalled {
				t.Error("Expected DebugKeys to be called")
			} else if !tc.expectDebugKeysCalled && m.debugKeysCalled {
				t.Error("Expected DebugKeys NOT to be called")
			}
		})
	}
}

// Benchmark routing performance
func BenchmarkRouter_DebugKeysRouting(b *testing.B) {
	m := &mockExecuter{}
	r := NewRouter(m, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Route([]string{"debug-keys"})
	}
}

func BenchmarkRouter_DebugKeysWithArgs(b *testing.B) {
	m := &mockExecuter{}
	r := NewRouter(m, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Route([]string{"debug-keys", "raw", "output.txt"})
	}
}
