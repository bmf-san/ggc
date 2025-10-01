// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	commandregistry "github.com/bmf-san/ggc/v7/cmd/command"
	"github.com/bmf-san/ggc/v7/cmd/templates"
)

// Helper provides help message functionality.
type Helper struct {
	outputWriter io.Writer
}

// NewHelper creates a new Helper.
func NewHelper() *Helper {
	return &Helper{
		outputWriter: os.Stdout,
	}
}

// ShowHelp shows the main help message.
func (h *Helper) ShowHelp() {
	helpMsg, err := templates.RenderMainHelp()
	if err != nil {
		_, _ = fmt.Fprintf(h.outputWriter, "Error: %v\n", err)
		return
	}
	_, _ = fmt.Fprint(h.outputWriter, helpMsg)
}

// ShowCommandHelp shows help message for a command.
func (h *Helper) ShowCommandHelp(data templates.HelpData) {
	helpMsg, err := templates.RenderCommandHelp(data)
	if err != nil {
		_, _ = fmt.Fprintf(h.outputWriter, "Error: %v\n", err)
		return
	}
	_, _ = fmt.Fprint(h.outputWriter, helpMsg)
}

func (h *Helper) renderCommandFromRegistry(name string, usageOverride []string, descriptionOverride string) {
	h.renderCommandFromRegistryWithFilter(name, usageOverride, descriptionOverride, nil)
}

func (h *Helper) renderCommandFromRegistryWithFilter(name string, usageOverride []string, descriptionOverride string, filter func(commandregistry.SubcommandInfo) bool) {
	info, ok := commandregistry.Find(name)
	if !ok {
		usage := usageOverride
		if len(usage) == 0 {
			usage = []string{fmt.Sprintf("ggc %s", name)}
		}
		h.ShowCommandHelp(templates.HelpData{
			Usage:       strings.Join(usage, " | "),
			Description: fmt.Sprintf("No help available for '%s'", name),
		})
		return
	}

	data := buildHelpData(&info, usageOverride, descriptionOverride, filter)
	h.ShowCommandHelp(data)
}

func buildHelpData(info *commandregistry.Info, usageOverride []string, descriptionOverride string, filter func(commandregistry.SubcommandInfo) bool) templates.HelpData {
	usageList := usageOverride
	if len(usageList) == 0 {
		if filter != nil {
			usageList = collectSubcommandUsages(info, filter)
		}
		if len(usageList) == 0 {
			usageList = info.Usage
		}
	}
	if len(usageList) == 0 {
		usageList = []string{fmt.Sprintf("ggc %s", info.Name)}
	}
	usage := strings.Join(uniqueStrings(usageList), " | ")

	description := descriptionOverride
	if description == "" {
		description = info.Summary
	}

	examples := buildExamples(info, filter)
	if len(examples) == 0 {
		examples = uniqueStrings(usageList)
	}

	return templates.HelpData{
		Usage:       usage,
		Description: description,
		Examples:    examples,
	}
}

func collectSubcommandUsages(info *commandregistry.Info, filter func(commandregistry.SubcommandInfo) bool) []string {
	var usages []string
	for _, sub := range info.Subcommands {
		if sub.Hidden {
			continue
		}
		if filter != nil && !filter(sub) {
			continue
		}
		usage := firstNonEmpty(sub.Usage, fmt.Sprintf("ggc %s", sub.Name))
		if usage != "" {
			usages = append(usages, usage)
		}
	}
	return uniqueStrings(usages)
}

func buildExamples(info *commandregistry.Info, filter func(commandregistry.SubcommandInfo) bool) []string {
	var examples []string
	examples = append(examples, info.Examples...)
	for _, sub := range info.Subcommands {
		if sub.Hidden {
			continue
		}
		if filter != nil && !filter(sub) {
			continue
		}
		if len(sub.Examples) > 0 {
			examples = append(examples, sub.Examples...)
			continue
		}
		usage := firstNonEmpty(sub.Usage, fmt.Sprintf("ggc %s", sub.Name))
		if usage == "" {
			continue
		}
		examples = append(examples, formatExample(usage, sub.Summary))
	}
	return uniqueStrings(examples)
}

func firstNonEmpty(values []string, fallback string) string {
	for _, v := range values {
		trimmed := strings.TrimSpace(v)
		if trimmed != "" {
			return trimmed
		}
	}
	return fallback
}

func formatExample(usage, summary string) string {
	usage = strings.TrimSpace(usage)
	if summary == "" {
		return usage
	}
	return fmt.Sprintf("%s  # %s", usage, summary)
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	var result []string
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		result = append(result, trimmed)
	}
	return result
}

// ShowAddHelp shows help message for add command.
func (h *Helper) ShowAddHelp() {
	h.renderCommandFromRegistry("add", nil, "")
}

// ShowBranchHelp shows help message for branch command.
func (h *Helper) ShowBranchHelp() {
	h.renderCommandFromRegistry("branch", []string{"ggc branch <command>"}, "")
}

// ShowCleanHelp shows help message for clean command.
func (h *Helper) ShowCleanHelp() {
	h.renderCommandFromRegistry("clean", []string{"ggc clean <command>"}, "")
}

// ShowCommitHelp shows help message for commit command.
func (h *Helper) ShowCommitHelp() {
	h.renderCommandFromRegistry("commit", nil, "Commit staged changes")
}

// ShowLogHelp shows help message for log command.
func (h *Helper) ShowLogHelp() {
	h.renderCommandFromRegistry("log", []string{"ggc log <command>"}, "Show commit logs")
}

// ShowPullHelp shows help message for pull command.
func (h *Helper) ShowPullHelp() {
	h.renderCommandFromRegistry("pull", []string{"ggc pull <command>"}, "Pull changes from remote")
}

// ShowPushHelp shows help message for push command.
func (h *Helper) ShowPushHelp() {
	h.renderCommandFromRegistry("push", []string{"ggc push <command>"}, "Push changes to remote")
}

// ShowRemoteHelp shows help message for remote command.
func (h *Helper) ShowRemoteHelp() {
	h.renderCommandFromRegistry("remote", []string{"ggc remote <command>"}, "Manage set of tracked repositories")
}

// ShowStashHelp shows help message for stash command.
func (h *Helper) ShowStashHelp() {
	h.renderCommandFromRegistry("stash", []string{"ggc stash [command]"}, "Stash changes")
}

// ShowHookHelp displays help information for hook commands.
func (h *Helper) ShowHookHelp() {
	h.renderCommandFromRegistry("hook", []string{"ggc hook [command]"}, "Manage Git hooks")
}

// ShowConfigHelp shows help message for config command.
func (h *Helper) ShowConfigHelp() {
	h.renderCommandFromRegistry("config", []string{"ggc config [command]"}, "Get, set, and list configuration values for ggc")
}

// ShowRestoreHelp shows help message for restore command.
func (h *Helper) ShowRestoreHelp() {
	h.renderCommandFromRegistry("restore", []string{"ggc restore [command]"}, "Restore working tree files")
}

// ShowStatusHelp shows help message for status command.
func (h *Helper) ShowStatusHelp() {
	h.renderCommandFromRegistry("status", []string{"ggc status [command]"}, "Show the working tree status")
}

// ShowTagHelp shows help message for tag command.
func (h *Helper) ShowTagHelp() {
	h.renderCommandFromRegistry("tag", []string{"ggc tag [command] [options]"}, "Create, list, delete and verify tags")
}

// ShowVersionHelp shows help message for Version command.
func (h *Helper) ShowVersionHelp() {
	h.renderCommandFromRegistry("version", nil, "Show current ggc version")
}

// ShowRebaseHelp shows help message for rebase command.
func (h *Helper) ShowRebaseHelp() {
	h.renderCommandFromRegistry("rebase", []string{"ggc rebase [interactive | <upstream> | continue | abort | skip]"}, "Rebase current branch onto another branch; supports interactive and common workflows")
}

// ShowResetHelp shows help message for reset command.
func (h *Helper) ShowResetHelp() {
	h.renderCommandFromRegistry("reset", nil, "Reset and clean")
}

// ShowListBranchesHelp displays help for the list branches command.
func (h *Helper) ShowListBranchesHelp() {
	filter := func(sub commandregistry.SubcommandInfo) bool {
		return strings.HasPrefix(sub.Name, "branch list ")
	}
	h.renderCommandFromRegistryWithFilter("branch", nil, "List local or remote branches", filter)
}

// ShowDeleteBranchHelp displays help for the delete branch command.
func (h *Helper) ShowDeleteBranchHelp() {
	filter := func(sub commandregistry.SubcommandInfo) bool {
		return sub.Name == "branch delete"
	}
	h.renderCommandFromRegistryWithFilter("branch", []string{"ggc branch delete <branch-name> [--force]"}, "Delete a branch", filter)
}

// ShowDeleteMergedBranchHelp displays help for the delete merged branch command.
func (h *Helper) ShowDeleteMergedBranchHelp() {
	filter := func(sub commandregistry.SubcommandInfo) bool {
		return sub.Name == "branch delete merged"
	}
	h.renderCommandFromRegistryWithFilter("branch", []string{"ggc branch delete merged"}, "Delete merged branches", filter)
}

// ShowDiffHelp displays help for the git diff command.
func (h *Helper) ShowDiffHelp() {
	h.renderCommandFromRegistry(
		"diff",
		[]string{"ggc diff [staged|unstaged|head] [options] [<commit> [<commit>]] [--] [<path>...]"},
		"Show changes between commits, the index, and the working tree",
	)
}

// ShowFetchHelp shows help message for fetch command.
func (h *Helper) ShowFetchHelp() {
	h.renderCommandFromRegistry("fetch", []string{"ggc fetch [subcommand]"}, "Download objects and refs from another repository")
}
