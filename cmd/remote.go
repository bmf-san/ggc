package cmd

import (
	"fmt"
	"os/exec"
)

type Remoteer struct {
	execCommand func(name string, arg ...string) *exec.Cmd
}

func NewRemoteer() *Remoteer {
	return &Remoteer{execCommand: exec.Command}
}

func (r *Remoteer) Remote(args []string) {
	if len(args) > 0 {
		switch args[0] {
		case "list":
			r.remoteList()
			return
		case "add":
			if len(args) < 3 {
				fmt.Println("Usage: ggc remote add <name> <url>")
				return
			}
			r.remoteAdd(args[1], args[2])
			return
		case "remove":
			if len(args) < 2 {
				fmt.Println("Usage: ggc remote remove <name>")
				return
			}
			r.remoteRemove(args[1])
			return
		case "set-url":
			if len(args) < 3 {
				fmt.Println("Usage: ggc remote set-url <name> <url>")
				return
			}
			r.remoteSetURL(args[1], args[2])
			return
		}
	}
	ShowRemoteHelp()
}

func (r *Remoteer) remoteList() {
	cmd := r.execCommand("git", "remote", "-v")
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error: failed to get git remote -v: %v\n", err)
		return
	}
	fmt.Print(string(out))
}

func (r *Remoteer) remoteAdd(name, url string) {
	cmd := r.execCommand("git", "remote", "add", name, url)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: failed to add remote: %v\n", err)
		return
	}
	fmt.Printf("Remote '%s' added\n", name)
}

func (r *Remoteer) remoteRemove(name string) {
	cmd := r.execCommand("git", "remote", "remove", name)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: failed to remove remote: %v\n", err)
		return
	}
	fmt.Printf("Remote '%s' removed\n", name)
}

func (r *Remoteer) remoteSetURL(name, url string) {
	cmd := r.execCommand("git", "remote", "set-url", name, url)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: failed to set remote URL: %v\n", err)
		return
	}
	fmt.Printf("Remote '%s' URL updated\n", name)
}

func ShowRemoteHelp() {
	fmt.Println("Usage: ggc remote list | ggc remote add <name> <url> | ggc remote remove <name> | ggc remote set-url <name> <url>")
}
