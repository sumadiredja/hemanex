package menu

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	helper "github.com/sumadiredja/hemanex/helper"
	"github.com/sumadiredja/hemanex/registry"
	"github.com/urfave/cli"
)

func BuildImage(c *cli.Context) error {
	var image_name, tag string
	cwd := c.Args().Get(0)
	var port = c.String("repository-port")
	tags := c.String("tags")
	r, err := registry.NewRegistry(c)
	if err != nil {
		return helper.CliErrorGen(err, 1)
	}
	var namespace = r.Namespace

	if c.String("namespace") != "" {
		namespace = c.String("namespace")
	}
	if tags == "" {
		return helper.ShowSubCommand("please provide image tags", c)
	}
	if cwd == "" {
		return helper.ShowSubCommand("please provide dockerfile path in the arguments", c)
	}

	tag_split := strings.Split(tags, ":")
	if len(tag_split) <= 1 {
		return helper.ShowSubCommand("please provide image tags", c)
	}
	image_name = tag_split[0]
	tag = tag_split[1]
	if tag == "" {
		return helper.ShowSubCommand("please provide correct image tags", c)
	}
	if port == "" {
		port = r.RepositoryPort
	}
	host := strings.Split(r.Host, "://")[1]
	command := fmt.Sprintf("docker build -t %s:%s/%s/%s:%s %s", host, port, namespace, image_name, tag, cwd)
	_, err = helper.RunSubProcess(command, "failed to build image")
	if err != nil {
		return err
	}

	helper.CliSuccessVerbose("Successfully built image " + image_name + ":" + tag + " with namespace " + namespace)

	return nil
}

func PushImage(c *cli.Context) error {
	var skipTls string
	var imgName = c.Args().Get(0)

	r, err := registry.NewRegistry(c)
	if err != nil {
		return helper.CliErrorGen(err, 1)
	}

	if imgName == "" {
		return helper.ShowSubCommand("please provide image name", c)
	}

	repository_port := helper.CheckFlagsStringExist(c.String("repository-port"), r.RepositoryPort)
	namespace := helper.CheckFlagsStringExist(c.String("namespace"), r.Namespace)

	var isInsecure = c.Bool("insecure-registry")

	if isInsecure {
		skipTls = " --tls-verify=false"
	}

	cmdLogin := fmt.Sprintf("docker login " + strings.Split(r.Host, "//")[1] + ":" + repository_port + " -u " + r.Username + " -p " + r.Password + skipTls)
	_, err = helper.RunSubProcess(cmdLogin, helper.ErrorCertificateResponse)
	if err != nil {
		return err
	}

	cmdPushImage := fmt.Sprintf("docker push " + strings.Split(r.Host, "//")[1] + ":" + repository_port + "/" + namespace + "/" + imgName + skipTls)
	_, err = helper.RunSubProcess(cmdPushImage, "image not found locally")
	if err != nil {
		return err
	}

	helper.CliSuccessVerbose("Successfully pushed image " + imgName + " to " + r.Host + " namespace " + namespace)

	return nil
}

func DeleteImageLocal(c *cli.Context) error {
	var err error
	var prefix string

	r, err := registry.NewRegistry(c)
	if err != nil {
		return helper.CliErrorGen(err, 1)
	}

	repository_port := helper.CheckFlagsStringExist(c.String("repository-port"), r.RepositoryPort)

	image_name := c.Args().Get(0)
	if image_name == "" {
		return helper.ShowSubCommand("please provide image name", c)
	}

	name_split := strings.Split(image_name, ":")
	if len(name_split) <= 1 {
		return helper.ShowSubCommand("please provide image tags", c)
	}
	if name_split[1] == "" {
		return helper.ShowSubCommand("please provide correct image tags", c)
	}

	var force = ""
	if c.Bool("force") {
		force = "-f "
	}

	prefix = strings.Split(r.Host, "://")[1] + ":" + repository_port

	cmd := fmt.Sprintf(`docker images --format="{{.Repository}} {{.Tag}}" | grep "%s" | grep "%s/%s" | grep "%s"`, prefix, r.Namespace, name_split[0], name_split[1])
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return helper.CliErrorGen(errors.New("No image found with name "+image_name), 1)
	}

	list := strings.Split(strings.ReplaceAll(string(out), "\rn", "\n"), "\n")
	if list[len(list)-1] == "" {
		list = list[0 : len(list)-1]
	}

	if len(list) == 1 {
		image := strings.ReplaceAll(list[0], " ", ":")
		cmd = fmt.Sprintf("docker rmi %s%s", force, image)
		if _, err := exec.Command("bash", "-c", cmd).Output(); err != nil {
			return helper.CliErrorGen(errors.New("error: error getting input from user"), 1)
		}
		helper.CliSuccessVerbose("Successfully deleted image " + image)
		return nil
	}

	var list_images []string
	for i := 0; i < len(list)-1; i++ {
		hehe := strings.Split(list[i], " ")
		image := strings.Join(hehe, ":")
		list_images = append(list_images, image)
		if c.Bool("all") {
			cmd = fmt.Sprintf("docker rmi %s%s", force, image)
			_, err = exec.Command("bash", "-c", cmd).Output()
			if err != nil {
				return helper.CliErrorGen(errors.New("error: error getting input from user"), 1)
			}
			helper.CliSuccessVerbose("Successfully deleted image " + image)
		}

	}

	if c.Bool("all") {
		return nil
	}

	template := &promptui.SelectTemplates{
		Active:   "🚀 {{ . | cyan }}",
		Selected: "Image Selected: {{ . | yellow }}",
	}

	prompt := promptui.Select{
		Label:     "Select Image",
		Items:     list_images,
		Templates: template,
		Searcher: func(input string, index int) bool {
			return strings.Contains(list_images[index], input)
		},
	}

	_, selected_image, err := prompt.Run()
	if err != nil {
		return helper.CliErrorGen(errors.New("error: error getting input from user"), 1)
	}

	cmd = fmt.Sprintf("docker rmi %s%s", force, selected_image)
	_, err = exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return helper.CliErrorGen(errors.New("error: error image not found"), 1)
	}
	helper.CliSuccessVerbose("Successfully deleted image " + selected_image)

	return nil
}
