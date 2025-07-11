// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// CommandInfo contiene información sobre un comando
type CommandInfo struct {
	Command     string
	Description string
}

var commands = []string{
	"add <file>",
	"add .",
	"add -p",
	"branch current",
	"branch checkout",
	"branch checkout-remote",
	"branch create",
	"branch delete",
	"branch delete-merged",
	"push current",
	"push force",
	"pull current",
	"pull rebase",
	"log simple",
	"log graph",
	"commit <message>",
	"commit allow-empty",
	"commit tmp",
	"commit amend <message>",
	"fetch --prune",
	"tag list",
	"tag annotated <tag> <message>",
	"tag delete <tag>",
	"tag show <tag>",
	"tag push",
	"tag create <tag>",
	"diff",
	"diff unstaged",
	"diff staged",
	"version",
	"clean files",
	"clean dirs",
	"clean-interactive",
	"reset-clean",
	"commit-push-interactive",
	"stash trash",
	"status",
	"status short",
	"rebase interactive",
	"remote list",
	"remote add <n> <url>",
	"remote remove <n>",
	"remote set-url <n> <url>",
	"add-commit-push",
	"pull-rebase-push",
	"stash-pull-pop",
	"quit",
}

// commandDescriptions mapea cada comando con su descripción
var commandDescriptions = map[string]string{
	"add <file>":                    "Añadir archivo específico al índice",
	"add .":                         "Añadir todos los cambios al índice",
	"add -p":                        "Añadir cambios interactivamente",
	"branch current":                "Mostrar nombre de la rama actual",
	"branch checkout":               "Cambiar a rama existente",
	"branch checkout-remote":        "Crear y cambiar a rama local desde remota",
	"branch create":                 "Crear y cambiar a nueva rama",
	"branch delete":                 "Eliminar rama local",
	"branch delete-merged":          "Eliminar ramas locales fusionadas",
	"push current":                  "Enviar rama actual al remoto",
	"push force":                    "Forzar envío de rama actual",
	"pull current":                  "Obtener rama actual del remoto",
	"pull rebase":                   "Obtener con rebase",
	"log simple":                    "Mostrar historial simple",
	"log graph":                     "Mostrar historial con gráfico",
	"commit <message>":              "Crear commit con mensaje",
	"commit allow-empty":            "Crear commit vacío",
	"commit tmp":                    "Crear commit temporal",
	"commit amend <message>":        "Modificar commit anterior",
	"fetch --prune":                 "Obtener y limpiar referencias obsoletas",
	"tag list":                      "Listar todas las etiquetas",
	"tag annotated <tag> <message>": "Crear etiqueta anotada",
	"tag delete <tag>":              "Eliminar etiqueta",
	"tag show <tag>":                "Mostrar información de etiqueta",
	"tag push":                      "Enviar etiquetas al remoto",
	"tag create <tag>":              "Crear etiqueta",
	"diff":                          "Mostrar diferencias",
	"diff unstaged":                 "Mostrar cambios no preparados",
	"diff staged":                   "Mostrar cambios preparados",
	"version":                       "Mostrar versión actual",
	"clean files":                   "Limpiar archivos no rastreados",
	"clean dirs":                    "Limpiar directorios no rastreados",
	"clean-interactive":             "Limpiar archivos interactivamente",
	"reset-clean":                   "Resetear y limpiar",
	"commit-push-interactive":       "Commit y push interactivo",
	"stash trash":                   "Eliminar stash",
	"status":                        "Mostrar estado del árbol de trabajo",
	"status short":                  "Mostrar estado conciso",
	"rebase interactive":            "Rebase interactivo",
	"remote list":                   "Listar repositorios remotos",
	"remote add <n> <url>":          "Añadir repositorio remoto",
	"remote remove <n>":             "Eliminar repositorio remoto",
	"remote set-url <n> <url>":      "Cambiar URL del repositorio remoto",
	"add-commit-push":               "Añadir, commit y push en una operación",
	"pull-rebase-push":              "Pull con rebase y push",
	"stash-pull-pop":                "Stash, pull y pop en secuencia",
	"quit":                          "Salir del modo interactivo",
}

// InteractiveUI provides an incremental search interactive UI for command selection.
// Returns the selected command as []string (nil if nothing selected)
func InteractiveUI() []string {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Println("Failed to set terminal to raw mode:", err)
		return nil
	}
	defer func() {
		if err := term.Restore(fd, oldState); err != nil {
			fmt.Fprintln(os.Stderr, "failed to restore terminal state:", err)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	selected := 0
	input := ""

	for {
		if _, err := os.Stdout.Write([]byte("\033[H\033[2J\033[H")); err != nil {
			fmt.Fprintln(os.Stderr, "failed to write clear screen sequence:", err)
		}
		fmt.Printf("Select a command (incremental search: type to filter, ctrl+n: down, ctrl+p: up, enter: execute, ctrl+c: quit)\n")
		fmt.Printf("\rSearch: %s\n\n", input)

		// Filtering
		filtered := []string{}
		for _, cmd := range commands {
			if strings.Contains(cmd, input) {
				filtered = append(filtered, cmd)
			}
		}
		if input == "" {
			fmt.Println("(Type to filter commands...)")
		} else {
			if len(filtered) == 0 {
				fmt.Println("  (No matching command)")
			}
			if selected >= len(filtered) {
				selected = len(filtered) - 1
			}
			if selected < 0 {
				selected = 0
			}

			// Encontrar el comando más largo para alineación
			maxCmdLen := 0
			for _, cmd := range filtered {
				if len(cmd) > maxCmdLen {
					maxCmdLen = len(cmd)
				}
			}

			for i, cmd := range filtered {
				description := commandDescriptions[cmd]
				if description == "" {
					description = "Sin descripción"
				}

				// Formatear con alineación
				padding := strings.Repeat(" ", maxCmdLen-len(cmd))
				if i == selected {
					fmt.Printf("\r> %s%s  %s\n", cmd, padding, description)
				} else {
					fmt.Printf("\r  %s%s  %s\n", cmd, padding, description)
				}
			}
		}
		fmt.Print("\n\r") // Ensure next output starts at left edge

		b, err := reader.ReadByte()
		if err != nil {
			continue
		}
		if b == 3 { // Ctrl+C in raw mode
			if err := term.Restore(fd, oldState); err != nil {
				fmt.Fprintln(os.Stderr, "failed to restore terminal state:", err)
			}
			fmt.Println("\nExiting...")
			os.Exit(0)
		} else if b == 13 { // Enter
			if len(filtered) > 0 {
				fmt.Printf("\nExecute: %s\n", filtered[selected])
				if err := term.Restore(fd, oldState); err != nil {
					fmt.Fprintln(os.Stderr, "failed to restore terminal state:", err)
				}
				// Placeholder detection
				cmdTemplate := filtered[selected]
				placeholders := extractPlaceholders(cmdTemplate)
				inputs := make(map[string]string)
				readerStdin := bufio.NewReader(os.Stdin)
				for _, ph := range placeholders {
					fmt.Print("\n\r") // Newline + carriage return
					fmt.Printf("Enter value for %s: ", ph)
					val, _ := readerStdin.ReadString('\n')
					val = strings.TrimSpace(val)
					inputs[ph] = val
				}
				// Placeholder replacement
				finalCmd := cmdTemplate
				for ph, val := range inputs {
					finalCmd = strings.ReplaceAll(finalCmd, "<"+ph+">", val)
				}
				args := []string{"ggc"}
				args = append(args, strings.Fields(finalCmd)...)
				return args
			}
			break
		} else if b == 16 { // Ctrl+p
			if selected > 0 {
				selected--
			}
		} else if b == 14 { // Ctrl+n
			if selected < len(filtered)-1 {
				selected++
			}
		} else if b == 127 || b == 8 { // Backspace
			if len(input) > 0 {
				input = input[:len(input)-1]
			}
		} else if b >= 32 && b <= 126 { // Printable ASCII
			input += string(b)
		}
	}
	return nil
}

// Extract <...> placeholders from a string
func extractPlaceholders(s string) []string {
	var res []string
	start := -1
	for i, c := range s {
		if c == '<' {
			start = i + 1
		} else if c == '>' && start != -1 {
			res = append(res, s[start:i])
			start = -1
		}
	}
	return res
}
