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
	reader := bufio.NewReader(os.Stdin)
	selected := []string{}
	for {
		fmt.Println("\033[1;36m削除するファイルを番号で選択（スペース区切り, all:全選択, none:全解除, 例: 1 3 5）:\033[0m")
		for i, f := range files {
			fmt.Printf("  [\033[1;33m%d\033[0m] %s\n", i+1, f)
		}
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			fmt.Println("キャンセルしました")
			return
		}
		if input == "all" {
			selected = files
			break
		}
		if input == "none" {
			selected = []string{}
			continue
		}
		indices := strings.Fields(input)
		tmp := []string{}
		valid := true
		for _, idx := range indices {
			n, err := strconv.Atoi(idx)
			if err != nil || n < 1 || n > len(files) {
				fmt.Printf("\033[1;31m無効な番号: %s\033[0m\n", idx)
				valid = false
				break
			}
			tmp = append(tmp, files[n-1])
		}
		if !valid {
			continue
		}
		selected = tmp
		if len(selected) == 0 {
			fmt.Println("\033[1;33m何も選択されませんでした\033[0m")
			continue
		}
		fmt.Printf("\033[1;32m選択したファイル: %v\033[0m\n", selected)
		fmt.Print("このファイルを削除しますか？ (y/n): ")
		ans, _ := reader.ReadString('\n')
		ans = strings.TrimSpace(ans)
		if ans == "y" || ans == "Y" {
			break
		}
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
