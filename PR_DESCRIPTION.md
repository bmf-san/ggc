# テストカバレッジの大幅改善とGitパッケージ構造のリファクタリング

## 📊 概要

このPRでは、プロジェクト全体のテストカバレッジを大幅に改善し、gitパッケージの構造をより一貫性のある設計にリファクタリングしました。

## 🎯 主な成果

### テストカバレッジの改善

| パッケージ | Before | After | 改善幅 |
|-----------|--------|-------|--------|
| **cmd** | 54.4% | **60.2%** | +5.8% |
| **router** | 78.3% | **93.5%** | +15.2% 🌟 |
| **config** | 78.3% | **82.2%** | +3.9% |
| **git** | 85.2% | **89.9%** | +4.7% |
| **全体** | ~65% | **69.3%** | +4.3% |

### 新規作成されたテストファイル

#### cmdパッケージ (7つの新しいテストファイル)
- `cmd/diff_test.go` - Diffコマンドの包括的テスト
- `cmd/remote_test.go` - Remoteコマンドの基本テスト
- `cmd/status_test.go` - Statusコマンドとユーティリティ関数のテスト
- `cmd/reset_test.go` - Resetコマンドのテスト
- `cmd/restore_test.go` - Restoreコマンドのテスト
- `cmd/stash_test.go` - Stashコマンドのテスト
- `cmd/tag_test.go` - Tagコマンドとユーティリティメソッドのテスト

#### gitパッケージ (新しい構造に対応)
- `git/rev-list_test.go` - rev-listコマンド関数の包括的テスト
- `git/ls-files_test.go` - ls-filesコマンド関数のテスト

## 🏗️ リファクタリング内容

### gitパッケージの構造改善

**Before: 一貫性のない設計**
```
util.go (混在)
├── ListFiles() - git ls-files
├── GetUpstreamBranchName() - git rev-parse
└── GetAheadBehindCount() - git rev-list

rev-parse.go
├── GetCurrentBranch()
├── GetBranchName()
├── RevParseVerify()
└── GetCommitHash()

tag.go
└── GetTagCommit() - git rev-list (分散)
```

**After: コマンドベースの一貫した設計**
```
rev-parse.go (全てのrev-parseコマンド)
├── GetCurrentBranch()
├── GetBranchName()
├── RevParseVerify()
├── GetCommitHash()
└── GetUpstreamBranchName() ← 移動

rev-list.go (全てのrev-listコマンド)
├── GetAheadBehindCount() ← 移動
└── GetTagCommit() ← 移動

ls-files.go (全てのls-filesコマンド)
└── ListFiles() ← 移動

util.go → 削除 (全ての関数を適切なファイルに移動)
```

### 改善されたテスト品質

#### 1. **エラーケースの包括的カバレッジ**
```go
// 正常系テスト
func TestClient_Diff(t *testing.T) { /* ... */ }

// エラー系テスト (新規追加)
func TestClient_Diff_Error(t *testing.T) {
    client := &Client{
        execCommand: func(name string, args ...string) *exec.Cmd {
            return exec.Command("false") // 失敗するコマンド
        },
    }
    _, err := client.Diff()
    if err == nil {
        t.Error("Expected Diff to return an error")
    }
}
```

#### 2. **エイリアス機能の完全テスト (router)**
```go
func TestRouter_WithAliases(t *testing.T) {
    tests := []struct {
        name     string
        alias    config.Alias
        args     []string
        expected []string
    }{
        {
            name: "SimpleAlias",
            alias: config.Alias{
                Type:  config.SimpleAlias,
                Value: "status --short",
            },
            args:     []string{},
            expected: []string{"status", "--short"},
        },
        // ... 他のテストケース
    }
}
```

#### 3. **設定管理の詳細テスト (config)**
```go
func TestFlattenMapDirect(t *testing.T) {
    input := map[string]interface{}{
        "aliases": map[string]interface{}{
            "st": "status",
            "co": "checkout",
        },
    }
    result := flattenMap(input, "")
    // 詳細な検証...
}
```

## 🔧 技術的改善

### 1. **testutilパッケージの適切な除外**
```makefile
# Makefile
cover:
	go test $$(go list ./... | grep -v testutil) -coverprofile=coverage.out
	go tool cover -func=coverage.out
```

### 2. **モックの活用と標準化**
```go
// testutil.NewMockGitClient()を統一使用
mockClient := testutil.NewMockGitClient()
tagger := NewTagger(mockClient)
```

### 3. **テーブル駆動テストの活用**
```go
tests := []struct {
    name           string
    args           []string
    expectedOutput string
    wantArgs       []string
}{
    // 複数のテストケースを効率的に管理
}
```

## ✅ 品質保証

### テスト実行結果
```bash
=== 全てのテストが成功 ===
ok github.com/bmf-san/ggc/v5/cmd    1.185s coverage: 60.2%
ok github.com/bmf-san/ggc/v5/config 0.735s coverage: 82.2%
ok github.com/bmf-san/ggc/v5/git    1.598s coverage: 89.9%
ok github.com/bmf-san/ggc/v5/router 1.470s coverage: 93.5%
```

### カバレッジの100%達成関数
- **git/diff.go**: 全関数100%
- **git/status.go**: 全関数100%
- **git/stash.go**: 全関数100%
- **git/rebase.go**: 全関数100%
- **git/rev-parse.go**: 全関数100%
- **git/rev-list.go**: 全関数100%
- **git/ls-files.go**: 全関数100%

## 🚀 今後の展望

1. **mainパッケージの改善** (現在42.9%) - エントリーポイントの詳細テスト
2. **インタラクティブ機能のテスト強化** - UI操作の詳細テスト
3. **統合テストの追加** - パッケージ間連携のテスト

## 📋 変更ファイル一覧

### 新規作成 (10ファイル)
- `cmd/diff_test.go`
- `cmd/remote_test.go`
- `cmd/status_test.go`
- `cmd/reset_test.go`
- `cmd/restore_test.go`
- `cmd/stash_test.go`
- `cmd/tag_test.go`
- `git/rev-list.go`
- `git/rev-list_test.go`
- `git/ls-files.go`
- `git/ls-files_test.go`

### 更新 (8ファイル)
- `Makefile` - testutil除外設定
- `router/router_test.go` - エイリアステスト追加
- `config/config_test.go` - 設定関数テスト追加
- `git/rev-parse.go` - GetUpstreamBranchName追加
- `git/rev-parse_test.go` - 新関数テスト追加
- `git/diff_test.go` - エラーケーステスト追加
- `git/status_test.go` - エラーケーステスト追加
- `git/stash_test.go` - エラーケーステスト追加
- `git/rebase_test.go` - 包括的テスト追加
- `git/tag.go` - GetTagCommit関数削除
- `git/tag_test.go` - 重複テスト削除
- `internal/testutil/git_client.go` - ドキュメント追加

### 削除 (2ファイル)
- `git/util.go` - 関数を適切なファイルに移動
- `git/util_test.go` - 対応するテストも移動

---

このPRにより、プロジェクトの**テスト品質**と**コード構造**が大幅に改善され、より保守しやすく信頼性の高いコードベースとなりました。
