package main

import (
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v8/cmd/command"
)

func TestBuildTemplateData(t *testing.T) {
	input := []command.Info{
		{
			Name:    "branch",
			Summary: "Branch operations",
			Subcommands: []command.SubcommandInfo{
				{Name: "branch checkout remote", Summary: "Checkout branch"},
				{Name: "branch delete merged", Summary: "Delete branches"},
				{Name: "branch list local", Summary: "List local"},
				{Name: "branch list remote", Summary: "List remote"},
				{Name: "branch rename <old> <new>", Summary: "Rename branch"},
			},
		},
		{
			Name:    "add",
			Summary: "Add files",
			Subcommands: []command.SubcommandInfo{
				{Name: "add interactive", Summary: "Interactive add"},
				{Name: "add patch", Summary: "Patch add"},
				{Name: "add <file>", Summary: "File add"},
			},
		},
		{
			Name:   "hidden",
			Hidden: true,
		},
	}

	data := buildTemplateData(input)

	if data == nil {
		t.Fatalf("buildTemplateData returned nil")
	}

	if len(data.TopLevel) != 2 {
		t.Fatalf("expected 2 top-level commands, got %d", len(data.TopLevel))
	}

	if data.TopLevelList != "add branch" {
		t.Fatalf("unexpected top level list: %q", data.TopLevelList)
	}

	branch := data.CommandMap["branch"]
	if branch == nil {
		t.Fatalf("branch command missing in map")
	}

	if !branch.IncludeInCase {
		t.Fatalf("branch should be included in bash case list")
	}

	if !strings.Contains(branch.SubcommandList, "checkout") {
		t.Fatalf("branch subcommand list missing checkout: %q", branch.SubcommandList)
	}

	if branch.Subcommand("checkout") == nil {
		t.Fatalf("expected checkout subcommand present")
	}

	checkout := branch.Subcommand("checkout")
	if checkout == nil {
		t.Fatalf("checkout subcommand lookup failed")
	}

	if checkout.KeywordList != "remote" {
		t.Fatalf("expected checkout keyword list to be remote, got %q", checkout.KeywordList)
	}

	// checkout keywords should be skipped for generic keyword generation
	for _, sub := range branch.KeywordSubcommands {
		if sub.Name == "checkout" {
			t.Fatalf("checkout should not appear in keyword subcommands")
		}
	}

	deleteSub := branch.Subcommand("delete")
	if deleteSub == nil {
		t.Fatalf("delete subcommand not found")
	}
	if deleteSub.KeywordList != "merged" {
		t.Fatalf("expected delete keyword 'merged', got %q", deleteSub.KeywordList)
	}

	listSub := branch.Subcommand("list")
	if listSub == nil {
		t.Fatalf("list subcommand not found")
	}
	if listSub.KeywordList != "local remote" {
		t.Fatalf("expected list keywords 'local remote', got %q", listSub.KeywordList)
	}

	add := data.CommandMap["add"]
	if add == nil {
		t.Fatalf("add command missing in map")
	}

	if add.IncludeInCase {
		t.Fatalf("add should be excluded from bash case list")
	}

	if add.SubcommandList != "interactive patch" {
		t.Fatalf("unexpected add subcommand list: %q", add.SubcommandList)
	}

	if data.BranchCheckoutKeywordList != "remote" {
		t.Fatalf("expected branch checkout keyword list 'remote', got %q", data.BranchCheckoutKeywordList)
	}
}

func TestEscapeFunctions(t *testing.T) {
	if got := escapeZsh("it's"); got != "it'\\''s" {
		t.Errorf("escapeZsh = %q, want %q", got, "it'\\''s")
	}
	if got := escapeFish(`say "hi"`); got != `say \"hi\"` {
		t.Errorf("escapeFish = %q, want %q", got, `say \"hi\"`)
	}
	if got := escapeBash(`say "hi"`); got != `say \"hi\"` {
		t.Errorf("escapeBash = %q, want %q", got, `say \"hi\"`)
	}
}

func TestJoin(t *testing.T) {
	if got := join([]string{"a", "b", "c"}, " "); got != "a b c" {
		t.Errorf("join = %q", got)
	}
	if got := join(nil, " "); got != "" {
		t.Errorf("join nil = %q", got)
	}
}

func TestContainsFold(t *testing.T) {
	if !containsFold([]string{"Foo", "bar"}, "foo") {
		t.Error("containsFold should find case-insensitive match")
	}
	if containsFold([]string{"foo", "bar"}, "baz") {
		t.Error("containsFold should return false for missing value")
	}
	if containsFold(nil, "foo") {
		t.Error("containsFold nil should return false")
	}
}

func TestHasKeywords(t *testing.T) {
	noKw := []SubcommandData{{Name: "a"}, {Name: "b"}}
	if hasKeywords(noKw) {
		t.Error("hasKeywords should be false when no keywords")
	}
	withKw := []SubcommandData{{Name: "a", Keywords: []string{"kw"}}}
	if !hasKeywords(withKw) {
		t.Error("hasKeywords should be true when keywords present")
	}
	skipped := []SubcommandData{{Name: "a", Keywords: []string{"kw"}, SkipKeyword: true}}
	if hasKeywords(skipped) {
		t.Error("hasKeywords should be false when SkipKeyword=true")
	}
}

func TestFallbackSummary(t *testing.T) {
	if got := fallbackSummary("", "  name  "); got != "name" {
		t.Errorf("fallbackSummary empty summary = %q", got)
	}
	if got := fallbackSummary("desc", "name"); got != "desc" {
		t.Errorf("fallbackSummary with summary = %q", got)
	}
}

func TestSanitizeSubcommandName(t *testing.T) {
	tests := []struct{ in, want string }{
		{"branch", "branch"},
		{"", ""},
		{".", ""},
		{"<file>", ""},
		{"[opt]", ""},
		{"(opt)", ""},
		{"  trim  ", "trim"},
	}
	for _, tt := range tests {
		if got := sanitizeSubcommandName(tt.in); got != tt.want {
			t.Errorf("sanitizeSubcommandName(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestSanitizeKeyword(t *testing.T) {
	tests := []struct{ in, want string }{
		{"push", "push"},
		{"", ""},
		{"<arg>", ""},
		{"a|b", ""},
		{"[opt]", ""},
		{"(p)", ""},
		{".", ""},
		{"..", ""},
		{",push,", "push"},
	}
	for _, tt := range tests {
		if got := sanitizeKeyword(tt.in); got != tt.want {
			t.Errorf("sanitizeKeyword(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestSubcommand_NotFound(t *testing.T) {
	cmd := &CommandData{
		Subcommands: []SubcommandData{{Name: "create"}},
	}
	if got := cmd.Subcommand("nonexistent"); got != nil {
		t.Errorf("Subcommand('nonexistent') = %v, want nil", got)
	}
}

func TestSubcommandBy(t *testing.T) {
	cmd := &CommandData{
		Subcommands: []SubcommandData{{Name: "delete"}},
	}
	if got := subcommandBy(cmd, "delete"); got == nil || got.Name != "delete" {
		t.Errorf("subcommandBy('delete') = %v, want a SubcommandData with Name='delete'", got)
	}
	if got := subcommandBy(cmd, "missing"); got != nil {
		t.Errorf("subcommandBy('missing') = %v, want nil", got)
	}
}

func TestNeedsHandler(t *testing.T) {
	// Has subcommands → true
	cmd := &CommandData{
		Name:        "branch",
		Subcommands: []SubcommandData{{Name: "create"}},
	}
	if !needsHandler(cmd) {
		t.Error("needsHandler with subcommands should return true")
	}

	// Named command (cmdBranch) → true
	bare := &CommandData{Name: cmdBranch}
	if !needsHandler(bare) {
		t.Errorf("needsHandler for %q should return true", cmdBranch)
	}

	// No subcommands, not a special name → false
	other := &CommandData{Name: "log"}
	if needsHandler(other) {
		t.Error("needsHandler for plain command should return false")
	}
}

func TestSplitSubcommandTokens(t *testing.T) {
	cmd := &command.Info{Name: "branch", Aliases: []string{"br"}}

	// Hidden → false
	hidden := &command.SubcommandInfo{Name: "branch create", Hidden: true}
	if _, _, ok := splitSubcommandTokens(cmd, hidden); ok {
		t.Error("expected false for hidden subcommand")
	}

	// Single token → false
	single := &command.SubcommandInfo{Name: "branch"}
	if _, _, ok := splitSubcommandTokens(cmd, single); ok {
		t.Error("expected false for single-token name")
	}

	// Mismatched command name → false
	mismatch := &command.SubcommandInfo{Name: "tag create"}
	if _, _, ok := splitSubcommandTokens(cmd, mismatch); ok {
		t.Error("expected false for mismatched command name")
	}

	// Matches via alias → success
	alias := &command.SubcommandInfo{Name: "br create"}
	name, extra, ok := splitSubcommandTokens(cmd, alias)
	if !ok {
		t.Error("expected true for alias match")
	}
	if name != "create" {
		t.Errorf("name = %q, want 'create'", name)
	}
	if len(extra) != 0 {
		t.Errorf("extra = %v, want empty", extra)
	}

	// Empty sanitized name → false
	badName := &command.SubcommandInfo{Name: "branch <file>"}
	if _, _, ok := splitSubcommandTokens(cmd, badName); ok {
		t.Error("expected false for unsanitizable subcommand name")
	}
}

func TestAddKeyword_Deduplication(t *testing.T) {
	s := &SubcommandData{}
	s.addKeyword("push")
	s.addKeyword("push") // duplicate → should not add
	if len(s.Keywords) != 1 {
		t.Errorf("addKeyword dedup: len = %d, want 1", len(s.Keywords))
	}
}

func TestShouldIncludeInCase(t *testing.T) {
	// No subcommands → false
	if shouldIncludeInCase("branch", false) {
		t.Error("shouldIncludeInCase(hasSubcommands=false) should return false")
	}
	// cmdAdd with subcommands → false
	if shouldIncludeInCase(cmdAdd, true) {
		t.Error("shouldIncludeInCase(cmdAdd) should return false")
	}
	// Other command with subcommands → true
	if !shouldIncludeInCase("branch", true) {
		t.Error("shouldIncludeInCase(branch, true) should return true")
	}
}

func TestEnsureSubcommandEntry_UpdatesSummary(t *testing.T) {
	subMap := map[string]*SubcommandData{
		"create": {Name: "create", Summary: ""},
	}
	entry := ensureSubcommandEntry(subMap, "create", "Create a new branch")
	if entry.Summary != "Create a new branch" {
		t.Errorf("expected summary to be updated, got %q", entry.Summary)
	}
}
