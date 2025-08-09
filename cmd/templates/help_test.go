package templates

import (
	"strings"
	"testing"
)

func TestSelectLogo(t *testing.T) {
	// テスト用の関数として selectLogo を public にする必要がないので、
	// RenderMainHelp を通してテストします
	result, err := RenderMainHelp()
	if err != nil {
		t.Fatalf("RenderMainHelp() should not return error: %v", err)
	}

	// ロゴが含まれていることを確認
	hasLogo := strings.Contains(result, "__ _") || strings.Contains(result, "╔═════════════════════╗")
	if !hasLogo {
		t.Error("Result should contain either full logo or small logo")
	}
}

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


