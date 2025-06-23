package cmd

import (
	"bytes"
	"testing"
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
