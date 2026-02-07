package interactive

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// deriveArgsFromDescription extracts arguments from a description string.
// Returns the args portion (everything after the command name).
func deriveArgsFromDescription(description string) []string {
	parts := strings.Fields(description)
	if len(parts) > 1 {
		return parts[1:]
	}
	return nil
}

// collectPlaceholders extracts unique placeholders from a list of arguments.
func collectPlaceholders(args []string) []string {
	var placeholders []string
	seen := make(map[string]bool)
	for _, arg := range args {
		for _, ph := range extractPlaceholders(arg) {
			if !seen[ph] {
				seen[ph] = true
				placeholders = append(placeholders, ph)
			}
		}
	}
	return placeholders
}

// replacePlaceholdersInArgs replaces placeholders in each argument with their values.
func replacePlaceholdersInArgs(args []string, inputs map[string]string) []string {
	resolvedArgs := make([]string, len(args))
	for i, arg := range args {
		resolved := arg
		for ph, val := range inputs {
			resolved = strings.ReplaceAll(resolved, "<"+ph+">", val)
		}
		resolvedArgs[i] = resolved
	}
	return resolvedArgs
}

// resolveStepPlaceholders resolves placeholders in a workflow step's arguments.
// Each argument is processed individually, preserving multiword placeholder values as single arguments.
func resolveStepPlaceholders(ui *UI, step WorkflowStep) ([]string, bool) {
	// If Args is empty, derive from Description
	args := step.Args
	if len(args) == 0 {
		args = deriveArgsFromDescription(step.Description)
	}

	// Extract unique placeholders from all args
	placeholders := collectPlaceholders(args)
	if len(placeholders) == 0 {
		return args, false
	}

	// Get user input for each placeholder
	inputs, canceled := interactiveInputForWorkflow(ui, placeholders)
	if canceled {
		return nil, true
	}

	return replacePlaceholdersInArgs(args, inputs), false
}

// interactiveInputForWorkflow provides interactive input for placeholders during workflow execution
func interactiveInputForWorkflow(ui *UI, placeholders []string) (map[string]string, bool) {
	if ui != nil && ui.handler != nil {
		return interactiveInputForWorkflowUI(ui, placeholders)
	}
	scanner := bufio.NewScanner(os.Stdin)
	return interactiveInputForWorkflowScanner(scanner, placeholders)
}

func interactiveInputForWorkflowUI(ui *UI, placeholders []string) (map[string]string, bool) {
	inputs := make(map[string]string)
	for i, ph := range placeholders {
		ui.write("\n")
		if len(placeholders) > 1 {
			ui.write("%s[%d/%d]%s ",
				ui.colors.BrightBlue+ui.colors.Bold,
				i+1, len(placeholders),
				ui.colors.Reset)
		}
		ui.write("%s? %s%s%s: ",
			ui.colors.BrightGreen,
			ui.colors.BrightWhite+ui.colors.Bold,
			ph,
			ui.colors.Reset)

		value, canceled := ui.readPlaceholderInput()
		if canceled {
			return nil, true
		}
		if strings.TrimSpace(value) == "" {
			return nil, true
		}

		inputs[ph] = value
		ui.write("%s✓ %s%s: %s%s%s\n",
			ui.colors.BrightGreen,
			ui.colors.BrightBlue,
			ph,
			ui.colors.BrightYellow+ui.colors.Bold,
			value,
			ui.colors.Reset)
	}
	return inputs, false
}

func interactiveInputForWorkflowScanner(scanner *bufio.Scanner, placeholders []string) (map[string]string, bool) {
	inputs := make(map[string]string)
	for i, ph := range placeholders {
		if len(placeholders) > 1 {
			fmt.Printf("\n[%d/%d] ", i+1, len(placeholders))
		} else {
			fmt.Print("\n")
		}

		fmt.Printf("? %s: ", ph)

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Printf("Input error: %v\n", err)
			}
			return nil, true
		}
		value := strings.TrimSpace(scanner.Text())

		if value == "" {
			fmt.Printf("Operation canceled\n")
			return nil, true
		}

		inputs[ph] = value
		fmt.Printf("✓ %s: %s\n", ph, value)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Scanner error: %v\n", err)
		return nil, true
	}

	return inputs, false
}
