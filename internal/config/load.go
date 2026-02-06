package config

import (
	"fmt"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v3"
)

// getConfigPaths returns possible configuration file paths in order of priority
func (cm *Manager) getConfigPaths() []string {
	homeDir, _ := os.UserHomeDir()

	return []string{
		filepath.Join(homeDir, ".ggcconfig.yaml"),               // Home directory
		filepath.Join(homeDir, ".config", "ggc", "config.yaml"), // XDG config
	}
}

// Load loads configuration from the first available config file
func (cm *Manager) Load() error {
	return cm.LoadWithFileOps(OSFileOps{})
}

// LoadWithFileOps loads configuration with custom file operations (for testing)
func (cm *Manager) LoadWithFileOps(fileOps FileOps) error {
	paths := cm.getConfigPaths()

	for _, path := range paths {
		if _, err := fileOps.Stat(path); err == nil {
			cm.configPath = path
			return cm.loadFromFileWithOps(path, fileOps)
		}
	}

	if err := cm.config.Validate(); err != nil {
		return err
	}
	cm.configPath = paths[0]
	return nil
}

// loadFromFile loads configuration from a specific file
func (cm *Manager) loadFromFile(path string) error {
	return cm.loadFromFileWithOps(path, OSFileOps{})
}

// loadFromFileWithOps loads configuration from a specific file with custom file operations
func (cm *Manager) loadFromFileWithOps(path string, fileOps FileOps) error {
	data, err := fileOps.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	config := getDefaultConfig(cm.gitClient)
	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	cm.syncFromGitConfig()
	cm.config = config
	return nil
}

func (cm *Manager) syncFromCommandName(command string) {
	value, err := cm.gitClient.ConfigGetGlobal(command)
	if err != nil || value == "" {
		return
	}

	switch command {
	case "core.editor":
		cm.config.Default.Editor = value
	case "merge.tool":
		cm.config.Default.MergeTool = value
	case "init.defaultBranch":
		cm.config.Default.Branch = value
	case "color.ui":
		cm.config.UI.Color = value == "true" || value == "auto"
	case "core.pager":
		cm.config.UI.Pager = value != "cat"
	case "fetch.auto":
		cm.config.Behavior.AutoFetch = value == "true"
	}
}

// syncFromGitConfig imports relevant Git config values into your app config
func (cm *Manager) syncFromGitConfig() {
	commands := []string{
		"core.editor",
		"merge.tool",
		"init.defaultBranch",
		"color.ui",
		"core.pager",
	}

	for _, command := range commands {
		cm.syncFromCommandName(command)
	}
}

// LoadConfig loads and saves the configuration file.
// Returns an error if loading or saving fails.
func (cm *Manager) LoadConfig() error {
	if err := cm.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	if err := cm.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	return nil
}
