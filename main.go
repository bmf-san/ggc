package main

import (
	"fmt"
	"os"

	"github.com/bmf-san/gcl/cmd"
	"github.com/bmf-san/gcl/router"
)

const version = "v0.1.0-beta"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println("gcl version", version)
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
