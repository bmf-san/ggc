// Package main is the entry point for the ggc CLI tool.
package main

import (
	"os"
	"runtime/debug"
	"strings"

	"github.com/bmf-san/ggc/v6/cmd"
	"github.com/bmf-san/ggc/v6/config"
	"github.com/bmf-san/ggc/v6/git"
	"github.com/bmf-san/ggc/v6/router"
)

var (
	version string
	commit  string
)

// GetVersionInfo returns the version information
func GetVersionInfo() (string, string) {
	// Prefer ldflags-injected values when available
	if version != "" || commit != "" {
		return version, commit
	}

	// Fallback for `go install`: use module build info
	if bi, ok := debug.ReadBuildInfo(); ok {
		v := bi.Main.Version
		// Treat test/dev builds as unset
		if v == "(devel)" {
			v = ""
		}
		var rev string
		for _, s := range bi.Settings {
			if s.Key == "vcs.revision" {
				if len(s.Value) >= 7 {
					rev = s.Value[:7]
				} else {
					rev = s.Value
				}
				break
			}
		}
		return v, rev
	}

	return "", ""
}

// RunApp contains the main application logic, separated for testability.
// This function initializes all components and routes the provided arguments.
func RunApp(args []string) {
	cm := config.NewConfigManager(git.NewClient())
	cm.LoadConfig()
	cmd.SetVersionGetter(GetVersionInfo)
	c := cmd.NewCmd(git.NewClient())
	// Cache default remote in tagger to avoid repeated config loads.
	if r := strings.TrimSpace(cm.GetConfig().Integration.Github.DefaultRemote); r != "" {
		c.SetDefaultRemote(r)
	}
	r := router.NewRouter(c, cm)
	r.Route(args)
}

func main() {
	RunApp(os.Args[1:])
}
