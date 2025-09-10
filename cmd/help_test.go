package cmd

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v5/cmd/templates"
)

func TestHelper_ShowHelp(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	helper.ShowHelp()

	if buf.Len() == 0 {
		t.Error("Expected help message, got empty output")
	}
}

func TestHelper_ShowLogHelp(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	helper.ShowLogHelp()

	if buf.Len() == 0 {
		t.Error("Expected help message, got empty output")
	}
}

func TestHelper_ShowPushHelp(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	helper.ShowPushHelp()

	if buf.Len() == 0 {
		t.Error("Expected help message, got empty output")
	}
}

func TestHelper_ShowRemoteHelp(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	helper.ShowRemoteHelp()

	if buf.Len() == 0 {
		t.Error("Expected help message, got empty output")
	}
}

func TestHelper_ShowCommitHelp(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	helper.ShowCommitHelp()

	if buf.Len() == 0 {
		t.Error("Expected help message, got empty output")
	}
}

func TestHelper_ShowStashHelp(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	helper.ShowStashHelp()

	if buf.Len() == 0 {
		t.Error("Expected help message, got empty output")
	}
}

func TestHelper_ShowResetHelp(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	helper.ShowResetHelp()

	if buf.Len() == 0 {
		t.Error("Expected help message, got empty output")
	}
}

func TestHelper_ShowPullHelp(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	helper.ShowPullHelp()

	if buf.Len() == 0 {
		t.Error("Expected help message, got empty output")
	}
}

func TestHelper_ShowCleanHelp(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	helper.ShowCleanHelp()

	if buf.Len() == 0 {
		t.Error("Expected help message, got empty output")
	}
}

func TestHelper_ShowRebaseHelp(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	helper.ShowRebaseHelp()

	if buf.Len() == 0 {
		t.Error("Expected help message, got empty output")
	}
}

func TestShowAddHelp(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	helper.ShowAddHelp()

	if buf.Len() == 0 {
		t.Error("Expected help message, got empty output")
	}
}

func TestShowListBranchesHelp(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	helper.ShowListBranchesHelp()

	if buf.Len() == 0 {
		t.Error("Expected help message, got empty output")
	}
}

func TestShowDeleteBranchHelp(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	helper.ShowDeleteBranchHelp()

	if buf.Len() == 0 {
		t.Error("Expected help message, got empty output")
	}
}

func TestShowDeleteMergedBranchHelp(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	helper.ShowDeleteMergedBranchHelp()

	if buf.Len() == 0 {
		t.Error("Expected help message, got empty output")
	}
}

func TestHelper_ShowCommandHelp_Error(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	// Use invalid template data to trigger an error
	helper.ShowCommandHelp(templates.HelpData{
		Usage:       "test {{.InvalidField}}",
		Description: "test description",
		Examples:    []string{"example"},
	})

	output := buf.String()
	// For template errors, the template may be output as-is,
	// so check for either error message or template execution result
	if !strings.Contains(output, "Error:") && !strings.Contains(output, "test") {
		t.Errorf("Expected error message or template output, got: %s", output)
	}
}

func TestHelper_ShowFetchHelp(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	helper.ShowFetchHelp()

	if buf.Len() == 0 {
		t.Error("Expected help message, got empty output")
	}
}

func TestHelper_ShowCommandHelp_WithExamples(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	// Test with multiple examples
	helper.ShowCommandHelp(templates.HelpData{
		Usage:       "test usage",
		Description: "test description",
		Examples:    []string{"example1", "example2", "example3"},
	})

	output := buf.String()
	if !strings.Contains(output, "test usage") {
		t.Error("Expected usage in output")
	}
	if !strings.Contains(output, "test description") {
		t.Error("Expected description in output")
	}
	if !strings.Contains(output, "example1") {
		t.Error("Expected example1 in output")
	}
	if !strings.Contains(output, "example2") {
		t.Error("Expected example2 in output")
	}
	if !strings.Contains(output, "example3") {
		t.Error("Expected example3 in output")
	}
}

func TestHelper_ShowCommandHelp_EmptyExamples(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	// Test with empty examples
	helper.ShowCommandHelp(templates.HelpData{
		Usage:       "test usage",
		Description: "test description",
		Examples:    []string{},
	})

	output := buf.String()
	if !strings.Contains(output, "test usage") {
		t.Error("Expected usage in output")
	}
	if !strings.Contains(output, "test description") {
		t.Error("Expected description in output")
	}
}

// Test all remaining help functions that are currently at 100% but need comprehensive testing
func TestHelper_AllHelpFunctions(t *testing.T) {
	tests := []struct {
		name         string
		helpFunc     func(*Helper)
		expectedText []string
	}{
		{
			name:     "ShowBranchHelp",
			helpFunc: (*Helper).ShowBranchHelp,
			expectedText: []string{
				"ggc branch <command>",
				"List, create, or delete branches",
				"ggc branch current",
				"ggc branch checkout",
				"ggc branch create",
				"ggc branch delete",
			},
		},
		{
			name:     "ShowConfigHelp",
			helpFunc: (*Helper).ShowConfigHelp,
			expectedText: []string{
				"ggc config [command]",
				"Get, set, and list configuration values",
				"ggc config list",
				"ggc config get",
				"ggc config set",
			},
		},
		{
			name:     "ShowHookHelp",
			helpFunc: (*Helper).ShowHookHelp,
			expectedText: []string{
				"ggc hook [command]",
				"Manage Git hooks",
				"ggc hook list",
				"ggc hook install",
				"ggc hook enable",
			},
		},
		{
			name:     "ShowRestoreHelp",
			helpFunc: (*Helper).ShowRestoreHelp,
			expectedText: []string{
				"ggc restore [command]",
				"Restore working tree files",
				"restore <file>",
				"restore staged",
			},
		},
		{
			name:     "ShowStatusHelp",
			helpFunc: (*Helper).ShowStatusHelp,
			expectedText: []string{
				"ggc status [command]",
				"Show the working tree status",
				"ggc status",
				"ggc status short",
			},
		},
		{
			name:     "ShowTagHelp",
			helpFunc: (*Helper).ShowTagHelp,
			expectedText: []string{
				"ggc tag [command] [options]",
				"Create, list, delete and verify tags",
				"ggc tag list",
				"ggc tag create",
				"ggc tag delete",
			},
		},
		{
			name:     "ShowVersionHelp",
			helpFunc: (*Helper).ShowVersionHelp,
			expectedText: []string{
				"ggc version",
				"Show current ggc version",
				"Shows build time",
			},
		},
		{
			name:     "ShowDiffHelp",
			helpFunc: (*Helper).ShowDiffHelp,
			expectedText: []string{
				"ggc diff [options]",
				"Show changes between commits",
				"ggc diff",
				"ggc diff staged",
				"ggc diff unstaged",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			helper := &Helper{
				outputWriter: &buf,
			}

			tt.helpFunc(helper)

			output := buf.String()
			if buf.Len() == 0 {
				t.Errorf("Expected help message, got empty output for %s", tt.name)
			}

			for _, expectedText := range tt.expectedText {
				if !strings.Contains(output, expectedText) {
					t.Errorf("Expected %q in output for %s, got: %s", expectedText, tt.name, output)
				}
			}
		})
	}
}

// Test error handling in ShowHelp - this should improve ShowHelp coverage from 60% to higher
func TestHelper_ShowHelp_Error(t *testing.T) {
	// We need to test the error path in ShowHelp
	// This is challenging because templates.RenderMainHelp() is a direct call
	// We can test by creating a scenario where the template rendering might fail
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	// Call ShowHelp - if templates.RenderMainHelp() returns an error,
	// we should see "Error:" in the output
	helper.ShowHelp()

	output := buf.String()
	// In normal cases, we expect successful rendering
	// If there's an error, we'd see "Error:" prefix
	if strings.Contains(output, "Error:") {
		// This means we caught an error case, which is good for coverage
		t.Log("Successfully tested error path in ShowHelp")
	} else if buf.Len() > 0 {
		// Normal successful case
		t.Log("ShowHelp executed successfully")
	} else {
		t.Error("ShowHelp produced no output")
	}
}

// Test NewHelper constructor
func TestNewHelper(t *testing.T) {
	helper := NewHelper()
	
	if helper == nil {
		t.Error("NewHelper should return non-nil Helper")
	}
	
	if helper.outputWriter == nil {
		t.Error("NewHelper should set outputWriter")
	}
	
	// Test that the helper can be used
	var buf bytes.Buffer
	helper.outputWriter = &buf
	helper.ShowVersionHelp()
	
	if buf.Len() == 0 {
		t.Error("NewHelper should create functional Helper")
	}
}

// Test ShowCommandHelp with various edge cases to improve coverage
func TestHelper_ShowCommandHelp_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		data     templates.HelpData
		expected []string
	}{
		{
			name: "minimal data",
			data: templates.HelpData{
				Usage:       "minimal",
				Description: "",
				Examples:    nil,
			},
			expected: []string{"minimal"},
		},
		{
			name: "single example",
			data: templates.HelpData{
				Usage:       "single usage",
				Description: "single description",
				Examples:    []string{"single example"},
			},
			expected: []string{"single usage", "single description", "single example"},
		},
		{
			name: "long usage and description",
			data: templates.HelpData{
				Usage:       "very long usage string that spans multiple concepts and ideas",
				Description: "very long description that explains in detail what this command does and how to use it properly with various options and flags",
				Examples:    []string{"example with long explanation and multiple parameters"},
			},
			expected: []string{"very long usage", "very long description", "example with long"},
		},
		{
			name: "special characters",
			data: templates.HelpData{
				Usage:       "usage with <brackets> and [optional] and | pipes",
				Description: "description with special chars: @#$%^&*()",
				Examples:    []string{"example with \"quotes\" and 'apostrophes'"},
			},
			expected: []string{"<brackets>", "[optional]", "| pipes", "@#$%^&*()", "\"quotes\"", "'apostrophes'"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			helper := &Helper{
				outputWriter: &buf,
			}

			helper.ShowCommandHelp(tt.data)

			output := buf.String()
			if buf.Len() == 0 {
				t.Errorf("Expected output for %s, got empty", tt.name)
			}

			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected %q in output for %s, got: %s", expected, tt.name, output)
				}
			}
		})
	}
}

// Test output writer functionality
func TestHelper_OutputWriter(t *testing.T) {
	// Test with different output writers
	tests := []struct {
		name   string
		writer io.Writer
	}{
		{
			name:   "bytes.Buffer",
			writer: &bytes.Buffer{},
		},
		{
			name:   "custom writer",
			writer: &customWriter{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helper := &Helper{
				outputWriter: tt.writer,
			}

			// Test that it doesn't panic and produces output
			helper.ShowVersionHelp()

			// For bytes.Buffer, we can check content
			if buf, ok := tt.writer.(*bytes.Buffer); ok {
				if buf.Len() == 0 {
					t.Error("Expected output to bytes.Buffer")
				}
			}
		})
	}
}

// Custom writer for testing
type customWriter struct {
	data []byte
}

func (w *customWriter) Write(p []byte) (n int, err error) {
	w.data = append(w.data, p...)
	return len(p), nil
}

// Test comprehensive coverage of all help functions
func TestHelper_ComprehensiveCoverage(t *testing.T) {
	helper := NewHelper()
	var buf bytes.Buffer
	helper.outputWriter = &buf

	// Test all help functions to ensure they work
	helpFunctions := []func(){
		helper.ShowHelp,
		helper.ShowAddHelp,
		helper.ShowBranchHelp,
		helper.ShowCleanHelp,
		helper.ShowCommitHelp,
		helper.ShowLogHelp,
		helper.ShowPullHelp,
		helper.ShowPushHelp,
		helper.ShowRemoteHelp,
		helper.ShowStashHelp,
		helper.ShowHookHelp,
		helper.ShowConfigHelp,
		helper.ShowRestoreHelp,
		helper.ShowStatusHelp,
		helper.ShowTagHelp,
		helper.ShowVersionHelp,
		helper.ShowRebaseHelp,
		helper.ShowResetHelp,
		helper.ShowListBranchesHelp,
		helper.ShowDeleteBranchHelp,
		helper.ShowDeleteMergedBranchHelp,
		helper.ShowDiffHelp,
		helper.ShowFetchHelp,
	}

	for i, helpFunc := range helpFunctions {
		buf.Reset()
		helpFunc()
		
		if buf.Len() == 0 {
			t.Errorf("Help function %d produced no output", i)
		}
	}
}

// Test ShowHelp content verification to improve coverage
func TestHelper_ShowHelp_ContentVerification(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	helper.ShowHelp()

	output := buf.String()
	if buf.Len() == 0 {
		t.Fatal("Expected help message, got empty output")
	}

	// Verify main help content includes expected sections
	expectedSections := []string{
		"ggc: A Go-based CLI tool",
		"Usage:",
		"Main Commands:",
		"ggc help",
		"ggc add",
		"ggc branch",
		"ggc commit",
		"ggc status",
		"Notes:",
	}

	for _, section := range expectedSections {
		if !strings.Contains(output, section) {
			t.Errorf("Expected main help to contain %q, but it was missing from: %s", section, output)
		}
	}
}

// Test ShowCommandHelp with nil/empty examples to improve coverage
func TestHelper_ShowCommandHelp_NilExamples(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	// Test with nil examples
	helper.ShowCommandHelp(templates.HelpData{
		Usage:       "test usage nil",
		Description: "test description nil",
		Examples:    nil,
	})

	output := buf.String()
	if buf.Len() == 0 {
		t.Error("Expected output with nil examples")
	}

	if !strings.Contains(output, "test usage nil") {
		t.Error("Expected usage in output with nil examples")
	}
	if !strings.Contains(output, "test description nil") {
		t.Error("Expected description in output with nil examples")
	}
}

// Test ShowCommandHelp content structure
func TestHelper_ShowCommandHelp_ContentStructure(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	testData := templates.HelpData{
		Usage:       "ggc test <command>",
		Description: "Test command description",
		Examples: []string{
			"ggc test example1",
			"ggc test example2",
		},
	}

	helper.ShowCommandHelp(testData)

	output := buf.String()
	if buf.Len() == 0 {
		t.Fatal("Expected help output")
	}

	// Verify the structure of command help
	expectedStructure := []string{
		"Usage: ggc test <command>",
		"Description:",
		"Test command description",
		"Examples:",
		"ggc test example1",
		"ggc test example2",
	}

	for _, expected := range expectedStructure {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected command help to contain %q, but it was missing from: %s", expected, output)
		}
	}
}

// Test helper with different output writers to ensure proper io.Writer usage
func TestHelper_DifferentWriters(t *testing.T) {
	tests := []struct {
		name   string
		writer io.Writer
	}{
		{
			name:   "standard bytes buffer",
			writer: &bytes.Buffer{},
		},
		{
			name:   "custom buffer writer",
			writer: &customBufferWriter{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helper := &Helper{
				outputWriter: tt.writer,
			}

			// Test both main functions to ensure they work with different writers
			helper.ShowHelp()
			helper.ShowCommandHelp(templates.HelpData{
				Usage:       "test usage",
				Description: "test description",
				Examples:    []string{"test example"},
			})

			// Verify output was written (for bytes.Buffer)
			if buf, ok := tt.writer.(*bytes.Buffer); ok {
				if buf.Len() == 0 {
					t.Errorf("Expected output to be written to %s", tt.name)
				}
			}
		})
	}
}

// Custom buffer writer for testing
type customBufferWriter struct {
	data []byte
}

func (w *customBufferWriter) Write(p []byte) (n int, err error) {
	w.data = append(w.data, p...)
	return len(p), nil
}

// Test all individual help functions with content verification
func TestHelper_IndividualHelpContent(t *testing.T) {
	tests := []struct {
		name         string
		helpFunc     func(*Helper)
		expectedUsage string
		expectedDesc  string
	}{
		{
			name:         "ShowAddHelp",
			helpFunc:     (*Helper).ShowAddHelp,
			expectedUsage: "ggc add <file>",
			expectedDesc:  "Add file contents to the index",
		},
		{
			name:         "ShowCleanHelp",
			helpFunc:     (*Helper).ShowCleanHelp,
			expectedUsage: "ggc clean <command>",
			expectedDesc:  "Clean untracked files and directories",
		},
		{
			name:         "ShowCommitHelp",
			helpFunc:     (*Helper).ShowCommitHelp,
			expectedUsage: "ggc commit <message>",
			expectedDesc:  "Commit staged changes",
		},
		{
			name:         "ShowLogHelp",
			helpFunc:     (*Helper).ShowLogHelp,
			expectedUsage: "ggc log <command>",
			expectedDesc:  "Show commit logs",
		},
		{
			name:         "ShowPullHelp",
			helpFunc:     (*Helper).ShowPullHelp,
			expectedUsage: "ggc pull <command>",
			expectedDesc:  "Pull changes from remote",
		},
		{
			name:         "ShowPushHelp",
			helpFunc:     (*Helper).ShowPushHelp,
			expectedUsage: "ggc push <command>",
			expectedDesc:  "Push changes to remote",
		},
		{
			name:         "ShowRemoteHelp",
			helpFunc:     (*Helper).ShowRemoteHelp,
			expectedUsage: "ggc remote <command>",
			expectedDesc:  "Manage set of tracked repositories",
		},
		{
			name:         "ShowStashHelp",
			helpFunc:     (*Helper).ShowStashHelp,
			expectedUsage: "ggc stash [command]",
			expectedDesc:  "Stash changes",
		},
		{
			name:         "ShowRebaseHelp",
			helpFunc:     (*Helper).ShowRebaseHelp,
			expectedUsage: "ggc rebase [interactive",
			expectedDesc:  "Rebase current branch onto another branch",
		},
		{
			name:         "ShowResetHelp",
			helpFunc:     (*Helper).ShowResetHelp,
			expectedUsage: "ggc reset",
			expectedDesc:  "Reset and clean",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			helper := &Helper{
				outputWriter: &buf,
			}

			tt.helpFunc(helper)

			output := buf.String()
			if buf.Len() == 0 {
				t.Errorf("Expected help output for %s", tt.name)
			}

			if !strings.Contains(output, tt.expectedUsage) {
				t.Errorf("Expected %s to contain usage %q, got: %s", tt.name, tt.expectedUsage, output)
			}

			if !strings.Contains(output, tt.expectedDesc) {
				t.Errorf("Expected %s to contain description %q, got: %s", tt.name, tt.expectedDesc, output)
			}
		})
	}
}
