package main

import (
	"os"

	"github.com/beauknowssoftware/makehcl/internal/cmd"
)

const (
	failureStatusCode = 1
)

func main() {
	if err := cmd.Exec(); err != nil {
		os.Exit(failureStatusCode)
	}
}
