package helper

import (
	"fmt"
	"syscall"

	b64 "encoding/base64"

	"golang.org/x/term"
)

func GetInputOrFlags(flags string) string {
	var user_input string
	if flags != "" {
		user_input = flags
	} else {
		fmt.Print("Enter Nexus Host: ")
		fmt.Scan(&user_input)
	}
	return user_input
}

func GetPassword(flags string) (string, error) {
	var bytePassword []byte

	if flags != "" {
		bytePassword = []byte(flags)
	} else {
		var err error
		fmt.Print("Enter Nexus Password: ")
		bytePassword, err = term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			return "", CliErrorGen(err, 1)
		}
	}

	return b64.StdEncoding.EncodeToString(bytePassword), nil
}
