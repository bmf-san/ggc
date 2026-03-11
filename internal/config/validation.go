package config

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (c *Config) validateBranch() error {
	branch := c.Default.Branch
	if strings.TrimSpace(branch) == "" || strings.Contains(branch, " ") {
		return &ValidationError{"default.branch", branch, "must not contain spaces or be empty"}
	}
	return nil
}

func (c *Config) validateEditor() error {
	editor := strings.TrimSpace(c.Default.Editor)
	bin := parseEditorBinary(editor)
	if validEditorPath(bin) {
		return nil
	}
	if _, err := exec.LookPath(bin); err != nil {
		return &ValidationError{"default.editor", editor, "command not found in PATH or invalid path"}
	}
	return nil
}

func parseEditorBinary(editor string) string {
	if editor == "" {
		return ""
	}
	// Support basic quoted paths or first token before whitespace
	if (strings.HasPrefix(editor, "\"") && strings.Count(editor, "\"") >= 2) || (strings.HasPrefix(editor, "'") && strings.Count(editor, "'") >= 2) {
		q := editor[0:1]
		if idx := strings.Index(editor[1:], q); idx >= 0 {
			return editor[1 : 1+idx]
		}
	}
	if i := strings.IndexAny(editor, " \t"); i > 0 {
		return editor[:i]
	}
	return editor
}

func validEditorPath(bin string) bool {
	if !strings.ContainsAny(bin, "/\\") {
		return false
	}
	_, err := os.Stat(bin)
	return err == nil
}

func (c *Config) validateConfirmDestructive() error {
	val := c.Behavior.ConfirmDestructive
	valid := map[string]bool{"simple": true, "always": true, "never": true}
	if !valid[val] {
		return &ValidationError{"behavior.confirm-destructive", val, "must be one of: simple, always, never"}
	}
	return nil
}

// validateGitDefaultRemote validates git default remote name format
func (c *Config) validateGitDefaultRemote() error {
	remote := c.Git.DefaultRemote
	if remote == "" {
		return nil
	}

	if !gitRemoteNameCharsRe.MatchString(remote) || strings.Contains(remote, " ") {
		return &ValidationError{
			Field:   "git.default-remote",
			Value:   remote,
			Message: "Remote may contain letters, digits, ., _, -, and / only",
		}
	}

	// Additional structural checks: no leading/trailing '.' or '/', and no empty/unsafe segments
	if strings.HasPrefix(remote, "/") || strings.HasSuffix(remote, "/") || strings.HasPrefix(remote, ".") || strings.HasSuffix(remote, ".") || strings.Contains(remote, "//") || strings.Contains(remote, "..") {
		return &ValidationError{
			Field:   "git.default-remote",
			Value:   remote,
			Message: "Remote must not start/end with '.' or '/', nor contain '..' or '//'",
		}
	}

	return nil
}

// validateWorkflows validates all workflows defined in the configuration.
// Workflow names must be non-empty and contain no spaces. Each workflow must
// have at least one step, and each step must be a non-empty command string that
// does not contain shell metacharacters.
//
// Two placeholder syntaxes are permitted within step strings:
//   - Interactive placeholder syntax: <name> (angle brackets, e.g. "commit <message>").
//     These are stripped before the metacharacter check via angleBracketPlaceholderRe.
//   - Alias-style positional placeholders: {0}, {1}, … (curly braces, e.g. "commit {0}").
//     These are stripped by defaultValidator.validateCommand via aliasPlaceholderPattern.
func (c *Config) validateWorkflows() error {
	for name, steps := range c.Workflows {
		if strings.TrimSpace(name) == "" || strings.Contains(name, " ") {
			return &ValidationError{
				Field:   "workflows." + name,
				Value:   name,
				Message: "workflow names must not be empty or contain spaces",
			}
		}
		if len(steps) == 0 {
			return &ValidationError{
				Field:   "workflows." + name,
				Value:   name,
				Message: "workflow must have at least one step",
			}
		}
		for i, step := range steps {
			if strings.TrimSpace(step) == "" {
				return &ValidationError{
					Field:   fmt.Sprintf("workflows.%s[%d]", name, i),
					Value:   step,
					Message: "step command must not be empty",
				}
			}
			// Strip <placeholder> tokens (interactive workflow syntax) before
			// the metacharacter check so that e.g. "commit <message>" is valid.
			// Note: alias-style placeholders like {0} are stripped by
			// defaultValidator.validateCommand, so both forms are permitted
			// in workflow step commands.
			cleaned := angleBracketPlaceholderRe.ReplaceAllString(step, "")
			if err := defaultValidator.validateCommand(cleaned); err != nil {
				return &ValidationError{
					Field:   fmt.Sprintf("workflows.%s[%d]", name, i),
					Value:   step,
					Message: err.Error(),
				}
			}
		}
	}
	return nil
}

// Validate is a function that handles validation operations
func (c *Config) Validate() error {
	if err := c.validateBranch(); err != nil {
		return err
	}
	if err := c.validateEditor(); err != nil {
		return err
	}
	if err := c.validateConfirmDestructive(); err != nil {
		return err
	}
	if err := c.validateGitDefaultRemote(); err != nil {
		return err
	}
	if err := c.validateAliases(); err != nil {
		return err
	}
	if err := c.validateKeybindings(); err != nil {
		return err
	}
	if err := c.validateWorkflows(); err != nil {
		return err
	}
	return nil
}
