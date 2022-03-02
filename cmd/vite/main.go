package main

import (
	"encoding/json"
	"fmt"
	"github.com/vite-cloud/vite/core/domain/metrics"
	"os"
)

func main() {
	m, err := metrics.Gather()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	s, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(string(s))

	//os.Exit(cli.New().Run(os.Args[1:]))
}
