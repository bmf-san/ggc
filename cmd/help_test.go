package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestHelper_ShowHelp_Output(t *testing.T) {
	var buf bytes.Buffer
	helper := &Helper{writer: &buf}

	helper.ShowHelp()

	output := buf.String()
	if !strings.Contains(output, "ggc: A Go-based CLI tool to streamline Git operations") {
		t.Errorf("help output is not what was expected: %s", output)
	}
	if !strings.Contains(output, "ggc add <file>") {
		t.Errorf("command list is not displayed: %s", output)
	}
	if !strings.Contains(output, "Examples:") {
		t.Errorf("examples are not displayed: %s", output)
	}
}
