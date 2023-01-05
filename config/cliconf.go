package config

import (
	"fmt"

	menu "hemanex/menu"

	"github.com/urfave/cli"
)

func CliConfig(app *cli.App) *cli.App {
	app.Name = "hemanex"
	app.Usage = "Manage Docker Private Registry on Nexus"
	app.Version = "1.0.0-beta"
	app.Authors = []cli.Author{
		{
			Name:  "Robby Hemawan P",
			Email: "robby.pramudito@btpn.com",
		},
		{
			Name:  "Hamdani Fadhli",
			Email: "hamdani.fadhli@btpn.com",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "configure",
			Aliases: []string{"c"},
			Usage:   "Configure Nexus Credentials",
			Action: func(c *cli.Context) error {
				return menu.SetNexusCredentials(c)
			},
		},
		{
			Name:    "test",
			Usage:   "testing root",
			Aliases: []string{"t"},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "insecure-registry, k",
					Usage: "Turn on insecure registries",
				},
			},
			Action: func(c *cli.Context) error {
				return menu.CheckToml(c)
			},
		},
		{
			Name:    "image",
			Aliases: []string{"img"},
			Usage:   "Manage Docker Images",
			Subcommands: []cli.Command{
				{
					Name:    "list",
					Aliases: []string{"ls"},
					Usage:   "List all images in repository",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "insecure-registry, k",
							Usage: "Turn on insecure registries",
						},
					},
					Action: func(c *cli.Context) error {
						return menu.ListImages(c)
					},
				},
				{
					Name:    "tags",
					Usage:   "Display all image tags",
					Aliases: []string{"tg"},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "List tags by image name",
						},
						cli.BoolFlag{
							Name:  "insecure-registry, k",
							Usage: "Turn on insecure registries",
						},
					},
					Action: func(c *cli.Context) error {
						return menu.ListTagsByImage(c)
					},
				},
				{
					Name:    "info",
					Usage:   "Show image details",
					Aliases: []string{"if"},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "name, n",
						},
						cli.StringFlag{
							Name: "tag, t",
						},
						cli.BoolFlag{
							Name:  "insecure-registry, k",
							Usage: "Turn on insecure registries",
						},
					},
					Action: func(c *cli.Context) error {
						return menu.ShowImageInfo(c)
					},
				},
				{
					Name:    "delete",
					Aliases: []string{"del"},
					Usage:   "Delete an image",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "name, n",
						},
						cli.StringFlag{
							Name: "tag, t",
						},
						cli.StringFlag{
							Name: "keep, kp",
						},
						cli.BoolFlag{
							Name:  "insecure-registry, k",
							Usage: "Turn on insecure registries",
						},
					},
					Action: func(c *cli.Context) error {
						return menu.DeleteImage(c)
					},
				},
				{
					Name:    "size",
					Aliases: []string{"sz"},
					Usage:   "Show total size of image including all tags",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "name, n",
						},
						cli.BoolFlag{
							Name:  "insecure-registry, k",
							Usage: "Turn on insecure registries",
						},
					},
					Action: func(c *cli.Context) error {
						return menu.ShowTotalImageSize(c)
					},
				},
			},
		},
	}
	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Fprintf(c.App.Writer, "Wrong command %q !", command)
	}

	return app
}
