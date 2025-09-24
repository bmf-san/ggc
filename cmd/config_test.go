package cmd

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v6/internal/testutil"
)

// Mock config manager for testing
type mockConfigManager struct {
	configs   map[string]any
	loadError error
	getError  error
	setError  error
}

func (m *mockConfigManager) Load() error {
	if m.loadError != nil {
		return m.loadError
	}
	return nil
}

func (m *mockConfigManager) List() map[string]any {
	return m.configs
}

func (m *mockConfigManager) Get(key string) (any, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	if value, exists := m.configs[key]; exists {
		return value, nil
	}
	return nil, errors.New("key not found")
}

func (m *mockConfigManager) Set(key string, value any) error {
	if m.setError != nil {
		return m.setError
	}
	if m.configs == nil {
		m.configs = make(map[string]any)
	}
	m.configs[key] = value
	return nil
}

// Mock function to replace config.NewConfigManager
var mockNewConfigManager func() interface {
	Load() error
	List() map[string]any
	Get(string) (any, error)
	Set(string, any) error
}

func TestConfigurer_Config(t *testing.T) {
	cases := []struct {
		name           string
		args           []string
		mockConfig     *mockConfigManager
		expectedOutput []string
		notContains    []string
	}{
		{
			name:           "config no args shows help",
			args:           []string{},
			mockConfig:     &mockConfigManager{},
			expectedOutput: []string{"Usage:"},
		},
		{
			name:           "config list",
			args:           []string{"list"},
			mockConfig:     &mockConfigManager{},
			expectedOutput: []string{"ui.color"},
		},
		{
			name: "config get non-existing key",
			args: []string{"get", "nonexistent.key"},
			mockConfig: &mockConfigManager{
				configs:  map[string]any{},
				getError: errors.New("key not found"),
			},
			expectedOutput: []string{
				"failed to get config value: ",
			},
		},
		{
			name: "config get missing key argument",
			args: []string{"get"},
			mockConfig: &mockConfigManager{
				configs: map[string]any{},
			},
			expectedOutput: []string{
				"must provide key to get (arg missing)",
			},
		},
		{
			name: "config set string value",
			args: []string{"set", "user.name", "jane.doe"},
			mockConfig: &mockConfigManager{
				configs: make(map[string]any),
			},
			expectedOutput: []string{
				"Set user.name = jane.doe",
			},
		},
		{
			name: "config set boolean value",
			args: []string{"set", "feature.enabled", "true"},
			mockConfig: &mockConfigManager{
				configs: make(map[string]any),
			},
			expectedOutput: []string{
				"Set feature.enabled = true",
			},
		},
		{
			name: "config set integer value",
			args: []string{"set", "timeout", "300"},
			mockConfig: &mockConfigManager{
				configs: make(map[string]any),
			},
			expectedOutput: []string{
				"Set timeout = 300",
			},
		},
		{
			name: "config set float value",
			args: []string{"set", "ratio", "1.5"},
			mockConfig: &mockConfigManager{
				configs: make(map[string]any),
			},
			expectedOutput: []string{
				"Set ratio = 1.5",
			},
		},
		{
			name: "config set missing arguments",
			args: []string{"set"},
			mockConfig: &mockConfigManager{
				configs: make(map[string]any),
			},
			expectedOutput: []string{
				"must provide key && value to set (arg(s) missing)",
			},
		},
		{
			name: "config set missing value argument",
			args: []string{"set", "key"},
			mockConfig: &mockConfigManager{
				configs: make(map[string]any),
			},
			expectedOutput: []string{
				"must provide key && value to set (arg(s) missing)",
			},
		},
		{
			name: "config set with error",
			args: []string{"set", "readonly.key", "value"},
			mockConfig: &mockConfigManager{
				configs:  make(map[string]any),
				setError: errors.New("cannot set readonly key"),
			},
			expectedOutput: []string{
				"failed to set config value: ",
			},
		},
		{
			name:           "config invalid command shows help",
			args:           []string{"invalid"},
			mockConfig:     &mockConfigManager{},
			expectedOutput: []string{"Usage:"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			// Set up mock for LoadConfig method
			originalMockNewConfigManager := mockNewConfigManager
			mockNewConfigManager = func() interface {
				Load() error
				List() map[string]any
				Get(string) (any, error)
				Set(string, any) error
			} {
				return tc.mockConfig
			}
			defer func() {
				mockNewConfigManager = originalMockNewConfigManager
			}()

			c := &Configurer{
				gitClient:    testutil.NewMockGitClient(),
				outputWriter: &buf,
				helper:       NewHelper(),
				execCommand:  exec.Command,
			}
			c.helper.outputWriter = &buf

			c.Config(tc.args)

			output := buf.String()
			for _, expected := range tc.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("expected output to contain %q, got %q", expected, output)
				}
			}

			for _, notExpected := range tc.notContains {
				if strings.Contains(output, notExpected) {
					t.Errorf("expected output to NOT contain %q, got %q", notExpected, output)
				}
			}
		})
	}
}

func TestConfigurer_LoadConfig(t *testing.T) {
	cases := []struct {
		name           string
		mockConfig     *mockConfigManager
		expectedOutput []string
		expectNil      bool
	}{
		{
			name: "load config success",
			mockConfig: &mockConfigManager{
				configs: map[string]any{"key": "value"},
			},
			expectedOutput: []string{},
			expectNil:      false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			// Set up mock for LoadConfig method
			originalMockNewConfigManager := mockNewConfigManager
			mockNewConfigManager = func() interface {
				Load() error
				List() map[string]any
				Get(string) (any, error)
				Set(string, any) error
			} {
				return tc.mockConfig
			}
			defer func() {
				mockNewConfigManager = originalMockNewConfigManager
			}()

			c := &Configurer{
				gitClient:    testutil.NewMockGitClient(),
				outputWriter: &buf,
				helper:       NewHelper(),
				execCommand:  exec.Command,
			}

			result := c.LoadConfig()

			output := buf.String()
			for _, expected := range tc.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("expected output to contain %q, got %q", expected, output)
				}
			}

			if tc.expectNil && result != nil {
				t.Error("expected LoadConfig to return nil, but got non-nil result")
			}
			if !tc.expectNil && result == nil {
				t.Error("expected LoadConfig to return non-nil, but got nil result")
			}
		})
	}
}

func TestFormatValue(t *testing.T) {
	cases := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name:     "string value",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "boolean true",
			input:    true,
			expected: "true",
		},
		{
			name:     "boolean false",
			input:    false,
			expected: "false",
		},
		{
			name:     "int value",
			input:    42,
			expected: "42",
		},
		{
			name:     "int8 value",
			input:    int8(127),
			expected: "127",
		},
		{
			name:     "int16 value",
			input:    int16(32767),
			expected: "32767",
		},
		{
			name:     "int32 value",
			input:    int32(2147483647),
			expected: "2147483647",
		},
		{
			name:     "int64 value",
			input:    int64(9223372036854775807),
			expected: "9223372036854775807",
		},
		{
			name:     "uint value",
			input:    uint(42),
			expected: "42",
		},
		{
			name:     "uint8 value",
			input:    uint8(255),
			expected: "255",
		},
		{
			name:     "uint16 value",
			input:    uint16(65535),
			expected: "65535",
		},
		{
			name:     "uint32 value",
			input:    uint32(4294967295),
			expected: "4294967295",
		},
		{
			name:     "uint64 value",
			input:    uint64(18446744073709551615),
			expected: "18446744073709551615",
		},
		{
			name:     "float32 value",
			input:    float32(3.14),
			expected: "3.14",
		},
		{
			name:     "float64 value",
			input:    float64(2.718281828),
			expected: "2.718281828",
		},
		{
			name:     "map value",
			input:    map[string]any{"key": "value"},
			expected: "map[key:value]",
		},
		{
			name:     "slice value",
			input:    []string{"a", "b", "c"},
			expected: "[a b c]",
		},
		{
			name:     "nil value",
			input:    nil,
			expected: "<nil>",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := formatValue(tc.input)
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestParseValue(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected any
	}{
		{
			name:     "boolean true",
			input:    "true",
			expected: true,
		},
		{
			name:     "boolean false",
			input:    "false",
			expected: false,
		},
		{
			name:     "integer value",
			input:    "42",
			expected: int64(42),
		},
		{
			name:     "negative integer",
			input:    "-123",
			expected: int64(-123),
		},
		{
			name:     "float value",
			input:    "3.14",
			expected: float64(3.14),
		},
		{
			name:     "negative float",
			input:    "-2.718",
			expected: float64(-2.718),
		},
		{
			name:     "string value",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "string that looks like bool but isn't",
			input:    "truthy",
			expected: "truthy",
		},
		{
			name:     "string that looks like number but isn't",
			input:    "123abc",
			expected: "123abc",
		},
		{
			name:     "scientific notation",
			input:    "1e10",
			expected: float64(1e10),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := parseValue(tc.input)
			if result != tc.expected {
				t.Errorf("expected %v (type %T), got %v (type %T)", tc.expected, tc.expected, result, result)
			}
		})
	}
}
