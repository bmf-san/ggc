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
		t.Errorf("ヘルプ出力が想定と異なります: %s", output)
	}
	if !strings.Contains(output, "ggc add <file>") {
		t.Errorf("コマンド一覧が出力されていません: %s", output)
	}
	if !strings.Contains(output, "Examples:") {
		t.Errorf("Examplesが出力されていません: %s", output)
	}
}
