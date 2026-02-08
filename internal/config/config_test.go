package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"go.yaml.in/yaml/v3"

	"github.com/bmf-san/ggc/v7/internal/testutil"
)

// MockFileOps implements FileOps for testing
type MockFileOps struct {
	// files maps file paths to their content.
	files map[string][]byte
	// dirs tracks directory existence.
	dirs map[string]bool
}

func NewMockFileOps() *MockFileOps {
	return &MockFileOps{
		files: make(map[string][]byte),
		dirs:  map[string]bool{"/": true, ".": true},
	}
}

func (m *MockFileOps) ReadFile(filename string) ([]byte, error) {
	if data, ok := m.files[filename]; ok {
		return data, nil
	}
	return nil, &os.PathError{Op: "open", Path: filename, Err: os.ErrNotExist}
}

func (m *MockFileOps) WriteFile(filename string, data []byte, perm os.FileMode) error {
	m.files[filename] = data
	return nil
}

func (m *MockFileOps) Stat(name string) (os.FileInfo, error) {
	if _, ok := m.files[name]; ok {
		return &mockFileInfo{name: name, size: int64(len(m.files[name]))}, nil
	}
	if m.dirs[name] {
		return &mockFileInfo{name: name, isDir: true}, nil
	}
	return nil, &os.PathError{Op: "stat", Path: name, Err: os.ErrNotExist}
}

func (m *MockFileOps) MkdirAll(path string, perm os.FileMode) error {
	m.dirs[path] = true
	return nil
}

func (m *MockFileOps) CreateTemp(dir, pattern string) (TempFile, error) {
	if !m.dirs[dir] && dir != "." && dir != "/" {
		return nil, &os.PathError{Op: "createtemp", Path: dir, Err: os.ErrNotExist}
	}
	name := dir + "/temp_" + pattern
	return &mockTempFile{name: name, fs: m}, nil
}

func (m *MockFileOps) Remove(name string) error {
	delete(m.files, name)
	return nil
}

func (m *MockFileOps) Rename(oldpath, newpath string) error {
	if data, ok := m.files[oldpath]; ok {
		m.files[newpath] = data
		delete(m.files, oldpath)
	}
	return nil
}

func (m *MockFileOps) Chmod(name string, mode os.FileMode) error {
	return nil // No-op for testing
}

type mockFileInfo struct {
	name  string
	size  int64
	isDir bool
}

func (m *mockFileInfo) Name() string       { return m.name }
func (m *mockFileInfo) Size() int64        { return m.size }
func (m *mockFileInfo) Mode() os.FileMode  { return 0644 }
func (m *mockFileInfo) ModTime() time.Time { return time.Now() }
func (m *mockFileInfo) IsDir() bool        { return m.isDir }
func (m *mockFileInfo) Sys() interface{}   { return nil }

type mockTempFile struct {
	name string
	data []byte
	fs   *MockFileOps
}

func (m *mockTempFile) Write(p []byte) (n int, err error) {
	m.data = append(m.data, p...)
	return len(p), nil
}

func (m *mockTempFile) Close() error {
	m.fs.files[m.name] = m.data
	return nil
}

func (m *mockTempFile) Name() string { return m.name }

// newTestConfigManager creates a config manager for testing without executing git commands
func newTestConfigManager() *Manager {
	mockClient := testutil.NewMockGitClient()
	return NewConfigManager(mockClient)
}

// TestGetDefaultConfig tests the default configuration values
func TestGetDefaultConfig(t *testing.T) {
	// Use mock client to avoid executing real git commands
	mockClient := testutil.NewMockGitClient()
	config := getDefaultConfig(mockClient)

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

	if config.Git.DefaultRemote != "origin" {
		t.Errorf("Expected default remote to be 'origin', got %s", config.Git.DefaultRemote)
	}
}

// TestNewConfigManager tests the creation of a new config manager
func TestNewConfigManager(t *testing.T) {
	mockClient := testutil.NewMockGitClient()
	cm := NewConfigManager(mockClient)

	t.Run("manager_creation", func(t *testing.T) {
		if cm == nil {
			t.Fatal("Expected config manager to be created")
		}
	})

	t.Run("config_initialization", func(t *testing.T) {
		if cm == nil {
			t.Skip("Skipping due to nil manager")
		}
		if cm.config == nil {
			t.Fatal("Expected config to be initialized")
		}
	})

	t.Run("config_path_initialization", func(t *testing.T) {
		if cm == nil {
			t.Skip("Skipping due to nil manager")
		}
		if cm.configPath != "" {
			t.Fatalf("Expected configPath to be empty initially, got: %s", cm.configPath)
		}
	})

	t.Run("default_branch", func(t *testing.T) {
		if cm == nil || cm.config == nil {
			t.Skip("Skipping due to nil manager or config")
		}
		if cm.config.Default.Branch != "main" {
			t.Errorf("Expected default branch to be 'main', got %s", cm.config.Default.Branch)
		}
	})
}

// TestGetConfigPaths tests the configuration path resolution
func TestGetConfigPaths(t *testing.T) {
	cm := newTestConfigManager()
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
	configPath := filepath.Join(tempDir, "config.yaml")

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
git:
  default-remote: "upstream"
`

	err := os.WriteFile(configPath, []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cm := newTestConfigManager()
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
	if cm.config.Git.DefaultRemote != "upstream" {
		t.Errorf("Expected git default remote to be 'upstream', got %s", cm.config.Git.DefaultRemote)
	}
}

// TestLoad tests the Load method with no config file
func TestLoad(t *testing.T) {
	// Note: This test uses the real filesystem because Load() method uses getConfigPaths()
	// which depends on os.UserHomeDir(). This is acceptable as it only reads (doesn't write)
	// and doesn't modify the system state.
	cm := newTestConfigManager()

	err := cm.Load()
	// Load should succeed even when no config file exists (uses default config)
	if err != nil {
		// If error occurs, it should be because config file doesn't exist, which is expected
		t.Logf("Load returned expected error (no config file): %v", err)
	}

	// The config path should be set to the first path from getConfigPaths()
	if cm.configPath == "" {
		t.Error("Expected config path to be set after Load()")
	}

	// Verify that default config is loaded
	if cm.config == nil {
		t.Error("Expected config to be loaded with defaults")
	}
}

// TestSave tests saving configuration to file
func TestSave(t *testing.T) {
	mockFS := NewMockFileOps()
	configPath := "/test/config.yaml"

	// Create directory in mock filesystem
	err := mockFS.MkdirAll("/test", 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	cm := newTestConfigManager()
	cm.configPath = configPath

	cm.config.Default.Branch = "development"
	cm.config.UI.Color = false
	cm.config.Aliases["test"] = "help"

	err = cm.SaveWithFileOps(mockFS)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Check that file was created in mock filesystem
	_, err = mockFS.Stat(configPath)
	if err != nil {
		t.Fatalf("Config file was not created: %v", err)
	}

	data, err := mockFS.ReadFile(configPath)
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

// TestSaveDoesNotWriteOnInvalidConfig ensures Save validates before writing
// and does not leave a config file on disk when validation fails.
func TestSaveDoesNotWriteOnInvalidConfig(t *testing.T) {
	mockFS := NewMockFileOps()
	configPath := "/test/invalid-config.yaml"

	// Create directory in mock filesystem
	err := mockFS.MkdirAll("/test", 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	cm := newTestConfigManager()
	cm.configPath = configPath

	// Force an invalid editor so validation fails
	cm.config.Default.Editor = "this-editor-should-not-exist-xyz"

	err = cm.SaveWithFileOps(mockFS)
	if err == nil {
		t.Fatal("expected Save to fail validation, got nil error")
	}

	// Check that no config file was written to mock filesystem
	if _, statErr := mockFS.Stat(configPath); statErr == nil {
		t.Fatal("expected no config file to be written, but file exists")
	}
}

// TestGetValueByPath tests getting values using dot notation
func TestGetValueByPath(t *testing.T) {
	cm := newTestConfigManager()

	testCases := []struct {
		path     string
		expected any
	}{
		{"default.branch", "main"},
		{"default.editor", "vim"},
		{"ui.color", true},
		{"behavior.auto-push", false},
		{"git.default-remote", "origin"},
	}

	for _, tc := range testCases {
		value, err := cm.getValueByPath(cm.config, tc.path)
		if err != nil {
			t.Errorf("Failed to get value for path %s: %v", tc.path, err)
			continue
		}

		if value != tc.expected {
			t.Errorf("Expected value for path %s to be %v, got %v", tc.path, tc.expected, value)
		}
	}
}

// TestGetValueByPathErrors tests error cases for getValueByPath
func TestGetValueByPathErrors(t *testing.T) {
	cm := NewConfigManager(testutil.NewMockGitClient())

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
	cm := newTestConfigManager()

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

		if actualValue != tc.expected {
			t.Errorf("Expected value for path %s to be %v, got %v", tc.path, tc.expected, actualValue)
		}
	}
}

// TestGet tests the Get method
func TestGet(t *testing.T) {
	cm := newTestConfigManager()

	value, err := cm.Get("default.branch")
	if err != nil {
		t.Fatalf("Failed to get value: %v", err)
	}

	if value != "main" {
		t.Errorf("Expected 'main', got %s", value)
	}

	if _, err := cm.Get("invalid..path"); err == nil {
		t.Error("Expected error for config path with empty segment, got nil")
	} else if !strings.Contains(err.Error(), "segment") {
		t.Errorf("unexpected error for invalid path: %v", err)
	}
}

// TestSet tests the Set method
func TestSet(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	cm := newTestConfigManager()
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

	// Test validation: empty path segment
	err = cm.Set("interactive..profile", "default")
	if err == nil {
		t.Error("Expected error when setting path with empty segment, but got nil")
	} else if !strings.Contains(err.Error(), "segment") {
		t.Errorf("unexpected error for double dot path: %v", err)
	}

	// Test error case: invalid value type
	err = cm.Set("ui.color", "not_a_boolean")
	if err == nil {
		t.Error("Expected error when setting invalid boolean value, but got nil")
	}
}

// TestList tests the List method
func TestList(t *testing.T) {
	cm := newTestConfigManager()
	list := cm.List()

	expectedKeys := []string{
		"default.branch",
		"default.editor",
		"ui.color",
		"behavior.auto-push",
		"git.default-remote",
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
	cm := newTestConfigManager()
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
	cm := newTestConfigManager()
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

	cm := NewConfigManager(testutil.NewMockGitClient())
	// Use mock git client to avoid real git config operations
	cm.gitClient = testutil.NewMockGitClient()
	_ = cm.LoadConfig()

	if cm.config == nil {
		t.Error("Expected config to be loaded")
	}
}

// TestGetConfig tests the GetConfig method
func TestGetConfig(t *testing.T) {
	cm := newTestConfigManager()
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

	cm := NewConfigManager(testutil.NewMockGitClient())
	// Use mock git client to avoid real git config operations
	cm.gitClient = testutil.NewMockGitClient()
	err = cm.loadFromFile(configPath)
	if err == nil {
		t.Error("Expected error when loading invalid YAML")
	}
}

func TestLoadConfigDoesNotOverwriteMalformedFile(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".ggcconfig.yaml")
	invalidYAML := "broken: [yaml\n"
	if err := os.WriteFile(configPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("failed to write malformed config: %v", err)
	}

	originalHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", tempDir); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}
	defer func() {
		if err := os.Setenv("HOME", originalHome); err != nil {
			t.Fatalf("failed to restore HOME: %v", err)
		}
	}()

	cm := NewConfigManager(testutil.NewMockGitClient())
	cm.gitClient = testutil.NewMockGitClient()
	cm.LoadConfig()

	got, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config after LoadConfig: %v", err)
	}
	if string(got) != invalidYAML {
		t.Fatalf("malformed config was overwritten.\nwant: %q\ngot:  %q", invalidYAML, string(got))
	}
}

// TestTypeConversion tests type conversion in setValueByPath
func TestTypeConversion(t *testing.T) {
	cm := newTestConfigManager()

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
		cfg.Git.DefaultRemote = "origin"

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

	t.Run("Invalid interactive keybinding", func(t *testing.T) {
		cfg := &Config{}
		cfg.Default.Branch = "main"
		cfg.Default.Editor = "cat"
		cfg.Behavior.ConfirmDestructive = "never"
		cfg.Interactive.Keybindings.DeleteWord = "Shift+A"

		err := cfg.Validate()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "unsupported key binding format") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("Context keybindings missing in partial configuration", func(t *testing.T) {
		cfg := &Config{}
		cfg.Default.Branch = "main"
		cfg.Default.Editor = "vim"
		cfg.Behavior.ConfirmDestructive = "never"
		cfg.Interactive.Contexts.Input.Keybindings = map[string]interface{}{
			"move_up": "Ctrl+P",
		}

		err := cfg.Validate()
		if err == nil {
			t.Fatal("expected validation error for missing context keybindings")
		}
		if !strings.Contains(err.Error(), "interactive.contexts.") || !strings.Contains(err.Error(), "keybindings map is missing") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("Invalid interactive profile", func(t *testing.T) {
		cfg := &Config{}
		cfg.Default.Branch = "main"
		cfg.Default.Editor = "cat"
		cfg.Behavior.ConfirmDestructive = "never"
		cfg.Interactive.Profile = "custom"

		err := cfg.Validate()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "interactive.profile") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestConfig_ParseAlias(t *testing.T) {
	config := &Config{
		Aliases: map[string]interface{}{
			"st":         "status",
			"acp":        []interface{}{"add", "commit", "push"},
			"sync":       []interface{}{"pull", "add", "commit", "push"},
			"deploy":     []interface{}{"branch checkout {0}", "push {0}"},
			"deploy-msg": "commit -m '{0}'",
		},
	}

	tests := []struct {
		name                 string
		aliasName            string
		wantType             AliasType
		wantCommands         []string
		wantMaxPositionalArg int
		wantPlaceholderCount int
		wantError            bool
	}{
		{
			name:                 "simple alias",
			aliasName:            "st",
			wantType:             SimpleAlias,
			wantCommands:         []string{"status"},
			wantMaxPositionalArg: -1,
			wantPlaceholderCount: 0,
			wantError:            false,
		},
		{
			name:                 "sequence alias - short",
			aliasName:            "acp",
			wantType:             SequenceAlias,
			wantCommands:         []string{"add", "commit", "push"},
			wantMaxPositionalArg: -1,
			wantPlaceholderCount: 0,
			wantError:            false,
		},
		{
			name:                 "sequence alias - long",
			aliasName:            "sync",
			wantType:             SequenceAlias,
			wantCommands:         []string{"pull", "add", "commit", "push"},
			wantMaxPositionalArg: -1,
			wantPlaceholderCount: 0,
			wantError:            false,
		},
		{
			name:                 "simple alias with placeholder",
			aliasName:            "deploy-msg",
			wantType:             SimpleAlias,
			wantCommands:         []string{"commit -m '{0}'"},
			wantMaxPositionalArg: 0,
			wantPlaceholderCount: 1,
			wantError:            false,
		},
		{
			name:                 "sequence alias with placeholders",
			aliasName:            "deploy",
			wantType:             SequenceAlias,
			wantCommands:         []string{"branch checkout {0}", "push {0}"},
			wantMaxPositionalArg: 0,
			wantPlaceholderCount: 1,
			wantError:            false,
		},
		{
			name:                 "non-existent alias",
			aliasName:            "nonexistent",
			wantType:             SimpleAlias,
			wantCommands:         nil,
			wantMaxPositionalArg: -1,
			wantPlaceholderCount: 0,
			wantError:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alias, err := config.ParseAlias(tt.aliasName)

			if tt.wantError {
				if err == nil {
					t.Errorf("ParseAlias() error = %v, wantErr %v", err, tt.wantError)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseAlias() unexpected error = %v", err)
				return
			}

			if alias.Type != tt.wantType {
				t.Errorf("ParseAlias() type = %v, want %v", alias.Type, tt.wantType)
			}

			if len(alias.Commands) != len(tt.wantCommands) {
				t.Errorf("ParseAlias() commands length = %v, want %v", len(alias.Commands), len(tt.wantCommands))
			}

			for i, cmd := range alias.Commands {
				if i < len(tt.wantCommands) && cmd != tt.wantCommands[i] {
					t.Errorf("ParseAlias() command[%d] = %v, want %v", i, cmd, tt.wantCommands[i])
				}
			}

			if alias.MaxPositionalArg != tt.wantMaxPositionalArg {
				t.Errorf("ParseAlias() MaxPositionalArg = %v, want %v", alias.MaxPositionalArg, tt.wantMaxPositionalArg)
			}

			if len(alias.Placeholders) != tt.wantPlaceholderCount {
				t.Errorf("ParseAlias() placeholder count = %v, want %v", len(alias.Placeholders), tt.wantPlaceholderCount)
			}
		})
	}
}

func TestAnalyzePlaceholders(t *testing.T) {
	tests := []struct {
		name                 string
		commands             []string
		wantPlaceholders     map[string]struct{}
		wantMaxPositionalArg int
		wantError            bool
		wantErrorMsg         string
	}{
		{
			name:                 "no placeholders",
			commands:             []string{"status", "branch"},
			wantPlaceholders:     map[string]struct{}{},
			wantMaxPositionalArg: -1,
			wantError:            false,
		},
		{
			name:                 "single positional placeholder",
			commands:             []string{"branch checkout {0}"},
			wantPlaceholders:     map[string]struct{}{"0": {}},
			wantMaxPositionalArg: 0,
			wantError:            false,
		},
		{
			name:                 "multiple positional placeholders",
			commands:             []string{"commit -m '{0}'", "push {1}"},
			wantPlaceholders:     map[string]struct{}{"0": {}, "1": {}},
			wantMaxPositionalArg: 1,
			wantError:            false,
		},
		{
			name:                 "named placeholders",
			commands:             []string{"branch checkout {env}", "push {branch}"},
			wantPlaceholders:     map[string]struct{}{"env": {}, "branch": {}},
			wantMaxPositionalArg: -1,
			wantError:            false,
		},
		{
			name:                 "mixed placeholders",
			commands:             []string{"commit -m '{0} on {env}'"},
			wantPlaceholders:     map[string]struct{}{"0": {}, "env": {}},
			wantMaxPositionalArg: 0,
			wantError:            false,
		},
		{
			name:         "invalid placeholder with shell chars",
			commands:     []string{"echo {0; rm -rf /}"},
			wantError:    true,
			wantErrorMsg: "placeholder contains unsafe characters",
		},
		{
			name:         "empty placeholder",
			commands:     []string{"echo {}"},
			wantError:    true,
			wantErrorMsg: "empty placeholder",
		},
		{
			name:         "placeholder with spaces",
			commands:     []string{"echo {hello world}"},
			wantError:    true,
			wantErrorMsg: "placeholder contains invalid character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			placeholders, maxArg, err := analyzePlaceholders(tt.commands)

			if tt.wantError {
				if err == nil {
					t.Errorf("analyzePlaceholders() expected error containing %q, got nil", tt.wantErrorMsg)
					return
				}
				if !strings.Contains(err.Error(), tt.wantErrorMsg) {
					t.Errorf("analyzePlaceholders() error = %q, want error containing %q", err.Error(), tt.wantErrorMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("analyzePlaceholders() unexpected error = %v", err)
				return
			}

			if maxArg != tt.wantMaxPositionalArg {
				t.Errorf("analyzePlaceholders() maxArg = %v, want %v", maxArg, tt.wantMaxPositionalArg)
			}

			if len(placeholders) != len(tt.wantPlaceholders) {
				t.Errorf("analyzePlaceholders() placeholders count = %v, want %v", len(placeholders), len(tt.wantPlaceholders))
			}

			for placeholder := range tt.wantPlaceholders {
				if _, exists := placeholders[placeholder]; !exists {
					t.Errorf("analyzePlaceholders() missing placeholder %q", placeholder)
				}
			}
		})
	}
}

func TestValidatePlaceholder(t *testing.T) {
	tests := []struct {
		name        string
		placeholder string
		wantError   bool
		wantErrMsg  string
	}{
		{
			name:        "valid positional",
			placeholder: "0",
			wantError:   false,
		},
		{
			name:        "valid named",
			placeholder: "env",
			wantError:   false,
		},
		{
			name:        "valid with underscore",
			placeholder: "env_name",
			wantError:   false,
		},
		{
			name:        "valid with hyphen",
			placeholder: "env-name",
			wantError:   false,
		},
		{
			name:        "empty placeholder",
			placeholder: "",
			wantError:   true,
			wantErrMsg:  "empty placeholder",
		},
		{
			name:        "semicolon injection",
			placeholder: "0; rm -rf /",
			wantError:   true,
			wantErrMsg:  "placeholder contains unsafe characters",
		},
		{
			name:        "pipe injection",
			placeholder: "env | cat",
			wantError:   true,
			wantErrMsg:  "placeholder contains unsafe characters",
		},
		{
			name:        "command substitution",
			placeholder: "$(whoami)",
			wantError:   true,
			wantErrMsg:  "placeholder contains unsafe characters",
		},
		{
			name:        "backtick injection",
			placeholder: "`whoami`",
			wantError:   true,
			wantErrMsg:  "placeholder contains unsafe characters",
		},
		{
			name:        "space in placeholder",
			placeholder: "hello world",
			wantError:   true,
			wantErrMsg:  "placeholder contains invalid character",
		},
		{
			name:        "special characters",
			placeholder: "test@domain",
			wantError:   true,
			wantErrMsg:  "placeholder contains invalid character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePlaceholder(tt.placeholder)

			if tt.wantError {
				if err == nil {
					t.Errorf("validatePlaceholder() expected error containing %q, got nil", tt.wantErrMsg)
					return
				}
				if !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("validatePlaceholder() error = %q, want error containing %q", err.Error(), tt.wantErrMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("validatePlaceholder() unexpected error = %v", err)
			}
		})
	}
}

func TestConfig_ParseAliasPlaceholderEdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		aliases    map[string]interface{}
		aliasName  string
		wantError  bool
		wantErrMsg string
	}{
		{
			name: "alias with malformed placeholder - missing closing brace",
			aliases: map[string]interface{}{
				"bad": "commit -m '{0'",
			},
			aliasName: "bad",
			wantError: false, // Regex won't match, so no error but no placeholders detected
		},
		{
			name: "alias with malformed placeholder - missing opening brace",
			aliases: map[string]interface{}{
				"bad": "commit -m '0}'",
			},
			aliasName: "bad",
			wantError: false, // Regex won't match, so no error but no placeholders detected
		},
		{
			name: "alias with invalid placeholder containing injection",
			aliases: map[string]interface{}{
				"malicious": "echo {0; rm -rf /}",
			},
			aliasName:  "malicious",
			wantError:  true,
			wantErrMsg: "placeholder contains unsafe characters",
		},
		{
			name: "sequence alias with mixed valid and invalid placeholders",
			aliases: map[string]interface{}{
				"mixed": []interface{}{"status {0}", "echo {bad|placeholder}"},
			},
			aliasName:  "mixed",
			wantError:  true,
			wantErrMsg: "placeholder contains unsafe characters",
		},
		{
			name: "alias with very long placeholder name",
			aliases: map[string]interface{}{
				"long": "echo {this_is_a_very_long_placeholder_name_that_should_still_be_valid_123}",
			},
			aliasName: "long",
			wantError: false,
		},
		{
			name: "alias with placeholder containing unicode",
			aliases: map[string]interface{}{
				"unicode": "echo {测试}",
			},
			aliasName:  "unicode",
			wantError:  true,
			wantErrMsg: "placeholder contains invalid character",
		},
		{
			name: "sequence alias with non-sequential placeholders",
			aliases: map[string]interface{}{
				"nonseq": []interface{}{"echo {0}", "echo {2}"},
			},
			aliasName: "nonseq",
			wantError: false, // This should be valid - gaps are allowed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				Aliases: tt.aliases,
			}

			_, err := config.ParseAlias(tt.aliasName)

			if tt.wantError {
				if err == nil {
					t.Errorf("ParseAlias() expected error containing %q, got nil", tt.wantErrMsg)
					return
				}
				if !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("ParseAlias() error = %q, want error containing %q", err.Error(), tt.wantErrMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseAlias() unexpected error = %v", err)
			}
		})
	}
}

func TestAnalyzePlaceholdersEdgeCases(t *testing.T) {
	tests := []struct {
		name                 string
		commands             []string
		wantPlaceholders     map[string]struct{}
		wantMaxPositionalArg int
		wantError            bool
		wantErrorMsg         string
	}{
		{
			name:         "command with nested braces",
			commands:     []string{"echo '{message: {0}}'"},
			wantError:    true,
			wantErrorMsg: "placeholder contains unsafe characters",
		},
		{
			name:                 "command with multiple placeholders on same line",
			commands:             []string{"echo {0} {1} {0}"},
			wantPlaceholders:     map[string]struct{}{"0": {}, "1": {}},
			wantMaxPositionalArg: 1,
			wantError:            false,
		},
		{
			name:                 "command with placeholder containing numbers and letters",
			commands:             []string{"echo {arg0}", "echo {arg1}"},
			wantPlaceholders:     map[string]struct{}{"arg0": {}, "arg1": {}},
			wantMaxPositionalArg: -1,
			wantError:            false,
		},
		{
			name:         "command with placeholder containing spaces",
			commands:     []string{"echo {hello world}"},
			wantError:    true,
			wantErrorMsg: "placeholder contains invalid character",
		},
		{
			name:         "command with placeholder containing special characters",
			commands:     []string{"echo {test@domain.com}"},
			wantError:    true,
			wantErrorMsg: "placeholder contains invalid character",
		},
		{
			name:                 "command with high numbered positional placeholder",
			commands:             []string{"echo {9}"},
			wantPlaceholders:     map[string]struct{}{"9": {}},
			wantMaxPositionalArg: 9,
			wantError:            false,
		},
		{
			name:                 "empty command list",
			commands:             []string{},
			wantPlaceholders:     map[string]struct{}{},
			wantMaxPositionalArg: -1,
			wantError:            false,
		},
		{
			name:         "command with just braces no content",
			commands:     []string{"echo {}"},
			wantError:    true,
			wantErrorMsg: "empty placeholder",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			placeholders, maxArg, err := analyzePlaceholders(tt.commands)

			if tt.wantError {
				if err == nil {
					t.Errorf("analyzePlaceholders() expected error containing %q, got nil", tt.wantErrorMsg)
					return
				}
				if !strings.Contains(err.Error(), tt.wantErrorMsg) {
					t.Errorf("analyzePlaceholders() error = %q, want error containing %q", err.Error(), tt.wantErrorMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("analyzePlaceholders() unexpected error = %v", err)
				return
			}

			if maxArg != tt.wantMaxPositionalArg {
				t.Errorf("analyzePlaceholders() maxArg = %v, want %v", maxArg, tt.wantMaxPositionalArg)
			}

			if len(placeholders) != len(tt.wantPlaceholders) {
				t.Errorf("analyzePlaceholders() placeholders count = %v, want %v", len(placeholders), len(tt.wantPlaceholders))
			}

			for placeholder := range tt.wantPlaceholders {
				if _, exists := placeholders[placeholder]; !exists {
					t.Errorf("analyzePlaceholders() missing placeholder %q", placeholder)
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
		{
			name: "sequence alias with semicolon injection blocked at config load",
			aliases: map[string]interface{}{
				"malicious": []interface{}{"status; echo pwned"},
			},
			wantError: true,
			errorMsg:  "unsafe shell metacharacters",
		},
		{
			name: "sequence alias with pipe injection blocked at config load",
			aliases: map[string]interface{}{
				"malicious": []interface{}{"status | cat"},
			},
			wantError: true,
			errorMsg:  "unsafe shell metacharacters",
		},
		{
			name: "sequence alias with ampersand injection blocked at config load",
			aliases: map[string]interface{}{
				"malicious": []interface{}{"status && echo pwned"},
			},
			wantError: true,
			errorMsg:  "unsafe shell metacharacters",
		},
		{
			name: "sequence alias with command substitution blocked at config load",
			aliases: map[string]interface{}{
				"malicious": []interface{}{"status $(whoami)"},
			},
			wantError: true,
			errorMsg:  "unsafe shell metacharacters",
		},
		{
			name: "sequence alias with backtick injection blocked at config load",
			aliases: map[string]interface{}{
				"malicious": []interface{}{"`whoami`"},
			},
			wantError: true,
			errorMsg:  "unsafe shell metacharacters",
		},
		{
			name: "sequence alias with redirection blocked at config load",
			aliases: map[string]interface{}{
				"malicious": []interface{}{"status > /tmp/out"},
			},
			wantError: true,
			errorMsg:  "unsafe shell metacharacters",
		},
		{
			name: "sequence alias with invalid command blocked at config load",
			aliases: map[string]interface{}{
				"bad": []interface{}{"notacommand"},
			},
			wantError: true,
			errorMsg:  "not a valid ggc command",
		},
		{
			name: "sequence alias with mixed valid and invalid blocked at config load",
			aliases: map[string]interface{}{
				"bad": []interface{}{"status", "invalidcmd", "branch"},
			},
			wantError: true,
			errorMsg:  "not a valid ggc command",
		},
		{
			name: "sequence alias with valid commands and args passes validation",
			aliases: map[string]interface{}{
				"good": []interface{}{"branch current", "status short", "diff staged"},
			},
			wantError: false,
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

func TestManagerLoadWithKeybindingConfig(t *testing.T) {
	mockFS := NewMockFileOps()
	homeDir := filepath.Join(os.TempDir(), "ggc-home-phase2")
	t.Setenv("HOME", homeDir)
	configPath := filepath.Join(homeDir, ".ggcconfig.yaml")
	mockFS.files[configPath] = []byte(`interactive:
  profile: vi
  keybindings:
    delete_word: "Ctrl+W"
  contexts:
    input:
      keybindings:
        move_up:
          - "Ctrl+P"
          - "Ctrl+N"
  darwin:
    keybindings:
      move_down: "Ctrl+J"
  terminals:
    wezterm:
      keybindings:
        move_to_end: "Ctrl+L"
`)

	cm := newTestConfigManager()
	if err := cm.LoadWithFileOps(mockFS); err != nil {
		t.Fatalf("LoadWithFileOps returned error: %v", err)
	}

	if cm.configPath != configPath {
		t.Fatalf("configPath = %s, want %s", cm.configPath, configPath)
	}

	cfg := cm.GetConfig()
	if cfg.Interactive.Profile != "vi" {
		t.Fatalf("profile = %s, want vi", cfg.Interactive.Profile)
	}

	moveUp, ok := cfg.Interactive.Contexts.Input.Keybindings["move_up"]
	if !ok {
		t.Fatalf("expected move_up context keybinding to be present")
	}
	seq, ok := moveUp.([]interface{})
	if !ok || len(seq) != 2 {
		t.Fatalf("unexpected move_up bindings: %#v", moveUp)
	}

	if down := cfg.Interactive.Darwin.Keybindings["move_down"]; down != "Ctrl+J" {
		t.Fatalf("darwin move_down = %v, want Ctrl+J", down)
	}

	termCfg, ok := cfg.Interactive.Terminals["wezterm"]
	if !ok {
		t.Fatalf("expected wezterm terminal keybindings")
	}
	if termCfg.Keybindings["move_to_end"] != "Ctrl+L" {
		t.Fatalf("terminal override missing, got %#v", termCfg.Keybindings)
	}
}

func TestManagerSaveWithKeybindingValidation(t *testing.T) {
	cm := newTestConfigManager()
	cm.configPath = filepath.Join(os.TempDir(), "ggc", "config.yaml")
	cm.config.Interactive.Keybindings.DeleteWord = "Shift+A"
	mockFS := NewMockFileOps()

	err := cm.SaveWithFileOps(mockFS)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported key binding format") {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, exists := mockFS.files[cm.configPath]; exists {
		t.Fatal("config file should not be written when validation fails")
	}
}

func TestManagerSaveWithKeybindingOverrides(t *testing.T) {
	cm := newTestConfigManager()
	cm.configPath = filepath.Join(os.TempDir(), "ggc", "config.yaml")
	cm.config.Interactive.Keybindings.DeleteWord = "Ctrl+W"
	cm.config.Interactive.Contexts.Input.Keybindings = map[string]interface{}{
		"move_up": []interface{}{"Ctrl+P", "Ctrl+N"},
	}
	cm.config.Interactive.Contexts.Results.Keybindings = map[string]interface{}{}
	cm.config.Interactive.Contexts.Search.Keybindings = map[string]interface{}{}
	cm.config.Interactive.Darwin.Keybindings = map[string]interface{}{
		"move_down": "Ctrl+J",
	}
	cm.config.Interactive.Terminals = map[string]KeybindingsConfig{
		"wezterm": {Keybindings: map[string]interface{}{"move_to_end": "Ctrl+L"}},
	}

	mockFS := NewMockFileOps()
	if err := cm.SaveWithFileOps(mockFS); err != nil {
		t.Fatalf("SaveWithFileOps returned error: %v", err)
	}

	data, exists := mockFS.files[cm.configPath]
	if !exists {
		t.Fatal("expected config file to be written")
	}

	saved := string(data)
	if !strings.Contains(saved, "Ctrl+P") || !strings.Contains(saved, "Ctrl+N") {
		t.Fatalf("saved config missing context overrides: %s", saved)
	}
	if !strings.Contains(saved, "wezterm") || !strings.Contains(saved, "Ctrl+L") {
		t.Fatalf("saved config missing terminal override: %s", saved)
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
		cm := NewConfigManager(testutil.NewMockGitClient())
		// Use mock git client to avoid real git config operations
		cm.gitClient = testutil.NewMockGitClient()
		cm.configPath = "/nonexistent/directory/config.yaml"

		// This should not panic, even if file operations fail
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("LoadConfig should not panic, but got: %v", r)
			}
		}()

		_ = cm.LoadConfig()

		// Verify that config is still accessible (default config should be loaded)
		config := cm.GetConfig()
		if config == nil {
			t.Error("Config should not be nil after LoadConfig, even with invalid path")
		}
	})

	t.Run("LoadConfig handles missing directory gracefully", func(t *testing.T) {
		cm := NewConfigManager(testutil.NewMockGitClient())
		// Use mock git client to avoid real git config operations
		cm.gitClient = testutil.NewMockGitClient()
		cm.configPath = "/definitely/nonexistent/path/config.yaml"

		// Should not panic or crash
		_ = cm.LoadConfig()

		// Should still have a valid config object
		if cm.GetConfig() == nil {
			t.Error("Should have valid config even with missing directory")
		}
	})
}

// TestFlattenMap tests the flattenMap function by testing aliases flattening
func TestFlattenMap(t *testing.T) {
	cm := newTestConfigManager()

	// Test flattening a map with aliases
	cfg := cm.GetConfig()
	cfg.Aliases = map[string]interface{}{
		"st":   "status",
		"br":   "branch",
		"sync": []interface{}{"pull current", "push current"},
	}

	result := cm.List()

	// Check that aliases exist in the result (they're stored as a nested map in interface{})
	aliasesValue, exists := result["aliases"]
	if !exists {
		t.Error("Expected 'aliases' key to exist in flattened result")
		return
	}

	// Convert to map for checking
	aliasesMap, ok := aliasesValue.(map[string]interface{})
	if !ok {
		t.Errorf("Expected aliases to be a map, got %T", aliasesValue)
		return
	}

	// Check individual aliases
	if aliasesMap["st"] != "status" {
		t.Errorf("Expected aliases['st'] to be 'status', got %v", aliasesMap["st"])
	}
	if aliasesMap["br"] != "branch" {
		t.Errorf("Expected aliases['br'] to be 'branch', got %v", aliasesMap["br"])
	}
	if aliasesMap["sync"] == nil {
		t.Error("Expected aliases['sync'] to have a value")
	}
}

// TestFlattenMapDirect tests the flattenMap function directly
func TestFlattenMapDirect(t *testing.T) {
	cm := newTestConfigManager()

	// Create a test map directly
	testMap := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"nested": map[string]interface{}{
			"subkey": "subvalue",
		},
	}

	// Use reflection to call flattenMap directly
	value := reflect.ValueOf(testMap)
	result := make(map[string]any)

	cm.flattenMap(value, "test", result)

	// Check that the map was flattened
	if result["test.key1"] != "value1" {
		t.Errorf("Expected test.key1 to be 'value1', got %v", result["test.key1"])
	}
	if result["test.key2"] != "value2" {
		t.Errorf("Expected test.key2 to be 'value2', got %v", result["test.key2"])
	}

	// Check nested value (should be stored as-is since it's interface{})
	nested, exists := result["test.nested"]
	if !exists {
		t.Error("Expected test.nested to exist")
	}
	if nested == nil {
		t.Error("Expected test.nested to have a value")
	}
}

// TestParseEditorBinary tests the parseEditorBinary function
func TestParseEditorBinary(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "simple binary",
			input:    "vim",
			expected: "vim",
		},
		{
			name:     "binary with args",
			input:    "vim -n",
			expected: "vim",
		},
		{
			name:     "quoted path with spaces",
			input:    "\"/usr/local/bin/code\" --wait",
			expected: "/usr/local/bin/code",
		},
		{
			name:     "single quoted path",
			input:    "'/usr/local/bin/sublime text' --wait",
			expected: "/usr/local/bin/sublime text",
		},
		{
			name:     "path with tab",
			input:    "emacs\t-nw",
			expected: "emacs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseEditorBinary(tt.input)
			if result != tt.expected {
				t.Errorf("parseEditorBinary(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestValidEditorPath tests the validEditorPath function
func TestValidEditorPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "simple binary name",
			input:    "vim",
			expected: false, // No path separators
		},
		{
			name:     "relative path",
			input:    "./vim",
			expected: false, // File doesn't exist
		},
		{
			name:     "absolute path that doesn't exist",
			input:    "/nonexistent/path/editor",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validEditorPath(tt.input)
			if result != tt.expected {
				t.Errorf("validEditorPath(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestSyncFromCommandName tests the syncFromCommandName function
func TestSyncFromCommandName(t *testing.T) {
	// Create a mock client that returns specific values
	mockClient := testutil.NewMockGitClient()

	// Test with a config manager
	cm := &Manager{
		config:    getDefaultConfig(mockClient),
		gitClient: mockClient,
	}

	// Test syncing core.editor
	cm.syncFromCommandName("core.editor")

	// Test syncing merge.tool
	cm.syncFromCommandName("merge.tool")

	// Test syncing init.defaultBranch
	cm.syncFromCommandName("init.defaultBranch")

	// Test syncing color.ui
	cm.syncFromCommandName("color.ui")

	// Test syncing core.pager
	cm.syncFromCommandName("core.pager")

	// Test syncing fetch.auto
	cm.syncFromCommandName("fetch.auto")

	// Test syncing push.default
	cm.syncFromCommandName("push.default")

	// Test with unknown command
	cm.syncFromCommandName("unknown.command")
}

// TestLoadConfig tests the LoadConfig function
func TestLoadConfigErrorHandling(t *testing.T) {
	cm := newTestConfigManager()

	// This should not panic and should handle errors gracefully
	_ = cm.LoadConfig()

	// Config should still be accessible
	cfg := cm.GetConfig()
	if cfg == nil {
		t.Error("Expected config to be available after LoadConfig")
	}
}

// TestWriteTempConfig tests error cases in writeTempConfig
func TestWriteTempConfig(t *testing.T) {
	cm := newTestConfigManager()

	// Test with invalid directory
	_, err := cm.writeTempConfig("/nonexistent/directory", []byte("test"))
	if err == nil {
		t.Error("Expected error when writing to nonexistent directory")
	}

	// Test with valid directory
	tmpDir := t.TempDir()
	tmpFile, err := cm.writeTempConfig(tmpDir, []byte("test content"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(tmpFile); err != nil {
		t.Errorf("Expected temp file to be created, got error: %v", err)
	}

	// Verify content
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Errorf("Failed to read temp file: %v", err)
	}
	if string(data) != "test content" {
		t.Errorf("Expected 'test content', got %s", string(data))
	}
}

// TestReplaceConfigFile tests the replaceConfigFile function
func TestReplaceConfigFile(t *testing.T) {
	tempDir := t.TempDir()
	cm := newTestConfigManager()

	// Create source file
	srcFile := filepath.Join(tempDir, "source.yaml")
	err := os.WriteFile(srcFile, []byte("test content"), 0600)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Set destination path
	destFile := filepath.Join(tempDir, "dest.yaml")
	cm.configPath = destFile

	// Test replacing file
	err = cm.replaceConfigFile(srcFile)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify destination file exists
	if _, err := os.Stat(destFile); err != nil {
		t.Errorf("Expected destination file to exist after replace, got error: %v", err)
	}

	// Verify content
	content, err := os.ReadFile(destFile)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}
	if string(content) != "test content" {
		t.Errorf("Expected content 'test content', got %q", string(content))
	}

	// Verify source file no longer exists
	if _, err := os.Stat(srcFile); !os.IsNotExist(err) {
		t.Error("Expected source file to be removed after rename")
	}
}

// TestSyncToGitConfigErrors tests error handling in syncToGitConfig
func TestSyncToGitConfigErrors(t *testing.T) {
	cm := newTestConfigManager()

	// This should not panic even if git operations fail
	err := cm.syncToGitConfig()
	// We expect no error with mock client, but test that it doesn't panic
	if err != nil {
		t.Logf("syncToGitConfig returned error: %v", err)
	}
}

// TestNavigateOneLevel tests different path navigation scenarios
func TestNavigateOneLevel(t *testing.T) {
	cm := newTestConfigManager()
	cfg := cm.GetConfig()

	// Test navigating to struct field
	value := reflect.ValueOf(cfg).Elem()
	result, err := cm.navigateOneLevel(value, "default", []string{"default"})
	if err != nil {
		t.Errorf("Expected no error navigating to default, got %v", err)
	}
	if !result.IsValid() {
		t.Error("Expected valid result when navigating to default")
	}

	// Test navigating to non-existent field
	_, err = cm.navigateOneLevel(value, "nonexistent", []string{"nonexistent"})
	if err == nil {
		t.Error("Expected error when navigating to non-existent field")
	}

	// Test navigating into non-struct, non-map
	stringValue := reflect.ValueOf("test")
	_, err = cm.navigateOneLevel(stringValue, "field", []string{"field"})
	if err == nil {
		t.Error("Expected error when navigating into string value")
	}
}

// TestSetMapValue tests setting values in maps
func TestSetMapValue(t *testing.T) {
	cm := newTestConfigManager()

	// Create a test map
	testMap := make(map[string]interface{})
	mapValue := reflect.ValueOf(testMap)

	// Test setting a string value
	err := cm.setMapValue(mapValue, "testkey", "testvalue")
	if err != nil {
		t.Errorf("Expected no error setting map value, got %v", err)
	}

	if testMap["testkey"] != "testvalue" {
		t.Errorf("Expected testkey to be 'testvalue', got %v", testMap["testkey"])
	}

	// Test setting incompatible type (this should work since it's interface{})
	err = cm.setMapValue(mapValue, "intkey", 42)
	if err != nil {
		t.Errorf("Expected no error setting int value, got %v", err)
	}
}
