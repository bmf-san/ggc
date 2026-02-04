// Command-line tool that generates shell completion scripts from the command registry.
package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/bmf-san/ggc/v7/cmd/command"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

func main() {
	data := buildTemplateData(command.DefaultRegistry.VisibleCommands())

	templates := map[string]string{
		"bash": "templates/bash.tmpl",
		"zsh":  "templates/zsh.tmpl",
		"fish": "templates/fish.tmpl",
	}

	writers := map[string]string{
		"bash": filepath.Join("tools", "completions", "ggc.bash"),
		"zsh":  filepath.Join("tools", "completions", "ggc.zsh"),
		"fish": filepath.Join("tools", "completions", "ggc.fish"),
	}

	funcMap := template.FuncMap{
		"join":         join,
		"escapeZsh":    escapeZsh,
		"escapeFish":   escapeFish,
		"escapeBash":   escapeBash,
		"hasKeywords":  hasKeywords,
		"needsHandler": needsHandler,
		"subcommandBy": subcommandBy,
	}

	for name, tmplPath := range templates {
		if err := renderTemplate(funcMap, tmplPath, writers[name], data); err != nil {
			fmt.Fprintf(os.Stderr, "error generating %s completions: %v\n", name, err)
			os.Exit(1)
		}
	}

	fmt.Println("Shell completions regenerated successfully")
}

func renderTemplate(funcMap template.FuncMap, templatePath, dest string, data *TemplateData) error {
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(funcMap).ParseFS(templateFS, templatePath)
	if err != nil {
		return err
	}

	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	if err := tmpl.Execute(file, data); err != nil {
		return err
	}

	return nil
}
