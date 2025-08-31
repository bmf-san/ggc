// Package main is the entry point for the ggc CLI tool.
package main

import (
	"os"

	"github.com/bmf-san/ggc/v4/cmd"
	"github.com/bmf-san/ggc/v4/config"
	"github.com/bmf-san/ggc/v4/git"
	"github.com/bmf-san/ggc/v4/router"
)

var (
	version string
	commit  string
)

// GetVersionInfo returns the version information
func GetVersionInfo() (string, string) {
	return version, commit
}

func main() {
	cm := config.NewConfigManager()
	cm.LoadConfig()
	cmd.SetVersionGetter(GetVersionInfo)
	c := cmd.NewCmd(git.NewClient())
	r := router.NewRouter(c, cm)
	r.Route(os.Args[1:])
}
