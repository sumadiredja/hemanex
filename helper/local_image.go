package helper

import (
	"errors"
	"fmt"
	"os/exec"
)

func DeleteImageCommand(image_name string, force string) error {
	var cmd = fmt.Sprintf("docker rmi %s%s", force, image_name)
	if _, err := exec.Command("bash", "-c", cmd).Output(); err != nil {
		return CliErrorGen(errors.New("error: error getting input from user"), 1)
	}
	CliSuccessVerbose("Successfully deleted image " + image_name)
	return nil
}
