# gcl CLI ツール設計ドキュメント

## 概要

`gcl` は Git 操作を効率化するための CLI ツールであり、Go 標準ライブラリのみを用いて実装される。複雑化しがちなシェルスクリプトやエイリアスの代替として、保守性・拡張性に優れた構成を目指す。

---

## 要件

### 機能要件

* Git 操作の簡略化（push/pull/rebase/stash/branch/log など）
* 現在のブランチ名取得、対話的ブランチ選択などのユーティリティ機能
* 複数の Git 操作を統合する複合コマンド（例：add + commit + push）
* 対話式 UI による選択や入力（ブランチ選択、ファイル選択、メッセージ入力）

### 非機能要件

* Go 標準ライブラリのみを使用
* クロスプラットフォーム対応（Unix系 / Windows）
* 最小限の依存構成、ビルドサイズの抑制

---

## CLIコマンド構成

```
gcl
├── branch
│   ├── current             # 現在のブランチ名を表示
│   ├── checkout            # ローカルブランチを選択して checkout（対話式）
│   ├── checkout-remote     # リモートから新規ブランチ checkout（対話式）
│   ├── delete              # ローカルブランチを選んで削除（対話式）
│   └── delete-merged       # マージ済みブランチを一括削除
├── pull
│   ├── current             # 現在のブランチを pull
│   └── rebase              # rebase付き pull
├── push
│   ├── current             # 現在のブランチを push
│   └── force               # HEAD を強制 push
├── stash
│   └── trash               # git add . && stash
├── log
│   ├── simple              # git log
│   └── graph               # git log --graph
├── commit
│   ├── allow-empty         # 空コミット
│   ├── tmp                 # 一時コミット（メッセージ "tmp"）
│   └── push [-i]           # 対話的 add → commit → push（ファイル選択・メッセージ入力）
├── fetch
│   └── --prune             # fetch --prune（旧 update）
├── clean
│   ├── files               # git clean -f
│   └── dirs                # git clean -d
├── reset
│   └── clean               # reset --hard HEAD + clean -fd
├── rebase
│   └── interactive         # HEAD~N まで対話的 rebase（件数入力 + 編集）
```

---

## ディレクトリ構成

```
gcl/
├── main.go                  # CLIルーター（エントリポイント）
├── router/
│   └── router.go            # os.Args を使ったコマンド分岐ロジック
├── cmd/
│   ├── branch.go            # 各コマンドのエントリ処理
│   ├── commit.go
│   └── ...
├── git/
│   ├── branch.go            # Git操作のロジックラッパー（内部処理）
│   ├── push.go
│   └── ...
├── ui/
│   └── prompt.go            # 対話UI関連処理（ブランチ・ファイル選択、入力プロンプト）
```

---

## 実装方針

* コマンドライン引数は `os.Args` を使って自前でパース
* 引数に応じたハンドラ関数を `cmd/` 以下に配置
* Git 操作は基本的に `exec.Command("git", ...)` を用いて `git` CLI を実行
* 現在のブランチ取得など共通処理は `git/` に集約
* 対話式選択や入力が必要な箇所は `ui/` に分離（`fmt.Scanln`, `select`, `bufio.NewReader`）

---

## シェル補完対応

* `gcl` コマンドは標準ライブラリのみで構築されるため、補完機能（tab補完）は自動では有効にならない。
* 補完対応のため、手動で **bash/zsh 用の補完スクリプト** を用意する。
* bash では `complete -F _gcl_completions gcl` により補完可能。
* 補完スクリプトは今後、`tools/completions/gcl.bash` 等に配置・管理予定。

---

## 今後の拡張余地

* `--dry-run` オプションなどの実装
* コマンドの実行ログ出力
* `.gclconfig` によるカスタム設定
* テスト用のモック実装切替（例：`git.FakeExec`）

---

以上が、標準ライブラリベースで構築する `gcl` CLI ツールの初期設計である。
