// Package config provides a base configuration schema for ggc.
package config

import (
	"regexp"

	"github.com/bmf-san/ggc/v7/pkg/git"
)

var (
	// Allow slashes; additional structural checks are applied separately
	gitRemoteNameCharsRe = regexp.MustCompile(`^[A-Za-z0-9._/\-]+$`)
	configPathSegmentRe  = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

	aliasPlaceholderPattern = regexp.MustCompile(`\{([^}]*)\}`)
)

// Config represents the complete configuration structure
type Config struct {
	Meta struct {
		Version       string `yaml:"version"`
		Commit        string `yaml:"commit"`
		CreatedAt     string `yaml:"created-at"`
		ConfigVersion string `yaml:"config-version"`
	} `yaml:"meta"`

	Default struct {
		Branch    string `yaml:"branch"`
		Editor    string `yaml:"editor"`
		MergeTool string `yaml:"merge-tool"`
	} `yaml:"default"`

	UI struct {
		Color bool `yaml:"color"`
		Pager bool `yaml:"pager"`
	} `yaml:"ui"`

	Interactive struct {
		Profile string `yaml:"profile,omitempty"`

		Keybindings struct {
			DeleteWord         string `yaml:"delete_word"`
			ClearLine          string `yaml:"clear_line"`
			DeleteToEnd        string `yaml:"delete_to_end"`
			MoveToBeginning    string `yaml:"move_to_beginning"`
			MoveToEnd          string `yaml:"move_to_end"`
			MoveUp             string `yaml:"move_up"`
			MoveDown           string `yaml:"move_down"`
			MoveLeft           string `yaml:"move_left"`
			MoveRight          string `yaml:"move_right"`
			AddToWorkflow      string `yaml:"add_to_workflow"`
			ToggleWorkflowView string `yaml:"toggle_workflow_view"`
			ClearWorkflow      string `yaml:"clear_workflow"`
			WorkflowCreate     string `yaml:"workflow_create"`
			WorkflowDelete     string `yaml:"workflow_delete"`
			SoftCancel         string `yaml:"soft_cancel"`
		} `yaml:"keybindings"`

		Contexts struct {
			Input   KeybindingsConfig `yaml:"input,omitempty"`
			Results KeybindingsConfig `yaml:"results,omitempty"`
			Search  KeybindingsConfig `yaml:"search,omitempty"`
		} `yaml:"contexts,omitempty"`

		Darwin  KeybindingsConfig `yaml:"darwin,omitempty"`
		Linux   KeybindingsConfig `yaml:"linux,omitempty"`
		Windows KeybindingsConfig `yaml:"windows,omitempty"`

		Terminals map[string]KeybindingsConfig `yaml:"terminals,omitempty"`
	} `yaml:"interactive"`

	Behavior struct {
		AutoPush           bool   `yaml:"auto-push"`
		ConfirmDestructive string `yaml:"confirm-destructive"`
		AutoFetch          bool   `yaml:"auto-fetch"`
		StashBeforeSwitch  bool   `yaml:"stash-before-switch"`
	} `yaml:"behavior"`

	Aliases map[string]interface{} `yaml:"aliases"`

	Git struct {
		DefaultRemote string `yaml:"default-remote"`
	} `yaml:"git"`
}

// Manager handles configuration loading, saving, and operations
type Manager struct {
	config     *Config
	configPath string
	gitClient  git.ConfigOps
}

// NewConfigManager creates a new configuration manager with the provided git client
func NewConfigManager(gitClient git.ConfigOps) *Manager {
	return &Manager{
		config:    getDefaultConfig(gitClient),
		gitClient: gitClient,
	}
}

// GetConfig returns the current configuration
func (cm *Manager) GetConfig() *Config {
	return cm.config
}

// getDefaultConfig returns the default configuration values
func getDefaultConfig(gitClient git.ConfigOps) *Config {
	config := &Config{
		Aliases: make(map[string]interface{}),
	}

	// Set default values
	config.Default.Branch = "main"
	config.Default.Editor = "vim"
	config.Default.MergeTool = "vimdiff"

	config.UI.Color = true
	config.UI.Pager = true

	config.Behavior.AutoPush = false
	config.Behavior.ConfirmDestructive = "simple"
	config.Behavior.AutoFetch = true
	config.Behavior.StashBeforeSwitch = true

	config.Git.DefaultRemote = "origin"

	// Set meta values using gitClient
	if version, err := gitClient.GetVersion(); err == nil {
		config.Meta.Version = version
	} else {
		config.Meta.Version = "dev"
	}

	if commit, err := gitClient.GetCommitHash(); err == nil {
		config.Meta.Commit = commit
	} else {
		config.Meta.Commit = "unknown"
	}

	// Only set ConfigVersion if it's empty (preserving original updateMeta behavior)
	if config.Meta.ConfigVersion == "" {
		config.Meta.ConfigVersion = "1.0"
	}

	return config
}
