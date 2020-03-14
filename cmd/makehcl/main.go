package main

import (
	"fmt"
	"os"

	"github.com/beauknowssoftware/makehcl/internal/cmd"
)

func main() {
	if err := cmd.Exec(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
