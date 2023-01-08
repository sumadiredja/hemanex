package menu

import (
	"fmt"
	"strings"

	"github.com/sumadiredja/hemanex/registry"

	helper "github.com/sumadiredja/hemanex/helper"

	"github.com/urfave/cli"
)

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
		return helper.ShowSubCommand("please provide image name", c)
	}

	imgName = helper.CheckFlagsStringExist(c.String("namespace"), r.Namespace) + "/" + imgName

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
		return helper.ShowSubCommand("please provide image name and tag", c)
	}

	imgName = helper.CheckFlagsStringExist(c.String("namespace"), r.Namespace) + "/" + imgName

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
	var images = c.Args()
	var isDifferentNamespace = c.Bool("namespace")

	if len(images) == 0 {
		return helper.ShowSubCommand("please provide image name", c)
	} else {
		r, err := registry.NewRegistry(c)
		if err != nil {
			return helper.CliErrorGen(err, 1)
		}
		if !isDifferentNamespace {
			for _, image := range images {
				imageName := r.Namespace + "/" + strings.Split(image, ":")[0]
				imageTag := strings.Split(image, ":")[1]
				err = r.DeleteImageByTag(imageName, imageTag)

				if err != nil {
					return helper.CliErrorGen(err, 1)
				}
			}
		} else {
			for _, image := range images {
				imageName := strings.Split(image, ":")[0]
				imageTag := strings.Split(image, ":")[1]
				err = r.DeleteImageByTag(imageName, imageTag)

				if err != nil {
					return helper.CliErrorGen(err, 1)
				}
			}
		}
	}
	return nil
}

func DeleteImageKeep(c *cli.Context) error {
	var imgName = c.Args().Get(0)
	var remains = c.Int("remains")
	var keepTag = c.String("keep-tag")

	if imgName == "" {
		return helper.ShowSubCommand("please provide image name", c)
	} else {
		r, err := registry.NewRegistry(c)
		if err != nil {
			return helper.CliErrorGen(err, 1)
		}
		var imgName = helper.CheckFlagsStringExist(c.String("namespace"), r.Namespace) + "/" + imgName
		if keepTag == "" {
			if remains == 0 {
				return helper.ShowSubCommand("please provide image tag or how many images you want to keep", c)
			} else {
				tags, err := r.ListTagsByImage(imgName)
				compareStringNumber := func(str1, str2 string) bool {
					return helper.ExtractNumberFromString(str1) < helper.ExtractNumberFromString(str2)
				}
				helper.Compare(compareStringNumber).Sort(tags)
				if err != nil {
					return helper.CliErrorGen(err, 1)
				}
				if len(tags) >= remains {
					for _, tag := range tags[:len(tags)-1] {
						fmt.Printf("%s:%s image will be deleted ...\n", imgName, tag)
						err = r.DeleteImageByTag(imgName, tag)
						if err != nil {
							return helper.CliErrorGen(err, 1)
						}
					}
				} else {
					fmt.Printf("Only %d images are available\n", len(tags))
				}
			}
		} else {
			var keepIndex int
			tags, err := r.ListTagsByImage(imgName)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			for i, tag := range tags {
				if tag == keepTag {
					keepIndex = i
				}
			}
			temp := tags[keepIndex]
			tags[keepIndex] = tags[len(tags)-1]
			tags[len(tags)-1] = temp
			deletedTags := tags[:len(tags)-1]
			for _, tag := range deletedTags {
				err = r.DeleteImageByTag(imgName, tag)
				if err != nil {
					return helper.CliErrorGen(err, 1)
				}
			}
		}
	}
	return nil
}

func ShowTotalImageSize(c *cli.Context) error {
	var imgName = c.String("name")
	var totalSize (int64) = 0

	if imgName == "" {
		return helper.ShowSubCommand("please provide image name", c)
	} else {
		r, err := registry.NewRegistry(c)
		if err != nil {
			return helper.CliErrorGen(err, 1)
		}
		imgName = helper.CheckFlagsStringExist(c.String("namespace"), r.Namespace) + "/" + imgName

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
