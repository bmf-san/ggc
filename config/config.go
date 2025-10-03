// Package config provides a base configuration schema for ggc.
package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"go.yaml.in/yaml/v3"

	"github.com/bmf-san/ggc/v7/git"
)

// TempFile interface for temporary file operations
type TempFile interface {
	Write([]byte) (int, error)
	Close() error
	Name() string
}

// FileOps interface for file operations (for testability)
type FileOps interface {
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
	Stat(name string) (os.FileInfo, error)
	MkdirAll(path string, perm os.FileMode) error
	CreateTemp(dir, pattern string) (TempFile, error)
	Remove(name string) error
	Rename(oldpath, newpath string) error
	Chmod(name string, mode os.FileMode) error
}

// OSFileOps implements FileOps using real OS operations
type OSFileOps struct{}

// ReadFile reads a file from the filesystem
func (OSFileOps) ReadFile(filename string) ([]byte, error) { return os.ReadFile(filename) }

// WriteFile writes data to a file
func (OSFileOps) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

// Stat returns file information
func (OSFileOps) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }

// MkdirAll creates directories recursively
func (OSFileOps) MkdirAll(path string, perm os.FileMode) error { return os.MkdirAll(path, perm) }

// CreateTemp creates a temporary file
func (OSFileOps) CreateTemp(dir, pattern string) (TempFile, error) {
	return os.CreateTemp(dir, pattern)
}

// Remove removes a file
func (OSFileOps) Remove(name string) error { return os.Remove(name) }

// Rename renames a file
func (OSFileOps) Rename(oldpath, newpath string) error { return os.Rename(oldpath, newpath) }

// Chmod changes file permissions
func (OSFileOps) Chmod(name string, mode os.FileMode) error { return os.Chmod(name, mode) }

// KeybindingsConfig represents a configuration section with keybindings
type KeybindingsConfig struct {
	Keybindings map[string]interface{} `yaml:"keybindings,omitempty"`
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
			AddToWorkflow      string `yaml:"add_to_workflow"`
			ToggleWorkflowView string `yaml:"toggle_workflow_view"`
			ClearWorkflow      string `yaml:"clear_workflow"`
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
	gitClient  git.ConfigOps
}

var (
	// Accept classic GitHub tokens like ghp_, gho_, ghu_, ghs_, ghr_
	githubTokenClassicRe = regexp.MustCompile(`^gh[opusr]_[A-Za-z0-9]{20,250}$`)
	// Accept fine-grained PATs starting with github_pat_
	githubTokenFineRe = regexp.MustCompile(`^github_pat_[A-Za-z0-9_-]{20,255}$`)
	// Accept GitLab tokens with optional glpat- prefix
	gitlabTokenRe = regexp.MustCompile(`^(glpat-)?[A-Za-z0-9_-]{20,100}$`)
	// Allow slashes; additional structural checks are applied separately
	gitRemoteNameCharsRe = regexp.MustCompile(`^[A-Za-z0-9._/\-]+$`)
)

// NewConfigManager creates a new configuration manager with the provided git client
func NewConfigManager(gitClient git.ConfigOps) *Manager {
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

func (c *Config) validateBranch() error {
	branch := c.Default.Branch
	if strings.TrimSpace(branch) == "" || strings.Contains(branch, " ") {
		return &ValidationError{"default.branch", branch, "must not contain spaces or be empty"}
	}
	return nil
}

func (c *Config) validateEditor() error {
	editor := strings.TrimSpace(c.Default.Editor)
	bin := parseEditorBinary(editor)
	if validEditorPath(bin) {
		return nil
	}
	if _, err := exec.LookPath(bin); err != nil {
		return &ValidationError{"default.editor", editor, "command not found in PATH or invalid path"}
	}
	return nil
}

func parseEditorBinary(editor string) string {
	if editor == "" {
		return ""
	}
	// Support basic quoted paths or first token before whitespace
	if (strings.HasPrefix(editor, "\"") && strings.Count(editor, "\"") >= 2) || (strings.HasPrefix(editor, "'") && strings.Count(editor, "'") >= 2) {
		q := editor[0:1]
		if idx := strings.Index(editor[1:], q); idx >= 0 {
			return editor[1 : 1+idx]
		}
	}
	if i := strings.IndexAny(editor, " \t"); i > 0 {
		return editor[:i]
	}
	return editor
}

func validEditorPath(bin string) bool {
	if !strings.ContainsAny(bin, "/\\") {
		return false
	}
	_, err := os.Stat(bin)
	return err == nil
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
	if err := c.validateGithubToken(); err != nil {
		return err
	}

	if err := c.validateGitlabToken(); err != nil {
		return err
	}

	return c.validateGithubRemote()
}

// validateGithubToken validates GitHub token format
func (c *Config) validateGithubToken() error {
	ghToken := c.Integration.Github.Token
	if ghToken != "" {
		if !githubTokenClassicRe.MatchString(ghToken) && !githubTokenFineRe.MatchString(ghToken) {
			return &ValidationError{
				Field:   "integration.github.token",
				Value:   "[REDACTED]",
				Message: "GitHub token must be a valid classic (gh[pousr]_) or fine-grained (github_pat_) token",
			}
		}
	}
	return nil
}

// validateGitlabToken validates GitLab token format
func (c *Config) validateGitlabToken() error {
	glToken := c.Integration.Gitlab.Token
	if glToken != "" {
		if !gitlabTokenRe.MatchString(glToken) {
			return &ValidationError{
				Field:   "integration.gitlab.token",
				Value:   "[REDACTED]",
				Message: "GitLab token must be 20-100 characters (alphanumeric, _ or -), with optional glpat- prefix",
			}
		}
	}
	return nil
}

// validateGithubRemote validates GitHub remote name format
func (c *Config) validateGithubRemote() error {
	remote := c.Integration.Github.DefaultRemote
	if remote == "" {
		return nil
	}

	if !gitRemoteNameCharsRe.MatchString(remote) || strings.Contains(remote, " ") {
		return &ValidationError{
			Field:   "integration.github.default-remote",
			Value:   remote,
			Message: "Remote may contain letters, digits, ., _, -, and / only",
		}
	}

	// Additional structural checks: no leading/trailing '.' or '/', and no empty/unsafe segments
	if strings.HasPrefix(remote, "/") || strings.HasSuffix(remote, "/") || strings.HasPrefix(remote, ".") || strings.HasSuffix(remote, ".") || strings.Contains(remote, "//") || strings.Contains(remote, "..") {
		return &ValidationError{
			Field:   "integration.github.default-remote",
			Value:   remote,
			Message: "Remote must not start/end with '.' or '/', nor contain '..' or '//'",
		}
	}

	return nil
}

func (c *Config) validateAliases() error {
	for name, value := range c.Aliases {
		if err := validateAliasName(name); err != nil {
			return err
		}
		if err := validateAliasValue(name, value); err != nil {
			return err
		}
	}
	return nil
}

func validateAliasName(name string) error {
	if strings.TrimSpace(name) == "" || strings.Contains(name, " ") {
		return &ValidationError{"aliases." + name, name, "alias names must not contain spaces"}
	}
	return nil
}

func validateAliasValue(name string, value interface{}) error {
	switch v := value.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			return &ValidationError{"aliases." + name, v, "alias command cannot be empty"}
		}
		return nil
	case []interface{}:
		return validateAliasSequence(name, v)
	default:
		return &ValidationError{Field: "aliases." + name, Value: value, Message: "alias must be either a string or array of strings"}
	}
}

func validateAliasSequence(name string, seq []interface{}) error {
	if len(seq) == 0 {
		return &ValidationError{"aliases." + name, seq, "alias sequence cannot be empty"}
	}
	for i, cmd := range seq {
		cmdStr, ok := cmd.(string)
		if !ok {
			return &ValidationError{Field: fmt.Sprintf("aliases.%s[%d]", name, i), Value: cmd, Message: "sequence commands must be strings"}
		}
		if strings.TrimSpace(cmdStr) == "" {
			return &ValidationError{Field: fmt.Sprintf("aliases.%s[%d]", name, i), Value: cmdStr, Message: "command in sequence cannot be empty"}
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
	if err := c.validateKeybindings(); err != nil {
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
	updaters := map[string]func(string){
		"core.editor":        func(v string) { cm.config.Default.Editor = v },
		"merge.tool":         func(v string) { cm.config.Default.MergeTool = v },
		"init.defaultBranch": func(v string) { cm.config.Default.Branch = v },
		"color.ui":           func(v string) { cm.config.UI.Color = v == "true" || v == "auto" },
		"core.pager":         func(v string) { cm.config.UI.Pager = v != "cat" },
		"fetch.auto":         func(v string) { cm.config.Behavior.AutoFetch = v == "true" },
	}
	if f, ok := updaters[command]; ok {
		f(value)
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

	if err := cm.syncDefaultSettings(config); err != nil {
		return err
	}

	if err := cm.syncUISettings(config); err != nil {
		return err
	}

	if err := cm.syncBehaviorSettings(config); err != nil {
		return err
	}

	return cm.syncAliases(config)
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

// syncAliases syncs alias settings to git config
func (cm *Manager) syncAliases(config *Config) error {
	for alias, value := range config.Aliases {
		if cmdStr, ok := value.(string); ok {
			if err := cm.gitClient.ConfigSetGlobal(fmt.Sprintf("alias.%s", alias), cmdStr); err != nil {
				return fmt.Errorf("failed to set git alias.%s: %w", alias, err)
			}
		}
	}

	return nil
}

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
	// On Windows, os.Rename requires the target not exist, so remove it first
	// On Unix, os.Rename is atomic and replaces the target
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

	parent, err := cm.navigateToParent(current, parts)
	if err != nil {
		return err
	}

	return cm.setFinalValue(parent, parts[len(parts)-1], value)
}

// navigateToParent navigates to the parent of the target field
func (cm *Manager) navigateToParent(current reflect.Value, parts []string) (reflect.Value, error) {
	for i, part := range parts[:len(parts)-1] {
		var err error
		current, err = cm.navigateOneLevel(current, part, parts[:i+1])
		if err != nil {
			return reflect.Value{}, err
		}
	}
	return current, nil
}

// navigateOneLevel navigates one level into a struct or map
func (cm *Manager) navigateOneLevel(current reflect.Value, part string, pathSoFar []string) (reflect.Value, error) {
	if current.Kind() == reflect.Ptr {
		current = current.Elem()
	}

	switch current.Kind() {
	case reflect.Struct:
		field, found := cm.findFieldByYamlTag(current.Type(), current, part)
		if !found {
			return reflect.Value{}, fmt.Errorf("field '%s' not found", strings.Join(pathSoFar, "."))
		}
		return field, nil

	case reflect.Map:
		mapValue := current.MapIndex(reflect.ValueOf(part))
		if !mapValue.IsValid() {
			return reflect.Value{}, fmt.Errorf("key '%s' not found", strings.Join(pathSoFar, "."))
		}
		return mapValue, nil

	default:
		return reflect.Value{}, fmt.Errorf("cannot navigate into %s", current.Kind())
	}
}

// setFinalValue sets the final value in the target location
func (cm *Manager) setFinalValue(current reflect.Value, lastPart string, value any) error {
	if current.Kind() == reflect.Ptr {
		current = current.Elem()
	}

	switch current.Kind() {
	case reflect.Struct:
		return cm.setStructField(current, lastPart, value)
	case reflect.Map:
		return cm.setMapValue(current, lastPart, value)
	default:
		return fmt.Errorf("cannot set value in %s", current.Kind())
	}
}

// setStructField sets a field value in a struct
func (cm *Manager) setStructField(current reflect.Value, fieldName string, value any) error {
	field, found := cm.findFieldByYamlTag(current.Type(), current, fieldName)
	if !found || !field.CanSet() {
		return fmt.Errorf("field '%s' not found or cannot be set", fieldName)
	}

	newValue := reflect.ValueOf(value)
	if !newValue.Type().ConvertibleTo(field.Type()) {
		return fmt.Errorf("cannot convert %s to %s", newValue.Type(), field.Type())
	}

	field.Set(newValue.Convert(field.Type()))
	return nil
}

// setMapValue sets a value in a map
func (cm *Manager) setMapValue(current reflect.Value, key string, value any) error {
	if current.Type().Key().Kind() != reflect.String {
		return fmt.Errorf("map key must be string")
	}

	newValue := reflect.ValueOf(value)
	if !newValue.Type().ConvertibleTo(current.Type().Elem()) {
		return fmt.Errorf("cannot convert %s to %s", newValue.Type(), current.Type().Elem())
	}

	current.SetMapIndex(reflect.ValueOf(key), newValue.Convert(current.Type().Elem()))
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
		cm.flattenStruct(value, prefix, result)
	case reflect.Map:
		cm.flattenMap(value, prefix, result)
	}
}

func (cm *Manager) flattenStruct(value reflect.Value, prefix string, result map[string]any) {
	structType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := structType.Field(i)

		fieldName := fieldType.Name
		if yamlTag := fieldType.Tag.Get("yaml"); yamlTag != "" {
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
}

func (cm *Manager) flattenMap(value reflect.Value, prefix string, result map[string]any) {
	for _, mapKey := range value.MapKeys() {
		mapValue := value.MapIndex(mapKey)
		key := mapKey.String()
		if prefix != "" {
			key = prefix + "." + mapKey.String()
		}
		result[key] = mapValue.Interface()
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

// validateKeybindings validates the keybinding configuration
func (c *Config) validateKeybindings() error {
	// Validate profile selection
	if err := c.validateProfile(); err != nil {
		return err
	}

	// Validate global keybindings
	bindings := map[string]string{
		"delete_word":          c.Interactive.Keybindings.DeleteWord,
		"clear_line":           c.Interactive.Keybindings.ClearLine,
		"delete_to_end":        c.Interactive.Keybindings.DeleteToEnd,
		"move_to_beginning":    c.Interactive.Keybindings.MoveToBeginning,
		"move_to_end":          c.Interactive.Keybindings.MoveToEnd,
		"move_up":              c.Interactive.Keybindings.MoveUp,
		"move_down":            c.Interactive.Keybindings.MoveDown,
		"add_to_workflow":      c.Interactive.Keybindings.AddToWorkflow,
		"toggle_workflow_view": c.Interactive.Keybindings.ToggleWorkflowView,
		"clear_workflow":       c.Interactive.Keybindings.ClearWorkflow,
		"soft_cancel":          c.Interactive.Keybindings.SoftCancel,
	}

	for action, keyStr := range bindings {
		// Empty bindings are allowed (will use defaults)
		if keyStr == "" {
			continue
		}
		if err := parseKeyBinding(keyStr); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("interactive.keybindings.%s", action),
				Value:   keyStr,
				Message: err.Error(),
			}
		}
	}

	// Validate context-specific keybindings
	if err := c.validateContextKeybindings(); err != nil {
		return err
	}

	// Validate platform-specific keybindings
	if err := c.validatePlatformKeybindings(); err != nil {
		return err
	}

	return nil
}

// validateProfile validates the profile selection
func (c *Config) validateProfile() error {
	profile := c.Interactive.Profile
	if profile == "" {
		return nil // Empty profile is allowed (defaults to "default")
	}

	validProfiles := map[string]bool{
		"default":  true,
		"emacs":    true,
		"vi":       true,
		"readline": true,
	}

	if !validProfiles[profile] {
		return &ValidationError{
			Field:   "interactive.profile",
			Value:   profile,
			Message: "must be one of: default, emacs, vi, readline",
		}
	}
	return nil
}

// validateContextKeybindings validates context-specific keybindings
func (c *Config) validateContextKeybindings() error {
	contexts := map[string]map[string]interface{}{
		"input":   c.Interactive.Contexts.Input.Keybindings,
		"results": c.Interactive.Contexts.Results.Keybindings,
		"search":  c.Interactive.Contexts.Search.Keybindings,
	}

	for contextName, bindings := range contexts {
		if bindings == nil {
			continue
		}
		for action, value := range bindings {
			if err := validateKeybindingValue(fmt.Sprintf("interactive.contexts.%s.keybindings.%s", contextName, action), value); err != nil {
				return err
			}
		}
	}
	return nil
}

// validatePlatformKeybindings validates platform and terminal specific keybindings
func (c *Config) validatePlatformKeybindings() error {
	platforms := map[string]map[string]interface{}{
		"darwin":  c.Interactive.Darwin.Keybindings,
		"linux":   c.Interactive.Linux.Keybindings,
		"windows": c.Interactive.Windows.Keybindings,
	}

	for platformName, bindings := range platforms {
		if bindings == nil {
			continue
		}
		for action, value := range bindings {
			if err := validateKeybindingValue(fmt.Sprintf("interactive.%s.keybindings.%s", platformName, action), value); err != nil {
				return err
			}
		}
	}

	// Validate terminal-specific keybindings
	if c.Interactive.Terminals != nil {
		for termName, termConfig := range c.Interactive.Terminals {
			if termConfig.Keybindings == nil {
				continue
			}
			for action, value := range termConfig.Keybindings {
				if err := validateKeybindingValue(fmt.Sprintf("interactive.terminals.%s.keybindings.%s", termName, action), value); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// validateKeybindingValue validates a keybinding value (string or array of strings)
func validateKeybindingValue(fieldPath string, value interface{}) error {
	switch v := value.(type) {
	case string:
		if v == "" {
			return nil // Empty is allowed
		}
		if err := parseKeyBinding(v); err != nil {
			return &ValidationError{
				Field:   fieldPath,
				Value:   v,
				Message: err.Error(),
			}
		}
	case []interface{}:
		for i, item := range v {
			itemStr, ok := item.(string)
			if !ok {
				return &ValidationError{
					Field:   fmt.Sprintf("%s[%d]", fieldPath, i),
					Value:   item,
					Message: "keybinding array items must be strings",
				}
			}
			if itemStr != "" {
				if err := parseKeyBinding(itemStr); err != nil {
					return &ValidationError{
						Field:   fmt.Sprintf("%s[%d]", fieldPath, i),
						Value:   itemStr,
						Message: err.Error(),
					}
				}
			}
		}
	default:
		return &ValidationError{
			Field:   fieldPath,
			Value:   value,
			Message: "keybinding must be a string or array of strings",
		}
	}
	return nil
}

// parseKeyBinding validates key binding strings.
// This simple validation is implemented here to avoid a circular import:
// importing the full keybinding parser from the 'cmd' (interactive UI) package
// would cause a circular dependency, since that package depends on 'config'.
func parseKeyBinding(keyStr string) error { //nolint:revive // parsing multiple legacy formats
	s := strings.TrimSpace(keyStr)
	if s == "" {
		return fmt.Errorf("empty key binding")
	}

	// Basic validation - check for supported formats
	sLower := strings.ToLower(s)

	// Accept ctrl+<key>, ^<key>, or c-<key> formats
	if (strings.HasPrefix(sLower, "ctrl+") && len(s) >= 6) ||
		(strings.HasPrefix(s, "^") && len(s) == 2) ||
		(strings.HasPrefix(sLower, "c-") && len(s) == 3) {
		return nil
	}

	return fmt.Errorf("unsupported key binding format: %s (supported: 'ctrl+<key>', '^<key>', 'c-<key>')", keyStr)
}
