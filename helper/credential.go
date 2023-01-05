package helper

import (
	"errors"
	"fmt"
	"os"
)

func GetCredentialPath() (string, error) {
	var CREDENTIALS_FILE string

	if IsWindows() {
		if _, err := os.Stat(os.Getenv("AppData") + "\\hemanex"); os.IsNotExist(err) {
			if os.Mkdir(os.Getenv("AppData")+"\\hemanex", 0666) != nil {
				return "", CliErrorGen(errors.New(fmt.Sprintf("Error: %s", err)), 1)
			}

			if os.Chmod(os.Getenv("AppData")+"\\hemanex", 0666) != nil {
				return "", CliErrorGen(errors.New(fmt.Sprintf("Error: %s", err)), 1)
			}
		}
		CREDENTIALS_FILE = os.Getenv("AppData") + "\\hemanex" + "\\.credentials"
	} else {
		if _, err := os.Stat("/opt/hemanex"); os.IsNotExist(err) {
			if os.Mkdir("/opt/hemanex", 0755) != nil {
				return "", CliErrorGen(errors.New(fmt.Sprintf("Error: please run as superuser (sudo) \nExample : sudo hemanex login <flags>")), 1)
			}

			if os.Chmod("/opt/hemanex", 0755) != nil {
				return "", CliErrorGen(errors.New(fmt.Sprintf("Error: please run as superuser (sudo) \nExample : sudo hemanex login <flags>")), 1)
			}
		}
		CREDENTIALS_FILE = "/opt/hemanex/.credentials"
	}

	return CREDENTIALS_FILE, nil

}
