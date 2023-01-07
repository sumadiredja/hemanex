package menu

import (
	b64 "encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	helper "github.com/sumadiredja/hemanex/helper"
	"github.com/sumadiredja/hemanex/registry"
	validator "github.com/sumadiredja/hemanex/validator"
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
