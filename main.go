package main

import (
	"os"

	config "hemanex/config"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app = config.CliConfig(app)
	app.Run(os.Args)
}
