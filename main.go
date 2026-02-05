// Package main is the entry point for the ggc CLI tool.
package main

import (
	"os"
	"runtime/debug"
	"strings"

	"github.com/bmf-san/ggc/v7/cmd"
	"github.com/bmf-san/ggc/v7/internal/config"
	"github.com/bmf-san/ggc/v7/pkg/git"
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
func RunApp(args []string) error {
	client := git.NewClient()
	cm := config.NewConfigManager(client)
	if err := cm.LoadConfig(); err != nil {
		// Continue with default config on error
		_, _ = os.Stderr.WriteString("Warning: " + err.Error() + "\n")
	}
	cmd.SetVersionGetter(GetVersionInfo)
	c := cmd.NewCmd(client, cm)
	// Cache default remote in tagger to avoid repeated config loads.
	if r := strings.TrimSpace(cm.GetConfig().Git.DefaultRemote); r != "" {
		c.SetDefaultRemote(r)
	}
	return c.Execute(args)
}

func main() {
	if err := RunApp(os.Args[1:]); err != nil {
		_, _ = os.Stderr.WriteString("Error: " + err.Error() + "\n")
		os.Exit(1)
	}
}
