package helper

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func DeleteImageCommand(image_name string, force string) error {
	var cmd = fmt.Sprintf("docker rmi %s%s", force, image_name)
	if _, err := exec.Command("bash", "-c", cmd).Output(); err != nil {
		return CliErrorGen(errors.New("error: error getting input from user"), 1)
	}
	CliSuccessVerbose("Successfully deleted image " + image_name)
	return nil
}

func GetAllNoneImageID() (string, error) {
	var imageIDs = ""
	var cmd = "docker images --format '{{.Repository}} {{.ID}}' | grep '<none>'"
	res, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", CliErrorGen(errors.New("error: error getting input from user"), 1)
	}

	list_id := strings.Split(strings.ReplaceAll(string(res), "\rn", "\n"), "\n")
	// if list_id[len(list_id)-1] == "" {
	// 	list_id = list_id[0 : len(list_id)-1]
	// }
	for i := 0; i < len(list_id)-1; i++ {
		if i == 0 {
			imageIDs = strings.Split(list_id[i], " ")[1]
		} else {
			imageIDs = imageIDs + " " + strings.Split(list_id[i], " ")[1]
		}
	}
	return imageIDs, nil
}
