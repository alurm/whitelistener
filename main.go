package main

import (
	_ "embed"
	"fmt"
	"os"
)

//go:embed README.md
var readme string

func main() {
	if err := whitelisten(); err != nil {
		fmt.Fprintf(os.Stderr, "An error: %v.\n\n%s", err, readme)
		os.Exit(1)
	}
}
