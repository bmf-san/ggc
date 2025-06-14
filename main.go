package main

import (
	"os"

	"github.com/bmf-san/gcl/router"
)

func main() {
	router.Route(os.Args)
}
