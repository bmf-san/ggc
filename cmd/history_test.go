package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v8/internal/config"
	"github.com/bmf-san/ggc/v8/internal/history"
)

// newCmdForHistoryTest builds a Cmd wired to a buffer for output capture
// and an isolated history store rooted in the test's temp dir, so each
// test gets a deterministic starting state without touching the
// per-user history file.
func newCmdForHistoryTest(t *testing.T) (*Cmd, *bytes.Buffer, *history.Store) {
	t.Helper()
	mockClient := &mockGitClient{}
	cm := config.NewConfigManager(mockClient)
	cmd, err := NewCmd(mockClient, cm)
	if err != nil {
		t.Fatalf("NewCmd: %v", err)
	}
	buf := &bytes.Buffer{}
	cmd.outputWriter = buf

	store := &history.Store{Path: filepath.Join(t.TempDir(), "history.jsonl")}
	prev := history.Default()
	history.SetDefault(store)
	t.Cleanup(func() { history.SetDefault(prev) })
	return cmd, buf, store
}

func TestHistory_NoArgsShowsRecent(t *testing.T) {
	cmd, buf, store := newCmdForHistoryTest(t)
	for _, c := range []string{"status", "commit", "push"} {
		if err := store.Append(c, nil, c); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	cmd.History(nil)

	got := buf.String()
	for _, want := range []string{"status", "commit", "push"} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q: %s", want, got)
		}
	}
}

func TestHistory_NumericShortcut(t *testing.T) {
	cmd, buf, store := newCmdForHistoryTest(t)
	for _, c := range []string{"a", "b", "c", "d"} {
		_ = store.Append(c, nil, c)
	}
	cmd.History([]string{"2"})

	got := buf.String()
	if !strings.Contains(got, "c") || !strings.Contains(got, "d") {
		t.Errorf("expected last 2 (c, d), got: %s", got)
	}
	if strings.Contains(got, "\ta\n") || strings.Contains(got, "\tb\n") {
		t.Errorf("should not include older entries: %s", got)
	}
}

func TestHistory_LastSubcommand(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		seedN    int
		want     string
		wantSkip string
	}{
		{"default count when omitted", []string{"last"}, 3, "cmd-0", ""},
		{"explicit positive", []string{"last", "2"}, 4, "cmd-3", "cmd-1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd, buf, store := newCmdForHistoryTest(t)
			for i := 0; i < tc.seedN; i++ {
				name := "cmd-" + string(rune('0'+i))
				_ = store.Append(name, nil, name)
			}
			cmd.History(tc.args)
			out := buf.String()
			if !strings.Contains(out, tc.want) {
				t.Errorf("want %q in output, got: %s", tc.want, out)
			}
			if tc.wantSkip != "" && strings.Contains(out, "\t"+tc.wantSkip+"\n") {
				t.Errorf("should not include %q, got: %s", tc.wantSkip, out)
			}
		})
	}
}

func TestHistory_LastInvalidArg(t *testing.T) {
	cases := []struct {
		name string
		arg  string
	}{
		{"non-numeric", "abc"},
		{"zero", "0"},
		{"negative", "-1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd, buf, _ := newCmdForHistoryTest(t)
			cmd.History([]string{"last", tc.arg})
			if !strings.Contains(buf.String(), "Error:") {
				t.Errorf("want error for %q, got: %s", tc.arg, buf.String())
			}
		})
	}
}

func TestHistory_Search(t *testing.T) {
	cmd, buf, store := newCmdForHistoryTest(t)
	_ = store.Append("commit", []string{"-m", "feat: add"}, `commit -m "feat: add"`)
	_ = store.Append("checkout", []string{"main"}, "checkout main")
	_ = store.Append("status", nil, "status")

	cmd.History([]string{"search", "commit"})
	out := buf.String()
	if !strings.Contains(out, "commit") {
		t.Errorf("want commit hit, got: %s", out)
	}
	if strings.Contains(out, "\tstatus\n") {
		t.Errorf("should not include status: %s", out)
	}
}

func TestHistory_SearchMultiWord(t *testing.T) {
	cmd, buf, store := newCmdForHistoryTest(t)
	_ = store.Append("commit", []string{"-m", "feat: add"}, `commit -m "feat: add"`)
	_ = store.Append("status", nil, "status")

	cmd.History([]string{"search", "feat:", "add"})
	out := buf.String()
	if !strings.Contains(out, "feat: add") {
		t.Errorf("multi-word search should join args, got: %s", out)
	}
}

func TestHistory_SearchMissingPattern(t *testing.T) {
	cmd, buf, _ := newCmdForHistoryTest(t)
	cmd.History([]string{"search"})
	if !strings.Contains(buf.String(), "Usage:") {
		t.Errorf("want Usage hint, got: %s", buf.String())
	}
}

func TestHistory_Clear(t *testing.T) {
	cmd, buf, store := newCmdForHistoryTest(t)
	_ = store.Append("status", nil, "status")
	cmd.History([]string{"clear"})

	if !strings.Contains(buf.String(), "cleared") {
		t.Errorf("want confirmation, got: %s", buf.String())
	}
	got, err := store.ReadAll()
	if err != nil || len(got) != 0 {
		t.Errorf("history not cleared: %+v %v", got, err)
	}
}

func TestHistory_UnknownSubcommand(t *testing.T) {
	cmd, buf, _ := newCmdForHistoryTest(t)
	cmd.History([]string{"bogus"})
	if !strings.Contains(buf.String(), "Usage:") {
		t.Errorf("want Usage hint, got: %s", buf.String())
	}
}
