package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/bmf-san/gcl/git"
)

func Clean(args []string) {
	if len(args) > 0 {
		switch args[0] {
		case "files":
			err := git.CleanFiles()
			if err != nil {
				fmt.Println("エラー:", err)
			}
			return
		case "dirs":
			err := git.CleanDirs()
			if err != nil {
				fmt.Println("エラー:", err)
			}
			return
		}
	}
	ShowCleanHelp()
}

func ShowCleanHelp() {
	fmt.Println("使用例: gcl clean files | gcl clean dirs")
}

// 対話的に削除候補ファイルを選択してクリーン
func CleanInteractive() {
	cmd := exec.Command("git", "clean", "-nd") // dry-runで候補取得
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("エラー: git clean -nd の取得に失敗しました: %v\n", err)
		return
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	files := []string{}
	for _, line := range lines {
		if strings.HasPrefix(line, "Would remove ") {
			files = append(files, strings.TrimPrefix(line, "Would remove "))
		}
	}
	if len(files) == 0 {
		fmt.Println("削除候補ファイルはありません")
		return
	}
	fmt.Println("削除するファイルを番号で選択（スペース区切り, 例: 1 3 5）:")
	for i, f := range files {
		fmt.Printf("  [%d] %s\n", i+1, f)
	}
	fmt.Print("> ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		fmt.Println("キャンセルしました")
		return
	}
	indices := strings.Fields(input)
	selected := []string{}
	for _, idx := range indices {
		n, err := strconv.Atoi(idx)
		if err != nil || n < 1 || n > len(files) {
			fmt.Printf("無効な番号: %s\n", idx)
			return
		}
		selected = append(selected, files[n-1])
	}
	if len(selected) == 0 {
		fmt.Println("何も選択されませんでした")
		return
	}
	// git clean -f -- <file1> <file2> ...
	args := append([]string{"clean", "-f", "--"}, selected...)
	cleanCmd := exec.Command("git", args...)
	cleanCmd.Stdout = os.Stdout
	cleanCmd.Stderr = os.Stderr
	if err := cleanCmd.Run(); err != nil {
		fmt.Printf("エラー: git clean に失敗しました: %v\n", err)
		return
	}
	fmt.Println("選択したファイルを削除しました")
}
