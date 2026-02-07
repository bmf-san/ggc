package config

import (
	"fmt"
	"path/filepath"
	"runtime"

	"go.yaml.in/yaml/v3"
)

// Save writes the configuration using restrictive permissions to prevent token disclosure.
func (cm *Manager) Save() error {
	return cm.SaveWithFileOps(OSFileOps{})
}

// SaveWithFileOps saves configuration with custom file operations (for testing)
func (cm *Manager) SaveWithFileOps(fileOps FileOps) error {
	dir := filepath.Dir(cm.configPath)
	if err := fileOps.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	data, err := yaml.Marshal(cm.config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	if err := cm.config.Validate(); err != nil {
		return fmt.Errorf("cannot save invalid config: %w", err)
	}
	tmpName, err := cm.writeTempConfigWithOps(dir, data, fileOps)
	if err != nil {
		return err
	}
	if err := cm.replaceConfigFileWithOps(tmpName, fileOps); err != nil {
		return err
	}
	cm.hardenPermissionsWithOps(cm.configPath, fileOps)
	return cm.syncToGitConfig()
}

func (cm *Manager) writeTempConfig(dir string, data []byte) (string, error) {
	return cm.writeTempConfigWithOps(dir, data, OSFileOps{})
}

func (cm *Manager) writeTempConfigWithOps(dir string, data []byte, fileOps FileOps) (string, error) {
	tmpFile, err := fileOps.CreateTemp(dir, ".ggcconfig-*.tmp")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpName := tmpFile.Name()
	if runtime.GOOS != "windows" {
		_ = fileOps.Chmod(tmpName, 0600)
	}
	if _, err := tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		_ = fileOps.Remove(tmpName)
		return "", fmt.Errorf("failed to write temp config file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		_ = fileOps.Remove(tmpName)
		return "", fmt.Errorf("failed to close temp config file: %w", err)
	}
	return tmpName, nil
}

func (cm *Manager) replaceConfigFile(tmpName string) error {
	return cm.replaceConfigFileWithOps(tmpName, OSFileOps{})
}

func (cm *Manager) replaceConfigFileWithOps(tmpName string, fileOps FileOps) error {
	// On Windows versions prior to 10 1903, os.Rename requires the target not exist, so remove it first.
	// Modern Windows (10 1903+) and Unix support atomic replacement with os.Rename, but we remove the target for compatibility.
	if runtime.GOOS == "windows" {
		_ = fileOps.Remove(cm.configPath)
	}

	// Perform atomic rename - this is the only file system operation that matters
	if err := fileOps.Rename(tmpName, cm.configPath); err != nil {
		_ = fileOps.Remove(tmpName)
		return fmt.Errorf("failed to replace config file: %w", err)
	}
	return nil
}

func (cm *Manager) hardenPermissionsWithOps(path string, fileOps FileOps) {
	if runtime.GOOS != "windows" {
		_ = fileOps.Chmod(path, 0600)
	}
}

// syncToGitConfig synchronizes relevant config values TO Git's global configuration
func (cm *Manager) syncToGitConfig() error {
	config := cm.GetConfig()

	if err := cm.syncDefaultSettings(config); err != nil {
		return err
	}

	if err := cm.syncUISettings(config); err != nil {
		return err
	}

	if err := cm.syncBehaviorSettings(config); err != nil {
		return err
	}

	return cm.syncAliasSettings(config)
}

// syncDefaultSettings syncs default editor, merge tool, and branch settings
func (cm *Manager) syncDefaultSettings(config *Config) error {
	if config.Default.Editor != "" {
		if err := cm.gitClient.ConfigSetGlobal("core.editor", config.Default.Editor); err != nil {
			return fmt.Errorf("failed to set git editor: %w", err)
		}
	}

	if config.Default.MergeTool != "" {
		if err := cm.gitClient.ConfigSetGlobal("merge.tool", config.Default.MergeTool); err != nil {
			return fmt.Errorf("failed to set git merge tool: %w", err)
		}
	}

	if config.Default.Branch != "" {
		if err := cm.gitClient.ConfigSetGlobal("init.defaultBranch", config.Default.Branch); err != nil {
			return fmt.Errorf("failed to set git default branch: %w", err)
		}
	}

	return nil
}

// syncUISettings syncs color and pager settings
func (cm *Manager) syncUISettings(config *Config) error {
	colorValue := "false"
	if config.UI.Color {
		colorValue = "true"
	}
	if err := cm.gitClient.ConfigSetGlobal("color.ui", colorValue); err != nil {
		return fmt.Errorf("failed to set git color: %w", err)
	}

	if !config.UI.Pager {
		if err := cm.gitClient.ConfigSetGlobal("core.pager", "cat"); err != nil {
			return fmt.Errorf("failed to set git pager: %w", err)
		}
	}

	return nil
}

// syncBehaviorSettings syncs autofetch settings
func (cm *Manager) syncBehaviorSettings(config *Config) error {
	autoFetchValue := "false"
	if config.Behavior.AutoFetch {
		autoFetchValue = "true"
	}
	if err := cm.gitClient.ConfigSetGlobal("fetch.auto", autoFetchValue); err != nil {
		return fmt.Errorf("failed to set git autofetch: %w", err)
	}
	return nil
}

// syncAliasSettings syncs alias settings to git config
func (cm *Manager) syncAliasSettings(config *Config) error {
	for alias, value := range config.Aliases {
		if cmdStr, ok := value.(string); ok {
			if err := cm.gitClient.ConfigSetGlobal(fmt.Sprintf("alias.%s", alias), cmdStr); err != nil {
				return fmt.Errorf("failed to set git alias.%s: %w", alias, err)
			}
		}
	}

	return nil
}
