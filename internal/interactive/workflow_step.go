package interactive

import (
	"fmt"
	"strings"
)

// WorkflowStep represents a single step in a workflow
type WorkflowStep struct {
	ID          int      `json:"id"`
	Command     string   `json:"command"`
	Args        []string `json:"args"`
	Description string   `json:"description"`
}

// String returns a string representation of the workflow step
func (ws *WorkflowStep) String() string {
	if ws.Description != "" {
		return fmt.Sprintf("[%d] %s", ws.ID, ws.Description)
	}

	cmdStr := ws.Command
	if len(ws.Args) > 0 {
		cmdStr += " " + strings.Join(ws.Args, " ")
	}
	return fmt.Sprintf("[%d] %s", ws.ID, cmdStr)
}
