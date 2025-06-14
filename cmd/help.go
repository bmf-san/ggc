package cmd

import "fmt"

func ShowHelp() {
	fmt.Println(`gcl: Git操作を効率化するGo製CLIツール

Usage:
  gcl <command> [subcommand] [options]

Commands:
  add <file>                ファイルをステージング
  branch current           現在のブランチ名を表示
  branch checkout          ローカルブランチを選択して checkout（対話式）
  branch checkout-remote   リモートから新規ブランチ checkout（対話式）
  branch delete            ローカルブランチを選んで削除（対話式）
  branch delete-merged     マージ済みブランチを一括削除
  pull current             現在のブランチを pull
  pull rebase              rebase付き pull
  push current             現在のブランチを push
  push force               HEAD を強制 push
  stash trash              git add . && stash
  log simple               git log --oneline
  log graph                git log --graph
  commit allow-empty       空コミット
  commit tmp               一時コミット（メッセージ "tmp"）
  commit push [-i]         対話的 add → commit → push（ファイル選択・メッセージ入力）
  fetch --prune            fetch --prune
  clean files              git clean -f
  clean dirs               git clean -d
  reset clean              reset --hard HEAD + clean -fd
  rebase interactive       HEAD~N まで対話的 rebase（件数入力 + 編集）

Examples:
  gcl add .
  gcl branch current
  gcl push current
  gcl pull current
  gcl log simple
  gcl commit allow-empty
`)
}
