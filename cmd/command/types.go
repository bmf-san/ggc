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
		Category:  c.Category,
		Summary:   c.Summary,
		Hidden:    c.Hidden,
		HandlerID: c.HandlerID,
	}
	if len(c.Aliases) > 0 {
		clone.Aliases = append([]string(nil), c.Aliases...)
	}
	if len(c.Usage) > 0 {
		clone.Usage = append([]string(nil), c.Usage...)
	}
	if len(c.Examples) > 0 {
		clone.Examples = append([]string(nil), c.Examples...)
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
	clone := SubcommandInfo{
		Name:    s.Name,
		Summary: s.Summary,
		Hidden:  s.Hidden,
	}
	if len(s.Usage) > 0 {
		clone.Usage = append([]string(nil), s.Usage...)
	}
	if len(s.Examples) > 0 {
		clone.Examples = append([]string(nil), s.Examples...)
	}
	return clone
}
