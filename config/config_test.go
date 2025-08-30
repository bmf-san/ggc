package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestGetDefaultConfig tests the default configuration values
func TestGetDefaultConfig(t *testing.T) {
	config := getDefaultConfig()

	// test default values
	if config.Default.Branch != "main" {
		t.Errorf("Expected default branch to be 'main', got %s", config.Default.Branch)
	}
	if config.Default.Editor != "vim" {
		t.Errorf("Expected default editor to be 'vim', got %s", config.Default.Editor)
	}
	if config.Default.MergeTool != "vimdiff" {
		t.Errorf("Expected default merge tool to be 'vimdiff', got %s", config.Default.MergeTool)
	}

	// Test UI defaults
	if !config.UI.Color {
		t.Error("Expected UI color to be true")
	}
	if !config.UI.Pager {
		t.Error("Expected UI pager to be true")
	}

	// Test behavior defaults
	if config.Behavior.AutoPush {
		t.Error("Expected auto-push to be false")
	}
	if config.Behavior.ConfirmDestructive != "simple" {
		t.Errorf("Expected confirm-destructive to be 'simple', got %s", config.Behavior.ConfirmDestructive)
	}
	if !config.Behavior.AutoFetch {
		t.Error("Expected auto-fetch to be true")
	}
	if !config.Behavior.StashBeforeSwitch {
		t.Error("Expected stash-before-switch to be true")
	}

	if config.Integration.Github.DefaultRemote != "origin" {
		t.Errorf("Expected default remote to be 'origin', got %s", config.Integration.Github.DefaultRemote)
	}
}

// TestNewConfigManager tests the creation of a new config manager
func TestNewConfigManager(t *testing.T) {
	cm := NewConfigManager()

	if cm == nil {
		t.Fatal("Expected config manager to be created")
	}
	if cm.config == nil {
		t.Fatal("Expected config to be initialized")
	}
	if cm.configPath != "" {
		t.Errorf("Expected configPath to be empty initially, got %s", cm.configPath)
	}

	if cm.config.Default.Branch != "main" {
		t.Errorf("Expected default branch to be 'main', got %s", cm.config.Default.Branch)
	}
}

// TestGetConfigPaths tests the configuration path resolution
func TestGetConfigPaths(t *testing.T) {
	cm := NewConfigManager()
	paths := cm.getConfigPaths()

	if len(paths) != 2 {
		t.Errorf("Expected 2 config paths, got %d", len(paths))
	}

	homeDir, _ := os.UserHomeDir()
	expectedPaths := []string{
		filepath.Join(homeDir, ".ggcconfig.yaml"),
		filepath.Join(homeDir, ".config", "ggc", "config.yaml"),
	}

	for i, expected := range expectedPaths {
		if paths[i] != expected {
			t.Errorf("Expected path %d to be %s, got %s", i, expected, paths[i])
		}
	}
}

// TestLoadFromFile tests loading configuration from a file
func TestLoadFromFile(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.yaml")

	testConfig := `
default:
  branch: "develop"
  editor: "nano"
  merge-tool: "meld"
ui:
  color: false
  pager: false
behavior:
  auto-push: true
  confirm-destructive: "never"
  auto-fetch: false
  stash-before-switch: false
aliases:
  s: "status"
  c: "commit"
integration:
  github:
    token: "test-token"
    default-remote: "upstream"
  gitlab:
    token: "gitlab-token"
`

	err := os.WriteFile(configPath, []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cm := NewConfigManager()
	err = cm.loadFromFile(configPath)
	if err != nil {
		t.Fatalf("Failed to load config from file: %v", err)
	}

	if cm.config.Default.Branch != "develop" {
		t.Errorf("Expected branch to be 'develop', got %s", cm.config.Default.Branch)
	}
	if cm.config.Default.Editor != "nano" {
		t.Errorf("Expected editor to be 'nano', got %s", cm.config.Default.Editor)
	}
	if cm.config.UI.Color {
		t.Error("Expected color to be false")
	}
	if cm.config.Behavior.AutoPush != true {
		t.Error("Expected auto-push to be true")
	}
	if cm.config.Aliases["s"] != "status" {
		t.Errorf("Expected alias 's' to be 'status', got %s", cm.config.Aliases["s"])
	}
	if cm.config.Integration.Github.Token != "test-token" {
		t.Errorf("Expected github token to be 'test-token', got %s", cm.config.Integration.Github.Token)
	}
}

// TestLoad tests the Load method with no config file
func TestLoad(t *testing.T) {
	cm := NewConfigManager()

	originalHome := os.Getenv("HOME")
	tempDir := t.TempDir()
	if err := os.Setenv("HOME", tempDir); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}
	defer func() {
		if err := os.Setenv("HOME", originalHome); err != nil {
			t.Fatalf("failed to restore HOME: %v", err)
		}
	}()

	err := cm.Load()
	if err != nil {
		t.Fatalf("Load should not fail when no config file exists: %v", err)
	}

	expectedPath := filepath.Join(tempDir, ".ggcconfig.yaml")
	if cm.configPath != expectedPath {
		t.Errorf("Expected config path to be %s, got %s", expectedPath, cm.configPath)
	}
}

// TestSave tests saving configuration to file
func TestSave(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-save.yaml")

	cm := NewConfigManager()
	cm.configPath = configPath

	cm.config.Default.Branch = "development"
	cm.config.UI.Color = false
	cm.config.Aliases["test"] = "help"

	err := cm.Save()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read saved config: %v", err)
	}

	var loadedConfig Config
	err = yaml.Unmarshal(data, &loadedConfig)
	if err != nil {
		t.Fatalf("Failed to unmarshal saved config: %v", err)
	}

	if loadedConfig.Default.Branch != "development" {
		t.Errorf("Expected saved branch to be 'development', got %s", loadedConfig.Default.Branch)
	}
	if loadedConfig.UI.Color {
		t.Error("Expected saved color to be false")
	}
	if loadedConfig.Aliases["test"] != "help" {
		t.Errorf("Expected saved alias to be 'help', got %s", loadedConfig.Aliases["test"])
	}
}

// TestGetValueByPath tests getting values using dot notation
func TestGetValueByPath(t *testing.T) {
	cm := NewConfigManager()

	testCases := []struct {
		path     string
		expected any
	}{
		{"default.branch", "main"},
		{"default.editor", "vim"},
		{"ui.color", true},
		{"behavior.auto-push", false},
		{"integration.github.default-remote", "origin"},
	}

	for _, tc := range testCases {
		value, err := cm.getValueByPath(cm.config, tc.path)
		if err != nil {
			t.Errorf("Failed to get value for path %s: %v", tc.path, err)
			continue
		}

		if !reflect.DeepEqual(value, tc.expected) {
			t.Errorf("Expected value for path %s to be %v, got %v", tc.path, tc.expected, value)
		}
	}
}

// TestGetValueByPathErrors tests error cases for getValueByPath
func TestGetValueByPathErrors(t *testing.T) {
	cm := NewConfigManager()

	testCases := []string{
		"nonexistent.field",
		"default.nonexistent",
		"aliases.nonexistent",
		"ui.color.invalid", // trying to navigate into a bool
	}

	for _, path := range testCases {
		_, err := cm.getValueByPath(cm.config, path)
		if err == nil {
			t.Errorf("Expected error for invalid path %s", path)
		}
	}
}

// TestSetValueByPath tests setting values using dot notation
func TestSetValueByPath(t *testing.T) {
	cm := NewConfigManager()

	testCases := []struct {
		path     string
		value    any
		expected any
	}{
		{"default.branch", "develop", "develop"},
		{"default.editor", "emacs", "emacs"},
		{"ui.color", false, false},
		{"behavior.auto-push", true, true},
		{"aliases.new", "new-command", "new-command"},
	}

	for _, tc := range testCases {
		err := cm.setValueByPath(cm.config, tc.path, tc.value)
		if err != nil {
			t.Errorf("Failed to set value for path %s: %v", tc.path, err)
			continue
		}

		actualValue, err := cm.getValueByPath(cm.config, tc.path)
		if err != nil {
			t.Errorf("Failed to get value after setting for path %s: %v", tc.path, err)
			continue
		}

		if !reflect.DeepEqual(actualValue, tc.expected) {
			t.Errorf("Expected value for path %s to be %v, got %v", tc.path, tc.expected, actualValue)
		}
	}
}

// TestGet tests the Get method
func TestGet(t *testing.T) {
	cm := NewConfigManager()

	value, err := cm.Get("default.branch")
	if err != nil {
		t.Fatalf("Failed to get value: %v", err)
	}

	if value != "main" {
		t.Errorf("Expected 'main', got %s", value)
	}
}

// TestSet tests the Set method
func TestSet(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-set.yaml")

	cm := NewConfigManager()
	cm.configPath = configPath

	err := cm.Set("default.branch", "develop")
	if err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}

	value, err := cm.Get("default.branch")
	if err != nil {
		t.Fatalf("Failed to get value after setting: %v", err)
	}

	if value != "develop" {
		t.Errorf("Expected 'develop', got %s", value)
	}

	// Test error case: invalid path
	err = cm.Set("invalid.nonexistent.path", "value")
	if err == nil {
		t.Error("Expected error when setting invalid path, but got nil")
	}

	// Test error case: invalid value type
	err = cm.Set("ui.color", "not_a_boolean")
	if err == nil {
		t.Error("Expected error when setting invalid boolean value, but got nil")
	}
}

// TestList tests the List method
func TestList(t *testing.T) {
	cm := NewConfigManager()
	list := cm.List()

	expectedKeys := []string{
		"default.branch",
		"default.editor",
		"ui.color",
		"behavior.auto-push",
		"integration.github.default-remote",
	}

	for _, key := range expectedKeys {
		if _, exists := list[key]; !exists {
			t.Errorf("Expected key %s to exist in list: %s", key, list)
		}
	}

	if list["default.branch"] != "main" {
		t.Errorf("Expected default.branch to be 'main', got %v", list["default.branch"])
	}
	if list["ui.color"] != true {
		t.Errorf("Expected ui.color to be true, got %v", list["ui.color"])
	}

	if !strings.Contains(stringifyAnyMap(list), "aliases") {
		t.Errorf("Expected list to contain aliases, got %v", list)
	}
}

// TestFindFieldByYamlTag tests the findFieldByYamlTag method
func TestFindFieldByYamlTag(t *testing.T) {
	cm := NewConfigManager()
	configType := reflect.TypeOf(*cm.config)
	configValue := reflect.ValueOf(*cm.config)

	field, found := cm.findFieldByYamlTag(configType, configValue, "default")
	if !found {
		t.Error("Expected to find 'default' field by YAML tag")
	}
	if field.Type().Name() != "" { // Anonymous struct
		defaultValue := field.Interface()
		defaultType := reflect.TypeOf(defaultValue)
		branchField, branchFound := cm.findFieldByYamlTag(defaultType, reflect.ValueOf(defaultValue), "branch")
		if !branchFound {
			t.Error("Expected to find 'branch' field in default struct")
		}
		if branchField.String() != "main" {
			t.Errorf("Expected branch to be 'main', got %s", branchField.String())
		}
	}

	_, found = cm.findFieldByYamlTag(configType, configValue, "Default")
	if !found {
		t.Error("Expected to find 'Default' field by name")
	}

	_, found = cm.findFieldByYamlTag(configType, configValue, "nonexistent")
	if found {
		t.Error("Expected not to find nonexistent field")
	}
}

// TestFlattenConfig tests the flattenConfig method
func TestFlattenConfig(t *testing.T) {
	cm := NewConfigManager()
	result := make(map[string]any)

	cm.flattenConfig(cm.config, "", result)

	expectedKeys := []string{
		"default.branch",
		"default.editor",
		"ui.color",
		"behavior.auto-push",
	}

	for _, key := range expectedKeys {
		if _, exists := result[key]; !exists {
			t.Errorf("Expected key %s to exist in flattened config", key)
		}
	}

	if !strings.Contains(stringifyAnyMap(result), "aliases") {
		t.Errorf("Expected list to contain aliases, got %v", result)
	}

	result2 := make(map[string]any)
	cm.flattenConfig(cm.config.Default, "test", result2)

	if _, exists := result2["test.branch"]; !exists {
		t.Error("Expected key 'test.branch' to exist with prefix")
	}
}

// TestLoadConfig tests the LoadConfig method
func TestLoadConfig(t *testing.T) {
	tempDir := t.TempDir()

	originalHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", tempDir); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}
	defer func() {
		if err := os.Setenv("HOME", originalHome); err != nil {
			t.Fatalf("failed to restore HOME: %v", err)
		}
	}()

	cm := NewConfigManager()
	cm.LoadConfig()

	if cm.config == nil {
		t.Error("Expected config to be loaded")
	}
}

// TestGetConfig tests the GetConfig method
func TestGetConfig(t *testing.T) {
	cm := NewConfigManager()
	config := cm.GetConfig()

	if config == nil {
		t.Fatal("Expected config to be returned")
	}

	if config != cm.config {
		t.Error("Expected GetConfig to return the same config instance")
	}
}

// TestConfigStructTags tests that all struct fields have proper YAML tags
func TestConfigStructTags(t *testing.T) {
	config := &Config{}
	configType := reflect.TypeOf(*config)

	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		yamlTag := field.Tag.Get("yaml")
		if yamlTag == "" {
			t.Errorf("Field %s should have a yaml tag", field.Name)
		}
	}
}

// TestInvalidYAMLHandling tests handling of invalid YAML
func TestInvalidYAMLHandling(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "invalid-config.yaml")

	invalidYAML := `
default:
  branch: "main"
  editor: "vim"
ui:
  color: invalid_boolean
  pager: [this, is, invalid]
`

	err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid config: %v", err)
	}

	cm := NewConfigManager()
	err = cm.loadFromFile(configPath)
	if err == nil {
		t.Error("Expected error when loading invalid YAML")
	}
}

// TestTypeConversion tests type conversion in setValueByPath
func TestTypeConversion(t *testing.T) {
	cm := NewConfigManager()

	// Test setting string value to string field
	err := cm.setValueByPath(cm.config, "default.branch", "test")
	if err != nil {
		t.Errorf("Failed to set string value: %v", err)
	}

	// Test setting bool value to bool field
	err = cm.setValueByPath(cm.config, "ui.color", false)
	if err != nil {
		t.Errorf("Failed to set bool value: %v", err)
	}

	// Test type conversion error
	err = cm.setValueByPath(cm.config, "ui.color", "invalid_bool")
	if err == nil {
		t.Error("Expected error when setting invalid type")
	}
}

func TestConfig_Validate(t *testing.T) {
	t.Run("Valid config", func(t *testing.T) {
		cfg := &Config{}
		cfg.Default.Branch = "main"
		cfg.Default.Editor = "vim"
		cfg.Default.MergeTool = "meld"
		cfg.UI.Color = true
		cfg.UI.Pager = false
		cfg.Behavior.AutoPush = true
		cfg.Behavior.ConfirmDestructive = "simple"
		cfg.Behavior.AutoFetch = true
		cfg.Behavior.StashBeforeSwitch = true
		cfg.Aliases = map[string]any{"st": "status"}
		cfg.Integration.Github.Token = "ghp_1234567890asdflkasfdasf"
		cfg.Integration.Github.DefaultRemote = "origin"
		cfg.Integration.Gitlab.Token = "glpat-abc123asdlfkasjdflasfdasdf"

		err := cfg.Validate()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("Invalid confirm-destructive", func(t *testing.T) {
		cfg := &Config{}
		cfg.Behavior.ConfirmDestructive = "maybe"
		cfg.Default.Branch = "main"

		err := cfg.Validate()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "invalid value") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Invalid alias name", func(t *testing.T) {
		cfg := &Config{}
		cfg.Behavior.ConfirmDestructive = "never"
		cfg.Default.Branch = "main"
		cfg.Default.Editor = "vim"
		cfg.Default.Branch = "main"
		cfg.Aliases = map[string]any{"invalid alias": "status"}

		err := cfg.Validate()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "invalid value") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Invalid GitHub token", func(t *testing.T) {
		cfg := &Config{}
		cfg.Behavior.ConfirmDestructive = "always"
		cfg.Default.Editor = "vim"
		cfg.Integration.Github.Token = "bad-token"
		cfg.Default.Branch = "main"

		err := cfg.Validate()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "invalid value") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Invalid GitLab token", func(t *testing.T) {
		cfg := &Config{}
		cfg.Behavior.ConfirmDestructive = "always"
		cfg.Default.Editor = "vim"
		cfg.Integration.Gitlab.Token = "bad-token"
		cfg.Default.Branch = "main"

		err := cfg.Validate()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "invalid value") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestConfig_ParseAlias(t *testing.T) {
	config := &Config{
		Aliases: map[string]interface{}{
			"st":   "status",
			"acp":  []interface{}{"add", "commit", "push"},
			"sync": []interface{}{"pull", "add", "commit", "push"},
		},
	}

	tests := []struct {
		name         string
		aliasName    string
		wantType     AliasType
		wantCommands []string
		wantError    bool
	}{
		{
			name:         "simple alias",
			aliasName:    "st",
			wantType:     SimpleAlias,
			wantCommands: []string{"status"},
			wantError:    false,
		},
		{
			name:         "sequence alias - short",
			aliasName:    "acp",
			wantType:     SequenceAlias,
			wantCommands: []string{"add", "commit", "push"},
			wantError:    false,
		},
		{
			name:         "sequence alias - long",
			aliasName:    "sync",
			wantType:     SequenceAlias,
			wantCommands: []string{"pull", "add", "commit", "push"},
			wantError:    false,
		},
		{
			name:         "non-existent alias",
			aliasName:    "nonexistent",
			wantType:     SimpleAlias,
			wantCommands: nil,
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := config.ParseAlias(tt.aliasName)

			if tt.wantError {
				if err == nil {
					t.Errorf("ParseAlias() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseAlias() unexpected error: %v", err)
				return
			}

			if parsed.Type != tt.wantType {
				t.Errorf("ParseAlias() type = %v, want %v", parsed.Type, tt.wantType)
			}

			if len(parsed.Commands) != len(tt.wantCommands) {
				t.Errorf("ParseAlias() commands length = %v, want %v", len(parsed.Commands), len(tt.wantCommands))
				return
			}

			for i, cmd := range parsed.Commands {
				if cmd != tt.wantCommands[i] {
					t.Errorf("ParseAlias() commands[%d] = %v, want %v", i, cmd, tt.wantCommands[i])
				}
			}
		})
	}
}

func TestConfig_IsAlias(t *testing.T) {
	config := &Config{
		Aliases: map[string]interface{}{
			"st":  "status",
			"acp": []interface{}{"add", "commit", "push"},
		},
	}

	tests := []struct {
		name      string
		aliasName string
		want      bool
	}{
		{"existing simple alias", "st", true},
		{"existing sequence alias", "acp", true},
		{"non-existing alias", "nonexistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := config.IsAlias(tt.aliasName); got != tt.want {
				t.Errorf("IsAlias() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetAliasCommands(t *testing.T) {
	config := &Config{
		Aliases: map[string]interface{}{
			"st":    "status",
			"acp":   []interface{}{"add", "commit", "push"},
			"quick": []interface{}{"status", "add", "commit"},
		},
	}

	tests := []struct {
		name         string
		aliasName    string
		wantCommands []string
		wantError    bool
	}{
		{
			name:         "simple alias",
			aliasName:    "st",
			wantCommands: []string{"status"},
			wantError:    false,
		},
		{
			name:         "sequence alias",
			aliasName:    "acp",
			wantCommands: []string{"add", "commit", "push"},
			wantError:    false,
		},
		{
			name:         "another sequence alias",
			aliasName:    "quick",
			wantCommands: []string{"status", "add", "commit"},
			wantError:    false,
		},
		{
			name:         "non-existent alias",
			aliasName:    "nonexistent",
			wantCommands: nil,
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commands, err := config.GetAliasCommands(tt.aliasName)

			if tt.wantError {
				if err == nil {
					t.Errorf("GetAliasCommands() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("GetAliasCommands() unexpected error: %v", err)
				return
			}

			if len(commands) != len(tt.wantCommands) {
				t.Errorf("GetAliasCommands() commands length = %v, want %v", len(commands), len(tt.wantCommands))
				return
			}

			for i, cmd := range commands {
				if cmd != tt.wantCommands[i] {
					t.Errorf("GetAliasCommands() commands[%d] = %v, want %v", i, cmd, tt.wantCommands[i])
				}
			}
		})
	}
}

func TestConfig_validateAliases(t *testing.T) {
	tests := []struct {
		name      string
		aliases   map[string]interface{}
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid simple aliases",
			aliases: map[string]interface{}{
				"st": "status",
				"br": "branch",
			},
			wantError: false,
		},
		{
			name: "valid sequence aliases",
			aliases: map[string]interface{}{
				"acp":  []interface{}{"add", "commit", "push"},
				"sync": []interface{}{"pull", "status"},
			},
			wantError: false,
		},
		{
			name: "mixed valid aliases",
			aliases: map[string]interface{}{
				"st":  "status",
				"acp": []interface{}{"add", "commit", "push"},
			},
			wantError: false,
		},
		{
			name: "invalid alias name with space",
			aliases: map[string]interface{}{
				"bad name": "status",
			},
			wantError: true,
			errorMsg:  "alias names must not contain spaces",
		},
		{
			name: "empty alias name",
			aliases: map[string]interface{}{
				"": "status",
			},
			wantError: true,
			errorMsg:  "alias names must not contain spaces",
		},
		{
			name: "empty sequence alias",
			aliases: map[string]interface{}{
				"empty": []interface{}{},
			},
			wantError: true,
			errorMsg:  "alias sequence cannot be empty",
		},
		{
			name: "invalid alias type",
			aliases: map[string]interface{}{
				"bad": 123,
			},
			wantError: true,
			errorMsg:  "alias must be either a string or array of strings",
		},
		{
			name: "non-string command in sequence",
			aliases: map[string]interface{}{
				"bad": []interface{}{"add", 123, "push"},
			},
			wantError: true,
			errorMsg:  "sequence commands must be strings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				Aliases: tt.aliases,
			}

			err := config.validateAliases()

			if tt.wantError {
				if err == nil {
					t.Errorf("validateAliases() expected error, got nil")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("validateAliases() error = %v, want to contain %v", err.Error(), tt.errorMsg)
				}
			} else if err != nil {
				t.Errorf("validateAliases() unexpected error: %v", err)
			}
		})
	}
}

func TestConfig_GetAllAliases(t *testing.T) {
	config := &Config{
		Aliases: map[string]interface{}{
			"st":    "status",
			"acp":   []interface{}{"add", "commit", "push"},
			"quick": []interface{}{"status", "add"},
		},
	}

	aliases := config.GetAllAliases()

	expectedCount := 3
	if len(aliases) != expectedCount {
		t.Errorf("GetAllAliases() returned %d aliases, want %d", len(aliases), expectedCount)
	}

	// Check simple alias
	if parsed, ok := aliases["st"]; !ok {
		t.Errorf("GetAllAliases() missing 'st' alias")
	} else if parsed.Type != SimpleAlias {
		t.Errorf("GetAllAliases() 'st' alias type = %v, want %v", parsed.Type, SimpleAlias)
	}

	if parsed, ok := aliases["acp"]; !ok {
		t.Errorf("GetAllAliases() missing 'acp' alias")
	} else if parsed.Type != SequenceAlias {
		t.Errorf("GetAllAliases() 'acp' alias type = %v, want %v", parsed.Type, SequenceAlias)
	} else if len(parsed.Commands) != 3 {
		t.Errorf("GetAllAliases() 'acp' alias commands length = %v, want 3", len(parsed.Commands))
	}
}

func stringifyAnyMap(m map[string]any) string {
	var b strings.Builder
	b.WriteString("{")
	first := true
	for key, val := range m {
		if !first {
			b.WriteString(" ")
		}
		first = false
		b.WriteString(fmt.Sprintf("%s: %s", key, stringifyValue(val)))
	}
	b.WriteString("}")
	return b.String()
}

func stringifyValue(val any) string {
	switch v := val.(type) {
	case string:
		return v
	case []any:
		strs := make([]string, len(v))
		for i, item := range v {
			strs[i] = stringifyValue(item)
		}
		return "[" + strings.Join(strs, " ") + "]"
	case []string:
		return "[" + strings.Join(v, " ") + "]"
	case map[string]any:
		return stringifyAnyMap(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// TestManagerLoadConfig tests the LoadConfig method error paths
func TestManagerLoadConfig(t *testing.T) {
	t.Run("LoadConfig with invalid path executes without panic", func(t *testing.T) {
		cm := NewConfigManager()
		cm.configPath = "/nonexistent/directory/config.yaml"

		// This should not panic, even if file operations fail
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("LoadConfig should not panic, but got: %v", r)
			}
		}()

		cm.LoadConfig()

		// Verify that config is still accessible (default config should be loaded)
		config := cm.GetConfig()
		if config == nil {
			t.Error("Config should not be nil after LoadConfig, even with invalid path")
		}
	})

	t.Run("LoadConfig handles missing directory gracefully", func(t *testing.T) {
		cm := NewConfigManager()
		cm.configPath = "/definitely/nonexistent/path/config.yaml"

		// Should not panic or crash
		cm.LoadConfig()

		// Should still have a valid config object
		if cm.GetConfig() == nil {
			t.Error("Should have valid config even with missing directory")
		}
	})
}
