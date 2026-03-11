package keybindings

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bmf-san/ggc/v8/internal/config"
)

// ── KeyBindingMap getter methods ────────────────────────────────────────────

func TestKeyBindingMap_GetDeleteToEndByte(t *testing.T) {
	km := DefaultKeyBindingMap()
	got := km.GetDeleteToEndByte()
	want := ctrl('k')
	if got != want {
		t.Errorf("GetDeleteToEndByte() = %d, want %d", got, want)
	}
}

func TestKeyBindingMap_GetMoveToBeginningByte(t *testing.T) {
	km := DefaultKeyBindingMap()
	got := km.GetMoveToBeginningByte()
	want := ctrl('a')
	if got != want {
		t.Errorf("GetMoveToBeginningByte() = %d, want %d", got, want)
	}
}

func TestKeyBindingMap_GetMoveToEndByte(t *testing.T) {
	km := DefaultKeyBindingMap()
	got := km.GetMoveToEndByte()
	want := ctrl('e')
	if got != want {
		t.Errorf("GetMoveToEndByte() = %d, want %d", got, want)
	}
}

func TestKeyBindingMap_GetMoveUpByte(t *testing.T) {
	km := DefaultKeyBindingMap()
	got := km.GetMoveUpByte()
	want := ctrl('p')
	if got != want {
		t.Errorf("GetMoveUpByte() = %d, want %d", got, want)
	}
}

func TestKeyBindingMap_GetMoveDownByte(t *testing.T) {
	km := DefaultKeyBindingMap()
	got := km.GetMoveDownByte()
	want := ctrl('n')
	if got != want {
		t.Errorf("GetMoveDownByte() = %d, want %d", got, want)
	}
}

func TestKeyBindingMap_GetAddToWorkflowByte(t *testing.T) {
	km := DefaultKeyBindingMap()
	got := km.GetAddToWorkflowByte()
	// Tab is ASCII 9; getFirstControlByte won't find a ctrl keystroke in a raw seq,
	// so it returns the fallback (9).
	if got != 9 {
		t.Errorf("GetAddToWorkflowByte() = %d, want 9 (Tab)", got)
	}
}

func TestKeyBindingMap_GetToggleWorkflowViewByte(t *testing.T) {
	km := DefaultKeyBindingMap()
	got := km.GetToggleWorkflowViewByte()
	want := ctrl('t')
	if got != want {
		t.Errorf("GetToggleWorkflowViewByte() = %d, want %d", got, want)
	}
}

func TestKeyBindingMap_GetClearWorkflowByte(t *testing.T) {
	km := DefaultKeyBindingMap()
	got := km.GetClearWorkflowByte()
	// ClearWorkflow uses NewCharKeyStroke('c') which is a raw seq, not ctrl.
	// So it falls back to 'c'.
	if got != 'c' {
		t.Errorf("GetClearWorkflowByte() = %d (%q), want %d (%q)", got, got, byte('c'), 'c')
	}
}

func TestKeyBindingMap_GetDeleteToEndByte_EmptySlice(t *testing.T) {
	km := &KeyBindingMap{DeleteToEnd: []KeyStroke{}}
	got := km.GetDeleteToEndByte()
	want := ctrl('k') // fallback
	if got != want {
		t.Errorf("GetDeleteToEndByte (empty) = %d, want %d", got, want)
	}
}

// ── KeyStrokeKind.String ────────────────────────────────────────────────────

func TestKeyStrokeKind_String(t *testing.T) {
	tests := []struct {
		kind KeyStrokeKind
		want string
	}{
		{KeyStrokeCtrl, "Ctrl"},
		{KeyStrokeAlt, "Alt"},
		{KeyStrokeRawSeq, "RawSeq"},
		{KeyStrokeFnKey, "FnKey"},
		{KeyStrokeKind(99), "Unknown"},
	}
	for _, tt := range tests {
		got := tt.kind.String()
		if got != tt.want {
			t.Errorf("KeyStrokeKind(%d).String() = %q, want %q", tt.kind, got, tt.want)
		}
	}
}

// ── KeyStroke constructors ──────────────────────────────────────────────────

func TestNewEnterKeyStroke(t *testing.T) {
	ks := NewEnterKeyStroke()
	if ks.Kind != KeyStrokeRawSeq {
		t.Errorf("NewEnterKeyStroke().Kind = %v, want KeyStrokeRawSeq", ks.Kind)
	}
	if len(ks.Seq) != 1 || ks.Seq[0] != 13 {
		t.Errorf("NewEnterKeyStroke().Seq = %v, want [13]", ks.Seq)
	}
}

func TestNewSpaceKeyStroke(t *testing.T) {
	ks := NewSpaceKeyStroke()
	if ks.Kind != KeyStrokeRawSeq {
		t.Errorf("NewSpaceKeyStroke().Kind = %v, want KeyStrokeRawSeq", ks.Kind)
	}
	if len(ks.Seq) != 1 || ks.Seq[0] != 32 {
		t.Errorf("NewSpaceKeyStroke().Seq = %v, want [32]", ks.Seq)
	}
}

// ── KeyStroke.String ────────────────────────────────────────────────────────

func TestKeyStroke_String_Ctrl(t *testing.T) {
	ks := NewCtrlKeyStroke('a')
	got := ks.String()
	if got != "Ctrl+a" {
		t.Errorf("KeyStroke.String() for ctrl = %q, want %q", got, "Ctrl+a")
	}
}

func TestKeyStroke_String_Alt_WithName(t *testing.T) {
	ks := KeyStroke{Kind: KeyStrokeAlt, Name: "backspace"}
	got := ks.String()
	if got != "Alt+backspace" {
		t.Errorf("KeyStroke.String() for Alt+name = %q, want %q", got, "Alt+backspace")
	}
}

func TestKeyStroke_String_Alt_WithRune(t *testing.T) {
	ks := KeyStroke{Kind: KeyStrokeAlt, Rune: 'x'}
	got := ks.String()
	if got != "Alt+x" {
		t.Errorf("KeyStroke.String() for Alt+rune = %q, want %q", got, "Alt+x")
	}
}

func TestKeyStroke_String_FnKey(t *testing.T) {
	ks := KeyStroke{Kind: KeyStrokeFnKey, Name: "F1"}
	got := ks.String()
	if got != "F1" {
		t.Errorf("KeyStroke.String() for FnKey = %q, want 'F1'", got)
	}
}

func TestKeyStroke_String_Unknown(t *testing.T) {
	ks := KeyStroke{Kind: KeyStrokeKind(999)}
	got := ks.String()
	if got != "Unknown" {
		t.Errorf("KeyStroke.String() for unknown = %q, want 'Unknown'", got)
	}
}

// ── validateKeyStroke ────────────────────────────────────────────────────────

func TestValidateKeyStroke_Ctrl_Valid(t *testing.T) {
	ks := NewCtrlKeyStroke('w')
	if err := validateKeyStroke(ks); err != nil {
		t.Errorf("unexpected error for valid ctrl keystroke: %v", err)
	}
}

func TestValidateKeyStroke_Ctrl_InvalidRune(t *testing.T) {
	ks := KeyStroke{Kind: KeyStrokeCtrl, Rune: 0}
	if err := validateKeyStroke(ks); err == nil {
		t.Error("expected error for ctrl with zero rune")
	}
}

func TestValidateKeyStroke_Alt_ValidRune(t *testing.T) {
	ks := KeyStroke{Kind: KeyStrokeAlt, Rune: 'b'}
	if err := validateKeyStroke(ks); err != nil {
		t.Errorf("unexpected error for valid alt keystroke: %v", err)
	}
}

func TestValidateKeyStroke_Alt_Empty(t *testing.T) {
	ks := KeyStroke{Kind: KeyStrokeAlt}
	if err := validateKeyStroke(ks); err == nil {
		t.Error("expected error for alt keystroke with no rune or name")
	}
}

func TestValidateKeyStroke_RawSeq_Valid(t *testing.T) {
	ks := NewEnterKeyStroke()
	if err := validateKeyStroke(ks); err != nil {
		t.Errorf("unexpected error for valid raw seq: %v", err)
	}
}

func TestValidateKeyStroke_RawSeq_Empty(t *testing.T) {
	ks := KeyStroke{Kind: KeyStrokeRawSeq, Seq: []byte{}}
	if err := validateKeyStroke(ks); err == nil {
		t.Error("expected error for raw seq with empty sequence")
	}
}

func TestValidateKeyStroke_FnKey_Valid(t *testing.T) {
	ks := KeyStroke{Kind: KeyStrokeFnKey, Name: "F1"}
	if err := validateKeyStroke(ks); err != nil {
		t.Errorf("unexpected error for valid fn key: %v", err)
	}
}

func TestValidateKeyStroke_FnKey_NoName(t *testing.T) {
	ks := KeyStroke{Kind: KeyStrokeFnKey}
	if err := validateKeyStroke(ks); err == nil {
		t.Error("expected error for fn key with no name")
	}
}

func TestValidateKeyStroke_Unknown(t *testing.T) {
	ks := KeyStroke{Kind: KeyStrokeKind(999)}
	if err := validateKeyStroke(ks); err == nil {
		t.Error("expected error for unknown keystroke kind")
	}
}

// ── ContextManager: SetContext, ForceEnvironment ─────────────────────────────

func TestContextManager_SetContext(t *testing.T) {
	resolver := NewKeyBindingResolver(&config.Config{})
	RegisterBuiltinProfiles(resolver)
	cm := NewContextManager(resolver)

	var transitions [][2]Context
	cm.RegisterContextCallback(ContextResults, func(from, to Context) {
		transitions = append(transitions, [2]Context{from, to})
	})

	// SetContext to new context
	cm.SetContext(ContextResults)
	if cm.GetCurrentContext() != ContextResults {
		t.Errorf("SetContext: current = %v, want %v", cm.GetCurrentContext(), ContextResults)
	}
	if len(transitions) != 1 {
		t.Errorf("expected 1 transition callback, got %d", len(transitions))
	}

	// SetContext to same context should be no-op
	cm.SetContext(ContextResults)
	if len(transitions) != 1 {
		t.Errorf("SetContext same context should not fire callback, got %d transitions", len(transitions))
	}

	// Stack should be unmodified
	if len(cm.GetContextStack()) != 0 {
		t.Errorf("SetContext should not modify stack, got %v", cm.GetContextStack())
	}
}

func TestContextManager_ForceEnvironment(t *testing.T) {
	resolver := NewKeyBindingResolver(&config.Config{})
	RegisterBuiltinProfiles(resolver)
	cm := NewContextManager(resolver)

	// Should not panic
	cm.ForceEnvironment("darwin", "xterm-256color")
}

func TestContextManager_ForceEnvironment_NilCM(t *testing.T) {
	var cm *ContextManager
	// Should not panic
	cm.ForceEnvironment("linux", "xterm")
}

// ── ContextTransitionAnimator ────────────────────────────────────────────────

func TestContextTransitionAnimator_FadeAndSlide(t *testing.T) {
	cta := NewContextTransitionAnimator()
	cta.SetDuration(0) // no sleep in tests

	cta.SetStyle("fade")
	cta.AnimateTransition(ContextGlobal, ContextResults)

	cta.SetStyle("slide")
	cta.AnimateTransition(ContextGlobal, ContextInput)
}

func TestContextTransitionAnimator_Disable(t *testing.T) {
	cta := NewContextTransitionAnimator()
	cta.Disable()
	// Should return early without doing anything
	cta.AnimateTransition(ContextGlobal, ContextResults)
	if cta.enabled {
		t.Error("expected disabled animator")
	}
}

func TestContextTransitionAnimator_Enable(t *testing.T) {
	cta := NewContextTransitionAnimator()
	cta.Disable()
	cta.Enable()
	if !cta.enabled {
		t.Error("expected enabled animator")
	}
}

func TestContextTransitionAnimator_RegisterAnimation(t *testing.T) {
	cta := NewContextTransitionAnimator()
	cta.RegisterAnimation(func(from, to Context) {})
	cta.RegisterAnimation(func(from, to Context) {})
	if len(cta.animations) != 2 {
		t.Errorf("expected 2 registered animations, got %d", len(cta.animations))
	}
}

func TestContextTransitionAnimator_SetDuration(t *testing.T) {
	cta := NewContextTransitionAnimator()
	cta.SetDuration(500 * time.Millisecond)
	if cta.duration != 500*time.Millisecond {
		t.Errorf("duration = %v, want 500ms", cta.duration)
	}
}

// ── KeybindingExport.ToYAML ──────────────────────────────────────────────────

func TestKeybindingExport_ToYAML(t *testing.T) {
	export := &KeybindingExport{
		Profile:     "default",
		Keybindings: map[string]string{"delete_word": "ctrl+w"},
		Metadata: ExportMetadata{
			Version:  "v8.0.0",
			Platform: "darwin",
			Terminal: "xterm-256color",
		},
	}

	yaml, err := export.ToYAML()
	if err != nil {
		t.Fatalf("ToYAML() error: %v", err)
	}
	if !strings.Contains(yaml, "profile:") {
		t.Errorf("ToYAML() missing 'profile:', got:\n%s", yaml)
	}
}

// ── Profile pure functions ───────────────────────────────────────────────────

func TestGetAllProfiles(t *testing.T) {
	profiles := GetAllProfiles()
	if len(profiles) != 4 {
		t.Fatalf("GetAllProfiles() returned %d profiles, want 4", len(profiles))
	}
}

func TestGetAllContexts_Count(t *testing.T) {
	contexts := GetAllContexts()
	if len(contexts) != 4 {
		t.Fatalf("GetAllContexts() returned %d contexts, want 4", len(contexts))
	}
}

func TestProfile_IsValid(t *testing.T) {
	for _, p := range GetAllProfiles() {
		if !p.IsValid() {
			t.Errorf("Profile(%q).IsValid() = false, want true", p)
		}
	}
	if Profile("invalid").IsValid() {
		t.Error("Profile(\"invalid\").IsValid() = true, want false")
	}
}

func TestContext_IsValid(t *testing.T) {
	for _, c := range GetAllContexts() {
		if !c.IsValid() {
			t.Errorf("Context(%q).IsValid() = false, want true", c)
		}
	}
	if Context("invalid").IsValid() {
		t.Error("Context(\"invalid\").IsValid() = true, want false")
	}
}

func TestKeyBindingProfile_GetAllActions(t *testing.T) {
	p := NewKeyBindingProfile("test", "Test profile")
	p.SetGlobalBinding("action1", []KeyStroke{NewCtrlKeyStroke('a')})
	p.SetContextBinding(ContextInput, "action2", []KeyStroke{NewCtrlKeyStroke('b')})

	actions := p.GetAllActions()
	if len(actions) != 2 {
		t.Errorf("GetAllActions() returned %d actions, want 2: %v", len(actions), actions)
	}
}

func TestKeyBindingProfile_Clone(t *testing.T) {
	p := NewKeyBindingProfile("orig", "Original")
	p.SetGlobalBinding("action1", []KeyStroke{NewCtrlKeyStroke('a')})
	p.SetContextBinding(ContextInput, "action2", []KeyStroke{NewCtrlKeyStroke('b')})

	clone := p.Clone()
	if clone.Name != p.Name {
		t.Errorf("Clone().Name = %q, want %q", clone.Name, p.Name)
	}
	if len(clone.Global) != len(p.Global) {
		t.Errorf("Clone().Global len = %d, want %d", len(clone.Global), len(p.Global))
	}
	// Mutations to clone should not affect original
	clone.Global["action1"] = []KeyStroke{NewCtrlKeyStroke('z')}
	if clone.Global["action1"][0].Rune == p.Global["action1"][0].Rune {
		t.Error("Clone is not a deep copy: modifying clone affected original")
	}
}

func TestGetAllProfilesBuiltin(t *testing.T) {
	profiles := GetAllProfilesBuiltin()
	if len(profiles) != 4 {
		t.Fatalf("GetAllProfilesBuiltin() returned %d profiles, want 4", len(profiles))
	}
}

// ── ValidateProfile / ValidateAllBuiltinProfiles ─────────────────────────────

func TestValidateAllBuiltinProfiles(t *testing.T) {
	if err := ValidateAllBuiltinProfiles(); err != nil {
		t.Fatalf("ValidateAllBuiltinProfiles() error: %v", err)
	}
}

func TestValidateProfile_Nil(t *testing.T) {
	if err := ValidateProfile(nil); err == nil {
		t.Error("expected error for nil profile")
	}
}

func TestValidateProfile_EmptyName(t *testing.T) {
	p := NewKeyBindingProfile("", "desc")
	if err := ValidateProfile(p); err == nil {
		t.Error("expected error for empty profile name")
	}
}

func TestValidateProfile_EmptyDescription(t *testing.T) {
	p := NewKeyBindingProfile("test", "")
	if err := ValidateProfile(p); err == nil {
		t.Error("expected error for empty description")
	}
}

func TestValidateProfile_NilContexts(t *testing.T) {
	p := &KeyBindingProfile{Name: "test", Description: "desc", Contexts: nil}
	if err := ValidateProfile(p); err == nil {
		t.Error("expected error for nil contexts")
	}
}

// ── GetProfileStatistics / CompareProfiles ────────────────────────────────────

func TestGetProfileStatistics_Nil(t *testing.T) {
	stats := GetProfileStatistics(nil)
	if len(stats) != 0 {
		t.Errorf("GetProfileStatistics(nil) returned non-empty map: %v", stats)
	}
}

func TestGetProfileStatistics_WithProfile(t *testing.T) {
	p := NewKeyBindingProfile("test", "Test")
	p.SetContextBinding(ContextInput, "move_up", []KeyStroke{NewCtrlKeyStroke('p')})
	stats := GetProfileStatistics(p)
	if stats["profile_name"] != "test" {
		t.Errorf("stats[profile_name] = %v, want %q", stats["profile_name"], "test")
	}
	if stats["total_context_bindings"].(int) != 1 {
		t.Errorf("stats[total_context_bindings] = %v, want 1", stats["total_context_bindings"])
	}
}

func TestCompareProfiles_NilInputs(t *testing.T) {
	result := CompareProfiles(nil, nil)
	if result["error"] == nil {
		t.Error("expected error key in comparison result for nil inputs")
	}
}

func TestCompareProfiles_TwoProfiles(t *testing.T) {
	p1 := NewKeyBindingProfile("p1", "Profile 1")
	p1.SetContextBinding(ContextInput, "move_up", []KeyStroke{NewCtrlKeyStroke('p')})
	p2 := NewKeyBindingProfile("p2", "Profile 2")
	p2.SetContextBinding(ContextResults, "move_down", []KeyStroke{NewCtrlKeyStroke('n')})

	result := CompareProfiles(p1, p2)
	if result["profile1_name"] != "p1" {
		t.Errorf("profile1_name = %v, want p1", result["profile1_name"])
	}
}

// ── ProfileSwitcher ──────────────────────────────────────────────────────────

func newTestSwitcher() *ProfileSwitcher {
	resolver := NewKeyBindingResolver(nil)
	RegisterBuiltinProfiles(resolver)
	return NewProfileSwitcher(resolver, nil)
}

func TestProfileSwitcher_GetCurrentProfile(t *testing.T) {
	ps := newTestSwitcher()
	if ps.GetCurrentProfile() != ProfileDefault {
		t.Errorf("GetCurrentProfile() = %v, want %v", ps.GetCurrentProfile(), ProfileDefault)
	}
}

func TestProfileSwitcher_GetAvailableProfiles(t *testing.T) {
	ps := newTestSwitcher()
	profiles := ps.GetAvailableProfiles()
	if len(profiles) != 4 {
		t.Fatalf("GetAvailableProfiles() returned %d profiles, want 4", len(profiles))
	}
}

func TestProfileSwitcher_CanSwitchTo_Valid(t *testing.T) {
	ps := newTestSwitcher()
	ok, err := ps.CanSwitchTo(ProfileEmacs)
	if err != nil {
		t.Fatalf("CanSwitchTo(emacs) error: %v", err)
	}
	if !ok {
		t.Error("CanSwitchTo(emacs) = false, want true")
	}
}

func TestProfileSwitcher_CanSwitchTo_NotRegistered(t *testing.T) {
	resolver := NewKeyBindingResolver(nil) // no profiles registered
	ps := NewProfileSwitcher(resolver, nil)
	_, err := ps.CanSwitchTo(ProfileEmacs)
	if err == nil {
		t.Error("expected error for unregistered profile")
	}
}

func TestProfileSwitcher_GetProfileComparison(t *testing.T) {
	ps := newTestSwitcher()
	result, err := ps.GetProfileComparison(ProfileEmacs)
	if err != nil {
		t.Fatalf("GetProfileComparison(emacs) error: %v", err)
	}
	if result["profile1_name"] == nil {
		t.Error("expected profile1_name in comparison result")
	}
}

func TestProfileSwitcher_ShowCurrentProfileCommand(t *testing.T) {
	ps := newTestSwitcher()
	got := ShowCurrentProfileCommand(ps)
	if !strings.Contains(got, "default") {
		t.Errorf("ShowCurrentProfileCommand() = %q, want to contain 'default'", got)
	}
}

func TestHandleProfileSwitchCommand_List(t *testing.T) {
	ps := newTestSwitcher()
	if err := HandleProfileSwitchCommand(ps, "list"); err != nil {
		t.Fatalf("HandleProfileSwitchCommand(list) error: %v", err)
	}
}

func TestHandleProfileSwitchCommand_Unknown(t *testing.T) {
	ps := newTestSwitcher()
	if err := HandleProfileSwitchCommand(ps, "bogus"); err == nil {
		t.Error("expected error for unknown subcommand")
	}
}

func TestHandleProfileSwitchCommand_Empty(t *testing.T) {
	ps := newTestSwitcher()
	if err := HandleProfileSwitchCommand(ps, ""); err == nil {
		t.Error("expected error for empty command")
	}
}

func TestHandleProfileSwitchCommand_Preview(t *testing.T) {
	ps := newTestSwitcher()
	if err := HandleProfileSwitchCommand(ps, "preview emacs"); err != nil {
		t.Fatalf("HandleProfileSwitchCommand(preview emacs) error: %v", err)
	}
}

func TestHandleProfileSwitchCommand_Compare(t *testing.T) {
	ps := newTestSwitcher()
	if err := HandleProfileSwitchCommand(ps, "compare emacs"); err != nil {
		t.Fatalf("HandleProfileSwitchCommand(compare emacs) error: %v", err)
	}
}

func TestHandleProfileSwitchCommand_PreviewNoArg(t *testing.T) {
	ps := newTestSwitcher()
	if err := HandleProfileSwitchCommand(ps, "preview"); err == nil {
		t.Error("expected error for preview without arg")
	}
}

// ─── DetectTerminal: cover all TERM_PROGRAM and TERM branches ─────────────────

func TestDetectTerminal_TermProgram(t *testing.T) {
	t.Setenv("TERM", "")

	cases := []struct {
		prog string
		want string
	}{
		{"iTerm.app", "iterm"},
		{"Apple_Terminal", "terminal"},
		{"vscode", "vscode"},
		{"Hyper", "hyper"},
	}
	for _, c := range cases {
		t.Setenv("TERM_PROGRAM", c.prog)
		if got := DetectTerminal(); got != c.want {
			t.Errorf("TERM_PROGRAM=%q: want %q, got %q", c.prog, c.want, got)
		}
	}
}

func TestDetectTerminal_TERM(t *testing.T) {
	t.Setenv("TERM_PROGRAM", "")

	cases := []struct {
		term string
		want string
	}{
		{"tmux-256color", "tmux"},
		{"screen-256color", "screen"},
		{"xterm-256color", "xterm"},
		{"alacritty", "alacritty"},
		{"kitty", "kitty"},
		{"wezterm", "wezterm"},
		{"konsole-256color", "konsole"},
		{"gnome-256color", "gnome-terminal"},
		{"rxvt-unicode", "rxvt"},
		{"dumb", "dumb"},
		{"", "generic"},
	}
	for _, c := range cases {
		t.Setenv("TERM", c.term)
		if got := DetectTerminal(); got != c.want {
			t.Errorf("TERM=%q: want %q, got %q", c.term, c.want, got)
		}
	}
}

// ─── GetTerminalCapabilities: cover all missing terminal branches ─────────────

func TestGetTerminalCapabilities_StandardTerminals(t *testing.T) {
	for _, terminal := range []string{"xterm", "gnome-terminal", "konsole"} {
		caps := GetTerminalCapabilities(terminal)
		if !caps["alt_keys"] {
			t.Errorf("%s: expected alt_keys=true", terminal)
		}
		if caps["mouse"] {
			t.Errorf("%s: expected mouse=false", terminal)
		}
	}
}

func TestGetTerminalCapabilities_Multiplexers(t *testing.T) {
	for _, terminal := range []string{"tmux", "screen"} {
		caps := GetTerminalCapabilities(terminal)
		if !caps["alt_keys"] {
			t.Errorf("%s: expected alt_keys=true", terminal)
		}
		if caps["mouse"] {
			t.Errorf("%s: expected mouse=false", terminal)
		}
	}
}

func TestGetTerminalCapabilities_MacTerminal(t *testing.T) {
	caps := GetTerminalCapabilities("terminal")
	if !caps["alt_keys"] {
		t.Error("terminal: expected alt_keys=true")
	}
	if caps["mouse"] {
		t.Error("terminal: expected mouse=false")
	}
}

func TestGetTerminalCapabilities_Generic(t *testing.T) {
	caps := GetTerminalCapabilities("generic")
	if !caps["alt_keys"] {
		t.Error("generic: expected alt_keys=true")
	}
	if caps["function_keys"] {
		t.Error("generic: expected function_keys=false")
	}
	if caps["mouse"] {
		t.Error("generic: expected mouse=false")
	}
}

func TestGetTerminalCapabilities_Unknown(t *testing.T) {
	caps := GetTerminalCapabilities("unknown-terminal-xyz")
	// default branch — same as generic
	if caps["mouse"] {
		t.Error("unknown terminal: expected mouse=false")
	}
}

// ─── GetPlatformSpecificKeyBindings: cover linux, windows, default ────────────

func TestGetPlatformSpecificKeyBindings_Linux(t *testing.T) {
	bindings := GetPlatformSpecificKeyBindings("linux")
	ks, ok := bindings["delete_word"]
	if !ok {
		t.Fatal("linux: expected delete_word binding")
	}
	if len(ks) == 0 {
		t.Error("linux: expected at least one delete_word keystroke")
	}
}

func TestGetPlatformSpecificKeyBindings_BSD(t *testing.T) {
	bindings := GetPlatformSpecificKeyBindings("bsd")
	if _, ok := bindings["delete_word"]; !ok {
		t.Error("bsd: expected delete_word binding")
	}
}

func TestGetPlatformSpecificKeyBindings_Unix(t *testing.T) {
	bindings := GetPlatformSpecificKeyBindings("unix")
	if _, ok := bindings["delete_word"]; !ok {
		t.Error("unix: expected delete_word binding")
	}
}

func TestGetPlatformSpecificKeyBindings_Windows(t *testing.T) {
	// Windows branch exists but currently returns empty map
	bindings := GetPlatformSpecificKeyBindings("windows")
	_ = bindings // no panic = pass
}

func TestGetPlatformSpecificKeyBindings_Default(t *testing.T) {
	bindings := GetPlatformSpecificKeyBindings("unknown-platform")
	if len(bindings) != 0 {
		t.Error("unknown platform: expected empty bindings")
	}
}

// ─── GetTerminalSpecificKeyBindings: cover all terminal branches ──────────────

func TestGetTerminalSpecificKeyBindings_All(t *testing.T) {
	for _, terminal := range []string{"tmux", "screen", "iterm", "alacritty", "kitty", "wezterm", "other"} {
		bindings := GetTerminalSpecificKeyBindings(terminal)
		_ = bindings // all return empty map — no panic = pass
	}
}

// ─── formatKeystrokeForExport: cover RawSeq cases, FnKey, unknown ─────────────

func formatExporter() *KeybindingExporter {
	cfg := &config.Config{}
	cfg.Interactive.Profile = "emacs"
	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)
	return NewKeybindingExporter(resolver)
}

func TestFormatKeystrokeForExport_RawSeq_SingleByte(t *testing.T) {
	ke := formatExporter()
	cases := []struct {
		seq  []byte
		want string
	}{
		{[]byte{9}, "tab"},
		{[]byte{13}, "enter"},
		{[]byte{27}, "esc"},
		{[]byte{32}, "space"},
	}
	for _, c := range cases {
		ks := NewRawKeyStroke(c.seq)
		got := ke.formatKeystrokeForExport(ks)
		if got != c.want {
			t.Errorf("seq=%v: want %q, got %q", c.seq, c.want, got)
		}
	}
}

func TestFormatKeystrokeForExport_RawSeq_ArrowKeys(t *testing.T) {
	ke := formatExporter()
	cases := []struct {
		seq  []byte
		want string
	}{
		{[]byte{27, 91, 65}, "up"},
		{[]byte{27, 91, 66}, "down"},
		{[]byte{27, 91, 67}, "right"},
		{[]byte{27, 91, 68}, "left"},
	}
	for _, c := range cases {
		ks := NewRawKeyStroke(c.seq)
		got := ke.formatKeystrokeForExport(ks)
		if got != c.want {
			t.Errorf("seq=%v: want %q, got %q", c.seq, c.want, got)
		}
	}
}

func TestFormatKeystrokeForExport_RawSeq_Raw(t *testing.T) {
	ke := formatExporter()
	ks := NewRawKeyStroke([]byte{0xff, 0x01})
	got := ke.formatKeystrokeForExport(ks)
	if !strings.HasPrefix(got, "raw:") {
		t.Errorf("expected raw: prefix, got %q", got)
	}
}

func TestFormatKeystrokeForExport_FnKey(t *testing.T) {
	ke := formatExporter()
	ks := KeyStroke{Kind: KeyStrokeFnKey, Name: "F1"}
	got := ke.formatKeystrokeForExport(ks)
	if got != "f1" {
		t.Errorf("want 'f1', got %q", got)
	}
}

func TestFormatKeystrokeForExport_Unknown(t *testing.T) {
	ke := formatExporter()
	ks := KeyStroke{Kind: KeyStrokeKind(99), Rune: 'x'}
	got := ke.formatKeystrokeForExport(ks)
	if !strings.HasPrefix(got, "unknown:") {
		t.Errorf("expected 'unknown:' prefix, got %q", got)
	}
}

// ─── KeybindingExport.ToYAML: cover Contexts and Platform branches ────────────

func TestKeybindingExport_ToYAML_WithContextsAndPlatform(t *testing.T) {
	export := &KeybindingExport{
		Profile: "test",
		Keybindings: map[string]string{
			"delete_word": "ctrl+w",
		},
		Contexts: map[string]map[string]string{
			"input": {"complete": "tab"},
		},
		Platform: map[string]map[string]string{
			"darwin": {"delete_word": "alt+backspace"},
		},
		Metadata: ExportMetadata{
			Version:    "v1",
			ExportedAt: time.Now(),
			ExportedBy: "test",
			Platform:   "darwin",
			Terminal:   "iterm",
			DeltaFrom:  "base-profile",
		},
	}
	yaml, err := export.ToYAML()
	if err != nil {
		t.Fatalf("ToYAML() error: %v", err)
	}
	for _, want := range []string{"contexts:", "darwin:", "delta_from:"} {
		if !strings.Contains(yaml, want) {
			t.Errorf("ToYAML() output missing %q", want)
		}
	}
}

// ─── KeybindingImporter: parseImportFile via Import(InputFile=...) ────────────

func minimalImportYAML() []byte {
	return []byte(`profile: emacs
keybindings:
  delete_word: "ctrl+w"
metadata:
  exported_at: "2024-01-01T00:00:00Z"
  exported_by: "test"
  version: "v1"
  platform: "darwin"
  terminal: "iterm"
`)
}

func testImporter() *KeybindingImporter {
	cfg := &config.Config{}
	cfg.Interactive.Profile = "emacs"
	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)
	return NewKeybindingImporter(resolver)
}

func TestKeybindingImporter_Import_InputFile(t *testing.T) {
	dir := t.TempDir()
	fpath := filepath.Join(dir, "keybindings.yaml")
	if err := os.WriteFile(fpath, minimalImportYAML(), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	ki := testImporter()
	if err := ki.Import(ImportOptions{InputFile: fpath}); err != nil {
		t.Fatalf("Import(InputFile) error: %v", err)
	}
}

func TestKeybindingImporter_Import_InputFile_NotFound(t *testing.T) {
	ki := testImporter()
	err := ki.Import(ImportOptions{InputFile: "/nonexistent/path/keybindings.yaml"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestKeybindingImporter_Import_Interactive(t *testing.T) {
	ki := testImporter()
	err := ki.Import(ImportOptions{Data: minimalImportYAML(), Interactive: true})
	if err != nil {
		t.Fatalf("Import(Interactive) error: %v", err)
	}
}

func TestKeybindingImporter_Import_BackupConfig(t *testing.T) {
	ki := testImporter()
	err := ki.Import(ImportOptions{Data: minimalImportYAML(), BackupConfig: true})
	if err != nil {
		t.Fatalf("Import(BackupConfig) error: %v", err)
	}
}

// ─── resolver_user.go: cover clear_line, delete_to_end, move_left, move_right ─

func testResolver() *KeyBindingResolver {
	cfg := &config.Config{}
	cfg.Interactive.Profile = "emacs"
	resolver := NewKeyBindingResolver(cfg)
	RegisterBuiltinProfiles(resolver)
	return resolver
}

func TestApplyUserEditingAction_ClearLine(t *testing.T) {
	r := testResolver()
	km := DefaultKeyBindingMap()
	ks := []KeyStroke{NewTabKeyStroke()}
	if !r.applyUserEditingAction(km, "clear_line", ks) {
		t.Error("expected applyUserEditingAction to return true for clear_line")
	}
	if len(km.ClearLine) == 0 {
		t.Error("expected ClearLine to be set")
	}
}

func TestApplyUserEditingAction_DeleteToEnd(t *testing.T) {
	r := testResolver()
	km := DefaultKeyBindingMap()
	ks := []KeyStroke{NewTabKeyStroke()}
	if !r.applyUserEditingAction(km, "delete_to_end", ks) {
		t.Error("expected applyUserEditingAction to return true for delete_to_end")
	}
	if len(km.DeleteToEnd) == 0 {
		t.Error("expected DeleteToEnd to be set")
	}
}

func TestApplyUserNavigationAction_MoveLeft(t *testing.T) {
	r := testResolver()
	km := DefaultKeyBindingMap()
	ks := []KeyStroke{NewTabKeyStroke()}
	if !r.applyUserNavigationAction(km, "move_left", ks) {
		t.Error("expected applyUserNavigationAction to return true for move_left")
	}
	if len(km.MoveLeft) == 0 {
		t.Error("expected MoveLeft to be set")
	}
}

func TestApplyUserNavigationAction_MoveRight(t *testing.T) {
	r := testResolver()
	km := DefaultKeyBindingMap()
	ks := []KeyStroke{NewTabKeyStroke()}
	if !r.applyUserNavigationAction(km, "move_right", ks) {
		t.Error("expected applyUserNavigationAction to return true for move_right")
	}
	if len(km.MoveRight) == 0 {
		t.Error("expected MoveRight to be set")
	}
}

func TestParseUserBindingValue_Slice(t *testing.T) {
	r := testResolver()
	got := r.parseUserBindingValue([]interface{}{"ctrl+w", "ctrl+u"})
	if len(got) != 2 {
		t.Errorf("expected 2 keystrokes from slice input, got %d", len(got))
	}
}
