package main

import (
	"convert/cli"
	"os"
)

func main() {
	app := &cli.App{OutStream: os.Stdout, ErrStream: os.Stderr}
	os.Exit(app.Run(os.Args))
}
