package cmd

import (
	"reflect"
	"testing"
)

func TestExtractPlaceholders(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "no placeholders",
			input: "simple command",
			want:  []string{},
		},
		{
			name:  "single placeholder",
			input: "add <file>",
			want:  []string{"file"},
		},
		{
			name:  "multiple placeholders",
			input: "remote add <n> <url>",
			want:  []string{"n", "url"},
		},
		{
			name:  "empty placeholder",
			input: "command <>",
			want:  []string{""},
		},
		{
			name:  "invalid format",
			input: "command <incomplete",
			want:  []string{},
		},
		{
			name:  "nested placeholder",
			input: "command <<nested>>",
			want:  []string{"nested"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractPlaceholders(tt.input)
			// Compare without distinguishing between nil and empty slice
			if len(tt.want) == 0 && len(got) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractPlaceholders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandDescriptions(t *testing.T) {
	// Verificar que todos los comandos tienen descripciones
	for _, cmd := range commandList {
		description := cmd.Description
		if description == "" {
			t.Errorf("Command '%s' has no description", cmd.Command)
		}
	}

	// Verificar que no hay descripciones para comandos que no existen
	// (ya no aplica porque solo hay una lista)
}

func TestCommandDescriptionsContent(t *testing.T) {
	// Verificar que algunas descripciones clave están presentes
	expectedDescriptions := map[string]string{
		"add <file>":       "Add specific file to index",
		"status":           "Show working tree status",
		"commit <message>": "Create commit with message",
		"quit":             "Exit interactive mode",
	}

	for cmdStr, expectedDesc := range expectedDescriptions {
		found := false
		for _, cmd := range commandList {
			if cmd.Command == cmdStr {
				if cmd.Description != expectedDesc {
					t.Errorf("Description for '%s' is '%s', expected '%s'", cmdStr, cmd.Description, expectedDesc)
				}
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Command '%s' not found in commandList", cmdStr)
		}
	}
}
