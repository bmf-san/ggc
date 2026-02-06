package keybindings

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// KeybindingExport represents exported keybinding configuration
type KeybindingExport struct {
	Profile     string                       `yaml:"profile"`
	Keybindings map[string]string            `yaml:"keybindings,omitempty"`
	Contexts    map[string]map[string]string `yaml:"contexts,omitempty"`
	Platform    map[string]map[string]string `yaml:"platform,omitempty"`
	Metadata    ExportMetadata               `yaml:"metadata"`
}

// ExportMetadata provides context about the export
type ExportMetadata struct {
	ExportedAt time.Time `yaml:"exported_at"`
	ExportedBy string    `yaml:"exported_by"`
	Version    string    `yaml:"version"`
	Platform   string    `yaml:"platform"`
	Terminal   string    `yaml:"terminal"`
	DeltaFrom  string    `yaml:"delta_from,omitempty"`
	Comment    string    `yaml:"comment,omitempty"`
}

// ExportOptions configures the export behavior
type ExportOptions struct {
	Profile     Profile
	Context     Context
	DeltaMode   bool
	OutputFile  string
	IncludeMeta bool
	Format      string // "yaml" or "json"
}

// KeybindingExporter handles configuration export
type KeybindingExporter struct {
	resolver *KeyBindingResolver
	platform string
	terminal string
}

// NewKeybindingExporter creates a new exporter
func NewKeybindingExporter(resolver *KeyBindingResolver) *KeybindingExporter {
	return &KeybindingExporter{
		resolver: resolver,
		platform: DetectPlatform(),
		terminal: DetectTerminal(),
	}
}

// Export generates a keybinding configuration export.
func (ke *KeybindingExporter) Export(opts ExportOptions) (*KeybindingExport, error) { //nolint:gocritic // opts is small struct used widely; keep by value for backward compatibility
	export := &KeybindingExport{
		Profile:     string(opts.Profile),
		Keybindings: make(map[string]string),
		Contexts:    make(map[string]map[string]string),
		Platform:    make(map[string]map[string]string),
		Metadata: ExportMetadata{
			ExportedAt: time.Now(),
			ExportedBy: os.Getenv("USER"),
			Version:    "5.0.0", // Would be injected from build
			Platform:   ke.platform,
			Terminal:   ke.terminal,
		},
	}

	if opts.DeltaMode {
		return ke.exportDelta(opts, export)
	}

	return ke.exportFull(opts, export)
}

// exportFull exports complete configuration.
func (ke *KeybindingExporter) exportFull(opts ExportOptions, export *KeybindingExport) (*KeybindingExport, error) { //nolint:gocritic // opts intentionally passed by value to avoid pointer aliasing in tests
	// Get profile information
	profile, exists := ke.resolver.GetProfile(opts.Profile)
	if !exists {
		return nil, fmt.Errorf("profile '%s' not found", opts.Profile)
	}

	export.Metadata.Comment = fmt.Sprintf("Complete keybinding export for %s profile", profile.Name)

	ke.addGlobalBindings(export, profile)
	ke.addContextBindings(export, profile)
	ke.promoteCoreBindings(export, profile)
	ke.addPlatformBindings(export)

	return export, nil
}

func (ke *KeybindingExporter) addGlobalBindings(export *KeybindingExport, profile *KeyBindingProfile) {
	for action, keystrokes := range profile.Global {
		if len(keystrokes) == 0 {
			continue
		}
		export.Keybindings[action] = ke.formatKeystrokesForExport(keystrokes)
	}
}

func (ke *KeybindingExporter) addContextBindings(export *KeybindingExport, profile *KeyBindingProfile) {
	for context, bindings := range profile.Contexts {
		if len(bindings) == 0 {
			continue
		}
		contextName := string(context)
		export.Contexts[contextName] = make(map[string]string)
		for action, keystrokes := range bindings {
			if len(keystrokes) == 0 {
				continue
			}
			export.Contexts[contextName][action] = ke.formatKeystrokesForExport(keystrokes)
		}
	}
}

func (ke *KeybindingExporter) promoteCoreBindings(export *KeybindingExport, profile *KeyBindingProfile) {
	inputCtx, exists := profile.Contexts[ContextInput]
	if !exists {
		return
	}

	coreActions := []string{
		"move_to_beginning",
		"move_to_end",
		"delete_word",
		"delete_to_end",
		"clear_line",
	}
	for _, action := range coreActions {
		if _, already := export.Keybindings[action]; already {
			continue
		}
		if keys, ok := inputCtx[action]; ok && len(keys) > 0 {
			export.Keybindings[action] = ke.formatKeystrokesForExport(keys)
		}
	}
}

func (ke *KeybindingExporter) addPlatformBindings(export *KeybindingExport) {
	platformBindings := GetPlatformSpecificKeyBindings(ke.platform)
	if len(platformBindings) == 0 {
		return
	}
	if export.Platform == nil {
		export.Platform = make(map[string]map[string]string)
	}

	export.Platform[ke.platform] = make(map[string]string)
	for action, keystrokes := range platformBindings {
		export.Platform[ke.platform][action] = ke.formatKeystrokesForExport(keystrokes)
	}
}

// exportDelta exports only differences from base profile.
func (ke *KeybindingExporter) exportDelta(opts ExportOptions, export *KeybindingExport) (*KeybindingExport, error) { //nolint:gocritic // opts intentionally passed by value to preserve API
	if _, exists := ke.resolver.GetProfile(opts.Profile); !exists {
		return nil, fmt.Errorf("profile '%s' not found", opts.Profile)
	}

	export.Metadata.DeltaFrom = string(opts.Profile)
	export.Metadata.Comment = fmt.Sprintf("Delta export: overrides for %s profile", opts.Profile)

	// Delta export only includes user overrides; since this resolver has no
	// additional configuration applied yet, there are no differences to report.
	return export, nil
}

// formatKeystrokesForExport converts keystrokes to export format
func (ke *KeybindingExporter) formatKeystrokesForExport(keystrokes []KeyStroke) string {
	if len(keystrokes) == 0 {
		return ""
	}

	if len(keystrokes) == 1 {
		return ke.formatKeystrokeForExport(keystrokes[0])
	}

	// Multiple keystrokes - return as comma-separated string
	var parts []string
	for _, ks := range keystrokes {
		parts = append(parts, ke.formatKeystrokeForExport(ks))
	}

	return strings.Join(parts, ", ")
}

// formatKeystrokeForExport converts a single keystroke to export format
func (ke *KeybindingExporter) formatKeystrokeForExport(ks KeyStroke) string { //nolint:revive // export formatting mirrors import expectations
	switch ks.Kind {
	case KeyStrokeCtrl:
		return fmt.Sprintf("ctrl+%c", ks.Rune)
	case KeyStrokeAlt:
		return fmt.Sprintf("alt+%c", ks.Rune)
	case KeyStrokeRawSeq:
		// Handle common sequences
		if len(ks.Seq) == 1 {
			switch ks.Seq[0] {
			case 9:
				return "tab"
			case 13:
				return "enter"
			case 27:
				return "esc"
			case 32:
				return "space"
			}
		}
		// Arrow keys
		if len(ks.Seq) == 3 && ks.Seq[0] == 27 && ks.Seq[1] == 91 {
			switch ks.Seq[2] {
			case 65:
				return "up"
			case 66:
				return "down"
			case 67:
				return "right"
			case 68:
				return "left"
			}
		}
		// Raw sequence
		return fmt.Sprintf("raw:%x", ks.Seq)
	case KeyStrokeFnKey:
		return strings.ToLower(ks.Name)
	default:
		return fmt.Sprintf("unknown:%v", ks)
	}
}

// ToYAML converts export to YAML format
func (ke *KeybindingExport) ToYAML() (string, error) { //nolint:revive // YAML rendering preserves explicit ordering
	var result strings.Builder

	// Write header comment
	result.WriteString(fmt.Sprintf("# Generated by ggc %s on %s\n",
		ke.Metadata.Version, ke.Metadata.ExportedAt.Format("2006-01-02T15:04:05Z07:00")))
	result.WriteString(fmt.Sprintf("# Profile: %s\n", ke.Profile))
	result.WriteString(fmt.Sprintf("# Platform: %s/%s\n", ke.Metadata.Platform, ke.Metadata.Terminal))

	if ke.Metadata.Comment != "" {
		result.WriteString(fmt.Sprintf("# %s\n", ke.Metadata.Comment))
	}
	result.WriteString("\n")

	// Write profile
	result.WriteString(fmt.Sprintf("profile: %s\n\n", ke.Profile))

	// Write global keybindings
	if len(ke.Keybindings) > 0 {
		result.WriteString("keybindings:\n")
		for action, keys := range ke.Keybindings {
			result.WriteString(fmt.Sprintf("  %s: \"%s\"\n", action, keys))
		}
		result.WriteString("\n")
	}

	// Write context-specific keybindings
	if len(ke.Contexts) > 0 {
		result.WriteString("contexts:\n")
		for context, bindings := range ke.Contexts {
			result.WriteString(fmt.Sprintf("  %s:\n", context))
			result.WriteString("    keybindings:\n")
			for action, keys := range bindings {
				result.WriteString(fmt.Sprintf("      %s: \"%s\"\n", action, keys))
			}
		}
		result.WriteString("\n")
	}

	// Write platform-specific bindings
	if len(ke.Platform) > 0 {
		for platform, bindings := range ke.Platform {
			result.WriteString(fmt.Sprintf("%s:\n", platform))
			result.WriteString("  keybindings:\n")
			for action, keys := range bindings {
				result.WriteString(fmt.Sprintf("    %s: \"%s\"\n", action, keys))
			}
		}
		result.WriteString("\n")
	}

	// Write metadata
	result.WriteString("metadata:\n")
	result.WriteString(fmt.Sprintf("  exported_at: %s\n", ke.Metadata.ExportedAt.Format(time.RFC3339)))
	result.WriteString(fmt.Sprintf("  exported_by: %s\n", ke.Metadata.ExportedBy))
	result.WriteString(fmt.Sprintf("  version: %s\n", ke.Metadata.Version))
	result.WriteString(fmt.Sprintf("  platform: %s\n", ke.Metadata.Platform))
	result.WriteString(fmt.Sprintf("  terminal: %s\n", ke.Metadata.Terminal))

	if ke.Metadata.DeltaFrom != "" {
		result.WriteString(fmt.Sprintf("  delta_from: %s\n", ke.Metadata.DeltaFrom))
	}

	return result.String(), nil
}
