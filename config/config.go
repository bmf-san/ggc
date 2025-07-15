// Package config provides a base configuration schema for ggc.
package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
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

	Behavior struct {
		AutoPush           bool   `yaml:"auto-push"`
		ConfirmDestructive string `yaml:"confirm-destructive"`
		AutoFetch          bool   `yaml:"auto-fetch"`
		StashBeforeSwitch  bool   `yaml:"stash-before-switch"`
	} `yaml:"behavior"`

	Aliases map[string]string `yaml:"aliases"`

	Integration struct {
		Github struct {
			Token         string `yaml:"token"`
			DefaultRemote string `yaml:"default-remote"`
		} `yaml:"github"`
		Gitlab struct {
			Token string `yaml:"token"`
		} `yaml:"gitlab"`
	} `yaml:"integration"`
}

// Manager handles configuration loading, saving, and operations
type Manager struct {
	config     *Config
	configPath string
}

// NewConfigManager creates a new configuration manager
func NewConfigManager() *Manager {
	return &Manager{
		config: getDefaultConfig(),
	}
}

func getGitVersion() string {
	cmd := exec.Command("git", "describe", "--tags", "--always", "--dirty")
	output, err := cmd.Output()
	if err != nil {
		return "dev"
	}
	return strings.TrimSpace(string(output))
}

func getGitCommit() string {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

func (c *Config) updateMeta() {
	version, commit := getGitVersion(), getGitCommit()

	c.Meta.Version = version
	c.Meta.Commit = commit

	if c.Meta.ConfigVersion == "" {
		c.Meta.ConfigVersion = "1.0"
	}
}

// getDefaultConfig returns the default configuration values
func getDefaultConfig() *Config {
	config := &Config{
		Aliases: make(map[string]string),
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

	// Default aliases
	config.Aliases["st"] = "status"
	config.Aliases["co"] = "checkout"
	config.Aliases["br"] = "branch"
	config.Aliases["ci"] = "commit"

	config.Integration.Github.DefaultRemote = "origin"

	config.updateMeta()

	return config
}

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
	paths := cm.getConfigPaths()

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			cm.configPath = path
			return cm.loadFromFile(path)
		}
	}

	cm.configPath = paths[0]
	return nil
}

// loadFromFile loads configuration from a specific file
func (cm *Manager) loadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	config := getDefaultConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	cm.syncFromGitConfig()
	cm.config = config
	return nil
}

// setGitConfig sets a git configuration value using git config command
func (cm *Manager) setGitConfig(key, value string) error {
	cmd := exec.Command("git", "config", "--global", key, value)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run git config: %w", err)
	}
	return nil
}

// getGitConfig retrieves a git configuration value
func (cm *Manager) getGitConfig(key string) (string, error) {
	cmd := exec.Command("git", "config", "--global", key)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git config: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func (cm *Manager) syncFromCommandName(command string) {
	if value, err := cm.getGitConfig(command); err == nil && value != "" {
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
		case "push.default":
			cm.config.Behavior.ConfirmDestructive = value
		}
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

// syncToGitConfig synchronizes relevant config values TO Git's global configuration
func (cm *Manager) syncToGitConfig() error {
	config := cm.GetConfig()

	if config.Default.Editor != "" {
		if err := cm.setGitConfig("core.editor", config.Default.Editor); err != nil {
			return fmt.Errorf("failed to set git editor: %w", err)
		}
	}

	if config.Default.MergeTool != "" {
		if err := cm.setGitConfig("merge.tool", config.Default.MergeTool); err != nil {
			return fmt.Errorf("failed to set git merge tool: %w", err)
		}
	}

	if config.Default.Branch != "" {
		if err := cm.setGitConfig("init.defaultBranch", config.Default.Branch); err != nil {
			return fmt.Errorf("failed to set git default branch: %w", err)
		}
	}

	colorValue := "false"
	if config.UI.Color {
		colorValue = "true"
	}
	if err := cm.setGitConfig("color.ui", colorValue); err != nil {
		return fmt.Errorf("failed to set git color: %w", err)
	}

	autoFetchValue := "false"
	if config.Behavior.AutoFetch {
		autoFetchValue = "true"
	}
	if err := cm.setGitConfig("fetch.auto", autoFetchValue); err != nil {
		return fmt.Errorf("failed to set git autofetch: %w", err)
	}

	if err := cm.setGitConfig("push.default", config.Behavior.ConfirmDestructive); err != nil {
		return fmt.Errorf("failed to set git push default: %w", err)
	}

	if !config.UI.Pager {
		if err := cm.setGitConfig("core.pager", "cat"); err != nil {
			return fmt.Errorf("failed to set git pager: %w", err)
		}
	}

	if len(config.Aliases) != 0 {
		for alias, command := range config.Aliases {
			if err := cm.setGitConfig(fmt.Sprintf("alias.%s", alias), command); err != nil {
				return fmt.Errorf("failed to set git alias.%s: %w", alias, err)
			}
		}
	}

	return nil
}

// Save saves the current configuration to file
func (cm *Manager) Save() error {
	dir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(cm.config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return cm.syncToGitConfig()
}

// Get retrieves a configuration value by key path (e.g., "ui.color", "default.branch")
func (cm *Manager) Get(key string) (any, error) {
	return cm.getValueByPath(cm.config, key)
}

// Set sets a configuration value by key path
func (cm *Manager) Set(key string, value any) error {
	if err := cm.setValueByPath(cm.config, key, value); err != nil {
		return err
	}
	return cm.Save()
}

// List returns all configuration keys and values
func (cm *Manager) List() map[string]any {
	result := make(map[string]any)
	cm.flattenConfig(cm.config, "", result)
	return result
}

// getValueByPath retrieves a value using dot notation path
func (cm *Manager) getValueByPath(obj any, path string) (any, error) {
	parts := strings.Split(path, ".")
	current := reflect.ValueOf(obj)

	for _, part := range parts {
		if current.Kind() == reflect.Ptr {
			current = current.Elem()
		}

		switch current.Kind() {
		case reflect.Struct:
			field, found := cm.findFieldByYamlTag(current.Type(), current, part)
			if !found {
				return nil, fmt.Errorf("field '%s' not found", part)
			}
			current = field

		case reflect.Map:
			mapValue := current.MapIndex(reflect.ValueOf(part))
			if !mapValue.IsValid() {
				return nil, fmt.Errorf("key '%s' not found", part)
			}
			current = mapValue

		default:
			return nil, fmt.Errorf("cannot navigate into %s", current.Kind())
		}
	}

	return current.Interface(), nil
}

// findFieldByYamlTag finds a struct field by its YAML tag or field name
func (cm *Manager) findFieldByYamlTag(structType reflect.Type, structValue reflect.Value, tagName string) (reflect.Value, bool) {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)

		// Check YAML tag first
		yamlTag := field.Tag.Get("yaml")
		if yamlTag != "" {
			yamlName := strings.Split(yamlTag, ",")[0]
			if yamlName == tagName {
				return fieldValue, true
			}
		}

		// Fallback to field name (case-insensitive)
		if strings.EqualFold(field.Name, tagName) {
			return fieldValue, true
		}
	}
	return reflect.Value{}, false
}

// setValueByPath sets a value using dot notation path
func (cm *Manager) setValueByPath(obj any, path string, value any) error {
	parts := strings.Split(path, ".")
	current := reflect.ValueOf(obj)

	// Navigate to the parent of the target field
	for i, part := range parts[:len(parts)-1] {
		if current.Kind() == reflect.Ptr {
			current = current.Elem()
		}

		switch current.Kind() {
		case reflect.Struct:
			field, found := cm.findFieldByYamlTag(current.Type(), current, part)
			if !found {
				return fmt.Errorf("field '%s' not found", strings.Join(parts[:i+1], "."))
			}
			current = field

		case reflect.Map:
			mapValue := current.MapIndex(reflect.ValueOf(part))
			if !mapValue.IsValid() {
				return fmt.Errorf("key '%s' not found", strings.Join(parts[:i+1], "."))
			}
			current = mapValue

		default:
			return fmt.Errorf("cannot navigate into %s", current.Kind())
		}
	}

	// Set the final value
	lastPart := parts[len(parts)-1]
	if current.Kind() == reflect.Ptr {
		current = current.Elem()
	}

	switch current.Kind() {
	case reflect.Struct:
		field, found := cm.findFieldByYamlTag(current.Type(), current, lastPart)
		if !found || !field.CanSet() {
			return fmt.Errorf("field '%s' not found or cannot be set", lastPart)
		}

		newValue := reflect.ValueOf(value)
		if !newValue.Type().ConvertibleTo(field.Type()) {
			return fmt.Errorf("cannot convert %s to %s", newValue.Type(), field.Type())
		}

		field.Set(newValue.Convert(field.Type()))

	case reflect.Map:
		if current.Type().Key().Kind() != reflect.String {
			return fmt.Errorf("map key must be string")
		}

		newValue := reflect.ValueOf(value)
		if !newValue.Type().ConvertibleTo(current.Type().Elem()) {
			return fmt.Errorf("cannot convert %s to %s", newValue.Type(), current.Type().Elem())
		}

		current.SetMapIndex(reflect.ValueOf(lastPart), newValue.Convert(current.Type().Elem()))

	default:
		return fmt.Errorf("cannot set value in %s", current.Kind())
	}

	return nil
}

// flattenConfig converts nested config to flat key-value pairs
func (cm *Manager) flattenConfig(obj any, prefix string, result map[string]any) {
	value := reflect.ValueOf(obj)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Struct:
		structType := value.Type()
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			fieldType := structType.Field(i)

			yamlTag := fieldType.Tag.Get("yaml")
			fieldName := fieldType.Name
			if yamlTag != "" {
				fieldName = strings.Split(yamlTag, ",")[0]
			}

			key := fieldName
			if prefix != "" {
				key = prefix + "." + fieldName
			}

			if field.Kind() == reflect.Struct || (field.Kind() == reflect.Map && field.Type().Elem().Kind() == reflect.String) {
				cm.flattenConfig(field.Interface(), key, result)
			} else {
				result[key] = field.Interface()
			}
		}

	case reflect.Map:
		for _, mapKey := range value.MapKeys() {
			mapValue := value.MapIndex(mapKey)
			key := mapKey.String()
			if prefix != "" {
				key = prefix + "." + mapKey.String()
			}
			result[key] = mapValue.Interface()
		}
	}
}

// LoadConfig loads and saves the configuration file.
func (cm *Manager) LoadConfig() {
	if err := cm.Load(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to save config: %v\n", err)
	}
	if err := cm.Save(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to save config: %v\n", err)
	}
}

// GetConfig returns the current configuration
func (cm *Manager) GetConfig() *Config {
	return cm.config
}
