package main

import (
	"github.com/vite-cloud/vite/core/handler/cli"
	"os"
)

func main() {
	os.Exit(cli.New().Run(os.Args[1:]))
}
