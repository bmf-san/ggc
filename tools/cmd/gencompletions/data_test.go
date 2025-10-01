package main

import (
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v6/cmd/command"
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
