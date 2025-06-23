// Package main is the entry point for the ggc CLI tool.
package main

import (
	"os"

	"github.com/bmf-san/ggc/cmd"
	"github.com/bmf-san/ggc/router"
)

func main() {
	c := cmd.NewCmd()
	r := router.NewRouter(c)
	r.Route(os.Args[1:])
}
