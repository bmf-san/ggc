package cmd

import "fmt"

func ShowHelp() {
	fmt.Print(`gcl: Git操作を効率化するGo製CLIツール

Usage:
  gcl <コマンド> [サブコマンド] [オプション]

主なコマンド:
  gcl add <file>              ファイルをステージング
  gcl branch current          現在のブランチ名を表示
  gcl branch checkout         対話的にブランチ切替
  gcl push current            現在のブランチをpush
  gcl push force              現在のブランチを強制push
  gcl pull current            現在のブランチをpull
  gcl pull rebase             rebase付きpull
  gcl log simple              シンプルなログ表示
  gcl log graph               グラフ付きログ表示
  gcl commit allow-empty      空コミット
  gcl commit tmp              一時コミット
  gcl fetch --prune           prune付きfetch
  gcl clean files             ファイルのクリーン
  gcl clean dirs              ディレクトリのクリーン
  gcl reset clean             リセット＋クリーン
  gcl commit-push             対話的にadd/commit/push一括実行

Examples:
  gcl add .
  gcl branch current
  gcl branch checkout
  gcl push current
  gcl push force
  gcl pull current
  gcl pull rebase
  gcl log simple
  gcl log graph
  gcl commit allow-empty
  gcl commit tmp
  gcl fetch --prune
  gcl clean files
  gcl clean dirs
  gcl reset clean
  gcl commit-push
`)
}
