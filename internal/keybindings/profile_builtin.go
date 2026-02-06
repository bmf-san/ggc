package keybindings

// RegisterBuiltinProfiles registers all built-in profiles with the resolver
func RegisterBuiltinProfiles(resolver *KeyBindingResolver) {
	resolver.RegisterProfile(ProfileDefault, CreateDefaultProfile())
	resolver.RegisterProfile(ProfileEmacs, CreateEmacsProfile())
	resolver.RegisterProfile(ProfileVi, CreateViProfile())
	resolver.RegisterProfile(ProfileReadline, CreateReadlineProfile())
}

// GetAllProfilesBuiltin returns all available profile names
func GetAllProfilesBuiltin() []Profile {
	return []Profile{ProfileDefault, ProfileEmacs, ProfileVi, ProfileReadline}
}

// GetProfileDescription returns a description for a profile
func GetProfileDescription(profile Profile) string {
	switch profile {
	case ProfileDefault:
		return "Default keybindings compatible with legacy behavior"
	case ProfileEmacs:
		return "Comprehensive Emacs-style keybindings with authentic GNU Emacs behavior"
	case ProfileVi:
		return "Vi-style modal keybindings adapted for command-line interface with insert and normal modes"
	case ProfileReadline:
		return "Comprehensive GNU Readline compatible keybindings for authentic bash-like CLI experience"
	default:
		return "Unknown profile"
	}
}

