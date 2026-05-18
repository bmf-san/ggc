package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bmf-san/ggc/v8/internal/history"
)

// History displays or searches the persisted ggc command history.
func (c *Cmd) History(args []string) {
	// Default: show last 20
	if len(args) == 0 {
		c.showLast(20)
		return
	}
	switch args[0] {
	case "last":
		if len(args) >= 2 {
			n, err := strconv.Atoi(args[1])
			if err != nil || n <= 0 {
				fmt.Fprintln(c.outputWriter, "Error: invalid number for last")
				return
			}
			c.showLast(n)
			return
		}
		c.showLast(20)
	case "search":
		if len(args) >= 2 {
			pattern := strings.Join(args[1:], " ")
			lines, err := history.Search(pattern)
			if err != nil {
				fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
				return
			}
			for _, l := range lines {
				fmt.Fprintln(c.outputWriter, l)
			}
			return
		}
		fmt.Fprintln(c.outputWriter, "Usage: ggc history search <pattern>")
	default:
		// allow numeric immediate e.g., `ggc history 10`
		if n, err := strconv.Atoi(args[0]); err == nil && n > 0 {
			c.showLast(n)
			return
		}
		fmt.Fprintln(c.outputWriter, "Usage: ggc history [last N|search <pattern>]")
	}
}

func (c *Cmd) showLast(n int) {
	lines, err := history.ReadLast(n)
	if err != nil {
		fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
		return
	}
	for _, l := range lines {
		fmt.Fprintln(c.outputWriter, l)
	}
}
