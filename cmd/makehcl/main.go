package main

import (
	"os"

	"github.com/beauknowssoftware/makehcl/internal/cmd"
)

func main() {
	if err := cmd.Exec(); err != nil {
		os.Exit(1)
	}
}
