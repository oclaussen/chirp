package main

import (
	"os"

	"github.com/oclaussen/chirp/pkg/command"
)

func main() {
	os.Exit(command.Execute())
}
