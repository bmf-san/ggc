package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/bmf-san/ggc/git"
)

func Clean(args []string) {
	if len(args) > 0 {
		switch args[0] {
		case "files":
			err := git.CleanFiles()
			if err != nil {
				fmt.Println("Error:", err)
			}
			return
		case "dirs":
			err := git.CleanDirs()
			if err != nil {
				fmt.Println("Error:", err)
			}
			return
		}
	}
	ShowCleanHelp()
}

func ShowCleanHelp() {
	fmt.Println("Usage: ggc clean files | ggc clean dirs")
}

// Interactively select files to clean
func CleanInteractive() {
	cmd := exec.Command("git", "clean", "-nd") // get candidates with dry-run
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error: failed to get candidates with git clean -nd: %v\n", err)
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
		fmt.Println("No files to clean.")
		return
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\033[1;36mSelect files to delete by number (space separated, all: select all, none: deselect all, e.g. 1 3 5):\033[0m")
		for i, f := range files {
			fmt.Printf("  [\033[1;33m%d\033[0m] %s\n", i+1, f)
		}
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			fmt.Println("Cancelled.")
			return
		}
		if input == "all" {
			args := append([]string{"clean", "-f", "--"}, files...)
			cleanCmd := exec.Command("git", args...)
			cleanCmd.Stdout = os.Stdout
			cleanCmd.Stderr = os.Stderr
			if err := cleanCmd.Run(); err != nil {
				fmt.Printf("Error: failed to clean files: %v\n", err)
				return
			}
			fmt.Println("Selected files deleted.")
			break
		}
		if input == "none" {
			continue
		}
		indices := strings.Fields(input)
		tmp := []string{}
		valid := true
		for _, idx := range indices {
			n, err := strconv.Atoi(idx)
			if err != nil || n < 1 || n > len(files) {
				fmt.Printf("\033[1;31mInvalid number: %s\033[0m\n", idx)
				valid = false
				break
			}
			tmp = append(tmp, files[n-1])
		}
		if !valid {
			continue
		}
		if len(tmp) == 0 {
			fmt.Println("\033[1;33mNothing selected.\033[0m")
			continue
		}
		fmt.Printf("\033[1;32mSelected files: %v\033[0m\n", tmp)
		fmt.Print("Delete these files? (y/n): ")
		ans, _ := reader.ReadString('\n')
		ans = strings.TrimSpace(ans)
		if ans == "y" || ans == "Y" {
			args := append([]string{"clean", "-f", "--"}, tmp...)
			cleanCmd := exec.Command("git", args...)
			cleanCmd.Stdout = os.Stdout
			cleanCmd.Stderr = os.Stderr
			if err := cleanCmd.Run(); err != nil {
				fmt.Printf("Error: failed to clean files: %v\n", err)
				return
			}
			fmt.Println("Selected files deleted.")
			break
		}
	}
}
