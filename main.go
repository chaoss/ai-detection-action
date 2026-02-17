package main

import (
	"os"

	"github.com/chaoss/ai-detection-action/cmd"
)

func main() {
	os.Exit(cmd.Run(os.Args[1:], os.Stdout, os.Stderr))
}
