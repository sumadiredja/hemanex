package menu

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"hemanex/registry"

	helper "hemanex/helper"

	b64 "encoding/base64"

	"github.com/estebangarcia21/subprocess"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
)

const (
	CREDENTIALS_TEMPLATES = `# Nexus Credentials
nexus_host = "{{ .Host }}"
nexus_host_port = "{{ .NexusPort }}"
nexus_repository = "{{ .Repository }}"
nexus_repository_port = "{{ .Port }}"
nexus_namespace = "{{ .Namespace }}"
nexus_username = "{{ .Username }}"
nexus_password = "{{ .Password }}"
`
)

func SetNexusCredentials(c *cli.Context) error {
	var hostname, repository, username, password, namespace, nexus_host_port string
	var CREDENTIALS_FILE string
	var err error
	var tmpl *template.Template
	var f *os.File
	var port string = c.String("repository-port")
	var isInsecure bool = c.Bool("insecure-registry")
	var skipTls string

	// CREDENTIALS_FILE, err := helper.GetCredentialPath()
	if CREDENTIALS_FILE, err = helper.GetCredentialPath(); err != nil {
		return err
	}

	hostname = helper.GetInputOrFlags(c.String("nexus-host"), "Host", func(input string) error {
		if len(strings.Split(input, "//")) == 2 {
			return nil
		}
		return errors.New("Please provide https:// or http://")
	})
	nexus_host_port = helper.GetInputOrFlags(c.String("host-port"), "Host Port", func(input string) error {
		_, err := strconv.Atoi(input)
		if err != nil {
			return errors.New("Invalid Port")
		}
		return nil
	})

	repository = helper.GetInputOrFlags(c.String("repository-name"), "Repository Name", func(input string) error {
		if len(input) != 0 {
			return nil
		}
		return errors.New("Please provide the repository-name")
	})
	port = helper.GetInputOrFlags(c.String("repository-port"), "Repository Port", func(input string) error {
		_, err := strconv.Atoi(input)
		if err != nil {
			return errors.New("Invalid Port")
		}
		return nil
	})
	namespace = helper.GetInputOrFlags(c.String("namespace"), "Namespace", func(input string) error {
		if len(input) != 0 {
			return nil
		}
		return errors.New("Please provide the namespace")
	})
	username = helper.GetInputOrFlags(c.String("username"), "Username", func(input string) error {
		if len(input) != 0 {
			return nil
		}
		return errors.New("Please provide the password")
	})

	if password, err = helper.GetPassword(c.String("password"), func(input string) error {
		if len(input) < 6 {
			return errors.New("Password must have more than 6 characters")
		}
		return nil
	}); err != nil {
		return err
	}

	if !c.Bool("ignore-confirmation") {
		prompt := promptui.Prompt{
			Label:     "Are you sure to login with this credentials",
			IsConfirm: true,
		}

		_, err = prompt.Run()

		if err != nil {
			fmt.Print(helper.CliErrorGen(fmt.Errorf("User not logged in : %v\n", "user decide to cancel the login"), 1))
			os.Exit(1)
		}
	}

	data := struct {
		Host       string
		NexusPort  string
		Port       string
		Username   string
		Password   string
		Repository string
		Namespace  string
	}{
		hostname,
		nexus_host_port,
		port,
		username,
		password,
		repository,
		namespace,
	}

	decodePassword, _ := b64.StdEncoding.DecodeString(password)

	if isInsecure {
		skipTls = " --tls-verify=false"
	}

	cmdLogin := fmt.Sprintf("docker login " + strings.Split(hostname, "//")[1] + ":" + port + " -u " + username + " -p " + string(decodePassword) + skipTls)
	login := subprocess.New(cmdLogin, subprocess.Shell)
	if err = login.Exec(); err != nil {
		return helper.CliErrorGen(err, 1)
	}

	if login.ExitCode() != 0 {
		return helper.CliErrorGen(fmt.Errorf("Error: Cannot login to Nexus,\nautorization failed or registry is using self signed certificate\n\nif the registry self signed\nplease add the registry to docker daemon.json.\nplease read this https://docs.docker.com/registry/insecure/\n\nif you using podman please provide -k flag"), 1)
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

	helper.CliSuccessVerbose("nexus user configured")
	return nil
}

func CheckToml(c *cli.Context) error {
	r, err := registry.NewRegistry(c)
	if err != nil {
		return helper.CliErrorGen(err, 1)

	}
	fmt.Println(r.Host, r.Password, r.Repository, r.Username, r.Namespace)
	// fmt.Println(c.Bool("insecure-registry"))
	// fmt.Println(r.Host)
	return nil
}

func GetRepository(c *cli.Context) error {
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
			Namespace  string
		}{
			r.Host,
			r.Username,
			b64.StdEncoding.EncodeToString([]byte(r.Password)),
			c.String("repository-name"),
			r.Namespace,
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

	if c.String("namespace-name") != "" {
		data := struct {
			Host       string
			Username   string
			Password   string
			Repository string
			Namespace  string
		}{
			r.Host,
			r.Username,
			b64.StdEncoding.EncodeToString([]byte(r.Password)),
			r.Repository,
			c.String("namespace-name"),
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

		helper.CliSuccessVerbose(fmt.Sprintf("Namespace changed to %s", c.String("namespace-name")))
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
		cli.ShowSubcommandHelp(c)
		return nil
	}
	if cwd == "" {
		cli.ShowSubcommandHelp(c)
		return nil
	}

	tag_split := strings.Split(tags, ":")
	if len(tag_split) <= 1 {
		cli.ShowSubcommandHelp(c)
		return nil
	}
	image_name = tag_split[0]
	tag = tag_split[1]
	if tag == "" {
		cli.ShowSubcommandHelp(c)
		return nil
	}
	if port == "" {
		port = r.Port
	}
	host := strings.Split(r.Host, "://")[1]
	command := fmt.Sprintf("docker build -t %s:%s/%s/%s:%s %s", host, port, namespace, image_name, tag, cwd)
	s := subprocess.New(command, subprocess.Shell)

	if _ = s.Exec(); s.ExitCode() != 0 {
		return helper.CliErrorGen(errors.New("Failed to build image"), 1)

	}

	helper.CliSuccessVerbose("Successfully built image " + image_name + ":" + tag + " with namespace " + namespace)

	return nil
}

func PushImage(c *cli.Context) error {
	var imgName = c.Args().Get(0)
	var port = c.String("repository-port")
	var isInsecure = c.Bool("insecure-registry")
	var skipTls string
	var namespace string

	if imgName == "" {
		cli.ShowSubcommandHelp(c)
		return nil
	}

	r, err := registry.NewRegistry(c)
	if err != nil {
		return helper.CliErrorGen(err, 1)
	}

	namespace = r.Namespace
	if c.String("namespace") != "" {
		namespace = c.String("namespace")
	}

	if port == "" {
		port = r.Port
	}

	if isInsecure {
		skipTls = " --tls-verify=false"
	}

	cmdLogin := fmt.Sprintf("docker login " + strings.Split(r.Host, "//")[1] + ":" + port + " -u " + r.Username + " -p " + r.Password + skipTls)
	login := subprocess.New(cmdLogin, subprocess.Shell)
	if err = login.Exec(); login.ExitCode() != 0 {
		return helper.CliErrorGen(fmt.Errorf("Error: Cannot login to Nexus,\nautorization failed or registry is using self signed certificate\n\nif the registry self signed\nplease add the registry to docker daemon.json.\nplease read this https://docs.docker.com/registry/insecure/\n\nif you using podman please provide -k flag"), 1)
	}

	cmdPushImage := fmt.Sprintf("docker push " + strings.Split(r.Host, "//")[1] + ":" + port + "/" + namespace + "/" + imgName + skipTls)
	pushImage := subprocess.New(cmdPushImage, subprocess.Shell)

	if err = pushImage.Exec(); pushImage.ExitCode() != 0 {
		return helper.CliErrorGen(fmt.Errorf("Error: Cannot push to Nexus,\nautorization failed or registry is using self signed certificate\n\nif the registry self signed\nplease add the registry to docker daemon.json.\nplease read this https://docs.docker.com/registry/insecure/\n\nif you using podman please provide -k flag"), 1)
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
	var port = r.Port

	if image_name == "" {
		cli.ShowSubcommandHelp(c)
		return nil
	}
	name_split := strings.Split(image_name, ":")
	img_name = name_split[0]
	if len(name_split) <= 1 {
		cli.ShowSubcommandHelp(c)
		return nil
	}
	tag = name_split[1]
	if tag == "" {
		cli.ShowSubcommandHelp(c)
		return nil
	}
	if c.Bool("force") {
		force = "-f "
	}
	if c.String("repository-port") != "" {
		port = c.String("repository-port")
	}

	prefix = fmt.Sprintf("%s:%s", strings.Split(r.Host, "://")[1], port)

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
			return helper.CliErrorGen(errors.New("Error: Error getting input from user"), 1)
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
				return helper.CliErrorGen(errors.New("Error: Error getting input from user"), 1)
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

	cmd = fmt.Sprintf("docker rmi %s%s", force, selected_image)
	_, err = exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return helper.CliErrorGen(errors.New("Error: Error getting input from user"), 1)
	}
	helper.CliSuccessVerbose("Successfully deleted image " + selected_image)

	return nil
}
