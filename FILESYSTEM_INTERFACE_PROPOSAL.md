# Config Package ファイルシステム依存性の改善提案

## 🚨 現在の問題

`config` パッケージのテストが以下の副作用を持つOS操作を直接実行している：

### 1. ファイルI/O操作
```go
// config/config.go
func (cm *Manager) loadFromFile(path string) error {
    data, err := os.ReadFile(path)  // 直接OS呼び出し
    // ...
}

func (cm *Manager) Save() error {
    err := os.MkdirAll(dir, 0700)  // 直接OS呼び出し
    // ...
}
```

### 2. テストでの副作用
```go
// config/config_test.go
func TestLoadFromFile(t *testing.T) {
    tempDir := t.TempDir()  // 実際のファイルシステム操作
    configPath := filepath.Join(tempDir, "test-config.yaml")

    // 実際のファイル作成
    err := os.WriteFile(configPath, []byte(testConfig), 0644)

    // 実際の環境変数操作
    originalHome := os.Getenv("HOME")
    if err := os.Setenv("HOME", tempDir); err != nil {
        // ...
    }
}
```

## 🎯 改善提案: FileSystem インターフェースの導入

### 1. ファイルシステム抽象化インターフェース

```go
// config/filesystem.go
package config

import (
    "io"
    "os"
    "time"
)

// FileSystem abstracts file system operations for testing
type FileSystem interface {
    ReadFile(filename string) ([]byte, error)
    WriteFile(filename string, data []byte, perm os.FileMode) error
    MkdirAll(path string, perm os.FileMode) error
    Remove(name string) error
    Rename(oldpath, newpath string) error
    Stat(name string) (os.FileInfo, error)
    CreateTemp(dir, pattern string) (File, error)
    Chmod(name string, mode os.FileMode) error
}

// File abstracts file operations
type File interface {
    io.WriteCloser
    Name() string
}

// OSFileSystem implements FileSystem using real OS operations
type OSFileSystem struct{}

func (fs *OSFileSystem) ReadFile(filename string) ([]byte, error) {
    return os.ReadFile(filename)
}

func (fs *OSFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
    return os.WriteFile(filename, data, perm)
}

func (fs *OSFileSystem) MkdirAll(path string, perm os.FileMode) error {
    return os.MkdirAll(path, perm)
}

func (fs *OSFileSystem) Remove(name string) error {
    return os.Remove(name)
}

func (fs *OSFileSystem) Rename(oldpath, newpath string) error {
    return os.Rename(oldpath, newpath)
}

func (fs *OSFileSystem) Stat(name string) (os.FileInfo, error) {
    return os.Stat(name)
}

func (fs *OSFileSystem) CreateTemp(dir, pattern string) (File, error) {
    return os.CreateTemp(dir, pattern)
}

func (fs *OSFileSystem) Chmod(name string, mode os.FileMode) error {
    return os.Chmod(name, mode)
}
```

### 2. メモリ内ファイルシステム (テスト用)

```go
// config/memory_filesystem.go
package config

import (
    "bytes"
    "errors"
    "io"
    "os"
    "path/filepath"
    "strings"
    "time"
)

// MemoryFileSystem implements FileSystem in memory for testing
type MemoryFileSystem struct {
    files map[string][]byte
    dirs  map[string]bool
}

func NewMemoryFileSystem() *MemoryFileSystem {
    return &MemoryFileSystem{
        files: make(map[string][]byte),
        dirs:  make(map[string]bool),
    }
}

func (mfs *MemoryFileSystem) ReadFile(filename string) ([]byte, error) {
    data, exists := mfs.files[filename]
    if !exists {
        return nil, &os.PathError{Op: "open", Path: filename, Err: os.ErrNotExist}
    }
    return data, nil
}

func (mfs *MemoryFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
    // Ensure directory exists
    dir := filepath.Dir(filename)
    if !mfs.dirs[dir] && dir != "." {
        return &os.PathError{Op: "open", Path: filename, Err: os.ErrNotExist}
    }

    mfs.files[filename] = make([]byte, len(data))
    copy(mfs.files[filename], data)
    return nil
}

func (mfs *MemoryFileSystem) MkdirAll(path string, perm os.FileMode) error {
    mfs.dirs[path] = true
    return nil
}

func (mfs *MemoryFileSystem) Remove(name string) error {
    if _, exists := mfs.files[name]; exists {
        delete(mfs.files, name)
        return nil
    }
    return &os.PathError{Op: "remove", Path: name, Err: os.ErrNotExist}
}

func (mfs *MemoryFileSystem) Rename(oldpath, newpath string) error {
    data, exists := mfs.files[oldpath]
    if !exists {
        return &os.PathError{Op: "rename", Path: oldpath, Err: os.ErrNotExist}
    }
    mfs.files[newpath] = data
    delete(mfs.files, oldpath)
    return nil
}

func (mfs *MemoryFileSystem) Stat(name string) (os.FileInfo, error) {
    if _, exists := mfs.files[name]; exists {
        return &memoryFileInfo{name: filepath.Base(name), size: int64(len(mfs.files[name]))}, nil
    }
    if mfs.dirs[name] {
        return &memoryFileInfo{name: filepath.Base(name), isDir: true}, nil
    }
    return nil, &os.PathError{Op: "stat", Path: name, Err: os.ErrNotExist}
}

func (mfs *MemoryFileSystem) CreateTemp(dir, pattern string) (File, error) {
    name := filepath.Join(dir, "temp_"+pattern)
    return &memoryFile{name: name, fs: mfs}, nil
}

func (mfs *MemoryFileSystem) Chmod(name string, mode os.FileMode) error {
    // Memory filesystem doesn't need permission handling
    return nil
}

// memoryFile implements File interface
type memoryFile struct {
    name   string
    buffer bytes.Buffer
    fs     *MemoryFileSystem
}

func (mf *memoryFile) Write(p []byte) (n int, error) {
    return mf.buffer.Write(p)
}

func (mf *memoryFile) Close() error {
    mf.fs.files[mf.name] = mf.buffer.Bytes()
    return nil
}

func (mf *memoryFile) Name() string {
    return mf.name
}

// memoryFileInfo implements os.FileInfo
type memoryFileInfo struct {
    name  string
    size  int64
    isDir bool
}

func (mfi *memoryFileInfo) Name() string       { return mfi.name }
func (mfi *memoryFileInfo) Size() int64        { return mfi.size }
func (mfi *memoryFileInfo) Mode() os.FileMode  { return 0644 }
func (mfi *memoryFileInfo) ModTime() time.Time { return time.Now() }
func (mfi *memoryFileInfo) IsDir() bool        { return mfi.isDir }
func (mfi *memoryFileInfo) Sys() interface{}   { return nil }
```

### 3. Config Manager の修正

```go
// config/config.go
type Manager struct {
    config     *Config
    configPath string
    gitClient  git.Clienter
    fs         FileSystem  // 新規追加
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(gitClient git.Clienter) *Manager {
    return &Manager{
        config:    getDefaultConfig(gitClient),
        gitClient: gitClient,
        fs:        &OSFileSystem{}, // デフォルトは実際のファイルシステム
    }
}

// NewConfigManagerWithFS creates a new configuration manager with custom filesystem
func NewConfigManagerWithFS(gitClient git.Clienter, fs FileSystem) *Manager {
    return &Manager{
        config:    getDefaultConfig(gitClient),
        gitClient: gitClient,
        fs:        fs,
    }
}

// loadFromFile を修正
func (cm *Manager) loadFromFile(path string) error {
    data, err := cm.fs.ReadFile(path) // os.ReadFile → cm.fs.ReadFile
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

// Save を修正
func (cm *Manager) Save() error {
    dir := filepath.Dir(cm.configPath)
    if err := cm.fs.MkdirAll(dir, 0700); err != nil { // os.MkdirAll → cm.fs.MkdirAll
        return fmt.Errorf("failed to create config directory: %w", err)
    }

    data, err := yaml.Marshal(cm.config)
    if err != nil {
        return fmt.Errorf("failed to marshal config: %w", err)
    }

    if err := cm.config.Validate(); err != nil {
        return fmt.Errorf("cannot save invalid config: %w", err)
    }

    tmpName, err := cm.writeTempConfig(dir, data)
    if err != nil {
        return err
    }

    if err := cm.replaceConfigFile(tmpName); err != nil {
        return err
    }

    cm.hardenPermissions(cm.configPath)
    return cm.syncToGitConfig()
}

// 他の関数も同様に修正...
```

### 4. 改善されたテスト

```go
// config/config_test.go
func TestLoadFromFile_WithMemoryFS(t *testing.T) {
    // メモリ内ファイルシステムを使用
    memFS := NewMemoryFileSystem()
    mockClient := testutil.NewMockGitClient()

    // メモリ内にテスト設定ファイルを作成
    testConfig := `
default:
  branch: "develop"
  editor: "nano"
ui:
  color: false
`
    configPath := "/test/config.yaml"
    memFS.MkdirAll("/test", 0755)
    memFS.WriteFile(configPath, []byte(testConfig), 0644)

    // カスタムファイルシステムでConfig Managerを作成
    cm := NewConfigManagerWithFS(mockClient, memFS)

    // ファイル読み込みテスト（副作用なし）
    err := cm.loadFromFile(configPath)
    if err != nil {
        t.Fatalf("Failed to load config: %v", err)
    }

    // 設定値の検証
    if cm.config.Default.Branch != "develop" {
        t.Errorf("Expected branch 'develop', got %s", cm.config.Default.Branch)
    }
}

func TestSave_WithMemoryFS(t *testing.T) {
    memFS := NewMemoryFileSystem()
    mockClient := testutil.NewMockGitClient()

    cm := NewConfigManagerWithFS(mockClient, memFS)
    cm.configPath = "/test/config.yaml"

    // ディレクトリを作成
    memFS.MkdirAll("/test", 0755)

    // 設定を変更
    cm.config.Default.Branch = "main"
    cm.config.UI.Color = false

    // 保存（副作用なし）
    err := cm.Save()
    if err != nil {
        t.Fatalf("Failed to save config: %v", err)
    }

    // ファイルが作成されたか確認
    data, err := memFS.ReadFile("/test/config.yaml")
    if err != nil {
        t.Fatalf("Config file not created: %v", err)
    }

    // 内容を検証
    if !strings.Contains(string(data), "branch: main") {
        t.Error("Config file doesn't contain expected branch setting")
    }
}
```

## 🚀 改善効果

### 1. **テストの分離性**
- 実際のファイルシステムに依存しない
- 環境変数の変更が不要
- テスト間の相互影響を排除

### 2. **テスト実行速度**
- メモリ内操作で高速化
- 一時ファイル作成/削除のオーバーヘッド削減

### 3. **テストの信頼性**
- 並列テスト実行時の競合状態を回避
- CI/CD環境での一貫した動作

### 4. **エラーケースのテスト容易性**
- ファイルシステムエラーの模擬が簡単
- 権限エラーなどの再現が可能

## 📋 実装手順

1. **FileSystem インターフェースの定義**
2. **OSFileSystem の実装**
3. **MemoryFileSystem の実装**
4. **Manager の修正（依存性注入）**
5. **テストの書き換え**
6. **既存機能の動作確認**

## 🔍 考慮事項

### 利点
- ✅ 副作用のない純粋なユニットテスト
- ✅ 高速なテスト実行
- ✅ 複雑なエラーケースの再現が容易
- ✅ CI/CD環境での安定性向上

### 欠点
- ❌ 実装コードの複雑性が若干増加
- ❌ インターフェース維持のオーバーヘッド
- ❌ 実際のファイルシステムとの差異リスク

### 結論
**メリットがデメリットを大幅に上回るため、実装を強く推奨**

この改善により、`config_test.go`は真の意味でのユニットテストになり、副作用のない高品質なテストスイートを実現できます。
