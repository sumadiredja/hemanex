package helper

import (
	"fmt"
	"os"

	b64 "encoding/base64"

	"github.com/manifoldco/promptui"
)

func CheckFlagsStringExist(flags string, repository_data string) string {
	if flags != "" {
		return flags
	}
	return repository_data
}

func GetInputOrFlags(flags string, input string, validation func(input string) error) string {
	var user_input string
	if flags != "" {
		user_input = flags
	} else {
		prompt := promptui.Prompt{
			Label:    "Enter Nexus " + input,
			Validate: validation,
		}
		result, err := prompt.Run()

		if err != nil {
			fmt.Print(CliErrorGen(fmt.Errorf("error : %v", "error getting input from user"), 1))
			os.Exit(1)
		}

		user_input = result
	}
	return user_input

}

func GetPassword(flags string, validation func(input string) error) (string, error) {
	var bytePassword []byte

	if flags != "" {
		bytePassword = []byte(flags)
	} else {
		var err error
		prompt := promptui.Prompt{
			Label:    "Password",
			Validate: validation,
			Mask:     '*',
		}

		result, err := prompt.Run()

		if err != nil {
			fmt.Print(CliErrorGen(fmt.Errorf("error : %v", "error getting input from user"), 1))
			os.Exit(1)
		}
		bytePassword = []byte(result)
	}

	return b64.StdEncoding.EncodeToString(bytePassword), nil
}
