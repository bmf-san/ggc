// Package main is the entry point for the ggc CLI tool.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime/debug"

	"github.com/bmf-san/ggc/v8/cmd"
	"github.com/bmf-san/ggc/v8/internal/config"
	"github.com/bmf-san/ggc/v8/internal/git"
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
	// Bind a signal-aware context so Ctrl+C cancels any running git subprocess.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	client := git.NewClient().WithContext(ctx)
	cm := config.NewConfigManager(client)
	if err := cm.LoadConfig(); err != nil {
		if config.IsWarning(err) {
			// Non-fatal: persist step failed but config was loaded OK.
			_, _ = os.Stderr.WriteString("Warning: " + err.Error() + "\n")
		} else {
			// Fatal: config file is malformed or invalid.
			return err
		}
	}
	cmd.SetVersionGetter(GetVersionInfo)
	c, err := cmd.NewCmd(client, cm)
	if err != nil {
		return err
	}
	return c.Execute(args)
}

func main() {
	if err := RunApp(os.Args[1:]); err != nil {
		writeCLIError(os.Stderr, err, os.Getenv("GGC_VERBOSE") == "1")
		os.Exit(1)
	}
}

// writeCLIError renders a terminal-facing error consistently across the CLI.
//
// For *git.OpError we print a two-line summary (what failed, then the
// underlying message). The raw git command is only shown when GGC_VERBOSE=1
// because it can be long and is usually noise in normal use. Non-git errors
// keep their historical single-line format so we don't churn existing tests
// or user expectations.
func writeCLIError(w io.Writer, err error, verbose bool) {
	var opErr *git.OpError
	if errors.As(err, &opErr) {
		_, _ = fmt.Fprintf(w, "Error: %s failed\n  %s\n", opErr.Op, opErr.Err)
		if verbose && opErr.Command != "" {
			_, _ = fmt.Fprintf(w, "  command: %s\n", opErr.Command)
		}
		return
	}
	_, _ = fmt.Fprintf(w, "Error: %s\n", err.Error())
}
