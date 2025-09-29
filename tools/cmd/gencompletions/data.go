package main

import (
	"sort"
	"strings"

	"github.com/bmf-san/ggc/v6/cmd/command"
)

type TemplateData struct {
	Commands                  []*CommandData
	TopLevel                  []string
	TopLevelList              string
	CommandMap                map[string]*CommandData
	BranchCheckoutKeywordList string
}

type CommandData struct {
	Name               string
	Summary            string
	Subcommands        []SubcommandData
	SubcommandList     string
	KeywordSubcommands []SubcommandData
	IncludeInCase      bool
}

type SubcommandData struct {
	Name        string
	Summary     string
	Keywords    []string
	KeywordList string
	SkipKeyword bool
}

func (c *CommandData) Subcommand(name string) *SubcommandData {
	for i := range c.Subcommands {
		if c.Subcommands[i].Name == name {
			return &c.Subcommands[i]
		}
	}
	return nil
}

func hasKeywords(subs []SubcommandData) bool {
	for _, sub := range subs {
		if len(sub.Keywords) > 0 && !sub.SkipKeyword {
			return true
		}
	}
	return false
}

func needsHandler(cmd *CommandData) bool {
	if len(cmd.Subcommands) > 0 {
		return true
	}
	if hasKeywords(cmd.Subcommands) {
		return true
	}
	switch cmd.Name {
	case "branch", "add", "rebase":
		return true
	}
	return false
}

func subcommandBy(cmd *CommandData, name string) *SubcommandData {
	return cmd.Subcommand(name)
}

func join(items []string, sep string) string {
	return strings.Join(items, sep)
}

func escapeZsh(s string) string {
	return strings.ReplaceAll(s, "'", "'\\''")
}

func escapeFish(s string) string {
	return strings.ReplaceAll(s, "\"", "\\\"")
}

func escapeBash(s string) string {
	return strings.ReplaceAll(s, "\"", "\\\"")
}

func buildTemplateData(cmds []command.Info) *TemplateData {
	data := &TemplateData{
		Commands:   make([]*CommandData, 0, len(cmds)),
		CommandMap: make(map[string]*CommandData),
	}

	var topLevel []string
	for i := range cmds {
		cmd := &cmds[i]
		if cmd.Hidden {
			continue
		}
		topLevel = append(topLevel, cmd.Name)
		cmdData := buildCommandData(cmd)
		data.Commands = append(data.Commands, cmdData)
		data.CommandMap[cmdData.Name] = cmdData
	}

	sort.Strings(topLevel)
	data.TopLevel = topLevel
	data.TopLevelList = strings.Join(topLevel, " ")

	sort.Slice(data.Commands, func(i, j int) bool {
		return data.Commands[i].Name < data.Commands[j].Name
	})

	if branchCmd, ok := data.CommandMap["branch"]; ok {
		if checkout := branchCmd.Subcommand("checkout"); checkout != nil {
			if len(checkout.Keywords) > 0 {
				data.BranchCheckoutKeywordList = strings.Join(checkout.Keywords, " ")
			}
			checkout.SkipKeyword = true
		}
	}

	return data
}

func buildCommandData(cmd *command.Info) *CommandData {
	subMap := collectSubcommands(cmd)

	subPointers := make([]*SubcommandData, 0, len(subMap))
	for _, sub := range subMap {
		sort.Strings(sub.Keywords)
		sub.KeywordList = strings.Join(sub.Keywords, " ")
		sub.SkipKeyword = shouldSkipKeyword(cmd.Name, sub.Name)
		subPointers = append(subPointers, sub)
	}

	sort.Slice(subPointers, func(i, j int) bool {
		return subPointers[i].Name < subPointers[j].Name
	})

	subcommands := make([]SubcommandData, len(subPointers))
	keywordSubs := make([]SubcommandData, 0)
	subNames := make([]string, len(subPointers))
	for i, sub := range subPointers {
		subcommands[i] = *sub
		subNames[i] = sub.Name
		if len(sub.Keywords) > 0 && !sub.SkipKeyword {
			keywordSubs = append(keywordSubs, *sub)
		}
	}

	return &CommandData{
		Name:               cmd.Name,
		Summary:            fallbackSummary(cmd.Summary, cmd.Name),
		Subcommands:        subcommands,
		SubcommandList:     strings.Join(subNames, " "),
		KeywordSubcommands: keywordSubs,
		IncludeInCase:      shouldIncludeInCase(cmd.Name, len(subcommands) > 0),
	}
}

func collectSubcommands(cmd *command.Info) map[string]*SubcommandData {
	subMap := make(map[string]*SubcommandData)
	for i := range cmd.Subcommands {
		sub := &cmd.Subcommands[i]
		name, tail, ok := splitSubcommandTokens(cmd, sub)
		if !ok {
			continue
		}

		entry := ensureSubcommandEntry(subMap, name, sub.Summary)
		keywords := extractKeywords(tail)
		for _, kw := range keywords {
			entry.addKeyword(kw)
		}
	}
	return subMap
}

func splitSubcommandTokens(cmd *command.Info, sub *command.SubcommandInfo) (string, []string, bool) {
	if sub.Hidden {
		return "", nil, false
	}
	tokens := strings.Fields(sub.Name)
	if len(tokens) < 2 {
		return "", nil, false
	}
	if !strings.EqualFold(tokens[0], cmd.Name) && !containsFold(cmd.Aliases, tokens[0]) {
		return "", nil, false
	}

	name := sanitizeSubcommandName(tokens[1])
	if name == "" {
		return "", nil, false
	}
	return name, tokens[2:], true
}

func ensureSubcommandEntry(subMap map[string]*SubcommandData, name, summary string) *SubcommandData {
	entry, ok := subMap[name]
	if !ok {
		entry = &SubcommandData{
			Name:    name,
			Summary: fallbackSummary(summary, name),
		}
		subMap[name] = entry
		return entry
	}
	if entry.Summary == "" {
		entry.Summary = fallbackSummary(summary, name)
	}
	return entry
}

func (s *SubcommandData) addKeyword(keyword string) {
	for _, existing := range s.Keywords {
		if existing == keyword {
			return
		}
	}
	s.Keywords = append(s.Keywords, keyword)
}

func fallbackSummary(summary, name string) string {
	if strings.TrimSpace(summary) == "" {
		return strings.TrimSpace(name)
	}
	return summary
}

func sanitizeSubcommandName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	if name == "." {
		return ""
	}
	if strings.HasPrefix(name, "<") || strings.HasPrefix(name, "[") || strings.HasPrefix(name, "(") {
		return ""
	}
	return name
}

func extractKeywords(tokens []string) []string {
	var out []string
	for _, token := range tokens {
		keyword := sanitizeKeyword(token)
		if keyword == "" {
			continue
		}
		out = append(out, keyword)
	}
	return out
}

func sanitizeKeyword(token string) string {
	token = strings.TrimSpace(token)
	token = strings.Trim(token, ",")
	if token == "" {
		return ""
	}
	if strings.HasPrefix(token, "<") || strings.Contains(token, "|") || strings.HasPrefix(token, "[") || strings.HasPrefix(token, "(") {
		return ""
	}
	if token == "." || token == ".." {
		return ""
	}
	return token
}

func shouldSkipKeyword(commandName, subcommandName string) bool {
	if commandName == "branch" && subcommandName == "checkout" {
		return true
	}
	return false
}

func shouldIncludeInCase(commandName string, hasSubcommands bool) bool {
	if !hasSubcommands {
		return false
	}
	if commandName == "add" {
		return false
	}
	return true
}

func containsFold(values []string, candidate string) bool {
	for _, value := range values {
		if strings.EqualFold(value, candidate) {
			return true
		}
	}
	return false
}
