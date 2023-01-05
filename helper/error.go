package helper

import (
	"errors"
	"fmt"

	"github.com/urfave/cli"
)

const colorRed = "\033[0;31m"
const colorYellow = "\033[0;33m"
const colorGreen = "\033[0;32m"

const colorNone = "\033[0m"

func CliErrorGen(err error, status int) *cli.ExitError {
	return cli.NewExitError(errors.New(fmt.Sprintf("%s%s%s", colorRed, err, colorNone)), status)
}

func CliCriticalGen(err error, status int) *cli.ExitError {
	return cli.NewExitError(errors.New(fmt.Sprintf("%s[Critial] %s%s", colorRed, err, colorNone)), status)
}

func CliInfoVerbose(message string) {
	fmt.Printf("%s[Info] %s%s\n", colorYellow, message, colorNone)
}

func CliSuccessVerbose(message string) {
	fmt.Printf("%s[Success] %s%s\n", colorGreen, message, colorGreen)
}
