package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/bmf-san/ggc/git"
)

type Logger struct {
	execCommand func(name string, arg ...string) *exec.Cmd
	logSimple   func() error
}

func NewLogger() *Logger {
	return &Logger{
		execCommand: exec.Command,
		logSimple:   git.LogSimple,
	}
}

func (l *Logger) Log(args []string) {
	if len(args) > 0 {
		switch args[0] {
		case "simple":
			err := l.logSimple()
			if err != nil {
				fmt.Println("Error:", err)
			}
			return
		case "graph":
			err := l.logGraph()
			if err != nil {
				fmt.Println("Error:", err)
			}
			return
		}
	}
	ShowLogHelp()
}

func (l *Logger) logGraph() error {
	cmd := l.execCommand("git", "log", "--graph")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ShowLogHelp() {
	fmt.Println("Usage: ggc log simple | ggc log graph")
}

// For backward compatibility
func Log(args []string) {
	NewLogger().Log(args)
}
