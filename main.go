package main

import (
	"fmt"
	"os"

	"github.com/bmf-san/ggc/cmd"
	"github.com/bmf-san/ggc/router"
)

const version = "v1.0.2"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println("ggc version", version)
		return
	}
	if len(os.Args) == 1 {
		args := cmd.InteractiveUI()
		if args != nil {
			router.Route(args)
		}
		return
	}
	router.Route(os.Args)
}
