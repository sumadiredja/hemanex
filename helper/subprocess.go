package helper

import (
	"fmt"

	"github.com/estebangarcia21/subprocess"
)

const ErrorCertificateResponse = "error: Cannot login to Nexus,\nautorization failed, registry not found or registry is using self signed certificate\n\nif the registry self signed\nplease add the registry to docker daemon.json.\nplease read this https://docs.docker.com/registry/insecure/\n\nif you using podman please provide -k flag"

func RunSubProcess(command string, error_message string) (*subprocess.Subprocess, error) {
	process := subprocess.New(command, subprocess.Shell)
	if err := process.Exec(); process.ExitCode() != 0 || err != nil {
		return process, CliErrorGen(fmt.Errorf(error_message), 1)
	}
	return process, nil
}
