// Package config provides a base configuration schema for ggc.
package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/bmf-san/ggc/v4/git"
	"gopkg.in/yaml.v3"
)

// GitConfigExecutor interface for git config operations (for testing)
type GitConfigExecutor interface {
	ConfigSetGlobal(key, value string) error
	ConfigGetGlobal(key string) (string, error)
}

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

	Aliases map[string]interface{} `yaml:"aliases"`

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

// AliasType represents the type of alias
type AliasType int

const (
	// SimpleAlias represents a simple string alias
	SimpleAlias AliasType = iota
	// SequenceAlias represents an array string alias
	SequenceAlias
)

// ParsedAlias represents a parsed alias with its type and commands
type ParsedAlias struct {
	Type     AliasType
	Commands []string
}

// Manager handles configuration loading, saving, and operations
type Manager struct {
	config     *Config
	configPath string
	gitClient  git.Clienter
}

// NewConfigManager creates a new configuration manager
func NewConfigManager() *Manager {
	gitClient := git.NewClient()
	return &Manager{
		config:    getDefaultConfig(gitClient),
		gitClient: gitClient,
	}
}

// ValidationError creates a new error manager for validation operations
type ValidationError struct {
	Field   string
	Value   any
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid value for '%s': %v (%s)", e.Field, e.Value, e.Message)
}

// Validator creates an interface for validating config
type Validator interface {
	Validate() error
}

func (c *Config) validateBranch() error {
	branch := c.Default.Branch
	if strings.TrimSpace(branch) == "" || strings.Contains(branch, " ") {
		return &ValidationError{"default.branch", branch, "must not contain spaces or be empty"}
	}
	return nil
}

func (c *Config) validateEditor() error {
	editor := c.Default.Editor
	_, err := exec.LookPath(editor)
	if err != nil {
		return &ValidationError{"default.editor", editor, "command not found in PATH"}
	}
	return nil
}

func (c *Config) validateConfirmDestructive() error {
	val := c.Behavior.ConfirmDestructive
	valid := map[string]bool{"simple": true, "always": true, "never": true}
	if !valid[val] {
		return &ValidationError{"behavior.confirm-destructive", val, "must be one of: simple, always, never"}
	}
	return nil
}

func (c *Config) validateIntegrationTokens() error {
	ghToken := c.Integration.Github.Token
	if ghToken != "" && !strings.HasPrefix(ghToken, "ghp_") {
		return &ValidationError{
			Field:   "integration.github.token",
			Value:   ghToken,
			Message: "GitHub token must start with 'ghp_'",
		}
	}

	glToken := c.Integration.Gitlab.Token
	if glToken != "" && len(glToken) < 20 {
		return &ValidationError{
			Field:   "integration.gitlab.token",
			Value:   glToken,
			Message: "GitLab token seems too short to be valid",
		}
	}

	if remote := c.Integration.Github.DefaultRemote; remote != "" {
		if strings.Contains(remote, " ") || strings.Contains(remote, "/") {
			return &ValidationError{
				Field:   "integration.github.default-remote",
				Value:   remote,
				Message: "Default remote must be a valid Git remote name (no spaces or slashes)",
			}
		}
	}

	return nil
}

func (c *Config) validateAliases() error {
	for name, value := range c.Aliases {
		if strings.TrimSpace(name) == "" || strings.Contains(name, " ") {
			return &ValidationError{"aliases." + name, name, "alias names must not contain spaces"}
		}

		switch v := value.(type) {
		case string:
			// Simple alias validation
			if strings.TrimSpace(v) == "" {
				return &ValidationError{"aliases." + name, v, "alias command cannot be empty"}
			}

		case []interface{}:
			// Sequence alias validation
			if len(v) == 0 {
				return &ValidationError{"aliases." + name, v, "alias sequence cannot be empty"}
			}
			for i, cmd := range v {
				cmdStr, ok := cmd.(string)
				if !ok {
					return &ValidationError{
						Field:   fmt.Sprintf("aliases.%s[%d]", name, i),
						Value:   cmd,
						Message: "sequence commands must be strings",
					}
				}
				if strings.TrimSpace(cmdStr) == "" {
					return &ValidationError{
						Field:   fmt.Sprintf("aliases.%s[%d]", name, i),
						Value:   cmdStr,
						Message: "command in sequence cannot be empty",
					}
				}
			}

		default:
			return &ValidationError{
				Field:   "aliases." + name,
				Value:   value,
				Message: "alias must be either a string or array of strings",
			}
		}
	}
	return nil
}

// Validate is a function that handles validation operations
func (c *Config) Validate() error {
	if err := c.validateBranch(); err != nil {
		return err
	}
	if err := c.validateEditor(); err != nil {
		return err
	}
	if err := c.validateConfirmDestructive(); err != nil {
		return err
	}
	if err := c.validateIntegrationTokens(); err != nil {
		return err
	}
	if err := c.validateAliases(); err != nil {
		return err
	}
	return nil
}

// ParseAlias parses an alias value and returns its type and commands
func (c *Config) ParseAlias(name string) (*ParsedAlias, error) {
	value, exists := c.Aliases[name]
	if !exists {
		return nil, fmt.Errorf("alias '%s' not found", name)
	}

	switch v := value.(type) {
	case string:
		return &ParsedAlias{
			Type:     SimpleAlias,
			Commands: []string{v},
		}, nil

	case []interface{}:
		commands := make([]string, len(v))
		for i, cmd := range v {
			cmdStr, ok := cmd.(string)
			if !ok {
				return nil, fmt.Errorf("invalid command type in alias '%s'", name)
			}
			commands[i] = cmdStr
		}
		return &ParsedAlias{
			Type:     SequenceAlias,
			Commands: commands,
		}, nil

	default:
		return nil, fmt.Errorf("invalid alias type for '%s'", name)
	}
}

// IsAlias checks if a given name is an alias
func (c *Config) IsAlias(name string) bool {
	_, exists := c.Aliases[name]
	return exists
}

// GetAliasCommands returns the commands for a given alias
func (c *Config) GetAliasCommands(name string) ([]string, error) {
	parsed, err := c.ParseAlias(name)
	if err != nil {
		return nil, err
	}
	return parsed.Commands, nil
}

// GetAllAliases returns all aliases with their parsed commands
func (c *Config) GetAllAliases() map[string]*ParsedAlias {
	result := make(map[string]*ParsedAlias)
	for name := range c.Aliases {
		if parsed, err := c.ParseAlias(name); err == nil {
			result[name] = parsed
		}
	}
	return result
}

// Note: getGitVersion, getGitCommit, and updateMeta functions removed
// to eliminate direct git command execution. Meta values are now set
// manually in getDefaultConfig() to avoid side effects in tests.

// getDefaultConfig returns the default configuration values
func getDefaultConfig(gitClient git.Clienter) *Config {
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

	config.Integration.Github.DefaultRemote = "origin"

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

	if err := cm.config.Validate(); err != nil {
		return err
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

	config := getDefaultConfig(cm.gitClient)
	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	cm.syncFromGitConfig()
	cm.config = config
	return nil
}

func (cm *Manager) syncFromCommandName(command string) {
	if value, err := cm.gitClient.ConfigGetGlobal(command); err == nil && value != "" {
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

	colorValue := "false"
	if config.UI.Color {
		colorValue = "true"
	}
	if err := cm.gitClient.ConfigSetGlobal("color.ui", colorValue); err != nil {
		return fmt.Errorf("failed to set git color: %w", err)
	}

	autoFetchValue := "false"
	if config.Behavior.AutoFetch {
		autoFetchValue = "true"
	}
	if err := cm.gitClient.ConfigSetGlobal("fetch.auto", autoFetchValue); err != nil {
		return fmt.Errorf("failed to set git autofetch: %w", err)
	}

	if err := cm.gitClient.ConfigSetGlobal("push.default", config.Behavior.ConfirmDestructive); err != nil {
		return fmt.Errorf("failed to set git push default: %w", err)
	}

	if !config.UI.Pager {
		if err := cm.gitClient.ConfigSetGlobal("core.pager", "cat"); err != nil {
			return fmt.Errorf("failed to set git pager: %w", err)
		}
	}

	for alias, value := range config.Aliases {
		if cmdStr, ok := value.(string); ok {
			if err := cm.gitClient.ConfigSetGlobal(fmt.Sprintf("alias.%s", alias), cmdStr); err != nil {
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

	if err := cm.config.Validate(); err != nil {
		return fmt.Errorf("cannot save invalid config: %w", err)
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
	if err := cm.config.Validate(); err != nil {
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

			if field.Kind() == reflect.Struct || (field.Kind() == reflect.Map && field.Type().Elem().Kind() != reflect.Interface) {
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
		_, _ = fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
	}
	if err := cm.Save(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to save config: %v\n", err)
	}
}

// GetConfig returns the current configuration
func (cm *Manager) GetConfig() *Config {
	return cm.config
}
