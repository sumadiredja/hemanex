package helper

import (
	"fmt"

	"github.com/urfave/cli"
)

const colorRed = "\033[0;31m"
const colorYellow = "\033[0;33m"
const colorGreen = "\033[0;32m"

const colorNone = "\033[0m"

func ShowSubCommand(err_message string, c *cli.Context) error {
	fmt.Println(colorRed + err_message + colorNone)
	fmt.Println()
	err := cli.ShowSubcommandHelp(c)
	if err != nil {
		return CliErrorGen(err, 1)
	}
	return nil
}

func CliErrorGen(err error, status int) *cli.ExitError {
	return cli.NewExitError(fmt.Errorf("%s%s%s", colorRed, err, colorNone), status)
}

func CliCriticalGen(err error, status int) *cli.ExitError {

	return cli.NewExitError(fmt.Errorf("%s[Critial] %s%s", colorRed, err, colorNone), status)
}

func CliInfoVerbose(message string) {
	fmt.Printf("%s[Info] %s%s\n", colorYellow, message, colorNone)
}

func CliSuccessVerbose(message string) {
	fmt.Printf("%s[Success] %s%s\n", colorGreen, message, colorNone)
}
