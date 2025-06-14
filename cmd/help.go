package cmd

import "fmt"

func ShowHelp() {
	fmt.Print(`gcl: Git操作を効率化するGo製CLIツール

Usage:
  gcl <command> [subcommand] [options]

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
`)
}
