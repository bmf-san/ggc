package cmd

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v7/pkg/git"
)

// mockDiffClient implements git.DiffReader with argument capture.
type mockDiffClient struct {
	diffArgs []string
	output   string
	err      error
}

func (m *mockDiffClient) Diff() (string, error)       { return "", nil }
func (m *mockDiffClient) DiffStaged() (string, error) { return "", nil }
func (m *mockDiffClient) DiffHead() (string, error)   { return "", nil }
func (m *mockDiffClient) DiffWith(args []string) (string, error) {
	m.diffArgs = append([]string(nil), args...)
	if m.err != nil {
		return "", m.err
	}
	return m.output, nil
}

var _ git.DiffReader = (*mockDiffClient)(nil)

func newTestDiffer(client git.DiffReader, buf *bytes.Buffer) *Differ {
	helper := NewHelper()
	helper.outputWriter = buf
	return &Differ{gitClient: client, outputWriter: buf, helper: helper}
}

func TestDiffer_Diff_DefaultsToHead(t *testing.T) {
	buf := &bytes.Buffer{}
	mockClient := &mockDiffClient{output: "DIFF"}
	differ := newTestDiffer(mockClient, buf)

	differ.Diff(nil)

	if want := []string{"HEAD"}; !slices.Equal(mockClient.diffArgs, want) {
		t.Fatalf("expected git args %v, got %v", want, mockClient.diffArgs)
	}

	if buf.String() != "DIFF" {
		t.Fatalf("expected diff output to be written, got %q", buf.String())
	}
}

func TestDiffer_Diff_UnstagedMode(t *testing.T) {
	buf := &bytes.Buffer{}
	mockClient := &mockDiffClient{}
	differ := newTestDiffer(mockClient, buf)

	differ.Diff([]string{"unstaged"})

	if len(mockClient.diffArgs) != 0 {
		t.Fatalf("expected no additional args for unstaged diff, got %v", mockClient.diffArgs)
	}
}

func TestDiffer_Diff_StagedWithPath(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(filePath, []byte("content"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	buf := &bytes.Buffer{}
	mockClient := &mockDiffClient{}
	differ := newTestDiffer(mockClient, buf)

	differ.Diff([]string{"staged", filePath})

	want := []string{"--staged", "--", filePath}
	if !slices.Equal(mockClient.diffArgs, want) {
		t.Fatalf("expected git args %v, got %v", want, mockClient.diffArgs)
	}
}

func TestDiffer_Diff_WithFlagsAndPath(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(filePath, []byte("content"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	buf := &bytes.Buffer{}
	mockClient := &mockDiffClient{}
	differ := newTestDiffer(mockClient, buf)

	differ.Diff([]string{"--stat", filePath})

	want := []string{"--stat", "HEAD", "--", filePath}
	if !slices.Equal(mockClient.diffArgs, want) {
		t.Fatalf("expected git args %v, got %v", want, mockClient.diffArgs)
	}
}

func TestDiffer_Diff_CommitRange(t *testing.T) {
	buf := &bytes.Buffer{}
	mockClient := &mockDiffClient{}
	differ := newTestDiffer(mockClient, buf)

	differ.Diff([]string{"abc123", "def456"})

	want := []string{"abc123", "def456"}
	if !slices.Equal(mockClient.diffArgs, want) {
		t.Fatalf("expected git args %v, got %v", want, mockClient.diffArgs)
	}
}

func TestDiffer_Diff_CommitWithPath(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(filePath, []byte("content"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	buf := &bytes.Buffer{}
	mockClient := &mockDiffClient{}
	differ := newTestDiffer(mockClient, buf)

	differ.Diff([]string{"abc123", filePath})

	want := []string{"abc123", "--", filePath}
	if !slices.Equal(mockClient.diffArgs, want) {
		t.Fatalf("expected git args %v, got %v", want, mockClient.diffArgs)
	}
}

func TestDiffer_Diff_DoubleDashForPaths(t *testing.T) {
	buf := &bytes.Buffer{}
	mockClient := &mockDiffClient{}
	differ := newTestDiffer(mockClient, buf)

	differ.Diff([]string{"--", "deleted.txt"})

	want := []string{"HEAD", "--", "deleted.txt"}
	if !slices.Equal(mockClient.diffArgs, want) {
		t.Fatalf("expected git args %v, got %v", want, mockClient.diffArgs)
	}
}

func TestDiffer_Diff_InvalidFlagShowsHelp(t *testing.T) {
	buf := &bytes.Buffer{}
	mockClient := &mockDiffClient{}
	differ := newTestDiffer(mockClient, buf)

	differ.Diff([]string{"--unknown"})

	output := buf.String()
	if !strings.Contains(output, "unknown flag") {
		t.Fatalf("expected unknown flag message, got %q", output)
	}
	if !strings.Contains(output, "Usage:") {
		t.Fatalf("expected help output, got %q", output)
	}
}

func TestDiffer_Diff_ConflictingFlags(t *testing.T) {
	buf := &bytes.Buffer{}
	mockClient := &mockDiffClient{}
	differ := newTestDiffer(mockClient, buf)

	differ.Diff([]string{"--name-only", "--name-status"})

	output := buf.String()
	if !strings.Contains(output, "cannot be combined") {
		t.Fatalf("expected conflict message, got %q", output)
	}
}

func TestDiffer_Diff_StagedRejectsCommits(t *testing.T) {
	buf := &bytes.Buffer{}
	mockClient := &mockDiffClient{}
	differ := newTestDiffer(mockClient, buf)

	differ.Diff([]string{"staged", "abc123"})

	output := buf.String()
	if !strings.Contains(output, "does not accept commit arguments") {
		t.Fatalf("expected staged error, got %q", output)
	}
}

func TestDiffer_Diff_HeadModeAllowsPaths(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(filePath, []byte("content"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	buf := &bytes.Buffer{}
	mockClient := &mockDiffClient{}
	differ := newTestDiffer(mockClient, buf)

	differ.Diff([]string{"head", filePath})

	want := []string{"HEAD", "--", filePath}
	if !slices.Equal(mockClient.diffArgs, want) {
		t.Fatalf("expected git args %v, got %v", want, mockClient.diffArgs)
	}
}

func TestDiffer_Diff_GitErrorSurfaced(t *testing.T) {
	buf := &bytes.Buffer{}
	mockClient := &mockDiffClient{err: errors.New("git diff failed")}
	differ := newTestDiffer(mockClient, buf)

	differ.Diff(nil)

	if !strings.Contains(buf.String(), "git diff failed") {
		t.Fatalf("expected git error in output, got %q", buf.String())
	}
}

func TestParseDiffArgs_PathClassification(t *testing.T) {
	fakeExists := func(target string) func(string) bool {
		return func(path string) bool { return path == target }
	}

	opts, err := parseDiffArgs([]string{"--stat", "file.txt"}, fakeExists("file.txt"))
	if err != nil {
		t.Fatalf("parseDiffArgs returned error: %v", err)
	}

	if !opts.stat {
		t.Fatalf("expected --stat to be set")
	}
	if want := []string{"file.txt"}; !slices.Equal(opts.paths, want) {
		t.Fatalf("expected paths %v, got %v", want, opts.paths)
	}
	if len(opts.commits) != 0 {
		t.Fatalf("expected no commits, got %v", opts.commits)
	}
}

func TestParseDiffArgs_CommitDetection(t *testing.T) {
	opts, err := parseDiffArgs([]string{"abc123"}, func(string) bool { return false })
	if err != nil {
		t.Fatalf("parseDiffArgs returned error: %v", err)
	}

	if want := []string{"abc123"}; !slices.Equal(opts.commits, want) {
		t.Fatalf("expected commits %v, got %v", want, opts.commits)
	}
}

func TestParseDiffArgs_TooManyCommits(t *testing.T) {
	_, err := parseDiffArgs([]string{"a", "b", "c"}, func(string) bool { return false })
	if err == nil || !strings.Contains(err.Error(), "too many commit arguments") {
		t.Fatalf("expected too many commit arguments error, got %v", err)
	}
}

func TestParseDiffArgs_NameStatusFlag(t *testing.T) {
	opts, err := parseDiffArgs([]string{"--name-status"}, func(string) bool { return false })
	if err != nil {
		t.Fatalf("parseDiffArgs returned error: %v", err)
	}

	if !opts.nameStatus {
		t.Fatalf("expected name-status flag to be set")
	}
	if opts.nameOnly {
		t.Fatalf("did not expect name-only to be set")
	}
}

func TestParseDiffArgs_RepeatedDoubleDash(t *testing.T) {
	opts, err := parseDiffArgs([]string{"--", "--", "file"}, func(string) bool { return false })
	if err != nil {
		t.Fatalf("parseDiffArgs returned error: %v", err)
	}

	if want := []string{"--", "file"}; !slices.Equal(opts.paths, want) {
		t.Fatalf("expected paths %v, got %v", want, opts.paths)
	}
}

func TestDiffMode_String(t *testing.T) {
	cases := []struct {
		mode diffMode
		want string
	}{
		{diffModeDefault, "default"},
		{diffModeUnstaged, "unstaged"},
		{diffModeStaged, "staged"},
		{diffModeHead, "head"},
	}

	for _, tc := range cases {
		if got := tc.mode.String(); got != tc.want {
			t.Fatalf("mode %v: expected %q, got %q", tc.mode, tc.want, got)
		}
	}
}

func TestMapMode(t *testing.T) {
	cases := map[string]diffMode{
		"unstaged": diffModeUnstaged,
		"staged":   diffModeStaged,
		"head":     diffModeHead,
		"":         diffModeDefault,
		"other":    diffModeDefault,
	}

	for input, want := range cases {
		if got := mapMode(input); got != want {
			t.Fatalf("mapMode(%q) expected %v, got %v", input, want, got)
		}
	}
}

func TestModeArg(t *testing.T) {
	if got := modeArg(&diffOptions{mode: diffModeStaged}); !slices.Equal(got, []string{"--staged"}) {
		t.Fatalf("expected staged mode arg, got %v", got)
	}

	if got := modeArg(&diffOptions{mode: diffModeUnstaged}); len(got) != 0 {
		t.Fatalf("expected no args for unstaged/default modes, got %v", got)
	}
}

func TestSummaryArgs(t *testing.T) {
	cases := []struct {
		name string
		opts diffOptions
		want []string
	}{
		{name: "stat only", opts: diffOptions{stat: true}, want: []string{"--stat"}},
		{name: "name only", opts: diffOptions{nameOnly: true}, want: []string{"--name-only"}},
		{name: "name status", opts: diffOptions{nameStatus: true}, want: []string{"--name-status"}},
		{name: "combo", opts: diffOptions{stat: true, nameOnly: true}, want: []string{"--stat", "--name-only"}},
	}

	for _, tc := range cases {
		if got := summaryArgs(&tc.opts); !slices.Equal(got, tc.want) {
			t.Fatalf("%s: expected %v, got %v", tc.name, tc.want, got)
		}
	}
}

func TestCommitArgs(t *testing.T) {
	cases := []struct {
		name string
		opts diffOptions
		want []string
	}{
		{name: "explicit commits", opts: diffOptions{commits: []string{"a", "b"}}, want: []string{"a", "b"}},
		{name: "default fallback", opts: diffOptions{mode: diffModeDefault}, want: []string{"HEAD"}},
		{name: "head fallback", opts: diffOptions{mode: diffModeHead}, want: []string{"HEAD"}},
		{name: "staged no commits", opts: diffOptions{mode: diffModeStaged}, want: nil},
	}

	for _, tc := range cases {
		if got := commitArgs(&tc.opts); !slices.Equal(got, tc.want) {
			t.Fatalf("%s: expected %v, got %v", tc.name, tc.want, got)
		}
	}
}

func TestDefaultPathExists(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "exists.txt")
	if err := os.WriteFile(file, []byte("content"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	cases := []struct {
		name string
		path string
		want bool
	}{
		{name: "empty", path: "", want: false},
		{name: "flag", path: "--name-only", want: false},
		{name: "missing", path: filepath.Join(dir, "missing.txt"), want: false},
		{name: "exists", path: file, want: true},
	}

	for _, tc := range cases {
		if got := defaultPathExists(tc.path); got != tc.want {
			t.Fatalf("%s: expected %v, got %v", tc.name, tc.want, got)
		}
	}
}

func TestDiffer_Diff_NameStatusFlag(t *testing.T) {
	buf := &bytes.Buffer{}
	mockClient := &mockDiffClient{}
	differ := newTestDiffer(mockClient, buf)

	differ.Diff([]string{"--name-status"})

	if want := []string{"--name-status", "HEAD"}; !slices.Equal(mockClient.diffArgs, want) {
		t.Fatalf("expected git args %v, got %v", want, mockClient.diffArgs)
	}
}
