// Package main is the entry point for the ggc CLI tool.
package main

import (
	"os"

	"github.com/bmf-san/ggc/cmd"
	"github.com/bmf-san/ggc/config"
	"github.com/bmf-san/ggc/router"
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
	config.NewConfigManager().LoadConfig()
	cmd.SetVersionGetter(GetVersionInfo)
	c := cmd.NewCmd()
	r := router.NewRouter(c)
	r.Route(os.Args[1:])
}
