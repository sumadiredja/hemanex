package menu

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sumadiredja/hemanex/registry"

	helper "github.com/sumadiredja/hemanex/helper"
	validator "github.com/sumadiredja/hemanex/validator"

	b64 "encoding/base64"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
)

const (
	CREDENTIALS_TEMPLATES = `# Nexus Credentials
nexus_host = "{{ .Host }}"
nexus_host_port = "{{ .NexusPort }}"
nexus_repository = "{{ .Repository }}"
nexus_repository_port = "{{ .RepositoryPort }}"
nexus_namespace = "{{ .Namespace }}"
nexus_username = "{{ .Username }}"
nexus_password = "{{ .Password }}"
`
)

func SetNexusCredentials(c *cli.Context) error {
	var hostname, repository, username, password, namespace, nexus_host_port, port string
	var CREDENTIALS_FILE string
	var err error
	var skipTls string

	isInsecure := c.Bool("insecure-registry")

	if CREDENTIALS_FILE, err = helper.GetCredentialPath(); err != nil {
		return err
	}

	hostname = helper.GetInputOrFlags(c.String("nexus-host"), "Host", validator.HostnameValidator)
	nexus_host_port = helper.GetInputOrFlags(c.String("nexus-host-port"), "Host Port", validator.IsNumber("Invalid Port"))
	repository = helper.GetInputOrFlags(c.String("repository-name"), "Repository Name", validator.IsNotEmpty("repository name"))
	port = helper.GetInputOrFlags(c.String("repository-port"), "Repository Port", validator.IsNumber("Invalid Port"))
	namespace = helper.GetInputOrFlags(c.String("namespace"), "Namespace", validator.IsNotEmpty("namespace"))
	username = helper.GetInputOrFlags(c.String("username"), "Username", validator.IsNotEmpty("username"))

	if password, err = helper.GetPassword(c.String("password"), validator.PasswordValidator); err != nil {
		return err
	}

	if !c.Bool("ignore-confirmation") {
		prompt := promptui.Prompt{
			Label:     "Are you sure to login with this credentials",
			IsConfirm: true,
		}

		_, err = prompt.Run()

		if err != nil {
			fmt.Print(helper.CliErrorGen(fmt.Errorf("user not logged in : %v", "user decide to cancel the login"), 1))
			os.Exit(1)
		}
	}

	data := helper.CredentialsDataType{
		Host:           hostname,
		NexusPort:      nexus_host_port,
		Repository:     repository,
		RepositoryPort: port,
		Namespace:      namespace,
		Username:       username,
		Password:       password,
	}

	decodePassword, _ := b64.StdEncoding.DecodeString(password)

	if isInsecure {
		skipTls = " --tls-verify=false"
	}

	cmdLogin := fmt.Sprintf("docker login " + strings.Split(hostname, "//")[1] + ":" + port + " -u " + username + " -p " + string(decodePassword) + skipTls)
	_, err = helper.RunSubProcess(cmdLogin, helper.ErrorCertificateResponse)

	if err != nil {
		return err
	}

	err = helper.CredentialsWritter(data, CREDENTIALS_FILE, CREDENTIALS_TEMPLATES, "nexus user configured")
	if err != nil {
		return err
	}
	return nil
}

func CheckToml(c *cli.Context) error {
	r, err := registry.NewRegistry(c)
	if err != nil {
		return helper.CliErrorGen(err, 1)

	}
	fmt.Println(r.Host, r.Password, r.Repository, r.Username, r.Namespace)

	return nil
}

func GetRepository(c *cli.Context) error {
	var r registry.Registry
	var err error
	var CREDENTIALS_FILE string

	if r, err = registry.NewRegistry(c); err != nil {
		return helper.CliErrorGen(err, 1)
	}

	if CREDENTIALS_FILE, err = helper.GetCredentialPath(); err != nil {
		return err
	}

	if c.String("repository-name") != "" {
		data := helper.CredentialsDataType{
			Host:           r.Host,
			NexusPort:      r.NexusPort,
			Repository:     c.String("repository-name"),
			RepositoryPort: r.RepositoryPort,
			Namespace:      r.Namespace,
			Username:       r.Username,
			Password:       b64.StdEncoding.EncodeToString([]byte(r.Password)),
		}

		err := helper.CredentialsWritter(data, CREDENTIALS_FILE, CREDENTIALS_TEMPLATES, fmt.Sprintf("Repository changed to %s", c.String("repository-name")))
		if err != nil {
			return err
		}
		return nil
	}
	helper.CliInfoVerbose(fmt.Sprintf("Currently working in %s repository", r.Repository))
	return nil
}

func GetNamespace(c *cli.Context) error {
	var r registry.Registry
	var err error
	var CREDENTIALS_FILE string

	if r, err = registry.NewRegistry(c); err != nil {
		return helper.CliErrorGen(err, 1)
	}

	if CREDENTIALS_FILE, err = helper.GetCredentialPath(); err != nil {
		return err
	}

	if c.String("namespace-name") != "" {
		data := helper.CredentialsDataType{
			Host:           r.Host,
			NexusPort:      r.NexusPort,
			Repository:     r.Repository,
			RepositoryPort: r.RepositoryPort,
			Namespace:      c.String("namespace-name"),
			Username:       r.Username,
			Password:       b64.StdEncoding.EncodeToString([]byte(r.Password)),
		}

		err := helper.CredentialsWritter(data, CREDENTIALS_FILE, CREDENTIALS_TEMPLATES, fmt.Sprintf("Namespace changed to %s", c.String("namespace-name")))
		if err != nil {
			return err
		}
		return nil
	}
	helper.CliInfoVerbose(fmt.Sprintf("Current namespace is %s", r.Namespace))
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
		return helper.ShowSubCommand("please provide image name", c)
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
		return helper.ShowSubCommand("please provide image name and tag", c)
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
		return helper.ShowSubCommand("please provide image name", c)
	} else {
		r, err := registry.NewRegistry(c)
		if err != nil {
			return helper.CliErrorGen(err, 1)
		}
		if tag == "" {
			if keep == 0 {
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
				if len(tags) >= keep {
					for _, tag := range tags[:len(tags)-keep] {
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
		return helper.ShowSubCommand("please provide image name", c)
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
	var tag, img_name string
	var err error
	var force = ""
	var prefix string

	image_name := c.Args().Get(0)
	r, err := registry.NewRegistry(c)
	if err != nil {
		return helper.CliErrorGen(err, 1)
	}
	var port = r.RepositoryPort

	if image_name == "" {
		return helper.ShowSubCommand("please provide image name", c)
	}
	name_split := strings.Split(image_name, ":")
	img_name = name_split[0]
	if len(name_split) <= 1 {
		return helper.ShowSubCommand("please provide image tags", c)
	}
	tag = name_split[1]
	if tag == "" {
		return helper.ShowSubCommand("please provide correct image tags", c)
	}
	if c.Bool("force") {
		force = "-f "
	}
	if c.String("repository-port") != "" {
		port = c.String("repository-port")
	}

	prefix = strings.Split(r.Host, "://")[1] + ":" + port

	// cmd := exec.Command("docker", "images")
	cmd := fmt.Sprintf(`docker images --format="{{.Repository}} {{.Tag}}" | grep "%s" | grep "%s/%s" | grep "%s"`, prefix, r.Namespace, img_name, tag)
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
		_, err = exec.Command("bash", "-c", cmd).Output()
		if err != nil {
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
