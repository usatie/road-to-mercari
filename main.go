package main

import (
	"os"

	"github.com/usatie/convert/cli"
)

func main() {
	app := &cli.App{OutStream: os.Stdout, ErrStream: os.Stderr}
	os.Exit(app.Run(os.Args))
}
