package helper

import (
	"fmt"
	"os"
)

func GetCredentialPath() (string, error) {
	var CREDENTIALS_FILE string

	if IsWindows() {
		if _, err := os.Stat(os.Getenv("AppData") + "\\hemanex"); os.IsNotExist(err) {
			if os.Mkdir(os.Getenv("AppData")+"\\hemanex", 0666) != nil {
				return "", CliErrorGen(fmt.Errorf("Error: %s", err), 1)
			}

			if os.Chmod(os.Getenv("AppData")+"\\hemanex", 0666) != nil {
				return "", CliErrorGen(fmt.Errorf("Error: %s", err), 1)
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
