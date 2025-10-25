package command

import (
	"reflect"
	"testing"
)

func TestValidateAll(t *testing.T) {
	if err := ValidateAll(); err != nil {
		t.Fatalf("ValidateAll() returned error: %v", err)
	}
}

func TestValidate_DuplicateCommand(t *testing.T) {
	commands := []Info{
		{Name: "test", Summary: "one", HandlerID: "test"},
		{Name: "test", Summary: "two", HandlerID: "test"},
	}

	if err := Validate(commands); err == nil {
		t.Fatalf("expected duplicate command validation to fail")
	}
}

func TestValidate_MissingSummary(t *testing.T) {
	commands := []Info{{Name: "test", HandlerID: "test"}}
	if err := Validate(commands); err == nil {
		t.Fatalf("expected missing summary validation failure")
	}
}

func TestValidate_DuplicateSubcommand(t *testing.T) {
	commands := []Info{
		{
			Name:      "test",
			Summary:   "ok",
			HandlerID: "test",
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

func TestAllReturnsCopy(t *testing.T) {
	cmds := All()
	if len(cmds) == 0 {
		t.Fatal("expected registry to contain commands")
	}

	cmds[0].Name = "mutated"
	if registry[0].Name == "mutated" {
		t.Fatalf("mutating All() result modified registry")
	}

	if len(cmds[0].Subcommands) > 0 {
		original := registry[0].Subcommands[0].Name
		cmds[0].Subcommands[0].Name = "changed"
		if registry[0].Subcommands[0].Name != original {
			t.Fatalf("mutating subcommands altered registry")
		}
	}
}

func TestFind(t *testing.T) {
	if _, ok := Find("help"); !ok {
		t.Fatalf("expected to find help command")
	}

	if _, ok := Find("debug-keys"); !ok {
		t.Fatalf("expected to find debug-keys command")
	}

	if _, ok := Find("HELP"); !ok {
		t.Fatalf("expected case-insensitive find")
	}

	if _, ok := Find("does-not-exist"); ok {
		t.Fatalf("expected lookup miss")
	}
}

func TestVisibleCommands(t *testing.T) {
	cmds := VisibleCommands()
	if len(cmds) == 0 {
		t.Fatal("expected visible commands to be returned")
	}

	for _, cmd := range cmds {
		if cmd.Hidden {
			t.Fatalf("visible commands should not include hidden entries: %+v", cmd)
		}
	}

	originalName := registry[0].Name
	cmds[0].Name = "mutated"
	if registry[0].Name == "mutated" {
		t.Fatalf("modifying VisibleCommands result mutated registry")
	}
	registry[0].Name = originalName

	if len(cmds[0].Subcommands) > 0 {
		origSub := registry[0].Subcommands[0].Name
		cmds[0].Subcommands[0].Name = "changed"
		if registry[0].Subcommands[0].Name != origSub {
			t.Fatalf("modifying subcommands in VisibleCommands result mutated registry")
		}
	}
}

func TestVisibleCommands_WithHiddenCommands(t *testing.T) {
	hidden := Info{Name: "__hidden_test__", Summary: "hidden", Hidden: true}
	visible := Info{Name: "__visible_test__", Summary: "visible", HandlerID: "visible"}
	registry = append(registry, hidden, visible)
	defer func() {
		registry = registry[:len(registry)-2]
	}()

	cmds := VisibleCommands()
	for _, cmd := range cmds {
		if cmd.Name == hidden.Name {
			t.Fatalf("hidden command %q should not appear in results", hidden.Name)
		}
	}

	foundVisible := false
	for _, cmd := range cmds {
		if cmd.Name == visible.Name {
			foundVisible = true
			break
		}
	}
	if !foundVisible {
		t.Fatalf("expected visible test command %q in results", visible.Name)
	}
}

func TestValidate_EmptyCommandName(t *testing.T) {
	commands := []Info{{Name: " \t", Summary: "desc", HandlerID: "handler"}}
	if err := Validate(commands); err == nil {
		t.Fatalf("expected validation failure for empty command name")
	}
}

func TestValidate_MissingHandlerID(t *testing.T) {
	commands := []Info{{Name: "test", Summary: "desc"}}
	if err := Validate(commands); err == nil {
		t.Fatalf("expected validation failure for missing handler ID")
	}
}

func TestValidate_HiddenCommandWithoutHandlerID(t *testing.T) {
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
	commands := []Info{{
		Name:      "test",
		Summary:   "desc",
		HandlerID: "handler",
		Subcommands: []SubcommandInfo{
			{Name: "   ", Summary: "desc"},
		},
	}}
	if err := Validate(commands); err == nil {
		t.Fatalf("expected validation failure for empty subcommand name")
	}
}

func TestValidate_MissingSubcommandSummary(t *testing.T) {
	commands := []Info{{
		Name:      "test",
		Summary:   "desc",
		HandlerID: "handler",
		Subcommands: []SubcommandInfo{
			{Name: "child", Summary: ""},
		},
	}}
	if err := Validate(commands); err == nil {
		t.Fatalf("expected validation failure for missing subcommand summary")
	}
}

func TestCloneEmptySubcommands(t *testing.T) {
	original := Info{Name: "test", Summary: "desc", HandlerID: "handler"}
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
	original := Info{
		Name:      "test",
		Summary:   "desc",
		HandlerID: "handler",
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
