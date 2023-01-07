package main

import (
	"os"

	config "hemanex/config"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app = config.CliConfig(app)
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
