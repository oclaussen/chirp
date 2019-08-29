package main

import (
	"os"

	"github.com/oclaussen/chirp/pkg/command"
)

func main() {
	cmd := command.NewCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
