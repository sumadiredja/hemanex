package menu

import (
	"fmt"
	"html/template"
	"os"

	"hemanex/registry"

	helper "hemanex/helper"

	b64 "encoding/base64"

	"github.com/urfave/cli"
)

const (
	CREDENTIALS_TEMPLATES = `# Nexus Credentials
nexus_host = "{{ .Host }}"
nexus_username = "{{ .Username }}"
nexus_password = "{{ .Password }}"
nexus_repository = "{{ .Repository }}"
`
)

func SetNexusCredentials(c *cli.Context) error {
	var hostname, repository, username, password string
	var CREDENTIALS_FILE string
	var err error
	var tmpl *template.Template
	var f *os.File

	// CREDENTIALS_FILE, err := helper.GetCredentialPath()
	if CREDENTIALS_FILE, err = helper.GetCredentialPath(); err != nil {
		return err
	}

	hostname = helper.GetInputOrFlags(c.String("nexus-host"))
	repository = helper.GetInputOrFlags(c.String("repository-name"))
	username = helper.GetInputOrFlags(c.String("username"))
	if password, err = helper.GetPassword(c.String("password")); err != nil {
		return err
	}

	data := struct {
		Host       string
		Username   string
		Password   string
		Repository string
	}{
		hostname,
		username,
		password,
		repository,
	}

	if tmpl, err = template.New(CREDENTIALS_FILE).Parse(CREDENTIALS_TEMPLATES); err != nil {
		return helper.CliErrorGen(err, 1)
	}

	if f, err = os.Create(CREDENTIALS_FILE); err != nil {
		return helper.CliErrorGen(err, 1)
	}

	if err = tmpl.Execute(f, data); err != nil {
		return helper.CliErrorGen(err, 1)
	}

	if err = os.Chmod(CREDENTIALS_FILE, 0666); err != nil {
		return helper.CliErrorGen(err, 1)
	}

	helper.CliSuccessVerbose("nexus user configured")
	return nil
}

func CheckToml(c *cli.Context) error {
	r, err := registry.NewRegistry(c)
	if err != nil {
		return helper.CliErrorGen(err, 1)

	}
	fmt.Println(r.Host, r.Password, r.Repository, r.Username)
	// fmt.Println(c.Bool("insecure-registry"))
	// fmt.Println(r.Host)
	return nil
}

func GetNamespace(c *cli.Context) error {
	var r registry.Registry
	var err error
	var CREDENTIALS_FILE string
	var tmpl *template.Template
	var f *os.File

	if r, err = registry.NewRegistry(c); err != nil {
		return helper.CliErrorGen(err, 1)
	}

	if CREDENTIALS_FILE, err = helper.GetCredentialPath(); err != nil {
		return err
	}

	if c.String("repository-name") != "" {
		data := struct {
			Host       string
			Username   string
			Password   string
			Repository string
		}{
			r.Host,
			r.Username,
			b64.StdEncoding.EncodeToString([]byte(r.Password)),
			c.String("repository-name"),
		}

		if tmpl, err = template.New(CREDENTIALS_FILE).Parse(CREDENTIALS_TEMPLATES); err != nil {
			return helper.CliErrorGen(err, 1)
		}

		if f, err = os.Create(CREDENTIALS_FILE); err != nil {
			return helper.CliErrorGen(err, 1)
		}

		if err = tmpl.Execute(f, data); err != nil {
			return helper.CliErrorGen(err, 1)
		}

		helper.CliSuccessVerbose(fmt.Sprintf("Repository changed to %s", c.String("repository-name")))
		return nil
	}
	helper.CliInfoVerbose(fmt.Sprintf("Currently working in %s repository", r.Repository))
	return nil
}

func ListImages(c *cli.Context) error {
	r, err := registry.NewRegistry(c)
	if err != nil {
		return helper.CliErrorGen(err, 1)
	}
	images, err := r.ListImages()
	if err != nil {
		return helper.CliErrorGen(err, 1)
	}
	for _, image := range images {
		fmt.Println(image)
	}
	helper.CliInfoVerbose(fmt.Sprintf("Total images: %d\n", len(images)))
	return nil
}

func ListTagsByImage(c *cli.Context) error {
	var imgName = c.String("name")
	r, err := registry.NewRegistry(c)
	if err != nil {
		return helper.CliErrorGen(err, 1)
	}
	if imgName == "" {
		cli.ShowSubcommandHelp(c)
	}
	tags, err := r.ListTagsByImage(imgName)

	compareStringNumber := func(str1, str2 string) bool {
		return helper.ExtractNumberFromString(str1) < helper.ExtractNumberFromString(str2)
	}
	helper.Compare(compareStringNumber).Sort(tags)

	if err != nil {
		return helper.CliErrorGen(err, 1)
	}
	for _, tag := range tags {
		fmt.Println(tag)
	}
	helper.CliInfoVerbose(fmt.Sprintf("There are %d images for %s\n", len(tags), imgName))
	return nil
}

func ShowImageInfo(c *cli.Context) error {
	var imgName = c.String("name")
	var tag = c.String("tag")
	r, err := registry.NewRegistry(c)
	if err != nil {
		return helper.CliErrorGen(err, 1)
	}
	if imgName == "" || tag == "" {
		cli.ShowSubcommandHelp(c)
	}
	manifest, err := r.ImageManifest(imgName, tag)
	if err != nil {
		return helper.CliErrorGen(err, 1)
	}
	fmt.Printf("Image: %s:%s\n", imgName, tag)
	fmt.Printf("Size: %d\n", manifest.Config.Size)
	fmt.Println("Layers:")
	for _, layer := range manifest.Layers {
		fmt.Printf("\t%s\t%d\n", layer.Digest, layer.Size)
	}
	return nil
}

func DeleteImage(c *cli.Context) error {
	var imgName = c.String("name")
	var tag = c.String("tag")
	var keep = c.Int("keep")
	if imgName == "" {
		fmt.Fprintf(c.App.Writer, "You should specify the image name\n")
		cli.ShowSubcommandHelp(c)
	} else {
		r, err := registry.NewRegistry(c)
		if err != nil {
			return helper.CliErrorGen(err, 1)
		}
		if tag == "" {
			if keep == 0 {
				fmt.Fprintf(c.App.Writer, "You should either specify the tag or how many images you want to keep\n")
				cli.ShowSubcommandHelp(c)
			} else {
				tags, err := r.ListTagsByImage(imgName)
				compareStringNumber := func(str1, str2 string) bool {
					return helper.ExtractNumberFromString(str1) < helper.ExtractNumberFromString(str2)
				}
				helper.Compare(compareStringNumber).Sort(tags)
				if err != nil {
					return helper.CliErrorGen(err, 1)
				}
				if len(tags) >= keep {
					for _, tag := range tags[:len(tags)-keep] {
						fmt.Printf("%s:%s image will be deleted ...\n", imgName, tag)
						r.DeleteImageByTag(imgName, tag)
					}
				} else {
					fmt.Printf("Only %d images are available\n", len(tags))
				}
			}
		} else {
			err = r.DeleteImageByTag(imgName, tag)
			if err != nil {
				return helper.CliErrorGen(err, 1)
			}
		}
	}
	return nil
}

func ShowTotalImageSize(c *cli.Context) error {
	var imgName = c.String("name")
	var totalSize (int64) = 0

	if imgName == "" {
		cli.ShowSubcommandHelp(c)
	} else {
		r, err := registry.NewRegistry(c)
		if err != nil {
			return helper.CliErrorGen(err, 1)
		}

		tags, err := r.ListTagsByImage(imgName)
		if err != nil {
			return helper.CliErrorGen(err, 1)
		}

		for _, tag := range tags {
			manifest, err := r.ImageManifest(imgName, tag)
			if err != nil {
				return helper.CliErrorGen(err, 1)
			}

			sizeInfo := make(map[string]int64)

			for _, layer := range manifest.Layers {
				sizeInfo[layer.Digest] = layer.Size
			}

			for _, size := range sizeInfo {
				totalSize += size
			}
		}
		fmt.Printf("%d %s\n", totalSize, imgName)
	}
	return nil
}
