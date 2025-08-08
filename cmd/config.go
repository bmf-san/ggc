package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/bmf-san/ggc/config"
)

// Configureer handles config operations.
type Configureer struct {
	outputWriter io.Writer
	helper       *Helper
	execCommand  func(string, ...string) *exec.Cmd
}

// NewConfigureer creates a new Configureer instance.
func NewConfigureer() *Configureer {
	return &Configureer{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		execCommand:  exec.Command,
	}
}

// LoadConfig executes loads the configuration.
func (c *Configureer) LoadConfig() *config.Manager {
	cm := config.NewConfigManager()
	if err := cm.Load(); err != nil {
		_, _ = fmt.Fprintf(c.outputWriter, "failed to load config: %s", err)
		return nil
	}
	return cm
}

func parseAliasValue(v interface{}) ([]string, error) {
	switch val := v.(type) {
	case string:
		return []string{val}, nil
	case []interface{}:
		var result []string
		for _, item := range val {
			str, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("non-string in alias list: %v", item)
			}
			result = append(result, str)
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unexpected type: %T", v)
	}
}

func formatAliasValue(commands []string) string {
	if len(commands) == 1 {
		return commands[0]
	}
	return fmt.Sprintf("[%s]", strings.Join(commands, " -> "))
}

// Config executes config command operations with the given arguments.
func (c *Configureer) Config(args []string) {
	if len(args) == 0 {
		c.helper.ShowConfigHelp()
		return
	}
	switch args[0] {
	case "list":
		cm := c.LoadConfig()
		configs := cm.List()

		keys := make([]string, 0, len(configs))
		for key := range configs {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			val := configs[key]

			if key == "aliases" {
				if aliasMap, ok := val.(map[string]any); ok {
					for aliasName, raw := range aliasMap {
						commands, err := parseAliasValue(raw)
						if err != nil {
							_, _ = fmt.Fprintf(c.outputWriter, "%-30s = <invalid alias: %v>\n", "aliases."+aliasName, err)
							continue
						}
						formatted := formatAliasValue(commands)
						_, _ = fmt.Fprintf(c.outputWriter, "%-30s = %s\n", "aliases."+aliasName, formatted)
					}
					continue
				}
			}

			_, _ = fmt.Fprintf(c.outputWriter, "%-30s = %s\n", key, formatValue(val))
		}
		return
	case "get":
		if len(args) < 2 {
			_, _ = fmt.Fprintf(c.outputWriter, "must provide key to get (arg missing)\n")
			return
		}

		cm := c.LoadConfig()

		value, err := cm.Get(args[1])
		if err != nil {
			_, _ = fmt.Fprintf(c.outputWriter, "failed to get config value: %s", err)
		}

		_, _ = fmt.Fprintf(c.outputWriter, "%s\n", formatValue(value))
		return
	case "set":
		if len(args) < 3 {
			_, _ = fmt.Fprintf(c.outputWriter, "must provide key && value to set (arg(s) missing)\n")
			return
		}

		cm := c.LoadConfig()

		value := parseValue(args[2])
		if err := cm.Set(args[1], value); err != nil {
			_, _ = fmt.Fprintf(c.outputWriter, "failed to set config value: %s", err)
		}

		_, _ = fmt.Fprintf(c.outputWriter, "Set %s = %s\n", args[1], formatValue(value))
		return
	default:
		c.helper.ShowConfigHelp()
		return
	}
}

func formatValue(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%g", v)
	case map[string]any:
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func parseValue(value string) any {
	if b, err := strconv.ParseBool(value); err == nil {
		return b
	}
	if i, err := strconv.ParseInt(value, 10, 64); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}
	return value
}
