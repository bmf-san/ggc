package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/bmf-san/gcl/git"
)

// Complete handles dynamic completion for subcommands/args
func Complete(args []string) {
	if len(args) < 1 {
		return
	}
	switch args[0] {
	case "branch":
		branches, err := git.ListLocalBranches()
		if err != nil {
			return
		}
		for _, b := range branches {
			fmt.Println(b)
		}
	case "files":
		// git ls-files で管理下ファイル一覧を取得
		cmd := exec.Command("git", "ls-files")
		out, err := cmd.Output()
		if err != nil {
			return
		}
		files := strings.Split(strings.TrimSpace(string(out)), "\n")
		for _, f := range files {
			fmt.Println(f)
		}
	default:
		// 今後他の補完も追加可能
	}
}
