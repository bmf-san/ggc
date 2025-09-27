package command

// Info captures metadata for a top-level command.
type Info struct {
	Name        string
	Aliases     []string
	Category    Category
	Summary     string
	Usage       []string
	Examples    []string
	Hidden      bool
	Subcommands []SubcommandInfo
	HandlerID   string
}

// SubcommandInfo describes a subcommand surface under a top-level command.
type SubcommandInfo struct {
	Name     string
	Summary  string
	Usage    []string
	Examples []string
	Hidden   bool
}

func (c *Info) clone() Info {
	clone := Info{
		Name:      c.Name,
		Aliases:   append([]string(nil), c.Aliases...),
		Category:  c.Category,
		Summary:   c.Summary,
		Usage:     append([]string(nil), c.Usage...),
		Examples:  append([]string(nil), c.Examples...),
		Hidden:    c.Hidden,
		HandlerID: c.HandlerID,
	}
	if len(c.Subcommands) > 0 {
		clone.Subcommands = make([]SubcommandInfo, len(c.Subcommands))
		for i, sc := range c.Subcommands {
			clone.Subcommands[i] = (&sc).clone()
		}
	}
	return clone
}

func (s *SubcommandInfo) clone() SubcommandInfo {
	return SubcommandInfo{
		Name:     s.Name,
		Summary:  s.Summary,
		Usage:    append([]string(nil), s.Usage...),
		Examples: append([]string(nil), s.Examples...),
		Hidden:   s.Hidden,
	}
}
