package command

import (
	"reflect"
	"testing"
)

func TestNewRegistry_Validate(t *testing.T) {
	t.Parallel()
	reg := NewRegistry()
	if err := reg.Validate(); err != nil {
		t.Fatalf("NewRegistry().Validate() returned error: %v", err)
	}
}

func TestValidate_DuplicateCommand(t *testing.T) {
	t.Parallel()
	commands := []Info{
		{Name: "dup", Summary: "first"},
		{Name: "dup", Summary: "second"},
	}

	if err := Validate(commands); err == nil {
		t.Fatalf("expected duplicate command validation to fail")
	}
}

func TestValidate_MissingSummary(t *testing.T) {
	t.Parallel()
	commands := []Info{{Name: "test"}}
	if err := Validate(commands); err == nil {
		t.Fatalf("expected missing summary validation failure")
	}
}

func TestValidate_DuplicateSubcommand(t *testing.T) {
	t.Parallel()
	commands := []Info{
		{
			Name:    "test",
			Summary: "ok",
			Subcommands: []SubcommandInfo{
				{Name: "sub", Summary: "one"},
				{Name: "sub", Summary: "two"},
			},
		},
	}

	if err := Validate(commands); err == nil {
		t.Fatalf("expected duplicate subcommand validation to fail")
	}
}

func TestRegistry_All(t *testing.T) {
	t.Parallel()
	reg := NewRegistryWith([]Info{
		{Name: "cmd1", Summary: "summary one"},
		{Name: "cmd2", Summary: "summary two"},
	})

	cmds := reg.All()
	if len(cmds) != 2 {
		t.Fatalf("expected 2 commands, got %d", len(cmds))
	}

	// Verify defensive copy
	cmds[0].Name = "mutated"
	original := reg.All()
	if original[0].Name == "mutated" {
		t.Fatalf("mutating All() result modified registry")
	}
}

func TestRegistry_Find(t *testing.T) {
	t.Parallel()
	reg := NewRegistryWith([]Info{
		{Name: "help", Summary: "show help"},
		{Name: "other", Summary: "some other cmd"},
	})

	if _, ok := reg.Find("help"); !ok {
		t.Fatalf("expected to find help command")
	}

	if _, ok := reg.Find("HELP"); !ok {
		t.Fatalf("expected case-insensitive find")
	}

	if _, ok := reg.Find("does-not-exist"); ok {
		t.Fatalf("expected lookup miss")
	}
}

func TestRegistry_VisibleCommands(t *testing.T) {
	t.Parallel()
	reg := NewRegistryWith([]Info{
		{Name: "hidden", Summary: "hidden command", Hidden: true},
		{Name: "visible", Summary: "visible command"},
	})

	cmds := reg.VisibleCommands()
	if len(cmds) != 1 {
		t.Fatalf("expected 1 visible command, got %d", len(cmds))
	}

	if cmds[0].Name != "visible" {
		t.Fatalf("expected visible command, got %s", cmds[0].Name)
	}
}

func TestRegistry_Validate(t *testing.T) {
	t.Parallel()
	reg := NewRegistryWith([]Info{})

	if err := reg.Validate(); err != nil {
		t.Fatalf("expected valid registry, got error: %v", err)
	}

	invalidReg := NewRegistryWith([]Info{
		{Name: "invalid", Summary: ""},
	})

	if err := invalidReg.Validate(); err == nil {
		t.Fatalf("expected validation error for missing summary")
	}
}

func TestNewRegistry_All_ReturnsCopy(t *testing.T) {
	t.Parallel()
	reg := NewRegistry()
	cmds := reg.All()
	if len(cmds) == 0 {
		t.Fatal("expected registry to contain commands")
	}

	originalName := cmds[0].Name
	cmds[0].Name = "mutated"
	fresh := reg.All()
	if fresh[0].Name == "mutated" {
		t.Fatalf("mutating All() result modified registry")
	}
	if fresh[0].Name != originalName {
		t.Fatalf("expected original name %q, got %q", originalName, fresh[0].Name)
	}

	if len(cmds[0].Subcommands) > 0 {
		origSub := reg.All()[0].Subcommands[0].Name
		cmds[0].Subcommands[0].Name = "changed"
		if reg.All()[0].Subcommands[0].Name != origSub {
			t.Fatalf("mutating subcommands altered registry")
		}
	}
}

func TestNewRegistry_Find(t *testing.T) {
	t.Parallel()
	reg := NewRegistry()
	if _, ok := reg.Find("help"); !ok {
		t.Fatalf("expected to find help command")
	}

	if _, ok := reg.Find("debug-keys"); !ok {
		t.Fatalf("expected to find debug-keys command")
	}

	if _, ok := reg.Find("HELP"); !ok {
		t.Fatalf("expected case-insensitive find")
	}

	if _, ok := reg.Find("does-not-exist"); ok {
		t.Fatalf("expected lookup miss")
	}
}

func TestNewRegistry_VisibleCommands(t *testing.T) {
	t.Parallel()
	reg := NewRegistry()
	cmds := reg.VisibleCommands()
	if len(cmds) == 0 {
		t.Fatal("expected visible commands to be returned")
	}

	for _, cmd := range cmds {
		if cmd.Hidden {
			t.Fatalf("visible commands should not include hidden entries: %+v", cmd)
		}
	}

	originalName := cmds[0].Name
	cmds[0].Name = "mutated"
	fresh := reg.VisibleCommands()
	if fresh[0].Name == "mutated" {
		t.Fatalf("modifying VisibleCommands result mutated registry")
	}
	if fresh[0].Name != originalName {
		t.Fatalf("expected original name %q, got %q", originalName, fresh[0].Name)
	}

	if len(cmds[0].Subcommands) > 0 {
		origSub := reg.VisibleCommands()[0].Subcommands[0].Name
		cmds[0].Subcommands[0].Name = "changed"
		if reg.VisibleCommands()[0].Subcommands[0].Name != origSub {
			t.Fatalf("modifying subcommands in VisibleCommands result mutated registry")
		}
	}
}

func TestRegistry_VisibleCommands_ExcludesHidden(t *testing.T) {
	t.Parallel()
	reg := NewRegistryWith([]Info{
		{Name: "__hidden_test__", Summary: "hidden", Hidden: true},
		{Name: "__visible_test__", Summary: "visible"},
	})

	cmds := reg.VisibleCommands()
	for _, cmd := range cmds {
		if cmd.Name == "__hidden_test__" {
			t.Fatalf("hidden command should not appear in results")
		}
	}

	foundVisible := false
	for _, cmd := range cmds {
		if cmd.Name == "__visible_test__" {
			foundVisible = true
			break
		}
	}
	if !foundVisible {
		t.Fatalf("expected visible test command in results")
	}
}

func TestValidate_EmptyCommandName(t *testing.T) {
	t.Parallel()
	commands := []Info{{Name: "   ", Summary: "desc"}}
	if err := Validate(commands); err == nil {
		t.Fatalf("expected validation failure for empty command name")
	}
}

func TestValidate_ValidCommandPasses(t *testing.T) {
	t.Parallel()
	commands := []Info{{Name: "test", Summary: "desc"}}
	if err := Validate(commands); err != nil {
		t.Fatalf("expected valid command to pass validation, got: %v", err)
	}
}

func TestValidate_HiddenCommandWithoutHandlerID(t *testing.T) {
	t.Parallel()
	commands := []Info{{
		Name:    "hidden-test",
		Summary: "desc",
		Hidden:  true,
	}}
	if err := Validate(commands); err != nil {
		t.Fatalf("expected hidden command to be valid without handler ID, got %v", err)
	}
}

func TestValidate_EmptySubcommandName(t *testing.T) {
	t.Parallel()
	commands := []Info{{
		Name:    "test",
		Summary: "desc",
		Subcommands: []SubcommandInfo{
			{Name: "   ", Summary: "desc"},
		},
	}}
	if err := Validate(commands); err == nil {
		t.Fatalf("expected validation failure for empty subcommand name")
	}
}

func TestValidate_MissingSubcommandSummary(t *testing.T) {
	t.Parallel()
	commands := []Info{{
		Name:    "test",
		Summary: "desc",
		Subcommands: []SubcommandInfo{
			{Name: "child", Summary: ""},
		},
	}}
	if err := Validate(commands); err == nil {
		t.Fatalf("expected validation failure for missing subcommand summary")
	}
}

func TestCloneEmptySubcommands(t *testing.T) {
	t.Parallel()
	original := Info{Name: "test", Summary: "desc"}
	clone := (&original).clone()

	if len(clone.Subcommands) != 0 {
		t.Fatalf("expected clone to have no subcommands, got %d", len(clone.Subcommands))
	}

	clone.Name = "mutated"
	if original.Name == clone.Name {
		t.Fatalf("mutating clone should not affect original")
	}
}

func TestCloneNilSlices(t *testing.T) {
	t.Parallel()
	original := Info{
		Name:    "test",
		Summary: "desc",
	}
	clone := (&original).clone()

	clone.Aliases = append(clone.Aliases, "alias")
	clone.Usage = append(clone.Usage, "usage")
	clone.Examples = append(clone.Examples, "example")

	if len(original.Aliases) != 0 || len(original.Usage) != 0 || len(original.Examples) != 0 {
		t.Fatalf("mutating clone slices should not affect original")
	}

	if reflect.DeepEqual(original, clone) {
		t.Fatalf("expected clone to diverge after mutation, indicating defensive copy was not created")
	}
}

func TestNewRegistry(t *testing.T) {
	t.Parallel()
	reg := NewRegistry()
	cmds := reg.All()

	if len(cmds) == 0 {
		t.Fatal("expected NewRegistry to contain default commands")
	}

	// Verify it contains expected commands
	if _, ok := reg.Find("help"); !ok {
		t.Fatal("expected NewRegistry to contain help command")
	}
}

func TestNewRegistryWith(t *testing.T) {
	t.Parallel()
	customCmds := []Info{
		{Name: "custom1", Summary: "first custom"},
		{Name: "custom2", Summary: "second custom"},
	}

	reg := NewRegistryWith(customCmds)
	cmds := reg.All()

	if len(cmds) != 2 {
		t.Fatalf("expected 2 commands, got %d", len(cmds))
	}

	if cmds[0].Name != "custom1" || cmds[1].Name != "custom2" {
		t.Fatalf("expected custom commands, got %v", cmds)
	}
}

func TestCategoryOrder_KnownCategories(t *testing.T) {
	tests := []struct {
		cat  Category
		want int
	}{
		{CategoryBasics, 1},
		{CategoryBranch, 2},
		{CategoryCommit, 3},
		{CategoryRemote, 4},
		{CategoryStatus, 5},
		{CategoryCleanup, 6},
		{CategoryDiff, 7},
		{CategoryTag, 8},
		{CategoryConfig, 9},
		{CategoryHook, 10},
		{CategoryRebase, 11},
		{CategoryStash, 12},
		{CategoryUtility, 13},
	}
	for _, tt := range tests {
		if got := CategoryOrder(tt.cat); got != tt.want {
			t.Errorf("CategoryOrder(%q) = %d, want %d", tt.cat, got, tt.want)
		}
	}
}

func TestCategoryOrder_Unknown(t *testing.T) {
	if got := CategoryOrder(Category("unknown")); got != 999 {
		t.Errorf("CategoryOrder(unknown) = %d, want 999", got)
	}
}

func TestOrderedCategories(t *testing.T) {
	cats := OrderedCategories()
	if len(cats) != 13 {
		t.Fatalf("OrderedCategories() returned %d categories, want 13", len(cats))
	}
	if cats[0] != CategoryBasics {
		t.Errorf("first category = %q, want %q", cats[0], CategoryBasics)
	}
	if cats[len(cats)-1] != CategoryUtility {
		t.Errorf("last category = %q, want %q", cats[len(cats)-1], CategoryUtility)
	}
}
