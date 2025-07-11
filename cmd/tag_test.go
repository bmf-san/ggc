package cmd

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestTagger_Tag(t *testing.T) {
	cases := []struct {
		name           string
		args           []string
		expectedCmds   []string
		mockOutput     []byte
		mockError      error
		expectedOutput string
	}{
		{
			name:           "tag no args",
			args:           []string{},
			expectedCmds:   []string{"git tag"},
			mockOutput:     []byte("v1.0.0\nv1.1.0\nv2.0.0\n"),
			mockError:      nil,
			expectedOutput: "v1.0.0",
		},
		{
			name:           "tag list",
			args:           []string{"list"},
			expectedCmds:   []string{"git tag --sort=-version:refname"},
			mockOutput:     []byte("v2.0.0\nv1.1.0\nv1.0.0\n"),
			mockError:      nil,
			expectedOutput: "v2.0.0",
		},
		{
			name:           "tag list with pattern",
			args:           []string{"list", "v1.*"},
			expectedCmds:   []string{"git tag --sort=-version:refname -l v1.*"},
			mockOutput:     []byte("v1.1.0\nv1.0.0\n"),
			mockError:      nil,
			expectedOutput: "v1.1.0",
		},
		{
			name:           "tag create",
			args:           []string{"create", "v1.2.0"},
			expectedCmds:   []string{"git tag v1.2.0"},
			mockOutput:     []byte(""),
			mockError:      nil,
			expectedOutput: "Tag 'v1.2.0' created successfully",
		},
		{
			name:           "tag create with commit",
			args:           []string{"create", "v1.2.0", "abc123"},
			expectedCmds:   []string{"git tag v1.2.0 abc123"},
			mockOutput:     []byte(""),
			mockError:      nil,
			expectedOutput: "Tag 'v1.2.0' created successfully",
		},
		{
			name:           "tag delete",
			args:           []string{"delete", "v1.0.0"},
			expectedCmds:   []string{"git tag -d v1.0.0"},
			mockOutput:     []byte(""),
			mockError:      nil,
			expectedOutput: "Tag 'v1.0.0' deleted successfully",
		},
		{
			name:           "tag delete multiple",
			args:           []string{"delete", "v1.0.0", "v1.1.0"},
			expectedCmds:   []string{"git tag -d v1.0.0", "git tag -d v1.1.0"},
			mockOutput:     []byte(""),
			mockError:      nil,
			expectedOutput: "Tag 'v1.0.0' deleted successfully",
		},
		{
			name:           "tag annotated",
			args:           []string{"annotated", "v1.2.0", "Release", "version", "1.2.0"},
			expectedCmds:   []string{"git tag -a v1.2.0 -m Release version 1.2.0"},
			mockOutput:     []byte(""),
			mockError:      nil,
			expectedOutput: "Annotated tag 'v1.2.0' created successfully",
		},
		{
			name:           "tag annotated no message",
			args:           []string{"annotated", "v1.2.0"},
			expectedCmds:   []string{"git tag -a v1.2.0"},
			mockOutput:     []byte(""),
			mockError:      nil,
			expectedOutput: "Annotated tag 'v1.2.0' created successfully",
		},
		{
			name:           "tag push all",
			args:           []string{"push"},
			expectedCmds:   []string{"git push origin --tags"},
			mockOutput:     []byte(""),
			mockError:      nil,
			expectedOutput: "Tags pushed successfully",
		},
		{
			name:           "tag push specific",
			args:           []string{"push", "v1.2.0"},
			expectedCmds:   []string{"git push origin v1.2.0"},
			mockOutput:     []byte(""),
			mockError:      nil,
			expectedOutput: "Tags pushed successfully",
		},
		{
			name:           "tag push with remote",
			args:           []string{"push", "v1.2.0", "upstream"},
			expectedCmds:   []string{"git push upstream v1.2.0"},
			mockOutput:     []byte(""),
			mockError:      nil,
			expectedOutput: "Tags pushed successfully",
		},
		{
			name:           "tag show",
			args:           []string{"show", "v1.2.0"},
			expectedCmds:   []string{"git show v1.2.0"},
			mockOutput:     []byte("tag v1.2.0\ncommit abc123\n"),
			mockError:      nil,
			expectedOutput: "tag v1.2.0",
		},
		{
			name:           "invalid command",
			args:           []string{"invalid"},
			expectedCmds:   nil,
			mockOutput:     nil,
			mockError:      nil,
			expectedOutput: "Usage: ggc tag",
		},
		{
			name:           "create tag no name",
			args:           []string{"create"},
			expectedCmds:   nil,
			mockOutput:     nil,
			mockError:      nil,
			expectedOutput: "Error: tag name is required",
		},
		{
			name:           "delete tag no name",
			args:           []string{"delete"},
			expectedCmds:   nil,
			mockOutput:     nil,
			mockError:      nil,
			expectedOutput: "Error: tag name(s) required",
		},
		{
			name:           "show tag no name",
			args:           []string{"show"},
			expectedCmds:   nil,
			mockOutput:     nil,
			mockError:      nil,
			expectedOutput: "Error: tag name is required",
		},
		{
			name:           "annotated tag no name",
			args:           []string{"annotated"},
			expectedCmds:   nil,
			mockOutput:     nil,
			mockError:      nil,
			expectedOutput: "Error: tag name is required",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			cmdIndex := 0

			tagger := &Tagger{
				outputWriter: &buf,
				helper:       NewHelper(),
				execCommand: func(_ string, args ...string) *exec.Cmd {
					if cmdIndex < len(tc.expectedCmds) {
						gotCmd := strings.Join(append([]string{"git"}, args...), " ")
						if gotCmd != tc.expectedCmds[cmdIndex] {
							t.Errorf("expected command %q, got %q", tc.expectedCmds[cmdIndex], gotCmd)
						}
					}
					cmdIndex++

					if tc.mockError != nil {
						return exec.Command("false")
					}
					return exec.Command("echo", string(tc.mockOutput))
				},
			}
			tagger.helper.outputWriter = &buf

			tagger.Tag(tc.args)

			output := buf.String()
			if !strings.Contains(output, tc.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tc.expectedOutput, output)
			}
		})
	}
}

func TestTagger_GetLatestTag(t *testing.T) {
	cases := []struct {
		name           string
		mockOutput     []byte
		mockError      error
		expectedResult string
		expectedError  bool
	}{
		{
			name:           "get latest tag success",
			mockOutput:     []byte("v1.2.0\n"),
			mockError:      nil,
			expectedResult: "v1.2.0",
			expectedError:  false,
		},
		{
			name:           "get latest tag error",
			mockOutput:     nil,
			mockError:      exec.ErrNotFound,
			expectedResult: "",
			expectedError:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tagger := &Tagger{
				execCommand: func(_ string, args ...string) *exec.Cmd {
					expectedCmd := "git describe --tags --abbrev=0"
					gotCmd := strings.Join(append([]string{"git"}, args...), " ")
					if gotCmd != expectedCmd {
						t.Errorf("expected command %q, got %q", expectedCmd, gotCmd)
					}

					if tc.mockError != nil {
						return exec.Command("false")
					}
					return exec.Command("echo", string(tc.mockOutput))
				},
			}

			result, err := tagger.GetLatestTag()

			if tc.expectedError && err == nil {
				t.Error("expected error but got nil")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tc.expectedResult {
				t.Errorf("expected result %q, got %q", tc.expectedResult, result)
			}
		})
	}
}

func TestTagger_TagExists(t *testing.T) {
	cases := []struct {
		name           string
		tagName        string
		mockOutput     []byte
		mockError      error
		expectedResult bool
	}{
		{
			name:           "tag exists",
			tagName:        "v1.2.0",
			mockOutput:     []byte("v1.2.0\n"),
			mockError:      nil,
			expectedResult: true,
		},
		{
			name:           "tag does not exist",
			tagName:        "v1.2.0",
			mockOutput:     []byte(""),
			mockError:      nil,
			expectedResult: false,
		},
		{
			name:           "command error",
			tagName:        "v1.2.0",
			mockOutput:     nil,
			mockError:      exec.ErrNotFound,
			expectedResult: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tagger := &Tagger{
				execCommand: func(_ string, args ...string) *exec.Cmd {
					expectedCmd := "git tag -l " + tc.tagName
					gotCmd := strings.Join(append([]string{"git"}, args...), " ")
					if gotCmd != expectedCmd {
						t.Errorf("expected command %q, got %q", expectedCmd, gotCmd)
					}

					if tc.mockError != nil {
						return exec.Command("false")
					}
					return exec.Command("echo", string(tc.mockOutput))
				},
			}

			result := tagger.TagExists(tc.tagName)

			if result != tc.expectedResult {
				t.Errorf("expected result %t, got %t", tc.expectedResult, result)
			}
		})
	}
}

func TestTagger_GetTagCommit(t *testing.T) {
	cases := []struct {
		name           string
		tagName        string
		mockOutput     []byte
		mockError      error
		expectedResult string
		expectedError  bool
	}{
		{
			name:           "get tag commit success",
			tagName:        "v1.2.0",
			mockOutput:     []byte("abc123def456\n"),
			mockError:      nil,
			expectedResult: "abc123def456",
			expectedError:  false,
		},
		{
			name:           "get tag commit error",
			tagName:        "v1.2.0",
			mockOutput:     nil,
			mockError:      exec.ErrNotFound,
			expectedResult: "",
			expectedError:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tagger := &Tagger{
				execCommand: func(_ string, args ...string) *exec.Cmd {
					expectedCmd := "git rev-list -n 1 " + tc.tagName
					gotCmd := strings.Join(append([]string{"git"}, args...), " ")
					if gotCmd != expectedCmd {
						t.Errorf("expected command %q, got %q", expectedCmd, gotCmd)
					}

					if tc.mockError != nil {
						return exec.Command("false")
					}
					return exec.Command("echo", string(tc.mockOutput))
				},
			}

			result, err := tagger.GetTagCommit(tc.tagName)

			if tc.expectedError && err == nil {
				t.Error("expected error but got nil")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tc.expectedResult {
				t.Errorf("expected result %q, got %q", tc.expectedResult, result)
			}
		})
	}
}
