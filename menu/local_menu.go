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

func ImageRetag(c *cli.Context) error {
	var err error
	var source_prefix, target_prefix string
	r, err := registry.NewRegistry(c)
	if err != nil {
		return helper.CliErrorGen(err, 1)
	}

	repository_port_source := helper.CheckFlagsStringExist(c.String("port-source"), r.RepositoryPort)
	repository_port_target := helper.CheckFlagsStringExist(c.String("port-target"), r.RepositoryPort)
	namespace_source := helper.CheckFlagsStringExist(c.String("namespace-source"), r.Namespace)
	namespace_target := helper.CheckFlagsStringExist(c.String("namespace-targer"), r.Namespace)

	source_image := c.Args().Get(0)
	if source_image == "" {
		return helper.ShowSubCommand("please provide source image name", c)
	}

	target_image := c.Args().Get(1)
	if target_image == "" {
		return helper.ShowSubCommand("please provide target image name", c)
	}

	source_split := strings.Split(source_image, ":")
	if len(source_split) <= 1 {
		return helper.ShowSubCommand("please provide image tags", c)
	}
	if source_split[1] == "" {
		return helper.ShowSubCommand("please provide correct image tags", c)
	}

	target_split := strings.Split(target_image, ":")
	if len(target_split) <= 1 {
		return helper.ShowSubCommand("please provide image tags", c)
	}
	if target_split[1] == "" {
		return helper.ShowSubCommand("please provide correct image tags", c)
	}

	source_prefix = strings.Split(r.Host, "://")[1] + ":" + repository_port_source
	target_prefix = strings.Split(r.Host, "://")[1] + ":" + repository_port_target

	cmd := fmt.Sprintf(`docker images --format="{{.Repository}} {{.Tag}}" | grep "%s" | grep "%s/%s" | grep "%s"`, source_prefix, namespace_source, source_split[0], source_split[1])
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return helper.CliErrorGen(errors.New("No image found with name "+source_image), 1)
	}

	list := strings.Split(strings.ReplaceAll(string(out), "\rn", "\n"), "\n")
	if list[len(list)-1] == "" {
		list = list[0 : len(list)-1]
	}

	if len(list) == 1 {
		image := strings.ReplaceAll(list[0], " ", ":")
		cmd = fmt.Sprintf("docker tag %s %s/%s/%s", image, target_prefix, namespace_target, target_image)
		if _, err := exec.Command("bash", "-c", cmd).Output(); err != nil {
			return helper.CliErrorGen(errors.New("error: error getting input from user"), 1)
		}
		helper.CliSuccessVerbose("Successfully created tag " + target_image + " from image " + source_image)
		return nil
	}

	var list_images []string
	for i := 0; i < len(list)-1; i++ {
		hehe := strings.Split(list[i], " ")
		image := strings.Join(hehe, ":")
		list_images = append(list_images, image)
	}

	template := &promptui.SelectTemplates{
		Active:   "🚀 {{ . | cyan }}",
		Selected: "Image Selected: {{ . | yellow }}",
	}

	prompt := promptui.Select{
		Label:     "Select Source Image",
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

	cmd = fmt.Sprintf("docker tag %s %s/%s/%s", selected_image, target_prefix, namespace_target, target_image)
	_, err = exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return helper.CliErrorGen(errors.New("error: error image not found"), 1)
	}
	helper.CliSuccessVerbose("Successfully created tag " + target_image + " from image " + source_image)

	return nil

}
