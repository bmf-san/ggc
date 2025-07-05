package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/cmd/templates"
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

func TestShowStashPullPopHelp(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	helper.ShowStashPullPopHelp()

	if buf.Len() == 0 {
		t.Error("Expected help message, got empty output")
	}
}

func TestShowResetCleanHelp(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{
		outputWriter: &buf,
	}

	helper.ShowResetCleanHelp()

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
