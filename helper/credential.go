package helper

import (
	"fmt"
	"html/template"
	"os"
)

type CredentialsDataType struct {
	Host           string
	NexusPort      string
	Repository     string
	RepositoryPort string
	Namespace      string
	Username       string
	Password       string
}

func CredentialsWritter(data CredentialsDataType, CREDENTIALS_FILE string, CREDENTIALS_TEMPLATES string, success_message string) error {
	var tmpl *template.Template
	var f *os.File
	var err error

	if tmpl, err = template.New(CREDENTIALS_FILE).Parse(CREDENTIALS_TEMPLATES); err != nil {
		return CliErrorGen(err, 1)
	}

	if f, err = os.Create(CREDENTIALS_FILE); err != nil {
		return CliErrorGen(err, 1)
	}

	if err = tmpl.Execute(f, data); err != nil {
		return CliErrorGen(err, 1)
	}

	CliSuccessVerbose(fmt.Sprintf(success_message))

	return nil
}

func GetCredentialPath() (string, error) {
	var CREDENTIALS_FILE string

	if IsWindows() {
		if _, err := os.Stat(os.Getenv("AppData") + "\\hemanex"); os.IsNotExist(err) {
			if os.Mkdir(os.Getenv("AppData")+"\\hemanex", 0666) != nil {
				return "", CliErrorGen(fmt.Errorf("error: %s", err), 1)
			}

			if os.Chmod(os.Getenv("AppData")+"\\hemanex", 0666) != nil {
				return "", CliErrorGen(fmt.Errorf("error: %s", err), 1)
			}
		}
		CREDENTIALS_FILE = os.Getenv("AppData") + "\\hemanex" + "\\.credentials"
	} else {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return "", CliErrorGen(fmt.Errorf("error: %s", err), 1)
		}

		var hemanexUserConfigDir string = userHomeDir + "/.config/hemanex"

		if _, err := os.Stat(hemanexUserConfigDir); os.IsNotExist(err) {
			if os.Mkdir(hemanexUserConfigDir, 0755) != nil {
				return "", CliErrorGen(fmt.Errorf("error: %s", err), 1)
			}

			if os.Chmod(hemanexUserConfigDir, 0755) != nil {
				return "", CliErrorGen(fmt.Errorf("error: %s", err), 1)
			}
		}
		CREDENTIALS_FILE = hemanexUserConfigDir + "/.credentials"
	}
	return CREDENTIALS_FILE, nil
}
