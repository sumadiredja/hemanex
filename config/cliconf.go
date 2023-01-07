package config

import (
	"fmt"

	helper "hemanex/helper"
	menu "hemanex/menu"

	"github.com/urfave/cli"
)

func CliConfig(app *cli.App) *cli.App {
	app.Name = "hemanex"
	app.Usage = "Manage Docker Private Registry on Nexus"
	app.Version = "1.0.0-alpha"
	app.Authors = []cli.Author{
		{
			Name:  "Robby Hemawan P <pramuditorh>",
			Email: "pramuditorh@gmail.com",
		},
		{
			Name:  "Hamdani Fadhli <ArleyB>",
			Email: "hamdanifadhli@gmail.com",
		},
		{
			Name:  "Vaghan Muhammad Sumadiredja <vaghansumadiredja>",
			Email: "vaghansumadiredja@gmail.com",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "login",
			Aliases: []string{"lgn"},
			Usage:   "Login to nexus repository",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "nexus-host, host",
					Usage: "Nexus hostname",
				},
				cli.StringFlag{
					Name:  "nexus-host-port, hp",
					Usage: "Nexus host port",
				},
				cli.StringFlag{
					Name:  "repository-name, r",
					Usage: "Nexus repository name",
				},
				cli.StringFlag{
					Name:  "username, u",
					Usage: "Nexus credentials username",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "Nexus credentials password",
				},
				cli.StringFlag{
					Name:  "namespace, n",
					Usage: "Nexus namespace",
				},
				cli.BoolFlag{
					Name:  "ignore-confirmation, i",
					Usage: "Ignore confirmation prompt",
				},
				cli.StringFlag{
					Name:  "repository-port, rp",
					Usage: "Nexus repository port",
				},
				cli.BoolFlag{
					Name:  "insecure-registry, k",
					Usage: "Turn on insecure registries \033[0;31mjust for podman\033[0m",
				},
			},
			Action: func(c *cli.Context) error {
				return menu.SetNexusCredentials(c)
			},
		},
		{
			Name:    "repository",
			Aliases: []string{"rp"},
			Usage:   "Change nexus repository config",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "repository-name, r",
					Usage: "Nexus repository name",
				},
			},
			Action: func(c *cli.Context) error {
				return menu.GetRepository(c)
			},
		},
		{
			Name:    "namespace",
			Aliases: []string{"ns"},
			Usage:   "Change build namespace",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "namespace-name, n",
					Usage: "Nexus namespace",
				},
			},
			Action: func(c *cli.Context) error {
				return menu.GetNamespace(c)
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
			Name:    "build",
			Usage:   "Building image in nexus repository",
			Aliases: []string{"b"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "port, p",
					Usage: "Image port",
				},
				cli.StringFlag{
					Name:  "tags, t",
					Usage: "Image tags <image_name>:<tag>",
				},
				cli.StringFlag{
					Name:  "namespace, n",
					Usage: "Overwrite namespace",
				},
			},
			Action: func(c *cli.Context) error {
				return menu.BuildImage(c)
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
		{
			Name:    "image-local",
			Aliases: []string{"img-local"},
			Usage:   "Manage Docker Images",
			Subcommands: []cli.Command{
				{
					Name:    "delete",
					Aliases: []string{"del"},
					Usage:   "Delete an image usage: hemanex image-local delete <image_name>:<tag>",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "force, f",
							Usage: "Force delete container image",
						},
						cli.StringFlag{
							Name:  "port, p",
							Usage: "Nexus repository port",
						},
						cli.BoolFlag{
							Name:  "all, a",
							Usage: "Delete all matched images",
						},
					},
					Action: func(c *cli.Context) error {
						return menu.DeleteImageLocal(c)
					},
				},
			},
		},
		{
			Name:    "push",
			Usage:   "Push an image to Nexus Repository",
			Aliases: []string{"pu"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "port, p",
				},
				cli.BoolFlag{
					Name:  "insecure-registry, k",
					Usage: "Turn on insecure registries",
				},
				cli.StringFlag{
					Name:  "namespace, n",
					Usage: "Overwrite namespace",
				},
			},
			Action: func(c *cli.Context) error {
				return menu.PushImage(c)
			},
		},
	}
	app.CommandNotFound = func(c *cli.Context, command string) {
		// fmt.Fprintf(c.App.Writer, "Wrong command %q!\n", command)

		fmt.Println(helper.CliCriticalGen(fmt.Errorf("command %q not found", command), 1))
	}

	return app
}
