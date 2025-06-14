# gcl CLI ツール

## 概要

`gcl` は、Git 操作を効率化するための Go 製 CLI ツールです。シェルスクリプトやエイリアスの代替として、保守性・拡張性に優れた構成を目指しています。Go 標準ライブラリのみを使用し、最小限の依存で動作します。

## 特徴
- よく使う Git 操作（add/push/pull/branch/log など）を簡単なコマンドで実行
- 複数の Git 操作を統合した複合コマンド（今後実装予定）
- 対話式 UI によるブランチ・ファイル選択やメッセージ入力（今後実装予定）
- Go 標準ライブラリのみで実装

## インストール

### 1. make build でバイナリ生成

```sh
git clone <このリポジトリURL>
make build
```

`gcl` バイナリをパスの通ったディレクトリに配置してください。

### 2. go install でグローバルインストール

```sh
go install .
```

- `$GOBIN`（通常は `$HOME/go/bin`）に `gcl` バイナリがインストールされます。
- `$GOBIN` が `PATH` に含まれていれば、どこからでも `gcl` コマンドが使えます。
- もし `PATH` が通っていない場合は、以下を追加してください：

```sh
export PATH=$PATH:$(go env GOBIN)
# または
export PATH=$PATH:$HOME/go/bin
```

## 使い方

```sh
gcl <コマンド> [サブコマンド] [オプション]
```

### 主なコマンド例

- ファイルをステージング: `gcl add <file>`
- ブランチ名表示: `gcl branch current`
- 現在のブランチを push: `gcl push current`
- 現在のブランチを強制push: `gcl push force`
- 現在のブランチを pull: `gcl pull current`
- rebase付きpull: `gcl pull rebase`
- ログ表示: `gcl log simple` / `gcl log graph`
- 空コミット: `gcl commit allow-empty`
- 一時コミット: `gcl commit tmp`
- fetch --prune: `gcl fetch --prune`
- クリーン: `gcl clean files` / `gcl clean dirs`

### コマンド一覧

```
gcl
├── add <file>                # ファイルをステージング
├── branch
│   ├── current             # 現在のブランチ名を表示
│   ├── checkout            # ローカルブランチを選択して checkout（対話式, 今後実装）
│   ├── checkout-remote     # リモートから新規ブランチ checkout（対話式, 今後実装）
│   ├── delete              # ローカルブランチを選んで削除（対話式, 今後実装）
│   └── delete-merged       # マージ済みブランチを一括削除（今後実装）
├── pull
│   ├── current             # 現在のブランチを pull
│   └── rebase              # rebase付き pull
├── push
│   ├── current             # 現在のブランチを push
│   └── force               # HEAD を強制 push
├── stash
│   └── trash               # git add . && stash（今後実装）
├── log
│   ├── simple              # git log --oneline
│   └── graph               # git log --graph
├── commit
│   ├── allow-empty         # 空コミット
│   ├── tmp                 # 一時コミット
│   └── push [-i]           # 対話的 add → commit → push（今後実装）
├── fetch
│   └── --prune             # fetch --prune
├── clean
│   ├── files               # git clean -f
│   └── dirs                # git clean -d
├── reset
│   └── clean               # reset --hard HEAD + clean -fd（今後実装）
├── rebase
│   └── interactive         # HEAD~N まで対話的 rebase（今後実装）
```

## ディレクトリ構成

```
main.go                  # エントリポイント
router/                  # コマンド分岐ロジック
cmd/                     # 各コマンドのエントリ処理
  ├── add.go
  ├── branch.go
  ├── commit.go
  ├── help.go
  ├── log.go
  ├── pull.go
  ├── push.go
  ├── fetch.go
  ├── clean.go
  ...
git/                     # Git操作のラッパー
  ├── branch.go
  ├── commit.go
  ├── log.go
  ├── pull.go
  ├── push.go
  ├── fetch.go
  ├── clean.go
  ...
```

## 補完スクリプト

標準ライブラリのみで構築されているため自動補完はありませんが、bash/zsh 用の補完スクリプト（サブコマンドまで）を `tools/completions/gcl.bash` などに配置予定です。

## 今後の拡張
- `--dry-run` オプション
- コマンド実行ログ出力
- `.gclconfig` によるカスタム設定
- テスト用のモック実装切替
- 複合コマンドや対話UIの実装

---

ご意見・ご要望は Issue までお願いします。
