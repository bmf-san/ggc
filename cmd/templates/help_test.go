package templates

import (
	"fmt"
	"strings"
	"testing"
)

func TestRenderMainHelp(t *testing.T) {
	tests := []struct {
		name     string
		logo     string
		expected []string
	}{
		{
			name: "full logo",
			logo: Logo,
			expected: []string{
				"ggc: A Go-based CLI tool to streamline Git operations",
				"Usage:",
			},
		},
		{
			name: "small logo",
			logo: SmallLogo,
			expected: []string{
				"ggc: A Go-based CLI tool",
				"Usage:",
				"ggc <command>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RenderMainHelp()
			if err != nil {
				t.Fatalf("RenderMainHelp() should not return error: %v", err)
			}
			for _, want := range tt.expected {
				if !strings.Contains(result, want) {
					t.Errorf("expected help output to contain %q", want)
				}
			}
		})
	}
}

func TestLogoConstants(t *testing.T) {
	if Logo == "" {
		t.Error("Logo constant should not be empty")
	}

	if SmallLogo == "" {
		t.Error("SmallLogo constant should not be empty")
	}

	if !strings.Contains(Logo, "_") {
		t.Error("Logo should contain ASCII art characters")
	}

	if !strings.Contains(SmallLogo, "ggc") {
		t.Error("SmallLogo should contain 'ggc'")
	}
}

func TestRenderCommandHelp(t *testing.T) {
	data := HelpData{
		Usage:       "ggc test [options]",
		Description: "This is a test command",
		Examples:    []string{"ggc test --help", "ggc test --verbose"},
	}

	result, err := RenderCommandHelp(data)

	if err != nil {
		t.Errorf("RenderCommandHelp should not return error: %v", err)
	}

	if result == "" {
		t.Error("RenderCommandHelp should return non-empty string")
	}

	if !strings.Contains(result, data.Usage) {
		t.Error("RenderCommandHelp should contain usage in output")
	}

	if !strings.Contains(result, data.Description) {
		t.Error("RenderCommandHelp should contain description in output")
	}

	for _, example := range data.Examples {
		if !strings.Contains(result, example) {
			t.Errorf("RenderCommandHelp should contain example '%s' in output", example)
		}
	}
}

func TestRenderCommandHelp_EmptyData(t *testing.T) {
	data := HelpData{}

	result, err := RenderCommandHelp(data)

	if err != nil {
		t.Errorf("RenderCommandHelp should not return error for empty data: %v", err)
	}

	if result == "" {
		t.Error("RenderCommandHelp should return non-empty string even for empty data")
	}
}

func TestRenderCommandHelp_NoExamples(t *testing.T) {
	data := HelpData{
		Usage:       "ggc test",
		Description: "Test command",
		Examples:    []string{},
	}

	result, err := RenderCommandHelp(data)

	if err != nil {
		t.Errorf("RenderCommandHelp should not return error: %v", err)
	}

	if !strings.Contains(result, data.Usage) {
		t.Error("RenderCommandHelp should contain usage in output")
	}

	if !strings.Contains(result, data.Description) {
		t.Error("RenderCommandHelp should contain description in output")
	}
}

func TestRenderCommandHelp_SpecialCharacters(t *testing.T) {
	data := HelpData{
		Usage:       "ggc test --option=\"value with spaces\"",
		Description: "Test command with special characters: &<>\"'",
		Examples:    []string{"ggc test --file=\"my file.txt\"", "ggc test --pattern='*.go'"},
	}

	result, err := RenderCommandHelp(data)

	if err != nil {
		t.Errorf("RenderCommandHelp should not return error for special characters: %v", err)
	}

	if !strings.Contains(result, data.Usage) {
		t.Error("RenderCommandHelp should contain usage with special characters")
	}

	if !strings.Contains(result, data.Description) {
		t.Error("RenderCommandHelp should contain description with special characters")
	}

	for _, example := range data.Examples {
		if !strings.Contains(result, example) {
			t.Errorf("RenderCommandHelp should contain example '%s' in output", example)
		}
	}
}

func TestRenderCommandHelp_LongText(t *testing.T) {
	longUsage := strings.Repeat("ggc test --very-long-option-name ", 10)
	longDescription := strings.Repeat("This is a very long description that spans multiple lines and contains a lot of text to test how the template handles long content. ", 5)
	longExample := strings.Repeat("ggc test --example-with-very-long-parameters ", 3)

	data := HelpData{
		Usage:       longUsage,
		Description: longDescription,
		Examples:    []string{longExample},
	}

	result, err := RenderCommandHelp(data)

	if err != nil {
		t.Errorf("RenderCommandHelp should not return error for long text: %v", err)
	}

	if !strings.Contains(result, data.Usage) {
		t.Error("RenderCommandHelp should contain long usage text")
	}

	if !strings.Contains(result, data.Description) {
		t.Error("RenderCommandHelp should contain long description text")
	}

	if !strings.Contains(result, longExample) {
		t.Error("RenderCommandHelp should contain long example text")
	}
}

func TestRenderCommandHelp_MultilineText(t *testing.T) {
	data := HelpData{
		Usage:       "ggc test\n--option value\n--another-option",
		Description: "Test command\nwith multiline\ndescription",
		Examples:    []string{"ggc test\n--verbose", "ggc test\n--quiet"},
	}

	result, err := RenderCommandHelp(data)

	if err != nil {
		t.Errorf("RenderCommandHelp should not return error for multiline text: %v", err)
	}

	if !strings.Contains(result, "ggc test") {
		t.Error("RenderCommandHelp should contain usage text")
	}

	if !strings.Contains(result, "Test command") {
		t.Error("RenderCommandHelp should contain description text")
	}

	if !strings.Contains(result, "--verbose") {
		t.Error("RenderCommandHelp should contain example text")
	}
}

func TestRenderCommandHelp_SingleExample(t *testing.T) {
	data := HelpData{
		Usage:       "ggc single",
		Description: "Command with single example",
		Examples:    []string{"ggc single --help"},
	}

	result, err := RenderCommandHelp(data)

	if err != nil {
		t.Errorf("RenderCommandHelp should not return error for single example: %v", err)
	}

	if !strings.Contains(result, data.Usage) {
		t.Error("RenderCommandHelp should contain usage")
	}

	if !strings.Contains(result, data.Description) {
		t.Error("RenderCommandHelp should contain description")
	}

	if !strings.Contains(result, data.Examples[0]) {
		t.Error("RenderCommandHelp should contain the single example")
	}
}

func TestRenderCommandHelp_ManyExamples(t *testing.T) {
	examples := make([]string, 10)
	for i := 0; i < 10; i++ {
		examples[i] = fmt.Sprintf("ggc test --example-%d", i+1)
	}

	data := HelpData{
		Usage:       "ggc test [options]",
		Description: "Command with many examples",
		Examples:    examples,
	}

	result, err := RenderCommandHelp(data)

	if err != nil {
		t.Errorf("RenderCommandHelp should not return error for many examples: %v", err)
	}

	for i, example := range examples {
		if !strings.Contains(result, example) {
			t.Errorf("RenderCommandHelp should contain example %d: '%s'", i+1, example)
		}
	}
}
