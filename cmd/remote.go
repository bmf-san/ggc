package cmd

import (
	"fmt"
	"os/exec"
)

func Remote(args []string) {
	if len(args) > 0 {
		switch args[0] {
		case "list":
			remoteList()
			return
		case "add":
			if len(args) < 3 {
				fmt.Println("Usage: gcl remote add <name> <url>")
				return
			}
			remoteAdd(args[1], args[2])
			return
		case "remove":
			if len(args) < 2 {
				fmt.Println("Usage: gcl remote remove <name>")
				return
			}
			remoteRemove(args[1])
			return
		case "set-url":
			if len(args) < 3 {
				fmt.Println("Usage: gcl remote set-url <name> <url>")
				return
			}
			remoteSetURL(args[1], args[2])
			return
		}
	}
	ShowRemoteHelp()
}

func remoteList() {
	cmd := exec.Command("git", "remote", "-v")
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error: failed to get git remote -v: %v\n", err)
		return
	}
	fmt.Print(string(out))
}

func remoteAdd(name, url string) {
	cmd := exec.Command("git", "remote", "add", name, url)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: failed to add remote: %v\n", err)
		return
	}
	fmt.Printf("Remote '%s' added\n", name)
}

func remoteRemove(name string) {
	cmd := exec.Command("git", "remote", "remove", name)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: failed to remove remote: %v\n", err)
		return
	}
	fmt.Printf("Remote '%s' removed\n", name)
}

func remoteSetURL(name, url string) {
	cmd := exec.Command("git", "remote", "set-url", name, url)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: failed to set remote URL: %v\n", err)
		return
	}
	fmt.Printf("Remote '%s' URL updated\n", name)
}

func ShowRemoteHelp() {
	fmt.Println("Usage: gcl remote list | gcl remote add <name> <url> | gcl remote remove <name> | gcl remote set-url <name> <url>")
}
