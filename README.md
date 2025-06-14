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

| gclコマンド例                | 実際に実行されるgitコマンド           | 説明                       |
|-----------------------------|--------------------------------------|----------------------------|
| gcl add <file>              | git add <file>                       | ファイルをステージング     |
| gcl branch current          | git rev-parse --abbrev-ref HEAD       | 現在のブランチ名を表示     |
| gcl branch checkout         | git branch ... → git checkout <選択>   | 対話的にブランチ切替       |
| gcl push current            | git push origin <branch>              | 現在のブランチをpush       |
| gcl push force              | git push --force origin <branch>      | 現在のブランチを強制push   |
| gcl pull current            | git pull origin <branch>              | 現在のブランチをpull       |
| gcl pull rebase             | git pull --rebase origin <branch>     | rebase付きpull             |
| gcl log simple              | git log --oneline                     | シンプルなログ表示         |
| gcl log graph               | git log --graph                       | グラフ付きログ表示         |
| gcl commit allow-empty      | git commit --allow-empty -m ...        | 空コミット                 |
| gcl commit tmp              | git commit -m "tmp"                   | 一時コミット               |
| gcl fetch --prune           | git fetch --prune                     | prune付きfetch             |
| gcl clean files             | git clean -f                          | ファイルのクリーン         |
| gcl clean dirs              | git clean -d                          | ディレクトリのクリーン     |
| gcl reset clean             | git reset --hard HEAD; git clean -fd  | リセット＋クリーン         |

### 主なコマンド例

- gcl add .
- gcl branch current
- gcl branch checkout
- gcl push current
- gcl push force
- gcl pull current
- gcl pull rebase
- gcl log simple
- gcl log graph
- gcl commit allow-empty
- gcl commit tmp
- gcl fetch --prune
- gcl clean files
- gcl clean dirs
- gcl reset clean

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
