package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/bmf-san/ggc/v8/internal/config"
)

// Execute executes the command with alias resolution.
// This is the main entry point that handles both aliases and regular commands.
//
// It returns a non-nil error if executing an alias fails, such as when alias parsing
// or placeholder processing encounters an error; interactive mode and regular commands
// do not cause Execute to return an error.
func (c *Cmd) Execute(args []string) error {
	if len(args) == 0 {
		c.Interactive()
		return nil
	}

	cmdName, cmdArgs := args[0], args[1:]

	// Check if this is an alias
	if c.configManager != nil && c.configManager.GetConfig().IsAlias(cmdName) {
		return c.executeAlias(cmdName, cmdArgs)
	}

	// Regular command
	return c.Route(args)
}

// executeAlias resolves and executes an alias command identified by name.
// It returns an error if the alias cannot be parsed or if placeholder
// processing fails (e.g., insufficient arguments for placeholders).
func (c *Cmd) executeAlias(name string, args []string) error {
	cfg := c.configManager.GetConfig()
	alias, err := cfg.ParseAlias(name)
	if err != nil {
		return fmt.Errorf("error parsing alias: %w", err)
	}

	switch alias.Type {
	case config.SimpleAlias:
		return c.executeSimpleAlias(alias, args, name)
	case config.SequenceAlias:
		return c.executeSequenceAlias(alias, args, name)
	}
	return nil
}

// executeSimpleAlias executes a simple (single-command) alias.
func (c *Cmd) executeSimpleAlias(alias *config.ParsedAlias, args []string, name string) error {
	processedCommands, err := c.processPlaceholders(alias, args, name)
	if err != nil {
		return err
	}

	command := tokenize(processedCommands[0])
	if len(alias.Placeholders) == 0 {
		// No placeholders, forward user arguments
		return c.Route(append([]string{command[0]}, args...))
	}
	// Placeholders were processed, use the processed command
	return c.Route(command)
}

// executeSequenceAlias executes a sequence alias (multiple commands in order).
func (c *Cmd) executeSequenceAlias(alias *config.ParsedAlias, args []string, name string) error {
	processedCommands, err := c.processPlaceholders(alias, args, name)
	if err != nil {
		return err
	}

	for _, cmd := range processedCommands {
		_, _ = fmt.Fprintf(c.outputWriter, "Executing: %s\n", cmd)
		command := tokenize(cmd)
		if err := c.Route(command); err != nil {
			return err
		}
	}
	return nil
}

// processPlaceholders resolves placeholders in the given alias commands using the
// supplied args and returns a slice of fully-expanded command strings.
//
// When an alias defines no placeholders, simple aliases return their configured
// commands as-is, while sequence aliases reject any provided arguments and
// return an error. When placeholders are present, this function validates that
// enough arguments have been supplied for the highest positional placeholder
// index in use and returns an error if the requirements are not met.
func (c *Cmd) processPlaceholders(alias *config.ParsedAlias, args []string, aliasName string) ([]string, error) {
	// If no placeholders are used, handle arguments appropriately
	if len(alias.Placeholders) == 0 {
		if alias.Type == config.SequenceAlias && len(args) > 0 {
			return nil, fmt.Errorf("sequence alias '%s' does not accept arguments (got %s)", aliasName, strings.Join(args, " "))
		}
		// For simple aliases without placeholders, arguments are forwarded as usual
		return alias.Commands, nil
	}

	// Validate that we have enough arguments for positional placeholders.
	// Note: MaxPositionalArg is 0-indexed (the highest placeholder index used),
	// so if MaxPositionalArg = 0, we need at least 1 argument (for {0}).
	if alias.MaxPositionalArg >= 0 && len(args) <= alias.MaxPositionalArg {
		return nil, fmt.Errorf("alias '%s' requires at least %d argument(s), got %d",
			aliasName, alias.MaxPositionalArg+1, len(args))
	}

	// Process each command
	processedCommands := make([]string, len(alias.Commands))
	for i, cmd := range alias.Commands {
		processed := cmd

		// Replace positional placeholders {0}, {1}, etc.
		for j := 0; j <= alias.MaxPositionalArg; j++ {
			placeholder := fmt.Sprintf("{%d}", j)
			if strings.Contains(processed, placeholder) {
				if j < len(args) {
					processed = strings.ReplaceAll(processed, placeholder, args[j])
				}
			}
		}

		// Replace well-known named placeholders. These do not consume args;
		// they are derived from the current repo / config / time.
		processed = c.replaceNamedPlaceholders(processed)

		processedCommands[i] = processed
	}

	return processedCommands, nil
}

// replaceNamedPlaceholders substitutes a small, well-known set of named
// placeholders that are derived from environment rather than arguments.
// Supported names:
//
//	{branch} - current branch name
//	{remote} - default remote (from config, fallback "origin")
//	{date}   - today's date in YYYY-MM-DD (local time)
//
// Unknown named placeholders are left untouched so that future additions do
// not silently drop user text. This function is cheap: it short-circuits
// when no "{" is present.
func (c *Cmd) replaceNamedPlaceholders(s string) string {
	if !strings.Contains(s, "{") {
		return s
	}
	replacements := map[string]func() string{
		"{branch}": func() string {
			if c.gitClient == nil {
				return ""
			}
			name, err := c.gitClient.GetCurrentBranch()
			if err != nil {
				return ""
			}
			return name
		},
		"{remote}": func() string {
			if c.configManager == nil {
				return "origin"
			}
			r := strings.TrimSpace(c.configManager.GetConfig().Git.DefaultRemote)
			if r == "" {
				return "origin"
			}
			return r
		},
		"{date}": func() string {
			return time.Now().Format("2006-01-02")
		},
	}
	for placeholder, resolve := range replacements {
		if strings.Contains(s, placeholder) {
			s = strings.ReplaceAll(s, placeholder, resolve())
		}
	}
	return s
}
