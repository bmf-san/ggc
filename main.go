package main

import (
	"os"

	"github.com/bmf-san/ggc/cmd"
	"github.com/bmf-san/ggc/router"
)

func main() {
	c := cmd.NewCmd(os.Stdout)
	r := router.NewRouter(c)
	r.Route(os.Args)
}
