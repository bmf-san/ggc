package cmd

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/bmf-san/ggc/v8/internal/config"
	"github.com/bmf-san/ggc/v8/internal/testutil"
)

func TestReplaceNamedPlaceholders(t *testing.T) {
	mock := testutil.NewMockGitClient()
	cm := config.NewConfigManager(mock)
	c := &Cmd{
		gitClient:     mock,
		configManager: cm,
		outputWriter:  io.Discard,
	}

	today := time.Now().Format("2006-01-02")

	cases := []struct {
		in   string
		want string
	}{
		{"no placeholders here", "no placeholders here"},
		{"checkout {branch}", "checkout main"},
		{"push {remote} HEAD", "push origin HEAD"},
		{"tag release-{date}", "tag release-" + today},
		{"echo {unknown}", "echo {unknown}"}, // preserved
		{"{branch}-{date}", "main-" + today},
	}

	for _, tc := range cases {
		got := c.replaceNamedPlaceholders(tc.in)
		if got != tc.want {
			t.Errorf("replaceNamedPlaceholders(%q) = %q; want %q", tc.in, got, tc.want)
		}
	}
}

func TestReplaceNamedPlaceholders_FastPath(t *testing.T) {
	c := &Cmd{}
	in := "no braces at all"
	if got := c.replaceNamedPlaceholders(in); got != in {
		t.Errorf("fast path should return input unchanged; got %q", got)
	}
}

func TestReplaceNamedPlaceholders_NilClient(t *testing.T) {
	c := &Cmd{} // no gitClient, no configManager
	got := c.replaceNamedPlaceholders("cmd {branch} {remote}")
	// branch → "", remote → "origin" (fallback)
	if !strings.Contains(got, "origin") {
		t.Errorf("expected fallback 'origin', got %q", got)
	}
}
