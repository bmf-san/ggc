package keybindings

import (
	"fmt"
	"strings"
)

// ContextualMapApplier applies resolved keybindings to interested consumers.
type ContextualMapApplier interface {
	ApplyContextualKeybindings(*ContextualKeyBindingMap)
}

// ProfileSwitcher manages runtime profile switching functionality
type ProfileSwitcher struct {
	resolver       *KeyBindingResolver
	currentProfile Profile
	applier        ContextualMapApplier
}

// NewProfileSwitcher creates a new profile switcher
func NewProfileSwitcher(resolver *KeyBindingResolver, applier ContextualMapApplier) *ProfileSwitcher {
	return &ProfileSwitcher{
		resolver:       resolver,
		currentProfile: ProfileDefault,
		applier:        applier,
	}
}

// SwitchProfile switches to a new profile at runtime
func (ps *ProfileSwitcher) SwitchProfile(newProfile Profile) error {
	if _, exists := ps.resolver.GetProfile(newProfile); !exists {
		return fmt.Errorf("profile %s not found", newProfile)
	}

	ps.resolver.ClearCache()

	newContextualMap, err := ps.resolver.ResolveContextual(newProfile)
	if err != nil {
		return fmt.Errorf("failed to resolve profile %s: %w", newProfile, err)
	}

	if ps.applier != nil {
		ps.applier.ApplyContextualKeybindings(newContextualMap)
	}

	oldProfile := ps.currentProfile
	ps.currentProfile = newProfile

	fmt.Printf("Switched keybinding profile from %s to %s\n", oldProfile, newProfile)

	return nil
}

// GetCurrentProfile returns the currently active profile
func (ps *ProfileSwitcher) GetCurrentProfile() Profile {
	return ps.currentProfile
}

// GetAvailableProfiles returns all available profiles for switching
func (ps *ProfileSwitcher) GetAvailableProfiles() []Profile {
	return GetAllProfilesBuiltin()
}

// CanSwitchTo checks if switching to a profile is possible
func (ps *ProfileSwitcher) CanSwitchTo(profile Profile) (bool, error) {
	if _, exists := ps.resolver.GetProfile(profile); !exists {
		return false, fmt.Errorf("profile %s not registered", profile)
	}

	profileDef, _ := ps.resolver.GetProfile(profile)
	if err := ValidateProfile(profileDef); err != nil {
		return false, fmt.Errorf("profile %s validation failed: %w", profile, err)
	}

	return true, nil
}

// PreviewProfile returns a preview of what keybindings would be active with the new profile
func (ps *ProfileSwitcher) PreviewProfile(profile Profile) (*ContextualKeyBindingMap, error) {
	if _, exists := ps.resolver.GetProfile(profile); !exists {
		return nil, fmt.Errorf("profile %s not found", profile)
	}

	tempResolver := NewKeyBindingResolver(ps.resolver.userConfig)
	RegisterBuiltinProfiles(tempResolver)

	return tempResolver.ResolveContextual(profile)
}

// GetProfileComparison compares current profile with another profile
func (ps *ProfileSwitcher) GetProfileComparison(otherProfile Profile) (map[string]interface{}, error) {
	currentProfileDef, exists := ps.resolver.GetProfile(ps.currentProfile)
	if !exists {
		return nil, fmt.Errorf("current profile %s not found", ps.currentProfile)
	}

	otherProfileDef, exists := ps.resolver.GetProfile(otherProfile)
	if !exists {
		return nil, fmt.Errorf("comparison profile %s not found", otherProfile)
	}

	return CompareProfiles(currentProfileDef, otherProfileDef), nil
}

// ReloadCurrentProfile reloads the current profile (useful for config changes)
func (ps *ProfileSwitcher) ReloadCurrentProfile() error {
	return ps.SwitchProfile(ps.currentProfile)
}

type profileSwitchHandler func(*ProfileSwitcher, []string) error

var profileSwitchCommandHandlers = map[string]profileSwitchHandler{
	"list":    handleProfileListCommand,
	"switch":  handleProfileSwitchCommand,
	"preview": handleProfilePreviewCommand,
	"compare": handleProfileCompareCommand,
	"reload":  handleProfileReloadCommand,
}

// HandleProfileSwitchCommand processes profile switching commands
func HandleProfileSwitchCommand(switcher *ProfileSwitcher, command string) error {
	parts := strings.Fields(strings.TrimSpace(command))
	if len(parts) == 0 {
		return fmt.Errorf("no command provided")
	}

	subcommand := parts[0]
	args := parts[1:]

	handler, ok := profileSwitchCommandHandlers[subcommand]
	if !ok {
		return fmt.Errorf("unknown subcommand: %s", subcommand)
	}

	return handler(switcher, args)
}

func handleProfileListCommand(switcher *ProfileSwitcher, _ []string) error {
	profiles := switcher.GetAvailableProfiles()
	fmt.Println("Available profiles:")
	for _, profile := range profiles {
		currentMarker := ""
		if profile == switcher.GetCurrentProfile() {
			currentMarker = " (current)"
		}
		fmt.Printf("  - %s%s\n", profile, currentMarker)
	}

	return nil
}

func handleProfileSwitchCommand(switcher *ProfileSwitcher, args []string) error {
	profile, err := requireProfileArg(args, "switch <profile>")
	if err != nil {
		return err
	}

	return switcher.SwitchProfile(profile)
}

func handleProfilePreviewCommand(switcher *ProfileSwitcher, args []string) error {
	profile, err := requireProfileArg(args, "preview <profile>")
	if err != nil {
		return err
	}

	preview, err := switcher.PreviewProfile(profile)
	if err != nil {
		return err
	}

	fmt.Printf("Preview for profile %s:\n", profile)
	for ctx, mapBinding := range preview.Contexts {
		fmt.Printf("  Context: %s\n", ctx)
		fmt.Printf("    move_up                 %-20s Move up one line\n", FormatKeyStrokesForDisplay(mapBinding.MoveUp))
		fmt.Printf("    move_down               %-20s Move down one line\n", FormatKeyStrokesForDisplay(mapBinding.MoveDown))
		fmt.Printf("    move_to_beginning       %-20s Move to line beginning\n", FormatKeyStrokesForDisplay(mapBinding.MoveToBeginning))
		fmt.Printf("    move_to_end             %-20s Move to line end\n", FormatKeyStrokesForDisplay(mapBinding.MoveToEnd))
		fmt.Printf("    delete_word             %-20s Delete previous word\n", FormatKeyStrokesForDisplay(mapBinding.DeleteWord))
		fmt.Printf("    delete_to_end           %-20s Delete to line end\n", FormatKeyStrokesForDisplay(mapBinding.DeleteToEnd))
		fmt.Printf("    clear_line              %-20s Clear entire line\n", FormatKeyStrokesForDisplay(mapBinding.ClearLine))
	}

	return nil
}

func handleProfileCompareCommand(switcher *ProfileSwitcher, args []string) error {
	profile, err := requireProfileArg(args, "compare <profile>")
	if err != nil {
		return err
	}

	comparison, err := switcher.GetProfileComparison(profile)
	if err != nil {
		return err
	}

	fmt.Printf("Comparison between current profile (%s) and %s:\n", switcher.GetCurrentProfile(), profile)
	for category, value := range comparison {
		fmt.Printf("  %s: %v\n", category, value)
	}

	return nil
}

func handleProfileReloadCommand(switcher *ProfileSwitcher, _ []string) error {
	return switcher.ReloadCurrentProfile()
}

func requireProfileArg(args []string, usage string) (Profile, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("usage: %s", usage)
	}

	return Profile(args[0]), nil
}

// ShowCurrentProfileCommand returns a string representing the current profile status
func ShowCurrentProfileCommand(switcher *ProfileSwitcher) string {
	return fmt.Sprintf("Current profile: %s", switcher.GetCurrentProfile())
}

// RuntimeProfileSwitcher enables switching profiles without restart
type RuntimeProfileSwitcher struct {
	resolver        *KeyBindingResolver
	currentProfile  Profile
	contextManager  *ContextManager
	switchCallbacks []func(Profile, Profile)
}

// NewRuntimeProfileSwitcher creates a new runtime profile switcher
func NewRuntimeProfileSwitcher(resolver *KeyBindingResolver, contextManager *ContextManager) *RuntimeProfileSwitcher {
	return &RuntimeProfileSwitcher{
		resolver:        resolver,
		currentProfile:  ProfileDefault,
		contextManager:  contextManager,
		switchCallbacks: make([]func(Profile, Profile), 0),
	}
}

// SwitchProfile changes the active profile at runtime
func (rps *RuntimeProfileSwitcher) SwitchProfile(newProfile Profile) error {
	// Validate profile exists
	if _, exists := rps.resolver.GetProfile(newProfile); !exists {
		return fmt.Errorf("profile '%s' not found", newProfile)
	}

	oldProfile := rps.currentProfile
	rps.currentProfile = newProfile

	// Clear resolver cache to force re-resolution with new profile
	rps.resolver.ClearCache()

	// Notify callbacks
	for _, callback := range rps.switchCallbacks {
		callback(oldProfile, newProfile)
	}

	fmt.Printf("Switched from profile '%s' to '%s'\n", oldProfile, newProfile)
	return nil
}

// GetCurrentProfile returns the currently active profile
func (rps *RuntimeProfileSwitcher) GetCurrentProfile() Profile {
	return rps.currentProfile
}

// RegisterSwitchCallback registers a callback for profile switches
func (rps *RuntimeProfileSwitcher) RegisterSwitchCallback(callback func(Profile, Profile)) {
	rps.switchCallbacks = append(rps.switchCallbacks, callback)
}

// CycleProfile cycles through available profiles
func (rps *RuntimeProfileSwitcher) CycleProfile() error {
	profiles := []Profile{ProfileDefault, ProfileEmacs, ProfileVi, ProfileReadline}

	currentIndex := 0
	for i, p := range profiles {
		if p == rps.currentProfile {
			currentIndex = i
			break
		}
	}

	nextIndex := (currentIndex + 1) % len(profiles)
	return rps.SwitchProfile(profiles[nextIndex])
}
