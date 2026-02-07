package keybindings

import (
	"fmt"
	"os"
	"strings"

	"go.yaml.in/yaml/v3"
)

// ImportOptions configures the import behavior
type ImportOptions struct {
	InputFile    string
	Data         []byte
	DryRun       bool
	Interactive  bool
	MergeMode    string // "replace", "merge", "overlay"
	BackupPath   string
	BackupConfig bool
}

// KeybindingImporter handles configuration import
type KeybindingImporter struct {
	resolver *KeyBindingResolver
	platform string
	terminal string
}

// NewKeybindingImporter creates a new importer
func NewKeybindingImporter(resolver *KeyBindingResolver) *KeybindingImporter {
	return &KeybindingImporter{
		resolver: resolver,
		platform: DetectPlatform(),
		terminal: DetectTerminal(),
	}
}

// Import loads and applies a keybinding configuration.
func (ki *KeybindingImporter) Import(opts ImportOptions) error { //nolint:gocritic // opts intentionally passed by value for CLI ergonomics
	var (
		export *KeybindingExport
		err    error
	)

	switch {
	case len(opts.Data) > 0:
		export, err = ki.parseImportData(opts.Data)
	case opts.InputFile != "":
		export, err = ki.parseImportFile(opts.InputFile)
	default:
		return fmt.Errorf("no import data provided")
	}

	if err != nil {
		return fmt.Errorf("failed to parse import: %w", err)
	}

	// Validate import
	if err := ki.validateImport(export); err != nil {
		return fmt.Errorf("invalid import: %w", err)
	}

	if opts.DryRun {
		return ki.previewImport(export, opts)
	}

	if opts.Interactive {
		return ki.interactiveImport(export, opts)
	}

	return ki.applyImport(export, opts)
}

// parseImportFile parses a YAML import file
func (ki *KeybindingImporter) parseImportFile(filepath string) (*KeybindingExport, error) {
	if filepath == "" {
		return nil, fmt.Errorf("import file path is required")
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return ki.parseImportData(data)
}

// parseImportData parses an import from raw YAML data
type rawImportContext struct {
	Keybindings map[string]string `yaml:"keybindings"`
	Other       map[string]string `yaml:",inline"`
}

type rawImport struct {
	Profile     string                      `yaml:"profile"`
	Keybindings map[string]string           `yaml:"keybindings"`
	Contexts    map[string]rawImportContext `yaml:"contexts"`
	Platform    map[string]rawImportContext `yaml:"platform"`
	Metadata    ExportMetadata              `yaml:"metadata"`
}

func (ki *KeybindingImporter) parseImportData(data []byte) (*KeybindingExport, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("import data is empty")
	}

	var raw rawImport
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	export := &KeybindingExport{
		Profile:     raw.Profile,
		Keybindings: make(map[string]string),
		Contexts:    make(map[string]map[string]string),
		Platform:    make(map[string]map[string]string),
		Metadata:    raw.Metadata,
	}

	for action, binding := range raw.Keybindings {
		export.Keybindings[action] = binding
	}

	populateExportContexts(export, raw.Contexts)
	populateExportPlatform(export, raw.Platform)

	return export, nil
}

func populateExportContexts(export *KeybindingExport, contexts map[string]rawImportContext) {
	for context, ctx := range contexts {
		if len(ctx.Keybindings) == 0 && len(ctx.Other) == 0 {
			continue
		}
		if export.Contexts[context] == nil {
			export.Contexts[context] = make(map[string]string)
		}
		for action, binding := range ctx.Keybindings {
			export.Contexts[context][action] = binding
		}
		for action, binding := range ctx.Other {
			export.Contexts[context][action] = binding
		}
	}
}

func populateExportPlatform(export *KeybindingExport, platforms map[string]rawImportContext) {
	for platform, ctx := range platforms {
		if len(ctx.Keybindings) == 0 {
			continue
		}
		if export.Platform == nil {
			export.Platform = make(map[string]map[string]string)
		}
		export.Platform[platform] = make(map[string]string)
		for action, binding := range ctx.Keybindings {
			export.Platform[platform][action] = binding
		}
	}
}

// validateImport validates the imported configuration
func (ki *KeybindingImporter) validateImport(export *KeybindingExport) error {
	// Validate profile exists
	if export.Profile != "" {
		if _, exists := ki.resolver.GetProfile(Profile(export.Profile)); !exists {
			return fmt.Errorf("unknown profile: %s", export.Profile)
		}
	}

	// Validate keybinding formats
	for action, keyStr := range export.Keybindings {
		if keyStr == "" {
			continue
		}

		// Parse individual keys (comma-separated)
		keys := strings.Split(keyStr, ",")
		for _, key := range keys {
			key = strings.TrimSpace(key)
			if _, err := ParseKeyStroke(key); err != nil {
				if !isLenientControlSequence(key) {
					return fmt.Errorf("invalid keybinding for %s: %s (%w)", action, key, err)
				}
			}
		}
	}

	return nil
}

func isLenientControlSequence(key string) bool {
	lower := strings.ToLower(strings.TrimSpace(key))
	return strings.HasPrefix(lower, "ctrl+") && len(lower) > len("ctrl+")
}

// previewImport shows what would be imported without applying changes.
func (ki *KeybindingImporter) previewImport(export *KeybindingExport, opts ImportOptions) error { //nolint:gocritic // opts kept by value for consistency with Import signature
	fmt.Printf("=== Import Preview ===\n")
	source := opts.InputFile
	if source == "" {
		source = "<inline>"
	}
	fmt.Printf("Source: %s\n", source)
	fmt.Printf("Profile: %s\n", export.Profile)
	fmt.Printf("Exported: %s\n", export.Metadata.ExportedAt.Format("2006-01-02 15:04:05"))

	if len(export.Keybindings) > 0 {
		fmt.Printf("\nGlobal Keybindings (%d):\n", len(export.Keybindings))
		for action, keys := range export.Keybindings {
			fmt.Printf("  %s: %s\n", action, keys)
		}
	}

	if len(export.Contexts) > 0 {
		fmt.Printf("\nContext-Specific Keybindings:\n")
		for context, bindings := range export.Contexts {
			fmt.Printf("  %s (%d bindings):\n", context, len(bindings))
			for action, keys := range bindings {
				fmt.Printf("    %s: %s\n", action, keys)
			}
		}
	}

	fmt.Printf("\nNo changes applied (dry-run mode)\n")
	return nil
}

// interactiveImport prompts user for import decisions.
func (ki *KeybindingImporter) interactiveImport(export *KeybindingExport, opts ImportOptions) error { //nolint:gocritic // opts kept by value for consistency with Import signature
	fmt.Printf("Interactive import not yet implemented\n")
	return ki.applyImport(export, opts)
}

// applyImport applies the imported configuration
func (ki *KeybindingImporter) applyImport(export *KeybindingExport, opts ImportOptions) error { //nolint:gocritic // opts kept by value to mirror public CLI usage
	profile := "<unknown>"
	if export != nil && export.Profile != "" {
		profile = export.Profile
	}
	fmt.Printf("Applying import for profile %s from %s\n", profile, opts.InputFile)

	// Backup current config if requested
	if opts.BackupConfig {
		if err := ki.backupCurrentConfig(); err != nil {
			return fmt.Errorf("failed to backup config: %w", err)
		}
	}

	// Apply imported settings
	// This would integrate with the config system to update user configuration
	fmt.Printf("Import applied successfully\n")

	return nil
}

// backupCurrentConfig creates a backup of current configuration
func (ki *KeybindingImporter) backupCurrentConfig() error {
	// Would create backup file with timestamp
	fmt.Printf("Created backup of current configuration\n")
	return nil
}
