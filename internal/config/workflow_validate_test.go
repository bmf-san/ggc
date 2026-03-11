package config

import (
	"strings"
	"testing"
)

func TestValidateWorkflows_Valid(t *testing.T) {
	tests := []struct {
		name      string
		workflows map[string][]string
	}{
		{
			name:      "nil workflows",
			workflows: nil,
		},
		{
			name:      "empty workflows map",
			workflows: map[string][]string{},
		},
		{
			name: "simple workflow",
			workflows: map[string][]string{
				"acp": {"add .", "commit", "push current"},
			},
		},
		{
			name: "workflow with angle bracket placeholder",
			workflows: map[string][]string{
				"deploy": {"add .", "commit <message>", "push current"},
			},
		},
		{
			name: "workflow with curly brace placeholder",
			workflows: map[string][]string{
				"ci": {"commit {0}"},
			},
		},
		{
			name: "multiple workflows",
			workflows: map[string][]string{
				"acp":    {"add .", "commit <message>", "push current"},
				"rebase": {"fetch origin", "rebase origin main"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{Workflows: tt.workflows}
			if err := c.validateWorkflows(); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateWorkflows_Invalid(t *testing.T) {
	tests := []struct {
		name      string
		workflows map[string][]string
		wantMsg   string
	}{
		{
			name: "name with space",
			workflows: map[string][]string{
				"bad name": {"add ."},
			},
			wantMsg: "must not be empty or contain spaces",
		},
		{
			name: "empty steps slice",
			workflows: map[string][]string{
				"empty": {},
			},
			wantMsg: "at least one step",
		},
		{
			name: "blank step string",
			workflows: map[string][]string{
				"ws": {"   "},
			},
			wantMsg: "must not be empty",
		},
		{
			name: "step with shell injection via semicolon",
			workflows: map[string][]string{
				"bad": {"commit; rm -rf /"},
			},
			wantMsg: "unsafe shell metacharacters",
		},
		{
			name: "step with pipe",
			workflows: map[string][]string{
				"pipe": {"status | grep modified"},
			},
			wantMsg: "unsafe shell metacharacters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{Workflows: tt.workflows}
			err := c.validateWorkflows()
			if err == nil {
				t.Error("expected error, got nil")
				return
			}
			if !strings.Contains(err.Error(), tt.wantMsg) {
				t.Errorf("error %q did not contain %q", err.Error(), tt.wantMsg)
			}
		})
	}
}
